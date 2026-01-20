# Implementation Progress

## Commit 1: Config + CLI surface for OpenCode backend

- Added OpenCode config fields to `internal/config/config.go`:
  - `OpenCodeServerURL` - URL for OpenCode server
  - `OpenCodeUsername` - Username for basic auth
  - `OpenCodePassword` - Password for basic auth
  - `OpenCodeModelID` - Model ID (default: glm-4.7)

- Added CLI flags and environment variable fallbacks in `cmd/ralph/main.go`:
  - `--backend` - Backend selection: cli or opencode
  - `--opencode-url` - Server URL (env: OPENCODE_SERVER_URL)
  - `--opencode-user` - Username (env: OPENCODE_SERVER_USERNAME, default: opencode)
  - `--opencode-pass` - Password (env: OPENCODE_SERVER_PASSWORD)
  - `--opencode-model` - Model ID (env: OPENCODE_MODEL_ID, default: glm-4.7)

- Default max calls set to 10 when backend is `opencode` (vs 3 for cli)

- Updated help text with new backend options section
