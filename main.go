package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ktr0731/go-fuzzyfinder"
)

// CoAuthor represents a contributor that can be added to commits
type CoAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Alias string `json:"-"` // Not in JSON, filled from the map key
}

// Validate checks if the CoAuthor has valid fields
func (c *CoAuthor) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if c.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	if !strings.Contains(c.Email, "@") {
		return fmt.Errorf("email must contain '@' character")
	}

	return nil
}

// Config holds all available co-authors
type Config struct {
	CoAuthorsMap map[string]CoAuthor `json:"coauthors"`
	CoAuthors    []CoAuthor          `json:"-"` // This will be filled after loading
}

var (
	configPath string
	config     Config
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "pair",
		Short: "Manage Git commit co-authors",
		Long:  `A CLI tool to manage co-authors for Git commits based on a JSON configuration file.`,
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", getDefaultConfigPath(), "config file path")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all available co-authors",
		Run:   listCoAuthors,
		Aliases: []string{"ls"},
	}

	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show currently active co-authors",
		Run:   showActiveCoAuthors,
		Aliases: []string{"s"},
	}

	var addCmd = &cobra.Command{
		Use:   "add [alias or index]",
		Short: "Add a co-author to Git commits by alias or index",
		Args:  cobra.ExactArgs(1),
		RunE:  addCoAuthor,
		Aliases: []string{"a"},
	}

	var removeCmd = &cobra.Command{
		Use:   "remove [alias or index]",
		Short: "Remove a co-author from Git commits by alias or index",
		Args:  cobra.ExactArgs(1),
		RunE:  removeCoAuthor,
		Aliases: []string{"rm"},
	}

	var clearCmd = &cobra.Command{
		Use:   "clear",
		Short: "Clear all active co-authors",
		Run:   clearCoAuthors,
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new config file with sample co-authors",
		Run:   initConfig,
		Aliases: []string{"i"},
	}

	// Command for multi-select add
	var multiselectCmd = &cobra.Command{
		Use:   "multiselect",
		Short: "Interactively select multiple co-authors using fzf",
		Long:  `Use fuzzy finder (fzf) to interactively select multiple co-authors from your config`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			if len(config.CoAuthors) == 0 {
				return fmt.Errorf("no co-authors found in config")
			}

			// Get active co-authors
			templatePath, err := getCurrentTemplate()
			if err != nil {
				return err
			}

			activeCoAuthors, err := parseActiveCoAuthors(templatePath)
			if err != nil {
				return err
			}

			// Get co-authors using fzf
			selectedCoAuthors, err := selectMultipleCoAuthors(config.CoAuthors, activeCoAuthors)
			if err != nil {
				return err
			}

			// Add the selected co-authors
			for _, coAuthor := range selectedCoAuthors {
				activeCoAuthors = append(activeCoAuthors, coAuthor)
				fmt.Printf("Adding co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
			}

			// Update git template
			if err := updateTemplate(activeCoAuthors); err != nil {
				return err
			}

			fmt.Printf("Added %d co-authors to your commit template\n", len(selectedCoAuthors))
			return nil
		},
		Aliases: []string{"ms"},
	}

	// Command for multi-select remove
	var multiunselectCmd = &cobra.Command{
		Use:   "multiunselect",
		Short: "Interactively remove multiple co-authors using fzf",
		Long:  `Use fuzzy finder (fzf) to interactively remove multiple active co-authors`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get current template
			templatePath, err := getCurrentTemplate()
			if err != nil {
				return err
			}

			// Get active co-authors
			activeCoAuthors, err := parseActiveCoAuthors(templatePath)
			if err != nil {
				return err
			}

			if len(activeCoAuthors) == 0 {
				return fmt.Errorf("no active co-authors found")
			}

			// Select co-authors to remove
			toRemove, err := selectMultipleActiveCoAuthors(activeCoAuthors)
			if err != nil {
				return err
			}

			// Create new list without the removed co-authors
			var newActiveCoAuthors []CoAuthor
			for _, author := range activeCoAuthors {
				shouldRemove := false
				for _, remove := range toRemove {
					if author.Email == remove.Email {
						shouldRemove = true
						fmt.Printf("Removing co-author: %s <%s>\n", remove.Name, remove.Email)
						break
					}
				}
				if !shouldRemove {
					newActiveCoAuthors = append(newActiveCoAuthors, author)
				}
			}

			// Update git template
			if err := updateTemplate(newActiveCoAuthors); err != nil {
				return err
			}

			fmt.Printf("Removed %d co-authors from your commit template\n", len(toRemove))
			return nil
		},
	}

	// Update the select command in main()
	var selectCmd = &cobra.Command{
		Use:   "select",
		Short: "Interactively select a co-author using fzf",
		Long:  `Use fuzzy finder (fzf) to interactively select a co-author from your config`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(); err != nil {
				return err
			}

			if len(config.CoAuthors) == 0 {
				return fmt.Errorf("no co-authors found in config")
			}

			// Get active co-authors
			templatePath, err := getCurrentTemplate()
			if err != nil {
				return err
			}

			activeCoAuthors, err := parseActiveCoAuthors(templatePath)
			if err != nil {
				return err
			}

			// Get co-author using fzf
			coAuthor, err := selectCoAuthorWithFzf(config.CoAuthors, activeCoAuthors)
			if err != nil {
				return err
			}

			// Add the selected co-author
			activeCoAuthors = append(activeCoAuthors, coAuthor)

			// Update git template
			if err := updateTemplate(activeCoAuthors); err != nil {
				return err
			}

			fmt.Printf("Added co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
			return nil
		},
	}

	// Add a new command for interactively removing co-authors
	var unselectCmd = &cobra.Command{
		Use:   "unselect",
		Short: "Interactively remove a co-author using fzf",
		Long:  `Use fuzzy finder (fzf) to interactively remove an active co-author`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get current template
			templatePath, err := getCurrentTemplate()
			if err != nil {
				return err
			}

			// Get active co-authors
			activeCoAuthors, err := parseActiveCoAuthors(templatePath)
			if err != nil {
				return err
			}

			if len(activeCoAuthors) == 0 {
				return fmt.Errorf("no active co-authors found")
			}

			toRemove, err := selectActiveCoAuthorWithFzf(activeCoAuthors)
			if err != nil {
				return err
			}

			// Create new list without the removed co-author
			var newActiveCoAuthors []CoAuthor
			for _, author := range activeCoAuthors {
				if author.Email != toRemove.Email {
					newActiveCoAuthors = append(newActiveCoAuthors, author)
				}
			}

			// Update git template
			if err := updateTemplate(newActiveCoAuthors); err != nil {
				return err
			}

			fmt.Printf("Removed co-author: %s <%s>\n", toRemove.Name, toRemove.Email)
			return nil
		},
	}

	rootCmd.AddCommand(listCmd, showCmd, addCmd, removeCmd, clearCmd, initCmd, selectCmd, unselectCmd, multiselectCmd, multiunselectCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Add this function:
// selectActiveCoAuthorWithFzf presents a fuzzy finder interface for selecting an active co-author to remove
func selectActiveCoAuthorWithFzf(activeCoAuthors []CoAuthor) (CoAuthor, error) {
	// Check if fzf is installed
	if _, err := exec.LookPath("fzf"); err != nil {
		return CoAuthor{}, fmt.Errorf("fzf is not installed: %w", err)
	}

	if len(activeCoAuthors) == 0 {
		return CoAuthor{}, fmt.Errorf("no active co-authors to select from")
	}

	// Create input for fzf
	var input bytes.Buffer
	for i, author := range activeCoAuthors {
		input.WriteString(fmt.Sprintf("%d: %s <%s>\n", i, author.Name, author.Email))
	}

	// Create a temporary file to feed input to fzf
	tempFile, err := os.CreateTemp("", "pair-fzf-input")
	if err != nil {
		return CoAuthor{}, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the input data to the temp file
	if _, err := tempFile.Write(input.Bytes()); err != nil {
		tempFile.Close()
		return CoAuthor{}, fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempFile.Close()

	// Start fzf with input from the temp file
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cat %s | fzf --ansi --height 40%% --layout=reverse --prompt 'Select co-author to remove: '", tempFile.Name()))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	outputBytes, err := cmd.Output()

	// Handle command errors
	if err != nil {
		// User canceled selection (exit code 130)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return CoAuthor{}, fmt.Errorf("selection canceled")
		}
		return CoAuthor{}, fmt.Errorf("fzf selection failed: %w", err)
	}

	// Parse the selected line
	selection := strings.TrimSpace(string(outputBytes))
	if selection == "" {
		return CoAuthor{}, fmt.Errorf("no co-author selected")
	}

	// Extract index from selection
	var index int
	if _, err := fmt.Sscanf(selection, "%d:", &index); err != nil {
		return CoAuthor{}, fmt.Errorf("failed to parse selection: %w", err)
	}

	if index < 0 || index >= len(activeCoAuthors) {
		return CoAuthor{}, fmt.Errorf("invalid selection index")
	}

	return activeCoAuthors[index], nil
}

// saveConfig saves the config back to the file
func saveConfig() error {
	// Convert our internal structure back to the expected JSON format
	tempConfig := struct {
		CoAuthorsMap map[string]struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"coauthors"`
	}{
		CoAuthorsMap: make(map[string]struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}),
	}

	for alias, coauthor := range config.CoAuthorsMap {
		tempConfig.CoAuthorsMap[alias] = struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			Name:  coauthor.Name,
			Email: coauthor.Email,
		}
	}

	data, err := json.MarshalIndent(tempConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// selectCoAuthorWithFzf presents a fuzzy finder interface for selecting a co-author
func selectCoAuthorWithFzf(coAuthors []CoAuthor, activeCoAuthors []CoAuthor) (CoAuthor, error) {
	// Check if fzf is installed
	if _, err := exec.LookPath("fzf"); err != nil {
		return CoAuthor{}, fmt.Errorf("fzf is not installed: %w", err)
	}

	// Filter out already active co-authors
	var availableCoAuthors []CoAuthor
	for _, author := range coAuthors {
		isActive := false
		for _, active := range activeCoAuthors {
			if author.Email == active.Email {
				isActive = true
				break
			}
		}
		if !isActive {
			availableCoAuthors = append(availableCoAuthors, author)
		}
	}

	// Check if there are any available co-authors left
	if len(availableCoAuthors) == 0 {
		return CoAuthor{}, fmt.Errorf("all co-authors are already active")
	}

	// Create input for fzf
	var input bytes.Buffer
	for i, author := range availableCoAuthors {
		input.WriteString(fmt.Sprintf("%d: %s (%s) <%s>\n", i, author.Name, author.Alias, author.Email))
	}

	// Create a temporary file to feed input to fzf
	tempFile, err := os.CreateTemp("", "pair-fzf-input")
	if err != nil {
		return CoAuthor{}, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the input data to the temp file
	if _, err := tempFile.Write(input.Bytes()); err != nil {
		tempFile.Close()
		return CoAuthor{}, fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempFile.Close()

	// Start fzf with input from the temp file
	// cmd := exec.Command("sh", "-c", fmt.Sprintf("cat %s | fzf --ansi --height 40%% --layout=reverse --prompt 'Select co-author: '", tempFile.Name()))
	cmd := exec.Command(fmt.Sprintf("cat %s | fzf --ansi --height 40%% --layout=reverse --prompt 'Select co-author: '", tempFile.Name()))

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	outputBytes, err := cmd.Output()

	// Handle command errors
	if err != nil {
		// User canceled selection (exit code 130)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return CoAuthor{}, fmt.Errorf("selection canceled")
		}
		return CoAuthor{}, fmt.Errorf("fzf selection failed: %w", err)
	}

	// Parse the selected line
	selection := strings.TrimSpace(string(outputBytes))
	if selection == "" {
		return CoAuthor{}, fmt.Errorf("no co-author selected")
	}

	// Extract index from selection
	var index int
	if _, err := fmt.Sscanf(selection, "%d:", &index); err != nil {
		return CoAuthor{}, fmt.Errorf("failed to parse selection: %w", err)
	}

	if index < 0 || index >= len(availableCoAuthors) {
		return CoAuthor{}, fmt.Errorf("invalid selection index")
	}

	return availableCoAuthors[index], nil
}

func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home dir, use current directory
		return ".pair.json"
	}
	return filepath.Join(homeDir, ".pair.json")
}

// loadConfig loads and parses the config file with the new structure
func loadConfig() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	// Define a temporary struct for unmarshaling with the right structure
	var tempConfig struct {
		CoAuthorsMap map[string]struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"coauthors"`
	}

	if err := json.Unmarshal(data, &tempConfig); err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	// Convert the map to our internal structure
	config.CoAuthorsMap = make(map[string]CoAuthor)
	config.CoAuthors = make([]CoAuthor, 0, len(tempConfig.CoAuthorsMap))

	for alias, details := range tempConfig.CoAuthorsMap {
		coauthor := CoAuthor{
			Name:  details.Name,
			Email: details.Email,
			Alias: alias,
		}

		if err := coauthor.Validate(); err != nil {
			return fmt.Errorf("invalid co-author '%s': %w", alias, err)
		}

		config.CoAuthorsMap[alias] = coauthor
		config.CoAuthors = append(config.CoAuthors, coauthor)
	}

	return nil
}

// selectMultipleCoAuthors allows selecting multiple co-authors at once using fzf
func selectMultipleCoAuthors(coAuthors []CoAuthor, activeCoAuthors []CoAuthor) ([]CoAuthor, error) {
	// Check if fzf is installed
	if _, err := exec.LookPath("fzf"); err != nil {
		return nil, fmt.Errorf("fzf is not installed: %w", err)
	}

	// Filter out already active co-authors
	var availableCoAuthors []CoAuthor
	for _, author := range coAuthors {
		isActive := false
		for _, active := range activeCoAuthors {
			if author.Email == active.Email {
				isActive = true
				break
			}
		}
		if !isActive {
			availableCoAuthors = append(availableCoAuthors, author)
		}
	}

	// Check if there are any available co-authors left
	if len(availableCoAuthors) == 0 {
		return nil, fmt.Errorf("all co-authors are already active")
	}

	// Create input for fzf
	var input bytes.Buffer
	for i, author := range availableCoAuthors {
		input.WriteString(fmt.Sprintf("%d: %s (%s) <%s>\n", i, author.Name, author.Alias, author.Email))
	}

	// Create a temporary file to feed input to fzf
	tempFile, err := os.CreateTemp("", "pair-fzf-input")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the input data to the temp file
	if _, err := tempFile.Write(input.Bytes()); err != nil {
		tempFile.Close()
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempFile.Close()

	// Start fzf with input from the temp file and multi-select enabled
	cmd := exec.Command("sh", "-c", fmt.Sprintf("cat %s | fzf --ansi --multi --height 40%% --layout=reverse --prompt 'Select co-authors (TAB to select multiple): '", tempFile.Name()))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	outputBytes, err := cmd.Output()

	// Handle command errors
	if err != nil {
		// User canceled selection (exit code 130)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return nil, fmt.Errorf("selection canceled")
		}
		return nil, fmt.Errorf("fzf selection failed: %w", err)
	}

	// Parse the selected lines
	selection := strings.TrimSpace(string(outputBytes))
	if selection == "" {
		return nil, fmt.Errorf("no co-authors selected")
	}

	var selectedAuthors []CoAuthor
	lines := strings.Split(selection, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var index int
		if _, err := fmt.Sscanf(line, "%d:", &index); err != nil {
			return nil, fmt.Errorf("failed to parse selection: %w", err)
		}

		if index < 0 || index >= len(availableCoAuthors) {
			return nil, fmt.Errorf("invalid selection index")
		}

		selectedAuthors = append(selectedAuthors, availableCoAuthors[index])
	}

	return selectedAuthors, nil
}

func selectMultipleActiveCoAuthors(activeCoAuthors []CoAuthor) ([]CoAuthor, error) {
    if len(activeCoAuthors) == 0 {
        return nil, fmt.Errorf("no active co-authors to select from")
    }

    // Run the fuzzy finder and get selected indices
    indices, err := fuzzyfinder.FindMulti(
        activeCoAuthors,
        func(i int) string {
            return fmt.Sprintf("%s <%s>", activeCoAuthors[i].Name, activeCoAuthors[i].Email)
        },
        fuzzyfinder.WithPromptString("Select co-authors to remove (TAB to select multiple):"),
        fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
            if i == -1 {
                return ""
            }
            return fmt.Sprintf("Name: %s\nEmail: %s\n",
                activeCoAuthors[i].Name,
                activeCoAuthors[i].Email)
        }),
    )

    if err != nil {
        // User canceled selection
        if err == fuzzyfinder.ErrAbort {
            return nil, fmt.Errorf("selection canceled")
        }
        return nil, fmt.Errorf("fuzzy finder error: %w", err)
    }

    // Convert selected indices to co-authors
    var selectedCoAuthors []CoAuthor
    for _, idx := range indices {
        selectedCoAuthors = append(selectedCoAuthors, activeCoAuthors[idx])
    }

    return selectedCoAuthors, nil
}

