package plan

import (
	"bufio"
	"strings"
)

// Task represents a checklist item parsed from a plan file.
// Active is owned by the UI; plan parsing leaves it false.
type Task struct {
	Text      string
	Completed bool
	Active    bool
}

// Phase groups tasks under a section header (e.g. "## Phase 1: ...").
type Phase struct {
	Name      string
	Tasks     []Task
	Completed bool
}

// ParseTasks extracts a flat task list from plan content (backwards compatibility).
func ParseTasks(data string) []Task {
	phases := ParsePhases(data)
	var tasks []Task
	for _, phase := range phases {
		tasks = append(tasks, phase.Tasks...)
	}
	return tasks
}

// ParsePhases extracts tasks grouped by phase from plan file content.
// Supports multiple plan formats:
// - REFACTOR_PLAN.md: ## Phase N: ... headers
// - IMPLEMENTATION_PLAN.md: ## Phase N: ... or ### N) atomic commit headers
// - @fix_plan.md: ## Critical Fixes, ## High Priority, ## Medium Priority, etc.
func ParsePhases(data string) []Phase {
	var phases []Phase
	var currentPhase *Phase

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Detect phase/section headers.
		if isPhaseHeader(trimmed) {
			header := extractPhaseHeader(trimmed)
			if currentPhase != nil && len(currentPhase.Tasks) > 0 {
				phases = append(phases, *currentPhase)
			}
			currentPhase = &Phase{Name: header, Tasks: []Task{}}
			continue
		}

		// Parse checkbox items: - [ ] or - [x]
		if strings.HasPrefix(trimmed, "- [") {
			completed := strings.HasPrefix(trimmed, "- [x]") || strings.HasPrefix(trimmed, "- [X]")

			// Extract task text (skip "- [ ] " or "- [x] ")
			text := ""
			if len(trimmed) > 6 {
				text = strings.TrimSpace(trimmed[6:])
			}
			if text == "" {
				continue
			}

			task := Task{Text: text, Completed: completed}
			if currentPhase != nil {
				currentPhase.Tasks = append(currentPhase.Tasks, task)
			} else {
				// No phase yet, create a default one.
				currentPhase = &Phase{Name: "Tasks", Tasks: []Task{task}}
			}
		}
	}

	// Don't forget the last phase.
	if currentPhase != nil && len(currentPhase.Tasks) > 0 {
		phases = append(phases, *currentPhase)
	}

	// Update phase completion status.
	for i := range phases {
		allComplete := true
		for _, task := range phases[i].Tasks {
			if !task.Completed {
				allComplete = false
				break
			}
		}
		phases[i].Completed = allComplete
	}

	return phases
}

func isPhaseHeader(line string) bool {
	lower := strings.ToLower(line)

	// ## Phase N: ... (REFACTOR_PLAN.md, IMPLEMENTATION_PLAN.md)
	if strings.HasPrefix(line, "## ") {
		header := strings.TrimPrefix(line, "## ")
		headerLower := strings.ToLower(header)

		// Phase headers.
		if strings.Contains(headerLower, "phase") {
			return true
		}

		// Fix plan priority headers.
		if strings.HasPrefix(headerLower, "critical") ||
			strings.HasPrefix(headerLower, "high priority") ||
			strings.HasPrefix(headerLower, "medium priority") ||
			strings.HasPrefix(headerLower, "low priority") ||
			strings.HasPrefix(headerLower, "testing") ||
			strings.HasPrefix(headerLower, "nice to have") {
			return true
		}

		// Verification/Success criteria sections.
		if strings.Contains(headerLower, "verification") ||
			strings.Contains(headerLower, "success criteria") {
			return true
		}
	}

	// ### N) Atomic commit headers (IMPLEMENTATION_PLAN.md)
	if strings.HasPrefix(line, "### ") {
		header := strings.TrimPrefix(line, "### ")
		// Check for numbered headers like "1) Config..." or "2) OpenCode..."
		if len(header) >= 2 && header[0] >= '1' && header[0] <= '9' && header[1] == ')' {
			return true
		}
	}

	// ## Atomic Commits section header.
	if strings.HasPrefix(lower, "## atomic") {
		return true
	}

	return false
}

func extractPhaseHeader(line string) string {
	if strings.HasPrefix(line, "### ") {
		header := strings.TrimPrefix(line, "### ")
		// For "1) Config..." make it "Step 1: Config..."
		if len(header) >= 2 && header[0] >= '1' && header[0] <= '9' && header[1] == ')' {
			return "Step " + string(header[0]) + ":" + header[2:]
		}
		return header
	}

	if strings.HasPrefix(line, "## ") {
		return strings.TrimPrefix(line, "## ")
	}

	return line
}
