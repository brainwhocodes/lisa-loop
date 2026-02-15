package msg

import (
	"time"

	"github.com/brainwhocodes/lisa-loop/internal/loop"
	"github.com/brainwhocodes/lisa-loop/internal/tui/plan"
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

// AutoStartLoopMsg triggers the first controller run when monitor mode launches.
type AutoStartLoopMsg struct{}

// ControllerEventMsg wraps events from the loop controller.
type ControllerEventMsg struct {
	Event loop.LoopEvent
}

// ControllerDoneMsg is sent when the controller Run() command returns.
// The controller itself continues to emit fine-grained progress via ControllerEventMsg.
type ControllerDoneMsg struct {
	Err error
}

// PlanLoadedMsg is sent when a plan file is loaded (or fails to load).
type PlanLoadedMsg struct {
	Filename string
	Tasks    []plan.Task
	Phases   []plan.Phase
	Err      error
}

// DiffDebounceFiredMsg is emitted after a debounce interval and triggers a git diff refresh.
type DiffDebounceFiredMsg struct {
	Seq int
}

// GitDiffLoadedMsg is emitted when git diff output has been collected.
type GitDiffLoadedMsg struct {
	Seq        int
	NameStatus string
	Patch      string
	Err        error
	At         time.Time
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
