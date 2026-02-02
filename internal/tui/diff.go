package tui

import (
	"path/filepath"
	"strings"
	"time"
)

type pendingChange struct {
	Path      string
	Tool      string
	Status    string
	UpdatedAt time.Time
	Verified  bool
}

func looksLikeFileTarget(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	// Heuristic: anything with a path separator or a file extension is "file-ish".
	if strings.Contains(s, "/") || strings.Contains(s, string(filepath.Separator)) {
		return true
	}
	ext := filepath.Ext(s)
	return ext != ""
}

func isDiffRelevantTool(toolName string) bool {
	switch toolName {
	case "write", "apply_patch", "edit", "patch":
		return true
	default:
		return false
	}
}
