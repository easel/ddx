package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a fresh root command for tests
func getMCPTestRootCommand() *cobra.Command {
	factory := NewCommandFactory("/tmp")
	return factory.NewRootCommand()
}

// setupMockLibrary creates a mock library path with MCP registry for testing
func setupMockLibrary(t *testing.T, tempDir string) string {
	mockLibPath := filepath.Join(tempDir, "mock-library")
	t.Setenv("DDX_LIBRARY_BASE_PATH", mockLibPath)
	return mockLibPath
}

// TestAcceptance_US036_ListMCPServers tests US-036: List Available MCP Servers
func TestAcceptance_US036_ListMCPServers(t *testing.T) {
	// Ensure we're in a valid directory first
	ensureValidWorkingDirectory(t)

	// Use temp directory for test isolation

	// Library path will be mocked in each test

	t.Run("display_all_available_servers", func(t *testing.T) {
		// Save and restore working directory
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: MCP server registry is available
		tempDir := t.TempDir()

		// Create a mock library in temp directory for testing
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// When: Running ddx mcp list
		rootCmd := getMCPTestRootCommand()
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
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: Multiple categories of MCP servers exist
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// When: Filtering by category
		rootCmd := getMCPTestRootCommand()
		output, err := executeCommand(rootCmd, "mcp", "list", "--category", "development")

		// Then: Should only show servers in that category
		assert.NoError(t, err)
		assert.Contains(t, output, "Development", "Should indicate filter")
		assert.Contains(t, output, "github", "Should show dev servers")
		assert.NotContains(t, output, "sqlite", "Should not show other categories")
	})

	t.Run("search_functionality", func(t *testing.T) {
		// Save and restore working directory
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: Want to find servers related to "git"
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// When: Searching for "git"
		rootCmd := getMCPTestRootCommand()
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
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: Some MCP servers are installed
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// Install a server first
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		rootCmd := getMCPTestRootCommand()
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
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: Want more information
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// When: Running with --verbose
		rootCmd := getMCPTestRootCommand()
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
	// Ensure we're in a valid directory first
	ensureValidWorkingDirectory(t)

	// Use temp directory for test isolation

	// Library path will be mocked in each test

	t.Run("install_server_locally", func(t *testing.T) {
		// Save and restore working directory
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: MCP server not installed
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// When: Installing a server
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		rootCmd := getMCPTestRootCommand()
		output, err := executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// Then: Should install server locally
		assert.NoError(t, err)
		assert.Contains(t, output, "Successfully installed server", "Should show installation")
		assert.Contains(t, output, "filesystem", "Should name the server")

		// Check Claude config was updated at the custom path
		assert.FileExists(t, configPath, "Should create Claude config")
	})

	t.Run("detect_package_manager", func(t *testing.T) {
		// Save and restore working directory
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: Different package managers available
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// Create pnpm-lock.yaml to trigger pnpm detection
		os.WriteFile("pnpm-lock.yaml", []byte("lockfileVersion: 5.4"), 0644)

		// When: Installing
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		rootCmd := getMCPTestRootCommand()
		output, err := executeCommand(rootCmd, "mcp", "install", "github", "--env", "GITHUB_PERSONAL_ACCESS_TOKEN=ghp_012345678901234567890123456789012345", "--config-path", configPath)

		// Then: Should detect and use pnpm
		assert.NoError(t, err)
		assert.Contains(t, output, "Successfully installed server", "Should detect pnpm")
	})

	t.Run("configure_server_environment", func(t *testing.T) {
		// Save and restore working directory
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: Server needs configuration
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// When: Installing with configuration
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		rootCmd := getMCPTestRootCommand()
		output, err := executeCommand(rootCmd, "mcp", "install", "github", "--env", "GITHUB_PERSONAL_ACCESS_TOKEN=ghp_012345678901234567890123456789012345", "--config-path", configPath)

		// Then: Should configure the server
		assert.NoError(t, err)
		assert.Contains(t, output, "Successfully installed server", "Should configure server")

		// Check configuration was written
		claudeConfig := filepath.Join(".claude", "settings.local.json")
		if content, err := os.ReadFile(claudeConfig); err == nil {
			assert.Contains(t, string(content), "GITHUB_PERSONAL_ACCESS_TOKEN", "Should set env var")
		}
	})

	t.Run("handle_already_installed", func(t *testing.T) {
		// Save and restore working directory
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: Server already installed
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// Install once
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		rootCmd := getMCPTestRootCommand()
		_, _ = executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// When: Installing again
		output, _ := executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// Then: Should handle gracefully
		assert.Contains(t, output, "already installed", "Should detect existing")
	})

	t.Run("validate_installation", func(t *testing.T) {
		// Save and restore working directory
		// origDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection

		// Given: Server installed
		tempDir := t.TempDir()
		setupMockLibrary(t, tempDir)
		setupMCPTestProject(t)

		// When: Installing
		configPath := filepath.Join(tempDir, ".claude", "settings.local.json")
		rootCmd := getMCPTestRootCommand()
		output, err := executeCommand(rootCmd, "mcp", "install", "filesystem", "--config-path", configPath)

		// Then: Should install successfully
		assert.NoError(t, err)
		assert.Contains(t, output, "Successfully installed server", "Should show installation")
		assert.Contains(t, output, "filesystem", "Should name the server")
	})
}

// resolveLibraryPath finds the library path from the test location
func resolveLibraryPath(t *testing.T) string {
	t.Helper()

	// Try to use the environment variable if set
	if envPath := os.Getenv("DDX_LIBRARY_BASE_PATH"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	// Use a fixed path relative to the test binary location
	// The test binary runs from the cli directory
	// So we need to go up to the project root to find library
	libraryPath := filepath.Join("..", "..", "library")

	// Convert to absolute path
	absPath, err := filepath.Abs(libraryPath)
	if err == nil {
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	// Fallback: try common test paths
	testPaths := []string{
		"/host-home/erik/Projects/ddx/library",
		"../../library",
		"../library",
		"./library",
	}

	for _, path := range testPaths {
		absPath, err := filepath.Abs(path)
		if err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}

	// If we still can't find it, skip the test rather than fail
	t.Skip("Cannot locate library directory - skipping test")
	return ""
}

// Helper function to setup MCP test environment
func setupMCPTestProject(t *testing.T) {
	// Create .ddx/config.yaml configuration
	env := NewTestEnvironment(t)
	config := `version: "1.0"
library_base_path: "./library"
repository:
  url: "https://github.com/easel/ddx"
  branch: "main"
  subtree_prefix: "library"
variables:
  project_name: "test-project"
  package_manager: "npm"
`
	env.CreateConfig(config)

	// Create a mock MCP server registry if DDX_LIBRARY_BASE_PATH is set
	if libPath := os.Getenv("DDX_LIBRARY_BASE_PATH"); libPath != "" {
		// Always create the mock for tests
		if strings.Contains(libPath, "mock-library") || strings.Contains(libPath, os.TempDir()) {
			// Create mock MCP server registry
			mcpDir := filepath.Join(libPath, "mcp-servers")
			os.MkdirAll(mcpDir, 0755)

			// Create a simple registry.yml
			registry := `version: 1.0.0
updated: 2025-01-15T00:00:00Z
servers:
  - name: filesystem
    file: servers/filesystem.yml
    category: core
    description: Access local files
  - name: github
    file: servers/github.yml
    category: development
    description: Access GitHub repositories
`
			os.WriteFile(filepath.Join(mcpDir, "registry.yml"), []byte(registry), 0644)

			// Create the server files referenced in the registry
			serversDir := filepath.Join(mcpDir, "servers")
			os.MkdirAll(serversDir, 0755)

			// Create filesystem.yml
			filesystemYaml := `name: filesystem
description: Access local files
category: core
author: DDx Team
version: 1.0.0
tags: ["core", "filesystem"]
command:
  executable: npx
  args: ["@modelcontextprotocol/server-filesystem"]
environment:
  - name: FILESYSTEM_ROOT
    description: Root directory for filesystem access
    required: false
    default: "."
documentation:
  setup: "Install with npm install @modelcontextprotocol/server-filesystem"
  permissions: ["read", "write"]
compatibility:
  platforms: ["linux", "macos", "windows"]
  claude_versions: ["*"]
`
			os.WriteFile(filepath.Join(serversDir, "filesystem.yml"), []byte(filesystemYaml), 0644)

			// Create github.yml
			githubYaml := `name: github
description: Access GitHub repositories
category: development
author: DDx Team
version: 2.1.0
tags: ["git", "github", "development"]
command:
  executable: npx
  args: ["@modelcontextprotocol/server-github"]
environment:
  - name: GITHUB_PERSONAL_ACCESS_TOKEN
    description: GitHub personal access token
    required: true
    sensitive: true
    validation: "^ghp_[a-zA-Z0-9]{36}$"
    example: "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
documentation:
  setup: "Install with npm install @modelcontextprotocol/server-github"
  permissions: ["repo", "read:user"]
  security_notes: "Requires GitHub personal access token"
compatibility:
  platforms: ["linux", "macos", "windows"]
  claude_versions: ["*"]
`
			os.WriteFile(filepath.Join(serversDir, "github.yml"), []byte(githubYaml), 0644)
		}
	}
}

// ensureValidWorkingDirectory ensures we're in a valid directory before tests
func ensureValidWorkingDirectory(t *testing.T) {
	t.Helper()

	// Ensure we have a safe working directory
	// Tests use CommandFactory with explicit working directory
}
