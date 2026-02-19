package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJitterStrategy(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exponential := NewJitter(baseDelay, maxDelay)
		assert.NotNil(t, exponential)

		assert.Equal(t, baseDelay, exponential.baseDelay)
		assert.Equal(t, maxDelay, exponential.maxDelay)
	})

	t.Run("InitStrategyWithDefaultValDelay", func(t *testing.T) {
		exponential := NewJitter(0, 0)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultMinBackoff, exponential.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, exponential.maxDelay)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		exponential := NewJitter(0, -10*time.Second)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultMinBackoff, exponential.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, exponential.maxDelay)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		maxDelay := 1 * time.Second
		jitter := NewJitter(baseDelay, maxDelay)

		assert.Equal(t, baseDelay, jitter.maxDelay)
	})

	t.Run("SuccessGetNext", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		jitter := NewJitter(baseDelay, maxDelay)

		// For attempt 2:
		// Exponent limit = 1s * 2^2 = 4s.
		// Full Jitter should return a value in the range [0, 4s).
		for i := 0; i < 100; i++ {
			delay := jitter.Next(2)

			assert.GreaterOrEqual(t, int64(delay), int64(0))
			assert.Less(t, int64(delay), int64(4*time.Second))
		}
	})

	t.Run("SuccessGetNextAndCheckRandomNextDelay", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		jitter := NewJitter(baseDelay, maxDelay)

		attempt := int64(5)
		first := jitter.Next(attempt)
		second := jitter.Next(attempt)

		assert.NotEqual(t, first, second)
	})

	t.Run("OverflowProtection", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 1 * time.Hour

		jitter := NewJitter(baseDelay, maxDelay)

		delay := jitter.Next(100)
		assert.LessOrEqual(t, delay, jitter.maxDelay)
		assert.GreaterOrEqual(t, delay, time.Duration(0))
	})
}
