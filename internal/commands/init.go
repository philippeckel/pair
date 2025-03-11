package commands

import (
	"encoding/json"
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/spf13/cobra"
	"os"
)

func initConfig(cmd *cobra.Command, args []string) {
	// Check if the config file already exists
	if _, err := os.Stat(config.ConfigPath); err == nil {
		fmt.Printf("Config file already exists at %s. Use --config to specify a different path.\n", config.ConfigPath)
		return
	}

	// Create sample config with the new format
	sampleConfig := struct {
		CoAuthorsMap map[string]struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"coauthors"`
	}{
		CoAuthorsMap: map[string]struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			"jane": {
				Name:  "Jane Doe",
				Email: "jane.doe@example.com",
			},
			"john": {
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
		},
	}

	data, err := json.MarshalIndent(sampleConfig, "", "  ")
	if err != nil {
		fmt.Printf("Error creating sample config: %v\n", err)
		return
	}

	if err := os.WriteFile(config.ConfigPath, data, 0644); err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

	fmt.Printf("Created sample config file at %s\n", config.ConfigPath)
	fmt.Println("You can now use aliases to add co-authors, e.g.:")
	fmt.Println("  pair add john")
	fmt.Println("  pair add jane")
	fmt.Println("Or use interactive selection with:")
	fmt.Println("  pair select")
}
