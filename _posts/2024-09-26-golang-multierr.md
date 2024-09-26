---
layout: post
title:  Joining errors in Go
date:   2024-09-26
author: Paschalis Ts
tags:   [go, errors]
mathjax: false
description: "err1+err2=??"
---

I recently realized that the stdlib `errors` package in Go supports _joining_
multiple errors in addition to the usual direct `%w` wrapping.

I haven't really seen this used much in the wild; I think most people either
refactor to avoid multiple errors, return an []error or use
[uber-go/multierr](https://github.com/uber-go/multierr). Let's go have a look!

> I started drafting a longer version of this post, but it almost blew up to
> become a fully-fledged guide of "Proper use of errors in Go".
> While this bigger post _may_ be coming in the future, I think this smaller
> one is useful enough to share on its own.

### Joining errors

You can join multiple errors in two ways. They have _slightly_ different
semantics under the hood (look at the Appendix section if you care), but they
both work in a similar way.

The first one is by using multiple `%w` verbs on the same error.

```
var (
    ErrRelayOrientation = errors.New("bad relay orientation")
    ErrCosmicRayBitflip = errors.New("cosmic ray bitflip")
    ErrStickyPlunger    = errors.New("sticky sensor plunger")
)

err1 := fmt.Errorf("G-switch failed: %w %w %w", ErrRelayOrientation, ErrCosmicRayBitflip, ErrStickyPlunger)

// 2009/11/10 23:00:00 G-switch failed: bad relay orientation cosmic ray bitflip sticky sensor plunger
log.Fatal(err1)
```

The second one uses the `errors.Join` function introduced in Go 1.20.
The function takes in a variadic error argument, discards any nil values, and
wraps the rest of the provided errors. The message is formatted by joining the
strings obtained by calling each argument's Error() method, separated by a
newline.

```
err2 :=  errors.Join(
    ErrRelayOrientation,
    ErrCosmicRayBitflip,
    ErrStickyPlunger,
)

// 2009/11/10 23:00:00 bad relay orientation
// cosmic ray bitflip
// sticky sensor plunger
log.Fatal(err2)
```

## How to use them?

Both error wrapping variants ultimately form a _tree_ of errors. The ways to
inspect that tree are the `errors.Is` and `errors.As` functions. Both of these
examine the tree in a pre-order, depth-first traversal by successively
unwrapping every node found.

```
func Is(err, target error) bool
func As(err error, target any) bool
```

The `errors.Is` function examines the input error's tree, looking for a leaf
that matches the target argument and reports if it finds a match. In our case,
this can look for a leaf node that matches a specific joined error.

```
ok := errors.Is(err1, ErrStickyPlunger)
fmt.Println(ok) // true
```

On the other hand, `errors.As` examines the input error's tree, looking for a
leaf that can be _assigned to the type_ of the target argument. Think of it as
an analog to json.Unmarshal.

```
var engineErr *EngineError
ok = errors.As(err2, &engineErr)
fmt.Println(ok) // false
``` 

So, to summarize:

* `errors.Is` checks if a specific error is part of the error tree
* `errors.As` checks if the error tree contains an error that can be assigned to a target type

## The catch

So far so good! We can use both types of wrapping, both direct, single error
wrapping as well as joining to form a tree. Now that we've seen that, let's
explore how to inspect the original contents of that tree on another part of
the codebase.

But there's a _slight_ complication here. Let's try to call errors.Unwrap()
directly on any of the two joined errors created above.

```
fmt.Println(errors.Unwrap(err1)) // nil
fmt.Println(errors.Unwrap(err2)) // nil
```

So, why `nil`?! What's going on?  How can I get the original errors slice and
inspect it? Turns out, that the two 'varieties' of wrapping implement a
_different_ Unwrap method.

```
Unwrap() error
Unwrap() []error
```

The documentation of `errors.Unwrap` method clearly states that it only calls
the first one and does not unwrap errors returned by Join. There have been
[multiple](https://github.com/golang/go/issues/53435#issuecomment-1191752789)
[discussions](https://github.com/golang/go/issues/57358) on golang/go about
allowing a more straightforward way to unwrap joined errors, but there has been
no consensus. The way to achieve it right now is to either use `errors.As` or
an inline interface cast to get access to the second Unwrap implementation.


```
var joinedErrors interface{ Unwrap() []error }


// You can use errors.As to make sure that the alternate Unwrap() implementation is available
if errors.As(err1, &joinedErrors) {
	for _, e := range joinedErrors.Unwrap() {
		fmt.Println("-", e)
	}
}

// Or do it more directly with an inline cast
if uw, ok := err2.(interface{ Unwrap() []error }); ok {
	for _, e := range uw.Unwrap() {
		fmt.Println("~", e)
	}
}
```

So, it's an extra little step, but with either of these techniques you'll be
able to retrieve the original slice of errors. My inspiration for this was
following along the [Crafting Interpreters](https://craftinginterpreters.com/introduction.html)
book; when implementing the language's lexer/scanner, I wanted to keep gather
all encountered errors and report them in one go.

### Outro

And that's all for today! If you have any comments, remarks or ideas, feel free
to reach out to me on [X/Twitter](https://twitter.com/tpaschalis_) or
[Mastodon](https://m.tpaschalis.me/@tpaschalis)!

Oh, and you can play around with the code samples in the [Go Playground](https://go.dev/play/p/7qhSZWCthtW).

Until next time, bye!


<br>
<br>
<br>
<br>
<br>
<br>
<br>
<br>
<br>
<br>
<br>
<br>

### Appendix 1

In the example that we saw, the underlying types were as follows for these three errors.

So, while multiple `%w` verbs give out a slice of errors, directly, the errors.Join function wraps them in a `joinError` struct.

```
type joinError struct {
	errs []error
}
```

```
var (
	ErrRelayOrientation = errors.New("bad relay orientation")
	ErrCosmicRayBitflip = errors.New("cosmic ray bitflip")
	ErrStickyPlunger    = errors.New("sticky sensor plunger")
)

err1 := fmt.Errorf("G-switch failed: %w\n%w\n%w", ErrRelayOrientation, ErrCosmicRayBitflip, ErrStickyPlunger)

err2 := fmt.Errorf("G-switch failed: %w", errors.Join(
	ErrRelayOrientation,
	ErrCosmicRayBitflip,
	ErrStickyPlunger,
))

    err3 := errors.Join(
	ErrRelayOrientation,
	ErrCosmicRayBitflip,
	ErrStickyPlunger,
)

// &fmt.wrapErrors{msg:"bad relay orientation\ncosmic ray bitflip\nsticky sensor plunger", errs:[]error{(*errors.errorString)(0xc00009c050), (*errors.errorString)(0xc00009c060), (*errors.errorString)(0xc00009c070)}}
// &fmt.wrapError{msg:"G-switch failed: bad relay orientation\ncosmic ray bitflip\nsticky sensor plunger", err:(*errors.joinError)(0xc0000be000)}
// &errors.joinError{errs:[]error{(*errors.errorString)(0xc0000140a0), (*errors.errorString)(0xc0000140b0), (*errors.errorString)(0xc0000140c0)}}
```

