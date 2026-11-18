package backoff

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("DefaultStepWithoutOptions", func(t *testing.T) {
		assert.Equal(t, DefaultStep, applyOptions(nil).step)
	})

	t.Run("CustomStep", func(t *testing.T) {
		assert.Equal(t, 1.5, applyOptions([]Option{WithStep(1.5)}).step)
	})

	t.Run("StepDoesNotGrowFallsBackToDefault", func(t *testing.T) {
		for _, step := range []float64{1, 0.5, 0, -2, math.NaN()} {
			assert.Equal(t, DefaultStep, applyOptions([]Option{WithStep(step)}).step, "step %v", step)
		}
	})

	t.Run("LastOptionWins", func(t *testing.T) {
		assert.Equal(t, 3.0, applyOptions([]Option{WithStep(1.5), WithStep(3)}).step)
	})
}

func TestNormalizeBounds(t *testing.T) {
	t.Parallel()

	t.Run("ValidValuesKept", func(t *testing.T) {
		baseDelay, maxDelay := normalizeBounds(time.Second, 10*time.Second)
		assert.Equal(t, time.Second, baseDelay)
		assert.Equal(t, 10*time.Second, maxDelay)
	})

	t.Run("NonPositiveValuesGetDefaults", func(t *testing.T) {
		baseDelay, maxDelay := normalizeBounds(0, -10*time.Second)
		assert.Equal(t, DefaultMinBackoff, baseDelay)
		assert.Equal(t, DefaultMaxBackoff, maxDelay)
	})

	t.Run("MaxLiftedToBase", func(t *testing.T) {
		baseDelay, maxDelay := normalizeBounds(5*time.Second, time.Second)
		assert.Equal(t, 5*time.Second, baseDelay)
		assert.Equal(t, 5*time.Second, maxDelay)
	})
}

func TestExpLimit(t *testing.T) {
	t.Parallel()

	t.Run("GrowsFromBase", func(t *testing.T) {
		assert.Equal(t, time.Second, expLimit(time.Second, time.Minute, 2, 0))
		assert.Equal(t, 8*time.Second, expLimit(time.Second, time.Minute, 2, 3))
	})

	t.Run("NegativeAttemptCountsAsZero", func(t *testing.T) {
		assert.Equal(t, time.Second, expLimit(time.Second, time.Minute, 2, -5))
	})

	t.Run("CappedAtMax", func(t *testing.T) {
		assert.Equal(t, time.Minute, expLimit(time.Second, time.Minute, 2, 10))
	})

	t.Run("HugeAttemptDoesNotOverflow", func(t *testing.T) {
		maxDelay := time.Duration(math.MaxInt64)

		for _, attempt := range []int{62, 63, 64, 1000, math.MaxInt} {
			assert.Equal(t, maxDelay, expLimit(time.Second, maxDelay, 2, attempt), "attempt %d", attempt)
		}
	})
}
