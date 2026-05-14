package tui

import "github.com/bfbarkhouse-redpanda/terraview/internal/config"

// dirsSavedMsg is emitted by DirEditor when the user saves (ctrl+s).
type dirsSavedMsg struct {
	entries []config.WorkspaceEntry
}

// dirsCancelledMsg is emitted by DirEditor when the user cancels (esc).
type dirsCancelledMsg struct{}

// tagsUpdatedMsg is emitted by TagEditor when the user saves (ctrl+s).
type tagsUpdatedMsg struct {
	dir  string
	tags []string
}

// tagsCancelledMsg is emitted by TagEditor when the user cancels (esc).
type tagsCancelledMsg struct{}
