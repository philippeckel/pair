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
	Use:     "add [alias or index]",
	Short:   "Add a co-author to Git commits by alias or index",
	Args:    cobra.ExactArgs(1),
	RunE:    addCoAuthor,
	Aliases: []string{"a"},
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

		// Get flag value
		multi, _ := cmd.Flags().GetBool("multiple")

		if multi {
			// Multiple selection mode
			coAuthors, err := selectMultipleCoAuthors(config.Config.CoAuthors, activeCoAuthors)
			if err != nil {
				return err
			}

			// Add the selected co-authors
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
					fmt.Printf("Adding co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
				}
			}

			// Update git template
			if err := gittemplate.UpdateTemplate(activeCoAuthors); err != nil {
				return err
			}

			fmt.Printf("Added co-authors to your commit template\n")
		} else {
			// Single selection mode
			coAuthor, err := selectCoAuthorWithFzf(config.Config.CoAuthors, activeCoAuthors)
			if err != nil {
				return err
			}

			// Add the selected co-author
			activeCoAuthors = append(activeCoAuthors, coAuthor)

			// Update git template
			if err := gittemplate.UpdateTemplate(activeCoAuthors); err != nil {
				return err
			}

			fmt.Printf("Added co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
		}

		return nil
	},
	Aliases: []string{"s"},
}

func init() {
	// Add the multiple flag to select command
	selectCmd.Flags().BoolP("multiple", "m", false, "Enable multiple selection mode")
	// Make sure unselectCmd has the flag too
	unselectCmd.Flags().BoolP("multiple", "m", false, "Enable multiple selection mode")
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
