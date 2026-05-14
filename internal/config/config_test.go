package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad_WithTags(t *testing.T) {
	configDirOverride = t.TempDir()
	t.Cleanup(func() { configDirOverride = "" })

	cfg := Config{Workspaces: []WorkspaceEntry{
		{Dir: "/foo/bar", Tags: []string{"prod", "aws"}},
		{Dir: "/baz/qux"},
	}}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Workspaces) != 2 {
		t.Fatalf("expected 2 workspaces, got %d", len(loaded.Workspaces))
	}
	if loaded.Workspaces[0].Dir != "/foo/bar" {
		t.Fatalf("unexpected dir: %s", loaded.Workspaces[0].Dir)
	}
	if len(loaded.Workspaces[0].Tags) != 2 || loaded.Workspaces[0].Tags[0] != "prod" {
		t.Fatalf("unexpected tags: %v", loaded.Workspaces[0].Tags)
	}
	if len(loaded.Workspaces[1].Tags) != 0 {
		t.Fatalf("expected no tags, got %v", loaded.Workspaces[1].Tags)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	configDirOverride = t.TempDir()
	t.Cleanup(func() { configDirOverride = "" })

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load on missing file should not error, got: %v", err)
	}
	if len(cfg.Workspaces) != 0 {
		t.Fatalf("expected empty config, got %v", cfg.Workspaces)
	}
}

func TestLoad_LegacyMigration(t *testing.T) {
	configDirOverride = t.TempDir()
	t.Cleanup(func() { configDirOverride = "" })

	// Write old-format config with "directories" key
	path := filepath.Join(configDirOverride, "config.json")
	legacy := `{"directories":["/alpha","/beta"]}`
	if err := os.WriteFile(path, []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Workspaces) != 2 {
		t.Fatalf("expected 2 migrated workspaces, got %d", len(cfg.Workspaces))
	}
	if cfg.Workspaces[0].Dir != "/alpha" {
		t.Fatalf("unexpected dir: %s", cfg.Workspaces[0].Dir)
	}
	if len(cfg.Workspaces[0].Tags) != 0 {
		t.Fatalf("expected no tags after migration, got %v", cfg.Workspaces[0].Tags)
	}
}

func TestSave_ProducesNewFormat(t *testing.T) {
	configDirOverride = t.TempDir()
	t.Cleanup(func() { configDirOverride = "" })

	cfg := Config{Workspaces: []WorkspaceEntry{{Dir: "/x", Tags: []string{"staging"}}}}
	if err := Save(cfg); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(configDirOverride, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["workspaces"]; !ok {
		t.Fatal("saved config should have 'workspaces' key, not 'directories'")
	}
}
