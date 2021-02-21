---
layout: post
title:  Go finalizers
date:   2021-01-24
author: Paschalis Ts
tags:   [golang, code]
mathjax: false
description: "If you find any good use-cases please let me know"
---

In today's edition of *semi-obscure Go features that you probably shouldn't use*, we have Finalizers! (All code snippets will point to the Go 1.16 release branch)

I didn't even realize that Go has finalizers until I read [this tweet](https://twitter.com/rakyll/status/1343350504349257728) by Jaana Dogan; she has expressed her [thoughts](https://twitter.com/rakyll/status/787735060573138944) on finalizers again in the past.


## What are finalizers?

A [finalizer](https://en.wikipedia.org/wiki/Finalizer) is a special method that is executed during object destruction.

The term 'finalizer' is primarily used in Garbage-collected languages; in contrast to the 'destructor' term used to describe methods called for finalization in languages with manual memory management such as C++.

The use of finalizers is more widespread in Java and C#; I suspect it's because these languages give their users more control over the GC process and properties with knobs and tweaks.

Nevertheless, there are a number well-documented of disadvantages to the use of Finalizers.

## Finalizers in Go

The ability to use finalizers in Go is implemented in the `runtime` package; more specifically the [`runtime.SetFinalizer`](https://golang.org/pkg/runtime/#SetFinalizer) method.

Some example uses of finalizers in the language itself include [network file descriptors](https://github.com/golang/go/blob/3e06467282c6d5678a6273747658c04314e013ef/src/net/fd_posix.go#L32) or when [creating a new process](https://github.com/golang/go/blob/3e06467282c6d5678a6273747658c04314e013ef/src/os/exec.go#L30)
cd 
In the first case a finalizer is set when the address is being set so that `Close()` will always be called, and in cases where `Close()` is called manually, the finalizer is unset first. To unset a finalizer we can just call `runtime.SetFinalizer(obj, nil)`.

In the second cases, we make sure that the process will always call `Release()`.

## Experiment with finalizers!

Let's write a dead-simple example use of a variable with a finalizer
```go
package main

import (
	"log"
	"runtime"
	"time"
)

type foo struct {
	somefield bool
}

func myfunc() {
	f := new(foo)

	runtime.SetFinalizer(f, func(f *foo) {
		log.Print("I'm being garbace collected!")
	})
}

func main() {
	myfunc()
	runtime.GC()

	// Give the scheduler the time to switch to
	// the finalizer goroutine before exiting.
	time.Sleep(10 * time.Millisecond)
}
```

So you write your code, then `go run main.go` and... nothing happens! Why is that? 

Well, as it's stated by the documentation, *there is no guarantee that finalizers will run before a program exits*, so typically they are useful only for releasing non-memory resources associated with an object during a long-running program.
It's also not guaranteed that a finalizer will run if the size of *obj is zero bytes or package-level variables allocated in `init()` functions, as they might be linker-allocated, not heap-allocated.

What we *do* know is that a finalizer *may* run as soon as an object becomes *unreachable*. 

Another example from the documentation where Finalizers will fail to run is objects stored in global variables or that can be traced from global variables' pointers, which remain reachable.


All Finalizers will run in a single goroutine, sequentially, and in dependency order. That means given two objects A and B where A points at B, only A's finalizer is allowed to run, and once A is collected then the finalizer for B can run. This introduces a new failure, as cyclic structures with finalizers are not guaranteed to run and be GC'ed, because there's no obvious ordering that can untangle this dependency. Then, the documentation which states *if a finalizer must run for a long time, it should do so by starting a new goroutine* makes more sense.

## Sample use-cases
After this short overveiw of how Finalizers are set in our and implemented in the Go runtime, let's see (and rebuff) some of their usecases.




## Parting words
All in all, I don't think that if finalizers were to be removed from Go "2", many people would miss them.

## More reading
Here are some other good resources to read about the usage of finalizers in Go.

- https://crawshaw.io/blog/tragedy-of-finalizers
- https://lk4d4.darth.io/posts/finalizers/