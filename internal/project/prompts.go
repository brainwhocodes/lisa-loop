package project

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Default template paths (can be overridden)
var (
	TemplateDir                  = ""
	ImplementationPlanPromptFile = "IMPLEMENTATION_PLAN_PROMPT.md"
	AgentsPromptFile             = "AGENTS_PROMPT.md"
	FixPlanPromptFile            = "FIX_PLAN_PROMPT.md"
	RefactorPlanPromptFile       = "REFACTOR_PLAN_PROMPT.md"
)

// getDefaultTemplateDir returns the default template directory based on the binary location
func getDefaultTemplateDir() string {
	// Try to find templates relative to the executable
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		// Check ../templates (installed location)
		templatesDir := filepath.Join(exeDir, "..", "templates")
		if _, err := os.Stat(templatesDir); err == nil {
			return templatesDir
		}
		// Check ./templates (development location)
		templatesDir = filepath.Join(exeDir, "templates")
		if _, err := os.Stat(templatesDir); err == nil {
			return templatesDir
		}
	}

	// Try relative to source file (for development)
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		srcDir := filepath.Dir(thisFile)
		templatesDir := filepath.Join(srcDir, "..", "..", "templates")
		if _, err := os.Stat(templatesDir); err == nil {
			return templatesDir
		}
	}

	return ""
}

// loadPromptTemplate loads a prompt template from file
func loadPromptTemplate(filename string) (string, error) {
	// Try custom template directory first
	if TemplateDir != "" {
		customPath := filepath.Join(TemplateDir, filename)
		if data, err := os.ReadFile(customPath); err == nil {
			return string(data), nil
		}
	}

	// Try ~/.lisa/templates
	if home, err := os.UserHomeDir(); err == nil {
		homePath := filepath.Join(home, ".lisa", "templates", filename)
		if data, err := os.ReadFile(homePath); err == nil {
			return string(data), nil
		}
	}

	// Try default template directory
	if defaultDir := getDefaultTemplateDir(); defaultDir != "" {
		defaultPath := filepath.Join(defaultDir, filename)
		if data, err := os.ReadFile(defaultPath); err == nil {
			return string(data), nil
		}
	}

	return "", fmt.Errorf("template %s not found (checked custom dir, ~/.lisa/templates, and default locations)", filename)
}

// GetImplementationPlanPrompt loads the implementation plan system prompt
func GetImplementationPlanPrompt() (string, error) {
	return loadPromptTemplate(ImplementationPlanPromptFile)
}

// GetAgentsPrompt loads the agents system prompt
func GetAgentsPrompt() (string, error) {
	return loadPromptTemplate(AgentsPromptFile)
}

// BuildImplementationPlanPrompt creates the full prompt for generating IMPLEMENTATION_PLAN.md
func BuildImplementationPlanPrompt(prdContent string) string {
	systemPrompt, err := GetImplementationPlanPrompt()
	if err != nil {
		// Fall back to minimal prompt if template can't be loaded
		systemPrompt = "Generate an IMPLEMENTATION_PLAN.md with phases and checklist tasks based on this PRD."
	}

	return systemPrompt + "\n\n---\n\nHere is the PRD:\n\n" + prdContent
}

// BuildAgentsPrompt creates the full prompt for generating AGENTS.md
func BuildAgentsPrompt(prdContent string) string {
	systemPrompt, err := GetAgentsPrompt()
	if err != nil {
		// Fall back to minimal prompt if template can't be loaded
		systemPrompt = "Generate an AGENTS.md file with project overview, tech stack, and development guidelines based on this PRD."
	}

	return systemPrompt + "\n\n---\n\nHere is the PRD:\n\n" + prdContent
}

// GetFixPlanPrompt loads the fix plan system prompt
func GetFixPlanPrompt() (string, error) {
	return loadPromptTemplate(FixPlanPromptFile)
}

// BuildFixPlanPrompt creates the full prompt for generating @fix_plan.md from specs
func BuildFixPlanPrompt(specsContent string) string {
	systemPrompt, err := GetFixPlanPrompt()
	if err != nil {
		// Fall back to minimal prompt if template can't be loaded
		systemPrompt = "Generate a @fix_plan.md with prioritized checklist of fixes and improvements based on these specs."
	}

	return systemPrompt + "\n\n---\n\nHere are the specifications:\n\n" + specsContent
}

// GetRefactorPlanPrompt loads the refactor plan system prompt
func GetRefactorPlanPrompt() (string, error) {
	return loadPromptTemplate(RefactorPlanPromptFile)
}

// BuildRefactorPlanPrompt creates the full prompt for generating REFACTOR_PLAN.md
func BuildRefactorPlanPrompt(refactorContent string) string {
	systemPrompt, err := GetRefactorPlanPrompt()
	if err != nil {
		// Fall back to minimal prompt if template can't be loaded
		systemPrompt = "Generate a REFACTOR_PLAN.md with phased refactoring tasks based on this refactoring document."
	}

	return systemPrompt + "\n\n---\n\nHere is the refactoring document:\n\n" + refactorContent
}
