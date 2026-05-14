package tui

import (
	"strings"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tagEditorPhase int

const (
	tagEditorList  tagEditorPhase = iota // navigating current tags
	tagEditorInput                       // typing a new tag to add
)

// TagEditor is a modal for editing tags on an existing workspace.
type TagEditor struct {
	dir    string
	tags   []string
	cursor int
	phase  tagEditorPhase
	input  tagInput
	width  int
}

// NewTagEditor constructs a TagEditor for the given workspace.
// pool is the set of all known tags for autocomplete suggestions.
func NewTagEditor(ws state.Workspace, pool []string) TagEditor {
	tags := make([]string, len(ws.Tags))
	copy(tags, ws.Tags)
	return TagEditor{
		dir:   ws.Dir,
		tags:  tags,
		input: newTagInput(pool),
	}
}

// Update handles key events for the tag editor modal.
func (m TagEditor) Update(msg tea.Msg) (TagEditor, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch m.phase {
		case tagEditorInput:
			switch {
			case key.Matches(keyMsg, keys.Enter):
				tags := parseTags(m.input.Value())
				m.tags = append(m.tags, tags...)
				m.input.Reset()
				m.phase = tagEditorList
				return m, nil
			case key.Matches(keyMsg, keys.Esc):
				m.input.Reset()
				m.phase = tagEditorList
				return m, nil
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd

		default: // tagEditorList
			switch {
			case key.Matches(keyMsg, keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}
			case key.Matches(keyMsg, keys.Down):
				if m.cursor < len(m.tags)-1 {
					m.cursor++
				}
			case key.Matches(keyMsg, keys.Delete):
				if len(m.tags) > 0 {
					m.tags = append(m.tags[:m.cursor], m.tags[m.cursor+1:]...)
					if m.cursor >= len(m.tags) && m.cursor > 0 {
						m.cursor--
					}
				}
			case key.Matches(keyMsg, keys.Enter), key.Matches(keyMsg, keys.Tab):
				m.phase = tagEditorInput
				m.input.Focus()
				return m, textinput.Blink
			case key.Matches(keyMsg, keys.Save):
				tags := make([]string, len(m.tags))
				copy(tags, m.tags)
				return m, func() tea.Msg { return tagsUpdatedMsg{dir: m.dir, tags: tags} }
			case key.Matches(keyMsg, keys.Esc):
				return m, func() tea.Msg { return tagsCancelledMsg{} }
			}
		}
	}
	return m, nil
}

// View renders the tag editor as a centered modal.
func (m TagEditor) View() string {
	var sb strings.Builder
	sb.WriteString(styleTitle.Render("Edit Tags") + "\n")
	sb.WriteString(styleMuted.Render(shortenPath(m.dir)) + "\n\n")

	if len(m.tags) == 0 {
		sb.WriteString(styleMuted.Render("no tags") + "\n")
	} else {
		sb.WriteString(styleHeader.Render("current tags:") + "\n")
		for i, t := range m.tags {
			line := "  " + t
			if i == m.cursor && m.phase == tagEditorList {
				line = styleSelected.Render(line)
			}
			sb.WriteString(line + "\n")
		}
	}

	sb.WriteString("\n")
	if m.phase == tagEditorInput {
		sb.WriteString(styleHeader.Render("add tag:") + "\n")
		sb.WriteString(m.input.View() + "\n\n")
		sb.WriteString(styleMuted.Render("comma-separated  enter=add  esc=back"))
	} else {
		sb.WriteString(styleMuted.Render("enter=add tag  d=remove  ctrl+s=save  esc=cancel"))
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue).
		Padding(1, 2).
		Width(m.width).
		Render(sb.String())
}
