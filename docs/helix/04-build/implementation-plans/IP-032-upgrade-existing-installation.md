---
title: "Implementation Plan - US-032 Upgrade Existing Installation"
type: implementation-plan
user_story_id: US-032
feature_id: FEAT-004
workflow_phase: build
artifact_type: implementation-plan
tags:
  - helix/build
  - helix/artifact/implementation
  - installation
  - upgrade
  - self-update
related:
  - "[[US-032-upgrade-existing-installation]]"
  - "[[SD-004-cross-platform-installation]]"
  - "[[TS-004-installation-test-specification]]"
status: draft
priority: P1
created: 2025-01-22
updated: 2025-01-22
---

# Implementation Plan: US-032 Upgrade Existing Installation

## User Story Reference

**US-032**: As a developer with DDX already installed, I want to upgrade to newer versions so that I can access the latest features and bug fixes.

**Acceptance Criteria**:
- Self-update to latest version
- Self-update to specific version
- Preserve user configurations
- Rollback capability on failure
- Version compatibility checks

## Component Mapping

**Primary Component**: UpgradeManager (from SD-004)
**Supporting Components**:
- BinaryDistributor (download new versions)
- InstallationValidator (verify upgrades)
- InstallationManager (coordination)

## Implementation Strategy

### Overview
Create a `ddx self-update` command that can upgrade DDX to newer versions while preserving user configurations and providing rollback capability.

### Key Changes
1. Create new CLI command `cli/cmd/self_update.go`
2. Implement version comparison and detection
3. Add backup and rollback functionality
4. Integrate with GitHub releases for version checking
5. Preserve user configurations during upgrade

## Detailed Implementation Steps

### Step 1: Create Self-Update Command

**File**: `cli/cmd/self_update.go`

**Implementation**:
```go
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// selfUpdateCmd represents the self-update command
var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update DDX to the latest version",
	Long: `Update DDX to the latest version or a specific version.

Examples:
  ddx self-update                    # Update to latest version
  ddx self-update --version v1.2.3  # Update to specific version
  ddx self-update --check           # Check for updates without installing
  ddx self-update --force           # Force update even if current version is latest`,
	RunE: selfUpdateCommand,
}

func init() {
	selfUpdateCmd.Flags().String("version", "", "Target version to update to (default: latest)")
	selfUpdateCmd.Flags().Bool("check", false, "Check for updates without installing")
	selfUpdateCmd.Flags().Bool("force", false, "Force update even if current version is latest")
	selfUpdateCmd.Flags().Bool("rollback", false, "Rollback to previous version")
	selfUpdateCmd.Flags().Bool("backup", true, "Create backup before update")
	selfUpdateCmd.Flags().Bool("dry-run", false, "Show what would be done without making changes")
}

func selfUpdateCommand(cmd *cobra.Command, args []string) error {
	targetVersion, _ := cmd.Flags().GetString("version")
	checkOnly, _ := cmd.Flags().GetBool("check")
	force, _ := cmd.Flags().GetBool("force")
	rollback, _ := cmd.Flags().GetBool("rollback")
	backup, _ := cmd.Flags().GetBool("backup")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	updater := NewSelfUpdater(cmd.OutOrStdout())

	opts := UpdateOptions{
		TargetVersion: targetVersion,
		CheckOnly:     checkOnly,
		Force:         force,
		Rollback:      rollback,
		CreateBackup:  backup,
		DryRun:        dryRun,
	}

	return updater.Update(opts)
}

// UpdateOptions contains options for self-update
type UpdateOptions struct {
	TargetVersion string // Target version (empty = latest)
	CheckOnly     bool   // Only check for updates
	Force         bool   // Force update even if current is latest
	Rollback      bool   // Rollback to previous version
	CreateBackup  bool   // Create backup before update
	DryRun        bool   // Show what would be done
}

// SelfUpdater manages DDX self-updates
type SelfUpdater struct {
	out io.Writer
}

// NewSelfUpdater creates a new self-updater
func NewSelfUpdater(out io.Writer) *SelfUpdater {
	return &SelfUpdater{out: out}
}

