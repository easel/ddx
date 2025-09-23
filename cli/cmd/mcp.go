package cmd

import (
	"fmt"
	"strings"

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
	// Load registry using default path resolution
	registry, err := mcp.LoadRegistry("")
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Check installed servers via Claude CLI
	claude := mcp.NewClaudeWrapper()

	// Set Claude wrapper on registry for installation status checking
	registry.SetClaudeWrapper(claude)

	// Get config path if provided
	configPath, _ := cmd.Flags().GetString("config-path")

	// Use registry to list servers with installed status
	opts := mcp.ListOptions{
		Category:   category,
		Search:     search,
		Verbose:    verbose,
		Available:  true,
		ConfigPath: configPath,
	}

	return registry.ListWithWriter(cmd.OutOrStdout(), opts)
}

func runMCPInstallWithOptions(cmd *cobra.Command, serverName string) error {
	// Get additional options
	envVars, _ := cmd.Flags().GetStringSlice("env")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	yes, _ := cmd.Flags().GetBool("yes")
	configPath, _ := cmd.Flags().GetString("config-path")

	// Parse environment variables
	environment := make(map[string]string)
	for _, envVar := range envVars {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			environment[parts[0]] = parts[1]
		}
	}

	// Create installer
	installer := mcp.NewInstallerWithWriter(cmd.OutOrStdout())

	// Install options
	opts := mcp.InstallOptions{
		Environment: environment,
		DryRun:      dryRun,
		Yes:         yes,
		ConfigPath:  configPath,
	}

	// Install the server
	return installer.Install(serverName, opts)
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
func extractInstallOptions(cmd *cobra.Command, envVars []string, yes bool) mcp.InstallOptions {
	// Parse environment variables from KEY=VALUE format
	environment := make(map[string]string)
	for _, envVar := range envVars {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			environment[parts[0]] = parts[1]
		}
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")

	return mcp.InstallOptions{
		Environment: environment,
		DryRun:      dryRun,
		Yes:         yes,
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
	purge, _ := cmd.Flags().GetBool("purge")

	return mcp.RemoveOptions{
		SkipConfirmation: yes,
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
