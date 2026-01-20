# Ralph Codex - AI Development Loop with TUI

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Version](https://img.shields.io/badge/version-1.0.0-blue)
[![GitHub Issues](https://img.shields.io/github/issues/brainwhocodes/ralph-codex)](https://github.com/brainwhocodes/ralph-codex/issues)

> **Autonomous AI development loop with Charm TUI and intelligent exit detection**

Ralph Codex is a modern Go implementation of Geoffrey Huntley's autonomous development technique, featuring a beautiful terminal user interface powered by Charm libraries. It enables continuous autonomous development cycles where Codex iteratively improves your project until completion, with built-in safeguards to prevent infinite loops and API overuse.

**Features**:
- 🎨 **Modern TUI** - Beautiful terminal interface with real-time task progress and logs
- 🔄 **Codex Integration** - Powered by Codex CLI
- ⚡ **Session Continuity** - Preserve context across loop iterations
- 🛡️ **Circuit Breaker** - Prevent runaway loops with advanced error detection
- 📊 **Real-time Monitoring** - Live status, task checklist, and integrated logs
- 🎯 **Intelligent Exit** - Task-based completion with automatic detection

## Quick Start

### Install Ralph

```bash
# Build from source
git clone https://github.com/brainwhocodes/ralph-codex.git
cd ralph-codex
make build
make install
```

### Initialize a Project

Ralph supports three initialization modes based on your project setup:

#### Implementation Mode (New Projects)

Start with a Product Requirements Document to generate an implementation plan:

```bash
# Create your PRD
echo "# My Project\n\nBuild a CLI tool that..." > PRD.md

# Initialize - generates IMPLEMENTATION_PLAN.md and AGENTS.md
ralph init

# Start autonomous development
ralph --monitor
```

#### Fix Mode (Existing Projects)

Start with specification documents to generate a fix plan:

```bash
# Add specs to specs/ folder
mkdir specs
echo "# API Spec\n\n..." > specs/api.md

# Initialize - generates @fix_plan.md
ralph init --mode fix

# Start autonomous development
ralph --monitor
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
ralph init --mode refactor

# Start autonomous refactoring
ralph --monitor
```

### Running Ralph

```bash
# Basic run (3 iterations)
ralph

# With TUI monitoring
ralph --monitor

# Custom iteration limit
ralph --monitor --calls 5

# Verbose output
ralph --monitor --verbose
```

### Legacy Project Setup

```bash
# Import existing PRD/specification
ralph --command import --source my-requirements.md --import-name my-project

# Create blank project with templates
ralph --command setup --name my-project
```

## TUI Keybindings

### Navigation
- `q` / `Ctrl+C` - Quit Ralph Codex
- `?` - Toggle help screen

### Loop Control
- `r` - Run / Restart loop
- `p` - Pause / Resume loop

### Views
- `l` - Toggle log view
- `c` - Show circuit breaker status
- `R` - Reset circuit breaker

## Commands

### run (default)

Run the autonomous development loop.

```bash
ralph                           # Run with defaults (3 iterations)
ralph run                       # Explicit run command
ralph --monitor                 # With TUI interface
ralph --monitor --calls 5       # Custom iteration limit
ralph --verbose                 # Verbose output
```

**Options:**
| Option | Description | Default |
|--------|-------------|---------|
| `--project <path>` | Project directory | `.` |
| `--prompt <file>` | Prompt file | `PROMPT.md` |
| `--calls <n>` | Max loop iterations | `3` |
| `--timeout <sec>` | Codex timeout | `600` |
| `--monitor` | Enable TUI monitoring | `false` |
| `--verbose` | Verbose output | `false` |

### init

Initialize a Ralph project from PRD.md, specs/, or REFACTOR.md.

```bash
ralph init                        # Auto-detect mode
ralph init --mode implementation  # Force implementation mode
ralph init --mode fix             # Force fix mode
ralph init --mode refactor        # Force refactor mode
ralph init --verbose              # Verbose output
```

### setup

Create a new Ralph-managed project.

```bash
ralph setup --name my-project                    # Create new project
ralph setup --name my-project --description "…"  # With description for Codex
ralph setup --init                               # Initialize in current directory
ralph setup --name my-project --git=false        # Skip git init
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
ralph import --source requirements.md                    # Auto-detect project name
ralph import --source spec.md --import-name my-project   # Specify project name
```

**Options:**
| Option | Description |
|--------|-------------|
| `--source <file>` | Source file to import (required) |
| `--import-name <name>` | Project name (auto-detected if empty) |

### status

Show current project status.

```bash
ralph status
```

### reset-circuit

Reset the circuit breaker state.

```bash
ralph reset-circuit
```

### help / version

```bash
ralph help      # Show help
ralph version   # Show version
```

## Skills

Ralph includes Codex skills for common development tasks. Skills are bash scripts installed to `~/.codex/skills/`.

### Ralph-Specific Skills

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

- **Codex Backend** - Autonomous development loop powered by Codex CLI
- **Task-Based Completion** - Automatically detects when all tasks in `@fix_plan.md` are complete
- **Session Continuity** - Preserves context across loop iterations with automatic session management
- **Loop Management** - Built-in iteration limits with configurable max loops
- **Live Monitoring** - Real-time dashboard showing loop status, task progress, and integrated logs
- **Task Management** - Structured approach with prioritized task lists and progress tracking
- **Circuit Breaker** - Advanced error detection with two-stage filtering and automatic recovery

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
- [CLAUDE.md](CLAUDE.md) - Development guidelines

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

- **Go 1.21+** - For building Ralph
- **Codex CLI** - `npm install -g @anthropic/claude-code` (for CLI backend)
- **Git** - Version control
- **Standard Unix tools** - grep, date, etc.

## Testing

Run tests with Makefile:
```bash
make test          # All tests
make test-verbose  # Verbose output
make test-coverage # Coverage report
make lint          # Code linting
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for complete contributor guide.

## License

This project is licensed under MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Ralph technique](https://ghuntley.com/ralph/) created by Geoffrey Huntley
- Powered by [Codex](https://openai.com/codex)
- TUI built with [Charm libraries](https://charm.sh/) (Bubble Tea, Lipgloss)
- Community feedback and contributions

## Related Projects

- [Codex](https://openai.com/codex) - The AI coding assistant that powers Ralph
- [Aider](https://github.com/paul-gauthier/aider) - AI pair programming tool
