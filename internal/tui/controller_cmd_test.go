package tui

import (
	"context"
	"testing"

	"github.com/brainwhocodes/lisa-loop/internal/loop"
	tuimsg "github.com/brainwhocodes/lisa-loop/internal/tui/msg"
	tea "github.com/charmbracelet/bubbletea"
)

type fakeController struct {
	runCalled bool
	runErr    error
}

func (f *fakeController) Run(ctx context.Context) error {
	f.runCalled = true
	return f.runErr
}

func (f *fakeController) Pause()                              {}
func (f *fakeController) Resume()                             {}
func (f *fakeController) SetEventCallback(loop.EventCallback) {}

func TestRunStartsControllerViaCmdAndCompletesOnDoneMsg(t *testing.T) {
	fc := &fakeController{}
	model := Model{
		state:      StateInitializing,
		controller: fc,
	}

	// Start
	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if cmd == nil {
		t.Fatalf("expected a command on 'r'")
	}
	if newModel.(Model).state != StateRunning {
		t.Fatalf("expected StateRunning after 'r', got %v", newModel.(Model).state)
	}

	// Execute the cmd (simulates Bubble Tea running the command)
	done := cmd()
	doneMsg, ok := done.(tuimsg.ControllerDoneMsg)
	if !ok {
		t.Fatalf("expected ControllerDoneMsg, got %T", done)
	}
	if doneMsg.Err != nil {
		t.Fatalf("expected nil error, got %v", doneMsg.Err)
	}
	if !fc.runCalled {
		t.Fatalf("expected controller Run() to be called")
	}

	// Apply completion message
	afterDone, _ := newModel.(Model).Update(doneMsg)
	if afterDone.(Model).state != StateComplete {
		t.Fatalf("expected StateComplete after ControllerDoneMsg, got %v", afterDone.(Model).state)
	}
}

func TestAutoStartMsgStartsController(t *testing.T) {
	fc := &fakeController{}
	model := Model{
		state:      StateInitializing,
		controller: fc,
	}

	newModel, cmd := model.Update(tuimsg.AutoStartLoopMsg{})
	if cmd == nil {
		t.Fatalf("expected a command on AutoStartLoopMsg")
	}
	if newModel.(Model).state != StateRunning {
		t.Fatalf("expected StateRunning after auto-start, got %v", newModel.(Model).state)
	}
}

func TestQuitCancelsControllerContext(t *testing.T) {
	cancelled := false
	model := Model{
		cancel: func() { cancelled = true },
		ctx:    context.Background(),
	}

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if !newModel.(Model).quitting {
		t.Fatalf("expected quitting=true")
	}
	if !cancelled {
		t.Fatalf("expected cancel to be called on quit")
	}
}

func TestRestartCancelsPreviousRun(t *testing.T) {
	cancelled := false
	fc := &fakeController{}
	model := Model{
		state:      StatePaused,
		controller: fc,
		cancel:     func() { cancelled = true },
		ctx:        context.Background(),
	}

	newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	if cmd == nil {
		t.Fatalf("expected a command on restart")
	}
	if !cancelled {
		t.Fatalf("expected previous run to be cancelled on restart")
	}
	if newModel.(Model).state != StateRunning {
		t.Fatalf("expected StateRunning after restart, got %v", newModel.(Model).state)
	}
}
