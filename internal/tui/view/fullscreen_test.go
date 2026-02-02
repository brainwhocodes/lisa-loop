package view

import (
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestPadToFullScreen_PadsLinesAndHeight(t *testing.T) {
	out := PadToFullScreen("hi\nthere", 10, 4, lipgloss.Color("#000000"))
	plain := stripANSI(out)
	lines := strings.Split(plain, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if len(line) != 10 {
			t.Fatalf("expected line %d width 10, got %d (%q)", i, len(line), line)
		}
	}
}

var ansiRE = regexp.MustCompile("\x1b\\[[0-9;]*[A-Za-z]")

func stripANSI(s string) string {
	return ansiRE.ReplaceAllString(s, "")
}
