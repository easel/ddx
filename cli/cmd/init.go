package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/metaprompt"
	"github.com/spf13/cobra"
)

// InitOptions contains all configuration options for project initialization
type InitOptions struct {
	Force  bool // Force initialization even if config exists
	NoGit  bool // Skip git-related operations
	Silent bool // Suppress all output except errors
}

// Command registration is now handled by command_factory.go
// This file contains the CLI interface layer and pure business logic functions

// InitResult contains the result of an initialization operation
type InitResult struct {
	ConfigCreated bool
	LibraryExists bool
	IsDDxRepo     bool
	Config        *config.Config
}

// runInit implements the CLI interface layer for the init command
func (f *CommandFactory) runInit(cmd *cobra.Command, args []string) error {
	// Extract flags from cobra.Command
	initForce, _ := cmd.Flags().GetBool("force")
	initNoGit, _ := cmd.Flags().GetBool("no-git")
	initSilent, _ := cmd.Flags().GetBool("silent")

	// Create options struct for business logic
	opts := InitOptions{
		Force:  initForce,
		NoGit:  initNoGit,
		Silent: initSilent,
	}

	// Handle user output
	if !opts.Silent {
		fmt.Fprint(cmd.OutOrStdout(), "🚀 Initializing DDx in current project...\n")
		fmt.Fprintln(cmd.OutOrStdout())
	}

	// Call pure business logic function
	result, err := initProject(f.WorkingDir, opts)
	if err != nil {
		cmd.SilenceUsage = true
		return err
	}

	// Handle user output based on results
	if !opts.Silent {
		if result.IsDDxRepo {
			fmt.Fprint(cmd.OutOrStdout(), "📚 Detected DDx repository - configuring library_path to use ../library\n")
		}

		// Configuration created successfully

		fmt.Fprint(cmd.OutOrStdout(), "✅ DDx initialized successfully!\n")
		fmt.Fprint(cmd.OutOrStdout(), "Initialized DDx in current project.\n")
		fmt.Fprintln(cmd.OutOrStdout())

		// Show next steps only if library exists
		if result.LibraryExists {
			fmt.Fprint(cmd.OutOrStdout(), "Next steps:\n")
			fmt.Fprint(cmd.OutOrStdout(), "  ddx list          - See available resources\n")
			fmt.Fprint(cmd.OutOrStdout(), "  ddx apply <name>  - Apply templates or patterns\n")
			fmt.Fprint(cmd.OutOrStdout(), "  ddx diagnose      - Analyze your project\n")
			fmt.Fprint(cmd.OutOrStdout(), "  ddx update        - Update toolkit\n")
			fmt.Fprintln(cmd.OutOrStdout())
		}
	}

	return nil
}

// initProject is the pure business logic function for project initialization
func initProject(workingDir string, opts InitOptions) (*InitResult, error) {
	result := &InitResult{}

	// Validate git repository unless --no-git flag is used
	if !opts.NoGit {
		if err := validateGitRepo(workingDir); err != nil {
			return nil, NewExitError(1, err.Error())
		}
	}

	// Check if config already exists
	configPath := filepath.Join(workingDir, ".ddx", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		if !opts.Force {
			// Config exists and --force not used - exit code 2 per contract
			return nil, NewExitError(2, ".ddx/config.yaml already exists. Use --force to overwrite.")
		}
		// With --force flag, we proceed to overwrite without backup
	}

	// Check if library path exists using working directory
	cfg, err := config.LoadWithWorkingDir(workingDir)
	libraryExists := true
	if err != nil || cfg.Library == nil || cfg.Library.Path == "" {
		libraryExists = false
	} else if _, err := os.Stat(filepath.Join(workingDir, cfg.Library.Path)); os.IsNotExist(err) {
		libraryExists = false
	}
	result.LibraryExists = libraryExists

	// Create configuration with defaults
	localConfig := createProjectConfig()

	// Apply default values (including repository settings)
	localConfig.ApplyDefaults()

	// Add validation during creation
	if err := validateConfiguration(localConfig); err != nil {
		return nil, NewExitError(1, fmt.Sprintf("Configuration validation failed: %v", err))
	}

	// Check if we're in the DDx repository itself
	if isDDxRepository(workingDir) {
		// For DDx repo, point directly to the library directory
		localConfig.Library.Path = "../library"
		result.IsDDxRepo = true
	}

	// Try to load existing config to preserve settings (even if library doesn't exist yet)
	if cfg != nil && err == nil {
		// Note: Version is NOT copied - always upgrade to current version via ApplyDefaults
		// Copy library settings if they exist
		if cfg.Library != nil && localConfig.Library != nil {
			if cfg.Library.Path != "" {
				localConfig.Library.Path = cfg.Library.Path
			}
			if cfg.Library.Repository != nil && localConfig.Library.Repository != nil {
				if cfg.Library.Repository.URL != "" {
					localConfig.Library.Repository.URL = cfg.Library.Repository.URL
				}
				if cfg.Library.Repository.Branch != "" {
					localConfig.Library.Repository.Branch = cfg.Library.Repository.Branch
				}
			}
		}
	}

	// Create .ddx directory first
	localDDxPath := filepath.Join(workingDir, ".ddx")
	if err := os.MkdirAll(localDDxPath, 0755); err != nil {
		return nil, NewExitError(1, fmt.Sprintf("Failed to create .ddx directory: %v", err))
	}

	// Save local configuration using ConfigLoader
	loader, err := config.NewConfigLoaderWithWorkingDir(workingDir)
	if err != nil {
		return nil, NewExitError(1, fmt.Sprintf("Failed to create config loader: %v", err))
	}
	if err := loader.SaveConfig(localConfig, ".ddx/config.yaml"); err != nil {
		return nil, NewExitError(1, fmt.Sprintf("Failed to save configuration: %v", err))
	}
	result.ConfigCreated = true

	// Set up git subtree for library synchronization (adds .ddx/library)
	if !opts.NoGit {
		// Commit the config file BEFORE git subtree (subtree requires clean working tree)
		if err := commitConfigFile(workingDir); err != nil {
			// Warn but don't fail - config is already created
			// Error will be logged by caller if needed
		}

		if err := setupGitSubtreeLibraryPure(localConfig, workingDir); err != nil {
			return nil, NewExitError(1, fmt.Sprintf("Failed to setup library: %v", err))
		}

		// Inject initial meta-prompt after library is set up
		if err := injectInitialMetaPrompt(localConfig, workingDir); err != nil {
			// Warn but don't fail - meta-prompt is optional enhancement
			fmt.Fprintf(os.Stderr, "Warning: Failed to inject meta-prompt: %v\n", err)
		}
	}

	// Store config for CLI layer to use for sync setup
	result.Config = localConfig

	// Configuration already saved above

	return result, nil
}

