package tui

import (
	"testing"

	"github.com/brainwhocodes/lisa-loop/internal/tui/markdown"
	"github.com/brainwhocodes/lisa-loop/internal/tui/transcript"
)

func TestAddOutputLine_AssistantCumulativeUpsertsTranscript(t *testing.T) {
	m := Model{
		md:         markdown.New(),
		transcript: transcript.New(10),
	}

	m.addOutputLine("I'll", "agent_message")
	m.addOutputLine("I'll continue", "agent_message")
	m.addOutputLine("I'll continue fixing", "agent_message")

	items := m.transcript.Items()
	if len(items) != 1 {
		t.Fatalf("expected 1 transcript item, got %d: %#v", len(items), items)
	}
	if items[0].Role != transcript.RoleAssistant || items[0].Kind != transcript.KindMessage {
		t.Fatalf("unexpected transcript item: %#v", items[0])
	}
	if items[0].Body != "I'll continue fixing" {
		t.Fatalf("expected latest cumulative body, got %q", items[0].Body)
	}
}

func TestAddOutputLine_AssistantExactDuplicateDoesNotDuplicateTranscript(t *testing.T) {
	m := Model{
		md:         markdown.New(),
		transcript: transcript.New(10),
	}

	m.addOutputLine("Hello", "agent_message")
	m.addOutputLine("Hello", "agent_message") // duplicate

	items := m.transcript.Items()
	if len(items) != 1 {
		t.Fatalf("expected 1 transcript item, got %d: %#v", len(items), items)
	}
}
