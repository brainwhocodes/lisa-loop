# refactor(tui/state): make screen state explicit and remove redundant flags

## Code changes
- [ ] Introduce a typed screen enum (e.g., `type Screen int`) with values:
  - [ ] Split, Tasks, Output, Logs, Help, Circuit, Error
- [ ] Make routing a single source of truth:
  - [ ] remove or deprecate `activeView` (string)
  - [ ] remove or deprecate `helpVisible` (bool)
  - [ ] keep `viewMode` only if it becomes the enum, not a second source of truth
- [ ] Centralize toggle logic in helpers (ideally pure functions).

## Tests to add/update
- [ ] Update routing tests from commit 01 to assert screen enum transitions.
- [ ] Add test: error state does not break returning from help (match current behavior).

## Manual UX verification
- [ ] `go test ./...`
- [ ] Toggle all views; verify return-to-split behavior remains identical.
- [ ] Verify error view still appears when `m.err != nil` and not in help view.

