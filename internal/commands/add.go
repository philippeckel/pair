package commands

import (
	"fmt"
	"github.com/philippeckel/pair/internal/config"
	"github.com/philippeckel/pair/internal/gittemplate"
	"github.com/spf13/cobra"
	"strings"
)

func addCoAuthor(cmd *cobra.Command, args []string) error {
	if err := config.LoadConfig(); err != nil {
		return err
	}

	// Get current git user name and email for self-check
	userName, userEmail, err := getGitUserInfo()
	if err != nil {
		return fmt.Errorf("failed to get git user info: %w", err)
	}

	// Get current template
	templatePath, err := gittemplate.GetCurrentTemplate()
	if err != nil {
		return err
	}

	// Get active co-authors
	activeCoAuthors, err := gittemplate.ParseActiveCoAuthors(templatePath)
	if err != nil {
		return err
	}

	// Track if any co-authors were successfully added
	added := false
	// Track warnings to display at the end
	var warnings []string

	// Process each co-author identifier provided in args
	for _, identifier := range args {
		coAuthor, _, err := findCoAuthorByAliasOrIndex(identifier)
		if err != nil {
			// Return error for non-existent co-authors
			return fmt.Errorf("error with '%s': %w", identifier, err)
		}

		// Check if attempting to add yourself as co-author
		if strings.EqualFold(coAuthor.Email, userEmail) || strings.EqualFold(coAuthor.Name, userName) {
			warnings = append(warnings, fmt.Sprintf("Cannot add yourself as a co-author: %s <%s>", coAuthor.Name, coAuthor.Email))
			continue
		}

		// Check if co-author is already active
		alreadyActive := false
		for _, active := range activeCoAuthors {
			if active.Email == coAuthor.Email {
				alreadyActive = true
				warnings = append(warnings, fmt.Sprintf("Co-author already active: %s <%s>", coAuthor.Name, coAuthor.Email))
				break
			}
		}

		// Add co-author if not already active
		if !alreadyActive {
			activeCoAuthors = append(activeCoAuthors, coAuthor)
			fmt.Printf("Adding co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
			added = true
		}
	}

	// Display all collected warnings
	for _, warning := range warnings {
		fmt.Println(warning)
	}

	// Only update the template if at least one co-author was added
	if added {
		if err := gittemplate.UpdateTemplate(activeCoAuthors); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no co-authors were added")
	}

	return nil
}
