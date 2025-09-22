package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/easel/ddx/internal/config"
	"github.com/spf13/cobra"
)

// Command registration is now handled by command_factory.go
// This file only contains the prompts subcommand implementations

// runPromptsList implements the prompts list command
func runPromptsList(cmd *cobra.Command, args []string) error {
	// Get library path
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	promptsDir := filepath.Join(libPath, "prompts")

	// Check if prompts directory exists
	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		fmt.Fprintln(cmd.OutOrStdout(), "No prompts directory found")
		return nil
	}

	// Get search filter
	searchFilter, _ := cmd.Flags().GetString("search")
	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Fprintln(cmd.OutOrStdout(), "Available prompts:")
	fmt.Fprintln(cmd.OutOrStdout())

	// Walk through prompts directory
	err = filepath.Walk(promptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == promptsDir {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(promptsDir, path)
		if err != nil {
			return err
		}

		// Skip hidden files
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Apply search filter if provided
		if searchFilter != "" && !strings.Contains(strings.ToLower(path), strings.ToLower(searchFilter)) {
			return nil
		}

		// Print directories and markdown files
		if info.IsDir() {
			fmt.Fprintf(cmd.OutOrStdout(), "📁 %s/\n", relPath)
		} else if strings.HasSuffix(path, ".md") {
			// Show full filename with extension in verbose mode
			if verbose {
				fmt.Fprintf(cmd.OutOrStdout(), "  📝 %s\n", relPath)
			} else {
				// Remove .md extension for display
				name := strings.TrimSuffix(relPath, ".md")
				fmt.Fprintf(cmd.OutOrStdout(), "  📝 %s\n", name)
			}
		}

		return nil
	})

	return err
}

// runPromptsShow implements the prompts show command
func runPromptsShow(cmd *cobra.Command, args []string) error {
	promptName := args[0]

	// Get library path
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	// Try different paths for the prompt
	possiblePaths := []string{
		filepath.Join(libPath, "prompts", promptName+".md"),
		filepath.Join(libPath, "prompts", promptName),
		filepath.Join(libPath, "prompts", promptName, "README.md"),
	}

	var promptPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			promptPath = path
			break
		}
	}

	if promptPath == "" {
		return fmt.Errorf("prompt not found: %s", promptName)
	}

	// Read and display the prompt
	content, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Errorf("failed to read prompt: %w", err)
	}

	fmt.Fprint(cmd.OutOrStdout(), string(content))
	return nil
}
