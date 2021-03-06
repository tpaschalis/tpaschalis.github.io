---
layout: post
title:  Log your console history using COMMAND_PROMPT
date:   2018-06-01
author: Paschalis Ts
tags:   [bash, sysadm]
mathjax: false
description: "Log your bash history using COMMAND_PROMPT"  
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

I first found the specific command on [this](https://news.ycombinator.com/item?id=14103688) HackerNews post, and don't know the original author. There's a heated discussion over there too, regarding the security implications of using a tool like that. I understand the `ln -sf /dev/null ~/.bash_history` crowd, and that's what I would advocate for in a production setting, but the command has been really useful for some cases such as:

* Documenting the setup of a new system
* Creating a tutorial for someone
* Writing automation scripts
* Debugging someone's workflow
* [Visualizing](https://medium.com/@DavisKnuckles/visualizing-your-bash-history-using-python-7aa84ef244c0) your `.bash_history` and finding patterns you *should* automate
* Keep notes when I'm toying with a new tool or environment


Here's how it looks like in practice!

```
tpaschalis-Lenovo-G50-30:~:% ll ~/.logs 

-rw-rw-r--. 1 tpaschalis tpaschalis     44 May 16 11:07 bash-history-2018-05-16.log
-rw-rw-r--. 1 tpaschalis tpaschalis 110921 May 21 17:26 bash-history-2018-05-21.log
-rw-rw-r--. 1 tpaschalis tpaschalis 101026 May 22 19:07 bash-history-2018-05-22.log
-rw-rw-r--. 1 tpaschalis tpaschalis  99762 May 23 17:46 bash-history-2018-05-23.log
-rw-rw-r--. 1 tpaschalis tpaschalis  37052 May 24 18:08 bash-history-2018-05-24.log
-rw-rw-r--. 1 tpaschalis tpaschalis   3618 May 25 15:38 bash-history-2018-05-25.log
-rw-rw-r--. 1 tpaschalis tpaschalis    222 May 29 10:35 bash-history-2018-05-29.log
-rw-rw-r--. 1 tpaschalis tpaschalis   7275 May 30 18:04 bash-history-2018-05-30.log

tpaschalis-Lenovo-G50-30:~:% more bash-history-2018-05-30.log

...
...
2018-05-15.13:24:02 <some-dir>tpaschalis/foo   435  mkdir ../src/qtest
2018-05-15.13:24:05 <some-dir>tpaschalis/foo   436  touch ../src/qtest/qtest.go
2018-05-15.13:24:07 <some-dir>tpaschalis/foo   437  vi ../src/qtest/qtest.go
2018-05-15.13:24:40 <some-dir>tpaschalis/foo   438  go build qtest
2018-05-15.13:24:43 <some-dir>tpaschalis/foo   439  ./qtest
2018-05-15.13:24:46 <some-dir>tpaschalis/foo   440  go clean qtest
2018-05-15.13:24:50 <some-dir>tpaschalis/foo   441  cd ~
2018-05-15.13:24:51 <some-dir>tpaschalis/foo   442  ll	
```
