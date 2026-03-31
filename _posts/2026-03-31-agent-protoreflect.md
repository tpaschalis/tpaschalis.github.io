---
layout: post
title:  Protobuf reflection for agent-enabled, self-describing APIs
date:   2026-03-31
author: Paschalis Ts
tags:   [code, api]
description: "A house of mirrors"
---

Whichever side you are on the LLM turf wars, one thing is clear. APIs are back,
baby! It's funny how fashion comes and goes, but Unix philosophy tools and
wiring APIs together are what moves the industry right now (whether it moves it
to the _right_ direction is an entirely different question).

A new bottleneck arises, how to make agents understand and use new APIs.
The truth is, not everyone gives their API reference the same polish and
attention that [Stripe](https://docs.stripe.com/api) does. Does _your_ team
offer an OpenAPI spec and make sure it's up-to-date? If that's a yes, great
work, but certainly that's _not_ the case for most.

One thing I'm really bullish on is Protobuf APIs, and especially [Protobuf reflection](https://grpc.io/docs/guides/reflection/).

At work, we've been using [ConnectRPC](https://connectrpc.com/) for a couple of
years now and we've been really happy with that choice. But the unexpected
killer feature recently has been how it allows for a self-describing API that
agents can browse and understand with a few network calls.

### A real-world example.

Here's a small example from a fresh Claude Code session.

```
$ claude
❯ Reach out to the API in localhost:8080 using grpcurl and use protobuf reflection.
Give me an overview of the available services and methods I can use to manage pipelines.

⏺ Bash(grpcurl -plaintext localhost:8080 list)
  ⎿  collector.v1.CollectorService
     discovery.v1.DiscoveryService
     instrumentation.v1.InstrumentationService
     … +3 lines (ctrl+o to expand)

⏺ Bash(grpcurl -plaintext localhost:8080 describe pipeline.v1.PipelineService)
  ⎿  pipeline.v1.PipelineService is a service:
     service PipelineService {
       rpc BulkDeletePipelines ( .pipeline.v1.BulkDeletePipelinesRequest ) returns ( .pipeline.v1.BulkDeletePipelinesResponse );
     … +13 lines (ctrl+o to expand)

⏺ Here's an overview of the pipeline.v1.PipelineService available at localhost:8080:

  CRUD Operations

  ┌─────────────────────┬──────────────────────────────────────────────────────────────┐
  │       Method        │                         Description                          │
  ├─────────────────────┼──────────────────────────────────────────────────────────────┤
  │ CreatePipeline      │ Create a single pipeline, returns the created Pipeline       │
  ├─────────────────────┼──────────────────────────────────────────────────────────────┤
  │ GetPipeline         │ Retrieve a pipeline by ID                                    │
  ├─────────────────────┼──────────────────────────────────────────────────────────────┤
  │ GetPipelineID       │ Look up a pipeline's ID (likely by name or other identifier) │
  ├─────────────────────┼──────────────────────────────────────────────────────────────┤
  │ UpdatePipeline      │ Update an existing pipeline                                  │
  ├─────────────────────┼──────────────────────────────────────────────────────────────┤
  │ UpsertPipeline      │ Create or update a pipeline (insert-or-update)               │
  ├─────────────────────┼──────────────────────────────────────────────────────────────┤
  │ DeletePipeline      │ Delete a single pipeline                                     │
  ├─────────────────────┼──────────────────────────────────────────────────────────────┤
  │ BulkDeletePipelines │ Delete multiple pipelines at once                            │
  └─────────────────────┴──────────────────────────────────────────────────────────────┘

[.......]

  You can drill deeper into any method's request/response types with e.g. grpcurl -plaintext localhost:8080 describe
  pipeline.v1.CreatePipelineRequest.
```

These are the very same APIs that power the Fleet Management UI in Grafana
Cloud. We didn't have to build anything new here.

With a couple sentences as input, Claude was able to dive and understand that
API, what its inputs and outputs are, its various versions, how services have
evolved and so on. The agent can focus on the specific services it needs, and
avoid burning context from reading lengthy, inconclusive or outdated
documentation about features it doesn't necessarily need. Furthermore, the
truth is objective, what is defined in Protobuf files; there's no room for
subjective or misleading documentation, falling out of date or trying random
things until the agent gets a 200 - OK.

As we're building more advanced APIs that operate on higher abstraction layers,
agents get all these new powers for ~free. For example,  Fleet Management now
supports a new pair of APIs for `Discovery` and `Instrumentation`. These allow
onboarding services for telemetry using a slick and shiny UI in the [Instrumentation Hub](https://grafana.com/blog/instrumentation-hub-a-guided-scalable-way-to-roll-out-your-observability-coverage-without-losing-control/).
I'm proud of that UI jazz, but you know what? You can also tell Claude 'use
this API to discover all services and instrument X/Y/Z'. And to be honest, it
might be a better experience for some use cases.

### Interested? Try it out!

If you want to poke around with an example service,
[here's](https://github.com/grpc/grpc-go/tree/master/examples/features/reflection)
a great starting point.

Let's start an example server
```bash
$ git clone https://github.com/grpc/grpc-go
$ cd grpc-go/examples/features/reflection/
server listening at [::]:50051
```

Then use `grpcurl` to poke around
```bash
$ grpcurl -plaintext localhost:50051 list

grpc.examples.echo.Echo
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
helloworld.Greeter

$ grpcurl -plaintext localhost:50051 list helloworld.Greeter

helloworld.Greeter.SayHello
```

Hmm, looks like this server exposes four services. The `Greeter` service has a
`SayHello` RPC. Let's inspect how that looks like

```
$ grpcurl -plaintext localhost:50051 describe helloworld.Greeter

helloworld.Greeter is a service:
service Greeter {
  rpc SayHello ( .helloworld.HelloRequest ) returns ( .helloworld.HelloReply );
}

$ grpcurl -plaintext localhost:50051 describe helloworld.HelloRequest

helloworld.HelloRequest is a message:
message HelloRequest {
  string name = 1;
}
```

This should be enough information to call the rpc, let's go ahead and do that.

```bash
$ grpcurl -plaintext -format text -d 'name: "gRPCurl"' \
  localhost:50051 helloworld.Greeter.SayHello

message: "Hello gRPCurl"
```

### Outro

So that's all for today. A small brain-dump on something I'm excited about.

What's your take on agent-first APIs? If your team is working with protos, give
this a spin!
