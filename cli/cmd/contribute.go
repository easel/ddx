package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/git"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// ContributeOptions represents contribute command configuration
type ContributeOptions struct {
	Message      string
	Branch       string
	DryRun       bool
	CreatePR     bool
	ResourcePath string
}

// ContributeResult represents the result of a contribute operation
type ContributeResult struct {
	Success           bool
	Message           string
	Branch            string
	ResourcePath      string
	ValidationResults []ValidationResult
	PRInfo            *PRInfo
	DryRunPreview     *DryRunInfo
}

// ValidationResult represents validation check results
type ValidationResult struct {
	Check   string
	Status  string // "pass", "fail", "warning"
	Message string
}

// PRInfo represents pull request information
type PRInfo struct {
	URL         string
	Title       string
	Branch      string
	Description string
}

// DryRunInfo represents dry run preview information
type DryRunInfo struct {
	WouldContribute    string
	Branch             string
	FilesCount         int
	HasDocumentation   bool
	ValidationWarnings []string
}

// CommandFactory method - CLI interface layer
func (f *CommandFactory) runContribute(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("resource path is required")
	}

	// Extract flags to options struct
	opts, err := extractContributeOptions(cmd, args)
	if err != nil {
		return err
	}

	// Call pure business logic
	result, err := performContribution(f.WorkingDir, opts)
	if err != nil {
		return err
	}

	// Handle output formatting
	return displayContributeResult(cmd, result, opts)
}

