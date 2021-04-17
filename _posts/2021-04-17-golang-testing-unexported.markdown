---
layout: post
title:  Testing unexported methods in Go
date:   2021-04-17
author: Paschalis Ts
tags:   [golang, testing]
mathjax: false
description: ""
---

<meta http-equiv='Content-Type' content='text/html; charset=utf-8' />

Today's post will be a short one. So, I recently came across a cool-but-not-immediately-obvious-for-me fact in Go. You can see it [here](https://twitter.com/tpaschalis_/status/1383512838468104197) in the form of a `#golang pop quiz` tweet.

It was brought to my attention my [Emmanuel](https://twitter.com/odeke_et) when he was reviewing a CL for the Go language. At the time I was wondering how to continue testing an unexported method when moving a package tests from `package p` to `package p_test`.

I didn't really want to resort to [linknames](https://tpaschalis.github.io/golang-linknames/), and I didn't think that exporting it was a good solution either. So, he just commented

> We can create a fresh file "export_test.go" in which you export such symbols, given that they can ONLY be used in _test.go files

🤯

*So it seems that Go will NOT export symbols from test files, as I'd expected.* (as they start with a capital letter *and* belong in the same package!)

So here's how this deals with the problem described above, and allows to keep having access to methods, constants, or other test data.

```go
-- p.go --
package p

type T struct {}

func unexportedFunc() {}
func (t T) unexportedMethod() {}

-- export_test.go --
package p
func UnexportedFuncForTest() {
    return unexportedFunc()
}
func UnexportedMethodForTest(t) {
    return t.unexportedMethod()
}

-- p_test.go --
package p_test

import . "p"

UnexportedFuncForTest()
foo := T{}
UnexportedMethodForTest(foo)
```

See you around!
