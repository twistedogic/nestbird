package watcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
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
// It executes `netbird status --json` and checks if the netbirdIp field is non-empty.
func (n *NetBird) IsConnected(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "netbird", "status", "--json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("netbird status failed: %w (stderr: %s)", err, stderr.String())
	}

	return parseConnectedState(stdout.Bytes())
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

// netbirdStatus is the minimal shape of `netbird status --json` output.
type netbirdStatus struct {
	NetbirdIP string `json:"netbirdIp"`
}

// parseConnectedState parses JSON from `netbird status --json` and returns true
// when the netbirdIp field is non-empty.
func parseConnectedState(data []byte) (bool, error) {
	var status netbirdStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return false, fmt.Errorf("failed to parse netbird status JSON: %w", err)
	}
	return status.NetbirdIP != "", nil
}
