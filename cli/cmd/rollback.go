package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// rollbackCmd represents the rollback command
func (f *CommandFactory) newRollbackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback problematic updates to restore previous working state",
		Long: `Rollback problematic updates to restore previous working state.

The rollback command helps you recover from problematic updates by reverting
to a previous working version. It manages backup points and provides safe
restoration options.

Examples:
  # Rollback to the most recent backup point
  ddx rollback

  # List available backup points
  ddx rollback --list

  # Rollback to a specific version
  ddx rollback --to v1.2.0

  # Preview changes before rollback
  ddx rollback --preview

Features:
  • Automatic backup creation before rollback
  • Multiple rollback point management
  • Preview mode to see changes
  • State validation after rollback
  • Recovery instructions on failure`,
		RunE: f.runRollback,
	}

	// Add flags
	cmd.Flags().Bool("list", false, "List available rollback points")
	cmd.Flags().String("to", "", "Rollback to specific version")
	cmd.Flags().Bool("preview", false, "Preview changes without executing rollback")

	return cmd
}

// runRollback executes the rollback command
func (f *CommandFactory) runRollback(cmd *cobra.Command, args []string) error {
	// Color formatters
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)

	// Check if project is initialized
	if !isInitialized() {
		return fmt.Errorf("project is not initialized with DDx. Run 'ddx init' first")
	}

	// Get flags
	listFlag, _ := cmd.Flags().GetBool("list")
	toFlag, _ := cmd.Flags().GetString("to")
	previewFlag, _ := cmd.Flags().GetBool("preview")

	// Initialize backup manager
	backupManager := NewBackupManager()

	// Handle list command
	if listFlag {
		return f.listBackupPoints(cmd, backupManager)
	}

	// Get available backup points
	backupPoints, err := backupManager.ListBackupPoints()
	if err != nil {
		return fmt.Errorf("failed to list backup points: %w", err)
	}

	if len(backupPoints) == 0 {
		red.Fprintln(cmd.ErrOrStderr(), "❌ No backup points available")
		yellow.Fprintln(cmd.ErrOrStderr(), "💡 Backup points are created automatically during updates.")
		yellow.Fprintln(cmd.ErrOrStderr(), "   Run 'ddx update --backup' to create a backup point.")
		return fmt.Errorf("no backup points available for rollback")
	}

	// Determine target version
	var targetVersion string
	if toFlag != "" {
		targetVersion = toFlag
		// Validate the version exists
		found := false
		for _, bp := range backupPoints {
			if bp.Version == targetVersion {
				found = true
				break
			}
		}
		if !found {
			red.Fprintf(cmd.ErrOrStderr(), "❌ Version '%s' not found in backup points\n", targetVersion)
			fmt.Fprintln(cmd.ErrOrStderr(), "\nAvailable versions:")
			for _, bp := range backupPoints {
				fmt.Fprintf(cmd.ErrOrStderr(), "  • %s (%s)\n", bp.Version, bp.Timestamp.Format("2006-01-02 15:04:05"))
			}
			return fmt.Errorf("version '%s' not found", targetVersion)
		}
	} else {
		// Use most recent backup point
		targetVersion = backupPoints[0].Version
	}

	// Handle preview mode
	if previewFlag {
		return f.previewRollback(cmd, backupManager, targetVersion)
	}

	// Execute rollback
	cyan.Fprintln(cmd.OutOrStdout(), "🔄 Starting rollback operation...")
	fmt.Fprintf(cmd.OutOrStdout(), "   Target version: %s\n", targetVersion)
	fmt.Fprintln(cmd.OutOrStdout())

	// Create backup before rollback
	preRollbackVersion := fmt.Sprintf("pre-rollback-%d", time.Now().Unix())
	yellow.Fprintln(cmd.OutOrStdout(), "📦 Creating backup before rollback...")
	if err := backupManager.CreateBackup(preRollbackVersion, "Backup before rollback operation"); err != nil {
		red.Fprintf(cmd.ErrOrStderr(), "❌ Failed to create pre-rollback backup: %v\n", err)
		return f.provideRecoveryInstructions(cmd, err)
	}
	green.Fprintln(cmd.OutOrStdout(), "✅ Pre-rollback backup created")

	// Perform rollback
	yellow.Fprintln(cmd.OutOrStdout(), "🔄 Restoring from backup...")
	if err := backupManager.RestoreFromBackup(targetVersion); err != nil {
		red.Fprintf(cmd.ErrOrStderr(), "❌ Rollback failed: %v\n", err)
		return f.provideRecoveryInstructions(cmd, err)
	}

	// Validate state after rollback
	yellow.Fprintln(cmd.OutOrStdout(), "🔍 Validating system integrity...")
	if err := f.validateSystemIntegrity(); err != nil {
		red.Fprintf(cmd.ErrOrStderr(), "❌ System integrity validation failed: %v\n", err)
		return f.provideRecoveryInstructions(cmd, fmt.Errorf("rollback validation failed: %w", err))
	} else {
		green.Fprintln(cmd.OutOrStdout(), "✅ System integrity validated")
	}

	green.Fprintln(cmd.OutOrStdout(), "✅ Rollback completed successfully!")
	fmt.Fprintf(cmd.OutOrStdout(), "🔄 Project restored to version: %s\n", targetVersion)
	fmt.Fprintln(cmd.OutOrStdout())

	return nil
}

