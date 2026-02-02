package tui

import (
	"github.com/brainwhocodes/lisa-loop/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// Re-export styles from internal/tui/style to keep the tui package surface stable.

// Log level constants
const (
	LogLevelDebug   = style.LogLevelDebug
	LogLevelInfo    = style.LogLevelInfo
	LogLevelWarn    = style.LogLevelWarn
	LogLevelError   = style.LogLevelError
	LogLevelSuccess = style.LogLevelSuccess
)

// Spinner frame constants
var (
	BrailleSpinnerFrames = style.BrailleSpinnerFrames
)

// Border styles
var (
	BorderNormal  = style.BorderNormal
	BorderRounded = style.BorderRounded

	StyleBox        = style.StyleBox
	StyleBoxRounded = style.StyleBoxRounded
	StyleBoxError   = style.StyleBoxError
)

// Divider styles
var (
	StyleDivider       = style.StyleDivider
	DividerChar        = style.DividerChar
	StyleDividerSubtle = style.StyleDividerSubtle
	DividerCharSubtle  = style.DividerCharSubtle
)

// Progress bar styles
var (
	StyleProgressEmpty  = style.StyleProgressEmpty
	StyleProgressFilled = style.StyleProgressFilled
	StyleProgressBar    = style.StyleProgressBar
)

// Circuit breaker styles
var (
	StyleCircuitClosed   = style.StyleCircuitClosed
	StyleCircuitHalfOpen = style.StyleCircuitHalfOpen
	StyleCircuitOpen     = style.StyleCircuitOpen
)

// Error panel styles
var (
	StyleErrorPanel = style.StyleErrorPanel
	StyleErrorTitle = style.StyleErrorTitle
	StyleErrorStack = style.StyleErrorStack
)

// Spinner styles
var (
	StyleSpinner       = style.StyleSpinner
	StyleSpinnerActive = style.StyleSpinnerActive
)

// Header styles
var (
	StyleHeader      = style.StyleHeader
	StyleBrandPrefix = style.StyleBrandPrefix
	StyleBrandName   = style.StyleBrandName
	StyleDiagonal    = style.StyleDiagonal
	StyleHeaderMeta  = style.StyleHeaderMeta

	MetaDotSeparator = style.MetaDotSeparator
)

// Status bar styles
var (
	StyleStatus             = style.StyleStatus
	StyleStatusInitializing = style.StyleStatusInitializing
	StyleStatusRunning      = style.StyleStatusRunning
	StyleStatusPaused       = style.StyleStatusPaused
	StyleStatusError        = style.StyleStatusError
	StyleStatusComplete     = style.StyleStatusComplete
)

// Task styles
var (
	StyleTaskCompleted     = style.StyleTaskCompleted
	StyleTaskInProgress    = style.StyleTaskInProgress
	StyleTaskPending       = style.StyleTaskPending
	StyleTaskTextCompleted = style.StyleTaskTextCompleted
	StyleTaskTextActive    = style.StyleTaskTextActive
	StyleTaskTextPending   = style.StyleTaskTextPending
)

// Text styles
var (
	StyleTextBase     = style.StyleTextBase
	StyleTextSelected = style.StyleTextSelected
	StyleTextMuted    = style.StyleTextMuted
	StyleTextSubtle   = style.StyleTextSubtle

	StyleInfoMsg    = style.StyleInfoMsg
	StyleErrorMsg   = style.StyleErrorMsg
	StyleSuccessMsg = style.StyleSuccessMsg
	StyleWarningMsg = style.StyleWarningMsg
)

// Footer/help styles
var (
	StyleHelpKey  = style.StyleHelpKey
	StyleHelpDesc = style.StyleHelpDesc
	StyleFooter   = style.StyleFooter
)

// Pane styles
var (
	StylePane        = style.StylePane
	StylePaneHeader  = style.StylePaneHeader
	StylePaneContent = style.StylePaneContent
)

// Reasoning styles
var (
	StyleReasoning       = style.StyleReasoning
	StyleReasoningHeader = style.StyleReasoningHeader
)

func StyledLogEntry(level, message string) string {
	return style.StyledLogEntry(level, message)
}

func FormatTaskIcon(completed, active bool, spinnerFrame string) string {
	return style.FormatTaskIcon(completed, active, spinnerFrame)
}

// Compile-time assertion: we still expose lipgloss.Style values (not pointers).
var _ lipgloss.Style = StyleHeader
