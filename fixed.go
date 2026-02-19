package backoff

import (
	"time"
)

// Fixed implements a static backoff strategy.
// It always returns the same duration, regardless of the attempt count
// or previous delay state. This is useful for simple polling or
// scenarios where exponential growth is not desired.
type Fixed struct {
	delay time.Duration
}

// NewFixed initializes a new Fixed strategy with a constant delay.
//
// Arguments:
//   - delay: The duration to be returned on every call to Next.
//     If delay is less than or equal to zero, it defaults to 1 second.
func NewFixed(delay time.Duration) *Fixed {
	if delay <= 0 {
		delay = DefaultMinBackoff
	}

	return &Fixed{delay: delay}
}

// Next returns the fixed delay duration.
//
// Arguments:
//   - arg: This argument is ignored by the Fixed strategy as the
//     delay does not change based on attempts or state.
//
// Returns:
//
//	The constant time.Duration configured during initialization.
func (f *Fixed) Next(_ int64) time.Duration {
	return f.delay
}
