package cmd

import (
	"fmt"
	"os"

	"github.com/rasalas/envlint/internal/config"
	"github.com/rasalas/envlint/internal/gitcheck"
	"github.com/rasalas/envlint/internal/term"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "doctor",
		Short: "Check project setup for envlint",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor()
		},
	})
}

func runDoctor() error {
	fmt.Fprintf(term.W, "\n  %senvlint doctor%s\n", term.Primary, term.Reset)

	problems := 0

	// Check .env.example exists
	cfg := config.Default()
	if loaded, err := config.Load(); err == nil {
		cfg = loaded
	}

	fmt.Fprintln(term.W)
	fmt.Fprintf(term.W, "  %sFiles%s\n\n", term.Bold, term.Reset)

	if _, err := os.Stat(cfg.Example); err == nil {
		term.Pass(cfg.Example + " found")
	} else {
		term.Fail(cfg.Example + " not found")
		problems++
	}

	for _, envFile := range cfg.EnvFiles {
		if _, err := os.Stat(envFile); err == nil {
			term.Pass(envFile + " found")
		} else {
			term.WarnDetail(envFile, "not found (may be expected)")
		}
	}

	// Check config
	fmt.Fprintln(term.W)
	fmt.Fprintf(term.W, "  %sConfig%s\n\n", term.Bold, term.Reset)

	if _, err := config.Load(); err == nil {
		term.Pass(".envlint.toml found")
	} else {
		term.Info("no .envlint.toml (using defaults)")
	}

	// Git checks
	fmt.Fprintln(term.W)
	fmt.Fprintf(term.W, "  %sGit%s\n\n", term.Bold, term.Reset)

	if gitcheck.InGitRepo() {
		term.Pass("inside git repository")

		if gitcheck.GitignoreExists() {
			term.Pass(".gitignore found")

			if gitcheck.ContainsPattern(".env") {
				term.Pass(".env is in .gitignore")
			} else {
				term.Fail(".env is NOT in .gitignore — secrets may leak!")
				problems++
			}
		} else {
			term.Fail("no .gitignore found")
			problems++
		}
	} else {
		term.Info("not inside a git repository")
	}

	// Summary
	fmt.Fprintln(term.W)
	if problems == 0 {
		term.Pass("Everything looks good.")
	} else {
		fmt.Fprintf(term.W, "  %s%d problem(s) found — see above.%s\n", term.Red, problems, term.Reset)
	}
	fmt.Fprintln(term.W)

	return nil
}
