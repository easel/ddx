package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAcceptance_US013_RollbackProblematicUpdates tests US-013: Rollback Problematic Updates
func TestAcceptance_US013_RollbackProblematicUpdates(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		given    func(t *testing.T) string                 // Setup conditions
		when     func(t *testing.T, dir string) error      // Execute action
		then     func(t *testing.T, dir string, err error) // Verify outcome
	}{
		{
			name:     "simple_rollback",
			scenario: "Rollback to previous working version when a problematic update was applied",
			given: func(t *testing.T) string {
				// Given: a problematic update was applied
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Setup initial state
				createTestProject(t, workDir)
				createBackupPoint(t, workDir, "v1.0.0")

				// Simulate problematic state
				createProblematicState(t, workDir)

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx rollback`
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "rollback")
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: the system reverts to the previous working version
				assert.NoError(t, err)
				assertSystemAtVersion(t, dir, "v1.0.0")
				assertProblematicStateRemoved(t, dir)
			},
		},
		{
			name:     "specific_version_rollback",
			scenario: "Rollback to a specific version when multiple versions exist",
			given: func(t *testing.T) string {
				// Given: multiple versions exist
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				createTestProject(t, workDir)
				createBackupPoint(t, workDir, "v1.0.0")
				createBackupPoint(t, workDir, "v1.1.0")
				createBackupPoint(t, workDir, "v1.2.0")

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx rollback --to v1.1.0`
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "rollback", "--to", "v1.1.0")
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: I can rollback to the specific version
				assert.NoError(t, err)
				assertSystemAtVersion(t, dir, "v1.1.0")
			},
		},
		{
			name:     "list_rollback_options",
			scenario: "View available rollback points",
			given: func(t *testing.T) string {
				// Given: I want to see options
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				createTestProject(t, workDir)
				createBackupPoint(t, workDir, "v1.0.0")
				createBackupPoint(t, workDir, "v1.1.0")

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx rollback --list`
				rootCmd := getTestRootCommand()
				output, err := executeCommand(rootCmd, "rollback", "--list")

				// Store output for verification
				t.Setenv("TEST_OUTPUT", output)
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: available rollback points are displayed
				assert.NoError(t, err)
				output := os.Getenv("TEST_OUTPUT")
				assert.Contains(t, output, "v1.0.0")
				assert.Contains(t, output, "v1.1.0")
				assert.Contains(t, output, "TIMESTAMP")
				assert.Contains(t, output, "DESCRIPTION")
			},
		},
		{
			name:     "rollback_preview",
			scenario: "Preview changes before rollback",
			given: func(t *testing.T) string {
				// Given: I'm unsure about rollback
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				createTestProject(t, workDir)
				createBackupPoint(t, workDir, "v1.0.0")
				makeCurrentChanges(t, workDir)

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx rollback --preview`
				rootCmd := getTestRootCommand()
				output, err := executeCommand(rootCmd, "rollback", "--preview")

				t.Setenv("TEST_OUTPUT", output)
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: I see what changes will be reverted
				assert.NoError(t, err)
				output := os.Getenv("TEST_OUTPUT")
				assert.Contains(t, output, "Changes that will be reverted")
				assert.Contains(t, output, "preview")
				// Should not actually rollback in preview mode
				assertCurrentChangesStillExist(t, dir)
			},
		},
		{
			name:     "backup_before_rollback",
			scenario: "Create backup before rollback operation",
			given: func(t *testing.T) string {
				// Given: I initiate rollback
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				createTestProject(t, workDir)
				createBackupPoint(t, workDir, "v1.0.0")
				makeCurrentChanges(t, workDir)

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: rollback starts
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "rollback")
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: a backup is created before the rollback
				assert.NoError(t, err)
				assertBackupCreatedBeforeRollback(t, dir)
			},
		},
		{
			name:     "state_validation_after_rollback",
			scenario: "Validate system integrity after rollback",
			given: func(t *testing.T) string {
				// Given: rollback completes
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				createTestProject(t, workDir)
				createBackupPoint(t, workDir, "v1.0.0")

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I check state after rollback
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "rollback")
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: the system validates integrity after rollback
				assert.NoError(t, err)
				assertSystemIntegrityValidated(t, dir)
			},
		},
		{
			name:     "rollback_failure_recovery",
			scenario: "Handle rollback failure with clear recovery instructions",
			given: func(t *testing.T) string {
				// Given: rollback will fail
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				createTestProject(t, workDir)
				createCorruptedBackupPoint(t, workDir, "v1.0.0")

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: error occurs during rollback
				rootCmd := getTestRootCommand()
				output, err := executeCommand(rootCmd, "rollback")

				t.Setenv("TEST_OUTPUT", output)
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: clear recovery instructions are provided
				assert.Error(t, err)
				output := os.Getenv("TEST_OUTPUT")
				assert.Contains(t, output, "Recovery Instructions")
				assert.Contains(t, output, "backup integrity")
				assert.Contains(t, output, "manual restoration")
			},
		},
		{
			name:     "no_backup_available",
			scenario: "Handle case when no backup points exist",
			given: func(t *testing.T) string {
				// Given: no backup points exist
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				createTestProject(t, workDir)
				// Don't create any backup points

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx rollback`
				rootCmd := getTestRootCommand()
				output, err := executeCommand(rootCmd, "rollback")

				t.Setenv("TEST_OUTPUT", output)
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: clear message about no backup points
				assert.Error(t, err)
				output := os.Getenv("TEST_OUTPUT")
				assert.Contains(t, output, "no backup points")
				assert.Contains(t, output, "available")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore working directory
			origDir, _ := os.Getwd()
			defer os.Chdir(origDir)

			dir := tt.given(t)
			err := tt.when(t, dir)
			tt.then(t, dir, err)
		})
	}
}

// Helper functions for test setup and verification

func createTestProject(t *testing.T, dir string) {
	// Create a basic DDx project structure
	ddxContent := `name: test-project
version: 1.0.0
library:
  repository: "https://github.com/example/ddx"
  branch: "main"
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ddx.yml"), []byte(ddxContent), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, ".ddx"), 0755))
}

func createBackupPoint(t *testing.T, dir, version string) {
	// Create a backup point directory structure
	backupDir := filepath.Join(dir, ".ddx", "backups", version)
	require.NoError(t, os.MkdirAll(backupDir, 0755))

	// Create backup metadata
	metadata := `version: ` + version + `
timestamp: 2025-01-20T10:00:00Z
description: Backup before update to ` + version + `
`
	require.NoError(t, os.WriteFile(filepath.Join(backupDir, "metadata.yml"), []byte(metadata), 0644))

	// Create some backed up content - full config format
	backupConfig := `name: test-project
version: ` + version + `
library:
  repository: "https://github.com/example/ddx"
  branch: "main"
`
	require.NoError(t, os.WriteFile(filepath.Join(backupDir, ".ddx.yml"), []byte(backupConfig), 0644))
}

func createProblematicState(t *testing.T, dir string) {
	// Create state that represents a problematic update
	problematicContent := `name: test-project
version: 2.0.0-broken
library:
  repository: "https://github.com/example/ddx"
  branch: "broken-branch"
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ddx.yml"), []byte(problematicContent), 0644))
}

func createCorruptedBackupPoint(t *testing.T, dir, version string) {
	// Create a backup point that will fail during restoration
	backupDir := filepath.Join(dir, ".ddx", "backups", version)
	require.NoError(t, os.MkdirAll(backupDir, 0755))

	// Create valid metadata but corrupt backup content that will fail restoration
	metadata := `version: ` + version + `
timestamp: 2025-01-20T10:00:00Z
description: Corrupted backup for testing failure recovery
`
	require.NoError(t, os.WriteFile(filepath.Join(backupDir, "metadata.yml"), []byte(metadata), 0644))

	// Create a backup file that's corrupted (missing required content)
	require.NoError(t, os.WriteFile(filepath.Join(backupDir, ".ddx.yml"), []byte("corrupted content"), 0644))
}

func makeCurrentChanges(t *testing.T, dir string) {
	// Make some current changes that would be lost in rollback
	newContent := `name: test-project
version: 1.1.0
library:
  repository: "https://github.com/example/ddx"
  branch: "main"
new_field: "new value"
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".ddx.yml"), []byte(newContent), 0644))
}

