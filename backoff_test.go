package backoff

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExportAllStrategiesToCSV(t *testing.T) {
	baseDelay := 1 * time.Second
	maxDelay := 30 * time.Second

	attempts := 20
	samples := 100

	strategies := map[string]Backoff{
		"Fixed":              NewFixed(baseDelay),
		"Exponential":        NewExponential(baseDelay, maxDelay),
		"EqualJitter":        NewEqualJitter(baseDelay, maxDelay),
		"DecorrelatedJitter": NewDecorrelatedJitter(baseDelay, maxDelay),
	}

	file, err := os.Create("backoff_results.csv")
	assert.NoError(t, err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"Strategy", "Attempt", "DelaySeconds"})

	for name, strategy := range strategies {
		t.Run(name, func(t *testing.T) {
			for s := 0; s < samples; s++ {
				var currentInput int64 = 0

				for i := 0; i < attempts; i++ {
					delay := strategy.Next(currentInput)

					assert.GreaterOrEqual(t, int64(delay), int64(0), "Delay should not be negative")
					assert.LessOrEqual(t, int64(delay), int64(maxDelay), "Delay should not exceed maxDelay")

					err = writer.Write([]string{
						name,
						strconv.Itoa(i),
						fmt.Sprintf("%.4f", delay.Seconds()),
					})
					assert.NoError(t, err)

					if name == "DecorrelatedJitter" {
						currentInput = int64(delay)
					} else {
						currentInput = int64(i + 1)
					}
				}
			}
		})
	}
}

func TestBackoff_TimelineSimulation(t *testing.T) {
	baseDelay := 500 * time.Millisecond
	maxDelay := 10 * time.Second
	simulationDuration := 30 * time.Second
	clientsCount := 100

	strategies := map[string]func() Backoff{
		"Fixed":              func() Backoff { return NewFixed(baseDelay) },
		"Exponential":        func() Backoff { return NewExponential(baseDelay, maxDelay) },
		"EqualJitter":        func() Backoff { return NewEqualJitter(baseDelay, maxDelay) },
		"DecorrelatedJitter": func() Backoff { return NewDecorrelatedJitter(baseDelay, maxDelay) },
	}

	file, err := os.Create("api_calls_timeline.csv")
	assert.NoError(t, err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"Strategy", "TimestampMs"})

	for name, factory := range strategies {
		t.Run(name, func(t *testing.T) {
			for c := 0; c < clientsCount; c++ {
				strategy := factory()
				var currentTime time.Duration = 0
				var currentInput int64 = 0

				// Первый вызов происходит сразу (в 0 мс)
				writer.Write([]string{name, "0"})

				for {
					delay := strategy.Next(currentInput)
					currentTime += delay

					if currentTime > simulationDuration {
						break
					}

					// Записываем момент времени вызова в миллисекундах
					writer.Write([]string{
						name,
						strconv.FormatInt(currentTime.Milliseconds(), 10),
					})

					// Обновляем входные данные для следующего шага
					if name == "DecorrelatedJitter" {
						currentInput = int64(delay)
					} else {
						currentInput++
					}
				}
			}
		})
	}
}
