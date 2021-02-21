---
layout: post
title:  3 tips for (slightly) better Go iota enums
date:   2020-10-24
author: Paschalis Ts
tags:   [golang]
mathjax: false
description: ""  
---

Golang has no native enum types; the most idiomatic way of representing an enumerator is to use [constants](https://golang.org/pkg/os/#pkg-constants), many times along with [iota](https://golang.org/ref/spec#Iota). 

Here are three tips that can make for (just slightly) better Go iota enums.

## Add one more 'count' element
I first saw this trick in the Go standard library, eg. when enumerating [Goroutine states](https://github.com/golang/go/blob/93810ac1f4574e1e2a79ea156781bafaf8b8ebe0/src/cmd/trace/trace.go#L471), [register bounds](https://github.com/golang/go/blob/de932da453f68b8fc04e9c2ab25136748173c806/src/cmd/compile/internal/ssa/op.go#L372) or when [marking GC roots](https://github.com/golang/go/blob/master/src/runtime/mgcmark.go#L18).

In the most usual case it allows to build simpler validation methods.

```go
type method uint8

const (
    Get method = iota
    Post
    Patch
    ...
    methodCount
)
```

```go 
func (m method) IsValid() bool {
    return m < methodCount
}
```

In other cases it may allow optimizations like using fixed-size arrays instead of slices.

## Don't use [...]string{}; try a map instead
I've seen people use `[...]string{` when implementing the Stringer interface for their enum type.

Let's say you're building a Fighting game in Go, and you want to represent character states.
```go
type State int

const (
    Standing State = iota
    Walking
    Crouched
)

func (s State) String() string {
    return [...]string{"Stand", "Walk", "Crouch"}[s]
}
```

What's wrong with this snippet? First off, it will happily accept `fmt.Println(State(-1))` and `fmt.Println(State(100))` and panic. It's also a little harder to maintain since you need to keep the element order in mind.

Well how about
```go
var stateNames = map[State]string{
    Standing: "Stand",
    Walking: "Walk",
    Crouched: "Crouch",
}

func (s State) String() string {
    // just return stateNames[s], or
    v, ok := stateNames[s]
    if !ok {
        // do stuff
        return ""
    }
    return v
}
```

The map also provides constant access time, instead of linear in the case of the string array.

## Use uint8 instead of int
First off, why not use a `uint8` instead of an `int`? The 'iota' works the same, you get simpler validation, plus a small performance improvement, almost for free.

## Bonus - Try unexported enum types
Also, in the spirit of 'making invalid states unrepresentable', think about un-exporting the enumerator type. The types themselves can be exported, as well as their methods, but this will disallow people from instantiating their own `State(100)`.
