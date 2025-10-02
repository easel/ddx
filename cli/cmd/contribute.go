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
	// Extract flags to options struct
	opts, err := extractContributeOptions(cmd)
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

	// Check if there are uncommitted changes in the DDx library
	libraryPath := getResourcePath(workingDir, opts.ResourcePath)
	hasChanges, err := checkForChangesInDir(workingDir, libraryPath)
	if err != nil {
		return nil, err
	}

	if !hasChanges {
		result.Success = false
		result.Message = "No changes detected in .ddx/library"
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
func extractContributeOptions(cmd *cobra.Command) (*ContributeOptions, error) {
	opts := &ContributeOptions{
		ResourcePath: "library", // Always contribute from .ddx/library
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
	// Check actual git subtree (library is at .ddx/library not .ddx)
	ddxPath := ".ddx/library"

	// Change to working directory if specified (git commands need to run in the repo)
	if workingDir != "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return false, fmt.Errorf("failed to get current directory: %w", err)
		}
		defer func() { _ = os.Chdir(currentDir) }() // Restore after check

		if err := os.Chdir(workingDir); err != nil {
			return false, fmt.Errorf("failed to change to working directory: %w", err)
		}
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
		// Use timestamp-based branch name
		opts.Branch = fmt.Sprintf("contrib-%d", time.Now().Unix())
	}

	return nil
}

func checkForChangesInDir(workingDir, fullPath string) (bool, error) {
	// Check git for uncommitted changes
	ddxPath := ".ddx"

	// Change to working directory if specified (git commands need to run in the repo)
	if workingDir != "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return false, fmt.Errorf("failed to get current directory: %w", err)
		}
		defer func() { _ = os.Chdir(currentDir) }() // Restore after check

		if err := os.Chdir(workingDir); err != nil {
			return false, fmt.Errorf("failed to change to working directory: %w", err)
		}
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
	// Determine contribution branch
	contributionBranch := opts.Branch
	if contributionBranch == "" {
		// Default to "contributions" branch for community contributions
		contributionBranch = "contributions"
	}

	// Change to working directory for git operations
	if workingDir != "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		defer func() { _ = os.Chdir(currentDir) }()

		if err := os.Chdir(workingDir); err != nil {
			return nil, fmt.Errorf("failed to change to working directory: %w", err)
		}
	}

	// Execute git subtree push to contribute changes
	prefix := ".ddx/library"
	repoURL := cfg.Library.Repository.URL

	err := git.SubtreePush(prefix, repoURL, contributionBranch)
	if err != nil {
		// Wrap git error with user-friendly message
		return nil, wrapContributionError(err)
	}

	// Build success result
	result := &ContributeResult{
		Success:      true,
		Message:      "Contribution submitted successfully!",
		Branch:       contributionBranch,
		ResourcePath: opts.ResourcePath,
	}

	// Generate PR instructions if requested
	if opts.CreatePR {
		result.PRInfo = generatePRInstructions(cfg, contributionBranch, opts)
	}

	return result, nil
}

// generatePRInstructions creates PR information for the user
func generatePRInstructions(cfg *config.Config, branch string, opts *ContributeOptions) *PRInfo {
	repoURL := strings.TrimSuffix(cfg.Library.Repository.URL, ".git")
	baseBranch := cfg.Library.Repository.Branch
	if baseBranch == "" {
		baseBranch = "master"
	}

	compareURL := fmt.Sprintf("%s/compare/%s...%s", repoURL, baseBranch, branch)

	return &PRInfo{
		URL:         compareURL,
		Title:       opts.Message,
		Branch:      branch,
		Description: "Visit the URL above to create a pull request",
	}
}

