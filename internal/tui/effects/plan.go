package effects

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/brainwhocodes/lisa-loop/internal/tui/msg"
	"github.com/brainwhocodes/lisa-loop/internal/tui/plan"
)

// ReadFile abstracts filesystem reads for testability.
type ReadFile func(path string) ([]byte, error)

// OSReadFile is the production file reader.
var OSReadFile ReadFile = os.ReadFile

// LoadPlan reads and parses a plan file in a Bubble Tea command.
func LoadPlan(filename string, readFile ReadFile) tea.Cmd {
	return func() tea.Msg {
		if readFile == nil {
			readFile = OSReadFile
		}
		b, err := readFile(filename)
		if err != nil {
			return msg.PlanLoadedMsg{Filename: filename, Err: err}
		}
		content := string(b)
		return msg.PlanLoadedMsg{
			Filename: filename,
			Phases:   plan.ParsePhases(content),
			Tasks:    plan.ParseTasks(content),
			Err:      nil,
		}
	}
}
