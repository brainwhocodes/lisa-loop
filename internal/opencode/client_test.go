package opencode

import (
	"encoding/json"
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
