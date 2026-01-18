# Commit 09 — TUI polish, shell removal & documentation

## Intent
- Polish the Charm TUI with styling, error handling, and UX improvements.
- Remove deprecated shell scripts.
- Add comprehensive tests and update documentation.

## Shell Scripts to Remove

After TUI is stable, delete these files:

```
ralph_loop.sh           # Replaced by internal/loop/
ralph_monitor.sh        # Replaced by internal/tui/
ralph_import.sh         # Replaced by internal/project/import.go
setup.sh                # Replaced by internal/project/setup.go
install.sh              # Replaced by go install / Makefile
uninstall.sh            # Replaced by rm binary
lib/circuit_breaker.sh  # Replaced by internal/circuit/
lib/response_analyzer.sh # Replaced by internal/analysis/
lib/date_utils.sh       # Replaced by Go stdlib time package
```

## Scope / touched areas
- `internal/tui/**` (styling, error states, help screen)
- `internal/codex/**` (error handling, retry logic)
- `Makefile` (build, install, test, lint)
- `README.md`, `docs/`
- Delete: all shell scripts listed above

## Steps (atomic)

### Phase 1: TUI Polish
1. Create `internal/tui/styles.go`:
   - Define color palette using Lip Gloss
   - Status bar styles (running=green, paused=yellow, error=red)
   - Box borders, headers, dividers

2. Enhance status view:
   - Animated spinner during Codex execution
   - Progress bar for rate limit countdown
   - Circuit breaker state indicator

3. Add interactive forms with Huh:
   - `ralph setup` → project name prompt, template selection
   - `ralph import` → file picker, confirmation
   - Settings editor (calls/hour, timeout, etc.)

4. Improve error display:
   - Error panel with stack trace (collapsible)
   - Retry button for transient failures
   - "Codex not found" with install instructions

### Phase 2: Keybindings & Help
5. Finalize keybindings:
   - `q` / `ctrl+c` — quit
   - `r` — run next loop
   - `p` — pause/resume
   - `l` — toggle log panel
   - `c` — show circuit breaker status
   - `R` — reset circuit breaker
   - `?` — help screen
   - `tab` — cycle focus between panels

6. Create help screen:
   - Keybinding reference
   - Version info (`ralph --version`)
   - Links to docs

### Phase 3: Testing
7. Unit tests for each package:
   - `internal/codex/runner_test.go` — mock `exec.Command`
   - `internal/codex/jsonl_test.go` — parse fixtures
   - `internal/circuit/breaker_test.go` — state transitions
   - `internal/analysis/response_test.go` — JSON + text parsing
   - `internal/loop/ratelimit_test.go` — counter logic

8. TUI tests with `teatest`:
   - Model state transitions
   - Keybinding responses
   - View rendering snapshots

9. Integration tests:
   - Full loop with mocked Codex binary
   - Project setup creates expected files
   - Circuit breaker trips after N failures

### Phase 4: Build & Install
10. Create `Makefile`:
    ```makefile
    build:
        go build -o ralph ./cmd/ralph
    
    install:
        go install ./cmd/ralph
    
    test:
        go test ./...
    
    test-integration:
        go test -tags=integration ./...
    
    lint:
        golangci-lint run
    
    clean:
        rm -f ralph
    ```

11. Add `.goreleaser.yml` for cross-platform binaries:
    - macOS (amd64, arm64)
    - Linux (amd64, arm64)
    - Windows (amd64)

### Phase 5: Documentation
12. Update `README.md`:
    - New installation: `go install` or download binary
    - Quick start with TUI commands
    - Remove all shell script references

13. Create `docs/tui.md`:
    - Keybindings reference
    - Configuration options
    - Troubleshooting guide

14. Update `docs/codex.md`:
    - Remove shell-specific instructions
    - Document Go package usage for developers

### Phase 6: Cleanup
15. Delete shell scripts:
    - Remove files listed above
    - Update `.gitignore` if needed
    - Remove `node_modules/bats` (no longer needed for shell tests)

16. Archive shell tests:
    - Move `tests/test_*.sh` to `tests/legacy/` or delete
    - Keep only Go tests

## Tests (MUST RUN)
- `go test ./...`
- `go test -tags=integration ./...`
- `golangci-lint run`

## Gating criteria
- [ ] All unit tests pass (>80% coverage)
- [ ] Integration tests pass
- [ ] Lint passes with no warnings
- [ ] TUI works on macOS and Linux
- [ ] README is updated
- [ ] Shell scripts removed

## Acceptance check
- `go install ./cmd/ralph` works
- `ralph` launches polished TUI
- `ralph setup`, `ralph import` work via TUI forms
- Help screen shows all keybindings
- No shell script dependencies remain
- CI builds cross-platform binaries

## Rollback note
- Keep shell scripts in `legacy/` branch if needed for emergency fallback
