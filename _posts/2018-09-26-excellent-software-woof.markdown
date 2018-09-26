---
layout: post
title:  What is good software? Part 1 - Woof
date:   2018-09-26
author: Paschalis Ts
tags:   [software, industry, unix]
mathjax: false
description: "What constitutes 'good' software?"  
---


Wikipedia states that 
> "The term *engineering* is derived from the Latin ingenium, meaning "cleverness" and ingeniare, meaning "to contrive, devise".

And as Software Engineers, we often think, discuss, and fight about the 'cleverness' and 'excellence' of software, in the very same way people would argue about MJ vs LBJ, Ronaldo vs Messi or Pizza vs Tacos.   
And while there is definitely no "metric" which you could use to quantify that, there's a piece of software that I've been recently using that I truly believe is *excellent*.

That software is [woof](http://www.home.unix-ag.org/simon/woof.html), by Simon Budig. It's a 600-lines Python script from 2009, used to.. send  a file over a network. The code is completely grok-able, and it's something that most people could reproduce given a spec sheet, but bear with me! It's 2018 and it's still *surprisingly unconvenient* to reliably and quickly transfer some little data from one point to another. Most of the times, you'd try one of the bazillion ways that your shop has sanctioned, end up trying to share credentials, and in the end either do it via USB or plain-old unencrypted chain-email attachments, right? :P 

Woof (Web Offer One File), creates a simple stupid webserver, and serves the file or directory you want. When the file is tranferred, `woof` quits. 

Poof. *Simple. Beautiful. Genius*. No maintaining software, long-running daemons, transferring keys, agreeing on installed software/stack, or using three different interfaces.

```
$ woof filename

Now serving on http://193.169.10.112:8080/filename
```

The enemy of all beauty and simplicity are corner cases, but when you actually plan and design, things become easier!

* Woof has an option to 'offer itself', so the other person can send you something back
* If a directory is specified, a tar of that directory is served
* It allows to customize the IP/port (to share under specific network), and how many times your file can be shared
* You can set values of preset configurations, for easier access.
* If you call it by `woof -U` it allows you to select and upload directly from your browser


I'd like to believe this is a prime example of Unix-philosophy, and good overall design

[x] It has a specified goal, it gets it done; it's clean, fast, reliable, stable   
[x] Lacks 'bloat' and gazillions of features   
[x] It is portable   
[x] Is easy to understand and use   
[x] Needs no maintenance   

It might not win any awards, use the sexiest tech, or give you another 10 imaginary CV points, but it warms a good engineer's heart, and points where we should strive towards.

PS1. I'm going to reproduce it in Golang this weekend

PS2. Before you ask, yes, I do love GNU coreutils, and you should check [Gow - Gnu on Windows](https://github.com/bmatzelle/gow/wiki)


