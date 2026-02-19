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
		assert.Equal(t, 1*time.Second, fixed.delay)
	})

	t.Run("InitStrategyWithDelay", func(t *testing.T) {
		expectDelay := 500 * time.Millisecond
		fixed := NewFixed(expectDelay)
		assert.NotNil(t, fixed)
		assert.Equal(t, expectDelay, fixed.delay)
	})

	t.Run("SuccessGetNext", func(t *testing.T) {
		expectDelay := 500 * time.Millisecond
		fixed := NewFixed(expectDelay)
		assert.NotNil(t, fixed)
		assert.Equal(t, expectDelay, fixed.delay)

		delay := fixed.Next(97911111)
		assert.Equal(t, expectDelay, delay)
	})
}
