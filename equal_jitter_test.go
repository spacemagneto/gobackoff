package backoff

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEqualJitter(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		jitter := NewEqualJitter(baseDelay, maxDelay)
		assert.NotNil(t, jitter)

		assert.Equal(t, baseDelay, jitter.baseDelay)
		assert.Equal(t, maxDelay, jitter.maxDelay)
		assert.Equal(t, DefaultStep, jitter.step)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		jitter := NewEqualJitter(0, -10*time.Second, WithStep(0))

		assert.Equal(t, DefaultMinBackoff, jitter.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, jitter.maxDelay)
		assert.Equal(t, DefaultStep, jitter.step)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		jitter := NewEqualJitter(baseDelay, time.Second)

		assert.Equal(t, baseDelay, jitter.maxDelay)
	})

	t.Run("EachCallInRangeOfItsAttempt", func(t *testing.T) {
		baseDelay := 100 * time.Millisecond
		maxDelay := 3 * time.Second

		for trial := 0; trial < 50; trial++ {
			jitter := NewEqualJitter(baseDelay, maxDelay)

			for attempt := 0; attempt < 10; attempt++ {
				limit := expLimit(baseDelay, maxDelay, DefaultStep, attempt)
				delay := jitter.Next()

				assert.GreaterOrEqual(t, delay, limit/2, "attempt %d", attempt)
				assert.LessOrEqual(t, delay, limit, "attempt %d", attempt)
			}
		}
	})

	t.Run("ResetStartsOverFromFirstAttempt", func(t *testing.T) {
		jitter := NewEqualJitter(time.Second, 10*time.Second)
		for i := 0; i < 5; i++ {
			jitter.Next()
		}

		jitter.Reset()
		assert.Equal(t, 0, jitter.attempt)

		delay := jitter.Next()
		assert.GreaterOrEqual(t, delay, 500*time.Millisecond)
		assert.LessOrEqual(t, delay, time.Second)
	})

	t.Run("FormulaInRange", func(t *testing.T) {
		// limit 4s: range [2s, 4s]
		for i := 0; i < 50; i++ {
			delay := equalJitter(4 * time.Second)
			assert.GreaterOrEqual(t, delay, 2*time.Second)
			assert.LessOrEqual(t, delay, 4*time.Second)
		}
	})

	t.Run("LimitAtMaxDurationDoesNotOverflow", func(t *testing.T) {
		delay := equalJitter(time.Duration(math.MaxInt64))
		assert.GreaterOrEqual(t, delay, time.Duration(math.MaxInt64/2))
	})
}
