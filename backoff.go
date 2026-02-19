package backoff

import "time"

const (
	DefaultMinBackoff = 1 * time.Second
	DefaultMaxBackoff = 20 * time.Second
	DefaultStep       = 2.0
)

type Backoff interface {
	Next(int64) time.Duration
}