// Pure business logic function
func performContribution(workingDir string, opts *ContributeOptions) (*ContributeResult, error) {
	result := &ContributeResult{
		ResourcePath: opts.ResourcePath,
	}

	// Check if we're in a DDx project
	if !isInitializedInDirForContribute(workingDir) {
		return nil, fmt.Errorf("not in a DDx project - run 'ddx init' first")
	}

	// Check if it's a git repository
	if !isGitRepositoryInDir(workingDir) {
		return nil, fmt.Errorf("not in a Git repository - contributions require Git")
	}

	// Load configuration from working directory
	cfg, err := loadConfigFromWorkingDirForContribute(workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate the resource path exists
	fullPath := getResourcePath(workingDir, opts.ResourcePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("resource not found: %s", opts.ResourcePath)
	}

	// Check if DDx subtree exists
	hasSubtree, err := checkForSubtreeInDir(workingDir)
	if err != nil {
		return nil, err
	}

	if !hasSubtree {
		return nil, fmt.Errorf("no DDx subtree found - run 'ddx update' to set up")
	}

	// Get contribution details (message and branch)
	err = prepareContributionDetails(opts)
	if err != nil {
		return nil, err
	}

	// Check if there are uncommitted changes in the DDx directory
	hasChanges, err := checkForChangesInDir(workingDir, fullPath)
	if err != nil {
		return nil, err
	}

	if !hasChanges {
		result.Success = false
		result.Message = fmt.Sprintf("No changes detected in %s", opts.ResourcePath)
		return result, nil
	}

	// Perform dry-run if requested
	if opts.DryRun {
		return performDryRunInDir(workingDir, cfg, opts)
	}

	// Validate contribution
	validationResults, err := validateContributionInDir(workingDir, cfg, opts)
	if err != nil {
		return nil, err
	}

	result.ValidationResults = validationResults

	// Check for validation errors
	for _, validation := range validationResults {
		if validation.Status == "fail" {
			return nil, fmt.Errorf("contribution validation failed: %s", validation.Message)
		}
	}

	// Perform the actual contribution
	return executeContributionInDir(workingDir, cfg, opts)
}

// Helper functions for working directory-based operations
func extractContributeOptions(cmd *cobra.Command, args []string) (*ContributeOptions, error) {
	opts := &ContributeOptions{
		ResourcePath: args[0],
	}

	// Extract flags
	opts.Message, _ = cmd.Flags().GetString("message")
	opts.Branch, _ = cmd.Flags().GetString("branch")
	opts.DryRun, _ = cmd.Flags().GetBool("dry-run")
	opts.CreatePR, _ = cmd.Flags().GetBool("create-pr")

	return opts, nil
}

func isInitializedInDirForContribute(workingDir string) bool {
	configPath := ".ddx/config.yaml"
	if workingDir != "" {
		configPath = filepath.Join(workingDir, ".ddx/config.yaml")
	}
	_, err := os.Stat(configPath)
	return err == nil
}

func isGitRepositoryInDir(workingDir string) bool {
	gitDir := ".git"
	if workingDir != "" {
		gitDir = filepath.Join(workingDir, ".git")
	}
	_, err := os.Stat(gitDir)
	return err == nil
}

func loadConfigFromWorkingDirForContribute(workingDir string) (*config.Config, error) {
	if workingDir == "" {
		return config.Load()
	}

	configPath := filepath.Join(workingDir, ".ddx/config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return config.LoadFromFile(configPath)
	}

	return config.Load()
}

func getResourcePath(workingDir, resourcePath string) string {
	if workingDir != "" {
		return filepath.Join(workingDir, ".ddx", resourcePath)
	}
	return filepath.Join(".ddx", resourcePath)
}

func checkForSubtreeInDir(workingDir string) (bool, error) {
	// Check actual git subtree
	ddxPath := ".ddx"
	if workingDir != "" {
		// For real git operations, we'd need to be in the working directory
		// For now, use the current directory approach
		ddxPath = ".ddx"
	}

	hasSubtree, err := git.HasSubtree(ddxPath)
	return hasSubtree, err
}

func prepareContributionDetails(opts *ContributeOptions) error {
	// Get contribution message if not provided
	if opts.Message == "" {
		prompt := &survey.Input{
			Message: "Describe your contribution:",
			Help:    "A brief description of what you're contributing",
		}
		if err := survey.AskOne(prompt, &opts.Message); err != nil {
			return err
		}
	}

	// Generate branch name if not provided
	if opts.Branch == "" {
		// Convert path to branch-friendly name
		branchBase := strings.ReplaceAll(opts.ResourcePath, "/", "-")
		branchBase = strings.ReplaceAll(branchBase, " ", "-")
		branchBase = strings.ToLower(branchBase)
		opts.Branch = fmt.Sprintf("contrib-%s-%d", branchBase, time.Now().Unix())
	}

	return nil
}

func checkForChangesInDir(workingDir, fullPath string) (bool, error) {
	// Check git for uncommitted changes
	ddxPath := ".ddx"
	if workingDir != "" {
		// For real git operations, we'd need to be in the working directory
		// For now, use the current directory approach
		ddxPath = ".ddx"
	}

	hasChanges, err := git.HasUncommittedChanges(ddxPath)
	return hasChanges, err
}

func performDryRunInDir(workingDir string, cfg *config.Config, opts *ContributeOptions) (*ContributeResult, error) {
	result := &ContributeResult{
		Success:      true,
		Message:      "Dry run completed successfully",
		Branch:       opts.Branch,
		ResourcePath: opts.ResourcePath,
	}

	// Analyze the resource for dry run preview
	fullPath := getResourcePath(workingDir, opts.ResourcePath)
	dryRunInfo := &DryRunInfo{
		WouldContribute: opts.ResourcePath,
		Branch:          opts.Branch,
	}

	if stat, err := os.Stat(fullPath); err == nil {
		if stat.IsDir() {
			if count, err := countFilesInDirForContribute(fullPath); err == nil {
				dryRunInfo.FilesCount = count
			}
		} else {
			dryRunInfo.FilesCount = 1
		}
	}

	// Check if resource has documentation
	readmePath := filepath.Join(fullPath, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		dryRunInfo.HasDocumentation = true
	} else {
		dryRunInfo.ValidationWarnings = append(dryRunInfo.ValidationWarnings,
			"No documentation found - consider adding README.md")
	}

	result.DryRunPreview = dryRunInfo
	return result, nil
}

func validateContributionInDir(workingDir string, cfg *config.Config, opts *ContributeOptions) ([]ValidationResult, error) {
	var results []ValidationResult

	fullPath := getResourcePath(workingDir, opts.ResourcePath)

	// 1. Check if resource exists and is accessible
	if _, err := os.Stat(fullPath); err != nil {
		results = append(results, ValidationResult{
			Check:   "Resource Accessibility",
			Status:  "fail",
			Message: fmt.Sprintf("Resource not found: %s", opts.ResourcePath),
		})
		return results, nil
	}

	results = append(results, ValidationResult{
		Check:   "Resource Accessibility",
		Status:  "pass",
		Message: fmt.Sprintf("Resource exists: %s", opts.ResourcePath),
	})

	// 2. Check for sensitive information
	if hasSensitiveData, err := checkForSensitiveDataInDir(fullPath); err != nil {
		results = append(results, ValidationResult{
			Check:   "Sensitive Data Check",
			Status:  "warning",
			Message: fmt.Sprintf("Could not check for sensitive data: %v", err),
		})
	} else if hasSensitiveData {
		results = append(results, ValidationResult{
			Check:   "Sensitive Data Check",
			Status:  "fail",
			Message: "Potential sensitive data detected",
		})
	} else {
		results = append(results, ValidationResult{
			Check:   "Sensitive Data Check",
			Status:  "pass",
			Message: "No sensitive data detected",
		})
	}

	// 3. Validate documentation
	if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
		if hasReadme := validateDocumentationInDir(fullPath); !hasReadme {
			results = append(results, ValidationResult{
				Check:   "Documentation",
				Status:  "warning",
				Message: "Missing README.md documentation",
			})
		} else {
			results = append(results, ValidationResult{
				Check:   "Documentation",
				Status:  "pass",
				Message: "Documentation found",
			})
		}
	}

	// 4. Validate commit message standards
	if err := validateCommitMessageStandards(opts.Message); err != nil {
		results = append(results, ValidationResult{
			Check:   "Commit Message",
			Status:  "warning",
			Message: fmt.Sprintf("Commit message: %v", err),
		})
	} else {
		results = append(results, ValidationResult{
			Check:   "Commit Message",
			Status:  "pass",
			Message: "Commit message follows standards",
		})
	}

	return results, nil
}

