package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rasalas/envlint/internal/env"
	"github.com/rasalas/envlint/internal/term"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Generate .env.example from existing .env",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit()
		},
	})
}

func runInit() error {
	envPath := ".env"
	if envFlag != "" {
		envPath = envFlag
	}
	examplePath := ".env.example"
	if exampleFlag != "" {
		examplePath = exampleFlag
	}

	// Check if .env exists
	if _, err := os.Stat(envPath); err != nil {
		return fmt.Errorf("%s not found", envPath)
	}

	// Check if .env.example already exists
	if _, err := os.Stat(examplePath); err == nil {
		return fmt.Errorf("%s already exists — remove it first or use a different name", examplePath)
	}

	// Read the original .env file line by line and strip values
	in, err := os.Open(envPath)
	if err != nil {
		return err
	}
	defer in.Close()

	entries, err := env.ParseFile(envPath)
	if err != nil {
		return err
	}
	entryMap := env.ParseEntries(entries)

	var out strings.Builder
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Pass through empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			out.WriteString(line)
			out.WriteString("\n")
			continue
		}

		eqIdx := strings.IndexByte(trimmed, '=')
		if eqIdx < 0 {
			out.WriteString(line)
			out.WriteString("\n")
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		entry, exists := entryMap[key]

		// Keep comments but strip values
		comment := ""
		if exists && entry.Comment != "" {
			comment = " # " + entry.Comment
		}

		out.WriteString(key + "=" + comment + "\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := os.WriteFile(examplePath, []byte(out.String()), 0644); err != nil {
		return err
	}

	fmt.Fprintf(term.W, "\n  %s✓%s Generated %s from %s (%d keys)\n\n",
		term.Green, term.Reset, examplePath, envPath, len(entries))

	return nil
}
