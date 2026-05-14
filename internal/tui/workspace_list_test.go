package tui

import (
	"testing"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	tea "github.com/charmbracelet/bubbletea"
)

func makeWorkspaces(dirs ...string) []state.Workspace {
	ws := make([]state.Workspace, len(dirs))
	for i, d := range dirs {
		ws[i] = state.Workspace{Dir: d, Resources: []state.Resource{{Type: "aws_instance", Name: "x"}}}
	}
	return ws
}

func TestWorkspaceList_NavigateDown(t *testing.T) {
	wl := NewWorkspaceList(makeWorkspaces("/a", "/b", "/c"))
	wl.focused = true

	wl, _ = wl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if wl.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", wl.cursor)
	}
}

func TestWorkspaceList_NavigateUp(t *testing.T) {
	wl := NewWorkspaceList(makeWorkspaces("/a", "/b", "/c"))
	wl.focused = true
	wl.cursor = 2

	wl, _ = wl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if wl.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", wl.cursor)
	}
}

func TestWorkspaceList_ClampAtEnd(t *testing.T) {
	wl := NewWorkspaceList(makeWorkspaces("/a", "/b"))
	wl.focused = true
	wl.cursor = 1

	wl, _ = wl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if wl.cursor != 1 {
		t.Fatalf("cursor should stay at 1, got %d", wl.cursor)
	}
}

func TestWorkspaceList_NotFocused_NoNavigation(t *testing.T) {
	wl := NewWorkspaceList(makeWorkspaces("/a", "/b"))
	wl.focused = false

	wl, _ = wl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if wl.cursor != 0 {
		t.Fatalf("unfocused list should not navigate, got cursor %d", wl.cursor)
	}
}

func TestWorkspaceList_Selected(t *testing.T) {
	wl := NewWorkspaceList(makeWorkspaces("/a", "/b"))
	wl.focused = true

	sel := wl.Selected()
	if sel == nil || sel.Dir != "/a" {
		t.Fatalf("expected /a, got %v", sel)
	}

	wl, _ = wl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	sel = wl.Selected()
	if sel == nil || sel.Dir != "/b" {
		t.Fatalf("expected /b after nav, got %v", sel)
	}
}

func TestWorkspaceList_Empty_SelectedNil(t *testing.T) {
	wl := NewWorkspaceList(nil)
	if wl.Selected() != nil {
		t.Fatal("expected nil Selected() for empty list")
	}
}

func TestWorkspaceList_SetWorkspaces_ClampsCursor(t *testing.T) {
	wl := NewWorkspaceList(makeWorkspaces("/a", "/b", "/c"))
	wl.cursor = 2
	wl.SetWorkspaces(makeWorkspaces("/a"))
	if wl.cursor != 0 {
		t.Fatalf("cursor should clamp to 0, got %d", wl.cursor)
	}
}
