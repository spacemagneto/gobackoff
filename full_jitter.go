package backoff

import (
	"math/rand/v2"
	"time"
)

type FullJitter struct {
	baseDelay, maxDelay time.Duration
}

func NewFullJitter(delay, max time.Duration) *FullJitter {
	if delay <= 0 {
		delay = DefaultMinBackoff
	}

	if max <= 0 {
		max = DefaultMaxBackoff
	}

	if max < delay {
		max = delay
	}

	return &FullJitter{baseDelay: delay, maxDelay: max}
}

func (f *FullJitter) Next(attempt int64) time.Duration {
	// sleep = random_between(0, min(cap, base * 2^attempt))
	limit := int64(f.baseDelay) * (1 << uint64(attempt))

	if limit > int64(f.maxDelay) || limit <= 0 {
		limit = int64(f.maxDelay)
	}

	if limit == 0 {
		return 0
	}

	return time.Duration(rand.Int64N(limit + 1))
}
