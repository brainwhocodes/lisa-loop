package views

import (
	"strings"

	"github.com/brainwhocodes/ralph-codex/internal/tui"
)

// Keybinding represents a keybinding entry
type Keybinding struct {
	Key  string
	Desc string
}

// KeybindingSection represents a section of keybindings
type KeybindingSection struct {
	Title string
	Keys  []Keybinding
}

// HelpViewModel manages help view
type HelpViewModel struct {
	sections []KeybindingSection
}

// NewHelpViewModel creates a new help view model
func NewHelpViewModel() *HelpViewModel {
	return &HelpViewModel{
		sections: GetKeybindings(),
	}
}

// View renders help view
func (m *HelpViewModel) View() string {
	return RenderHelp(m.sections)
}

// GetKeybindings returns all keybindings organized by section
func GetKeybindings() []KeybindingSection {
	return []KeybindingSection{
		{
			Title: "Navigation",
			Keys: []Keybinding{
				{Key: "q", Desc: "Quit Ralph Codex"},
				{Key: "Ctrl+c", Desc: "Force quit"},
				{Key: "Tab", Desc: "Switch between views"},
				{Key: "?", Desc: "Show/hide this help"},
			},
		},
		{
			Title: "Log View",
			Keys: []Keybinding{
				{Key: "↑ / k", Desc: "Scroll up"},
				{Key: "↓ / j", Desc: "Scroll down"},
				{Key: "PgUp", Desc: "Page up"},
				{Key: "PgDn", Desc: "Page down"},
				{Key: "Home", Desc: "Jump to top"},
				{Key: "End", Desc: "Jump to bottom"},
				{Key: "c", Desc: "Clear logs"},
			},
		},
		{
			Title: "Loop Control",
			Keys: []Keybinding{
				{Key: "p", Desc: "Pause/resume loop"},
				{Key: "s", Desc: "Skip current iteration"},
				{Key: "r", Desc: "Reset circuit breaker"},
				{Key: "Ctrl+r", Desc: "Force reset session"},
			},
		},
		{
			Title: "Status Information",
			Keys: []Keybinding{
				{Key: "Circuit: CLOSED", Desc: "Normal operation"},
				{Key: "Circuit: HALF_OPEN", Desc: "Monitoring recovery"},
				{Key: "Circuit: OPEN", Desc: "Halted - check logs"},
			},
		},
	}
}

// RenderHelp renders the help screen
func RenderHelp(sections []KeybindingSection) string {
	var builder strings.Builder

	builder.WriteString(tui.StyleHeader.Render("\n  Ralph Codex - Help\n\n"))

	for _, section := range sections {
		builder.WriteString(tui.StyleHelpKey.Render("\n  " + section.Title + ":\n"))
		builder.WriteString("\n")

		maxKeyLen := 0
		for _, kb := range section.Keys {
			if len(kb.Key) > maxKeyLen {
				maxKeyLen = len(kb.Key)
			}
		}

		for _, kb := range section.Keys {
			keyPad := strings.Repeat(" ", maxKeyLen-len(kb.Key))
			builder.WriteString(tui.StyleHelpKey.Render("    "+kb.Key+keyPad) + "  ")
			builder.WriteString(tui.StyleHelpDesc.Render(kb.Desc) + "\n")
		}
	}

	builder.WriteString("\n")
	builder.WriteString(tui.StyleHelpDesc.Render("  Press any key to close help\n\n"))

	return builder.String()
}

// GetKeybindingKey returns the key for a keybinding
func (kb *Keybinding) GetKey() string {
	return kb.Key
}

// GetKeybindingDesc returns the description for a keybinding
func (kb *Keybinding) GetDesc() string {
	return kb.Desc
}

// GetSectionTitle returns the title for a section
func (ks *KeybindingSection) GetTitle() string {
	return ks.Title
}

// GetSectionKeys returns the keys for a section
func (ks *KeybindingSection) GetKeys() []Keybinding {
	return ks.Keys
}

// GetSections returns all sections
func (m *HelpViewModel) GetSections() []KeybindingSection {
	return m.sections
}

// FindKeybinding finds a keybinding by key name
func (m *HelpViewModel) FindKeybinding(key string) (*Keybinding, int, int) {
	for i, section := range m.sections {
		for j, kb := range section.Keys {
			if kb.Key == key {
				return &kb, i, j
			}
		}
	}
	return nil, -1, -1
}
