package tui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bfbarkhouse-redpanda/terraview/internal/config"
	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewMode int

const (
	viewSplit viewMode = iota
	viewSummary
)

type focusedPanel int

const (
	focusWorkspaces focusedPanel = iota
	focusResources
)

// App is the root Bubble Tea model.
type App struct {
	config        config.Config
	workspaces    []state.Workspace
	workspaceList WorkspaceList
	resourceList  ResourceList
	summaryView   SummaryView
	dirEditor     DirEditor
	showDirEditor bool
	tagEditor     TagEditor
	showTagEditor bool
	filterBar     FilterBar
	showFilterBar bool
	viewMode      viewMode
	focus         focusedPanel
	width         int
	height        int
}

// New constructs the App model from a config, loading all workspaces from disk.
func New(cfg config.Config) App {
	workspaces := loadAll(cfg.Workspaces)
	pool := allTagsFromWorkspaces(workspaces)
	wl := NewWorkspaceList(workspaces)
	wl.focused = true
	rl := NewResourceList()
	if sel := wl.Selected(); sel != nil {
		rl.SetWorkspace(sel)
	}
	return App{
		config:        cfg,
		workspaces:    workspaces,
		workspaceList: wl,
		resourceList:  rl,
		summaryView:   NewSummaryView(workspaces),
		dirEditor:     NewDirEditor(cfg.Workspaces, pool),
		filterBar:     NewFilterBar(pool),
		focus:         focusWorkspaces,
	}
}

func loadAll(entries []config.WorkspaceEntry) []state.Workspace {
	ws := make([]state.Workspace, len(entries))
	for i, e := range entries {
		ws[i] = state.Load(e.Dir, e.Tags)
	}
	return ws
}

func (m *App) syncWorkspaces(entries []config.WorkspaceEntry) {
	m.config.Workspaces = entries
	m.workspaces = loadAll(entries)
	filtered := m.filteredWorkspaces()
	m.workspaceList.SetWorkspaces(filtered)
	m.summaryView.SetWorkspaces(filtered)
	m.filterBar.UpdatePool(allTagsFromWorkspaces(m.workspaces))
	if sel := m.workspaceList.Selected(); sel != nil {
		m.resourceList.SetWorkspace(sel)
		m.resourceList.cursor = 0
	} else {
		m.resourceList.SetWorkspace(nil)
	}
}

func (m App) filteredWorkspaces() []state.Workspace {
	if len(m.filterBar.activeTags) == 0 {
		return m.workspaces
	}
	var out []state.Workspace
	for _, ws := range m.workspaces {
		if hasAllTags(ws.Tags, m.filterBar.activeTags) {
			out = append(out, ws)
		}
	}
	return out
}

