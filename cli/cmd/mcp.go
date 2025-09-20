package cmd

import (
	"fmt"
	"strings"

	"github.com/easel/ddx/internal/mcp"
	"github.com/spf13/cobra"
)

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

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Manage MCP servers for Claude Code and Desktop",
	Long: `MCP (Model Context Protocol) server management for Claude.

The mcp command allows you to discover, install, configure, and manage
MCP servers that extend Claude's capabilities with external tools and data sources.

Examples:
  ddx mcp list                    # List available MCP servers
  ddx mcp install github          # Install GitHub MCP server
  ddx mcp status                  # Check status of installed servers
  ddx mcp remove github           # Remove an installed server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand, show help
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)

	// Add subcommands
	mcpCmd.AddCommand(newMCPListCommand())
	mcpCmd.AddCommand(newMCPInstallCommand())
	mcpCmd.AddCommand(newMCPConfigureCommand())
	mcpCmd.AddCommand(newMCPRemoveCommand())
	mcpCmd.AddCommand(newMCPStatusCommand())
	mcpCmd.AddCommand(newMCPUpdateCommand())

	// Global flags for MCP commands
	mcpCmd.PersistentFlags().String("registry", "", "Custom registry URL or path")
	mcpCmd.PersistentFlags().Bool("no-cache", false, "Bypass registry cache")
	mcpCmd.PersistentFlags().String("claude-type", "auto", "Target Claude type (code/desktop/auto)")
}

// newMCPListCommand creates the list subcommand
func newMCPListCommand() *cobra.Command {
	var (
		category  string
		search    string
		installed bool
		available bool
		verbose   bool
		format    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available MCP servers",
		Long: `List all available MCP servers from the registry.

Filter by category, search by name or description, and see installation status.`,
		Example: `  ddx mcp list                     # List all servers
  ddx mcp list --category database  # List database servers
  ddx mcp list --search git         # Search for git-related servers
  ddx mcp list --installed          # Show only installed servers`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract options from CLI flags
			opts := extractListOptions(category, search, installed, available, verbose, format)

			// Create and invoke registry service
			registry, err := mcp.LoadRegistry("")
			if err != nil {
				return fmt.Errorf("loading registry: %w", err)
			}

			// Use the command's output writer for proper output capture in tests
			return registry.ListWithWriter(cmd.OutOrStdout(), opts)
		},
	}

	cmd.Flags().StringVarP(&category, "category", "c", "", "Filter by category")
	cmd.Flags().StringVarP(&search, "search", "s", "", "Search term")
	cmd.Flags().BoolVarP(&installed, "installed", "i", false, "Show only installed servers")
	cmd.Flags().BoolVarP(&available, "available", "a", false, "Show only available servers")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed information")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table/json/yaml)")

	return cmd
}

// newMCPInstallCommand creates the install subcommand
func newMCPInstallCommand() *cobra.Command {
	var (
		env        []string
		yes        bool
		configPath string
		// noBackup   bool // TODO: implement backup feature
		// dryRun     bool  // TODO: implement dry-run feature
	)

	cmd := &cobra.Command{
		Use:   "install <server-name>",
		Short: "Install an MCP server",
		Long: `Install and configure an MCP server for Claude.

This command will:
1. Download the server definition
2. Prompt for required configuration
3. Update your Claude configuration
4. Verify the installation`,
		Example: `  ddx mcp install github                          # Interactive installation
  ddx mcp install postgres --env DATABASE_URL=... # With environment variable
  ddx mcp install github --yes                    # Skip confirmations`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName := args[0]

			// Extract options from CLI flags
			opts := extractInstallOptions(cmd, env, yes, configPath)

			// Create and invoke installer service
			installer := mcp.NewInstallerWithWriter(cmd.OutOrStdout())
			return installer.Install(serverName, opts)
		},
	}

	cmd.Flags().StringArrayVarP(&env, "env", "e", nil, "Environment variables (KEY=VALUE)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompts")
	cmd.Flags().StringVarP(&configPath, "config-path", "p", "", "Custom config file path")
	cmd.Flags().Bool("no-backup", false, "Skip configuration backup")
	cmd.Flags().Bool("dry-run", false, "Show what would be done")

	return cmd
}

