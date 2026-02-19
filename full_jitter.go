package backoff

import (
	"math/rand/v2"
	"time"
)

// FullJitter implements an exponential backoff strategy with a fully randomized delay.
// This strategy calculates the maximum possible delay for the current attempt
// and returns a random value between 0 and that maximum.
//
// Full Jitter is highly recommended for large-scale distributed systems because
// it provides the best protection against "thundering herd" issues by ensuring
// maximum client desynchronization.
//
// Reference: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
type FullJitter struct {
	baseDelay, maxDelay time.Duration
}

// NewFullJitter initializes a new FullJitter strategy.
// If delay or max are less than or equal to zero, package defaults are used.
// If max is less than delay, max is set to delay.
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

// Next calculates the next delay using the Full Jitter formula:
// sleep = random_between(0, min(cap, base * 2^attempt))
//
// Arguments:
//   - attempt: The current retry attempt number (1, 2, 3...).
//
// Returns:
//
//	A randomized time.Duration between min delay and the current exponential limit.
//	The result is guaranteed to be capped by maxDelay.
func (f *FullJitter) Next(attempt int64) time.Duration {
	// sleep = random_between(0, min(cap, base * 2^attempt))
	limit := int64(f.baseDelay) * (1 << uint64(attempt))

	if limit > int64(f.maxDelay) || limit <= 0 {
		limit = int64(f.maxDelay)
	}

	if limit == 0 {
		return f.baseDelay
	}

	return time.Duration(rand.Int64N(limit + 1))
}
