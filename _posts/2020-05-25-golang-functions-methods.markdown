---
layout: post
title:  Go Functions vs Methods
date:   2020-05-25
author: Paschalis Ts
tags:   [golang, internals]
mathjax: false
description: "What's the difference"
---

## Appetizer
If you've wondered whether a certain piece of Go code should be declared as a function or as a method on a type, you've probably ended up reading the either [this](https://dave.cheney.net/2016/03/19/should-methods-be-declared-on-t-or-t) excellent post by Dave Cheney or the [Go spec](https://golang.org/ref/spec#Method_declarations) which states that :

*The type of a method is the type of a function with the receiver as first argument.*

So Go handles methods as functions where the first formal parameter is the receiver. Let's see how this is implemented under the hood in the [Go 1.14 release branch](https://github.com/golang/go/tree/release-branch.go1.14).

## Main Course

The `types` package is responsible for declaring the data types and implements the algorithms for type-checking of Go packages.

If we dive in [`src/go/types/call.go`](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/call.go), we can find all the typechecking code of function call and selector expressions along with the following two methods
```go 
func (check *Checker) call(x *operand, e *ast.CallExpr) exprKind {
func (check *Checker) selector(x *operand, e *ast.SelectorExpr) {
```

Looking more carefully in the selector method (as a method call *does* include a selector expression), we can see [the following code snippet](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/call.go#L405)
```go
func (check *Checker) selector(x *operand, e *ast.SelectorExpr) {
...
		// the receiver type becomes the type of the first function
		// argument of the method expression's function type
		var params []*Var
		sig := m.typ.(*Signature)
		if sig.params != nil {
			params = sig.params.vars
		}
		x.mode = value
		x.typ = &Signature{
			params:   NewTuple(append([]*Var{NewVar(token.NoPos, check.pkg, "", x.typ)}, params...)...),
			results:  sig.results,
			variadic: sig.variadic,
		}

        check.addDeclDep(m)
...
}
```

Essentially, it *does what the spec describes*; it substitutes the parameters of the method with a `NewTuple` which takes in the `NoPos` (zeroth argument) token, of the receiver `x.typ` type, and then appends all of the other function parameters!

You might have noticed that this selector method is unexported and thus internal to the `types` package.

Then where, is it used? As it currently stands, in two places. 

The first of these is in the [`exprInternal`](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/expr.go#L1267) method which contains the core logic for type checking all expressions
```
func (check *Checker) exprInternal(x *operand, e ast.Expr, hint Type) exprKind {
 ...
    switch e := e.(type) {
 ...
	case *ast.SelectorExpr:
		check.selector(x, e)
 ...
}
```
which is then called by [`rawExpr`](ttps://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/expr.go#L981) on various places where the Checker needs to, well check an expression.

The second place is in the [`typInternal`](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/typexpr.go#L235) method, which drives the checking of various types 
```go 
func (check *Checker) typInternal(e ast.Expr, def *Named) Type {
	switch e := e.(type) {
	case *ast.BadExpr:
 ...
	case *ast.SelectorExpr:
		var x operand
        check.selector(&x, e)
        
        switch x.mode {
		case typexpr:
			typ := x.typ
			def.setUnderlying(typ)
			return typ
 ...
}
```
and is then called by [`definedType`](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/typexpr.go#L138) where the Checker needs to either [declare](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/decl.go#L542) a type or [retrieve](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/typexpr.go#L118) the type of an expression.


## Dessert

So that's pretty much it! The code that makes up Go itself is surprisingly easy to dive into; if you want to learn more, I'd recommend looking at how the `Checker` object is [defined](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/check.go#L70) and [initialized](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/check.go#L175), as well as how the `type` package starts [evaluating](https://github.com/golang/go/blob/f758dabf52d38333985e24762f9b53a29e2e7da0/src/go/types/eval.go#L24) expressions.

I hope this was an interesting read, and provided some waypoints to start looking into the Go codebase.

Any comments, corrections and advice are highly welcomed; you'll probably make my day by reaching out, so don't hesitate to do so!

Until next time, bye!
