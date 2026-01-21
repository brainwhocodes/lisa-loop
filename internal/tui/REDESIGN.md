# TUI Redesign Plan

## Current Issues
1. No live streaming of Codex output during task execution
2. Activity log is limited and doesn't show actual work being done
3. Tasks show checkboxes but no real-time progress
4. No visibility into what Codex is thinking/doing

## Proposed Layout

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  Ralph Codex                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│  [refactor] RUNNING  Loop: 2  Calls: 3/10 [████████░░]  Circuit: CLOSED     │
├──────────────────────────────────┬──────────────────────────────────────────┤
│  Tasks (3/8 complete)            │  Live Output                             │
│  ────────────────────            │  ──────────                              │
│  [x] Extract shared utilities    │  > Reading src/utils/icsParser.ts...     │
│  [x] Add unit tests for parser   │  > Found 3 functions to refactor         │
│  [▸] Refactor paymentStore  ◀──  │  > Creating src/utils/paymentUtils.ts    │
│  [ ] Extract calendar hooks      │  > Writing extractAmount function...     │
│  [ ] Move inline styles to CSS   │  > Adding tests for new utilities...     │
│  [ ] Add PaymentService          │                                          │
│  [ ] Integration tests           │  [Reasoning]                             │
│  [ ] Documentation               │  I'll extract the amount parsing logic   │
│                                  │  into a separate utility function...     │
├──────────────────────────────────┴──────────────────────────────────────────┤
│  r Run  p Pause  l Logs  c Circuit  t Tasks  ? Help  q Quit                 │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Key Features

### 1. Split Pane Layout
- **Left pane**: Task list with current task highlighted
- **Right pane**: Live streaming output from Codex

### 2. Live Output Streaming
- Parse `codex exec --json` JSONL output in real-time
- Show `reasoning` items in a dimmed/italic style
- Show `agent_message` items as main output
- Show tool calls (file reads, writes, etc.)

### 3. Task Progress
- `[ ]` - Pending
- `[▸]` - In progress (with spinner)
- `[x]` - Completed
- `[!]` - Failed/blocked

### 4. Output Types to Display
From codex exec --json:
- `item.completed` with `type: "reasoning"` → Show in reasoning section
- `item.completed` with `type: "agent_message"` → Show in main output
- Tool usage events → Show as "Reading file...", "Writing file...", etc.

### 5. View Modes
- **Default**: Split view (tasks + output)
- **l**: Full logs view
- **t**: Full tasks view
- **o**: Full output view (maximize right pane)

### 6. Status Bar Improvements
- Show current mode badge: `[implement]`, `[refactor]`, `[fix]`
- Show prompt file being used
- Show elapsed time for current task

## Implementation Steps

1. **Create new layout components**
   - `renderSplitView()` - Main split pane layout
   - `renderTaskPane()` - Left pane with tasks
   - `renderOutputPane()` - Right pane with live output

2. **Add output streaming**
   - Create channel for Codex output events
   - Parse JSONL in controller
   - Send events to TUI via messages

3. **New message types**
   - `CodexOutputMsg` - Raw output line from Codex
   - `CodexReasoningMsg` - Reasoning/thinking output
   - `CodexToolCallMsg` - Tool usage (read, write, exec)
   - `TaskProgressMsg` - Task started/completed

4. **Model updates**
   - Add `outputLines []string` for live output
   - Add `reasoningLines []string` for reasoning
   - Add `currentTaskOutput string` for current task's output
   - Add `viewMode string` for toggling views

5. **Controller integration**
   - Pipe codex exec stdout to TUI
   - Parse JSONL events
   - Send appropriate messages

## File Changes Needed

- `internal/tui/model.go` - Add new fields, update View()
- `internal/tui/messages.go` - New message types (create new file)
- `internal/tui/views.go` - New view rendering functions (create new file)
- `internal/tui/program.go` - Update initialization
- `internal/loop/controller.go` - Add output streaming to TUI
