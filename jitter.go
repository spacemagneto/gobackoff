package backoff

import (
	"math/rand/v2"
	"time"
)

type Jitter struct {
	baseDelay, maxDelay time.Duration
}

func NewJitter(delay, max time.Duration) *Jitter {
	if delay <= 0 {
		delay = DefaultMinBackoff
	}

	if max <= 0 {
		max = DefaultMaxBackoff
	}

	if max < delay {
		max = delay
	}

	return &Jitter{baseDelay: delay, maxDelay: max}
}

func (j *Jitter) Next(attempt int64) time.Duration {
	// sleep = random_between(0, min(cap, base * 2 ** attempt))
	// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
	limit := int64(j.baseDelay) * (1 << uint64(attempt))

	if limit > int64(j.maxDelay) || limit <= 0 {
		limit = int64(j.maxDelay)
	}

	return time.Duration(rand.Int64N(limit))
}
