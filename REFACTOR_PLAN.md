# TUI Refactor Plan (Bubble Tea + Lipgloss)

You are refactoring an existing Bubble Tea + Lip Gloss TUI. The goal is to improve structure, readability, and testability while preserving behavior and UX.

Non-negotiables:
- No behavior regressions unless explicitly called out as a bug fix with before/after evidence.
- Keep the public surface stable (CLI flags, config, keybindings, screen flows, output formatting) unless required for correctness.
- Small, atomic changes (each commit reviewable + revertible).
- Build stays green: `go test ./...` passes after every commit.
- No big-bang rewrites. Incremental refactors only.
- Prefer stdlib; add deps only if clearly justified.

## Current State (Repo Audit)

Entry points / program creation:
- `cmd/lisa/main.go`: `runWithMonitor(...)` builds a loop `Controller`, then calls `tui.NewProgram(...)` and `program.Run()`.
- `internal/tui/program.go`: `Program.Run()` constructs `tea.NewProgram(p.model, tea.WithAltScreen(), tea.WithMouseCellMotion())`.
- The main Bubble Tea model is `internal/tui/model.go` (`type Model struct{ ... }`).
- Screens are currently represented by `viewMode` (`split/tasks/output/logs/help/circuit`) plus error handling via `m.err` and `StateError`.

Message flow:
- User input: `Model.Update(tea.KeyMsg)` handles:
  - quit (`q`, `Ctrl+C`, `Ctrl+Q`)
  - run (`r`) which currently spawns a goroutine (`go m.runController()`)
  - pause/resume (`p`) calls `m.controller.Pause()/Resume()`
  - view toggles (`l/t/o/?/c`) mutate `viewMode` and `activeView`
- Loop/controller input:
  - `internal/tui/program.go` sets `controller.SetEventCallback(...)` and uses `program.Send(ControllerEventMsg{...})`.
  - `Model.Update(ControllerEventMsg)` mutates UI state from loop events (status/logs/output/tool/analysis/preflight/outcome).
- Animation: `Model.Init()` schedules `tea.Tick(...)`; `TickMsg` increments `m.tick` and re-schedules.

Rendering:
- `Model.View()` chooses a renderer based on `viewMode`, then pads to full screen and applies a global background.
- Rendering helpers are spread across:
  - `internal/tui/views.go`: header/status bar/split panes/tasks/output/logs/preflight panels and helpers.
  - `internal/tui/model.go`: error view, circuit view, rate-limit progress, and padding/background wrapper.
  - `internal/tui/keybindings.go`: help screen renderer and keybinding copy.
- Styles:
  - `internal/tui/styles.go`: many package-level `lipgloss.NewStyle()` vars + helper functions.
  - `internal/tui/theme.go`: palette + a `Theme` struct that is not currently the main way styles are consumed.
- Duplicated patterns: width/height clamping and repeated header/footer composition across views.

IO & side effects:
- Plan file read/parsing:
  - `internal/tui/program.go` reads plan files via `os.ReadFile` (in `loadTasksForMode` / `parsePhasesFromData`).
  - `Model.reloadTasks()` calls `loadTasksForMode(...)` directly from `Update` (filesystem IO in `Update`).
- Loop execution:
  - `r` key starts controller by creating a context and spawning a goroutine (`go m.runController()`).
- Logging: appended to `m.logs` with a max length cap.

Concurrency:
- `tea.Tick` drives animation; Bubble Tea calls `Update` serially.
- Controller callbacks call `program.Send` from controller goroutines (expected).
- Risk: `runController()` currently mutates `m.state` / `m.err` from a goroutine (outside `Update`) which is a data race and violates the Bubble Tea update model.
- Note: the help screen copy in `internal/tui/keybindings.go` includes backend/options text that may not match the current CLI flags; treat copy changes as a bug fix (and lock down with tests) if updated.

## Refactor Targets (Preferred Architecture)

Move toward:
- Model composition: small root model that routes to sub-models (screens/components).
- Message taxonomy: explicit `tea.Msg` structs in a dedicated `msg` package (no ad-hoc strings).
- Side effects via commands: IO and controller execution in `tea.Cmd` functions in an `effects` package; keep `Update` mostly pure.
- Styling system: central style/theme package; create styles once (no per-frame `lipgloss.NewStyle()`).
- Layout helpers: header/footer wrappers, centered/boxed panels, common row renderers; keep `View()` functions short.
- Explicit state transitions: a typed screen enum/router; remove implicit flags (`activeView` + `helpVisible` + `viewMode` overlap).
- Error handling: errors represented by messages, displayed consistently; avoid mutating model from goroutines.

