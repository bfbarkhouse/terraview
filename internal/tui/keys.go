package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Tab     key.Binding
	Enter   key.Binding
	Refresh key.Binding
	Summary key.Binding
	Edit    key.Binding
	Delete  key.Binding
	Save    key.Binding
	Tag     key.Binding
	Filter  key.Binding
	Open    key.Binding
	Help    key.Binding
	Quit    key.Binding
	Esc     key.Binding
}

var keys = keyMap{
	Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Left:    key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
	Right:   key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
	Tab:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch panel")),
	Enter:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select/add")),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
	Summary: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "summary")),
	Edit:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit dirs")),
	Delete:  key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
	Save:    key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
	Tag:     key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "edit tags")),
	Filter:  key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "filter")),
	Open:    key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open shell")),
	Help:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Esc:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel/back")),
}
