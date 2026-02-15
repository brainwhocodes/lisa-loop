package opencode

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brainwhocodes/lisa-loop/internal/config"
)

// OutputCallback is called for streaming output events
type OutputCallback func(event map[string]interface{})

// Runner executes prompts using the OpenCode server API
type Runner struct {
	client         *Client
	outputCallback OutputCallback
	verbose        bool
	timeout        time.Duration
	sessionID      string // Cached session ID to avoid repeated file reads
	server         *Server
	cfg            config.Config
	// Track last emitted content to avoid duplicates (SSE sends cumulative updates)
	lastReasoning  string
	lastMessage    string
	messageParts   map[string]string
	reasoningParts map[string]string
	messageOrder   []string
	reasoningOrder []string
	// Context tracking for auto-save
	contextTracker *ContextTracker
	archiver       *SessionArchiver
	loopNumber     int
}

// NewRunner creates a new OpenCode runner from config
func NewRunner(cfg config.Config) *Runner {
	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Minute // Default timeout for streaming (long for npm tests, etc.)
	}

	runner := &Runner{
		verbose:        cfg.Verbose,
		timeout:        timeout,
		cfg:            cfg,
		contextTracker: NewContextTracker(cfg.OpenCodeModelID),
		archiver:       NewSessionArchiver(cfg.ProjectPath),
	}

	// If server URL is provided, use it directly
	if cfg.OpenCodeServerURL != "" {
		clientCfg := Config{
			ServerURL: cfg.OpenCodeServerURL,
			Username:  cfg.OpenCodeUsername,
			Password:  cfg.OpenCodePassword,
			ModelID:   cfg.OpenCodeModelID,
			Timeout:   timeout,
		}
		runner.client = NewClient(clientCfg)
	}
	// Otherwise, we'll start a managed server on first Run()

	return runner
}

// SetOutputCallback sets the callback for streaming output
func (r *Runner) SetOutputCallback(cb OutputCallback) {
	r.outputCallback = cb
}

// Run executes a prompt and returns the output, session ID, and any error
func (r *Runner) Run(prompt string) (output string, sessionID string, err error) {
	// Start managed server if needed
	if r.client == nil {
		if err := r.startManagedServer(); err != nil {
			return "", "", fmt.Errorf("failed to start managed server: %w", err)
		}
	}

	// Each loop/task run uses a fresh OpenCode session to keep iteration context isolated.
	r.emitEvent("message", map[string]interface{}{
		"content": "Creating new session for this loop...",
	})

	sessionID, err = r.client.CreateSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}

	r.emitEvent("message", map[string]interface{}{
		"content": fmt.Sprintf("Session created: %s", shortSessionID(sessionID)),
	})

	// Cache only for diagnostics/getters; do not reuse for future runs.
	r.sessionID = sessionID

	r.emitEvent("message", map[string]interface{}{
		"content": "Sending prompt to OpenCode...",
	})

	// Reset tracking for new message (SSE sends cumulative updates)
	r.lastReasoning = ""
	r.lastMessage = ""
	r.messageParts = make(map[string]string)
	r.reasoningParts = make(map[string]string)
	r.messageOrder = nil
	r.reasoningOrder = nil
	r.contextTracker.Reset()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	r.emitEvent("message", map[string]interface{}{
		"content": fmt.Sprintf("Connecting to SSE stream (timeout: %v)...", r.timeout),
	})

	// Send the message with SSE streaming
	result, err := r.client.SendMessageStreaming(ctx, sessionID, prompt, func(event SSEEvent) {
		// Parse and forward SSE events in TUI-compatible format
		r.handleSSEEvent(sessionID, event)
	})

	if err != nil {
		r.emitEvent("message.error", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		return "", sessionID, fmt.Errorf("failed to send message: %w", err)
	}

	r.emitEvent("message", map[string]interface{}{
		"content": fmt.Sprintf("Received response (%d chars)", len(result.Content)),
	})

	if err != nil {
		r.emitEvent("message.error", map[string]interface{}{
			"session_id": sessionID,
			"error":      err.Error(),
		})
		return "", sessionID, fmt.Errorf("failed to send message: %w", err)
	}

	content := result.Content
	r.emitEvent("message.received", map[string]interface{}{
		"session_id": sessionID,
		"content":    content,
	})

	// Emit the response as a message event (for TUI compatibility)
	r.emitEvent("message", map[string]interface{}{
		"type": "message",
		"text": content,
	})

	// Update context tracking with token usage from result
	usage := r.contextTracker.Update(result.PromptTokens, result.CompletionTokens, result.WasCompacted)

	// Emit context usage event for TUI
	r.emitEvent("context.usage", map[string]interface{}{
		"prompt_tokens":     usage.PromptTokens,
		"completion_tokens": usage.CompletionTokens,
		"total_tokens":      usage.TotalTokens,
		"context_limit":     usage.ContextLimit,
		"usage_percent":     usage.UsagePercent,
		"threshold_reached": usage.ThresholdReached,
		"was_compacted":     usage.WasCompacted,
	})

	// Check if we need to auto-save and start new session
	if usage.ThresholdReached {
		r.emitEvent("lifecycle", map[string]interface{}{
			"type":   "context_threshold",
			"status": "saving",
		})

		if archivePath, err := r.saveAndRotateSession(sessionID, usage, "threshold"); err != nil {
			r.emitEvent("message.error", map[string]interface{}{
				"error": fmt.Sprintf("Failed to save session: %v", err),
			})
		} else {
			r.emitEvent("lifecycle", map[string]interface{}{
				"type":         "context_threshold",
				"status":       "saved",
				"archive_path": archivePath,
			})
		}
	} else if result.WasCompacted {
		// OpenCode already compacted - save for our records
		r.emitEvent("lifecycle", map[string]interface{}{
			"type":   "session_compacted",
			"status": "detected",
		})
		if _, err := r.saveAndRotateSession(sessionID, usage, "compacted"); err != nil {
			r.emitEvent("message.error", map[string]interface{}{
				"error": fmt.Sprintf("Failed to save compacted session: %v", err),
			})
		}
	}

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
	event["type"] = eventType  // For codex event parser compatibility
	event["event"] = eventType // Legacy field

	r.outputCallback(event)
}

