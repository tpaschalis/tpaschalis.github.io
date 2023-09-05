---
layout: post
title:  Avoiding papercuts
date:   2023-09-05
author: Paschalis Ts
tags:   [meta, programming]
mathjax: false
description: "Why deal with a thousand papercuts when you can just throw the pile away?"
---

```vim
nnoremap dgg <nop>
command References GoReferrers
```

Focus is so important in our craft. I believe everyone can recall _that_ magical half-hour with no distractions where everything suddenly fell into place.
Getting "into the zone" is hard, and if you ask most people, they'll likely say that the biggest killer is external stimuli; coworkers hopping over for a quick question, Slack notifications or even the temperature not being quite right.

Yesterday I realized there's another parameter in play, that is my own practice and habits. Let me explain:

My main editor/ide/development environment is Vim.

I have my Vim set up ~just the way I like. Most of my daily use is muscle memory at this point, but there are always a couple of minor annoyances.
First, when hitting `gd` to **g**o to **d**efinition, I'd sometimes hit `dgg` instead (in vimspeak, that means delete from the cursor to the beginning of the file).
I must have truncated and re-fixed around two bazillion files this year, with results ranging from having to stop what I was doing and spam the Undo button a bunch of times, to realizing it much later when compiling failed.
Second, I use the `:GoReferrers` command _a lot_. I saw a teammate had remapped it to `:References` instead, which is both easier to type/autocomplete and removes one brain 'hop' from "what are the _references_ to this object?".

It doesn't sound much in the grand scheme of things, but I'm sure if I look, I'll find more examples of tiny papercuts that had me groaning again and again but actually would also just take 10 seconds to fix.

So, take some time to fix your own practice; I promise that the time spent will more than make up for the gradual buildup of grievances these issues are bringing into your day!
Is it a mouse wheel that needs cleaning? A USB-C cable that keeps disconnecting? A flaky test that shows up at the most unfortunate moment? Stop what you're doing, fix it and you'll never have to worry about it again.

Until next time, bye!
