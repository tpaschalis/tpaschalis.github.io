---
layout: post
title:  Curious bug with VS Code and Windows PATH.
date:   2018-07-25
author: Paschalis Ts
tags:   [fluff, windows, vscode]
mathjax: false
description: "How a curious Windows bug robbed me of some free time."  
---


Hey y'all.   
Small little (fluff) post of not much importance.   
Just how a curious Windows bug robbed me of 4-5 hours' worth of free time. ^^


## The start
It's no secret that I like `vi`. 

I won't try to evangelize, or tell you *why* it's the best thing since sliced bread, the pinnacle of human intellect, a minimalist's haven, the embodiment of KISS, or that it cumulates all that is good and just in the Open-Source community.


I've been working with Golang for a few weeks, on my familiar Linux+Vi combo, and following [Jon Calhoun's](https://github.com/joncalhoun) excellent [gophercises](https://gophercises.com). 
He uses VS Code with some Go extensions, and things like automatic import handling, or function definitions, made it appear appealing.


So I thought to myself, why not, let's take a shot developing on Windows! I haven't used an IDE in *years*. It can't be that bad!


## Trouble on the horizon.

I already had a portable copy of VS Code installed on my computer (Portable Apps make Windows 15% more enjoyable). I happily fired it up, got to the Extensions tab, and installed the official Go extension from Microsoft. 

After reloading VS Code, and ensuring the extension was installed, I (naively) tried to open my `main.go` file.

Alas, a foe appeared.
> 3221225477  

<img src="/images/vscode_errors12.png" style="height: 90%; width: 90%; object-fit: contain" />


## Headaches.

I thought of various scenarios that could cause issues, so I tried case-by-case examination.

- Should I just reinstall my extension?
- Is it any other extension that caused the problem?
- Is it that I'm running a portable version?
- Should I delete the `%appdata%` contents regarding VS Code (each time I did a change)?
- Should I try just setting it up via the installer, like a it like a normal human being?
- Should I try the 32-bit version?
- Should I try the insider edition?
- Should it be started as administrator?
- Should it be in a different paths, and not contain spaces?
- Is it some Electron issue, or that I don't have Chrome installed?

*Of course*, nothing worked, and *it* kept appearing.

> Extension Host terminated unexpectedly. Code : 3221225477   Signal null


There were a *lot* of people reporting the same in GitHub issues.
For measure, the error code provided over 15k results, on Google.

After searching countless of similar issues, with users mainly saying

>   "I HAVE THE SAME PROBLEM PLS HALP"   
or     
>   "Have you tried turning it on and off?"

lo and behold...

[mawenzy](https://github.com/mawenzy), had a comment, buried among (literally hundreds) of others, providing a probable solution.

<img src="/images/vscode_error_solution.png" style="height: 90%; width: 90%; object-fit: contain" />


The thing is, I've been working with the Windows console for some weeks. I have *never* had an issue with this. I don't know if I should have caught it, but I wouldn't think to look there by myself, until much later on.

***On one of the "components" that make up the Windows 10 `%PATH%` from the GUI, there was a stray `;` *in the end one line*. I removed it and everything worked perfectly. Put it back and nothing worked again.***


## End 
Well, off to write some Go code (as I originally intended). I'll probably spend some more of my free time (when I find some), banging my head against the wall to see *why* this happend, and probably report back here.


See you soon (with something better).
 
