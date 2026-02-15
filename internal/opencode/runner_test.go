package opencode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/brainwhocodes/lisa-loop/internal/config"
)

func TestHandleSSEEvent_MessageDeltaNotDuplicatedAfterEmptyStarter(t *testing.T) {
	r := &Runner{}
	var got []map[string]interface{}
	r.SetOutputCallback(func(event map[string]interface{}) {
		got = append(got, event)
	})

	r.handleSSEEvent("session-1", SSEEvent{Type: "message.part.updated", Properties: mustMarshalJSON(map[string]interface{}{
		"part": map[string]interface{}{"id": "p1", "type": "text"},
	})})
	r.handleSSEEvent("session-1", SSEEvent{Type: "message.part.updated", Properties: mustMarshalJSON(map[string]interface{}{
		"part": map[string]interface{}{"id": "p1", "type": "text", "delta": "Hi"},
	})})

	if len(got) != 1 {
		t.Fatalf("expected exactly 1 emitted message event, got %d", len(got))
	}
	if got[0]["type"] != "message" {
		t.Fatalf("expected type=message, got %v", got[0]["type"])
	}
	if got[0]["content"] != "Hi" {
		t.Fatalf("expected non-duplicated content 'Hi', got %v", got[0]["content"])
	}
}

func TestHandleSSEEvent_MessageDeltaAggregates(t *testing.T) {
	r := &Runner{}
	var got []map[string]interface{}
	r.SetOutputCallback(func(event map[string]interface{}) {
		got = append(got, event)
	})

	first := SSEEvent{Type: "message.part.updated", Properties: mustMarshalJSON(map[string]interface{}{
		"part": map[string]interface{}{"id": "p1", "type": "text", "text": "Hel"},
	})}
	second := SSEEvent{Type: "message.part.updated", Properties: mustMarshalJSON(map[string]interface{}{
		"part": map[string]interface{}{"id": "p1", "type": "text", "delta": "lo"},
	})}

	r.handleSSEEvent("session-1", first)
	r.handleSSEEvent("session-1", second)

	if len(got) != 2 {
		t.Fatalf("expected 2 emitted events, got %d", len(got))
	}
	if got[1]["type"] != "message" {
		t.Fatalf("expected type=message, got %v", got[1]["type"])
	}
	if got[1]["content"] != "Hello" {
		t.Fatalf("expected aggregated message 'Hello', got %v", got[1]["content"])
	}
}

func TestHandleSSEEvent_ReasoningDeltaAggregates(t *testing.T) {
	r := &Runner{}
	var got []map[string]interface{}
	r.SetOutputCallback(func(event map[string]interface{}) {
		got = append(got, event)
	})

	r.handleSSEEvent("session-1", SSEEvent{Type: "message.part.updated", Properties: mustMarshalJSON(map[string]interface{}{
		"part": map[string]interface{}{"id": "r1", "type": "reasoning", "text": "Think"},
	})})
	r.handleSSEEvent("session-1", SSEEvent{Type: "message.part.updated", Properties: mustMarshalJSON(map[string]interface{}{
		"part": map[string]interface{}{"id": "r1", "type": "reasoning", "delta": "ing"},
	})})

	if len(got) != 2 {
		t.Fatalf("expected 2 emitted events, got %d", len(got))
	}
	item, ok := got[1]["item"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected item payload, got %T", got[1]["item"])
	}
	if item["text"] != "Thinking" {
		t.Fatalf("expected aggregated reasoning 'Thinking', got %v", item["text"])
	}
}

func TestHandleSSEEvent_SessionStatusErrorIsForwarded(t *testing.T) {
	r := &Runner{}
	var got map[string]interface{}
	r.SetOutputCallback(func(event map[string]interface{}) {
		got = event
	})

	props, err := json.Marshal(SessionStatusProps{
		SessionID: "session-1",
		Status: struct {
			Type    string "json:\"type\""
			Attempt int    "json:\"attempt,omitempty\""
			Message string "json:\"message,omitempty\""
		}{Type: "error", Message: "boom"},
	})
	if err != nil {
		t.Fatalf("marshal props: %v", err)
	}

	r.handleSSEEvent("session-1", SSEEvent{Type: "session.status", Properties: props})

	if got == nil {
		t.Fatal("expected lifecycle event")
	}
	if got["type"] != "lifecycle" {
		t.Fatalf("expected lifecycle event type, got %v", got["type"])
	}
	if got["status"] != "error" {
		t.Fatalf("expected status=error, got %v", got["status"])
	}
	if got["message"] != "boom" {
		t.Fatalf("expected error message forwarded, got %v", got["message"])
	}
}

