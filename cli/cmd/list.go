package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/cobra"
)

var (
	listType    string
	listSearch  string
	listVerbose bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates, patterns, and configurations",
	Long: `List all available resources in the DDx toolkit.

You can filter by type or search for specific items.`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listType, "type", "t", "", "Filter by type (templates|patterns|configs|prompts|scripts)")
	listCmd.Flags().StringVarP(&listSearch, "search", "s", "", "Search for specific items")
	listCmd.Flags().BoolVar(&listVerbose, "verbose", false, "Show verbose output with additional details")
}

func runList(cmd *cobra.Command, args []string) error {
	ddxHome := getDDxHome()

	// Check if DDx is installed
	if _, err := os.Stat(ddxHome); os.IsNotExist(err) {
		cmd.PrintErr("❌ DDx not found. Please run the installation script first.")
		return nil
	}

	cmd.Println("📋 Available DDx Resources")
	cmd.Println()

	// Define resource types to list
	resourceTypes := []string{"templates", "patterns", "configs", "prompts", "scripts"}

	// Filter by type if specified (either via flag or argument)
	filterType := listType
	if len(args) > 0 {
		filterType = args[0]
	}

	if filterType != "" {
		resourceTypes = []string{filterType}
	}

	for _, resourceType := range resourceTypes {
		resourcePath := filepath.Join(ddxHome, resourceType)

		if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(resourcePath)
		if err != nil {
			continue
		}

		if len(entries) == 0 {
			continue
		}

		// Filter entries based on search term
		var filteredEntries []os.DirEntry
		for _, entry := range entries {
			if listSearch == "" || strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(listSearch)) {
				filteredEntries = append(filteredEntries, entry)
			}
		}

		if len(filteredEntries) == 0 {
			continue
		}

		// Print section header
		caser := cases.Title(language.English)
		cmd.Printf("%s:\n", caser.String(resourceType))

		for _, entry := range filteredEntries {
			itemPath := filepath.Join(resourcePath, entry.Name())

			// Get item info
			info := getResourceInfo(itemPath, entry)

			if entry.IsDir() {
				cmd.Printf("  📁 %s", entry.Name())
			} else {
				cmd.Printf("  📄 %s", entry.Name())
			}

			if info != "" {
				cmd.Printf(" - %s", info)
			}
			cmd.Println()
		}
		cmd.Println()
	}

	// Show usage examples
	cmd.Println("Usage examples:")
	cmd.Println("  ddx apply nextjs           # Apply Next.js template")
	cmd.Println("  ddx apply error-handling   # Apply error handling pattern")
	cmd.Println("  ddx list --type templates  # Show only templates")
	cmd.Println("  ddx list --search react    # Search for react-related items")

	return nil
}

// getResourceInfo returns descriptive information about a resource
func getResourceInfo(path string, entry os.DirEntry) string {
	// Try to read description from README or description file
	if entry.IsDir() {
		readmePath := filepath.Join(path, "README.md")
		if content, err := os.ReadFile(readmePath); err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					if len(line) > 60 {
						return line[:60] + "..."
					}
					return line
				}
			}
		}

		// Count items in directory
		if entries, err := os.ReadDir(path); err == nil {
			return fmt.Sprintf("%d items", len(entries))
		}
	} else {
		// For files, show size or type
		if stat, err := entry.Info(); err == nil {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			switch ext {
			case ".md":
				return "Markdown document"
			case ".yml", ".yaml":
				return "Configuration file"
			case ".sh":
				return "Shell script"
			case ".py":
				return "Python script"
			case ".js":
				return "JavaScript file"
			case ".go":
				return "Go source"
			default:
				if stat.Size() < 1024 {
					return fmt.Sprintf("%d bytes", stat.Size())
				} else if stat.Size() < 1024*1024 {
					return fmt.Sprintf("%.1f KB", float64(stat.Size())/1024)
				} else {
					return fmt.Sprintf("%.1f MB", float64(stat.Size())/(1024*1024))
				}
			}
		}
	}

	return ""
}
