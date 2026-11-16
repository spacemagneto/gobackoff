package backoff

import (
	"time"
)

// Exponential implements deterministic exponential backoff strategy.
// No randomness: every chain gets same sequence of delays.
//
// Delay follows formula: min(cap, base * step^attempt), attempts counted from 0 so first delay equals base
// Fits baseline use or cases where client synchronization isnt risk under contention prefer jitter strategy
type Exponential struct {
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

// NewExponential creates Exponential strategy.
// use pkg def if baseDelay or maxDelay is zero or negative
// raises maxDelay to match baseDelay if maxDelay is smaller
// growth factor def to DefaultStep unless WithStep is passed
func NewExponential(baseDelay, maxDelay time.Duration, opts ...Option) *Exponential {
	baseDelay, maxDelay = normalizeBounds(baseDelay, maxDelay)
	o := applyOptions(opts)

	return &Exponential{baseDelay: baseDelay, maxDelay: maxDelay, step: o.step}
}

// Next returns delay for current attempt and advances to next one
func (e *Exponential) Next() time.Duration {
	delay := expLimit(e.baseDelay, e.maxDelay, e.step, e.attempt)

	// counter stop at cap further growth changes nothing
	// and endless retry loop would eventually overflow it
	if delay < e.maxDelay {
		e.attempt++
	}

	return delay
}

// Reset restarts chain from base delay
func (e *Exponential) Reset() {
	e.attempt = 0
}

// ExponentialDelay is stateless version of Exponential for callers tracking attempt number themselves
// arguments are normalized same way as in NewExponential and WithStep negative attempt counts as 0
func ExponentialDelay(baseDelay, maxDelay time.Duration, step float64, attempt int) time.Duration {
	baseDelay, maxDelay = normalizeBounds(baseDelay, maxDelay)

	return expLimit(baseDelay, maxDelay, normalizeStep(step), attempt)
}
