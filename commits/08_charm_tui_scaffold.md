# Commit 08 — Charm TUI scaffold (replace shell scripts)

## Intent
- Replace shell-based Ralph management with a Go TUI using Charm libraries.
- Port all functionality from `ralph_loop.sh` (1539 lines), `ralph_monitor.sh`, `ralph_import.sh`, `setup.sh`, and `lib/*.sh`.

## Shell Scripts Being Replaced

| Script | Lines | Purpose |
|--------|-------|--------|
| `ralph_loop.sh` | 1539 | Main loop orchestration |
| `ralph_monitor.sh` | 126 | Live status dashboard |
| `ralph_import.sh` | 626 | PRD conversion |
| `setup.sh` | 36 | Project scaffolding |
| `install.sh` | 294 | Global installation |
| `uninstall.sh` | 195 | Cleanup |
| `lib/circuit_breaker.sh` | 329 | Stagnation detection |
| `lib/response_analyzer.sh` | 704 | Output analysis |
| `lib/date_utils.sh` | 54 | Cross-platform dates |

## Scope / touched areas
- `cmd/ralph/main.go` (CLI entrypoint)
- `internal/` (all Go packages below)
- `go.mod`, `go.sum`
- Deprecate: all `.sh` files except test helpers

## Go Package Structure

```
cmd/ralph/main.go
internal/
├── codex/
│   ├── runner.go      # Execute codex exec with timeout
│   ├── jsonl.go       # Parse JSONL output, extract thread_id
│   └── session.go     # Thread ID persistence (.codex_session_id)
├── loop/
│   ├── controller.go  # Main loop logic, graceful exit
│   ├── ratelimit.go   # Calls per hour tracking
│   └── context.go     # Build loop context for prompts
├── circuit/
│   └── breaker.go     # CLOSED/HALF_OPEN/OPEN states
├── analysis/
│   ├── response.go    # Analyze Codex output (JSON + text)
│   └── signals.go     # Exit signal detection
├── project/
│   ├── setup.go       # Create new project (ralph-setup)
│   └── import.go      # PRD conversion (ralph-import)
├── tui/
│   ├── model.go       # Bubble Tea model
│   ├── keys.go        # Keybindings
│   └── views/
│       ├── status.go  # Current loop, Codex status
│       ├── logs.go    # Scrollable log viewport
│       └── help.go    # Keybinding help
└── state/
    └── files.go       # Manage all state files
```

## State Files to Manage

```
.call_count              # Rate limit counter
.last_reset              # Hour of last rate reset
.codex_session_id        # Codex thread ID
.ralph_session           # Ralph session metadata (JSON)
.ralph_session_history   # Session transitions (JSON array)
.exit_signals            # Completion signal tracking (JSON)
.circuit_breaker_state   # Circuit breaker state (JSON)
.circuit_breaker_history # Circuit breaker transitions (JSON)
.response_analysis       # Last analysis result (JSON)
status.json              # External monitoring
progress.json            # Live execution progress
```

## Steps (atomic)

### Phase 1: Core Infrastructure
1. Initialize Go module with Charm dependencies:
   - `github.com/charmbracelet/bubbletea`
   - `github.com/charmbracelet/lipgloss`
   - `github.com/charmbracelet/bubbles`
   - `github.com/charmbracelet/huh`
   - `github.com/charmbracelet/log`

2. Create `internal/state/files.go`:
   - Read/write all JSON state files
   - Atomic file updates (write to temp, rename)

3. Create `internal/codex/runner.go`:
   - `Execute(prompt, sessionID, timeout) (output, threadID, error)`
   - Build args: `codex exec --json --prompt "..." [--resume --thread-id ID]`
   - Spawn subprocess with context timeout
   - Stream stdout/stderr

4. Create `internal/codex/jsonl.go`:
   - Parse JSONL line by line
   - Extract `thread_id` from `thread.started` event
   - Extract final message from text/message events

### Phase 2: Loop Control
5. Create `internal/loop/ratelimit.go`:
   - Track calls per hour (default 100)
   - Wait for reset with countdown

6. Create `internal/circuit/breaker.go`:
   - States: CLOSED → HALF_OPEN → OPEN
   - Triggers: 3 loops no progress, 5 loops same error
   - `RecordResult(loopNum, filesChanged, hasErrors)`
   - `ShouldHalt() bool`

7. Create `internal/analysis/response.go`:
   - Detect JSON vs text format
   - Parse `RALPH_STATUS` blocks
   - Detect completion keywords, test-only loops
   - Calculate confidence score

8. Create `internal/loop/controller.go`:
   - Main loop: check rate limit → check circuit → execute → analyze
   - Graceful exit detection (test saturation, completion signals)
   - Session management

### Phase 3: TUI
9. Create `internal/tui/model.go`:
   - States: `initializing`, `running`, `paused`, `complete`, `error`
   - Bubble Tea `Init()`, `Update()`, `View()`

10. Create `internal/tui/views/`:
    - `status.go`: loop count, calls made, circuit state
    - `logs.go`: scrollable viewport with recent logs
    - `help.go`: keybinding reference

11. Create `cmd/ralph/main.go`:
    - Parse flags: `--continue`, `--timeout`, `--calls`, `--verbose`
    - Subcommands: `setup`, `import`, `status`, `reset-circuit`
    - Run Bubble Tea program

### Phase 4: Project Commands
12. Create `internal/project/setup.go`:
    - Create dirs: `specs/`, `src/`, `logs/`, `docs/generated/`
    - Copy templates, git init

13. Create `internal/project/import.go`:
    - PRD → PROMPT.md + @fix_plan.md conversion
    - Call Codex to transform content

## Tests (MUST RUN)
- `go test ./...`
- `go test -tags=integration ./...`
- `golangci-lint run`

## Gating criteria
- [ ] All packages compile
- [ ] Unit tests pass
- [ ] TUI renders in iTerm2, Terminal.app
- [ ] Can run one Codex loop successfully

## Acceptance check
- `ralph` launches TUI, displays status
- `ralph setup myproject` creates project structure
- Loop executes Codex and updates display
- Circuit breaker halts on stagnation

## Rollback note
- Keep shell scripts in repo until TUI is stable
- Add `--legacy` flag to fall back to shell if needed
