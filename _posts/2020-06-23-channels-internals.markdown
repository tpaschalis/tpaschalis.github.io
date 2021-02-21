---
layout: post
title:  What happens when you `ch <- foo` in Go? 
date:   2020-06-23
author: Paschalis Ts
tags:   [golang, internals]
mathjax: false
description: ""
---

The selling point of anyone learning about Go, *has* to be the concurrency model.

Actually, I'd have a hard time believing someone would pitch using Go for a project and not mentioning *channels*, the mysterious messenger that allows goroutines to communicate in a thread-safe way.

So, after exploring what [goroutines](https://tpaschalis.github.io/goroutines-size/) and [defers](https://tpaschalis.github.io/defer-internals/), it's time we took a look under the hood, to see what is actually behind a channel, and what happens from the moment you `make(chan string)` until you receive your values from a different goroutine.

As before, for this exploration I'll be using the [Go 1.14 release branch](https://github.com/golang/go/tree/release-branch.go1.14), so all code snippets will point there. So grab some coffee or a cold beer and let's get started!




## Resources

https://stackoverflow.com/questions/19621149/how-are-go-channels-implemented
https://docs.google.com/document/d/1yIAYmbvL3JxOKOjuCyon7JhW4cSv1wy5hC0ApeGMV9s/pub
https://codeburst.io/diving-deep-into-the-golang-channels-549fd4ed21a8
https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html




## Scheduler Resources
https://morsmachine.dk/go-scheduler
