# Commit 05 — Codex SDK runner + bash integration

## Intent
- Add Codex SDK runner using `@openai/codex-sdk` for programmatic control and thread resume.

## Scope / touched areas
- `package.json` (add dependency: `@openai/codex-sdk`)
- `src/` (new small runner CLI, Node >= 18)
- `ralph_loop.sh` (dispatch to Node runner)
- `tests/**` (mock runner behavior; do not hit network)

## Steps (atomic)
1. Add `@openai/codex-sdk` dependency.
2. Create a minimal Node runner (single responsibility):
   - Inputs: prompt file path, working directory, optional thread id file, output paths.
   - Behavior:
     - `const codex = new Codex();`
     - `startThread()` if no prior thread id; else `resumeThread(id)`.
     - `await thread.run(promptText)` and write:
       - final message to a file
       - thread id to state file
   - Keep it small; no extra features yet.
3. Wire `ralph_loop.sh`:
   - Call `node <runner>` and read the output file as “agent output”.
   - Respect `--continue` by supplying thread id file path.
   - If Node runner is missing or fails with “module not found”, fall back to a warning.
4. Tests:
   - Add a test-mode switch (env var) for the Node runner to emit deterministic output without calling Codex.
   - BATS tests validate:
     - thread id persistence
     - resume path used on next loop
     - log file created

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
- `ralph` runs through the loop pipeline (logs + analyzer) using the Node runner output.
- `--continue` reuses the prior thread id and calls `resumeThread(...)`.

## Rollback note
- Revert Node runner + dependency if issues arise.
