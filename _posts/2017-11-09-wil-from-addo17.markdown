---
layout: post
title:  What I Learned from All Day Devops 2017 - Part 1
date:   2017-11-09
author: Paschalis Ts
tags:   [wil, sysadmin]
mathjax: false
description: "Notes, thoughts, and rambings from watching way too many hours of talks."  
---




ADDO 2017 is over, and it's already been over a week! This past weekend I had some time to catch up with a few more sessions, and decided to try and complete this little write-up. It contains some notes, ideas, and scattered thoughts, that I tried to somewhat structure. Hope you enjoy!

### What's All Day DevOps?
All Day DevOps describes itself as 'the world's largest DevOps conference'. 

Over the course of 24 hours, 100 45-minute sessions, plus four keynote talks were presented, simultaneously across five tracks. These 'tracks' or 'categories' were

* CI/CD
* Cultural Transformation
* DevSecOps
* Dev Ops in Government
* Modern Infrastructure and Monitoring

One of the refreshing things was their 'no vendor pitch' policy. I felt that most of the speakers were actual DevOps practitioners, and not some guy/gal from the Marketing Department that showed up with a scripted monologue to sell software licenses. I would prefer a sprinkle of more in-depth technical talks, but I'm pretty sure it was a great decision.
Also, most of the presenters were kind enough to hang out in the Slack channels, chat with the community and try to answer questions afterwards.  
But enough with sounding like a promoter; this post is to share What I Learned!

I tried to keep up with a couple of tracks when they were live, and then VOD whatever else sounded interesting, so I focused on *DevSecOps*, *Modern Infrastructure* and *CI/CD*.

### Hey I want to start some of that modern/buzz/hype stuff in my team too!
"Building Technical and Organizational Confidence Through Automated Deployments" was a great introductory talk by Mieke Deenen. Looking past all the hype and buzz, Mieke talked about implementing new ideas and processes in a traditional organization.

Their takeaways are

* Start with one small application with low impact.
* Start with a Project Initiation Document (PID), and justify your business case (Benefits vs Costs).
* Get support from upper management.
* Negotiate new rules, learn them and bend them.
* Automate legacy systems, to make time for your new ideas.
* Find balance between old and new way, building 'bridges'. Take small steps, use campfire culture.

### Modern Security Operations

> The most important phase of development is planning  

When implementing a 'DevOps lifecycle', Dev/Sec/Ops teams need to be on the same page. 
DevOps is not just bringing Dev and Ops guys in the same room. It's a philosophy that uses modern tools such as effective versioning/source control, CI/CD, proactive monitoring, testing, red-teams, automation, and allows teams to iterate, innovate and proceed in a faster pace.

> It's 30 times cheaper to fix security defects in dev vs prod.

In a "ABN AMRO Transforms with CI/CD", I liked their five simple core principles, that can be applied to most teams.  
Automate *all* repetitive tasks. Integrate quickly and often. Everyone is equally responsible. Keep changes small. Get continuous feedback.

The value of DevOps comes from implementing a KISS approach, using automation to minimize human errata, as well as tie the security state of the whole stack environment to the application build status.

"Tyro Payments: Securing Australia's Newest Bank", shared success (and failure) stories, in an IT-centered banking company. The takeaway is that security is an iterative process, an all-around effort, that covers infrastructure, code architecture, continuous validation, active monitoring.

A good overview of the modern software-building industry was given by Derek Weeks in "The Data Behind DevSecOps".
 Building software is like building a car, in the sense that you're relying on components built, released and updated by others. This 'production chain' is chaotic, in the sense that for very large, complex projects, it's difficult to 'follow' the security-side of their interrelationship. 
> In 2016, the defect download ratio (known vulnerability) for Java components was 1-in-18  
> Over 30% of official images in Docker Hub contain high priority security vulnerabilities

So, use mature tools, audit your own house, avoid dependency hell and technical debt.

Finally, I'd also recommend "Security In The Land of Microservices" as a good talk about security and management of interrelated components.

In a more technical note:

