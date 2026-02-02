# refactor(tui/plan): extract plan parsing/loading into internal/tui/plan

## Code changes
- [ ] Create `internal/tui/plan` package with pure functions:
  - [ ] `ParsePhases(data string) []plan.Phase`
  - [ ] `ParseTasks(data string) []plan.Task` (flat list, backwards compat)
  - [ ] `IsPhaseHeader(line string) bool` / `ExtractPhaseHeader(line string) string` (if kept)
- [ ] Keep old behavior:
  - [ ] supported headers stay supported (phase headers, atomic headers, fix-plan priority headers, verification/success criteria headers)
  - [ ] checkbox parsing stays `- [ ]` / `- [x]` (case-insensitive)
- [ ] Wire `internal/tui/program.go` to call `plan` package.

## Tests to add/update
- [ ] Add table tests for each supported plan format:
  - [ ] refactor plan phases
  - [ ] implementation plan phases + `### N) ...` headers
  - [ ] fix plan priority sections
- [ ] Add regression tests for tricky header detection cases (blank lines, leading/trailing spaces).
- [ ] Update `internal/tui/program_test.go` expectations if types move packages.

## Manual UX verification
- [ ] `go test ./...`
- [ ] With each mode (implement/refactor/fix), verify:
  - [ ] plan file is detected
  - [ ] task count matches file
  - [ ] header phase progress indicator renders

