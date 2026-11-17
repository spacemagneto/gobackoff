package backoff

import (
	"math/rand/v2"
	"time"
)

// EqualJitter implements exponential backoff where half of delay is fixed
// and half is random, keeping delay at least half of current limit
//
// # Offers higher lower bound than FullJitter while still spreading clients apart
//
// Reference: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
type EqualJitter struct {
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

// NewEqualJitter creates EqualJitter strategy
// use pkg def if baseDelay or maxDelay is zero or negative and raises maxDelay to baseDelay if smaller
// growth factor def to DefaultStep unless WithStep is passed
func NewEqualJitter(baseDelay, maxDelay time.Duration, opts ...Option) *EqualJitter {
	baseDelay, maxDelay = normalizeBounds(baseDelay, maxDelay)
	options := applyOptions(opts)

	return &EqualJitter{baseDelay: baseDelay, maxDelay: maxDelay, step: options.step}
}

// Next returns random delay in [limit/2, limit], where limit = min(cap, base * step^attempt) and advances to next attempt
func (e *EqualJitter) Next() time.Duration {
	limit := expLimit(e.baseDelay, e.maxDelay, e.step, e.attempt)

	// counter stops at cap, same reason as in exponential
	if limit < e.maxDelay {
		e.attempt++
	}

	return equalJitter(limit)
}

// Reset restarts chain from first attempt
func (e *EqualJitter) Reset() {
	e.attempt = 0
}

func equalJitter(limit time.Duration) time.Duration {
	half := limit / 2

	// random part covers limit-half instead of half: with odd limit top nanosecond would be lost
	// limit-half is at most MaxInt64/2+1, so +1 never overflows
	return half + time.Duration(rand.Int64N(int64(limit-half)+1))
}
