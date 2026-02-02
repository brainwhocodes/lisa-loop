package tui

// Screen is the active UI screen/page.
// This is the single source of truth for routing (replaces activeView/helpVisible/viewMode).
type Screen int

const (
	ScreenSplit Screen = iota
	ScreenTasks
	ScreenOutput
	ScreenLogs
	ScreenHelp
	ScreenCircuit
)
