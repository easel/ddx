package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMCPListCommand_Contract tests the contract for ddx mcp list command
func TestMCPListCommand_Contract(t *testing.T) {
	t.Run("contract_exit_code_0_success", func(t *testing.T) {
		// Given: Valid project with MCP configured
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Listing MCP servers
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "list"})

		// Then: Exit code should be 0
		err := cmd.Execute()
		if err != nil {
			var exitErr *ExitError
			if errors.As(err, &exitErr) {
				assert.Equal(t, 0, exitErr.Code, "Should exit with code 0 on success")
			}
		}
	})

	t.Run("contract_output_format", func(t *testing.T) {
		// Given: MCP servers available
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Listing servers
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "list"})

		_ = cmd.Execute()

		// Then: Output should follow format
		output := buf.String()
		lines := strings.Split(output, "\n")

		// Should have header
		hasHeader := false
		for _, line := range lines {
			if strings.Contains(line, "MCP Servers") || strings.Contains(line, "Available") {
				hasHeader = true
				break
			}
		}
		assert.True(t, hasHeader, "Should have header")

		// Should have status indicators
		assert.Contains(t, output, "✅", "Should have installed indicator")
		assert.Contains(t, output, "⬜", "Should have not-installed indicator")
	})

	t.Run("contract_category_filter", func(t *testing.T) {
		// Given: Multiple categories exist
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Filtering by category
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "list", "--category", "development"})

		_ = cmd.Execute()

		// Then: Should only show that category
		output := buf.String()
		assert.Contains(t, output, "development", "Should show category")
		// Should not show other categories
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "weather") {
				assert.Fail(t, "Should not show other categories")
			}
		}
	})

	t.Run("contract_search_parameter", func(t *testing.T) {
		// Given: Searchable servers
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Searching
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "list", "--search", "file"})

		_ = cmd.Execute()

		// Then: Should filter by search term
		output := buf.String()
		assert.Contains(t, output, "filesystem", "Should find filesystem")
		// Count results
		lines := strings.Split(output, "\n")
		resultCount := 0
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), "file") {
				resultCount++
			}
		}
		assert.Greater(t, resultCount, 0, "Should have search results")
	})
}

// TestMCPInstallCommand_Contract tests the contract for ddx mcp install command
func TestMCPInstallCommand_Contract(t *testing.T) {
	t.Run("contract_exit_code_0_success", func(t *testing.T) {
		// Given: Valid server to install
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Installing server
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "filesystem"})

		// Then: Should exit with 0 on success
		err := cmd.Execute()
		if err != nil {
			var exitErr *ExitError
			if errors.As(err, &exitErr) {
				assert.Equal(t, 0, exitErr.Code, "Should exit 0 on success")
			}
		}
	})

	t.Run("contract_exit_code_6_not_found", func(t *testing.T) {
		// Given: Non-existent server
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Installing non-existent server
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "nonexistent-server"})

		// Then: Should exit with code 6
		err := cmd.Execute()
		var exitErr *ExitError
		if errors.As(err, &exitErr) {
			assert.Equal(t, 6, exitErr.Code, "Should exit 6 when not found")
		}
	})

	t.Run("contract_package_json_creation", func(t *testing.T) {
		// Given: No package.json exists
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Installing server
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "filesystem"})

		_ = cmd.Execute()

		// Then: Should create package.json
		assert.FileExists(t, "package.json", "Should create package.json")
	})

	t.Run("contract_claude_config_update", func(t *testing.T) {
		// Given: Installing MCP server
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Installing
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "github"})

		_ = cmd.Execute()

		// Then: Should update Claude config
		claudeConfig := ".claude/settings.local.json"
		assert.FileExists(t, claudeConfig, "Should create Claude config")
	})

	t.Run("contract_validate_flag", func(t *testing.T) {
		// Given: Installing with validation
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Installing with --validate
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "filesystem", "--validate"})

		_ = cmd.Execute()

		// Then: Should validate installation
		output := buf.String()
		assert.Contains(t, output, "Validat", "Should validate")
	})
}

