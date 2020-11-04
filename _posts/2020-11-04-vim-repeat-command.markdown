---
layout: post
title:  Vim's 'single-repeat' dot command
date:   2020-11-04
author: Paschalis Ts
tags:   [vim]
mathjax: false
description: ""  
---

Okay, so first things first, let me get this off my chest.

Vim, I've been cheating on you. I know that I've said [you're all the IDE I could ask for](https://tpaschalis.github.io/vim-go-setup/) and that you were there for me when I was with a 2GB-of-ram-craptop, when I was [writing LaTeX](https://tpaschalis.github.io/macos-latex-vim/) or when I logged in that 20+ year old server.

But here's *one* of the many reasons why I'm never giving up on you.

## The 'repeat' command
When in Normal mode, you can command vim to 'repeat' *the last change* with the `.` key. 

That's all you have to know to start using it, but you can always consult the manpages by `:help single-repeat`.

The simplest of examples : say that you copy a line with `yy`. Move to some other place and paste with `p`. Move to another place, and repeat your last command (the paste) with `.`. Or use `5.` to paste it 5 times.

Or you want to indent and initialize some struct values as ints. That would be `i<tab>int ` and then `j^.j^.` to get to the start of the next line and repeat the command for these lines as well.

Or you want to change a word for substitutions, you can `cw<new-word>` and then move around and change similar words using the dot.

## Limitations
Unfortunately, the dot command does not repeat *Normal* sequences such as moving around or replaces. That means that you can't easily avoid sequences like the `j.j.` above. This is solvable by recording and using a [macro](https://vim.fandom.com/wiki/Macros) `:help recording`.

Also, to repeat an Ex command, use ` @:` for the same effect

Until next time!