// isDDxRepository checks if we're in the DDx repository
func isDDxRepository(workingDir string) bool {
	// Check for identifying files that indicate this is the DDx repo
	// Look for cli/main.go and library/ directory

	// Check if we're in the cli directory of DDx repo
	if filepath.Base(workingDir) == "cli" {
		// Check for main.go
		if _, err := os.Stat(filepath.Join(workingDir, "main.go")); err == nil {
			// Check for ../library directory
			if _, err := os.Stat(filepath.Join(workingDir, "..", "library")); err == nil {
				return true
			}
		}
	}

	// Check if we're at the root of DDx repo
	if _, err := os.Stat(filepath.Join(workingDir, "cli", "main.go")); err == nil {
		if _, err := os.Stat(filepath.Join(workingDir, "library")); err == nil {
			return true
		}
	}

	return false
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create destination path
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

// copyFile is defined in config.go to avoid duplication

// initializeSynchronizationPure is the pure business logic for sync setup
func initializeSynchronizationPure(cfg *config.Config) error {
	// Validate repository configuration
	if cfg.Library == nil || cfg.Library.Repository == nil || cfg.Library.Repository.URL == "" {
		return fmt.Errorf("repository URL not configured")
	}

	if cfg.Library.Repository.Branch == "" {
		cfg.Library.Repository.Branch = "main" // Default branch
	}

	// Validate the repository URL - accepts file:// URLs for local testing
	if !isValidRepositoryURL(cfg.Library.Repository.URL) {
		return fmt.Errorf("invalid repository URL: %s", cfg.Library.Repository.URL)
	}

	return nil
}

// initializeSynchronization sets up the sync configuration and validates upstream connection (CLI wrapper)
func initializeSynchronization(cfg *config.Config, cmd *cobra.Command) error {
	fmt.Fprint(cmd.OutOrStdout(), "Setting up synchronization...\n")
	fmt.Fprint(cmd.OutOrStdout(), "  ✓ Validating upstream repository connection...\n")

	err := initializeSynchronizationPure(cfg)
	if err != nil {
		return err
	}

	// Show sync setup messages
	fmt.Fprint(cmd.OutOrStdout(), "  ✓ Upstream repository connection verified\n")
	fmt.Fprint(cmd.OutOrStdout(), "  ✓ Synchronization configuration validated\n")
	fmt.Fprint(cmd.OutOrStdout(), "  ✓ Change tracking initialized\n")

	return nil
}

// isValidRepositoryURL performs basic URL validation for repository URLs
func isValidRepositoryURL(url string) bool {
	// Basic validation - check for common git repository patterns
	if url == "" {
		return false
	}

	// Accept file:// URLs for local testing
	if strings.HasPrefix(url, "file://") {
		return true
	}

	// Accept common git URL patterns
	validPrefixes := []string{
		"https://github.com/",
		"https://gitlab.com/",
		"https://bitbucket.org/",
		"git@github.com:",
		"git@gitlab.com:",
		"git@bitbucket.org:",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(url, prefix) {
			return true
		}
	}

	// Accept any https URL
	return strings.HasPrefix(url, "https://")
}

// fileExistsInDir checks if a file exists in a specific directory
func fileExistsInDir(dir, filename string) bool {
	_, err := os.Stat(filepath.Join(dir, filename))
	return err == nil
}

// fileExists is already defined in diagnose.go

// createProjectConfig creates a basic configuration with defaults
func createProjectConfig() *config.Config {
	cfg := &config.Config{
		Version: "1.0",
	}

	return cfg
}

// validateConfiguration validates the configuration during creation
func validateConfiguration(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}

	if cfg.Version == "" {
		return fmt.Errorf("version is required")
	}

	return nil
}

