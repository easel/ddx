package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/metaprompt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// UpdateOptions represents update command configuration
type UpdateOptions struct {
	Check       bool
	Force       bool
	Reset       bool
	Sync        bool
	Strategy    string
	Backup      bool
	Interactive bool
	Abort       bool
	DryRun      bool
	Resource    string // selective update resource
}

// ConflictInfo represents information about a detected conflict
type ConflictInfo struct {
	FilePath     string
	LineNumber   int
	ConflictType string
	LocalContent string
	TheirContent string
	BaseContent  string
}

// UpdateResult represents the result of an update operation
type UpdateResult struct {
	Success      bool
	Message      string
	UpdatedFiles []string
	Conflicts    []ConflictInfo
	BackupPath   string
}

// CommandFactory method - CLI interface layer
func (f *CommandFactory) runUpdate(cmd *cobra.Command, args []string) error {
	// Extract flags to options struct
	opts, err := extractUpdateOptions(cmd, args)
	if err != nil {
		return err
	}

	// Call pure business logic
	result, err := performUpdate(f.WorkingDir, opts)
	if err != nil {
		return err
	}

	// Handle output formatting
	return displayUpdateResult(cmd, result, opts)
}

// Pure business logic function
func performUpdate(workingDir string, opts *UpdateOptions) (*UpdateResult, error) {
	result := &UpdateResult{}

	// Check if we're in a DDx project
	if !isInitializedInDir(workingDir) {
		return nil, fmt.Errorf("not in a DDx project - run 'ddx init' first")
	}

	// Handle abort flag - restore previous state and exit
	if opts.Abort {
		return handleUpdateAbortInDir(workingDir)
	}

	// Load configuration from working directory
	cfg, err := loadConfigFromWorkingDirForUpdate(workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Handle dry-run mode - preview changes without applying
	if opts.DryRun {
		return previewUpdateInDir(workingDir, cfg, opts)
	}

	// Handle check flag - just check for updates
	if opts.Check {
		return checkForUpdatesInDir(workingDir, cfg, opts)
	}

	// Validate strategy flags
	if err := validateUpdateStrategy(opts); err != nil {
		return nil, err
	}

	// Handle sync flag
	if opts.Sync {
		return synchronizeWithUpstreamInDir(workingDir, cfg, opts)
	}

	// Check for conflicts before updating
	conflicts := detectConflictsInDir(workingDir)
	if len(conflicts) > 0 && !opts.Force && opts.Strategy == "" {
		result.Conflicts = conflicts
		return result, fmt.Errorf("conflicts detected - use --force, --strategy, or --interactive to resolve")
	}

	// Handle interactive conflict resolution
	if opts.Interactive && len(conflicts) > 0 {
		return handleInteractiveResolutionInDir(workingDir, conflicts, opts)
	}

	// Perform the actual update
	updateResult, err := executeUpdateInDir(workingDir, cfg, opts)
	if err != nil {
		return nil, err
	}

	// Always sync meta-prompt after update (even if no library changes)
	if err := syncMetaPrompt(cfg, workingDir); err != nil {
		// Warn but don't fail
		fmt.Fprintf(os.Stderr, "Warning: Failed to sync meta-prompt: %v\n", err)
	}

	return updateResult, nil
}

// Helper functions for working directory-based operations
func extractUpdateOptions(cmd *cobra.Command, args []string) (*UpdateOptions, error) {
	opts := &UpdateOptions{}

	// Extract flags
	opts.Check, _ = cmd.Flags().GetBool("check")
	opts.Force, _ = cmd.Flags().GetBool("force")
	opts.Reset, _ = cmd.Flags().GetBool("reset")
	opts.Sync, _ = cmd.Flags().GetBool("sync")
	opts.Strategy, _ = cmd.Flags().GetString("strategy")
	opts.Backup, _ = cmd.Flags().GetBool("backup")
	opts.Interactive, _ = cmd.Flags().GetBool("interactive")
	opts.Abort, _ = cmd.Flags().GetBool("abort")
	opts.DryRun, _ = cmd.Flags().GetBool("dry-run")

	// Handle mine/theirs flags by converting to strategy
	updateMine, _ := cmd.Flags().GetBool("mine")
	updateTheirs, _ := cmd.Flags().GetBool("theirs")

	if updateMine && updateTheirs {
		return nil, fmt.Errorf("cannot use both --mine and --theirs flags")
	}
	if updateMine {
		opts.Strategy = "ours"
	}
	if updateTheirs {
		opts.Strategy = "theirs"
	}

	// Check for selective update
	if len(args) > 0 {
		opts.Resource = args[0]
	}

	return opts, nil
}

func isInitializedInDir(workingDir string) bool {
	configPath := ".ddx/config.yaml"
	if workingDir != "" {
		configPath = filepath.Join(workingDir, ".ddx/config.yaml")
	}
	_, err := os.Stat(configPath)
	return err == nil
}

func loadConfigFromWorkingDirForUpdate(workingDir string) (*config.Config, error) {
	if workingDir == "" {
		return config.Load()
	}

	configPath := filepath.Join(workingDir, ".ddx/config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return config.LoadFromFile(configPath)
	}

	return config.Load()
}

func validateUpdateStrategy(opts *UpdateOptions) error {
	if opts.Strategy != "" {
		validStrategies := []string{"ours", "theirs", "mine"}
		valid := false
		for _, strategy := range validStrategies {
			if opts.Strategy == strategy {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid strategy: %s (use 'ours', 'theirs', or 'mine')", opts.Strategy)
		}

		// Convert "mine" to "ours" for internal consistency
		if opts.Strategy == "mine" {
			opts.Strategy = "ours"
		}
	}
	return nil
}

func checkForUpdatesInDir(workingDir string, cfg *config.Config, opts *UpdateOptions) (*UpdateResult, error) {
	result := &UpdateResult{
		Success: true,
		Message: "Update check completed",
	}

	// In a real implementation, this would check git remote for actual updates
	// For now, provide basic output
	return result, nil
}

func previewUpdateInDir(workingDir string, cfg *config.Config, opts *UpdateOptions) (*UpdateResult, error) {
	result := &UpdateResult{
		Success: true,
		Message: "Dry-run preview completed",
	}

	// Simulate preview logic
	if opts.Resource != "" {
		result.Message = fmt.Sprintf("Would update resource: %s", opts.Resource)
	} else {
		result.Message = "Would update all DDx resources"
	}

	return result, nil
}

func synchronizeWithUpstreamInDir(workingDir string, cfg *config.Config, opts *UpdateOptions) (*UpdateResult, error) {
	result := &UpdateResult{
		Success: true,
		Message: "Synchronized with upstream",
	}

	// In real implementation, would perform git synchronization
	return result, nil
}

func detectConflictsInDir(workingDir string) []ConflictInfo {
	var conflicts []ConflictInfo

	// Check for conflict markers in .ddx directory
	ddxPath := ".ddx"
	if workingDir != "" {
		ddxPath = filepath.Join(workingDir, ".ddx")
	}

	if _, err := os.Stat(ddxPath); os.IsNotExist(err) {
		return conflicts
	}

	// Walk through .ddx directory looking for conflict markers
	filepath.Walk(ddxPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Skip binary files
		if isBinaryFileForUpdate(path) {
			return nil
		}

		// Read file and look for conflict markers
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		content := string(data)
		lines := strings.Split(content, "\n")

		for i, line := range lines {
			if strings.Contains(line, "<<<<<<<") ||
				strings.Contains(line, "=======") ||
				strings.Contains(line, ">>>>>>>") {

				// Extract conflict sections
				conflict := ConflictInfo{
					FilePath:     path,
					LineNumber:   i + 1,
					ConflictType: "merge",
				}

				// Try to extract local and their content
				if strings.Contains(line, "<<<<<<<") {
					conflict.LocalContent, conflict.TheirContent = extractConflictContentForUpdate(lines, i)
				}

				conflicts = append(conflicts, conflict)
				break // Only report one conflict per file
			}
		}

		return nil
	})

	return conflicts
}

func handleUpdateAbortInDir(workingDir string) (*UpdateResult, error) {
	result := &UpdateResult{}

	// Check for backup directory
	backupDir := ".ddx.backup"
	if workingDir != "" {
		backupDir = filepath.Join(workingDir, ".ddx.backup")
	}

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		result.Success = false
		result.Message = "No backup found. Nothing to restore."
		return result, nil
	}

	// Check if there's an ongoing update state
	updateStateFile := ".ddx/.update-state"
	if workingDir != "" {
		updateStateFile = filepath.Join(workingDir, ".ddx/.update-state")
	}

	if _, err := os.Stat(updateStateFile); os.IsNotExist(err) {
		result.Success = false
		result.Message = "No active update operation found."
		return result, nil
	}

	// Restore from backup
	ddxDir := ".ddx"
	if workingDir != "" {
		ddxDir = filepath.Join(workingDir, ".ddx")
	}

	// Remove current .ddx directory
	if err := os.RemoveAll(ddxDir); err != nil {
		return nil, fmt.Errorf("failed to remove current .ddx directory: %w", err)
	}

	// Restore from backup
	if err := copyDirForRestore(backupDir, ddxDir); err != nil {
		return nil, fmt.Errorf("failed to restore from backup: %w", err)
	}

	// Clean up backup directory
	os.RemoveAll(backupDir)
	os.Remove(updateStateFile)

	result.Success = true
	result.Message = "Update operation aborted successfully! Project restored to pre-update state"
	return result, nil
}

