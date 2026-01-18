# Commit 01 — Add Codex integration plan docs

## Intent
- Document how Ralph will use Codex CLI (`codex exec`) and Codex SDK (`@openai/codex-sdk`).
- Create a stable reference before wiring code.

## Scope / touched areas
- `README.md`
- `docs/` (new: `docs/codex.md` or similar)
- (Optional) `CONTRIBUTING.md`

## Steps (atomic)
1. Add `docs/codex.md`:
   - Codex CLI install + `codex exec` basics.
   - Note `--json` JSONL streaming and why we’ll parse it.
   - SDK install + minimal usage (`Codex`, `startThread()`, `thread.run()`, `resumeThread()`).
2. Update `README.md`:
   - Update to reflect Codex as the agent backend.
   - Link to `docs/codex.md`.
3. (Optional) Add “dev note” in `CONTRIBUTING.md` describing mocking strategy (no real Codex calls in tests).

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
- Docs clearly explain Codex CLI and SDK usage.
- No runtime behavior changes; test suite remains green.

## Rollback note
- Revert doc changes only (no code touched).
