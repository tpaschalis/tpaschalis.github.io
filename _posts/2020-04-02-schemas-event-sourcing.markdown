---
layout: post
title:  DB Schemas and Event Sourcing
date:   2020-04-02
author: Paschalis Ts
tags:   [software, db]
mathjax: false
description: ""
---

Speaking of [kick-ass pairs](https://tpaschalis.github.io/gc-and-concurrency/), there's one more that I don't see mentioned enough; database schemas and event sourcing.

In a standard relational database the schema can be an integral part of the story the data is narrating. Just reading the schema explains relationships between tables, hint towards most used access patterns or one-to-many relationships which may not be apparent at first.

If you haven't heard of event sourcing, [this article](https://martinfowler.com/eaaDev/EventSourcing.html) by Martin Fowler can explain it much better than me.

## The idea 

In microservice architectures, the trend is to move away from a central, monolithic database towards a dedicated per-service database. While there are performance and reliability gains, it also means that data is now scattered around and deduplicated between services. Thus, the focus shifts towards maintaining consistent views of both the deduplicated data, *and* the model under which the data is stored.

Logging every database schema change is the simplest way to 'version control' your storage, enables a more agile mindset. This way, the schema is not just a hindrance, but can assist efficient storage and access patterns of data, changing whenever the requirements change.

In the simplest case all you have to do is record each change in your schema in a timestamped log, along with instructions on how to revert each step, instead of relying on ad-hoc queries.

- Applications will easily be able to find their current schema, and how to reach a desired state.
- Easy replication of production and development environments
- Faster testing of new schema versions, or rollbacks to debug with an earlier one
- Easier consistency between different views of the data in different services


## Implementations
Of course, this is not a novel idea; [Microsoft SQL Server](https://docs.microsoft.com/en-us/sql/relational-databases/track-changes/track-data-changes-sql-server) offers it, as well as [Oracle](https://docs.oracle.com/cd/E24628_01/em.121/e27046/change_management.htm#EMLCM11769), and [Postgres](https://wiki.postgresql.org/wiki/Audit_trigger) but it's not as ubiquitous as I think it should be.

I especially liked waht Laravel offers in [Migrations](https://laravel.com/docs/5.8/migrations). Each *migration* class contains two methods, `up` and `down`, which describe the operation that we wish to apply to the db schema and its reversal.

Migrations are created with a timestamp and stored in a `migrations` table; this way the application can keep track of migrations it has run, and the state it has found itself in, and solve the problem of not knowing where to pick up from. A neat, simple and effective implementation!

See you around!