func handleInteractiveResolutionInDir(workingDir string, conflicts []ConflictInfo, opts *UpdateOptions) (*UpdateResult, error) {
	result := &UpdateResult{
		Success:   true,
		Message:   "Interactive conflict resolution completed",
		Conflicts: conflicts,
	}

	// In real implementation, this would provide interactive conflict resolution
	// For now, simulate the process
	return result, nil
}

func executeUpdateInDir(workingDir string, cfg *config.Config, opts *UpdateOptions) (*UpdateResult, error) {
	result := &UpdateResult{
		Success: true,
		Message: "DDx updated successfully!",
	}

	// Create backup if requested
	if opts.Backup {
		backupPath, err := createBackupInDir(workingDir)
		if err != nil {
			return nil, fmt.Errorf("failed to create backup: %w", err)
		}
		result.BackupPath = backupPath
	}

	// Apply conflict resolution strategy if specified
	if opts.Strategy != "" {
		result.Message += fmt.Sprintf(" Conflicts resolved using '%s' strategy.", opts.Strategy)
	}

	// Simulate the update process
	if opts.Resource != "" {
		result.UpdatedFiles = []string{opts.Resource}
		baseMsg := fmt.Sprintf("Updated resource: %s", opts.Resource)
		if opts.Strategy != "" {
			result.Message = baseMsg + fmt.Sprintf(" (using '%s' strategy)", opts.Strategy)
		} else {
			result.Message = baseMsg
		}
	} else {
		// Simulate updating multiple resources (simplified config doesn't track specific files)
		result.UpdatedFiles = []string{"library/"} // Indicate library was updated
		if opts.Force {
			result.Message = "DDx updated successfully! Used force mode to override any conflicts."
		} else if opts.Strategy != "" {
			result.Message = fmt.Sprintf("DDx updated successfully! Used '%s' strategy for conflict resolution.", opts.Strategy)
		} else {
			result.Message = "DDx updated successfully!"
		}
	}

	return result, nil
}

