package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/spf13/cobra"
)

func clearCoAuthors(cmd *cobra.Command, args []string) {
	if err := gittemplate.ClearTemplate(); err != nil {
		fmt.Printf("Error clearing co-authors: %v\n", err)
		return
	}
	fmt.Println("All co-authors have been cleared")
}
