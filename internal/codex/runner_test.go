package codex

import (
	"os"
	"testing"
)

func TestParseJSONLLineValid(t *testing.T) {
	line := `{"event": "thread.started", "thread_id": "test-123"}`

	event, err := ParseJSONLLine(line)

	if err != nil {
		t.Errorf("ParseJSONLLine() error = %v", err)
	}

	if EventType(event) != "thread.started" {
		t.Errorf("EventType() = %s, want thread.started", EventType(event))
	}

	if ThreadID(event) != "test-123" {
		t.Errorf("ThreadID() = %s, want test-123", ThreadID(event))
	}
}

func TestParseJSONLLineMessage(t *testing.T) {
	line := `{"type": "message", "text": "Hello world"}`

	event, err := ParseJSONLLine(line)

	if err != nil {
		t.Errorf("ParseJSONLLine() error = %v", err)
	}

	if MessageType(event) != "message" {
		t.Errorf("MessageType() = %s, want message", MessageType(event))
	}

	if MessageText(event) != "Hello world" {
		t.Errorf("MessageText() = %s, want Hello world", MessageText(event))
	}
}

func TestParseJSONLStream(t *testing.T) {
	lines := []string{
		`{"event": "thread.started", "thread_id": "thread-abc"}`,
		`{"type": "message", "text": "First message"}`,
		`{"type": "message", "text": "Second message"}`,
	}

	threadID, message, _ := ParseJSONLStream(lines)

	if threadID != "thread-abc" {
		t.Errorf("ParseJSONLStream() threadID = %s, want thread-abc", threadID)
	}

	if message != "First message\nSecond message" {
		t.Errorf("ParseJSONLStream() message = %s, want both messages", message)
	}
}

func TestSessionSaveLoad(t *testing.T) {
	// Setup: create temporary dir
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	// Test: Load non-existent session
	os.Chdir(tmpDir)
	id, err := LoadSessionID()

	if err != nil {
		t.Errorf("LoadSessionID() error = %v, want nil", err)
	}

	if id != "" {
		t.Errorf("LoadSessionID() id = %s, want empty", id)
	}

	// Test: Save and verify
	testID := "session-test-456"
	err = SaveSessionID(testID)

	if err != nil {
		t.Errorf("SaveSessionID() error = %v, want nil", err)
	}

	loaded, _ := LoadSessionID()

	if loaded != testID {
		t.Errorf("LoadSessionID() loaded = %s, want %s", loaded, testID)
	}
}
