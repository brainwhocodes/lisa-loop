package tui

import (
	"strconv"
	"strings"
)

// decodeEscapes converts common escaped sequences (e.g. "\\n", "\\t", "\\uXXXX")
// into their literal equivalents. This is primarily to make streamed SSE content
// readable when upstream sends text with escaped newlines.
func decodeEscapes(s string) string {
	if s == "" || !strings.Contains(s, "\\") {
		return s
	}

	// Attempt to interpret the string as a Go-like quoted literal by wrapping it
	// in quotes and escaping any unescaped quote characters.
	out := s
	for i := 0; i < 2; i++ {
		decoded, ok := tryUnquoteEscapes(out)
		if !ok || decoded == out {
			break
		}
		out = decoded
	}
	return out
}

func tryUnquoteEscapes(s string) (string, bool) {
	// Fast-path: no obvious escape patterns.
	if !strings.Contains(s, "\\n") &&
		!strings.Contains(s, "\\t") &&
		!strings.Contains(s, "\\r") &&
		!strings.Contains(s, "\\u") &&
		!strings.Contains(s, "\\\"") &&
		!strings.Contains(s, "\\\\") {
		return s, false
	}

	b := make([]byte, 0, len(s)+2)
	b = append(b, '"')

	escaped := false
	for i := 0; i < len(s); i++ {
		c := s[i]

		// Escape unescaped quotes so the wrapper string is valid.
		if c == '"' && !escaped {
			b = append(b, '\\', '"')
			escaped = false
			continue
		}

		b = append(b, c)

		if escaped {
			escaped = false
		} else if c == '\\' {
			escaped = true
		}
	}

	b = append(b, '"')

	decoded, err := strconv.Unquote(string(b))
	if err != nil {
		return s, false
	}
	return decoded, true
}
