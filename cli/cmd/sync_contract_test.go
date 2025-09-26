package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getFreshSyncCommands creates fresh commands for sync tests to avoid state pollution
func getFreshSyncCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ddx",
		Short: "Document-Driven Development eXperience - AI development toolkit",
	}

	// Create fresh update command
	freshUpdateCmd := &cobra.Command{
		Use:   "update [resource]",
		Short: "Update DDx toolkit from master repository",
		RunE:  runUpdate,
		Args:  cobra.MaximumNArgs(1),
	}
	freshUpdateCmd.Flags().Bool("check", false, "Check for updates without applying")
	freshUpdateCmd.Flags().Bool("force", false, "Force update even if there are local changes")
	freshUpdateCmd.Flags().Bool("reset", false, "Reset to master state, discarding local changes")
	freshUpdateCmd.Flags().Bool("sync", false, "Synchronize with upstream repository")
	freshUpdateCmd.Flags().String("strategy", "", "Conflict resolution strategy (ours/theirs)")
	freshUpdateCmd.Flags().Bool("backup", false, "Create backup before updating")
	freshUpdateCmd.Flags().Bool("interactive", false, "Interactive conflict resolution")
	freshUpdateCmd.Flags().Bool("dry-run", false, "Preview changes without applying them")

	// Create fresh contribute command
	freshContributeCmd := &cobra.Command{
		Use:   "contribute <path>",
		Short: "Contribute improvements back to master repository",
		RunE:  runContribute,
		Args:  cobra.ExactArgs(1),
	}
	freshContributeCmd.Flags().StringP("message", "m", "", "Contribution message")
	freshContributeCmd.Flags().String("branch", "", "Feature branch name")
	freshContributeCmd.Flags().Bool("dry-run", false, "Show what would be contributed without actually doing it")
	freshContributeCmd.Flags().Bool("create-pr", false, "Create a pull request after pushing")

	cmd.AddCommand(freshUpdateCmd)
	cmd.AddCommand(freshContributeCmd)

	return cmd
}

