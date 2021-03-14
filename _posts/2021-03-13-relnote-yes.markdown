---
layout: post
title:  How the Go team could track what to include in release notes
date:   2021-03-13
author: Paschalis Ts
tags:   [golang, foss]
mathjax: false
description: ""
---

Release notes can sometimes be exciting to read. 

Condensing the work since the last release into a couple of paragraphs, announcing new exciting features or deprecating older ones, communicating bugfixes or architectural decisions, making important announcements.. Come to think of it, the couple of times that I've had to *write* them, wasn't so bad at all! 

Unfortunately, the current trend is release notes becoming a mix of *Bug fixes*, *Made ur app faster*, *Added new feature, won't tell you what it is*, which can sound like generalities at best and condescending or patronizing at worst; usually like something written just to fill an arbitrary word limit in the last five minutes before a release.

Here's what's currently listed in the "What's New" section for a handful of popular applications in the Google Play Store.

```
- Thanks for choosing Chrome! This release includes stability and performance improvements.

- Every week we polish up the Pinterest app to make it faster and better than ever. 
Tell us if you like this newest version at http://help.pinterest.com/contact

- Get the best experience for enjoying recent hits and timeless classics with our latest 
Netflix update for your phone and tablet.

- We update the Uber app as often as possible to help make it faster and more reliable 
for you. This version includes several bug fixes and performance improvements.

- We’re always making changes and improvements to Spotify. To make sure you don’t miss 
a thing, just keep your Updates turned on.

- For new features, look for in-product education & notifications sharing the feature 
and how to use it! (FYI this was YouTube, as it doesn't even mention the product's name)
```

The Opera browser, on the other hand has something more reminiscent of actual release notes.
```
What's New
Thanks for choosing Opera! This version includes improvements to Flow, 
the share dialog and the built-in video player.

More changes:
- Chromium 87
- Updated onboarding
- Adblocker improvements
- Various fixes and stability improvements
```


Just to make things clear *I'm not bashing these fellow developers at all*. [Here's](https://github.com/beatlabs/patron/releases) the release history of a project I'm helping maintain; our release notes can be just as vague sometimes. 

Writing helpful, informative (and even fun!) release notes is time consuming and has little benefit non-technical folks. It's also hard to keep track of what's changed since the last release, and deciding what's important and what's not.

How would *you* do it?

## The Go team solution
So, how is Go team approaching this problem? A typical Go release in the past three years may contain from 1.6k to 2.3k commits.

```
from -> to   commits
1.15 -> 1.16 1695
1.14 -> 1.15 1651
1.13 -> 1.14 1754
1.12 -> 1.13 1646
1.11 -> 1.12 1682
1.10 -> 1.11 2268
1.9  -> 1.10 1996
1.8  -> 1.9  2157
```

How do you keep track of what was important, what someone reading the release notes may need to know?

I set to find out, after [Emmanuel](https://twitter.com/odeke_et) (a great person, and one of the best ambassadors the Go community could wish for), added a mysterious comment on one of my [latest CLs](https://go-review.googlesource.com/c/go/+/284136) that read `RELNOTE=yes`.

The [`build`](https://github.com/golang/build) repo holds Go's continuous build and release infrastructure; it also contains the [`relnote` tool](https://github.com/golang/build/blob/master/cmd/relnote/relnote.go) that gathers and summarizes Gerrit changes (CLs) which are marked with RELNOTE annotations. The earliest reference of this idea I could find is [this CL](https://go-review.googlesource.com/c/build/+/30697) from Brad Fitzpatrick, back in October 2016.

So, any time a commit is merged (or close to merging) where someone thinks it may be useful to include in the release notes, they can leave a `RELNOTE=yes` or `RELNOTES=yes` comment. All these CLs are then gathered to be reviewed by the release author. Here's the actual Gerrit API query:
```
query := fmt.Sprintf(`status:merged branch:master since:%s (comment:"RELNOTE" OR comment:"RELNOTES")`
```

Of course, this is not a tool that will automatically generate something you can publish, but it's a pretty good alternative to sieving a couple thousands of commits manually.

I love the simplicity; I feel that it embodies the Go way of doing things. I feel that if my team at work tried to find a solution, we'd come up with something much more complex, fragile and unmaintainable than this. The tool doesn't even support time ranges as input; since Go releases are roughly once every six months, here's how it decides which commits to include

```go
// Releases are every 6 months. Walk forward by 6 month increments to next release.
cutoff := time.Date(2016, time.August, 1, 00, 00, 00, 0, time.UTC)
now := time.Now()
for cutoff.Before(now) {
    cutoff = cutoff.AddDate(0, 6, 0)
}

// Previous release was 6 months earlier.
cutoff = cutoff.AddDate(0, -6, 0)
```

## In action!
Here's me running the tool, and a small part of the output.

```bash
$ git clone https://github.com/golang/build
$ cd build/cmd/relnote
$ go build .
$ ./relnote
...
...
  https://golang.org/cl/268020: os: avoid allocation in File.WriteString
reflect
  https://golang.org/cl/266197: reflect: add Method.IsExported and StructField.IsExported methods
  https://golang.org/cl/281233: reflect: add VisibleFields function
syscall
  https://golang.org/cl/295371: syscall: do not overflow key memory in GetQueuedCompletionStatus
unicode
  https://golang.org/cl/280493: unicode: correctly handle negative runes
```


## Parting words
That's all for today! I hope that my change will find its way on the Go 1.17 release notes; if not I'm happy that I learned something new! 

I'm not sure if the `relnote` tool is still being actively used, but I think it would be fun to learn more about what goes into packaging a Go release.

Until next time, bye!
