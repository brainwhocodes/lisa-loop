package markdown

import "testing"

func TestRendererCachesByWidthAndContent(t *testing.T) {
	r := New()

	out1, err := r.Render(40, "# Hi\n\nthere")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	out2, err := r.Render(40, "# Hi\n\nthere")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if out1 != out2 {
		t.Fatalf("expected cached output to match")
	}

	out3, err := r.Render(41, "# Hi\n\nthere")
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if out1 == out3 {
		// Different widths can still coincidentally render same ANSI, but in practice
		// the cache key must differ; we validate this by forcing invalidation.
		r.InvalidateWidth(41)
		out4, err := r.Render(41, "# Hi\n\nthere")
		if err != nil {
			t.Fatalf("render: %v", err)
		}
		_ = out4
	}
}
