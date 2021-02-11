---
layout: post
title:  CLI configuration; Flags or Env? Why not both?
date:   2021-02-11
author: Paschalis Ts
tags:   [golang, code]
mathjax: false
description: ""
---

## Config files vs EnvVars vs Command Line flags
One of Go's most popular uses has to be command-line applications. In the [2019 user survey](https://blog.golang.org/survey2019-results) 62% of developers have used Go to develop CLI apps. It's no wonder, the standard library makes it easy to work with flags and arguments, plus there are lots of opinionated tools like [spf13/cobra](https://github.com/spf13/cobra) or [mitchellh/cli](https://github.com/mitchellh/cli) which provide a structure early on.

One of the decisions I've faced in the past is whether to expose configuration using a `config.json/config.yaml` or `.env` file, environment variables, or command-line flags. 

I've gradually moved away from using files; they can provide structured configuration, but they get easily outdated, usually are a mess to version control and often are just gitignored, so I don't think they're worth it for simpler applications with no more than a dozen arguments, or nested configuration/.

On the other hand, environment variables are universal, handy for local development and easy to use in a Dockerized deployments, or even in Kubernetes if you're using something like Helm. Finally, in my opinion command-line flags should be overriding everything. I have never explicitly used `-foo=bar` and expected something else to happen other than `foo` to be set to `bar`.

## The best of both worlds?
I recently had to build something similar, and was pretty happy with what I came up with, so I'm sharing in case it helps any one of you. The code contains a preset default value, which can be overriden with an env var, which in turn can be overriden by a flag, which is the ultimate decider.

```go
defaultInputFile := "/tmp/access.log"
defaultThreshold := 500

envInputFile, ok := os.LookupEnv("CFG_INPUT_FILE")
if ok {
    defaultInputFile = envInputFile
}
envThreshold, ok := os.LookupEnv("CFG_THRESHOLD")
if ok {
    if thr, err := strconv.Atoi(envThreshold); err != nil {
        defaultThreshold = thr
    } // else, you could fail here
}

inputFile := flag.String("cfg-input-file", defaultInputFile, "Choose the log file to consume.\nDefaults to '/tmp/access.log' or the value of the CFG_INPUT_FILE env var, if it is set")
threshold := flag.Int("cfg-threshold", defaultThreshold, "Choose the alerting threshold.\nDefaults to 500 or the value of the CFG_THRESHOLD env var, if it is set")

flag.Parse()
```

Here's how it looks in action!

```shell
➜ unset $CFG_INPUT_FILE
➜ unset $CFG_THRESHOLD
➜ go run main.go --help
Usage of main:
  -cfg-input-file string
    	Choose the log file to consume.
    	Defaults to '/tmp/access.log' or the value of the CFG_INPUT_FILE env var, if it is set (default "/tmp/access.log")
  -cfg-threshold int
    	Choose the alerting threshold.
    	Defaults to 500 or the value of the CFG_THRESHOLD env var, if it is set (default 500)

➜  export CFG_INPUT_FILE=/dev/urandom
➜  export CFG_THRESHOLD=3
➜  go run main.go --help
Usage of main:
  -input-file string
    	Choose the log file which to consume.
    	Defaults to '/tmp/access.log' or the value of the CFG_INPUT_FILE env var, if it is set (default "/dev/urandom")
  -threshold int
    	Choose the alerting threshold.
    	Defaults to 500 or the value of the CFG_THRESHOLD env var, if it is set (default 3)
```

Finally, if you're exposing the code as a package for other developers, I think that the Builder pattern lends itself nicely to building and validating configuration.

```go
cfg, err := tp.NewConfigBuilder().
    WithInputFile(*inputFile).
    WithThreshold(*threshold).
    WithPollingDuration(*polling).
    Build()

if err != nil {
    return nil, err
}
```

## Parting words
If all this sounded interesting, you should check out [Harvester](https://github.com/beatlabs/harvester) the Open-Source configuration library we've built at [Beat](http://thebeat.co/en) (as well as other of [Sotiris Mantziaris' works](https://github.com/mantzas)).

Harvester is a powerful solution which helps to set up and monitor configuration values, dynamically reconfigure your application, all inside your Go code. It's being actively developed and used by *dozens* of our microservices every day in production, so why not try it for yourself!?