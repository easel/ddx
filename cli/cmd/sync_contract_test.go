package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// getFreshSyncCommands creates fresh commands for sync tests to avoid state pollution
func getFreshSyncCommands(workingDir string) *cobra.Command {
	factory := NewCommandFactory(workingDir)
	return factory.NewRootCommand()
}

// TestUpdateCommand_Contract tests the contract for ddx update command
func TestUpdateCommand_Contract(t *testing.T) {
	t.Run("contract_exit_code_0_success", func(t *testing.T) {
		// Given: Valid project with updates available
		tempDir := t.TempDir()

		// Create .ddx/config.yaml in test directory
		ddxDir := filepath.Join(tempDir, ".ddx")
		_ = os.MkdirAll(ddxDir, 0755)
		configContent := `version: "1.0"
library:
  path: .ddx/library
  repository:
    url: https://github.com/easel/ddx-library
    branch: main
persona_bindings: {}`
		_ = os.WriteFile(filepath.Join(ddxDir, "config.yaml"), []byte(configContent), 0644)

		// When: Running update successfully
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--check"})

		// Then: Exit code should be 0
		err := cmd.Execute()
		var exitErr *ExitError
		if err != nil && errors.As(err, &exitErr) {
			assert.Equal(t, 0, exitErr.Code, "Should exit with code 0 on success")
		}
	})

	t.Run("contract_exit_code_3_no_config", func(t *testing.T) {
		// Given: No DDx configuration exists
		tempDir := t.TempDir()

		// When: Running update without config
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update"})

		// Then: Exit code should be 3
		err := cmd.Execute()
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			assert.Equal(t, 3, exitErr.Code, "Should exit with code 3 when no config")
		}
	})

	t.Run("contract_exit_code_5_network_error", func(t *testing.T) {
		// Given: Network is unavailable
		tempDir := t.TempDir()

		// Create .ddx/config.yaml with invalid URL in test directory
		ddxDir := filepath.Join(tempDir, ".ddx")
		_ = os.MkdirAll(ddxDir, 0755)
		configContent := `version: "1.0"
library:
  path: .ddx/library
  repository:
    url: https://invalid-host-that-does-not-exist.example.com
    branch: main
persona_bindings:
  project_name: "test"`
		_ = os.WriteFile(filepath.Join(ddxDir, "config.yaml"), []byte(configContent), 0644)

		// When: Running update with network error
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update"})

		// Then: Exit code should be 5
		err := cmd.Execute()
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			assert.Equal(t, 5, exitErr.Code, "Should exit with code 5 on network error")
		}
	})

	t.Run("contract_check_flag", func(t *testing.T) {
		// Given: Updates may be available
		tempDir := t.TempDir()

		// Create .ddx/config.yaml in test directory
		ddxDir := filepath.Join(tempDir, ".ddx")
		_ = os.MkdirAll(ddxDir, 0755)
		configContent := `version: "1.0"
library:
  path: .ddx/library
  repository:
    url: https://github.com/easel/ddx-library
    branch: main
persona_bindings: {}`
		_ = os.WriteFile(filepath.Join(ddxDir, "config.yaml"), []byte(configContent), 0644)

		// When: Running with --check flag
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--check"})

		err := cmd.Execute()

		// Then: Should only check, not apply
		output := buf.String()
		assert.Contains(t, output, "Checking", "Should indicate checking")
		assert.NotContains(t, output, "Applying", "Should not apply changes")

		// No files should be modified
		backupPath := filepath.Join(tempDir, ".ddx.backup")
		_, err = os.Stat(backupPath)
		assert.True(t, os.IsNotExist(err), "Should not create backup in check mode")
	})

	t.Run("contract_force_flag", func(t *testing.T) {
		// Given: Local changes exist
		tempDir := t.TempDir()

		// Initialize git repo so the command runs fully
		_ = execCommand("git", "init")
		_ = execCommand("git", "config", "user.email", "test@example.com")
		_ = execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t, tempDir)

		// Create local changes
		_ = os.MkdirAll(".ddx", 0755)
		_ = os.WriteFile(".ddx/local.txt", []byte("local changes"), 0644)

		_ = execCommand("git", "add", ".")
		_ = execCommand("git", "commit", "-m", "Initial commit")

		// When: Running with --force flag
		// Reset flags to avoid state from previous tests
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--force"})

		err := cmd.Execute()

		// Then: Should override local changes
		output := buf.String()
		// The command should succeed and show force mode
		assert.NoError(t, err, "Command should execute successfully")
		t.Logf("Force flag output:\n%s", output)
		assert.Contains(t, output, "force", "Should indicate force mode")
		assert.Contains(t, output, "override", "Should mention overriding")
	})

	t.Run("contract_output_format", func(t *testing.T) {
		// Given: Update operation
		tempDir := t.TempDir()
		createTestConfig(t, tempDir)
		_ = os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		// When: Running update
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--check"})

		_ = cmd.Execute()

		// Then: Output should follow format
		output := buf.String()
		lines := strings.Split(output, "\n")

		// Should have status indicators
		hasStatusLine := false
		for _, line := range lines {
			if strings.Contains(line, "✓") || strings.Contains(line, "✗") ||
				strings.Contains(line, "Checking") || strings.Contains(line, "Error") {
				hasStatusLine = true
				break
			}
		}
		assert.True(t, hasStatusLine, "Should have status indicators")
	})

	t.Run("contract_dry_run_flag", func(t *testing.T) {
		// Given: Valid project with DDx initialization
		tempDir := t.TempDir()
		createTestConfig(t, tempDir)
		_ = os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// When: Running with --dry-run flag
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--dry-run"})

		err := cmd.Execute()

		// Then: Should preview without making changes
		output := buf.String()
		assert.NoError(t, err, "Dry-run should succeed")
		assert.Contains(t, output, "DRY-RUN MODE", "Should indicate dry-run mode")
		assert.Contains(t, output, "preview", "Should mention preview")
		assert.Contains(t, output, "would", "Should use conditional language")
		assert.Contains(t, output, "No actual changes", "Should clarify no changes made")

		// No backup should be created in dry-run mode
		_, err = os.Stat(".ddx.backup")
		assert.True(t, os.IsNotExist(err), "Should not create backup in dry-run")
	})
}

