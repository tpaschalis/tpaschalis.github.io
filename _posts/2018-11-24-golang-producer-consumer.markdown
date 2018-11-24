---
layout: post
title:  Using the Producer-Consumer concurrency pattern in Go.
date:   2018-11-24
author: Paschalis Ts
tags:   [golang, code, concurrency]
mathjax: false
description: "Go makes concurrency fun, and simple, and that makes me happy!"  
---

### Intro
I've been writing new solutions to existing problems using Go.  
[CS50](https://www.edx.org/course/cs50s-introduction-computer-science-harvardx-cs50x) exercise sets, or [`psets`](http://cs50.tv/2014/fall/#about,psets) are interesting little problems, and while the problems themselves are not huge challenges, they're fun to solve, using a language's special features to provide idiomatic and efficient solutions.


One of the problems was 'breaking' a DES hash, to 'decode' a 5-character password. The bruteforce approach is quite slow, but using simple concurrency, there was a 5x speedup, *without adding much complexity*. Here's the solution in a [gist](https://gist.github.com/tpaschalis/223045b4c50490fc950dfbc2d98d0a4f). If you have any suggestions for improvements, you can leave a comment or feel free to fork it.

### The Pattern in Go

This specific problem was a good opportunity to implement a ["Producer-Consumer"](https://en.wikipedia.org/wiki/Producer%E2%80%93consumer_problem) [pattern](https://www.cs.cornell.edu/courses/cs3110/2010fa/lectures/lec18.html).


Being familiar with `OpenMP's` idea of parallel computing, I had an overview of the various common pitfalls, and how to overcome them (resource sharing is hard, man). Go provides a much simpler interface to deal with them, and forces you into a different mindset, eg. *Do not communicate by sharing memory; instead, share memory by communicating.*

Without sacrificing too many details, here's an outline of how I implemented the "Producer-Consumer" on this case

```go
func produce(in chan string) {
    defer close(in)
    // producedPassword = somestuff....
    in <- producedPassword
}

func consume(in, out chan string, wg *sync.WaitGroup) {
    defer wg.Done()

    for s := range in {
        currentHash := crypt(s)
        // success condition
        if currentHash == hashToCrack {
            fmt.Println(currentHash,"\n")
            out <- s
        }
    }
}

func stop(out chan string, wg *sync.WaitGroup){
    wg.Wait()
    close(out)
}

func main() {

    in, out := make(chan string), make(chan string)
    wg := &sync.WaitGroup{}
    go produce(in)

    //for i:=0; i<runtime.NumCPU();i++ {
    for i:=0; i<20; i++ {
        wg.Add(1)
        go consume(in, out, wg)
    }
    go stop (out, wg)
    fmt.Println(<-out)
}
```

* So, in short. You build two `string` channels and a `WaitGroup`.  
* You can then fire up goroutines to `produce` your passwords in parallel, along the `in` channel. 
* You can spawn a separate set of goroutines, to `consume` this `in` channel.
* The `WaitGroup` is used to sync this job.
* If the hash is cracked, the successful result is passed along to the `out` channel.
* Function `stop` waits until all consumers are finished, and then closes that chanel, to signal that there was nothing found. Otherwise, the program would be blocked forever from reading the `out` channel.

### Test it yourself!
A more bare-bones implementation that you can compile right away and start picking apart, is the one below. (HINT : You will need to add a WaitGroup at some point. Can you guess when?)

```go
package main

import "fmt" 

var fin = make(chan bool)
var stream = make(chan int)

func produce() {
    for i := 0; i < 100000; i++ {
        stream <- i
    }
    fin <- true
}

func consume() {
    for {
        data := <-stream
        fmt.Println(data)
    }
}

func main() {
    go produce()
    go consume()
    <-fin
}
```

Thanks for your time, and keep on tinkering!




&nbsp;

&nbsp;

&nbsp;

&nbsp;

Some *nice* resources :   
[0] https://blog.golang.org/share-memory-by-communicating   
[1] https://blog.golang.org/pipelines  
[2] https://github.com/golang/go/wiki/LearnConcurrency  
[3] https://stackoverflow.com/questions/48233009/using-concurrency-in-nested-for-loop-brute-force   
[4] https://stackoverflow.com/questions/11075876/what-is-the-neatest-idiom-for-producer-consumer-in-go   
[5] https://medium.com/@trevor4e/learning-gos-concurrency-through-illustrations-8c4aff603b3   