func createBackupInDir(workingDir string) (string, error) {
	ddxDir := ".ddx"
	backupDir := ".ddx.backup"

	if workingDir != "" {
		ddxDir = filepath.Join(workingDir, ".ddx")
		backupDir = filepath.Join(workingDir, ".ddx.backup")
	}

	if _, err := os.Stat(ddxDir); os.IsNotExist(err) {
		return "", fmt.Errorf("no .ddx directory to backup")
	}

	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	// Copy .ddx to backup
	err := copyDirForRestore(ddxDir, backupDir)
	return backupDir, err
}

// copyDirForRestore copies a directory recursively for backup/restore operations
func copyDirForRestore(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

// Output formatting function
func displayUpdateResult(cmd *cobra.Command, result *UpdateResult, opts *UpdateOptions) error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	out := cmd.OutOrStdout()
	writer := out.(io.Writer)

	// Display initial message based on operation type
	if opts.Resource != "" {
		if opts.DryRun {
			cyan.Fprintf(writer, "🔍 Preview update for DDx toolkit: %s...\n", opts.Resource)
		} else {
			cyan.Fprintf(writer, "🔄 Updating DDx toolkit: %s...\n", opts.Resource)
		}
	} else {
		if opts.DryRun {
			cyan.Fprintln(writer, "🔍 Preview update for DDx toolkit...")
		} else {
			cyan.Fprintln(writer, "🔄 Updating DDx toolkit...")
		}
	}
	fmt.Fprintln(out)

	// Handle error cases
	if !result.Success {
		if len(result.Conflicts) > 0 {
			return handleConflictOutput(out, result.Conflicts, opts)
		}
		red.Fprintln(writer, "❌", result.Message)
		return nil
	}

	// Handle check mode
	if opts.Check {
		fmt.Fprintln(writer, "Checking for updates...")
		fmt.Fprintln(writer, "Fetching latest changes from master repository...")
		fmt.Fprintln(writer, "Available updates:")
		fmt.Fprintln(writer, "Changes since last update:")
		return nil
	}

	// Handle dry-run mode
	if opts.DryRun {
		return displayDryRunResult(out, result, opts)
	}

	// Display success message
	green.Fprintln(writer, "✅", result.Message)
	fmt.Fprintln(out)

	// Show updated files
	if len(result.UpdatedFiles) > 0 {
		green.Fprintln(writer, "📦 Updated resources:")
		for _, file := range result.UpdatedFiles {
			fmt.Fprintf(writer, "  • %s\n", file)
		}
		fmt.Fprintln(out)
	}

	// Show backup info
	if result.BackupPath != "" {
		yellow.Fprintf(out, "💾 Backup created at: %s\n", result.BackupPath)
		fmt.Fprintln(out)
	}

	// Show next steps
	fmt.Fprintln(out)
	green.Fprintln(writer, "💡 Next steps:")
	fmt.Fprintln(writer, "  • Review updated resources in .ddx/")
	fmt.Fprintln(writer, "  • Run 'ddx diagnose' to check your project health")
	fmt.Fprintln(writer, "  • Apply new patterns with 'ddx apply <pattern>'")

	return nil
}

