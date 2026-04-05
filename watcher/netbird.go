package watcher

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	upTimeout = 30 * time.Second
)

// NetBird represents the NetBird CLI interface.
type NetBird struct {
	setupKeyFile string
}

// NewNetBird creates a new NetBird CLI wrapper.
func NewNetBird(setupKeyFile string) *NetBird {
	return &NetBird{setupKeyFile: setupKeyFile}
}

// IsConnected checks if NetBird is currently connected.
// It executes `netbird status` and parses the output.
func (n *NetBird) IsConnected(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "netbird", "status")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Command failed - check if it's a connection error
		errMsg := stderr.String()
		if strings.Contains(errMsg, "connection") || strings.Contains(errMsg, "network") {
			return false, nil // Disconnected due to network issues
		}
		return false, fmt.Errorf("netbird status failed: %w", err)
	}

	output := stdout.String()
	// Parse output - look for indicators of connected state
	// NetBird status typically shows "Connected" or shows peer info when connected
	return parseConnectedState(output), nil
}

// upArgs returns the arguments for the netbird up command.
func (n *NetBird) upArgs() []string {
	args := []string{"up"}
	if n.setupKeyFile != "" {
		args = append(args, "--setup-key-file="+n.setupKeyFile)
	}
	return args
}

// Up attempts to bring up the NetBird connection.
func (n *NetBird) Up(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "netbird", n.upArgs()...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("netbird up failed: %w (stderr: %s)", err, stderr.String())
	}

	return nil
}

// parseConnectedState checks if the status output indicates a connected state.
// This is a simple heuristic - NetBird outputs vary, but typically includes
// connection status in the output.
func parseConnectedState(output string) bool {
	output = strings.ToLower(output)
	// Check for "disconnected" first - if present, we're not connected
	if strings.Contains(output, "disconnected") {
		return false
	}
	// Check for common "connected" indicators
	indicators := []string{
		"connection status: connected",
		"peers:", // Peers section appears when connected
	}
	for _, indicator := range indicators {
		if strings.Contains(output, indicator) {
			return true
		}
	}
	return false
}