// listBackupPoints displays available backup points
func (f *CommandFactory) listBackupPoints(cmd *cobra.Command, backupManager *BackupManager) error {
	backupPoints, err := backupManager.ListBackupPoints()
	if err != nil {
		return fmt.Errorf("failed to list backup points: %w", err)
	}

	if len(backupPoints) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No backup points available")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Available rollback points:")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-25s %s\n", "VERSION", "TIMESTAMP", "DESCRIPTION")
	fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-25s %s\n", strings.Repeat("-", 20), strings.Repeat("-", 25), strings.Repeat("-", 40))

	for _, bp := range backupPoints {
		timestamps := bp.Timestamp.Format("2006-01-02 15:04:05")
		descriptions := bp.Description
		if len(descriptions) > 40 {
			descriptions = descriptions[:37] + "..."
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-25s %s\n", bp.Version, timestamps, descriptions)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "Use 'ddx rollback --to <version>' to rollback to a specific version")

	return nil
}

// previewRollback shows what changes will be reverted
func (f *CommandFactory) previewRollback(cmd *cobra.Command, backupManager *BackupManager, targetVersion string) error {
	yellow := color.New(color.FgYellow)
	cyan := color.New(color.FgCyan)

	cyan.Fprintln(cmd.OutOrStdout(), "🔍 Rollback Preview")
	fmt.Fprintf(cmd.OutOrStdout(), "   Target version: %s\n", targetVersion)
	fmt.Fprintln(cmd.OutOrStdout())

	// Get current state
	currentConfig, err := f.getCurrentConfig()
	if err != nil {
		return fmt.Errorf("failed to read current configuration: %w", err)
	}

	// Get target state
	targetConfig, err := backupManager.GetBackupConfig(targetVersion)
	if err != nil {
		return fmt.Errorf("failed to read target configuration: %w", err)
	}

	// Show differences
	yellow.Fprintln(cmd.OutOrStdout(), "📝 Changes that will be reverted:")
	fmt.Fprintln(cmd.OutOrStdout())

	// Compare configurations
	if currentConfig.Version != targetConfig.Version {
		fmt.Fprintf(cmd.OutOrStdout(), "  Version: %s → %s\n", currentConfig.Version, targetConfig.Version)
	}

	if currentConfig.Library.Repository != targetConfig.Library.Repository {
		fmt.Fprintf(cmd.OutOrStdout(), "  Repository: %s → %s\n", currentConfig.Library.Repository, targetConfig.Library.Repository)
	}

	if currentConfig.Library.Branch != targetConfig.Library.Branch {
		fmt.Fprintf(cmd.OutOrStdout(), "  Branch: %s → %s\n", currentConfig.Library.Branch, targetConfig.Library.Branch)
	}

	fmt.Fprintln(cmd.OutOrStdout())
	yellow.Fprintln(cmd.OutOrStdout(), "📋 This is a preview only. No changes will be made.")
	fmt.Fprintln(cmd.OutOrStdout(), "   Run 'ddx rollback --to "+targetVersion+"' to execute the rollback.")

	return nil
}

// provideRecoveryInstructions provides clear recovery instructions on failure
func (f *CommandFactory) provideRecoveryInstructions(cmd *cobra.Command, originalErr error) error {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)

	fmt.Fprintln(cmd.ErrOrStderr())
	red.Fprintln(cmd.ErrOrStderr(), "🚨 Rollback Operation Failed")
	fmt.Fprintln(cmd.ErrOrStderr())

	yellow.Fprintln(cmd.ErrOrStderr(), "📋 Recovery Instructions:")
	fmt.Fprintln(cmd.ErrOrStderr(), "")
	fmt.Fprintln(cmd.ErrOrStderr(), "1. Check backup integrity:")
	fmt.Fprintln(cmd.ErrOrStderr(), "   ls -la .ddx/backups/")
	fmt.Fprintln(cmd.ErrOrStderr(), "")
	fmt.Fprintln(cmd.ErrOrStderr(), "2. Attempt manual restoration:")
	fmt.Fprintln(cmd.ErrOrStderr(), "   cp .ddx/backups/<version>/.ddx.yml .ddx.yml")
	fmt.Fprintln(cmd.ErrOrStderr(), "")
	fmt.Fprintln(cmd.ErrOrStderr(), "3. Use recovery mode:")
	fmt.Fprintln(cmd.ErrOrStderr(), "   ddx init --force --recover")
	fmt.Fprintln(cmd.ErrOrStderr(), "")
	fmt.Fprintln(cmd.ErrOrStderr(), "4. Contact support with error logs:")
	fmt.Fprintf(cmd.ErrOrStderr(), "   Original error: %v\n", originalErr)
	fmt.Fprintln(cmd.ErrOrStderr(), "")
	fmt.Fprintln(cmd.ErrOrStderr(), "5. Follow disaster recovery plan:")
	fmt.Fprintln(cmd.ErrOrStderr(), "   Refer to project documentation for backup procedures")

	return originalErr
}