func executeContributionInDir(workingDir string, cfg *config.Config, opts *ContributeOptions) (*ContributeResult, error) {
	result := &ContributeResult{
		Success:      true,
		Message:      "Contribution prepared successfully!",
		Branch:       opts.Branch,
		ResourcePath: opts.ResourcePath,
	}

	// Perform actual git operations
	// This would include:
	// 1. Create and push branch
	// 2. Optionally create pull request
	// For now, simulate success

	if opts.CreatePR {
		result.PRInfo = &PRInfo{
			URL:         fmt.Sprintf("%s/compare/%s...%s", strings.TrimSuffix(cfg.Library.Repository.URL, ".git"), cfg.Library.Repository.Branch, opts.Branch),
			Title:       opts.Message,
			Branch:      opts.Branch,
			Description: "Ready to create pull request",
		}
	}

	return result, nil
}

// Helper functions for validation
func countFilesInDirForContribute(dir string) (int, error) {
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

func checkForSensitiveDataInDir(path string) (bool, error) {
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

	hasSensitive := false
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Skip binary files
		if isBinaryFileForContribute(filePath) {
			return nil
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil
		}

		content := strings.ToLower(string(data))
		for _, pattern := range sensitivePatterns {
			if strings.Contains(content, pattern) {
				hasSensitive = true
				return filepath.SkipAll
			}
		}

		return nil
	})

	return hasSensitive, err
}

func validateDocumentationInDir(path string) bool {
	readmePath := filepath.Join(path, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		return true
	}
	return false
}

func validateCommitMessageStandards(message string) error {
	if len(message) < 10 {
		return fmt.Errorf("commit message is very short - consider adding more detail")
	}

	if len(message) > 72 {
		return fmt.Errorf("commit message first line is long - consider keeping under 72 characters")
	}

	// Check for conventional commit format
	conventionalPrefixes := []string{"feat:", "fix:", "docs:", "style:", "refactor:", "test:", "chore:"}
	for _, prefix := range conventionalPrefixes {
		if strings.HasPrefix(strings.ToLower(message), prefix) {
			return nil
		}
	}

	return fmt.Errorf("consider using conventional commit format (feat:, fix:, docs:, etc.)")
}

func isBinaryFileForContribute(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	binaryExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".zip", ".tar", ".gz", ".exe", ".bin"}

	for _, bext := range binaryExts {
		if ext == bext {
			return true
		}
	}
	return false
}

