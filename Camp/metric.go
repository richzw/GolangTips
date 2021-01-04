package main

import (
	"sort"
	"sync"
	"time"
)

const SlidingWindowSize = 60

// Refer to https://github.com/Netflix/Hystrix/blob/master/hystrix-core/src/main/java/com/netflix/hystrix/metric/consumer/BucketedRollingCounterStream.java
// Implement sliding window through bucket
// Buckets: key is timestamp and granularity is second
//          bucket is the time duration of client requests
type WindowMetric struct {
	Buckets map[int64]*Bucket
	Mutex   sync.RWMutex
}

type Bucket struct {
	Durations []time.Duration
}

func NewWindowMetric() *WindowMetric {
	wm := &WindowMetric{
		Buckets: make(map[int64]*Bucket),
		Mutex:   sync.RWMutex{},
	}
	return wm
}

type Durations []time.Duration

func (c Durations) Len() int           { return len(c) }
func (c Durations) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Durations) Less(i, j int) bool { return c[i] < c[j] }

// TODO: cached durations after computed to prevent dup calc
func (wm *WindowMetric) SortDurations() Durations {
	var durations Durations
	now := time.Now()

	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	for ts, b := range wm.Buckets {
		if ts >= now.Unix() - SlidingWindowSize {
			for _, d := range b.Durations {
				durations = append(durations, d)
			}
		}
	}

	sort.Sort(durations)

	return durations
}

func (wm *WindowMetric) getBucketSlow(now time.Time) *Bucket {
	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	wm.Buckets[now.Unix()] = &Bucket{}
	return wm.Buckets[now.Unix()]
}

func (wm *WindowMetric) getBucket() *Bucket {
	wm.Mutex.RLock()
	now := time.Now()
	bucket, exists := wm.Buckets[now.Unix()]
	wm.Mutex.RUnlock()

	// Similar to sync.Once to do double check
	if !exists {
		bucket = wm.getBucketSlow(now)
	}

	return bucket
}

func (wm *WindowMetric) del() {
	now := time.Now()

	for ts := range wm.Buckets {
		if ts <= now.Unix() - SlidingWindowSize {
			delete(wm.Buckets, ts)
		}
	}
}

// Add appends the time.Duration given to the current time bucket.
func (wm *WindowMetric) Add(duration time.Duration) {
	b := wm.getBucket()

	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	b.Durations = append(b.Durations, duration)
	wm.del()
}

// Computes the average timing in the last one minute.
func (wm *WindowMetric) Mean() uint32 {
	sortedDurations := wm.SortDurations()
	length := sortedDurations.Len()
	if length == 0 {
		return 0
	}

	var sum time.Duration
	for _, d := range sortedDurations {
		sum += d
	}

	return uint32(sum.Nanoseconds()/int64(length)) / 1000000
}

