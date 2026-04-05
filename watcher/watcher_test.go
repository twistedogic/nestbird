package watcher

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

var discardLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError + 10}))

// mockNetBird is a controllable NetBirdClient for testing.
type mockNetBird struct {
	connectedResults []bool
	connIdx          atomic.Int32
	upErrors         []error
	upIdx            atomic.Int32
	upCalls          atomic.Int32
}

func (m *mockNetBird) IsConnected(_ context.Context) (bool, error) {
	idx := int(m.connIdx.Add(1)) - 1
	if idx >= len(m.connectedResults) {
		return m.connectedResults[len(m.connectedResults)-1], nil
	}
	return m.connectedResults[idx], nil
}

func (m *mockNetBird) Up(_ context.Context) error {
	m.upCalls.Add(1)
	idx := int(m.upIdx.Add(1)) - 1
	if idx >= len(m.upErrors) {
		return m.upErrors[len(m.upErrors)-1]
	}
	return m.upErrors[idx]
}

func newTestWatcher(mock *mockNetBird) *Watcher {
	w := NewWatcherWithClient(discardLogger, mock, time.Millisecond, 10*time.Millisecond)
	w.pollInterval = time.Millisecond
	return w
}

// TestWatcherDetectsDisconnection verifies that the watcher detects a disconnection
// and calls Up to reconnect (task 5.1).
func TestWatcherDetectsDisconnection(t *testing.T) {
	mock := &mockNetBird{
		connectedResults: []bool{false, true},
		upErrors:         []error{nil},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	newTestWatcher(mock).Run(ctx)

	if mock.upCalls.Load() == 0 {
		t.Error("expected Up to be called when disconnected, but it was not")
	}
}

// TestWatcherReconnectsAutomatically verifies that after a failed Up, the watcher
// retries and eventually reconnects (task 5.2).
func TestWatcherReconnectsAutomatically(t *testing.T) {
	errFail := errors.New("netbird up failed")

	mock := &mockNetBird{
		connectedResults: []bool{false, false, true},
		upErrors:         []error{errFail, nil},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	newTestWatcher(mock).Run(ctx)

	calls := mock.upCalls.Load()
	if calls < 2 {
		t.Errorf("expected at least 2 Up calls (1 fail + 1 success), got %d", calls)
	}
}

// TestWatcherBackoffUnderRepeatedFailures verifies that the backoff attempt counter
// increments with each consecutive failure (task 5.3).
func TestWatcherBackoffUnderRepeatedFailures(t *testing.T) {
	errFail := errors.New("netbird up failed")

	const failCount = 4
	upErrors := make([]error, failCount)
	for i := range upErrors {
		upErrors[i] = errFail
	}
	upErrors = append(upErrors, nil)

	mock := &mockNetBird{
		connectedResults: []bool{false},
		upErrors:         upErrors,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	newTestWatcher(mock).Run(ctx)

	calls := mock.upCalls.Load()
	if int(calls) < failCount+1 {
		t.Errorf("expected %d Up calls, got %d", failCount+1, calls)
	}
}

// TestWatcherGracefulShutdown verifies that cancelling the context stops the watcher
// cleanly without panics or hangs (task 5.4).
func TestWatcherGracefulShutdown(t *testing.T) {
	mock := &mockNetBird{
		connectedResults: []bool{true},
		upErrors:         []error{nil},
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)
		newTestWatcher(mock).Run(ctx)
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("watcher did not stop within 2 seconds after context cancellation")
	}
}
