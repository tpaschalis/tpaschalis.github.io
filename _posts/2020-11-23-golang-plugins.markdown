---
layout: post
title:  Go plugins
date:   2020-11-23
author: Paschalis Ts
tags:   [golang, code]
mathjax: false
description: "Don't do this kids"
---

For today's edition of [things you should never use in Go, but are fun to know about](https://tpaschalis.github.io/golang-linknames/), we have Go *plugins*! As is the case with CGO, they make use of the dynamic linker which rules out static binaries.

Plugins were introduced back in early 2017 with Go 1.8, and allow loading code dynamically at run time. A special [build mode](https://golang.org/cmd/go/#hdr-Build_modes) enables compiling packages into shared object (.so) libraries, and the [plugin package](https://golang.org/pkg/plugin/) implements loading and symbol resolution.

So in one sentence, a plugin is a `main` package with exported functions and variables, that has been built with 
```bash 
go build -buildmode=plugin
```

The helper package contains only two types and two methods; 
- `Plugin` which represents a loaded Go plugin, with
- `plugin.Open` which opens a Go plugin and
- `plugin.Lookup` which searches for a symbol via its name, in a threadsafe manner

plus 
- `Symbol` which is a pointer to an exported variable or a function

That last sentence also means that we cannot load *constants*, *interfaces* or *types* from plugins. In case we want our loaded objects to use a method, we will have to provide the interface that it implements.

## Constraints

Unfortunately, there's a bunch of constraints that make working with plugins an unpleasant experience.

- Your code and plugin must be compiled with the exact same compiler version
- Your code and plugin must be compiled with the same `GOPATH` variable
- As we mentioned, code using plugins will not be statically linked
- Any packages imported by both the code and plugin must have the exact same version
- They don't work in Windows
- Their size is *greatly* inflated
- Cannot use with `vendor` folder
- If you're trying to debug with Delve, you might run into issues, as your plugins need to be compiled with the exact same flags

With all their shortcomings, Go plugins have seen some real-world uses: [Tyk](https://tyk.io/docs/plugins/supported-languages/golang/) and [Kong](https://docs.konghq.com/enterprise/2.1.x/go/) used them so clients can customize and extend their services, while [gosh](https://github.com/vladimirvivien/gosh) uses them to build interactive console-based shell programs.

Finally, Hashicorp's [go-plugin](https://github.com/hashicorp/go-plugin) tries to provide similar functionality, but I don't know how it fares in the real world.


## Show me some code!

Let's read through the following silly geometry package. It contains an exported variable, a couple of functions that convert between degrees and radians, as well as two types: `circle` and `ellipse` with support for calculating their `Area()` and `Circumference()`.
```go
package main

import "math"

var Pi = 3.14159

func Deg2Rad(deg float64) float64 {
	return deg * Pi / 180.
}

func Rad2Deg(rad float64) float64 {
	return rad * 180. / Pi
}

var Circle circle       // exported to be used by the caller
var Ellipse ellipse

type Shape interface {
	Area() float64
	Circ() float64
}

type circle struct {
	cx, cy, r float64
}

type ellipse struct {
	cx, cy, a, b float64
}

func DefCircle(cx, cy, r float64) {
	Circle = circle{cx, cy, r}
}

func DefEllipse(cx, cy, a, b float64) {
	Ellipse = ellipse{cx, cy, a, b}
}

func (c circle) Area() float64 {
	return Pi * c.r * c.r
}

func (c circle) Circ() float64 {
	return 2 * Pi * c.r
}

func (e ellipse) Area() float64 {
	return Pi * e.a * e.b
}

func (e ellipse) Circ() float64 {
	return Pi * (3*(e.a+e.b) - math.Sqrt((3*e.a+e.b)*(e.a+3*e.b)))
}
```

We can compile this package as a plugin by running
```bash
go build -buildmode=plugin -o plugin/geo.so
```

To reuse it in a different package, all we have to do is 
```go
// Load the plugin
p, err := plugin.Open("plugin/geo.so")  

// Look up symbols using their names
piSymbol, err := p.Lookup("Pi")         
r2dSymbol, err := p.Lookup("Rad2Deg")

// Cast the address contents to the correct type
piValue, ok := *piSym.(*float64)            
rad2deg, ok := r2dSym.(func(float64) float64)

// Done!
fmt.Println(piValue, rad2deg(20.0))
```

Here's a snippet that reuses all of the symbols we defined on the `geo.so` plugin above. As we mentioned, we cannot load types, or interfaces, so we have to define the interface (shape) that we want the inferred types (circle, ellipse) to implement.

This is the same reason why we're not returning a `circle` from the `NewCircle` function, but we act on an exported `Circle` variable as you'll notice.

```go
package main

import (
	"fmt"
	"plugin"
)

type shape interface {
	Area() float64
	Circ() float64
}

func main() {
	p, _ := plugin.Open("plugin/geo.so")

	piSym, _ := p.Lookup("Pi")
	piValue := *piSym.(*float64)
	fmt.Println("Stored Pi value is :", piValue)

	r2dSym, _ := p.Lookup("Rad2Deg")
	rad2deg := r2dSym.(func(float64) float64)
	fmt.Println("1 rad to degrees = ", rad2deg(1.))

	circleSym, _ := p.Lookup("Circle")
	circle := circleSym.(shape)
	ellipseSym, _ := p.Lookup("Ellipse")
	ellipse := ellipseSym.(shape)

	defCircleSym, _ := p.Lookup("DefCircle")
	defCircle := defCircleSym.(func(float64, float64, float64))
	defCircle(3.0, 5.0, 10.0)
	fmt.Println("The circle area is :", circle.Area())
	fmt.Println("The circle circumference is :", circle.Circ())

	defEllipseSym, _ := p.Lookup("DefEllipse")
	defEllipse := defEllipseSym.(func(float64, float64, float64, float64))
	defEllipse(0.0, 5.0, 10.0, 13.0)
	fmt.Println("The ellipse area is :", ellipse.Area())
	fmt.Println("The ellipse circumference is : ~", ellipse.Circ())
```

## Outro

That’s all about plugins; make sure to *not* use them, and discourage your co-workers from doing so. They are against Go's philosophy as a statically-linked language, they are clunky, need special environment for building and maintaining, and can panic unexpectedly.

Until next time, bye!

## Resources
- https://golang.org/pkg/plugin/
- https://golang.org/cmd/go/#hdr-Build_modes
- https://golang.org/doc/go1.8
- https://www.reddit.com/r/golang/comments/b6h8qq/is_anyone_actually_using_go_plugins/
- https://medium.com/@alperkose/things-to-avoid-while-using-golang-plugins-f34c0a636e8
- https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9
