package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// State represents TUI state
type State int

const (
	StateInitializing State = iota
	StateRunning
	StatePaused
	StateComplete
	StateError
)

func (s State) String() string {
	switch s {
	case StateInitializing:
		return "Initializing"
	case StateRunning:
		return "Running"
	case StatePaused:
		return "Paused"
	case StateComplete:
		return "Complete"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// LoopUpdateMsg is sent when loop controller updates
type LoopUpdateMsg struct {
	LoopNumber int
	CallsUsed  int
	Status     string
}

// LogMsg is sent to add a log entry
type LogMsg struct {
	Message string
	Level   string // INFO, WARN, ERROR, SUCCESS
}

// StateChangeMsg is sent to change TUI state
type StateChangeMsg struct {
	State State
}

// StatusMsg is sent to update status text
type StatusMsg struct {
	Status string
}

// TickMsg is sent periodically for animations
type TickMsg time.Time

// Model represents main TUI model
type Model struct {
	state        State
	status       string
	loopNumber   int
	maxCalls     int
	callsUsed    int
	circuitState string
	logs         []string
	activeView   string
	quitting     bool
	err          error
	helpVisible  bool
	startTime    time.Time
	tick         int // Animation tick counter
}

// Init initializes model
func (m Model) Init() tea.Cmd {
	m.startTime = time.Now()
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// styledLogEntry returns a styled log entry
func styledLogEntry(level string, message string) string {
	switch level {
	case "INFO":
		return StyleLog.Render("[INFO] " + message)
	case "WARN":
		return StyleLog.Render("[WARN] " + message)
	case "ERROR":
		return StyleErrorMsg.Render("[ERROR] " + message)
	case "SUCCESS":
		return StyleStatusComplete.Render("[SUCCESS] " + message)
	default:
		return StyleLog.Render(message)
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlQ:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyRunes:
			switch msg.String() {
			case "q":
				m.quitting = true
				return m, tea.Quit

			case "r":
				if m.state != StateRunning {
					m.state = StateRunning
					return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
						return StateChangeMsg{State: StateRunning}
					})
				}

			case "p":
				if m.state == StateRunning {
					m.state = StatePaused
				} else if m.state == StatePaused {
					m.state = StateRunning
				}
				return m, nil

			case "l":
				if m.activeView == "status" {
					m.activeView = "logs"
				} else {
					m.activeView = "status"
				}
				return m, nil

			case "?":
				m.helpVisible = !m.helpVisible
				if m.helpVisible {
					m.activeView = "help"
				} else {
					m.activeView = "status"
				}
				return m, nil

			case "c":
				if m.activeView == "status" {
					m.activeView = "circuit"
				} else {
					m.activeView = "status"
				}
				return m, nil

			case "R":
				// Reset circuit breaker - send message to controller
				// For now, just show a log entry
				m.circuitState = "CLOSED"
				formattedLog := styledLogEntry("INFO", "Circuit breaker reset")
				m.logs = append(m.logs, formattedLog)
				return m, nil
			}
		}

	case LoopUpdateMsg:
		m.loopNumber = msg.LoopNumber
		m.callsUsed = msg.CallsUsed
		m.status = msg.Status
		return m, nil

	case LogMsg:
		formattedLog := styledLogEntry(msg.Level, msg.Message)
		m.logs = append(m.logs, formattedLog)
		if len(m.logs) > 500 {
			m.logs = m.logs[len(m.logs)-500:]
		}
		return m, nil

	case StateChangeMsg:
		m.state = msg.State
		return m, nil

	case StatusMsg:
		m.status = msg.Status
		return m, nil

	case tea.WindowSizeMsg:
		// Handle window resize (to be implemented)
		return m, nil

	case TickMsg:
		// Increment tick for animations
		m.tick++
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return TickMsg(t)
		})
	}

	return m, nil
}

// View renders TUI
func (m Model) View() string {
	if m.quitting {
		return StyleInfoMsg.Render("\nGoodbye!\n")
	}

	// Show error view if there's an error
	if m.err != nil && m.activeView != "help" {
		return m.renderErrorView()
	}

	var content string

	switch m.activeView {
	case "status":
		content = m.renderStatusView()
	case "logs":
		content = m.renderLogsView()
	case "help":
		content = m.renderHelpView()
	case "circuit":
		content = m.renderCircuitView()
	default:
		content = m.renderStatusView()
	}

	return content
}

