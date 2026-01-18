package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid simple name",
			input:   "my-project",
			wantErr: false,
		},
		{
			name:    "valid camelCase",
			input:   "MyProject",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
		},
		{
			name:    "name with slash",
			input:   "my/project",
			wantErr: true,
		},
		{
			name:    "name with backslash",
			input:   "my\\project",
			wantErr: true,
		},
		{
			name:    "name with asterisk",
			input:   "my*project",
			wantErr: true,
		},
		{
			name:    "name with question mark",
			input:   "my?project",
			wantErr: true,
		},
		{
			name:    "name with colon",
			input:   "my:project",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestExtractProjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple filename",
			input:    "project.md",
			expected: "project",
		},
		{
			name:     "with path",
			input:    "/path/to/PRD_my_project.md",
			expected: "my-project",
		},
		{
			name:     "with PRD prefix",
			input:    "PRD_test_project.md",
			expected: "test-project",
		},
		{
			name:     "with spec prefix",
			input:    "spec_awesome_app.md",
			expected: "awesome-app",
		},
		{
			name:     "with requirements prefix",
			input:    "REQUIREMENTS_my_app.md",
			expected: "my-app",
		},
		{
			name:     "with suffix",
			input:    "app_prd.md",
			expected: "app",
		},
		{
			name:     "full path with extension",
			input:    "/home/user/specs/README.md",
			expected: "readme",
		},
		{
			name:     "underscores to hyphens",
			input:    "my_test_project.md",
			expected: "my-test-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("extractProjectName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSupportedFormat(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "markdown file",
			filename: "test.md",
			expected: true,
		},
		{
			name:     "text file",
			filename: "test.txt",
			expected: true,
		},
		{
			name:     "json file",
			filename: "test.json",
			expected: true,
		},
		{
			name:     "yaml file",
			filename: "test.yaml",
			expected: true,
		},
		{
			name:     "yml file",
			filename: "test.yml",
			expected: true,
		},
		{
			name:     "pdf file (not supported)",
			filename: "test.pdf",
			expected: false,
		},
		{
			name:     "docx file (not supported)",
			filename: "test.docx",
			expected: false,
		},
		{
			name:     "no extension",
			filename: "test",
			expected: false,
		},
		{
			name:     "uppercase extension",
			filename: "test.MD",
			expected: true,
		},
		{
			name:     "mixed case extension",
			filename: "test.Md",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSupportedFormat(tt.filename)
			if result != tt.expected {
				t.Errorf("IsSupportedFormat(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestSupportedFormats(t *testing.T) {
	formats := SupportedFormats()

	if len(formats) == 0 {
		t.Error("SupportedFormats() returned empty list")
	}

	expectedFormats := []string{".md", ".txt", ".json", ".yaml", ".yml"}
	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SupportedFormats() missing expected format: %s", expected)
		}
	}
}

func TestSetup(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test-" + filepath.Base(tmpDir)

	opts := SetupOptions{
		ProjectName: projectName,
		TemplateDir: "",
		WithGit:     false,
		Verbose:     false,
	}

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	result, err := Setup(opts)
	if err != nil {
		t.Fatalf("Setup() error = %v", err)
	}

	if !result.Success {
		t.Error("Setup() returned unsuccessful result")
	}

	expectedPath := filepath.Join(tmpDir, projectName)
	if result.ProjectPath != expectedPath {
		// On some systems paths might differ, check basename instead
		if filepath.Base(result.ProjectPath) != filepath.Base(expectedPath) {
			t.Errorf("Setup() ProjectPath = %v, want %v", result.ProjectPath, expectedPath)
		}
	}

	// Check that expected files were created
	expectedFiles := []string{
		"PROMPT.md",
		"@fix_plan.md",
		"@AGENT.md",
		"README.md",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(result.ProjectPath, file)
		if _, err := os.Stat(fullPath); err != nil {
			t.Errorf("Setup() did not create expected file: %s", file)
		}
	}

	// Check that expected directories were created
	expectedDirs := []string{
		"src",
		"examples",
		"specs",
		"docs",
		"docs/generated",
		"logs",
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(result.ProjectPath, dir)
		if info, err := os.Stat(fullPath); err != nil {
			t.Errorf("Setup() did not create expected directory: %s", dir)
		} else if !info.IsDir() {
			t.Errorf("Setup() created %s but it's not a directory", dir)
		}
	}
}

func TestSetupInvalidName(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := "test/invalid"

	opts := SetupOptions{
		ProjectName: projectName,
		TemplateDir: "",
		WithGit:     false,
		Verbose:     false,
	}

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	_, err := Setup(opts)
	if err == nil {
		t.Error("Setup() expected error for invalid project name, got nil")
	}
}

func TestSetupEmptyName(t *testing.T) {
	opts := SetupOptions{
		ProjectName: "",
		TemplateDir: "",
		WithGit:     false,
		Verbose:     false,
	}

	_, err := Setup(opts)
	if err == nil {
		t.Error("Setup() expected error for empty project name, got nil")
	}
}
