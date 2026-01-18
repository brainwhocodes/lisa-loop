package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/brainwhocodes/ralph-codex/internal/codex"
)

// Program wraps the Bubble Tea program
type Program struct {
	model Model
}

// NewProgram creates a new TUI program
func NewProgram(config codex.Config) *Program {
	model := Model{
		state:      StateInitializing,
		status:     "Ready to start",
		loopNumber: 0,
		maxCalls:   config.MaxCalls,
		callsUsed:  0,
		logs:       []string{},
		activeView: "status",
		quitting:   false,
		err:        nil,
	}

	return &Program{
		model: model,
	}
}

// Run starts the TUI program
func (p *Program) Run() error {
	program := tea.NewProgram(p.model)
	_, err := program.Run()
	return err
}
