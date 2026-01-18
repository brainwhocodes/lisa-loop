package codex

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/brainwhocodes/ralph-codex/internal/state"
)

// Runner executes Codex commands
type Runner struct {
	config Config
}

// NewRunner creates a new Codex runner
func NewRunner(config Config) *Runner {
	return &Runner{config: config}
}

// Run executes a Codex command
func (r *Runner) Run(prompt string) (output string, threadID string, err error) {
	switch r.config.Backend {
	case "cli":
		output, threadID, err = r.runCLI(prompt)
	case "sdk":
		return "", "", fmt.Errorf("SDK backend not yet implemented")
	default:
		return "", "", fmt.Errorf("unknown backend: %s", r.config.Backend)
	}

	return output, threadID, err
}

// runCLI executes Codex CLI
func (r *Runner) runCLI(prompt string) (string, string, error) {
	cmd := r.buildCLICommand(prompt)

	if r.config.Verbose {
		fmt.Printf("Executing: %s\n", cmd.String())
	}

	// Run command
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("codex execution failed: %w", err)
	}

	output := string(outputBytes)

	// Parse JSONL output
	threadID, message, _ := ParseJSONLStream(strings.Split(output, "\n"))

	// Save session ID if we got one
	if threadID != "" {
		if err := SaveSessionID(threadID); err != nil {
			return output, threadID, fmt.Errorf("failed to save session ID: %w", err)
		}
	}

	// Return message content instead of full output (or thread ID if only that)
	if message != "" {
		return message, threadID, nil
	}

	return output, threadID, nil
}

// buildCLICommand builds Codex CLI command
func (r *Runner) buildCLICommand(prompt string) *exec.Cmd {
	args := []string{
		"exec",
		"--json",
		"--skip-git-repo-check",
	}

	// Add thread ID if session exists
	if id, err := LoadSessionID(); err == nil && id != "" {
		args = append(args, "--resume", "--thread-id", id)
	}

	cmd := exec.Command("codex", args...)

	// Write prompt to stdin
	cmd.Stdin = strings.NewReader(prompt)

	return cmd
}

// Event represents a single JSONL event from Codex
type Event map[string]interface{}

// ParseJSONLLine parses a single JSONL line
func ParseJSONLLine(line string) (Event, error) {
	if line == "" || strings.TrimSpace(line) == "" {
		return nil, nil
	}

	var event Event
	err := json.Unmarshal([]byte(line), &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// EventType extracts the "event" field from an event
func EventType(event Event) string {
	if val, ok := event["event"].(string); ok {
		return val
	}
	return ""
}

// ThreadID extracts the "thread_id" field from an event
func ThreadID(event Event) string {
	if val, ok := event["thread_id"].(string); ok {
		return val
	}
	return ""
}

// MessageType extracts the "type" field from an event
func MessageType(event Event) string {
	if val, ok := event["type"].(string); ok {
		return val
	}
	return ""
}

// MessageText extracts the "text" field from an event
func MessageText(event Event) string {
	if val, ok := event["text"].(string); ok {
		return val
	}
	return ""
}

// ParseJSONLStream parses a complete JSONL stream from a reader
// Returns: threadID, accumulated message text, all events parsed
func ParseJSONLStream(lines []string) (threadID string, message string, events []Event) {
	events = make([]Event, 0, len(lines))
	b := strings.Builder{}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		event, err := ParseJSONLLine(line)
		if err != nil {
			continue
		}

		events = append(events, event)

		if EventType(event) == "thread.started" {
			tid := ThreadID(event)
			if tid != "" {
				threadID = tid
			}
		}

		if MessageType(event) == "message" || MessageType(event) == "text" {
			msg := MessageText(event)
			if msg != "" {
				b.WriteString(msg)
				b.WriteString("\n")
			}
		}
	}

	return threadID, strings.TrimSpace(b.String()), events
}

// IsJSONL checks if a line looks like JSON (starts with { or [)
func IsJSONL(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")
}

// LoadSessionID loads Codex session ID from .codex_session_id
func LoadSessionID() (string, error) {
	return state.LoadCodexSession()
}

// SaveSessionID saves Codex session ID to .codex_session_id atomically
func SaveSessionID(id string) error {
	return state.SaveCodexSession(id)
}

// NewSession creates a new session by clearing the session ID
func NewSession() error {
	if err := SaveSessionID(""); err != nil {
		return err
	}
	return nil
}

// SessionExists checks if a session ID file exists
func SessionExists() bool {
	id, err := LoadSessionID()
	if err != nil {
		return false
	}
	return id != ""
}

// SessionAgeHours calculates the age of the session in hours
func SessionAgeHours() (int, error) {
	if !SessionExists() {
		return 0, nil
	}

	sessionFile := ".codex_session_id"
	info, err := os.Stat(sessionFile)
	if err != nil {
		return 0, err
	}

	age := time.Since(info.ModTime()).Hours()
	return int(age), nil
}

// IsSessionExpired checks if the session has expired based on age
func IsSessionExpired(expiryHours int) bool {
	if expiryHours <= 0 {
		return false
	}

	age, err := SessionAgeHours()
	if err != nil {
		return false
	}

	return age >= expiryHours
}

// SessionMetadata represents session metadata
type SessionMetadata struct {
	ID        string
	CreatedAt time.Time
	LastUsed  time.Time
}

// LoadSessionMetadata loads session metadata from .ralph_session
func LoadSessionMetadata() (*SessionMetadata, error) {
	sess, err := state.LoadRalphSession()
	if err != nil {
		return nil, err
	}

	meta := &SessionMetadata{
		ID: "",
	}

	if id, ok := sess["id"].(string); ok {
		meta.ID = id
	}

	if created, ok := sess["created_at"].(string); ok {
		meta.CreatedAt, _ = time.Parse(time.RFC3339, created)
	}

	if lastUsed, ok := sess["last_used"].(string); ok {
		meta.LastUsed, _ = time.Parse(time.RFC3339, lastUsed)
	}

	return meta, nil
}

// SaveSessionMetadata saves session metadata atomically
func SaveSessionMetadata(meta *SessionMetadata) error {
	sess := map[string]interface{}{
		"id":         meta.ID,
		"created_at": meta.CreatedAt.Format(time.RFC3339),
		"last_used":  meta.LastUsed.Format(time.RFC3339),
	}

	return state.SaveRalphSession(sess)
}
