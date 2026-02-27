package lint

// Severity represents the severity of a lint issue.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Issue represents a single lint finding.
type Issue struct {
	Rule     string   `json:"rule"`
	Key      string   `json:"key"`
	Severity Severity `json:"severity"`
	Detail   string   `json:"detail,omitempty"`
	LineNum  int      `json:"line,omitempty"`
}

// Result holds all lint findings.
type Result struct {
	Issues    []Issue `json:"issues"`
	totalKeys int
}

// AddIssue appends an issue to the result.
func (r *Result) AddIssue(issue Issue) {
	r.Issues = append(r.Issues, issue)
}

// SetTotalKeys sets the total number of keys checked.
func (r *Result) SetTotalKeys(n int) {
	r.totalKeys = n
}

// TotalKeys returns the total number of keys.
func (r Result) TotalKeys() int {
	return r.totalKeys
}

// ErrorCount returns the number of error-level issues.
func (r Result) ErrorCount() int {
	n := 0
	for _, issue := range r.Issues {
		if issue.Severity == SeverityError {
			n++
		}
	}
	return n
}

// WarnCount returns the number of warning-level issues.
func (r Result) WarnCount() int {
	n := 0
	for _, issue := range r.Issues {
		if issue.Severity == SeverityWarning {
			n++
		}
	}
	return n
}

// HasErrors returns true if there are any error-level issues.
func (r Result) HasErrors() bool {
	return r.ErrorCount() > 0
}

// ByRule returns issues matching the given rule name.
func (r Result) ByRule(rule string) []Issue {
	var out []Issue
	for _, issue := range r.Issues {
		if issue.Rule == rule {
			out = append(out, issue)
		}
	}
	return out
}

// ValueIssues returns all issues that are not missing-key or extra-key.
func (r Result) ValueIssues() []Issue {
	var out []Issue
	for _, issue := range r.Issues {
		if issue.Rule != "missing-key" && issue.Rule != "extra-key" {
			out = append(out, issue)
		}
	}
	return out
}

// PromoteWarnings changes all warnings to errors.
func (r *Result) PromoteWarnings() {
	for i := range r.Issues {
		if r.Issues[i].Severity == SeverityWarning {
			r.Issues[i].Severity = SeverityError
		}
	}
}

// JSONOutput is the JSON-serializable output format.
type JSONOutput struct {
	Valid    bool    `json:"valid"`
	Total   int     `json:"total"`
	Errors  int     `json:"errors"`
	Warns   int     `json:"warnings"`
	Issues  []Issue `json:"issues"`
}

// ToJSON converts the result to a JSON-friendly struct.
func (r Result) ToJSON() JSONOutput {
	issues := r.Issues
	if issues == nil {
		issues = []Issue{}
	}
	return JSONOutput{
		Valid:  !r.HasErrors(),
		Total:  r.totalKeys,
		Errors: r.ErrorCount(),
		Warns:  r.WarnCount(),
		Issues: issues,
	}
}
