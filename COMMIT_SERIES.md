# Refactor Commit Series (Atomic, Reviewable)

This file defines the intended sequence of small refactor commits for the Bubble Tea + Lipgloss TUI (`internal/tui/`).

Rules for every commit:
- Title format: `refactor(<area>): <short summary>`
- Preserve behavior and UX (keybindings, flows, formatting) unless explicitly called out as a bug fix with before/after evidence.
- Keep the public surface stable (`tui.NewProgram`, `Program.Run`, CLI flags).
- No blocking in `Update`; side effects happen in `tea.Cmd`.
- `go test ./...` passes after every commit.
- Use the checklist in `checklists/NN_*.md` for that commit.

## 01) refactor(tui/tests): add baseline invariants for routing + stable view sections

What changed:
- Add/extend tests that lock down key routing to screens and stable render sections (header/footer/status) at fixed widths.
- Add reducer-level tests around output/reasoning dedup (pure functions only).

Why it's safe:
- Tests only; no production logic changes.

Tests added/updated:
- `internal/tui/*_test.go` (new tests; keep them width-stable and not dependent on terminal quirks).

Manual check:
- `lisa --monitor` launches; `q` exits; `l/t/o/?/c` toggles still work.

Checklist: `checklists/01_refactor-tui-tests-baseline.md`

## 02) refactor(tui/msg): centralize tea.Msg taxonomy in internal/tui/msg

What changed:
- Create `internal/tui/msg` package.
- Move all message types currently split across `internal/tui/messages.go` and `internal/tui/model.go` into `msg`.
- Replace ad-hoc stringly-typed message payloads with explicit structs where it reduces ambiguity.

Why it's safe:
- Mechanical type move + import updates; no logic changes.
- Add compile-time guarantees by making message types explicit and discoverable.

Tests added/updated:
- Update any tests referencing moved message types.

Manual check:
- Launch TUI and verify controller events still update output/logs/status.

Checklist: `checklists/02_refactor-tui-msg-taxonomy.md`

## 03) refactor(tui/effects): run controller via tea.Cmd (no goroutine mutation)

What changed:
- Replace `go m.runController()` with a `tea.Cmd` that runs the controller and returns `msg.ControllerDone{Err error}` (or similar).
- Ensure all model mutation (state/errors) happens in `Update`.
- Keep pause/resume behavior and event callbacks unchanged.

Why it's safe:
- Preserves behavior but fixes a correctness issue: model writes from a goroutine are a data race.
- Bubble Tea pattern-compliant: side effects return messages; `Update` remains the only mutator.

Tests added/updated:
- Add unit tests verifying:
  - `r` schedules a command (does not mutate state from a goroutine).
  - controller completion/error transitions state via a message.

Manual check:
- Start loop with `r`, pause/resume with `p`, exit with `q`.

Checklist: `checklists/03_refactor-tui-effects-controller-cmd.md`

## 04) refactor(tui/plan): extract plan parsing/loading into internal/tui/plan

What changed:
- Move `parsePhasesFromData`, header detection, and task extraction into `internal/tui/plan` as pure functions.
- Keep IO wrapper thin (either in `plan` or `effects`) and preserve current header rules.

Why it's safe:
- Parsing becomes testable as pure functions; behavior locked via tests.

Tests added/updated:
- Add tests for all supported plan formats:
  - `REFACTOR_PLAN.md` (`## Phase N: ...`)
  - `IMPLEMENTATION_PLAN.md` (`## Phase ...` and `### N) ...`)
  - `@fix_plan.md` priority headers.

Manual check:
- With a real project, ensure tasks/phase progress shown in header and task pane match the plan file.

Checklist: `checklists/04_refactor-tui-plan-parsing.md`

## 05) refactor(tui/effects): move filesystem IO out of Update (reload tasks via Cmd)

What changed:
- Convert `Model.reloadTasks()` into a command-driven flow:
  - `Update` receives `msg.Preflight` and schedules `effects.LoadPlan(...)`.
  - `effects.LoadPlan` reads the plan file and returns a `msg.PlanLoaded{...}`.
- Add injection point (function var or small interface) for file reads so tests do not touch disk.

Why it's safe:
- Eliminates filesystem IO from `Update` while preserving the same outcome.
- Makes plan reload deterministic and testable.

Tests added/updated:
- Tests for:
  - preflight triggers reload command
  - plan reload merges completed in-memory tasks consistently (same as current behavior)

Manual check:
- Modify plan file while TUI is running (or between loops) and verify task list refreshes after the next preflight.

Checklist: `checklists/05_refactor-tui-effects-plan-io.md`

## 06) refactor(tui/state): make screen state explicit and remove redundant flags

What changed:
- Make the active screen/view the single source of truth (typed enum).
- Remove or deprecate overlapping fields (`activeView`, `helpVisible`) once routing is explicit.
- Centralize key handling for view toggles into small helpers (pure functions if possible).

Why it's safe:
- Routing logic becomes easier to reason about; behavior locked down by baseline tests from commit 01.

Tests added/updated:
- Routing tests updated to assert screen enum transitions.

Manual check:
- Toggle each view twice returns to split; error/help interactions still behave the same.

Checklist: `checklists/06_refactor-tui-state-router.md`

## 07) refactor(tui/view): extract reusable render components + layout helpers

What changed:
- Introduce `internal/tui/view` helpers:
  - header renderer
  - status bar renderer
  - footer renderer
  - boxed/padded panel helpers
- Keep per-screen view functions short and composable.

Why it's safe:
- Pure rendering refactor; output locked via snapshot-ish tests for stable sections.

Tests added/updated:
- Snapshot-style tests for header/footer/status at fixed widths and known model state.

Manual check:
- Resize terminal (if supported) and confirm layout remains acceptable.

Checklist: `checklists/07_refactor-tui-view-components.md`

## 08) refactor(tui/style): centralize theme/styles and avoid per-frame style construction

What changed:
- Introduce `internal/tui/style.Theme` that owns:
  - all `lipgloss.Style` instances
  - layout constants (padding, min widths)
- Root model holds a single Theme instance; views use it instead of package-level globals.

Why it's safe:
- Styles are constructed once (already mostly true today); this makes it explicit and consistent.
- Rendering output should remain identical; tests from commit 07 should catch diffs.

Tests added/updated:
- Theme construction tests (no nils, expected constants).
- Update snapshot tests if needed (prefer no change).

Manual check:
- Visually confirm colors and spacing match previous output; no flicker introduced.

Checklist: `checklists/08_refactor-tui-style-theme.md`

## Optional follow-ups (only if treated as bug fixes with evidence)

- `refactor(tui/help): align help screen copy with actual CLI flags`
  - Evidence to capture before change: help screen currently lists backend/options that differ from `cmd/lisa/main.go` flags.
  - After change: help screen reflects the real supported backends/flags; update tests to assert correct copy.
