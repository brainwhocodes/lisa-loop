# refactor(tui/msg): centralize tea.Msg taxonomy in internal/tui/msg

## Code changes
- [ ] Create `internal/tui/msg` package.
- [ ] Move all existing msg types into `msg`:
  - [ ] types currently in `internal/tui/messages.go`
  - [ ] types currently in `internal/tui/model.go` (LoopUpdateMsg, LogMsg, etc.)
- [ ] Update imports and switch statements to use `msg.X` types.
- [ ] Avoid import cycles:
  - [ ] `msg` must not import `tui` or `view/style` packages.
- [ ] Keep message payloads explicit (no stringly-typed variants) where it improves clarity without behavior change.

## Tests to add/update
- [ ] Update existing tests to reference moved types (compile + pass).
- [ ] Add a tiny compile-time sanity test (optional) that key message types exist in `msg` and are used by `tui.Model.Update`.

## Manual UX verification
- [ ] `go test ./...`
- [ ] Run `lisa --monitor` and verify controller events still update output/logs/status.

