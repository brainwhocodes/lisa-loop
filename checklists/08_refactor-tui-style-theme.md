# refactor(tui/style): centralize theme/styles and avoid per-frame style construction

## Code changes
- [ ] Create `internal/tui/style` package:
  - [ ] `Theme` struct holding semantic styles and layout constants
  - [ ] `DefaultTheme()` constructs all styles once
- [ ] Root model stores a theme instance; view helpers consume it.
- [ ] Migrate from package-level style vars:
  - [ ] either replace them outright, or keep thin compatibility wrappers during transition
- [ ] Ensure no style creation happens per frame in `View()` code paths (except trivial, justified cases).

## Tests to add/update
- [ ] Theme construction test:
  - [ ] verify non-zero styles/constants are set
  - [ ] verify palette matches current defaults (if intended)
- [ ] Ensure existing header/footer/status tests remain passing with no output diffs.

## Manual UX verification
- [ ] `go test ./...`
- [ ] Run TUI and confirm colors, padding, and borders match previous behavior.

