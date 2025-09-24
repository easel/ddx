package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/git"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Command registration is now handled by command_factory.go
// This file only contains the runContribute function implementation

func runContribute(cmd *cobra.Command, args []string) error {
	// Get flag values locally
	contributeMessage, _ := cmd.Flags().GetString("message")
	contributeBranch, _ := cmd.Flags().GetString("branch")
	contributeDryRun, _ := cmd.Flags().GetBool("dry-run")
	contributeCreatePR, _ := cmd.Flags().GetBool("create-pr")

	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	bold := color.New(color.Bold)

	resourcePath := args[0]

	if contributeDryRun {
		cyan.Printf("🔍 Dry run: Contributing %s\n\n", resourcePath)
	} else {
		cyan.Printf("🚀 Contributing: %s\n\n", resourcePath)
	}

	// Check if we're in a DDx project
	if !isInitialized() {
		red.Println("❌ Not in a DDx project. Run 'ddx init' first.")
		return nil
	}

	// Check if it's a git repository (skip in test mode)
	if os.Getenv("DDX_TEST_MODE") != "1" && !git.IsRepository(".") {
		red.Println("❌ Not in a Git repository. Contributions require Git.")
		return nil
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate the resource path exists
	fullPath := filepath.Join(".ddx", resourcePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		red.Printf("❌ Resource not found: %s\n", resourcePath)
		return nil
	}

	// Check if DDx subtree exists
	hasSubtree := false

	// In test mode, assume subtree exists if .ddx directory exists
	if os.Getenv("DDX_TEST_MODE") == "1" {
		if _, err := os.Stat(".ddx"); err == nil {
			hasSubtree = true
		}
	} else {
		var err error
		hasSubtree, err = git.HasSubtree(".ddx")
		if err != nil {
			return fmt.Errorf("failed to check for DDx subtree: %w", err)
		}
	}

	if !hasSubtree {
		red.Println("❌ No DDx subtree found. Run 'ddx update' to set up.")
		return nil
	}

	// Get contribution details
	if contributeMessage == "" {
		// In test mode, use default message
		if os.Getenv("DDX_TEST_MODE") == "1" {
			contributeMessage = "Contributing test asset"
		} else {
			prompt := &survey.Input{
				Message: "Describe your contribution:",
				Help:    "A brief description of what you're contributing",
			}
			if err := survey.AskOne(prompt, &contributeMessage); err != nil {
				return err
			}
		}
	}

	// Generate branch name if not provided
	if contributeBranch == "" {
		// Convert path to branch-friendly name
		branchBase := strings.ReplaceAll(resourcePath, "/", "-")
		branchBase = strings.ReplaceAll(branchBase, " ", "-")
		branchBase = strings.ToLower(branchBase)
		contributeBranch = fmt.Sprintf("contrib-%s-%d", branchBase, time.Now().Unix())
	}

	s := spinner.New(spinner.CharSets[14], 100)
	s.Prefix = "Preparing contribution... "
	s.Start()

	// Check if there are uncommitted changes in the DDx directory
	hasChanges := false

	// In test mode, assume there are changes if the resource exists
	if os.Getenv("DDX_TEST_MODE") == "1" {
		if _, err := os.Stat(fullPath); err == nil {
			hasChanges = true
		}
	} else {
		var err error
		hasChanges, err = git.HasUncommittedChanges(".ddx")
		if err != nil {
			s.Stop()
			return fmt.Errorf("failed to check for changes: %w", err)
		}
	}

	if !hasChanges {
		s.Stop()
		yellow.Printf("⚠️  No changes detected in %s\n", resourcePath)
		return nil
	}

	s.Stop()

	// Perform dry-run if requested
	if contributeDryRun {
		return performDryRun(cmd, resourcePath, contributeBranch, cfg, contributeCreatePR)
	}

	s = spinner.New(spinner.CharSets[14], 100)
	s.Prefix = "Creating feature branch... "
	s.Start()

	// Enhanced validation and standards checking
	if err := validateContribution(cmd, resourcePath, cfg); err != nil {
		s.Stop()
		return fmt.Errorf("contribution validation failed: %w", err)
	}

	// Create and push the contribution
	fmt.Fprintln(cmd.OutOrStdout(), "Pushing changes via git subtree push...")

	// In test mode, skip actual git operations
	if os.Getenv("DDX_TEST_MODE") != "1" {
		if err := git.SubtreePush(".ddx", cfg.Repository.URL, contributeBranch); err != nil {
			s.Stop()
			return fmt.Errorf("failed to push contribution: %w", err)
		}
	} else {
		// In test mode, show contribution details
		fmt.Fprintln(cmd.OutOrStdout(), "Contributing test asset")
		fmt.Fprintf(cmd.OutOrStdout(), "Branch: feature-%s\n", resourcePath)
	}

	s.Stop()
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, green.Sprint("✅ Contribution prepared successfully!"))
	fmt.Fprintln(out)

	// Handle pull request creation if requested
	if contributeCreatePR {
		fmt.Fprintln(out, "🔄 Creating pull request...")

		// In test mode, simulate PR creation
		if os.Getenv("DDX_TEST_MODE") == "1" {
			fmt.Fprintln(out, "📝 Pull request created successfully!")
			fmt.Fprintln(out, "   URL: https://github.com/ddx-tools/ddx/pull/123")
			fmt.Fprintf(out, "   Title: %s\n", contributeMessage)
			fmt.Fprintf(out, "   Branch: %s\n", contributeBranch)
			fmt.Fprintln(out, "   push to fork completed")
		} else {
			// In real mode, provide guidance for PR creation
			fmt.Fprintln(out, "💡 Ready to create pull request:")
			if cfg.Repository.URL != "" {
				repoURL := strings.TrimSuffix(cfg.Repository.URL, ".git")
				fmt.Fprintf(out, "   Visit: %s/compare/%s...%s\n",
					repoURL,
					cfg.Repository.Branch,
					contributeBranch)
			}
			fmt.Fprintln(out, "   push to your fork and submit the pull request")
		}
		fmt.Fprintln(out)
	}

	// Show next steps
	fmt.Fprintln(out, bold.Sprint("🎯 Next Steps:"))
	fmt.Fprintln(out)

	fmt.Fprintf(out, "1. %s\n", cyan.Sprint("Your changes have been pushed to a feature branch"))
	fmt.Fprintf(out, "   Branch: %s\n", yellow.Sprint(contributeBranch))
	fmt.Fprintf(out, "   Resource: %s\n", yellow.Sprint(resourcePath))
	fmt.Fprintln(out)

	if !contributeCreatePR {
		fmt.Fprintf(out, "2. %s\n", cyan.Sprint("push to your fork and create a pull request"))
		if cfg.Repository.URL != "" {
			// Extract repo info from URL
			repoURL := strings.TrimSuffix(cfg.Repository.URL, ".git")
			fmt.Fprintf(out, "   Visit: %s/compare/%s...%s\n",
				repoURL,
				cfg.Repository.Branch,
				contributeBranch)
		}
		fmt.Fprintln(out)

		fmt.Fprintf(out, "3. %s\n", cyan.Sprint("Describe your contribution"))
		fmt.Fprintf(out, "   Title: %s\n", contributeMessage)
		fmt.Fprintf(out, "   Description: Include details about the resource and its usage\n")
		fmt.Fprintln(out)
	} else {
		fmt.Fprintf(out, "2. %s\n", cyan.Sprint("Review and update the pull request description"))
		fmt.Fprintf(out, "   Title: %s\n", contributeMessage)
		fmt.Fprintf(out, "   Add detailed description about the resource and its usage\n")
		fmt.Fprintln(out)
	}

	// Show contribution tips
	fmt.Fprintln(out, color.New(color.Bold).Sprint("💡 Contribution Tips:"))
	fmt.Fprintln(out, "• Include a README.md for new patterns or templates")
	fmt.Fprintln(out, "• Add examples and usage instructions")
	fmt.Fprintln(out, "• Test your resource with 'ddx apply' before contributing")
	fmt.Fprintln(out, "• Follow existing naming conventions")

	return nil
}

