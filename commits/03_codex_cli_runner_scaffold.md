# Commit 03 — Codex CLI runner (`codex exec`) scaffold

## Intent
- Implement Codex CLI runner by invoking `codex exec` in non-interactive mode.
- Keep it minimal: single-run invocation + capture final message.

## Scope / touched areas
- `ralph_loop.sh` (`execute_codex()` implementation)
- `lib/` (optional helper for command building)
- `tests/helpers/mocks.*` (add mock `codex` binary behavior)
- `tests/unit/**` and/or `tests/integration/**`

## Steps (atomic)
1. Add `execute_codex_cli()`:
   - Use `codex exec "<prompt>"` with safe defaults:
     - Prefer least privilege (start with read-only / approvals unless Ralph requires write).
     - Allow overriding via env vars (e.g., `RALPH_CODEX_SANDBOX`, `RALPH_CODEX_FULL_AUTO`).
2. Implement prompt passing:
   - Read from the same `PROMPT.md` Ralph uses today.
   - Avoid fragile quoting: pass content via a temp file → construct a single argument robustly (document approach in code comments).
3. Capture Codex output:
   - Per docs, `codex exec` prints only final agent message to `stdout`; progress goes to `stderr`.
   - Write stdout to the per-loop log file in `logs/` (parallel to Claude logs).
4. Add dependency check:
   - If `codex` is missing, print actionable instructions (`npm i -g @openai/codex`).
5. Tests:
   - Add a mock `codex` executable in test PATH that records argv and emits deterministic output.
   - Validate: correct subcommand `exec` used; expected flags appear; logs written.

## Tests (MUST RUN)
- Unit: `npm run test:unit` *(assumption; else `npm test`)*
- Integration: `npm run test:integration` *(assumption; else `npm test`)*
- E2E: `npm run test:e2e` *(assumption; else N/A + `npm test`)*
- Lint/Typecheck: `npm run lint` *(assumption)*, `npm run typecheck` *(assumption)*

## Gating criteria (checkboxes)
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] E2E tests pass
- [ ] Lint/Typecheck pass (if applicable)

## Acceptance check
- With a mocked `codex`, `ralph` performs one iteration and writes `logs/codex_output_*.log` (or equivalent).
- All tests remain green.

## Rollback note
- Revert only the new dispatch branch and the mock additions.
