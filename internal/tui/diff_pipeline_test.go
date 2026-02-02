package tui

import (
	"strings"
	"testing"

	"github.com/brainwhocodes/lisa-loop/internal/loop"
	tuimsg "github.com/brainwhocodes/lisa-loop/internal/tui/msg"
)

func TestDiffPipeline_OutcomeTriggersDebounceThenLoadsGitDiff(t *testing.T) {
	model := Model{
		exec: func(name string, args ...string) ([]byte, error) {
			if name != "git" {
				t.Fatalf("unexpected exec name: %s", name)
			}
			a := strings.Join(args, " ")
			if strings.Contains(a, "--name-status") {
				return []byte("M\tinternal/tui/model.go\n"), nil
			}
			if strings.Contains(a, "--patch") {
				return []byte("diff --git a/x b/x\n+hi\n"), nil
			}
			return []byte{}, nil
		},
	}

	evt := loop.LoopEvent{
		Type: loop.EventTypeOutcome,
		Outcome: &loop.LoopOutcome{
			Success:       true,
			FilesModified: 1,
		},
	}

	m1, cmd := model.Update(tuimsg.ControllerEventMsg{Event: evt})
	if cmd == nil {
		t.Fatalf("expected debounce cmd")
	}
	if m1.(Model).diffSeq != 1 || !m1.(Model).diffPending {
		t.Fatalf("expected diffSeq=1 and diffPending=true, got seq=%d pending=%v", m1.(Model).diffSeq, m1.(Model).diffPending)
	}

	// Simulate debounce firing.
	m2, loadCmd := m1.(Model).Update(tuimsg.DiffDebounceFiredMsg{Seq: 1})
	if loadCmd == nil {
		t.Fatalf("expected load git diff cmd")
	}

	msg := loadCmd()
	loaded, ok := msg.(tuimsg.GitDiffLoadedMsg)
	if !ok {
		t.Fatalf("expected GitDiffLoadedMsg, got %T", msg)
	}
	m3, _ := m2.(Model).Update(loaded)

	if m3.(Model).diffPending {
		t.Fatalf("expected diffPending=false after load")
	}
	if m3.(Model).gitDiffPatch == "" || m3.(Model).gitDiffNameStatus == "" {
		t.Fatalf("expected git diff fields to be populated")
	}
}

func TestDiffPipeline_IgnoresStaleSeq(t *testing.T) {
	model := Model{diffSeq: 2}
	m1, cmd := model.Update(tuimsg.DiffDebounceFiredMsg{Seq: 1})
	if cmd != nil {
		t.Fatalf("expected nil cmd for stale seq")
	}
	if m1.(Model).diffSeq != 2 {
		t.Fatalf("expected diffSeq unchanged")
	}
}
