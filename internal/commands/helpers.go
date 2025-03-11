package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/models"
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
