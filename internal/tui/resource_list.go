package tui

import (
	"strings"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// ResourceList is the right panel showing resources for the selected workspace.
type ResourceList struct {
	workspace *state.Workspace
	cursor    int
	focused   bool
	width     int
	height    int
}

// NewResourceList constructs an empty ResourceList.
func NewResourceList() ResourceList {
	return ResourceList{}
}

// SetWorkspace updates the displayed workspace and resets the cursor.
func (m *ResourceList) SetWorkspace(ws *state.Workspace) {
	m.workspace = ws
	m.cursor = 0
}

// Update handles key events when the panel is focused.
func (m ResourceList) Update(msg tea.Msg) (ResourceList, tea.Cmd) {
	if !m.focused || m.workspace == nil {
		return m, nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		count := len(m.workspace.Resources)
		switch {
		case key.Matches(keyMsg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(keyMsg, keys.Down):
			if m.cursor < count-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

// View renders the resource list panel.
func (m ResourceList) View() string {
	border := styleBlurredBorder
	if m.focused {
		border = styleFocusedBorder
	}

	if m.workspace == nil {
		body := styleHeader.Render("RESOURCES") + "\n\n" +
			styleMuted.Render("select a workspace")
		return border.Width(m.width).Height(m.height).Render(body)
	}

	var sb strings.Builder
	title := shortenPath(m.workspace.Dir)
	count := len(m.workspace.Resources)
	sb.WriteString(styleHeader.Render(title+" — "+pluralize(count, "resource")) + "\n\n")

	if m.workspace.Err != nil {
		sb.WriteString(styleRedDot + " " + styleMuted.Render("no state file"))
	} else if count == 0 {
		sb.WriteString(styleMuted.Render("no resources in state"))
	} else {
		for i, r := range m.workspace.Resources {
			line := styleGreenDot + " " + r.DisplayName()
			if i == m.cursor && m.focused {
				line = styleSelected.Render(styleGreenDot + " " + r.DisplayName())
			}
			sb.WriteString(line + "\n")
		}
	}

	return border.Width(m.width).Height(m.height).Render(sb.String())
}
