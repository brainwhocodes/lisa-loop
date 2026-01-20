package codex

import (
	"testing"
)

func TestParseEventNil(t *testing.T) {
	result := ParseEvent(nil)
	if result != nil {
		t.Error("ParseEvent(nil) should return nil")
	}
}

func TestParseEventMessage(t *testing.T) {
	tests := []struct {
		name         string
		event        Event
		expectedType string
		expectedText string
	}{
		{
			name: "message with content string",
			event: Event{
				"type":    "message",
				"content": "Hello world",
			},
			expectedType: "message",
			expectedText: "Hello world",
		},
		{
			name: "message with content array",
			event: Event{
				"type": "message",
				"content": []interface{}{
					map[string]interface{}{"text": "First part"},
					map[string]interface{}{"text": " Second part"},
				},
			},
			expectedType: "message",
			expectedText: "First part Second part",
		},
		{
			name: "unknown type with text field uses default handler",
			event: Event{
				"type": "some_other_type",
				"text": "Direct text",
			},
			expectedType: "message",
			expectedText: "Direct text",
		},
		{
			name: "message with message field fallback",
			event: Event{
				"type":    "unknown",
				"message": "Fallback message",
			},
			expectedType: "message",
			expectedText: "Fallback message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseEvent(tt.event)
			if result == nil {
				t.Fatal("ParseEvent() returned nil")
			}
			if result.Type != tt.expectedType {
				t.Errorf("Type = %v, expected %v", result.Type, tt.expectedType)
			}
			if result.Text != tt.expectedText {
				t.Errorf("Text = %v, expected %v", result.Text, tt.expectedText)
			}
		})
	}
}

func TestParseEventContentBlockDelta(t *testing.T) {
	event := Event{
		"type": "content_block_delta",
		"delta": map[string]interface{}{
			"text": "Streaming text chunk",
		},
	}

	result := ParseEvent(event)
	if result == nil {
		t.Fatal("ParseEvent() returned nil")
	}
	if result.Type != "delta" {
		t.Errorf("Type = %v, expected delta", result.Type)
	}
	if result.Text != "Streaming text chunk" {
		t.Errorf("Text = %v, expected 'Streaming text chunk'", result.Text)
	}
}

func TestParseEventToolUse(t *testing.T) {
	tests := []struct {
		name           string
		event          Event
		expectedType   string
		expectedTool   string
		expectedTarget string
		expectedStatus string
	}{
		{
			name: "tool_use with file_path",
			event: Event{
				"type": "tool_use",
				"name": "Read",
				"arguments": map[string]interface{}{
					"file_path": "/path/to/file.go",
				},
			},
			expectedType:   "tool_call",
			expectedTool:   "Read",
			expectedTarget: "/path/to/file.go",
			expectedStatus: "started",
		},
		{
			name: "tool_use with command",
			event: Event{
				"type": "tool_use",
				"name": "Bash",
				"arguments": map[string]interface{}{
					"command": "go build ./...",
				},
			},
			expectedType:   "tool_call",
			expectedTool:   "Bash",
			expectedTarget: "go build ./...",
			expectedStatus: "started",
		},
		{
			name: "tool_use with input field",
			event: Event{
				"type": "tool_use",
				"name": "Write",
				"input": map[string]interface{}{
					"path": "/path/to/output.txt",
				},
			},
			expectedType:   "tool_call",
			expectedTool:   "Write",
			expectedTarget: "/path/to/output.txt",
			expectedStatus: "started",
		},
		{
			name: "tool_use with long command truncation",
			event: Event{
				"type": "tool_use",
				"name": "Bash",
				"arguments": map[string]interface{}{
					"command": "echo 'This is a very long command that exceeds fifty characters and should be truncated'",
				},
			},
			expectedType:   "tool_call",
			expectedTool:   "Bash",
			expectedTarget: "echo 'This is a very long command that exceeds fif...", // 50 chars + ...
			expectedStatus: "started",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseEvent(tt.event)
			if result == nil {
				t.Fatal("ParseEvent() returned nil")
			}
			if result.Type != tt.expectedType {
				t.Errorf("Type = %v, expected %v", result.Type, tt.expectedType)
			}
			if result.ToolName != tt.expectedTool {
				t.Errorf("ToolName = %v, expected %v", result.ToolName, tt.expectedTool)
			}
			if result.ToolTarget != tt.expectedTarget {
				t.Errorf("ToolTarget = %v, expected %v", result.ToolTarget, tt.expectedTarget)
			}
			if result.ToolStatus != tt.expectedStatus {
				t.Errorf("ToolStatus = %v, expected %v", result.ToolStatus, tt.expectedStatus)
			}
		})
	}
}

