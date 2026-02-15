package tui

import (
	"fmt"
	"strings"

	"github.com/brainwhocodes/lisa-loop/internal/config"
	"github.com/brainwhocodes/lisa-loop/internal/tui/transcript"
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
// Format: Charm LISA SAX ♫ ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱ mode • loop 2 • 3/10
func (m Model) renderHeader(width int) string {
	// Brand prefix and name
	brandPrefix := StyleBrandPrefix.Render("Charm")
	brandName := GradientText("LISA", Dolly, Charple)

	// Animated SAX with musical notes when running
	var saxAnim string
	if m.state == StateRunning {
		frame := m.tick % len(SaxNotes)
		saxAnim = " " + StyleBrandPrefix.Render(SaxNotes[frame])
	}

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

	// Phase-aware task progress
	if len(m.phases) > 0 {
		currentPhaseIdx := m.getCurrentPhaseIndex()
		if currentPhaseIdx >= 0 && currentPhaseIdx < len(m.phases) {
			phase := m.phases[currentPhaseIdx]
			completed := 0
			for _, t := range phase.Tasks {
				if t.Completed {
					completed++
				}
			}
			metaParts = append(metaParts, fmt.Sprintf("P%d:%d/%d", currentPhaseIdx+1, completed, len(phase.Tasks)))
		}
	} else if len(m.tasks) > 0 {
		// Fallback to flat task progress
		completed := 0
		for _, t := range m.tasks {
			if t.Completed {
				completed++
			}
		}
		metaParts = append(metaParts, fmt.Sprintf("%d/%d", completed, len(m.tasks)))
	}

	metadata := StyleHeaderMeta.Render(strings.Join(metaParts, MetaDotSeparator))

	// Calculate space for diagonal separators (include SAX animation width)
	brandWidth := lipgloss.Width(brandPrefix) + 1 + lipgloss.Width(brandName) + lipgloss.Width(saxAnim) + 1
	metaWidth := lipgloss.Width(metadata) + 1
	diagWidth := width - brandWidth - metaWidth - 2

	if diagWidth < 3 {
		diagWidth = 3
	}

	diagonals := StyleDiagonal.Render(DiagonalSeparator(diagWidth))

	// Assemble header
	leftPart := fmt.Sprintf("%s %s%s ", brandPrefix, brandName, saxAnim)
	headerContent := leftPart + diagonals + " " + metadata

	return StyleHeader.Width(width).Render(headerContent)
}

// renderStatusBar renders the status bar with state indicator, spinner, and circuit state
// Format: ● running  STATUS: WORKING  tasks: 2  files: 3     circuit closed
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

	// Add analysis status in the middle if available
	var midStatus string
	if m.analysisStatus != "" {
		var statusStyle lipgloss.Style
		switch m.analysisStatus {
		case "COMPLETE":
			statusStyle = StyleSuccessMsg
		case "BLOCKED":
			statusStyle = StyleErrorMsg
		default:
			statusStyle = StyleInfoMsg
		}
		midStatus = StyleTextMuted.Render(" │ ") + statusStyle.Render(m.analysisStatus)

		if m.tasksCompleted > 0 {
			midStatus += StyleTextMuted.Render(fmt.Sprintf(" tasks:%d", m.tasksCompleted))
		}
		if m.filesModified > 0 {
			midStatus += StyleTextMuted.Render(fmt.Sprintf(" files:%d", m.filesModified))
		}
		if m.exitSignal {
			midStatus += StyleSuccessMsg.Render(" EXIT")
		}
	}

	// Context usage indicator (before circuit)
	var contextIndicator string
	if m.contextLimit > 0 {
		usagePct := int(m.contextUsagePercent * 100)
		var ctxStyle lipgloss.Style
		switch {
		case m.contextThreshold:
			ctxStyle = StyleErrorMsg // Red when threshold reached
		case m.contextUsagePercent >= 0.6:
			ctxStyle = StyleWarningMsg // Yellow when > 60%
		default:
			ctxStyle = StyleTextMuted // Normal
		}
		// Mini progress bar for context
		barWidth := 10
		filled := int(m.contextUsagePercent * float64(barWidth))
		if filled > barWidth {
			filled = barWidth
		}
		empty := barWidth - filled
		contextIndicator = StyleTextMuted.Render("ctx ") +
			ctxStyle.Render(fmt.Sprintf("%d%%", usagePct)) +
			StyleTextMuted.Render(" [") +
			ctxStyle.Render(strings.Repeat("█", filled)) +
			StyleProgressEmpty.Render(strings.Repeat("░", empty)) +
			StyleTextMuted.Render("]")
		if m.contextWasCompacted {
			contextIndicator += StyleWarningMsg.Render(" ⟳")
		}
	}

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

	rightStatus := ""
	if contextIndicator != "" {
		rightStatus = contextIndicator + StyleTextMuted.Render("  ")
	}
	rightStatus += StyleTextMuted.Render("circuit ") + circuitStyle.Render(circuitState)

	// Calculate padding between left and right
	leftWidth := lipgloss.Width(leftStatus) + lipgloss.Width(midStatus)
	rightWidth := lipgloss.Width(rightStatus)
	paddingWidth := width - leftWidth - rightWidth - 4
	if paddingWidth < 1 {
		paddingWidth = 1
	}

	statusContent := " " + leftStatus + midStatus + strings.Repeat(" ", paddingWidth) + rightStatus

	return StyleStatus.Width(width).Render(statusContent)
}

