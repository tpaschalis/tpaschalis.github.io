---
layout: post
title: What are the limits of Go channels, and how fast are they?
date:   2020-06-29
author: Paschalis Ts
tags:   [golang, internals]
mathjax: false
description: ""
---

As part of my upcoming talk at the [GoWay conference](http://goway.io/), I learned a lot about how channels work under the hood. Let's highlight some of the most important numbers!

In this post, all code will point to the [Go 1.14 release branch](https://github.com/golang/go/tree/release-branch.go1.14).

## Intro

Channels are the main form of thread-safe communication between goroutines; but *how fast is that communication*? And what are the limitations that the language imposes on them?

## Limitations

The limitations on the `chan` type are [pretty straightforward](https://github.com/golang/go/blob/release-branch.go1.14/src/runtime/chan.go) to decode.
```go
func makechan(t *chantype, size int) *hchan {
	elem := t.elem

	// compiler checks this but be safe.
	if elem.size >= 1<<16 {
		throw("makechan: invalid channel element type")
	}
	if hchanSize%maxAlign != 0 || elem.align > maxAlign {
		throw("makechan: bad alignment")
	}

	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}
	...
	...
}
```

First off, the maximum message size (or channel type) is 2^16 bytes, or 64 kilobytes.

Secondly, the same limitations that exist for slices or maps are applied here; a check against `maxAlloc` as well as a check using `math.MulUintptr`.

In a 64-bit Unix-like system [`maxAlloc`](https://github.com/golang/go/blob/67d894ee652a3c6fd0a883a33b86686371b96a0e/src/runtime/malloc.go#L217) defines the maximum allowed allocation at 2^47 bytes (~140 Terabytes) while the pointer multiplication ensures a maximum buffer size so that the channel type multiplied by this buffer doesn't overflow our 64-bits.

Here are the maximum buffer sizes allowed by the compiler for some basic data types.
```go
    // Max buffer size allowed by compiler
    ch1 := make(chan struct{}, 1<<63-1)

    ch2 := make(chan bool, 1<<48-96)

    ch3 := make(chan int16, 1<<47-48)

    ch4 := make(chan int32, 1<<46-24)

    ch5 := make(chan int64, 1<<45-12)

    ch6 := make(chan complex128, 1<<44-6)
```


## Speed

On the topic of speed, the actual limit has to do with the price of goroutine context-switching, which should be consistently ≤ 200ns, as the passing-around of the data is not really expensive.

Using the following simple benchmark we can get a measure on the upper limit of the send/receive channel operations.
```go
func BenchmarkUnbufferedChannelEmptyStruct(b *testing.B) {
    ch := make(chan struct{})
    go func() {
        for {
            <-ch
        }
    }()
    for i := 0; i < b.N; i++ {
        ch <- struct{}{}
    }
}
```

I repeated the tests with a buffered channel, and sending a single byte instead of an empty struct. 

On a 2019 MBP, Intel Core i5 @ 1,4 GHz, the results can be seen in the following table.

```
BenchmarkBufferedChannelEmptyStruct-8         23657084            49.9 ns/op
BenchmarkBufferedChannelOneByte-8             21230530            54.6 ns/op      18.31 MB/s

BenchmarkUnbufferedChannelEmptyStruct-8        6075384            177 ns/op
BenchmarkUnbufferedChannelOneByte-8            6341457            184 ns/op       5.44 MB/s
```

First off, we notice a 4x decline when using unbuffered channels due to their blocking behavior slowing down the communication between goroutines.

On this machine, the upper limit is thus *~18-20 million messages per second* using buffered channels and *~5 million messages per second* when using unbuffered ones.

The transfer rate of 18.31 MB/s might seem low, but it's also constrained by the small size of the message type.

## Improvements

I find it a little unlikely that you'll be hitting this kind of limits for passing around messages.

But if you actually do, you can always ensure thread-safety using Mutexes. The Lock/Unlock operation of a Mutex is about 5x faster, and might not require moving the data itself around.

Finally, and as a measuring stick, the naive copying of data from/to a memory address, takes about ~0.5ns

```
func BenchmarkMutexLockUnlock(b *testing.B) {
    var mux sync.Mutex
    for i := 0; i < b.N; i++ {
        mux.Lock()
        mux.Unlock()
    }

}

func BenchmarkNaiveCopy(b *testing.B) {
    from := make([]byte, b.N)
    to := make([]byte, b.N)
    b.ReportAllocs()
    b.ResetTimer()
    b.SetBytes(1)
    copy(to, from)
}

BenchmarkMutexLockUnlock-8                	96631381	        11.8 ns/op
BenchmarkNaiveCopy-8                      	1000000000	         0.548 ns/op	1823.95 MB/s	       0 B/op	       0 allocs/op
```


## Outro
That's all for today! I plan on continuing with the articles dissecting the Go internals, so check back soon for more!
