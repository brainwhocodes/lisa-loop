package transcript

import "time"

type Role string

const (
	RoleUser      Role = "USER"
	RoleAssistant Role = "ASSISTANT"
	RoleTool      Role = "TOOL"
	RoleSystem    Role = "SYSTEM"
)

type Kind string

const (
	KindMessage   Kind = "message"
	KindReasoning Kind = "reasoning"
	KindToolCall  Kind = "tool"
	KindDiff      Kind = "diff"
	KindNotice    Kind = "notice"
	KindError     Kind = "error"
)

type Item struct {
	ID    string
	At    time.Time
	Role  Role
	Kind  Kind
	Title string
	Body  string
}

// Buffer is a simple append-mostly log with O(1) replace-by-ID.
type Buffer struct {
	capacity int
	items    []Item
	index    map[string]int
}

func New(capacity int) *Buffer {
	if capacity <= 0 {
		capacity = 500
	}
	return &Buffer{
		capacity: capacity,
		items:    make([]Item, 0, min(capacity, 64)),
		index:    make(map[string]int),
	}
}

func (b *Buffer) Items() []Item {
	return b.items
}

func (b *Buffer) Len() int {
	return len(b.items)
}

func (b *Buffer) Upsert(it Item) {
	if it.ID == "" {
		// No stable ID: treat as always-new.
		b.Append(it)
		return
	}
	if i, ok := b.index[it.ID]; ok && i >= 0 && i < len(b.items) {
		b.items[i] = it
		return
	}
	b.Append(it)
}

func (b *Buffer) Append(it Item) {
	b.items = append(b.items, it)
	if it.ID != "" {
		b.index[it.ID] = len(b.items) - 1
	}
	if len(b.items) > b.capacity {
		// Drop from the front; rebuild index (capacity is small).
		trim := len(b.items) - b.capacity
		b.items = b.items[trim:]
		b.rebuildIndex()
	}
}

func (b *Buffer) rebuildIndex() {
	b.index = make(map[string]int, len(b.items))
	for i, it := range b.items {
		if it.ID != "" {
			b.index[it.ID] = i
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
