package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/philippeckel/pair/internal/models"
	"github.com/spf13/cobra"
)

func showActiveCoAuthors(cmd *cobra.Command, args []string) {
	templatePath, err := gittemplate.GetCurrentTemplate()
	if err != nil {
		fmt.Printf("Error getting current git template: %v\n", err)
		return
	}

	if templatePath == "" {
		fmt.Println("No commit template is currently set. No active co-authors.")
		return
	}

	activeCoAuthors, err := gittemplate.ParseActiveCoAuthors(templatePath)
	if err != nil {
		fmt.Printf("Error parsing active co-authors: %v\n", err)
		return
	}

	if len(activeCoAuthors) == 0 {
		fmt.Println("No active co-authors found.")
		return
	}

	// Load config to get aliases
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Warning: Could not load config for aliases: %v\n", err)
		// Continue without aliases (getAlias will return empty strings)
	}

	// Use the extracted helper function
	renderCoAuthorTable("Active co-authors:", activeCoAuthors, func(author models.CoAuthor) string {
		// Find alias from config if available
		for configAlias, configAuthor := range config.Config.CoAuthorsMap {
			if configAuthor.Email == author.Email {
				return configAlias
			}
		}
		return ""
	})
}
