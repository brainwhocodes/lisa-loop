package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// renderSplitView renders the main split pane layout (tasks top, output bottom)
func (m Model) renderSplitView() string {
	width := m.width
	height := m.height
	if width < 80 {
		width = 80
	}
	if height < 24 {
		height = 24
	}

	// Header height + status bar + footer
	headerHeight := 1
	statusHeight := 1
	footerHeight := 1
	contentHeight := height - headerHeight - statusHeight - footerHeight - 2

	// Split content vertically: 40% tasks, 60% output
	topHeight := (contentHeight * 40) / 100
	bottomHeight := contentHeight - topHeight - 1 // 1 for divider

	// Render header
	header := m.renderHeader(width)

	// Render status bar
	statusBar := m.renderStatusBar(width)

	// Render top pane (tasks)
	topPane := m.renderTaskPane(width, topHeight)

	// Render horizontal divider
	divider := m.renderHorizontalDivider(width)

	// Render bottom pane (output)
	bottomPane := m.renderOutputPane(width, bottomHeight)

	// Render footer
	footer := m.renderFooter(width)

	// Join everything vertically
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		statusBar,
		topPane,
		divider,
		bottomPane,
		footer,
	)
}

// renderHeader renders the Crush-style header with gradient text and diagonal separators
// Format: Charm RALPH ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱ mode • loop 2 • 3/10
func (m Model) renderHeader(width int) string {
	// Brand prefix and name
	brandPrefix := StyleBrandPrefix.Render("Charm")
	brandName := GradientText("RALPH", Dolly, Charple)

	// Build metadata on right side
	var metaParts []string

	// Mode
	modeName := string(m.projectMode)
	if modeName == "" {
		modeName = "ready"
	}
	metaParts = append(metaParts, modeName)

	// Loop number
	metaParts = append(metaParts, fmt.Sprintf("loop %d", m.loopNumber))

	// Task progress
	completed := 0
	for _, t := range m.tasks {
		if t.Completed {
			completed++
		}
	}
	if len(m.tasks) > 0 {
		metaParts = append(metaParts, fmt.Sprintf("%d/%d", completed, len(m.tasks)))
	}

	metadata := StyleHeaderMeta.Render(strings.Join(metaParts, MetaDotSeparator))

	// Calculate space for diagonal separators
	brandWidth := lipgloss.Width(brandPrefix) + 1 + lipgloss.Width(brandName) + 1
	metaWidth := lipgloss.Width(metadata) + 1
	diagWidth := width - brandWidth - metaWidth - 2

	if diagWidth < 3 {
		diagWidth = 3
	}

	diagonals := StyleDiagonal.Render(DiagonalSeparator(diagWidth))

	// Assemble header
	leftPart := fmt.Sprintf("%s %s ", brandPrefix, brandName)
	headerContent := leftPart + diagonals + " " + metadata

	return StyleHeader.Copy().Width(width).Render(headerContent)
}

// renderStatusBar renders the status bar with state indicator, spinner, and circuit state
// Format: ● running                                    circuit closed
func (m Model) renderStatusBar(width int) string {
	// State indicator with icon
	var stateIcon, stateText string
	var stateStyle lipgloss.Style

	switch m.state {
	case StateInitializing:
		stateIcon = IconInProgress
		stateText = "initializing"
		stateStyle = StyleInfoMsg
	case StateRunning:
		// Animated spinner for running
		stateIcon = BrailleSpinnerFrames[m.tick%len(BrailleSpinnerFrames)]
		stateText = "running"
		stateStyle = StyleSuccessMsg
	case StatePaused:
		stateIcon = IconPending
		stateText = "paused"
		stateStyle = StyleWarningMsg
	case StateComplete:
		stateIcon = IconCheck
		stateText = "complete"
		stateStyle = StyleSuccessMsg
	case StateError:
		stateIcon = IconError
		stateText = "error"
		stateStyle = StyleErrorMsg
	}

	leftStatus := stateStyle.Render(stateIcon + " " + stateText)

	// Circuit state on right
	circuitState := m.circuitState
	if circuitState == "" {
		circuitState = "closed"
	}
	circuitState = strings.ToLower(circuitState)

	var circuitStyle lipgloss.Style
	switch circuitState {
	case "closed":
		circuitStyle = StyleCircuitClosed
	case "half_open":
		circuitStyle = StyleCircuitHalfOpen
		circuitState = "half-open"
	case "open":
		circuitStyle = StyleCircuitOpen
	default:
		circuitStyle = StyleTextMuted
	}

	rightStatus := StyleTextMuted.Render("circuit ") + circuitStyle.Render(circuitState)

	// Calculate padding between left and right
	leftWidth := lipgloss.Width(leftStatus)
	rightWidth := lipgloss.Width(rightStatus)
	paddingWidth := width - leftWidth - rightWidth - 4
	if paddingWidth < 1 {
		paddingWidth = 1
	}

	statusContent := " " + leftStatus + strings.Repeat(" ", paddingWidth) + rightStatus

	return StyleStatus.Copy().Width(width).Render(statusContent)
}

