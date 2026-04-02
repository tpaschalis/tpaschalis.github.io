---
layout: post
title:  Agent-first APIs using  protobuf reflection
date:   2026-03-31
author: Paschalis Ts
tags:   [code, api]
description: "A house of mirrors"
---

Whichever side you are on the LLM turf wars, one thing is clear: what's old is new again.
Wiring APIs together with unix philosophy tools and small CLIs are what moves
the industry right now (whether it moves it to the _right_ direction is an
entirely different question).

A new bottleneck arises: how to make agents understand and use new APIs.
The truth is, not everyone gives their reference pages the same polish and
attention that [Stripe](https://docs.stripe.com/api) does. Does _your_ team
offer an OpenAPI spec and make sure it's up-to-date?

One thing I'm really bullish on is Protobuf APIs, and especially [Protobuf reflection](https://grpc.io/docs/guides/reflection/)
for building _self-describing APIs_ that agents can browse and understand with
a few network calls.

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

These are the very same endpoints that power our UI. We didn't have to build
anything new.

With a couple sentences as input, Claude was able to dive and understand that
API, what its inputs and outputs are, how services have evolved and so on.
It can use the same functionality without any browser automation.
The agent can focus on what it needs, and avoid burning context from reading
lengthy documentation about other features.
Furthermore, the truth is objective; specs drift, documentation gets outdated,
but what is defined in Protobuf files is what the server receives and sends
back. There's no need for trial-and-error.

As we're combining these APIs into new ones that operate on higher abstraction
layers, agents also get these new powers almost for free. For example, we now
supports a new pair of APIs for `Discovery` and `Instrumentation`. These allow
onboarding services for telemetry using a slick and shiny UI in the
[Instrumentation Hub](https://grafana.com/blog/instrumentation-hub-a-guided-scalable-way-to-roll-out-your-observability-coverage-without-losing-control/).
I'm proud of that UI jazz, for hand-holding new users but you know what?
Now, anyone can also tell Claude 'look here, run discovery and instrument
X/Y/Z' and call it a day.

### So what?

I think this is an interesting lightweight alternative to MCPs. While MCP can
be more descriptive for complex usage patterns or multi-step actions, for
tighter APIs this approach requires no special infrastructure or MCP server to
maintain, and no need to write new tool definitions since the agent can
discover what's available on demand.
Pair it with the authentication you're likely already using, and you're off to
a good start.

This also works really well with ConnectRPC's approach of Protobuf + HTTP/JSON
compatibility. After agents sketch out the endpoints, they don't need gRPC
tooling, they can construct payloads and use curl for all subsequent calls.

Finally, it's a natural pairing with "always-on" agents. Say you have an agent
that's running for weeks, and then you deploy a new version of the API. Using
this paradigm the effort can automatically adapt without having to be
interrupted of its current workflows.

### Interested? Try it out!

If you want to poke around with an example service,
[here's](https://github.com/grpc/grpc-go/tree/master/examples/features/reflection)
a great starting point.

Let's start the server
```bash
$ git clone https://github.com/grpc/grpc-go
$ cd grpc-go/examples/features/reflection/
$ go run server/main.go
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
`SayHello` RPC.
Let's inspect what that looks like.

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

This should be enough information to call it, let's go ahead and do that.

```bash
$ grpcurl -plaintext -format text -d 'name: "gRPCurl"' \
  localhost:50051 helloworld.Greeter.SayHello

message: "Hello gRPCurl"
```

### Outro

So that's all for today. A small brain-dump on something I'm excited about.

What's your take on agents using APIs? Any other approaches that worked well for you? Hit me up on [Bluesky](https://bsky.app/profile/tpaschalis.me) and let me know!


<br><br><br>

_Special thanks to Apostolis Mpostanis and William Dumont for their feedback on this post_.
