package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/easel/ddx/internal/config"
	"github.com/spf13/cobra"
)

// Command registration is now handled by command_factory.go
// This file only contains the runList function implementation

// Resource represents a single asset for JSON output
type Resource struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Path        string   `json:"path"`
	IsDirectory bool     `json:"is_directory"`
	Size        int64    `json:"size,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// ListResponse represents the complete JSON response
type ListResponse struct {
	Resources []Resource     `json:"resources"`
	Summary   map[string]int `json:"summary"`
	Filter    string         `json:"filter,omitempty"`
	Type      string         `json:"type,omitempty"`
}

func runList(cmd *cobra.Command, args []string) error {
	// Get flag values
	filterValue, _ := cmd.Flags().GetString("filter")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	treeOutput, _ := cmd.Flags().GetBool("tree")

	// Get library path using the centralized library resolution
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	// Check if library exists
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		if jsonOutput {
			response := ListResponse{
				Resources: []Resource{},
				Summary:   map[string]int{},
			}
			jsonData, _ := json.MarshalIndent(response, "", "  ")
			cmd.Println(string(jsonData))
			return nil
		}
		cmd.PrintErr("❌ DDx library not found. Please check your configuration.")
		return nil
	}

	// Define resource types to list
	resourceTypes := []string{"templates", "patterns", "configs", "prompts", "scripts"}

	// Filter by type if specified via argument
	var filterType string
	if len(args) > 0 {
		filterType = args[0]
		resourceTypes = []string{filterType}
	}

	// Load configuration for resource filtering
	cfg, configErr := config.Load()
	var resourceFilter *config.ResourceFilterEngine
	if configErr == nil {
		resourceFilter = config.NewResourceFilterEngine(cfg, libPath)
	}

	// Collect all resources
	var allResources []Resource
	summary := make(map[string]int)

	for _, resourceType := range resourceTypes {
		var filteredPaths []string

		if resourceFilter != nil {
			// Use resource filtering if configuration is available
			discoveredPaths, err := resourceFilter.DiscoverResourcesWithFilter(resourceType)
			if err != nil {
				// Fall back to manual discovery if filtering fails
				discoveredPaths = discoverResourcesManually(libPath, resourceType)
			}
			filteredPaths = discoveredPaths
		} else {
			// Fall back to manual discovery if no config
			filteredPaths = discoverResourcesManually(libPath, resourceType)
		}

		var categoryResources []Resource
		for _, itemPath := range filteredPaths {
			// Apply additional text filter if specified
			if filterValue != "" {
				relPath := strings.TrimPrefix(itemPath, filepath.Join(libPath, resourceType)+"/")
				if !strings.Contains(strings.ToLower(relPath), strings.ToLower(filterValue)) {
					continue
				}
			}

			info, err := os.Stat(itemPath)
			if err != nil {
				continue
			}

			description := getResourceInfo(itemPath, &dirEntryWrapper{info})

			var size int64
			if !info.IsDir() {
				size = info.Size()
			}

			// Get relative name for display
			relPath := strings.TrimPrefix(itemPath, filepath.Join(libPath, resourceType)+"/")
			if relPath == itemPath {
				// Fallback to basename if prefix trimming didn't work
				relPath = filepath.Base(itemPath)
			}

			resource := Resource{
				Name:        relPath,
				Type:        resourceType,
				Description: description,
				Path:        itemPath,
				IsDirectory: info.IsDir(),
				Size:        size,
				Tags:        extractTags(itemPath, &dirEntryWrapper{info}),
			}

			categoryResources = append(categoryResources, resource)
		}

		if len(categoryResources) > 0 {
			allResources = append(allResources, categoryResources...)
			summary[resourceType] = len(categoryResources)
		}
	}

	// Output results
	if jsonOutput {
		response := ListResponse{
			Resources: allResources,
			Summary:   summary,
			Filter:    filterValue,
			Type:      filterType,
		}
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		cmd.Println(string(jsonData))
		return nil
	}

	// Tree output
	if treeOutput {
		return displayTreeOutput(cmd, allResources, filterValue)
	}

	// Human-readable output
	if len(allResources) == 0 {
		cmd.Println("📋 No DDx resources found")
		if filterValue != "" {
			cmd.Printf("No resources match filter: '%s'\n", filterValue)
		}
		return nil
	}

	cmd.Println("📋 Available DDx Resources")
	if filterValue != "" {
		cmd.Printf("Filtered by: '%s'\n", filterValue)
	}
	cmd.Println()

	// Show summary if listing all types
	if filterType == "" && len(summary) > 1 {
		cmd.Println("Summary:")
		caser := cases.Title(language.English)
		for resourceType, count := range summary {
			cmd.Printf("  %s: %d items\n", caser.String(resourceType), count)
		}
		cmd.Println()
	}

	// Group resources by type for display
	resourcesByType := make(map[string][]Resource)
	for _, resource := range allResources {
		resourcesByType[resource.Type] = append(resourcesByType[resource.Type], resource)
	}

	// Display each category
	caser := cases.Title(language.English)
	for _, resourceType := range resourceTypes {
		resources, exists := resourcesByType[resourceType]
		if !exists || len(resources) == 0 {
			continue
		}

		cmd.Printf("%s:\n", caser.String(resourceType))
		for _, resource := range resources {
			if resource.IsDirectory {
				cmd.Printf("  📁 %s", resource.Name)
			} else {
				cmd.Printf("  📄 %s", resource.Name)
			}

			if resource.Description != "" {
				cmd.Printf(" - %s", resource.Description)
			}
			cmd.Println()
		}
		cmd.Println()
	}

	// Show usage examples
	cmd.Println("Usage examples:")
	cmd.Println("  ddx apply nextjs           # Apply Next.js template")
	cmd.Println("  ddx apply error-handling   # Apply error handling pattern")
	cmd.Println("  ddx list templates         # Show only templates")
	cmd.Println("  ddx list --filter react    # Search for react-related items")
	cmd.Println("  ddx list --json            # Output as JSON")

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

// extractTags extracts tags from resource metadata or filename
func extractTags(path string, entry os.DirEntry) []string {
	var tags []string

	// Extract tags from filename patterns
	name := strings.ToLower(entry.Name())

	// Common technology tags
	if strings.Contains(name, "react") || strings.Contains(name, "jsx") {
		tags = append(tags, "react")
	}
	if strings.Contains(name, "vue") {
		tags = append(tags, "vue")
	}
	if strings.Contains(name, "angular") {
		tags = append(tags, "angular")
	}
	if strings.Contains(name, "nextjs") || strings.Contains(name, "next") {
		tags = append(tags, "nextjs")
	}
	if strings.Contains(name, "python") || strings.Contains(name, "py") {
		tags = append(tags, "python")
	}
	if strings.Contains(name, "go") || strings.Contains(name, "golang") {
		tags = append(tags, "go")
	}
	if strings.Contains(name, "javascript") || strings.Contains(name, "js") {
		tags = append(tags, "javascript")
	}
	if strings.Contains(name, "typescript") || strings.Contains(name, "ts") {
		tags = append(tags, "typescript")
	}
	if strings.Contains(name, "docker") {
		tags = append(tags, "docker")
	}
	if strings.Contains(name, "api") || strings.Contains(name, "rest") {
		tags = append(tags, "api")
	}
	if strings.Contains(name, "auth") {
		tags = append(tags, "authentication")
	}
	if strings.Contains(name, "test") {
		tags = append(tags, "testing")
	}
	if strings.Contains(name, "claude") || strings.Contains(name, "ai") {
		tags = append(tags, "ai")
	}

	return tags
}

// discoverResourcesManually discovers resources without filtering
func discoverResourcesManually(libPath, resourceType string) []string {
	var resources []string
	resourcePath := filepath.Join(libPath, resourceType)

	if _, err := os.Stat(resourcePath); os.IsNotExist(err) {
		return resources
	}

	filepath.Walk(resourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == resourcePath {
			return nil
		}
		resources = append(resources, path)
		return nil
	})

	return resources
}

// dirEntryWrapper wraps os.FileInfo to implement os.DirEntry interface
type dirEntryWrapper struct {
	info os.FileInfo
}

func (d *dirEntryWrapper) Name() string {
	return d.info.Name()
}

func (d *dirEntryWrapper) IsDir() bool {
	return d.info.IsDir()
}

func (d *dirEntryWrapper) Type() os.FileMode {
	return d.info.Mode().Type()
}

func (d *dirEntryWrapper) Info() (os.FileInfo, error) {
	return d.info, nil
}

// displayTreeOutput displays resources in tree format
func displayTreeOutput(cmd *cobra.Command, resources []Resource, filter string) error {
	if len(resources) == 0 {
		cmd.Println("📋 No DDx resources found")
		if filter != "" {
			cmd.Printf("No resources match filter: '%s'\n", filter)
		}
		return nil
	}

	cmd.Println("📋 DDx Resources Tree")
	if filter != "" {
		cmd.Printf("Filtered by: '%s'\n", filter)
	}
	cmd.Println()

	// Group resources by type
	resourcesByType := make(map[string][]Resource)
	for _, resource := range resources {
		resourcesByType[resource.Type] = append(resourcesByType[resource.Type], resource)
	}

	// Sort types
	types := []string{"prompts", "templates", "patterns", "configs", "scripts", "workflows"}
	for _, resourceType := range types {
		typeResources, exists := resourcesByType[resourceType]
		if !exists || len(typeResources) == 0 {
			continue
		}

		// Display type header
		cmd.Printf("📁 %s\n", resourceType)

		// Build tree structure for this type
		treeNodes := buildTreeStructure(typeResources)
		displayTreeNodes(cmd, treeNodes, "")
		cmd.Println()
	}

	return nil
}

// TreeNode represents a node in the resource tree
type TreeNode struct {
	Name        string
	IsDirectory bool
	Children    map[string]*TreeNode
	Resource    *Resource
}

// buildTreeStructure builds a tree structure from flat resource list
func buildTreeStructure(resources []Resource) map[string]*TreeNode {
	root := make(map[string]*TreeNode)

	for _, resource := range resources {
		parts := strings.Split(resource.Name, "/")
		current := root

		// Build path through tree
		for i, part := range parts {
			if current[part] == nil {
				current[part] = &TreeNode{
					Name:        part,
					IsDirectory: i < len(parts)-1 || resource.IsDirectory,
					Children:    make(map[string]*TreeNode),
				}
			}

			// If this is the last part, store the resource
			if i == len(parts)-1 {
				current[part].Resource = &resource
			}

			current = current[part].Children
		}
	}

	return root
}

// displayTreeNodes recursively displays tree nodes
func displayTreeNodes(cmd *cobra.Command, nodes map[string]*TreeNode, prefix string) {
	// Sort node names for consistent output
	var names []string
	for name := range nodes {
		names = append(names, name)
	}
	sort.Strings(names)

	for i, name := range names {
		node := nodes[name]
		isLast := i == len(names)-1

		// Choose tree characters
		var connector, nextPrefix string
		if isLast {
			connector = "└── "
			nextPrefix = prefix + "    "
		} else {
			connector = "├── "
			nextPrefix = prefix + "│   "
		}

		// Display node
		icon := "📄"
		if node.IsDirectory {
			icon = "📁"
		}

		displayName := name
		if node.Resource != nil && node.Resource.Description != "" {
			displayName += " - " + node.Resource.Description
		}

		cmd.Printf("%s%s%s %s\n", prefix, connector, icon, displayName)

		// Display children
		if len(node.Children) > 0 {
			displayTreeNodes(cmd, node.Children, nextPrefix)
		}
	}
}
