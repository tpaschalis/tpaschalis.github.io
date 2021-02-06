---
layout: post
title:  Using runes for direct comparisons
date:   2021-02-04
author: Paschalis Ts
tags:   [golang, code]
mathjax: false
description: ""
---

Many people have described Go as "C for the 21st century". One of the things that reminds me of this metaphor is the handling of text using runes.

For a *great* introduction on Runes, UTF-8 Strings and generally text in Go, I recommend that you refer to Donovan and Kernighan's '[Blue Book](https://www.amazon.com/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440)'. I feel that we can use runes to enhance performance in some straightforward ways.

I encountered this specific example in the Go codebase while attempting to solve issue [#44052](https://github.com/golang/go/issues/44052) with [CL 288712](https://go-review.googlesource.com/c/go/+/288712/). Long story short, it's about validating the [major version suffix](https://golang.org/ref/mod#go-mod-file-ident) for Go modules. 

The way this check is implemented in Go modules is by directly comparing the input using runes, 'simulating' the behavior of `strconv.Atoi` using `r >= '0' && r <= '9'`. 

To be fair, that's exactly what [Atoi](https://golang.org/src/strconv/atoi.go?s=5658:5690#L214) and [ParseInt/ParseUint](https://golang.org/src/strconv/atoi.go?s=5658:5690#L214) do in the background.

## Let's compare some code!

The rule that we'll try to implement is the following

> For a final path element of the form `/vN` where (...) N must not begin with a leading zero, must not be /v1, and must not contain any dots. 

So let's say you want to validate this given rule.. How would you go about it?

## Regex!
Oh, so you want to have *two* problems on your hands. Jokes aside, not necessarily a bad solution! 

Here's how someone might go about solving this
```go
var re = regexp.MustCompile(`^(\/v[2-9]|\/v([1-9][0-9]+))$`)

func parseUsingRegex(p string) bool {
    return re.MatchString(p)
}
```

Looks simple, right? Except for that write once/read never regex, that you *hope* is correct (but write a bunch of tests to make sure anyway).

## Stdlib
Why not use the tools in the standard library to implement these rules in a more straightforward way?

```go
func parseUsingStdlib(p string) bool {
    if !strings.HasPrefix(p, "/v") {
        return false
    }

    if strings.HasPrefix(p[2:], "0") {
        return false
    }

    n, err := strconv.Atoi(p[2:])
    if err != nil {
        return false
    }

    if n < 2 {
        return false
    }

    return true
}
```

I'll admit, this can be read at a glance, and is pretty easy to expand.

## Runes
Finally, here's (almost) what [rsc used](https://github.com/golang/mod/commit/5d307ac8d37c05b4a8ce233dfc138f2bc5783c7b) in the Modules toolset.

```go
func parseUsingRunes(p string) bool {
    if p[:2] != "/v" || p[2] == '0' || p == "/v1" {
        return false
    }

    i := len(p)
    for i > 2 {
        if !('0' <= p[i-1] && p[i-1] <= '9') {
            return false
        }
        i--
    }

    return true
}
```

I feel this reads more like C. Such handling make sense if we are certain that we're dealing with ASCII, but we *do* lose some of the previous solution's readability. What do we gain in return?

## Benchmarks

So, how do these three examples stack against each other? Take five seconds to guess...
<br>
<br>
<br>
<br>

Well, I ran the benchmark presented below; I expected *some* difference in performance, but not in that scale.

```go
var c1, c2, c3, c4, c5, c6, c7 bool

func BenchmarkRegex(b *testing.B) {
    var b1, b2, b3, b4, b5, b6, b7 bool
    for i := 0; i < b.N; i++ {
        b1 = parseUsingRegex("/v0")
        b2 = parseUsingRegex("/v1")
        b3 = parseUsingRegex("/v02")
        b4 = parseUsingRegex("/v3.0")
        b5 = parseUsingRegex("/v3b")
        b6 = parseUsingRegex("/v3")
        b7 = parseUsingRegex("/v42")
    }

    c1, c2, c3, c4, c5, c6, c7 = b1, b2, b3, b4, b5, b6, b7
}
```

For these seven cases, the rune comparison outperforms the 'stdlib' solution by 8.5x, and the 'regex' solution by 45.5x.

```
BenchmarkRegex-8    	 1447944	       816 ns/op
BenchmarkStdlib-8   	 7559529	       153 ns/op
BenchmarkRunes-8    	64005085	       17.9 ns/op
```


## Parting words
I personally have a soft-spot for such 'old-school' code, but I understand that possible performance gains do not usually outweigh the loss of readability.

This example might be a little simplistic, and in general there are a lot of pitfalls when working with UTF-8 strings. 

In any case, what I'd like to underline is, when you're trying to squeeze some more performance out of some critical path, don't hesitate to use the language's lower level building blocks.

Until next time, bye!