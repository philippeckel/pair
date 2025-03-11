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

	identifier := args[0]
	coAuthor, _, err := findCoAuthorByAliasOrIndex(identifier)
	if err != nil {
		return err
	}

	// Validate the co-author
	if err := coAuthor.Validate(); err != nil {
		return fmt.Errorf("invalid co-author: %w", err)
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

	// Check if co-author is already active
	for _, active := range activeCoAuthors {
		if active.Email == coAuthor.Email {
			return fmt.Errorf("co-author '%s' is already active", coAuthor.Name)
		}
	}

	// Add the co-author
	activeCoAuthors = append(activeCoAuthors, coAuthor)

	// Update template
	if err := gittemplate.UpdateTemplate(activeCoAuthors); err != nil {
		return err
	}

	fmt.Printf("Added co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
	return nil
}
