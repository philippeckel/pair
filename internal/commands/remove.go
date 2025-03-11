package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/spf13/cobra"
)

func removeCoAuthor(cmd *cobra.Command, args []string) error {
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

	if len(activeCoAuthors) == 0 {
		return fmt.Errorf("no active co-authors to remove")
	}

	identifier := args[0]
	var indexToRemove = -1

	// Try to parse as index in the active co-authors list
	var index int
	if _, err := fmt.Sscanf(identifier, "%d", &index); err == nil {
		if index >= 0 && index < len(activeCoAuthors) {
			indexToRemove = index
		}
	}

	// If not found by index, try finding by alias in the config
	if indexToRemove == -1 {
		coAuthor, _, err := findCoAuthorByAliasOrIndex(identifier)
		if err != nil {
			return err
		}

		// Find this co-author in the active list
		for i, active := range activeCoAuthors {
			if active.Email == coAuthor.Email {
				indexToRemove = i
				break
			}
		}

		if indexToRemove == -1 {
			return fmt.Errorf("co-author '%s' is not currently active", coAuthor.Name)
		}
	}

	// Remove the co-author
	removedAuthor := activeCoAuthors[indexToRemove]
	activeCoAuthors = append(activeCoAuthors[:indexToRemove], activeCoAuthors[indexToRemove+1:]...)

	// Update template
	if err := gittemplate.UpdateTemplate(activeCoAuthors); err != nil {
		return err
	}

	fmt.Printf("Removed co-author: %s <%s>\n", removedAuthor.Name, removedAuthor.Email)
	return nil
}
