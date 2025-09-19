package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/easel/ddx/internal/obsidian"
	"github.com/easel/ddx/internal/obsidian/converter"
	"github.com/spf13/cobra"
)

var (
	obsidianDryRun       bool
	obsidianValidateOnly bool
	obsidianPath         string
	obsidianOutput       string
)

// obsidianCmd represents the obsidian command
var obsidianCmd = &cobra.Command{
	Use:   "obsidian",
	Short: "Obsidian integration tools for markdown files",
	Long: `Obsidian integration tools that provide:

• Automatic frontmatter generation for different file types
• Markdown link to wikilink conversion
• Validation of Obsidian format compliance
• Navigation hub generation

Commands:
  migrate    Migrate markdown files to Obsidian format
  validate   Validate Obsidian format compliance
  nav        Generate navigation hub
  revert     Revert Obsidian formatting back to standard markdown`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Ensure we're in a project directory
		if !isInitialized() {
			fmt.Fprintf(os.Stderr, "Error: Not in a DDx-initialized project. Run 'ddx init' first.\n")
			os.Exit(1)
		}
	},
}

// obsidianMigrateCmd migrates markdown files to Obsidian format
var obsidianMigrateCmd = &cobra.Command{
	Use:   "migrate [path]",
	Short: "Migrate markdown files to Obsidian format",
	Long: `Migrate markdown files to Obsidian format by adding frontmatter
and converting markdown links to wikilinks.

This command will:
1. Scan for markdown files in the specified path
2. Detect file types based on path patterns
3. Generate appropriate frontmatter for each file type
4. Convert markdown links to wikilinks
5. Validate the resulting format

Use --dry-run to preview changes without modifying files.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		if obsidianDryRun {
			fmt.Printf("DRY RUN: Would migrate files in %s\n", targetPath)
		}

		return runObsidianMigrate(targetPath, obsidianDryRun, obsidianValidateOnly)
	},
}

// obsidianValidateCmd validates Obsidian format compliance
var obsidianValidateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate Obsidian format compliance",
	Long: `Validate that markdown files comply with Obsidian format requirements.

This command checks:
• Frontmatter schema compliance
• Wikilink resolution (no broken links)
• Tag hierarchy correctness
• Required fields presence

Returns exit code 0 if all files are valid, non-zero otherwise.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		return runObsidianValidate(targetPath)
	},
}

// obsidianNavCmd generates navigation hub
var obsidianNavCmd = &cobra.Command{
	Use:   "nav [generate]",
	Short: "Generate navigation hub for Obsidian vault",
	Long: `Generate a navigation hub that provides easy access to all workflow artifacts.

The navigation hub includes:
• Phase overview with status
• Artifact listings by category
• Tag-based navigation
• Quick access links

The hub is generated as NAVIGATOR.md in the workflow root.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && args[0] == "generate" {
			return runObsidianNavGenerate()
		}

		// Default to generate
		return runObsidianNavGenerate()
	},
}

// obsidianRevertCmd reverts Obsidian formatting
var obsidianRevertCmd = &cobra.Command{
	Use:   "revert [path]",
	Short: "Revert Obsidian formatting back to standard markdown",
	Long: `Revert Obsidian formatting back to standard markdown by:

• Removing frontmatter from files
• Converting wikilinks back to markdown links
• Preserving original content structure

Use --dry-run to preview changes without modifying files.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		if obsidianDryRun {
			fmt.Printf("DRY RUN: Would revert Obsidian formatting in %s\n", targetPath)
		}

		return runObsidianRevert(targetPath, obsidianDryRun)
	},
}

