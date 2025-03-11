package gittemplate

import (
	"os"
	"strings"

	"github.com/philippeckel/pair/internal/models"
)

// ParseActiveCoAuthors extracts co-authors from the current git template
func ParseActiveCoAuthors(templatePath string) ([]models.CoAuthor, error) {
	var activeCoAuthors []models.CoAuthor

	if templatePath == "" {
		return activeCoAuthors, nil
	}

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Co-authored-by:") {
			// Extract name and email from format: "Co-authored-by: Name <email>"
			authorInfo := strings.TrimSpace(strings.TrimPrefix(line, "Co-authored-by:"))
			parts := strings.Split(authorInfo, "<")
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				email := strings.TrimSpace(strings.TrimSuffix(parts[1], ">"))
				activeCoAuthors = append(activeCoAuthors, models.CoAuthor{Name: name, Email: email})
			}
		}
	}

	return activeCoAuthors, nil
}
