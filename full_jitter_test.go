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

		fullJitter := NewFullJitter(baseDelay, maxDelay)
		assert.NotNil(t, fullJitter)

		assert.Equal(t, baseDelay, fullJitter.baseDelay)
		assert.Equal(t, maxDelay, fullJitter.maxDelay)
	})

	t.Run("InitStrategyWithDefaultValDelay", func(t *testing.T) {
		fullJitter := NewFullJitter(0, 0)
		assert.NotNil(t, fullJitter)

		assert.Equal(t, DefaultMinBackoff, fullJitter.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, fullJitter.maxDelay)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		fullJitter := NewFullJitter(0, -10*time.Second)
		assert.NotNil(t, fullJitter)

		assert.Equal(t, DefaultMinBackoff, fullJitter.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, fullJitter.maxDelay)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		maxDelay := 1 * time.Second
		fullJitter := NewFullJitter(baseDelay, maxDelay)

		assert.Equal(t, baseDelay, fullJitter.maxDelay)
	})

	t.Run("ReturnsNextInStrictRange", func(t *testing.T) {
		baseDelay := 100 * time.Millisecond
		maxDelay := 10 * time.Second
		fullJitter := NewFullJitter(baseDelay, maxDelay)

		// For attempt = 3:
		// limit = 100ms * 2^3 = 800ms
		// Expected range [0, 800ms]
		maxExpected := 800 * time.Millisecond

		for i := 0; i < 15; i++ {
			delay := fullJitter.Next(3)
			assert.True(t, delay >= 0)
			assert.True(t, delay <= maxExpected)
		}
	})

	t.Run("CheckOverflow", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 100 * time.Second
		fullJitter := NewFullJitter(baseDelay, maxDelay)

		delay := fullJitter.Next(63)
		assert.LessOrEqual(t, delay, maxDelay)
		assert.GreaterOrEqual(t, delay, time.Duration(0))
	})

	t.Run("CheckMinDelay", func(t *testing.T) {
		fullJitter := &FullJitter{baseDelay: 0, maxDelay: 0}

		delay := fullJitter.Next(5)
		assert.Equal(t, fullJitter.baseDelay, delay)
	})
}
