package backoff

import (
	"time"
)

// Fixed implements static backoff strategy
// always returns same duration need only for simple polling or when delay growth isnt needed
type Fixed struct {
	delay time.Duration
}

// NewFixed creates Fixed strategy with constant delay
// Falls back to DefaultMinBackoff if delay is zero or negative
func NewFixed(delay time.Duration) *Fixed {
	if delay <= 0 {
		delay = DefaultMinBackoff
	}

	return &Fixed{delay: delay}
}

// Next returns the fixed delay
func (f *Fixed) Next() time.Duration {
	return f.delay
}

// Reset mock, implement only next
func (f *Fixed) Reset() {}
