package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twistedogic/nestbird/watcher"
)

const (
	baseBackoff = 1 * time.Minute
	maxBackoff  = 5 * time.Minute
)

func main() {
	setupKeyFile := flag.String("setup-key-file", "", "path to NetBird setup key file (overrides NETBIRD_SETUP_KEY_FILE)")
	flag.Parse()

	if *setupKeyFile == "" {
		*setupKeyFile = os.Getenv("NETBIRD_SETUP_KEY_FILE")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal, shutting down", "signal", sig)
		cancel()
	}()

	w := watcher.NewWatcher(logger, baseBackoff, maxBackoff, *setupKeyFile)
	w.Run(ctx)

	logger.Info("nestbird exited")
}
