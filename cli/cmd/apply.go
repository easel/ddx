package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/templates"
)

var (
	applyPath   string
	applyDryRun bool
)

var applyCmd = &cobra.Command{
	Use:   "apply <resource>",
	Short: "Apply a specific template, pattern, or configuration",
	Long: `Apply a DDx resource to your project.

Resources can be:
• Templates (complete project setups)
• Patterns (code examples and best practices)  
• Configurations (tool configs like ESLint, Prettier)
• Prompts (AI prompts and instructions)
• Scripts (automation and setup scripts)

Examples:
  ddx apply nextjs              # Apply Next.js template
  ddx apply error-handling      # Apply error handling patterns
  ddx apply prompts/claude      # Apply Claude AI prompts
  ddx apply scripts/hooks       # Install git hooks`,
	Args: cobra.ExactArgs(1),
	RunE: runApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)
	
	applyCmd.Flags().StringVarP(&applyPath, "path", "p", ".", "Target path for application")
	applyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false, "Show what would be applied without making changes")
}

func runApply(cmd *cobra.Command, args []string) error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	gray := color.New(color.FgHiBlack)

	resourceName := args[0]
	
	cyan.Printf("🎯 Applying resource: %s\n\n", resourceName)

	// Check if DDx is initialized
	if !isInitialized() {
		red.Println("❌ DDx not initialized. Run 'ddx init' first.")
		return nil
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	s := spinner.New(spinner.CharSets[14], 100)
	s.Prefix = "Loading resource... "
	s.Start()

	// Find the resource
	resourceInfo, err := findResource(resourceName)
	if err != nil {
		s.Stop()
		return err
	}

	if resourceInfo == nil {
		s.Stop()
		red.Printf("❌ Resource '%s' not found\n", resourceName)
		
		// Show available resources
		yellow.Println("\n💡 Available resources:")
		if err := showAvailableResources(); err != nil {
			return err
		}
		return nil
	}

	s.Suffix = fmt.Sprintf(" Found %s: %s", resourceInfo.Type, resourceInfo.Name)

	if applyDryRun {
		s.Stop()
		yellow.Println("🔍 Dry run mode - showing what would be applied:")
		return showDryRun(resourceInfo, applyPath, cfg)
	}

	// Apply the resource
	s.Suffix = fmt.Sprintf(" Applying %s...", resourceInfo.Name)
	
	if err := applyResource(resourceInfo, applyPath, cfg); err != nil {
		s.Stop()
		return fmt.Errorf("failed to apply resource: %w", err)
	}

	s.Stop()
	green.Printf("✅ Successfully applied %s!\n", resourceInfo.Name)

	// Show what was applied
	fmt.Println()
	gray.Printf("Applied %s: %s\n", resourceInfo.Type, resourceInfo.Name)
	if resourceInfo.Description != "" {
		gray.Printf("Description: %s\n", resourceInfo.Description)
	}
	gray.Printf("Target: %s\n", applyPath)

	return nil
}

// ResourceInfo represents information about a DDx resource
type ResourceInfo struct {
	Name        string
	Type        string
	Path        string
	Description string
}

// findResource locates a resource by name
func findResource(resourceName string) (*ResourceInfo, error) {
	// Resource directories to search
	resourceDirs := map[string]string{
		"templates": "templates",
		"patterns":  "patterns", 
		"configs":   "configs",
		"prompts":   "prompts",
		"scripts":   "scripts",
	}

	// First try exact match
	for resourceType, dir := range resourceDirs {
		resourcePath := filepath.Join(".ddx", dir, resourceName)
		if _, err := os.Stat(resourcePath); err == nil {
			return &ResourceInfo{
				Name: resourceName,
				Type: resourceType,
				Path: resourcePath,
				Description: getResourceDescription(resourcePath),
			}, nil
		}
	}

	// Try partial match
	for resourceType, dir := range resourceDirs {
		dirPath := filepath.Join(".ddx", dir)
		if entries, err := os.ReadDir(dirPath); err == nil {
			for _, entry := range entries {
				if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(resourceName)) {
					resourcePath := filepath.Join(dirPath, entry.Name())
					return &ResourceInfo{
						Name: entry.Name(),
						Type: resourceType,
						Path: resourcePath,
						Description: getResourceDescription(resourcePath),
					}, nil
				}
			}
		}
	}

	// Try nested search (e.g., "prompts/claude")
	if strings.Contains(resourceName, "/") {
		resourcePath := filepath.Join(".ddx", resourceName)
		if _, err := os.Stat(resourcePath); err == nil {
			parts := strings.Split(resourceName, "/")
			return &ResourceInfo{
				Name: parts[len(parts)-1],
				Type: parts[0],
				Path: resourcePath,
				Description: getResourceDescription(resourcePath),
			}, nil
		}
	}

	return nil, nil
}

