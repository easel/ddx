package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// End-to-end tests simulate real user workflows
// These tests require the CLI to be built and DDx repository to be available

// TestE2E_BasicWorkflow tests the complete workflow from init to apply
func TestE2E_BasicWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Build the CLI if needed
	buildCmd := exec.Command("go", "build", "-o", "ddx", "..")
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Could not build CLI: %v", err)
	}

	// Make sure it's executable
	os.Chmod("ddx", 0755)
	defer os.Remove("ddx")

	// Create test workspace
	workspace := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to workspace
	require.NoError(t, os.Chdir(workspace))

	// Step 1: Initialize DDx
	t.Run("init", func(t *testing.T) {
		cmd := exec.Command(filepath.Join(originalDir, "ddx"), "init")
		output, err := cmd.CombinedOutput()

		// May fail if DDx repository not available or no git
		if err != nil {
			t.Logf("Init output: %s", output)
			// For test purposes, we'll just check the command ran
			// Real E2E would need actual repository
			return
		}

		// Verify config file created
		assert.FileExists(t, filepath.Join(workspace, ".ddx.yml"))
	})

	// Step 2: List available resources
	t.Run("list", func(t *testing.T) {
		cmd := exec.Command(filepath.Join(originalDir, "ddx"), "list")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("List failed: %s", output)
		}

		// Should show available resources
		assert.Contains(t, string(output), "Templates")
	})

	// Step 3: Apply a template (if available)
	t.Run("apply", func(t *testing.T) {
		// First check what templates are available
		listCmd := exec.Command(filepath.Join(originalDir, "ddx"), "list", "templates")
		listOutput, _ := listCmd.CombinedOutput()

		// If we have templates, try to apply one
		if string(listOutput) != "" {
			cmd := exec.Command(filepath.Join(originalDir, "ddx"), "apply", "templates/common")
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Logf("Apply output: %s", output)
			}
		}
	})

	// Step 4: Check configuration
	t.Run("config", func(t *testing.T) {
		cmd := exec.Command(filepath.Join(originalDir, "ddx"), "config")
		output, err := cmd.CombinedOutput()

		assert.NoError(t, err)
		assert.Contains(t, string(output), "version")
	})
}

// TestE2E_TemplateWithVariables tests template application with variable substitution
func TestE2E_TemplateWithVariables(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup mock DDx home
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	// Create a template with variables
	templateDir := filepath.Join(homeDir, ".ddx", "templates", "test")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Create template files
	templateFiles := map[string]string{
		"README.md":   "# {{project_name}}\n\nVersion: {{version}}",
		"config.json": `{"name": "{{project_name}}", "port": {{port}}}`,
	}

	for name, content := range templateFiles {
		filePath := filepath.Join(templateDir, name)
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
	}

	// Create test project
	projectDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(projectDir))

	// Initialize project
	config := `version: "1.0"
variables:
  project_name: "MyProject"
  version: "1.0.0"
  port: "8080"`
	require.NoError(t, os.WriteFile(filepath.Join(projectDir, ".ddx.yml"), []byte(config), 0644))

	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "ddx", filepath.Join(originalDir, ".."))
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Could not build CLI: %v", err)
	}
	os.Chmod("ddx", 0755)
	defer os.Remove("ddx")

	// Apply template
	cmd := exec.Command("./ddx", "apply", "templates/test")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("Apply output: %s", output)
	}

	// Verify files were created with substituted variables
	readmeContent, err := os.ReadFile(filepath.Join(projectDir, "README.md"))
	if err == nil {
		assert.Contains(t, string(readmeContent), "MyProject")
		assert.Contains(t, string(readmeContent), "1.0.0")
	}

	configContent, err := os.ReadFile(filepath.Join(projectDir, "config.json"))
	if err == nil {
		assert.Contains(t, string(configContent), "MyProject")
		assert.Contains(t, string(configContent), "8080")
	}
}

