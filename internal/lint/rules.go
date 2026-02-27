package lint

import (
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/rasalas/envlint/internal/env"
)

// checkMissingKeys reports keys in example that are missing from env.
func checkMissingKeys(example, actual map[string]env.Entry, opts Options) []Issue {
	var issues []Issue
	for key, ex := range example {
		if isIgnored(key, opts) {
			continue
		}
		if _, ok := actual[key]; !ok {
			detail := ""
			if ex.Required || ex.Value != "" || isExplicitlyRequired(key, opts) {
				detail = "required"
			}
			issues = append(issues, Issue{
				Rule:     "missing-key",
				Key:      key,
				Severity: SeverityError,
				Detail:   detail,
			})
		}
	}
	return issues
}

// checkExtraKeys reports keys in env that are not in example.
func checkExtraKeys(example, actual map[string]env.Entry, opts Options) []Issue {
	var issues []Issue
	for key := range actual {
		if isIgnored(key, opts) {
			continue
		}
		if _, ok := example[key]; !ok {
			issues = append(issues, Issue{
				Rule:     "extra-key",
				Key:      key,
				Severity: SeverityWarning,
			})
		}
	}
	return issues
}

// checkRequiredEmpty reports required keys that have empty values.
func checkRequiredEmpty(example, actual map[string]env.Entry, opts Options) []Issue {
	var issues []Issue
	for key, ex := range example {
		if isIgnored(key, opts) {
			continue
		}
		act, ok := actual[key]
		if !ok {
			continue // handled by checkMissingKeys
		}
		// Skip if value contains a variable reference (treated as non-empty)
		if act.IsRef {
			continue
		}
		isReq := ex.Required || ex.Value != "" || isExplicitlyRequired(key, opts)
		if isReq && strings.TrimSpace(act.Value) == "" {
			issues = append(issues, Issue{
				Rule:     "required-empty",
				Key:      key,
				Severity: SeverityError,
				Detail:   "required but empty",
				LineNum:  act.LineNum,
			})
		}
	}
	return issues
}

// checkURLFormat validates keys containing "URL" have valid URL values.
func checkURLFormat(actual map[string]env.Entry, opts Options) []Issue {
	if !opts.StrictURLs {
		return nil
	}
	var issues []Issue
	for key, entry := range actual {
		if isIgnored(key, opts) {
			continue
		}
		if !containsCI(key, "URL") {
			continue
		}
		val := strings.TrimSpace(entry.Value)
		if val == "" || entry.IsRef {
			continue
		}
		u, err := url.Parse(val)
		if err != nil || u.Scheme == "" || u.Host == "" {
			issues = append(issues, Issue{
				Rule:     "invalid-url",
				Key:      key,
				Severity: SeverityError,
				Detail:   "invalid URL format",
				LineNum:  entry.LineNum,
			})
		}
	}
	return issues
}

// checkPortFormat validates keys containing "PORT" have valid port numbers.
func checkPortFormat(actual map[string]env.Entry, opts Options) []Issue {
	if !opts.StrictPorts {
		return nil
	}
	var issues []Issue
	for key, entry := range actual {
		if isIgnored(key, opts) {
			continue
		}
		if !containsCI(key, "PORT") {
			continue
		}
		val := strings.TrimSpace(entry.Value)
		if val == "" || entry.IsRef {
			continue
		}
		port, err := strconv.Atoi(val)
		if err != nil || port < 1 || port > 65535 {
			issues = append(issues, Issue{
				Rule:     "invalid-port",
				Key:      key,
				Severity: SeverityError,
				Detail:   "must be 1-65535, got " + strconv.Quote(val),
				LineNum:  entry.LineNum,
			})
		}
	}
	return issues
}

// checkEmailFormat validates keys containing "EMAIL" have a basic email format.
func checkEmailFormat(actual map[string]env.Entry, opts Options) []Issue {
	var issues []Issue
	for key, entry := range actual {
		if isIgnored(key, opts) {
			continue
		}
		if !containsCI(key, "EMAIL") {
			continue
		}
		val := strings.TrimSpace(entry.Value)
		if val == "" || entry.IsRef {
			continue
		}
		if !strings.Contains(val, "@") || !strings.Contains(val, ".") {
			issues = append(issues, Issue{
				Rule:     "invalid-email",
				Key:      key,
				Severity: SeverityWarning,
				Detail:   "invalid email format",
				LineNum:  entry.LineNum,
			})
		}
	}
	return issues
}

// checkBooleanFormat validates keys with boolean-like names have boolean values.
func checkBooleanFormat(actual map[string]env.Entry, opts Options) []Issue {
	var issues []Issue
	boolValues := map[string]bool{
		"true": true, "false": true,
		"1": true, "0": true,
		"yes": true, "no": true,
		"on": true, "off": true,
	}
	for key, entry := range actual {
		if isIgnored(key, opts) {
			continue
		}
		if !isBooleanKey(key) {
			continue
		}
		val := strings.TrimSpace(entry.Value)
		if val == "" || entry.IsRef {
			continue
		}
		if !boolValues[strings.ToLower(val)] {
			issues = append(issues, Issue{
				Rule:     "invalid-boolean",
				Key:      key,
				Severity: SeverityWarning,
				Detail:   "expected boolean value",
				LineNum:  entry.LineNum,
			})
		}
	}
	return issues
}

// isBooleanKey checks if a key name suggests a boolean value.
func isBooleanKey(key string) bool {
	upper := strings.ToUpper(key)
	return strings.Contains(upper, "ENABLED") ||
		strings.Contains(upper, "ACTIVE") ||
		strings.HasPrefix(upper, "IS_") ||
		strings.Contains(upper, "DISABLE") ||
		strings.Contains(upper, "DEBUG")
}

func containsCI(s, substr string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substr))
}

func isIgnored(key string, opts Options) bool {
	return slices.Contains(opts.IgnoreKeys, key)
}

func isExplicitlyRequired(key string, opts Options) bool {
	return slices.Contains(opts.RequiredKeys, key)
}
