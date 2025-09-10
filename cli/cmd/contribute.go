package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/git"
)

var (
	contributeMessage string
	contributeBranch  string
)

var contributeCmd = &cobra.Command{
	Use:   "contribute <path>",
	Short: "Contribute improvements back to master repository",
	Long: `Contribute local improvements back to the master DDx repository.

This command:
• Creates a feature branch in the DDx subtree
• Commits your changes with a descriptive message
• Pushes to your fork (if configured)
• Provides instructions for creating a pull request

Examples:
  ddx contribute patterns/my-pattern
  ddx contribute prompts/claude/new-prompt.md
  ddx contribute scripts/setup/my-script.sh`,
	Args: cobra.ExactArgs(1),
	RunE: runContribute,
}

func init() {
	rootCmd.AddCommand(contributeCmd)
	
	contributeCmd.Flags().StringVarP(&contributeMessage, "message", "m", "", "Contribution message")
	contributeCmd.Flags().StringVar(&contributeBranch, "branch", "", "Feature branch name")
}

func runContribute(cmd *cobra.Command, args []string) error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	bold := color.New(color.Bold)

	resourcePath := args[0]

	cyan.Printf("🚀 Contributing: %s\n\n", resourcePath)

	// Check if we're in a DDx project
	if !isInitialized() {
		red.Println("❌ Not in a DDx project. Run 'ddx init' first.")
		return nil
	}

	// Check if it's a git repository
	if !git.IsRepository(".") {
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
	hasSubtree, err := git.HasSubtree(".ddx")
	if err != nil {
		return fmt.Errorf("failed to check for DDx subtree: %w", err)
	}

	if !hasSubtree {
		red.Println("❌ No DDx subtree found. Run 'ddx update' to set up.")
		return nil
	}

	// Get contribution details
	if contributeMessage == "" {
		prompt := &survey.Input{
			Message: "Describe your contribution:",
			Help:    "A brief description of what you're contributing",
		}
		if err := survey.AskOne(prompt, &contributeMessage); err != nil {
			return err
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
	hasChanges, err := git.HasUncommittedChanges(".ddx")
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to check for changes: %w", err)
	}

	if !hasChanges {
		s.Stop()
		yellow.Printf("⚠️  No changes detected in %s\n", resourcePath)
		return nil
	}

	s.Suffix = " Creating feature branch..."

	// Create and push the contribution
	if err := git.SubtreePush(".ddx", cfg.Repository.URL, contributeBranch); err != nil {
		s.Stop()
		return fmt.Errorf("failed to push contribution: %w", err)
	}

	s.Stop()
	green.Println("✅ Contribution prepared successfully!")
	fmt.Println()

	// Show next steps
	bold.Println("🎯 Next Steps:")
	fmt.Println()
	
	fmt.Printf("1. %s\n", cyan.Sprint("Your changes have been pushed to a feature branch"))
	fmt.Printf("   Branch: %s\n", yellow.Sprint(contributeBranch))
	fmt.Printf("   Resource: %s\n", yellow.Sprint(resourcePath))
	fmt.Println()
	
	fmt.Printf("2. %s\n", cyan.Sprint("Create a Pull Request"))
	if cfg.Repository.URL != "" {
		// Extract repo info from URL
		repoURL := strings.TrimSuffix(cfg.Repository.URL, ".git")
		fmt.Printf("   Visit: %s/compare/%s...%s\n", 
			repoURL, 
			cfg.Repository.Branch, 
			contributeBranch)
	}
	fmt.Println()
	
	fmt.Printf("3. %s\n", cyan.Sprint("Describe your contribution"))
	fmt.Printf("   Title: %s\n", contributeMessage)
	fmt.Printf("   Description: Include details about the resource and its usage\n")
	fmt.Println()

	// Show contribution tips
	fmt.Println(color.New(color.Bold).Sprint("💡 Contribution Tips:"))
	fmt.Println("• Include a README.md for new patterns or templates")
	fmt.Println("• Add examples and usage instructions")
	fmt.Println("• Test your resource with 'ddx apply' before contributing")
	fmt.Println("• Follow existing naming conventions")

	return nil
}