package watcher

import (
	"context"
	"log/slog"
	"time"
)

const (
	pollInterval = 5 * time.Minute
)

// NetBirdClient is the interface for interacting with the NetBird CLI.
type NetBirdClient interface {
	IsConnected(ctx context.Context) (bool, error)
	Up(ctx context.Context) error
}

// Watcher monitors NetBird connection status and reconnects when needed.
type Watcher struct {
	netbird      NetBirdClient
	backoff      *Backoff
	logger       *slog.Logger
	pollInterval time.Duration
}

// NewWatcher creates a new NetBird connection watcher.
func NewWatcher(logger *slog.Logger, baseBackoff, maxBackoff time.Duration, setupKeyFile string) *Watcher {
	return &Watcher{
		netbird:      NewNetBird(setupKeyFile),
		backoff:      NewBackoff(baseBackoff, maxBackoff),
		logger:       logger,
		pollInterval: pollInterval,
	}
}

// NewWatcherWithClient creates a new Watcher with a custom NetBirdClient (for testing).
func NewWatcherWithClient(logger *slog.Logger, client NetBirdClient, baseBackoff, maxBackoff time.Duration) *Watcher {
	return &Watcher{
		netbird:      client,
		backoff:      NewBackoff(baseBackoff, maxBackoff),
		logger:       logger,
		pollInterval: pollInterval,
	}
}

// Run starts the watcher loop. It blocks until the context is cancelled.
func (w *Watcher) Run(ctx context.Context) {
	w.logger.Info("Starting nestbird watcher", "pollInterval", w.pollInterval)

	for {
		connected, err := w.checkConnection(ctx)
		if err != nil {
			w.logger.Warn("Error checking connection", "error", err)
		}

		if connected {
			w.backoff.Reset()
			w.logger.Info("Connected")
			w.waitWithCancel(ctx, w.pollInterval)
		} else {
			w.handleDisconnected(ctx)
		}

		// Check if we've been signalled to stop
		select {
		case <-ctx.Done():
			w.logger.Info("Watcher stopped")
			return
		default:
		}
	}
}

// checkConnection verifies the current NetBird connection status.
func (w *Watcher) checkConnection(ctx context.Context) (bool, error) {
	return w.netbird.IsConnected(ctx)
}

// handleDisconnected attempts to reconnect with exponential backoff.
func (w *Watcher) handleDisconnected(ctx context.Context) {
	w.logger.Warn("Disconnected, attempting reconnect")

	for {
		// Try to reconnect
		err := w.netbird.Up(ctx)
		if err == nil {
			w.logger.Info("Reconnected")
			return
		}

		w.logger.Warn("Reconnection failed, retrying", "error", err)

		// Calculate backoff
		wait := w.backoff.Duration()
		w.logger.Info("Retrying connection", "wait", wait.Round(time.Second))

		// Wait with ability to be cancelled
		if !w.waitWithCancel(ctx, wait) {
			return // Context cancelled
		}
	}
}

// waitWithCancel waits for the specified duration but returns early if context is cancelled.
func (w *Watcher) waitWithCancel(ctx context.Context, duration time.Duration) bool {
	select {
	case <-ctx.Done():
		return false
	case <-time.After(duration):
		return true
	}
}
