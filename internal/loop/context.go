package loop

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadFixPlan loads remaining tasks from @fix_plan.md
func LoadFixPlan() ([]string, error) {
	data, err := os.ReadFile("@fix_plan.md")
	if err != nil {
		return nil, fmt.Errorf("failed to read @fix_plan.md: %w", err)
	}

	var tasks []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "- [") || strings.HasPrefix(line, "- [x") {
			taskText := strings.TrimSpace(line[4:])
			if strings.HasPrefix(taskText, "[") {
				taskText = strings.TrimSpace(taskText[3:])
			}
			if taskText != "" {
				tasks = append(tasks, taskText)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse @fix_plan.md: %w", err)
	}

	return tasks, nil
}

// GetPrompt loads the main prompt from PROMPT.md
func GetPrompt() (string, error) {
	data, err := os.ReadFile("PROMPT.md")
	if err != nil {
		return "", fmt.Errorf("failed to read PROMPT.md: %w", err)
	}

	return string(data), nil
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
		prompts := []string{"PROMPT.md", "@fix_plan.md", "@AGENT.md", ".git"}
		found := false

		for _, promptFile := range prompts {
			if _, err := os.Stat(filepath.Join(pwd, promptFile)); err == nil {
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

	return "", fmt.Errorf("could not find project root (no PROMPT.md, @fix_plan.md, @AGENT.md, or .git found)")
}

// CheckProjectRoot verifies we're in a valid Ralph project
func CheckProjectRoot() error {
	prompts := []string{"PROMPT.md", "@fix_plan.md"}
	found := 0

	for _, promptFile := range prompts {
		if _, err := os.Stat(promptFile); err == nil {
			found++
		}
	}

	if found < 2 {
		return fmt.Errorf("not a valid Ralph project (found %d/%d required files)", found, len(prompts))
	}

	return nil
}