// NewSession clears the current session and starts fresh
func (r *Runner) NewSession() error {
	r.sessionID = "" // Clear cached session
	r.contextTracker.Reset()
	return ClearSession()
}

// SetLoopNumber sets the current loop number for archiving
func (r *Runner) SetLoopNumber(loopNum int) {
	r.loopNumber = loopNum
}

// GetContextUsage returns current context usage stats
func (r *Runner) GetContextUsage() ContextUsage {
	return r.contextTracker.GetUsage()
}

// saveAndRotateSession saves the current session and creates a new one
func (r *Runner) saveAndRotateSession(sessionID string, usage ContextUsage, reason string) (string, error) {
	// Create archive
	archive := SessionArchive{
		SessionID:        sessionID,
		ModelID:          r.cfg.OpenCodeModelID,
		PromptTokens:     usage.PromptTokens,
		CompletionTokens: usage.CompletionTokens,
		LoopNumber:       r.loopNumber,
		SavedAt:          time.Now(),
		Reason:           reason,
	}

	// Save to file
	archivePath, err := r.archiver.Save(archive)
	if err != nil {
		return "", err
	}

	// Create new session
	r.emitEvent("message", map[string]interface{}{
		"content": fmt.Sprintf("Context at %.0f%%, creating new session...", usage.UsagePercent*100),
	})

	newSessionID, err := r.client.CreateSession()
	if err != nil {
		return archivePath, fmt.Errorf("failed to create new session: %w", err)
	}

	// Save new session ID
	if err := SaveSessionID(newSessionID); err != nil {
		return archivePath, fmt.Errorf("failed to save new session ID: %w", err)
	}

	// Update runner state
	r.sessionID = newSessionID
	r.contextTracker.Reset()

	r.emitEvent("message", map[string]interface{}{
		"content": fmt.Sprintf("New session created: %s", shortSessionID(newSessionID)),
	})

	return archivePath, nil
}

// GetSessionID returns the current session ID
func (r *Runner) GetSessionID() (string, error) {
	return LoadSessionID()
}

// Stop shuts down the managed server if running
func (r *Runner) Stop() error {
	if r.server != nil {
		return r.server.Stop()
	}
	return nil
}

