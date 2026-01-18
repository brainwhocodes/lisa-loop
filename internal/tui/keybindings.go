package tui

import (
	"fmt"
	"strings"
)

// Keybinding represents a keybinding help entry
type Keybinding struct {
	Key         string
	Description string
}

// KeybindingSection represents a section of keybindings
type KeybindingSection struct {
	Title string
	Keys  []Keybinding
}

// GetKeybindingHelp returns formatted help text for all keybindings
func GetKeybindingHelp() string {
	sections := []KeybindingSection{
		{
			Title: "Navigation",
			Keys: []Keybinding{
				{"q / Ctrl+C", "Quit Ralph Codex"},
				{"?", "Toggle help screen"},
			},
		},
		{
			Title: "Loop Control",
			Keys: []Keybinding{
				{"r", "Run / Restart loop"},
				{"p", "Pause / Resume loop"},
			},
		},
		{
			Title: "Views",
			Keys: []Keybinding{
				{"l", "Toggle log view"},
				{"c", "Show circuit breaker status"},
				{"R", "Reset circuit breaker"},
			},
		},
		{
			Title: "CLI Options",
			Keys: []Keybinding{
				{"--monitor", "Enable integrated TUI monitoring"},
				{"--verbose", "Verbose output"},
				{"--backend cli", "Use CLI backend"},
				{"--backend sdk", "Use SDK backend"},
			},
		},
		{
			Title: "Project Options",
			Keys: []Keybinding{
				{"--project <path>", "Set project directory"},
				{"--prompt <file>", "Set prompt file"},
			},
		},
		{
			Title: "Rate Limiting",
			Keys: []Keybinding{
				{"--calls <num>", "Max API calls per hour"},
				{"--timeout <sec>", "Codex execution timeout"},
			},
		},
		{
			Title: "Project Commands",
			Keys: []Keybinding{
				{"setup --name <proj>", "Create new project"},
				{"import --source <file>", "Import PRD/document"},
				{"status", "Show project status"},
				{"reset-circuit", "Reset circuit breaker"},
			},
		},
		{
			Title: "Troubleshooting",
			Keys: []Keybinding{
				{"Press 'r' after error", "Retry failed operation"},
				{"Press 'R' after loop", "Reset circuit if stuck"},
				{"Check logs with 'l'", "View detailed execution logs"},
			},
		},
	}

	var builder strings.Builder

	for _, section := range sections {
		builder.WriteString(StyleHeader.Render(section.Title))
		builder.WriteString("\n\n")

		for _, keybinding := range section.Keys {
			builder.WriteString(fmt.Sprintf("  %s %s\n",
				StyleHelpKey.Render(keybinding.Key),
				StyleHelpDesc.Render(keybinding.Description)))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// RenderHelpScreen returns the full help screen
func (m Model) renderHelpView() string {
	header := StyleHeader.Render("Ralph Codex - Help")

	version := StyleHelpDesc.Render("Version 1.0.0")

	divider := StyleDivider.Render(DividerChar)

	helpContent := GetKeybindingHelp()

	footer := fmt.Sprintf(`
%s
Press '?' to return to status view
`,
		StyleInfoMsg.Render("Tip: Use --monitor flag for TUI mode"))

	return header + "\n" + version + "\n" + divider + "\n\n" + helpContent + footer
}
