package gittemplate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/philippeckel/pair/internal/models"
)

// GetCurrentTemplate returns the path to the current git commit template
func GetCurrentTemplate() (string, error) {
	cmd := exec.Command("git", "config", "--get", "commit.template")
	output, err := cmd.Output()
	if err != nil {
		// It's ok if the template doesn't exist yet
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// UpdateTemplate writes a new git template with the given co-authors
func UpdateTemplate(activeCoAuthors []models.CoAuthor) error {
	// Create a temporary file for the commit template
	tempDir := os.TempDir()
	templatePath := filepath.Join(tempDir, "git_commit_template")

	var content strings.Builder
	content.WriteString("\n\n") // Leave space for commit message
	content.WriteString("# Co-authors:\n")

	for _, author := range activeCoAuthors {
		content.WriteString(fmt.Sprintf("Co-authored-by: %s <%s>\n", author.Name, author.Email))
	}

	if err := os.WriteFile(templatePath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	// Set the commit.template configuration
	cmd := exec.Command("git", "config", "--global", "commit.template", templatePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set git commit template: %w", err)
	}

	return nil
}
