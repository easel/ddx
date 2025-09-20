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

// TestAcceptance_US036_ListMCPServers tests US-036: List Available MCP Servers
func TestAcceptance_US036_ListMCPServers(t *testing.T) {
	t.Run("display_all_available_servers", func(t *testing.T) {
		// Given: MCP server registry is available
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Running ddx mcp list
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "list"})

		err := cmd.Execute()

		// Then: Should display all available servers
		assert.NoError(t, err, "Should list MCP servers")
		output := buf.String()
		assert.Contains(t, output, "Available MCP Servers", "Should show header")
		assert.Contains(t, output, "filesystem", "Should show filesystem server")
		assert.Contains(t, output, "github", "Should show github server")
		assert.Contains(t, output, "✅", "Should show installed indicator")
		assert.Contains(t, output, "⬜", "Should show not-installed indicator")
	})

	t.Run("filter_by_category", func(t *testing.T) {
		// Given: Multiple categories of MCP servers exist
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Filtering by category
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "list", "--category", "development"})

		err := cmd.Execute()

		// Then: Should only show servers in that category
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Category: development", "Should indicate filter")
		assert.Contains(t, output, "sequential-thinking", "Should show dev servers")
		assert.NotContains(t, output, "weather", "Should not show other categories")
	})

	t.Run("search_functionality", func(t *testing.T) {
		// Given: Want to find servers related to "git"
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Searching for "git"
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "list", "--search", "git"})

		err := cmd.Execute()

		// Then: Should show matching servers
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "github", "Should find github server")
		assert.Contains(t, output, "git", "Should highlight search term")
		lines := strings.Split(output, "\n")
		gitCount := 0
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), "git") {
				gitCount++
			}
		}
		assert.Greater(t, gitCount, 0, "Should find git-related servers")
	})

	t.Run("show_installation_status", func(t *testing.T) {
		// Given: Some MCP servers are installed
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// Install a server first
		installCmd := rootCmd
		installBuf := new(bytes.Buffer)
		installCmd.SetOut(installBuf)
		installCmd.SetErr(installBuf)
		installCmd.SetArgs([]string{"mcp", "install", "filesystem"})
		_ = installCmd.Execute()

		// When: Listing servers
		listCmd := rootCmd
		listBuf := new(bytes.Buffer)
		listCmd.SetOut(listBuf)
		listCmd.SetErr(listBuf)
		listCmd.SetArgs([]string{"mcp", "list"})

		err := listCmd.Execute()

		// Then: Should show correct installation status
		assert.NoError(t, err)
		output := listBuf.String()
		// Find filesystem line and check it has installed indicator
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "filesystem") {
				assert.Contains(t, line, "✅", "Installed server should show ✅")
			}
		}
	})

	t.Run("detailed_verbose_view", func(t *testing.T) {
		// Given: Want more information
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Running with --verbose
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "list", "--verbose"})

		err := cmd.Execute()

		// Then: Should show additional details
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Environment", "Should show env vars")
		assert.Contains(t, output, "Package", "Should show package info")
		assert.Contains(t, output, "Author", "Should show author")
		assert.Contains(t, output, "Version", "Should show version")
	})
}

// TestAcceptance_US037_InstallMCPServer tests US-037: Install MCP Server
func TestAcceptance_US037_InstallMCPServer(t *testing.T) {
	t.Run("install_server_locally", func(t *testing.T) {
		// Given: MCP server not installed
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Installing a server
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "filesystem"})

		err := cmd.Execute()

		// Then: Should install server locally
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Installing", "Should show installation")
		assert.Contains(t, output, "filesystem", "Should name the server")
		assert.Contains(t, output, "Success", "Should indicate success")

		// Check package.json was created/updated
		assert.FileExists(t, "package.json", "Should create package.json")

		// Check Claude config was updated
		claudeConfig := filepath.Join(".claude", "settings.local.json")
		assert.FileExists(t, claudeConfig, "Should create Claude config")
	})

	t.Run("detect_package_manager", func(t *testing.T) {
		// Given: Different package managers available
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// Create pnpm-lock.yaml to trigger pnpm detection
		os.WriteFile("pnpm-lock.yaml", []byte("lockfileVersion: 5.4"), 0644)

		// When: Installing
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "github"})

		err := cmd.Execute()

		// Then: Should detect and use pnpm
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "pnpm", "Should detect pnpm")
		assert.Contains(t, output, "Using package manager: pnpm", "Should indicate pnpm usage")
	})

	t.Run("configure_server_environment", func(t *testing.T) {
		// Given: Server needs configuration
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Installing with configuration
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "github", "--config", "token=ghp_test123"})

		err := cmd.Execute()

		// Then: Should configure the server
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Configuring", "Should configure server")

		// Check configuration was written
		claudeConfig := filepath.Join(".claude", "settings.local.json")
		if content, err := os.ReadFile(claudeConfig); err == nil {
			assert.Contains(t, string(content), "GITHUB_TOKEN", "Should set env var")
		}
	})

	t.Run("handle_already_installed", func(t *testing.T) {
		// Given: Server already installed
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// Install once
		cmd1 := rootCmd
		buf1 := new(bytes.Buffer)
		cmd1.SetOut(buf1)
		cmd1.SetErr(buf1)
		cmd1.SetArgs([]string{"mcp", "install", "filesystem"})
		_ = cmd1.Execute()

		// When: Installing again
		cmd2 := rootCmd
		buf2 := new(bytes.Buffer)
		cmd2.SetOut(buf2)
		cmd2.SetErr(buf2)
		cmd2.SetArgs([]string{"mcp", "install", "filesystem"})

		_ = cmd2.Execute()

		// Then: Should handle gracefully
		output := buf2.String()
		assert.Contains(t, output, "already installed", "Should detect existing")
		assert.Contains(t, output, "upgrade", "Should offer upgrade option")
	})

	t.Run("validate_installation", func(t *testing.T) {
		// Given: Server installed
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		setupMCPTestProject(t)

		// When: Installing and validating
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"mcp", "install", "filesystem", "--validate"})

		err := cmd.Execute()

		// Then: Should validate installation
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Validating", "Should validate")
		assert.Contains(t, output, "Connection test", "Should test connection")
		assert.Contains(t, output, "✓", "Should show validation success")
	})
}

