package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/philippeckel/pair/internal/models"
	"github.com/spf13/cobra"
)

func clearCoAuthors(cmd *cobra.Command, args []string) {
	// Update with empty authors list
	if err := gittemplate.UpdateTemplate([]models.CoAuthor{}); err != nil {
		fmt.Printf("Error clearing co-authors: %v\n", err)
		return
	}
	fmt.Println("All co-authors have been cleared")
}
