package tui

import (
	"regexp"
	"strings"
	"testing"

	"github.com/brainwhocodes/lisa-loop/internal/loop"
	"github.com/brainwhocodes/lisa-loop/internal/tui/msg"
	tea "github.com/charmbracelet/bubbletea"
)

func TestViewToggles_TasksAndOutput(t *testing.T) {
	model := Model{
		activeView: "status",
		viewMode:   ViewModeSplit,
	}

	// Toggle tasks
	msg := teaKey('t')
	newModel, _ := model.Update(msg)
	if newModel.(Model).viewMode != ViewModeTasks {
		t.Fatalf("expected tasks view, got %q", newModel.(Model).viewMode)
	}
	newModel, _ = newModel.(Model).Update(msg)
	if newModel.(Model).viewMode != ViewModeSplit {
		t.Fatalf("expected split view after second toggle, got %q", newModel.(Model).viewMode)
	}

	// Toggle output
	msg = teaKey('o')
	newModel, _ = model.Update(msg)
	if newModel.(Model).viewMode != ViewModeOutput {
		t.Fatalf("expected output view, got %q", newModel.(Model).viewMode)
	}
	newModel, _ = newModel.(Model).Update(msg)
	if newModel.(Model).viewMode != ViewModeSplit {
		t.Fatalf("expected split view after second toggle, got %q", newModel.(Model).viewMode)
	}
}

func TestRenderStableSections_HeaderFooterStatus(t *testing.T) {
	model := Model{
		state:        StatePaused,
		projectMode:  loop.ModeRefactor,
		loopNumber:   2,
		width:        80,
		height:       24,
		circuitState: "CLOSED",
	}

	header := stripANSI(model.renderHeader(model.width))
	if !strings.Contains(header, "Charm") || !strings.Contains(header, "LISA") {
		t.Fatalf("header should contain brand text, got: %q", header)
	}
	if !strings.Contains(header, "loop 2") {
		t.Fatalf("header should contain loop metadata, got: %q", header)
	}

	footer := stripANSI(model.renderFooter(model.width))
	for _, want := range []string{"r", "run", "p", "pause", "q", "quit"} {
		if !strings.Contains(footer, want) {
			t.Fatalf("footer should contain %q, got: %q", want, footer)
		}
	}

	status := stripANSI(model.renderStatusBar(model.width))
	if !strings.Contains(status, "paused") {
		t.Fatalf("status bar should contain state text, got: %q", status)
	}
	if !strings.Contains(status, "circuit") {
		t.Fatalf("status bar should contain circuit label, got: %q", status)
	}
}

func TestOutputDedup_CumulativeAgentMessageReplacesLastLine(t *testing.T) {
	var m Model

	m.addOutputLine("I'll", "agent_message")
	if len(m.outputLines) != 1 {
		t.Fatalf("expected 1 output line, got %d", len(m.outputLines))
	}

	// SSE-style cumulative updates should replace, not append.
	m.addOutputLine("I'll continue", "agent_message")
	m.addOutputLine("I'll continue fixing", "agent_message")

	if len(m.outputLines) != 1 {
		t.Fatalf("expected 1 output line after cumulative updates, got %d: %#v", len(m.outputLines), m.outputLines)
	}
	if got := m.outputLines[0]; got != "I'll continue fixing" {
		t.Fatalf("expected final line to be the latest cumulative value, got %q", got)
	}
}

func TestOutputDedup_ExactDuplicateSkipped(t *testing.T) {
	var m Model
	m.addOutputLine("hello", "raw")
	m.addOutputLine("hello", "raw")
	if len(m.outputLines) != 1 {
		t.Fatalf("expected 1 output line after exact duplicate, got %d", len(m.outputLines))
	}
}

func TestReasoningDedup_ReplacesCumulativeReasoning(t *testing.T) {
	var m Model
	m.addReasoningLine("thinking...")
	m.addReasoningLine("thinking... more")
	m.addReasoningLine("thinking... more and more")

	if len(m.reasoningLines) != 1 {
		t.Fatalf("expected reasoning lines to be replaced (len=1), got %d: %#v", len(m.reasoningLines), m.reasoningLines)
	}
	if got := m.reasoningLines[0]; got != "thinking... more and more" {
		t.Fatalf("expected latest reasoning to win, got %q", got)
	}
}

func TestToolDedup_DoesNotAppendDuplicateToolCall(t *testing.T) {
	model := Model{}

	start := msg.CodexToolCallMsg{Tool: "read", Target: "file.txt", Status: "started"}
	newModel, _ := model.Update(start)
	newModel, _ = newModel.(Model).Update(start) // duplicate

	if got := len(newModel.(Model).outputLines); got != 1 {
		t.Fatalf("expected 1 tool line after duplicate, got %d: %#v", got, newModel.(Model).outputLines)
	}
}

func teaKey(r rune) tea.Msg {
	// Keep this helper tiny and local to tests.
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

var ansiRE = regexp.MustCompile("\x1b\\[[0-9;]*[A-Za-z]")

func stripANSI(s string) string {
	return ansiRE.ReplaceAllString(s, "")
}
