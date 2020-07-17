---
layout: post
title: Reaching the ceiling of single-instance Go.
date:   2020-07-18
author: Paschalis Ts
tags:   [golang, conference, transcript]
mathjax: false
description: ""
---

## Intro

Last week I had the honor of presenting *Reaching the ceiling of single-instance Go* as part of [GoWay 2020](https://goway.io).

The talk explored how our industry has been trying to modularize deployments in smaller and smaller units, what are the challenges of vertical scaling, and presented our results on exploring the limitations that Go has when scaling up. I've gradually published some of these results in previous blogposts, but some of these benchmarks are brand new.

If you're interested in chatting a bit more about scaling Go, or you think it would be a good idea to have this presentation on a conference or meetup, send an email or [hit me up on Twitter](https://twitter.com/tpaschalis_)!

## Slides

Here are the slides for the presentation. Thanks to the Beat marketing team for helping out with touchups!
[Google Sheets Slides link](https://docs.google.com/presentation/d/1p1C6mAXQTowNGDFtZ54PQewoC_Pe_VVq1rt4vZItB6c/edit?usp=sharing). 

I have another version that's more company-neutral, as well. If you have trouble accessing the slides, again, don't hesitate to reach out!

## Key Points

*Hey, sounds nice, but I don't have much time. How about you give me some numbers so I can rub them in my coworkers face?*

Well, here you Go, some bulletpoints you can reference, no attribution required!

- Go is a reasonable choice for scaling up.

- Regarding data; the hard limits on maximum allocation for slices/maps are circa ~2^47 (10^14) elements, and 2^47 bytes (~140 Terabytes).
- Regarding data; actual ceiling is memory size as it's not very easy to work with larger-than-memory datasets. You also pay in GC pauses and heap-related delays.

- Regarding goroutines; they're reminiscent of green threads. Their overhead is 2kb of stack memory which can grow and shrink as needed. The context switching is handled in the language by swapping some registers around and should be ≤200 nsec.
- Regarding goroutines; the hard limit on maximum stack is 1GB for 64-bit systems. You can launch tens of millions of goroutines (>50M in a mid-range laptop), but the constrained is the timeshare on the CPU. After 2xCPU threads, more and more goroutines will wait for their turn to use the CPU and perform some work.
- Regarding goroutines; actual ceiling is having no access to lower-level primitives (eg. NUMA awareness) for manual optimizations, the max 1GB stack, slow GC and slow scheduling.
- Regarding defers; the maximum amount of defers you can fit in a 1GB stack frame is ~4.8 million. Opening 700k files with one defer for each on an 2014 laptop takes <4.5 seconds, ~220MB of RAM.

- Regarding net/http; all requests are launched into new goroutines before reaching your handlers. Thread safety is not guaranteed, so you should pay attention. net/http can handle up to ~90-100k requests per second on a 4-core machine. For comparison, StackOverflow, the 40th most visited website worldwide peaks at about 4k req/sec.

- Regarding the Garbage Collector; has increased by leaps and bounds since the Go 1.4/1.5 days. Current iteration is a concurrent mark-and-sweep GC.
- Regarding the Garbage Collector; Service-Level Objectives dictate its behavior. Sub-nanosecond pauses (max 500µsec pause/cycle, usually lower than 100μsec), 25% of available CPU cores used, GC cost kept in linear proportion to the allocation cost. 

- Regarding channels; maximum size of message (channel data type) is 64kb. Maximum buffer size and allocation are similar to the data limits i.e. circa ~2^47 (10^14) elements, and 2^47 bytes (~140 Terabytes).
- Regarding channels; maximum ~18-20 million messages/sec for buffered, ~5.5 million messages/sec for unbuffered ones. If that's too restrictive, a Mutex lock/unlock is ~4 times faster.
- Regarding channels; limited by memory size, as allocation happens up-front for buffered channels.

## Transcript

### Intro
Hello everyone! Let me introduce myself. My name is Paschalis Tsilias, I'm a physicist-turned-software engineer, currently working as a Backend Engineer at Beat.

I'm here to talk about "Reaching the ceiling of single-instance Go" and I hope we can spark a new conversation about Go and vertical scaling.

So today, I'd like to start with what motivated me to talk about this subject; briefly mention the challenges of vertical scaling as we all know them, discuss some of Go's inner workings, and present our results.


### Motivation

Mainframe → VPS → VMs → Containers → K8S pods → AWS Lambda

Our industry has been steadily moving towards modularizing deployments and orchestration in smaller and smaller units. From the 1960’s enormous mainframes, to the on-premise server racks starting around the 90s, the explosion of virtualization, and the rise of Cloud providers and Virtual Private Servers.

In parallel, things like chroot, jails and LXC were being worked on, leading to the current state of Docker and Kubernetes.

The industry used solutions like Docker Swarm or Kubernetes to try to make sense of the whole orchestration process; and enable cloud-native deployments and microservices. Naturally, almost no-one deploys on bare metal anymore, and for good reason -- cloud is now the default; the promise of no upfront cost and magic autoscaling, both to save money and handle usage spikes is still very appealing.

But we have moved to yet another layer of minification, with things like serverless, lambda, or azure/gcp functions, where the boundaries of microservices and functions are a little fuzzy.

Hardcore vertical scaling is nowadays mostly prevalent in research facilities HPC and supercomputers; horizontal scaling is getting to be the first thing that comes to mind. And currently, we’re now back to counting CPUs in milli-cpu and memory in megabytes.

I have great respect for people who pushed their contemporary tech to the their limits. For example, the *whole Pokemon Red game was 11 MB*. So, I just wanted to understand the strengths and limitations of our stack. 

Go is a naturally great choice for distributed computing scenarios; so let’s find out how Go fares on a single-machine contexts as well!

### Challenges of vertical scaling

Of course, vertical scaling has some challenges; there’s a reason why our industry has tried to move away from it

- Hard-er limits on scaling possibility
- Greater risk of widespread outages and hardware failures
- Proneness to building snowflakes
- Vendor lock-in / Cost mean very little room for adjustment to bad decision-making

But, the thing is that some of these are not problems inherent to vertical scaling; but may be engineering culture problems, so let’s not forget some of the advantages that come with it.
StackExchange is a great success-story for vertical scaling. The 40th most visited website worldwide runs on less than a dozen web servers, and just 4 database plus 2 cache instances.


### Go 
So, it’s time we dive into Go, and how we can reach that ceiling.

Along these slides you’ll find references to the Go codebase, so that you can read for yourselves afterwards, as well as external resources, which were helpful in our investigation. 

### Go - Data
First off, let’s talk about our usual suspect; data!
Data is one of the few points where there’s some hardcoded limits that dictate the Go behavior.

Go is pretty reasonable in the way that it handles data, there’s little overhead to the structs that pass data around, but some details are a little hidden; it’s for example it might not be immediately apparent if some data is allocated on the stack or the heap without digging around.And talking about the heap, slices are common paths for inefficient allocation -- unless the compiler knows the size of the slice at compile time, the backing arrays are heap allocated. 

So what are those limits? Trying to allocate a very big slice, an array, a channel or a map might produce an error like `makeslice: len out of range`. It’s easy to dig into the code and see how the constraints are implemented it’s a two-part check. 

The first check is a comparison with maxAlloc,the maximum allowed allocation, which on a 64-bit linux system is 2 to the 47th power.

The second check is a pointer multiplication, so that the datatype and length don’t overflow 64 bytes.

Here, you can see the maximum slice length for some of the basic data types.

So these are the hard limits; Around 10 to the 12th elements and 140 Terabytes.

Other than that, we pay in larger GC pauses, and delays due to heap allocations; we also are limited by the size of our memory as Go is not very easy to work with larger-than-RAM datasets.

So, what do we do when we want to work with data of such gigantic proportions? Well, split the workload! Which leads us to .. 

### Go - Goroutines

One of the things that might be a bottleneck on scaling up is concurrency; after all there’s so many cores a single machine may have.

Let’s see how the Go concurrency model, using Goroutines can help us scale up.
Goroutines are a little different than CPU threads, OS threads or Green threads that exist in other languages. They are based on old ideas like SCP and have a number of advantages
1) Τhey are materialized in Go code themselves.
2) They are extremely lightweight, having an overhead only 2kb of stack which grows and shrinks as required
3) The scheduler is able to pool goroutines, and turns IO/Blocking work into CPU-bound work.
4) Context switching is cheap; the scheduler takes care of it, just swapping a few register around
5) The scheduler is aware of the goroutine states as well, and can make informed, smart decisions about their scheduling.