// renderHorizontalDivider renders a subtle horizontal divider
func (m Model) renderHorizontalDivider(width int) string {
	return StyleDivider.Render(strings.Repeat(DividerChar, width))
}

// renderTaskPane renders the tasks pane with Crush-style icons
// Shows only current phase tasks with phase header
func (m Model) renderTaskPane(width, height int) string {
	var lines []string

	// Check for phases
	if len(m.phases) == 0 {
		// Fallback to flat task list if no phases
		return m.renderFlatTaskPane(width, height)
	}

	// Get current phase (auto-advance if current is complete)
	currentPhaseIdx := m.getCurrentPhaseIndex()
	if currentPhaseIdx < 0 || currentPhaseIdx >= len(m.phases) {
		currentPhaseIdx = 0
	}
	phase := m.phases[currentPhaseIdx]

	// Phase header with progress indicator
	phaseCompleted := 0
	for _, t := range phase.Tasks {
		if t.Completed {
			phaseCompleted++
		}
	}

	// Animated phase indicator when running
	var phaseIcon string
	if phase.Completed {
		phaseIcon = StyleTaskCompleted.Render(IconCheck)
	} else if m.state == StateRunning {
		spinnerFrame := BrailleSpinnerFrames[m.tick%len(BrailleSpinnerFrames)]
		phaseIcon = StyleTaskInProgress.Render(spinnerFrame)
	} else {
		phaseIcon = StyleTaskPending.Render(IconInProgress)
	}

	// Phase header: "● Phase 1: Foundation [2/4]"
	phaseHeader := fmt.Sprintf(" %s %s [%d/%d]", phaseIcon, phase.Name, phaseCompleted, len(phase.Tasks))
	lines = append(lines, StyleTextBase.Render(phaseHeader))
	lines = append(lines, "")

	// Render phase tasks with icons
	for i, task := range phase.Tasks {
		if i >= height-4 { // Leave room for header and summary
			remaining := len(phase.Tasks) - i
			lines = append(lines, StyleTextSubtle.Render(fmt.Sprintf(" ... %d more", remaining)))
			break
		}
		lines = append(lines, m.renderPhaseTaskLine(task, i, currentPhaseIdx, width-2))
	}

	// Phase summary with overall progress
	lines = append(lines, "")
	totalPhases := len(m.phases)
	completedPhases := 0
	for _, p := range m.phases {
		if p.Completed {
			completedPhases++
		}
	}
	summary := fmt.Sprintf(" Phase %d/%d • %d/%d tasks", currentPhaseIdx+1, totalPhases, phaseCompleted, len(phase.Tasks))
	lines = append(lines, StyleTextMuted.Render(summary))

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

// renderFlatTaskPane renders tasks without phase grouping (fallback)
func (m Model) renderFlatTaskPane(width, height int) string {
	var lines []string

	completed := 0
	for _, t := range m.tasks {
		if t.Completed {
			completed++
		}
	}

	if len(m.tasks) == 0 {
		lines = append(lines, StyleTextMuted.Render(" No tasks loaded"))
	} else {
		for i, task := range m.tasks {
			if i >= height-2 {
				remaining := len(m.tasks) - i
				lines = append(lines, StyleTextSubtle.Render(fmt.Sprintf(" ... %d more", remaining)))
				break
			}
			lines = append(lines, m.renderTaskLine(task, i, width-2))
		}
		lines = append(lines, "")
		summary := fmt.Sprintf(" %d of %d complete", completed, len(m.tasks))
		lines = append(lines, StyleTextMuted.Render(summary))
	}

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

// getCurrentPhaseIndex returns the index of the first incomplete phase
func (m Model) getCurrentPhaseIndex() int {
	for i, phase := range m.phases {
		if !phase.Completed {
			return i
		}
	}
	// All complete, return last
	if len(m.phases) > 0 {
		return len(m.phases) - 1
	}
	return 0
}

// renderPhaseTaskLine renders a task line within a phase context
func (m Model) renderPhaseTaskLine(task Task, taskIdx, phaseIdx int, maxWidth int) string {
	// Calculate global task index for active tracking
	globalIdx := 0
	for i := 0; i < phaseIdx; i++ {
		globalIdx += len(m.phases[i].Tasks)
	}
	globalIdx += taskIdx

	isActive := globalIdx == m.activeTaskIdx && m.state == StateRunning
	text := task.Text

	// Truncate if needed
	if maxWidth > 10 && len(text) > maxWidth-6 {
		text = text[:maxWidth-9] + "..."
	}

	var icon string
	var textStyle lipgloss.Style

	if task.Completed {
		icon = StyleTaskCompleted.Render(IconCheck)
		textStyle = StyleTaskTextCompleted
	} else if isActive {
		spinnerFrame := BrailleSpinnerFrames[m.tick%len(BrailleSpinnerFrames)]
		icon = StyleTaskInProgress.Render(spinnerFrame)
		textStyle = StyleTaskTextActive
	} else {
		icon = StyleTaskPending.Render(IconPending)
		textStyle = StyleTaskTextPending
	}

	return " " + icon + " " + textStyle.Render(text)
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

// backendDisplayName returns a display-friendly name for the backend
func (m Model) backendDisplayName() string {
	cfg := config.Config{Backend: m.backend}
	return cfg.BackendDisplayName()
}

// renderOutputPane renders the output pane with live agent output
func (m Model) renderOutputPane(width, height int) string {
	var lines []string

	// Calculate space for reasoning (if any)
	reasoningHeight := 0
	if len(m.reasoningLines) > 0 {
		reasoningHeight = 3 // divider + thinking line + padding
	}

	// Show reasoning at the top if we have it
	if len(m.reasoningLines) > 0 {
		// Get the latest reasoning (last line only to avoid clutter)
		latestReasoning := m.reasoningLines[len(m.reasoningLines)-1]
		// Reasoning can be multi-line; keep it compact in split view.
		latestReasoning = strings.ReplaceAll(latestReasoning, "\n", " ")
		// Truncate if too long
		if len(latestReasoning) > width-20 {
			latestReasoning = latestReasoning[:width-23] + "..."
		}
		// Animated thinking indicator
		thinkAnim := ThinkingWave[m.tick%len(ThinkingWave)]
		lines = append(lines, StyleReasoning.Render(" ["+thinkAnim+"] "+latestReasoning))
		lines = append(lines, "")
	}

	if len(m.outputLines) == 0 && len(m.reasoningLines) == 0 {
		lines = append(lines, StyleTextMuted.Render(fmt.Sprintf(" Waiting for %s output...", m.backendDisplayName())))
	} else if len(m.outputLines) > 0 {
		// Flatten multi-line output entries so escaped newlines render as actual lines.
		var flat []string
		for _, ol := range m.outputLines {
			flat = append(flat, strings.Split(ol, "\n")...)
		}

		// Show most recent output lines that fit.
		maxLines := height - reasoningHeight - 2
		if maxLines < 1 {
			maxLines = 1
		}
		start := 0
		if len(flat) > maxLines {
			start = len(flat) - maxLines
		}

		for i := start; i < len(flat); i++ {
			line := flat[i]
			// Truncate long lines
			if len(line) > width-4 {
				line = line[:width-7] + "..."
			}
			lines = append(lines, " "+line)
		}
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

// renderFooter renders the Crush-style footer with keybindings
// Format: r run • p pause • l logs • c circuit • ? help • q quit
func (m Model) renderFooter(width int) string {
	var parts []string
	for _, b := range footerBindings() {
		parts = append(parts, StyleHelpKey.Render(b.Key)+" "+StyleHelpDesc.Render(b.Description))
	}

	footerContent := " " + strings.Join(parts, StyleTextSubtle.Render(MetaDotSeparator))

	return StyleFooter.Width(width).Render(footerContent)
}

// renderTasksFullView renders tasks in full screen mode
// Shows all phases with current phase expanded
func (m Model) renderTasksFullView() string {
	width := m.width
	if width < 60 {
		width = 60
	}

	header := m.renderHeader(width)

	var lines []string

	// If we have phases, show phase-organized view
	if len(m.phases) > 0 {
		currentPhaseIdx := m.getCurrentPhaseIndex()

		// Overall progress
		totalTasks := 0
		completedTasks := 0
		for _, phase := range m.phases {
			for _, t := range phase.Tasks {
				totalTasks++
				if t.Completed {
					completedTasks++
				}
			}
		}

		progressPct := 0
		if totalTasks > 0 {
			progressPct = (completedTasks * 100) / totalTasks
		}
		lines = append(lines, StyleTextBase.Render(fmt.Sprintf(" Overall: %d/%d (%d%%)", completedTasks, totalTasks, progressPct)))

		// Progress bar
		barWidth := 40
		filledCount := 0
		if totalTasks > 0 {
			filledCount = (completedTasks * barWidth) / totalTasks
		}
		emptyCount := barWidth - filledCount
		taskProgressBar := " " + StyleProgressFilled.Render(strings.Repeat("█", filledCount)) +
			StyleProgressEmpty.Render(strings.Repeat("░", emptyCount))
		lines = append(lines, taskProgressBar)
		lines = append(lines, "")

		// Render each phase
		globalTaskIdx := 0
		for phaseIdx, phase := range m.phases {
			// Phase header with icon
			var phaseIcon string
			if phase.Completed {
				phaseIcon = StyleTaskCompleted.Render(IconCheck)
			} else if phaseIdx == currentPhaseIdx && m.state == StateRunning {
				spinnerFrame := BrailleSpinnerFrames[m.tick%len(BrailleSpinnerFrames)]
				phaseIcon = StyleTaskInProgress.Render(spinnerFrame)
			} else if phaseIdx == currentPhaseIdx {
				phaseIcon = StyleTaskPending.Render(IconInProgress)
			} else {
				phaseIcon = StyleTextMuted.Render(IconPending)
			}

			// Count phase progress
			phaseCompleted := 0
			for _, t := range phase.Tasks {
				if t.Completed {
					phaseCompleted++
				}
			}

			phaseHeader := fmt.Sprintf(" %s %s [%d/%d]", phaseIcon, phase.Name, phaseCompleted, len(phase.Tasks))
			if phaseIdx == currentPhaseIdx {
				lines = append(lines, StyleTextBase.Render(phaseHeader))
			} else {
				lines = append(lines, StyleTextMuted.Render(phaseHeader))
			}

			// Show tasks for current phase, collapse others
			if phaseIdx == currentPhaseIdx {
				for taskIdx, task := range phase.Tasks {
					lines = append(lines, m.renderPhaseTaskLine(task, taskIdx, phaseIdx, width-4))
					globalTaskIdx++
				}
			} else {
				globalTaskIdx += len(phase.Tasks)
			}
			lines = append(lines, "")
		}
	} else {
		// Fallback: flat task list
		completed := 0
		for _, t := range m.tasks {
			if t.Completed {
				completed++
			}
		}

		progressPct := 0
		if len(m.tasks) > 0 {
			progressPct = (completed * 100) / len(m.tasks)
		}
		lines = append(lines, StyleTextBase.Render(fmt.Sprintf(" Progress: %d/%d (%d%%)", completed, len(m.tasks), progressPct)))

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

		for i, task := range m.tasks {
			lines = append(lines, m.renderTaskLine(task, i, width-4))
		}
	}

	content := strings.Join(lines, "\n")

	footer := StyleFooter.Width(width).Render(
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

	tabs := m.renderOutputTabs(width)

	// Layout
	const headerHeight = 1
	const footerHeight = 1
	const tabsHeight = 1
	contentHeight := height - headerHeight - tabsHeight - footerHeight - 2
	if contentHeight < 8 {
		contentHeight = 8
	}

	content := m.renderOutputTabContent(width, contentHeight)

	footer := StyleFooter.Width(width).Render(
		fmt.Sprintf(" %s return%s%s tabs%s%s reasoning%s%s quit",
			StyleHelpKey.Render("o"),
			StyleTextSubtle.Render(MetaDotSeparator),
			StyleHelpKey.Render("[ ]"),
			StyleTextSubtle.Render(MetaDotSeparator),
			StyleHelpKey.Render("y"),
			StyleTextSubtle.Render(MetaDotSeparator),
			StyleHelpKey.Render("q")),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		tabs,
		"",
		content,
		"",
		footer,
	)
}

func (m Model) renderOutputTabs(width int) string {
	tabs := []OutputTab{OutputTabTranscript, OutputTabDiffs, OutputTabReasoning}

	var parts []string
	for _, t := range tabs {
		label := t.String()
		if t == m.outputTab {
			parts = append(parts, lipgloss.NewStyle().
				Foreground(Pepper).
				Background(Guac).
				Padding(0, 1).
				Bold(true).
				Render(label))
		} else {
			parts = append(parts, lipgloss.NewStyle().
				Foreground(Squid).
				Background(BBQ).
				Padding(0, 1).
				Render(label))
		}
	}

	line := lipgloss.JoinHorizontal(lipgloss.Top, parts...)
	return lipgloss.NewStyle().Width(width).Padding(0, 1).Render(line)
}

func (m Model) renderOutputTabContent(width, height int) string {
	switch m.outputTab {
	case OutputTabDiffs:
		return m.renderDiffTab(width, height)
	case OutputTabReasoning:
		return m.renderReasoningTab(width, height)
	default:
		return m.renderTranscriptTab(width, height)
	}
}

func (m Model) renderTranscriptTab(width, height int) string {
	// If the structured transcript is not populated yet, fall back to the legacy output.
	if m.transcript == nil || m.transcript.Len() == 0 {
		return m.renderOutputPane(width, height)
	}

	items := m.transcript.Items()

	max := height
	start := 0
	if len(items) > max {
		start = len(items) - max
	}

	lines := make([]string, 0, max)
	for i := start; i < len(items); i++ {
		it := items[i]
		lines = append(lines, " "+m.formatTranscriptLine(width-2, it))
		if len(lines) >= height {
			break
		}
	}

	for len(lines) < height {
		lines = append(lines, "")
	}

	return lipgloss.NewStyle().Width(width).Height(height).Padding(0, 1).Render(strings.Join(lines[:height], "\n"))
}

func (m Model) formatTranscriptLine(maxWidth int, it transcript.Item) string {
	ts := ""
	if !it.At.IsZero() {
		ts = it.At.Format("15:04:05")
	}

	roleStyle := StyleTextMuted
	switch it.Role {
	case transcript.RoleAssistant:
		roleStyle = StyleSuccessMsg
	case transcript.RoleTool:
		roleStyle = StyleInfoMsg
	case transcript.RoleUser:
		roleStyle = StyleTextSelected
	case transcript.RoleSystem:
		roleStyle = StyleTextMuted
	}

	role := roleStyle.Render(string(it.Role))

	body := strings.TrimSpace(it.Body)
	body = strings.ReplaceAll(body, "\n", " ")
	if body == "" && it.Title != "" {
		body = it.Title
	}

	line := fmt.Sprintf("%s %s %s", ts, role, body)
	if maxWidth > 10 && lipgloss.Width(line) > maxWidth {
		// Best-effort truncation for plain text; ANSI width can be imperfect but should be acceptable here.
		runes := []rune(line)
		if len(runes) > maxWidth-1 {
			line = string(runes[:maxWidth-1]) + "…"
		}
	}
	return line
}

func (m Model) renderDiffTab(width, height int) string {
	// Prefer git diff output; fall back to pending changes if git is unavailable.
	var md strings.Builder
	md.WriteString("## Diffs\n\n")

	if m.diffPending {
		md.WriteString("Status: **pending** (collecting `git diff`)\n\n")
	}

	if m.diffErr != nil {
		md.WriteString("Status: **error**\n\n")
		md.WriteString("```text\n")
		md.WriteString(m.diffErr.Error())
		md.WriteString("\n```\n\n")
	}

	if m.gitDiffNameStatus == "" && m.gitDiffPatch == "" {
		if len(m.pendingChanges) == 0 {
			md.WriteString("_No changes detected._\n")
		} else {
			md.WriteString("### Touched files (unverified)\n\n")
			for _, pc := range m.pendingChanges {
				md.WriteString(fmt.Sprintf("- `%s` (%s %s)\n", pc.Path, pc.Tool, pc.Status))
			}
		}
	} else {
		md.WriteString("### Changed files\n\n")
		md.WriteString("```text\n")
		md.WriteString(m.gitDiffNameStatus)
		md.WriteString("\n```\n\n")

		md.WriteString("### Patch\n\n")
		md.WriteString("```diff\n")
		md.WriteString(m.gitDiffPatch)
		md.WriteString("\n```\n")
	}

	rendered := md.String()
	if m.md != nil {
		out, err := m.md.Render(width-2, rendered)
		if err == nil {
			rendered = out
		}
	}

	// Clamp to height.
	lines := strings.Split(rendered, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}

	return lipgloss.NewStyle().Width(width).Height(height).Padding(0, 1).Render(strings.Join(lines, "\n"))
}

func (m Model) renderReasoningTab(width, height int) string {
	reasoning := m.currentReasoning
	if reasoning == "" && len(m.reasoningLines) > 0 {
		reasoning = m.reasoningLines[len(m.reasoningLines)-1]
	}

	var md strings.Builder
	md.WriteString("## Reasoning\n\n")
	if reasoning == "" {
		md.WriteString("_No reasoning yet._\n")
	} else {
		if m.reasoningExpanded {
			md.WriteString("_Showing full reasoning (`y` to collapse)._")
		} else {
			md.WriteString("_Showing truncated reasoning (`y` to expand)._")
		}
		md.WriteString("\n\n")

		// Keep rendering fast and readable by default; allow full expansion via toggle.
		if !m.reasoningExpanded {
			const maxRunes = 1200
			r := []rune(reasoning)
			if len(r) > maxRunes {
				reasoning = string(r[:maxRunes]) + "\n…"
			}
		}

		md.WriteString("```text\n")
		md.WriteString(reasoning)
		md.WriteString("\n```\n")
	}

	rendered := md.String()
	if m.md != nil {
		out, err := m.md.Render(width-2, rendered)
		if err == nil {
			rendered = out
		}
	}

	lines := strings.Split(rendered, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}

	return lipgloss.NewStyle().Width(width).Height(height).Padding(0, 1).Render(strings.Join(lines, "\n"))
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

	footer := StyleFooter.Width(width).Render(
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