func assertSystemAtVersion(t *testing.T, dir, expectedVersion string) {
	// Verify system is at expected version
	configPath := filepath.Join(dir, ".ddx.yml")
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "version: "+expectedVersion)
}

func assertProblematicStateRemoved(t *testing.T, dir string) {
	// Verify problematic state is no longer present
	configPath := filepath.Join(dir, ".ddx.yml")
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.NotContains(t, string(content), "broken")
	assert.NotContains(t, string(content), "broken-branch")
}

func assertCurrentChangesStillExist(t *testing.T, dir string) {
	// Verify current changes still exist (for preview mode)
	configPath := filepath.Join(dir, ".ddx.yml")
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "new_field")
}

func assertBackupCreatedBeforeRollback(t *testing.T, dir string) {
	// Verify backup was created before rollback operation
	backupDirs, err := os.ReadDir(filepath.Join(dir, ".ddx", "backups"))
	require.NoError(t, err)

	// Should have at least one backup (original + pre-rollback backup)
	assert.GreaterOrEqual(t, len(backupDirs), 2)

	// Check for pre-rollback backup
	foundPreRollbackBackup := false
	for _, entry := range backupDirs {
		if strings.Contains(entry.Name(), "pre-rollback") || strings.Contains(entry.Name(), "before-rollback") {
			foundPreRollbackBackup = true
			break
		}
	}
	assert.True(t, foundPreRollbackBackup, "Should create backup before rollback")
}

