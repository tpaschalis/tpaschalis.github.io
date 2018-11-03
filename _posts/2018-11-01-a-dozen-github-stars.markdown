---
layout: post
title:  Show HN/Reddit&#58; Goof – Easily share a file over a network
date:   2018-11-01
author: Paschalis Ts
tags:   [code, project, golang, public]
mathjax: false
description: "I showed another of my toy projects to the world, on HN and Reddit. Here's what went down."  
---

[I got my first GitHub stars...](https://github.com/tpaschalis/goof) About a dozen of them, actually, along with the fork!!

I'm pretty happy about what happened this past week, and I wanted to share!   
I haven't been able to wrap up personal projects as much as I'd like, and even moreso open-source them. And this bite-sized utility was a good start to encourage me to polish and show the world some more complex stuff I've been shelving.

## What happened?
Last weekend, I finally had *juust a little* free time, and ported the ["woof"](http://www.home.unix-ag.org/simon/woof.html) Python script to Golang. It was something that I and some other colleagues had been using in a day-to-day basis.

> [Goof](https://github.com/tpaschalis/goof) (Go Offer File) is a small utility to easily serve a single file or directory over a network.    
> All you need is Goof, and a command line, or a web browser

Finally, I was pretty satisfied with the first feature set. I knew the code was not without faults, and decided to put it in front of other people's judgment. 

Of course, the point was *not* to show off, or pretending I created a snowflake. I just meant to get honest feedback, some criticism and hints on how to continue improving. Since I've started working on Golang in a more serious way, I want to learn how to write a more idiomatic style of code, as well as avoid picking up bad habits.

On October 20th, I posted it on Hacker News : [Show HN: Goof – Easily offer a file/directory over a network](https://news.ycombinator.com/item?id=18260330).   
It spent about three days in the Show HN: section, with time on the section's frontpage.

<img src="/images/show-hn-1.png" style='height: 80%; width: 80%; object-fit: contain'/>

On October 22nd, I decided to also post it on the Golang subreddit for some more feedback : [Goof - Easily offer a file/directory over a network](https://www.reddit.com/r/golang/comments/9qgb6m/goof_easily_offer_a_filedirectory_over_a_network/).   
It spent about a day and a half on the front page, and about three days where someone could find it within a click. It got about 2.5k "views", (I don't exactly know how Reddit counts views, but I think it's more or less how many unique people had reddit post rendered on their screen).

<img src="/images/reddit-1.png" style='height: 80%; width: 80%; object-fit: contain'/>

<img src="/images/reddit-3.png" style='height: 80%; width: 80%; object-fit: contain'/>

Finally, I showed it on a local tech community Slack, where some more experienced Gophers hang around from time to time.


## Here are the stats!
From GitHub traffic, we can see the two spikes, and the number of times the repo was visited.

From Hacker News, I got about 40 views, and the first 5 stars, where Reddit accounted for about 80 views and the other 6 stars, as well as one fork.
Also, there were about five clones from these two sources.

<img src="/images/gh-1.png" style='height: 80%; width: 80%; object-fit: contain'/>

<img src="/images/gh-2.png" style='height: 80%; width: 80%; object-fit: contain'/>

I understand that these are not the spectacular, mind-blowing numbers people might be used to, but for me, it was my first small project whose code I've put out in the wild, and I hope that some people could have found it to be useful!

About 130 people actually checked out the source code via the GitHub interface, (as I assume that the other people who cloned/forked the code did as well).

I'd like to think that nearly 10% of the people who visited the repo left a star.

## What did I learn..?

* First of all, simple comments can give a huge confidence boost and warm feelings; Make sure to reach out to OSS contributors more often :)   

One of the comments I got was : 
> (also this is a cool app and I will be using it)

and another person stating
> (generally speaking, it all seems good)   

 
* Secondly, I realized the importance of a README.md, a proper one with Introduction, Examples, Usage, Build Instructions or Design Decisions).  
It's the front-facing side of the project, where people will devote about 4 seconds to see what's this about. If this fails, people won't bother figuring it out themselves.    
My README.md could have been better, but I think it served that purpose well enough.

* Third, I actually had a Release! I did that after the first day of the Reddit post, but in hindsight it should have been done earlier, to allow people to try the binary's functionality directly, without bothering to install Go, or build it from source.    

## What about Go?

Here's the feedback I got. (Jotting it down here to force myself to fix it during the following weekend :p)
* Re-think the pros and cons of `http.Fileserver` and `io.Copy`.
* I could avoid starting a new server each time; Someone suggested using a mutex to lock, and then shutdown the server when `count` reaches zero.
* I should somehow 'extract' the functionality to facilitate testing; maybe using structs 
* I should probably get used to inline error handling. I.e. from `err = myZip.Close()` to `if err := myzip.Close(); err != nil {`, and generally handle error checking/returns a little better.
* Use a struct to do away with all the pointer-ing.
* `panic` in `main()` is not a great idea, as there's no `recover`
* Have some better error logging so the user can trace back issues.
* `defer` the closing of a file and the zip writer
* Maybe build a custom http server, that would be fun!


## What's next!?

Hmmm. I don't know, but I'd really, really like to hit that three-digit mark on a larger project, and maybe get on the *Trending* tab. But that's probably a goal for 2019!
