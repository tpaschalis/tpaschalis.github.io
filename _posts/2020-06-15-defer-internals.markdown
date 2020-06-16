---
layout: post
title:  What is a defer? And how many can you run?
date:   2020-06-15
author: Paschalis Ts
tags:   [golang, internals]
mathjax: false
description: "Look at them Go; in a few minutes!"
---

Defer is one of my favorite Go features.

It offers *predictability* and simplifies the way that we interact with the host system.

So, naturally I got curious and attempted to look under the hood and find out how it's implemented. Grab some coffee and let's go!

In this post, all code will point to the [Go 1.14 release branch](https://github.com/golang/go/tree/release-branch.go1.14).


## Intro -- TL;DR
Quoting Tour of Go *a defer statement defers the execution of a function until the surrounding function returns*.

In the Go runtime, defers are handled like [goroutines](https://tpaschalis.github.io/goroutines-size/) or channels -- as constructs of the language itself. Multiple defers are stacked on the *defer chain*, and executed in LIFO order, as seen here.

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

    fd      unsafe.Pointer // funcdata for the function associated with the frame
    varp    uintptr        // value of varp for the stack frame
    framepc uintptr
}
```


### What's an "open-coded" defer?
As we explore the `_defer` struct, we come across the `openDefer` field, which specifies whether a defer is *open-coded*, which is short for *not being called in a for-loop*.

This concept was introduced in Go just this year (February 2020) in [CL 202340](https://go-review.googlesource.com/c/go/+/202340/) and launched with Go 1.14. The design document is available [here](https://github.com/golang/proposal/blob/master/design/34481-opencoded-defers.md) and is a great read.

Due to their predictable nature, the cost of these kind of defers was *greatly* lowered via [inlining](https://en.wikipedia.org/wiki/Inline_expansion) machine code and storing some extra data about the function they will be calling. 

Benchmarks such as the following has lead some people to declare defers an *almost zero-cost abstraction*, which while an exaggeration, is not so far from the truth.
```
Cost of defer statement  [ go test -run NONE -bench BenchmarkDefer$ runtime ]
  With normal (stack-allocated) defers only:         35.4  ns/op
  With open-coded defers:                             5.6  ns/op
  Cost of function call alone (remove defer keyword): 4.4  ns/op
```


Here's a short example from `defer_test.go` showing an open-coded defer and a non-open-coded one.
```go
func TestOpenAndNonOpenDefers(t *testing.T) {
    // f() is a more complicated function that is recover()'ed  
    for {
        defer f()   // <-- non open-coded defer
    }
    defer f()       // <-- open-coded defer
}
```

Currently, there's a limit of 8 open-coded defers in a function (defined in [`maxOpenDefers`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/cmd/compile/internal/gc/ssa.go#L34)) as this optimization is meant for smaller functions; after this point, inlining is disallowed and the older implementation is used.

```go
if Curfn.Func.numDefers > maxOpenDefers {
    // Don't allow open-coded defers if there are more than 8 defers in the 
    // function, since we use a single byte to record active defers.
    Curfn.Func.SetOpenCodedDeferDisallowed(true)
}
```

### Creation and Execution of a defer

What happens when we call `defer` from our code?

When the compiler encounters a defer statement, it will turn it into a [`deferproc`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/panic.go#L218) or [`deferprocStack`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/panic.go#L271) call at that specific point, as well as a [`deferreturn`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/panic.go#L528) at the return point of the function.

Let's see this in action! We can just compile and use the *go tool* command to inspect the generated code.

```go
// `go build main.go`
// `go tool objdump -S main > compiler-generated.s`
package main

func main() {
    for i := 0; i < 5; i++ {
        defer f()
    }
}

func f() {}
```

Viewing the compiler-generated code, we spot the following section
```asm
  0x1057003    ...    MOVQ AX, 0x8(SP)
  0x1057008    ...    CALL runtime.deferproc(SB)    <--- defer created
  0x105700d    ...    TESTL AX, AX
  0x105700f    ...    JNE 0x1057013
  0x1057011    ...    JMP 0x1056fe1
  0x1057013    ...    NOPL
  0x1057014    ...    CALL runtime.deferreturn(SB)  <--- defer returned
  0x1057019    ...    MOVQ 0x18(SP), BP
```


### Let's dig deeper

So what's actually happening in `deferproc` and `deferreturn`?

As we see, `deferproc` gets the current goroutine, and uses [`newdefer`](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/panic.go#L387) for the actual creation. After some checks, the new defer is allocated and added on the defer chain.

```go 
func deferproc(siz int32, fn *funcval) {
	gp := getg()
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
func newdefer(siz int32) *_defer {
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

A little later, `deferreturn` will be called at the return point of the 'parent' function.

It gets the current goroutine, and checks whether there are any deferred functions. If a deferred function is detected, a call to `runtime.jmpdefer` will be executed, which will jump to the deferred function, and execute it as if it had been called at that point. 

The `deferreturn` is called again and again, until there are no more deferred functions, and `jmpdefer` can flip the Program Counter over to the current function.

```go
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

Finally, we can see that the `freedefer` is responsible for cleaning up all defers as they're being executed.

```go
// Free the given defer.
// The defer cannot be used after this call.
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

Items placed on the `deferpool` will be cleaned up by the GC as seen [here](https://github.com/golang/go/blob/73f86d2a78423f26323e7acf52bc489fb3e7fcbc/src/runtime/mgc.go#L2241).

## How many defers could we run?

What do you think will happen by running the following piece of code?

```go
package main

import "fmt"

func main() {
    f()
}

func f() {
    defer f()
}
```

If you read my [previous piece on goroutines](https://tpaschalis.github.io/goroutines-size/), you could guess that the *defer chain* of the main goroutine, will eventually run out of stack size; and you'd be right! After a few moments you'd be hit by

```bash
#  64-bit stack frames have max size of 1GB
runtime: goroutine stack exceeds 1000000000-byte limit
runtime: sp=0xc020108378 stack=[0xc020108000, 0xc040108000]
fatal error: stack overflow
```

If we include a simple counter, we see that we can fit nearly 4.8 million defers (4'793'476 to be exact) per stack frame.

## How costly are defers?

Well, it's easy to find out! 

Many of the operations that you'll encounter in your daily work are constrained not by Go, but by the host system, such as hardware specs, or kernel limits.

For example, the Linux kernel will limit the max file descriptors, you can check on this limit using `cat /proc/sys/fs/file-max`.

The limit is enforced [here](https://github.com/torvalds/linux/blob/cb8e59cc87201af93dfbb6c3dccc8fcad72a09c2/fs/file_table.c#L134) and should work out to about ~10^5.
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


So, let's try to measure defer performance in a straightforward task :
- Open as many files as possible 
- Allocate one defer for each, using the following code snippet. 
- Read one file at random (to avoid any optimizations of files closing before their time)
- Execute all those defers

```go
// for i in {0..N}; do touch "file-${i}.txt" ; done

func main() {
	N := 250_000
	for i := 0; i < N; i++ {
		file, err := os.Open("../testdata/" + "file-" + strconv.Itoa(i) + ".txt")
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
		fmt.Println("Reading file", r)
	}

    fmt.Println("Exiting...")
}
```

On a mid-end MacBook Pro, opening 250k files and as executing as many defers is completed in less than a second. It consumes about 80MB of real memory and 500MB of virtual memory.

On an older Linux laptop, opening 700k files and executing the defers takes 4.5sec, and consumes about 226MB of real memory and 860MB of virtual memory.

```bash
$ time go run main.go
Opened all files; length of array is : 250000
Exiting...
go run main.go  0.90s user 2.02s system 107% cpu 2.715 total

$ time go run main.go
Opened all files; length of array is : 700000
Exiting...
go run main.go  4.51s user 2.75s system 105% cpu 2.125 total
```

I hope you agree with me in saying that the defers themselves are *preetty cheap*.

## Outro 

That's all from me today! 

We saw the `_defer` struct itself, and explained what's an open-coded defer and the latest Go 1.14 optimizations.

We examined how defer calls are translated to machine code, and how defers are created, scheduled and executed.

Finally, we ran a few benchmarks to see the overhead of defers; on a 64-bit system you can fit ~480k "empty" defers per stack frame. Honestly it's a little unlikely that they'll be causing any performance headaches on their own.

I hope you learned something new, and have some waypoints to start digging around the Go source.

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