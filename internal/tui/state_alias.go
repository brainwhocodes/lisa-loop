package tui

import "github.com/brainwhocodes/lisa-loop/internal/tui/msg"

// State is re-exported for backwards-compat within the tui package.
// Internally, we centralize the enum in internal/tui/msg so other packages can
// reference it without importing the root tui package.
type State = msg.State

const (
	StateInitializing = msg.StateInitializing
	StateRunning      = msg.StateRunning
	StatePaused       = msg.StatePaused
	StateComplete     = msg.StateComplete
	StateError        = msg.StateError
)
