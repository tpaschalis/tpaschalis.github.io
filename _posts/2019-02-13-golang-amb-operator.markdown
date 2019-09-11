---
layout: post
title:  The `amb` operator in Go
date:   2019-02-13
author: Paschalis Ts
tags:   [go, code]
mathjax: false
description: "As if CS wasn't 'ambiguous' enough :P"
---


A few days ago, [this](http://www.randomhacks.net/2005/10/11/amb-operator/) article from 2005 was featured on HN.    

It was about the `amb` operator, which many people knew from [SICP](https://mitpress.mit.edu/sites/default/files/sicp/full-text/book/book-Z-H-28.html). It was my first time hearing about it, and as I haven't found the time to go through the Wizard Book yet, I meant to make it work in Go.

## What is the `amb` operator?

The `amb` operator was first proposed by John McCarthy (the man behind LISP) in a 1963 paper ["A Basis for a Mathematical Theory of Computation"](www-formal.stanford.edu/jmc/basis1.pdf) <sup><sub>(which is a pretty *metal* name for a paper in my opinion)</sub></sup>. The operator itself belongs to the realm of "nondeterministic computation", and `amb` stands for (surprise) *ambiguity*. 

In this paper McCarthy defines "*Ambiguous functions* : Functions whose values are incompletely specified. May be useful in proving facts about functions where certain details are irrelevant to the statement being proved."

***In short, `amb` expects a set of values as input. These are used as candidate values, to find and return the ones which satisfy a constraint.***

[Rosetta Code](https://www.rosettacode.org/wiki/Amb) has another good explanation. *The Amb operator takes a variable number of expressions (or values if that's simpler in the language) and yields a correct one which will satisfy a constraint in some future computation, thereby avoiding failure. Essentially Amb(x, y, z) splits the computation into three possible futures: a future in which the value x is yielded, a future in which the value y is yielded and a future in which the value z is yielded. The future which leads to a successful subsequent computation is chosen.*

SICP's [(Chapter 4)](http://mitpress.mit.edu/sites/default/files/sicp/full-text/book/book-Z-H-28.html#%_sec_4.3.1) which among other things, deals with non-deterministic computing, has certainly contributed to its popularity as the book has influenced a whole generation of US Computer Scientists.


A couple of [examples](http://community.schemewiki.org/?amb) might make things even clearer

```
 (let ((a (amb 1 2 10)) 
       (b (amb 0 1 2))) 
   (require (< a b)) 
   (list a b)) 

 (let ((c (amb 2 4 8)) 
       (d (amb 3 7 10))) 
   (require (eq 40 (* c d)) 
   (list c d)) 
```

In the first case, the code should return the list `(1,2)`, as they're the only values where the constraint `a<b` is met.   
In the same fashion, the second example should return the list `(4, 10)`, as they're the only pair of numbers whose product equals `40`.

In the HN discussion, there are comments stating that this can be simply defined as a list compherension, such as

```
[(c, d) | c <- [2, 4, 8], d <- [3, 7, 10], c * d == 40]
```

But the initial idea behind the operator is that of abstracted candidates (i.e. not necessarily a list), but that a search space can be an arbitrary call stack, which we'll try to explore.

One can find people who have put `amb` to creative use, such as [solving sudokus](http://continuation.passing.style/2016/02/21/solving-sudoku-amb/), to [combine concurrent sequences](http://introtorx.com/Content/v1.0.10621.0/12_CombiningSequences.html), or to find [word chains](http://www.fantascienza.net/leonardo/ar/amb_chain.html). If all this sounds intriguing, the [Mozart Programming System](http://mozart.github.io/) with the Oz language or [Icon](https://en.wikipedia.org/wiki/Icon_(programming_language)), which features "goal-directed execution"  might also be of interest.



## Where's the code?!

It's our turn to use Go and provide such an implementation! I have to admit that generics might have made our lives somewhat easier, to accustom different constraints down the line, but let's keep it simple.

Our 'constraint' is to pick `string` fragments from four sets, `A`, `B`, `C`, `D`, to construct the word "*camaroptera*", which is a bird species

```
A := []string{"bar", "foo", "cam", "kam"}
B := []string{"aropter", "ar", "ra", "baz"}
C := []string{"qux", "opter", "amar", "ra"}
D := []string{"a", "aropteraa", "ham", "optera"}
```

The correct result is obviously `A[2], B[1], C[1], D[0]`, so let's get to the action!

In many implementations, the `amb` operator is bound to the concept of [Continuation](https://en.wikipedia.org/wiki/Continuation), but in Go land, we'll use *channels* and *goroutines* instead.


Here's the code in action (also available as a [Gist](https://gist.github.com/tpaschalis/48147b36c0816a13999fa3283fe3b34f)):

```
package main

import "fmt"
import "sync"

func ambString(str []string) chan []string {
	c := make(chan []string)
	go func() {
		for _, s := range str {
			c <- []string{s}
		}
		close(c)
	}()
	return c
}

func amb(str []string, chanIn chan []string) chan []string {
	chanOut := make(chan []string)
	go func() {
		var w sync.WaitGroup
		for frag := range chanIn {
			w.Add(1)
			go func(frag []string) {
				for s := range ambString(str) {
					// Constraing block start -- Check if fragment matches the keyword
					tmp := join(append(s, frag...))
					if len(tmp) <= len(keyword) {
						// Possible candidates pass on further along on `chanOut`
						if tmp == keyword[len(keyword)-len(tmp):] {
							fmt.Println(tmp, keyword[len(keyword)-len(tmp):], s, frag)
							chanOut <- append(s, frag...)
						}
					}
					// Constraint block end
				}
				w.Done()
			}(frag)
		}
		w.Wait()
		close(chanOut)
	}()
	return chanOut
}

const keyword = "camaroptera"

func main() {
	A := []string{"bar", "foo", "cam", "kam"}
	B := []string{"aropter", "ar", "ra", "baz"}
	C := []string{"qux", "opter", "amar", "ra"}
	D := []string{"a", "aropteraa", "ham", "optera"}

	c := amb(A, amb(B, amb(C, ambString(D))))

	for s := range c {
		_, _ = c, s
		fmt.Println("************ RESULTS ***************")
		fmt.Println(join(s))
		fmt.Println("A[", indexOf(s[0], A), "]\nB[", indexOf(s[1], B), "]\nC[", indexOf(s[2], C), "]\nD[", indexOf(s[3], D), "]")
	}
}

func join(s []string) string {
	var res string
	for _, str := range s {
		res += str
	}
	return res
}

func indexOf(element string, input []string) int {
	for k, v := range input {
		if element == v {
			return k
		}
	}
	return -1
}

```

## Let's walk things through!

The main dish are the `ambString` and `amb` functions. 

The former reads in a string slice, and sends its contents along the *input channel* or `chanIn`.

The latter is the function where the operator is implemented. It uses a pair of goroutines to concurrently search through the candidate values and checks for ones that might satisfy our constraint "in the future". Those who might be suitable are passed along in the *output channel* or `chanOut` to participate in the next round of checks. The synchronization of these goroutines is achieved by the use of a `waitGroup`. There are two main components here, each of them plays a different role; the goroutines/iterators and the constraint block.

To rephrase, the operator, works from the end to the beginning and tries to find candidates from `D`, `C`, `B` and then `A`, who could construct our keyword `camaroptera`. These candidates are passed further along in channels, fullfulling the role of continuations, until the "True" candidate set is found and returned.


Here's the output of the code
```
optera optera [opter] [a]
aroptera aroptera [ar] [opter a]
camaroptera camaroptera [cam] [ar opter a]
amaroptera amaroptera [amar] [optera]
************ RESULTS ***************
camaroptera
A[ 2 ]
B[ 1 ]
C[ 1 ]
D[ 0 ]
```


## Conclusion

I hope that by now I have either piqued your interest or explained something in a clear way.

Unfortunately, I don't think the operator has seen much use in the real-world, although it is a neat theoretical idea. The whole concept of (naive) *backtracking* adds a O(2^n) term to the algorithmic complexity (and then some more on the mental outline of someone who will maintain the code), and perceived gains in productivity and elegance might not be enough offset.

Feel free to contact me via e-mail or Twitter for any comments or corrections. Until next time!
