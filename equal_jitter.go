package backoff

import (
	"math/rand/v2"
	"time"
)

// EqualJitter implements a strategy that combines exponential backoff with a
// randomized component, ensuring that the delay is always at least half of the
// current exponential limit.
//
// This strategy provides a more predictable lower bound for wait times compared
// to Full Jitter while still preventing synchronization between clients.
//
// Reference: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
type EqualJitter struct {
	baseDelay, maxDelay time.Duration
}

// NewEqualJitter initializes a new EqualJitter strategy.
// If delay or max are less than or equal to zero, package defaults are used.
// If max is less than delay, max is set to delay.
func NewEqualJitter(delay, max time.Duration) *EqualJitter {
	if delay <= 0 {
		delay = DefaultMinBackoff
	}

	if max <= 0 {
		max = DefaultMaxBackoff
	}

	if max < delay {
		max = delay
	}

	return &EqualJitter{baseDelay: delay, maxDelay: max}
}

// Next calculates the next delay using the Equal Jitter formula:
// temp = min(cap, base * 2^attempt)
// sleep = temp/2 + random_between(0, temp/2)
//
// Arguments:
//   - attempt: The current retry attempt number (1, 2, 3...).
//
// Returns:
//
//	A randomized time.Duration between limit/2 and limit.
func (e *EqualJitter) Next(attempt int64) time.Duration {
	// temp = min(cap, base * 2 ** attempt)
	// sleep = temp / 2 + random_between(0, temp / 2)
	// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
	limit := int64(e.baseDelay) * (1 << uint64(attempt))

	if limit > int64(e.maxDelay) || limit <= 0 {
		limit = int64(e.maxDelay)
	}

	temp := limit / 2

	var delay int64
	if temp > 0 {
		delay = rand.Int64N(temp)
	}

	return time.Duration(temp + delay)
}