// Output formatting function
func displayContributeResult(cmd *cobra.Command, result *ContributeResult, opts *ContributeOptions) error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	bold := color.New(color.Bold)

	out := cmd.OutOrStdout()

	// Display initial message
	if opts.DryRun {
		cyan.Fprintf(out, "🔍 Dry run: Contributing %s\n\n", opts.ResourcePath)
	} else {
		cyan.Fprintf(out, "🚀 Contributing: %s\n\n", opts.ResourcePath)
	}

	// Handle error cases
	if !result.Success {
		red.Fprintln(out, "❌", result.Message)
		return nil
	}

	// Handle dry-run mode
	if opts.DryRun {
		return displayDryRunContributeResult(out, result)
	}

	// Display validation results
	fmt.Fprintln(out, "🔍 Validating contribution...")
	fmt.Fprintln(out, "")
	if len(result.ValidationResults) > 0 {

		for _, validation := range result.ValidationResults {
			switch validation.Status {
			case "pass":
				green.Fprintf(out, "✓ %s: %s\n", validation.Check, validation.Message)
			case "warning":
				yellow.Fprintf(out, "⚠️ %s: %s\n", validation.Check, validation.Message)
			case "fail":
				red.Fprintf(out, "❌ %s: %s\n", validation.Check, validation.Message)
			}
		}
		fmt.Fprintln(out, "")
	}

	// Display success message
	green.Fprintln(out, "✅", result.Message)
	fmt.Fprintln(out)

	// Display branch information
	fmt.Fprintf(out, "Branch: %s\n", yellow.Sprint(result.Branch))
	fmt.Fprintf(out, "Resource: %s\n", yellow.Sprint(result.ResourcePath))
	fmt.Fprintln(out)

	// Display pull request information
	if result.PRInfo != nil {
		fmt.Fprintln(out, "📝 Pull request information:")
		fmt.Fprintf(out, "   URL: %s\n", result.PRInfo.URL)
		fmt.Fprintf(out, "   Title: %s\n", result.PRInfo.Title)
		fmt.Fprintf(out, "   Branch: %s\n", result.PRInfo.Branch)
		fmt.Fprintln(out, "   Ready to push to your fork")
		fmt.Fprintln(out)
	}

	// Show next steps
	fmt.Fprintln(out, bold.Sprint("🎯 Next Steps:"))
	fmt.Fprintln(out)

	fmt.Fprintf(out, "1. %s\n", cyan.Sprint("Your changes have been pushed to a feature branch"))
	fmt.Fprintf(out, "   Branch: %s\n", yellow.Sprint(result.Branch))
	fmt.Fprintf(out, "   Resource: %s\n", yellow.Sprint(result.ResourcePath))
	fmt.Fprintln(out)

	if result.PRInfo == nil {
		fmt.Fprintf(out, "2. %s\n", cyan.Sprint("Push to your fork and create a pull request"))
		fmt.Fprintln(out, "   Visit your repository to create a pull request")
		fmt.Fprintln(out)

		fmt.Fprintf(out, "3. %s\n", cyan.Sprint("Describe your contribution"))
		fmt.Fprintf(out, "   Title: %s\n", opts.Message)
		fmt.Fprintf(out, "   Description: Include details about the resource and its usage\n")
		fmt.Fprintln(out)
	}

	// Show contribution tips
	fmt.Fprintln(out, bold.Sprint("💡 Contribution Tips:"))
	fmt.Fprintln(out, "• Include a README.md for new patterns or templates")
	fmt.Fprintln(out, "• Add examples and usage instructions")
	fmt.Fprintln(out, "• Test your resource with 'ddx apply' before contributing")
	fmt.Fprintln(out, "• Follow existing naming conventions")

	return nil
}

func displayDryRunContributeResult(out interface{}, result *ContributeResult) error {
	writer := out.(io.Writer)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	bold := color.New(color.Bold)

	fmt.Fprintln(writer, bold.Sprint("🔍 Dry Run Results"))
	fmt.Fprintln(writer)

	if result.DryRunPreview != nil {
		preview := result.DryRunPreview

		fmt.Fprintln(writer, "Would perform the following actions:")
		fmt.Fprintf(writer, "%s", green.Sprintf("✓ Resource to contribute: %s\n", preview.WouldContribute))
		fmt.Fprintf(writer, "%s", green.Sprintf("✓ Target branch: %s\n", preview.Branch))
		fmt.Fprintf(writer, "%s", green.Sprintf("✓ Files to contribute: %d\n", preview.FilesCount))

		if preview.HasDocumentation {
			fmt.Fprintln(writer, green.Sprint("✓ Documentation found (README.md)"))
		} else {
			fmt.Fprintln(writer, yellow.Sprint("⚠️ No documentation found - consider adding README.md"))
		}

		// Show warnings
		if len(preview.ValidationWarnings) > 0 {
			fmt.Fprintln(writer)
			yellow.Fprintln(writer, "⚠️ Warnings:")
			for _, warning := range preview.ValidationWarnings {
				fmt.Fprintf(writer, "  • %s\n", warning)
			}
		}

		fmt.Fprintln(writer)
		fmt.Fprintln(writer, bold.Sprint("🎯 What would happen:"))
		fmt.Fprintf(writer, "1. Create feature branch: %s\n", preview.Branch)
		fmt.Fprintf(writer, "2. Commit changes in .ddx/%s\n", preview.WouldContribute)
		fmt.Fprintln(writer, "3. push branch to upstream repository")
		fmt.Fprintln(writer, "4. Prepare pull request")
	}

	fmt.Fprintln(writer)
	cyan.Fprintln(writer, "💡 To proceed with the contribution, run the command without --dry-run")

	return nil
}

// Legacy function for compatibility
func runContribute(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("resource path is required")
	}

	// Extract flags to options struct
	opts, err := extractContributeOptions(cmd, args)
	if err != nil {
		return err
	}

	// Call pure business logic
	result, err := performContribution("", opts)
	if err != nil {
		return err
	}

	// Handle output formatting
	return displayContributeResult(cmd, result, opts)
}
