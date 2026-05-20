package tui

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/bfbarkhouse-redpanda/terraview/internal/config"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dirEditorPhase int

const (
	phaseBrowse    dirEditorPhase = iota
	phasePathInput                // user is typing a new directory path
	phaseTagInput                 // user is typing tags for the pending path
	phaseScanConfirm              // user is confirming results of a recursive scan
)

// DirEditor is a modal overlay for adding/removing workspace directories.
type DirEditor struct {
	entries     []config.WorkspaceEntry
	cursor      int
	pathInput   textinput.Model
	tagInput    tagInput
	pendingDir  string
	recursive   bool
	scannedDirs []string
	phase       dirEditorPhase
	width       int
}

// NewDirEditor constructs a DirEditor pre-populated with entries.
// tagPool is the set of all known tags offered for autocomplete.
func NewDirEditor(entries []config.WorkspaceEntry, tagPool []string) DirEditor {
	copied := make([]config.WorkspaceEntry, len(entries))
	copy(copied, entries)

	pi := textinput.New()
	pi.Placeholder = "/path/to/terraform/dir"
	pi.CharLimit = 256

	return DirEditor{
		entries:   copied,
		pathInput: pi,
		tagInput:  newTagInput(tagPool),
	}
}

// Entries returns a copy of the current entry list.
func (m DirEditor) Entries() []config.WorkspaceEntry {
	result := make([]config.WorkspaceEntry, len(m.entries))
	copy(result, m.entries)
	return result
}

// Update handles keystrokes for the dir editor modal.
func (m DirEditor) Update(msg tea.Msg) (DirEditor, tea.Cmd) {
	var cmd tea.Cmd

	switch m.phase {
	case phaseTagInput:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(keyMsg, keys.Enter):
				tags := parseTags(m.tagInput.Value())
				m.entries = append(m.entries, config.WorkspaceEntry{Dir: m.pendingDir, Tags: tags})
				m.tagInput.Reset()
				m.pendingDir = ""
				m.phase = phaseBrowse
				return m, nil
			case key.Matches(keyMsg, keys.Esc):
				m.entries = append(m.entries, config.WorkspaceEntry{Dir: m.pendingDir})
				m.tagInput.Reset()
				m.pendingDir = ""
				m.phase = phaseBrowse
				return m, nil
			}
		}
		m.tagInput, cmd = m.tagInput.Update(msg)
		return m, cmd

	case phasePathInput:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(keyMsg, keys.Enter):
				val := strings.TrimSpace(m.pathInput.Value())
				m.pathInput.Reset()
				if val != "" {
					expanded := expandPath(val)
					if m.recursive {
						m.scannedDirs = scanForStateFiles(expanded)
						m.phase = phaseScanConfirm
						return m, nil
					}
					m.pendingDir = expanded
					m.phase = phaseTagInput
					m.tagInput.Focus()
					return m, textinput.Blink
				}
				m.phase = phaseBrowse
				return m, nil
			case key.Matches(keyMsg, keys.Tab):
				m.recursive = !m.recursive
				return m, nil
			case key.Matches(keyMsg, keys.Esc):
				m.pathInput.Reset()
				m.phase = phaseBrowse
				return m, nil
			}
		}
		m.pathInput, cmd = m.pathInput.Update(msg)
		return m, cmd

	case phaseScanConfirm:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(keyMsg, keys.Enter):
				for _, d := range m.scannedDirs {
					m.entries = append(m.entries, config.WorkspaceEntry{Dir: d})
				}
				m.scannedDirs = nil
				m.recursive = false
				m.phase = phaseBrowse
				return m, nil
			case key.Matches(keyMsg, keys.Esc):
				m.scannedDirs = nil
				m.recursive = false
				m.phase = phaseBrowse
				return m, nil
			}
		}
		return m, nil

	default: // phaseBrowse
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch {
			case key.Matches(keyMsg, keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}
			case key.Matches(keyMsg, keys.Down):
				if m.cursor < len(m.entries)-1 {
					m.cursor++
				}
			case key.Matches(keyMsg, keys.Delete):
				if len(m.entries) > 0 {
					m.entries = append(m.entries[:m.cursor], m.entries[m.cursor+1:]...)
					if m.cursor >= len(m.entries) && m.cursor > 0 {
						m.cursor--
					}
				}
			case key.Matches(keyMsg, keys.Enter):
				m.phase = phasePathInput
				m.pathInput.Focus()
				cmd = textinput.Blink
			case key.Matches(keyMsg, keys.Save):
				entries := m.Entries()
				return m, func() tea.Msg { return dirsSavedMsg{entries: entries} }
			case key.Matches(keyMsg, keys.Esc):
				return m, func() tea.Msg { return dirsCancelledMsg{} }
			}
		}
		return m, cmd
	}
}

// View renders the dir editor as a centered modal.
func (m DirEditor) View() string {
	var sb strings.Builder
	sb.WriteString(styleTitle.Render("Edit State Files") + "\n\n")

	for i, e := range m.entries {
		line := shortenPath(e.Dir)
		if len(e.Tags) > 0 {
			line += " " + styleMuted.Render(strings.Join(e.Tags, " "))
		}
		if i == m.cursor && m.phase == phaseBrowse {
			line = styleSelected.Render(line)
		}
		sb.WriteString(line + "\n")
	}
	if len(m.entries) == 0 {
		sb.WriteString(styleMuted.Render("no directories configured") + "\n")
	}

	sb.WriteString("\n")
	switch m.phase {
	case phasePathInput:
		sb.WriteString(m.pathInput.View() + "\n\n")
		mode := "single"
		if m.recursive {
			mode = "recursive"
		}
		sb.WriteString(styleMuted.Render(fmt.Sprintf("enter=confirm  tab=mode:%s  esc=cancel", mode)))
	case phaseScanConfirm:
		sb.WriteString(styleMuted.Render(fmt.Sprintf("Found %d state file(s):", len(m.scannedDirs))) + "\n")
		for _, d := range m.scannedDirs {
			sb.WriteString("  " + shortenPath(d) + "\n")
		}
		if len(m.scannedDirs) == 0 {
			sb.WriteString(styleMuted.Render("  (none found)") + "\n")
		}
		sb.WriteString("\n" + styleMuted.Render("enter=add all  esc=cancel"))
	case phaseTagInput:
		sb.WriteString(styleMuted.Render("Tags for "+shortenPath(m.pendingDir)+":") + "\n")
		sb.WriteString(m.tagInput.View() + "\n\n")
		sb.WriteString(styleMuted.Render("comma-separated  enter=add  esc=skip tags"))
	default:
		sb.WriteString(styleMuted.Render("enter=add  d=delete  ctrl+s=save  esc=cancel"))
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBlue).
		Padding(1, 2).
		Width(m.width).
		Render(sb.String())
}

// scanForStateFiles walks root recursively and returns directories containing a terraform.tfstate file.
func scanForStateFiles(root string) []string {
	var dirs []string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if !d.IsDir() && d.Name() == "terraform.tfstate" {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	return dirs
}

// parseTags splits a comma-separated string into trimmed, non-empty tag strings.
func parseTags(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
