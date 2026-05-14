package tui

import (
	"testing"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	tea "github.com/charmbracelet/bubbletea"
)

func TestResourceList_SetWorkspace(t *testing.T) {
	rl := NewResourceList()
	ws := state.Workspace{
		Dir: "/foo",
		Resources: []state.Resource{
			{Type: "aws_instance", Name: "worker"},
			{Type: "aws_vpc", Name: "main"},
		},
	}
	rl.SetWorkspace(&ws)
	if rl.workspace == nil {
		t.Fatal("workspace should be set")
	}
	if len(rl.workspace.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(rl.workspace.Resources))
	}
}

func TestResourceList_Navigation(t *testing.T) {
	rl := NewResourceList()
	ws := state.Workspace{
		Dir: "/foo",
		Resources: []state.Resource{
			{Type: "aws_instance", Name: "worker"},
			{Type: "aws_vpc", Name: "main"},
		},
	}
	rl.SetWorkspace(&ws)
	rl.focused = true

	rl, _ = rl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if rl.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", rl.cursor)
	}

	// Can't go past end
	rl, _ = rl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if rl.cursor != 1 {
		t.Fatalf("cursor should stay at 1, got %d", rl.cursor)
	}
}

func TestResourceList_SetWorkspace_ResetsCursor(t *testing.T) {
	rl := NewResourceList()
	ws := state.Workspace{Dir: "/foo", Resources: []state.Resource{
		{Type: "aws_instance", Name: "a"},
		{Type: "aws_vpc", Name: "b"},
	}}
	rl.SetWorkspace(&ws)
	rl.focused = true
	rl, _ = rl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if rl.cursor != 1 {
		t.Fatalf("cursor should be 1, got %d", rl.cursor)
	}

	ws2 := state.Workspace{Dir: "/bar", Resources: []state.Resource{{Type: "aws_s3_bucket", Name: "x"}}}
	rl.SetWorkspace(&ws2)
	if rl.cursor != 0 {
		t.Fatalf("SetWorkspace should reset cursor to 0, got %d", rl.cursor)
	}
}

func TestResourceList_NilWorkspace_ViewNoPanic(t *testing.T) {
	rl := NewResourceList()
	rl.SetWorkspace(nil)
	// Should not panic
	_ = rl.View()
}
