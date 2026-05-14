package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Resource represents a single Terraform-managed resource from a state file.
type Resource struct {
	Type string
	Name string
}

// DisplayName returns the canonical "type.name" string shown in the UI.
func (r Resource) DisplayName() string {
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}

// Workspace holds the parsed state for one Terraform directory.
type Workspace struct {
	Dir       string
	Tags      []string
	Resources []Resource
	Err       error // non-nil if state file is missing or unreadable
}

// HasResources reports whether the workspace has any deployed resources.
func (w Workspace) HasResources() bool {
	return len(w.Resources) > 0
}

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

// Load reads and parses the terraform.tfstate file in dir, attaching tags.
func Load(dir string, tags []string) Workspace {
	ws := Workspace{Dir: dir, Tags: tags}
	path := filepath.Join(expandPath(dir), "terraform.tfstate")

	data, err := os.ReadFile(path)
	if err != nil {
		ws.Err = err
		return ws
	}

	var raw struct {
		Resources []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"resources"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		ws.Err = err
		return ws
	}

	for _, r := range raw.Resources {
		ws.Resources = append(ws.Resources, Resource{Type: r.Type, Name: r.Name})
	}
	return ws
}
