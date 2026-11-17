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

// Backoff provides delays for single retry chain
//
// Strategy tracks its own state (attempt counter or last delay)
// so caller just asks for next delay without track any
// One instance serves single chain and isnt safe for concurrent use
// create new instance per chain, or call Reset when chain finish
type Backoff interface {
	// Next returns the delay before the next retry.
	// It never exceeds the cap set at initialization.
	Next() time.Duration

	// Reset starts the chain over: the next call to Next returns the first delay again.
	Reset()
}