// renderHorizontalDivider renders a subtle horizontal divider
func (m Model) renderHorizontalDivider(width int) string {
	return StyleDivider.Render(strings.Repeat(DividerChar, width))
}

// renderTaskPane renders the tasks pane with Crush-style icons
func (m Model) renderTaskPane(width, height int) string {
	var lines []string

	// Count completed tasks
	completed := 0
	for _, t := range m.tasks {
		if t.Completed {
			completed++
		}
	}

	if len(m.tasks) == 0 {
		lines = append(lines, StyleTextMuted.Render(" No tasks loaded"))
	} else {
		// Render tasks with icons
		for i, task := range m.tasks {
			if i >= height-2 { // Leave room for summary
				remaining := len(m.tasks) - i
				lines = append(lines, StyleTextSubtle.Render(fmt.Sprintf(" ... %d more", remaining)))
				break
			}
			lines = append(lines, m.renderTaskLine(task, i, width-2))
		}

		// Task summary
		lines = append(lines, "")
		summary := fmt.Sprintf(" %d of %d complete", completed, len(m.tasks))
		lines = append(lines, StyleTextMuted.Render(summary))
	}

	// Pad to fill height
	for len(lines) < height {
		lines = append(lines, "")
	}

	content := strings.Join(lines[:height], "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(0, 1).
		Render(content)
}

// renderTaskLine renders a single task with Crush-style icons
func (m Model) renderTaskLine(task Task, index int, maxWidth int) string {
	isActive := index == m.activeTaskIdx && m.state == StateRunning
	text := task.Text

	// Truncate if needed
	if maxWidth > 10 && len(text) > maxWidth-6 {
		text = text[:maxWidth-9] + "..."
	}

	// Get icon based on state
	var icon string
	var textStyle lipgloss.Style

	if task.Completed {
		icon = StyleTaskCompleted.Render(IconCheck)
		textStyle = StyleTaskTextCompleted
	} else if isActive {
		// Use animated spinner for active task
		spinnerFrame := BrailleSpinnerFrames[m.tick%len(BrailleSpinnerFrames)]
		icon = StyleTaskInProgress.Render(spinnerFrame)
		textStyle = StyleTaskTextActive
	} else {
		icon = StyleTaskPending.Render(IconPending)
		textStyle = StyleTaskTextPending
	}

	return " " + icon + " " + textStyle.Render(text)
}

// renderOutputPane renders the output pane with live Codex output
func (m Model) renderOutputPane(width, height int) string {
	var lines []string

	if len(m.outputLines) == 0 {
		lines = append(lines, StyleTextMuted.Render(" Waiting for Codex output..."))
	} else {
		// Combine output lines and render as markdown
		outputContent := strings.Join(m.outputLines, "\n")
		rendered := renderMarkdown(outputContent, width-4)

		// Split rendered content into lines
		renderedLines := strings.Split(rendered, "\n")

		// Show most recent lines that fit
		maxLines := height - 6 // Leave room for reasoning section
		start := 0
		if len(renderedLines) > maxLines {
			start = len(renderedLines) - maxLines
		}

		for i := start; i < len(renderedLines); i++ {
			lines = append(lines, " "+renderedLines[i])
		}
	}

	// Add reasoning section if we have reasoning
	if len(m.reasoningLines) > 0 {
		lines = append(lines, "")
		// Subtle divider before reasoning
		lines = append(lines, StyleDividerSubtle.Render(strings.Repeat(DividerCharSubtle, width-4)))

		// Combine and render reasoning with markdown support
		reasoningContent := strings.Join(m.reasoningLines[max(0, len(m.reasoningLines)-3):], "\n")
		renderedReasoning := renderMarkdownFallback(reasoningContent)

		lines = append(lines, StyleReasoningHeader.Render(" thinking: ")+StyleReasoning.Render(renderedReasoning))
	}

	// Pad to fill height
	for len(lines) < height {
		lines = append(lines, "")
	}

	content := strings.Join(lines[:height], "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(0, 1).
		Render(content)
}

// renderMarkdown renders markdown content using glamour
func renderMarkdown(content string, width int) string {
	if content == "" {
		return ""
	}

	// Create a glamour renderer with dark style
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		// Fallback: manually handle **bold** -> ANSI bold
		return renderMarkdownFallback(content)
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		return renderMarkdownFallback(content)
	}

	return strings.TrimSpace(rendered)
}

// renderMarkdownFallback provides basic markdown rendering when glamour fails
func renderMarkdownFallback(content string) string {
	// Handle **bold** markers manually
	result := content
	for {
		start := strings.Index(result, "**")
		if start == -1 {
			break
		}
		end := strings.Index(result[start+2:], "**")
		if end == -1 {
			break
		}
		end += start + 2
		boldText := result[start+2 : end]
		// Use lipgloss bold styling
		styled := lipgloss.NewStyle().Bold(true).Render(boldText)
		result = result[:start] + styled + result[end+2:]
	}
	return result
}

// renderFooter renders the Crush-style footer with keybindings
// Format: r run • p pause • l logs • c circuit • ? help • q quit
func (m Model) renderFooter(width int) string {
	bindings := []struct {
		key  string
		desc string
	}{
		{"r", "run"},
		{"p", "pause"},
		{"l", "logs"},
		{"c", "circuit"},
		{"t", "tasks"},
		{"o", "output"},
		{"?", "help"},
		{"q", "quit"},
	}

	var parts []string
	for _, b := range bindings {
		parts = append(parts, StyleHelpKey.Render(b.key)+" "+StyleHelpDesc.Render(b.desc))
	}

	footerContent := " " + strings.Join(parts, StyleTextSubtle.Render(MetaDotSeparator))

	return StyleFooter.Copy().Width(width).Render(footerContent)
}

// renderTasksFullView renders tasks in full screen mode
func (m Model) renderTasksFullView() string {
	width := m.width
	height := m.height
	if width < 60 {
		width = 60
	}
	if height < 20 {
		height = 20
	}

	header := m.renderHeader(width)

	// Task progress header
	completed := 0
	for _, t := range m.tasks {
		if t.Completed {
			completed++
		}
	}

	var lines []string

	progressPct := 0
	if len(m.tasks) > 0 {
		progressPct = (completed * 100) / len(m.tasks)
	}
	lines = append(lines, StyleTextBase.Render(fmt.Sprintf(" Progress: %d/%d (%d%%)", completed, len(m.tasks), progressPct)))

	// Progress bar
	barWidth := 40
	filledCount := 0
	if len(m.tasks) > 0 {
		filledCount = (completed * barWidth) / len(m.tasks)
	}
	emptyCount := barWidth - filledCount
	taskProgressBar := " " + StyleProgressFilled.Render(strings.Repeat("█", filledCount)) +
		StyleProgressEmpty.Render(strings.Repeat("░", emptyCount))
	lines = append(lines, taskProgressBar)
	lines = append(lines, "")
	lines = append(lines, StyleDivider.Render(strings.Repeat(DividerChar, width-4)))
	lines = append(lines, "")

	// All tasks
	for i, task := range m.tasks {
		lines = append(lines, m.renderTaskLine(task, i, width-4))
	}

	content := strings.Join(lines, "\n")

	footer := StyleFooter.Copy().Width(width).Render(
		fmt.Sprintf(" %s return%s%s quit",
			StyleHelpKey.Render("t"),
			StyleTextSubtle.Render(MetaDotSeparator),
			StyleHelpKey.Render("q")),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		content,
		"",
		footer,
	)
}

// renderOutputFullView renders output in full screen mode
func (m Model) renderOutputFullView() string {
	width := m.width
	height := m.height
	if width < 60 {
		width = 60
	}
	if height < 20 {
		height = 20
	}

	header := m.renderHeader(width)

	var lines []string

	if len(m.outputLines) == 0 {
		lines = append(lines, StyleTextMuted.Render(" Waiting for Codex output..."))
	} else {
		// Combine and render output as markdown
		outputContent := strings.Join(m.outputLines, "\n")
		rendered := renderMarkdown(outputContent, width-4)

		renderedLines := strings.Split(rendered, "\n")
		maxLines := height - 6
		start := 0
		if len(renderedLines) > maxLines {
			start = len(renderedLines) - maxLines
		}

		for i := start; i < len(renderedLines); i++ {
			lines = append(lines, " "+renderedLines[i])
		}
	}

	// Reasoning section
	if len(m.reasoningLines) > 0 {
		lines = append(lines, "")
		lines = append(lines, StyleDividerSubtle.Render(strings.Repeat(DividerCharSubtle, width-4)))
		reasoningContent := strings.Join(m.reasoningLines, "\n")
		lines = append(lines, StyleReasoningHeader.Render(" thinking: ")+StyleReasoning.Render(reasoningContent))
	}

	content := strings.Join(lines, "\n")

	footer := StyleFooter.Copy().Width(width).Render(
		fmt.Sprintf(" %s return%s%s quit",
			StyleHelpKey.Render("o"),
			StyleTextSubtle.Render(MetaDotSeparator),
			StyleHelpKey.Render("q")),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		content,
		"",
		footer,
	)
}

// renderLogsFullView renders logs in full screen mode
func (m Model) renderLogsFullView() string {
	width := m.width
	height := m.height
	if width < 60 {
		width = 60
	}
	if height < 20 {
		height = 20
	}

	header := m.renderHeader(width)

	var lines []string

	if len(m.logs) == 0 {
		lines = append(lines, StyleTextMuted.Render(" No log entries yet..."))
	} else {
		// Show all available logs
		maxLines := height - 6
		start := 0
		if len(m.logs) > maxLines {
			start = len(m.logs) - maxLines
		}

		for i := start; i < len(m.logs); i++ {
			lines = append(lines, " "+m.logs[i])
		}
	}

	content := strings.Join(lines, "\n")

	footer := StyleFooter.Copy().Width(width).Render(
		fmt.Sprintf(" %s return%s%d entries",
			StyleHelpKey.Render("l"),
			StyleTextSubtle.Render(MetaDotSeparator),
			len(m.logs)),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		content,
		"",
		footer,
	)
}
