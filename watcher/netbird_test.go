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
		name        string
		input       []byte
		expected    bool
		expectError bool
	}{
		{
			name:     "connected with netbirdIp",
			input:    []byte(`{"netbirdIp":"100.64.0.1"}`),
			expected: true,
		},
		{
			name:     "disconnected with empty netbirdIp",
			input:    []byte(`{"netbirdIp":""}`),
			expected: false,
		},
		{
			name:     "disconnected with missing netbirdIp field",
			input:    []byte(`{}`),
			expected: false,
		},
		{
			name:        "invalid JSON",
			input:       []byte(`not json`),
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseConnectedState(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("parseConnectedState() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("parseConnectedState() unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("parseConnectedState() = %v, want %v", result, tt.expected)
			}
		})
	}
}
