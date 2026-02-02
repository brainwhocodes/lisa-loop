package tui

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/brainwhocodes/lisa-loop/internal/codex"
	"github.com/brainwhocodes/lisa-loop/internal/loop"
	"github.com/brainwhocodes/lisa-loop/internal/tui/effects"
	"github.com/brainwhocodes/lisa-loop/internal/tui/markdown"
	"github.com/brainwhocodes/lisa-loop/internal/tui/msg"
	"github.com/brainwhocodes/lisa-loop/internal/tui/plan"
	"github.com/brainwhocodes/lisa-loop/internal/tui/transcript"
)

// Program wraps the Bubble Tea program
type Program struct {
	model      Model
	controller Controller
}

// NewProgram creates a new TUI program
// If explicitMode is provided (non-empty), use it instead of auto-detecting
func NewProgram(config codex.Config, controller *loop.Controller, explicitMode ...loop.ProjectMode) *Program {
	var projectMode loop.ProjectMode
	if len(explicitMode) > 0 && explicitMode[0] != "" {
		projectMode = explicitMode[0]
	} else {
		projectMode = loop.DetectProjectMode()
	}
	planInfo := loadTasksForMode(projectMode)

	// Determine initial state and status based on loaded files
	initialState := StateInitializing
	initialStatus := "Ready to start"
	var initialErr error
	var logs []string

	// Validate project mode
	if projectMode == loop.ModeUnknown {
		initialState = StateError
		initialStatus = "Invalid project - no mode detected"
		initialErr = fmt.Errorf("could not detect project mode")
		logs = append(logs, formatLog("ERROR", "No valid project mode detected. Need PRD.md+IMPLEMENTATION_PLAN.md, REFACTOR_PLAN.md, or PROMPT.md+@fix_plan.md"))
	} else {
		// Try to load prompt for the detected mode
		prompt, err := loop.GetPromptForMode(projectMode)
		if err != nil {
			initialState = StateError
			initialStatus = "Failed to load prompt"
			initialErr = err
			logs = append(logs, formatLog("ERROR", err.Error()))
		} else {
			logs = append(logs, formatLog("INFO", fmt.Sprintf("Loaded %s mode", projectMode)))
			logs = append(logs, formatLog("INFO", fmt.Sprintf("Prompt size: %d bytes", len(prompt))))
		}

		// Log plan file status
		if planInfo.Filename != "" {
			logs = append(logs, formatLog("INFO", fmt.Sprintf("Loaded %d tasks from %s", len(planInfo.Tasks), planInfo.Filename)))
		} else {
			logs = append(logs, formatLog("WARN", "No plan file found"))
		}
	}

	model := Model{
		state:          initialState,
		status:         initialStatus,
		loopNumber:     0,
		maxCalls:       config.MaxCalls,
		callsUsed:      0,
		circuitState:   "CLOSED",
		logs:           logs,
		screen:         ScreenSplit,
		quitting:       false,
		err:            initialErr,
		startTime:      time.Now(),
		width:          80,
		height:         24,
		tasks:          planInfo.Tasks,
		phases:         planInfo.Phases,
		currentPhase:   findFirstIncompletePhase(planInfo.Phases),
		planFile:       planInfo.Filename,
		projectMode:    projectMode,
		activity:       "",
		controller:     controller,
		readFile:       effects.OSReadFile,
		exec:           effects.OSExec,
		md:             markdown.New(),
		transcript:     transcript.New(500),
		outputTab:      OutputTabTranscript,
		activeTaskIdx:  -1,
		backend:        config.Backend,
		outputLines:    []string{},
		reasoningLines: []string{},
	}

	return &Program{
		model:      model,
		controller: controller,
	}
}

// formatLog formats a log entry with timestamp
func formatLog(level, message string) string {
	return fmt.Sprintf("[%s] %s: %s", time.Now().Format("15:04:05"), level, message)
}

// findFirstIncompletePhase returns the index of the first incomplete phase
func findFirstIncompletePhase(phases []Phase) int {
	for i, phase := range phases {
		if !phase.Completed {
			return i
		}
	}
	// All phases complete, return last one
	if len(phases) > 0 {
		return len(phases) - 1
	}
	return 0
}

// PlanFileInfo holds info about loaded plan file
type PlanFileInfo struct {
	Filename string
	Tasks    []Task
	Phases   []Phase
}

// loadTasksForMode reads tasks from the plan file for the given mode
func loadTasksForMode(mode loop.ProjectMode) PlanFileInfo {
	planFile := loop.GetPlanFileForMode(mode)
	if planFile == "" {
		return PlanFileInfo{
			Filename: "",
			Tasks:    []Task{},
			Phases:   []Phase{},
		}
	}

	data, err := os.ReadFile(planFile)
	if err != nil {
		return PlanFileInfo{
			Filename: "",
			Tasks:    []Task{},
			Phases:   []Phase{},
		}
	}

	phases := plan.ParsePhases(string(data))
	tasks := plan.ParseTasks(string(data))
	return PlanFileInfo{
		Filename: planFile,
		Tasks:    tasks,
		Phases:   phases,
	}
}

// Run starts the TUI program
func (p *Program) Run() error {
	program := tea.NewProgram(
		p.model,
		tea.WithAltScreen(),       // Full-screen alternate buffer mode
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Set up controller event callback to send messages to the TUI
	if p.controller != nil {
		p.controller.SetEventCallback(func(event loop.LoopEvent) {
			program.Send(msg.ControllerEventMsg{Event: event})
		})
	}

	_, err := program.Run()
	return err
}
