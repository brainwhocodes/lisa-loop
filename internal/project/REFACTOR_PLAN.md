# Refactor Plan

## Overview
Reduce duplication in Codex invocation, make command execution and template resolution testable, and tighten parsing/error handling in `internal/project` without changing user-facing behavior or breaking existing tests.

## Current State Analysis
- Codex execution is duplicated between `internal/project/init.go:generateWithCodex` and `internal/project/setup.go:generateTemplatesWithCodex`, with different prompt delivery and output parsing.
- `internal/project/setup.go:executeCommand` is a stub, so `initGitRepo` reports success without running git commands.
- Template resolution in `internal/project/prompts.go:loadPromptTemplate` relies on global state (`TemplateDir`) and implicit search order, which makes tests and overrides brittle.
- `internal/project/import.go:parseSourceContent` resets section builders on headings, which can drop earlier content when multiple headings appear.
- Error messages are inconsistent and often omit context (which file/action failed), making failures harder to diagnose.

## Phase 1: Low-Risk Foundation
### Goals
- Make command execution real and observable.
- Add tests around the most fragile helpers before refactoring behavior.

### Tasks
- [ ] Task 1: Implement `internal/project/setup.go:executeCommand` using `os/exec`, returning errors that include the command and exit status.
- [ ] Task 2: Add a `CommandRunner` interface in `internal/project/runner.go` and inject it into `internal/project/setup.go:initGitRepo` for test doubles.
- [ ] Task 3: Add unit tests for `internal/project/init.go:readSpecsFolder` covering empty folders, non-markdown files, and multi-file ordering.
- [ ] Task 4: Add unit tests for `internal/project/prompts.go:loadPromptTemplate` using temp dirs to verify custom, home, and default resolution order.

### Verification
- [ ] All tests pass (`go test ./...`)
- [ ] No functionality changes in generated files for existing inputs
- [ ] Code review completed

## Phase 2: Core Refactoring
### Goals
- Consolidate Codex execution and JSONL parsing into a single helper.
- Normalize output assembly and error handling across init/setup flows.

### Tasks
- [ ] Task 1: Extract `internal/project/codex.go:RunCodex(prompt string, opts CodexOptions)` and route `internal/project/init.go:generateWithCodex` through it.
- [ ] Task 2: Replace `internal/project/setup.go:generateTemplatesWithCodex` with the shared helper while preserving its working directory behavior.
- [ ] Task 3: Add `internal/project/codex.go:ParseCodexJSONL(r io.Reader)` with a larger scanner buffer and use it in all Codex callers.
- [ ] Task 4: Add tests for JSONL parsing covering `message`, `content_block_delta`, and `assistant` content array events.

### Verification
- [ ] All tests pass (`go test ./...`)
- [ ] Performance unchanged for large prompts (spot check with verbose output)
- [ ] Init and Setup flows produce identical files to pre-refactor runs

## Phase 3: Cleanup and Polish
### Goals
- Preserve all imported content.
- Make template resolution explicit and configurable without breaking defaults.

### Tasks
- [ ] Task 1: Update `internal/project/import.go:parseSourceContent` to accumulate section content instead of resetting on each heading; add tests with interleaved headings.
- [ ] Task 2: Introduce `internal/project/prompts.go:TemplateResolver` and route `Get*Prompt` helpers through it while keeping `TemplateDir` as a backward-compatible default.
- [ ] Task 3: Normalize errors in `internal/project/init.go`, `internal/project/setup.go`, and `internal/project/import.go` to include file paths and action context.
- [ ] Task 4: Remove dead or duplicated helpers uncovered by consolidation (e.g., legacy prompt-loading paths).

### Verification
- [ ] All tests pass (`go test ./...`)
- [ ] Documentation updated if public behavior or configuration changes
- [ ] No dead code remaining

## Rollback Plan
- Revert the commit for the impacted phase; each task is independently committable.
- Restore original Codex execution functions if the shared helper changes output.
- Keep new tests even if code changes are reverted to preserve coverage.

## Success Criteria
- [ ] `internal/project/setup.go:initGitRepo` performs real git initialization with actionable errors on failure.
- [ ] All Codex invocations use one helper with consistent JSONL parsing.
- [ ] Template resolution is testable without relying on global state.
- [ ] Import parsing preserves all sections with repeated headings.
- [ ] Unit tests cover specs reading, template resolution, and JSONL parsing.