func TestParseEventToolResult(t *testing.T) {
	tests := []struct {
		name           string
		event          Event
		expectedType   string
		expectedTool   string
		expectedStatus string
	}{
		{
			name: "tool_result with name",
			event: Event{
				"type": "tool_result",
				"name": "Read",
			},
			expectedType:   "tool_result",
			expectedTool:   "Read",
			expectedStatus: "completed",
		},
		{
			name: "tool_result with nested tool_use",
			event: Event{
				"type": "tool_result",
				"tool_use": map[string]interface{}{
					"name": "Write",
				},
			},
			expectedType:   "tool_result",
			expectedTool:   "Write",
			expectedStatus: "completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseEvent(tt.event)
			if result == nil {
				t.Fatal("ParseEvent() returned nil")
			}
			if result.Type != tt.expectedType {
				t.Errorf("Type = %v, expected %v", result.Type, tt.expectedType)
			}
			if result.ToolName != tt.expectedTool {
				t.Errorf("ToolName = %v, expected %v", result.ToolName, tt.expectedTool)
			}
			if result.ToolStatus != tt.expectedStatus {
				t.Errorf("ToolStatus = %v, expected %v", result.ToolStatus, tt.expectedStatus)
			}
		})
	}
}

func TestParseEventItemCompleted(t *testing.T) {
	tests := []struct {
		name         string
		event        Event
		expectedType string
		expectedText string
		expectedTool string
	}{
		{
			name: "item.completed reasoning",
			event: Event{
				"type": "item.completed",
				"item": map[string]interface{}{
					"type": "reasoning",
					"text": "Thinking about the problem...",
				},
			},
			expectedType: "reasoning",
			expectedText: "Thinking about the problem...",
		},
		{
			name: "item.completed agent_message",
			event: Event{
				"type": "item.completed",
				"item": map[string]interface{}{
					"type": "agent_message",
					"text": "Here is my response",
				},
			},
			expectedType: "message",
			expectedText: "Here is my response",
		},
		{
			name: "item.completed message",
			event: Event{
				"type": "item.completed",
				"item": map[string]interface{}{
					"type": "message",
					"text": "Another message",
				},
			},
			expectedType: "message",
			expectedText: "Another message",
		},
		{
			name: "item.completed tool_call",
			event: Event{
				"type": "item.completed",
				"item": map[string]interface{}{
					"type": "tool_call",
					"name": "Read",
					"arguments": map[string]interface{}{
						"file_path": "/test.go",
					},
				},
			},
			expectedType: "tool_call",
			expectedTool: "Read",
		},
		{
			name: "item.completed function_call",
			event: Event{
				"type": "item.completed",
				"item": map[string]interface{}{
					"type": "function_call",
					"name": "Write",
					"arguments": map[string]interface{}{
						"file": "/output.txt",
					},
				},
			},
			expectedType: "tool_call",
			expectedTool: "Write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseEvent(tt.event)
			if result == nil {
				t.Fatal("ParseEvent() returned nil")
			}
			if result.Type != tt.expectedType {
				t.Errorf("Type = %v, expected %v", result.Type, tt.expectedType)
			}
			if tt.expectedText != "" && result.Text != tt.expectedText {
				t.Errorf("Text = %v, expected %v", result.Text, tt.expectedText)
			}
			if tt.expectedTool != "" && result.ToolName != tt.expectedTool {
				t.Errorf("ToolName = %v, expected %v", result.ToolName, tt.expectedTool)
			}
		})
	}
}

func TestParseEventAssistant(t *testing.T) {
	event := Event{
		"type": "assistant",
		"content": []interface{}{
			map[string]interface{}{"text": "Assistant response"},
		},
	}

	result := ParseEvent(event)
	if result == nil {
		t.Fatal("ParseEvent() returned nil")
	}
	if result.Type != "message" {
		t.Errorf("Type = %v, expected message", result.Type)
	}
	if result.Text != "Assistant response" {
		t.Errorf("Text = %v, expected 'Assistant response'", result.Text)
	}
}

