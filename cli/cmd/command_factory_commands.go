package cmd

import (
	"github.com/spf13/cobra"
)

// newInitCommand creates a fresh init command
func (f *CommandFactory) newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize DDx in current project",
		Long: `Initialize DDx in the current project.

This command:
• Creates a .ddx.yml configuration file
• Sets up git subtree for the DDx toolkit
• Downloads essential resources (prompts, templates, patterns)
• Configures project-specific settings

Examples:
  ddx init                  # Interactive setup
  ddx init -t nextjs        # Initialize with Next.js template
  ddx init --force          # Reinitialize existing project`,
		Args: cobra.NoArgs,
		RunE: f.runInit,
	}

	cmd.Flags().StringP("template", "t", "", "Use specific template")
	cmd.Flags().BoolP("force", "f", false, "Force initialization even if DDx already exists")
	cmd.Flags().Bool("no-git", false, "Skip git subtree setup")

	return cmd
}

// newListCommand creates a fresh list command
func (f *CommandFactory) newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [type]",
		Short:   "List available DDx resources",
		Aliases: []string{"ls"},
		Long: `List available DDx resources.

Resources include:
• Templates - Complete project setups
• Patterns - Reusable code patterns
• Prompts - AI interaction prompts
• Scripts - Automation scripts
• Configs - Tool configurations

Examples:
  ddx list              # List all resources
  ddx list templates    # List only templates
  ddx list patterns     # List only patterns`,
		Args: cobra.MaximumNArgs(1),
		RunE: f.runList,
	}

	cmd.Flags().BoolP("detailed", "d", false, "Show detailed information")
	cmd.Flags().StringP("filter", "f", "", "Filter resources by name")
	cmd.Flags().Bool("json", false, "Output results as JSON")
	cmd.Flags().Bool("tree", false, "Display resources in tree format")

	return cmd
}

// newDiagnoseCommand creates a fresh diagnose command
func (f *CommandFactory) newDiagnoseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Analyze project health and suggest improvements",
		Long: `Diagnose analyzes your project and provides recommendations.

This command checks:
• Project structure and configuration
• Development tool setup
• AI integration readiness
• Code quality metrics
• Missing configurations
• Potential improvements

The diagnosis helps identify:
• Configuration issues
• Missing dependencies
• Optimization opportunities
• Best practice violations`,
		Args: cobra.NoArgs,
		RunE: f.runDiagnose,
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show detailed diagnostic output")
	cmd.Flags().Bool("fix", false, "Automatically fix issues where possible")
	cmd.Flags().String("output", "", "Output format (json, yaml, markdown)")

	return cmd
}

// newUpdateCommand creates a fresh update command
func (f *CommandFactory) newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [resource]",
		Short: "Update DDx toolkit from master repository",
		Long: `Update the local DDx toolkit with the latest resources from the master repository.

This command:
• Pulls the latest changes from the master DDx repository
• Updates local resources while preserving customizations
• Uses git subtree for reliable version control
• Creates backups before making changes

You can optionally specify a specific resource to update:
  ddx update templates/nextjs  # Update only the nextjs template
  ddx update prompts           # Update all prompts`,
		Args: cobra.MaximumNArgs(1),
		RunE: runUpdate,
	}

	cmd.Flags().Bool("check", false, "Check for updates without applying")
	cmd.Flags().Bool("force", false, "Force update even if there are local changes")
	cmd.Flags().Bool("reset", false, "Reset to master state, discarding local changes")
	cmd.Flags().Bool("sync", false, "Synchronize with upstream repository")
	cmd.Flags().String("strategy", "", "Conflict resolution strategy (ours/theirs/mine)")
	cmd.Flags().Bool("backup", false, "Create backup before updating")
	cmd.Flags().Bool("interactive", false, "Interactive conflict resolution")
	cmd.Flags().Bool("abort", false, "Abort update and restore previous state")
	cmd.Flags().Bool("mine", false, "Use local changes in conflict resolution")
	cmd.Flags().Bool("theirs", false, "Use upstream changes in conflict resolution")
	cmd.Flags().Bool("dry-run", false, "Preview changes without applying them")

	return cmd
}

