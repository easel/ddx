package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/briandowns/spinner"
	"github.com/easel/ddx/internal/config"
	"github.com/spf13/cobra"
)

// Command registration is now handled by command_factory.go
// This file only contains the runApply function implementation

func runApply(cmd *cobra.Command, args []string) error {
	// Get flag values locally
	applyPath, _ := cmd.Flags().GetString("path")
	applyDryRun, _ := cmd.Flags().GetBool("dry-run")
	applyVars, _ := cmd.Flags().GetStringSlice("var")

	resourceName := args[0]

	// Sanitize resource name for display to prevent command injection in output
	displayName := resourceName
	if strings.ContainsAny(displayName, ";|&$`") {
		displayName = "[sanitized]"
	}

	cmd.Printf("🎯 Applying resource: %s\n\n", displayName)

	// Check if we can load configuration (either local or have DDx home)
	// This allows the command to work if either DDx is installed globally or locally initialized

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
		cmd.PrintErrf("❌ Resource '%s' not found\n", displayName)

		// Show available resources
		cmd.PrintErr("\n💡 Available resources:")
		if err := showAvailableResources(); err != nil {
			return err
		}
		// Exit code 6: Resource not found
		return &ExitError{Code: 6, Message: fmt.Sprintf("resource '%s' not found", displayName)}
	}

	s.Suffix = fmt.Sprintf(" Found %s: %s", resourceInfo.Type, resourceInfo.Name)

	if applyDryRun {
		s.Stop()
		cmd.Println("🔍 Dry run mode - showing what would be applied:")
		return showDryRun(resourceInfo, applyPath, cfg, cmd)
	}

	// Apply the resource
	s.Suffix = fmt.Sprintf(" Applying %s...", resourceInfo.Name)

	if err := applyResource(resourceInfo, applyPath, cfg, applyVars); err != nil {
		s.Stop()
		return fmt.Errorf("failed to apply resource: %w", err)
	}

	s.Stop()
	cmd.Printf("✅ Successfully applied %s!\n", resourceInfo.Name)

	// Show what was applied
	cmd.Println()
	cmd.Printf("Applied %s: %s\n", resourceInfo.Type, resourceInfo.Name)
	if resourceInfo.Description != "" {
		cmd.Printf("Description: %s\n", resourceInfo.Description)
	}
	cmd.Printf("Target: %s\n", applyPath)

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
	// Get library path using the centralized resolution
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return nil, fmt.Errorf("failed to get library path: %w", err)
	}

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
		resourcePath := filepath.Join(libPath, dir, resourceName)
		if _, err := os.Stat(resourcePath); err == nil {
			return &ResourceInfo{
				Name:        resourceName,
				Type:        resourceType,
				Path:        resourcePath,
				Description: getResourceDescription(resourcePath),
			}, nil
		}
	}

	// Try partial match
	for resourceType, dir := range resourceDirs {
		dirPath := filepath.Join(libPath, dir)
		if entries, err := os.ReadDir(dirPath); err == nil {
			for _, entry := range entries {
				if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(resourceName)) {
					resourcePath := filepath.Join(dirPath, entry.Name())
					return &ResourceInfo{
						Name:        entry.Name(),
						Type:        resourceType,
						Path:        resourcePath,
						Description: getResourceDescription(resourcePath),
					}, nil
				}
			}
		}
	}

	// Try nested search (e.g., "prompts/claude")
	if strings.Contains(resourceName, "/") {
		resourcePath := filepath.Join(libPath, resourceName)
		if _, err := os.Stat(resourcePath); err == nil {
			parts := strings.Split(resourceName, "/")
			return &ResourceInfo{
				Name:        parts[len(parts)-1],
				Type:        parts[0],
				Path:        resourcePath,
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
func applyResource(resource *ResourceInfo, targetPath string, cfg *config.Config, vars []string) error {
	stat, err := os.Stat(resource.Path)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		// Apply directory resource
		return applyDirectory(resource.Path, targetPath, cfg, vars)
	} else {
		// Apply single file resource
		return applySingleFile(resource.Path, targetPath, cfg, vars)
	}
}

// applyDirectory applies a directory resource
func applyDirectory(sourcePath, targetPath string, cfg *config.Config, vars []string) error {
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
		return processFile(path, targetFile, cfg, vars)
	})
}

