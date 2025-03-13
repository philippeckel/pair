package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	docOutputDir  string
	docOutputType string
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation",
	Long:  `Generate documentation for the pair command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Ensure output directory exists
		if err := os.MkdirAll(docOutputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
		rootCmd.DisableAutoGenTag = true

		switch strings.ToLower(docOutputType) {
		case "markdown", "md":
			return doc.GenMarkdownTree(rootCmd, docOutputDir)
		case "man":
			return doc.GenManTree(rootCmd, &doc.GenManHeader{
				Title:   "PAIR",
				Section: "1",
			}, docOutputDir)
		case "yaml":
			return doc.GenYamlTree(rootCmd, docOutputDir)
		case "rest":
			return doc.GenReSTTree(rootCmd, docOutputDir)
		default:
			return fmt.Errorf("unknown documentation type: %s (supported: markdown, man, yaml, rest)", docOutputType)
		}
	},
}

func init() {
	// Add the doc command to rootCmd
	rootCmd.AddCommand(docsCmd)

	// Set default output directory to ./doc
	defaultOutputDir := filepath.Join(".", "doc")
	// Add flags for customizing output
	docsCmd.Flags().StringVarP(&docOutputDir, "output-dir", "o", defaultOutputDir, "Directory to output documentation")
	docsCmd.Flags().StringVarP(&docOutputType, "type", "t", "markdown", "Documentation type: markdown, man, yaml, rest")

}