// newContributeCommand creates a fresh contribute command
func (f *CommandFactory) newContributeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contribute <path>",
		Short: "Contribute improvements back to master repository",
		Long: `Contribute local improvements back to the master DDx repository.

This command:
• Creates a feature branch in the DDx subtree
• Commits your changes with a descriptive message
• Pushes to your fork (if configured)
• Provides instructions for creating a pull request

Examples:
  ddx contribute patterns/my-pattern
  ddx contribute prompts/claude/new-prompt.md
  ddx contribute scripts/setup/my-script.sh`,
		Args: cobra.ExactArgs(1),
		RunE: runContribute,
	}

	cmd.Flags().StringP("message", "m", "", "Contribution message")
	cmd.Flags().String("branch", "", "Feature branch name")
	cmd.Flags().Bool("dry-run", false, "Show what would be contributed without actually doing it")
	cmd.Flags().Bool("create-pr", false, "Create a pull request after pushing")

	return cmd
}

// newConfigCommand creates a fresh config command
func (f *CommandFactory) newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure DDx settings",
		Long: `Configure DDx settings and preferences.

This command allows you to:
• View current configuration
• Modify settings interactively
• Set individual configuration values
• Reset to defaults

Examples:
  ddx config                    # Interactive configuration
  ddx config --show             # Display current config
  ddx config set key value      # Set specific value
  ddx config get key            # Get specific value`,
		RunE: runConfig,
	}

	cmd.Flags().Bool("show", false, "Display current configuration")
	cmd.Flags().Bool("show-files", false, "Display all config file locations")
	cmd.Flags().Bool("edit", false, "Edit configuration file directly")
	cmd.Flags().Bool("reset", false, "Reset to default configuration")
	cmd.Flags().Bool("wizard", false, "Run configuration wizard")
	cmd.Flags().Bool("validate", false, "Validate configuration")
	cmd.Flags().Bool("preview", false, "Preview resource selection")
	cmd.Flags().Bool("global", false, "Use global configuration")

	// Enhanced validation flags for US-022
	cmd.Flags().String("file", "", "Validate specific configuration file")
	cmd.Flags().Bool("verbose", false, "Detailed validation output")
	cmd.Flags().Bool("offline", false, "Skip network checks during validation")

	// Enhanced config show flags for US-024
	cmd.Flags().String("format", "yaml", "Output format (yaml, json, table)")
	cmd.Flags().Bool("only-overrides", false, "Show only overridden values")
	cmd.Flags().String("filter", "", "Filter by key pattern")

	return cmd
}

// newWorkflowCommand creates a fresh workflow command
func (f *CommandFactory) newWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage development workflows",
		Long: `Manage structured development workflows like HELIX.

Workflows provide:
• Structured development phases
• Gate criteria for phase transitions
• Automated enforcement of best practices
• Progress tracking and reporting

Examples:
  ddx workflow status           # Show current workflow state
  ddx workflow list             # List available workflows
  ddx workflow activate helix   # Activate HELIX workflow
  ddx workflow advance          # Move to next phase`,
		RunE: runWorkflow,
	}

	return cmd
}

// newPersonaCommand creates a fresh persona command
func (f *CommandFactory) newPersonaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "persona",
		Short: "Manage AI personas for consistent interactions",
		Long: `Manage AI personas for consistent role-based interactions.

Personas provide:
• Consistent AI behavior across team members
• Specialized expertise for different roles
• Reusable personality templates
• Project-specific persona bindings

Examples:
  ddx persona list              # List available personas
  ddx persona show reviewer     # Show persona details
  ddx persona bind code-reviewer strict-reviewer`,
		RunE: runPersona,
	}

	cmd.Flags().Bool("list", false, "List available personas")
	cmd.Flags().String("show", "", "Show details of a specific persona")
	cmd.Flags().String("bind", "", "Bind a persona to a role")
	cmd.Flags().String("role", "", "Role to bind persona to or filter by")
	cmd.Flags().String("tag", "", "Filter personas by tag")

	return cmd
}

