---
layout: post
title:  TDD experiment - The Setup.
date:   2020-02-11
author: Paschalis Ts
tags:   [tdd]
mathjax: false
description: ""
---

There's no better moment than a quiet weekend for a small side-project to begin. This time, the focus was not on the code itself, but on a challenge, and a small experiment. This post describes the experiment's first phase, the setup.

I've been supposedly doing *some* TDD, but a friend suggested completing a project using proper, dogmatic, hardcore TDD, just to see how it goes. That means no cutting corners, no cheating with writing a method or two and testing later, and *absolutely no* early refactors.

I decided to use Uncle Bob's rule verbtim:
```
- You are not allowed to write any production code unless it is to make a failing unit test pass.

- You are not allowed to write any more of a unit test than is sufficient to fail; and compilation failures are failures.

- You are not allowed to write any more production code than is sufficient to pass the one failing unit test.
```

And so, the *2020 TDD experiment* has begun. This is the first post, where the experiment is set up.

In a next post, I'll release the project (a CLI app, nothing groundbreaking), and gather some feedaback, while later I will hopefully write down my experience and observations. I don't believe in panaceas and silver bullets, but spending time with a tool in *deliberate practice* is in my opinion the best way to improve one's craft.


So, for this project, which is going to be a small-ish CLI application, *all commits will strictly correspond to TDD steps*.

Unfortunately, I don't have one of Uncle Bob's TDD hats
<center>
<img src="/images/uncle-bob-tdd-hat.png" style='height: 40%; width: 40%; object-fit: contain'/>
</center>


so settled for the next best thing

```bash
function red(){ git add --all && git commit -m ":x: $1" }
function green(){ git add --all && git commit -m ":white_check_mark: $1" }
function refactor(){ git add --all && git commit -m ":cyclone: $1" }
```

My commits now look like *this*, so I can keep track of what step I'm currently in and what I should do next.

<center>
<img src="/images/tdd-experiment-commits.png" style='height: 40%; width: 40%; object-fit: contain'/>
</center>

After a couple of days, I can admit it's hard to keep up (for now) and I'm continuously switching between
a) jesus, I'm wasting so much time, wth   
b) okay, *that* came out nice   

But I'll reserve judgement until I get the code to a more advanced state where TDD will be either a bigger hindrance or a bigger aid.



