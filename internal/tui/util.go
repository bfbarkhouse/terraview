package tui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bfbarkhouse-redpanda/terraview/internal/state"
)

// shortenPath replaces the home directory prefix with ~.
func shortenPath(p string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	if strings.HasPrefix(p, home) {
		return "~" + p[len(home):]
	}
	return p
}

// expandPath replaces a leading ~ with the user's home directory.
func expandPath(p string) string {
	if p != "~" && !strings.HasPrefix(p, "~/") {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	return home + p[1:]
}

// pluralize returns "1 word" or "N words".
func pluralize(n int, word string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// allTagsFromWorkspaces returns a sorted, deduplicated list of all tags
// across the given workspaces.
func allTagsFromWorkspaces(ws []state.Workspace) []string {
	seen := make(map[string]string) // lowercase key → original value (last-write-wins)
	for _, w := range ws {
		for _, t := range w.Tags {
			seen[strings.ToLower(t)] = t
		}
	}
	out := make([]string, 0, len(seen))
	for _, t := range seen {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// hasAllTags reports whether wsTags is a superset of filterTags (AND semantics).
// Comparison is case-insensitive.
func hasAllTags(wsTags, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	set := make(map[string]struct{}, len(wsTags))
	for _, t := range wsTags {
		set[strings.ToLower(t)] = struct{}{}
	}
	for _, f := range filterTags {
		if _, ok := set[strings.ToLower(f)]; !ok {
			return false
		}
	}
	return true
}
