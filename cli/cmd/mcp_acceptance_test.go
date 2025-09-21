package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAcceptance_US036_ListMCPServers tests US-036: List Available MCP Servers
func TestAcceptance_US036_ListMCPServers(t *testing.T) {
	// Save current directory before any changes
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Should get working directory")
	defer os.Chdir(originalDir)

	// Resolve library path before changing directories
	libraryPath := resolveLibraryPath(t)

	t.Run("display_all_available_servers", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: MCP server registry is available
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// When: Running ddx mcp list
		output, err := executeCommand(rootCmd, "mcp", "list")

		// Then: Should display all available servers
		assert.NoError(t, err, "Should list MCP servers")
		assert.Contains(t, output, "Available MCP Servers", "Should show header")
		assert.Contains(t, output, "filesystem", "Should show filesystem server")
		assert.Contains(t, output, "github", "Should show github server")
		// All servers start as not installed in a fresh project
		assert.Contains(t, output, "⬜", "Should show not-installed indicator")
	})

	t.Run("filter_by_category", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: Multiple categories of MCP servers exist
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// When: Filtering by category
		output, err := executeCommand(rootCmd, "mcp", "list", "--category", "development")

		// Then: Should only show servers in that category
		assert.NoError(t, err)
		assert.Contains(t, output, "Development", "Should indicate filter")
		assert.Contains(t, output, "github", "Should show dev servers")
		assert.NotContains(t, output, "sqlite", "Should not show other categories")
	})

	t.Run("search_functionality", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: Want to find servers related to "git"
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// When: Searching for "git"
		output, err := executeCommand(rootCmd, "mcp", "list", "--search", "git")

		// Then: Should show matching servers
		assert.NoError(t, err)
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
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: Some MCP servers are installed
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// Install a server first
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		_, _ = executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// When: Listing servers
		output, err := executeCommand(rootCmd, "mcp", "list")

		// Then: Should show correct installation status
		assert.NoError(t, err)
		// Find filesystem line and check it has installed indicator
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "filesystem") {
				assert.Contains(t, line, "✅", "Installed server should show ✅")
			}
		}
	})

	t.Run("detailed_verbose_view", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: Want more information
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// When: Running with --verbose
		output, err := executeCommand(rootCmd, "mcp", "list", "--verbose")

		// Then: Should show additional details
		assert.NoError(t, err)
		assert.Contains(t, output, "Environment", "Should show env vars")
		assert.Contains(t, output, "Package", "Should show package info")
		assert.Contains(t, output, "Author", "Should show author")
		assert.Contains(t, output, "Version", "Should show version")
	})
}

// TestAcceptance_US037_InstallMCPServer tests US-037: Install MCP Server
func TestAcceptance_US037_InstallMCPServer(t *testing.T) {
	// Save current directory before any changes
	originalDir, err := os.Getwd()
	require.NoError(t, err, "Should get working directory")
	defer os.Chdir(originalDir)

	// Resolve library path before changing directories
	libraryPath := resolveLibraryPath(t)

	t.Run("install_server_locally", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: MCP server not installed
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// When: Installing a server
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		output, err := executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// Then: Should install server locally
		assert.NoError(t, err)
		assert.Contains(t, output, "Installing", "Should show installation")
		assert.Contains(t, output, "filesystem", "Should name the server")
		assert.Contains(t, output, "successfully", "Should indicate success")

		// Check Claude config was updated at the custom path
		assert.FileExists(t, configPath, "Should create Claude config")
	})

	t.Run("detect_package_manager", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: Different package managers available
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// Create pnpm-lock.yaml to trigger pnpm detection
		os.WriteFile("pnpm-lock.yaml", []byte("lockfileVersion: 5.4"), 0644)

		// When: Installing
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		output, err := executeCommand(rootCmd, "mcp", "install", "github", "--env", "GITHUB_PERSONAL_ACCESS_TOKEN=ghp_012345678901234567890123456789012345", "--config-path", configPath)

		// Then: Should detect and use pnpm
		assert.NoError(t, err)
		assert.Contains(t, output, "pnpm", "Should detect pnpm")
		assert.Contains(t, output, "Using package manager: pnpm", "Should indicate pnpm usage")
	})

	t.Run("configure_server_environment", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: Server needs configuration
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// When: Installing with configuration
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		output, err := executeCommand(rootCmd, "mcp", "install", "github", "--env", "GITHUB_PERSONAL_ACCESS_TOKEN=ghp_012345678901234567890123456789012345", "--config-path", configPath)

		// Then: Should configure the server
		assert.NoError(t, err)
		assert.Contains(t, output, "Configuring", "Should configure server")

		// Check configuration was written
		claudeConfig := filepath.Join(".claude", "settings.local.json")
		if content, err := os.ReadFile(claudeConfig); err == nil {
			assert.Contains(t, string(content), "GITHUB_TOKEN", "Should set env var")
		}
	})

	t.Run("handle_already_installed", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: Server already installed
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// Install once
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		_, _ = executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// When: Installing again
		output, _ := executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// Then: Should handle gracefully
		assert.Contains(t, output, "already installed", "Should detect existing")
		assert.Contains(t, output, "upgrade", "Should offer upgrade option")
	})

	t.Run("validate_installation", func(t *testing.T) {
		// Save and restore working directory
		origDir, _ := os.Getwd()
		defer os.Chdir(origDir)

		// Given: Server installed
		tempDir := t.TempDir()
		os.Chdir(tempDir)
		t.Setenv("DDX_LIBRARY_BASE_PATH", libraryPath)
		setupMCPTestProject(t)

		// When: Installing
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		output, err := executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// Then: Should install successfully
		assert.NoError(t, err)
		assert.Contains(t, output, "Installing", "Should show installation")
		assert.Contains(t, output, "filesystem", "Should name the server")
	})
}

// resolveLibraryPath finds the library path from the test location
func resolveLibraryPath(t *testing.T) string {
	t.Helper()

	// Get the current working directory (should be cli/cmd when tests run)
	pwd, err := os.Getwd()
	require.NoError(t, err, "Should get working directory")

	// Navigate to project root and find library
	// Tests run from cli directory, so go up one level to find library
	projectRoot := filepath.Dir(pwd)
	libraryPath := filepath.Join(projectRoot, "library")

	// Verify the library exists
	if _, err := os.Stat(libraryPath); os.IsNotExist(err) {
		// If not found, try going up one more level (in case we're in cli/cmd)
		projectRoot = filepath.Dir(projectRoot)
		libraryPath = filepath.Join(projectRoot, "library")
	}

	require.DirExists(t, libraryPath, "Library directory should exist")
	return libraryPath
}

// Helper function to setup MCP test environment
func setupMCPTestProject(t *testing.T) {
	// Create .ddx.yml configuration
	config := `
name: test-project
package_manager: npm
mcp:
  servers: []
`
	err := os.WriteFile(".ddx.yml", []byte(config), 0644)
	require.NoError(t, err, "Should create config file")
}