* During the final hour I was happy to see a `Vault` demo by Fabian Lim, live-streaming [a simple hands-on implementation]( https://3jmaster.github.io/addo2017/). Also, check `Goldfish` out.
* Microservices-slash-container based infrastructure has gone past the 'buzzword hype', and is no longer seen as a panacea. Their merits *and* pitfalls, as well as their place in production are better understood. Huseyin Babal had a *great* talk about best practices in Microservice Architectures, in a wide range of examples.
* `InSpec` is a neat tool to automate and 'turn infrastructure testing, compliance and security requirements into code', by the guys behind Chef.
* "Testing Docker Images Security" showed best practices for Docker security.
Keywords to help with either static analysis or secure implementations of containers
`clair`, `canchore`, `dagda`, `lynis`, `TwistLock`, `Aqua`.
* `Escrow` is a tool by Under Armour that will be open-sourced in the next few months, similar to Vault's functionality.
* `Foreman`, `OpenSCAP`
* One of the more interesting ideas that floated around in the chat was to `Predefine correct process flow, to meter logs against for compliance`.
* [OWASP Top10]() is a good place to start.
* RASP (Runtime application self-protection).
* Watch out for new tools that are maturing, such as the evolution of `SAST` vs `IAST` vs `DAST` [Static, Interactive, Dynamic] Application security testing.

### Modern Infrastructure

A result of DevOps getting Developers and Operations/SysAdmins guys together is many ideas 'crossing fields'. Similar to the Dev side of the equation integrating new practices regarding security and architecture, the Operation/SysAdmin side has been keeping up.

Immutable Infrastructure can be compared to a programmer using a more functional approach to minimize side effects.
One (paraphrased) quote that stayed with me, was
> If you need to SSH to your server you're doing something wrong. It's okay for logging, debugging, developing, but on prod, SSH might as well be disabled. [...] SSH just means you're getting hacked by someone. [...] SSH is a red-alert at Clever Cloud

Getting to such a level of abstraction with regard to your infrastructure requires tremendous planning, a deep understanding of what you're trying to accomplish.

Immutable Infrastructure means that once a server is deployed, it's never modified, merely killed and replaced with a new instance. This makes your 'server herd' easier to manage, much more secure, plus you don't need to maintain state anymore. When there's an updated version of the server, trash the old version, let your automated factory deploy a new one. Manually created, mutable infrastructure "leads to fragile, hard-to-reproduce snowflakes".
> 80% of Outages are Caused by Changes 

And while decoupling runtime from the services and infrastructure running in them requires adding a whole level of abstraction, it seems necessary in complex systems. 
Modern production systems are *much more* than a typical three-tier Host-App-Database stack. Traceability of issues can become a nightmare, and such systems become chaotic. It's not possible for *anyone* at a corporation the size of AWS or even Snap to have a deterministic understanding of their stack, in the way we understand how an Apache server works. When working in such environments, thinking in terms of "Good vs Bad" is erroneous, simply because it's not possible to have all the variables. 
In many cases "There is No Root Cause" for stuff breaking down.

> Emergence refers to the existence or formation of collective behaviors - what parts of a system do together that they would not do alone.
Properties and behaviors of systems arise from both the fine structures that compose those systems and the interrelationship between the systems' discrete parts.

Instead of thinking in term of {Actions, Causal, Prevent, Problems, Failure, Factors}, we can start thinking with {Behavior, Interaction, Properties, Pattern, Dependencies}.

Especially when considering the high rate of change in a startup environment
> High Complexity + Dramatic Change Vectors = Emergent Behavior

[Cynefin]() is a way of quantifying your understanding of complex systems. Knowledge and practice moves patterns towards more favorable 'quadrants'; where complacency moves you backwards.

Cue "Infrastructure as Code". Not just in the sense of automating specific provisioning tasks, but actually dealing with your whole infrastructure as code, using Versioning, CI/CD, proper pipelines. Each time, you start from the same validated state, and the nodes are built once and deployed into a production-like environment everywhere. 
Infrastructure as Code has to do with  
Resource Provisioning --> Configuration Management --> Monitoring and Performance --> Governance and Compliance --> Resource Optimization.  
It also simplifies auditability, and a way of keeping track of your 'herd'. Core infrastructure may change in a slower pace than your network, the application code, but to deal with scaling issues, this automation is necessary.

`CloudFormation` by AWS and `Terraform` by hashicorp are two of the tools that are being used.

People are still worried about working on the cloud instead of using their own bare metal, and for good reason. Cloud APIs are hard, and one solution is for security policies to be embedded as code. That allows peer reviewing, NIST controls in code, deployments validation in design-time, not posthumously, dry-running to help with code reviews. 

Garbage Collection is something I wouldn't think as a priority, but sounds like a very real problem. We think of Garbage Collection as something that happens automatically, but it is actually *not* free, and is not always optimal by default. Tooling includes `GCeasy.io`, `IBM Pattern Modelling for Java`, `HP JMeter`, `Google Garbage Cat`, `Eclipse MAT`, `Heap Hero` .

Some other bulletpoints;

* Safe Rolling of infrastructure and Canary Release for testing.
* Cross Stack References vs Nested Stacks.
* Model-based testing.
* `Consul` is something to watch out for.
* Keywords, keywords, keywords... `Prometheus`, `Sleuth`, `NewRelic`, `AppDynamic`, `DynaTrace`, `Zipkin`, `Spring Boot`, `gelf`.
* API Gateways, `Kong`, `Tyk`.
* Event sourcing, saving the 'state-changing' events of an entity in a time-series format.
* `Sozu`.
* `CloudTrail`.
* `ChaoSlingr`.

Well, the important part is that "Work in Progress: In DevOps, 'Done' is a Quantum State." Just keep moving forward :D :D

### That's all for Part 1. 
PS. Keeping secrets is still a pain in the butt. 
