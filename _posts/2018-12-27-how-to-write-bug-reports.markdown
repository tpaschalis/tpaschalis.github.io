---
layout: post
title:  Writing a correct defect/bug report
date:   2018-12-27
author: Paschalis Ts
tags:   [bugs, software, industry]
mathjax: false
description: "Because a bug never reported is never getting fixed"  
---

I'm very grateful for testers, especially good ones. The kind of people that tirelessly follow test scenarios, find new, creative ways to break things, propose innovative solutions, and catch game-changing bugs before they reach production. I feel that it's not that unusual for testers to understand how systems glue together better than some engineers, actually.

Truth be told, testing and debugging are skills that are gained on the job, but here's a few basic guidelines so that your team can start developing better practices in testing and reporting defects. But first...

### What is a defect?

There are different definitions of a 'defect' in the industry, some more right than others. A *good* definition in my opinion is the following.

> A defect is any deviation between a system's designed and actual functionality.

In the case of software, a defect is any difference between the initial design and specification and its actual behavior, which may be a cause of a bug, a logic error, or a design mishap. The defect itself would not and should not differentiate such cases. 

### How to write a 'good' bug report

The ideal bug report should be [short, clear, concise](http://sscce.org/), and contain *a single defect per ticket*. If your system/team uses a template, *use that* and do not improvise. Its essential components are the following :

* A unique identifier   
This identifier will be used to refer to, investigate, and cooperate on an issue, rather than stating "that login bug Brad mentioned on his email last September".
* Platform/OS/Version      
Describe the way your system is set up.    
What version is the system and its dependencies on? Are you on Windows, MacOS, or Linux? What distribution? Did you follow the 'official' installation steps? Did you just `apt-get` or did you compile from source?
* Context   
A descriptive title and TL;DR is invaluable to the engineering team, who will be able to quickly relate to what you're talking about
* Steps to reproduce    
*It works on my machine*, but backwards, because, well, it doesn't.
* Expected Outcome   
How did you expect the system to behave?
* Actual Outcome, with proof   
What did the system do instead? It's vital that you include screenshots, logs, session history, or any other information that will help the developer investigate the issue.


In some cases, there are some other things you can do as a tester to make your case heard.
* (Optional) Contact information    
Especially in FOSS, it can help the developer communicate back to you, if he requires more data or test cases, or information on your setup  
* (Optional) A [minimal, complete, verifiable example](https://stackoverflow.com/help/mcve)   
* (Optional) Proposed solution     
If you're pretty sure that you have a good hint for the engineering team, go ahead and tell them! It's better that you have a very good understanding of the system and its dynamic nature.
* (Optional) References to similar issues    
* (Optional) Screenshots   


Finally, finding a bug, or tracking down a defect can be very challenging and satisfying. That doesn't mean you can to be rude to the developers, proving how much smarter you are, and that you'd have already a fix present. Again, especially in FOSS, where people volunteer their off days and weekends so that everyone can enjoy their quality work at no cost. If you wish to submit an issue, make sure that you're at least polite, and as helpful as possible to resolve issues.


### Severity vs Priority

It *might* be your responsibility as a tester, to assign a "Severity" or a "Priority" value for the issue you're reporting. Many people think that these two are synonymous terms, but this is *not* the case, even though they are usually correlated.

The first one is associated with standards and is driven by functionality and the system's technical aspect, while the second one is associated with scheduling and project management, and is driven by business value and the customer's requirements.

* ***Severity*** : Is a metric that classifies a defect *based on the impact it has on the system's functionality*.   
In other words, it defines the defect's ability to 'break' the system, disrupt parts of a service it provides and cause issues like financial loss, loss of data, compliance failures or disrupted workflows. A high-severity (aka a "Show Stopper") issue will many times be directly related different parts of a system, will have business impact, or be prioritized by the Technical Lead.


* ***Priority*** : Is a metric that classifies a defect *based on the urgency in resolving it*.   
In an ideal world, our systems should be "bug free", but since that cannot really be the case, we have to find out which issues deserve prior attention, to be quickly fixed and schedule our workflow accordingly. This scheduling order might be related to customer needs, to the project's milestones or it might be the result of a risk/cost assessment. A high-priority (aka "Urgent") issue will often be a high-severity issue, the prerequisite to testing/fixing another part of the system, or the decision of a Project Manager.


To make things more clear, here are some possible examples of defects if you were running an e-shop,

* High Severity - High Priority bug
Customers are able to order items from our e-shop for free! :X    

* High Severity - Low Priority bug   
Our e-shop is broken for customers using IE6.

* Low Severity - High Priority   
The website's dark mode has low contrast and is bad for accessibility.

* Low Severity - Low Priority   
The company logo does not appear on our "About" page.







### How to design or implement a bug-reporting system

If you're responsible of creating a reporting system or updating the team practices on issue reporting, here's a few of things that make an issue tracking team and system great : 

* Have an easy way for users to submit them, it's your responsibility to filter them out later on.   
Users have limited amount of time to dedicate to your app/system/software before giving up, so why would you make it hard for them to help you? You want to make it as easy as possible for them to offer feedback, as it's *way better* to have a large amount of 'data points' to filter from, even duplicates, than having no feedback at all.

* Keep things as simple as you can, but not more   
As an example, there are systems where there was a hard distinction between 'major' and 'critical' bugs. While this might be true, and useful for tracking internal progress, in most cases it makes no difference for the user that reports the issue and just serves to complicate and confuse them.

* Provide incentives   
Just a simple name-drop of the initial reporter in the Release Notes, or a "Thank you" email when you close the issue makes it that much more possible that he'll offer his feedback again. If you can afford a bug bounty program, *do it*!.

* Have technical staff get involved   
Companies often outsource bug fixing to Junior hires(or worse, to non-developers), and don't see the value behind this process. Neglecting bug fixing means technical debt, which *will* catch up on you down the road. Having your technical team actively involved in bug reporting, fixing, and customer feedback is invaluable, since it 'tightens' that [feedback loop](https://en.wikipedia.org/wiki/Feedback), and provides immense business value.

* UNIQUE IDS    
Have unique identification for each defect, which can be used to detect duplicates, and track their status directly in commit messages.

* Keep them in the clear of what will be fixed first and what will have to wait.     
Clear communication will help users understand your prioritization of the defect fixes and keeping them in the dark makes them much more likely to be frustrated.

* Get your unit tests on    
In some cases, having enough test coverage can immediately identify serious defects versus ones that have happened due to user errors or improper following of the test scenarios. 



### That's all!
I hope you have some quality time off the screen these holidays, have the company of your loved ones, and charge up the batteries to keep producing awesomeness.   
See ya!