func listCoAuthors(cmd *cobra.Command, args []string) {
	if err := loadConfig(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(config.CoAuthors) == 0 {
		fmt.Println("No co-authors found in config. Use 'pair init' to create a sample config.")
		return
	}

	fmt.Println("Available co-authors:")
	for i, author := range config.CoAuthors {
		fmt.Printf("%d: %s (%s) <%s>\n", i, author.Name, author.Alias, author.Email)
	}
}

func getCurrentTemplate() (string, error) {
	cmd := exec.Command("git", "config", "--get", "commit.template")
	output, err := cmd.Output()
	if err != nil {
		// It's ok if the template doesn't exist yet
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func parseActiveCoAuthors(templatePath string) ([]CoAuthor, error) {
	var activeCoAuthors []CoAuthor

	if templatePath == "" {
		return activeCoAuthors, nil
	}

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Co-authored-by:") {
			// Extract name and email from format: "Co-authored-by: Name <email>"
			authorInfo := strings.TrimSpace(strings.TrimPrefix(line, "Co-authored-by:"))
			parts := strings.Split(authorInfo, "<")
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				email := strings.TrimSpace(strings.TrimSuffix(parts[1], ">"))
				activeCoAuthors = append(activeCoAuthors, CoAuthor{Name: name, Email: email})
			}
		}
	}

	return activeCoAuthors, nil
}

