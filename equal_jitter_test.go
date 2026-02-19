package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEqualJitter(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exponential := NewEqualJitter(baseDelay, maxDelay)
		assert.NotNil(t, exponential)

		assert.Equal(t, baseDelay, exponential.baseDelay)
		assert.Equal(t, maxDelay, exponential.maxDelay)
	})

	t.Run("InitStrategyWithDefaultValDelay", func(t *testing.T) {
		exponential := NewEqualJitter(0, 0)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultMinBackoff, exponential.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, exponential.maxDelay)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		exponential := NewEqualJitter(0, -10*time.Second)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultMinBackoff, exponential.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, exponential.maxDelay)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		maxDelay := 1 * time.Second
		jitter := NewEqualJitter(baseDelay, maxDelay)

		assert.Equal(t, baseDelay, jitter.maxDelay)
	})

	t.Run("SuccessGetNext", func(t *testing.T) {
		equalJitter := NewEqualJitter(1*time.Second, 20*time.Second)
		delay := equalJitter.Next(1)
		assert.True(t, delay >= 0)
	})

	t.Run("ReturnsNextInStrictRange", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second
		equalJitter := NewEqualJitter(baseDelay, maxDelay)

		// For attempt = 2:
		// limit = 1s * 2^2 = 4s
		// temp = 4s / 2 = 2s
		// The range should be [2s, 4s)

		minExpected := 2 * time.Second
		maxExpected := 4 * time.Second

		for i := 0; i < 15; i++ {
			delay := equalJitter.Next(2)
			assert.True(t, delay >= minExpected, "Delay %v should be >= %v", delay, minExpected)
			assert.True(t, delay < maxExpected, "Delay %v should be < %v", delay, maxExpected)
		}
	})

	t.Run("GetNextWithMaxDelay", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 3 * time.Second
		equalJitter := NewEqualJitter(baseDelay, maxDelay)

		// Attempt 10: limit greater than max (3s)
		// temp = 3s / 2 = 1.5s
		// Range: [1.5s, 3s)

		minExpected := 1500 * time.Millisecond

		for i := 0; i < 15; i++ {
			delay := equalJitter.Next(10)
			assert.GreaterOrEqual(t, delay, minExpected)
			assert.Less(t, delay, maxDelay)
		}
	})
}