func TestHandleSSEEvent_SessionDiffEmitsToolUse(t *testing.T) {
	r := &Runner{}
	var got []map[string]interface{}
	r.SetOutputCallback(func(event map[string]interface{}) {
		got = append(got, event)
	})

	props := mustMarshalJSON(map[string]interface{}{
		"sessionID": "session-1",
		"diff": []map[string]interface{}{
			{"file": "internal/opencode/runner.go", "additions": 10, "deletions": 2},
			{"file": "README.md", "additions": 1, "deletions": 0},
		},
	})

	r.handleSSEEvent("session-1", SSEEvent{Type: "session.diff", Properties: props})

	if len(got) != 2 {
		t.Fatalf("expected 2 tool_use events, got %d", len(got))
	}
	for i, ev := range got {
		if ev["type"] != "tool_use" {
			t.Fatalf("event[%d] expected tool_use, got %v", i, ev["type"])
		}
		if ev["name"] != "apply_patch" {
			t.Fatalf("event[%d] expected apply_patch name, got %v", i, ev["name"])
		}
		if ev["status"] != "completed" {
			t.Fatalf("event[%d] expected completed status, got %v", i, ev["status"])
		}
	}
}

func TestRun_CreatesNewSessionEachCall(t *testing.T) {
	var createSessionCalls int32
	var sessionSeq int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/session":
			call := atomic.AddInt32(&createSessionCalls, 1)
			id := fmt.Sprintf("session-%d", atomic.AddInt32(&sessionSeq, 1))
			if call <= 0 {
				t.Fatalf("invalid create session call count")
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(CreateSessionResponse{ID: id, Slug: id})

		case r.Method == http.MethodGet && r.URL.Path == "/global/event":
			w.Header().Set("Content-Type", "text/event-stream")
			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Skip("Flusher not supported")
				return
			}
			fmt.Fprint(w, "event: message.updated\n")
			fmt.Fprint(w, "data: {\"properties\":{\"info\":{\"id\":\"msg-1\",\"sessionID\":\"session-1\",\"role\":\"assistant\"}}}\n\n")
			fmt.Fprint(w, "event: message.updated\n")
			fmt.Fprint(w, "data: {\"properties\":{\"info\":{\"id\":\"msg-2\",\"sessionID\":\"session-2\",\"role\":\"assistant\"}}}\n\n")
			fmt.Fprint(w, "event: session.status\n")
			fmt.Fprint(w, "data: {\"properties\":{\"sessionID\":\"session-1\",\"status\":{\"type\":\"idle\"}}}\n\n")
			fmt.Fprint(w, "event: session.status\n")
			fmt.Fprint(w, "data: {\"properties\":{\"sessionID\":\"session-2\",\"status\":{\"type\":\"idle\"}}}\n\n")
			flusher.Flush()

		case r.Method == http.MethodPost && (r.URL.Path == "/session/session-1/prompt_async" || r.URL.Path == "/session/session-2/prompt_async"):
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	r := NewRunner(config.Config{OpenCodeServerURL: server.URL, Timeout: 5})

	_, sid1, err := r.Run("first")
	if err != nil {
		t.Fatalf("first run failed: %v", err)
	}
	_, sid2, err := r.Run("second")
	if err != nil {
		t.Fatalf("second run failed: %v", err)
	}

	if sid1 == sid2 {
		t.Fatalf("expected unique session IDs per run, got same id %q", sid1)
	}
	if atomic.LoadInt32(&createSessionCalls) != 2 {
		t.Fatalf("expected 2 session creations, got %d", createSessionCalls)
	}
}

func TestShortSessionID(t *testing.T) {
	if got := shortSessionID("session"); got != "session" {
		t.Fatalf("expected unchanged short session id, got %q", got)
	}
	if got := shortSessionID("1234567890123456"); got != "123456789012..." {
		t.Fatalf("expected truncated session id, got %q", got)
	}
}
