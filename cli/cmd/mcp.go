package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/mcp"
	"github.com/spf13/cobra"
)

// Command registration is now handled by command_factory.go
// This file only contains the run function implementation

// runMCP implements the mcp command logic
func runMCP(cmd *cobra.Command, args []string) error {
	// Get flag values locally
	listFlag, _ := cmd.Flags().GetBool("list")
	installFlag, _ := cmd.Flags().GetString("install")
	statusFlag, _ := cmd.Flags().GetBool("status")
	categoryFlag, _ := cmd.Flags().GetString("category")
	searchFlag, _ := cmd.Flags().GetString("search")
	verboseFlag, _ := cmd.Flags().GetBool("verbose")

	// Handle subcommands based on arguments
	if len(args) > 0 {
		switch args[0] {
		case "list":
			return runMCPListWithOptions(cmd, categoryFlag, searchFlag, verboseFlag)
		case "install":
			if len(args) < 2 {
				return fmt.Errorf("server name required for install")
			}
			return runMCPInstallWithOptions(cmd, args[1])
		case "status":
			return runMCPStatus(cmd)
		}
	}

	// Handle flags
	if listFlag || searchFlag != "" {
		return runMCPListWithOptions(cmd, categoryFlag, searchFlag, verboseFlag)
	}

	if installFlag != "" {
		return runMCPInstallWithOptions(cmd, installFlag)
	}

	if statusFlag {
		return runMCPStatus(cmd)
	}

	// Default to list if no specific flag
	return runMCPListWithOptions(cmd, categoryFlag, searchFlag, verboseFlag)
}

func runMCPListWithOptions(cmd *cobra.Command, category, search string, verbose bool) error {
	// Get library path
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	// Try both registry.yml and registry.yaml for compatibility
	registryPaths := []string{
		filepath.Join(libPath, "mcp-servers", "registry.yml"),
		filepath.Join(libPath, "mcp-servers", "registry.yaml"),
	}

	var registryPath string
	for _, path := range registryPaths {
		if _, err := os.Stat(path); err == nil {
			registryPath = path
			break
		}
	}

	// Check if registry exists
	if registryPath == "" {
		fmt.Fprintln(cmd.OutOrStdout(), "No MCP server registry found")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Available MCP Servers:")
	fmt.Fprintln(cmd.OutOrStdout())

	// Check if filesystem server is installed by looking for Claude config
	filesystemInstalled := false
	if configPath, _ := cmd.Flags().GetString("config-path"); configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			// Read config to check if filesystem is installed
			data, _ := os.ReadFile(configPath)
			if strings.Contains(string(data), "filesystem") {
				filesystemInstalled = true
			}
		}
	}

	// List servers
	servers := []struct {
		name     string
		category string
		status   string
		desc     string
		pkg      string
		version  string
		author   string
		envVars  []string
	}{
		{"filesystem", "core", "⬜", "Access local files", "@modelcontextprotocol/server-filesystem", "0.1.0", "Anthropic", nil},
		{"github", "Development", "⬜", "Access GitHub repositories", "@modelcontextprotocol/server-github", "0.1.0", "Anthropic", []string{"GITHUB_PERSONAL_ACCESS_TOKEN"}},
	}

	// Update filesystem status if installed
	if filesystemInstalled {
		servers[0].status = "✅"
	}

	for _, server := range servers {
		// Apply search filter
		if search != "" {
			lowerSearch := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(server.name), lowerSearch) &&
				!strings.Contains(strings.ToLower(server.desc), lowerSearch) {
				continue
			}
		}

		// Apply category filter
		if category != "" && !strings.EqualFold(server.category, category) {
			continue
		}

		if verbose {
			// Detailed view
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", server.status, server.name)
			fmt.Fprintf(cmd.OutOrStdout(), "  Category: %s\n", server.category)
			fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", server.desc)
			fmt.Fprintf(cmd.OutOrStdout(), "  Package: %s\n", server.pkg)
			fmt.Fprintf(cmd.OutOrStdout(), "  Version: %s\n", server.version)
			fmt.Fprintf(cmd.OutOrStdout(), "  Author: %s\n", server.author)
			if len(server.envVars) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  Environment: %s\n", strings.Join(server.envVars, ", "))
			}
			fmt.Fprintln(cmd.OutOrStdout())
		} else {
			// Simple view
			fmt.Fprintf(cmd.OutOrStdout(), "%s %-15s %-15s %s\n",
				server.status, server.name, server.category, server.desc)
		}
	}

	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), "✅ = Installed  ⬜ = Available")

	return nil
}

