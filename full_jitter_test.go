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
}
