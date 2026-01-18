package tui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
)

// TestModelTick tests that tick counter increments properly
func TestModelTick(t *testing.T) {
	model := Model{
		state:     StateRunning,
		startTime: time.Now(),
	}

	// Initialize model
	cmd := model.Init()
	if cmd == nil {
		t.Fatal("Init should return a command")
	}

	// Send a few tick messages
	for i := 0; i < 5; i++ {
		msg := TickMsg(time.Now())
		newModel, newCmd := model.Update(msg)
		if newCmd == nil {
			t.Fatal("Tick should return a command")
		}

		// Check tick was incremented
		if newModel.(Model).tick != i+1 {
			t.Errorf("Expected tick %d, got %d", i+1, newModel.(Model).tick)
		}

		model = newModel.(Model)
	}
}

// TestModelInit tests model initialization
func TestModelInit(t *testing.T) {
	model := Model{
		state:     StateInitializing,
		startTime: time.Now(),
	}

	cmd := model.Init()

	if cmd == nil {
		t.Error("Init should return a tick command")
	}

	if model.startTime.IsZero() {
		t.Error("Start time should be set")
	}
}

// TestModelUpdateQuit tests quit keybinding
func TestModelQuit(t *testing.T) {
	model := Model{
		state: StateRunning,
	}

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	newModel, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Quit should return Quit command")
	}

	if !newModel.(Model).quitting {
		t.Error("Quitting flag should be set")
	}

	// Test 'q' key
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, cmd = model.Update(msg)

	if cmd == nil {
		t.Error("Quit should return Quit command")
	}

	if !newModel.(Model).quitting {
		t.Error("Quitting flag should be set")
	}
}

// TestModelTogglePause tests pause/resume toggle
func TestModelTogglePause(t *testing.T) {
	model := Model{
		state: StateRunning,
	}

	// Pause
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
	newModel, _ := model.Update(msg)

	if newModel.(Model).state != StatePaused {
		t.Error("State should be paused")
	}

	// Resume
	newModel, _ = newModel.(Model).Update(msg)

	if newModel.(Model).state != StateRunning {
		t.Error("State should be running")
	}
}

// TestModelToggleView tests view toggling
func TestModelToggleView(t *testing.T) {
	model := Model{
		activeView: "status",
	}

	// Toggle to logs
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	newModel, _ := model.Update(msg)

	if newModel.(Model).activeView != "logs" {
		t.Errorf("Expected logs view, got %s", newModel.(Model).activeView)
	}

	// Toggle back to status
	newModel, _ = newModel.(Model).Update(msg)

	if newModel.(Model).activeView != "status" {
		t.Errorf("Expected status view, got %s", newModel.(Model).activeView)
	}
}

// TestModelToggleHelp tests help toggle
func TestModelToggleHelp(t *testing.T) {
	model := Model{
		activeView:  "status",
		helpVisible: false,
	}

	// Show help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, _ := model.Update(msg)

	if !newModel.(Model).helpVisible {
		t.Error("Help should be visible")
	}

	if newModel.(Model).activeView != "help" {
		t.Errorf("Expected help view, got %s", newModel.(Model).activeView)
	}

	// Hide help
	newModel, _ = newModel.(Model).Update(msg)

	if newModel.(Model).helpVisible {
		t.Error("Help should be hidden")
	}

	if newModel.(Model).activeView != "status" {
		t.Errorf("Expected status view, got %s", newModel.(Model).activeView)
	}
}

// TestModelLogMsg tests log message handling
func TestModelLogMsg(t *testing.T) {
	model := Model{
		logs: []string{},
	}

	msg := LogMsg{
		Message: "Test message",
		Level:   "INFO",
	}

	newModel, _ := model.Update(msg)

	logs := newModel.(Model).logs
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	if !contains(logs[0], "Test message") {
		t.Errorf("Log should contain message, got: %s", logs[0])
	}
}

// TestModelStateString tests State.String() method
func TestStateString(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateInitializing, "Initializing"},
		{StateRunning, "Running"},
		{StatePaused, "Paused"},
		{StateComplete, "Complete"},
		{StateError, "Error"},
		{State(99), "Unknown"},
	}

	for _, test := range tests {
		result := test.state.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

// TestRenderRateLimitProgress tests progress bar rendering
func TestRenderRateLimitProgress(t *testing.T) {
	tests := []struct {
		name          string
		callsUsed     int
		maxCalls      int
		shouldContain string
	}{
		{"zero calls", 0, 100, "Calls: 0/100"},
		{"half calls", 50, 100, "Calls: 50/100"},
		{"max calls", 100, 100, "Calls: 100/100"},
		{"over limit", 150, 100, "Calls: 150/100"},
		{"zero max", 10, 0, "Calls: 0/0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			model := Model{
				callsUsed: test.callsUsed,
				maxCalls:  test.maxCalls,
			}

			result := model.renderRateLimitProgress()

			if !contains(result, test.shouldContain) {
				t.Errorf("Expected to contain %s, got: %s", test.shouldContain, result)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && contains(s[1:], substr) || s[0:len(substr)] == substr)
}
