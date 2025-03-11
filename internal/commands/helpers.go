package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/models"
)

func findCoAuthorByAliasOrIndex(identifier string) (models.CoAuthor, int, error) {
	// First try to parse as index
	var index int
	var author models.CoAuthor

	if _, err := fmt.Sscanf(identifier, "%d", &index); err == nil {
		// It's an index
		if index < 0 || index >= len(config.Config.CoAuthors) {
			return author, -1, fmt.Errorf("invalid co-author index")
		}
		return config.Config.CoAuthors[index], index, nil
	}

	// Then try as alias
	if coauthor, exists := config.Config.CoAuthorsMap[identifier]; exists {
		// Find the index in the slice for compatibility
		for i, author := range config.Config.CoAuthors {
			if author.Alias == identifier {
				return coauthor, i, nil
			}
		}
		return coauthor, -1, nil
	}

	return author, -1, fmt.Errorf("no co-author found with alias '%s'", identifier)
}
