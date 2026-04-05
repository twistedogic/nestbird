package watcher

import (
	"reflect"
	"testing"
)

func TestNetBirdUpArgs(t *testing.T) {
	tests := []struct {
		name         string
		setupKeyFile string
		want         []string
	}{
		{
			name:         "no setup key file",
			setupKeyFile: "",
			want:         []string{"up"},
		},
		{
			name:         "with setup key file",
			setupKeyFile: "/etc/netbird/setup.key",
			want:         []string{"up", "--setup-key-file=/etc/netbird/setup.key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNetBird(tt.setupKeyFile)
			got := n.upArgs()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("upArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConnectedState(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		expected bool
	}{
		{
			name:     "Connected status",
			output:   "Connection status: Connected\n",
			expected: true,
		},
		{
			name:     "Connected (lowercase)",
			output:   "connection status: connected\n",
			expected: true,
		},
		{
			name:     "Peers section indicates connected",
			output:   "Peers: \n  - peer1: 10.0.0.1\n",
			expected: true,
		},
		{
			name:     "Disconnected status",
			output:   "Connection status: Disconnected\n",
			expected: false,
		},
		{
			name:     "Disconnected (lowercase)",
			output:   "disconnected\n",
			expected: false,
		},
		{
			name:     "Empty output",
			output:   "",
			expected: false,
		},
		{
			name:     "Unknown status",
			output:   "Some unknown status\n",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseConnectedState(tt.output)
			if result != tt.expected {
				t.Errorf("parseConnectedState(%q) = %v, want %v", tt.output, result, tt.expected)
			}
		})
	}
}
