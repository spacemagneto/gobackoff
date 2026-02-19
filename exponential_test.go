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
		assert.Equal(t, float64(2), exponential.step)
	})

	t.Run("InitStrategyWithCustomStep", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second
		step := 1.5

		exponential := NewExponentialWithStep(baseDelay, maxDelay, step)
		assert.NotNil(t, exponential)

		assert.Equal(t, baseDelay, exponential.baseDelay)
		assert.Equal(t, maxDelay, exponential.maxDelay)
		assert.Equal(t, step, exponential.step)
	})

	t.Run("InitStrategyWithInvalidStep", func(t *testing.T) {
		exp := NewExponentialWithStep(time.Second, time.Minute, -0.5)
		assert.Equal(t, 2.0, exp.step)
	})

	t.Run("InitStrategyWithInvalidAgrs", func(t *testing.T) {
		exp := NewExponential(0, -1*time.Second)

		assert.Equal(t, DefaultMinBackoff, exp.baseDelay)
		assert.Equal(t, DefaultMaxBackoff, exp.maxDelay)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		maxDelay := 1 * time.Second
		exp := NewExponential(baseDelay, maxDelay)

		assert.Equal(t, baseDelay, exp.maxDelay)
	})

	t.Run("SuccessGetNext", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exp := NewExponentialWithStep(baseDelay, maxDelay, 2)
		delay := exp.Next(2)
		assert.Equal(t, 4*time.Second, delay)
	})

	t.Run("SuccessGetNextWithGetMaxDelay", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exp := NewExponentialWithStep(baseDelay, maxDelay, 2)
		delay := exp.Next(4)
		assert.Equal(t, maxDelay, delay)
	})
}
