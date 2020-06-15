---
layout: post
title:  What is a defer?<br>How are they implemented and how many can you run?
date:   2020-06-15
author: Paschalis Ts
tags:   [golang, internals]
mathjax: false
description: "Look at them Go; in a few minutes!"
---

Defers are one of my favorite Go features.

They offer *predictability* and routinely simplify the way that our code interacts with the host system.

So, naturally I got curious and attempted look under the hood to find out how they're implemented. Grab some coffee and let's go!

For this exploration I'll be using the [Go 1.14 release branch](https://github.com/golang/go/tree/release-branch.go1.14), so all code snippets will point there.


## Intro -- TL;DR
In the Go runtime, defers are handled like [goroutines](https://tpaschalis.github.io/goroutines-size/) or channels -- as constructs of the language itself.

Multiple defers are stacked on the *defer chain*, and executed them in LIFO order, as in this example below.

```go
func main() {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")
}
// $ go run main.go
// 3
// 2
// 1
```

There are countless tutorials and guides on *using* defers out there, so why don't we just dive into the Go code directly?

## The _defer struct

The docstring informs us that the `_defer` struct is just an entry on *a list of deferred calls*. Some of these entries will live on the stack and some on the heap. Here's the entire struct definition, with some comments omitted for brevity.

```go
// A _defer holds an entry on the list of deferred calls.
// Some defers will be allocated on the stack and some on the heap.
// All defers are logically part of the stack, so write barriers to
// initialize them are not required. All defers must be manually scanned,
// and for heap defers, marked.
type _defer struct {
	siz     int32 // includes both arguments and results
	started bool
    heap    bool
    
	// openDefer indicates that this _defer is for a frame with open-coded defers
	openDefer bool
	sp        uintptr  // stack pointer at time of defer
	pc        uintptr  // program counter at time of defer
	fn        *funcval // can be nil for open-coded defers
	_panic    *_panic  // panic that is running defer
	link      *_defer

	fd   unsafe.Pointer // funcdata for the function associated with the frame
	varp uintptr        // value of varp for the stack frame
    // framepc is the current program counter associated with the stack frame. 
    // Together, with sp above (the stack pointer associated with the frame),
	// framepc/sp can be used to continue a stack trace via gentraceback().
	framepc uintptr
}
```


### What's an Open-Coded defer?
As we explore the `_defer` struct, we come across the `openDefer` field, which specifies whether a defer is *open-coded*.

This concept was introduced in Go just this year (February 2020) in [CL 202340](https://go-review.googlesource.com/c/go/+/202340/) and launched with Go 1.14.
Open-coded is short for "not being called as part of a loop"; the cost of these kind of defers was *greatly* lowered via [inlining](https://en.wikipedia.org/wiki/Inline_expansion) machine code and storing some extra `FUNCDATA` due to their predictable nature. 

Benchmarks such as the following has lead some people to declare defers an *almost zero-cost abstraction*, which while an exaggeration, proves a point.
```
Cost of defer statement  [ go test -run NONE -bench BenchmarkDefer$ runtime ]
  With normal (stack-allocated) defers only:         35.4  ns/op
  With open-coded defers:                             5.6  ns/op
  Cost of function call alone (remove defer keyword): 4.4  ns/op
```


Here's a short example from `defer_test.go`, which tests the behavior of an open-coded and a non-open-coded defer
```go
func TestOpenAndNonOpenDefers(t *testing.T) {
	for {
        // f() is a more complicated function that is recover()'ed
        defer f()   // <-- non open-coded defer
    }
    defer f()       // <-- open-coded defer
}
```

Currently, there's a limit of 8 open-coded defers in a function (defined in [`maxOpenDefers`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/cmd/compile/internal/gc/ssa.go#L34)) as this optimization is meant for smaller functions; after this point, inlined defers are disallowed and the older implementation is used.

```go
if Curfn.Func.numDefers > maxOpenDefers {
    // Don't allow open-coded defers if there are more than 8 defers in the 
    // function, since we use a single byte to record active defers.
    Curfn.Func.SetOpenCodedDeferDisallowed(true)
}
```

### Creation and Execution of a defer

So what happens when we call `defer` from our code?

When the compiler encounters a defer statement, it will turn it into a [`deferproc`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/panic.go#L218) or [`deferprocStack`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/panic.go#L271) call at that specific point, as well as a [`deferreturn`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/panic.go#L528) at the return point of the function.

Let's see this with a couple of examples! We can just compile the source code and use the *go tool* command to inspect the generated code, like `go tool objdump -S main > compiler-generated.s`.

```go 
package main

func main() {
    for i := 0; i < 5; i++ {
        defer f()
    }
}

func f() {}
```

If we inspect the code, we can see the following section
```asm
...
  0x1057003    4889442408    MOVQ AX, 0x8(SP)
  0x1057008    e8a303fdff    CALL runtime.deferproc(SB)    <--- defer is created
  0x105700d    85c0          TESTL AX, AX
  0x105700f    7502          JNE 0x1057013
  0x1057011    ebce          JMP 0x1056fe1
  0x1057013    90            NOPL
  0x1057014    e8570cfdff    CALL runtime.deferreturn(SB)  <--- defer is created
  0x1057019    488b6c2418    MOVQ 0x18(SP), BP
...
```


### Let's dig deeper

So what does actually happen in `deferproc` and `deferreturn`?

As we see, `deferproc` gets the current goroutine, and uses [`newdefer`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/panic.go#L387) for the actual creation. After some checks, the newly created defer is allocated and added on the defer chain.

```go 
// Create a new deferred function fn with siz bytes of arguments.
// The compiler turns a defer statement into a call to this.
//go:nosplit
func deferproc(siz int32, fn *funcval) { // arguments of fn follow fn
	gp := getg()
	if gp.m.curg != gp {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}
...
	d := newdefer(siz)
	if d._panic != nil {
		throw("deferproc: d.panic != nil after newdefer")
	}
	d.link = gp._defer
	gp._defer = d
	d.fn = fn
	d.pc = callerpc
	d.sp = sp
...
	return0()
}
```

```go
// Allocate a Defer, usually using per-P pool.
// Each defer must be released with freedefer.  The defer is not
// added to any defer chain yet.
func newdefer(siz int32) *_defer {
	var d *_defer
	sc := deferclass(uintptr(siz))
	gp := getg()
...
	if d == nil {
		// Allocate new defer+args.
		systemstack(func() {
			total := roundupsize(totaldefersize(uintptr(siz)))
			d = (*_defer)(mallocgc(total, deferType, true))
		})
...
			d.siz = siz
			d.link = gp._defer
			gp._defer = d
			return d
...
}
```

As we mentioned, the deferreturn is placed at the return point of the 'parent' function.
It gets the current goroutine, and checks whether there are any deferred functions. If a deferred function is detected, a call to `runtime.jmpdefer` will be executed, which will jump to the deferred function, and will execute it to appear like it has been called at that point. The deferreturn is called again and again, until there are no more deferred functions, that is until `jmpdefer` can flip the Program Counter over to the current function.

```go
// Run a deferred function if there is one.
// The compiler inserts a call to this at the end of any
// function which calls defer.
// If there is a deferred function, this will call runtime·jmpdefer,
// which will jump to the deferred function such that it appears
// to have been called by the caller of deferreturn at the point
// just before deferreturn was called. The effect is that deferreturn
// is called again and again until there are no more deferred functions.
//
// Declared as nosplit, because the function should not be preempted once we start
// modifying the caller's frame in order to reuse the frame to call the deferred
// function.
//
// The single argument isn't actually used - it just has its address
// taken so it can be matched against pending defers.
//go:nosplit
func deferreturn(arg0 uintptr) {
	gp := getg()
	d := gp._defer
...
	if d.openDefer {
		done := runOpenDeferFrame(gp, d)
		if !done {
			throw("unfinished open-coded defers in deferreturn")
		}
		gp._defer = d.link
		freedefer(d)
		return
	}
	...
	fn := d.fn
	d.fn = nil
	gp._defer = d.link
	freedefer(d)

_ = fn.fn
	jmpdefer(fn, uintptr(unsafe.Pointer(&arg0)))
}
```

Subsequently, we can see that the `freedefer` is responsible for cleaning up all defers as they're being executed.

```go
// Free the given defer.
// The defer cannot be used after this call.
//
// This must not grow the stack because there may be a frame without a
// stack map when this is called.
//
//go:nosplit
func freedefer(d *_defer) {
...
	if d.fn != nil {
		freedeferfn()
	}
	if !d.heap {
		return
	}
...
	// These lines used to be simply `*d = _defer{}` but that
	// started causing a nosplit stack overflow via typedmemmove.
	d.siz = 0
	d.started = false
	d.openDefer = false
...
	pp.deferpool[sc] = append(pp.deferpool[sc], d)
}
```

Items placed on the `deferpool` will be cleaned up by the Garbage Collector; you can examine how it's done [here](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/mgc.go#L2241).



## How many defers could we run?

What do you think will happen by running the following piece of code?

```go
package main

import "fmt"

func main() {
    f()
}

func f() {
    fmt.Println("Hey!")
    defer f()
}
```

If you read my [previous piece on goroutines](https://tpaschalis.github.io/goroutines-size/), you could guess that the *defer chain* of the main goroutine, will eventually run out of stack size; and you'd be right! After a few moments you'd be hit by

```
runtime: goroutine stack exceeds 1000000000-byte limit
runtime: sp=0xc020108378 stack=[0xc020108000, 0xc040108000]
fatal error: stack overflow
```

If we include a simple counter, we see that as of Go 1.14 we can fit nearly 480k defers (4793475 to be exact) per stack frame, on its 1GB max size.

## How costly are defers?

Well, it's easy to find out! 

Many of the operations that you'll encounter in your daily work are constrained not by Go, but by the host system, such as hardware specs, or kernel limits.

For example, sufficiently modern Linux kernels will limit the max file descriptors `/proc/sys/fs/file-max`.

The limit is enforced [here](https://github.com/torvalds/linux/blob/cb8e59cc87201af93dfbb6c3dccc8fcad72a09c2/fs/file_table.c#L134) and which should work out to around 590432
```c
void __init files_maxfiles_init(void)
{
	unsigned long n;
	unsigned long nr_pages = totalram_pages();
	unsigned long memreserve = (nr_pages - nr_free_pages()) * 3/2;

	memreserve = min(memreserve, nr_pages - 1);
	n = ((nr_pages - memreserve) * (PAGE_SIZE / 1024)) / 10;

	files_stat.max_files = max_t(unsigned long, n, NR_FILE);
}
```


So, let's try to measure defer performance in the straightforward task of opening as many files as possible and allocating one defer for each, using the following code snippet. We'll also read one file at random, just to avoid any optimizations of files closing before their time.

```go
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

var files []*os.File

func main() {
	N := 250_000
	for i := 0; i < N; i++ {
		n := strconv.Itoa(i)
		filename := "file-" + n + ".txt"
		file, err := os.Open("../testdata/" + filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		files = append(files, file)
	}

	fmt.Println("Opened all files; length of array is :", len(files))

	r := rand.Intn(N)
	scanner := bufio.NewScanner(files[r])
	for scanner.Scan() {
		fmt.Println("Reading file ")
	}
	//time.Sleep(120 * time.Second)
	fmt.Println("Exiting...")
}
```

Opening 250k files and as many defers is straightforward. On my system, it consumes about 80MB of real memory and about 5GB of Virtual memory.

```bash
$ time go run main.go
Opened all files; length of array is : 250000
Exiting...
go run main.go  0.90s user 2.02s system 107% cpu 2.715 total
```

I guess 


## Outro 

So that's all from me today! 

We saw the `_defer` struct itself, and examined what's an open-coded defer and the latest Go 1.14 optimizations.

We examined how defer calls are translated to machine code, and how defers are created, scheduled and executed.

Finally, we ran a few benchmarks to see the overhead of defers; it's a little unlikely that they'll be causing any performance headaches on their own.

I hope you learned something new, and have some waypoints in order to poke into the code of the Go language itself.

Until next time, bye!



<!--


The following is an *open-coded defer*, since it is outside a loop, and thus can be inlined directly into the compiler-generated code. 

```go 
package main

import "fmt"

func main() {
    defer f()
    defer f()
    defer f()
}

func f() {
    fmt.Println("hey!")
}
```

```asm
...
  0x10996e2             e87936f9ff            CALL runtime.deferprocStack(SB) <--- defer is created
  0x10996e7             85c0                  TESTL AX, AX
  0x10996e9             7549                  NE 0x1099734
                                p.fmtString(v.String(), verb)
  0x10996eb             488b442438            MOVQ 0x38(SP), AX
...
  0x1099719             e832e6ffff            CALL fmt.(*pp).fmtString(SB)
                                return
  0x109971e             90                    NOPL
  0x109971f             e83c3ef9ff            CALL runtime.deferreturn(SB) <--- defer is created
...
```

The handling of a non-open-coded defer (living in a loop), is very similar as well!


-->