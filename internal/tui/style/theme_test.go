package style

import "testing"

func TestDefaultTheme(t *testing.T) {
	th := DefaultTheme()
	if th == nil {
		t.Fatalf("expected non-nil theme")
	}
	if th.BgBase != Pepper {
		t.Fatalf("expected BgBase Pepper, got %q", th.BgBase)
	}
	if th.Primary != Charple {
		t.Fatalf("expected Primary Charple, got %q", th.Primary)
	}
}
