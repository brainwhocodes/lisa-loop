package transcript

import (
	"testing"
	"time"
)

func TestBufferUpsert_ReplacesByID(t *testing.T) {
	b := New(10)
	now := time.Now()

	b.Upsert(Item{ID: "x", At: now, Role: RoleAssistant, Kind: KindMessage, Body: "a"})
	b.Upsert(Item{ID: "x", At: now.Add(time.Second), Role: RoleAssistant, Kind: KindMessage, Body: "b"})

	if b.Len() != 1 {
		t.Fatalf("expected len=1, got %d", b.Len())
	}
	if got := b.Items()[0].Body; got != "b" {
		t.Fatalf("expected replacement body, got %q", got)
	}
}

func TestBufferCapacity_TrimsOldest(t *testing.T) {
	b := New(2)
	b.Append(Item{ID: "a"})
	b.Append(Item{ID: "b"})
	b.Append(Item{ID: "c"})

	if b.Len() != 2 {
		t.Fatalf("expected len=2, got %d", b.Len())
	}
	if b.Items()[0].ID != "b" || b.Items()[1].ID != "c" {
		t.Fatalf("unexpected items: %#v", b.Items())
	}

	// Ensure index is consistent after trim.
	b.Upsert(Item{ID: "b", Body: "updated"})
	if b.Items()[0].Body != "updated" {
		t.Fatalf("expected item b to be replaceable after trim")
	}
}
