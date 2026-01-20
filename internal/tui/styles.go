package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Log level constants
const (
	LogLevelInfo    = "INFO"
	LogLevelWarn    = "WARN"
	LogLevelError   = "ERROR"
	LogLevelSuccess = "SUCCESS"
)

// Spinner frame constants
var (
	// Braille spinner for status bar and active tasks
	BrailleSpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
)

// Border styles
var (
	BorderNormal  = lipgloss.NormalBorder()
	BorderRounded = lipgloss.RoundedBorder()

	// Subtle box style with Charmtone colors
	StyleBox = lipgloss.NewStyle().
			Border(BorderNormal, true, true, true, true).
			BorderForeground(Charcoal).
			Padding(1)

	StyleBoxRounded = lipgloss.NewStyle().
			Border(BorderRounded, true, true, true, true).
			BorderForeground(Charcoal).
			Padding(1)

	StyleBoxError = lipgloss.NewStyle().
			Border(BorderNormal, true, true, true, true).
			BorderForeground(Sriracha).
			Padding(1)
)

// Divider styles - subtle horizontal lines
var (
	StyleDivider = lipgloss.NewStyle().
			Foreground(Charcoal)

	DividerChar = "─"

	StyleDividerSubtle = lipgloss.NewStyle().
				Foreground(Iron)

	DividerCharSubtle = "┈"
)

// Progress bar styles
var (
	StyleProgressEmpty = lipgloss.NewStyle().
				Foreground(Oyster)

	StyleProgressFilled = lipgloss.NewStyle().
				Foreground(Guac)

	StyleProgressBar = lipgloss.NewStyle().
				Foreground(Guac)
)

// Circuit breaker styles
var (
	StyleCircuitClosed = lipgloss.NewStyle().
				Foreground(Guac).
				Bold(true)

	StyleCircuitHalfOpen = lipgloss.NewStyle().
				Foreground(Zest).
				Bold(true)

	StyleCircuitOpen = lipgloss.NewStyle().
				Foreground(Sriracha).
				Bold(true)
)

// Error panel styles
var (
	StyleErrorPanel = lipgloss.NewStyle().
			Foreground(Salt).
			Background(Sriracha).
			Padding(1, 2)

	StyleErrorTitle = lipgloss.NewStyle().
			Foreground(Salt).
			Bold(true).
			Underline(true)

	StyleErrorStack = lipgloss.NewStyle().
			Foreground(Smoke).
			Italic(true)
)

// Spinner styles
var (
	StyleSpinner = lipgloss.NewStyle().
			Foreground(Guac)

	StyleSpinnerActive = lipgloss.NewStyle().
				Foreground(Julep)
)

// Header styles - Crush-inspired
var (
	// Main header with brand gradient
	StyleHeader = lipgloss.NewStyle().
			Foreground(Salt).
			Background(BBQ).
			Padding(0, 1)

	// Brand text "Charm" style
	StyleBrandPrefix = lipgloss.NewStyle().
				Foreground(Dolly).
				Bold(false)

	// Brand name "RALPH" in gradient
	StyleBrandName = lipgloss.NewStyle().
			Foreground(Charple).
			Bold(true)

	// Diagonal separator
	StyleDiagonal = lipgloss.NewStyle().
			Foreground(Charple)

	// Header metadata (right side)
	StyleHeaderMeta = lipgloss.NewStyle().
			Foreground(Squid)

	// Dot separator for metadata
	MetaDotSeparator = " • "
)

// Status bar styles
var (
	StyleStatus = lipgloss.NewStyle().
			Foreground(Ash).
			Background(Pepper).
			Padding(0, 1)

	// Status badges
	StyleStatusInitializing = lipgloss.NewStyle().
				Foreground(Pepper).
				Background(Malibu).
				Padding(0, 1).
				Bold(true)

	StyleStatusRunning = lipgloss.NewStyle().
				Foreground(Pepper).
				Background(Guac).
				Padding(0, 1).
				Bold(true)

	StyleStatusPaused = lipgloss.NewStyle().
				Foreground(Pepper).
				Background(Zest).
				Padding(0, 1).
				Bold(true)

	StyleStatusError = lipgloss.NewStyle().
				Foreground(Salt).
				Background(Sriracha).
				Padding(0, 1).
				Bold(true)

	StyleStatusComplete = lipgloss.NewStyle().
				Foreground(Pepper).
				Background(Charple).
				Padding(0, 1).
				Bold(true)
)