func showActiveCoAuthors(cmd *cobra.Command, args []string) {
	templatePath, err := getCurrentTemplate()
	if err != nil {
		fmt.Printf("Error getting current git template: %v\n", err)
		return
	}

	if templatePath == "" {
		fmt.Println("No commit template is currently set. No active co-authors.")
		return
	}

	activeCoAuthors, err := parseActiveCoAuthors(templatePath)
	if err != nil {
		fmt.Printf("Error parsing active co-authors: %v\n", err)
		return
	}

	if len(activeCoAuthors) == 0 {
		fmt.Println("No active co-authors found.")
		return
	}

	fmt.Println("Active co-authors:")
	for i, author := range activeCoAuthors {
		fmt.Printf("%d: %s <%s>\n", i, author.Name, author.Email)
	}
}

// findCoAuthorByAlias finds a co-author by alias or index
func findCoAuthorByAliasOrIndex(identifier string) (CoAuthor, int, error) {
	// First try to parse as index
	var index int
	var author CoAuthor

	if _, err := fmt.Sscanf(identifier, "%d", &index); err == nil {
		// It's an index
		if index < 0 || index >= len(config.CoAuthors) {
			return author, -1, fmt.Errorf("invalid co-author index")
		}
		return config.CoAuthors[index], index, nil
	}

	// Then try as alias
	if coauthor, exists := config.CoAuthorsMap[identifier]; exists {
		// Find the index in the slice for compatibility
		for i, author := range config.CoAuthors {
			if author.Alias == identifier {
				return coauthor, i, nil
			}
		}
		return coauthor, -1, nil
	}

	return author, -1, fmt.Errorf("no co-author found with alias '%s'", identifier)
}

