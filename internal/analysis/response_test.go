package analysis

import (
	"testing"
)

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected OutputFormat
	}{
		{
			name:     "json object",
			input:    `{"type": "message", "text": "hello"}`,
			expected: FormatJSON,
		},
		{
			name:     "json array",
			input:    `[{"type": "message"}]`,
			expected: FormatJSON,
		},
		{
			name:     "json with leading whitespace",
			input:    `  {"type": "message"}`,
			expected: FormatJSON,
		},
		{
			name:     "json with leading newlines",
			input:    "\n\n{\"type\": \"message\"}",
			expected: FormatJSON,
		},
		{
			name:     "text with RALPH_STATUS",
			input:    "Some output\n---RALPH_STATUS---\nSTATUS: COMPLETE\n---END_RALPH_STATUS---",
			expected: FormatText,
		},
		{
			name:     "plain text",
			input:    "Hello, this is plain text output",
			expected: FormatText,
		},
		{
			name:     "empty string",
			input:    "",
			expected: FormatText,
		},
		{
			name:     "whitespace only",
			input:    "   \n\n   ",
			expected: FormatText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectFormat(tt.input)
			if result != tt.expected {
				t.Errorf("DetectFormat() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestParseRALPHStatus(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStatus string
		expectedTasks  int
		expectedFiles  int
		expectedExit   bool
	}{
		{
			name: "complete status block",
			input: `Some output text
---RALPH_STATUS---
STATUS: COMPLETE
TASKS_COMPLETED_THIS_LOOP: 3
FILES_MODIFIED: 5
TESTS_STATUS: PASSING
WORK_TYPE: implementation
EXIT_SIGNAL: true
RECOMMENDATION: Review changes
---END_RALPH_STATUS---
More text`,
			expectedStatus: "COMPLETE",
			expectedTasks:  3,
			expectedFiles:  5,
			expectedExit:   true,
		},
		{
			name: "blocked status",
			input: `---RALPH_STATUS---
STATUS: BLOCKED
TASKS_COMPLETED_THIS_LOOP: 0
FILES_MODIFIED: 0
TESTS_STATUS: FAILING
EXIT_SIGNAL: false
---END_RALPH_STATUS---`,
			expectedStatus: "BLOCKED",
			expectedTasks:  0,
			expectedFiles:  0,
			expectedExit:   false,
		},
		{
			name: "in_progress status",
			input: `---RALPH_STATUS---
STATUS: IN_PROGRESS
TASKS_COMPLETED_THIS_LOOP: 2
FILES_MODIFIED: 3
---END_RALPH_STATUS---`,
			expectedStatus: "IN_PROGRESS",
			expectedTasks:  2,
			expectedFiles:  3,
			expectedExit:   false,
		},
		{
			name:           "no status block",
			input:          "Just plain text without status block",
			expectedStatus: "UNKNOWN",
			expectedTasks:  0,
			expectedFiles:  0,
			expectedExit:   false,
		},
		{
			name:           "incomplete status block",
			input:          "---RALPH_STATUS---\nSTATUS: COMPLETE\nNo end marker",
			expectedStatus: "UNKNOWN",
			expectedTasks:  0,
			expectedFiles:  0,
			expectedExit:   false,
		},
		{
			name: "exit signal case insensitive",
			input: `---RALPH_STATUS---
STATUS: COMPLETE
EXIT_SIGNAL: TRUE
---END_RALPH_STATUS---`,
			expectedStatus: "COMPLETE",
			expectedExit:   true,
		},
		{
			name: "extra whitespace in values",
			input: `---RALPH_STATUS---
STATUS:    COMPLETE
TASKS_COMPLETED_THIS_LOOP:   5
---END_RALPH_STATUS---`,
			expectedStatus: "COMPLETE",
			expectedTasks:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseRALPHStatus(tt.input)

			if result.Status != tt.expectedStatus {
				t.Errorf("Status = %v, expected %v", result.Status, tt.expectedStatus)
			}
			if result.TasksCompleted != tt.expectedTasks {
				t.Errorf("TasksCompleted = %v, expected %v", result.TasksCompleted, tt.expectedTasks)
			}
			if result.FilesModified != tt.expectedFiles {
				t.Errorf("FilesModified = %v, expected %v", result.FilesModified, tt.expectedFiles)
			}
			if result.ExitSignal != tt.expectedExit {
				t.Errorf("ExitSignal = %v, expected %v", result.ExitSignal, tt.expectedExit)
			}
		})
	}
}

func TestDetectCompletionKeywords(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		minExpected int
	}{
		{
			name:        "multiple completion keywords",
			input:       "All done! Task completed successfully. We're all set and ready to review.",
			minExpected: 4, // done, completed, all set, ready (to review)
		},
		{
			name:        "single completion keyword",
			input:       "The task is done.",
			minExpected: 1,
		},
		{
			name:        "no completion keywords",
			input:       "Working on the task...",
			minExpected: 0,
		},
		{
			name:        "case insensitive",
			input:       "DONE and COMPLETE and FINISHED",
			minExpected: 3,
		},
		{
			name:        "partial matches",
			input:       "nothing to do here, no more work needed",
			minExpected: 2, // nothing to do, no more work
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectCompletionKeywords(tt.input)
			if result < tt.minExpected {
				t.Errorf("DetectCompletionKeywords() = %v, expected at least %v", result, tt.minExpected)
			}
		})
	}
}

func TestExtractErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name: "multiple errors",
			input: `Error: Failed to compile
Some output
ERROR: Network timeout
More output`,
			expectedCount: 2,
		},
		{
			name: "no errors",
			input: `Success! All tasks completed.
Everything is working correctly.`,
			expectedCount: 0,
		},
		{
			name: "json field false positive filtered",
			input: `{"is_error": false}
Error: Actual error`,
			expectedCount: 1,
		},
		{
			name:          "lowercase error prefix",
			input:         "error: something went wrong",
			expectedCount: 1,
		},
		{
			name:          "fatal error",
			input:         "Fatal: Critical failure",
			expectedCount: 1,
		},
		{
			name:          "exception",
			input:         "Exception occurred during processing",
			expectedCount: 1,
		},
		{
			name:          "empty input",
			input:         "",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractErrors(tt.input)
			if len(result) != tt.expectedCount {
				t.Errorf("ExtractErrors() returned %v errors, expected %v. Got: %v", len(result), tt.expectedCount, result)
			}
		})
	}
}

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		exitSignals    []string
		expectedFormat OutputFormat
		expectedExit   bool
		hasErrors      bool
	}{
		{
			name: "complete status with exit signal",
			input: `---RALPH_STATUS---
STATUS: COMPLETE
EXIT_SIGNAL: true
TESTS_STATUS: PASSING
---END_RALPH_STATUS---`,
			exitSignals:    []string{},
			expectedFormat: FormatText,
			expectedExit:   true,
			hasErrors:      false,
		},
		{
			name: "blocked status",
			input: `---RALPH_STATUS---
STATUS: BLOCKED
TESTS_STATUS: FAILING
---END_RALPH_STATUS---`,
			exitSignals:    []string{},
			expectedFormat: FormatText,
			expectedExit:   false,
			hasErrors:      true,
		},
		{
			name:           "json output",
			input:          `{"type": "message", "text": "Task complete"}`,
			exitSignals:    []string{},
			expectedFormat: FormatJSON,
			expectedExit:   false,
			hasErrors:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Analyze(tt.input, tt.exitSignals)
			if err != nil {
				t.Fatalf("Analyze() returned error: %v", err)
			}

			if result.Format != tt.expectedFormat {
				t.Errorf("Format = %v, expected %v", result.Format, tt.expectedFormat)
			}
			if result.ExitSignal != tt.expectedExit {
				t.Errorf("ExitSignal = %v, expected %v", result.ExitSignal, tt.expectedExit)
			}
			if result.HasErrors != tt.hasErrors {
				t.Errorf("HasErrors = %v, expected %v", result.HasErrors, tt.hasErrors)
			}
		})
	}
}

func TestCalculateConfidence(t *testing.T) {
	tests := []struct {
		name            string
		status          *RALPHStatus
		completionCount int
		output          string
		minConfidence   float64
		maxConfidence   float64
	}{
		{
			name: "high confidence complete",
			status: &RALPHStatus{
				Status:     "COMPLETE",
				ExitSignal: true,
			},
			completionCount: 5,
			output:          "done complete finished",
			minConfidence:   0.9,
			maxConfidence:   1.0,
		},
		{
			name: "low confidence blocked",
			status: &RALPHStatus{
				Status:     "BLOCKED",
				ExitSignal: false,
			},
			completionCount: 0,
			output:          "",
			minConfidence:   0.0,
			maxConfidence:   0.3,
		},
		{
			name: "medium confidence in progress",
			status: &RALPHStatus{
				Status:     "IN_PROGRESS",
				ExitSignal: false,
			},
			completionCount: 1,
			output:          "working",
			minConfidence:   0.4,
			maxConfidence:   0.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateConfidence(tt.status, tt.completionCount, tt.output)
			if result < tt.minConfidence || result > tt.maxConfidence {
				t.Errorf("calculateConfidence() = %v, expected between %v and %v", result, tt.minConfidence, tt.maxConfidence)
			}
		})
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		shouldError bool
	}{
		{
			name:        "simple number",
			input:       "42",
			expected:    42,
			shouldError: false,
		},
		{
			name:        "zero",
			input:       "0",
			expected:    0,
			shouldError: false,
		},
		{
			name:        "negative number",
			input:       "-5",
			expected:    -5,
			shouldError: false,
		},
		{
			name:        "number with trailing text",
			input:       "123 files",
			expected:    123,
			shouldError: false,
		},
		{
			name:        "non-numeric",
			input:       "abc",
			expected:    0,
			shouldError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    0,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseNumber(tt.input)
			if tt.shouldError && err == nil {
				t.Errorf("parseNumber(%q) expected error but got none", tt.input)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("parseNumber(%q) unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("parseNumber(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