Compare to threads in other stacks, like Java where a thread would take up 1MB, that’s 500 times more, and it also wouldn’t have the flexibility of the variable stack size. And on the other hand, having hundreds of thousands of OS threads would take you to kernel-tweaking-settings land, where weird things are happening.
I'd also like to mention that this scheduling is non-deterministic, and the preemption, the allocation of goroutines onto CPU threads is done on runtime.

So, how many of them you could launch? You can just make it rain; On a mid-end laptop you could go up to 50-60 million goroutines easily. 

The ceiling is having no access to lower-level primitives (like NUMA awareness) for further manual optimizations, and the 1GB of stack per goroutine.
You actually end up paying in slower GC, and well, slower scheduling when you need all your goroutines to actually do something interesting. As we mentioned, the context switching is actually cheap, and goroutines in the pool can be reused, so one estimate is the number of goroutines that are required until your CPU time is all tied up.

Our experiment had us create and launch two types of goroutines; with 1 and 10 millisecond workloads. 
Using pprof we visualize the execution and see until which point goroutines spend considerably more time idle than performing some actual work.

Our results show that after twice the CPU threads, the returns are diminishing.

### Go - Defer
Let’s move on with defers; the main way of controlling system resources like file descriptors, sockets and network connections in a predictable way.
Many people used to avoid defers due to their performance cost. In Go 1.14, this has changed; the cost of a certain kind of defers was greatly lowered via inlining them. 

