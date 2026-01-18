# Codex Backend

Ralph uses OpenAI's Codex CLI (`codex exec`) and Codex SDK (`@openai/codex-sdk`) as the AI agent backend for autonomous development loops.

## Overview

Ralph is an autonomous development loop system that uses OpenAI's Codex as its AI backend. It provides two integration options:

- **Codex CLI backend** - Uses `codex exec` for command-line based execution with JSONL streaming
- **Codex SDK backend** - Uses `@openai/codex-sdk` for programmatic, resumable sessions with fine-grained control

The system is designed to continuously iterate on your project until completion, with intelligent exit detection, rate limiting, and circuit breaker patterns to prevent infinite loops.

## Codex CLI Backend (Default)

The Codex CLI backend is the default execution method for Ralph. It uses the `codex exec` command with JSONL streaming output for real-time parsing of AI responses.

### Installation

Install the Codex CLI globally:

```bash
npm install -g @openai/codex
```

Verify installation:

```bash
codex --version
```

### Authentication

The Codex CLI uses the `CODEX_API_KEY` environment variable for authentication:

```bash
export CODEX_API_KEY="your-api-key-here"
```

Add to your shell profile for persistence:

```bash
# For bash
echo 'export CODEX_API_KEY="your-api-key-here"' >> ~/.bashrc
source ~/.bashrc

# For zsh
echo 'export CODEX_API_KEY="your-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

**Important**: `CODEX_API_KEY` only works with `codex exec` commands. The Codex SDK backend uses a different authentication method (see below).

### Basic Usage

```bash
# Run Ralph with Codex CLI backend (default)
ralph --monitor

# Explicitly select CLI backend
ralph --backend cli --monitor

# With custom prompt file
ralph --backend cli --prompt my-prompt.md

# With rate limiting
ralph --backend cli --calls 50 --monitor

# With verbose output
ralph --backend cli --verbose --monitor
```

### JSONL Output Parsing

The Codex CLI emits JSONL (JSON Lines) streaming output when invoked with the `--json` flag:

```bash
codex exec --json --prompt "your prompt"
```

Each line is a JSON object representing an event, tool call, or response. Ralph parses this JSONL stream to:

- Extract tool usage (file edits, bash commands, etc.)
- Track progress and completion indicators
- Detect errors and circuit breaker conditions
- Maintain session continuity

**Why JSONL?**
- Streamable parsing (process events as they arrive)
- Structured data for reliable state tracking
- Backward-compatible with Ralph's existing JSON-based analyzers

### Resume Support

Codex CLI supports session resumption:

```bash
codex exec --json --prompt "your prompt"
# ... loop iterations ...
codex exec --json --resume --thread-id <thread-id>
```

Ralph automatically handles session persistence and resumption when using `--agent codex-cli` with the `--continue` flag (enabled by default).

### Git Repository Requirements

Codex CLI requires running inside a git repository unless the `--skip-git-repo-check` flag is provided.

Ralph automatically adds this flag when needed:

```bash
# Codex automatically checks for git repo
ralph  # Will work in git-initialized projects

# For non-git projects, Ralph adds --skip-git-repo-check automatically
```

### Shell Quoting Safety

When passing large prompt files (like `PROMPT.md`) to `codex exec`, Ralph takes special care to avoid shell-quoting pitfalls:

```bash
# Safe: Use heredoc or process substitution
codex exec --json --prompt "$(cat PROMPT.md)"

# Unsafe: Direct string expansion with special characters
codex exec --json --prompt "$(< PROMPT.md)"  # May break on special chars
```

Ralph uses process substitution with proper escaping to ensure prompt content is passed correctly.

## Codex SDK Backend

The Codex SDK backend provides programmatic control over Codex with direct API access, enabling more advanced features and better error handling than the CLI wrapper.

### Installation

The Codex SDK requires Node.js 18+:

```bash
# Install in your project
npm install @openai/codex-sdk

# Or install globally if Ralph will use it system-wide
npm install -g @openai/codex-sdk
```

Verify Node.js version:

```bash
node --version  # Should be 18.0.0 or higher
```

### Authentication

The Codex SDK uses OpenAI's standard API key authentication:

```bash
export OPENAI_API_KEY="your-api-key-here"
```

Add to your shell profile:

```bash
# For bash
echo 'export OPENAI_API_KEY="your-api-key-here"' >> ~/.bashrc
source ~/.bashrc

