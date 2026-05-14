package tui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// tagInput is a text input with prefix-match autocomplete from a pool of known tags.
// Tab accepts the top suggestion. Enter and Esc are handled by the parent component.
type tagInput struct {
	input       textinput.Model
	pool        []string
	suggestions []string
}

func newTagInput(pool []string) tagInput {
	ti := textinput.New()
	ti.Placeholder = "tag name"
	ti.CharLimit = 64
	return tagInput{input: ti, pool: pool}
}

// Focus activates the cursor. Caller should also return textinput.Blink as the cmd.
func (m *tagInput) Focus() tea.Cmd {
	return m.input.Focus()
}

// Value returns the current input text.
func (m tagInput) Value() string {
	return m.input.Value()
}

// Reset clears the input and suggestions.
func (m *tagInput) Reset() {
	m.input.Reset()
	m.suggestions = nil
}

// Update handles key events. Tab accepts the top suggestion; all other keys
// are forwarded to the underlying textinput and suggestions are refreshed.
func (m tagInput) Update(msg tea.Msg) (tagInput, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, keys.Tab) && len(m.suggestions) > 0 {
			m.input.SetValue(m.suggestions[0])
			m.updateSuggestions()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.updateSuggestions()
	return m, cmd
}

func (m *tagInput) updateSuggestions() {
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

// View renders the input followed by up to 3 suggestion lines.
func (m tagInput) View() string {
	var sb strings.Builder
	sb.WriteString(m.input.View())
	for _, s := range m.suggestions {
		sb.WriteString("\n" + styleMuted.Render("  "+s))
	}
	return sb.String()
}