func TestParseEventLifecycle(t *testing.T) {
	event := Event{
		"type": "thread.started",
	}

	result := ParseEvent(event)
	if result == nil {
		t.Fatal("ParseEvent() returned nil")
	}
	if result.Type != "lifecycle" {
		t.Errorf("Type = %v, expected lifecycle", result.Type)
	}
	if result.RawType != "thread.started" {
		t.Errorf("RawType = %v, expected thread.started", result.RawType)
	}
}

func TestParseEventUnknownWithText(t *testing.T) {
	event := Event{
		"type": "unknown_event",
		"text": "Some text content",
	}

	result := ParseEvent(event)
	if result == nil {
		t.Fatal("ParseEvent() returned nil")
	}
	if result.Type != "message" {
		t.Errorf("Type = %v, expected message", result.Type)
	}
	if result.Text != "Some text content" {
		t.Errorf("Text = %v, expected 'Some text content'", result.Text)
	}
}

func TestExtractToolTarget(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected string
	}{
		{
			name: "file_path in arguments",
			data: map[string]interface{}{
				"arguments": map[string]interface{}{
					"file_path": "/path/to/file",
				},
			},
			expected: "/path/to/file",
		},
		{
			name: "path in arguments",
			data: map[string]interface{}{
				"arguments": map[string]interface{}{
					"path": "/another/path",
				},
			},
			expected: "/another/path",
		},
		{
			name: "filename in arguments",
			data: map[string]interface{}{
				"arguments": map[string]interface{}{
					"filename": "test.txt",
				},
			},
			expected: "test.txt",
		},
		{
			name: "file in arguments",
			data: map[string]interface{}{
				"arguments": map[string]interface{}{
					"file": "output.go",
				},
			},
			expected: "output.go",
		},
		{
			name: "command in arguments",
			data: map[string]interface{}{
				"arguments": map[string]interface{}{
					"command": "ls -la",
				},
			},
			expected: "ls -la",
		},
		{
			name: "path in input field",
			data: map[string]interface{}{
				"input": map[string]interface{}{
					"path": "/input/path",
				},
			},
			expected: "/input/path",
		},
		{
			name: "path in parameters field",
			data: map[string]interface{}{
				"parameters": map[string]interface{}{
					"file_path": "/param/path",
				},
			},
			expected: "/param/path",
		},
		{
			name: "no matching field",
			data: map[string]interface{}{
				"arguments": map[string]interface{}{
					"unrelated": "value",
				},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractToolTarget(tt.data)
			if result != tt.expected {
				t.Errorf("extractToolTarget() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestExtractTextFromContentArray(t *testing.T) {
	tests := []struct {
		name     string
		content  []interface{}
		expected string
	}{
		{
			name: "single text item",
			content: []interface{}{
				map[string]interface{}{"text": "Hello"},
			},
			expected: "Hello",
		},
		{
			name: "multiple text items",
			content: []interface{}{
				map[string]interface{}{"text": "Hello"},
				map[string]interface{}{"text": " World"},
			},
			expected: "Hello World",
		},
		{
			name: "mixed content",
			content: []interface{}{
				map[string]interface{}{"text": "Text"},
				map[string]interface{}{"image": "base64..."},
				map[string]interface{}{"text": " More"},
			},
			expected: "Text More",
		},
		{
			name:     "empty array",
			content:  []interface{}{},
			expected: "",
		},
		{
			name: "no text fields",
			content: []interface{}{
				map[string]interface{}{"type": "image"},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTextFromContentArray(tt.content)
			if result != tt.expected {
				t.Errorf("extractTextFromContentArray() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestProcessEventStream(t *testing.T) {
	events := []Event{
		{"type": "message", "content": "First"},
		{"type": "message", "content": "Second"},
		{"type": "tool_use", "name": "Read"},
	}

	var results []*ParsedEvent
	ProcessEventStream(events, func(parsed *ParsedEvent) {
		results = append(results, parsed)
	})

	if len(results) != 3 {
		t.Errorf("ProcessEventStream() processed %d events, expected 3", len(results))
	}

	if results[0].Text != "First" {
		t.Errorf("First event text = %v, expected 'First'", results[0].Text)
	}

	if results[2].ToolName != "Read" {
		t.Errorf("Third event tool = %v, expected 'Read'", results[2].ToolName)
	}
}
