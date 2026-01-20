package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSourceContent(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		wantPrompt   string
		wantFixPlan  string
		wantAgent    string
		wantWarnings int
	}{
		{
			name: "full content with all sections",
			content: `# Project Prompt
This is the prompt section.

## Tasks
- Task1
- Task2

### Agent Instructions
Build with: make
Run with: ./app`,
			wantPrompt:   "# Project Prompt\nThis is the prompt section.",
			wantFixPlan:  "## Tasks\n- Task1\n- Task2",
			wantAgent:    "### Agent Instructions\nBuild with: make\nRun with: ./app",
			wantWarnings: 0,
		},
		{
			name: "only prompt section",
			content: `# Development Instructions
Build a great application.`,
			wantPrompt:   "Build a great application.",
			wantFixPlan:  "# Ralph Development Task List\n\nUse this checklist to track progress. Mark complete items with [x].\n\n## Phase 1: Initial Setup\n- [ ] Create project structure\n- [ ] Set up development environment\n- [ ] Write initial tests\n- [ ] Implement core features\n\n## Phase 2: Development\n- [ ] Implement feature X\n- [ ] Implement feature Y\n- [ ] Write tests for X and Y\n\n## Phase 3: Polish\n- [ ] Add documentation\n- [ ] Performance optimization\n- [ ] Code review and cleanup\n\n## Phase 4: Deployment\n- [ ] Prepare for release\n- [ ] Deploy to production\n",
			wantAgent:    "# Build and Run Instructions\n\n## Building the Project\n\n```bash\ngo build ./...\n```\n\n## Running Tests\n\n```bash\ngo test ./...\n```\n\n## Running the Application\n\n```bash\n./<binary-name>\n```\n\n## Development Workflow\n\n1. Make changes to source code\n2. Run tests to verify\n3. Build the project\n4. Run to test manually\n5. Commit changes\n",
			wantWarnings: 2,
		},
		{
			name:         "empty content",
			content:      "",
			wantPrompt:   "# Development Instructions\n\nPlease specify development goals and rules.",
			wantFixPlan:  "# Ralph Development Task List\n\nUse this checklist to track progress. Mark complete items with [x].\n\n## Phase 1: Initial Setup\n- [ ] Create project structure\n- [ ] Set up development environment\n- [ ] Write initial tests\n- [ ] Implement core features\n\n## Phase 2: Development\n- [ ] Implement feature X\n- [ ] Implement feature Y\n- [ ] Write tests for X and Y\n\n## Phase 3: Polish\n- [ ] Add documentation\n- [ ] Performance optimization\n- [ ] Code review and cleanup\n\n## Phase 4: Deployment\n- [ ] Prepare for release\n- [ ] Deploy to production\n",
			wantAgent:    "# Build and Run Instructions\n\n## Building the Project\n\n```bash\ngo build ./...\n```\n\n## Running Tests\n\n```bash\ngo test ./...\n```\n\n## Running the Application\n\n```bash\n./<binary-name>\n```\n\n## Development Workflow\n\n1. Make changes to source code\n2. Run tests to verify\n3. Build the project\n4. Run to test manually\n5. Commit changes\n",
			wantWarnings: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt, fixPlan, agent, warnings, err := parseSourceContent(tt.content)
			if err != nil {
				t.Fatalf("parseSourceContent() error = %v", err)
			}

			if prompt != tt.wantPrompt {
				t.Errorf("parseSourceContent() prompt = %q, want %q", prompt, tt.wantPrompt)
			}

			if fixPlan != tt.wantFixPlan {
				t.Errorf("parseSourceContent() fixPlan = %q, want %q", fixPlan, tt.wantFixPlan)
			}

			if agent != tt.wantAgent {
				t.Errorf("parseSourceContent() agent = %q, want %q", agent, tt.wantAgent)
			}

			if len(warnings) != tt.wantWarnings {
				t.Errorf("parseSourceContent() warnings = %v, want %v", warnings, tt.wantWarnings)
			}
		})
	}
}

func TestImportPRD(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test PRD file
	prdContent := `# Test Project

## Prompt
Build a simple web app.

## Tasks
- Create homepage
- Add navigation

## Agent
npm install
npm start`

	prdPath := filepath.Join(tmpDir, "test-prd.md")
	if err := os.WriteFile(prdPath, []byte(prdContent), 0644); err != nil {
		t.Fatalf("Failed to create test PRD: %v", err)
	}

	// Import PRD
	opts := ImportOptions{
		SourcePath:    prdPath,
		ProjectName:   "test-import",
		OutputDir:     tmpDir,
		Verbose:       false,
		ConvertFormat: "",
	}

	result, err := ImportPRD(opts)
	if err != nil {
		t.Fatalf("ImportPRD() error = %v", err)
	}

	if !result.Success {
		t.Error("ImportPRD() returned unsuccessful result")
	}

	if result.ProjectName != "test-import" {
		t.Errorf("ImportPRD() ProjectName = %v, want test-import", result.ProjectName)
	}

	// Check that files were created
	expectedFiles := []string{
		"PROMPT.md",
		"@fix_plan.md",
		"@AGENT.md",
	}

	for _, filename := range expectedFiles {
		fullPath := filepath.Join(tmpDir, filename)
		if _, err := os.Stat(fullPath); err != nil {
			t.Errorf("ImportPRD() did not create expected file: %s", filename)
		}
	}
}