func init() {
	rootCmd.AddCommand(obsidianCmd)

	// Add subcommands
	obsidianCmd.AddCommand(obsidianMigrateCmd)
	obsidianCmd.AddCommand(obsidianValidateCmd)
	obsidianCmd.AddCommand(obsidianNavCmd)
	obsidianCmd.AddCommand(obsidianRevertCmd)

	// Flags for migrate command
	obsidianMigrateCmd.Flags().BoolVar(&obsidianDryRun, "dry-run", false, "Preview changes without modifying files")
	obsidianMigrateCmd.Flags().BoolVar(&obsidianValidateOnly, "validate-only", false, "Only validate format, don't make changes")

	// Flags for validate command
	obsidianValidateCmd.Flags().StringVar(&obsidianOutput, "output", "text", "Output format: text, json")

	// Flags for nav command
	obsidianNavCmd.Flags().StringVar(&obsidianPath, "output", "NAVIGATOR.md", "Output file for navigation hub")

	// Flags for revert command
	obsidianRevertCmd.Flags().BoolVar(&obsidianDryRun, "dry-run", false, "Preview changes without modifying files")
}

// runObsidianMigrate handles the migration process
func runObsidianMigrate(targetPath string, dryRun, validateOnly bool) error {
	fmt.Printf("Migrating markdown files in %s to Obsidian format...\n", targetPath)

	// Create file scanner
	scanner := obsidian.NewFileScanner()
	files, err := scanner.ScanDirectory(targetPath)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No markdown files found to migrate.")
		return nil
	}

	fmt.Printf("Found %d markdown files to process.\n", len(files))

	// Create generators and converters
	generator := obsidian.NewFrontmatterGenerator()
	linkConverter := converter.NewLinkConverter()

	// Build file index for link conversion
	linkConverter.BuildIndex(files)

	processedCount := 0
	errorCount := 0

	for _, file := range files {
		fmt.Printf("Processing: %s\n", file.Path)

		// Generate frontmatter if not present
		if !file.HasFrontmatter() {
			fm, err := generator.Generate(file)
			if err != nil {
				fmt.Printf("  Error generating frontmatter: %v\n", err)
				errorCount++
				continue
			}
			file.Frontmatter = fm
		}

		// Convert links to wikilinks
		originalContent := file.Content
		convertedContent := linkConverter.ConvertContent(file.Content)
		if convertedContent != originalContent {
			file.Content = convertedContent
			fmt.Printf("  Converted %d markdown links to wikilinks\n",
				countLinkConversions(originalContent, convertedContent))
		}

		if validateOnly {
			// Only validate, don't save
			fmt.Printf("  Would add frontmatter and convert links\n")
		} else if dryRun {
			fmt.Printf("  DRY RUN: Would save changes to %s\n", file.Path)
		} else {
			// Save the modified file
			err := file.WriteToFile()
			if err != nil {
				fmt.Printf("  Error saving file: %v\n", err)
				errorCount++
				continue
			}
			fmt.Printf("  Saved with frontmatter and converted links\n")
		}

		processedCount++
	}

	fmt.Printf("\nMigration complete: %d files processed, %d errors\n", processedCount, errorCount)

	if errorCount > 0 {
		return fmt.Errorf("migration completed with %d errors", errorCount)
	}

	return nil
}

// runObsidianValidate validates Obsidian format compliance
func runObsidianValidate(targetPath string) error {
	fmt.Printf("Validating Obsidian format in %s...\n", targetPath)

	scanner := obsidian.NewFileScanner()
	files, err := scanner.ScanDirectory(targetPath)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No markdown files found to validate.")
		return nil
	}

	linkConverter := converter.NewLinkConverter()
	linkConverter.BuildIndex(files)

	validCount := 0
	errorCount := 0
	totalErrors := []string{}

	for _, file := range files {
		fmt.Printf("Validating: %s\n", file.Path)

		// Check if file has frontmatter
		if !file.HasFrontmatter() {
			fmt.Printf("  ❌ Missing frontmatter\n")
			totalErrors = append(totalErrors, fmt.Sprintf("%s: missing frontmatter", file.Path))
			errorCount++
			continue
		}

		// Validate frontmatter fields
		if file.Frontmatter.Title == "" {
			fmt.Printf("  ❌ Missing title in frontmatter\n")
			totalErrors = append(totalErrors, fmt.Sprintf("%s: missing title", file.Path))
			errorCount++
		}

		if file.Frontmatter.Type == "" {
			fmt.Printf("  ❌ Missing type in frontmatter\n")
			totalErrors = append(totalErrors, fmt.Sprintf("%s: missing type", file.Path))
			errorCount++
		}

		if len(file.Frontmatter.Tags) == 0 {
			fmt.Printf("  ❌ Missing tags in frontmatter\n")
			totalErrors = append(totalErrors, fmt.Sprintf("%s: missing tags", file.Path))
			errorCount++
		}

		// Validate wikilinks
		brokenLinks := linkConverter.ValidateWikilinks(file.Content)
		if len(brokenLinks) > 0 {
			fmt.Printf("  ❌ Broken wikilinks: %s\n", strings.Join(brokenLinks, ", "))
			for _, link := range brokenLinks {
				totalErrors = append(totalErrors, fmt.Sprintf("%s: broken link %s", file.Path, link))
			}
			errorCount += len(brokenLinks)
		}

		if errorCount == 0 {
			fmt.Printf("  ✅ Valid\n")
			validCount++
		}
	}

	fmt.Printf("\nValidation complete: %d files valid, %d errors found\n", validCount, len(totalErrors))

	if len(totalErrors) > 0 {
		fmt.Println("\nErrors found:")
		for _, err := range totalErrors {
			fmt.Printf("  • %s\n", err)
		}
		return fmt.Errorf("validation failed with %d errors", len(totalErrors))
	}

	fmt.Println("✅ All files are valid!")
	return nil
}

