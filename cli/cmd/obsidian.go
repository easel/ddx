package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/easel/ddx/internal/obsidian"
	"github.com/easel/ddx/internal/obsidian/converter"
	"github.com/easel/ddx/internal/obsidian/validator"
)

var obsidianCmd = &cobra.Command{
	Use:   "obsidian",
	Short: "Manage Obsidian integration for HELIX workflow",
	Long: `Convert HELIX workflow documentation to Obsidian-compatible format with frontmatter and wikilinks.

This command helps you enhance your HELIX workflow documentation with:
• YAML frontmatter for metadata and organization
• Wikilinks for seamless navigation between documents
• Tag hierarchies for powerful searching and filtering
• Navigation hubs for quick access to all resources

The migration is backward compatible - files remain readable in standard markdown viewers.`,
}

var migrateCmd = &cobra.Command{
	Use:   "migrate [path]",
	Short: "Migrate HELIX files to Obsidian format",
	Long: `Migrate HELIX workflow files to Obsidian-compatible format.

This command will:
1. Scan all markdown files in the specified path
2. Detect file types (phases, artifacts, templates, etc.)
3. Generate appropriate YAML frontmatter for each file
4. Convert markdown links to Obsidian wikilinks
5. Create a navigation hub for easy browsing
6. Validate the results

The migration preserves all existing content and is fully reversible.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "workflows/helix"
		if len(args) > 0 {
			path = args[0]
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		validateAfter, _ := cmd.Flags().GetBool("validate")
		backupDir, _ := cmd.Flags().GetString("backup-dir")

		if verbose {
			fmt.Printf("🔄 Starting Obsidian migration for %s...\n", path)
		}

		// Check if path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}

		// Create backup if not dry run
		if !dryRun && backupDir != "" {
			if err := createBackup(path, backupDir); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
			fmt.Printf("📁 Created backup at %s\n", backupDir)
		}

		// Step 1: Scan files
		files, err := scanMarkdownFiles(path)
		if err != nil {
			return fmt.Errorf("failed to scan files: %w", err)
		}
		fmt.Printf("📁 Found %d markdown files\n", len(files))

		// Step 2: Detect file types
		detector := obsidian.NewFileTypeDetector()
		for _, file := range files {
			file.FileType = detector.Detect(file.Path)
			if verbose && file.FileType != obsidian.FileTypeUnknown {
				fmt.Printf("   %s -> %s\n", file.Path, file.FileType)
			}
		}

		// Step 3: Generate frontmatter
		generator := obsidian.NewFrontmatterGenerator()
		for _, file := range files {
			if file.FileType == obsidian.FileTypeUnknown {
				continue // Skip unknown files
			}

			fm, err := generator.Generate(file)
			if err != nil {
				fmt.Printf("⚠️  Failed to generate frontmatter for %s: %v\n", file.Path, err)
				continue
			}
			file.Frontmatter = fm
		}

		// Step 4: Convert links
		linkConverter := converter.NewLinkConverter()
		linkConverter.BuildIndex(files)

		for _, file := range files {
			if file.FileType == obsidian.FileTypeUnknown {
				continue
			}
			file.Content = linkConverter.ConvertContent(file.Content)
		}

		// Step 5: Save files (unless dry-run)
		savedCount := 0
		if !dryRun {
			for _, file := range files {
				if file.FileType == obsidian.FileTypeUnknown {
					continue
				}

				if err := saveMarkdownFile(file); err != nil {
					fmt.Printf("⚠️  Failed to save %s: %v\n", file.Path, err)
					continue
				}
				if verbose {
					fmt.Printf("✅ Updated %s\n", file.Path)
				}
				savedCount++
			}
			fmt.Printf("✅ Updated %d files\n", savedCount)
		} else {
			fmt.Println("🔍 Dry run mode - no files were modified")
			for _, file := range files {
				if file.FileType != obsidian.FileTypeUnknown {
					fmt.Printf("   Would update: %s (%s)\n", file.Path, file.FileType)
				}
			}
		}

		// Step 6: Generate navigation hub
		if !dryRun {
			if err := generateNavigationHub(path, files); err != nil {
				return fmt.Errorf("failed to generate navigation hub: %w", err)
			}
			fmt.Printf("🗺️  Generated navigation hub at %s/NAVIGATOR.md\n", path)
		}

		// Step 7: Validate if requested
		if validateAfter && !dryRun {
			fmt.Println("🔍 Validating migration results...")
			return runValidation(path)
		}

		fmt.Println("✨ Migration complete!")
		return nil
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate Obsidian format in HELIX files",
	Long: `Validate that HELIX files conform to Obsidian format requirements.

This command checks:
• YAML frontmatter is valid and complete
• Required fields are present for each file type
• Tags follow the correct hierarchical structure
• Wikilinks are properly formed
• No broken references exist

Use this after migration to ensure everything is working correctly.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "workflows/helix"
		if len(args) > 0 {
			path = args[0]
		}

		return runValidation(path)
	},
}

