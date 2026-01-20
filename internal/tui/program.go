package tui

import (
	"bufio"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/brainwhocodes/ralph-codex/internal/codex"
	"github.com/brainwhocodes/ralph-codex/internal/loop"
)

// Program wraps the Bubble Tea program
type Program struct {
	model      Model
	controller *loop.Controller
}

// NewProgram creates a new TUI program
func NewProgram(config codex.Config, controller *loop.Controller) *Program {
	planInfo := loadTasks()

	model := Model{
		state:         StateInitializing,
		status:        "Ready to start",
		loopNumber:    0,
		maxCalls:      config.MaxCalls,
		callsUsed:     0,
		circuitState:  "CLOSED",
		logs:          []string{},
		activeView:    "status",
		quitting:      false,
		err:           nil,
		startTime:     time.Now(),
		width:         80,
		height:        24,
		tasks:         planInfo.Tasks,
		planFile:      planInfo.Filename,
		activity:      "",
		controller:    controller,
		activeTaskIdx: -1,
	}

	return &Program{
		model:      model,
		controller: controller,
	}
}

// PlanFileInfo holds info about loaded plan file
type PlanFileInfo struct {
	Filename string
	Tasks    []Task
}

// loadTasks reads tasks from plan files in order of preference:
// 1. IMPLEMENTATION_PLAN.md (implementation mode)
// 2. REFACTOR_PLAN.md (refactor mode)
// 3. @fix_plan.md (fix mode)
func loadTasks() PlanFileInfo {
	// Try plan files in order of preference
	planFiles := []string{
		"IMPLEMENTATION_PLAN.md",
		"REFACTOR_PLAN.md",
		"@fix_plan.md",
	}

	for _, planFile := range planFiles {
		data, err := os.ReadFile(planFile)
		if err != nil {
			continue
		}

		tasks := parseTasksFromData(string(data))
		return PlanFileInfo{
			Filename: planFile,
			Tasks:    tasks,
		}
	}

	return PlanFileInfo{
		Filename: "",
		Tasks:    []Task{},
	}
}

// parseTasksFromData extracts checklist tasks from plan file content
func parseTasksFromData(data string) []Task {
	var tasks []Task
	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse checkbox items: - [ ] or - [x]
		if strings.HasPrefix(line, "- [") {
			completed := strings.HasPrefix(line, "- [x]") || strings.HasPrefix(line, "- [X]")

			// Extract task text (skip "- [ ] " or "- [x] ")
			text := ""
			if len(line) > 6 {
				text = strings.TrimSpace(line[6:])
			}

			if text != "" {
				tasks = append(tasks, Task{
					Text:      text,
					Completed: completed,
				})
			}
		}
	}

	return tasks
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
			program.Send(ControllerEventMsg{Event: event})
		})
	}

	_, err := program.Run()
	return err
}
