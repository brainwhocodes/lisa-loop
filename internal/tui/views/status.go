package views

import (
	"fmt"
	"strings"

	"github.com/brainwhocodes/ralph-codex/internal/tui"
)

// Spinner frames for animation
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// GetSpinnerFrame returns the current spinner frame
func GetSpinnerFrame(tick int) string {
	return spinnerFrames[tick%len(spinnerFrames)]
}

// Render renders status view with live updates
func Render(loopNum int, callsMade int, callsRemaining int, circuitState string, codexStatus string, tick int) string {
	header := tui.StyleHeader.Render("Ralph Codex - Status")

	// Status badge based on circuit state
	var statusBadge string
	switch circuitState {
	case "CLOSED":
		statusBadge = tui.StyleCircuitClosed.Render(circuitState)
	case "HALF_OPEN":
		statusBadge = tui.StyleCircuitHalfOpen.Render(circuitState)
	case "OPEN":
		statusBadge = tui.StyleCircuitOpen.Render(circuitState)
	default:
		statusBadge = circuitState
	}

	// Progress bar for rate limit
	progressBar := renderRateLimitProgress(callsMade, callsRemaining)

	// Animated spinner if running
	spinner := ""
	if codexStatus == "RUNNING" {
		spinner = fmt.Sprintf(" %s", tui.StyleSpinner.Render(GetSpinnerFrame(tick)))
	}

	status := tui.StyleStatus.Render(
		fmt.Sprintf("Loop: %d%s | %s | Circuit: %s",
			loopNum, spinner, progressBar, statusBadge),
	)

	statusDetail := fmt.Sprintf("Codex Status: %s", codexStatus)

	return header + "\n" + status + "\n" + statusDetail
}

// renderRateLimitProgress creates a styled progress bar for rate limiting
func renderRateLimitProgress(callsMade, callsRemaining int) string {
	total := callsMade + callsRemaining
	if total == 0 {
		return "Calls: 0/0"
	}

	progress := float64(callsMade) / float64(total)
	width := 20
	filled := int(progress * float64(width))

	emptyWidth := width - filled
	if emptyWidth < 0 {
		emptyWidth = 0
	}

	bar := fmt.Sprintf("Calls: %d/%d [%s%s]",
		callsMade, total,
		tui.StyleProgressBar.Render(strings.Repeat("█", filled)),
		strings.Repeat("░", emptyWidth),
	)

	return bar
}

// UpdateProgressBar returns a progress bar string (legacy, kept for compatibility)
func UpdateProgressBar(progress float64) string {
	width := 40
	filled := int(progress * float64(width))

	bar := "[" + tui.StyleHelpKey.Render(strings.Repeat("=", filled)) +
		tui.StyleHelpDesc.Render(strings.Repeat(" ", width-filled)) + "]"

	return bar
}

// FormatCircuitState returns styled circuit state badge
func FormatCircuitState(state string) string {
	switch state {
	case "CLOSED":
		return tui.StyleStatusRunning.Render("CLOSED")
	case "HALF_OPEN":
		return tui.StyleStatusPaused.Render("HALF_OPEN")
	case "OPEN":
		return tui.StyleStatusError.Render("OPEN")
	default:
		return state
	}
}

// FormatWorkType returns styled work type badge
func FormatWorkType(workType string) string {
	if workType == "" {
		return workType
	}

	var badge string
	switch workType {
	case "IMPLEMENTATION":
		badge = tui.StyleStatusRunning.Render(workType)
	case "TESTING":
		badge = tui.StyleHelpKey.Render(workType)
	case "DOCUMENTATION":
		badge = tui.StyleHelpDesc.Render(workType)
	case "REFACTORING":
		badge = tui.StyleInfoMsg.Render(workType)
	default:
		badge = workType
	}

	return badge
}
