---
layout: post
title:  Deep vs Shallow Go interfaces
date:   2025-03-21
author: Paschalis Ts
tags:   [meta, blog,jekyll]
description: "Not to be confused with Shallow vs Deep copy :P"
---

I recently read [A Philosophy of Software Design](https://www.amazon.com/Philosophy-Software-Design-John-Ousterhout/dp/1732102201) by John Ousterhout (of Tcl/Tk, Raft, Sprite fame).

One of the core concepts explored in this book is the distinction between
"deep" vs "shallow" modules (in the author's terms a module is any kind of
abstraction, separated into the user-facing interface and and the underlying
implementation).

The author argues that "the best modules are those that provide powerful
functionality yet have simple interfaces". The argument is _not_ about absolute
size, but rather the ratio of utility afforded by the abstraction to the size
the abstraction itself.

In our case, the main mechanism for composable abstractions in Go is the
`interface` type, so let's examine the concept under this lens.


<figure>
<center>
	<img src="/images/deep-vs-shallow.png" style='height: 70%; width: 70%; object-fit: contain'/> 
	<figcaption>(Deep versus shallow interfaces)</figcaption>
</center>
</figure>

## A deep interface

To me, maybe the best example of a deep interface is `io.Reader`.

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

It couldn't possibly get _any_ smaller than that, right? It's simple enough
that you won't ever need to look it up again. Searching the Go standard library, one will find
<a href="https://cs.opensource.google/search?q=Read%5C(%5Cw%2B%5Cs%5C%5B%5C%5Dbyte%5C)&ss=go%2Fgo">numerous implementations</a>
including reading from files, from network connections, compressors, ciphers
and more.

This abstraction is both easy to understand and use; the underlying
implementation can be buffered, might allow reading from streams or remote
locations like an S3 bucket. But crucially, consumers of this API don't need to
worry about _how_ reading happens — implementation can be deep and non-trivial,
but a user doesn't have to care. Furthermore, it allows for very little
ambiguity when reasoning about what the code does.

## A 'shallow' interface

On the other hand, an example of a shallow interface I've used recently is from
the [redis-go](https://github.com/redis/go-redis) client.

I've trimmed it down for the purposes of this post, but you can see it 
[here](https://github.com/redis/go-redis/blob/11efd6a01ebf1c3a0f8cc41d0a1d54c5afbae26f/commands.go#L160-L230)
in its entirety. It contains 45 methods _and_ uses 19 other interfaces as
extensions for a total of ~200 methods.

```go
type Cmdable interface {
	Pipeline() Pipeliner
	Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
	TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error)
	TxPipeline() Pipeliner
	Command(ctx context.Context) *CommandsInfoCmd
	CommandList(ctx context.Context, filter *FilterBy) *StringSliceCmd
	CommandGetKeys(ctx context.Context, commands ...interface{}) *StringSliceCmd
	CommandGetKeysAndFlags(ctx context.Context, commands ...interface{}) *KeyFlagsCmd
	Info(ctx context.Context, section ...string) *StringCmd
	LastSave(ctx context.Context) *IntCmd
	Save(ctx context.Context) *StatusCmd
	Shutdown(ctx context.Context) *StatusCmd
	ShutdownSave(ctx context.Context) *StatusCmd
	ShutdownNoSave(ctx context.Context) *StatusCmd
    ...
    ...
	StringCmdable
	StreamCmdable
	TimeseriesCmdable
	JSONCmdable
}
```

While the functionality provided by Redis is much larger than just _reading_,
each of these methods has a much shallower implementation; they do exactly one
thing, and they're small enough you could possibly replicate them just by their
name and arguments. The _ratio_ of the functionality provided to the size of
the abstraction is very different than before.

```go
func (c cmdable) CommandGetKeys(ctx context.Context, commands ...interface{}) *StringSliceCmd {
	args := make([]interface{}, 2+len(commands))
	args[0] = "command"
	args[1] = "getkeys"
	copy(args[2:], commands)
	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}
```

This also shifts the responsibility of doing the right thing towards the
consumer, as they have to understand the nuances between individual methods.
In a code review, this makes it harder to reason about what happens at a
glance.

```go
func (c cmdable) Get(ctx context.Context, key string) *StringCmd
func (c cmdable) MGet(ctx context.Context, keys ...string) *SliceCmd
func (c cmdable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
func (c cmdable) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
func (c cmdable) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd
```

## Comparison

So, is this another post criticizing other dev practices? Not at all! As
always, things exist in a spectrum.

As a developer, it can often feel more natural to write shallower interfaces as
they can more closely map what the system affords. Similar 'shallow' examples are the
[aws-sdk-go's session Options](https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/session#Options) or
[Viper's](https://github.com/spf13/viper/blob/d319333b0ffd91a9681feaf2d202ce9332df8ecc/viper.go)
public API.

* It takes less time to think up-front about about how the user will consume it
* It makes methods smaller and easier to test
* It maps more closely to the mental map of the system itself

Conversely, the Reader implementation has its own benefits.

* It allows for much easier composability into a `ReadWriter` or a `ReadCloser`
* Can be easily retrofitted to other use cases
* Requires no state checks to use properly

## Outro

And that's all for today! If you have any comments, remarks or ideas, feel free
to reach out to me on [Bluesky](https://bsky.app/profile/tpaschalis.me).

What are your favorite interfaces? Any specific one that you think touches the
Platonic ideal? Any that disgusts you beyond imagination and makes you wanna
quit and go grow tomatoes in the countryside? Let me know!

Until next time, bye!

<br>
