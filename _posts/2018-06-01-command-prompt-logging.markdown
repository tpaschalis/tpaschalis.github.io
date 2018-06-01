---
layout: post
title:  Log your console history using COMMAND\_PROMPT
date:   2018-06-01
author: Paschalis Ts
tags:   [bash]
mathjax: false
description: "Log your bash history using COMMAND\_PROMPT"  
---


As per the [GNU Bash doc](http://www.gnu.org/software/bash/manual/bashref.html)
```
PROMPT_COMMAND
    If set, the value is interpreted as a command to execute before
    the printing of each primary prompt ($PS1).
```

The contents of this environment variable will be evaluated as a `bash` command just before your shell displays each prompt.

With this in mind, there's a ton of cool stuff you could do! You could use it like `PS1` to spice up your terminal with additional [colors/styles](https://wiki.archlinux.org/index.php/Bash/Prompt_customization), to display git repo status (a la oh-my-zsh), the exit codes [or the timing](https://github.com/jichu4n/bash-command-timer) of the last command run, change the terminal's title bar... The possibilities are endless!

What I've been using it though, is for logging my terminal history on my dev machine.
```
export PROMPT_COMMAND='if [ "$(id -u)" -ne 0 ]; then echo "$(date "+%Y-%m-%d.%H:%M:%S") $(pwd) $(history 1)" >> ~/.logs/bash-history-$(date "+%Y-%m-%d").log; fi'
```

This prompt checks that the current user is not `root` and then appends the current date, working directory and command to a new file for each day, under `~/.logs`, where it's available for `grep`.

I first found the specific command on [this](https://news.ycombinator.com/item?id=14103688) HackerNews post, and don't know the original author. There's a heated discussion over there too, regarding the security implications of using a tool like that. I understand the `ln -sf /dev/null ~/.bash_history` crowd, and that's what I would advocate in a production setting, but the command has been really useful for some cases such as:

* Documenting the setup of a new system
* Creating a tutorial for someone
* Writing automation scripts
* Debugging someone's workflow
* [Visualizing](https://medium.com/@DavisKnuckles/visualizing-your-bash-history-using-python-7aa84ef244c0) your bash\_history and finding patterns you *should* furtherly automate
* Keep notes when I'm toying with a new tool or environment


Here's how it looks like in practice!

```
tpaschalis-Lenovo-G50-30:~:% ll ~/.logs 
total 24K
-rw-rw-r-- 1 tpaschalis tpaschalis 816 Feb 12 20:44 bash-history-2018-02-12.log
-rw-rw-r-- 1 tpaschalis tpaschalis 724 May 24 21:48 bash-history-2018-05-24.log
```
