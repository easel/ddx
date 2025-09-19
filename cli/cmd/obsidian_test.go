package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObsidianCommand_Help(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CLI test in short mode")
	}

	// Test that obsidian command exists and shows help
	cmd := exec.Command("ddx", "obsidian", "--help")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	outputStr := string(output)
	assert.Contains(t, outputStr, "obsidian")
	assert.Contains(t, outputStr, "migrate")
	assert.Contains(t, outputStr, "validate")
	assert.Contains(t, outputStr, "revert")
}

func TestObsidianMigrate_DryRun(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CLI test in short mode")
	}

	// Setup test directory
	testDir := t.TempDir()
	setupTestHelixWorkflow(t, testDir)

	// Run dry-run migration
	cmd := exec.Command("ddx", "obsidian", "migrate", "--dry-run", "--path", testDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	outputStr := string(output)
	// Verify output shows what would be done
	assert.Contains(t, outputStr, "Would add frontmatter to")
	assert.Contains(t, outputStr, "Would convert links in")
	assert.Contains(t, outputStr, "Would generate navigation hub")

	// Verify no files were actually modified
	assertNoFilesModified(t, testDir)
}

func TestObsidianMigrate_FullMigration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CLI test in short mode")
	}

	testDir := t.TempDir()
	setupTestHelixWorkflow(t, testDir)

	// Capture original state
	originalFiles := captureFileStates(t, testDir)

	// Run full migration
	cmd := exec.Command("ddx", "obsidian", "migrate", "--path", testDir)
	err := cmd.Run()
	require.NoError(t, err)

	// Verify all HELIX markdown files have frontmatter
	err = filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".md") && isHelixFile(path) {
			content, err := os.ReadFile(path)
			require.NoError(t, err)

			// Must start with frontmatter
			assert.True(t, strings.HasPrefix(string(content), "---\n"), "File missing frontmatter: %s", path)

			// Must contain required fields
			contentStr := string(content)
			assert.Contains(t, contentStr, "title:", "Missing title in %s", path)
			assert.Contains(t, contentStr, "type:", "Missing type in %s", path)
			assert.Contains(t, contentStr, "tags:", "Missing tags in %s", path)
			assert.Contains(t, contentStr, "created:", "Missing created in %s", path)
			assert.Contains(t, contentStr, "updated:", "Missing updated in %s", path)
		}
		return nil
	})
	require.NoError(t, err)

	// Verify navigation hub was created
	navPath := filepath.Join(testDir, "workflows/helix/HELIX-NAVIGATOR.md")
	assert.FileExists(t, navPath)

	// Verify some files were actually modified
	currentFiles := captureFileStates(t, testDir)
	assert.NotEqual(t, originalFiles, currentFiles, "Migration should have modified files")
}

func TestObsidianValidate_ValidFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CLI test in short mode")
	}

	testDir := t.TempDir()
	setupMigratedHelixWorkflow(t, testDir)

	// Run validation
	cmd := exec.Command("ddx", "obsidian", "validate", "--path", testDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	outputStr := string(output)
	// Should report no errors for valid files
	assert.Contains(t, outputStr, "✅")    // Success indicator
	assert.NotContains(t, outputStr, "❌") // No failure indicators
}

func TestObsidianValidate_InvalidFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CLI test in short mode")
	}

	testDir := t.TempDir()
	setupInvalidHelixFiles(t, testDir)

	cmd := exec.Command("ddx", "obsidian", "validate", "--path", testDir)
	output, err := cmd.CombinedOutput()

	// Should exit with error code
	assert.Error(t, err)

	outputStr := string(output)
	// Should report specific validation errors
	assert.Contains(t, outputStr, "❌") // Failure indicators
	// At least one of these error types should be present
	hasError := strings.Contains(outputStr, "Missing required field") ||
		strings.Contains(outputStr, "Invalid tag format") ||
		strings.Contains(outputStr, "Broken wikilink") ||
		strings.Contains(outputStr, "validation")
	assert.True(t, hasError, "Should report validation errors")
}

func TestObsidianRevert(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping CLI test in short mode")
	}

	testDir := t.TempDir()
	setupTestHelixWorkflow(t, testDir)

	// Capture original state
	originalFiles := captureFileStates(t, testDir)

	// Migrate first
	migrateCmd := exec.Command("ddx", "obsidian", "migrate", "--path", testDir)
	err := migrateCmd.Run()
	require.NoError(t, err)

	// Verify migration changed files
	migratedFiles := captureFileStates(t, testDir)
	assert.NotEqual(t, originalFiles, migratedFiles)

	// Revert migration
	revertCmd := exec.Command("ddx", "obsidian", "revert", "--path", testDir)
	err = revertCmd.Run()
	require.NoError(t, err)

	// Verify files restored to original state
	currentFiles := captureFileStates(t, testDir)
	assert.Equal(t, originalFiles, currentFiles)
}

