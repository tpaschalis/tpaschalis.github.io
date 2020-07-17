---
layout: post
title: I was on GoWay Fest!
date:   2020-07-16
author: Paschalis Ts
tags:   [golang, conference]
mathjax: false
description: ""
---

## Intro

This past week I took part in the 2020 [GoWay conference](https://goway.io). It was my first conference contribution in 3 years (and my first ever on Software, and not Physics).

This year's GoWayConf iteration was held online; the lineup was out of this world, with talks from Ellen Körbes, Dave Cheney, Mat Ryer, Mark Bates, and other kickass gophers. I feel very lucky to have been a part it all, and I've set a goal to participate *in at least one more* in 2020!

So I'd like to express my gratitude to Stanislav, Maria, Mike, Panagiotis, Kostas, Vangelis, Athina and Giorgos; you all helped in some way!

Here's what went down during my parts!

## The conference

About two months ago, a co-worker mentioned that he was helping organize a conference, and that there were a few slots for new talks. After a few days, I got an idea and after validating it with a couple of people, I hit him up with the proposal.

The idea sounded cool to some people, and my talk about *"Reaching the ceiling of single-instance Go"* was approved. After some discussions with the committee, they also asked whether I wanted to participate in a Panel Discussion around the *Adoption of Go*. I figured it was a fantastic opportunity to share our story at Beat, and learn how this process went for other large organizations, so I said yes!

## The Panel Discussion
Prior to this, I had never been on a panel discussion, so I didn't know what to expect.
Our little virtual room included [Federico Paolinelli](https://twitter.com/fedepaol), [Iskander Sharipov](https://twitter.com/quasilyte), [Mike Kabischev](https://twitter.com/mkabischev) and [Mat Ryer](https://twitter.com/matryer), along with roughly a hundred viewers.

We hit things off by comparing our experiences on adopting Go, in Intel, RedHat, Beat, and various other startups. What we all agreed on, is that the best way to kick-off adoption is to solve an actual problem. It can be a self-contained new business requirement, an infrastructure tool, or a simple CLI application. All you have to do to set the spark is say *look, there was this problem. I solved it in X days in Go and looks pretty cool, check it out!*. That's usually enough to get people going!

We mentioned some of Go's killer features that can help adoption. Namely, *gofmt* and generally the *great developer toolchain*, the compatibility promise, an easy deployment method, the fact that there's usually *one good way* of doing things, great performance, short feedback loops, built-in testing and benchmarking. 

What's also interesting is how easy is for newcomers to pick up Go! Surely, your code might look like PHP-esque or Java-ish for a few weeks, but the idiomatic way of writing good Go code is usually presented to you. For example, Beat doesn't hire exclusively people for their Go background, but rather for strong engineering principles.

Finally, and responding to questions from the crowd, we all agreed that a really underrated metric is *developer happiness*. All in all, Go's simplicity, straightforwardness and buzzing community play *a vital role* in having happy developers, who can solve challenging problems for you given the right toolset and independence.

Generally, I think the panel discussion was a success! There was great chemistry and crowd interaction, and no long awkward pauses like I feared ^^

## The Presentation

The presentation went a little ...less swimmingly. 
A few minutes before starting, the organizer proposed testing the screen sharing. And while it *did succeed* on previous tries, Firefox refused to display my slides. In panic, I tried to switch to Chrome, and logged in juuuust in time.

I swear I may have rehearsed the whole thing over 50 times, but stage fright took over, and I wasn't able to enjoy the whole thing as much as I'd like. If you ever get to know me personally, when I'm nervous I talk *really fast* and *mumble my words*. Which, are all problematic in the context of a presentation. :D 

Ultimately, I sprinted through 30 minutes of content in about ~16-17 minutes. The final nail in the coffin was that people complained some of the slides were blurry, probably due to some network issue. 

So, I'm convinced that Murphy's Law holds, that is *"Anything that can go wrong will go wrong"*.

Nevertheless, I had some good interactions with the crowd of ~70 people; good, thoughtful questions, ideas on new things to benchmark and how to improve current ones.

I'll be sharing soon a new blog post that will contain the slides along with a transcript of the presentation, and my key results on Golang and vertical scaling!

What I was most happy for was that I was able to meet [Felix Geisendörfer](https://twitter.com/felixge). He strikes me as one of the humblest, most passionate and smart engineers I've met! He had great words of encouragement, shared his own thoughts on my results, and pointed possible ways to move them forward. I found out later that he was one of the early NodeJS core maintainers, and has a *bunch* of kick-ass projects under his belt, with the most recent being [fgprof](https://github.com/felixge/fgprof), you should definitely check it out!

## Outro
All in all, it was a really great experience. I rubbed shoulders with some pretty great engineers, and came out with newfound motivation and desire to do more! I suppose I'll see some of y'all in a conference soon, hopefully a physical one!


