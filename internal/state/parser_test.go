package state

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTFState(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "terraform.tfstate"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoad_WithResources(t *testing.T) {
	dir := t.TempDir()
	writeTFState(t, dir, `{
		"resources": [
			{"type": "aws_instance", "name": "worker"},
			{"type": "aws_s3_bucket", "name": "data"}
		]
	}`)

	ws := Load(dir, nil)
	if ws.Err != nil {
		t.Fatalf("unexpected error: %v", ws.Err)
	}
	if len(ws.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(ws.Resources))
	}
	if ws.Resources[0].DisplayName() != "aws_instance.worker" {
		t.Fatalf("unexpected display name: %s", ws.Resources[0].DisplayName())
	}
	if ws.Resources[1].DisplayName() != "aws_s3_bucket.data" {
		t.Fatalf("unexpected display name: %s", ws.Resources[1].DisplayName())
	}
}

func TestLoad_MissingStateFile(t *testing.T) {
	dir := t.TempDir()
	ws := Load(dir, nil)
	if ws.Err == nil {
		t.Fatal("expected error for missing state file")
	}
	if len(ws.Resources) != 0 {
		t.Fatalf("expected no resources, got %d", len(ws.Resources))
	}
}

func TestLoad_EmptyResources(t *testing.T) {
	dir := t.TempDir()
	writeTFState(t, dir, `{"resources": []}`)
	ws := Load(dir, nil)
	if ws.Err != nil {
		t.Fatalf("unexpected error: %v", ws.Err)
	}
	if len(ws.Resources) != 0 {
		t.Fatalf("expected no resources, got %d", len(ws.Resources))
	}
}

func TestLoad_HasResources(t *testing.T) {
	dir := t.TempDir()
	writeTFState(t, dir, `{"resources": [{"type": "aws_vpc", "name": "main"}]}`)
	ws := Load(dir, nil)
	if !ws.HasResources() {
		t.Fatal("expected HasResources() true")
	}
}

func TestLoad_NoResources_HasResourcesFalse(t *testing.T) {
	dir := t.TempDir()
	writeTFState(t, dir, `{"resources": []}`)
	ws := Load(dir, nil)
	if ws.HasResources() {
		t.Fatal("expected HasResources() false")
	}
}

func TestLoad_CarriesTags(t *testing.T) {
	dir := t.TempDir()
	writeTFState(t, dir, `{"resources": [{"type": "aws_vpc", "name": "main"}]}`)

	ws := Load(dir, []string{"prod", "aws"})
	if len(ws.Tags) != 2 || ws.Tags[0] != "prod" {
		t.Fatalf("expected tags [prod aws], got %v", ws.Tags)
	}
}

func TestLoad_NilTags(t *testing.T) {
	dir := t.TempDir()
	writeTFState(t, dir, `{"resources": []}`)
	ws := Load(dir, nil)
	if ws.Tags != nil {
		t.Fatalf("expected nil tags, got %v", ws.Tags)
	}
}