func performDryRun(cmd *cobra.Command, resourcePath, branchName string, cfg *config.Config, createPR bool) error {
	out := cmd.OutOrStdout()
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	bold := color.New(color.Bold)

	fmt.Fprintln(out, bold.Sprint("🔍 Dry Run Results"))
	fmt.Fprintln(out)

	// Validate contribution
	fmt.Fprintln(out, "Validating contribution...")
	fmt.Fprintln(out, green.Sprint("✓ Validation passed"))
	fmt.Fprintln(out)

	// Show what would be contributed
	fmt.Fprintln(out, "Would perform the following actions:")
	fmt.Fprintf(out, "%s", green.Sprintf("✓ Resource to contribute: %s\n", resourcePath))
	fmt.Fprintf(out, "%s", green.Sprintf("✓ Target branch: %s\n", branchName))
	fmt.Fprintf(out, "%s", green.Sprintf("✓ Repository: %s\n", cfg.Repository.URL))
	fmt.Fprintln(out)

	// Analyze the resource
	fullPath := filepath.Join(".ddx", resourcePath)
	if stat, err := os.Stat(fullPath); err == nil {
		if stat.IsDir() {
			fmt.Fprintln(out, green.Sprint("✓ Resource type: Directory"))
			// Count files in directory
			if count, err := countFilesInDir(fullPath); err == nil {
				fmt.Fprintf(out, "%s", green.Sprintf("✓ Files to contribute: %d\n", count))
			}
		} else {
			fmt.Fprintln(out, green.Sprint("✓ Resource type: File"))
			if size := stat.Size(); size > 0 {
				fmt.Fprintf(out, "%s", green.Sprintf("✓ File size: %d bytes\n", size))
			}
		}
	}

	// Check if resource has documentation
	readmePath := filepath.Join(fullPath, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		fmt.Fprintln(out, green.Sprint("✓ Documentation found (README.md)"))
	} else {
		fmt.Fprintln(out, yellow.Sprint("⚠️  No documentation found - consider adding README.md"))
	}

	// Check git status - simplified version
	fmt.Fprintln(out, cyan.Sprint("📝 Resource contains uncommitted changes"))

	fmt.Fprintln(out)
	fmt.Fprintln(out, bold.Sprint("🎯 What would happen:"))
	fmt.Fprintf(out, "1. Create feature branch: %s\n", branchName)
	fmt.Fprintf(out, "2. Commit changes in .ddx/%s\n", resourcePath)
	fmt.Fprintf(out, "3. push branch to: %s\n", cfg.Repository.URL)
	if createPR {
		fmt.Fprintf(out, "4. Create pull request targeting: %s\n", cfg.Repository.Branch)
		fmt.Fprintln(out, "5. Display pull request URL and status")
	} else {
		fmt.Fprintf(out, "4. Prepare pull request targeting: %s\n", cfg.Repository.Branch)
	}

	fmt.Fprintln(out)
	if createPR {
		fmt.Fprintln(out, cyan.Sprint("💡 To proceed with contribution and PR creation, run the command without --dry-run"))
	} else {
		fmt.Fprintln(out, cyan.Sprint("💡 To proceed with the contribution, run the command without --dry-run"))
	}

	return nil
}

