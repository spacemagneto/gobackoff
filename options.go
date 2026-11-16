package backoff

import (
	"math"
	"time"
)

// Option for exponential strategies Exponential, FullJitter, EqualJitter
type Option func(*strategyOptions)

type strategyOptions struct {
	step float64
}

// WithStep sets growth multiplier between retrie
// Values less than or equal to 1 (and NaN) fall back to DefaultStep so delay still grows
func WithStep(step float64) Option {
	return func(o *strategyOptions) {
		o.step = step
	}
}

func applyOptions(opts []Option) strategyOptions {
	options := strategyOptions{step: DefaultStep}
	for _, opt := range opts {
		opt(&options)
	}

	options.step = normalizeStep(options.step)

	return options
}

// normalizeBounds replaces zero or negative values with defaults
// and raises cap to match base delay so range is never empty
func normalizeBounds(baseDelay, maxDelay time.Duration) (time.Duration, time.Duration) {
	if baseDelay <= 0 {
		baseDelay = DefaultMinBackoff
	}

	if maxDelay <= 0 {
		maxDelay = DefaultMaxBackoff
	}

	if maxDelay < baseDelay {
		maxDelay = baseDelay
	}

	return baseDelay, maxDelay
}

func normalizeStep(step float64) float64 {
	// NOTE: write as !(step > 1) on purpose
	// NaN fails every comparison and must land here too
	if !(step > 1) {
		return DefaultStep
	}

	return step
}

// expLimit returns min(maxDelay, baseDelay * step^attempt).
//
// comparison runs in float before int64 conversion out-of-range float-to-int64
// NOTE: behavior varies in Go, so checking sign after conversion isnt reliable
// Infinity/NaN also fail comparison and safely return maxDelay
func expLimit(baseDelay, maxDelay time.Duration, step float64, attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	limit := float64(baseDelay) * math.Pow(step, float64(attempt))
	if !(limit < float64(maxDelay)) {
		return maxDelay
	}

	return time.Duration(limit)
}
