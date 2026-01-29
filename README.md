# Lisa Codex - AI Development Loop with TUI

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Version](https://img.shields.io/badge/version-1.1.0-blue)
[![GitHub Issues](https://img.shields.io/github/issues/brainwhocodes/lisa-loop)](https://github.com/brainwhocodes/lisa-loop/issues)

> **Autonomous AI development loop with Charm TUI and intelligent exit detection**

Lisa Codex is a modern Go implementation of Geoffrey Huntley's autonomous development technique, featuring a beautiful terminal user interface powered by Charm libraries. It enables continuous autonomous development cycles where Codex iteratively improves your project until completion, with built-in safeguards to prevent infinite loops and API overuse.

**Features**:
- 🎨 **Modern TUI** - Beautiful terminal interface with real-time task progress and logs
- 🔄 **Dual Backend Support** - Codex CLI or OpenCode server backend
- 📋 **Preflight Checks** - Validates plan status before each loop iteration
- ⚡ **Session Continuity** - Preserve context across loop iterations
- 🛡️ **Circuit Breaker** - Prevent runaway loops with advanced error detection
- 📊 **Real-time Monitoring** - Live status, task checklist, and integrated logs
- 🎯 **Intelligent Exit** - Task-based completion with automatic detection
- ✅ **Flexible Task Formats** - Supports `- [ ]`, `* [ ]`, `1. [ ]`, and `[ ]` checklists

## Quick Start

### Install Lisa

```bash
# Build from source
git clone https://github.com/brainwhocodes/lisa-loop.git
cd lisa-loop
make build
make install
```

### Initialize a Project

Lisa supports three initialization modes based on your project setup:

#### Implementation Mode (New Projects)

Start with a Product Requirements Document to generate an implementation plan:

```bash
# Create your PRD
echo "# My Project\n\nBuild a CLI tool that..." > PRD.md

# Initialize - generates IMPLEMENTATION_PLAN.md and AGENTS.md
lisa init

# Start autonomous development
lisa --monitor
```

#### Fix Mode (Existing Projects)

Start with specification documents to generate a fix plan:

```bash
# Add specs to specs/ folder
mkdir specs
echo "# API Spec\n\n..." > specs/api.md

# Initialize - generates @fix_plan.md
lisa init --mode fix

# Start autonomous development
lisa --monitor
```

#### Refactor Mode (Code Refactoring)

Start with a refactoring goals document to generate a phased refactoring plan:

```bash
# Create REFACTOR.md with refactoring goals
cat > REFACTOR.md << 'EOF'
# Refactoring Goals

## Scope
- Refactor the API layer for better testability
- Extract shared utilities into a common package

## Goals
- Improve code coverage to 80%
- Reduce cyclomatic complexity
- Add dependency injection

## Constraints
- Maintain backwards compatibility
- No breaking API changes
EOF

# Initialize - generates REFACTOR_PLAN.md
lisa init --mode refactor

# Start autonomous refactoring
lisa --monitor
```

### Running Lisa

```bash
# Basic run (3 iterations)
ralph

# With TUI monitoring
lisa --monitor

# Custom iteration limit
lisa --monitor --calls 5

# Use OpenCode backend
lisa --monitor --backend opencode --opencode-url http://localhost:3000

# Structured logging (JSON)
lisa --log-format json

# Verbose output
lisa --monitor --verbose
```

### Backend Selection

Lisa supports two backends for AI execution:

#### Codex CLI (Default)
Uses the local Codex CLI for autonomous development:

```bash
# Default - uses Codex CLI
lisa --monitor
```

**Requirements:**
- Install Codex CLI: `npm install -g @anthropic/claude-code`
- Authenticate with `codex auth`

#### OpenCode Server
Connect to a self-hosted or cloud OpenCode server:

```bash
# Using flags
lisa --monitor --backend opencode --opencode-url http://localhost:3000 --opencode-pass mypassword

# Using environment variables
export OPENCODE_SERVER_URL=http://localhost:3000
export OPENCODE_SERVER_PASSWORD=mypassword
lisa --monitor --backend opencode
```

**Environment Variables:**
| Variable | Description | Default |
|----------|-------------|---------|
| `OPENCODE_SERVER_URL` | Server URL | - |
| `OPENCODE_SERVER_USERNAME` | Auth username | `opencode` |
| `OPENCODE_SERVER_PASSWORD` | Auth password | - |
| `OPENCODE_MODEL_ID` | Model ID | `glm-4.7` |

### Preflight Checks

Before each loop iteration, Lisa performs preflight checks:

1. **Plan Status** - Verifies remaining tasks in the plan file
2. **Circuit Breaker** - Checks if the circuit is OPEN (too many errors)
3. **Rate Limit** - Ensures API calls haven't exceeded the limit
4. **Max Loops** - Checks if iteration limit has been reached

If any check fails, the loop skips the backend call and exits with a clear reason:
```
Skipped: All tasks complete
Skipped: Circuit breaker is OPEN
Skipped: Rate limit exhausted (0 calls remaining)
```

### Legacy Project Setup

```bash
# Import existing PRD/specification
lisa --command import --source my-requirements.md --import-name my-project

# Create blank project with templates
lisa --command setup --name my-project
```

## TUI Keybindings

### Navigation
- `q` / `Ctrl+C` / `Ctrl+Q` - Quit Lisa Codex
- `?` - Toggle help screen

### Loop Control
- `r` - Run / Restart loop
- `p` - Pause / Resume loop

### Views
- `l` - Toggle log view
- `t` - Toggle tasks view
- `o` - Toggle output view
- `c` - Show circuit breaker status
- `R` - Reset circuit breaker

The TUI displays:
- **Header** - Mode, loop number, task progress
- **Status Bar** - Current state, circuit breaker status, context usage
- **Task Panel** - Current phase tasks with completion status
- **Output Panel** - Live agent output and reasoning

## Commands

### run (default)

Run the autonomous development loop.

```bash
lisa                           # Run with defaults (3 iterations)
lisa run                       # Explicit run command
lisa --monitor                 # With TUI interface
lisa --monitor --calls 5       # Custom iteration limit
lisa --verbose                 # Verbose output
```

**Options:**
| Option | Description | Default |
|--------|-------------|---------|
| `--project <path>` | Project directory | `.` |
| `--prompt <file>` | Prompt file | `PROMPT.md` |
| `--calls <n>` | Max loop iterations | `3` (10 for opencode) |
| `--timeout <sec>` | Codex timeout | `600` |
| `--monitor` | Enable TUI monitoring | `false` |
| `--verbose` | Verbose output | `false` |
| `--backend` | Backend: `cli` or `opencode` | `cli` |
| `--opencode-url` | OpenCode server URL | - |
| `--opencode-user` | OpenCode username | `opencode` |
| `--opencode-pass` | OpenCode password | - |
| `--opencode-model` | OpenCode model ID | `glm-4.7` |
| `--log-format` | Log format: `text`, `json`, `logfmt` | `text` |

### init

Initialize a Lisa project from PRD.md, specs/, or REFACTOR.md.

```bash
lisa init                        # Auto-detect mode
lisa init --mode implementation  # Force implementation mode
lisa init --mode fix             # Force fix mode
lisa init --mode refactor        # Force refactor mode
lisa init --verbose              # Verbose output
```

### setup

Create a new Lisa-managed project.

```bash
lisa setup --name my-project                    # Create new project
lisa setup --name my-project --description "…"  # With description for Codex
lisa setup --init                               # Initialize in current directory
lisa setup --name my-project --git=false        # Skip git init
```

**Options:**
| Option | Description |
|--------|-------------|
| `--name <name>` | Project name (required unless --init) |
| `--description <text>` | Project description for Codex |
| `--init` | Initialize in current directory |
| `--git` | Initialize git repository (default: true) |

### import

Import a PRD or specification document.

```bash
lisa import --source requirements.md                    # Auto-detect project name
lisa import --source spec.md --import-name my-project   # Specify project name
```

**Options:**
| Option | Description |
|--------|-------------|
| `--source <file>` | Source file to import (required) |
| `--import-name <name>` | Project name (auto-detected if empty) |

### status

Show current project status.

```bash
lisa status
```

### reset-circuit

Reset the circuit breaker state.

```bash
lisa reset-circuit
```

### help / version

```bash
lisa help      # Show help
lisa version   # Show version
```

## Skills

Lisa includes Codex skills for common development tasks. Skills are bash scripts installed to `~/.codex/skills/`.

### Lisa-Specific Skills

| Skill | Description |
|-------|-------------|
| `/ralph-init` | Initialize project from PRD.md, specs/, or REFACTOR.md |
| `/ralph-run` | Start the autonomous development loop |
| `/ralph-status` | Show project status and progress |
| `/ralph-reset` | Reset circuit breaker and session state |

### Generic Skills

| Skill | Description |
|-------|-------------|
| `/lint` | Run lint/typecheck (Go, Node, Python, Rust) |
| `/run-tests` | Run tests (Go, Node, Python, Rust) |

### Installing Skills

Skills are included in the `.codex/skills/` directory:

```bash
# Install skills globally
cp .codex/skills/* ~/.codex/skills/

# Or install manually
make install-skills
```

### Using Skills in Codex

In a Codex session, invoke skills with:

```
/ralph-status
/lint
/run-tests
```

Skills auto-detect project type and run appropriate commands.

## Circuit Breaker States

### CLOSED (Green)
Normal operation - all loop iterations execute

### HALF_OPEN (Yellow)
Monitoring mode - may pause if stagnation continues

### OPEN (Red)
Loop execution halted - press `R` to reset

## Features

- **Dual Backend Support** - Choose between Codex CLI (default) or OpenCode server backend
- **Preflight Validation** - Checks plan status, circuit breaker, and rate limits before each iteration
- **Task-Based Completion** - Automatically detects when all tasks are complete (supports multiple checklist formats)
- **Session Continuity** - Preserves context across loop iterations with automatic session management
- **Loop Management** - Built-in iteration limits with configurable max loops
- **Live Monitoring** - Real-time dashboard showing loop status, task progress, and integrated logs
- **Task Management** - Structured approach with prioritized task lists and progress tracking
- **Circuit Breaker** - Advanced error detection with three states (CLOSED/HALF_OPEN/OPEN) and automatic recovery
- **Structured Logging** - JSON and logfmt output formats for integration with logging systems
- **Event-Driven Architecture** - Preflight and outcome events for monitoring and observability

## Building from Source

```bash
# Build binary
make build

# Install to GOPATH
make install

# Run tests
make test

# Generate coverage report
make test-coverage

# Run linter
make lint
```

## Documentation

- [docs/tui.md](docs/tui.md) - TUI documentation and keybindings
- [docs/codex.md](docs/codex.md) - Codex integration guide
- [docs/opencode.md](docs/opencode.md) - OpenCode backend setup guide
- [AGENTS.md](AGENTS.md) - Agent development guidelines
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contributor guide
- [TESTING.md](TESTING.md) - Testing documentation

## Project Structure

### Implementation Mode (PRD-based)

```
my-project/
├── PRD.md                    # Product Requirements Document (input)
├── IMPLEMENTATION_PLAN.md    # Generated task checklist
├── AGENTS.md                 # Generated project guidance
├── src/                      # Source code
├── logs/                     # Execution logs
└── docs/                     # Documentation
```

### Fix Mode (Specs-based)

```
my-project/
├── specs/                    # Specification documents (input)
│   ├── api.md
│   ├── architecture.md
│   └── ...
├── @fix_plan.md              # Generated fix checklist
├── PROMPT.md                 # Development instructions
├── src/                      # Source code
├── logs/                     # Execution logs
└── docs/                     # Documentation
```

### Refactor Mode

```
my-project/
├── REFACTOR.md               # Refactoring goals (input)
├── REFACTOR_PLAN.md          # Generated phased refactoring plan
├── src/                      # Source code
├── logs/                     # Execution logs
└── docs/                     # Documentation
```

## System Requirements

- **Go 1.21+** - For building Lisa
- **Codex CLI** - `npm install -g @anthropic/claude-code` (for CLI backend)
- **OR OpenCode Server** - For OpenCode backend (self-hosted or cloud)
- **Git** - Version control
- **Standard Unix tools** - grep, date, etc.

### Backend Requirements

**Codex CLI Backend:**
- Codex CLI installed and authenticated
- OpenAI API access

**OpenCode Backend:**
- OpenCode server URL
- Authentication credentials
- Compatible model (default: glm-4.7)

## Testing

Run tests with Makefile:
```bash
make test          # All tests
make test-verbose  # Verbose output
make test-coverage # Coverage report
make lint          # Code linting
```

### Test Structure

```
├── internal/
│   ├── analysis/       # Response analysis tests
│   ├── circuit/        # Circuit breaker tests
│   ├── codex/          # Codex integration tests
│   ├── loop/           # Loop controller & preflight tests
│   ├── opencode/       # OpenCode client tests
│   ├── project/        # Project management tests
│   ├── runner/         # Backend runner tests
│   ├── state/          # State persistence tests
│   └── tui/            # TUI component tests
├── tests/
│   ├── fixtures/       # Test fixtures for all modes
│   │   ├── fix/
│   │   ├── implement/
│   │   └── refactor/
│   └── e2e_loop_test.go # End-to-end loop tests
```

### End-to-End Tests

The E2E tests verify the complete loop execution:

```bash
go test ./tests/... -v
```

Tests include:
- **Fix Mode** - Runs loop until 3 tasks complete
- **Implement Mode** - Runs loop until 6 tasks complete
- **Refactor Mode** - Runs loop until 4 tasks complete
- **Early Exit** - Verifies loop skips when all tasks already complete
- **Event Capture** - Validates preflight and outcome events

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for complete contributor guide.

## License

This project is licensed under MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Lisa technique](https://ghuntley.com/ralph/) created by Geoffrey Huntley
- Powered by [Codex](https://openai.com/codex) and [OpenCode](https://github.com/opencode-ai)
- TUI built with [Charm libraries](https://charm.sh/) (Bubble Tea, Lipgloss)
- Community feedback and contributions

## Related Projects

- [Codex](https://openai.com/codex) - The AI coding assistant that powers Lisa
- [Aider](https://github.com/paul-gauthier/aider) - AI pair programming tool
