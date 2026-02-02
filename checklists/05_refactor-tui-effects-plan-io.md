# refactor(tui/effects): move filesystem IO out of Update (reload tasks via Cmd)

## Code changes
- [ ] Replace direct `os.ReadFile` usage from `Update` paths with `effects.LoadPlan(...) tea.Cmd`.
- [ ] Add injection for tests:
  - [ ] `type ReadFile func(path string) ([]byte, error)` passed into effects, or a small `FS` interface.
- [ ] Convert `Model.reloadTasks()` to:
  - [ ] schedule a command
  - [ ] apply result in `Update` on `msg.PlanLoaded{...}`
- [ ] Preserve current merge behavior:
  - [ ] tasks completed in-memory remain completed even if file hasn't been updated yet.

## Tests to add/update
- [ ] Add unit tests using a fake `ReadFile`:
  - [ ] reload triggered by preflight msg
  - [ ] merge keeps completed tasks
  - [ ] error reading plan file logs a warning and preserves current tasks

## Manual UX verification
- [ ] `go test ./...`
- [ ] While TUI is open, edit the plan file; verify tasks refresh on next preflight event.

