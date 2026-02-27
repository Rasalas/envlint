package lint

import (
	"testing"

	"github.com/rasalas/envlint/internal/env"
)

func TestCheckIntegration(t *testing.T) {
	example := []env.Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db", Required: true},
		{Key: "REDIS_URL", Value: "redis://localhost:6379"},
		{Key: "APP_PORT", Value: "3000"},
		{Key: "DEBUG", Value: ""},
		{Key: "API_KEY", Value: "", Required: true},
	}

	envEntries := []env.Entry{
		{Key: "DATABASE_URL", Value: "postgres://prod/db"},
		{Key: "APP_PORT", Value: "abc"},
		{Key: "DEBUG", Value: "true"},
		{Key: "API_KEY", Value: ""},
		{Key: "LEGACY_FLAG", Value: "old"},
	}

	opts := Options{
		StrictURLs:  true,
		StrictPorts: true,
	}

	result := Check(example, envEntries, opts)

	if result.TotalKeys() != 6 {
		t.Errorf("expected 6 total keys, got %d", result.TotalKeys())
	}

	// Should have: missing REDIS_URL, invalid APP_PORT, empty API_KEY, extra LEGACY_FLAG
	if result.ErrorCount() < 2 {
		t.Errorf("expected at least 2 errors, got %d", result.ErrorCount())
	}

	missing := result.ByRule("missing-key")
	if len(missing) != 1 {
		t.Errorf("expected 1 missing key, got %d", len(missing))
	}

	extra := result.ByRule("extra-key")
	if len(extra) != 1 {
		t.Errorf("expected 1 extra key, got %d", len(extra))
	}
}

func TestCheckNoIssues(t *testing.T) {
	example := []env.Entry{
		{Key: "FOO", Value: "bar"},
	}
	envEntries := []env.Entry{
		{Key: "FOO", Value: "baz"},
	}
	opts := Options{}

	result := Check(example, envEntries, opts)
	if result.HasErrors() {
		t.Error("expected no errors")
	}
}

func TestPromoteWarnings(t *testing.T) {
	result := Result{
		Issues: []Issue{
			{Severity: SeverityWarning, Rule: "extra-key"},
			{Severity: SeverityError, Rule: "missing-key"},
		},
	}
	result.PromoteWarnings()

	for _, issue := range result.Issues {
		if issue.Severity != SeverityError {
			t.Errorf("expected all issues to be errors after promote, got %s", issue.Severity)
		}
	}
}
