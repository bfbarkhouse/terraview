package tui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FilterBar is the tag-filter row shown between the header and main content.
// Backspace on an empty input removes the last active tag.
// Enter adds the typed tag to the active set. Esc/f in App closes the bar.
type FilterBar struct {
	activeTags  []string
	input       textinput.Model
	pool        []string
	suggestions []string
	width       int
}

// NewFilterBar constructs an empty FilterBar with the given tag pool.
func NewFilterBar(pool []string) FilterBar {
	ti := textinput.New()
	ti.Placeholder = "tag name"
	ti.CharLimit = 64
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#e4e4e7"))
	ti.Focus() // pre-focus so Update accepts keystrokes; Blink cmd re-issued when opened
	return FilterBar{input: ti, pool: pool}
}

// UpdatePool replaces the autocomplete tag pool (called when workspaces change).
func (m *FilterBar) UpdatePool(pool []string) {
	m.pool = pool
}

// ActiveTags returns a copy of the currently active filter tags.
func (m FilterBar) ActiveTags() []string {
	out := make([]string, len(m.activeTags))
	copy(out, m.activeTags)
	return out
}

// Update handles messages for the filter bar.
func (m FilterBar) Update(msg tea.Msg) (FilterBar, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, keys.Tab) && len(m.suggestions) > 0:
			m.input.SetValue(m.suggestions[0])
			m.updateSuggestions()
			return m, nil
		case key.Matches(keyMsg, keys.Enter):
			for _, t := range parseTags(m.input.Value()) {
				if !containsTagCI(m.activeTags, t) {
					m.activeTags = append(m.activeTags, t)
				}
			}
			m.input.Reset()
			m.suggestions = nil
			return m, nil
		case keyMsg.Type == tea.KeyBackspace && m.input.Value() == "" && len(m.activeTags) > 0:
			m.activeTags = m.activeTags[:len(m.activeTags)-1]
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.updateSuggestions()
	return m, cmd
}

func (m *FilterBar) updateSuggestions() {
	val := strings.ToLower(strings.TrimSpace(m.input.Value()))
	if val == "" {
		m.suggestions = nil
		return
	}
	var out []string
	for _, t := range m.pool {
		tl := strings.ToLower(t)
		if strings.HasPrefix(tl, val) && tl != val {
			out = append(out, t)
		}
	}
	sort.Strings(out)
	if len(out) > 3 {
		out = out[:3]
	}
	m.suggestions = out
}

// View renders the filter bar: the input row, active tags, and a hint row.
func (m FilterBar) View() string {
	var sb strings.Builder
	sb.WriteString(styleMuted.Render("Filter: "))
	sb.WriteString(m.input.View())
	for _, t := range m.activeTags {
		sb.WriteString("  " + styleSelected.Render(" "+t+" ×"))
	}
	for _, s := range m.suggestions {
		sb.WriteString("\n" + styleMuted.Render("  tab="+s))
	}
	sb.WriteString("\n" + styleMuted.Render("enter=add tag  backspace=remove last  esc=close"))
	return sb.String()
}

func containsTagCI(tags []string, t string) bool {
	tl := strings.ToLower(t)
	for _, existing := range tags {
		if strings.ToLower(existing) == tl {
			return true
		}
	}
	return false
}