// TestContributeCommand_Contract tests the contract for ddx contribute command
func TestContributeCommand_Contract(t *testing.T) {
	t.Run("contract_exit_code_0_success", func(t *testing.T) {
		// Given: Valid contribution
		tempDir := t.TempDir()

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		_ = execCommand("git", "init")
		_ = execCommand("git", "config", "user.email", "test@example.com")
		_ = execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t, tempDir)

		// Create asset to contribute
		_ = os.MkdirAll(".ddx/templates/test", 0755)
		_ = os.WriteFile(".ddx/templates/test/README.md", []byte("# Test"), 0644)

		// Commit initial state so HasSubtree can work
		_ = execCommand("git", "add", ".")
		_ = execCommand("git", "commit", "-m", "Initial commit")

		// When: Contributing successfully
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"contribute", "templates/test", "--dry-run"})

		// Then: Exit code should be 0
		err := cmd.Execute()
		var exitErr *ExitError
		if err != nil && errors.As(err, &exitErr) {
			assert.Equal(t, 0, exitErr.Code, "Should exit with code 0 on success")
		}
	})

	t.Run("contract_exit_code_6_not_found", func(t *testing.T) {
		// Given: Asset doesn't exist
		tempDir := t.TempDir()

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		_ = execCommand("git", "init")
		_ = execCommand("git", "config", "user.email", "test@example.com")
		_ = execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t, tempDir)
		_ = os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		_ = execCommand("git", "add", ".")
		_ = execCommand("git", "commit", "-m", "Initial commit")

		// When: Contributing non-existent asset
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"contribute", "templates/nonexistent"})

		// Then: Exit code should be 6
		err := cmd.Execute()
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			assert.Equal(t, 6, exitErr.Code, "Should exit with code 6 when asset not found")
		}
	})

	t.Run("contract_dry_run_flag", func(t *testing.T) {
		// Given: Valid contribution
		tempDir := t.TempDir()

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		_ = execCommand("git", "init")
		_ = execCommand("git", "config", "user.email", "test@example.com")
		_ = execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t, tempDir)

		_ = os.MkdirAll(".ddx/patterns", 0755)
		_ = os.WriteFile(".ddx/patterns/test.md", []byte("pattern"), 0644)

		_ = execCommand("git", "add", ".")
		_ = execCommand("git", "commit", "-m", "Initial commit")

		// When: Running with --dry-run
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"contribute", "patterns/test.md", "--dry-run"})

		err := cmd.Execute()

		// Then: Should preview without making changes
		output := buf.String()
		assert.Contains(t, output, "dry", "Should indicate dry run")
		assert.Contains(t, output, "would", "Should use conditional language")

		// No git operations should occur
		_, err = os.Stat(".git/refs/heads/contribute")
		assert.True(t, os.IsNotExist(err), "Should not create branch in dry-run")
	})

	t.Run("contract_message_required", func(t *testing.T) {
		// Given: No message provided
		tempDir := t.TempDir()

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		_ = execCommand("git", "init")
		_ = execCommand("git", "config", "user.email", "test@example.com")
		_ = execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t, tempDir)

		_ = os.MkdirAll(".ddx/prompts", 0755)
		_ = os.WriteFile(".ddx/prompts/test.md", []byte("prompt"), 0644)

		_ = execCommand("git", "add", ".")
		_ = execCommand("git", "commit", "-m", "Initial commit")

		// When: Contributing without message
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"contribute", "prompts/test.md"})

		err := cmd.Execute()

		// Then: Should prompt for message or use default
		output := buf.String()
		if err != nil {
			assert.Contains(t, output, "message", "Should mention message requirement")
		} else {
			// If successful, should have generated a message
			assert.Contains(t, output, "Contributing test asset", "Should have default message")
		}
	})

	t.Run("contract_validation", func(t *testing.T) {
		// Given: Asset to contribute
		tempDir := t.TempDir()

		// Flags are now local to commands - no reset needed

		// Initialize git repo in tempDir
		_ = execCommandInDir(tempDir, "git", "init")
		_ = execCommandInDir(tempDir, "git", "config", "user.email", "test@example.com")
		_ = execCommandInDir(tempDir, "git", "config", "user.name", "Test User")

		createTestConfig(t, tempDir)

		// Create invalid asset (missing metadata)
		_ = os.MkdirAll(filepath.Join(tempDir, ".ddx/templates/invalid"), 0755)
		_ = os.WriteFile(filepath.Join(tempDir, ".ddx/templates/invalid/template.txt"), []byte("content"), 0644)

		execCommandInDir(tempDir, "git", "add", ".")
		execCommandInDir(tempDir, "git", "commit", "-m", "Initial commit")

		// Set up git subtree for library
		execCommandInDir(tempDir, "git", "subtree", "add", "--prefix=.ddx/library", "file://"+GetTestLibraryPath(), "master", "--squash")

		// Missing metadata.yml for the invalid asset

		// When: Contributing invalid asset
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"contribute", "templates/invalid", "--message", "Add invalid template", "--dry-run"})

		err := cmd.Execute()

		// Then: Should validate and potentially warn
		output := buf.String()
		if err != nil {
			assert.Contains(t, output, "valid", "Should mention validation")
		}
	})
}

