package config

import "testing"

func TestBackendDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		backend  string
		expected string
	}{
		{
			name:     "OpenCode backend",
			backend:  "opencode",
			expected: "OpenCode Server",
		},
		{
			name:     "Codex CLI backend",
			backend:  "cli",
			expected: "Codex CLI",
		},
		{
			name:     "Custom backend",
			backend:  "custom-agent",
			expected: "custom-agent",
		},
		{
			name:     "Empty backend defaults to OpenCode",
			backend:  "",
			expected: "OpenCode Server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Backend: tt.backend}
			if got := cfg.BackendDisplayName(); got != tt.expected {
				t.Errorf("BackendDisplayName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
