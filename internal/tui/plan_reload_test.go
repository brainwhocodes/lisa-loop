package tui

import (
	"errors"
	"testing"

	"github.com/brainwhocodes/lisa-loop/internal/loop"
	tuimsg "github.com/brainwhocodes/lisa-loop/internal/tui/msg"
)

func TestPreflightTriggersPlanReloadViaCmdAndMergesCompletion(t *testing.T) {
	model := Model{
		projectMode: loop.ModeRefactor,
		planFile:    "plan.md",
		readFile: func(path string) ([]byte, error) {
			return []byte("## Phase 1: X\n- [ ] Task A\n- [ ] Task B\n"), nil
		},
		tasks: []Task{
			{Text: "Task A", Completed: true},
		},
	}

	event := loop.LoopEvent{
		Type: loop.EventTypePreflight,
		Preflight: &loop.PreflightSummary{
			PlanFile:       "plan.md",
			TotalTasks:     2,
			RemainingCount: 2,
		},
	}

	newModel, cmd := model.Update(tuimsg.ControllerEventMsg{Event: event})
	if cmd == nil {
		t.Fatalf("expected a plan reload command on preflight")
	}

	loaded := cmd()
	pl, ok := loaded.(tuimsg.PlanLoadedMsg)
	if !ok {
		t.Fatalf("expected PlanLoadedMsg, got %T", loaded)
	}
	if pl.Err != nil {
		t.Fatalf("expected nil error, got %v", pl.Err)
	}

	afterLoad, _ := newModel.(Model).Update(pl)
	m := afterLoad.(Model)

	// Task A should remain completed even if the plan file hasn't been updated yet.
	foundA := false
	for _, task := range m.tasks {
		if task.Text == "Task A" {
			foundA = true
			if !task.Completed {
				t.Fatalf("expected Task A to remain completed after reload")
			}
		}
	}
	if !foundA {
		t.Fatalf("expected Task A to be present after reload")
	}

	if m.planFile != "plan.md" {
		t.Fatalf("expected planFile to be updated, got %q", m.planFile)
	}
}

func TestPlanReloadErrorDoesNotWipeTasks(t *testing.T) {
	readErr := errors.New("boom")
	model := Model{
		projectMode: loop.ModeRefactor,
		planFile:    "plan.md",
		readFile: func(path string) ([]byte, error) {
			return nil, readErr
		},
		tasks: []Task{
			{Text: "Task A", Completed: true},
		},
	}

	event := loop.LoopEvent{
		Type: loop.EventTypePreflight,
		Preflight: &loop.PreflightSummary{
			PlanFile:       "plan.md",
			TotalTasks:     1,
			RemainingCount: 0,
		},
	}

	newModel, cmd := model.Update(tuimsg.ControllerEventMsg{Event: event})
	if cmd == nil {
		t.Fatalf("expected a plan reload command on preflight")
	}

	loaded := cmd()
	pl, ok := loaded.(tuimsg.PlanLoadedMsg)
	if !ok {
		t.Fatalf("expected PlanLoadedMsg, got %T", loaded)
	}
	if pl.Err == nil {
		t.Fatalf("expected error, got nil")
	}

	afterLoad, _ := newModel.(Model).Update(pl)
	m := afterLoad.(Model)
	if len(m.tasks) != 1 || m.tasks[0].Text != "Task A" {
		t.Fatalf("expected tasks to be preserved on reload error, got %#v", m.tasks)
	}

	// Warning should be logged.
	foundWarn := false
	for _, l := range m.logs {
		if contains(l, "Could not reload tasks") {
			foundWarn = true
			break
		}
	}
	if !foundWarn {
		t.Fatalf("expected reload warning log to be present")
	}
}
