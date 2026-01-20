# Refactor Plan

## Overview
Stabilize loop exit analysis, reduce duplicated Codex and project-mode logic, and improve testability across loop, project, and codex packages without changing external behavior.

## Current State Analysis
- Codex execution is duplicated in `internal/codex/runner.go`, `internal/project/init.go:generateWithCodex`, and `internal/project/setup.go:generateTemplatesWithCodex`, with inconsistent streaming and JSONL parsing.
- Exit analysis exists in `internal/analysis/response.go` but is unused; `internal/loop/controller.go` has a TODO and relies on string heuristics for errors and file changes.
- Project-mode detection and root validation are implemented in multiple places (`internal/loop/context.go` and `internal/project/setup.go`), which risks divergence.
- `internal/project/setup.go:executeCommand` is a no-op, so git initialization and commits report success without executing.
- Event and log types are stringly-typed across `internal/loop/controller.go` and `internal/tui/*`, increasing the chance of drift.
- JSONL scanning in `internal/codex/runner.go` uses default `bufio.Scanner` limits and ignores `scanner.Err()`, which can truncate large Codex outputs.
- `internal/project/import.go:parseSourceContent` resets builders on repeated headings, which can drop earlier section content.
- `internal/loop/ratelimit.go:LoadState` hard-codes defaults rather than preserving configured limits.

## Phase 1: Low-Risk Foundation
### Goals
- Add tests around parsing and event boundaries.
- Reduce string drift and improve diagnostics without changing logic.

### Tasks
- [ ] Task 1: Introduce typed constants for loop event types and log levels in `internal/loop/events.go`, then replace string literals in `internal/loop/controller.go` and `internal/tui/model.go`.
- [ ] Task 2: Add unit tests for `internal/analysis/response.go` covering `DetectFormat`, `ParseRALPHStatus`, `DetectCompletionKeywords`, and `ExtractErrors`.
- [ ] Task 3: Expand `internal/codex/runner.go:runCLI` to set `scanner.Buffer` and check `scanner.Err()`; add tests in `internal/codex/runner_test.go` for large JSONL lines.
- [ ] Task 4: Add tests for `internal/codex/events.go:ParseEvent` covering `message`, `content_block_delta`, `tool_use`, and `tool_result` event variants.

### Verification
- [ ] All tests pass (`go test ./...`)
- [ ] No functionality changes in loop behavior
- [ ] Code review completed

## Phase 2: Core Refactoring
### Goals
- Centralize Codex execution and wire exit analysis into the loop.
- Make project mode detection and state handling single-sourced.

### Tasks
- [ ] Task 1: Integrate `internal/analysis.Analyze` into `internal/loop/controller.go:ExecuteLoop` to drive `hasErrors`, `exitSignal`, and `completion` decisions; persist exit signals using `internal/state` helpers.
- [ ] Task 2: Consolidate project-mode detection and root validation into `internal/project/mode.go`, then replace `internal/loop/context.go:DetectProjectMode/CheckProjectRoot` and `internal/project/setup.go:ValidateProject/GetProjectRoot` with the shared helpers.
- [ ] Task 3: Extract a shared Codex invocation helper (e.g., `internal/project/codex.go:RunCodex`) that wraps `internal/codex.Runner` and use it from `internal/project/init.go` and `internal/project/setup.go`.
- [ ] Task 4: Implement `internal/project/setup.go:executeCommand` using `os/exec` and inject a `CommandRunner` for testability; add tests around `initGitRepo`.
- [ ] Task 5: Update `internal/loop/ratelimit.go:LoadState` to preserve configured `maxCalls` and `resetHours` instead of hard-coded defaults.

### Verification
- [ ] All tests pass (`go test ./...`)
- [ ] Init/setup flows generate the same files as before (spot-check with fixtures)
- [ ] Loop exits reflect RALPH_STATUS and completion keywords consistently

## Phase 3: Cleanup and Polish
### Goals
- Remove dead code and fix edge-case parsing.
- Align docs and build tooling with the actual entrypoints.

### Tasks
- [ ] Task 1: Update `internal/project/import.go:parseSourceContent` to append on repeated headings; add tests covering interleaved sections.
- [ ] Task 2: Remove unused or redundant helpers in `internal/analysis/response.go` (or wire `calculateConfidence` and `ErrorMessages` into `Analyze`).
- [ ] Task 3: Remove unused parameters like `internal/loop/context.go:BuildContext(promptPath string)` or make them functional; update callers and tests.
- [ ] Task 4: Reconcile entrypoint references in `Makefile`, `README.md`, and `AGENTS.md` with the actual `cmd/` tree.

### Verification
- [ ] All tests pass (`go test ./...`)
- [ ] Documentation updated and consistent with repository layout
- [ ] No dead code remaining (`rg "TODO|unused"` yields no stale items)

## Rollback Plan
- Revert changes per phase if regressions appear; each task should be independently committable.
- Keep any new tests even if code changes are rolled back to preserve coverage.

## Success Criteria
- [ ] Loop exit behavior uses `internal/analysis` and recorded exit signals instead of string heuristics.
- [ ] Codex execution has a single shared helper across `internal/codex` and `internal/project`.
- [ ] Project-mode detection and root validation are implemented once and reused.
- [ ] Git initialization executes real commands with actionable errors.
- [ ] Parsing, rate-limit state, and project import edge cases are covered by tests.
