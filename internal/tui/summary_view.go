package tui

import (
	"path/filepath"
	"strings"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SummaryView is the full-width workspace cards overview (toggled with s).
type SummaryView struct {
	workspaces []state.Workspace
	cursor     int
	width      int
	height     int
}

// NewSummaryView constructs a SummaryView with the given workspaces.
func NewSummaryView(workspaces []state.Workspace) SummaryView {
	return SummaryView{workspaces: workspaces}
}

// SetWorkspaces replaces the workspace list, clamping cursor if needed.
func (m *SummaryView) SetWorkspaces(ws []state.Workspace) {
	m.workspaces = ws
	if len(ws) == 0 {
		m.cursor = 0
	} else if m.cursor >= len(ws) {
		m.cursor = len(ws) - 1
	}
}

// SelectedIndex returns the index of the focused card.
func (m SummaryView) SelectedIndex() int {
	return m.cursor
}

func (m SummaryView) cols() int {
	const cardTotalWidth = cardWidth + 2 // +2 for left/right borders
	if m.width < cardTotalWidth {
		return 1
	}
	return m.width / cardTotalWidth
}

// Update handles navigation between workspace cards.
func (m SummaryView) Update(msg tea.Msg) (SummaryView, tea.Cmd) {
	cols := m.cols()
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Left):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(keyMsg, keys.Right):
			if m.cursor < len(m.workspaces)-1 {
				m.cursor++
			}
		case key.Matches(keyMsg, keys.Up):
			if m.cursor >= cols {
				m.cursor -= cols
			}
		case key.Matches(keyMsg, keys.Down):
			if m.cursor+cols < len(m.workspaces) {
				m.cursor += cols
			}
		}
	}
	return m, nil
}

const cardWidth = 24

// View renders the workspace summary cards.
func (m SummaryView) View() string {
	var sb strings.Builder
	sb.WriteString(styleTitle.Render("Summary") + "\n\n")

	var cards []string
	for i, ws := range m.workspaces {
		name := filepath.Base(ws.Dir)
		dot := styleGreenDot
		if !ws.HasResources() {
			dot = styleRedDot
		}
		count := pluralize(len(ws.Resources), "resource")
		body := dot + " " + name + "\n" + styleMuted.Render(count)
		if len(ws.Tags) > 0 {
			body += "\n" + styleMuted.Render(strings.Join(ws.Tags, " "))
		}

		borderColor := colorBorder
		if i == m.cursor {
			borderColor = colorBlue
		}
		card := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1).
			Width(cardWidth).
			Render(body)
		cards = append(cards, card)
	}

	if len(cards) > 0 {
		cols := m.cols()
		var rows []string
		for i := 0; i < len(cards); i += cols {
			end := i + cols
			if end > len(cards) {
				end = len(cards)
			}
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cards[i:end]...))
		}
		sb.WriteString(lipgloss.JoinVertical(lipgloss.Left, rows...))
	} else {
		sb.WriteString(styleMuted.Render("no workspaces — press e to add"))
	}

	sb.WriteString("\n\n" + styleMuted.Render("←/→ navigate  enter=open  s=back"))
	return sb.String()
}