// validateSystemIntegrity validates the system state after rollback
func (f *CommandFactory) validateSystemIntegrity() error {
	// Check .ddx.yml exists and is valid
	configPath := ".ddx.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf(".ddx.yml configuration file missing")
	}

	// Validate YAML structure
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("cannot read .ddx.yml: %w", err)
	}

	var config struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Library struct {
			Repository string `yaml:"repository"`
			Branch     string `yaml:"branch"`
		} `yaml:"library"`
	}

	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("invalid YAML in .ddx.yml: %w", err)
	}

	// Basic validation
	if config.Name == "" {
		return fmt.Errorf("missing project name in configuration")
	}

	if config.Version == "" {
		return fmt.Errorf("missing version in configuration")
	}

	// Check .ddx directory structure
	if _, err := os.Stat(".ddx"); os.IsNotExist(err) {
		return fmt.Errorf(".ddx directory missing")
	}

	if _, err := os.Stat(".ddx/backups"); os.IsNotExist(err) {
		return fmt.Errorf(".ddx/backups directory missing")
	}

	return nil
}

// getCurrentConfig reads the current .ddx.yml configuration
func (f *CommandFactory) getCurrentConfig() (*Config, error) {
	configData, err := os.ReadFile(".ddx.yml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Config represents the DDx configuration structure
type Config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Library struct {
		Repository string `yaml:"repository"`
		Branch     string `yaml:"branch"`
	} `yaml:"library"`
}

// BackupManager handles backup operations
type BackupManager struct {
	backupDir string
}

// NewBackupManager creates a new backup manager
func NewBackupManager() *BackupManager {
	return &BackupManager{
		backupDir: ".ddx/backups",
	}
}

// BackupPoint represents a backup point
type BackupPoint struct {
	Version     string    `yaml:"version"`
	Timestamp   time.Time `yaml:"timestamp"`
	Description string    `yaml:"description"`
}

// ListBackupPoints lists all available backup points
func (bm *BackupManager) ListBackupPoints() ([]*BackupPoint, error) {
	if _, err := os.Stat(bm.backupDir); os.IsNotExist(err) {
		return []*BackupPoint{}, nil
	}

	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, err
	}

	var backupPoints []*BackupPoint
	for _, entry := range entries {
		if entry.IsDir() {
			metadataPath := filepath.Join(bm.backupDir, entry.Name(), "metadata.yml")
			if _, err := os.Stat(metadataPath); err == nil {
				data, err := os.ReadFile(metadataPath)
				if err != nil {
					continue
				}

				var bp BackupPoint
				if err := yaml.Unmarshal(data, &bp); err != nil {
					continue
				}

				backupPoints = append(backupPoints, &bp)
			}
		}
	}

	// Sort by timestamp (newest first)
	for i := 0; i < len(backupPoints)-1; i++ {
		for j := i + 1; j < len(backupPoints); j++ {
			if backupPoints[i].Timestamp.Before(backupPoints[j].Timestamp) {
				backupPoints[i], backupPoints[j] = backupPoints[j], backupPoints[i]
			}
		}
	}

	return backupPoints, nil
}

// CreateBackup creates a new backup point
func (bm *BackupManager) CreateBackup(version, description string) error {
	// Ensure backup directory exists
	if err := os.MkdirAll(bm.backupDir, 0755); err != nil {
		return err
	}

	backupPath := filepath.Join(bm.backupDir, version)
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return err
	}

	// Create metadata
	metadata := BackupPoint{
		Version:     version,
		Timestamp:   time.Now(),
		Description: description,
	}

	metadataData, err := yaml.Marshal(&metadata)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(backupPath, "metadata.yml"), metadataData, 0644); err != nil {
		return err
	}

	// Backup current .ddx.yml
	if _, err := os.Stat(".ddx.yml"); err == nil {
		configData, err := os.ReadFile(".ddx.yml")
		if err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(backupPath, ".ddx.yml"), configData, 0644); err != nil {
			return err
		}
	}

	return nil
}

// RestoreFromBackup restores from a specific backup point
func (bm *BackupManager) RestoreFromBackup(version string) error {
	backupPath := filepath.Join(bm.backupDir, version)

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup point '%s' not found", version)
	}

	// Restore .ddx.yml
	configBackupPath := filepath.Join(backupPath, ".ddx.yml")
	if _, err := os.Stat(configBackupPath); err == nil {
		configData, err := os.ReadFile(configBackupPath)
		if err != nil {
			return fmt.Errorf("failed to read backup config: %w", err)
		}

		if err := os.WriteFile(".ddx.yml", configData, 0644); err != nil {
			return fmt.Errorf("failed to restore config: %w", err)
		}
	}

	return nil
}

// GetBackupConfig reads the configuration from a backup point
func (bm *BackupManager) GetBackupConfig(version string) (*Config, error) {
	configPath := filepath.Join(bm.backupDir, version, ".ddx.yml")
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
