package tui

import (
	"testing"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
	tea "github.com/charmbracelet/bubbletea"
)

func makeWorkspaceWithTags(dir string, tags ...string) state.Workspace {
	return state.Workspace{Dir: dir, Tags: tags}
}

func TestTagEditor_DeleteTag(t *testing.T) {
	te := NewTagEditor(makeWorkspaceWithTags("/a", "prod", "staging"), nil)
	te.cursor = 0

	te, _ = te.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	if len(te.tags) != 1 || te.tags[0] != "staging" {
		t.Fatalf("expected [staging] after delete, got %v", te.tags)
	}
}

func TestTagEditor_Save_EmitsTagsUpdatedMsg(t *testing.T) {
	te := NewTagEditor(makeWorkspaceWithTags("/a", "prod"), nil)
	_, cmd := te.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if cmd == nil {
		t.Fatal("expected cmd from ctrl+s")
	}
	msg := cmd()
	upd, ok := msg.(tagsUpdatedMsg)
	if !ok {
		t.Fatalf("expected tagsUpdatedMsg, got %T", msg)
	}
	if upd.dir != "/a" {
		t.Fatalf("unexpected dir: %s", upd.dir)
	}
	if len(upd.tags) != 1 || upd.tags[0] != "prod" {
		t.Fatalf("unexpected tags: %v", upd.tags)
	}
}

func TestTagEditor_Cancel_EmitsTagsCancelledMsg(t *testing.T) {
	te := NewTagEditor(makeWorkspaceWithTags("/a"), nil)
	_, cmd := te.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected cmd from esc")
	}
	if _, ok := cmd().(tagsCancelledMsg); !ok {
		t.Fatal("expected tagsCancelledMsg")
	}
}

func TestTagEditor_AddTagViaInput(t *testing.T) {
	te := NewTagEditor(makeWorkspaceWithTags("/a", "prod"), nil)

	// Enter switches to input phase
	te, _ = te.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if te.phase != tagEditorInput {
		t.Fatalf("expected tagEditorInput, got %v", te.phase)
	}

	// Type new tag
	for _, r := range "aws" {
		te, _ = te.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Enter adds tag and returns to list
	te, _ = te.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if te.phase != tagEditorList {
		t.Fatalf("expected tagEditorList after adding tag, got %v", te.phase)
	}
	if len(te.tags) != 2 || te.tags[1] != "aws" {
		t.Fatalf("expected [prod aws], got %v", te.tags)
	}
}

func TestTagEditor_Navigation(t *testing.T) {
	te := NewTagEditor(makeWorkspaceWithTags("/a", "prod", "staging", "aws"), nil)

	te, _ = te.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if te.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", te.cursor)
	}
	te, _ = te.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if te.cursor != 0 {
		t.Fatalf("expected cursor 0, got %d", te.cursor)
	}
}
