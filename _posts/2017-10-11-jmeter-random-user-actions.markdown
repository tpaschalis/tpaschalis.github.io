---
layout: post
title:  Simulate Random User Actions in JMeter.
date:   2017-10-11
author: Paschalis Ts
tags:   [jmeter, sysadmin]
mathjax: false
description: "I wanted to randomize user actions for a JMeter test."  
---

I have been assigned a project of designing and running some stress tests using JMeter in a team.  
We would like to simulate an unusually large amount of users logging in and performing some actions in a web-based platform, similar to what happens a handful of times over the course of a year. This will allow us to monitor the platform in our own time, locate any bottlenecks and limitations, and try to fix them.  

After getting to know the basics of what our test would look like and how it could be implemented, we had to make it more realistic, by randomizing certain components in our scenarios. Users navigate the page in different and mysterious ways, some spam their clicks, some go away to get a coffee and return after a while.. You get the point.  
Well, here are some of the tools JMeter has available for this purpose! In my experience, working with CSV files to programmatically control this random behavior was very useful.

### Randomizing User Actions

There are many different ways to randomize user activity for JMeter.  
The `Aggregate Report` Listener is useful for most of these setups. The end result will have many of these tools combined, along with.  


1) Create different different `Thread Groups` with different `Number of Threads`  
Each one will contain a separate controller/sampler with the data that needs to be executed. We need to make sure to untick the "Run...consecutively" option, so they run them simultaneously.

<img src="/images/jmeter-random/Random1.png" style='height: 30%; width: 30%; object-fit: contain'/>

2) Use `Throughput Controllers`, with the `Percent Executions` attribute enabled.  
This allows us to use a single Thread Group, but needs separate controller/samplers for each case.

<img src="/images/jmeter-random/Random2.png" style='height: 30%; width: 30%; object-fit: contain'/>

3) Use `Switch Controllers`.  
We can provide a dynamic Switch Value, from some function/variable/website resposnse, to trigger subordinate samplers. This way, our result can be a `Random Variable`, or use a specific probability (eg choose randomly from a 011223334 string)

<img src="/images/jmeter-random/Random3.png" style='height: 30%; width: 30%; object-fit: contain'/>

4) Use a `Random Controller`.  
Instead of going in order through its sub-controllers and samplers, it picks one at random at each pass.

<img src="/images/jmeter-random/Random4.png" style='height: 30%; width: 30%; object-fit: contain'/>

5) Use a `Random Order Controller`.  
It *will* execute each child element at most once, but the order of execution of the nodes will be random.

<img src="/images/jmeter-random/Random5.png" style='height: 30%; width: 30%; object-fit: contain'/>

6) Use `CSV Data Set Config` to parameterize http requests.

<img src="/images/jmeter-random/Random6.png" style='height: 30%; width: 30%; object-fit: contain'/>

7) Use an `Interleave Controller` to alternate among each of the sub-controllers for each loop iteration.  
So, in our example, the requests will run as 1-4, 2-4, 3-4, 1-4 .... Remember to `Interleave across threads`.

<img src="/images/jmeter-random/Random7.png" style='height: 30%; width: 30%; object-fit: contain'/>


Finally, to properly simulate users, we need to add pauses between each action.
This can happen in two ways, using a constant value pause, or, in a more realistic scenario, use a variable/randomized pause duration.

<img src="/images/jmeter-random/TimePause.png" style='height: 30%; width: 30%; object-fit: contain'/>




