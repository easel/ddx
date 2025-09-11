package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	listType   string
	listSearch string
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
}

func runList(cmd *cobra.Command, args []string) error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	gray := color.New(color.FgHiBlack)
	bold := color.New(color.Bold)

	ddxHome := getDDxHome()

	// Check if DDx is installed
	if _, err := os.Stat(ddxHome); os.IsNotExist(err) {
		color.Red("❌ DDx not found. Please run the installation script first.")
		return nil
	}

	cyan.Println("📋 Available DDx Resources")
	fmt.Println()

	// Define resource types to list
	resourceTypes := []string{"templates", "patterns", "configs", "prompts", "scripts"}

	// Filter by type if specified
	if listType != "" {
		resourceTypes = []string{listType}
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
		bold.Printf("%s:\n", strings.Title(resourceType))

		for _, entry := range filteredEntries {
			itemPath := filepath.Join(resourcePath, entry.Name())

			// Get item info
			info := getResourceInfo(itemPath, entry)

			if entry.IsDir() {
				green.Printf("  📁 %s", entry.Name())
			} else {
				green.Printf("  📄 %s", entry.Name())
			}

			if info != "" {
				gray.Printf(" - %s", info)
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Show usage examples
	gray.Println("Usage examples:")
	gray.Println("  ddx apply nextjs           # Apply Next.js template")
	gray.Println("  ddx apply error-handling   # Apply error handling pattern")
	gray.Println("  ddx list --type templates  # Show only templates")
	gray.Println("  ddx list --search react    # Search for react-related items")

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
