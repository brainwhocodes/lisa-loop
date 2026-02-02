# refactor(tui/effects): run controller via tea.Cmd (no goroutine mutation)

## Code changes
- [ ] Create `internal/tui/effects` package for `tea.Cmd` builders.
- [ ] Replace `go m.runController()` with `effects.RunController(...) tea.Cmd`.
- [ ] Introduce explicit messages for controller lifecycle:
  - [ ] `msg.ControllerStarted{...}` (optional)
  - [ ] `msg.ControllerDone{Err error}` (required)
- [ ] Ensure `Update` is the only place that mutates `Model` fields (`state`, `err`, etc.).
- [ ] Preserve context cancellation semantics:
  - [ ] `r` creates a context (or uses one owned by model/root)
  - [ ] quitting cancels the context (if currently running)

## Tests to add/update
- [ ] Add a unit test asserting:
  - [ ] pressing `r` returns a non-nil command
  - [ ] controller completion message transitions state to complete/error as before
- [ ] If possible, add a race-test locally (optional): `go test -race ./internal/tui/...`

## Manual UX verification
- [ ] `go test ./...`
- [ ] Start loop: `r` (ensure UI remains responsive).
- [ ] Pause/resume: `p` toggles; controller receives pause/resume (no blocking).
- [ ] Quit: `q` exits cleanly.

