---
layout: post
title:  Error wrapping and errors.Join
date:   2024-09-22
author: Paschalis Ts
tags:   [golang]
mathjax: false
description: ""
---

Errors are the bread-and-butter of Go's flow handling. After working with them
for a few years, I _really_ cannot see why anyone would prefer using exceptions
(at least, when used properly).

The blue book [(The Go Programming Language)](https://www.gopl.io/) offers a
nice piece of advice for constructing useful errors:

> The fmt.Errorf function formats an error message using fmt.Sprintf and returns a new
> error value. We use it to build descriptive errors by successively prefixing additional context
> information to the original error message.

> When the error is ultimately handled by the program’s main function, it
> _should provide a clear causal chain from the root problem to the overall failure_,
> reminiscent of a NASA accident investigation:

```
genesis: crashed: no parachute: G-switch failed: bad relay orientation
```

## Usual error wrapping

In practice, here's how that could look using Go's `%w` verb for _wrapping_ errors.
Error wrapping allows you to not only concatenate the error _message_, but to
include the actual underlying errors, each with potentially its custom type,
fields and methods.

```golang
var ErrRelayOrientation = errors.New("bad relay orientation")

err := fmt.Errorf("G-switch failed: %w", ErrRelayOrientation)
err = fmt.Errorf("no parachute: %w", err)
err = fmt.Errorf("crashed: %w", err)
err = fmt.Errorf("genesis: %w", err)

// 2009/11/10 23:00:00 genesis: crashed: no parachute: G-switch failed: bad relay orientation
log.Fatal(err)
```

Then, the errors.Unwrap function allows you to 'peel back' the successive
errors to create a tree, and the errors.Is and errors.As functions can help
examine that tree, as well as the sub-tree of each of the children nodes in a
pre-order, depth-first traversal.

As an example, here's the (straightforward) unwrapping of the previous error;
we do not have any branching, so it's more of an error _chain_ than a tree.

```golang
// crashed: no parachute: G-switch failed: bad relay orientation
err = errors.Unwrap(err)

//no parachute: G-switch failed: bad relay orientation
err = errors.Unwrap(err)

// G-switch failed: bad relay orientation
err = errors.Unwrap(err)

// bad relay orientation
err = errors.Unwrap(err)

// <nil>
err = errors.Unwrap(err)
```


## Error joining

But what happens if the Genesis crash didn't have _a single_ root cause? How
could we track _multiple_ things that (potentially) went wrong, from different
subsystems?

There's two ways to _join_ errors. They have _slightly_ different semantics
under the hood (look at the Appendix section if you care), but they both work
in a similar way.

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

The second way uses the `errors.Join` function introduced in Go 1.20.
The function returns takes in a variadic argument of errors, discards any nil
values, and wraps the rest of the provided errors.
The message is formatted by joining the strings obtained by calling each
argument's Error() method, separated by a newline.

```
err2 := fmt.Errorf("G-switch failed: %w", errors.Join(
    ErrRelayOrientation,
    ErrCosmicRayBitflip,
    ErrStickyPlunger,
)

// 2009/11/10 23:00:00 G-switch failed: bad relay orientation
// cosmic ray bitflip
// sticky sensor plunger
log.Fatal(err2)
```

But, how to make use of these?
errors.Is, errors.As


If you wanted to get access to the underlying list of errors, unfortunately there's no straightforward way to go about it right now. Calling `errors.Unwrap()` on a multiple error will simply return nil.

But, however you create a multiple wrapped error,the result error implements a new `Unwrap() []error`
method. There have been [multiple](https://github.com/golang/go/issues/57358)
[discussions](https://github.com/golang/go/issues/53435#issuecomment-1191752789)
on whether to provide a helper for unwrapping multiple errors, but so far, none
have been accepted by the community.

The way to achieve it right now is either via explicit interface casting,
either using errors.As or simply inline. Here's an example (or a [playground
link](https://go.dev/play/p/G5vw-xdqBvm) if that's more to your liking).
```golang
	err := errors.Join(
		errors.New("bad relay orientation"),
		errors.New("cosmic ray bitflip"),
		errors.New("sticky sensor plunger"),
	)

	// You can't Unwrap directly, there's no %w error so this returns nil
	fmt.Println(errors.Unwrap(err))

	// You can use errors.As to make sure that the alternate Unwrap() implementation is in place
	var joinedErrors interface{ Unwrap() []error }
	if errors.As(err, &joinedErrors) {
		for _, e := range joinedErrors.Unwrap() {
			fmt.Println("-", e)
		}
	}

	// Or do it more directly with an inline cast
	if uw, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range uw.Unwrap() {
			fmt.Println("~", e)
		}
	}
```



### Outro

What's a G-switch? I came across [this](https://ntrs.nasa.gov/api/citations/20170002586/downloads/20170002586.pdf) resource.

https://appel.nasa.gov/2011/08/29/ata_4-6_genesis_launches-html/

I came across this as I was following through the wonderful [Crafting Interpeters by Robert Nystrom](https://craftinginterpreters.com/). I wanted to return a list of errors found during the parsing stage of the `lox` language the book builds on, but not stop execution. And then I found that since Go 1.XX, it supports a new errors.Join method!

### Appendix 1

In the example that we saw, the underlying types were as follows for the three errors.

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


