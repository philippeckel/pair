package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
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

	fmt.Println("Available co-authors:")
	for i, author := range config.Config.CoAuthors {
		fmt.Printf("%d: %s (%s) <%s>\n", i, author.Name, author.Alias, author.Email)
	}
}