// newMCPConfigureCommand creates the configure subcommand
func newMCPConfigureCommand() *cobra.Command {
	var (
		env       []string
		addEnv    []string
		removeEnv []string
		reset     bool
	)

	cmd := &cobra.Command{
		Use:   "configure <server-name>",
		Short: "Configure an installed MCP server",
		Long:  `Update configuration for an installed MCP server.`,
		Example: `  ddx mcp configure github --env GITHUB_TOKEN=new_token
  ddx mcp configure postgres --add-env POOL_SIZE=10
  ddx mcp configure github --reset`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName := args[0]

			// Extract options from CLI flags
			opts := extractConfigureOptions(env, addEnv, removeEnv, reset)

			// Create and invoke config manager service
			configManager := mcp.NewConfigManager("")
			return configManager.UpdateServer(serverName, opts, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringArrayVarP(&env, "env", "e", nil, "Set environment variables")
	cmd.Flags().StringArrayVar(&addEnv, "add-env", nil, "Add environment variables")
	cmd.Flags().StringArrayVar(&removeEnv, "remove-env", nil, "Remove environment variables")
	cmd.Flags().BoolVar(&reset, "reset", false, "Reset to defaults")

	return cmd
}

// newMCPRemoveCommand creates the remove subcommand
func newMCPRemoveCommand() *cobra.Command {
	var (
		yes bool
		// noBackup bool // TODO: implement backup feature
		// purge    bool // TODO: implement purge feature
	)

	cmd := &cobra.Command{
		Use:   "remove <server-name>",
		Short: "Remove an installed MCP server",
		Long:  `Remove an MCP server configuration from Claude.`,
		Example: `  ddx mcp remove github          # Remove with confirmation
  ddx mcp remove github --yes     # Skip confirmation
  ddx mcp remove github --purge   # Remove all related data`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverName := args[0]

			// Extract options from CLI flags
			opts := extractRemoveOptions(cmd)

			// Create and invoke config manager service
			configManager := mcp.NewConfigManager("")
			return configManager.RemoveServer(serverName, opts, cmd.OutOrStdout())
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation")
	cmd.Flags().Bool("no-backup", false, "Skip backup creation")
	cmd.Flags().Bool("purge", false, "Remove all related data")

	return cmd
}

// newMCPStatusCommand creates the status subcommand
func newMCPStatusCommand() *cobra.Command {
	var (
		check   bool
		verbose bool
		format  string
	)

	cmd := &cobra.Command{
		Use:   "status [server-name]",
		Short: "Show status of MCP servers",
		Long:  `Display the status of installed MCP servers.`,
		Example: `  ddx mcp status              # Show all servers
  ddx mcp status github       # Show specific server
  ddx mcp status --check      # Verify connectivity`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var serverName string
			if len(args) > 0 {
				serverName = args[0]
			}

			// Extract options from CLI flags
			opts := extractStatusOptions(cmd, serverName)

			// Create and invoke status checker service
			statusChecker := mcp.NewStatusCheckerWithWriter(cmd.OutOrStdout())
			return statusChecker.Check(opts)
		},
	}

	cmd.Flags().BoolVarP(&check, "check", "c", false, "Verify server connectivity")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed information")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format")

	return cmd
}

// newMCPUpdateCommand creates the update subcommand
func newMCPUpdateCommand() *cobra.Command {
	var (
		force  bool
		server string
		check  bool
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update MCP server registry",
		Long:  `Update the MCP server registry to get the latest available servers.`,
		Example: `  ddx mcp update              # Update registry
  ddx mcp update --check      # Check for updates
  ddx mcp update --force      # Force update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Extract options from CLI flags
			opts := extractUpdateOptions(cmd)

			// Create and invoke registry service
			registry, err := mcp.LoadRegistry("")
			if err != nil {
				return fmt.Errorf("loading registry: %w", err)
			}

			return registry.Update(opts)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force update")
	cmd.Flags().StringVarP(&server, "server", "s", "", "Update specific server")
	cmd.Flags().BoolVarP(&check, "check", "c", false, "Check for updates only")

	return cmd
}