// Helper functions for test setup

func setupTestHelixWorkflow(t *testing.T, dir string) {
	// Create basic HELIX structure for testing
	workflowDir := filepath.Join(dir, "workflows/helix")
	err := os.MkdirAll(workflowDir, 0755)
	require.NoError(t, err)

	phasesDir := filepath.Join(workflowDir, "phases/01-frame")
	err = os.MkdirAll(phasesDir, 0755)
	require.NoError(t, err)

	docsDir := filepath.Join(dir, "docs/01-frame/features")
	err = os.MkdirAll(docsDir, 0755)
	require.NoError(t, err)

	// Create sample coordinator file
	coordinatorPath := filepath.Join(workflowDir, "coordinator.md")
	content := "# HELIX Coordinator\n\nSee [Frame Phase](./phases/01-frame/README.md) to start."
	err = os.WriteFile(coordinatorPath, []byte(content), 0644)
	require.NoError(t, err)

	// Create sample phase file
	phasePath := filepath.Join(phasesDir, "README.md")
	phaseContent := "# Frame Phase\n\nDefine what you're building.\n\nNext: [Design Phase](../02-design/README.md)"
	err = os.WriteFile(phasePath, []byte(phaseContent), 0644)
	require.NoError(t, err)

	// Create sample feature file
	featurePath := filepath.Join(docsDir, "FEAT-001-test-feature.md")
	featureContent := "# Feature Specification: FEAT-001\n\n**Priority**: P1\n\nTest feature for migration."
	err = os.WriteFile(featurePath, []byte(featureContent), 0644)
	require.NoError(t, err)
}

func setupMigratedHelixWorkflow(t *testing.T, dir string) {
	setupTestHelixWorkflow(t, dir)

	// Add frontmatter to files to simulate successful migration
	files := []string{
		"workflows/helix/coordinator.md",
		"workflows/helix/phases/01-frame/README.md",
		"docs/01-frame/features/FEAT-001-test-feature.md",
	}

	for _, file := range files {
		path := filepath.Join(dir, file)
		content, err := os.ReadFile(path)
		require.NoError(t, err)

		// Add basic frontmatter
		frontmatter := `---
title: "Test File"
type: "test"
tags:
  - helix
  - test
created: 2025-01-18
updated: 2025-01-18
---

`
		newContent := frontmatter + string(content)
		err = os.WriteFile(path, []byte(newContent), 0644)
		require.NoError(t, err)
	}
}

func setupInvalidHelixFiles(t *testing.T, dir string) {
	setupTestHelixWorkflow(t, dir)

	// Create files with invalid frontmatter
	invalidPath := filepath.Join(dir, "workflows/helix/invalid.md")
	invalidContent := `---
title: "Missing closing"
type: coordinator
# Invalid YAML - missing closing ---

# Invalid File

This file has broken YAML frontmatter.
`
	err := os.WriteFile(invalidPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	// Create file with broken wikilink
	brokenLinkPath := filepath.Join(dir, "workflows/helix/broken-links.md")
	brokenContent := `---
title: "Broken Links"
type: coordinator
tags:
  - helix
created: 2025-01-18
updated: 2025-01-18
---

# Broken Links

See [[Nonexistent File]] for details.
`
	err = os.WriteFile(brokenLinkPath, []byte(brokenContent), 0644)
	require.NoError(t, err)
}

func assertNoFilesModified(t *testing.T, dir string) {
	// For dry-run, verify that original test files still have their original content
	coordinatorPath := filepath.Join(dir, "workflows/helix/coordinator.md")
	content, err := os.ReadFile(coordinatorPath)
	require.NoError(t, err)

	// Should not start with frontmatter
	assert.False(t, strings.HasPrefix(string(content), "---\n"), "File should not have been modified in dry-run")
}

func captureFileStates(t *testing.T, dir string) map[string]string {
	states := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".md") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Store relative path from dir
			relPath, _ := filepath.Rel(dir, path)
			states[relPath] = string(content)
		}
		return nil
	})
	require.NoError(t, err)

	return states
}

func isHelixFile(path string) bool {
	// Simple check for HELIX files in test context
	return strings.Contains(path, "/helix/") ||
		strings.Contains(path, "/docs/") ||
		strings.Contains(path, "workflows")
}
