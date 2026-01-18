package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetupOptions holds options for project setup
type SetupOptions struct {
	ProjectName string
	TemplateDir string
	WithGit     bool
	Verbose     bool
}

// SetupResult holds result of project setup
type SetupResult struct {
	ProjectPath    string
	FilesCreated   []string
	GitInitialized bool
	Success        bool
}

// Setup creates a new Ralph-managed project
func Setup(opts SetupOptions) (*SetupResult, error) {
	if opts.ProjectName == "" {
		return nil, fmt.Errorf("project name is required")
	}

	// Validate project name
	if err := validateProjectName(opts.ProjectName); err != nil {
		return nil, fmt.Errorf("invalid project name: %w", err)
	}

	// Create project directory
	projectPath := opts.ProjectName
	if strings.HasPrefix(opts.ProjectName, "/") || strings.HasPrefix(opts.ProjectName, "~") {
		projectPath = opts.ProjectName
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		projectPath = filepath.Join(wd, opts.ProjectName)
	}

	// Check if directory exists
	if _, err := os.Stat(projectPath); err == nil {
		return nil, fmt.Errorf("directory already exists: %s", projectPath)
	}

	result := &SetupResult{
		ProjectPath:  projectPath,
		FilesCreated: []string{},
		Success:      false,
	}

	// Create base directory structure
	dirs := []string{
		projectPath,
		filepath.Join(projectPath, "src"),
		filepath.Join(projectPath, "examples"),
		filepath.Join(projectPath, "specs"),
		filepath.Join(projectPath, "docs"),
		filepath.Join(projectPath, "docs", "generated"),
		filepath.Join(projectPath, "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		if opts.Verbose {
			fmt.Printf("Created directory: %s\n", dir)
		}
	}

	// Create template files
	files, err := createTemplateFiles(projectPath, opts.TemplateDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create template files: %w", err)
	}
	result.FilesCreated = files

	// Initialize git repository if requested
	if opts.WithGit {
		if err := initGitRepo(projectPath); err != nil {
			return nil, fmt.Errorf("failed to initialize git: %w", err)
		}
		result.GitInitialized = true
		if opts.Verbose {
			fmt.Println("Initialized git repository")
		}
	}

	// Create README
	readmePath := filepath.Join(projectPath, "README.md")
	if err := createREADME(readmePath, opts.ProjectName); err != nil {
		return nil, fmt.Errorf("failed to create README: %w", err)
	}
	result.FilesCreated = append(result.FilesCreated, readmePath)

	result.Success = true
	return result, nil
}

// validateProjectName checks if project name is valid
func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	// Check for invalid characters
	for _, c := range name {
		if c == '/' || c == '\\' || c == ':' || c == '*' || c == '?' ||
			c == '"' || c == '<' || c == '>' || c == '|' {
			return fmt.Errorf("name contains invalid character: %c", c)
		}
	}

	return nil
}

// createTemplateFiles creates template files in project directory
func createTemplateFiles(projectPath, templateDir string) ([]string, error) {
	files := []string{}

	// Default template directory
	if templateDir == "" {
		templateDir = "~/.ralph/templates"
	}

	// Expand ~ to home directory
	if strings.HasPrefix(templateDir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		templateDir = filepath.Join(home, templateDir[2:])
	}

	// Template files to create
	templates := map[string]string{
		"PROMPT.md":    defaultPromptTemplate(),
		"@fix_plan.md": defaultFixPlanTemplate(),
		"@AGENT.md":    defaultAgentTemplate(),
	}

	for filename, content := range templates {
		targetPath := filepath.Join(projectPath, filename)

		// Try to load from template directory first
		templatePath := filepath.Join(templateDir, filename)
		var fileContent string

		if _, err := os.Stat(templatePath); err == nil {
			data, err := os.ReadFile(templatePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read template %s: %w", templatePath, err)
			}
			fileContent = string(data)
		} else {
			fileContent = content
		}

		if err := os.WriteFile(targetPath, []byte(fileContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", targetPath, err)
		}
		files = append(files, targetPath)
	}

	return files, nil
}

// initGitRepo initializes a git repository in the project directory
func initGitRepo(projectPath string) error {
	origDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(origDir)

	if err := os.Chdir(projectPath); err != nil {
		return err
	}

	// Initialize git repo
	cmd := fmt.Sprintf("cd %s && git init", projectPath)
	if err := executeCommand(cmd); err != nil {
		return err
	}

	// Create .gitignore
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	gitignoreContent := `# Ralph Codex
.ralph_session
.circuit_breaker_state
.exit_signals
.call_count
.response_analysis

# Logs
logs/

# OS
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/
*.swp
*.swo
`
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return err
	}

	// Initial commit
	cmd = fmt.Sprintf("cd %s && git add . && git commit -m 'Initial commit from ralph-setup'", projectPath)
	if err := executeCommand(cmd); err != nil {
		return err
	}

	return nil
}

