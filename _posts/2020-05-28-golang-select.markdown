---
layout: post
title:  How is the select statement implemented in Go? 
date:   2020-05-28
author: Paschalis Ts
tags:   [golang, internals]
mathjax: false
description: ""
---


src/cmd/compile/internal/gc/select.go
src/runtime/select.go

src/cmd/compile/internal/gc/walk.go


Once you get to picking the Golang internals, you'll end up visiting the `gc(?)` package more and more.

There, you'll find the `walkstmt` function, which is used to traverse the AST tree.

```go
func walkstmt(n *Node) *Node {
	if n == nil {
		return n
	}

	setlineno(n)

	walkstmtlist(n.Ninit.Slice())

	switch n.Op {
	default:
		if n.Op == ONAME {
			yyerror("%v is not a top level statement", n.Sym)
		} else {
			yyerror("%v is not a top level statement", n.Op)
		}
		Dump("nottop", n)

 ...
	case OSELECT:
		walkselect(n)
 ...
    }

     return n
}
```



The walkselect function finds out about the possible `walkselectcases` and then walks through these statements using `walkstmtlist`, which calls back to the original `walkstmt` function.


```
func walkstmtlist(s []*Node) {
	for i := range s {
		s[i] = walkstmt(s[i])
	}
}```


```
func walkselect(sel *Node) {
	lno := setlineno(sel)
	if sel.Nbody.Len() != 0 {
		Fatalf("double walkselect")
	}

	init := sel.Ninit.Slice()
	sel.Ninit.Set(nil)

	init = append(init, walkselectcases(&sel.List)...)
	sel.List.Set(nil)

	sel.Nbody.Set(init)
	walkstmtlist(sel.Nbody.Slice())

	lineno = lno
}
```