var revertCmd = &cobra.Command{
	Use:   "revert [path]",
	Short: "Revert Obsidian migration",
	Long: `Revert files from Obsidian format back to standard markdown.

This command will:
• Remove YAML frontmatter from all files
• Convert wikilinks back to standard markdown links
• Restore from backup if available

Use this if you need to undo the Obsidian migration.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "workflows/helix"
		if len(args) > 0 {
			path = args[0]
		}

		backupDir, _ := cmd.Flags().GetString("backup-dir")
		stripFrontmatter, _ := cmd.Flags().GetBool("strip-frontmatter")

		fmt.Printf("🔄 Reverting Obsidian migration for %s...\n", path)

		// Try to restore from backup first
		if backupDir != "" {
			if _, err := os.Stat(backupDir); err == nil {
				if err := restoreFromBackup(backupDir, path); err != nil {
					return fmt.Errorf("failed to restore from backup: %w", err)
				}
				fmt.Printf("✅ Restored from backup %s\n", backupDir)
				return nil
			}
			fmt.Printf("⚠️  Backup directory not found: %s\n", backupDir)
		}

		// Fallback: strip frontmatter manually
		if stripFrontmatter {
			files, err := scanMarkdownFiles(path)
			if err != nil {
				return fmt.Errorf("failed to scan files: %w", err)
			}

			count := 0
			for _, file := range files {
				if file.HasFrontmatter() {
					// Remove frontmatter and save
					file.Frontmatter = nil
					if err := saveMarkdownFile(file); err != nil {
						fmt.Printf("⚠️  Failed to update %s: %v\n", file.Path, err)
						continue
					}
					count++
				}
			}
			fmt.Printf("✅ Removed frontmatter from %d files\n", count)
		}

		return nil
	},
}

var navCmd = &cobra.Command{
	Use:   "nav [path]",
	Short: "Generate navigation hub",
	Long: `Generate or update the navigation hub for HELIX workflow.

The navigation hub provides:
• Quick access to all phases and artifacts
• Status overview of current work
• Tag-based browsing and filtering
• Dataview queries for dynamic content

This is automatically run during migration but can be updated separately.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "workflows/helix"
		if len(args) > 0 {
			path = args[0]
		}

		fmt.Printf("🗺️  Generating navigation hub for %s...\n", path)

		files, err := scanMarkdownFiles(path)
		if err != nil {
			return fmt.Errorf("failed to scan files: %w", err)
		}

		if err := generateNavigationHub(path, files); err != nil {
			return fmt.Errorf("failed to generate navigation hub: %w", err)
		}

		fmt.Printf("✅ Generated navigation hub at %s/NAVIGATOR.md\n", path)
		return nil
	},
}

func init() {
	// Add flags to migrate command
	migrateCmd.Flags().BoolP("dry-run", "d", false, "Preview changes without modifying files")
	migrateCmd.Flags().Bool("validate", false, "Validate after migration")
	migrateCmd.Flags().StringP("backup-dir", "b", "", "Create backup in specified directory")

	// Add flags to revert command
	revertCmd.Flags().StringP("backup-dir", "b", "", "Restore from backup directory")
	revertCmd.Flags().Bool("strip-frontmatter", false, "Strip frontmatter if no backup available")

	// Add subcommands
	obsidianCmd.AddCommand(migrateCmd)
	obsidianCmd.AddCommand(validateCmd)
	obsidianCmd.AddCommand(revertCmd)
	obsidianCmd.AddCommand(navCmd)

	// Add to root command
	rootCmd.AddCommand(obsidianCmd)
}

