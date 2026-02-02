package msg

import (
	"time"

	"github.com/brainwhocodes/lisa-loop/internal/loop"
)

// LoopUpdateMsg is sent when loop controller updates.
type LoopUpdateMsg struct {
	LoopNumber int
	CallsUsed  int
	Status     string
}

// LogMsg is sent to add a log entry.
type LogMsg struct {
	Message string
	Level   string // INFO, WARN, ERROR, SUCCESS
}

// StateChangeMsg is sent to change TUI state.
type StateChangeMsg struct {
	State State
}

// StatusMsg is sent to update status text.
type StatusMsg struct {
	Status string
}

// TickMsg is sent periodically for animations.
type TickMsg time.Time

// ControllerEventMsg wraps events from the loop controller.
type ControllerEventMsg struct {
	Event loop.LoopEvent
}

// CodexOutputMsg represents a line of output from a backend stream.
type CodexOutputMsg struct {
	Line string
	Type string // "reasoning", "agent_message", "tool_call", "raw"
}

// CodexReasoningMsg represents reasoning/thinking output.
type CodexReasoningMsg struct {
	Text string
}

// CodexToolCallMsg represents a tool call event.
type CodexToolCallMsg struct {
	Tool   string // "read", "write", "exec", etc.
	Target string // file path or command
	Status string // "started", "completed"
}

// TaskStartedMsg indicates a task has started.
type TaskStartedMsg struct {
	TaskIndex int
	TaskText  string
}

// TaskCompletedMsg indicates a task has been completed.
type TaskCompletedMsg struct {
	TaskIndex int
	TaskText  string
}

// TaskFailedMsg indicates a task has failed.
type TaskFailedMsg struct {
	TaskIndex int
	TaskText  string
	Error     string
}

// ViewModeMsg changes the current view mode (legacy string-based mode).
type ViewModeMsg struct {
	Mode string // "split", "tasks", "output", "logs"
}

// PreflightMsg carries preflight check summary from loop controller.
type PreflightMsg struct {
	Mode           string
	PlanFile       string
	TotalTasks     int
	RemainingCount int
	RemainingTasks []string
	CircuitState   string
	RateLimitOK    bool
	CallsRemaining int
	ShouldSkip     bool
	SkipReason     string
}

// LoopOutcomeMsg carries loop iteration outcome from loop controller.
type LoopOutcomeMsg struct {
	Success        bool
	TasksCompleted int
	FilesModified  int
	TestsStatus    string
	ExitSignal     bool
	Error          string
}
