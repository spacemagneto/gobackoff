package backoff

import (
	"math"
	"time"
)

// Exponential implements a deterministic exponential backoff strategy.
// Unlike Jitter strategies, this implementation is entirely predictable
// and returns a fixed duration for a given attempt number.
//
// The delay follows the formula: baseDelay * (step ^ attempt).
// It is typically used as a baseline for comparison or in scenarios
// where client synchronization is not a significant risk.
type Exponential struct {
	baseDelay, maxDelay time.Duration
	step                float64
}

// NewExponential initializes a new Exponential strategy using the DefaultStep.
// If delay or max are less than or equal to zero, package defaults are used.
func NewExponential(delay, max time.Duration) *Exponential {
	return NewExponentialWithStep(delay, max, DefaultStep)
}

// NewExponentialWithStep initializes an Exponential strategy with a custom growth factor.
//
// Arguments:
//   - delay: The starting duration for the first attempt.
//   - max: The absolute upper bound for any calculated delay.
//   - step: The growth factor (e.g., 2.0 for binary exponential backoff).
//     If step <= 1.0, DefaultStep is used to ensure growth.
func NewExponentialWithStep(delay, max time.Duration, step float64) *Exponential {
	if delay <= 0 {
		delay = DefaultMinBackoff
	}

	if max <= 0 {
		max = DefaultMaxBackoff
	}

	if max < delay {
		max = delay
	}

	if step <= 1.0 {
		step = DefaultStep
	}

	return &Exponential{baseDelay: delay, maxDelay: max, step: step}
}

// Next calculates the delay for the given attempt.
//
// Arguments:
//   - attempt: The current retry attempt number (1, 2, 3...).
//
// Returns:
//
//	A time.Duration representing the calculated delay, capped by maxDelay.
//	If the calculation results in a value larger than the capacity of time.Duration
//	or causes an overflow, maxDelay is returned.
func (e *Exponential) Next(attempt int64) time.Duration {
	// calculate the delay: base * step^attempt
	delay := time.Duration(float64(e.baseDelay) * math.Pow(e.step, float64(attempt)))

	if delay > e.maxDelay || delay < 0 {
		return e.maxDelay
	}

	return delay
}
