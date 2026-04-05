package watcher

import (
	"math/rand"
	"time"
)

// Backoff implements exponential backoff with full jitter.
// Formula: wait = min(base*2^attempt, max) + random(0, wait/2)
type Backoff struct {
	base    time.Duration
	max     time.Duration
	attempt int
}

func NewBackoff(base, max time.Duration) *Backoff {
	return &Backoff{
		base:    base,
		max:     max,
		attempt: 0,
	}
}

// Duration returns the next wait duration and advances the attempt counter.
// Formula: wait = min(base * 2^attempt + random(0, base*2^attempt/2), max)
// This gives us ±50% jitter while never exceeding max.
func (b *Backoff) Duration() time.Duration {
	// Calculate base wait time with exponential backoff
	// Cap early to avoid overflow
	wait := b.base
	for i := 0; i < b.attempt; i++ {
		if wait >= b.max/2 {
			wait = b.max
			break
		}
		wait *= 2
	}

	// Cap at max before calculating jitter
	if wait > b.max {
		wait = b.max
	}

	// Add full jitter: random(0, wait/2)
	jitterMax := int64(wait / 2)
	if jitterMax < 1 {
		jitterMax = 1 // Ensure at least 1 nanosecond of jitter range
	}
	jitter := time.Duration(rand.Int63n(jitterMax))

	// Final cap (in case jitter pushes us over)
	result := wait + jitter
	if result > b.max {
		result = b.max
	}

	b.attempt++
	return result
}

// Reset restarts the backoff to initial state.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current attempt number (0-indexed).
func (b *Backoff) Attempt() int {
	return b.attempt
}