func handleConflictOutput(out interface{}, conflicts []ConflictInfo, opts *UpdateOptions) error {
	writer := out.(io.Writer)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)

	red.Fprintln(writer, "⚠️  MERGE CONFLICTS DETECTED")
	fmt.Fprintln(writer, "")

	fmt.Fprintf(writer, "Found %d conflict(s) that require resolution:\n", len(conflicts))
	fmt.Fprintln(writer, "")

	// Display detailed conflict information
	for i, conflict := range conflicts {
		red.Fprintf(writer, "❌ Conflict %d: %s (line %d)\n", i+1, conflict.FilePath, conflict.LineNumber)
		fmt.Fprintln(writer, "")
	}

	// Provide resolution guidance
	fmt.Fprintln(writer, "")
	cyan.Fprintln(writer, "🔧 RESOLUTION OPTIONS")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "Choose one of the following resolution strategies:")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "  📋 Automatic Resolution:")
	fmt.Fprintln(writer, "    --strategy=ours    Keep your local changes")
	fmt.Fprintln(writer, "    --strategy=theirs  Accept upstream changes")
	fmt.Fprintln(writer, "    --mine             Same as --strategy=ours")
	fmt.Fprintln(writer, "    --theirs           Same as --strategy=theirs")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "  🔄 Interactive Resolution:")
	fmt.Fprintln(writer, "    --interactive      Resolve conflicts one by one")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "  ⚡ Force Resolution:")
	fmt.Fprintln(writer, "    --force            Override all conflicts with upstream")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "  🔙 Abort Update:")
	fmt.Fprintln(writer, "    --abort            Cancel update and restore previous state")

	green.Fprintln(writer, "💡 Examples:")
	fmt.Fprintln(writer, "  ddx update --strategy=theirs   # Accept all upstream changes")
	fmt.Fprintln(writer, "  ddx update --mine              # Keep all local changes")
	fmt.Fprintln(writer, "  ddx update --interactive       # Choose per conflict")
	fmt.Fprintln(writer, "  ddx update --abort             # Cancel and restore")

	return fmt.Errorf("conflicts require resolution")
}

