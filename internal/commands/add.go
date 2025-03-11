package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/spf13/cobra"
)

func addCoAuthor(cmd *cobra.Command, args []string) error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	// Get current template
	templatePath, err := gittemplate.GetCurrentTemplate()
	if err != nil {
		return err
	}

	// Get active co-authors
	activeCoAuthors, err := gittemplate.ParseActiveCoAuthors(templatePath)
	if err != nil {
		return err
	}

	// Process each co-author identifier provided in args
	for _, identifier := range args {
		coAuthor, _, err := findCoAuthorByAliasOrIndex(identifier)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue // Continue processing other identifiers even if one fails
		}

		// Check if co-author is already active
		alreadyActive := false
		for _, active := range activeCoAuthors {
			if active.Email == coAuthor.Email {
				alreadyActive = true
				fmt.Printf("Co-author already active: %s <%s>\n", coAuthor.Name, coAuthor.Email)
				break
			}
		}

		// Add co-author if not already active
		if !alreadyActive {
			activeCoAuthors = append(activeCoAuthors, coAuthor)
			fmt.Printf("Adding co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
		}
	}

	// Update git template with all co-authors
	if err := gittemplate.UpdateTemplate(activeCoAuthors); err != nil {
		return err
	}

	return nil
}
