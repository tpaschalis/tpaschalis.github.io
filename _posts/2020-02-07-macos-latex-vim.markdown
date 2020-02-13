---
layout: post
title:  LaTeX + Vim + MacOS.
date:   2020-02-07
author: Paschalis Ts
tags:   [latex, macos, vim]
mathjax: false
description: "Ez pz"
---


I absolutely love LaTeX; and since I've started using a Mac for work, I wanted a setup for writing with it.

I settled with installing [MacTex](http://www.tug.org/mactex/), a free re-distribution of TeX Live which also includes Ghostcript and some mac-specific GUI applications. (Don't judge for not building from source ^^)

However, the goal was to keep using Vim; there's a bunch of plugins for that like [vimtex](https://github.com/lervag/vimtex) or [vim-latex](https://github.com/vim-latex/vim-latex), but a) I'm not going to learn or use all of their features, since I'm not writing on a day-to-day basis, and b) I didn't want to slow down my Vim for no reason. I will concede that vimtex looks very cool though, and I'll be checking it out in the near future.

So, thinking the Unix way, I thought what's my LaTeX workflow? I keep editor and document side-by-side and

a) Write  
b) Compile  
c) See PDF change  
d) GOTO a  

This is easily doable with Vim commands and a nice PDF viewer which will support auto-reloading (eg. [evince](https://wiki.gnome.org/Apps/Evince)).  
*Preview*, the default PDF viewer in MacOS, *does* support auto-reloading;

### Version 1
```vim
:! pdflatex %
```

The problem is that *Preview* will auto-reload the pdf only when you switch focus to the it, which had me alt-tab twice for every change.    
This is easily solved with the following command, where MacOS' `open -a` will open or focus the instance an application.

### Version 2
```vim
:! pdflatex % && open -a Preview && open -a iTerm
```

This way we get the compilation output, but have to press an extra enter key to return. We can solve this by using Vim's [silent](https://vimhelp.org/various.txt.html) option, which will suppress the message box that pops up and reports the result. The problem is that when you run something *silently*, if it outputs, you'll have to either use ^L or [redraw!](https://vimhelp.org/various.txt.html#CTRL-L), as your screen might be messed up.

So, let's wrap this up in a proper way; introduce a new command, *Silent*, which will silently call your arguments and then re-draw the screen; then map this to F5.
(Thanks to [this](https://vim.fandom.com/wiki/Avoiding_the_%22Hit_ENTER_to_continue%22_prompts) vim wiki article.)

### Version 3
```vim
" Custom Silent command that will call redraw
command! -nargs=+ Silent
\   execute 'silent ! <args>'
\ | redraw!

:map <F5> :Silent pdflatex % && open -a Preview && open -a iTerm
```

And voila, adding these to .vimrc covered 90% of my day-to-day workflow. How I wish all software was simple to use and modify this way!
