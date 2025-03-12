package commands

import (
	"fmt"
	"os"

	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pair",
	Short: "Manage Git commit co-authors",
	Long:  `A CLI tool to manage co-authors for Git commits based on a JSON configuration file.`,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all available co-authors",
	Run:     listCoAuthors,
	Aliases: []string{"ls"},
}

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show currently active co-authors",
	Run:     showActiveCoAuthors,
	Aliases: []string{"s"},
}

var addCmd = &cobra.Command{
	Use:     "add [alias or index]...",
	Short:   "Add one or more co-authors to Git commits by alias or index",
	Args:    cobra.MinimumNArgs(1),
	RunE:    addCoAuthor,
	Aliases: []string{"a"},
	Example: "pair add jane john",
}

var removeCmd = &cobra.Command{
	Use:     "remove [alias or index]",
	Short:   "Remove a co-author from Git commits by alias or index",
	Args:    cobra.ExactArgs(1),
	RunE:    removeCoAuthor,
	Aliases: []string{"rm"},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all active co-authors",
	Run:   clearCoAuthors,
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initialize a new config file with sample co-authors",
	Run:     initConfig,
	Aliases: []string{"i"},
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Interactively select co-authors using fuzzy finder",
	Long:  `Use fuzzy finder to interactively select co-authors from your config`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.LoadConfig(); err != nil {
			return err
		}

		if len(config.Config.CoAuthors) == 0 {
			return fmt.Errorf("no co-authors found in config")
		}

		// Get active co-authors
		templatePath, err := gittemplate.GetCurrentTemplate()
		if err != nil {
			return err
		}

		activeCoAuthors, err := gittemplate.ParseActiveCoAuthors(templatePath)
		if err != nil {
			return err
		}

		// Always use multiple selection mode as default
		coAuthors, err := selectMultipleCoAuthors(config.Config.CoAuthors, activeCoAuthors)
		if err != nil {
			return err
		}

		// Add the selected co-authors
		added := false
		for _, coAuthor := range coAuthors {
			// Check if already active
			alreadyActive := false
			for _, active := range activeCoAuthors {
				if active.Email == coAuthor.Email {
					alreadyActive = true
					fmt.Printf("Co-author already active: %s <%s>\n", coAuthor.Name, coAuthor.Email)
					break
				}
			}

			if !alreadyActive {
				activeCoAuthors = append(activeCoAuthors, coAuthor)
				fmt.Printf("Added co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
				added = true
			}
		}

		// Only update the template if at least one co-author was added
		if added {
			if err := gittemplate.UpdateTemplate(activeCoAuthors); err != nil {
				return err
			}
		} else {
			fmt.Println("No new co-authors were added")
		}

		return nil
	},
	Aliases: []string{"s"},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&config.ConfigPath, "config", "c",
		config.GetConfigPath(), "config file path")

	// Add all subcommands
	rootCmd.AddCommand(listCmd, showCmd, addCmd, removeCmd, clearCmd, initCmd, selectCmd, unselectCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
