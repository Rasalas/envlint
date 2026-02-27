package gitcheck

import (
	"os"
	"os/exec"
	"strings"
)

// InGitRepo returns true if the current directory is inside a git repository.
func InGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// IsIgnored checks if a file path is listed in .gitignore (tracked by git's check-ignore).
func IsIgnored(path string) bool {
	cmd := exec.Command("git", "check-ignore", "-q", path)
	return cmd.Run() == nil
}

// GitignoreExists returns true if a .gitignore file exists in the current directory.
func GitignoreExists() bool {
	_, err := os.Stat(".gitignore")
	return err == nil
}

// ContainsPattern checks if .gitignore contains a specific pattern.
func ContainsPattern(pattern string) bool {
	data, err := os.ReadFile(".gitignore")
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == pattern {
			return true
		}
	}
	return false
}
