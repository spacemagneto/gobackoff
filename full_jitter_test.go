package backoff

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFullJitterStrategy(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		jitter := NewFullJitter(baseDelay, maxDelay)
		assert.NotNil(t, jitter)

		assert.Equal(t, baseDelay, jitter.baseDelay)
		assert.Equal(t, maxDelay, jitter.maxDelay)
		assert.Equal(t, DefaultStep, jitter.step)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		jitter := NewFullJitter(0, -10*time.Second, WithStep(1))

		assert.Equal(t, DefaultMinBackoff, jitter.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, jitter.maxDelay)
		assert.Equal(t, DefaultStep, jitter.step)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		jitter := NewFullJitter(baseDelay, time.Second)

		assert.Equal(t, baseDelay, jitter.maxDelay)
	})

	t.Run("EachCallInRangeOfItsAttempt", func(t *testing.T) {
		baseDelay := 100 * time.Millisecond
		maxDelay := 2 * time.Second

		for trial := 0; trial < 50; trial++ {
			jitter := NewFullJitter(baseDelay, maxDelay)

			for attempt := 0; attempt < 10; attempt++ {
				limit := expLimit(baseDelay, maxDelay, DefaultStep, attempt)
				delay := jitter.Next()

				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, limit, "attempt %d", attempt)
			}
		}
	})

	t.Run("CustomStepChangesLimit", func(t *testing.T) {
		// third call: limit = 1s * 3^2 = 9s
		for trial := 0; trial < 50; trial++ {
			jitter := NewFullJitter(time.Second, time.Minute, WithStep(3))
			jitter.Next()
			jitter.Next()

			assert.LessOrEqual(t, jitter.Next(), 9*time.Second)
		}
	})

	t.Run("ResetStartsOverFromFirstAttempt", func(t *testing.T) {
		jitter := NewFullJitter(time.Second, 10*time.Second)
		for i := 0; i < 5; i++ {
			jitter.Next()
		}

		jitter.Reset()
		assert.Equal(t, 0, jitter.attempt)
		assert.LessOrEqual(t, jitter.Next(), time.Second)
	})

	t.Run("FormulaInRange", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			delay := fullJitter(800 * time.Millisecond)
			assert.GreaterOrEqual(t, delay, time.Duration(0))
			assert.LessOrEqual(t, delay, 800*time.Millisecond)
		}
	})

	t.Run("LimitAtMaxDurationDoesNotOverflow", func(t *testing.T) {
		delay := fullJitter(time.Duration(math.MaxInt64))
		assert.GreaterOrEqual(t, delay, time.Duration(0))
	})
}