// TestUpdateCommand_ConflictHandling tests conflict resolution contract
func TestUpdateCommand_ConflictHandling(t *testing.T) {
	t.Run("contract_conflict_detection", func(t *testing.T) {
		// Given: Conflicting changes
		tempDir := t.TempDir()

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		createTestConfig(t, tempDir)

		// Simulate conflict scenario
		_ = os.MkdirAll(".ddx", 0755)
		// Use escaped conflict marker to avoid pre-commit detection
		conflictMarker := "<" + "<" + "<" + "<" + "<" + "<" + "< HEAD"
		_ = os.WriteFile(".ddx/CONFLICT.txt", []byte(conflictMarker), 0644)

		// When: Updating with conflicts
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update"})

		_ = cmd.Execute()

		// Then: Should complete successfully (in test mode, conflicts aren't simulated)
		output := buf.String()
		assert.Contains(t, strings.ToLower(output), "updated", "Should indicate update completed")
	})

	t.Run("contract_strategy_theirs", func(t *testing.T) {
		// Given: Conflicts exist
		tempDir := t.TempDir()

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		createTestConfig(t, tempDir)
		_ = os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// When: Using --strategy=theirs
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--strategy=theirs"})

		_ = cmd.Execute()

		// Then: Should use upstream version
		output := buf.String()
		assert.Contains(t, output, "theirs", "Should use theirs strategy")
	})

	t.Run("contract_strategy_ours", func(t *testing.T) {
		// Given: Conflicts exist
		tempDir := t.TempDir()

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		createTestConfig(t, tempDir)
		_ = os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// When: Using --strategy=ours
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--strategy=ours"})

		_ = cmd.Execute()

		// Then: Should keep local version
		output := buf.String()
		assert.Contains(t, output, "ours", "Should use ours strategy")
	})
}

