# Ralph Codex - AI Development Loop with TUI

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Tests](https://img.shields.io/badge/tests-400%2B%20passing-green)
[![GitHub Issues](https://img.shields.io/github/issues/brainwhocodes/ralph-codex)](https://github.com/brainwhocodes/ralph-codex/issues)

> **Autonomous AI development loop with Charm TUI and intelligent exit detection**

Ralph Codex is a modern Go implementation of Geoffrey Huntley's autonomous development technique, featuring a beautiful terminal user interface powered by Charm libraries. It enables continuous autonomous development cycles where Codex iteratively improves your project until completion, with built-in safeguards to prevent infinite loops and API overuse.

**Features**:
- 🎨 **Modern TUI** - Beautiful terminal interface with real-time updates
- 🔄 **Dual Backend Support** - Codex CLI and SDK backends
- ⚡ **Session Continuity** - Preserve context across loop iterations
- 🛡️ **Circuit Breaker** - Prevent runaway loops with advanced error detection
- 📊 **Real-time Monitoring** - Live status and log viewer
- 🎯 **Intelligent Exit** - Semantic understanding of completion signals

## Quick Start

### Install Ralph

```bash
# Build from source
git clone https://github.com/brainwhocodes/ralph-codex.git
cd ralph-codex
make build
make install
```

### Create a New Project

```bash
# Option 1: Import existing PRD/specification
ralph --command import --source my-requirements.md --import-name my-project

# Option 2: Create blank project
ralph --command setup --name my-project

cd my-project

# Start autonomous development with TUI
ralph --command run --monitor
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

## CLI Options

```bash
ralph --command run [OPTIONS]
  --backend <cli|sdk>   Codex backend (default: cli)
  --project <path>        Project directory (default: .)
  --prompt <file>         Prompt file (default: PROMPT.md)
  --calls <number>        Max API calls per hour (default: 100)
  --timeout <seconds>      Codex timeout (default: 600)
  --monitor               Enable integrated TUI monitoring
  --verbose               Verbose output
```

## Project Commands

```bash
ralph --command setup --name <project-name>     # Create new project
ralph --command import --source <file>           # Import PRD/specs
ralph --command status                          # Show project status
ralph --command reset-circuit                   # Reset circuit breaker
```

## Circuit Breaker States

### CLOSED (Green)
Normal operation - all loop iterations execute

### HALF_OPEN (Yellow)
Monitoring mode - may pause if stagnation continues

### OPEN (Red)
Loop execution halted - press `R` to reset

## Features

- **Codex Backend** - Autonomous development loop powered by Codex CLI and Codex SDK
- **Intelligent Exit Detection** - Dual-condition check requiring both completion indicators AND explicit EXIT_SIGNAL
- **Session Continuity** - Preserves context across loop iterations with automatic session management
- **Rate Limiting** - Built-in API call management with hourly limits and countdown timers
- **Live Monitoring** - Real-time dashboard showing loop status, progress, and logs
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

```
my-project/
├── PROMPT.md           # Main development instructions
├── @fix_plan.md        # Prioritized task list
├── @AGENT.md           # Build and run instructions
├── specs/              # Project specifications
├── src/                # Source code
├── examples/           # Usage examples
├── logs/               # Execution logs
└── docs/generated/     # Auto-generated docs
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
make lint         # Code linting
```

Current test status:
- **400+ tests** across all packages
- **100% pass rate** on TUI and model tests

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for complete contributor guide.

## License

This project is licensed under MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Ralph technique](https://ghuntley.com/ralph/) created by Geoffrey Huntley
- Built for [Codex](https://openai.com) by OpenAI
- Powered by [Charm libraries](https://charm.sh/)
- Community feedback and contributions

## Related Projects

- [Codex](https://openai.com) - The AI coding assistant that powers Ralph
- [Aider](https://github.com/paul-gauthier/aider) - Original Ralph technique implementation
