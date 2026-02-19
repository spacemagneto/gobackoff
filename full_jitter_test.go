package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFullJitterStrategy(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exponential := NewFullJitter(baseDelay, maxDelay)
		assert.NotNil(t, exponential)

		assert.Equal(t, baseDelay, exponential.baseDelay)
		assert.Equal(t, maxDelay, exponential.maxDelay)
	})

	t.Run("InitStrategyWithDefaultValDelay", func(t *testing.T) {
		exponential := NewFullJitter(0, 0)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultMinBackoff, exponential.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, exponential.maxDelay)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		exponential := NewFullJitter(0, -10*time.Second)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultMinBackoff, exponential.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, exponential.maxDelay)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		maxDelay := 1 * time.Second
		jitter := NewFullJitter(baseDelay, maxDelay)

		assert.Equal(t, baseDelay, jitter.maxDelay)
	})

	t.Run("ReturnsNextInStrictRange", func(t *testing.T) {
		baseDelay := 100 * time.Millisecond
		maxDelay := 10 * time.Second
		jitter := NewFullJitter(baseDelay, maxDelay)

		// For attempt = 3:
		// limit = 100ms * 2^3 = 800ms
		// Expected range [0, 800ms]
		maxExpected := 800 * time.Millisecond

		for i := 0; i < 15; i++ {
			delay := jitter.Next(3)
			assert.True(t, delay >= 0)
			assert.True(t, delay <= maxExpected)
		}
	})
}
