package codex

// Config holds Codex runner configuration
type Config struct {
	Backend      string
	ProjectPath  string
	PromptPath   string
	MaxCalls     int
	Timeout      int
	Verbose      bool
	ResetCircuit bool
}
