package backoff

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecorrelatedJitter(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		jitter := NewDecorrelatedJitter(baseDelay, maxDelay)
		assert.NotNil(t, jitter)

		assert.Equal(t, baseDelay, jitter.baseDelay)
		assert.Equal(t, maxDelay, jitter.maxDelay)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		jitter := NewDecorrelatedJitter(0, -10*time.Second)

		assert.Equal(t, DefaultMinBackoff, jitter.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, jitter.maxDelay)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		jitter := NewDecorrelatedJitter(baseDelay, time.Second)

		assert.Equal(t, baseDelay, jitter.maxDelay)
	})

	t.Run("FirstCallReturnsBase", func(t *testing.T) {
		decoratedJitter := NewDecorrelatedJitter(time.Second, 10*time.Second)
		assert.Equal(t, time.Second, decoratedJitter.Next())
	})

	t.Run("EachCallInRangeOfPrevious", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		for trial := 0; trial < 50; trial++ {
			jitter := NewDecorrelatedJitter(baseDelay, maxDelay)

			prev := time.Duration(0)
			for i := 0; i < 20; i++ {
				high := min(maxDelay, max(baseDelay, prev*3))
				delay := jitter.Next()

				assert.GreaterOrEqual(t, delay, baseDelay)
				assert.LessOrEqual(t, delay, high, "call %d after %v", i, prev)

				prev = delay
			}
		}
	})

	t.Run("ResetForgetsPrevious", func(t *testing.T) {
		jitter := NewDecorrelatedJitter(time.Second, 10*time.Second)
		for i := 0; i < 5; i++ {
			jitter.Next()
		}

		jitter.Reset()
		assert.Equal(t, time.Second, jitter.Next())
	})

	t.Run("FormulaInRange", func(t *testing.T) {
		// prev 2s: range [1s, min(10s, 6s)] = [1s, 6s]
		for i := 0; i < 50; i++ {
			delay := decorrelatedJitter(time.Second, 10*time.Second, 2*time.Second)
			assert.GreaterOrEqual(t, delay, time.Second)
			assert.LessOrEqual(t, delay, 6*time.Second)
		}

		// prev 10s: prev*3 above cap, range [1s, 5s]
		for i := 0; i < 50; i++ {
			delay := decorrelatedJitter(time.Second, 5*time.Second, 10*time.Second)
			assert.GreaterOrEqual(t, delay, time.Second)
			assert.LessOrEqual(t, delay, 5*time.Second)
		}
	})

	t.Run("HugePreviousDoesNotOverflow", func(t *testing.T) {
		maxDelay := time.Duration(math.MaxInt64)
		delay := decorrelatedJitter(time.Second, maxDelay, maxDelay)

		assert.GreaterOrEqual(t, delay, time.Second)
	})
}
