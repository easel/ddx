package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAcceptance_US004_UpdateAssetsFromMaster tests US-004: Update Assets from Master
func TestAcceptance_US004_UpdateAssetsFromMaster(t *testing.T) {
	t.Run("pull_latest_changes", func(t *testing.T) {
		// Given: A project with DDx initialized and updates available
		tempDir := t.TempDir()
		os.Chdir(tempDir)

		// Initialize DDx first
		initCmd := rootCmd
		initBuf := new(bytes.Buffer)
		initCmd.SetOut(initBuf)
		initCmd.SetErr(initBuf)
		initCmd.SetArgs([]string{"init", "--no-git"})

		err := initCmd.Execute()
		require.NoError(t, err, "Should initialize DDx")

		// When: User runs ddx update
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update"})

		err = updateCmd.Execute()

		// Then: Latest changes are fetched from master repository
		assert.NoError(t, err, "Update should complete successfully")
		output := updateBuf.String()
		assert.Contains(t, output, "Checking for updates", "Should show update check")
		assert.Contains(t, output, "Fetching latest changes", "Should fetch changes")
	})

	t.Run("display_changelog", func(t *testing.T) {
		// Given: Updates are available
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Running update command
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update", "--check"})

		err := updateCmd.Execute()

		// Then: Changelog is displayed
		assert.NoError(t, err)
		output := updateBuf.String()
		assert.Contains(t, output, "Available updates:", "Should show available updates")
		assert.Contains(t, output, "Changes since last update:", "Should show changelog")
	})

	t.Run("handle_merge_conflicts", func(t *testing.T) {
		// Given: Local changes conflict with upstream
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// Create a local modification that will conflict
		conflictFile := filepath.Join(".ddx", "templates", "test.md")
		os.MkdirAll(filepath.Dir(conflictFile), 0755)
		os.WriteFile(conflictFile, []byte("local changes"), 0644)

		// When: Updating with conflicts
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update"})

		_ = updateCmd.Execute()

		// Then: Conflicts are detected and resolution options provided
		output := updateBuf.String()
		assert.Contains(t, output, "conflict", "Should detect conflicts")
		assert.Contains(t, output, "resolution", "Should provide resolution options")
	})

	t.Run("selective_update", func(t *testing.T) {
		// Given: Multiple assets available for update
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Updating specific asset only
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update", "templates/nextjs"})

		err := updateCmd.Execute()

		// Then: Only specified asset is updated
		assert.NoError(t, err)
		output := updateBuf.String()
		assert.Contains(t, output, "Updating templates/nextjs", "Should update specific asset")
		assert.NotContains(t, output, "prompts", "Should not update other assets")
	})

	t.Run("preserve_local_changes", func(t *testing.T) {
		// Given: Local modifications exist
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		localFile := filepath.Join(".ddx", "custom.md")
		localContent := "my local customization"
		os.WriteFile(localFile, []byte(localContent), 0644)

		// When: Running update
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update"})

		err := updateCmd.Execute()

		// Then: Local changes are preserved
		assert.NoError(t, err)
		content, _ := os.ReadFile(localFile)
		assert.Equal(t, localContent, string(content), "Local changes should be preserved")
	})

	t.Run("force_update", func(t *testing.T) {
		// Given: Local changes exist that user wants to override
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Running update with --force
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update", "--force"})

		err := updateCmd.Execute()

		// Then: Updates are applied overriding local changes
		assert.NoError(t, err)
		output := updateBuf.String()
		assert.Contains(t, output, "Force updating", "Should indicate force update")
	})

	t.Run("create_backup", func(t *testing.T) {
		// Given: Project ready for update
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Running update
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update"})

		err := updateCmd.Execute()

		// Then: Backup is created before changes
		assert.NoError(t, err)
		output := updateBuf.String()
		assert.Contains(t, output, "Creating backup", "Should create backup")
		// Check for backup directory or files
		assert.DirExists(t, ".ddx.backup", "Backup directory should exist")
	})
}

