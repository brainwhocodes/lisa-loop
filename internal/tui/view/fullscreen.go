package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// PadToFullScreen pads content to fill the entire terminal and applies a uniform background.
// This is pure rendering logic (no model dependencies) to keep Model.View readable and testable.
func PadToFullScreen(content string, width, height int, background lipgloss.Color) string {
	lines := strings.Split(content, "\n")

	// Pad each line to full width.
	padded := make([]string, 0, len(lines))
	for _, line := range lines {
		lineLen := lipgloss.Width(line)
		if lineLen < width {
			line = line + strings.Repeat(" ", width-lineLen)
		}
		padded = append(padded, line)
	}

	// Add empty lines to fill height.
	for len(padded) < height {
		padded = append(padded, strings.Repeat(" ", width))
	}

	// Apply background color to entire output.
	result := strings.Join(padded[:height], "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(background).
		Render(result)
}