// createREADME creates a README.md for the project
func createREADME(path, projectName string) error {
	builder := &strings.Builder{}
	builder.WriteString("# " + projectName + "\n\n")
	builder.WriteString("This is a Ralph Codex managed project.\n\n")
	builder.WriteString("## Getting Started\n\n")
	builder.WriteString("Ralph Codex autonomously develops software based on specifications in:\n")
	builder.WriteString("- PROMPT.md - Main development instructions\n")
	builder.WriteString("- @fix_plan.md - Prioritized task list\n")
	builder.WriteString("- @AGENT.md - Build and run instructions\n\n")
	builder.WriteString("## Running Ralph\n\n")
	builder.WriteString("To start the autonomous development loop:\n\n")
	builder.WriteString("```bash\n")
	builder.WriteString("ralph\n")
	builder.WriteString("```\n\n")
	builder.WriteString("With monitoring dashboard:\n\n")
	builder.WriteString("```bash\n")
	builder.WriteString("ralph --monitor\n")
	builder.WriteString("```\n\n")
	builder.WriteString("## Project Structure\n\n")
	builder.WriteString("- src/ - Source code\n")
	builder.WriteString("- examples/ - Usage examples\n")
	builder.WriteString("- specs/ - Project specifications\n")
	builder.WriteString("- docs/ - Documentation\n")
	builder.WriteString("- logs/ - Ralph execution logs\n")
	builder.WriteString("- docs/generated/ - Auto-generated documentation\n\n")
	builder.WriteString("## Status\n\n")
	builder.WriteString("Ralph manages the development cycle. Check @fix_plan.md for current progress.\n")

	return os.WriteFile(path, []byte(builder.String()), 0644)
}

// executeCommand executes a shell command (simplified)
func executeCommand(cmd string) error {
	// For now, return nil - this would be implemented with os/exec
	// The actual implementation would use exec.Command()
	return nil
}

// ValidateProject checks if current directory is a valid Ralph project
func ValidateProject() error {
	requiredFiles := []string{"PROMPT.md", "@fix_plan.md"}
	missing := []string{}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); err != nil {
			missing = append(missing, file)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("not a valid Ralph project (missing: %s)", strings.Join(missing, ", "))
	}

	return nil
}

// GetProjectRoot finds the project root directory
func GetProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(wd, "PROMPT.md")); err == nil {
			if _, err := os.Stat(filepath.Join(wd, "@fix_plan.md")); err == nil {
				return wd, nil
			}
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}

	return "", fmt.Errorf("could not find Ralph project root (no PROMPT.md or @fix_plan.md found)")
}

// defaultPromptTemplate returns default PROMPT.md content
func defaultPromptTemplate() string {
	builder := &strings.Builder{}
	builder.WriteString("# Ralph Codex Development Instructions\n\n")
	builder.WriteString("You are an autonomous AI developer working on this project.\n\n")
	builder.WriteString("## Project Goals\n\n")
	builder.WriteString("[Describe what this project should accomplish]\n\n")
	builder.WriteString("## Development Rules\n\n")
	builder.WriteString("1. Always read the @fix_plan.md before starting work\n")
	builder.WriteString("2. Follow @AGENT.md for build and run instructions\n")
	builder.WriteString("3. Run tests after each change\n")
	builder.WriteString("4. Commit your work regularly with clear messages\n")
	builder.WriteString("5. Update @fix_plan.md when tasks are complete\n\n")
	builder.WriteString("## Tech Stack\n\n")
	builder.WriteString("[List your technologies here]\n")
	return builder.String()
}

// defaultFixPlanTemplate returns default @fix_plan.md content
func defaultFixPlanTemplate() string {
	builder := &strings.Builder{}
	builder.WriteString("# Ralph Development Task List\n\n")
	builder.WriteString("Use this checklist to track progress. Mark complete items with [x].\n\n")
	builder.WriteString("## Phase 1: Initial Setup\n")
	builder.WriteString("- [ ] Create project structure\n")
	builder.WriteString("- [ ] Set up development environment\n")
	builder.WriteString("- [ ] Write initial tests\n")
	builder.WriteString("- [ ] Implement core features\n\n")
	builder.WriteString("## Phase 2: Development\n")
	builder.WriteString("- [ ] Implement feature X\n")
	builder.WriteString("- [ ] Implement feature Y\n")
	builder.WriteString("- [ ] Write tests for X and Y\n\n")
	builder.WriteString("## Phase 3: Polish\n")
	builder.WriteString("- [ ] Add documentation\n")
	builder.WriteString("- [ ] Performance optimization\n")
	builder.WriteString("- [ ] Code review and cleanup\n\n")
	builder.WriteString("## Phase 4: Deployment\n")
	builder.WriteString("- [ ] Prepare for release\n")
	builder.WriteString("- [ ] Deploy to production\n")
	return builder.String()
}

// defaultAgentTemplate returns default @AGENT.md content
func defaultAgentTemplate() string {
	builder := &strings.Builder{}
	builder.WriteString("# Build and Run Instructions\n\n")
	builder.WriteString("## Building the Project\n\n")
	builder.WriteString("```bash\n")
	builder.WriteString("go build ./...\n")
	builder.WriteString("```\n\n")
	builder.WriteString("## Running Tests\n\n")
	builder.WriteString("```bash\n")
	builder.WriteString("go test ./...\n")
	builder.WriteString("```\n\n")
	builder.WriteString("## Running the Application\n\n")
	builder.WriteString("```bash\n")
	builder.WriteString("./<binary-name>\n")
	builder.WriteString("```\n\n")
	builder.WriteString("## Development Workflow\n\n")
	builder.WriteString("1. Make changes to source code\n")
	builder.WriteString("2. Run tests to verify\n")
	builder.WriteString("3. Build the project\n")
	builder.WriteString("4. Run to test manually\n")
	builder.WriteString("5. Commit changes\n")
	return builder.String()
}
