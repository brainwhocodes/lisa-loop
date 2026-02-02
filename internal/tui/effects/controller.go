package effects

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/brainwhocodes/lisa-loop/internal/tui/msg"
)

// Runner is the minimal dependency needed to run the loop controller.
// It matches loop.Controller.Run and enables fakes in tests.
type Runner interface {
	Run(ctx context.Context) error
}

// RunController runs the controller in a Bubble Tea command and reports completion via msg.ControllerDoneMsg.
// The controller's granular progress is still delivered via msg.ControllerEventMsg callbacks.
func RunController(ctx context.Context, runner Runner) tea.Cmd {
	return func() tea.Msg {
		if runner == nil {
			return msg.ControllerDoneMsg{Err: nil}
		}
		if ctx == nil {
			ctx = context.Background()
		}
		return msg.ControllerDoneMsg{Err: runner.Run(ctx)}
	}
}
