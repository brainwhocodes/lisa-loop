package opencode

import (
	"encoding/json"
	"testing"
)

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
