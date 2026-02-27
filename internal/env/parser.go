package env

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{?[A-Za-z_][A-Za-z0-9_]*\}?`)

// ParseFile reads an env file and returns its entries.
func ParseFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	lineNum := 0
	var multilineKey string
	var multilineValue strings.Builder
	var multilineStart int

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Continue multiline value
		if multilineKey != "" {
			multilineValue.WriteString("\n")
			multilineValue.WriteString(line)
			if strings.HasSuffix(strings.TrimSpace(line), `"`) {
				val := multilineValue.String()
				// Strip surrounding quotes
				val = strings.TrimSuffix(val, `"`)
				entry := Entry{
					Key:     multilineKey,
					Value:   val,
					LineNum: multilineStart,
					IsRef:   refPattern.MatchString(val),
				}
				entries = append(entries, entry)
				multilineKey = ""
				multilineValue.Reset()
			}
			continue
		}

		trimmed := strings.TrimSpace(line)

		// Skip empty lines and full-line comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Split on first '='
		eqIdx := strings.IndexByte(trimmed, '=')
		if eqIdx < 0 {
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		rest := trimmed[eqIdx+1:]

		// Check for multiline start: value begins with " but doesn't end with "
		stripped := strings.TrimSpace(rest)
		if strings.HasPrefix(stripped, `"`) && !strings.HasSuffix(stripped, `"`) {
			multilineKey = key
			multilineStart = lineNum
			multilineValue.WriteString(strings.TrimPrefix(stripped, `"`))
			continue
		}

		value, comment := parseValueAndComment(rest)
		required := strings.Contains(strings.ToLower(comment), "required")

		entry := Entry{
			Key:      key,
			Value:    value,
			Comment:  comment,
			LineNum:  lineNum,
			Required: required,
			IsRef:    refPattern.MatchString(value),
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %w", path, err)
	}

	return entries, nil
}

// parseValueAndComment splits the raw value part into the actual value and any inline comment.
func parseValueAndComment(raw string) (string, string) {
	raw = strings.TrimSpace(raw)

	// Quoted value: find closing quote, rest is comment
	if strings.HasPrefix(raw, `"`) {
		end := strings.Index(raw[1:], `"`)
		if end >= 0 {
			value := raw[1 : end+1]
			rest := strings.TrimSpace(raw[end+2:])
			comment := ""
			if strings.HasPrefix(rest, "#") {
				comment = strings.TrimSpace(rest[1:])
			}
			return value, comment
		}
	}
	if strings.HasPrefix(raw, `'`) {
		end := strings.Index(raw[1:], `'`)
		if end >= 0 {
			value := raw[1 : end+1]
			rest := strings.TrimSpace(raw[end+2:])
			comment := ""
			if strings.HasPrefix(rest, "#") {
				comment = strings.TrimSpace(rest[1:])
			}
			return value, comment
		}
	}

	// Value starts with # → entire thing is a comment (empty value)
	if strings.HasPrefix(raw, "#") {
		return "", strings.TrimSpace(raw[1:])
	}

	// Unquoted: split on first # that has a space before it
	if idx := strings.Index(raw, " #"); idx >= 0 {
		value := strings.TrimSpace(raw[:idx])
		comment := strings.TrimSpace(raw[idx+2:])
		return value, comment
	}

	return raw, ""
}

// ParseEntries builds a map from key to Entry for quick lookup.
func ParseEntries(entries []Entry) map[string]Entry {
	m := make(map[string]Entry, len(entries))
	for _, e := range entries {
		m[e.Key] = e
	}
	return m
}
