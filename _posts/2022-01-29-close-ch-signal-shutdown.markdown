---
layout: post
title:  Signal multiple goroutines with close(ch)
date:   2022-01-29
author: Paschalis Ts
tags:   [golang, code]
mathjax: false
description: ""
---

So, this past week I learned something unexpected and fun at work. Well I learned many things, but this just ticked all the boxes of a great party trick; 
✅ short  
✅ unexpected  
✅ makes total sense if you think about it  

It's about what `close(ch)`. Do I have your attention?

### The code
Let's say you have a pool of small, long-running things that you want to stop by signalling a channel
```go
type foo struct {
	name string
	exit chan struct{}
}

func (f foo) Run() {
	fmt.Println("Starting", f.name)
	// do stuff
	<- f.exit
	fmt.Println("Stopping", f.name)
}
```

and you start up a bunch of those.
```go 
exit := make(chan struct{})
a := foo{"first", exit}
b := foo{"second", exit}
c := foo{"third", exit}
// .....
go a.Run()
go b.Run()
go c.Run()
```

_How would you stop all instances of `foo`?_

- Keep track of them and send an exit signal to each one (eg. `a.exit <- struct{}{}`)

This won't work as you'd expect; the shared channel means you don't know _which_ instance would receive the exit signal.

- Send N exit signals to the shared exit channel (`eg. for i := range pool; exit <- struct{}{}`)

This _could_ work, but requires that you keep meticulous track of how many instances are currently running. In the meantime, if an instance was already shut down, there would be no receiver and you could get a deadlock; or if a new instance was added, then it might be left running.

The best solution here is much simpler. Just call `close(ch)`!

The [Go language specification](https://go.dev/ref/spec#Close) mentions : 

> For a channel c, the built-in function close(c) records that no more values will be sent on the channel. [...]

> After calling close, and after any previously sent values have been received, _receive operations will return the zero value for the channel's type without blocking._

So `close(ch)` provides a concise way to unblock all receive operations and simultaneously provide that signal. It also is another reminder to ["make zero values useful"](https://dave.cheney.net/2013/01/19/what-is-the-zero-value-and-why-is-it-useful).

When I discussed this with a friend, he mentioned that he'd actually use a context to handle cancellation for long-running goroutines, but that's also what happens when we use [context.cancelCtx](https://github.com/golang/go/blob/release-branch.go1.17/src/context/context.go#L411).

So that's it for today! I love when I learn this kind of small and useful tidbits. Let me know if you have any other fun things around closing channels!

Until next time, bye!



## Acknowledgments
Thanks to [Robert Fratto](https://mobile.twitter.com/robertfratto) for showing me this pattern.  
Thanks to [Kostas Stamatakis](https://mobile.twitter.com/moukoublen) for pointing out the context package thing.