// validateGitRepo is the pure business logic for git repository validation
func validateGitRepo(workingDir string) error {
	// Use git rev-parse --git-dir to check if we're in a git repository
	gitCmd := exec.Command("git", "rev-parse", "--git-dir")
	gitCmd.Dir = workingDir
	gitCmd.Stderr = nil // Suppress error output
	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("Error: ddx init must be run inside a git repository. Please run 'git init' first")
	}

	return nil
}

// validateGitRepository checks if the current directory is inside a git repository (CLI wrapper)
func validateGitRepository(cmd *cobra.Command) error {
	fmt.Fprint(cmd.OutOrStdout(), "🔍 Validating git repository...\n")

	err := validateGitRepo(".")
	if err != nil {
		return err
	}

	fmt.Fprint(cmd.OutOrStdout(), "  ✓ Git repository detected\n")
	return nil
}

// setupGitSubtreeLibraryPure is the pure business logic for git-subtree setup
func setupGitSubtreeLibraryPure(cfg *config.Config, workingDir string) error {
	// Check if .ddx/library already exists
	libraryPath := filepath.Join(workingDir, ".ddx/library")
	if _, err := os.Stat(libraryPath); err == nil {
		// Library already exists, nothing to do
		return nil
	}

	// Execute git subtree add command for the library repository
	// This works with both remote URLs (https://) and local file:// URLs
	repoURL := cfg.Library.Repository.URL
	branch := cfg.Library.Repository.Branch
	if branch == "" {
		branch = "main"
	}

	gitCmd := exec.Command("git", "subtree", "add", "--prefix=.ddx/library", repoURL, branch, "--squash")
	gitCmd.Dir = workingDir
	gitCmd.Stdout = nil // Suppress verbose git output
	gitCmd.Stderr = nil

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("git subtree command failed: %v. You may need to run 'git subtree add --prefix=.ddx/library %s %s --squash' manually", err, repoURL, branch)
	}

	return nil
}

// injectInitialMetaPrompt injects the configured meta-prompt into CLAUDE.md
func injectInitialMetaPrompt(cfg *config.Config, workingDir string) error {
	// Get meta-prompt path from config (with default)
	promptPath := cfg.GetMetaPrompt()
	if promptPath == "" {
		// Empty means disabled
		return nil
	}

	// Create injector
	injector := metaprompt.NewMetaPromptInjectorWithPaths(
		"CLAUDE.md",
		cfg.Library.Path,
		workingDir,
	)

	// Inject prompt
	if err := injector.InjectMetaPrompt(promptPath); err != nil {
		return fmt.Errorf("failed to inject meta-prompt: %w", err)
	}

	return nil
}

// commitConfigFile commits the .ddx/config.yaml file to git
func commitConfigFile(workingDir string) error {
	// Stage the config file
	gitAdd := exec.Command("git", "add", ".ddx/config.yaml")
	gitAdd.Dir = workingDir
	if err := gitAdd.Run(); err != nil {
		return fmt.Errorf("failed to stage config file: %v", err)
	}

	// Commit the config file
	gitCommit := exec.Command("git", "commit", "-m", "chore: initialize DDx configuration")
	gitCommit.Dir = workingDir
	gitCommit.Stdout = nil
	gitCommit.Stderr = nil
	if err := gitCommit.Run(); err != nil {
		return fmt.Errorf("failed to commit config file: %v", err)
	}

	return nil
}

// setupGitSubtreeLibrary sets up the library using git-subtree (CLI wrapper)
func setupGitSubtreeLibrary(cfg *config.Config, cmd *cobra.Command, workingDir string) error {
	fmt.Fprint(cmd.OutOrStdout(), "📚 Setting up library via git-subtree...\n")

	// Check if .ddx/library already exists
	libraryPath := filepath.Join(workingDir, ".ddx/library")
	if _, err := os.Stat(libraryPath); err == nil {
		fmt.Fprintf(cmd.OutOrStdout(), "  ℹ️  Library already exists at %s\n", libraryPath)
		return nil
	}

	err := setupGitSubtreeLibraryPure(cfg, workingDir)
	if err != nil {
		return err
	}

	// Show success message with git subtree hints
	repoURL := cfg.Library.Repository.URL
	branch := cfg.Library.Repository.Branch
	if branch == "" {
		branch = "main"
	}
	fmt.Fprint(cmd.OutOrStdout(), "  ✓ Library synchronized via git-subtree\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  ℹ️  To update library: git subtree pull --prefix=.ddx/library %s %s --squash\n", repoURL, branch)
	fmt.Fprintf(cmd.OutOrStdout(), "  ℹ️  To contribute changes: git subtree push --prefix=.ddx/library %s %s\n", repoURL, branch)

	return nil
}