// getResourceDescription tries to extract description from README or other files
func getResourceDescription(resourcePath string) string {
	// Try README.md first
	readmePath := filepath.Join(resourcePath, "README.md")
	if content, err := os.ReadFile(readmePath); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				if len(line) > 100 {
					return line[:100] + "..."
				}
				return line
			}
		}
	}

	// If it's a single file, return file type
	if stat, err := os.Stat(resourcePath); err == nil && !stat.IsDir() {
		ext := strings.ToLower(filepath.Ext(resourcePath))
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
			return "Go source file"
		default:
			return "File"
		}
	}

	return ""
}

// applyResource applies a resource to the target path
func applyResource(resource *ResourceInfo, targetPath string, cfg *config.Config) error {
	stat, err := os.Stat(resource.Path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		// Apply directory resource
		return applyDirectory(resource.Path, targetPath, cfg)
	} else {
		// Apply single file resource
		return applySingleFile(resource.Path, targetPath, cfg)
	}
}

// applyDirectory applies a directory resource
func applyDirectory(sourcePath, targetPath string, cfg *config.Config) error {
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from source
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		targetFile := filepath.Join(targetPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetFile, info.Mode())
		}

		// Process file with variable substitution
		return processFile(path, targetFile, cfg)
	})
}

// applySingleFile applies a single file resource
func applySingleFile(sourcePath, targetPath string, cfg *config.Config) error {
	// If target is a directory, use the source filename
	if stat, err := os.Stat(targetPath); err == nil && stat.IsDir() {
		filename := filepath.Base(sourcePath)
		targetPath = filepath.Join(targetPath, filename)
	}

	return processFile(sourcePath, targetPath, cfg)
}

// processFile processes a file with variable substitution
func processFile(sourcePath, targetPath string, cfg *config.Config) error {
	// Read source content
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Apply variable substitution
	processedContent := cfg.ReplaceVariables(string(content))

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Write processed content
	return os.WriteFile(targetPath, []byte(processedContent), 0644)
}

// showDryRun shows what would be applied without making changes
func showDryRun(resource *ResourceInfo, targetPath string, cfg *config.Config) error {
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)
	gray := color.New(color.FgHiBlack)

	blue.Printf("\n📋 Dry Run Results:\n\n")
	fmt.Printf("Resource: %s\n", resource.Name)
	fmt.Printf("Type: %s\n", resource.Type)
	fmt.Printf("Source: %s\n", resource.Path)
	fmt.Printf("Target: %s\n\n", targetPath)

	stat, err := os.Stat(resource.Path)
	if err != nil {
		return err
	}

	yellow.Println("Files that would be created/updated:")
	
	if stat.IsDir() {
		return filepath.Walk(resource.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(resource.Path, path)
			if err != nil {
				return err
			}

			if relPath != "." && !info.IsDir() {
				targetFile := filepath.Join(targetPath, relPath)
				
				// Check if file exists
				if _, err := os.Stat(targetFile); err == nil {
					gray.Printf("  ~ %s (would update)\n", targetFile)
				} else {
					gray.Printf("  + %s (would create)\n", targetFile)
				}
			}
			return nil
		})
	} else {
		filename := filepath.Base(resource.Path)
		targetFile := filepath.Join(targetPath, filename)
		
		if _, err := os.Stat(targetFile); err == nil {
			gray.Printf("  ~ %s (would update)\n", targetFile)
		} else {
			gray.Printf("  + %s (would create)\n", targetFile)
		}
	}

	return nil
}

// showAvailableResources shows available resources organized by type
func showAvailableResources() error {
	resourceDirs := []string{"templates", "patterns", "configs", "prompts", "scripts"}
	
	for _, dir := range resourceDirs {
		dirPath := filepath.Join(".ddx", dir)
		if entries, err := os.ReadDir(dirPath); err == nil && len(entries) > 0 {
			fmt.Printf("\n%s:\n", strings.Title(dir))
			for _, entry := range entries {
				fmt.Printf("  • %s\n", entry.Name())
			}
		}
	}
	
	return nil
}