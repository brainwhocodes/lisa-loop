# Implementation Plan

## Goal
Enable Ralph's TUI loop to run up to 10 iterations using the OpenCode server API, with the default model set to the z.ai GLM 4.7 coding plan model. Add authentication handling per OpenCode server docs and define a commit-by-commit delivery plan with DONE.md updates.

## Source Notes
- OpenCode server auth uses HTTP basic auth with `OPENCODE_SERVER_PASSWORD`.
- Username defaults to `opencode`, or override with `OPENCODE_SERVER_USERNAME`.
- Applies to both `opencode serve` and `opencode web`.
- Models.dev entry for Z.AI Coding Plan GLM-4.7 uses model ID `glm-4.7` (provider ID `zhipuai-coding-plan`).

## Process Notes
- Each commit below must append a short summary to `DONE.md` and mark that commit's checklist item as complete.
- Before starting the next commit (including after compaction), review `IMPLEMENTATION_PLAN.md` and `DONE.md`.

## Atomic Commits

### 1) Config + CLI surface for OpenCode backend
- Add config fields for OpenCode server URL, auth username/password, and model ID.
- Add CLI flags and env fallbacks (e.g., `OPENCODE_SERVER_URL`, `OPENCODE_SERVER_USERNAME`, `OPENCODE_SERVER_PASSWORD`, `OPENCODE_MODEL_ID`).
- Establish backend name (e.g., `opencode`) and wire it through config plumbing.
- Default max calls to 10 when backend is `opencode`, while preserving existing defaults for other backends.
- Default model ID to `glm-4.7` for Z.AI Coding Plan (override via flag/env).
- Update help text and any docs that list supported backends/flags.
- DONE.md entry for this commit.

### 2) OpenCode server client + session persistence
- Add an HTTP client wrapper (new package) for `/session` and `/session/:id/message` endpoints.
- Implement basic auth headers and request timeouts.
- Persist session IDs in a dedicated state file (e.g., `.opencode_session_id`) with load/save helpers.
- Map OpenCode response payloads into the existing runner output format (message text + session ID).
- Unit tests using `httptest` for auth headers, session creation, message send, and error handling.
- DONE.md entry for this commit.

### 3) Runner integration + TUI loop behavior
- Extend `internal/codex.Runner` to route to the OpenCode backend when selected.
- Ensure the loop/controller uses the new backend without breaking CLI JSONL parsing paths.
- Validate loop iteration limit up to 10 in TUI display and rate limiter state.
- Add tests (or update existing ones) to cover the backend selection and loop count behavior.
- DONE.md entry for this commit.

### 4) Docs + examples for OpenCode usage
- Document OpenCode server setup, auth env vars, and example usage with `--backend opencode`.
- Document the default model setting: z.ai GLM 4.7 coding plan model, plus how to override.
- Include quick-start examples for running the TUI with 10-iteration loops.
- DONE.md entry for this commit.
