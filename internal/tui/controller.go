package tui

import (
	"context"

	"github.com/brainwhocodes/lisa-loop/internal/loop"
)

// Controller is the minimal surface the TUI needs from the loop controller.
// Using an interface keeps the UI testable without constructing a real loop.Controller.
type Controller interface {
	Run(ctx context.Context) error
	Pause()
	Resume()
	SetEventCallback(cb loop.EventCallback)
}
