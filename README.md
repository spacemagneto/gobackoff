[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# backoff

A Go package providing composable retry delay strategies based on the recommendations from the [AWS Architecture Blog](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/).

---

## Overview

When a service call fails, naive retry logic that immediately retries or retries at fixed intervals can make things worse under load. This package provides several well-known backoff strategies that help distribute retry pressure across clients, reducing contention and improving system resilience.

---

## Installation

```bash
go get github.com/spacemagneto/gobackoff
```

---

## Defaults

| Constant            | Value    | Description                                              |
|---------------------|----------|----------------------------------------------------------|
| `DefaultMinBackoff` | `1s`     | Starting delay if strategy is initialized with `<= 0`   |
| `DefaultMaxBackoff` | `20s`    | Upper bound — no strategy returns a value exceeding this |
| `DefaultStep`       | `2.0`    | Base factor for exponential growth                       |

---

## Interface

```go
type Backoff interface {
    Next(arg int64) time.Duration
}
```

All strategies implement the `Backoff` interface. The meaning of the `arg` parameter depends on the strategy:

| Strategy             | `arg` meaning                                                 |
|----------------------|---------------------------------------------------------------|
| `Exponential`        | Retry attempt number (1, 2, 3...)                             |
| `ExponentialJitter`  | Retry attempt number (1, 2, 3...)                             |
| `FullJitter`         | Retry attempt number (1, 2, 3...)                             |
| `EqualJitter`        | Retry attempt number (1, 2, 3...)                             |
| `DecorrelatedJitter` | Previous delay in nanoseconds (`time.Duration.Nanoseconds()`) |
| `Fixed`              | Ignored                                                       |

---

## Strategies

### Fixed

Always returns the same delay regardless of attempts or state. Useful for simple polling or scenarios where growth is not desired.

**Formula:** `sleep = delay`

```go
b := backoff.NewFixed(2 * time.Second)

b.Next(1) // 2s
b.Next(2) // 2s
b.Next(3) // 2s
```

---

### Exponential

A deterministic exponential backoff without any randomness. The delay grows by a fixed multiplier on every attempt. Predictable, but susceptible to thundering herd under high concurrency — consider a jitter strategy instead.

**Formula:** `sleep = min(cap, base * step^attempt)`

```go
b := backoff.NewExponential(1*time.Second, 20*time.Second)
// or with a custom step:
b := backoff.NewExponentialWithStep(1*time.Second, 20*time.Second, 3.0)

b.Next(1) // 2s
b.Next(2) // 4s
b.Next(3) // 8s
```

---

### ExponentialJitter

Exponential backoff with full jitter. Combines the gradual back-pressure of exponential growth with the desynchronization benefits of randomness. Unlike `FullJitter` (which always uses base `2` as the multiplier), this strategy accepts a configurable `step` — making it a flexible default for most retry scenarios.

**Formula:** `sleep = random_between(0, min(cap, base * step^attempt))`

```go
b := backoff.NewExponentialJitter(1*time.Second, 20*time.Second)
// or with a custom growth factor:
b := backoff.NewExponentialJitterWithStep(1*time.Second, 20*time.Second, 1.5)

b.Next(1) // random in [0s, 2s]
b.Next(2) // random in [0s, 4s]
b.Next(3) // random in [0s, 8s]
```

---

### FullJitter

Exponential backoff with a fully randomized delay. For each attempt the maximum possible delay is computed, and a random value in `[0, max]` is returned. Provides the best protection against thundering herd by maximizing desynchronization between clients.

**Formula:** `sleep = random_between(0, min(cap, base * 2^attempt))`

```go
b := backoff.NewFullJitter(1*time.Second, 20*time.Second)

b.Next(1) // random in [0s, 2s]
b.Next(2) // random in [0s, 4s]
b.Next(3) // random in [0s, 8s]
```

---

### EqualJitter

Exponential backoff where the delay is always at least half of the current exponential limit. Provides a more predictable lower bound than Full Jitter while still preventing client synchronization.

**Formula:**
```
temp  = min(cap, base * 2^attempt)
sleep = temp/2 + random_between(0, temp/2)
```

```go
b := backoff.NewEqualJitter(1*time.Second, 20*time.Second)

b.Next(1) // random in [1s, 2s]
b.Next(2) // random in [2s, 4s]
b.Next(3) // random in [4s, 8s]
```

---

### DecorrelatedJitter

A state-aware strategy where each delay is derived from the **previous delay** rather than the attempt number. This produces a "random walk" behavior with excellent desynchronization properties. AWS recommends this alongside Full Jitter as one of the most effective strategies.

**Formula:** `sleep = min(cap, random_between(base, previous_sleep * 3))`

> **Note:** Because this strategy is state-dependent, pass `0` on the first call and the result of the previous call's `.Nanoseconds()` on subsequent calls.

```go
b := backoff.NewDecorrelatedJitter(1*time.Second, 20*time.Second)

var prev int64
for i := 0; i < 5; i++ {
    d := b.Next(prev)
    prev = d.Nanoseconds()
    fmt.Println(d)
}
```

---

## Strategy Comparison

| Strategy             | Jitter | State-aware | Lower bound guaranteed | Best for                                  |
|----------------------|--------|-------------|------------------------|-------------------------------------------|
| `Fixed`              | ✗      | ✗           | ✓                      | Simple polling                            |
| `Exponential`        | ✗      | ✗           | ✓                      | Baseline / non-concurrent scenarios       |
| `ExponentialJitter`  | ✓      | ✗           | ✗                      | General purpose — configurable growth     |
| `FullJitter`         | ✓      | ✗           | ✗                      | High-concurrency, thundering herd         |
| `EqualJitter`        | ✓      | ✗           | ✓                      | Balance between predictability and spread |
| `DecorrelatedJitter` | ✓      | ✓           | ✓                      | Large-scale distributed systems           |

---

## Usage Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/spacemagneto/gobackoff"
)

func main() {
	b := backoff.NewFullJitter(500*time.Millisecond, 30*time.Second)

	for attempt := int64(1); attempt <= 5; attempt++ {
		delay := b.Next(attempt)
		fmt.Printf("Attempt %d: waiting %s\n", attempt, delay)
		time.Sleep(delay)

		// ... perform operation
	}
}
```

---

## References

- [Exponential Backoff and Jitter — AWS Architecture Blog](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)