// handleSSEEvent parses SSE events and emits them in TUI-compatible format
func (r *Runner) handleSSEEvent(sessionID string, event SSEEvent) {
	switch event.Type {
	case "message.part.updated":
		if r.messageParts == nil {
			r.messageParts = make(map[string]string)
		}
		if r.reasoningParts == nil {
			r.reasoningParts = make(map[string]string)
		}
		var props PartUpdatedProps
		if err := json.Unmarshal(event.Properties, &props); err == nil {
			part := props.Part
			switch part.Type {
			case "reasoning":
				reasoning := r.mergePartText(part.ID, part.Text, part.Delta, r.reasoningParts, &r.reasoningOrder)
				if reasoning != "" && reasoning != r.lastReasoning {
					r.lastReasoning = reasoning
					r.emitEvent("item.completed", map[string]interface{}{
						"item": map[string]interface{}{
							"type": "reasoning",
							"text": reasoning,
						},
					})
				}
			case "text":
				message := r.mergePartText(part.ID, part.Text, part.Delta, r.messageParts, &r.messageOrder)
				if message != "" && message != r.lastMessage {
					r.lastMessage = message
					r.emitEvent("message", map[string]interface{}{
						"type":    "message",
						"content": message,
					})
				}
			case "tool":
				// Extract tool name and target
				toolName := part.Tool
				if toolName == "" {
					toolName = "tool"
				}
				target := ""
				status := "started"
				if part.State != nil {
					if part.State.Status == "completed" || part.State.Status == "done" {
						status = "completed"
					}
					// Extract target from input
					if input := part.State.Input; input != nil {
						// Try various field names for file paths
						if fp, ok := input["filePath"].(string); ok {
							target = fp
						} else if fp, ok := input["file_path"].(string); ok {
							target = fp
						} else if fp, ok := input["path"].(string); ok {
							target = fp
						} else if cmd, ok := input["command"].(string); ok {
							if len(cmd) > 50 {
								target = cmd[:50] + "..."
							} else {
								target = cmd
							}
						}
					}
				}
				// Emit tool event
				r.emitEvent("tool_use", map[string]interface{}{
					"name":   toolName,
					"target": target,
					"status": status,
				})
			}
		}

	case "session.status":
		var props SessionStatusProps
		if err := json.Unmarshal(event.Properties, &props); err == nil {
			status := props.Status
			switch status.Type {
			case "busy":
				r.emitEvent("lifecycle", map[string]interface{}{
					"type":   "status",
					"status": "processing",
				})
			case "retry":
				r.emitEvent("lifecycle", map[string]interface{}{
					"type":    "status",
					"status":  "retry",
					"attempt": status.Attempt,
					"message": status.Message,
				})
			case "idle":
				r.emitEvent("lifecycle", map[string]interface{}{
					"type":   "status",
					"status": "complete",
				})
			case "error":
				r.emitEvent("lifecycle", map[string]interface{}{
					"type":    "status",
					"status":  "error",
					"message": status.Message,
				})
			}
		}

	case "session.diff":
		var props SessionDiffProps
		if err := json.Unmarshal(event.Properties, &props); err == nil {
			for _, d := range props.Diff {
				if d.File == "" {
					continue
				}
				r.emitEvent("tool_use", map[string]interface{}{
					"name":   "apply_patch",
					"target": d.File,
					"status": "completed",
				})
			}
		}
	}
}

func (r *Runner) mergePartText(partID, text, delta string, partMap map[string]string, order *[]string) string {
	if partID == "" {
		return ""
	}
	if partMap == nil {
		return ""
	}

	if _, exists := partMap[partID]; !exists {
		partMap[partID] = ""
		*order = append(*order, partID)
	}

	if text != "" {
		partMap[partID] = text
	} else if delta != "" {
		partMap[partID] += delta
	}

	if partMap[partID] == "" {
		return ""
	}

	combined := ""
	for _, id := range *order {
		combined += partMap[id]
	}
	return combined
}

func shortSessionID(id string) string {
	const keep = 12
	if len(id) <= keep {
		return id
	}
	return id[:keep] + "..."
}

// startManagedServer starts a child OpenCode server in the project directory
func (r *Runner) startManagedServer() error {
	// Emit startup message to TUI
	r.emitEvent("message", map[string]interface{}{
		"content": "Starting OpenCode server...",
	})

	r.server = NewServer(ServerConfig{
		ProjectDir: r.cfg.ProjectPath,
		Verbose:    r.verbose,
	})

	ctx := context.Background()
	if err := r.server.Start(ctx); err != nil {
		return err
	}

	// Create client connected to managed server
	r.client = NewClient(Config{
		ServerURL: r.server.URL(),
		Username:  r.cfg.OpenCodeUsername,
		Password:  r.cfg.OpenCodePassword,
		ModelID:   r.cfg.OpenCodeModelID,
		Timeout:   r.timeout,
	})

	// Emit server ready message to TUI
	r.emitEvent("message", map[string]interface{}{
		"content": fmt.Sprintf("OpenCode server ready at %s", r.server.URL()),
	})

	return nil
}
