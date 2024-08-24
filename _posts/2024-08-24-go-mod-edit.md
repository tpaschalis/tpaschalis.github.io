---
layout: post
title:  Let's stop editing go.mod manually
date:   2024-08-24
author: Paschalis Ts
tags:   [go, tooling]
mathjax: false
description: ""
---

I've spent most of my time at Grafana Labs working on [Grafana Agent](https://github.com/grafana/agent), the
project that eventually evolved to [Grafana Alloy](https://github.com/grafana/alloy), our distribution of the
OpenTelemetry Collector.

The nature of the project, building a swiss-knife for the modern observability
practitioner to meet them whenever their data is, meant juggling _a lot_ of
dependencies, on many different projects:

* Prometheus
* Loki
* OpenTelemetry
* Pyroscope
* more than two dozen external Prometheus exporters, and >30 service discovery mechanisms
* libraries for interacting with different cloud providers, the Kubernetes API
* our own libraries like [dskit](https://github.com/grafana/dskit) or [ckit](https://github.com/grafana/ckit)
* aaand a few more stuff.

Given our efforts of giving back to the Open Source community, we're
upstreaming most, if not all, of our changes to the projects that make Grafana
Alloy possible. Sometimes, to prove the feasibility of a new feature, or to
allow us to keep up with our regular release cadence (every 6 weeks), we've had
to temporarily use forks in replace statements. I've personally pushed tags
accidentally in our production CI, while we've had to guard against
incompatible dependencies.

Needless to say, all of these means I'm proudshamed to say that we've had one
of the [_most_ complex `go.mod` files out there](https://github.com/grafana/agent/blob/b9c3594b84e7be24fa936548782130e39636eb38/go.mod); 800 lines long and using most (if not all) of the keywords available.

So how do you go about editing go.mod to upgrade a dependency, or remove a
replace statement? Just pop it into an editor change your stuff, run `go mod tidy`
a bunch of times and hope everything works right? Well, that's what _I've_ been
doing this whole time, but I've come to realize, why not use `go mod edit`
directly instead?

`go mod edit` allows for reproducible changes (show your co-workers how you
came to the result), guards against accidentally breaking the go.mod file and
can be more easily integrated into other workflows, like an automated script.

The help (`go help mod edit`) command and [reference page](https://go.dev/ref/mod#go-mod-edit)
do a great job at explaining how it works, but here's a quick primer:

```
# Reformat the go.mod file without making other changes.
$ go mod edit -fmt

# Print out the go.mod file in text/json format instead of writing changes back
$ go mod edit [-print|-json]

# Point to the non-root go.mod file
$ go mod edit [flags] path/to/go.mod

# The command takes a list of _editing flags_, which are applied in order.
# Some flags edit the overall module name, go version or toolchain that's used
#
$ go mod edit [-module=modname|-go=version|-toolchain=name]

# While others are for managing dependencies
# These flags can be repeated multiple times.
$ go mod edit [-require=path@version|-droprequire=path]
$ go mod edit [-exclude=path@version|-dropexclude=path@version]
$ go mod edit [-replace=old[@v]=new[@v]|-dropreplace=old[@v]]
$ go mod edit [-retract=version|-dropretract=version]
```

#### Special mention -- pseudo versions!

Another pet peeve of mine while dealing with go.mod files is pseudo-version
handling.

Let's say I want to want to temporarily test out a dependency that doesn't
exist in a released version. Here's my usual workflow:

```
# Let's try to add the commit directly
  reading github.com/tpaschalis/prometheus/go.mod at revision v0.0.0-af5a7d1078cee78856d38a879d9ce33a6fdc10b7:
  unknown revision v0.0.0-af5a7d1078cee78856d38a879d9ce33a6fdc10b7

# Whoops, that was wrong, I need a timestamp. Let's add a random one
# Attempt #2; whoops looks like I need I need to trim down the SHA length
  pseudo-version "v0.0.0-20211119180816-af5a7d1078cee78856d38a879d9ce33a6fdc10b7"
  invalid: revision is longer than canonical (expected af5a7d1078ce)

# Ok, now I need to copy the correct timestamp
  pseudo-version "v0.0.0-20211119180816-af5a7d1078ce"
  invalid: does not match version-control timestamp (expected 20200316180026)

# Great, I can _finally_ try to test my custom dependency
  go: downloading github.com/tpaschalis/prometheus v0.0.0-20200316180026-af5a7d1078ce
```

But now, things can be much simpler:

```
$ go mod edit -replace=github.com/prometheus/prometheus=github.com/tpaschalis/prometheus@1b86d54c7facda4a8d3a4df8e143283a9b498492
$ go mod tidy
```

### Outro

And that's all for today! If you know of any interesting tricks or use cases for go mod edit, hit me up on [Twitter](https://twitter.com/tpaschalis_) or [Mastodon](https://m.tpaschalis.me/@tpaschalis)!

Until next time, bye!