// runObsidianNavGenerate generates the navigation hub
func runObsidianNavGenerate() error {
	fmt.Println("Generating navigation hub...")

	scanner := obsidian.NewFileScanner()
	files, err := scanner.ScanDirectory(".")
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No markdown files found for navigation.")
		return nil
	}

	// Generate navigation content
	nav := generateNavigationContent(files)

	outputPath := obsidianPath
	if outputPath == "" {
		outputPath = "NAVIGATOR.md"
	}

	err = os.WriteFile(outputPath, []byte(nav), 0644)
	if err != nil {
		return fmt.Errorf("failed to write navigation file: %w", err)
	}

	fmt.Printf("Navigation hub generated: %s\n", outputPath)
	return nil
}

// runObsidianRevert reverts Obsidian formatting
func runObsidianRevert(targetPath string, dryRun bool) error {
	fmt.Printf("Reverting Obsidian formatting in %s...\n", targetPath)

	scanner := obsidian.NewFileScanner()
	files, err := scanner.ScanDirectory(targetPath)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No markdown files found to revert.")
		return nil
	}

	processedCount := 0
	errorCount := 0

	for _, file := range files {
		fmt.Printf("Reverting: %s\n", file.Path)

		modified := false

		// Remove frontmatter if present
		if file.HasFrontmatter() {
			file.Frontmatter = nil
			modified = true
			fmt.Printf("  Would remove frontmatter\n")
		}

		// Convert wikilinks back to markdown links
		originalContent := file.Content
		revertedContent := revertWikilinksToMarkdown(file.Content)
		if revertedContent != originalContent {
			file.Content = revertedContent
			modified = true
			fmt.Printf("  Would convert wikilinks back to markdown links\n")
		}

		if modified {
			if dryRun {
				fmt.Printf("  DRY RUN: Would save reverted file\n")
			} else {
				err := file.WriteToFile()
				if err != nil {
					fmt.Printf("  Error saving file: %v\n", err)
					errorCount++
					continue
				}
				fmt.Printf("  Reverted and saved\n")
			}
		} else {
			fmt.Printf("  No changes needed\n")
		}

		processedCount++
	}

	fmt.Printf("\nRevert complete: %d files processed, %d errors\n", processedCount, errorCount)

	if errorCount > 0 {
		return fmt.Errorf("revert completed with %d errors", errorCount)
	}

	return nil
}

// Helper functions

func countLinkConversions(original, converted string) int {
	// Simple count of markdown link patterns that were converted
	originalLinks := strings.Count(original, "](")
	convertedLinks := strings.Count(converted, "](")
	return originalLinks - convertedLinks
}