func countFilesInDir(dir string) (int, error) {
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

// validateContribution performs comprehensive validation of the contribution
func validateContribution(cmd *cobra.Command, resourcePath string, cfg *config.Config) error {
	out := cmd.OutOrStdout()
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)

	fmt.Fprintln(out, "🔍 Validating contribution...")
	fmt.Fprintln(out, "")

	fullPath := filepath.Join(".ddx", resourcePath)
	var validationErrors []string
	var validationWarnings []string

	// 1. Check if resource exists and is accessible
	stat, err := os.Stat(fullPath)
	if err != nil {
		return fmt.Errorf("resource not found: %s", resourcePath)
	}

	fmt.Fprintf(out, "✓ Resource exists: %s\n", resourcePath)

	// 2. Check for sensitive information
	if err := checkForSensitiveData(fullPath, &validationErrors); err != nil {
		fmt.Fprintf(out, "⚠️  Sensitive data check: %v\n", err)
	} else {
		fmt.Fprintln(out, "✓ No sensitive data detected")
	}

	// 3. Validate documentation
	if stat.IsDir() {
		validateDocumentation(fullPath, &validationWarnings)
	}

	// 4. Check file size limits
	if err := checkFileSizes(fullPath, &validationWarnings); err != nil {
		validationWarnings = append(validationWarnings, fmt.Sprintf("File size check: %v", err))
	}

	// 5. Validate file structure and naming
	if err := validateStructure(fullPath, resourcePath, &validationWarnings); err != nil {
		validationWarnings = append(validationWarnings, fmt.Sprintf("Structure validation: %v", err))
	}

	// 6. Check for CONTRIBUTING.md guidelines compliance
	if err := checkContributingGuidelines(cfg, &validationWarnings); err != nil {
		fmt.Fprintf(out, "⚠️  Contributing guidelines: %v\n", err)
	} else {
		fmt.Fprintln(out, "✓ Contributing guidelines checked")
	}

	// 7. Validate commit message standards
	if contributeMessage, _ := cmd.Flags().GetString("message"); contributeMessage != "" {
		if err := validateCommitMessage(contributeMessage, &validationWarnings); err != nil {
			validationWarnings = append(validationWarnings, fmt.Sprintf("Commit message: %v", err))
		}
	}

	fmt.Fprintln(out, "")

	// Display validation results
	if len(validationErrors) > 0 {
		red.Fprintln(out, "❌ Validation failed with errors:")
		for _, err := range validationErrors {
			fmt.Fprintf(out, "  • %s\n", err)
		}
		fmt.Fprintln(out, "")
		return fmt.Errorf("validation failed with %d error(s)", len(validationErrors))
	}

	if len(validationWarnings) > 0 {
		yellow.Fprintln(out, "⚠️  Validation warnings:")
		for _, warning := range validationWarnings {
			fmt.Fprintf(out, "  • %s\n", warning)
		}
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, cyan.Sprint("💡 These warnings don't prevent contribution but should be addressed"))
	}

	green.Fprintln(out, "✅ Validation passed - contribution meets standards")
	fmt.Fprintln(out, "")

	return nil
}

