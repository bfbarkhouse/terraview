package tui

import (
	"strings"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// WorkspaceList is the left panel showing configured Terraform directories.
type WorkspaceList struct {
	workspaces []state.Workspace
	cursor     int
	focused    bool
	width      int
	height     int
}

// NewWorkspaceList constructs a WorkspaceList with the given workspaces.
func NewWorkspaceList(workspaces []state.Workspace) WorkspaceList {
	return WorkspaceList{workspaces: workspaces}
}

// SetWorkspaces replaces the workspace list, clamping cursor if needed.
func (m *WorkspaceList) SetWorkspaces(ws []state.Workspace) {
	m.workspaces = ws
	if len(ws) == 0 {
		m.cursor = 0
	} else if m.cursor >= len(ws) {
		m.cursor = len(ws) - 1
	}
}

// Selected returns the currently highlighted workspace, or nil if empty.
func (m *WorkspaceList) Selected() *state.Workspace {
	if len(m.workspaces) == 0 {
		return nil
	}
	return &m.workspaces[m.cursor]
}

// SelectedIndex returns the index of the currently highlighted workspace.
func (m *WorkspaceList) SelectedIndex() int {
	return m.cursor
}

// Update handles key events when the panel is focused.
func (m WorkspaceList) Update(msg tea.Msg) (WorkspaceList, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(keyMsg, keys.Down):
			if m.cursor < len(m.workspaces)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

// View renders the workspace list panel.
func (m WorkspaceList) View() string {
	var sb strings.Builder
	sb.WriteString(styleHeader.Render("WORKSPACES") + "\n\n")

	for i, ws := range m.workspaces {
		selected := i == m.cursor && m.focused
		name := shortenPath(ws.Dir)
		// Truncate long paths to fit panel
		if len(name) > m.width-6 && m.width > 10 {
			name = "…" + name[len(name)-(m.width-7):]
		}
		prefix := "  "
		if selected {
			prefix = "▶ "
		}
		var dot string
		if selected {
			dot = styleSelected.Render("●")
		} else if ws.HasResources() {
			dot = styleGreenDot
		} else {
			dot = styleRedDot
		}
		line := prefix + dot + " " + styleSelected.Render(name)
		tags := ""
		if len(ws.Tags) > 0 {
			tags = "  " + strings.Join(ws.Tags, " ")
		}
		if !selected {
			line = prefix + dot + " " + name
		}
		if tags != "" {
			if selected {
				tags = styleSelected.Render(tags)
			} else {
				tags = styleMuted.Render(tags)
			}
		}
		sb.WriteString(line + "\n")
		if tags != "" {
			sb.WriteString(tags + "\n")
		}
	}

	if len(m.workspaces) == 0 {
		sb.WriteString(styleMuted.Render("  no workspaces") + "\n")
	}

	sb.WriteString("\n" + styleMuted.Render("e = add/edit"))

	border := styleBlurredBorder
	if m.focused {
		border = styleFocusedBorder
	}
	return border.Width(m.width).Height(m.height).Render(sb.String())
}
