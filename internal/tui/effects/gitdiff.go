package effects

import (
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/brainwhocodes/lisa-loop/internal/tui/msg"
)

// Exec abstracts command execution for testability.
type Exec func(name string, args ...string) ([]byte, error)

// OSExec is the production command runner.
var OSExec Exec = func(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

// LoadGitDiff runs git to collect changed files and a unified patch.
// It is safe to run even when no changes exist (returns empty output).
func LoadGitDiff(seq int, execFn Exec) tea.Cmd {
	return func() tea.Msg {
		if execFn == nil {
			execFn = OSExec
		}

		nameStatus, nsErr := execFn("git", "diff", "--name-status")
		patch, pErr := execFn("git", "diff", "--patch", "--no-color")

		err := firstErr(nsErr, pErr)
		return msg.GitDiffLoadedMsg{
			Seq:        seq,
			NameStatus: strings.TrimSpace(string(nameStatus)),
			Patch:      strings.TrimRight(string(patch), "\n"),
			Err:        err,
			At:         time.Now(),
		}
	}
}

// DebounceDiff emits msg.DiffDebounceFiredMsg after the given delay.
func DebounceDiff(seq int, d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return msg.DiffDebounceFiredMsg{Seq: seq}
	})
}

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
