package tui

import "github.com/brainwhocodes/lisa-loop/internal/tui/plan"

// Task/Phase are re-exported for backwards-compat within the tui package.
// Parsing lives in internal/tui/plan; the UI owns the Active flag.
type Task = plan.Task
type Phase = plan.Phase