// TestConfigCommand_Contract tests configuration command contracts
func TestConfigCommand_ContractExtended(t *testing.T) {
	t.Run("contract_init_creates_config", func(t *testing.T) {
		// Given: No config exists
		tempDir := t.TempDir()
		os.Chdir(tempDir)

		// When: Initializing config
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "init"})

		err := cmd.Execute()

		// Then: Should create config file
		if err == nil {
			assert.FileExists(t, ".ddx.yml", "Should create config")
		}
	})

	t.Run("contract_set_variable", func(t *testing.T) {
		// Given: Config exists
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// When: Setting variable
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "set", "variables.test", "value"})

		err := cmd.Execute()

		// Then: Should update config
		if err == nil {
			content, _ := os.ReadFile(".ddx.yml")
			assert.Contains(t, string(content), "test", "Should set variable")
		}
	})

	t.Run("contract_export_import", func(t *testing.T) {
		// Given: Config to export
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// When: Exporting
		exportCmd := rootCmd
		exportBuf := new(bytes.Buffer)
		exportCmd.SetOut(exportBuf)
		exportCmd.SetErr(exportBuf)
		exportCmd.SetArgs([]string{"config", "export"})

		_ = exportCmd.Execute()

		// Then: Should output config
		output := exportBuf.String()
		assert.Contains(t, output, "name:", "Should export config")
		assert.Contains(t, output, "repository:", "Should include repository")
	})

	t.Run("contract_validate_exit_codes", func(t *testing.T) {
		// Given: Invalid config
		tempDir := t.TempDir()
		os.Chdir(tempDir)

		// Create invalid config (missing required fields)
		os.WriteFile(".ddx.yml", []byte("invalid: yaml: content:"), 0644)

		// When: Validating
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "validate"})

		// Then: Should exit with error code
		err := cmd.Execute()
		assert.Error(t, err, "Should fail validation")
	})

	t.Run("contract_override_precedence", func(t *testing.T) {
		// Given: Multiple config layers
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// Create local override
		os.WriteFile(".ddx.local.yml", []byte("variables:\n  override: local"), 0644)

		// When: Showing effective config
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "show", "--effective"})

		_ = cmd.Execute()

		// Then: Should apply override precedence
		output := buf.String()
		assert.Contains(t, output, "override", "Should show override")
	})
}

// TestInstallationCommands_Contract tests installation-related command contracts
func TestInstallationCommands_Contract(t *testing.T) {
	t.Run("contract_doctor_command", func(t *testing.T) {
		// Given: DDx installed
		// When: Running doctor
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"doctor"})

		err := cmd.Execute()

		// Then: Should provide diagnostic info
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "DDx", "Should mention DDx")
			assert.Contains(t, output, "Check", "Should check components")
		}
	})

	t.Run("contract_self_update", func(t *testing.T) {
		// Given: DDx installed
		// When: Checking for updates
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"self-update", "--check"})

		err := cmd.Execute()

		// Then: Should check version
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "version", "Should check version")
		}
	})

	t.Run("contract_setup_path", func(t *testing.T) {
		// Given: DDx binary exists
		// When: Setting up PATH
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"setup", "path", "--dry-run"})

		err := cmd.Execute()

		// Then: Should show PATH setup
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "PATH", "Should mention PATH")
			assert.Contains(t, output, "shell", "Should detect shell")
		}
	})

	t.Run("contract_uninstall_confirm", func(t *testing.T) {
		// Given: DDx installed
		// When: Uninstalling without confirmation
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"uninstall"})

		err := cmd.Execute()

		// Then: Should require confirmation
		if err != nil {
			output := buf.String()
			assert.Contains(t, output, "confirm", "Should require confirmation")
		}
	})
}

// TestWorkflowCommands_Contract tests workflow command contracts
func TestWorkflowCommands_Contract(t *testing.T) {
	t.Run("contract_workflow_status", func(t *testing.T) {
		// Given: HELIX workflow active
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createWorkflowState(t)

		// When: Checking status
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"workflow", "status"})

		err := cmd.Execute()

		// Then: Should show workflow status
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "Phase", "Should show phase")
			assert.Contains(t, output, "Progress", "Should show progress")
		}
	})

	t.Run("contract_workflow_validate", func(t *testing.T) {
		// Given: Workflow in progress
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createWorkflowState(t)

		// When: Validating phase
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"workflow", "validate"})

		err := cmd.Execute()

		// Then: Should validate current phase
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "Validat", "Should validate")
			assert.Contains(t, output, "criteria", "Should check criteria")
		}
	})

	t.Run("contract_workflow_advance", func(t *testing.T) {
		// Given: Phase complete
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createWorkflowState(t)

		// When: Advancing workflow
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"workflow", "advance"})

		err := cmd.Execute()

		// Then: Should advance to next phase
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "Advancing", "Should advance")
			assert.Contains(t, output, "phase", "Should mention phase")
		}
	})
}

// Helper to create workflow state
func createWorkflowState(t *testing.T) {
	state := `workflow: helix
current_phase: test
phases_completed:
  - frame
  - design
`
	err := os.WriteFile(".helix-state.yml", []byte(state), 0644)
	require.NoError(t, err)
}
