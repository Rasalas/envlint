package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestParseBasicKeyValue(t *testing.T) {
	path := writeTmp(t, "FOO=bar\nBAZ=123\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
	if entries[1].Key != "BAZ" || entries[1].Value != "123" {
		t.Errorf("unexpected entry: %+v", entries[1])
	}
}

func TestParseQuotedValues(t *testing.T) {
	path := writeTmp(t, `DB_URL="postgres://localhost/db"
SECRET='my secret value'
`)
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Value != "postgres://localhost/db" {
		t.Errorf("expected unquoted value, got %q", entries[0].Value)
	}
	if entries[1].Value != "my secret value" {
		t.Errorf("expected unquoted value, got %q", entries[1].Value)
	}
}

func TestParseComments(t *testing.T) {
	path := writeTmp(t, `# This is a comment
FOO=bar
BAZ=123 # inline comment
REQUIRED_KEY= # required
`)
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[1].Comment != "inline comment" {
		t.Errorf("expected inline comment, got %q", entries[1].Comment)
	}
	if !entries[2].Required {
		t.Error("expected REQUIRED_KEY to be required")
	}
}

func TestParseEmptyValues(t *testing.T) {
	path := writeTmp(t, "EMPTY=\nALSO_EMPTY=   \n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Value != "" {
		t.Errorf("expected empty value, got %q", entries[0].Value)
	}
}

func TestParseVariableReferences(t *testing.T) {
	path := writeTmp(t, `BASE_URL=http://localhost
FULL_URL=${BASE_URL}/api
OTHER=$HOME/path
`)
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].IsRef {
		t.Error("BASE_URL should not be a ref")
	}
	if !entries[1].IsRef {
		t.Error("FULL_URL should be a ref")
	}
	if !entries[2].IsRef {
		t.Error("OTHER should be a ref")
	}
}

func TestParseMultiline(t *testing.T) {
	path := writeTmp(t, `SINGLE=value
MULTI="line one
line two
line three"
AFTER=done
`)
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[1].Key != "MULTI" {
		t.Errorf("expected MULTI, got %s", entries[1].Key)
	}
	expected := "line one\nline two\nline three"
	if entries[1].Value != expected {
		t.Errorf("expected multiline value %q, got %q", expected, entries[1].Value)
	}
}

func TestParseSkipsInvalidLines(t *testing.T) {
	path := writeTmp(t, "valid=yes\nno_equals_here\nalso_valid=true\n")
	entries, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestParseEntries(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	m := ParseEntries(entries)
	if m["A"].Value != "1" || m["B"].Value != "2" {
		t.Error("ParseEntries map incorrect")
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}
