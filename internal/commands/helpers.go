package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/models"
	exec "os/exec"
	"strings"
)

// findCoAuthorByAliasOrIndex finds a co-author by alias or index
func findCoAuthorByAliasOrIndex(identifier string) (models.CoAuthor, int, error) {
	// Try to find by alias first
	if coAuthor, exists := config.Config.CoAuthorsMap[identifier]; exists {
		// Return -1 to indicate found by alias, not by index
		return coAuthor, -1, nil
	}

	// Try to parse as index
	var index int
	if _, err := fmt.Sscanf(identifier, "%d", &index); err == nil {
		// It's an index
		if index < 0 || index >= len(config.Config.CoAuthors) {
			return models.CoAuthor{}, -1, fmt.Errorf("invalid co-author index: %d", index)
		}
		return config.Config.CoAuthors[index], index, nil
	}

	return models.CoAuthor{}, -1, fmt.Errorf("no co-author found with alias '%s'", identifier)
}

// getGitUserInfo retrieves the current git user.name and user.email
func getGitUserInfo() (name string, email string, err error) {
	// Get user name
	cmdName := exec.Command("git", "config", "--get", "user.name")
	nameBytes, err := cmdName.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get git user.name: %w", err)
	}
	name = strings.TrimSpace(string(nameBytes))

	// Get user email
	cmdEmail := exec.Command("git", "config", "--get", "user.email")
	emailBytes, err := cmdEmail.Output()
	if err != nil {
		return name, "", fmt.Errorf("failed to get git user.email: %w", err)
	}
	email = strings.TrimSpace(string(emailBytes))

	return name, email, nil
}
