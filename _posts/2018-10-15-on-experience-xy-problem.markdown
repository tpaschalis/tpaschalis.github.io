---
layout: post
title:  On experience and the XY problem
date:   2018-10-15
author: Paschalis Ts
tags:   [fluff, software, experience]
mathjax: false
description: "A simple, silly story about a bug"  
---


## Setting
One unsuspecting morning, I checked into the office, and before even *making* coffee, I heard a colleague murmuring "why doesn't `ps aux` work on `$Server123`", half to me, half to himself.   
This colleague is probably the brightest sysadmin I've known; probably wouldn't be a reach to say he's one of the best around these parts.  

I sleepily ssh'ed, and tried coming back to him, saying that the same command worked for me. He didn't need almost any time to settle that it was a terminal emulator bug, which was fixed by restarting his PC, and the day went on normally.

## Conflict
Later, when I started thinking about this issue, I asked him how he determined that so fast. He said that other commands such as `w`, a bare `ps` worked, as well as redirection of the suspect command to a file. In the meantime, a big `find` would fail. The terminal emulator wouldn't scroll past line 12 or so, when a large amount of data was thrown on stdout. 


I thought this small story was very interesting because it demonstrates the power of experience, as well as the [XY](https://meta.stackexchange.com/questions/66377/what-is-the-xy-problem) [Problem](http://mywiki.wooledge.org/XyProblem),

I could *easily* imagine a well-meaning, but inexperienced developer panicking in the same situation, and either *giving up*, or developing weird, conspiratory scenarios in his head. 
* Why does it happen in my staging server, but works on my local VM??
* Is the `ps` binary broken?
* Is it a weird/rogue process that wants to hide itself by blocking `ps`?
* Oh crap, this issue is only on *my user*. have I been hacked??
* Oh hell, I can't even open vi to look at some logs
* *logs into stack overflow* Hi StackOverflow, `cat myfile | head -n 12` works, but `-n 13` freezes my server, what should I do?


I'm not saying it's a hard problem, but to arrive at a conclusion, or a satisfactory explanation of this, or a similarly simple problem one *could* need to piece together a strange puzzle; what's the difference between a shell, a terminal, a terminal emulator? What's the difference between an interactive session and a non-interactive one? What does TTY show, and what happens when you ssh somewhere?

## Resolution

This simple, silly story is an example how seemingly unimportant pieces of knowledge we've gathered over the years can tackle similarly simple, silly problems. 

It's also part of the *huge* impact of an experienced engineer can have on a team. Tunelling the enthusiasm and fresh ideas into the correct direction, steering away from potholes, and actually concentrating the efforts on directly on `X`, instead of asking `How to do Y?` can make all the difference in the world on a project