# For zsh
echo 'export OPENAI_API_KEY="your-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

> **Note**: The Codex SDK's `OPENAI_API_KEY` is separate from the CLI's `CODEX_API_KEY`. You may need both depending on which backend you use.

**Note**: The Codex SDK authentication is separate from the Codex CLI's `CODEX_API_KEY`.

### Basic Usage

```bash
# Run Ralph with Codex SDK backend
ralph --backend sdk --monitor

# With custom prompt
ralph --backend sdk --prompt my-prompt.md

# With rate limiting
ralph --backend sdk --calls 50 --monitor
```

### SDK API Overview

Ralph's Codex SDK backend uses these key components:

```javascript
import { Codex } from '@openai/codex-sdk';

// Initialize client
const codex = new Codex({
  apiKey: process.env.OPENAI_API_KEY
});

// Start a new session/thread
const thread = await codex.threads.create();

// Run an iteration
const run = await codex.threads.run(thread.id, {
  messages: [{ role: 'user', content: prompt }]
});

// Resume an existing session
const resumedRun = await codex.threads.resume(thread.id);

// Wait for completion
await codex.threads.wait(thread.id);
```

### Thread Persistence

When using the Codex SDK backend, Ralph:

1. Creates a new thread on first run
2. Persists the thread ID to `.ralph_thread_id`
3. Automatically resumes the thread on subsequent iterations
4. Manages thread lifecycle (reset on circuit breaker, completion, etc.)

Thread state is stored separately from Ralph's session state:

- `.ralph_session` - Ralph loop state (rate limits, circuit breaker, etc.)
- `.ralph_thread_id` - Codex SDK thread identifier
- `.ralph_thread_history` - Thread lifecycle events for debugging

### Advantages of SDK Backend

Compared to the CLI backend, the SDK backend offers:

- **Programmatic control** - Direct API access without subprocess overhead
- **Better error handling** - Structured exceptions vs. parsing stderr
- **Thread management** - Fine-grained control over thread lifecycle
- **Streaming support** - Event-driven progress tracking
- **Type safety** - TypeScript types for all API calls

### Disadvantages

- Requires Node.js 18+ runtime
- More complex integration (vs. simple CLI wrapper)
- SDK version compatibility concerns

## Comparison: Codex CLI vs Codex SDK

| Feature | Codex CLI | Codex SDK |
|---------|-----------|-----------|
| **Installation** | `npm install -g @openai/codex` | `npm install @openai/codex-sdk` |
| **Auth** | `CODEX_API_KEY` | `OPENAI_API_KEY` |
| **Node.js** | Not required (CLI is binary) | Required (18+) |
| **Execution** | Subprocess (`codex exec`) | Programmatic API calls |
| **Output** | JSONL streaming | Structured API responses |
| **Resume** | `--resume --thread-id` | `thread.resume()` |
| **Error Handling** | Parse stderr/stdout | JavaScript exceptions |
| **Thread Mgmt** | CLI-managed | SDK-controlled |
| **Use Case** | Simple CLI wrapper | Advanced, programmatic |

## Migration Path

### Default Behavior

Ralph uses Codex CLI as the default backend:

```bash
ralph  # Uses Codex CLI backend by default
ralph --monitor  # With integrated monitoring
```

### Backend Selection

To select between Codex CLI and Codex SDK, use the `--backend` flag:

1. **Install dependencies**:
   ```bash
   # For Codex CLI backend
   npm install -g @openai/codex
   export CODEX_API_KEY="your-key"

   # For Codex SDK backend
   npm install @openai/codex-sdk
   export OPENAI_API_KEY="your-key"
   ```

2. **Use the appropriate backend flag**:
   ```bash
   ralph --backend cli --monitor      # Codex CLI
   ralph --backend sdk --monitor      # Codex SDK
   ```

3. **Test with a sample project**:
   ```bash
   ralph-setup test-codex
   cd test-codex
   ralph --monitor
   ```

