package commands

import (
	"fmt"

	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/philippeckel/pair/internal/models"
	"github.com/spf13/cobra"
)

// UnselectCmd allows interactively removing co-authors
var unselectCmd = &cobra.Command{
	Use:   "unselect",
	Short: "Interactively remove co-authors using fuzzy finder",
	Long:  `Use fuzzy finder to interactively remove active co-authors`,
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return fmt.Errorf("no active co-authors found")
		}

		// Always use multiple selection
		toRemove, err := selectMultipleCoAuthorsToRemove(activeCoAuthors)
		if err != nil {
			return err
		}

		if len(toRemove) == 0 {
			return fmt.Errorf("no co-authors selected for removal")
		}

		// Create new list without the removed co-authors
		var newActiveCoAuthors []models.CoAuthor
		for _, author := range activeCoAuthors {
			shouldRemove := false
			for _, remove := range toRemove {
				if author.Email == remove.Email {
					shouldRemove = true
					fmt.Printf("Removing co-author: %s <%s>\n", remove.Name, remove.Email)
					break
				}
			}
			if !shouldRemove {
				newActiveCoAuthors = append(newActiveCoAuthors, author)
			}
		}

		// Update git template
		if err := gittemplate.UpdateTemplate(newActiveCoAuthors); err != nil {
			return err
		}

		fmt.Printf("Removed %d co-authors from your commit template\n", len(toRemove))
		return nil
	},
	Aliases: []string{"us"},
}