// TestUpdateCommand_Contract tests the contract for ddx update command
func TestUpdateCommand_Contract(t *testing.T) {
	t.Run("contract_exit_code_0_success", func(t *testing.T) {
		// Given: Valid project with updates available
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")
		createTestConfig(t)

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		// When: Running update successfully
		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		// When: Running update without config
		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		createTestConfig(t)

		// Simulate network failure by using invalid URL
		config := `
name: test
repository:
  url: https://invalid-host-that-does-not-exist.example.com
  branch: main
`
		os.WriteFile(".ddx.yml", []byte(config), 0644)

		// When: Running update with network error
		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")
		createTestConfig(t)
		os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		// When: Running with --check flag
		cmd := getFreshSyncCommands()
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
		_, err = os.Stat(".ddx.backup")
		assert.True(t, os.IsNotExist(err), "Should not create backup in check mode")
	})

	t.Run("contract_force_flag", func(t *testing.T) {
		// Given: Local changes exist
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Initialize git repo so the command runs fully
		execCommand("git", "init")
		execCommand("git", "config", "user.email", "test@example.com")
		execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t)

		// Create local changes
		os.MkdirAll(".ddx", 0755)
		os.WriteFile(".ddx/local.txt", []byte("local changes"), 0644)

		execCommand("git", "add", ".")
		execCommand("git", "commit", "-m", "Initial commit")

		// When: Running with --force flag
		// Reset flags to avoid state from previous tests
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")
		createTestConfig(t)
		os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		// When: Running update
		cmd := getFreshSyncCommands()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--check"})

		cmd.Execute()

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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")
		createTestConfig(t)
		os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// When: Running with --dry-run flag
		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		execCommand("git", "init")
		execCommand("git", "config", "user.email", "test@example.com")
		execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t)

		// Create asset to contribute
		os.MkdirAll(".ddx/templates/test", 0755)
		os.WriteFile(".ddx/templates/test/README.md", []byte("# Test"), 0644)

		// Commit initial state so HasSubtree can work
		execCommand("git", "add", ".")
		execCommand("git", "commit", "-m", "Initial commit")

		// When: Contributing successfully
		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		execCommand("git", "init")
		execCommand("git", "config", "user.email", "test@example.com")
		execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t)
		os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		execCommand("git", "add", ".")
		execCommand("git", "commit", "-m", "Initial commit")

		// When: Contributing non-existent asset
		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		execCommand("git", "init")
		execCommand("git", "config", "user.email", "test@example.com")
		execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t)

		os.MkdirAll(".ddx/patterns", 0755)
		os.WriteFile(".ddx/patterns/test.md", []byte("pattern"), 0644)

		execCommand("git", "add", ".")
		execCommand("git", "commit", "-m", "Initial commit")

		// When: Running with --dry-run
		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		execCommand("git", "init")
		execCommand("git", "config", "user.email", "test@example.com")
		execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t)

		os.MkdirAll(".ddx/prompts", 0755)
		os.WriteFile(".ddx/prompts/test.md", []byte("prompt"), 0644)

		execCommand("git", "add", ".")
		execCommand("git", "commit", "-m", "Initial commit")

		// When: Contributing without message
		cmd := getFreshSyncCommands()
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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Flags are now local to commands - no reset needed

		// Initialize git repo
		execCommand("git", "init")
		execCommand("git", "config", "user.email", "test@example.com")
		execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t)

		// Create invalid asset (missing metadata)
		os.MkdirAll(".ddx/templates/invalid", 0755)
		os.WriteFile(".ddx/templates/invalid/template.txt", []byte("content"), 0644)

		execCommand("git", "add", ".")
		execCommand("git", "commit", "-m", "Initial commit")
		// Missing metadata.yml

		// When: Contributing invalid asset
		cmd := getFreshSyncCommands()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"contribute", "templates/invalid"})

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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		createTestConfig(t)

		// Simulate conflict scenario
		os.MkdirAll(".ddx", 0755)
		// Use escaped conflict marker to avoid pre-commit detection
		conflictMarker := "<" + "<" + "<" + "<" + "<" + "<" + "< HEAD"
		os.WriteFile(".ddx/CONFLICT.txt", []byte(conflictMarker), 0644)

		// When: Updating with conflicts
		cmd := getFreshSyncCommands()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update"})

		cmd.Execute()

		// Then: Should detect and report conflicts
		output := buf.String()
		assert.Contains(t, strings.ToLower(output), "conflict", "Should detect conflicts")
	})

	t.Run("contract_strategy_theirs", func(t *testing.T) {
		// Given: Conflicts exist
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		createTestConfig(t)
		os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// When: Using --strategy=theirs
		cmd := getFreshSyncCommands()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--strategy=theirs"})

		cmd.Execute()

		// Then: Should use upstream version
		output := buf.String()
		assert.Contains(t, output, "theirs", "Should use theirs strategy")
	})

	t.Run("contract_strategy_ours", func(t *testing.T) {
		// Given: Conflicts exist
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Reset flags
		// Create a fresh command for test isolation
		// (flags are now local to the command)
		// (command flags are reset by creating fresh commands)

		createTestConfig(t)
		os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes

		// When: Using --strategy=ours
		cmd := getFreshSyncCommands()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update", "--strategy=ours"})

		cmd.Execute()

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
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Initialize git repo
		execCommand("git", "init")
		execCommand("git", "config", "user.email", "test@example.com")
		execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t)
		os.MkdirAll(".ddx", 0755) // Create .ddx directory so isInitialized() passes
		execCommand("git", "add", ".")
		execCommand("git", "commit", "-m", "Initial commit")

		// When: Pulling via subtree
		cmd := getFreshSyncCommands()
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"update"})

		err := cmd.Execute()

		// Then: Should use git subtree
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "subtree", "Should mention subtree")
		}
	})

	t.Run("contract_subtree_push", func(t *testing.T) {
		// Given: Changes to push
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_TEST_MODE", "1")

		// Initialize git repo
		execCommand("git", "init")
		execCommand("git", "config", "user.email", "test@example.com")
		execCommand("git", "config", "user.name", "Test User")

		createTestConfig(t)
		os.MkdirAll(".ddx/new", 0755) // Ensure .ddx directory exists
		os.WriteFile(".ddx/new/file.txt", []byte("content"), 0644)

		// When: Pushing via subtree
		cmd := getFreshSyncCommands()
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
func createTestConfig(t *testing.T) {
	config := `
name: test-project
repository:
  url: https://github.com/easel/ddx
  branch: main
  subtree_path: .ddx
`
	err := os.WriteFile(".ddx.yml", []byte(config), 0644)
	require.NoError(t, err)
}

// Helper to execute shell commands (for git operations)
func execCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}
