# Commit 04 — Codex CLI JSONL parsing + resume support

## Intent
- Add session continuity for Codex CLI using `codex exec --json` + `codex exec resume`.
- Persist a Codex thread/session ID in Ralph state files.

## Scope / touched areas
- `ralph_loop.sh` (codex-cli path: JSONL + resume logic)
- `lib/response_analyzer.sh` (if needed to support parsing Codex outputs)
- Project state files (new file e.g. `.codex_thread_id` or extend existing session state)
- `tests/**` (fixtures for JSONL + resume)

## Steps (atomic)
1. Add an opt-in (default-on for codex-cli) to run:
   - `codex exec --json "<prompt>"` to receive JSONL event stream.
2. Implement parsing helpers (bash + `jq`):
   - Extract `thread_id` from the `thread.started` event line.
   - Extract the final agent message from the last `item.*` event representing an agent message (store it in the same place Ralph expects to analyze).
3. Persist thread/session:
   - Store extracted `thread_id` to a stable state file (e.g., `.codex_thread_id`).
   - Add guardrails: if parsing fails, fall back to non-`--json` mode and log a warning.
4. Resume support:
   - When `--continue` is enabled and `.codex_thread_id` exists, call:
     - `codex exec resume <THREAD_ID> "<prompt>"`
5. Tests:
   - Add a JSONL fixture file that includes `thread.started` and a final agent message.
   - Update mock `codex` to emit JSONL for `--json`, and validate resume path uses `exec resume`.

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
- `ralph --continue` uses an existing thread ID and calls `codex exec resume ...`.
- The analyzer receives the extracted final agent message text.

## Rollback note
- Revert JSONL parsing + thread id persistence; keep basic `codex exec` runner.
