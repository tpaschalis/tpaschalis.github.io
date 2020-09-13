---
layout: post
title:  Go linknames; aka Access unexported symbols using this one weird trick, Go developers hate it!
date:   2020-09-13
author: Paschalis Ts
tags:   [golang, code]
mathjax: false
description: ""  
---

The Go codebase has a number of interesting tricks up its sleeve. While I wouldn't easily greenlight them in a production codebase, they are still nice to know about.

In today's edition, let's see one that (among other things) allows accessing unexported symbols from other packages. 

## The //go:linkname directive

Linknames are briefly mentioned in the [compile command](https://golang.org/cmd/compile/) docs. One of their uses in the Go tree is to enable refactoring and allow flexibility while maintaining the compatibility promise.

In contrast to other directives (eg. `//go:norace`) that affect the code blocks immediately under them, they can appear anywhere in a file, using the `//go:linkname localname [importpath.name]` format.

This linkname directive tells the compiler to 'rename' or 'link' the variable or function `localname` to `importpath.name`; in other words, whenever the code calls `importpath.name` it will reach out to `localname` instead. 

If the `importpath.name` is omitted, it will just make the localname symbol accessible to other packages, even if it doesn't start with a capital letter, so they can be called in assembly code.

Keep in mind you'll either have to provide a correct `go tool compile` command so that the compiler skips checking for partially defined functions with the `-complete` flag, or just add a dummy assembly (.s) file for the same effect.

## Linknames in action

Here's an example of linknames in action!

```
# foo/bar/bar.go
package bar

import "encoding/base64"
import "fmt"
import _ "unsafe"

//go:linkname encode main.bar_myencode
//go:linkname decode main.fromString

func Setup() {
    fmt.Println("starting..")
}

func encode(input []byte) string {
    return base64.StdEncoding.EncodeToString(input)
}

func decode(input string) (string, error) {
    res, err := base64.StdEncoding.DecodeString(input)
    if err != nil {
        return "", err
    }
    return string(res), nil
}

# foo/main.go
package main

import (
    "fmt"
    "foo/bar"
    _ "unsafe"
)

func bar_myencode(input []byte) string
func fromString(input string) (string, error)

func main() {
    bar.Setup()
    fmt.Println(bar_myencode([]byte("greetings")))
    fmt.Println(fromString("R28gPDMK"))
}
```


## How is it implemented?

Linkname directives are stored in a [`linkname` slice](https://github.com/golang/go/blob/dbc5602d18397d1841cb7b2e8974d472c15dee83/src/cmd/compile/internal/gc/noder.go#L230), which record the symbol's position on the file, as well as the local and remapped symbol names. The actual parsing of the directive happens [a few lines below](https://github.com/golang/go/blob/dbc5602d18397d1841cb7b2e8974d472c15dee83/src/cmd/compile/internal/gc/noder.go#L1555). This slice is ultimately used to add these symbols as nodes in the AST Node tree.

Linknamed functions are the only ones allowed to [not have a body](https://github.com/golang/go/blob/dbc5602d18397d1841cb7b2e8974d472c15dee83/src/cmd/compile/internal/gc/noder.go#L535), as we saw in our example. 

## How is the check for 'unsafe' implemented?
In [`noder.go`](https://github.com/golang/go/blob/dbc5602d18397d1841cb7b2e8974d472c15dee83/src/cmd/compile/internal/gc/noder.go#L251), there is a comparison with `imported_unsafe`.

That boolean is switched to 'true' if [importfile](https://github.com/golang/go/blob/dbc5602d18397d1841cb7b2e8974d472c15dee83/src/cmd/compile/internal/gc/main.go#L1139) encounters 'unsafe' in the import list. Unsafe is characterized as a 'pseudopackage' since it is mostly used for enabling features like this, which circumvent the type safety of Go programs, and not an actual package containing functionality.


## Outro
That's all about linknames; make sure to *not* use this new toy in your arsenal. 
Unless you have loads of experience, a robust testing infrastructure *and* you're working with low-level primitives or reaching for extreme performance optimizations, you shouldn't circumvent the type safety of your Go code by be messing with 'unsafe'.
 
Until next time!