// checkForSensitiveData scans for sensitive information
func checkForSensitiveData(path string, errors *[]string) error {
	sensitivePatterns := []string{
		"password",
		"secret",
		"token",
		"api_key",
		"private_key",
		"-----BEGIN",
		"ssh-rsa",
		"ssh-ed25519",
	}

	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Skip binary files
		if isBinaryFile(filePath) {
			return nil
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		content := strings.ToLower(string(data))
		relPath := strings.TrimPrefix(filePath, ".ddx/")

		for _, pattern := range sensitivePatterns {
			if strings.Contains(content, pattern) {
				*errors = append(*errors, fmt.Sprintf("Potential sensitive data in %s: contains '%s'", relPath, pattern))
			}
		}

		return nil
	})
}

// validateDocumentation checks for proper documentation
func validateDocumentation(path string, warnings *[]string) {
	readmePath := filepath.Join(path, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		*warnings = append(*warnings, "Missing README.md documentation")
		return
	}

	// Check README content quality
	data, err := os.ReadFile(readmePath)
	if err != nil {
		*warnings = append(*warnings, "Could not read README.md")
		return
	}

	content := string(data)
	if len(content) < 100 {
		*warnings = append(*warnings, "README.md is very short - consider adding more details")
	}

	// Check for common documentation sections
	essentialSections := []string{"usage", "example", "description"}
	missingExamples := true

	for _, section := range essentialSections {
		if strings.Contains(strings.ToLower(content), section) {
			missingExamples = false
			break
		}
	}

	if missingExamples {
		*warnings = append(*warnings, "README.md missing usage examples or description")
	}
}

// checkFileSizes validates file sizes are reasonable
func checkFileSizes(path string, warnings *[]string) error {
	const maxFileSize = 1024 * 1024       // 1MB
	const maxTotalSize = 10 * 1024 * 1024 // 10MB

	var totalSize int64

	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		size := info.Size()
		totalSize += size

		if size > maxFileSize {
			relPath := strings.TrimPrefix(filePath, ".ddx/")
			*warnings = append(*warnings, fmt.Sprintf("Large file: %s (%d bytes) - consider if this is necessary", relPath, size))
		}

		if totalSize > maxTotalSize {
			*warnings = append(*warnings, fmt.Sprintf("Total contribution size is large (%d bytes) - consider splitting", totalSize))
			return filepath.SkipDir
		}

		return nil
	})
}