// wrapContributionError wraps git errors with user-friendly messages
func wrapContributionError(err error) error {
	errMsg := err.Error()

	if strings.Contains(errMsg, "authentication") || strings.Contains(errMsg, "Authentication") {
		return fmt.Errorf("authentication required: %w\n\nConfigure git credentials with your GitHub token", err)
	}

	if strings.Contains(errMsg, "rejected") {
		return fmt.Errorf("push rejected: %w\n\nYour contribution conflicts with recent changes. Pull latest and retry", err)
	}

	if strings.Contains(errMsg, "no subtree found") {
		return fmt.Errorf("DDx library not found: %w\n\nRun 'ddx update' to setup the library first", err)
	}

	return fmt.Errorf("failed to push contribution: %w", err)
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
		_, _ = cyan.Fprintln(out, "🔍 Dry run: Contributing changes from .ddx/library")
	} else {
		_, _ = cyan.Fprintln(out, "🚀 Contributing changes from .ddx/library")
	}
	_, _ = fmt.Fprintln(out)

	// Handle error cases
	if !result.Success {
		_, _ = red.Fprintln(out, "❌", result.Message)
		return nil
	}

	// Handle dry-run mode
	if opts.DryRun {
		return displayDryRunContributeResult(out, result)
	}

	// Display validation results
	_, _ = fmt.Fprintln(out, "🔍 Validating contribution...")
	_, _ = fmt.Fprintln(out, "")
	if len(result.ValidationResults) > 0 {

		for _, validation := range result.ValidationResults {
			switch validation.Status {
			case "pass":
				_, _ = green.Fprintf(out, "✓ %s: %s\n", validation.Check, validation.Message)
			case "warning":
				_, _ = yellow.Fprintf(out, "⚠️ %s: %s\n", validation.Check, validation.Message)
			case "fail":
				_, _ = red.Fprintf(out, "❌ %s: %s\n", validation.Check, validation.Message)
			}
		}
		_, _ = fmt.Fprintln(out, "")
	}

	// Display success message
	_, _ = green.Fprintln(out, "✅", result.Message)
	_, _ = fmt.Fprintln(out)

	// Display branch information
	_, _ = fmt.Fprintf(out, "Branch: %s\n", yellow.Sprint(result.Branch))
	_, _ = fmt.Fprintf(out, "Resource: %s\n", yellow.Sprint(result.ResourcePath))
	_, _ = fmt.Fprintln(out)

	// Display pull request information
	if result.PRInfo != nil {
		_, _ = fmt.Fprintln(out, "📝 Pull request information:")
		_, _ = fmt.Fprintf(out, "   URL: %s\n", result.PRInfo.URL)
		_, _ = fmt.Fprintf(out, "   Title: %s\n", result.PRInfo.Title)
		_, _ = fmt.Fprintf(out, "   Branch: %s\n", result.PRInfo.Branch)
		_, _ = fmt.Fprintln(out, "   Ready to push to your fork")
		_, _ = fmt.Fprintln(out)
	}

	// Show next steps
	_, _ = fmt.Fprintln(out, bold.Sprint("🎯 Next Steps:"))
	_, _ = fmt.Fprintln(out)

	_, _ = fmt.Fprintf(out, "1. %s\n", cyan.Sprint("Your changes have been pushed to a feature branch"))
	_, _ = fmt.Fprintf(out, "   Branch: %s\n", yellow.Sprint(result.Branch))
	_, _ = fmt.Fprintln(out)

	if result.PRInfo == nil {
		_, _ = fmt.Fprintf(out, "2. %s\n", cyan.Sprint("Push to your fork and create a pull request"))
		_, _ = fmt.Fprintln(out, "   Visit your repository to create a pull request")
		_, _ = fmt.Fprintln(out)

		_, _ = fmt.Fprintf(out, "3. %s\n", cyan.Sprint("Describe your contribution"))
		_, _ = fmt.Fprintf(out, "   Title: %s\n", opts.Message)
		_, _ = fmt.Fprintf(out, "   Description: Include details about the resource and its usage\n")
		_, _ = fmt.Fprintln(out)
	}

	// Show contribution tips
	_, _ = fmt.Fprintln(out, bold.Sprint("💡 Contribution Tips:"))
	_, _ = fmt.Fprintln(out, "• Include a README.md for new patterns or templates")
	_, _ = fmt.Fprintln(out, "• Add examples and usage instructions")
	_, _ = fmt.Fprintln(out, "• Test your resource with 'ddx apply' before contributing")
	_, _ = fmt.Fprintln(out, "• Follow existing naming conventions")

	return nil
}

func displayDryRunContributeResult(out interface{}, result *ContributeResult) error {
	writer := out.(io.Writer)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	bold := color.New(color.Bold)

	_, _ = fmt.Fprintln(writer, bold.Sprint("🔍 Dry Run Results"))
	_, _ = fmt.Fprintln(writer)

	if result.DryRunPreview != nil {
		preview := result.DryRunPreview

		_, _ = fmt.Fprintln(writer, "Would perform the following actions:")
		_, _ = fmt.Fprintf(writer, "%s", green.Sprintf("✓ Resource to contribute: %s\n", preview.WouldContribute))
		_, _ = fmt.Fprintf(writer, "%s", green.Sprintf("✓ Target branch: %s\n", preview.Branch))
		_, _ = fmt.Fprintf(writer, "%s", green.Sprintf("✓ Files to contribute: %d\n", preview.FilesCount))

		if preview.HasDocumentation {
			_, _ = fmt.Fprintln(writer, green.Sprint("✓ Documentation found (README.md)"))
		} else {
			_, _ = fmt.Fprintln(writer, yellow.Sprint("⚠️ No documentation found - consider adding README.md"))
		}

		// Show warnings
		if len(preview.ValidationWarnings) > 0 {
			_, _ = fmt.Fprintln(writer)
			_, _ = yellow.Fprintln(writer, "⚠️ Warnings:")
			for _, warning := range preview.ValidationWarnings {
				_, _ = fmt.Fprintf(writer, "  • %s\n", warning)
			}
		}

		_, _ = fmt.Fprintln(writer)
		_, _ = fmt.Fprintln(writer, bold.Sprint("🎯 What would happen:"))
		_, _ = fmt.Fprintf(writer, "1. Create feature branch: %s\n", preview.Branch)
		_, _ = fmt.Fprintf(writer, "2. Commit changes in .ddx/%s\n", preview.WouldContribute)
		_, _ = fmt.Fprintln(writer, "3. push branch to upstream repository")
		_, _ = fmt.Fprintln(writer, "4. Prepare pull request")
	}

	_, _ = fmt.Fprintln(writer)
	_, _ = cyan.Fprintln(writer, "💡 To proceed with the contribution, run the command without --dry-run")

	return nil
}

// Legacy function for compatibility
func runContribute(cmd *cobra.Command, args []string) error {
	// Extract flags to options struct
	opts, err := extractContributeOptions(cmd)
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
