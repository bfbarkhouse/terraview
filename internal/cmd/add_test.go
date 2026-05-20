package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bfbarkhouse-redpanda/terraview/internal/cmd"
	"github.com/bfbarkhouse-redpanda/terraview/internal/config"
)

func TestResolvePath_Dot(t *testing.T) {
	got, err := cmd.ResolvePath(".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cwd, _ := os.Getwd()
	if got != cwd {
		t.Fatalf("expected %s, got %s", cwd, got)
	}
}

func TestResolvePath_Tilde(t *testing.T) {
	home, _ := os.UserHomeDir()
	got, err := cmd.ResolvePath("~/foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(home, "foo", "bar")
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestResolvePath_Relative(t *testing.T) {
	got, err := cmd.ResolvePath("foo/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("expected absolute path, got %s", got)
	}
}

func TestResolvePath_Absolute(t *testing.T) {
	got, err := cmd.ResolvePath("/tmp/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/tmp/foo" {
		t.Fatalf("expected /tmp/foo, got %s", got)
	}
}

func TestResolvePath_BareTilde(t *testing.T) {
	home, _ := os.UserHomeDir()
	got, err := cmd.ResolvePath("~")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != home {
		t.Fatalf("expected %s, got %s", home, got)
	}
}

func TestRunAdd_AddsWorkspace(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDir(t.TempDir())
	t.Cleanup(func() { config.SetConfigDir("") })

	err := cmd.RunAdd([]string{dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, _ := config.Load()
	if len(cfg.Workspaces) != 1 {
		t.Fatalf("expected 1 workspace, got %d", len(cfg.Workspaces))
	}
	if cfg.Workspaces[0].Dir != dir {
		t.Fatalf("expected dir %s, got %s", dir, cfg.Workspaces[0].Dir)
	}
	if len(cfg.Workspaces[0].Tags) != 0 {
		t.Fatalf("expected no tags, got %v", cfg.Workspaces[0].Tags)
	}
}

func TestRunAdd_WithTags(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDir(t.TempDir())
	t.Cleanup(func() { config.SetConfigDir("") })

	err := cmd.RunAdd([]string{"--tags", "prod,aws", dir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, _ := config.Load()
	if len(cfg.Workspaces[0].Tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", cfg.Workspaces[0].Tags)
	}
	if cfg.Workspaces[0].Tags[0] != "prod" || cfg.Workspaces[0].Tags[1] != "aws" {
		t.Fatalf("unexpected tags: %v", cfg.Workspaces[0].Tags)
	}
}

func TestRunAdd_Duplicate(t *testing.T) {
	dir := t.TempDir()
	config.SetConfigDir(t.TempDir())
	t.Cleanup(func() { config.SetConfigDir("") })

	// Add once
	if err := cmd.RunAdd([]string{dir}); err != nil {
		t.Fatalf("first add: %v", err)
	}
	// Add again — should succeed (exit 0) but not duplicate
	if err := cmd.RunAdd([]string{dir}); err != nil {
		t.Fatalf("duplicate add should not return error: %v", err)
	}

	cfg, _ := config.Load()
	if len(cfg.Workspaces) != 1 {
		t.Fatalf("expected 1 workspace after duplicate add, got %d", len(cfg.Workspaces))
	}
}

func TestRunAdd_NonExistentDir(t *testing.T) {
	config.SetConfigDir(t.TempDir())
	t.Cleanup(func() { config.SetConfigDir("") })

	err := cmd.RunAdd([]string{"/nonexistent/path/that/does/not/exist"})
	if err == nil {
		t.Fatal("expected error for non-existent directory, got nil")
	}
}

func TestRunAdd_NoArgs(t *testing.T) {
	err := cmd.RunAdd([]string{})
	if err == nil {
		t.Fatal("expected error for missing argument, got nil")
	}
}

func TestRunAdd_NotADirectory(t *testing.T) {
	f, err := os.CreateTemp("", "terraview-test-*")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	config.SetConfigDir(t.TempDir())
	t.Cleanup(func() { config.SetConfigDir("") })

	err = cmd.RunAdd([]string{f.Name()})
	if err == nil {
		t.Fatal("expected error for non-directory path, got nil")
	}
}
