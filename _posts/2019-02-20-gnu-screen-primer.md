---
layout: post
title:  GNU screen primer in 4 minutes.
date:   2019-02-20
author: Paschalis Ts
tags:   [linux, screen, tutorial]
mathjax: false
description: "GNU to the rescue"
---

It would be really interesting if one could define an abstract metric and rank software on how much $$$ it has help generated, or man-hours saved. A silly metric to be sure, but nevertheless, I'm pretty sure that `GNU Screen` would be up on top, with the best of them.

I recently had the opportunty to re-[introduce](https://xkcd.com/1053/) a colleague to this wonderful utility, and wrote this small 4-minute primer on it. There are literally thousands of "tutorials" on the internets on how to use `screen`, many better and many worse but a) this one is mine and b) this is a complete cheatsheet for anyone that just wants to *get shit done quickly*. 


### Intro - TL;DR

[GNU Screen](https://www.gnu.org/software/screen/) is a command-line *window manager*, that allows multiple separate login sessions inside a single terminal window. It allows to create, detach and re-attach "screens", where you can run processes independently from each other. Processes continue to run even if their "screen" is not currently visible, is detached from the user's terminal or the user has lost connectivity.

[Screen's](https://en.wikipedia.org/wiki/GNU_Screen) most common usage is to separate processes from a specific Unix shell to *long-lived sessions*, so that the process can run even when the user is disconnected. The user can log back at any time and resume the session, since applications won't even know the terminal has been detached.



### On host machine
These are commands that you can run from your normal shell to handle screen sessions

```
screen -ls                  # List screen sessions
screen -S <name>            # Create <name> screen session (mnemonic -S=string) 
screen -r <name>            # Resume <name> screen session
screen -r -d <name>         # Force detach and resume already-attached session
screen -X -S <name> quit    # Destroy <name> screen session
screen -X -S <name> <cmd>   # Run <cmd> in <name> screen session
```


### Inside screen session
The default "leader", the key that will escape your keystrokes from the terminal, and interact with the session itself is `Ctrl-A`, or `^A`.

Once inside a screen session, you can use the following basic commands.

```
^A d        # detach
            # Careful, do NOT use ^D, as it will terminate your session by default
^A n        # Goto next screen
^A p        # Goto previous screen
^A k        # Kill current screen
^A ?        # Get Help page/keybindings

echo $STY   # Get name of current screen session (unset variable if not in screen)
            # useful for bash scripts.
```

### Scrolling bar history

One of the things that jar people using `screen` for the first time is the inability to use the mouse to scroll back.

Do not fret! Simply add the following to your `~/.screenrc`.
```
# Enable sensible mouse scrolling and scroll bar history
termcapinfo xterm* ti@:te@
```

You could also simply use `| more` and the like, or `^A <Esc>` and then scroll normally with your arrow keys and mouse, pressing `<Enter><Enter>` to return to your prompt.


### First Gotchas

If it's your first time going over `screen`, here are a couple of beginner pitfalls.

* If you try to create a screen using `-S` and a name that already exists, it will terminate the current and create a new one
* First couple of times you will inadvertently use `Ctrl-D` instead of `Ctrl-A d`, which will *kill* the screen session instead of detaching it. That's why I have added `set -o ignoreeof` on my `~/.bashrc`.
* Speaking of, reattaching/resuming a screen session *will not* source your `~/.bashrc` or shell equivalent (you probably expected that, since it's the same login session)


### Closing words

There are peoiple who have been doing [*absolutely*](http://www.softpanorama.org/Utilities/Screen/screenrc_examples.shtml) [***crazy***](https://bbs.archlinux.org/viewtopic.php?id=55618) things with their `~/.screenrc` and customized setups. As a fellow nerd I understand and have myself indulged in that. But when you're ssh-ing to multiple machines per day, don't get ahead of yourself, it's maybe better to start by keeping things simple, and customize when needed.   
Also many people swear by [tmux](https://dominik.honnef.co/posts/2010/10/why_you_should_try_tmux_instead_of_screen/) (not me personally), so that's another cool thing to look at.

You can stop reading at this point, as this is enough information for everyday usage. If you want some extra juice though, feel free to continue.



<br>
<br>
<br>
<br>
<br>
<br>

### Windows inside a screen
Up to now, we've talked about *screens*. We created separate screens for each process we want to maintain. But as we mentioned, `screen` is a *window manager*, so of course you can use a single screen to maintain separate windows! It's exactly how you imagine it, and works just like your favorite graphical Desktop Environment. The following commands provide the "Alt-Tab" functionality that you're familiar with.

Once inside a screen session
```
^A c        # Create one new window
^A "        # Display window list (use arrow keys to navigate to another one)
^A A        # Rename current window
^A k        # Kill current window

^A n        # Move to next window
^A <space>

^A p        # Move to prev window
^A <backspace> 

^A <num>    # Move directly to <num> window

^A ^A       # Flip between two last active windows
```

Of course, as in a physical screen, or a graphical Desktop Environment, you display multiple windows, split your screen to dedicate space to each one, cycle between them. `GNU Screen` is no exception! You can find images of *really* cool setups using this, but after a point you're probably better off using i3wm.

```
^A S        # Split display horizontally
^A |        # Split display vertically
^A <tab>    # Cycle to the next region
^A X        # Remove current region
^A Q        # Remove all regions but not the current one.
^A :resize +<num>   # Resize current region by +<num> or -<num>
```

And if you're hardcore enough to use no mouse at all

```
^A <esc>    # Use default screen scrollback buffer
<space>     # Toggle selection to copy
^A ]        # Paste aforementioned selection
```
