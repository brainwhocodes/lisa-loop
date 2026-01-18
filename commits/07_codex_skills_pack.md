# Commit 07 — Add Codex skills pack in templates

## Intent
- Provide a small, repo-local “skills” pack for Codex workflows (run tests, summarize status, etc.).
- Update prompts/instructions to encourage using skills where appropriate.

## Scope / touched areas
- `templates/.codex/skills/` (new)
- `templates/PROMPT.md` (or equivalent template prompt)
- `templates/AGENTS.md`
- `tests/**` (template integrity)

## Steps (atomic)
1. Add `templates/.codex/skills/` with a minimal set:
   - `run-tests` (runs the canonical test command(s))
   - `lint` (runs lint/typecheck if present)
   - `status-summary` (prints/parses Ralph `status.json` + latest log pointers)
2. Keep each skill:
   - Small
   - Least-privilege by default
   - Explicit about commands it runs
3. Update `templates/PROMPT.md`:
   - Instruct the agent to use skills (when available) for consistent test execution and status reporting.
4. Update `setup.sh` / `ralph_import.sh` if needed so skills are copied.
5. Tests:
   - Validate skills directory is created/copied.
   - Validate no regressions to existing template tests.

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
- New projects contain `.codex/skills/*` and prompt guidance mentions them.
- Skills integrate seamlessly with Codex workflow.

## Rollback note
- Revert skills templates + prompt text changes if skill format needs revision.