## Proposed Structure (Target End-State)

This is an aspirational layout; we will migrate incrementally.

```
internal/tui/
  program.go              // tea.NewProgram + controller callback wiring (thin)
  root/                   // root model + routing
  screens/                // screen models (split, tasks, output, logs, help, circuit, error)
  msg/                    // tea.Msg taxonomy
  effects/                // tea.Cmd builders (run controller, read plan, timers)
  plan/                   // plan parsing/loading (pure parsing + small IO adapter)
  view/                   // reusable render components (header/footer/panels)
  style/                  // Theme + styles + layout constants
```

Stability rules:
- Keep `tui.NewProgram(...)` and `(*Program).Run()` stable.
- Preserve keybindings and screen flows exactly (unless a bug is fixed with evidence).
- Preserve output formatting (header/status panes/help/circuit) to the extent practical; where small diffs are unavoidable, lock them down with tests before change.

## Risk Areas

- Concurrency correctness:
  - Moving controller execution into `tea.Cmd` must preserve pause/resume behavior and event ordering.
- Rendering diffs:
  - Extracting view components can subtly change spacing/truncation; add targeted snapshot-ish tests for stable sections (header/footer/status) with fixed widths.
- Plan parsing:
  - `parsePhasesFromData` supports multiple plan formats; refactor must preserve header detection and task extraction rules.
- Dedup logic:
  - SSE cumulative output handling (`currentMessage`, `currentReasoning`, `seenMessages`) is easy to regress; add tests around these reducers before altering.

## Sequencing (Phases + Tasks)

Follow `COMMIT_SERIES.md` and use the per-commit checklist in `checklists/`.

### Phase 1: Safety Net (Baseline Tests + Documentation)
- [ ] Add baseline tests for view routing (keybindings -> viewMode) and stable render sections (header/footer/status) with fixed widths.
- [ ] Add reducer-level tests for output/reasoning deduplication behavior (pure logic only).
- [ ] Document current routing/state invariants in code comments near the root model.

### Phase 2: Message Taxonomy
- [ ] Create `internal/tui/msg` and move all TUI `tea.Msg` types into it (including those currently in `model.go` and `messages.go`).
- [ ] Replace magic strings (log levels, view names) with typed constants/enums where safe.

### Phase 3: Side Effects via Commands (Eliminate Goroutine Mutation)
- [ ] Replace `go m.runController()` with a `tea.Cmd` that runs the controller and returns messages; ensure all model mutation happens in `Update`.
- [ ] Move plan reload IO out of `Update` into `effects` with an injectable `ReadFile` function for tests.
- [ ] Ensure pause/resume still calls controller methods and does not block `Update`.

### Phase 4: Plan Parsing & Loading
- [ ] Extract plan parsing to `internal/tui/plan` as pure functions; keep compatibility with current header detection.
- [ ] Keep `loadTasksForMode` as a thin IO wrapper calling the pure parser (or move IO wrapper into `effects`).
- [ ] Expand tests for plan formats (phase headers, atomic headers, fix-plan sections).

### Phase 5: Explicit Screen State + Model Composition
- [ ] Introduce a typed screen enum/router in the root model (single source of truth).
- [ ] Extract help + circuit into sub-models first (lowest coupling) while preserving rendering and keys.
- [ ] Gradually extract split/tasks/output/logs screens into sub-models; root owns shared state (width/height/theme/controller buffers).

### Phase 6: View Components + Styling System
- [ ] Introduce `internal/tui/view` helpers for header/status/footer/panels and reuse them across screens.
- [ ] Introduce `internal/tui/style.Theme` holding styles + layout constants; build once and store in the root model.
- [ ] Reduce per-frame style construction in render functions (especially `lipgloss.NewStyle()` calls inside `View()` paths).

### Phase 7: Cleanup & Consolidation
- [ ] Delete or deprecate redundant fields (`activeView`, `helpVisible`) once router is the source of truth.
- [ ] Reduce `internal/tui/model.go` and `internal/tui/views.go` size by moving cohesive blocks into packages.
- [ ] Update docs (`internal/tui/REDESIGN.md` if needed) to match the new structure.

## Success Criteria

- `go test ./...` passes after every commit.
- Keybindings unchanged; screen transitions unchanged.
- No flicker/regressions in layout; resize behavior remains acceptable.
- No goroutine writes to the model; side effects happen via `tea.Cmd` and return explicit messages.
- Codebase has clearer boundaries: messages, effects, plan parsing, view components, styles.
