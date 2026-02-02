# refactor(tui/view): extract reusable render components + layout helpers

## Code changes
- [ ] Create `internal/tui/view` package for pure render helpers (no IO, no controller refs).
- [ ] Extract and reuse:
  - [ ] header renderer (brand + diagonals + metadata)
  - [ ] status bar renderer
  - [ ] footer renderer
  - [ ] common panel layout helpers (pad, clamp sizes, divider)
- [ ] Keep `Model.View()` readable:
  - [ ] delegate to screen view functions that call `view` helpers
- [ ] Avoid introducing per-frame `lipgloss.NewStyle()` calls in hot paths.

## Tests to add/update
- [ ] Add snapshot-ish tests for header/footer/status:
  - [ ] fixed width/height
  - [ ] deterministic model (tick=0)
  - [ ] assert key substrings; avoid full-string brittle snapshots unless necessary

## Manual UX verification
- [ ] `go test ./...`
- [ ] Run TUI, confirm:
  - [ ] no spacing regressions in header/status/footer
  - [ ] split pane still renders correctly
  - [ ] resize behavior remains acceptable

