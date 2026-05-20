package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bfbarkhouse-redpanda/terraview/internal/cmd"
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
