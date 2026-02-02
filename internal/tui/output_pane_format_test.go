package tui

import (
	"regexp"
	"strings"
	"testing"
)

func TestRenderOutputPane_RendersEscapedNewlinesAsLines(t *testing.T) {
	m := Model{
		width:  80,
		height: 20,
	}

	// Simulate an SSE-ish message that contains escaped newlines.
	m.addOutputLine("a\\nb\\nc", "agent_message")

	out := stripANSIForTest(m.renderOutputPane(80, 10))
	if strings.Contains(out, "\\n") {
		t.Fatalf("expected output not to contain literal \\\\n, got: %q", out)
	}
	if !strings.Contains(out, " a") || !strings.Contains(out, " b") || !strings.Contains(out, " c") {
		t.Fatalf("expected output to contain split lines, got: %q", out)
	}
}

var ansiREForTest = regexp.MustCompile("\x1b\\[[0-9;]*[A-Za-z]")

func stripANSIForTest(s string) string {
	return ansiREForTest.ReplaceAllString(s, "")
}
