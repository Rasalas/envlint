package lint

import (
	"testing"

	"github.com/rasalas/envlint/internal/env"
)

func TestCheckMissingKeys(t *testing.T) {
	example := map[string]env.Entry{
		"A": {Key: "A", Value: "x"},
		"B": {Key: "B", Value: ""},
		"C": {Key: "C", Value: "", Required: true},
	}
	actual := map[string]env.Entry{
		"A": {Key: "A", Value: "1"},
	}
	opts := Options{}

	issues := checkMissingKeys(example, actual, opts)
	if len(issues) != 2 {
		t.Fatalf("expected 2 missing keys, got %d", len(issues))
	}
	for _, issue := range issues {
		if issue.Rule != "missing-key" {
			t.Errorf("expected missing-key rule, got %s", issue.Rule)
		}
	}
}

func TestCheckExtraKeys(t *testing.T) {
	example := map[string]env.Entry{
		"A": {Key: "A"},
	}
	actual := map[string]env.Entry{
		"A": {Key: "A", Value: "1"},
		"B": {Key: "B", Value: "2"},
	}
	opts := Options{}

	issues := checkExtraKeys(example, actual, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 extra key, got %d", len(issues))
	}
	if issues[0].Key != "B" {
		t.Errorf("expected key B, got %s", issues[0].Key)
	}
}

func TestCheckRequiredEmpty(t *testing.T) {
	example := map[string]env.Entry{
		"A": {Key: "A", Value: "default", Required: true},
		"B": {Key: "B", Value: "", Required: false},
	}
	actual := map[string]env.Entry{
		"A": {Key: "A", Value: ""},
		"B": {Key: "B", Value: ""},
	}
	opts := Options{}

	issues := checkRequiredEmpty(example, actual, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 required-empty, got %d", len(issues))
	}
	if issues[0].Key != "A" {
		t.Errorf("expected key A, got %s", issues[0].Key)
	}
}

func TestCheckRequiredEmptyWithRef(t *testing.T) {
	example := map[string]env.Entry{
		"A": {Key: "A", Value: "required", Required: true},
	}
	actual := map[string]env.Entry{
		"A": {Key: "A", Value: "${OTHER}", IsRef: true},
	}
	opts := Options{}

	issues := checkRequiredEmpty(example, actual, opts)
	if len(issues) != 0 {
		t.Fatalf("expected 0 issues for ref value, got %d", len(issues))
	}
}

func TestCheckURLFormat(t *testing.T) {
	actual := map[string]env.Entry{
		"API_URL":  {Key: "API_URL", Value: "not-a-url"},
		"BASE_URL": {Key: "BASE_URL", Value: "https://example.com"},
		"NAME":     {Key: "NAME", Value: "hello"},
	}
	opts := Options{StrictURLs: true}

	issues := checkURLFormat(actual, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 URL issue, got %d", len(issues))
	}
	if issues[0].Key != "API_URL" {
		t.Errorf("expected API_URL, got %s", issues[0].Key)
	}
}

func TestCheckPortFormat(t *testing.T) {
	actual := map[string]env.Entry{
		"APP_PORT":  {Key: "APP_PORT", Value: "abc"},
		"SMTP_PORT": {Key: "SMTP_PORT", Value: "587"},
		"BAD_PORT":  {Key: "BAD_PORT", Value: "99999"},
	}
	opts := Options{StrictPorts: true}

	issues := checkPortFormat(actual, opts)
	if len(issues) != 2 {
		t.Fatalf("expected 2 port issues, got %d", len(issues))
	}
}

func TestCheckEmailFormat(t *testing.T) {
	actual := map[string]env.Entry{
		"ADMIN_EMAIL": {Key: "ADMIN_EMAIL", Value: "notanemail"},
		"USER_EMAIL":  {Key: "USER_EMAIL", Value: "user@example.com"},
	}
	opts := Options{}

	issues := checkEmailFormat(actual, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 email issue, got %d", len(issues))
	}
}

func TestCheckBooleanFormat(t *testing.T) {
	actual := map[string]env.Entry{
		"IS_ENABLED":   {Key: "IS_ENABLED", Value: "maybe"},
		"DEBUG":        {Key: "DEBUG", Value: "true"},
		"FEATURE_ACTIVE": {Key: "FEATURE_ACTIVE", Value: "yes"},
	}
	opts := Options{}

	issues := checkBooleanFormat(actual, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 boolean issue, got %d", len(issues))
	}
	if issues[0].Key != "IS_ENABLED" {
		t.Errorf("expected IS_ENABLED, got %s", issues[0].Key)
	}
}

func TestIgnoredKeys(t *testing.T) {
	example := map[string]env.Entry{
		"A": {Key: "A", Value: "x"},
	}
	actual := map[string]env.Entry{}
	opts := Options{IgnoreKeys: []string{"A"}}

	issues := checkMissingKeys(example, actual, opts)
	if len(issues) != 0 {
		t.Errorf("expected ignored key to produce no issues, got %d", len(issues))
	}
}

func TestExplicitlyRequiredKeys(t *testing.T) {
	example := map[string]env.Entry{
		"A": {Key: "A", Value: ""},
	}
	actual := map[string]env.Entry{
		"A": {Key: "A", Value: ""},
	}
	opts := Options{RequiredKeys: []string{"A"}}

	issues := checkRequiredEmpty(example, actual, opts)
	if len(issues) != 1 {
		t.Errorf("expected 1 required-empty for explicitly required key, got %d", len(issues))
	}
}
