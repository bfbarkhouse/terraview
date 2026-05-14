package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bfbarkhouse-redpanda/terraview/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func makeEntries(dirs ...string) []config.WorkspaceEntry {
	entries := make([]config.WorkspaceEntry, len(dirs))
	for i, d := range dirs {
		entries[i] = config.WorkspaceEntry{Dir: d}
	}
	return entries
}

func TestDirEditor_DeleteMiddleItem(t *testing.T) {
	de := NewDirEditor(makeEntries("/a", "/b", "/c"), nil)
	de.cursor = 1

	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	if len(de.entries) != 2 {
		t.Fatalf("expected 2 entries after delete, got %d", len(de.entries))
	}
	if de.entries[1].Dir != "/c" {
		t.Fatalf("expected /c at index 1, got %s", de.entries[1].Dir)
	}
}

func TestDirEditor_Delete_ClampsCursor(t *testing.T) {
	de := NewDirEditor(makeEntries("/a", "/b"), nil)
	de.cursor = 1

	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	if de.cursor != 0 {
		t.Fatalf("cursor should clamp to 0, got %d", de.cursor)
	}
}

func TestDirEditor_Entries_ReturnsCopy(t *testing.T) {
	de := NewDirEditor(makeEntries("/a", "/b"), nil)
	entries := de.Entries()
	entries[0].Dir = "/mutated"
	if de.entries[0].Dir == "/mutated" {
		t.Fatal("Entries() should return a copy, not the internal slice")
	}
}

func TestDirEditor_Save_EmitsDirsSavedMsg(t *testing.T) {
	de := NewDirEditor(makeEntries("/a", "/b"), nil)
	_, cmd := de.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	if cmd == nil {
		t.Fatal("expected a tea.Cmd from ctrl+s, got nil")
	}
	msg := cmd()
	saved, ok := msg.(dirsSavedMsg)
	if !ok {
		t.Fatalf("expected dirsSavedMsg, got %T", msg)
	}
	if len(saved.entries) != 2 || saved.entries[0].Dir != "/a" {
		t.Fatalf("unexpected entries in message: %v", saved.entries)
	}
}

func TestDirEditor_Navigation(t *testing.T) {
	de := NewDirEditor(makeEntries("/a", "/b", "/c"), nil)

	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if de.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", de.cursor)
	}

	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	if de.cursor != 0 {
		t.Fatalf("expected cursor 0, got %d", de.cursor)
	}
}

func TestDirEditor_TwoPhaseAdd_WithTags(t *testing.T) {
	de := NewDirEditor(nil, nil)

	// Enter browse → path input phase
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if de.phase != phasePathInput {
		t.Fatalf("expected phasePathInput, got %v", de.phase)
	}

	// Type a path
	for _, r := range "/mydir" {
		de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Confirm path → tag input phase
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if de.phase != phaseTagInput {
		t.Fatalf("expected phaseTagInput after confirming path, got %v", de.phase)
	}
	if de.pendingDir != "/mydir" {
		t.Fatalf("expected pendingDir=/mydir, got %q", de.pendingDir)
	}

	// Type tags
	for _, r := range "prod,aws" {
		de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Confirm tags → entry added, back to browse
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if de.phase != phaseBrowse {
		t.Fatalf("expected phaseBrowse after confirming tags, got %v", de.phase)
	}
	if len(de.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(de.entries))
	}
	if de.entries[0].Dir != "/mydir" {
		t.Fatalf("unexpected dir: %s", de.entries[0].Dir)
	}
	if len(de.entries[0].Tags) != 2 || de.entries[0].Tags[0] != "prod" {
		t.Fatalf("unexpected tags: %v", de.entries[0].Tags)
	}
}

func TestDirEditor_TwoPhaseAdd_EscSkipsTags(t *testing.T) {
	de := NewDirEditor(nil, nil)

	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter})
	for _, r := range "/mydir" {
		de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter}) // confirm path

	// Esc in tag phase → add with no tags
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if de.phase != phaseBrowse {
		t.Fatalf("expected phaseBrowse after esc in tag phase, got %v", de.phase)
	}
	if len(de.entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(de.entries))
	}
	if len(de.entries[0].Tags) != 0 {
		t.Fatalf("expected no tags after esc, got %v", de.entries[0].Tags)
	}
}

func TestDirEditor_RecursiveScan_FindsAndAdds(t *testing.T) {
	// Set up a temp directory tree with two tfstate files
	root := t.TempDir()
	dir1 := filepath.Join(root, "project-a")
	dir2 := filepath.Join(root, "nested", "project-b")
	if err := os.MkdirAll(dir1, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir2, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "terraform.tfstate"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir2, "terraform.tfstate"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}

	de := NewDirEditor(nil, nil)

	// Enter path input phase
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Toggle recursive mode
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyTab})
	if !de.recursive {
		t.Fatal("expected recursive=true after Tab")
	}

	// Type the root path
	for _, r := range root {
		de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	// Confirm — should scan and go to phaseScanConfirm
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if de.phase != phaseScanConfirm {
		t.Fatalf("expected phaseScanConfirm, got %v", de.phase)
	}
	if len(de.scannedDirs) != 2 {
		t.Fatalf("expected 2 scanned dirs, got %d: %v", len(de.scannedDirs), de.scannedDirs)
	}

	// Confirm add all
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if de.phase != phaseBrowse {
		t.Fatalf("expected phaseBrowse after confirming, got %v", de.phase)
	}
	if len(de.entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(de.entries))
	}
}

func TestDirEditor_RecursiveScan_EscCancels(t *testing.T) {
	root := t.TempDir()
	dir1 := filepath.Join(root, "proj")
	if err := os.MkdirAll(dir1, 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(dir1, "terraform.tfstate"), []byte(`{}`), 0o644)

	de := NewDirEditor(nil, nil)
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter})
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyTab}) // toggle recursive
	for _, r := range root {
		de, _ = de.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEnter}) // scan

	// Esc cancels
	de, _ = de.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if de.phase != phaseBrowse {
		t.Fatalf("expected phaseBrowse after esc, got %v", de.phase)
	}
	if len(de.entries) != 0 {
		t.Fatalf("expected no entries after cancel, got %d", len(de.entries))
	}
}
