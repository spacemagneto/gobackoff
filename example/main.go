package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	backoff "github.com/spacemagneto/gobackoff"
)

func main() {
	// Global timeout for the entire operation (including all retries)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use httpbin.org to simulate a 503 Service Unavailable error
	url := "https://httpbin.org/status/503"

	err := fetchDataWithRetry(ctx, url)
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)
	}
}

func fetchDataWithRetry(ctx context.Context, url string) error {
	// Initialize Full Jitter: 1s base, 15s max cap
	b := backoff.NewFullJitter(1*time.Second, 15*time.Second)
	client := &http.Client{}
	maxAttempts := 5

	for attempt := 0; attempt < maxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}

		fmt.Printf("[Attempt %d] Requesting %s...\n", attempt, url)
		resp, err := client.Do(req)

		// Check if we should retry:
		// 1. Network/DNS error (err != nil)
		// 2. Server-side error (5xx)
		// 3. Rate limiting (429)
		if err != nil || (resp != nil && (resp.StatusCode >= 500 || resp.StatusCode == 429)) {
			if resp != nil {
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
			}

			if attempt == maxAttempts-1 {
				return fmt.Errorf("exhausted retries: last error was %v", err)
			}

			delay := b.Next(int64(attempt))
			fmt.Printf("Retrying in %v due to error/status...\n", delay)

			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Success path
		if resp != nil {
			_ = resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				fmt.Println("Successfully received 200 OK")
				return nil
			}

			return fmt.Errorf("received non-retryable status: %d", resp.StatusCode)
		}
	}

	return fmt.Errorf("retry loop exhausted")
}