func (m Model) renderRateLimitProgress() string {
	if m.maxCalls == 0 {
		return "Calls: 0/0"
	}

	total := m.maxCalls
	progress := float64(m.callsUsed) / float64(total)
	if progress > 1.0 {
		progress = 1.0
	}

	width := 20
	filled := int(progress * float64(width))

	emptyWidth := width - filled
	if emptyWidth < 0 {
		emptyWidth = 0
	}

	bar := fmt.Sprintf("Calls: %d/%d [%s%s]",
		m.callsUsed, total,
		StyleProgressBar.Render(strings.Repeat("█", filled)),
		strings.Repeat("░", emptyWidth),
	)

	return bar
}

func (m Model) renderStatusView() string {
	header := StyleHeader.Render("Ralph Codex")

	// Circuit state badge
	circuitState := "CLOSED"
	if m.circuitState != "" {
		circuitState = m.circuitState
	}

	var circuitBadge string
	switch circuitState {
	case "CLOSED":
		circuitBadge = StyleCircuitClosed.Render(circuitState)
	case "HALF_OPEN":
		circuitBadge = StyleCircuitHalfOpen.Render(circuitState)
	case "OPEN":
		circuitBadge = StyleCircuitOpen.Render(circuitState)
	default:
		circuitBadge = circuitState
	}

	// Progress bar for rate limit
	progressBar := m.renderRateLimitProgress()

	// Animated spinner if running
	spinner := ""
	if m.state == StateRunning {
		spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		spinner = fmt.Sprintf(" %s", StyleSpinner.Render(spinnerFrames[m.tick%len(spinnerFrames)]))
	}

	statusLine := StyleStatus.Render(fmt.Sprintf("Loop: %d%s | %s | Circuit: %s",
		m.loopNumber, spinner, progressBar, circuitBadge))

	statusDetail := fmt.Sprintf("\n%s", m.status)

	elapsed := time.Since(m.startTime).Round(time.Second)
	elapsedLine := StyleHelpDesc.Render(fmt.Sprintf("\nElapsed: %s", elapsed))

	keybindings := fmt.Sprintf(`
 %s %s Run/Restart loop
 %s %s Pause/Resume
 %s %s Toggle log view
 %s %s Show help
 %s %s Quit
 `,
		StyleHelpKey.Render("r"), StyleHelpDesc.Render("-"),
		StyleHelpKey.Render("p"), StyleHelpDesc.Render("-"),
		StyleHelpKey.Render("l"), StyleHelpDesc.Render("-"),
		StyleHelpKey.Render("?"), StyleHelpDesc.Render("-"),
		StyleHelpKey.Render("q"), StyleHelpDesc.Render("Quit (or Ctrl+C)"))

	return header + "\n" + statusLine + statusDetail + elapsedLine + keybindings
}

func (m Model) renderLogsView() string {
	header := StyleHeader.Render("Logs - Press 'l' to return")

	logContent := ""
	start := len(m.logs) - 30
	if start < 0 {
		start = 0
	}

	for i := start; i < len(m.logs); i++ {
		logContent += m.logs[i] + "\n"
	}

	if len(logContent) == 0 {
		logContent = StyleHelpDesc.Render("No logs yet\n")
	}

	return header + "\n" + StyleLog.Render(logContent)
}

func (m Model) renderErrorView() string {
	header := StyleHeader.Render("Ralph Codex - Error")

	errorMsg := StyleErrorMsg.Render(fmt.Sprintf("\nError: %v\n", m.err))

	helpText := StyleHelpDesc.Render(`
Press 'r' to retry
Press 'q' to quit
`)

	return header + errorMsg + helpText
}

func (m Model) renderCircuitView() string {
	header := StyleHeader.Render("Circuit Breaker Status")

	// Current state badge
	circuitState := "CLOSED"
	if m.circuitState != "" {
		circuitState = m.circuitState
	}

	var stateBadge string
	var stateDesc string

	switch circuitState {
	case "CLOSED":
		stateBadge = StyleCircuitClosed.Render(circuitState)
		stateDesc = "Circuit is operational. Normal loop execution is allowed."
	case "HALF_OPEN":
		stateBadge = StyleCircuitHalfOpen.Render(circuitState)
		stateDesc = "Circuit is monitoring. Loop may be paused if no progress continues."
	case "OPEN":
		stateBadge = StyleCircuitOpen.Render(circuitState)
		stateDesc = "Circuit is open! Loop execution is halted due to repeated failures."
	default:
		stateBadge = circuitState
		stateDesc = "Unknown circuit state."
	}

	// Circuit info
	circuitInfo := fmt.Sprintf(`
%s

%s
%s

State Explanation:
  %s

Keybindings:
  %s Return to status
  %s Reset circuit breaker (if stuck)
`,
		StyleInfoMsg.Render("Current State"),
		stateBadge,
		StyleDivider.Render(DividerChar),
		StyleHelpDesc.Render(stateDesc),
		StyleHelpKey.Render("Esc / l"),
		StyleHelpKey.Render("R"))

	return header + circuitInfo
}