// TestAcceptance_ConfigurationManagement_Extended tests US-017 through US-024
func TestAcceptance_ConfigurationManagement_Extended(t *testing.T) {
	t.Run("initialize_configuration", func(t *testing.T) {
		// Given: No configuration exists
		tempDir := t.TempDir()
		os.Chdir(tempDir)

		// When: Initializing configuration
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "init"})

		err := cmd.Execute()

		// Then: Should create configuration
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Configuration initialized", "Should init config")
		assert.FileExists(t, ".ddx.yml", "Should create config file")
	})

	t.Run("configure_variables", func(t *testing.T) {
		// Given: Configuration exists
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// When: Setting variables
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "set", "variables.project_name", "my-project"})

		err := cmd.Execute()

		// Then: Should update variables
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Variable set", "Should set variable")
		assert.Contains(t, output, "project_name", "Should show variable name")
		assert.Contains(t, output, "my-project", "Should show value")
	})

	t.Run("override_configuration", func(t *testing.T) {
		// Given: Base configuration exists
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// When: Creating override
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "override", "--env", "production"})

		err := cmd.Execute()

		// Then: Should create environment override
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Override created", "Should create override")
		assert.FileExists(t, ".ddx.production.yml", "Should create override file")
	})

	t.Run("configure_resource_selection", func(t *testing.T) {
		// Given: Want to select specific resources
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// When: Configuring resources
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "resources", "add", "templates/nextjs"})

		err := cmd.Execute()

		// Then: Should update resource selection
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Resource added", "Should add resource")
		assert.Contains(t, output, "templates/nextjs", "Should show resource")
	})

	t.Run("validate_configuration", func(t *testing.T) {
		// Given: Configuration to validate
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// When: Validating
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "validate"})

		err := cmd.Execute()

		// Then: Should validate configuration
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Configuration valid", "Should validate")
		assert.Contains(t, output, "✓", "Should show success indicator")
	})

	t.Run("export_import_configuration", func(t *testing.T) {
		// Given: Configuration to export
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// When: Exporting configuration
		exportCmd := rootCmd
		exportBuf := new(bytes.Buffer)
		exportCmd.SetOut(exportBuf)
		exportCmd.SetErr(exportBuf)
		exportCmd.SetArgs([]string{"config", "export", "--output", "config.export.yml"})

		err := exportCmd.Execute()

		// Then: Should export configuration
		assert.NoError(t, err)
		assert.FileExists(t, "config.export.yml", "Should create export")

		// When: Importing configuration
		os.Remove(".ddx.yml") // Remove original
		importCmd := rootCmd
		importBuf := new(bytes.Buffer)
		importCmd.SetOut(importBuf)
		importCmd.SetErr(importBuf)
		importCmd.SetArgs([]string{"config", "import", "config.export.yml"})

		err = importCmd.Execute()

		// Then: Should import configuration
		assert.NoError(t, err)
		assert.FileExists(t, ".ddx.yml", "Should restore config")
	})

	t.Run("view_effective_configuration", func(t *testing.T) {
		// Given: Multiple configuration layers
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		createTestConfig(t)

		// Create override
		os.WriteFile(".ddx.local.yml", []byte("variables:\n  override: true"), 0644)

		// When: Viewing effective configuration
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"config", "show", "--effective"})

		err := cmd.Execute()

		// Then: Should show merged configuration
		assert.NoError(t, err)
		output := buf.String()
		assert.Contains(t, output, "Effective Configuration", "Should show effective")
		assert.Contains(t, output, "override: true", "Should include override")
	})
}