// TestAcceptance_US005_ContributeImprovements tests US-005: Contribute Improvements
func TestAcceptance_US005_ContributeImprovements(t *testing.T) {
	t.Run("contribute_new_template", func(t *testing.T) {
		// Given: User has created a new template
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// Create new template
		templatePath := filepath.Join(".ddx", "templates", "my-template", "README.md")
		os.MkdirAll(filepath.Dir(templatePath), 0755)
		os.WriteFile(templatePath, []byte("# My Template"), 0644)

		// When: Running contribute command
		contributeCmd := rootCmd
		contributeBuf := new(bytes.Buffer)
		contributeCmd.SetOut(contributeBuf)
		contributeCmd.SetErr(contributeBuf)
		contributeCmd.SetArgs([]string{"contribute", "templates/my-template", "-m", "Add new template"})

		err := contributeCmd.Execute()

		// Then: Contribution is prepared
		assert.NoError(t, err)
		output := contributeBuf.String()
		assert.Contains(t, output, "Preparing contribution", "Should prepare contribution")
		assert.Contains(t, output, "templates/my-template", "Should reference the template")
	})

	t.Run("validate_contribution", func(t *testing.T) {
		// Given: User wants to contribute changes
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Contributing with validation
		contributeCmd := rootCmd
		contributeBuf := new(bytes.Buffer)
		contributeCmd.SetOut(contributeBuf)
		contributeCmd.SetErr(contributeBuf)
		contributeCmd.SetArgs([]string{"contribute", "patterns/test", "--dry-run"})

		_ = contributeCmd.Execute()

		// Then: Contribution is validated
		output := contributeBuf.String()
		assert.Contains(t, output, "Validating contribution", "Should validate")
		assert.Contains(t, output, "Validation passed", "Should show validation status")
	})

	t.Run("create_pull_request", func(t *testing.T) {
		// Given: Valid contribution ready
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Contributing with PR creation
		contributeCmd := rootCmd
		contributeBuf := new(bytes.Buffer)
		contributeCmd.SetOut(contributeBuf)
		contributeCmd.SetErr(contributeBuf)
		contributeCmd.SetArgs([]string{"contribute", "prompts/new-prompt.md", "--create-pr"})

		_ = contributeCmd.Execute()

		// Then: Pull request instructions are provided
		output := contributeBuf.String()
		assert.Contains(t, output, "pull request", "Should mention PR")
		assert.Contains(t, output, "branch", "Should create branch")
	})
}

// TestAcceptance_US009_PullUpdatesFromUpstream tests US-009: Pull Updates from Upstream
func TestAcceptance_US009_PullUpdatesFromUpstream(t *testing.T) {
	t.Run("sync_with_upstream", func(t *testing.T) {
		// Given: Upstream has new commits
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Pulling updates
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update", "--sync"})

		err := updateCmd.Execute()

		// Then: Local repository is synchronized
		assert.NoError(t, err)
		output := updateBuf.String()
		assert.Contains(t, output, "Synchronizing with upstream", "Should sync")
		assert.Contains(t, output, "commits behind", "Should show sync status")
	})

	t.Run("handle_diverged_branches", func(t *testing.T) {
		// Given: Local and upstream have diverged
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// Simulate divergence
		// This would normally require git operations

		// When: Attempting to sync
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update"})

		_ = updateCmd.Execute()

		// Then: Divergence is detected and handled
		output := updateBuf.String()
		assert.Contains(t, strings.ToLower(output), "diverg", "Should detect divergence")
	})
}