// Update performs the self-update process
func (su *SelfUpdater) Update(opts UpdateOptions) error {
	if opts.Rollback {
		return su.rollbackToPrevious()
	}

	fmt.Fprintln(su.out, "🔄 DDX Self-Update")
	fmt.Fprintln(su.out)

	// Get current version
	currentVersion, err := su.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	fmt.Fprintf(su.out, "📦 Current version: %s\n", currentVersion)

	// Get target version
	targetVersion := opts.TargetVersion
	if targetVersion == "" {
		targetVersion, err = su.getLatestVersion()
		if err != nil {
			return fmt.Errorf("failed to get latest version: %w", err)
		}
	}

	fmt.Fprintf(su.out, "🎯 Target version: %s\n", targetVersion)

	// Version comparison
	if currentVersion == targetVersion && !opts.Force {
		fmt.Fprintln(su.out, "✅ You are already using the latest version!")
		return nil
	}

	if opts.CheckOnly {
		if currentVersion != targetVersion {
			fmt.Fprintf(su.out, "📢 Update available: %s → %s\n", currentVersion, targetVersion)
			fmt.Fprintln(su.out, "Run 'ddx self-update' to upgrade")
		}
		return nil
	}

	if opts.DryRun {
		fmt.Fprintln(su.out, "🔍 DRY RUN - would perform the following actions:")
		fmt.Fprintf(su.out, "  1. Backup current binary (%s)\n", currentVersion)
		fmt.Fprintf(su.out, "  2. Download new binary (%s)\n", targetVersion)
		fmt.Fprintln(su.out, "  3. Replace current binary")
		fmt.Fprintln(su.out, "  4. Verify new installation")
		return nil
	}

	// Perform the update
	return su.performUpdate(currentVersion, targetVersion, opts)
}

func (su *SelfUpdater) performUpdate(currentVersion, targetVersion string, opts UpdateOptions) error {
	// Get current binary path
	currentBinary, err := su.getCurrentBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to locate current binary: %w", err)
	}

	fmt.Fprintf(su.out, "📍 Current binary: %s\n", currentBinary)

	// Create backup if requested
	var backupPath string
	if opts.CreateBackup {
		backupPath, err = su.createBackup(currentBinary, currentVersion)
		if err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		fmt.Fprintf(su.out, "💾 Backup created: %s\n", backupPath)
	}

	// Download new version
	fmt.Fprintln(su.out, "⬇️  Downloading new version...")
	newBinary, err := su.downloadVersion(targetVersion)
	if err != nil {
		return fmt.Errorf("failed to download new version: %w", err)
	}
	defer os.Remove(newBinary) // Cleanup temp file

	// Verify new binary
	fmt.Fprintln(su.out, "🔍 Verifying new binary...")
	if err := su.verifyBinary(newBinary, targetVersion); err != nil {
		return fmt.Errorf("binary verification failed: %w", err)
	}

	// Replace binary atomically
	fmt.Fprintln(su.out, "🔄 Installing new version...")
	if err := su.replaceBinary(currentBinary, newBinary); err != nil {
		// Attempt rollback if backup exists
		if backupPath != "" {
			fmt.Fprintln(su.out, "❌ Update failed, attempting rollback...")
			if rollbackErr := su.restoreBackup(backupPath, currentBinary); rollbackErr != nil {
				return fmt.Errorf("update failed and rollback failed: %w (original error: %v)", rollbackErr, err)
			}
			fmt.Fprintln(su.out, "✅ Rollback completed")
		}
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	// Verify installation
	fmt.Fprintln(su.out, "✅ Verifying installation...")
	if err := su.verifyInstallation(targetVersion); err != nil {
		fmt.Fprintf(su.out, "⚠️  Installation verification failed: %v\n", err)
		fmt.Fprintln(su.out, "You may need to run 'ddx doctor' to doctor issues")
	}

	fmt.Fprintln(su.out)
	fmt.Fprintf(su.out, "🎉 Successfully updated DDX from %s to %s!\n", currentVersion, targetVersion)

	// Store backup info for potential rollback
	if backupPath != "" {
		su.storeBackupInfo(backupPath, currentVersion)
	}

	return nil
}