// Task/Todo styles - Crush icons
var (
	// Completed task: green checkmark
	StyleTaskCompleted = lipgloss.NewStyle().
				Foreground(Guac)

	// In-progress task: darker green dot with spinner
	StyleTaskInProgress = lipgloss.NewStyle().
				Foreground(Julep)

	// Pending task: muted bullet
	StyleTaskPending = lipgloss.NewStyle().
				Foreground(Squid)

	// Task text styles
	StyleTaskTextCompleted = lipgloss.NewStyle().
				Foreground(Smoke)

	StyleTaskTextActive = lipgloss.NewStyle().
				Foreground(Salt)

	StyleTaskTextPending = lipgloss.NewStyle().
				Foreground(Squid)
)

// Text styles
var (
	// Primary text
	StyleTextBase = lipgloss.NewStyle().
			Foreground(Ash)

	// Selected/highlighted text
	StyleTextSelected = lipgloss.NewStyle().
				Foreground(Salt)

	// Muted/secondary text
	StyleTextMuted = lipgloss.NewStyle().
			Foreground(Squid)

	// Subtle/hint text
	StyleTextSubtle = lipgloss.NewStyle().
			Foreground(Oyster)

	// Info message
	StyleInfoMsg = lipgloss.NewStyle().
			Foreground(Malibu).
			Bold(true)

	// Error message
	StyleErrorMsg = lipgloss.NewStyle().
			Foreground(Sriracha).
			Bold(true)

	// Success message
	StyleSuccessMsg = lipgloss.NewStyle().
			Foreground(Guac).
			Bold(true)

	// Warning message
	StyleWarningMsg = lipgloss.NewStyle().
			Foreground(Zest).
			Bold(true)
)

// Footer/help styles
var (
	StyleHelpKey = lipgloss.NewStyle().
			Foreground(Charple).
			Bold(true)

	StyleHelpDesc = lipgloss.NewStyle().
			Foreground(Squid)

	StyleFooter = lipgloss.NewStyle().
			Foreground(Squid).
			Background(Pepper).
			Padding(0, 1)
)

// Pane styles
var (
	StylePane = lipgloss.NewStyle().
			Background(Pepper)

	StylePaneHeader = lipgloss.NewStyle().
			Foreground(Ash).
			Bold(true)

	StylePaneContent = lipgloss.NewStyle().
			Foreground(Smoke)
)

// Reasoning/thinking styles
var (
	StyleReasoning = lipgloss.NewStyle().
			Foreground(Squid).
			Italic(true)

	StyleReasoningHeader = lipgloss.NewStyle().
				Foreground(Oyster).
				Italic(true)
)

// StyledLogEntry returns a styled log entry with Crush-style icons
func StyledLogEntry(level, message string) string {
	switch level {
	case LogLevelInfo:
		return StyleTextMuted.Render(IconInfo + " " + message)
	case LogLevelWarn:
		return StyleWarningMsg.Render(IconWarning + " " + message)
	case LogLevelError:
		return StyleErrorMsg.Render(IconError + " " + message)
	case LogLevelSuccess:
		return StyleSuccessMsg.Render(IconCheck + " " + message)
	default:
		return StyleTextMuted.Render("  " + message)
	}
}

// FormatTaskIcon returns the appropriate icon and style for a task state
func FormatTaskIcon(completed, active bool, spinnerFrame string) string {
	if completed {
		return StyleTaskCompleted.Render(IconCheck)
	}
	if active {
		return StyleTaskInProgress.Render(spinnerFrame)
	}
	return StyleTaskPending.Render(IconPending)
}
