---
layout: post
title:  Don't us.
date:   2020-07-18
author: Paschalis Ts
tags:   [meta, k8s]
mathjax: false
description: "Maybe worth less than 2 cents"
---

For the past ~9 months, I've been working with Kubernetes on a weekly basis, either for work, or personal pet projects.

I haven't dived as deep as I'd like, but I take every K8S-related task as an opportunity to learn, *read a lot*, and try to really understand how this complex beast works.

Here's a list of my (unsolicited) opinions on Kubernetes as of July 2020. Hopefully, I'll return back later in the year, and see how these opinions have changed.

Feel free to tell me why I'm wrong, by sending an email or reaching out on Twitter!

- K8S adds new failure modes. A lot of them.
- I'm not convinced a team of developers can reasonably run large-scale K8S without a dedicated DevOps team supporting to solve that 10% ugly edge cases.
- K8s networking can be a nightmare to debug. DNS = :( 
- K8S Jobs/Cronjobs are useful.
- K8S is often used for Resume-Driven Development.
- Part of the learning curve is understanding that K8S introduces a new way of thinking; one that looks like eventual consistency but for software deployments.
- You don't have to understand *all of it* to start reaping its benefits.
- K8S *solves real problems in deploying large-scale architectures*.
- K8S is *overused* by small teams and startups. That 8-person team can maybe spend more time building a product they can sell than fiddling with YAMLs.
- You just can't use K8S without a strong infrastructure; it calls for good CI/CD, logging, tracing and metrics.
- K8S namespacing is preeetty cool. 
- K8S + Helm = <3 
- Random typos like LoadBalancerIp instead of LoadBalancerIP haunt me.
- YAMLS aren't as bad as people make them.

