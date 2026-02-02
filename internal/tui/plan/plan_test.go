package plan

import "testing"

func TestParsePhases_RefactorPlanPhases(t *testing.T) {
	data := `
## Phase 1: Foundation
- [ ] Task A
- [x] Task B

## Phase 2: Polish
- [ ] Task C
`

	phases := ParsePhases(data)
	if len(phases) != 2 {
		t.Fatalf("expected 2 phases, got %d", len(phases))
	}

	if phases[0].Name != "Phase 1: Foundation" {
		t.Fatalf("unexpected phase[0] name: %q", phases[0].Name)
	}
	if len(phases[0].Tasks) != 2 {
		t.Fatalf("expected 2 tasks in phase[0], got %d", len(phases[0].Tasks))
	}
	if phases[0].Completed {
		t.Fatalf("expected phase[0] not completed")
	}
	if phases[0].Tasks[1].Text != "Task B" || !phases[0].Tasks[1].Completed {
		t.Fatalf("expected Task B to be completed, got: %#v", phases[0].Tasks[1])
	}

	if phases[1].Name != "Phase 2: Polish" {
		t.Fatalf("unexpected phase[1] name: %q", phases[1].Name)
	}
	if len(phases[1].Tasks) != 1 {
		t.Fatalf("expected 1 task in phase[1], got %d", len(phases[1].Tasks))
	}
}

func TestParsePhases_ImplementationAtomicHeaders(t *testing.T) {
	data := `
### 1) Setup
- [ ] Task A
### 2) Next
- [x] Task B
`

	phases := ParsePhases(data)
	if len(phases) != 2 {
		t.Fatalf("expected 2 phases, got %d", len(phases))
	}
	if phases[0].Name != "Step 1: Setup" {
		t.Fatalf("unexpected phase[0] name: %q", phases[0].Name)
	}
	if phases[1].Name != "Step 2: Next" {
		t.Fatalf("unexpected phase[1] name: %q", phases[1].Name)
	}
	if len(phases[1].Tasks) != 1 || phases[1].Tasks[0].Text != "Task B" || !phases[1].Tasks[0].Completed {
		t.Fatalf("expected Task B completed under phase[1], got: %#v", phases[1].Tasks)
	}
	if !phases[1].Completed {
		t.Fatalf("expected phase[1] completed (all tasks done)")
	}
}

func TestParsePhases_FixPlanPriorityHeaders(t *testing.T) {
	data := `
## Critical Fixes
- [ ] Fix A

## Medium Priority
- [ ] Fix B
`

	phases := ParsePhases(data)
	if len(phases) != 2 {
		t.Fatalf("expected 2 phases, got %d", len(phases))
	}
	if phases[0].Name != "Critical Fixes" {
		t.Fatalf("unexpected phase[0] name: %q", phases[0].Name)
	}
	if phases[1].Name != "Medium Priority" {
		t.Fatalf("unexpected phase[1] name: %q", phases[1].Name)
	}
}

func TestParsePhases_DefaultPhaseWhenNoHeader(t *testing.T) {
	data := `
- [ ] Lone task
`

	phases := ParsePhases(data)
	if len(phases) != 1 {
		t.Fatalf("expected 1 phase, got %d", len(phases))
	}
	if phases[0].Name != "Tasks" {
		t.Fatalf("expected default phase name \"Tasks\", got %q", phases[0].Name)
	}
	if len(phases[0].Tasks) != 1 || phases[0].Tasks[0].Text != "Lone task" {
		t.Fatalf("unexpected tasks: %#v", phases[0].Tasks)
	}
}

func TestParseTasks_FlattensPhases(t *testing.T) {
	data := `
## Phase 1: A
- [ ] Task A
## Phase 2: B
- [ ] Task B
`

	tasks := ParseTasks(data)
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Text != "Task A" || tasks[1].Text != "Task B" {
		t.Fatalf("unexpected task order: %#v", tasks)
	}
}
