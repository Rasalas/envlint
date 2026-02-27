package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Example != ".env.example" {
		t.Errorf("expected .env.example, got %s", cfg.Example)
	}
	if len(cfg.EnvFiles) != 1 || cfg.EnvFiles[0] != ".env" {
		t.Errorf("unexpected envFiles: %v", cfg.EnvFiles)
	}
	if !cfg.Rules.StrictURLs {
		t.Error("expected StrictURLs to be true by default")
	}
}

func TestLoadFrom(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envlint.toml")
	content := `
example = ".env.production"
envFiles = [".env", ".env.local"]

[rules]
requireAll = true
noExtra = true
strictUrls = false

[rules.required]
keys = ["API_KEY", "DB_URL"]

[rules.ignore]
keys = ["DEBUG"]
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFrom(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Example != ".env.production" {
		t.Errorf("expected .env.production, got %s", cfg.Example)
	}
	if len(cfg.EnvFiles) != 2 {
		t.Fatalf("expected 2 envFiles, got %d", len(cfg.EnvFiles))
	}
	if !cfg.Rules.RequireAll {
		t.Error("expected requireAll to be true")
	}
	if !cfg.Rules.NoExtra {
		t.Error("expected noExtra to be true")
	}
	if cfg.Rules.StrictURLs {
		t.Error("expected strictUrls to be false")
	}
	if len(cfg.Rules.Required.Keys) != 2 {
		t.Errorf("expected 2 required keys, got %d", len(cfg.Rules.Required.Keys))
	}
	if len(cfg.Rules.Ignore.Keys) != 1 {
		t.Errorf("expected 1 ignore key, got %d", len(cfg.Rules.Ignore.Keys))
	}
}

func TestLoadFromMissing(t *testing.T) {
	_, err := LoadFrom("/nonexistent/.envlint.toml")
	if err == nil {
		t.Error("expected error for nonexistent config")
	}
}