// TestAcceptance_US010_HandleUpdateConflicts tests US-010: Handle Update Conflicts
func TestAcceptance_US010_HandleUpdateConflicts(t *testing.T) {
	t.Run("detect_conflicts", func(t *testing.T) {
		// Given: Conflicting changes exist
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// Create conflicting file
		conflictFile := filepath.Join(".ddx", "conflict.txt")
		os.WriteFile(conflictFile, []byte("local version"), 0644)

		// When: Updating
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update"})

		_ = updateCmd.Execute()

		// Then: Conflicts are detected
		output := updateBuf.String()
		assert.Contains(t, strings.ToLower(output), "conflict", "Should detect conflicts")
	})

	t.Run("interactive_resolution", func(t *testing.T) {
		// Given: Conflicts need resolution
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Using interactive resolution
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update", "--interactive"})

		// Would need to mock user input here
		updateErr := updateCmd.Execute()

		// Then: Interactive options are provided
		output := updateBuf.String()
		if updateErr == nil || strings.Contains(output, "interactive") {
			assert.Contains(t, output, "Choose", "Should provide choices")
		}
	})

	t.Run("automatic_resolution_strategy", func(t *testing.T) {
		// Given: User wants automatic resolution
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Using automatic strategy
		updateCmd := rootCmd
		updateBuf := new(bytes.Buffer)
		updateCmd.SetOut(updateBuf)
		updateCmd.SetErr(updateBuf)
		updateCmd.SetArgs([]string{"update", "--strategy=ours"})

		err := updateCmd.Execute()

		// Then: Conflicts are resolved automatically
		assert.NoError(t, err)
		output := updateBuf.String()
		assert.Contains(t, output, "strategy", "Should apply strategy")
	})
}

// TestAcceptance_US011_ContributeChangesUpstream tests US-011: Contribute Changes Upstream
func TestAcceptance_US011_ContributeChangesUpstream(t *testing.T) {
	t.Run("prepare_contribution_branch", func(t *testing.T) {
		// Given: User has changes to contribute
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// Create changes
		newFile := filepath.Join(".ddx", "patterns", "new-pattern.md")
		os.MkdirAll(filepath.Dir(newFile), 0755)
		os.WriteFile(newFile, []byte("# New Pattern"), 0644)

		// When: Preparing contribution
		contributeCmd := rootCmd
		contributeBuf := new(bytes.Buffer)
		contributeCmd.SetOut(contributeBuf)
		contributeCmd.SetErr(contributeBuf)
		contributeCmd.SetArgs([]string{"contribute", "patterns/new-pattern.md"})

		err := contributeCmd.Execute()

		// Then: Feature branch is created
		assert.NoError(t, err)
		output := contributeBuf.String()
		assert.Contains(t, output, "branch", "Should create branch")
		assert.Contains(t, output, "feature", "Should be feature branch")
	})

	t.Run("validate_contribution_standards", func(t *testing.T) {
		// Given: Contribution needs validation
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Contributing
		contributeCmd := rootCmd
		contributeBuf := new(bytes.Buffer)
		contributeCmd.SetOut(contributeBuf)
		contributeCmd.SetErr(contributeBuf)
		contributeCmd.SetArgs([]string{"contribute", "templates/test"})

		_ = contributeCmd.Execute()

		// Then: Standards are validated
		output := contributeBuf.String()
		assert.Contains(t, output, "Checking contribution standards", "Should validate standards")
	})

	t.Run("push_to_fork", func(t *testing.T) {
		// Given: Contribution is ready
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupTestProject(t)

		// When: Pushing to fork
		contributeCmd := rootCmd
		contributeBuf := new(bytes.Buffer)
		contributeCmd.SetOut(contributeBuf)
		contributeCmd.SetErr(contributeBuf)
		contributeCmd.SetArgs([]string{"contribute", "prompts/test.md", "--push"})

		_ = contributeCmd.Execute()

		// Then: Changes are pushed to fork
		output := contributeBuf.String()
		assert.Contains(t, output, "push", "Should push changes")
		assert.Contains(t, output, "fork", "Should mention fork")
	})
}

// Helper function to setup a test project with DDx
func setupTestProject(t *testing.T) {
	// Create .ddx.yml configuration
	config := `
name: test-project
repository:
  url: https://github.com/ddx-tools/ddx
  branch: main
`
	err := os.WriteFile(".ddx.yml", []byte(config), 0644)
	require.NoError(t, err, "Should create config file")

	// Create .ddx directory structure
	os.MkdirAll(".ddx/templates", 0755)
	os.MkdirAll(".ddx/patterns", 0755)
	os.MkdirAll(".ddx/prompts", 0755)
}
