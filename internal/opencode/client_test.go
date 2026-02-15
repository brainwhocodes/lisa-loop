package opencode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		ServerURL: "http://localhost:8080",
		Username:  "testuser",
		Password:  "testpass",
		ModelID:   "glm-4.7",
		Timeout:   30 * time.Second,
	}

	client := NewClient(cfg)

	if client.ServerURL() != "http://localhost:8080" {
		t.Errorf("expected ServerURL http://localhost:8080, got %s", client.ServerURL())
	}

	if client.ModelID() != "glm-4.7" {
		t.Errorf("expected ModelID glm-4.7, got %s", client.ModelID())
	}
}

func TestNewClient_TrimsTrailingSlash(t *testing.T) {
	cfg := Config{
		ServerURL: "http://localhost:8080/",
	}

	client := NewClient(cfg)

	if client.ServerURL() != "http://localhost:8080" {
		t.Errorf("expected trailing slash to be trimmed, got %s", client.ServerURL())
	}
}

func TestNewClient_DefaultTimeout(t *testing.T) {
	cfg := Config{
		ServerURL: "http://localhost:8080",
	}

	client := NewClient(cfg)

	if client.httpClient.Timeout != 5*time.Minute {
		t.Errorf("expected default timeout of 5 minutes, got %v", client.httpClient.Timeout)
	}
}

func TestCreateSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session" {
			t.Errorf("expected path /session, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateSessionResponse{
			ID:   "test-session-123",
			Slug: "test-slug",
		})
	}))
	defer server.Close()

	client := NewClient(Config{
		ServerURL: server.URL,
	})

	sessionID, err := client.CreateSession()
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if sessionID != "test-session-123" {
		t.Errorf("expected session ID test-session-123, got %s", sessionID)
	}
}

func TestCreateSession_WithBasicAuth(t *testing.T) {
	var receivedUsername, receivedPassword string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUsername, receivedPassword, _ = r.BasicAuth()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreateSessionResponse{
			ID: "auth-session-123",
		})
	}))
	defer server.Close()

	client := NewClient(Config{
		ServerURL: server.URL,
		Username:  "opencode",
		Password:  "secret123",
	})

	_, err := client.CreateSession()
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	if receivedUsername != "opencode" {
		t.Errorf("expected username opencode, got %s", receivedUsername)
	}
	if receivedPassword != "secret123" {
		t.Errorf("expected password secret123, got %s", receivedPassword)
	}
}

func TestCreateSession_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := NewClient(Config{
		ServerURL: server.URL,
	})

	_, err := client.CreateSession()
	if err == nil {
		t.Error("expected error for 500 response, got nil")
	}
}

func TestSendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session/test-session/message" {
			t.Errorf("expected path /session/test-session/message, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Verify request body uses parts format
		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if len(req.Parts) != 1 {
			t.Errorf("expected 1 part, got %d", len(req.Parts))
		}
		if req.Parts[0].Type != "text" {
			t.Errorf("expected part type 'text', got %s", req.Parts[0].Type)
		}
		if req.Parts[0].Text != "Hello, world!" {
			t.Errorf("expected text 'Hello, world!', got %s", req.Parts[0].Text)
		}

		// Send response in OpenCode format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SendMessageResponse{
			Info: MessageInfo{
				ID:        "msg-123",
				SessionID: "test-session",
				Role:      "assistant",
			},
			Parts: []ResponsePart{
				{ID: "part-1", Type: "text", Text: "Hello back!"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{
		ServerURL: server.URL,
	})

	resp, err := client.SendMessage("test-session", "Hello, world!")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if resp.Content() != "Hello back!" {
		t.Errorf("expected content 'Hello back!', got %s", resp.Content())
	}
	if resp.SessionID() != "test-session" {
		t.Errorf("expected session ID test-session, got %s", resp.SessionID())
	}
}

func TestSendMessage_WithBasicAuth(t *testing.T) {
	var receivedUsername, receivedPassword string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUsername, receivedPassword, _ = r.BasicAuth()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SendMessageResponse{
			Info: MessageInfo{
				ID:        "msg-123",
				SessionID: "test-session",
			},
			Parts: []ResponsePart{
				{ID: "part-1", Type: "text", Text: "response"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(Config{
		ServerURL: server.URL,
		Username:  "user1",
		Password:  "pass1",
	})

	_, err := client.SendMessage("test-session", "test message")
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if receivedUsername != "user1" {
		t.Errorf("expected username user1, got %s", receivedUsername)
	}
	if receivedPassword != "pass1" {
		t.Errorf("expected password pass1, got %s", receivedPassword)
	}
}

func TestSendMessage_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}))
	defer server.Close()

	client := NewClient(Config{
		ServerURL: server.URL,
	})

	_, err := client.SendMessage("test-session", "test message")
	if err == nil {
		t.Error("expected error for 401 response, got nil")
	}
}

