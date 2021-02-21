---
layout: post
title:  Play MP3 files using Go - Part 1
date:   2020-10-24
author: Paschalis Ts
tags:   [golang, code]
mathjax: false
description: ""  
---

## Intro
I'm not a major audiophile, or keeping up to date with the latest releases and nowadays consume most of my music on Spotify. 

Nevertheless, I occasionally want to play some MP3 files on Mac OS; by default they launch using iTunes.

I have no use for iTunes, and instead of downloading a more minimal MP3 player, why not reuse some libraries and write a small CLI utility?

## Prerequisites
I'm using two of [Hajime Hoshi's](https://github.com/hajimehoshi)'s projects; the first is [go-mp3](https://github.com/hajimehoshi/go-mp3), which is a port of [Krister Lagerstrom's PDMP3](https://sites.google.com/a/kmlager.com/www/projects) (another MP3 decoder) from C to Go as well as [oto](https://github.com/hajimehoshi/oto) which is a lower level sound library which supports all major OSes.

## The code
(error handling is omitted for brevity)



## Next Steps
How about 

## Outro
*Isn't Open-Source just marvelous?* People showcasing their work, sharing their passion with a wide audience, and enabling building 