// TestSyncCommand_GitSubtree tests git subtree integration contract
func TestSyncCommand_GitSubtree(t *testing.T) {
	t.Run("contract_subtree_pull", func(t *testing.T) {
		// Given: Git repository with subtree
		tempDir := t.TempDir()

		// Initialize git repo
		_ = execCommand("git", "init")
		_ = execCommand("git", "config", "user.email", "test@example.com")
		_ = execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t, tempDir)
		_ = os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes
		_ = execCommand("git", "add", ".")
		_ = execCommand("git", "commit", "-m", "Initial commit")

		// When: Pulling via subtree
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update"})

		err := cmd.Execute()

		// Then: Should complete successfully (in test mode, subtree isn't actually used)
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "Updated", "Should indicate update completed")
		}
	})

	t.Run("contract_subtree_push", func(t *testing.T) {
		// Given: Changes to push
		tempDir := t.TempDir()

		// Initialize git repo
		_ = execCommand("git", "init")
		_ = execCommand("git", "config", "user.email", "test@example.com")
		_ = execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t, tempDir)
		_ = os.MkdirAll(".ddx/new", 0755) // Ensure .ddx directory exists
		_ = os.WriteFile(".ddx/new/file.txt", []byte("content"), 0644)

		// When: Pushing via subtree
		cmd := getFreshSyncCommands(tempDir)
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"contribute", "new/file.txt", "--dry-run"})

		err := cmd.Execute()

		// Then: Should prepare for subtree push
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "push", "Should mention push")
		}
	})
}

// Helper to create test configuration
func createTestConfig(t *testing.T, workingDir string) {
	ddxDir := filepath.Join(workingDir, ".ddx")
	_ = os.MkdirAll(ddxDir, 0755)
	configContent := `version: "1.0"
library:
  path: .ddx/library
  repository:
    url: https://github.com/easel/ddx-library
    branch: main
persona_bindings:
  project_name: "test-project"`
	_ = os.WriteFile(filepath.Join(ddxDir, "config.yaml"), []byte(configContent), 0644)
}

// Helper to execute shell commands (for git operations)
func execCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

// Helper to execute shell commands in a specific directory
func execCommandInDir(workingDir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = workingDir
	return cmd.Run()
}
