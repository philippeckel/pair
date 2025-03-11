package commands

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	fmt.Println("Active co-authors:")
	t.AppendHeader(table.Row{"Name", "Email"})

	for _, author := range activeCoAuthors {
		t.AppendRow([]interface{}{author.Name, author.Email})
	}
	t.Render()
}
