package opencode

import (
	"fmt"
	"time"

	"github.com/brainwhocodes/ralph-codex/internal/config"
)

// OutputCallback is called for streaming output events
type OutputCallback func(event map[string]interface{})

// Runner executes prompts using the OpenCode server API
type Runner struct {
	client         *Client
	outputCallback OutputCallback
	verbose        bool
}

// NewRunner creates a new OpenCode runner from config
func NewRunner(cfg config.Config) *Runner {
	clientCfg := Config{
		ServerURL: cfg.OpenCodeServerURL,
		Username:  cfg.OpenCodeUsername,
		Password:  cfg.OpenCodePassword,
		ModelID:   cfg.OpenCodeModelID,
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}

	return &Runner{
		client:  NewClient(clientCfg),
		verbose: cfg.Verbose,
	}
}

// SetOutputCallback sets the callback for streaming output
func (r *Runner) SetOutputCallback(cb OutputCallback) {
	r.outputCallback = cb
}

// Run executes a prompt and returns the output, session ID, and any error
func (r *Runner) Run(prompt string) (output string, sessionID string, err error) {
	// Try to load existing session
	sessionID, err = LoadSessionID()
	if err != nil {
		return "", "", fmt.Errorf("failed to load session: %w", err)
	}

	// Create new session if none exists
	if sessionID == "" {
		if r.verbose {
			fmt.Println("Creating new OpenCode session...")
		}

		sessionID, err = r.client.CreateSession()
		if err != nil {
			return "", "", fmt.Errorf("failed to create session: %w", err)
		}

		if err := SaveSessionID(sessionID); err != nil {
			return "", sessionID, fmt.Errorf("failed to save session ID: %w", err)
		}

		r.emitEvent("session.created", map[string]interface{}{
			"session_id": sessionID,
		})
	}

	if r.verbose {
		fmt.Printf("Using session: %s\n", sessionID)
		fmt.Printf("Sending prompt to OpenCode server...\n")
	}

	r.emitEvent("message.sending", map[string]interface{}{
		"session_id": sessionID,
	})

	// Send the message
	resp, err := r.client.SendMessage(sessionID, prompt)
	if err != nil {
		r.emitEvent("message.error", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		return "", sessionID, fmt.Errorf("failed to send message: %w", err)
	}

	content := resp.Content()
	r.emitEvent("message.received", map[string]interface{}{
		"session_id": sessionID,
		"content":    content,
	})

	// Emit the response as a message event (for TUI compatibility)
	r.emitEvent("message", map[string]interface{}{
		"type": "message",
		"text": content,
	})

	return content, sessionID, nil
}

// emitEvent sends an event to the output callback if set
func (r *Runner) emitEvent(eventType string, data map[string]interface{}) {
	if r.outputCallback == nil {
		return
	}

	event := make(map[string]interface{})
	for k, v := range data {
		event[k] = v
	}
	event["event"] = eventType

	r.outputCallback(event)
}

// NewSession clears the current session and starts fresh
func (r *Runner) NewSession() error {
	return ClearSession()
}

// GetSessionID returns the current session ID
func (r *Runner) GetSessionID() (string, error) {
	return LoadSessionID()
}
