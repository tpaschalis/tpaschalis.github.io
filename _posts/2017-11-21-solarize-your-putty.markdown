---
layout: post
title:  Solarize your PuTTY (or pick your own colorscheme).
date:   2017-11-21
author: Paschalis Ts
tags:   [tutorial, windows]
mathjax: true
description: "Modify the registry to avoid the horrendous PuTTY default colorscheme"  
---

Many people either prefer to or have to work with Windows. (IMHO Windows are on a right track to being more developer-friendly, especially with the Linux Subsystem, but well...).

The thing is, as a developer, you're probably spending a LOT of time ssh-ed into other machines, and as far as I know, there's no built-in ssh client into Windows.

Enter [PuTTY](https://www.chiark.greenend.org.uk/~sgtatham/putty/). It does *so many things* right, but the default colorscheme on low-brightness screen can be painful. Let's change this!

<figure>
	<img src="/images/puttycolor2.png" style='height: 100%; width: 100%; object-fit: contain'/> 
	<figcaption>(You can play around with your screen's brightness)</figcaption>
</figure>


### Get to the point!

Open `regedit` and `Ctrl+F` the `SimonTatham` key.

Right click, save the registry file somewhere, and open it with your favorite text editor.

Modify by pasting [your favorite colorscheme](https://github.com/altercation/solarized/blob/master/putty-colors-solarized/solarized_light.reg) in the appropriate section.

```
Windows Registry Editor Version 5.00

[HKEY_CURRENT_USER\Software\SimonTatham\PuTTY\Sessions\YOUR_SESSION_NAME]
"Colour0"="101,123,131"
"Colour1"="88,110,117"
"Colour2"="253,246,227"
"Colour3"="238,232,213"
"Colour4"="238,232,213"
"Colour5"="101,123,131"
"Colour6"="7,54,66"
"Colour7"="0,43,54"
"Colour8"="220,50,47"
"Colour9"="203,75,22"
"Colour10"="133,153,0"
"Colour11"="88,110,117"
"Colour12"="181,137,0"
"Colour13"="101,123,131"
"Colour14"="38,139,210"
"Colour15"="131,148,150"
"Colour16"="211,54,130"
"Colour17"="108,113,196"
"Colour18"="42,161,152"
"Colour19"="147,161,161"
"Colour20"="238,232,213"
"Colour21"="253,246,227"
```

Save the modified `.reg` file, right click on it and `Merge`. You can also keep the modified registry file to reproduce on another machine


### Hey, I use [X]..!
Of course, a bunch of different tools exist. But corporate environments restrict users (for good reasons, too). PuTTY is well established, and that means well audited and documented.

IMHO, PuTTY is a golden standard for modern software, in the sense that most modern software projects can learn a thing or two from it. It was released 18 years ago, comes bundled with a lot of useful stuff, supports several network protocols, and foremost *is FOSS, simple, and just works*. I'll make sure to buy Simon a beer if I ever meet him.