// Version and binary management methods
func (su *SelfUpdater) getCurrentVersion() (string, error) {
	// This would normally extract version from the binary or version command
	// For now, return a placeholder
	return "v1.0.0", nil
}

func (su *SelfUpdater) getLatestVersion() (string, error) {
	// Query GitHub API for latest release
	// For now, return a placeholder
	return "v1.1.0", nil
}

func (su *SelfUpdater) getCurrentBinaryPath() (string, error) {
	return exec.LookPath("ddx")
}

func (su *SelfUpdater) downloadVersion(version string) (string, error) {
	// Detect platform
	platform := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go architecture names to release names
	switch arch {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	// Construct download URL
	var extension, binaryName string
	if platform == "windows" {
		extension = "zip"
		binaryName = "ddx.exe"
	} else {
		extension = "tar.gz"
		binaryName = "ddx"
	}

	archiveName := fmt.Sprintf("ddx-%s-%s.%s", platform, arch, extension)
	downloadURL := fmt.Sprintf("https://github.com/easel/ddx/releases/download/%s/%s", version, archiveName)

	// Download to temp file
	tempDir := os.TempDir()
	tempArchive := filepath.Join(tempDir, archiveName)

	fmt.Fprintf(su.out, "📥 Downloading from: %s\n", downloadURL)

	// Use HTTP client to download (implementation details omitted for brevity)
	if err := su.downloadFile(downloadURL, tempArchive); err != nil {
		return "", err
	}

	// Extract binary
	extractedBinary := filepath.Join(tempDir, "ddx-new")
	if err := su.extractBinary(tempArchive, extractedBinary, binaryName); err != nil {
		return "", err
	}

	return extractedBinary, nil
}

func (su *SelfUpdater) verifyBinary(binaryPath, expectedVersion string) error {
	// Run version command on new binary
	cmd := exec.Command(binaryPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("binary is not executable: %w", err)
	}

	version := strings.TrimSpace(string(output))
	if !strings.Contains(version, expectedVersion) {
		return fmt.Errorf("version mismatch: expected %s, got %s", expectedVersion, version)
	}

	return nil
}

func (su *SelfUpdater) createBackup(currentBinary, version string) (string, error) {
	backupDir := filepath.Join(filepath.Dir(currentBinary), ".ddx-backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	backupName := fmt.Sprintf("ddx-%s.backup", version)
	backupPath := filepath.Join(backupDir, backupName)

	return backupPath, su.copyFile(currentBinary, backupPath)
}

func (su *SelfUpdater) replaceBinary(oldPath, newPath string) error {
	// Get permissions of old binary
	info, err := os.Stat(oldPath)
	if err != nil {
		return err
	}

	// Copy new binary to temporary location
	tempPath := oldPath + ".new"
	if err := su.copyFile(newPath, tempPath); err != nil {
		return err
	}

	// Set correct permissions
	if err := os.Chmod(tempPath, info.Mode()); err != nil {
		os.Remove(tempPath)
		return err
	}

	// Atomic replacement (platform-specific implementation needed for Windows)
	if runtime.GOOS == "windows" {
		// Windows requires special handling for replacing running executables
		return su.replaceOnWindows(oldPath, tempPath)
	} else {
		return os.Rename(tempPath, oldPath)
	}
}

func (su *SelfUpdater) replaceOnWindows(oldPath, newPath string) error {
	// Windows-specific replacement logic
	// This might involve using a helper process or delayed replacement
	// For now, implement a simple approach
	backup := oldPath + ".old"
	if err := os.Rename(oldPath, backup); err != nil {
		return err
	}

	if err := os.Rename(newPath, oldPath); err != nil {
		// Restore backup on failure
		os.Rename(backup, oldPath)
		return err
	}

	// Schedule cleanup of backup file
	os.Remove(backup)
	return nil
}

func (su *SelfUpdater) verifyInstallation(expectedVersion string) error {
	// Run version command to verify
	cmd := exec.Command("ddx", "version")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	version := strings.TrimSpace(string(output))
	if !strings.Contains(version, expectedVersion) {
		return fmt.Errorf("verification failed: expected %s, got %s", expectedVersion, version)
	}

	return nil
}

func (su *SelfUpdater) rollbackToPrevious() error {
	fmt.Fprintln(su.out, "🔄 Rolling back to previous version...")

	// Find the most recent backup
	backupInfo, err := su.getBackupInfo()
	if err != nil {
		return fmt.Errorf("no backup found for rollback: %w", err)
	}

	currentBinary, err := su.getCurrentBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to locate current binary: %w", err)
	}

	// Restore from backup
	if err := su.restoreBackup(backupInfo.Path, currentBinary); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Fprintf(su.out, "✅ Rolled back to version %s\n", backupInfo.Version)
	return nil
}

// Utility methods (implementations omitted for brevity)
func (su *SelfUpdater) downloadFile(url, dest string) error {
	// HTTP download implementation
	return nil
}

func (su *SelfUpdater) extractBinary(archive, dest, binaryName string) error {
	// Archive extraction implementation
	return nil
}

func (su *SelfUpdater) copyFile(src, dest string) error {
	// File copying implementation
	return nil
}

type BackupInfo struct {
	Path    string
	Version string
}

func (su *SelfUpdater) storeBackupInfo(path, version string) error {
	// Store backup metadata for rollback
	return nil
}

func (su *SelfUpdater) getBackupInfo() (*BackupInfo, error) {
	// Retrieve backup metadata
	return nil, nil
}

func (su *SelfUpdater) restoreBackup(backupPath, currentPath string) error {
	return su.copyFile(backupPath, currentPath)
}
```

### Step 2: Register Self-Update Command

**File**: `cli/cmd/command_factory.go`

**Integration**:
```go
func (f *CommandFactory) NewRootCommand() *cobra.Command {
	// ... existing code ...

	// Add self-update command
	rootCmd.AddCommand(selfUpdateCmd)

	// ... rest of commands ...
}
```

### Step 3: Version Management

**File**: `cli/cmd/version.go` (enhance existing)

**Add version comparison**:
```go
// Add version checking functionality
func checkForUpdates() {
	// Query GitHub API for latest version
	// Compare with current version
	// Display update notification
}

func init() {
	versionCmd.Flags().Bool("check", false, "Check for available updates")
}
```

## Integration Points

### Test Coverage
- **Primary Test**: `TestAcceptance_US032_UpgradeExistingInstallation`
- **Test Scenarios**:
  - Self-update to latest version
  - Self-update to specific version
  - Configuration preservation
  - Rollback functionality
- **Test Location**: `cli/cmd/installation_acceptance_test.go:403-456`

### Related Components
- **US-028**: Updated installation scripts create release assets
- **US-030**: Doctor command verifies successful upgrades
- **US-035**: Enhanced diagnostics detect upgrade issues

### Dependencies
- GitHub releases with proper versioning
- Network access for version checking and downloads
- Write permissions to binary location
- Backup storage space

## Success Criteria

✅ **Implementation Complete When**:
1. `TestAcceptance_US032_UpgradeExistingInstallation` passes
2. `ddx self-update` successfully upgrades to latest version
3. `ddx self-update --version X.Y.Z` upgrades to specific version
4. User configurations are preserved during upgrade
5. Rollback functionality works when upgrade fails
6. Version checking works (`ddx self-update --check`)

## Risk Mitigation

### High-Risk Areas
1. **Binary Replacement**: Atomic operations to prevent corruption
2. **Configuration Loss**: Backup and preserve user data
3. **Network Failures**: Graceful handling with retry logic
4. **Permission Issues**: Clear error messages and fallback options

### Edge Cases
1. **Running Binary Replacement**: Special handling for Windows
2. **Corrupted Downloads**: Checksum verification
3. **Incompatible Versions**: Version compatibility checks
4. **Disk Space**: Check available space before download

---

This implementation provides safe, reliable self-update functionality with comprehensive backup and rollback capabilities.