func TestImportPRDMissingSource(t *testing.T) {
	opts := ImportOptions{
		SourcePath:    "/nonexistent/file.md",
		ProjectName:   "test",
		OutputDir:     ".",
		Verbose:       false,
		ConvertFormat: "",
	}

	_, err := ImportPRD(opts)
	if err == nil {
		t.Error("ImportPRD() expected error for missing source, got nil")
	}
}

func TestImportPRDEmptySource(t *testing.T) {
	tmpDir := t.TempDir()

	opts := ImportOptions{
		SourcePath:    "",
		ProjectName:   "test",
		OutputDir:     tmpDir,
		Verbose:       false,
		ConvertFormat: "",
	}

	_, err := ImportPRD(opts)
	if err == nil {
		t.Error("ImportPRD() expected error for empty source path, got nil")
	}
}

func TestImportPRDInvalidOutputDir(t *testing.T) {
	tmpDir := t.TempDir()

	prdPath := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(prdPath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	opts := ImportOptions{
		SourcePath:    prdPath,
		ProjectName:   "test",
		OutputDir:     "/nonexistent/output",
		Verbose:       false,
		ConvertFormat: "",
	}

	_, err := ImportPRD(opts)
	if err == nil {
		t.Error("ImportPRD() expected error for invalid output dir, got nil")
	}
}

func TestGetConversionSummary(t *testing.T) {
	result := &ImportResult{
		ProjectName:   "test-project",
		ConvertedFrom: ".md",
		FilesCreated:  []string{"PROMPT.md", "@fix_plan.md", "@AGENT.md"},
		Success:       true,
		Warnings:      []string{"Warning 1", "Warning 2"},
	}

	summary := result.GetConversionSummary()
	if summary == "" {
		t.Error("GetConversionSummary() returned empty string")
	}

	if summary[0:8] != "Project:" {
		t.Errorf("GetConversionSummary() summary doesn't start with 'Project:', got: %s", summary[:8])
	}

	if summary[len(summary)-1:] != "\n" {
		t.Errorf("GetConversionSummary() summary doesn't end with newline, got: %q", summary[len(summary)-1:])
	}
}

func TestParseSourceContentInterleavedSections(t *testing.T) {
	// Test that interleaved/repeated section headings append content instead of resetting
	content := `# First Prompt Section
First prompt content.

## First Task Section
- Task A
- Task B

# Second Prompt Section
Second prompt content.

## Second Task Section
- Task C
- Task D

### Agent Section
Agent instructions here.`

	prompt, fixPlan, _, _, err := parseSourceContent(content)
	if err != nil {
		t.Fatalf("parseSourceContent() error = %v", err)
	}

	// Verify prompt contains both sections
	if !strings.Contains(prompt, "First prompt content") {
		t.Error("parseSourceContent() prompt missing 'First prompt content'")
	}
	if !strings.Contains(prompt, "Second prompt content") {
		t.Error("parseSourceContent() prompt missing 'Second prompt content' - repeated headings should append, not reset")
	}

	// Verify fix plan contains both task sections
	if !strings.Contains(fixPlan, "Task A") {
		t.Error("parseSourceContent() fixPlan missing 'Task A'")
	}
	if !strings.Contains(fixPlan, "Task C") {
		t.Error("parseSourceContent() fixPlan missing 'Task C' - repeated headings should append, not reset")
	}
}

func TestAutoDetectProjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "my-project.md",
			expected: "my-project",
		},
		{
			name:     "with path",
			input:    "PRD_my_app.md",
			expected: "my-app",
		},
		{
			name:     "all caps",
			input:    "MY_PROJECT.PRD",
			expected: "my-project",
		},
		{
			name:     "spaces and underscores",
			input:    "my_test_project_spec.MD",
			expected: "my-test-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			sourcePath := filepath.Join(tmpDir, tt.input)
			if err := os.WriteFile(sourcePath, []byte("# Test"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Create subdirectory as output
			outputDir := filepath.Join(tmpDir, "output")
			os.MkdirAll(outputDir, 0755)

			opts := ImportOptions{
				SourcePath:    sourcePath,
				ProjectName:   "", // Auto-detect
				OutputDir:     outputDir,
				Verbose:       false,
				ConvertFormat: "",
			}

			result, err := ImportPRD(opts)
			if err != nil {
				t.Fatalf("ImportPRD() error = %v", err)
			}

			if result.ProjectName != tt.expected {
				t.Errorf("ImportPRD() auto-detected ProjectName = %v, want %v", result.ProjectName, tt.expected)
			}
		})
	}
}
