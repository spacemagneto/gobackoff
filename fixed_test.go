package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedBackoffStrategy(t *testing.T) {
	t.Parallel()

	t.Run("InitStrategyWithoutDelay", func(t *testing.T) {
		fixed := NewFixed(0)
		assert.NotNil(t, fixed)
		assert.Equal(t, DefaultMinBackoff, fixed.delay)
	})

	t.Run("InitStrategyWithDelay", func(t *testing.T) {
		expectDelay := 500 * time.Millisecond
		fixed := NewFixed(expectDelay)
		assert.NotNil(t, fixed)
		assert.Equal(t, expectDelay, fixed.delay)
	})

	t.Run("NextAlwaysSame", func(t *testing.T) {
		expectDelay := 500 * time.Millisecond
		fixed := NewFixed(expectDelay)

		for i := 0; i < 5; i++ {
			assert.Equal(t, expectDelay, fixed.Next())
		}

		fixed.Reset()
		assert.Equal(t, expectDelay, fixed.Next())
	})
}