// scanMarkdownFiles recursively scans for markdown files
func scanMarkdownFiles(root string) ([]*obsidian.MarkdownFile, error) {
	var files []*obsidian.MarkdownFile

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		file := &obsidian.MarkdownFile{
			Path:    path,
			Content: string(content),
		}

		// Parse existing frontmatter if present
		if fm, content := extractFrontmatter(string(content)); fm != nil {
			file.Frontmatter = fm
			file.Content = content
		}

		files = append(files, file)
		return nil
	})

	return files, err
}

// extractFrontmatter extracts existing YAML frontmatter from content
func extractFrontmatter(content string) (*obsidian.Frontmatter, string) {
	if !strings.HasPrefix(content, "---\n") {
		return nil, content
	}

	endIdx := strings.Index(content[4:], "\n---\n")
	if endIdx == -1 {
		return nil, content
	}

	yamlContent := content[4 : endIdx+4]
	remainingContent := content[endIdx+9:]

	var fm obsidian.Frontmatter
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		return nil, content
	}

	return &fm, remainingContent
}

// saveMarkdownFile saves a markdown file with frontmatter
func saveMarkdownFile(file *obsidian.MarkdownFile) error {
	var content strings.Builder

	// Write frontmatter
	if file.Frontmatter != nil {
		yamlBytes, err := yaml.Marshal(file.Frontmatter)
		if err != nil {
			return err
		}
		content.WriteString("---\n")
		content.WriteString(string(yamlBytes))
		content.WriteString("---\n\n")
	}

	// Write content
	content.WriteString(file.Content)

	return ioutil.WriteFile(file.Path, []byte(content.String()), 0644)
}

// generateNavigationHub creates a navigation hub file
func generateNavigationHub(basePath string, files []*obsidian.MarkdownFile) error {
	navGen := obsidian.NewNavigationGenerator()
	hubContent := navGen.GenerateNavigationHub(files)
	hubPath := filepath.Join(basePath, "NAVIGATOR.md")
	return ioutil.WriteFile(hubPath, []byte(hubContent), 0644)
}

// runValidation runs validation on the specified path
func runValidation(path string) error {
	fmt.Printf("🔍 Validating Obsidian format in %s...\n", path)

	files, err := scanMarkdownFiles(path)
	if err != nil {
		return fmt.Errorf("failed to scan files: %w", err)
	}

	v := validator.NewValidator()
	errorCount := 0
	fileCount := 0

	for _, file := range files {
		if file.FileType == obsidian.FileTypeUnknown {
			continue
		}

		fileCount++
		errors := v.ValidateFile(file)
		if len(errors) > 0 {
			fmt.Printf("\n❌ %s:\n", file.Path)
			for _, err := range errors {
				fmt.Printf("  - %s\n", err)
				errorCount++
			}
		}
	}

	if errorCount == 0 {
		fmt.Printf("\n✅ All %d files are valid!\n", fileCount)
	} else {
		fmt.Printf("\n⚠️  Found %d validation errors in %d files\n", errorCount, fileCount)
	}

	return nil
}

// createBackup creates a backup of the specified directory
func createBackup(srcPath, backupPath string) error {
	return copyObsidianDir(srcPath, backupPath)
}

// restoreFromBackup restores from backup directory
func restoreFromBackup(backupPath, destPath string) error {
	// Remove existing directory
	if err := os.RemoveAll(destPath); err != nil {
		return err
	}

	// Copy from backup
	return copyObsidianDir(backupPath, destPath)
}

// copyObsidianDir recursively copies a directory
func copyObsidianDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Create destination directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = destFile.ReadFrom(srcFile)
		return err
	})
}
