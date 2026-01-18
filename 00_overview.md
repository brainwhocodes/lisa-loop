# Overview

## Goal
- Make `ralph` able to run loops using **Codex CLI** (`codex exec`) in addition to Claude Code.
- Add a **Codex SDK** backend (`@openai/codex-sdk`) for programmatic, resumable sessions/threads.
- Ship **Codex-friendly project scaffolding** (Codex config + skills + AGENTS.md) via Ralph templates so new Ralph projects “just work” with Codex.

## Repo context (what we’re converting)
- Ralph is primarily **Bash** with a modular layout (`ralph_loop.sh`, `lib/*`, `setup.sh`, `install.sh`, templates, tests).
- Tests are **BATS via npm** (repo has `package.json`, `tests/`, CI workflow).

## Assumptions
- There is an existing `npm test` (or equivalent) that runs the BATS suite; if per-scope scripts exist (`test:unit`, `test:integration`), we’ll use them; otherwise we’ll run the full suite for every gate.
- We keep the current Claude backend working while adding Codex support behind a flag (`--agent`). (Safer migration; avoids breaking existing users.)
- Node.js **>= 18** is available/required for the Codex SDK backend.
- Codex CLI is installed by users (or via optional installer help) as `@openai/codex`.

## Non-goals
- Rewriting Ralph from Bash into another language.
- Implementing a Codex IDE extension or MCP server.
- Running real Codex network calls in CI/unit tests (we’ll mock `codex` / SDK calls).

## Key risks
- **Output format differences**: Codex `exec --json` emits JSONL events; parsing must be robust and backwards-compatible with existing analyzers.
- **Session/thread continuity**: Codex CLI supports `exec resume`; SDK uses `thread.run()` with resumable thread IDs; state must be persisted safely.
- **Git requirement**: Codex may require running inside a git repo unless `--skip-git-repo-check`; Ralph already expects git but must handle edge cases cleanly.
- **Auth**: `CODEX_API_KEY` only works for `codex exec`; SDK auth expectations differ; docs and checks must be explicit.
- **Prompt size / quoting**: Passing large `PROMPT.md` content to `codex exec` must avoid shell-quoting pitfalls.

## Test commands (assumed; adjust to repo scripts)
- Install deps: `npm ci` (or `npm install`)
- Unit:
  - `npm run test:unit` *(assumption; else `npm test`)*
- Integration:
  - `npm run test:integration` *(assumption; else `npm test`)*
- E2E:
  - `npm run test:e2e` *(assumption; else N/A and run `npm test`)*
- Lint/Typecheck:
  - `npm run lint` *(assumption; if missing, skip)*
  - `npm run typecheck` *(assumption; if missing, skip)*

## Definition of Done
- [ ] `ralph --agent codex-cli` runs a loop iteration using `codex exec` (mockable in tests).
- [ ] `ralph --agent codex-sdk` runs a loop iteration using `@openai/codex-sdk` and persists/resumes thread IDs.
- [ ] Templates include Codex project config (`AGENTS.md`, `.codex/*`) and a small skills pack.
- [ ] Installer/docs clearly describe prerequisites, auth, and safety defaults.
- [ ] Unit + integration tests pass for every commit; no regression for Claude path.
