---
layout: post
title:  What's in a goroutine and how large is one?
date:   2020-05-23
author: Paschalis Ts
tags:   [golang, internals]
mathjax: false
description: "Look at them Go!"
---

I'm pretty sure that anyone learning Go has heard that *"goroutines are like lightweight threads"* and that *"it's okay to launch hundreds, thousands of goroutines"*. Some people learn that *"a goroutine takes up around 2 kilobytes"*, most likely referencing [the Go 1.4 release notes](https://golang.org/doc/go1.4#runtime), and even fewer learn that this represents its initial stack size.

And while all those statements are true, I'd like to show *why* is that, explore what exactly is a goroutine, how much space it takes, and provide starting points for anyone to poke around the Go internals.

For this exploration I'll be using the [Go 1.14 release branch](https://github.com/golang/go/tree/release-branch.go1.14), so all code snippets will point there.

## The Goroutine scheduler

The [Goroutine scheduler](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/proc.go#L19) is a work-stealing scheduler introduced back in Go 1.1 by Dmitry Vyukov and the Go team. Its design document is available [here](https://golang.org/s/go11sched) and discusses possible future improvements. There are lots of [great](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html) [resources](https://rakyll.org/scheduler/) to grok how it works in depth, but the main thing to understand is that it tries to manage **G's**, **M's** and **P's** ; goroutines, machine threads and processors.

A "G"  is simply a Golang goroutine.
An "M" is an OS thread that can be either executing something or idle.
A "P" can be thought as a CPU in the OS' scheduler; it represents the resources required to execute our Go code, such as a scheduler, or a memory allocator state.

These are represented in the runtime as structs of [`type g`](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/runtime2.go#L395), [`type m`](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/runtime2.go#L473), or [`type p`](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/runtime2.go#L552).

***The scheduler's main responsibility is to match up each G (the code we want to execute) to an M (where to execute it) and a P (the rights and resources to execute it).***

When an M stops executing our code, it returns its P to the idle P pool. To resume executing Go code, it must re-acquire it. Similarly, when a goroutine exits, the G object is returned to a pool of free Gs, and can later be reused for some other goroutine.

When starting a Goroutine, either firing up [main](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/asm_amd64.s#L216) or [in code](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/asm_amd64.s#L777), a `g` struct is initialized.

A new goroutine, i.e. an object of type `g` is created via the [`malg`](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/proc.go#L3353) function, 
```go
// Allocate a new g, with a stack big enough for stacksize bytes.
func malg(stacksize int32) *g {
	newg := new(g)      // <--- this is where it all starts 
	if stacksize >= 0 {
		stacksize = round2(_StackSystem + stacksize)
		systemstack(func() {
			newg.stack = stackalloc(uint32(stacksize))
		})
		newg.stackguard0 = newg.stack.lo + _StackGuard
        newg.stackguard1 = ^uintptr(0)
        ...
        ...
    }
    return newg
}
```

which is called from [newproc](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/proc.go#L3376) and [newproc1](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/proc.go#L3388).

```go
// Create a new g running fn with narg bytes of arguments starting
// at argp. callerpc is the address of the go statement that created
// this. The new g is put on the queue of g's waiting to run.
func newproc1(fn *funcval, argp unsafe.Pointer, narg int32, callergp *g, callerpc uintptr) {
  ...
    acquirem() // disable preemption because it can be holding p in a local var
	siz := narg
    siz = (siz + 7) &^ 7
  ...
	_p_ := _g_.m.p.ptr()
	newg := gfget(_p_)
	if newg == nil {
		newg = malg(_StackMin) // !!! <-- magic happens here
		casgstatus(newg, _Gidle, _Gdead)
		allgadd(newg) 
	}
  ...
}        

```

So now we're ready to dissect the goroutine itself!

## The Goroutine Object

The goroutine object is about 70 lines long. Let me remove the comments and clean it up a little.

```go
type g struct {
    stack            stack   
    stackguard0      uintptr 
    stackguard1      uintptr 
    _panic           *_panic 
    _defer           *_defer 
    m                *m      
    sched            gobuf
    syscallsp        uintptr        
    syscallpc        uintptr        
    stktopsp         uintptr        
    param            unsafe.Pointer 
    atomicstatus     uint32
    stackLock        uint32 
    goid             int64
    schedlink        guintptr
    waitsince        int64      
    waitreason       waitReason
    preempt          bool 
    preemptStop      bool 
    preemptShrink    bool 
    asyncSafePoint   bool
    paniconfault     bool 
    gcscandone       bool 
    throwsplit       bool 
    activeStackChans bool
    raceignore       int8     
    sysblocktraced   bool     
    sysexitticks     int64   
    traceseq         uint64   
    tracelastp       puintptr 
    lockedm          muintptr
    sig              uint32
    writebuf         []byte
    sigcode0         uintptr
    sigcode1         uintptr
    sigpc            uintptr
    gopc             uintptr         
    ancestors        *[]ancestorInfo 
    startpc          uintptr         
    racectx          uintptr
    waiting          *sudog        
    cgoCtxt          []uintptr     
    labels           unsafe.Pointer
    timer            *timer        
    selectDone       uint32        
    gcAssistBytes    int64
}
```

And that's all there really is to it! 

Let's try adding these numbers up; a `uintptr` is 64-bits, so 8 bytes in our architecture, same as an `int64`. Booleans are 1 byte long and a slice is just a pointer plus two integers. 

There are some more complex type such as `timer` (~70 bytes), `_panic` (~40 bytes), or `_defer` (~100 bytes), but I'm getting around ~600 bytes in total.

Hmm, seems a little fishy.. Where does the famous "2 kb" value come from?

Let's take a closer look to the first struct field and explore ...

## The Goroutine Stack

The first field of the `g` struct is a `stack` type. 
```go
type g struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// stackguard1 is the stack pointer compared in the C stack growth prologue.
...
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

```


The stack itself is nothing more than two values denoting where it begins and ends.
```go
type stack struct {
	lo uintptr
	hi uintptr
}
```

By this time, you're either probably wondering *Hmm, so what is the size of this stack?*, or you've already guessed that the 2 kilobytes refer exactly to this stack size!

***A goroutine [starts](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/proc.go#L3410) with a [2-kilobyte](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/stack.go#L72) minimum stack size which grows and shrinks as needed without the risk of ever running out.***

[This](https://dave.cheney.net/2013/06/02/why-is-a-goroutines-stack-infinite) excellent post by Dave Cheney explains how this works in more detail. Essentially, before executing any function Go checks whether the amount of stack required for the function it's about to execute is available; if not a call is made to [`runtime.morestack`](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/asm_amd64.s#L407) which allocates a new page and only then the function is executed. Finally, when that function exits, its return arguments are copied back to the original stack frame, and any unneeded stack space is released.

While the minimum stack size is defined as 2048 bytes, the Go runtime does also not allow goroutines to exceed [a maximum stack size](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/stack.go#L1031); this maximum depends on the architecture and is [1 GB for 64-bit and 250MB for 32-bit](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/proc.go#L120) systems.   
If this limit is reached a call to [`runtime.abort`](https://github.com/golang/go/blob/f296b7a6f045325a230f77e9bda1470b1270f817/src/runtime/asm_amd64.s#L450) will take place.

Exceeding this stack size is *very* easy with a recursive function; all you have to do is

```go
package main

func foo(i int) int {
	if i < 1e8 {
		return foo(i + 1)
	}
	return -1
}

func main() {
	foo(0)
}
```
And we can see that the application panics, the stack can no longer grow and the aforementioned `runtime.abort` is executed.

```bash
$ go run exceed-stack.go
runtime: goroutine stack exceeds 1000000000-byte limit
fatal error: stack overflow

runtime stack:
runtime.throw(0x1071ce1, 0xe)
	/usr/local/go/src/runtime/panic.go:774 +0x72
runtime.newstack()
	/usr/local/go/src/runtime/stack.go:1046 +0x6e9
runtime.morestack()
	/usr/local/go/src/runtime/asm_amd64.s:449 +0x8f

goroutine 1 [running]:
main.foo(0xffffdf, 0x0)
...
...
```


## So, how many can you run?

I'm using the script in the Appendix, which was copied from [here](https://stackoverflow.com/a/8534711). 

On a mid-end laptop, I'm able to launch *50 million* goroutines.   
As the number grows, there are two main concerns: memory usage (and you start swapping), and slower garbage collection.

```bash
$ ~ go run poc-goroutines-sizing.go

# 10 Thousand goroutines
Number of goroutines: 100000
Per goroutine:
  Memory: 2115.71 bytes
  Time:   1.404500 µs

# 1 Million goroutines
Number of goroutines: 1000000
Per goroutine:
  Memory: 2655.21 bytes
  Time:   1.518857 µs

# 3 Million goroutines
Number of goroutines: 3000000
Per goroutine:
  Memory: 2700.37 bytes
  Time:   1.637003 µs

# 6 Million goroutines
Number of goroutines: 6000000
Per goroutine:
  Memory: 2700.29 bytes
  Time:   2.541744 µs

# 9 Million goroutines
Number of goroutines: 9000000
Per goroutine:
  Memory: 2700.27 bytes
  Time:   2.857699 µs

# 12 Million goroutines
Number of goroutines: 12000000
Per goroutine:
  Memory: 2694.09 bytes
  Time:   3.232870 µs

# 50 Million goroutines
Number of goroutines: 50000000
Per goroutine:
  Memory: 2695.37 bytes
  Time:   5.098005 µs
```



## Outro

So more or less, that's all! 

There's the Goroutine scheduler which is how Go code is scheduled to run on the host. Then there are the Goroutines themselves, is the way that Go code is actually executed, and there's each goroutine's stack which grows and shrinks to accommodate the code execution.

I recommend skimming over [src/runtime/HACKING.md](https://github.com/golang/go/blob/release-branch.go1.14/src/runtime/HACKING.md) where many of the concepts and conventions of the code in the Golang runtime are explained in more detail.

I hope you learned something new, and have some waypoints in order to poke into the code of the Go language itself.

Until next time, bye!

## Resources

- https://stackoverflow.com/questions/8509152/max-number-of-goroutines
- https://medium.com/a-journey-with-go/go-how-does-the-goroutine-stack-size-evolve-447fc02085e5
- https://dave.cheney.net/2013/06/02/why-is-a-goroutines-stack-infinite
- https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html
- https://medium.com/@genchilu/if-a-goroutine-call-a-new-goroutine-which-one-would-scheduler-pick-up-first-890002dc54f8
- https://povilasv.me/go-scheduler/




## Appendix

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

var n = flag.Int("n", 3*1e6, "Number of goroutines to create")

var ch = make(chan byte)
var counter = 0

func f() {
	counter++
	<-ch // Block this goroutine
}

func main() {
	flag.Parse()
	if *n <= 0 {
		fmt.Fprintf(os.Stderr, "invalid number of goroutines")
		os.Exit(1)
	}

	// Limit the number of spare OS threads to just 1
	runtime.GOMAXPROCS(1)

	// Make a copy of MemStats
	var m0 runtime.MemStats
	runtime.ReadMemStats(&m0)

	t0 := time.Now().UnixNano()
	for i := 0; i < *n; i++ {
		go f()
	}
	runtime.Gosched()
	t1 := time.Now().UnixNano()
	runtime.GC()

	// Make a copy of MemStats
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	if counter != *n {
		fmt.Fprintf(os.Stderr, "failed to begin execution of all goroutines")
		os.Exit(1)
	}

	fmt.Printf("Number of goroutines: %d\n", *n)
	fmt.Printf("Per goroutine:\n")
	fmt.Printf("  Memory: %.2f bytes\n", float64(m1.Sys-m0.Sys)/float64(*n))
	fmt.Printf("  Time:   %f µs\n", float64(t1-t0)/float64(*n)/1e3)
}
```
