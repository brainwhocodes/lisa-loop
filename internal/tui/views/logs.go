package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	maxLogLines  = 500
	viewHeight   = 30
	scrollAmount = 10
)

// LogsViewModel manages scrollable log view
type LogsViewModel struct {
	logs      []string
	scrollPos int
	maxLines  int
}

// NewLogsViewModel creates a new logs view model
func NewLogsViewModel() *LogsViewModel {
	return &LogsViewModel{
		logs:      make([]string, 0, maxLogLines),
		scrollPos: 0,
		maxLines:  viewHeight,
	}
}

// View renders logs view with scrolling
func (m *LogsViewModel) View() string {
	if len(m.logs) == 0 {
		return "No logs available yet..."
	}

	// Calculate visible range
	start := len(m.logs) - m.maxLines - m.scrollPos
	if start < 0 {
		start = 0
	}

	end := start + m.maxLines
	if end > len(m.logs) {
		end = len(m.logs)
	}

	var builder strings.Builder
	for i := start; i < end; i++ {
		builder.WriteString(m.logs[i])
		builder.WriteString("\n")
	}

	return builder.String()
}

// AddLog adds a new log entry
func (m *LogsViewModel) AddLog(log string) {
	if len(m.logs) >= maxLogLines {
		m.logs = m.logs[1:]
	}
	m.logs = append(m.logs, log)
	m.ScrollToBottom()
}

// ScrollToBottom scrolls to most recent logs
func (m *LogsViewModel) ScrollToBottom() {
	m.scrollPos = 0
}

// ScrollUp scrolls up
func (m *LogsViewModel) ScrollUp() {
	m.scrollPos += scrollAmount
	if m.scrollPos > len(m.logs)-m.maxLines {
		m.scrollPos = len(m.logs) - m.maxLines
	}
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}
}

// ScrollDown scrolls down
func (m *LogsViewModel) ScrollDown() {
	m.scrollPos -= scrollAmount
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}
}

// PageUp scrolls up half a page
func (m *LogsViewModel) PageUp() {
	m.scrollPos += scrollAmount / 2
	if m.scrollPos > len(m.logs)-m.maxLines {
		m.scrollPos = len(m.logs) - m.maxLines
	}
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}
}

// PageDown scrolls down half a page
func (m *LogsViewModel) PageDown() {
	m.scrollPos -= scrollAmount / 2
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}
}

// Clear clears all logs
func (m *LogsViewModel) Clear() {
	m.logs = make([]string, 0, maxLogLines)
	m.scrollPos = 0
}

// GetLogs returns all logs
func (m *LogsViewModel) GetLogs() []string {
	return m.logs
}

// GetScrollPos returns current scroll position
func (m *LogsViewModel) GetScrollPos() int {
	return m.scrollPos
}

// SetMaxLines sets maximum number of lines to display
func (m *LogsViewModel) SetMaxLines(lines int) {
	m.maxLines = lines
}

// StyledLogEntry returns a styled log entry
func StyledLogEntry(level string, message string) string {
	switch level {
	case "INFO":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Render("[INFO] " + message)
	case "WARN":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Render("[WARN] " + message)
	case "ERROR":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Render("[ERROR] " + message)
	case "SUCCESS":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Render("[SUCCESS] " + message)
	default:
		return message
	}
}

// GetTotalLogCount returns total number of logs
func (m *LogsViewModel) GetTotalLogCount() int {
	return len(m.logs)
}

// GetVisibleRange returns range of visible log indices
func (m *LogsViewModel) GetVisibleRange() (start, end int) {
	start = len(m.logs) - m.maxLines - m.scrollPos
	if start < 0 {
		start = 0
	}

	end = start + m.maxLines
	if end > len(m.logs) {
		end = len(m.logs)
	}

	return start, end
}
