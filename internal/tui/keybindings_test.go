package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
)

// TestGetKeybindingHelp tests keybinding help generation
func TestGetKeybindingHelp(t *testing.T) {
	help := GetKeybindingHelp()

	// Check that all main sections are present
	expectedSections := []string{
		"Navigation",
		"Loop Control",
		"Views",
		"CLI Options",
		"Project Options",
		"Rate Limiting",
		"Project Commands",
		"Troubleshooting",
	}

	for _, section := range expectedSections {
		if !strings.Contains(help, section) {
			t.Errorf("Help should contain section '%s'", section)
		}
	}

	// Check that key keybindings are present
	expectedKeys := []string{
		"q",
		"?",
		"r",
		"p",
		"l",
		"c",
		"R",
		"--monitor",
		"--verbose",
		"--backend",
	}

	for _, key := range expectedKeys {
		if !strings.Contains(help, key) {
			t.Errorf("Help should contain key '%s'", key)
		}
	}
}

// TestModelToggleCircuitView tests circuit view toggle
func TestModelToggleCircuitView(t *testing.T) {
	model := Model{
		activeView: "status",
	}

	// Toggle to circuit view
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	newModel, _ := model.Update(msg)

	if newModel.(Model).activeView != "circuit" {
		t.Errorf("Expected circuit view, got %s", newModel.(Model).activeView)
	}

	// Toggle back to status
	newModel, _ = newModel.(Model).Update(msg)

	if newModel.(Model).activeView != "status" {
		t.Errorf("Expected status view, got %s", newModel.(Model).activeView)
	}
}

// TestModelResetCircuit tests circuit breaker reset
func TestModelResetCircuit(t *testing.T) {
	model := Model{
		circuitState: "OPEN",
		logs:         []string{},
	}

	// Reset circuit breaker
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'R'}}
	newModel, _ := model.Update(msg)

	if newModel.(Model).circuitState != "CLOSED" {
		t.Errorf("Expected CLOSED state, got %s", newModel.(Model).circuitState)
	}

	logs := newModel.(Model).logs
	if len(logs) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(logs))
	}

	if !contains(logs[0], "Circuit breaker reset") {
		t.Errorf("Log should contain reset message, got: %s", logs[0])
	}
}

// TestRenderCircuitView tests circuit breaker view rendering
func TestRenderCircuitView(t *testing.T) {
	tests := []struct {
		name          string
		circuitState  string
		shouldContain string
	}{
		{
			name:          "CLOSED state",
			circuitState:  "CLOSED",
			shouldContain: "operational",
		},
		{
			name:          "HALF_OPEN state",
			circuitState:  "HALF_OPEN",
			shouldContain: "monitoring",
		},
		{
			name:          "OPEN state",
			circuitState:  "OPEN",
			shouldContain: "halted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := Model{
				circuitState: tt.circuitState,
			}

			result := model.renderCircuitView()

			if !contains(result, "Circuit Breaker Status") {
				t.Error("View should contain header")
			}

			if !contains(result, tt.circuitState) {
				t.Errorf("View should contain state %s", tt.circuitState)
			}

			if !contains(result, tt.shouldContain) {
				t.Errorf("View should contain '%s'", tt.shouldContain)
			}
		})
	}
}

// TestRenderCircuitViewUnknownState tests unknown circuit state
func TestRenderCircuitViewUnknownState(t *testing.T) {
	model := Model{
		circuitState: "UNKNOWN",
	}

	result := model.renderCircuitView()

	if !contains(result, "UNKNOWN") {
		t.Error("View should contain UNKNOWN state")
	}

	if !contains(result, "Unknown circuit state") {
		t.Error("View should contain unknown state description")
	}
}
