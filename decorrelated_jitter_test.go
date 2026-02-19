package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecorrelatedJitter(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		decorrelatedJitter := NewDecorrelatedJitter(baseDelay, maxDelay)
		assert.NotNil(t, decorrelatedJitter)

		assert.Equal(t, baseDelay, decorrelatedJitter.baseDelay)
		assert.Equal(t, maxDelay, decorrelatedJitter.maxDelay)
	})

	t.Run("InitStrategyWithDefaultValDelay", func(t *testing.T) {
		decorrelatedJitter := NewDecorrelatedJitter(0, 0)
		assert.NotNil(t, decorrelatedJitter)

		assert.Equal(t, DefaultMinBackoff, decorrelatedJitter.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, decorrelatedJitter.maxDelay)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		decorrelatedJitter := NewDecorrelatedJitter(0, -10*time.Second)
		assert.NotNil(t, decorrelatedJitter)

		assert.Equal(t, DefaultMinBackoff, decorrelatedJitter.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, decorrelatedJitter.maxDelay)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		maxDelay := 1 * time.Second

		decorrelatedJitter := NewDecorrelatedJitter(baseDelay, maxDelay)

		assert.Equal(t, baseDelay, decorrelatedJitter.maxDelay)
	})

	t.Run("SuccessGetNext", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		decorrelatedJitter := NewDecorrelatedJitter(baseDelay, maxDelay)
		delay := decorrelatedJitter.Next(1)
		assert.True(t, delay >= 0)
	})

	t.Run("CheckMinDelay", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		decorrelatedJitter := NewDecorrelatedJitter(baseDelay, maxDelay)
		assert.Equal(t, baseDelay, decorrelatedJitter.Next(0))
	})

	t.Run("ReturnsNextInStrictRange", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		decorrelatedJitter := NewDecorrelatedJitter(baseDelay, maxDelay)

		// If the previous delay was 2 seconds:
		// high = min(10s, 2s * 3) = 6s
		// The range should be [baseDelay, high], i.e. [1s, 6s]
		minExpected := 1 * time.Second
		maxExpected := 6 * time.Second

		for i := 0; i < 15; i++ {
			delay := decorrelatedJitter.Next(int64(2 * time.Second))
			assert.True(t, delay >= minExpected, "Delay %v should be >= %v", delay, minExpected)
			assert.True(t, delay < maxExpected, "Delay %v should be < %v", delay, maxExpected)
		}
	})

	t.Run("ReturnsNextWithHighExceedsMaxDelay", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 5 * time.Second

		decorrelatedJitter := NewDecorrelatedJitter(baseDelay, maxDelay)

		// When transferring 10s: high = 10s * 3 = 30s.
		// 30s > 5s (maxDelay), so high will be 5s.
		// Expected range: [1s, 5s]
		for i := 0; i < 15; i++ {
			delay := decorrelatedJitter.Next(int64(10 * time.Second))
			assert.True(t, delay >= baseDelay)
			assert.True(t, delay <= maxDelay)
		}
	})
}
