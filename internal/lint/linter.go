package lint

import "github.com/rasalas/envlint/internal/env"

// Options configures linter behavior.
type Options struct {
	Strict      bool
	NoExtra     bool
	StrictURLs  bool
	StrictPorts bool
	RequiredKeys []string
	IgnoreKeys   []string
}

// Check runs all lint rules against the given env entries.
func Check(exampleEntries, envEntries []env.Entry, opts Options) Result {
	example := env.ParseEntries(exampleEntries)
	actual := env.ParseEntries(envEntries)

	var result Result

	// Count total unique keys
	allKeys := make(map[string]bool)
	for k := range example {
		allKeys[k] = true
	}
	for k := range actual {
		allKeys[k] = true
	}
	result.SetTotalKeys(len(allKeys))

	// Run all rules
	for _, issue := range checkMissingKeys(example, actual, opts) {
		result.AddIssue(issue)
	}

	if opts.NoExtra {
		for _, issue := range checkExtraKeys(example, actual, opts) {
			issue.Severity = SeverityError
			result.AddIssue(issue)
		}
	} else {
		for _, issue := range checkExtraKeys(example, actual, opts) {
			result.AddIssue(issue)
		}
	}

	for _, issue := range checkRequiredEmpty(example, actual, opts) {
		result.AddIssue(issue)
	}

	for _, issue := range checkURLFormat(actual, opts) {
		result.AddIssue(issue)
	}

	for _, issue := range checkPortFormat(actual, opts) {
		result.AddIssue(issue)
	}

	for _, issue := range checkEmailFormat(actual, opts) {
		result.AddIssue(issue)
	}

	for _, issue := range checkBooleanFormat(actual, opts) {
		result.AddIssue(issue)
	}

	return result
}