These optimizations take part in what they call “Open-Coded” defers, that is defers not living in a for-loop, thus they’re predictable and can be optimized. This has lead some people to say that defers are now almost a zero-cost abstraction, and while an exaggeration, there’s good reason for saying that as you can see in this benchmark.

The defers are stacked on a defer chain in a Last In First Out order, and like goroutines they are handled as part of the language itself. 

They can be weaved and live either on the stack or the heap. One can go about exploring the `_deferw struct itself, as well as the way that new defers are created and pushed onto the defer chain, as well as returned on the exit point of a function.

Our test consisted of opening a huge amount of files, and having one defer for each; we went on to see the performance of the defers.

The maximum number of defers you can fit in a 1GB stack frame (in a single goroutine) is ~4.8 million. Opening as many files as the linux kernel would allow on a 2014 linux laptop (~700k), creating and executing as many defers took <4.5 seconds, and ~200MB of memory.


### Go - net/http
It’s time we talked about net/http; Net/http is in my opinion one of Golang’s strongest cards.
Almost half of the workforce uses it directly; and many third-party packages are just wrapping it as well.

To be honest, I had a gut feeling that I wasn’t really constrained by net/http’s performance. I looked around to see if any fancy benchmarks already exist. The numbers are… more than enough. Remember that StackOverflow peaks at about 4k requests per second. Which is about the point, where these numbers stop losing their true meaning, and it’s just evolves into a silly contest. 

If this performance is not enough for your needs, frameworks like jlhttp and fasthttp, and brag about an order of magnitude better performance.

So what is net/http using under the hood?
Well, the code is easy to spot! All incoming requests are accepted and launched into a new goroutine before they reach your handlers. So the goroutine section can give us some hints; i.e. performance is pretty good.

Our experiment consisted of a web server with three routes.
One route just set an HTTP-OK header, another performed a light computation, and one served a small file.
We had an external server load perform a flood of requests each on one of three routes randomly, measuring latency and status codes up until to the point that 1% of our requests failed, or that p90 was over 200ms.

The results show that net/http would be able to handle StackOverflow’s PEAK 4.000 request per second.


### Go - GC

Garbage collection is to many programmers a black box; I think it’s fair to say that it sometimes is more of an art than a science.
Go has taken a different direction to other languages in regards to garbage collection, moving away from knob-tweaking and manual work; opting to prioritize low latency and simplicity.

Go being a GC-language means it’s not the best fit for real-time performance, like financial sector operations or audio mixing; but I’m fairly confident in saying that the GC is one of the best around. But, you get what they give you, there’s no room for more improvement unless you submit a change yourself.

The garbage collector has increased by leaps and bounds from 2014, when it was routine to have 40-50 millisecond pauses.

The current iteration works as a “Concurrent mark and sweep GC”.
There are Service Level Objectives in place, which dictate its behavior. 
Namely, 
- Maximum 500 usec STW pause per GC cycle (actually lower than that)
- 25% of CPU during GC phase
- GC cost kept in linear proportion to the allocation cost

Let me repeat that; 
These types of latencies are reaching the point where they can happen for non-GC reasons.
As Rick Hudson says: “You don't have to be faster than the bear, you just have to be faster than the guy next to you."



### Go - Channels/Mutexes
Main form of communicating in a thread-safe manner, and one of Go’s selling features.
It’s a great abstraction, allowing for asynchronous communication, and enable implementation of event-driven logic.
Channels can transfer a single type of message, and they come in two flavors, buffered and unbuffered; this ‘flavor’ dictates its behavior in regards to blocking on whether the message has been received or not.
Underneath, a channel contains a lock wrapped in a threadsafe queue; the struct itself is straightforward; the chan.go file contains the code for receiving, sending, as well as support for the ‘select’ structure, I recommend you check it out!

The main limitations, come from just 8 lines. These lines implement the limits on the message type size, the buffer size, as well as the maximum allocation.
The code is very similar to what we saw on the Data section. These hard compiler limits are 64kb for the message size, about 10^12 buffer size, and again, the maximum allocation on a 64-bit linux system comes at 17.5 Terabytes. The main limiting factor, is again, memory, as the send queue in buffered channels is allocated upfront, and there isn’t an easy way to work with larger-than-memory datasets.

One more thing I wished to explore, is how fast are channels really? As it happens, the channel struct itself has little overhead, and the main limiting factor is the goroutine context-switching.

The benchmarks show around 50ns for sending and receiving in an buffered channel, and around 180ns for sending and receiving in a buffered channel.
This discrepancy is due to the way that Go has to obtain the locks; and puts an upper limit to the number of messages you could send, which on a MacBook comes up to 20 million messages per second and 5 million messages per second for buffered and unbuffered channels respectively.

Just to have some sense of magnitude, a naive copy takes about half a nanosecond on the same machine, and a mutex lock-unlock is about 4 times faster, without having to move data around that much.
So if you ever start hitting those limits, it’s something to keep in mind.



### Outro 
That said, I think it’s time we wrap this up.

I hope that you’ll agree with me saying Go is a reasonable choice for scaling up.
My bold prediction is that its simplicity and performance will spearhead Go’s adoption into other fields in the next few years. High-Energy Physics is a field ripe for the taking, as Fortran developers are getting fewer, Python is too slow, and C++ is too dangerous/difficult.
The main limiting factor in scaling up is not the language, but the host system, and kernel, few hard limits exist.
Concurrency and HTTP serving are pretty strong.
GC is enough for most cases, but it’s where caution should be placed, by monitoring and profiling your large-scale applications.


