# Codex Configuration Directory

This directory contains configuration for the Codex AI agent integration.

## Files

- `config` - Main configuration file for Codex backend selection and options

## Configuration Options

### Backend Selection
- `backend=cli` - Use Codex CLI (default)
- `backend=sdk` - Use Codex SDK (requires Node.js >= 18)

### CLI Backend Settings
- `codex_exec_timeout` - Maximum execution time in seconds (default: 600)
- `codex_resume_session` - Continue previous session (default: true)

### SDK Backend Settings
- `sdk_runner_path` - Path to Node.js SDK runner (default: ../src/codex_runner.js)
- `sdk_output_file` - Temporary output file for SDK runner

### Loop Context
- `loop_context_enabled` - Enable loop context injection (default: true)
- `loop_context_lines` - Number of context lines to inject (default: 5)

### File Operations
- `skip_git_repo_check` - Skip git repository validation (default: false)

## Backend Usage

The Ralph loop reads this configuration and selects the appropriate backend:

```bash
# CLI backend (default)
ralph --backend cli

# SDK backend
ralph --backend sdk
```

## Safe-By-Default Configuration

This configuration is designed to be safe and predictable:
- No dangerous sandbox permissions
- File operations require explicit allowlist
- Timeouts prevent runaway execution
- Session continuity respects manual resets
