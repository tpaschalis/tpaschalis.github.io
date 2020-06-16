---
layout: post
title:  Experimenting with Go Generics (build a Map/Filter/Reduce!)
date:   2020-06-17
author: Paschalis Ts
tags:   [golang, generics]
mathjax: false
description: "Could this be it!"
---

## Intro

On June 16th, a refined generics design draft was published by the Go team.

The announcement post can be viewed [here](https://blog.golang.org/generics-next-step), while the actual proposal is [here](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md). The community's reaction is still to be seen, but the authors cite Philip Wadler's [Featherweight Go](https://arxiv.org/abs/2005.11710) as a source of influence. Personally, I'm very happy that the language it taking a methodical, careful approach to such a radical change.

As the authors mention, the actual implementation *if and when* they're accepted will look and work differently, but this is a time to theorize!

There were two more important things mentioned. 

- The *earliest possible date* for the launch of generics is August 2021, with Go 1.17. 
- There is now a type checker and a new version of playground supporting generics!

The document is quite large, and not so easy to grok. I'm a little skeptical towards generics, and currently slowly going through the document, but *I couldn't wait to play around and share some results!*

Let's see how generics could simplify some common operations in the Go language!

Thanks to reddit users [/u/Rican7](https://www.reddit.com/user/Rican7) and [/u/PaluMacil](https://www.reddit.com/user/PaluMacil) for providing some examples which helped kickstart things!

## Min - Max

The first thing that came to mind was, of course min/max. [Playground Link](https://go2goplay.golang.org/p/XhMBKX7gmqa).

```go
package main

import (
	"fmt"
)

type num interface {
	type int, int32, int64, uint, uint32, uint64, float32, float64
}

func min(type T num)(a, b T) T {
	if a > b {
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

The second most frequent operation, would be a simple sum. [Playground Link](https://go2goplay.golang.org/p/aAuFIZLOKNA)

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

The design draft introduces *a new predeclared type constraint: comparable*.

One more simple operation, is the Set. Here's a simple implementation that supports Add, Remove, Contains and Length operations. (Thanks to reddit user /u/Rican7)  [Playground Link](https://go2goplay.golang.org/p/JU_fAhn3Pfo).

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

You're probably thinking, wow hold on a second. Yes, it's true, a functional Map/Filter/Reduce implementation for all you functional programming lovers out there! [Playground Link](https://go2goplay.golang.org/p/-Cdr3jQ1RGS)
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

// Reduce reduces a []Tin to a single value using a reduction function.
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

It's getting late, and I've got to get up early tomorrow; I promise I'll try to experiment with some more common operations tomorrow, and keep you posted right here!