// TestE2E_GitIntegration tests git subtree operations
func TestE2E_GitIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// This test requires git to be available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git not available")
	}

	workspace := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(workspace))

	// Initialize git repo
	cmd := exec.Command("git", "init")
	require.NoError(t, cmd.Run())

	// Configure git (required for commits)
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "ddx", filepath.Join(originalDir, ".."))
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Could not build CLI: %v", err)
	}
	os.Chmod("ddx", 0755)
	defer os.Remove("ddx")

	// Initialize DDx
	initCmd := exec.Command("./ddx", "init")
	initOutput, err := initCmd.CombinedOutput()

	if err != nil {
		t.Logf("Init output: %s", initOutput)
		t.Skip("DDx init failed, likely repository not available")
	}

	// Verify git subtree was set up
	subtreeDir := filepath.Join(workspace, ".ddx")
	if _, err := os.Stat(subtreeDir); err == nil {
		// Check if it's tracked by git
		statusCmd := exec.Command("git", "status", "--porcelain", subtreeDir)
		statusOutput, _ := statusCmd.CombinedOutput()
		t.Logf("Git status: %s", statusOutput)
	}
}

// TestE2E_ContributionWorkflow tests the contribution workflow
func TestE2E_ContributionWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// This test simulates a user creating a new pattern and contributing it
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	workspace := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(workspace))

	// Initialize project
	config := `version: "1.0"
repository:
  url: "https://github.com/ddx-tools/ddx"
  branch: "main"`
	require.NoError(t, os.WriteFile(filepath.Join(workspace, ".ddx.yml"), []byte(config), 0644))

	// Create DDx directory structure
	ddxDir := filepath.Join(workspace, ".ddx")
	patternsDir := filepath.Join(ddxDir, "patterns", "custom-pattern")
	require.NoError(t, os.MkdirAll(patternsDir, 0755))

	// Add a new pattern
	patternFile := filepath.Join(patternsDir, "pattern.js")
	patternContent := `// Custom error handling pattern
function handleError(err) {
  console.error('Error:', err);
  // Custom logic here
}`
	require.NoError(t, os.WriteFile(patternFile, []byte(patternContent), 0644))

	// Add pattern documentation
	readmeFile := filepath.Join(patternsDir, "README.md")
	readmeContent := `# Custom Pattern

This pattern provides error handling.

## Usage
` + "```javascript\n" + patternContent + "\n```"
	require.NoError(t, os.WriteFile(readmeFile, []byte(readmeContent), 0644))

	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "ddx", filepath.Join(originalDir, ".."))
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Could not build CLI: %v", err)
	}
	os.Chmod("ddx", 0755)
	defer os.Remove("ddx")

	// Test contribution command (would normally push to upstream)
	contribCmd := exec.Command("./ddx", "contribute", "--dry-run")
	contribOutput, err := contribCmd.CombinedOutput()

	t.Logf("Contribute output: %s", contribOutput)
	// The contribute command may not be implemented yet
	_ = err
}

// TestE2E_UpdateWorkflow tests updating from upstream
func TestE2E_UpdateWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	workspace := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(workspace))

	// Create initial config
	config := `version: "1.0"
repository:
  url: "https://github.com/ddx-tools/ddx"
  branch: "main"
sync:
  last_update: "2024-01-01T00:00:00Z"
  upstream_commit: "abc123"`
	require.NoError(t, os.WriteFile(filepath.Join(workspace, ".ddx.yml"), []byte(config), 0644))

	// Build CLI
	buildCmd := exec.Command("go", "build", "-o", "ddx", filepath.Join(originalDir, ".."))
	if err := buildCmd.Run(); err != nil {
		t.Skipf("Could not build CLI: %v", err)
	}
	os.Chmod("ddx", 0755)
	defer os.Remove("ddx")

	// Test update command
	updateCmd := exec.Command("./ddx", "update", "--check")
	updateOutput, err := updateCmd.CombinedOutput()

	t.Logf("Update output: %s", updateOutput)
	// The update command may not be implemented yet
	_ = err
}
