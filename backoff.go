package backoff

import "time"

const (
	// DefaultMinBackoff is the starting delay used if a strategy
	// is initialized with a zero or negative duration.
	DefaultMinBackoff = 1 * time.Second

	// DefaultMaxBackoff is the upper limit for any delay.
	// No strategy will return a value exceeding this threshold.
	DefaultMaxBackoff = 20 * time.Second

	// DefaultStep is the base factor for exponential growth.
	// A value of 2.0 represents a standard binary exponential backoff.
	DefaultStep = 2.0
)

// Backoff defines the contract for various retry delay strategies.
// It encapsulates the logic for calculating wait times between retries.
type Backoff interface {
	// Next calculates the duration of the next delay.
	//
	// The meaning of the input argument depends on the implementation:
	//   - For deterministic or exponential-based strategies (Exponential, FullJitter, EqualJitter),
	//     it represents the current retry attempt number (1, 2, 3...).
	//   - For state-aware strategies (DecorrelatedJitter), it represents the
	//     previous delay duration in nanoseconds (time.Duration.Nanoseconds()).
	//   - For static strategies (Fixed), the argument is ignored.
	//
	// The returned duration is guaranteed to be capped by the maxDelay
	// configured during the strategy's initialization.
	Next(arg int64) time.Duration
}
