package config

// Config holds unified configuration for Lisa Codex
type Config struct {
	Backend      string
	ProjectPath  string
	PromptPath   string
	MaxCalls     int
	Timeout      int
	Verbose      bool
	ResetCircuit bool

	// OpenCode backend configuration
	OpenCodeServerURL string // URL for OpenCode server (env: OPENCODE_SERVER_URL)
	OpenCodeUsername  string // Username for OpenCode auth (env: OPENCODE_SERVER_USERNAME)
	OpenCodePassword  string // Password for OpenCode auth (env: OPENCODE_SERVER_PASSWORD)
	OpenCodeModelID   string // Model ID to use (env: OPENCODE_MODEL_ID, default: glm-4.7)
}

// BackendDisplayName returns a display-friendly name for the backend
func (c *Config) BackendDisplayName() string {
	switch c.Backend {
	case "opencode":
		return "OpenCode Server"
	case "cli":
		return "Codex CLI"
	default:
		if c.Backend != "" {
			return c.Backend
		}
		return "OpenCode Server"
	}
}