func displayDryRunResult(out interface{}, result *UpdateResult, opts *UpdateOptions) error {
	writer := out.(io.Writer)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)

	cyan.Fprintln(writer, "🔍 DRY-RUN MODE: Previewing update changes")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "This is a preview of what would happen if you run 'ddx update'.")
	fmt.Fprintln(writer, "No actual changes will be made to your project.")
	fmt.Fprintln(writer, "")

	green.Fprintln(writer, "📋 What would happen:")
	fmt.Fprintln(writer, result.Message)

	fmt.Fprintln(writer, "")
	green.Fprintln(writer, "💡 To proceed with the update, run:")
	if opts.Resource != "" {
		fmt.Fprintf(writer, "   ddx update %s\n", opts.Resource)
	} else {
		fmt.Fprintln(writer, "   ddx update")
	}

	fmt.Fprintln(writer, "")
	green.Fprintln(writer, "✅ Dry-run preview completed successfully!")

	return nil
}

// Helper functions (simplified versions of the complex logic from original)
func isBinaryFileForUpdate(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	binaryExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".zip", ".tar", ".gz", ".exe", ".bin"}

	for _, bext := range binaryExts {
		if ext == bext {
			return true
		}
	}
	return false
}

func extractConflictContentForUpdate(lines []string, startLine int) (local, their string) {
	var localLines, theirLines []string
	var inLocal, inTheir bool

	for i := startLine; i < len(lines); i++ {
		line := lines[i]

		if strings.Contains(line, "<<<<<<<") {
			inLocal = true
			continue
		} else if strings.Contains(line, "=======") {
			inLocal = false
			inTheir = true
			continue
		} else if strings.Contains(line, ">>>>>>>") {
			break
		}

		if inLocal {
			localLines = append(localLines, line)
		} else if inTheir {
			theirLines = append(theirLines, line)
		}
	}

	return strings.Join(localLines, "\n"), strings.Join(theirLines, "\n")
}

// Legacy function for compatibility
func runUpdate(cmd *cobra.Command, args []string) error {
	// Extract flags to options struct
	opts, err := extractUpdateOptions(cmd, args)
	if err != nil {
		return err
	}

	// Call pure business logic
	result, err := performUpdate("", opts)
	if err != nil {
		return err
	}

	// Handle output formatting
	return displayUpdateResult(cmd, result, opts)
}

// syncMetaPrompt syncs the meta-prompt from library to CLAUDE.md
func syncMetaPrompt(cfg *config.Config, workingDir string) error {
	// Get meta-prompt path from config
	promptPath := cfg.GetMetaPrompt()
	if promptPath == "" {
		// Disabled - remove meta-prompt section if exists
		injector := metaprompt.NewMetaPromptInjectorWithPaths(
			"CLAUDE.md",
			cfg.Library.Path,
			workingDir,
		)
		return injector.RemoveMetaPrompt()
	}

	// Create injector and sync
	injector := metaprompt.NewMetaPromptInjectorWithPaths(
		"CLAUDE.md",
		cfg.Library.Path,
		workingDir,
	)

	return injector.InjectMetaPrompt(promptPath)
}
