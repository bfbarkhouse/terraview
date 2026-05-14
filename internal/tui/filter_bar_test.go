package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestFilterBar_AddTag(t *testing.T) {
	fb := NewFilterBar(nil)

	for _, r := range "prod" {
		fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if len(fb.ActiveTags()) != 1 || fb.ActiveTags()[0] != "prod" {
		t.Fatalf("expected [prod], got %v", fb.ActiveTags())
	}
}

func TestFilterBar_BackspaceRemovesLastTag(t *testing.T) {
	fb := NewFilterBar(nil)

	// Add two tags
	for _, r := range "prod" {
		fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyEnter})
	for _, r := range "staging" {
		fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if len(fb.ActiveTags()) != 2 {
		t.Fatalf("expected 2 tags before backspace, got %v", fb.ActiveTags())
	}

	// Backspace with empty input removes last tag
	fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if len(fb.ActiveTags()) != 1 || fb.ActiveTags()[0] != "prod" {
		t.Fatalf("expected [prod] after backspace, got %v", fb.ActiveTags())
	}
}

func TestFilterBar_NoDuplicateTags(t *testing.T) {
	fb := NewFilterBar(nil)

	for _, r := range "prod" {
		fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyEnter})
	for _, r := range "prod" {
		fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	fb, _ = fb.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if len(fb.ActiveTags()) != 1 {
		t.Fatalf("expected no duplicate tags, got %v", fb.ActiveTags())
	}
}

func TestHasAllTags_ANDSemantics(t *testing.T) {
	tests := []struct {
		wsTags     []string
		filterTags []string
		want       bool
	}{
		{[]string{"prod", "aws"}, []string{"prod"}, true},
		{[]string{"prod", "aws"}, []string{"prod", "aws"}, true},
		{[]string{"prod"}, []string{"prod", "aws"}, false},
		{[]string{"prod", "aws"}, []string{"staging"}, false},
		{[]string{"prod"}, nil, true},
		{nil, []string{"prod"}, false},
	}
	for _, tt := range tests {
		got := hasAllTags(tt.wsTags, tt.filterTags)
		if got != tt.want {
			t.Errorf("hasAllTags(%v, %v) = %v, want %v", tt.wsTags, tt.filterTags, got, tt.want)
		}
	}
}
