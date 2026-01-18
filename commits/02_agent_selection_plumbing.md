# Commit 02 — Replace Claude with Codex backend

## Intent
- Remove Claude code and replace with Codex as the sole agent backend.
- Establish a single internal "execute agent" entrypoint using Codex.

## Scope / touched areas
- `ralph_loop.sh`
- (Possible) `lib/*` for shared helpers
- `tests/**` (add minimal coverage for the new flag behavior)

## Steps (atomic)
1. Remove `execute_claude_code()` and all Claude-specific code from `ralph_loop.sh`.
2. Replace with `execute_codex()` as the default and only agent backend.
3. Remove any `--agent` flag logic (Codex is now the only option).
4. Ensure `status.json` (or equivalent status output) reflects Codex as the agent.
5. Add/adjust BATS tests:
   - Verify Codex is used by default.
   - Remove any Claude-specific test cases.

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
- Running `ralph` uses Codex as the agent backend.
- All Claude-specific code has been removed.

## Rollback note
- Revert to previous commit if Codex integration fails.