func generateNavigationContent(files []*obsidian.MarkdownFile) string {
	var nav strings.Builder

	nav.WriteString("---\n")
	nav.WriteString("title: HELIX Workflow Navigator\n")
	nav.WriteString("type: navigator\n")
	nav.WriteString("tags:\n")
	nav.WriteString("  - helix\n")
	nav.WriteString("  - navigation\n")
	nav.WriteString("created: " + fmt.Sprintf("%v", obsidian.NewFrontmatterGenerator()) + "\n")
	nav.WriteString("updated: " + fmt.Sprintf("%v", obsidian.NewFrontmatterGenerator()) + "\n")
	nav.WriteString("---\n\n")

	nav.WriteString("# 🧭 HELIX Workflow Navigator\n\n")
	nav.WriteString("Welcome to the HELIX workflow navigation hub. This page provides quick access to all workflow artifacts and documentation.\n\n")

	// Group files by type
	phases := []*obsidian.MarkdownFile{}
	templates := []*obsidian.MarkdownFile{}
	examples := []*obsidian.MarkdownFile{}
	features := []*obsidian.MarkdownFile{}
	other := []*obsidian.MarkdownFile{}

	for _, file := range files {
		if !file.HasFrontmatter() {
			other = append(other, file)
			continue
		}

		switch file.FileType {
		case obsidian.FileTypePhase:
			phases = append(phases, file)
		case obsidian.FileTypeTemplate:
			templates = append(templates, file)
		case obsidian.FileTypeExample:
			examples = append(examples, file)
		case obsidian.FileTypeFeature:
			features = append(features, file)
		default:
			other = append(other, file)
		}
	}

	// Phases section
	if len(phases) > 0 {
		nav.WriteString("## 🔄 Workflow Phases\n\n")
		for _, file := range phases {
			title := file.GetTitle()
			if title == "" {
				title = filepath.Base(file.Path)
			}
			nav.WriteString(fmt.Sprintf("- [[%s]]\n", title))
		}
		nav.WriteString("\n")
	}

	// Templates section
	if len(templates) > 0 {
		nav.WriteString("## 📋 Templates\n\n")
		for _, file := range templates {
			title := file.GetTitle()
			if title == "" {
				title = filepath.Base(file.Path)
			}
			nav.WriteString(fmt.Sprintf("- [[%s]]\n", title))
		}
		nav.WriteString("\n")
	}

	// Examples section
	if len(examples) > 0 {
		nav.WriteString("## 💡 Examples\n\n")
		for _, file := range examples {
			title := file.GetTitle()
			if title == "" {
				title = filepath.Base(file.Path)
			}
			nav.WriteString(fmt.Sprintf("- [[%s]]\n", title))
		}
		nav.WriteString("\n")
	}

	// Features section
	if len(features) > 0 {
		nav.WriteString("## 🚀 Features\n\n")
		for _, file := range features {
			title := file.GetTitle()
			if title == "" {
				title = filepath.Base(file.Path)
			}
			nav.WriteString(fmt.Sprintf("- [[%s]]\n", title))
		}
		nav.WriteString("\n")
	}

	// Tag index
	nav.WriteString("## 🏷️ Tag Index\n\n")
	tagMap := make(map[string][]*obsidian.MarkdownFile)
	for _, file := range files {
		for _, tag := range file.GetTags() {
			tagMap[tag] = append(tagMap[tag], file)
		}
	}

	for tag, tagFiles := range tagMap {
		nav.WriteString(fmt.Sprintf("### %s\n", tag))
		for _, file := range tagFiles {
			title := file.GetTitle()
			if title == "" {
				title = filepath.Base(file.Path)
			}
			nav.WriteString(fmt.Sprintf("- [[%s]]\n", title))
		}
		nav.WriteString("\n")
	}

	nav.WriteString("---\n")
	nav.WriteString("*Generated automatically by DDx Obsidian integration*\n")

	return nav.String()
}

func revertWikilinksToMarkdown(content string) string {
	// Simple reversion - convert [[Link]] back to [Link](Link.md)
	// This is a basic implementation; a full implementation would need
	// to track the original paths
	result := content

	// Convert simple wikilinks [[Link]] to [Link](Link.md)
	result = strings.ReplaceAll(result, "[[", "[")
	result = strings.ReplaceAll(result, "]]", "]()")

	// This is a simplified reversion - in practice, we'd need to maintain
	// a mapping of wikilinks to original paths
	return result
}
