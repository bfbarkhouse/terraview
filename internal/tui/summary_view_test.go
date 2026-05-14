package tui

import (
	"testing"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	tea "github.com/charmbracelet/bubbletea"
)

func TestSummaryView_NavigateRight(t *testing.T) {
	sv := NewSummaryView([]state.Workspace{{Dir: "/a"}, {Dir: "/b"}, {Dir: "/c"}})

	sv, _ = sv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if sv.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", sv.cursor)
	}
}

func TestSummaryView_NavigateLeft(t *testing.T) {
	sv := NewSummaryView([]state.Workspace{{Dir: "/a"}, {Dir: "/b"}})
	sv.cursor = 1

	sv, _ = sv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	if sv.cursor != 0 {
		t.Fatalf("expected cursor 0, got %d", sv.cursor)
	}
}

func TestSummaryView_ClampAtStart(t *testing.T) {
	sv := NewSummaryView([]state.Workspace{{Dir: "/a"}, {Dir: "/b"}})

	sv, _ = sv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	if sv.cursor != 0 {
		t.Fatalf("cursor should stay at 0, got %d", sv.cursor)
	}
}

func TestSummaryView_ClampAtEnd(t *testing.T) {
	sv := NewSummaryView([]state.Workspace{{Dir: "/a"}, {Dir: "/b"}})
	sv.cursor = 1

	sv, _ = sv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if sv.cursor != 1 {
		t.Fatalf("cursor should stay at 1, got %d", sv.cursor)
	}
}

func TestSummaryView_SelectedIndex(t *testing.T) {
	sv := NewSummaryView([]state.Workspace{{Dir: "/a"}, {Dir: "/b"}, {Dir: "/c"}})
	sv.cursor = 2
	if sv.SelectedIndex() != 2 {
		t.Fatalf("expected 2, got %d", sv.SelectedIndex())
	}
}

func TestSummaryView_SetWorkspaces_ClampsCursor(t *testing.T) {
	sv := NewSummaryView([]state.Workspace{{Dir: "/a"}, {Dir: "/b"}, {Dir: "/c"}})
	sv.cursor = 2
	sv.SetWorkspaces([]state.Workspace{{Dir: "/a"}})
	if sv.cursor != 0 {
		t.Fatalf("cursor should clamp to 0, got %d", sv.cursor)
	}
}