## Template Support

Ralph's project templates include Codex-specific configuration:

### Standard Ralph Template

```
my-project/
├── PROMPT.md              # Development instructions
├── @fix_plan.md           # Task priorities
├── @AGENT.md              # Build/run instructions
├── specs/                 # Project specifications
└── src/                   # Source code
```

### Codex-Enhanced Template (Future)

When Codex support is fully integrated, templates will include:

```
my-project/
├── PROMPT.md              # Development instructions
├── @fix_plan.md           # Task priorities
├── @AGENT.md              # Build/run instructions
├── AGENTS.md              # Codex-specific guidance
├── .codex/                # Codex configuration (future)
│   ├── config.json        # Codex CLI/SDK settings
│   └── skills/            # Codex skills pack
├── specs/                 # Project specifications
└── src/                   # Source code
```

## Testing and Mocking

### No Real Codex Calls in CI/Unit Tests

Ralph mocks all Codex CLI and SDK calls in tests to avoid:

- API costs during test runs
- Rate limit issues in CI pipelines
- Flakiness from network dependencies
- Need for test API keys

### Mocking Strategy

**Codex CLI Mocking**:
```bash
# Test helper creates mock JSONL output
mock_codex_exec() {
    # Returns pre-canned JSONL events
    # Simulates tool calls, responses, errors
}
```

**Codex SDK Mocking**:
```javascript
// Test helper stubs SDK calls
mockCodexRun() {
    return {
        threadId: "mock-thread-123",
        messages: [...],
        status: "completed"
    };
}
```

### Testing Approach

Tests verify:

1. **Backend selection** - `--backend cli|sdk` correctly selects Codex variant
2. **Command construction** - Proper `codex exec` flags and quoting
3. **JSONL parsing** - Correct extraction of events from streaming output
4. **Thread persistence** - Thread IDs stored and retrieved correctly
5. **Resume logic** - Sessions resumed with correct state
6. **Error handling** - Circuit breaker triggered on Codex errors

See `CONTRIBUTING.md` for detailed testing guidelines.

## Troubleshooting

### Common Issues

**"codex: command not found"**
```bash
# Install Codex CLI
npm install -g @openai/codex
```

**"CODEX_API_KEY not set"**
```bash
export CODEX_API_KEY="your-key"
# Add to ~/.bashrc or ~/.zshrc for persistence
```

**"Thread not found" (SDK)**
```bash
# Thread may have expired or been deleted
ralph --reset-session  # Start fresh thread
```

**JSONL parsing errors**
- Check that prompt doesn't contain invalid JSON
- Verify `--json` flag is being passed
- Review `logs/ralph.log` for detailed parsing errors

**Session continuity issues**
- Check `.ralph_thread_id` exists and is valid
- Verify `--continue` flag is not disabled
- Use `ralph --status` to inspect session state

## Security Considerations

- **Never commit API keys** to version control
- Use environment variables or secret managers
- Rotate keys regularly
- Restrict key permissions (if possible)
- Audit key usage

## Performance Notes

### Codex CLI
- Subprocess overhead per iteration
- JSONL parsing is lightweight
- Resume via `--thread-id` is fast

### Codex SDK
- Direct API calls (no subprocess overhead)
- Persistent connection pooling
- Streaming responses reduce latency

### Rate Limiting

Ralph's rate limiting applies to both Codex backends:

```bash
# Limit API calls
ralph --calls 50 --monitor
```

## Future Enhancements

Planned Codex integration improvements:

- [ ] Skills pack support (custom Codex skills in templates)
- [ ] `.codex/` directory configuration
- [ ] Advanced thread management (parallel threads, etc.)
- [ ] Fine-grained tool permissions per backend
- [ ] Backend-specific exit conditions
- [ ] Per-backend session timeout settings
- [ ] Codex-specific output formatting

## See Also

- [README.md](README.md) - Ralph overview and quick start
- [CONTRIBUTING.md](CONTRIBUTING.md) - Development and testing guidelines
- [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) - Full development roadmap
- [00_overview.md](00_overview.md) - Codex integration project goals
- [CLAUDE.md](CLAUDE.md) - Technical specifications and quality standards