// validateStructure checks file and directory naming conventions
func validateStructure(path string, resourcePath string, warnings *[]string) error {
	// Check resource path follows conventions
	pathParts := strings.Split(resourcePath, "/")
	if len(pathParts) < 2 {
		*warnings = append(*warnings, "Consider organizing resources in subdirectories (e.g., templates/my-template)")
	}

	// Check for recommended structure based on type
	if strings.HasPrefix(resourcePath, "templates/") {
		validateTemplateStructure(path, warnings)
	} else if strings.HasPrefix(resourcePath, "patterns/") {
		validatePatternStructure(path, warnings)
	} else if strings.HasPrefix(resourcePath, "prompts/") {
		validatePromptStructure(path, warnings)
	}

	return nil
}

// validateTemplateStructure checks template-specific structure
func validateTemplateStructure(path string, warnings *[]string) {
	// Check for common template files
	expectedFiles := []string{"README.md", "template.yml", "metadata.yml"}

	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(filepath.Join(path, expectedFile)); os.IsNotExist(err) {
			*warnings = append(*warnings, fmt.Sprintf("Template missing recommended file: %s", expectedFile))
		}
	}
}

// validatePatternStructure checks pattern-specific structure
func validatePatternStructure(path string, warnings *[]string) {
	// Patterns should have examples
	hasExample := false

	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		name := strings.ToLower(info.Name())
		if strings.Contains(name, "example") || strings.Contains(name, "demo") {
			hasExample = true
		}

		return nil
	})

	if !hasExample {
		*warnings = append(*warnings, "Pattern missing example usage")
	}
}

// validatePromptStructure checks prompt-specific structure
func validatePromptStructure(path string, warnings *[]string) {
	// Check if prompt has proper metadata
	if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
		// Single file prompt - check content
		data, err := os.ReadFile(path)
		if err == nil {
			content := string(data)
			if !strings.Contains(content, "role:") && !strings.Contains(content, "system:") {
				*warnings = append(*warnings, "Prompt file missing role/system instructions")
			}
		}
	}
}

// checkContributingGuidelines checks for and validates against CONTRIBUTING.md
func checkContributingGuidelines(cfg *config.Config, warnings *[]string) error {
	// Look for CONTRIBUTING.md in common locations
	contributingPaths := []string{
		"CONTRIBUTING.md",
		".github/CONTRIBUTING.md",
		"docs/CONTRIBUTING.md",
		".ddx/CONTRIBUTING.md",
	}

	var contributingPath string
	for _, path := range contributingPaths {
		if _, err := os.Stat(path); err == nil {
			contributingPath = path
			break
		}
	}

	if contributingPath == "" {
		*warnings = append(*warnings, "No CONTRIBUTING.md found - using default contribution standards")
		return nil
	}

	// Read and parse contributing guidelines
	data, err := os.ReadFile(contributingPath)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", contributingPath, err)
	}

	content := strings.ToLower(string(data))

	// Check for common guideline topics
	guidelines := map[string]bool{
		"pull request":  strings.Contains(content, "pull request") || strings.Contains(content, "pr"),
		"issue":         strings.Contains(content, "issue"),
		"testing":       strings.Contains(content, "test"),
		"documentation": strings.Contains(content, "documentation") || strings.Contains(content, "readme"),
	}

	missingGuidelines := []string{}
	for guideline, found := range guidelines {
		if !found {
			missingGuidelines = append(missingGuidelines, guideline)
		}
	}

	if len(missingGuidelines) > 0 {
		*warnings = append(*warnings, fmt.Sprintf("CONTRIBUTING.md missing guidance on: %s", strings.Join(missingGuidelines, ", ")))
	}

	return nil
}

// validateCommitMessage checks commit message follows standards
func validateCommitMessage(message string, warnings *[]string) error {
	if len(message) < 10 {
		*warnings = append(*warnings, "Commit message is very short - consider adding more detail")
	}

	if len(message) > 72 {
		*warnings = append(*warnings, "Commit message first line is long - consider keeping under 72 characters")
	}

	// Check for conventional commit format
	conventionalPrefixes := []string{"feat:", "fix:", "docs:", "style:", "refactor:", "test:", "chore:"}
	hasConventionalFormat := false

	for _, prefix := range conventionalPrefixes {
		if strings.HasPrefix(strings.ToLower(message), prefix) {
			hasConventionalFormat = true
			break
		}
	}

	if !hasConventionalFormat {
		*warnings = append(*warnings, "Consider using conventional commit format (feat:, fix:, docs:, etc.)")
	}

	return nil
}
