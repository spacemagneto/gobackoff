package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialStrategy(t *testing.T) {
	t.Parallel()

	t.Run("InitStrategyWithDefaultStep", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exponential := NewExponential(baseDelay, maxDelay)
		assert.NotNil(t, exponential)

		assert.Equal(t, baseDelay, exponential.baseDelay)
		assert.Equal(t, maxDelay, exponential.maxDelay)
		assert.Equal(t, DefaultStep, exponential.step)
	})

	t.Run("InitStrategyWithCustomStep", func(t *testing.T) {
		exponential := NewExponential(time.Second, 10*time.Second, WithStep(1.5))
		assert.Equal(t, 1.5, exponential.step)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		exponential := NewExponential(0, -1*time.Second, WithStep(-0.5))

		assert.Equal(t, DefaultMinBackoff, exponential.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, exponential.maxDelay)
		assert.Equal(t, DefaultStep, exponential.step)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		exponential := NewExponential(baseDelay, time.Second)

		assert.Equal(t, baseDelay, exponential.maxDelay)
	})

	t.Run("SequenceStartsAtBaseAndStopsAtCap", func(t *testing.T) {
		exponential := NewExponential(time.Second, 10*time.Second)

		expected := []time.Duration{
			1 * time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second,
			10 * time.Second, 10 * time.Second,
		}
		for i, want := range expected {
			assert.Equal(t, want, exponential.Next(), "call %d", i)
		}
	})

	t.Run("SequenceWithCustomStep", func(t *testing.T) {
		exponential := NewExponential(time.Second, time.Minute, WithStep(1.5))

		assert.Equal(t, 1000*time.Millisecond, exponential.Next())
		assert.Equal(t, 1500*time.Millisecond, exponential.Next())
		assert.Equal(t, 2250*time.Millisecond, exponential.Next())
	})

	t.Run("AttemptStopsGrowingAtCap", func(t *testing.T) {
		exponential := NewExponential(time.Second, 4*time.Second)

		for i := 0; i < 100; i++ {
			exponential.Next()
		}

		// 1s, 2s, 4s counter reaches cap on attempt 2 and stops 3
		assert.Equal(t, 2, exponential.attempt)
	})

	t.Run("ResetStartsOverFromBase", func(t *testing.T) {
		exponential := NewExponential(time.Second, 10*time.Second)
		exponential.Next()
		exponential.Next()

		exponential.Reset()
		assert.Equal(t, time.Second, exponential.Next())
	})
}
