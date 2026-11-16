package backoff

import (
	"time"
)

const (
	// DefaultMinBackoff is starting delay when strategy gets zero or negative duration
	DefaultMinBackoff = 1 * time.Second

	// DefaultMaxBackoff is upper limit when strategy gets zero or negative cap
	DefaultMaxBackoff = 20 * time.Second

	// DefaultStep is multiplier for exponential growth (2.0 doubles delay each retry)
	DefaultStep = 2.0
)

// Backoff hands out delays for one chain of retries.
//
// Strategy keeps its own state (attempt counter or previous delay), so caller
// just asks for the next delay and never has to track anything. The flip side:
// one instance serves one chain and is not safe for concurrent use.
// Create a new instance per chain, or call Reset when the chain is over.
type Backoff interface {
	// Next returns the delay before the next retry.
	// It never exceeds the cap set at initialization.
	Next() time.Duration

	// Reset starts the chain over: the next call to Next returns the first delay again.
	Reset()
}
