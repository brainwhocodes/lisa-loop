package markdown

import (
	"hash/fnv"
	"sync"

	"github.com/charmbracelet/glamour"
)

type cacheKey struct {
	width int
	hash  uint64
}

// Renderer wraps Glamour with a tiny render cache to avoid re-rendering
// unchanged markdown on every tick/frame.
type Renderer struct {
	mu sync.Mutex

	// rendererByWidth avoids rebuilding the term renderer on resize.
	rendererByWidth map[int]*glamour.TermRenderer

	// cache stores rendered ANSI output keyed by width+content hash.
	cache map[cacheKey]string
}

func New() *Renderer {
	return &Renderer{
		rendererByWidth: make(map[int]*glamour.TermRenderer),
		cache:           make(map[cacheKey]string),
	}
}

func (r *Renderer) Render(width int, md string) (string, error) {
	if width <= 0 {
		width = 80
	}

	key := cacheKey{width: width, hash: hash64(md)}

	r.mu.Lock()
	if out, ok := r.cache[key]; ok {
		r.mu.Unlock()
		return out, nil
	}

	tr, ok := r.rendererByWidth[width]
	if !ok {
		var err error
		tr, err = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(width),
		)
		if err != nil {
			r.mu.Unlock()
			return "", err
		}
		r.rendererByWidth[width] = tr
	}
	r.mu.Unlock()

	out, err := tr.Render(md)
	if err != nil {
		return "", err
	}

	r.mu.Lock()
	r.cache[key] = out
	r.mu.Unlock()

	return out, nil
}

// InvalidateWidth clears cached renders for a specific width.
// Call this when terminal width changes materially.
func (r *Renderer) InvalidateWidth(width int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k := range r.cache {
		if k.width == width {
			delete(r.cache, k)
		}
	}
	delete(r.rendererByWidth, width)
}

func hash64(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
