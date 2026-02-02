package msg

// State represents the high-level TUI lifecycle state.
// Kept small and explicit so other packages (effects/screens) can emit state changes
// without importing the root tui package.
type State int

const (
	StateInitializing State = iota
	StateRunning
	StatePaused
	StateComplete
	StateError
)

func (s State) String() string {
	switch s {
	case StateInitializing:
		return "Initializing"
	case StateRunning:
		return "Running"
	case StatePaused:
		return "Paused"
	case StateComplete:
		return "Complete"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}
