package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestTagInput_TabCompletes(t *testing.T) {
	ti := newTagInput([]string{"prod", "production", "staging"})
	ti.Focus()

	// type "pro"
	for _, r := range "pro" {
		ti, _ = ti.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	if len(ti.suggestions) == 0 {
		t.Fatal("expected suggestions after typing 'pro'")
	}

	// tab should complete to first suggestion
	ti, _ = ti.Update(tea.KeyMsg{Type: tea.KeyTab})
	if ti.Value() != ti.suggestions[0] && ti.Value() != "prod" && ti.Value() != "production" {
		t.Fatalf("unexpected value after tab: %q", ti.Value())
	}
}

func TestTagInput_NoSuggestionsForExactMatch(t *testing.T) {
	ti := newTagInput([]string{"prod"})
	ti.Focus()
	for _, r := range "prod" {
		ti, _ = ti.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	// "prod" exactly matches the only tag — no suggestion (would be a no-op autocomplete)
	if len(ti.suggestions) != 0 {
		t.Fatalf("expected no suggestions for exact match, got %v", ti.suggestions)
	}
}

func TestTagInput_Reset(t *testing.T) {
	ti := newTagInput([]string{"prod"})
	ti.Focus()
	for _, r := range "pro" {
		ti, _ = ti.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	ti.Reset()
	if ti.Value() != "" {
		t.Fatalf("expected empty after Reset, got %q", ti.Value())
	}
	if len(ti.suggestions) != 0 {
		t.Fatalf("expected no suggestions after Reset")
	}
}

func TestTagInput_EmptyInputNoSuggestions(t *testing.T) {
	ti := newTagInput([]string{"prod", "staging"})
	if len(ti.suggestions) != 0 {
		t.Fatal("new tagInput should have no suggestions")
	}
}
