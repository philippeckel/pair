package commands

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/philippeckel/pair/internal/config"
	"github.com/spf13/cobra"
	"os"
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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	fmt.Println("Available co-authors:")
	t.AppendHeader(table.Row{"#", "Alias", "Name", "Email"})
	for i, author := range config.Config.CoAuthors {
		t.AppendRow([]interface{}{i, author.Alias, author.Name, author.Email})
	}
	t.Render()
}
