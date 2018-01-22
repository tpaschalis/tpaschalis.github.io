---
layout: post
title:  Show HN&#58; Tweetstream – Get a user's tweets, between two dates
date:   2017-03-09
author: Paschalis Ts
tags:   [code, twitter, project, public]
mathjax: false
description: "I showed one of my toy projects to the world, on HN. Here's what went down."  
---



I finally put one of my hobby projects forth to HackerNews!

### TL;DR What happened?
Yesterday, [tweetstream.space](http://www.tweetstream.space/) sat at the #3 spot of Show HN for half a day, while reaching at #2, and spending some days on the Show HN frontpage. It peaked at #39 of the main HN ranks.

During the first hours of it being up there, the site received ~2200 visits, had 66 timeouts :(, and was able to serve more than 97%  of the users requests. I received comments about bugs, recommendations, people poking around the website, along with a comment that included
> Cool app though!

### Give me the story!
A couple of days before going for my [mandatory military service]() in the Greek Army Forces (and as such, will be under the radar for a few months), I decided to submit [tweetstream.space](http://tweetstream.space) to HackerNews.

I have [some prior experience]() about putting my research, my work, to be reviewed by my peers, but well, that was Physics, something I felt more comfortable with.  
A little voice whispered  "What are you going to show these Silicon Valley hotshots that might be worth their time?", or "Aren't you afraid they're going to tell you it's shit?", but come on, it's not like I couldn't take some criticism, and improve from it!

It was kind of a exciting experience showing one of my products-slash-side projects to the world, getting some exposure and feedback. Sure, it's not something to brag about. A simple web application that indexes tweets will not be 'the talk of the town', like the guy who wrote a [hobby OS](https://news.ycombinator.com/item?id=13794879), but getting nearly 100 users per hour, especially the kind of people that can tell you what's wrong and what's right in your product is very different from showing it to just friends and family, and testing it yourself.

Well, to think of it, I should try to design an OS from scratch, sometime

### Anything else about the app?
The project was created using Flask, and deployed via Heroku. Once I get around tidying the backend, I will open-source the project on my GitHub. It certainly needs some polishing, and has been in my backlog for some months, I must get to it *pronto*!  
Credit goes to [Jefferson Henrique](https://github.com/Jefferson-Henrique) for the initial idea on how to fetch the tweets, as it bypasses the official Twitter API limitation, that provides only two-weeks-old tweets.   
There are no ads, no tracking, no records, no google metrics. I don't think I will add them, either.

The main takeout? I learned a couple of new things about input-testing and stress-testing, that can transfer to and benefit not only this but other projects as well.

### Here's a couple of pics!
<img src="/images/hnpost.png" style='height: 100%; width: 100%; object-fit: contain'/>
