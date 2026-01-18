# Commit 06 — Add Codex templates (AGENTS.md + .codex config)

## Intent
- Make newly-created Ralph projects Codex-ready by default (clear commands + safe config).
- Set up Codex templates as the default.

## Scope / touched areas
- `templates/` (new: `AGENTS.md`, `.codex/`)
- `setup.sh` (project creation copies new files)
- `ralph_import.sh` (PRD import also emits new files)
- `tests/integration/**` (template copying expectations)

## Steps (atomic)
1. Add `templates/AGENTS.md`:
   - Project-specific "how to build/test/lint" for Codex to read.
2. Add `templates/.codex/` directory:
   - Minimal config scaffold (safe-by-default; avoid dangerous sandbox by default).
3. Update `setup.sh`:
   - Ensure new projects include `AGENTS.md` + `.codex/` contents.
4. Update `ralph_import.sh` similarly.
5. Update/extend tests that validate project scaffolding:
   - Assert `AGENTS.md` exists.
   - Assert `.codex/` directory exists and contains expected files.

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
- `ralph-setup <project>` creates Codex-ready artifacts without breaking existing workflows.
- PRD import output includes the same Codex artifacts.

## Rollback note
- Revert template additions + template-copy logic only.
