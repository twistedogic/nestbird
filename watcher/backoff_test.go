package watcher

import (
	"testing"
	"time"
)

func TestNewBackoff(t *testing.T) {
	b := NewBackoff(1*time.Minute, 5*time.Minute)
	if b.attempt != 0 {
		t.Errorf("expected attempt 0, got %d", b.attempt)
	}
}

func TestBackoffDuration_ExponentialGrowth(t *testing.T) {
	base := 1 * time.Minute
	maxBackoff := 5 * time.Minute

	// Test each attempt level directly
	tests := []struct {
		attempt         int
		expectedBase    time.Duration
		expectedMin     time.Duration
		expectedMax     time.Duration
	}{
		{0, 1 * time.Minute, 30 * time.Second, 90 * time.Second},   // 1m ± 50%
		{1, 2 * time.Minute, 60 * time.Second, 180 * time.Second},  // 2m ± 50%
		{2, 4 * time.Minute, 120 * time.Second, 360 * time.Second}, // 4m ± 50%
		{3, 5 * time.Minute, 150 * time.Second, 300 * time.Second}, // 5m capped (base) ± 50% capped at 5m
		{4, 5 * time.Minute, 150 * time.Second, 300 * time.Second}, // 5m capped (base) ± 50% capped at 5m
	}

	for _, tt := range tests {
		// Sample multiple times to account for jitter
		var minDur, maxDur time.Duration
		for j := 0; j < 100; j++ {
			// Create fresh backoff for each sample to avoid attempt incrementing
			b := NewBackoff(base, maxBackoff)
			b.attempt = tt.attempt
			d := b.Duration()
			if j == 0 || d < minDur {
				minDur = d
			}
			if j == 0 || d > maxDur {
				maxDur = d
			}
		}

		if minDur < tt.expectedMin {
			t.Errorf("attempt %d: min duration %v < expected min %v", tt.attempt, minDur, tt.expectedMin)
		}
		if maxDur > tt.expectedMax {
			t.Errorf("attempt %d: max duration %v > expected max %v", tt.attempt, maxDur, tt.expectedMax)
		}
	}
}

func TestBackoffReset(t *testing.T) {
	b := NewBackoff(1*time.Minute, 5*time.Minute)

	// Advance a few times
	b.Duration()
	b.Duration()
	b.Duration()

	if b.Attempt() != 3 {
		t.Errorf("expected attempt 3, got %d", b.Attempt())
	}

	b.Reset()

	if b.Attempt() != 0 {
		t.Errorf("expected attempt 0 after reset, got %d", b.Attempt())
	}
}

func TestBackoffJitterVariance(t *testing.T) {
	b := NewBackoff(1*time.Minute, 5*time.Minute)

	// Collect durations - they should vary due to jitter
	durations := make(map[time.Duration]bool)
	for i := 0; i < 100; i++ {
		b.Reset()
		d := b.Duration()
		durations[d] = true
	}

	// With jitter, we should see multiple distinct values
	// Without jitter, we'd always see the same value
	if len(durations) < 5 {
		t.Errorf("expected diverse durations due to jitter, got %d unique values", len(durations))
	}
}

func TestBackoffNeverExceedsMax(t *testing.T) {
	base := 1 * time.Minute
	maxBackoff := 5 * time.Minute
	b := NewBackoff(base, maxBackoff)

	// Generate many durations at high attempt numbers
	for i := 0; i < 100; i++ {
		b.attempt = 100 // High attempt number
		d := b.Duration()
		if d > maxBackoff {
			t.Errorf("duration %v exceeds max %v", d, maxBackoff)
		}
	}
}

func TestBackoffDuration_Boundaries(t *testing.T) {
	base := 1 * time.Minute
	maxBackoff := 5 * time.Minute
	b := NewBackoff(base, maxBackoff)

	// Test that durations stay within reasonable bounds
	tests := []struct {
		attempt     int
		minExpected time.Duration
		maxExpected time.Duration
	}{
		{0, 30 * time.Second, 90 * time.Second},   // 1m ± 50%
		{1, 60 * time.Second, 180 * time.Second},  // 2m ± 50%
		{2, 120 * time.Second, 360 * time.Second}, // 4m ± 50%
	}

	for _, tt := range tests {
		b.attempt = tt.attempt
		d := b.Duration()

		if d < tt.minExpected || d > tt.maxExpected {
			t.Errorf("attempt %d: duration %v not in range [%v, %v]",
				tt.attempt, d, tt.minExpected, tt.maxExpected)
		}
	}
}
