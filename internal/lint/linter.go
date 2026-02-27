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
	result.addAll(checkMissingKeys(example, actual, opts))

	extras := checkExtraKeys(example, actual, opts)
	if opts.NoExtra {
		for i := range extras {
			extras[i].Severity = SeverityError
		}
	}
	result.addAll(extras)

	result.addAll(checkRequiredEmpty(example, actual, opts))
	result.addAll(checkURLFormat(actual, opts))
	result.addAll(checkPortFormat(actual, opts))
	result.addAll(checkEmailFormat(actual, opts))
	result.addAll(checkBooleanFormat(actual, opts))

	return result
}
