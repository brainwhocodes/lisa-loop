package tui

import "testing"

func TestDecodeEscapes_Newlines(t *testing.T) {
	in := "line1\\nline2"
	out := decodeEscapes(in)
	if out == in {
		t.Fatalf("expected decode to change string")
	}
	if out != "line1\nline2" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestDecodeEscapes_NoChange(t *testing.T) {
	in := "plain"
	out := decodeEscapes(in)
	if out != in {
		t.Fatalf("expected no change, got %q", out)
	}
}
