---
layout: post
title:  Daffodil, a Snowflake cousin
date:   2020-05-08
author: Paschalis Ts
tags:   [go, software]
mathjax: false
description: ""
---


Last weekend I built a distributed ID generator using Go. I named it [Daffodil](http://github.com/tpaschalis/daffodil), and it's just another [Snowflake](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake.html) cousin, like [IdGen](https://github.com/RobThree/IdGen), [Sonyflake](http://github.com/sony/sonyflake), or many others.

[As I said](https://twitter.com/tpaschalis_/status/1257660129459277825), one would think *hey, 2010 called they want their distributed algorithms back*, but here's what I learned


## Performance 
The [original Snowflake requirements](https://github.com/twitter-archive/snowflake/tree/scala_28), that of a large system 10 years ago were about 10k rps and sub-2ms response time per process.

It seems like today, Go can achieve this on a mid-range laptop *easily*. Under no circumstances do I mean that my silly project can scale to Twitter-like scale, but that modern tooling allows scaling to a respectful level with little effort.

## Motivation
So would anyone use Snowflake et al. in 2020?

For some use-cases, I find it very compelling to have time-ordered IDs which can also instantly provide insights like the following, with just some basic bit-shifting and arithmetic.
- How many IDs were created this past day?
- Did failed IDs come from a single node, or was it a general problem?
- What was the biggest burst of activity, and on which time window?

In the age that you can click-and-deploy solutions like ELK and Prometheus, this may not as strong a point as it was.   
You don't want to give this information to your end-users and competitors? You could just salt your IDs and rotate it every once in a while.


## Alternatives
So what are some of your other options?

- Autoincrements / Ticket Servers  
I've seen auto-increments be greatly (ab)-used, but the reason is simple. They *just work*, and you're probably underestimating how far they can bring you. On the other hand, it takes some work to run in a distributed fashion, and you might write-bottleneck your system just to generate IDs, which is less than ideal.

- UUIDs  
I personally love the idea of UUIDs, both server and client-side. They are standardized, they can mostly just plug-and-play, and can contain information about time and the node they were created in. Most complaints about UUIDs are about a) their bulkiness b) people treat them as string in their databases, degrading performance c) harder for manual debugging.

- Use hashes  
In many cases, hashing some content plus metadata can be used for IDs. They are easy to produce both server and client side, but most grievances are the same as UUIDs.

- Use specific tools  
    - [MongoDB’s ObjectId](https://docs.mongodb.com/manual/reference/method/ObjectId/)
    - [Cassandra Counters](https://docs.datastax.com/en/cql-oss/3.3/cql/cql_using/useCountersConcept.html)
    - [Apache Ignite](https://apacheignite.readme.io/v1.0/docs/atomic-types) using java.util.concurrent.atomic.AtomicLong
    - [Zookeeper's PERSISTENT_SEQUENTIAL](https://zookeeper.apache.org/doc/r3.3.3/api/org/apache/zookeeper/CreateMode.html#PERSISTENT_SEQUENTIAL)

