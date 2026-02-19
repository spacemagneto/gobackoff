package backoff

import (
	"math/rand/v2"
	"time"
)

// DecorrelatedJitter implements a jitter strategy where the next delay
// is based on the previous delay rather than the retry attempt number.
//
// This strategy is known for its "random walk" behavior, providing excellent
// desynchronization for clients in large-scale distributed systems.
//
// Reference: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
type DecorrelatedJitter struct {
	baseDelay time.Duration
	maxDelay  time.Duration
}

// NewDecorrelatedJitter initializes a new DecorrelatedJitter strategy.
// If delay or max are less than or equal to zero, package defaults are used.
// If max is less than delay, max is set to delay.
func NewDecorrelatedJitter(delay, max time.Duration) *DecorrelatedJitter {
	if delay <= 0 {
		delay = DefaultMinBackoff
	}

	if max <= 0 {
		max = DefaultMaxBackoff
	}

	if max < delay {
		max = delay
	}

	return &DecorrelatedJitter{baseDelay: delay, maxDelay: max}
}

// Next calculates the next delay using the Decorrelated Jitter formula:
// sleep = min(cap, random_between(base, previous_sleep * 3))
//
// Arguments:
//   - previousDelay: The duration of the last sleep in nanoseconds.
//     Pass 0 for the first attempt.
//
// Returns:
//
//	A randomized time.Duration between baseDelay and min(maxDelay, previousDelay * 3).
func (d *DecorrelatedJitter) Next(previousDelay int64) time.Duration {
	// sleep = min(cap, random_between(base, sleep * 3))
	high := time.Duration(previousDelay) * 3
	if high < d.baseDelay {
		high = d.baseDelay
	}

	if high > d.maxDelay {
		high = d.maxDelay
	}

	diff := int64(high - d.baseDelay)
	if diff <= 0 {
		return d.baseDelay
	}

	return d.baseDelay + time.Duration(rand.Int64N(diff+1))
}
