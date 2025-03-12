package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/models"
	"github.com/spf13/cobra"
)

func listCoAuthors(cmd *cobra.Command, args []string) {
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(config.Config.CoAuthors) == 0 {
		fmt.Println("No co-authors found in config. Use 'pair init' to create a sample config.")
		return
	}

	// Use the extracted helper function
	renderCoAuthorTable("Available co-authors:", config.Config.CoAuthors, func(author models.CoAuthor) string {
		return author.Alias
	})
}