// newMCPCommand creates a fresh mcp command
func (f *CommandFactory) newMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Manage Model Context Protocol servers",
		Long: `Manage Model Context Protocol (MCP) servers.

MCP servers provide:
• Extended capabilities for AI assistants
• Tool integrations (GitHub, Slack, etc.)
• Custom data sources and APIs
• Enhanced context awareness

Examples:
  ddx mcp list                  # List available MCP servers
  ddx mcp install github        # Install GitHub MCP server
  ddx mcp status                # Show installed servers`,
		RunE: runMCP,
	}

	cmd.Flags().Bool("list", false, "List available MCP servers")
	cmd.Flags().String("install", "", "Install an MCP server")
	cmd.Flags().Bool("status", false, "Show status of installed servers")
	cmd.Flags().String("category", "", "Filter by category")
	cmd.Flags().String("search", "", "Search for servers")
	cmd.Flags().Bool("verbose", false, "Show detailed information")
	cmd.Flags().StringSlice("env", []string{}, "Environment variables for server")
	cmd.Flags().String("config-path", "", "Path to Claude config file")
	cmd.Flags().Bool("dry-run", false, "Show what would be done without making changes")
	cmd.Flags().Bool("yes", false, "Skip confirmation prompts")

	return cmd
}

// newUninstallCommand creates a fresh uninstall command
func (f *CommandFactory) newUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall DDx from the system",
		Long: `Completely remove DDx from your system.

This command will:
• Remove the DDx binary
• Clean up configuration files
• Remove shell completions
• Optionally remove project files

Use with caution as this action cannot be undone.`,
		RunE: runUninstall,
	}

	cmd.Flags().Bool("keep-config", false, "Keep configuration files")
	cmd.Flags().Bool("keep-projects", false, "Keep .ddx directories in projects")
	cmd.Flags().BoolP("force", "f", false, "Skip confirmation prompts")

	return cmd
}

// newPromptsListCommand creates the prompts list subcommand
func (f *CommandFactory) newPromptsListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available prompts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPromptsList(cmd, args)
		},
	}
	cmd.Flags().String("search", "", "Search for prompts containing this text")
	return cmd
}

// newPromptsShowCommand creates the prompts show subcommand
func (f *CommandFactory) newPromptsShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <prompt-name>",
		Short: "Show a specific prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPromptsShow(cmd, args)
		},
	}
}

// newStatusCommand creates a fresh status command
func (f *CommandFactory) newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show version and status information",
		Long: `Show comprehensive version and status information for your DDX project.

This command displays:
- Current DDX version and commit hash
- Last update timestamp
- Local modifications to DDX resources
- Available upstream updates
- Change history and differences

Examples:
  ddx status                          # Show basic status
  ddx status --verbose                # Show detailed information
  ddx status --check-upstream         # Check for updates
  ddx status --changes                # List changed files
  ddx status --diff                   # Show differences
  ddx status --export manifest.yml    # Export version manifest`,
		RunE: runStatus,
	}

	cmd.Flags().Bool("check-upstream", false, "Check for upstream updates")
	cmd.Flags().Bool("changes", false, "Show list of changed files")
	cmd.Flags().Bool("diff", false, "Show differences between versions")
	cmd.Flags().String("export", "", "Export version manifest to file")

	return cmd
}

// newLogCommand creates a fresh log command
func (f *CommandFactory) newLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Show DDX asset history",
		Long: `Show commit history for DDX assets and resources.

This command displays the git log for your DDX resources, helping you
track changes, updates, and the evolution of your project setup.

Examples:
  ddx log                    # Show recent commit history
  ddx log -n 10              # Show last 10 commits
  ddx log --oneline          # Show compact format
  ddx log --since="1 week ago" # Show commits from last week`,
		RunE: runLog,
	}

	cmd.Flags().IntP("number", "n", 20, "Number of commits to show")
	cmd.Flags().Int("limit", 20, "Limit number of commits to show (same as --number)")
	cmd.Flags().Bool("oneline", false, "Show compact one-line format")
	cmd.Flags().Bool("diff", false, "Show changes in each commit")
	cmd.Flags().String("export", "", "Export history to file (format: .md, .json, .csv, .html)")
	cmd.Flags().String("since", "", "Show commits since date (e.g., '1 week ago')")
	cmd.Flags().String("author", "", "Filter by author")
	cmd.Flags().String("grep", "", "Filter by commit message")

	return cmd
}