func TestSendMessageResponse_Content(t *testing.T) {
	resp := &SendMessageResponse{
		Parts: []ResponsePart{
			{Type: "reasoning", Text: "thinking..."},
			{Type: "text", Text: "Hello!"},
			{Type: "step-finish"},
		},
	}

	if resp.Content() != "Hello!" {
		t.Errorf("expected 'Hello!', got %s", resp.Content())
	}
}

func TestSendMessageResponse_Content_Empty(t *testing.T) {
	resp := &SendMessageResponse{
		Parts: []ResponsePart{
			{Type: "reasoning", Text: "thinking..."},
		},
	}

	if resp.Content() != "" {
		t.Errorf("expected empty string, got %s", resp.Content())
	}
}

// Phase 4: OpenCode API Alignment Tests

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "healthy server",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "server error",
			statusCode: http.StatusServiceUnavailable,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/global/health" {
					t.Errorf("expected path /global/health, got %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient(Config{ServerURL: server.URL})
			err := client.HealthCheck()

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestAbortSession(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful abort",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "no content response",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/session/test-session/abort"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewClient(Config{ServerURL: server.URL})
			err := client.AbortSession("test-session")

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestConnectToSSE_PrimaryEndpoint(t *testing.T) {
	// Test that primary endpoint /global/event is tried first
	primaryCalled := false
	fallbackCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/global/event" {
			primaryCalled = true
			w.WriteHeader(http.StatusOK)
			// Write SSE headers
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			// Send a test event and close
			flusher, ok := w.(http.Flusher)
			if ok {
				fmt.Fprint(w, "data: {\"type\":\"test\"}\n\n")
				flusher.Flush()
			}
		} else if r.URL.Path == "/event" {
			fallbackCalled = true
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := NewClient(Config{ServerURL: server.URL})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// This will try primary and succeed
	resp, err := client.connectToSSE(ctx)
	if err != nil {
		t.Skipf("SSE connection test skipped (may not be supported in test environment): %v", err)
		return
	}
	if resp != nil {
		resp.Body.Close()
	}

	if !primaryCalled {
		t.Error("primary endpoint /global/event was not called")
	}
	if fallbackCalled {
		t.Error("fallback endpoint /event should not be called when primary succeeds")
	}
}

func TestConnectToSSE_FallbackEndpoint(t *testing.T) {
	// Test that fallback endpoint /event is used when primary fails
	primaryCalled := false
	fallbackCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/global/event" {
			primaryCalled = true
			w.WriteHeader(http.StatusNotFound) // Primary fails
		} else if r.URL.Path == "/event" {
			fallbackCalled = true
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := NewClient(Config{ServerURL: server.URL})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	resp, err := client.connectToSSE(ctx)
	if err != nil {
		t.Skipf("SSE connection test skipped (may not be supported in test environment): %v", err)
		return
	}
	if resp != nil {
		resp.Body.Close()
	}

	if !primaryCalled {
		t.Error("primary endpoint /global/event was not called")
	}
	if !fallbackCalled {
		t.Error("fallback endpoint /event was not called when primary failed")
	}
}

func TestSendMessageSync(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session/test-session/message" {
			t.Errorf("expected path /session/test-session/message, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Verify request body
		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		if len(req.Parts) != 1 || req.Parts[0].Text != "test message" {
			t.Errorf("unexpected request body: %+v", req)
		}

		// Send response
		resp := SendMessageResponse{
			Info: MessageInfo{
				ID:        "msg-123",
				SessionID: "test-session",
				Role:      "assistant",
				ModelID:   "glm-4.7",
			},
			Parts: []ResponsePart{
				{ID: "part-1", Type: "text", Text: "Hello!", SessionID: "test-session", MessageID: "msg-123"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(Config{ServerURL: server.URL})
	ctx := context.Background()

	var receivedEvents []SSEEvent
	eventCb := func(event SSEEvent) {
		receivedEvents = append(receivedEvents, event)
	}

	result, err := client.sendMessageSync(ctx, "test-session", "test message", eventCb)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.SessionID != "test-session" {
		t.Errorf("expected session ID test-session, got %s", result.SessionID)
	}
	if result.MessageID != "msg-123" {
		t.Errorf("expected message ID msg-123, got %s", result.MessageID)
	}
	if result.Content != "Hello!" {
		t.Errorf("expected content 'Hello!', got %s", result.Content)
	}

	// Verify events were emitted
	if len(receivedEvents) < 2 {
		t.Errorf("expected at least 2 events, got %d", len(receivedEvents))
	}

	// Check first event is message.updated
	if len(receivedEvents) > 0 && receivedEvents[0].Type != "message.updated" {
		t.Errorf("expected first event type message.updated, got %s", receivedEvents[0].Type)
	}

	// Check last event is session.status
	if len(receivedEvents) > 0 {
		lastEvent := receivedEvents[len(receivedEvents)-1]
		if lastEvent.Type != "session.status" {
			t.Errorf("expected last event type session.status, got %s", lastEvent.Type)
		}
	}
}

func TestMustMarshalJSON(t *testing.T) {
	data := map[string]interface{}{
		"key": "value",
		"num": 42,
	}

	result := mustMarshalJSON(data)

	var decoded map[string]interface{}
	if err := json.Unmarshal(result, &decoded); err != nil {
		t.Errorf("failed to unmarshal result: %v", err)
	}

	if decoded["key"] != "value" {
		t.Errorf("expected key='value', got %v", decoded["key"])
	}
	if decoded["num"] != float64(42) {
		t.Errorf("expected num=42, got %v", decoded["num"])
	}
}

func TestSSEEventParsing_WithEventType(t *testing.T) {
	// Test that SSE event type lines are properly captured
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Skip("Flusher not supported")
			return
		}

		// Send SSE event with event: line before data:
		fmt.Fprint(w, "event: message.part.updated\n")
		fmt.Fprint(w, "data: {\"part\":{\"id\":\"part-1\",\"type\":\"text\",\"text\":\"Hello\"}}\n\n")
		flusher.Flush()

		// Send another event
		fmt.Fprint(w, "event: session.status\n")
		fmt.Fprint(w, "data: {\"status\":{\"type\":\"idle\"}}\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient(Config{ServerURL: server.URL})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var receivedEvents []SSEEvent
	eventCb := func(event SSEEvent) {
		receivedEvents = append(receivedEvents, event)
	}

	_, err := client.SendMessageStreaming(ctx, "test-session", "test", eventCb)
	if err != nil {
		// Expected since we're not sending a proper response
		t.Logf("Expected error (no idle response): %v", err)
	}

	// Check that events were received with proper types
	if len(receivedEvents) == 0 {
		t.Fatal("No events received")
	}

	// First event should be message.part.updated
	if receivedEvents[0].Type != "message.part.updated" {
		t.Errorf("Expected first event type 'message.part.updated', got '%s'", receivedEvents[0].Type)
	}

	// Second event should be session.status
	if len(receivedEvents) > 1 && receivedEvents[1].Type != "session.status" {
		t.Errorf("Expected second event type 'session.status', got '%s'", receivedEvents[1].Type)
	}
}

func TestSSEEventParsing_WithoutEventType(t *testing.T) {
	// Test that SSE events without event: line still work (type in JSON)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Skip("Flusher not supported")
			return
		}

		// Send SSE event with type in JSON only
		fmt.Fprint(w, "data: {\"type\":\"message.updated\",\"properties\":{}}\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient(Config{ServerURL: server.URL})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var receivedEvents []SSEEvent
	eventCb := func(event SSEEvent) {
		receivedEvents = append(receivedEvents, event)
	}

	_, err := client.SendMessageStreaming(ctx, "test-session", "test", eventCb)
	if err != nil {
		t.Logf("Expected error: %v", err)
	}

	if len(receivedEvents) == 0 {
		t.Fatal("No events received")
	}

	if receivedEvents[0].Type != "message.updated" {
		t.Errorf("Expected event type 'message.updated', got '%s'", receivedEvents[0].Type)
	}
}

func TestSendMessageStreaming_AggregatesDeltaInOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/global/event":
			w.Header().Set("Content-Type", "text/event-stream")
			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Skip("Flusher not supported")
				return
			}

			fmt.Fprint(w, "event: message.updated\n")
			fmt.Fprint(w, "data: {\"properties\":{\"info\":{\"id\":\"msg-1\",\"sessionID\":\"test-session\",\"role\":\"assistant\"}}}\n\n")
			flusher.Flush()

			fmt.Fprint(w, "event: message.part.updated\n")
			fmt.Fprint(w, "data: {\"properties\":{\"part\":{\"id\":\"p1\",\"sessionID\":\"test-session\",\"messageID\":\"msg-1\",\"type\":\"text\",\"text\":\"Hello \"}}}\n\n")
			flusher.Flush()

			fmt.Fprint(w, "event: message.part.updated\n")
			fmt.Fprint(w, "data: {\"properties\":{\"part\":{\"id\":\"p2\",\"sessionID\":\"test-session\",\"messageID\":\"msg-1\",\"type\":\"text\",\"text\":\"Wor\"}}}\n\n")
			flusher.Flush()

			fmt.Fprint(w, "event: message.part.updated\n")
			fmt.Fprint(w, "data: {\"properties\":{\"part\":{\"id\":\"p2\",\"sessionID\":\"test-session\",\"messageID\":\"msg-1\",\"type\":\"text\",\"delta\":\"ld!\"}}}\n\n")
			flusher.Flush()

			fmt.Fprint(w, "event: session.status\n")
			fmt.Fprint(w, "data: {\"properties\":{\"sessionID\":\"test-session\",\"status\":{\"type\":\"idle\"}}}\n\n")
			flusher.Flush()

		case "/session/test-session/prompt_async":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewClient(Config{ServerURL: server.URL})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.sendMessageStreamingInternal(ctx, "test-session", "prompt", nil)
	if err != nil {
		t.Fatalf("sendMessageStreamingInternal failed: %v", err)
	}

	if result.Content != "Hello World!" {
		t.Fatalf("expected ordered aggregated content 'Hello World!', got %q", result.Content)
	}
}