// applySingleFile applies a single file resource
func applySingleFile(sourcePath, targetPath string, cfg *config.Config, vars []string) error {
	// If target is a directory, use the source filename
	if stat, err := os.Stat(targetPath); err == nil && stat.IsDir() {
		filename := filepath.Base(sourcePath)
		targetPath = filepath.Join(targetPath, filename)
	}

	return processFile(sourcePath, targetPath, cfg, vars)
}

// processFile processes a file with variable substitution and conflict handling
func processFile(sourcePath, targetPath string, cfg *config.Config, vars []string) error {
	// Read source content
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Parse runtime variables from --var flags
	runtimeVars := parseRuntimeVariables(vars)

	// Create a copy of the config with runtime variables merged
	mergedConfig := cfg.WithRuntimeVariables(runtimeVars)

	// Apply variable substitution
	processedContent := mergedConfig.ReplaceVariables(string(content))

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Check for file conflicts
	if _, err := os.Stat(targetPath); err == nil {
		// File exists - handle conflict
		if err := handleFileConflict(targetPath, processedContent); err != nil {
			return err
		}
	} else {
		// File doesn't exist - create it
		if err := os.WriteFile(targetPath, []byte(processedContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

// handleFileConflict handles conflicts when a target file already exists
func handleFileConflict(targetPath, newContent string) error {
	// Read existing content
	existingContent, err := os.ReadFile(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read existing file %s: %w", targetPath, err)
	}

	// If content is identical, no conflict
	if string(existingContent) == newContent {
		return nil
	}

	// For now, create a backup and overwrite
	// TODO: In the future, this should prompt the user for choices
	backupPath := targetPath + ".ddx-backup"
	if err := os.WriteFile(backupPath, existingContent, 0644); err != nil {
		return fmt.Errorf("failed to create backup %s: %w", backupPath, err)
	}

	// Overwrite with new content
	if err := os.WriteFile(targetPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write new content to %s: %w", targetPath, err)
	}

	fmt.Printf("⚠️  File conflict resolved: %s\n", targetPath)
	fmt.Printf("   Created backup: %s\n", backupPath)

	return nil
}

// parseRuntimeVariables parses --var flags into a map
func parseRuntimeVariables(vars []string) map[string]string {
	result := make(map[string]string)
	for _, v := range vars {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// showDryRun shows what would be applied without making changes
func showDryRun(resource *ResourceInfo, targetPath string, cfg *config.Config, cmd *cobra.Command) error {
	cmd.Printf("\n📋 Dry Run Results:\n\n")
	cmd.Printf("Would apply: %s\n", resource.Name)
	cmd.Printf("Type: %s\n", resource.Type)
	cmd.Printf("Source: %s\n", resource.Path)
	cmd.Printf("Target: %s\n\n", targetPath)

	stat, err := os.Stat(resource.Path)
	if err != nil {
		return err
	}

	cmd.Println("Files that would be created/updated:")

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
					cmd.Printf("  ~ %s (would update)\n", targetFile)
				} else {
					cmd.Printf("  + %s (would create)\n", targetFile)
				}
			}
			return nil
		})
	} else {
		filename := filepath.Base(resource.Path)
		targetFile := filepath.Join(targetPath, filename)

		if _, err := os.Stat(targetFile); err == nil {
			cmd.Printf("  ~ %s (would update)\n", targetFile)
		} else {
			cmd.Printf("  + %s (would create)\n", targetFile)
		}
	}

	return nil
}

// showAvailableResources shows available resources organized by type
func showAvailableResources() error {
	// Get library path using the centralized resolution
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	resourceDirs := []string{"templates", "patterns", "configs", "prompts", "scripts"}

	for _, dir := range resourceDirs {
		dirPath := filepath.Join(libPath, dir)
		if entries, err := os.ReadDir(dirPath); err == nil && len(entries) > 0 {
			caser := cases.Title(language.English)
			fmt.Printf("\n%s:\n", caser.String(dir))
			for _, entry := range entries {
				fmt.Printf("  • %s\n", entry.Name())
			}
		}
	}

	return nil
}
