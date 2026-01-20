# Ralph Codex Skills

Skills for the Ralph autonomous development loop.

## Ralph-Specific Skills

### /ralph-init

Initialize a Ralph project from PRD.md, specs/ folder, or REFACTOR.md.

```bash
ralph-init [--mode <implementation|fix|refactor>] [--verbose]
```

**Auto-detection:**
- If `REFACTOR.md` exists → refactor mode
- If `specs/` folder exists with .md files → fix mode
- If `PRD.md` exists → implementation mode

**Implementation mode** creates:
- `IMPLEMENTATION_PLAN.md` - Task checklist from PRD
- `AGENTS.md` - Project guidance and tech stack

**Fix mode** creates:
- `@fix_plan.md` - Prioritized fixes from specs

**Refactor mode** creates:
- `REFACTOR_PLAN.md` - Phased refactoring tasks

### /ralph-run

Start the Ralph autonomous development loop.

```bash
ralph-run [--monitor] [--calls <n>] [--verbose]
```

**Options:**
- `--monitor` - Enable TUI monitoring interface
- `--calls <n>` - Maximum loop iterations (default: 3)
- `--verbose` - Verbose output

### /ralph-status

Show current Ralph project status including:
- Project type (implementation, fix, or refactor mode)
- Task completion progress
- Circuit breaker state
- Active session info
- Recent log activity

```bash
ralph-status
```

### /ralph-reset

Reset Ralph circuit breaker and session state.

```bash
ralph-reset [--all]
```

**Options:**
- `--all` - Also clear logs and status.json

## Generic Skills

### /lint

Run lint and typecheck for the detected project type.

Supports: Go, Node.js, Python, Rust

### /run-tests

Run tests for the detected project type.

Supports: Go, Node.js, Python, Rust

## Workflow

### New Project (Implementation)

1. Create `PRD.md` with project requirements
2. Run `/ralph-init` to generate plan and agents
3. Run `/ralph-run --monitor` to start development loop

### Existing Project (Fix)

1. Add specs to `specs/` folder
2. Run `/ralph-init --mode fix` to generate fix plan
3. Run `/ralph-run --monitor` to start development loop

### Code Refactoring

1. Create `REFACTOR.md` with refactoring goals and scope
2. Run `/ralph-init --mode refactor` to generate refactor plan
3. Run `/ralph-run --monitor` to start refactoring loop

## Project Structure

### Implementation Mode
```
project/
├── PRD.md                    # Product Requirements (input)
├── IMPLEMENTATION_PLAN.md    # Generated task checklist
├── AGENTS.md                 # Generated project guidance
└── src/                      # Source code
```

### Fix Mode
```
project/
├── specs/                    # Specifications (input)
│   ├── api.md
│   └── ...
├── @fix_plan.md              # Generated fix checklist
├── PROMPT.md                 # Development instructions
└── src/                      # Source code
```

### Refactor Mode
```
project/
├── REFACTOR.md               # Refactoring goals (input)
├── REFACTOR_PLAN.md          # Generated phased plan
└── src/                      # Source code
```

## REFACTOR.md Format

Example structure for `REFACTOR.md`:

```markdown
# Refactoring Goals

## Scope
- Which parts of the codebase to refactor
- Files or modules to focus on

## Goals
- What improvements to achieve
- Technical debt to address
- Patterns to introduce or remove

## Constraints
- Backwards compatibility requirements
- Performance requirements
- Test coverage requirements

## Out of Scope
- What NOT to change
- Features to preserve as-is
```
