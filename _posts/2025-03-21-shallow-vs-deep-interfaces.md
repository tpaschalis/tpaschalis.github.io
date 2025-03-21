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
abstraction, separated into the user-facing interface and the underlying
implementation).

The author argues that "the best modules are those that provide powerful
functionality yet have simple interfaces". The argument is _not_ about absolute
size, but rather the ratio of utility afforded by the abstraction compared to
the size of the abstraction itself.

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
// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// ...
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.

type Reader interface {
	Read(p []byte) (n int, err error)
}
```

It couldn't possibly get _any_ smaller than that, right? It's simple enough
that you won't ever need to look it up again. Searching the Go standard library, one will find
<a href="https://cs.opensource.google/search?q=Read%5C(%5Cw%2B%5Cs%5C%5B%5C%5Dbyte%5C)&ss=go%2Fgo">numerous implementations</a>
including reading from files, from network connections, compressors, ciphers
and more.

This abstraction is both easy to understand and use; the docstring tells you
everything you, as a user, need to know. The underlying implementation can be
buffered, might allow reading from streams or remote locations like an S3
bucket. But crucially, consumers of this API don't need to worry about _how_
reading happens — implementation can be deep and non-trivial,
but a user doesn't have to care. Furthermore, it allows for very little
ambiguity when reasoning about what the code does.

## A shallow interface

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

While the functionality provided by Redis is much larger than just 'reading',
each of these methods has a much shallower implementation; they do exactly one
thing, and they're small enough you could possibly replicate them just by their
name and arguments. The ratio of the functionality provided to the size of
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
user, as they have to understand the nuances between individual methods.
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

So, is this another post criticizing other dev practices? Not really. As
always, things exist on a spectrum.

As a developer, it can often feel more natural to write shallower interfaces.
Similar 'shallow' examples (that are not strictly interfaces) are the
[aws-sdk-go's session Options](https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/session#Options) or
[Viper's](https://github.com/spf13/viper/blob/d319333b0ffd91a9681feaf2d202ce9332df8ecc/viper.go)
public API. Why?

* It makes methods smaller and easier to test
* It maps more closely to the mental map of the system itself
* It takes less time to think up-front about how the user will consume it
* Usually, it only affords a single implementation, maybe two, so it's easier
  to imagine how it will be used.

In contrast, `io.Reader` offers additional advantages:

* Can be easily retrofitted to other use cases
* Requires no state checks to use properly
* Interfaces like io.Reader tend to remain stable over time, while a shallower version would often grow to accomodate more and more features
* It allows for natural composability into other abstractions, like a `ReadWriter` or a `ReadCloser`

```
type ReadCloser interface {
	Reader
	Closer
}
```

So next time you're designing or reviewing an abstraction, pay some closer
attention. [How "deep" is your API](https://www.youtube.com/watch?v=XpqqjU7u5Yc)?
In what ways could you mold it into something simpler that hides complexity
from the user and reduces cognitive load?

For example, does that Redis client API _need_ five different methods for
saving and shutting down? Does a user of the client need to deal with both
running commands and getting meta-information around the DB connection and
runtime metrics at the same time? Are each of the datatypes different enough to
have their own interface? And do I as a reviewer, need to know beforehand
whether the code needs to `Ping`, `Echo` or `Hello`?

## Outro

And that's all for today! If you have any comments, remarks or ideas, feel free
to reach out to me on [Bluesky](https://bsky.app/profile/tpaschalis.me).

What are your favorite interfaces? Any specific one that you think touches the
Platonic ideal? Any that disgusts you beyond imagination and makes you wanna
quit and go grow tomatoes in the countryside? Let me know!

Until next time, bye!

<br>
