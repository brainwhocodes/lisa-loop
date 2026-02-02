# refactor(tui/tests): add baseline invariants for routing + stable view sections

## Code changes
- [ ] Add small helpers for tests only (if needed) to construct a `tui.Model` with deterministic fields (tick=0, width fixed).
- [ ] Avoid changing production behavior; if a helper must be exported, stop and reconsider.

## Tests to add/update
- [ ] Add routing invariants:
  - [ ] `l/t/o/?/c` toggle behavior (enter view, exit back to split).
  - [ ] `q` and `Ctrl+C` quit.
- [ ] Add stable render tests (fixed width/height):
  - [ ] header contains brand + mode + loop metadata.
  - [ ] footer contains the expected key list.
  - [ ] status bar contains state + circuit label (when set).
- [ ] Add reducer-level tests for:
  - [ ] cumulative agent message replacement logic
  - [ ] cumulative reasoning replacement logic
  - [ ] tool call dedup (same tool/target/status does not double-add)

## Manual UX verification
- [ ] `go test ./...`
- [ ] `lisa --monitor` launches and exits cleanly (`q`).
- [ ] Toggle each view twice: `l`, `t`, `o`, `?`, `c` (returns to split).
- [ ] Press `r` (start) then `p` (pause) then `p` (resume) with a controller wired (if available).

