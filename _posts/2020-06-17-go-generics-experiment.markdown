---
layout: post
title:  Experimenting with Go Generics
date:   2020-06-17
author: Paschalis Ts
tags:   [golang, generics]
mathjax: false
description: "Could this be it!"
---

## Intro

On June 16th, a refined generics design draft was published by the Go team.

The announcement post can be viewed [here](https://blog.golang.org/generics-next-step), while the actual draft is [here](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md). The jury is still out on this, but the authors cite Philip Wadler's [Featherweight Go](https://arxiv.org/abs/2005.11710) as a source of influence. Personally, I'm very happy that the language it taking a methodical, careful approach to such a radical change.

As the authors mention, the actual implementation *if and when* generics are accepted will look and work differently, but this is a good time to theorize!

There are ~two~ three more important things : 

- The *earliest possible date* for the launch is August 2021, with Go 1.17.
- There is now a type checker and a new version of playground supporting generics
- There were some actual working examples provided

The document is quite large, and not so easy to grok. I'm slowly going through the document, but I couldn't wait to run and share some code!

*All of these snippets were based on the design draft page*.

## Min - Max

The simplest of examples is, of course min/max. [Playground Link](https://go2goplay.golang.org/p/XhMBKX7gmqa)

```go
package main

import (
	"fmt"
)

type num interface {
	type int, int32, int64, uint, uint32, uint64, float32, float64
}

func min(type T num)(a, b T) T {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println(min(5, 1))
	fmt.Println(min(0.4, 0.2))
}
```

## Sum

Another trivial operation, would be a simple sum. [Playground Link](https://go2goplay.golang.org/p/aAuFIZLOKNA)

```go   
package main

import (
	"fmt"
)

type num interface {
	type int, int32, int64, uint, uint32, uint64, float32, float64
}

func sum(type T num)(s []T) T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

func main() {
	fmt.Println(sum([]int{1, 4, -9, 42}))
	fmt.Println(sum([]float64{1.618, -3.141, 2.718, }))
}
```


## Set

The design draft introduces *a new predeclared type constraint: comparable*, which includes types that can be compared using `==` and `!=`.

Here's the simple Set implementation that the authors published, which supports the Add, Remove, Contains and Length operations. [Playground Link](https://go2goplay.golang.org/p/JU_fAhn3Pfo)

```go
package main

import (
	"fmt"
)

type Set(type T comparable) struct {
	items map[T]struct{}
}

func MakeSet(type T comparable)(items []T) *Set(T) {
	set := new(Set(T))
	set.items = make(map[T]struct{})

	for _, v := range items {
		set.Add(v)
	}

	return set
}

func (s *Set(T)) Length() int {
	return len(s.items)
}

func (s *Set(T)) Add(it T) {
	s.items[it] = struct{}{}
}

func (s *Set(T)) Remove(it T) {
	delete(s.items, it)
}

func (s *Set(T)) Contains(it T) bool {
	_, ok := s.items[it]
	return ok
}

func main() {

	s1 := MakeSet([]int{0, 0, 1, 2})
	fmt.Println(s1)
	fmt.Println(s1.Contains(4))

	s1.Add(4)
	fmt.Println(s1.Contains(4))
	fmt.Println(s1.Length())
	s1.Remove(4)
	fmt.Println(s1.Contains(4))

	s2 := MakeSet([]string{"ATL", "LAX", "ORD", "DFW"})
	s2.Remove("ORD")
	s2.Remove("DFW")
	s2.Add("ATH")
	s2.Add("ATH")
	s2.Add("ATH")

	fmt.Println(s2.Length())
}
```



## Map - Filter - Reduce

You're probably thinking, wow hold on a second. A usable Map/Filter/Reduce implementation! [Playground Link](https://go2goplay.golang.org/p/-Cdr3jQ1RGS)
```go
package main

import (
	"fmt"
	"math/big"
)

type num interface {
	type int, int32, int64, uint, uint32, uint64, float32, float64
}

func Map(type Tin, Tout num)(s []Tin, f func(Tin) Tout) []Tout {
	res := make([]Tout, len(s))
	for i, v := range s { res[i] = f(v) }
	return res
}

func Reduce(type Tin, Tout num)(s []Tin, init Tout, f func(Tout, Tin) Tout) Tout {
	res := init
	for _, v := range s { res = f(res, v) }
	return res
}

func Filter(type T num)(s []T, f func(T) bool) []T {
	var res []T
	for _, v := range s {
		if f(v) { res = append(res, v) }
	}
	return res
}


func main() {
	n := []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("%+v\n", n)
	
	f := Map(n, func(i int) float64 { return float64(i)/10 })
	fmt.Printf("%+v\n", f)
	
	p := Reduce(f, 1, func(i, j float64) float64 {return i * j})
	fmt.Printf("%.5f\n", p)
	
	primes := Filter(n, func(i int) bool { return big.NewInt(int64(i)).ProbablyPrime(0) })
	fmt.Printf("%+v\n", primes)
}
```

## Outro

It's getting late, and I've got to get up early tomorrow. The design draft is there for you to read, but I hope these examples did pique your interest.

Until next time, bye!