package gitcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContainsPattern(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)

	content := ".env\n.env.local\n*.log\n"
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if !ContainsPattern(".env") {
		t.Error("expected .env to be found")
	}
	if !ContainsPattern(".env.local") {
		t.Error("expected .env.local to be found")
	}
	if ContainsPattern(".env.example") {
		t.Error("did not expect .env.example to be found")
	}
}

func TestGitignoreExists(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)

	if GitignoreExists() {
		t.Error("expected no .gitignore")
	}

	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(""), 0644)
	if !GitignoreExists() {
		t.Error("expected .gitignore to exist")
	}
}
