package tui

import (
	"github.com/brainwhocodes/lisa-loop/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// Re-export palette and iconography from internal/tui/style to keep the tui package surface stable.

// Colors
const (
	Charple  lipgloss.Color = style.Charple
	Dolly    lipgloss.Color = style.Dolly
	Julep    lipgloss.Color = style.Julep
	Zest     lipgloss.Color = style.Zest
	Butter   lipgloss.Color = style.Butter
	Pepper   lipgloss.Color = style.Pepper
	BBQ      lipgloss.Color = style.BBQ
	Charcoal lipgloss.Color = style.Charcoal
	Iron     lipgloss.Color = style.Iron
	Salt     lipgloss.Color = style.Salt
	Ash      lipgloss.Color = style.Ash
	Smoke    lipgloss.Color = style.Smoke
	Squid    lipgloss.Color = style.Squid
	Oyster   lipgloss.Color = style.Oyster
	Guac     lipgloss.Color = style.Guac
	Sriracha lipgloss.Color = style.Sriracha
	Malibu   lipgloss.Color = style.Malibu
	Coral    lipgloss.Color = style.Coral
)

// Icons
const (
	IconCheck       = style.IconCheck
	IconError       = style.IconError
	IconWarning     = style.IconWarning
	IconInfo        = style.IconInfo
	IconPending     = style.IconPending
	IconInProgress  = style.IconInProgress
	IconArrowRight  = style.IconArrowRight
	IconBorderThin  = style.IconBorderThin
	IconBorderThick = style.IconBorderThick
	IconDiagonal    = style.IconDiagonal
)

var (
	SaxNotes     = style.SaxNotes
	ThinkingWave = style.ThinkingWave
)

type Theme = style.Theme

func DefaultTheme() *Theme { return style.DefaultTheme() }

func GradientText(text string, from, to lipgloss.Color) string {
	return style.GradientText(text, from, to)
}

func DiagonalSeparator(width int) string { return style.DiagonalSeparator(width) }
