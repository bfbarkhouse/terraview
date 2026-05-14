package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorGreen  = lipgloss.Color("#4ade80")
	colorRed    = lipgloss.Color("#f87171")
	colorBlue   = lipgloss.Color("#7ec8e3")
	colorMuted  = lipgloss.Color("#555555")
	colorBorder = lipgloss.Color("#333333")

	styleFocusedBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorBlue)

	styleBlurredBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorBorder)

	styleGreenDot = lipgloss.NewStyle().Foreground(colorGreen).Render("●")
	styleRedDot   = lipgloss.NewStyle().Foreground(colorRed).Render("●")
	styleSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#facc15"))
	styleMuted  = lipgloss.NewStyle().Foreground(colorMuted)
	styleTitle  = lipgloss.NewStyle().Bold(true).Foreground(colorBlue)
	styleHeader = lipgloss.NewStyle().Foreground(colorMuted)
)
