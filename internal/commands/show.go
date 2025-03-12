package commands

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/spf13/cobra"
	"os"
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
		// Continue without aliases
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	fmt.Println("Active co-authors:")
	t.AppendHeader(table.Row{"#", "Alias", "Name", "Email"})

	for i, author := range activeCoAuthors {
		// Find alias from config if available
		alias := ""
		for configAlias, configAuthor := range config.Config.CoAuthorsMap {
			if configAuthor.Email == author.Email {
				alias = configAlias
				break
			}
		}

		t.AppendRow([]interface{}{i, alias, author.Name, author.Email})
	}
	t.Render()
}