// Init implements tea.Model.
func (m App) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.setSizes()
		return m, nil

	case dirsSavedMsg:
		if err := config.Save(config.Config{Workspaces: msg.entries}); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not save config: %v\n", err)
		}
		m.syncWorkspaces(msg.entries)
		m.showDirEditor = false
		return m, nil

	case dirsCancelledMsg:
		m.showDirEditor = false
		return m, nil

	case tagsUpdatedMsg:
		for i, e := range m.config.Workspaces {
			if e.Dir == msg.dir {
				m.config.Workspaces[i].Tags = msg.tags
				break
			}
		}
		if err := config.Save(m.config); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not save config: %v\n", err)
		}
		m.workspaces = loadAll(m.config.Workspaces)
		filtered := m.filteredWorkspaces()
		m.workspaceList.SetWorkspaces(filtered)
		m.summaryView.SetWorkspaces(filtered)
		m.filterBar.UpdatePool(allTagsFromWorkspaces(m.workspaces))
		if sel := m.workspaceList.Selected(); sel != nil {
			m.resourceList.SetWorkspace(sel)
		}
		m.showTagEditor = false
		return m, nil

	case tagsCancelledMsg:
		m.showTagEditor = false
		return m, nil

	case tea.KeyMsg:
		if m.showDirEditor {
			m.dirEditor, cmd = m.dirEditor.Update(msg)
			return m, cmd
		}

		if m.showTagEditor {
			m.tagEditor, cmd = m.tagEditor.Update(msg)
			return m, cmd
		}

		if m.showFilterBar {
			if key.Matches(msg, keys.Esc) {
				m.showFilterBar = false
				m.setSizes()
				return m, nil
			}
			prevLen := len(m.filterBar.activeTags)
			m.filterBar, cmd = m.filterBar.Update(msg)
			if len(m.filterBar.activeTags) != prevLen {
				filtered := m.filteredWorkspaces()
				m.workspaceList.SetWorkspaces(filtered)
				m.summaryView.SetWorkspaces(filtered)
			}
			return m, cmd
		}

		if key.Matches(msg, keys.Quit) {
			return m, tea.Quit
		}

		switch {
		case key.Matches(msg, keys.Edit):
			m.dirEditor = NewDirEditor(m.config.Workspaces, allTagsFromWorkspaces(m.workspaces))
			m.dirEditor.width = m.width / 2
			m.showDirEditor = true
			return m, nil

		case key.Matches(msg, keys.Tag):
			var ws *state.Workspace
			if m.viewMode == viewSummary {
				filtered := m.filteredWorkspaces()
				idx := m.summaryView.SelectedIndex()
				if idx >= 0 && idx < len(filtered) {
					ws = &filtered[idx]
				}
			} else if m.focus == focusWorkspaces {
				ws = m.workspaceList.Selected()
			}
			if ws != nil {
				te := NewTagEditor(*ws, allTagsFromWorkspaces(m.workspaces))
				te.width = m.width / 2
				m.tagEditor = te
				m.showTagEditor = true
			}
			return m, nil

		case key.Matches(msg, keys.Filter):
			m.filterBar.UpdatePool(allTagsFromWorkspaces(m.workspaces))
			cmd = m.filterBar.input.Focus()
			m.showFilterBar = true
			m.setSizes()
			return m, tea.Batch(cmd, textinput.Blink)

		case key.Matches(msg, keys.Summary):
			if m.viewMode == viewSplit {
				m.viewMode = viewSummary
				m.summaryView.cursor = m.workspaceList.SelectedIndex()
			} else {
				m.viewMode = viewSplit
			}
			return m, nil

		case key.Matches(msg, keys.Enter) && m.viewMode == viewSummary:
			m.workspaceList.cursor = m.summaryView.SelectedIndex()
			if sel := m.workspaceList.Selected(); sel != nil {
				m.resourceList.SetWorkspace(sel)
				m.resourceList.cursor = 0
			}
			m.viewMode = viewSplit
			return m, nil

		case key.Matches(msg, keys.Open):
			var dir string
			if m.viewMode == viewSummary {
				filtered := m.filteredWorkspaces()
				idx := m.summaryView.SelectedIndex()
				if idx >= 0 && idx < len(filtered) {
					dir = filtered[idx].Dir
				}
			} else {
				if sel := m.workspaceList.Selected(); sel != nil {
					dir = sel.Dir
				}
			}
			if dir != "" {
				shell := os.Getenv("SHELL")
				if shell == "" {
					shell = "/bin/sh"
				}
				c := exec.Command(shell)
				c.Dir = dir
				return m, tea.ExecProcess(c, nil)
			}
			return m, nil

		case key.Matches(msg, keys.Refresh) && m.viewMode == viewSplit:
			sel := m.workspaceList.Selected()
			if sel != nil {
				for i, ws := range m.workspaces {
					if ws.Dir == sel.Dir {
						m.workspaces[i] = state.Load(ws.Dir, ws.Tags)
						break
					}
				}
				filtered := m.filteredWorkspaces()
				m.workspaceList.SetWorkspaces(filtered)
				m.summaryView.SetWorkspaces(filtered)
				if sel2 := m.workspaceList.Selected(); sel2 != nil {
					m.resourceList.SetWorkspace(sel2)
				}
			}
			return m, nil

		case key.Matches(msg, keys.Tab) && m.viewMode == viewSplit:
			if m.focus == focusWorkspaces {
				m.focus = focusResources
				m.workspaceList.focused = false
				m.resourceList.focused = true
			} else {
				m.focus = focusWorkspaces
				m.workspaceList.focused = true
				m.resourceList.focused = false
			}
			return m, nil
		}

		if m.viewMode == viewSummary {
			m.summaryView, cmd = m.summaryView.Update(msg)
		} else {
			prevIdx := m.workspaceList.SelectedIndex()
			m.workspaceList, _ = m.workspaceList.Update(msg)
			m.resourceList, cmd = m.resourceList.Update(msg)
			if m.workspaceList.SelectedIndex() != prevIdx {
				if sel := m.workspaceList.Selected(); sel != nil {
					m.resourceList.SetWorkspace(sel)
					m.resourceList.cursor = 0
				}
			}
		}

	default:
		if m.showFilterBar {
			m.filterBar, cmd = m.filterBar.Update(msg)
			return m, cmd
		}
	}

	return m, cmd
}

// View implements tea.Model.
func (m App) View() string {
	if m.width == 0 {
		return ""
	}

	if m.showDirEditor {
		overlay := m.dirEditor.View()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, overlay)
	}

	if m.showTagEditor {
		overlay := m.tagEditor.View()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, overlay)
	}

	filterHint := "f=filter"
	if n := len(m.filterBar.activeTags); n > 0 && !m.showFilterBar {
		filterHint = fmt.Sprintf("f=filter(%d)", n)
	}
	header := styleTitle.Render("terraview") + "  " +
		styleMuted.Render(fmt.Sprintf("tab=focus  s=summary  e=edit  t=tags  %s  r=refresh  o=shell  q=quit", filterHint))

	var body string
	if m.viewMode == viewSummary {
		body = m.summaryView.View()
	} else {
		body = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.workspaceList.View(),
			m.resourceList.View(),
		)
	}

	parts := []string{header}
	if m.showFilterBar {
		parts = append(parts, m.filterBar.View())
	}
	parts = append(parts, body)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m *App) setSizes() {
	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth - 4
	panelHeight := m.height - 3
	if m.showFilterBar {
		panelHeight -= 2 // filter bar occupies 2 extra lines (input + hint)
	}

	m.workspaceList.width = leftWidth
	m.workspaceList.height = panelHeight
	m.resourceList.width = rightWidth
	m.resourceList.height = panelHeight
	m.summaryView.width = m.width
	m.summaryView.height = panelHeight
	m.dirEditor.width = m.width / 2
	m.tagEditor.width = m.width / 2
	m.filterBar.width = m.width
}
