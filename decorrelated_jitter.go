package backoff

import (
	"math/rand/v2"
	"time"
)

// DecorrelatedJitter implements jitter strategy where next delay
// depends on previous delay rather than attempt number.
//
// Produces random walk delay grows on average, but can also drop back toward base
// Fits desynchronization in large-scale distributed systems.
//
// Formula differs from AWS blog on purpose AWS caps after drawing:
// min(cap, random_between(base, prev*3)), so once prev*3 passes cap
// most draws collapse onto exact cap and clients synchronize again.
//
// Here cap applies to range: random_between(base, min(cap, prev*3)),
// keeping delays spread across whole range.
//
// Reference: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
type DecorrelatedJitter struct {
	baseDelay time.Duration
	maxDelay  time.Duration
	prev      time.Duration
}

// NewDecorrelatedJitter creates DecorrelatedJitter strategy
// use pkg def if baseDelay or maxDelay is zero or negative and raises maxDelay to baseDelay if small
func NewDecorrelatedJitter(baseDelay, maxDelay time.Duration) *DecorrelatedJitter {
	baseDelay, maxDelay = normalizeBounds(baseDelay, maxDelay)

	return &DecorrelatedJitter{baseDelay: baseDelay, maxDelay: maxDelay}
}

// Next returns random delay in [base, min(cap, previous*3)] and stores it
// First call after creation or Reset returns exact base
func (d *DecorrelatedJitter) Next() time.Duration {
	d.prev = decorrelatedJitter(d.baseDelay, d.maxDelay, d.prev)

	return d.prev
}

// Reset clears previous delay so next call returns base again
func (d *DecorrelatedJitter) Reset() {
	d.prev = 0
}

func decorrelatedJitter(baseDelay, maxDelay, prev time.Duration) time.Duration {
	// compare against cap/3 instead of multiplying first: prev*3 can overflow
	high := maxDelay
	if prev <= maxDelay/3 {
		high = prev * 3
	}

	if high < baseDelay {
		high = baseDelay
	}

	// base is at least 1ns, so high-base+1 dnever overflow
	return baseDelay + time.Duration(rand.Int64N(int64(high-baseDelay)+1))
}