func updateTemplate(activeCoAuthors []CoAuthor) error {
	// Create a temporary file for the commit template
	tempDir := os.TempDir()
	templatePath := filepath.Join(tempDir, "git_commit_template")

	var content strings.Builder
	content.WriteString("\n\n") // Leave space for commit message
	content.WriteString("# Co-authors:\n")

	for _, author := range activeCoAuthors {
		content.WriteString(fmt.Sprintf("Co-authored-by: %s <%s>\n", author.Name, author.Email))
	}

	if err := os.WriteFile(templatePath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	// Set the commit.template configuration
	cmd := exec.Command("git", "config", "--global", "commit.template", templatePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set git commit template: %w", err)
	}

	return nil
}

func addCoAuthor(cmd *cobra.Command, args []string) error {
	if err := loadConfig(); err != nil {
		return err
	}

	identifier := args[0]
	coAuthor, _, err := findCoAuthorByAliasOrIndex(identifier)
	if err != nil {
		return err
	}

	// Validate the co-author
	if err := coAuthor.Validate(); err != nil {
		return fmt.Errorf("invalid co-author: %w", err)
	}

	// Get current template
	templatePath, err := getCurrentTemplate()
	if err != nil {
		return err
	}

	// Get active co-authors
	activeCoAuthors, err := parseActiveCoAuthors(templatePath)
	if err != nil {
		return err
	}

	// Check if co-author is already active
	for _, active := range activeCoAuthors {
		if active.Email == coAuthor.Email {
			return fmt.Errorf("co-author '%s' is already active", coAuthor.Name)
		}
	}

	// Add the co-author
	activeCoAuthors = append(activeCoAuthors, coAuthor)

	// Update template
	if err := updateTemplate(activeCoAuthors); err != nil {
		return err
	}

	fmt.Printf("Added co-author: %s <%s>\n", coAuthor.Name, coAuthor.Email)
	return nil
}

func removeCoAuthor(cmd *cobra.Command, args []string) error {
	if err := loadConfig(); err != nil {
		return err
	}

	// Get current template
	templatePath, err := getCurrentTemplate()
	if err != nil {
		return err
	}

	// Get active co-authors
	activeCoAuthors, err := parseActiveCoAuthors(templatePath)
	if err != nil {
		return err
	}

	if len(activeCoAuthors) == 0 {
		return fmt.Errorf("no active co-authors to remove")
	}

	identifier := args[0]
	var indexToRemove = -1

	// Try to parse as index in the active co-authors list
	var index int
	if _, err := fmt.Sscanf(identifier, "%d", &index); err == nil {
		if index >= 0 && index < len(activeCoAuthors) {
			indexToRemove = index
		}
	}

	// If not found by index, try finding by alias in the config
	if indexToRemove == -1 {
		coAuthor, _, err := findCoAuthorByAliasOrIndex(identifier)
		if err != nil {
			return err
		}

		// Find this co-author in the active list
		for i, active := range activeCoAuthors {
			if active.Email == coAuthor.Email {
				indexToRemove = i
				break
			}
		}

		if indexToRemove == -1 {
			return fmt.Errorf("co-author '%s' is not currently active", coAuthor.Name)
		}
	}

	// Remove the co-author
	removedAuthor := activeCoAuthors[indexToRemove]
	activeCoAuthors = append(activeCoAuthors[:indexToRemove], activeCoAuthors[indexToRemove+1:]...)

	// Update template
	if err := updateTemplate(activeCoAuthors); err != nil {
		return err
	}

	fmt.Printf("Removed co-author: %s <%s>\n", removedAuthor.Name, removedAuthor.Email)
	return nil
}

func clearCoAuthors(cmd *cobra.Command, args []string) {
	// Update with empty authors list
	if err := updateTemplate([]CoAuthor{}); err != nil {
		fmt.Printf("Error clearing co-authors: %v\n", err)
		return
	}

	fmt.Println("All co-authors have been cleared")
}

// initConfig creates a new config file with the sample data
func initConfig(cmd *cobra.Command, args []string) {
	// Check if the config file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at %s. Use --config to specify a different path.\n", configPath)
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
			"john": {
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			"jane": {
				Name:  "Jane Smith",
				Email: "jane.smith@example.com",
			},
			"alex": {
				Name:  "Alex Johnson",
				Email: "alex.johnson@example.com",
			},
		},
	}

	data, err := json.MarshalIndent(sampleConfig, "", "  ")
	if err != nil {
		fmt.Printf("Error creating sample config: %v\n", err)
		return
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

	fmt.Printf("Created sample config file at %s\n", configPath)
	fmt.Println("You can now use aliases to add co-authors, e.g.:")
	fmt.Println("  pair add john")
	fmt.Println("  pair add jane")
	fmt.Println("Or use interactive selection with:")
	fmt.Println("  pair select")
}
