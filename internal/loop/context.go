package loop

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadFixPlan loads remaining tasks from plan files in order of preference:
// 1. IMPLEMENTATION_PLAN.md (implementation mode)
// 2. REFACTOR_PLAN.md (refactor mode)
// 3. @fix_plan.md (fix mode)
func LoadFixPlan() ([]string, error) {
	// Try plan files in order of preference
	planFiles := []string{
		"IMPLEMENTATION_PLAN.md",
		"REFACTOR_PLAN.md",
		"@fix_plan.md",
	}

	var data []byte
	var err error
	var planFile string

	for _, pf := range planFiles {
		data, err = os.ReadFile(pf)
		if err == nil {
			planFile = pf
			break
		}
	}

	if planFile == "" {
		return nil, fmt.Errorf("failed to read plan file (tried %s)", strings.Join(planFiles, ", "))
	}

	return parseTasksFromPlan(string(data), planFile)
}

// parseTasksFromPlan extracts checklist tasks from a plan file
func parseTasksFromPlan(content, filename string) ([]string, error) {
	var tasks []string
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Match checklist items: - [ ] or - [x]
		if strings.HasPrefix(line, "- [") {
			// Extract the checkbox state and text
			// Format: "- [ ] task" or "- [x] task"
			if len(line) >= 6 {
				checkbox := line[2:5] // "[ ]" or "[x]"
				taskText := strings.TrimSpace(line[5:])

				// Preserve the checkbox state in the task
				if checkbox == "[x]" || checkbox == "[X]" {
					tasks = append(tasks, "[x] "+taskText)
				} else {
					tasks = append(tasks, "[ ] "+taskText)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	return tasks, nil
}

// ProjectMode represents the type of Ralph project
type ProjectMode string

const (
	ModeImplement ProjectMode = "implement"
	ModeRefactor  ProjectMode = "refactor"
	ModeFix       ProjectMode = "fix"
	ModeUnknown   ProjectMode = "unknown"
)

// DetectProjectMode determines the project mode based on files present
func DetectProjectMode() ProjectMode {
	// Implement mode: PRD.md + IMPLEMENTATION_PLAN.md
	if fileExists("PRD.md") && fileExists("IMPLEMENTATION_PLAN.md") {
		return ModeImplement
	}
	// Refactor mode: REFACTOR_PLAN.md (no input file required)
	if fileExists("REFACTOR_PLAN.md") {
		return ModeRefactor
	}
	// Fix mode: PROMPT.md + @fix_plan.md
	if fileExists("PROMPT.md") && fileExists("@fix_plan.md") {
		return ModeFix
	}
	return ModeUnknown
}

// GetPlanFileForMode returns the plan file path for a given mode
func GetPlanFileForMode(mode ProjectMode) string {
	switch mode {
	case ModeImplement:
		return "IMPLEMENTATION_PLAN.md"
	case ModeRefactor:
		return "REFACTOR_PLAN.md"
	case ModeFix:
		return "@fix_plan.md"
	default:
		return ""
	}
}

// GetPromptForMode loads the appropriate prompt based on project mode
func GetPromptForMode(mode ProjectMode) (string, error) {
	switch mode {
	case ModeImplement:
		// Use PRD.md as the prompt for implement mode
		data, err := os.ReadFile("PRD.md")
		if err != nil {
			return "", fmt.Errorf("failed to read PRD.md: %w", err)
		}
		return string(data), nil

	case ModeRefactor:
		// Use REFACTOR_PLAN.md as context for refactor mode
		data, err := os.ReadFile("REFACTOR_PLAN.md")
		if err != nil {
			return "", fmt.Errorf("failed to read REFACTOR_PLAN.md: %w", err)
		}
		return string(data), nil

	case ModeFix:
		// Use PROMPT.md for fix mode
		data, err := os.ReadFile("PROMPT.md")
		if err != nil {
			return "", fmt.Errorf("failed to read PROMPT.md: %w", err)
		}
		return string(data), nil

	default:
		return "", fmt.Errorf("unknown project mode")
	}
}

// GetPrompt loads the main prompt based on detected project mode
func GetPrompt() (string, error) {
	mode := DetectProjectMode()
	if mode == ModeUnknown {
		return "", fmt.Errorf("could not detect project mode - need PRD.md, REFACTOR_PLAN.md, or PROMPT.md")
	}
	return GetPromptForMode(mode)
}

// BuildContext builds loop context for Codex
func BuildContext(promptPath string, loopNum int, remainingTasks []string, circuitState string, prevSummary string) (string, error) {
	var ctxBuilder strings.Builder

	ctxBuilder.WriteString(fmt.Sprintf("\n--- RALPH LOOP CONTEXT ---\n"))
	ctxBuilder.WriteString(fmt.Sprintf("Loop: %d\n", loopNum))
	ctxBuilder.WriteString(fmt.Sprintf("Circuit Breaker: %s\n", circuitState))

	if len(remainingTasks) > 0 && len(remainingTasks) <= 5 {
		ctxBuilder.WriteString(fmt.Sprintf("\nRemaining Tasks:\n"))
		for i, task := range remainingTasks {
			ctxBuilder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, task))
		}
	}

	if prevSummary != "" {
		ctxBuilder.WriteString(fmt.Sprintf("\nPrevious Loop Summary:\n%s\n", prevSummary))
	}

	ctxBuilder.WriteString(fmt.Sprintf("--- END LOOP CONTEXT ---\n\n"))

	return ctxBuilder.String(), nil
}

// InjectContext prepends context to prompt
func InjectContext(prompt string, ctx string) string {
	return ctx + prompt
}

// GetProjectRoot returns the project root directory
func GetProjectRoot() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Look for Ralph project markers
		markers := []string{
			"IMPLEMENTATION_PLAN.md",
			"REFACTOR_PLAN.md",
			"PRD.md",
			"REFACTOR.md",
			"AGENTS.md",
			"@fix_plan.md",
			"PROMPT.md",
			".git",
		}
		found := false

		for _, marker := range markers {
			if _, err := os.Stat(filepath.Join(pwd, marker)); err == nil {
				found = true
				break
			}
		}

		if found {
			return pwd, nil
		}

		parent := filepath.Dir(pwd)
		if parent == pwd {
			break
		}

		pwd = parent
	}

	return "", fmt.Errorf("could not find project root (no IMPLEMENTATION_PLAN.md, REFACTOR_PLAN.md, PRD.md, or .git found)")
}

// CheckProjectRoot verifies we're in a valid Ralph project
func CheckProjectRoot() error {
	// A valid Ralph project needs one of:
	// 1. PRD.md + IMPLEMENTATION_PLAN.md (implementation mode)
	// 2. REFACTOR.md + REFACTOR_PLAN.md (refactor mode)
	// 3. PROMPT.md + @fix_plan.md (fix mode / legacy)

	implementationMode := fileExists("PRD.md") && fileExists("IMPLEMENTATION_PLAN.md")
	refactorMode := fileExists("REFACTOR.md") && fileExists("REFACTOR_PLAN.md")
	fixMode := fileExists("PROMPT.md") && fileExists("@fix_plan.md")

	if !implementationMode && !refactorMode && !fixMode {
		return fmt.Errorf("not a valid Ralph project. Need one of: (PRD.md + IMPLEMENTATION_PLAN.md), (REFACTOR.md + REFACTOR_PLAN.md), or (PROMPT.md + @fix_plan.md)")
	}

	return nil
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