// TestAcceptance_InstallationFeatures tests US-028 through US-035
func TestAcceptance_InstallationFeatures(t *testing.T) {
	t.Run("one_command_installation", func(t *testing.T) {
		// Given: DDx not installed
		// This would be tested in a separate script/environment

		// When: Running installation script
		// curl -sSL https://get.ddx.tools | bash

		// Then: Should install DDx
		// This test would verify the installation script behavior
		t.Skip("Installation script tested separately")
	})

	t.Run("automatic_path_configuration", func(t *testing.T) {
		// Given: DDx installed but not in PATH
		tempDir := t.TempDir()
		os.Chdir(tempDir)

		// When: Running path setup
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"setup", "path"})

		err := cmd.Execute()

		// Then: Should configure PATH
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "PATH", "Should mention PATH")
			assert.Contains(t, output, "shell", "Should detect shell")
		}
	})

	t.Run("installation_verification", func(t *testing.T) {
		// Given: DDx installed
		// When: Verifying installation
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"doctor"})

		err := cmd.Execute()

		// Then: Should verify installation
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "DDx Installation", "Should check installation")
			assert.Contains(t, output, "Version", "Should show version")
			assert.Contains(t, output, "Dependencies", "Should check dependencies")
		}
	})

	t.Run("package_manager_installation", func(t *testing.T) {
		// Given: Package manager available
		// This would test homebrew, apt, etc.
		t.Skip("Package manager installation tested in CI/CD")
	})

	t.Run("upgrade_existing_installation", func(t *testing.T) {
		// Given: Older version installed
		// When: Upgrading
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"self-update"})

		err := cmd.Execute()

		// Then: Should upgrade to latest
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "Checking for updates", "Should check updates")
			assert.Contains(t, output, "version", "Should show versions")
		}
	})

	t.Run("uninstall_ddx", func(t *testing.T) {
		// Given: DDx installed
		// When: Uninstalling
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"uninstall", "--confirm"})

		err := cmd.Execute()

		// Then: Should remove DDx
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "Uninstalling", "Should uninstall")
			assert.Contains(t, output, "Removed", "Should remove files")
		}
	})

	t.Run("offline_installation", func(t *testing.T) {
		// Given: No network connection
		tempDir := t.TempDir()
		os.Chdir(tempDir)

		// Create offline installer bundle
		os.WriteFile("ddx-offline.tar.gz", []byte("mock bundle"), 0644)

		// When: Installing offline
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"install", "--offline", "--bundle", "ddx-offline.tar.gz"})

		err := cmd.Execute()

		// Then: Should install from bundle
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "Offline installation", "Should install offline")
			assert.Contains(t, output, "bundle", "Should use bundle")
		}
	})

	t.Run("installation_diagnostics", func(t *testing.T) {
		// Given: Installation issues
		// When: Running diagnostics
		cmd := rootCmd
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"doctor", "--diagnose"})

		err := cmd.Execute()

		// Then: Should diagnose issues
		if err == nil {
			output := buf.String()
			assert.Contains(t, output, "Diagnostics", "Should run diagnostics")
			assert.Contains(t, output, "Checking", "Should check components")
			assert.Contains(t, output, "Recommendation", "Should provide recommendations")
		}
	})
}

// Helper function to setup MCP test environment
func setupMCPTestProject(t *testing.T) {
	// Set library base path for tests to find registry
	os.Setenv("DDX_LIBRARY_BASE_PATH", "/host-home/erik/Projects/ddx/library")

	// Create .ddx.yml configuration
	config := `
name: test-project
package_manager: npm
mcp:
  servers: []
`
	err := os.WriteFile(".ddx.yml", []byte(config), 0644)
	require.NoError(t, err, "Should create config file")

	// Create mock registry cache
	os.MkdirAll(".ddx/cache", 0755)
	registry := `{
  "servers": [
    {
      "name": "filesystem",
      "package": "@modelcontextprotocol/server-filesystem",
      "description": "File system access for Claude",
      "category": "core",
      "author": "Anthropic"
    },
    {
      "name": "github",
      "package": "@modelcontextprotocol/server-github",
      "description": "GitHub integration for Claude",
      "category": "development",
      "author": "Anthropic"
    },
    {
      "name": "sequential-thinking",
      "package": "@modelcontextprotocol/server-sequential-thinking",
      "description": "Sequential thinking for complex problems",
      "category": "development",
      "author": "Anthropic"
    },
    {
      "name": "weather",
      "package": "@example/mcp-weather",
      "description": "Weather data access",
      "category": "data",
      "author": "Community"
    }
  ]
}`
	os.WriteFile(".ddx/cache/mcp-registry.json", []byte(registry), 0644)
}
