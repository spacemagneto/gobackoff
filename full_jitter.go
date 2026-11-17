package backoff

import (
	"math"
	"math/rand/v2"
	"time"
)

// FullJitter implements exponential backoff with fully random delay
// Computes exponential limit for each attempt and returns random value between 0 and that limit
//
// maximizes client spread to prevent thundering herd
// tradeoff is no lower bound, so delay can be close to zero.
//
// Reference: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
type FullJitter struct {
	baseDelay time.Duration
	maxDelay  time.Duration
	step      float64

	// NOTE: why i dont use mutex, atomic, why? isnt not solved race cond problem
	// state belongs to retry chain, isnt shared strategy
	// look if first goroutine got err Next return 1s,
	// second goroutine got error Next return 2s instead of 1s
	// under load counter hits cap immediately and all goroutines sleep max delay
	attempt int
}

// NewFullJitter creates FullJitter strategy.
// use pkg def if baseDelay or maxDelay is zero or negative and raises maxDelay to baseDelay if smaller.
// growth factor def to DefaultStep unless WithStep is passed
func NewFullJitter(baseDelay, maxDelay time.Duration, opts ...Option) *FullJitter {
	baseDelay, maxDelay = normalizeBounds(baseDelay, maxDelay)
	options := applyOptions(opts)

	return &FullJitter{baseDelay: baseDelay, maxDelay: maxDelay, step: options.step}
}

// Next returns rand delay in [0, min(cap, base * step^attempt)] and advances to next attempt
func (f *FullJitter) Next() time.Duration {
	limit := expLimit(f.baseDelay, f.maxDelay, f.step, f.attempt)

	// counter stops at cap, same reason as in exponential
	if limit < f.maxDelay {
		f.attempt++
	}

	return fullJitter(limit)
}

// Reset restarts chain from first attempt
func (f *FullJitter) Reset() {
	f.attempt = 0
}

func fullJitter(limit time.Duration) time.Duration {
	// limit+1 would overflow at the very top, and Int64 already covers [0, MaxInt64]
	if limit == math.MaxInt64 {
		return time.Duration(rand.Int64())
	}

	return time.Duration(rand.Int64N(int64(limit) + 1))
}
