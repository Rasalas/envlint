package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rasalas/envlint/internal/config"
	"github.com/rasalas/envlint/internal/env"
	"github.com/rasalas/envlint/internal/lint"
	"github.com/rasalas/envlint/internal/term"
	"github.com/spf13/cobra"
)

var (
	exampleFlag string
	envFlag     string
	strictFlag  bool
	formatFlag  string
	quietFlag   bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&exampleFlag, "example", "", "Path to example env file (default: .env.example)")
	rootCmd.PersistentFlags().StringVar(&envFlag, "env", "", "Path to env file to check (default: .env)")
	rootCmd.Flags().BoolVar(&strictFlag, "strict", false, "Treat warnings as errors")
	rootCmd.Flags().StringVar(&formatFlag, "format", "text", "Output format: text or json")
	rootCmd.Flags().BoolVar(&quietFlag, "quiet", false, "Only show errors")
}

var rootCmd = &cobra.Command{
	Use:   "envlint",
	Short: "Validate .env files against .env.example",
	Long:  "Check for missing keys, value formats, empty required fields, and more.",
	RunE:  runLint,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if exitErr, ok := err.(*exitError); ok {
			os.Exit(exitErr.code)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}

type exitError struct {
	code int
}

func (e *exitError) Error() string {
	return fmt.Sprintf("exit %d", e.code)
}

func runLint(cmd *cobra.Command, args []string) error {
	// Load config
	cfg := config.Default()
	if loaded, err := config.Load(); err == nil {
		cfg = loaded
	}

	// Determine file paths (flags override config)
	examplePath := cfg.Example
	if exampleFlag != "" {
		examplePath = exampleFlag
	}
	envPath := ".env"
	if len(cfg.EnvFiles) > 0 {
		envPath = cfg.EnvFiles[0]
	}
	if envFlag != "" {
		envPath = envFlag
	}

	// Parse files
	exampleEntries, err := env.ParseFile(examplePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return &exitError{code: 2}
	}

	envEntries, err := env.ParseFile(envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return &exitError{code: 2}
	}

	// Build linter options from config
	opts := lint.Options{
		Strict:     strictFlag || cfg.Rules.RequireAll,
		NoExtra:    cfg.Rules.NoExtra,
		StrictURLs: cfg.Rules.StrictURLs,
		StrictPorts: cfg.Rules.StrictPorts,
		RequiredKeys: cfg.Rules.Required.Keys,
		IgnoreKeys:   cfg.Rules.Ignore.Keys,
	}

	// Run linter
	result := lint.Check(exampleEntries, envEntries, opts)

	// Promote warnings to errors in strict mode
	if strictFlag {
		result.PromoteWarnings()
	}

	// Output
	switch formatFlag {
	case "json":
		return outputJSON(result)
	default:
		outputText(result, envPath, examplePath)
	}

	if result.ErrorCount() > 0 {
		return &exitError{code: 1}
	}
	return nil
}

func outputJSON(result lint.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result.ToJSON())
}

func outputText(result lint.Result, envPath, examplePath string) {
	term.Title(envPath, examplePath)

	// Group issues by category
	missing := result.ByRule("missing-key")
	extra := result.ByRule("extra-key")
	values := result.ValueIssues()

	if len(missing) > 0 {
		term.Header("Missing keys")
		for _, issue := range missing {
			if issue.Severity == lint.SeverityError {
				suffix := ""
				if issue.Detail != "" {
					suffix = "  " + term.Dim + "(" + issue.Detail + ")" + term.Reset
				}
				term.Fail(issue.Key + suffix)
			} else {
				term.Warn(issue.Key)
			}
		}
	}

	if len(extra) > 0 && !quietFlag {
		term.Header("Extra keys")
		for _, issue := range extra {
			term.Warn(issue.Key)
		}
	}

	if len(values) > 0 {
		term.Header("Value problems")
		for _, issue := range values {
			if issue.Severity == lint.SeverityError {
				term.FailDetail(issue.Key, issue.Detail)
			} else if !quietFlag {
				term.WarnDetail(issue.Key, issue.Detail)
			}
		}
	}

	total := result.TotalKeys()
	valid := total - result.ErrorCount() - result.WarnCount()
	if valid < 0 {
		valid = 0
	}
	term.Summary(valid, total, result.ErrorCount(), result.WarnCount())
	fmt.Fprintln(term.W)
}