func assertSystemIntegrityValidated(t *testing.T, dir string) {
	// Verify system integrity is validated after rollback
	// Check that .ddx.yml is valid and properly formatted
	configPath := filepath.Join(dir, ".ddx.yml")
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)

	// Basic validation - should be valid YAML
	assert.True(t, len(content) > 0)
	assert.Contains(t, string(content), "name:")
	assert.Contains(t, string(content), "version:")

	// Verify .ddx directory structure exists
	assert.DirExists(t, filepath.Join(dir, ".ddx"))
	assert.DirExists(t, filepath.Join(dir, ".ddx", "backups"))
}

// TestRollbackCommand_Help tests the rollback command help output
func TestRollbackCommand_Help(t *testing.T) {
	rootCmd := getTestRootCommand()
	output, err := executeCommand(rootCmd, "rollback", "--help")

	// Command should exist and provide help
	require.NoError(t, err)
	assert.Contains(t, output, "rollback")
	assert.Contains(t, output, "revert")
	assert.Contains(t, output, "--list")
	assert.Contains(t, output, "--to")
	assert.Contains(t, output, "--preview")
}

// TestRollbackCommand_Flags tests the rollback command flag parsing
func TestRollbackCommand_Flags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no_flags",
			args:    []string{"rollback"},
			wantErr: false, // Should work with defaults
		},
		{
			name:    "list_flag",
			args:    []string{"rollback", "--list"},
			wantErr: false,
		},
		{
			name:    "to_flag_with_value",
			args:    []string{"rollback", "--to", "v1.0.0"},
			wantErr: false,
		},
		{
			name:    "preview_flag",
			args:    []string{"rollback", "--preview"},
			wantErr: false,
		},
		{
			name:    "to_flag_without_value",
			args:    []string{"rollback", "--to"},
			wantErr: true,
		},
		{
			name:    "invalid_flag",
			args:    []string{"rollback", "--invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := getTestRootCommand()
			_, err := executeCommand(rootCmd, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				// Note: May error due to missing implementation, but not due to flag parsing
				// Check that error is not about unknown flags
				if err != nil {
					assert.NotContains(t, err.Error(), "unknown flag")
				}
			}
		})
	}
}