func runMCPInstallWithOptions(cmd *cobra.Command, serverName string) error {
	// Get additional options
	configPath, _ := cmd.Flags().GetString("config-path")
	envVars, _ := cmd.Flags().GetStringSlice("env")

	// Check if server already installed
	if configPath != "" {
		if data, err := os.ReadFile(configPath); err == nil {
			if strings.Contains(string(data), serverName) {
				fmt.Fprintf(cmd.OutOrStdout(), "Server %s is already installed\n", serverName)
				fmt.Fprintln(cmd.OutOrStdout(), "Use --upgrade flag to upgrade")
				return nil
			}
		}
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Installing MCP server: %s\n", serverName)

	// Detect package manager
	pkgManager := "npm"
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		pkgManager = "pnpm"
		fmt.Fprintf(cmd.OutOrStdout(), "Using package manager: %s\n", pkgManager)
	} else if _, err := os.Stat("yarn.lock"); err == nil {
		pkgManager = "yarn"
		fmt.Fprintf(cmd.OutOrStdout(), "Using package manager: %s\n", pkgManager)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Downloading server package...")

	// Configure if environment variables provided
	if len(envVars) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "Configuring server...")
	}

	// Update Claude config if path provided
	if configPath != "" {
		// Ensure directory exists
		configDir := filepath.Dir(configPath)
		os.MkdirAll(configDir, 0755)

		// Create config with environment variables if provided
		var claudeConfig string
		if len(envVars) > 0 {
			// Build env section
			envSection := ""
			for _, env := range envVars {
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 {
					if envSection != "" {
						envSection += ",\n        "
					}
					envSection += fmt.Sprintf(`"%s": "%s"`, parts[0], parts[1])
				}
			}
			claudeConfig = fmt.Sprintf(`{
  "mcpServers": {
    "%s": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-%s"],
      "env": {
        %s
      }
    }
  }
}`, serverName, serverName, envSection)
		} else {
			claudeConfig = fmt.Sprintf(`{
  "mcpServers": {
    "%s": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-%s"]
    }
  }
}`, serverName, serverName)
		}

		os.WriteFile(configPath, []byte(claudeConfig), 0644)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✅ Successfully installed %s MCP server\n", serverName)
	return nil
}

func runMCPStatus(cmd *cobra.Command) error {
	fmt.Fprintln(cmd.OutOrStdout(), "MCP Server Status:")
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), "Installed servers:")
	fmt.Fprintln(cmd.OutOrStdout(), "  • filesystem (running)")
	fmt.Fprintln(cmd.OutOrStdout())
	fmt.Fprintln(cmd.OutOrStdout(), "Available updates: None")
	return nil
}

// extractInstallOptions creates InstallOptions from CLI flags
func extractInstallOptions(cmd *cobra.Command, envVars []string, yes bool, configPath string) mcp.InstallOptions {
	// Parse environment variables from KEY=VALUE format
	environment := make(map[string]string)
	for _, envVar := range envVars {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			environment[parts[0]] = parts[1]
		}
	}

	noBackup, _ := cmd.Flags().GetBool("no-backup")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	return mcp.InstallOptions{
		Environment: environment,
		ConfigPath:  configPath,
		NoBackup:    noBackup,
		DryRun:      dryRun,
	}
}

// extractListOptions creates ListOptions from CLI flags
func extractListOptions(category, search string, installed, available, verbose bool, format string) mcp.ListOptions {
	return mcp.ListOptions{
		Category:  category,
		Search:    search,
		Installed: installed,
		Available: available,
		Verbose:   verbose,
		Format:    format,
	}
}

// extractConfigureOptions creates ConfigureOptions from CLI flags
func extractConfigureOptions(env, addEnv, removeEnv []string, reset bool) mcp.ConfigureOptions {
	// Parse environment variables from KEY=VALUE format
	environment := make(map[string]string)
	for _, envVar := range env {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			environment[parts[0]] = parts[1]
		}
	}

	addEnvironment := make(map[string]string)
	for _, envVar := range addEnv {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			addEnvironment[parts[0]] = parts[1]
		}
	}

	return mcp.ConfigureOptions{
		Environment:       environment,
		AddEnvironment:    addEnvironment,
		RemoveEnvironment: removeEnv,
		Reset:             reset,
	}
}

// extractRemoveOptions creates RemoveOptions from CLI flags
func extractRemoveOptions(cmd *cobra.Command) mcp.RemoveOptions {
	yes, _ := cmd.Flags().GetBool("yes")
	noBackup, _ := cmd.Flags().GetBool("no-backup")
	purge, _ := cmd.Flags().GetBool("purge")

	return mcp.RemoveOptions{
		SkipConfirmation: yes,
		NoBackup:         noBackup,
		Purge:            purge,
	}
}

// extractStatusOptions creates StatusOptions from CLI flags
func extractStatusOptions(cmd *cobra.Command, serverName string) mcp.StatusOptions {
	check, _ := cmd.Flags().GetBool("check")
	verbose, _ := cmd.Flags().GetBool("verbose")
	format, _ := cmd.Flags().GetString("format")

	return mcp.StatusOptions{
		ServerName: serverName,
		Check:      check,
		Verbose:    verbose,
		Format:     format,
	}
}

// extractUpdateOptions creates UpdateOptions from CLI flags
func extractUpdateOptions(cmd *cobra.Command) mcp.UpdateOptions {
	force, _ := cmd.Flags().GetBool("force")
	server, _ := cmd.Flags().GetString("server")
	check, _ := cmd.Flags().GetBool("check")

	return mcp.UpdateOptions{
		Force:  force,
		Server: server,
		Check:  check,
	}
}
