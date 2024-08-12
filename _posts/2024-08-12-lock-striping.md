---
layout: post
title:  Lock striping
date:   2024-08-12
author: Paschalis Ts
tags:   [code, go]
mathjax: false
description: "This could have a fancier name. Something with zebras? Or a race car?"
---

Let's discuss *lock striping*, a very simple technique to reduce lock contention that I haven't seen being discussed frequently. The technique is also known by a bunch of different names, let me know what _you_'ve heard it called!

## The problem

A common problem in concurrent code is how to take advantage of the constant lookup time of a hash-table like structure (a Map, HashMap, Dictionary, or whatever your favorite language calls it) in a thread-safe way.

Generally solutions to this fall somewhere in the spectrum between coarse-grained synchronization (eg. having a single lock for the entire map) and fine-grained synchronization (eg. row-based or per-access locking).

Solutions falling on the earlier part of the spectrum benefit from low overhead but end up with sequential access to the data as usage grows and grabbing/releasing locks in a hot path isn't cheap on the CPU; while solutions falling on the latter part theoretically allow for constant-time concurrent access but with  additional overhead that scales with the size of data.

Go itself contains a [sync.Map](https://pkg.go.dev/sync#Map) types, but it does come with very specific 

Generally, it's hard to predict what solution will perform better under concurrent access in your specific language, CPU architecture and system load. So as it's with all performance problems, the same advice applies _"First measure, then improve"_, and this technique allows just that.

## The idea

Lock striping is a middle-ground solution, which attempts to lower lock contention while also keeping memory stable.
It works by sharding a map in N stripes/sub-maps/buckets, with each stripe having its own lock.
Items are assigned to their stripe pseudorandomly by performing a bitmask or a modulo operation on a sticky characteristic of each item, such as an ID or a hash of its contents.

Then, to lookup an item, you first locate which stripe it belongs to, and obtain its lock, and read from that stripe's submap.

The size of the structure is configurable and thus allows the user to be in control of the tradeoff between access time and memory overhead.

## The solution

Let's look at a real-world example; [the `stripeSeries` struct](https://github.com/prometheus/prometheus/blob/5fd66ba8556053545fa1a1525aaaecfefb2c978a/tsdb/head.go#L1849-L1855) that Prometheus uses to store series in-memory. It was first introduced [seven (!) years ago](https://github.com/prometheus/prometheus/commit/c36d574290e378570d778111b99f1b0687168f6c) and works to this day more or less in the same way.

Here's a slimmed-down version of the struct to showcase the functionality

> I'm omitting `seriesLifecycleCallback` since it doesn't help with our demonstration.
> Similarly, the hashes also store a table of found conflicts, which we don't need here.
> Also, `chunks.HeadSeriesRef` is an alias for uint64; I'm using the latter for readability.

```go

// DefaultStripeSize is the default number of entries to allocate in the stripeSeries hash map.
DefaultStripeSize = 1 << 14

// StripeSize sets the number of entries in the hash map, it must be a power of 2.
// A larger StripeSize will allocate more memory up-front, but will increase performance when handling a large number of series.
// A smaller StripeSize reduces the memory allocated, but can decrease performance with large number of series.

type stripeSeries struct {
    size    int
	series  []map[uint64]*memSeries
	hashes  []map[uint64]*memSeries
	locks   []stripeLock
}

type stripeLock struct {
	sync.RWMutex
	// Padding to avoid multiple locks being on the same cache line.
	_ [40]byte
}
```

This struct holds two copies of series so that they can be looked up by _both_ their internal ID (on the `series` field), but also by the hash of their labels (on the `hashes` field). The lock itself contains a micro-optimization in the form of field padding to avoid CPU cache misses.

Initializing a stripe series is as easy as

```go
func newStripeSeries(stripeSize int, seriesCallback SeriesLifecycleCallback) *stripeSeries {
	s := &stripeSeries{
		size:                    stripeSize,
		series:                  make([]map[uint64]*memSeries, stripeSize),
		hashes:                  make([]map[uint64]*memSeries, stripeSize),
		locks:                   make([]stripeLock, stripeSize),
	}

	for i := range s.series {
		s.series[i] = map[uint64]*memSeries{}
	}
	for i := range s.hashes {
		s.hashes[i] = map[uint64]*memSeries{},
		}
	}
	return s
}
```

The code used to [compare-and-set](https://en.wikipedia.org/wiki/Compare-and-swap) values on this stripeSeries struct is pretty simple.
First, it selects the 'stripe' a series would be assigned to based on the hash of its labels and looks for an existing record. If it's already present, it returns the stored series directly.

Otherwise, it adds two new records for the series based on the hash of its labels and the series ID.

```go

func (s *stripeSeries) getOrSet(hash uint64, lset labels.Labels, createSeries func() *memSeries) (*memSeries, bool, error) {
	series := createSeries()

	i := hash & uint64(s.size-1)
	s.locks[i].Lock()

	if prev := s.hashes[i].get(hash, lset); prev != nil {
		s.locks[i].Unlock()
		return prev, false, nil
	}

	s.hashes[i].set(hash, series)
	s.locks[i].Unlock()

	// Setting the series in the s.hashes marks the creation of series
	// as any further calls to this methods would return that series.

	i = uint64(series.ref) & uint64(s.size-1)

	s.locks[i].Lock()
	s.series[i][series.ref] = series
	s.locks[i].Unlock()

	return series, true, nil
}
```

Then, to look up values by either their IDs or hashes, it's the reverse operation; get the stripe it's located at by bitmasking with the stripeSeries size, and reading simply reading from that sub-map.

```go
func (s *stripeSeries) getByID(id uint64) *memSeries {
	i := uint64(id) & uint64(s.size-1)

	s.locks[i].RLock()
	series := s.series[i][id]
	s.locks[i].RUnlock()

	return series
}

func (s *stripeSeries) getByHash(hash uint64, lset labels.Labels) *memSeries {
	i := hash & uint64(s.size-1)

	s.locks[i].RLock()
	series := s.hashes[i].get(hash, lset)
	s.locks[i].RUnlock()

	return series
}
```

### Outro

And that's about it! Let me know where you've seen this technique used in the wild, what other names you have for it, or other techniques you've used to provide concurrent, thread-safe access to maps!

Until next time, bye!
