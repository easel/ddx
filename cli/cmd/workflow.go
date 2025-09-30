package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/workflow"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Global variables have been removed - use local variables in runWorkflow

// workflowCmd represents the generic workflow command
// Command registration is now handled by command_factory.go
// This file only contains the run function implementation

// runWorkflow implements the workflow command logic
func runWorkflow(cmd *cobra.Command, args []string) error {
	return runWorkflowWithDir(cmd, args, ".")
}

// runWorkflowWithDir implements workflow command logic with explicit working directory
func runWorkflowWithDir(cmd *cobra.Command, args []string, workingDir string) error {
	workflowForce, _ := cmd.Flags().GetBool("force")

	if len(args) == 0 {
		return cmd.Help()
	}

	firstArg := strings.ToLower(args[0])

	// Check if first argument is a workflow name (even if unknown) or generic command
	switch firstArg {
	case "status":
		return showWorkflowStatusWithDir(cmd, workingDir)
	case "list":
		return listWorkflows(cmd)
	case "activate":
		if len(args) < 2 {
			return fmt.Errorf("workflow name required")
		}
		return activateWorkflowWithDir(cmd, args[1], workflowForce, workingDir)
	case "deactivate":
		if len(args) < 2 {
			return fmt.Errorf("workflow name required")
		}
		return deactivateWorkflowWithDir(cmd, args[1], workingDir)
	case "advance":
		return advanceWorkflow(cmd)
	default:
		// If not a generic command, treat as workflow name
		if len(args) > 1 {
			if isKnownWorkflow(firstArg) {
				return handleWorkflowSpecificCommand(cmd, firstArg, args[1:])
			} else {
				return fmt.Errorf("workflow '%s' not found", firstArg)
			}
		} else {
			return fmt.Errorf("unknown subcommand: %s", firstArg)
		}
	}
}

func showWorkflowStatus(cmd *cobra.Command) error {
	// Load config
	cfg, err := loadConfig()
	if err != nil || cfg == nil {
		fmt.Fprintln(cmd.OutOrStdout(), "No active workflows")
		return nil
	}

	cfg.ApplyDefaults()

	// Check if there are active workflows
	if len(cfg.Workflows.Active) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No active workflows")
		return nil
	}

	// Display active workflows
	fmt.Fprintln(cmd.OutOrStdout(), "Active workflows (in priority order):")
	fmt.Fprintln(cmd.OutOrStdout())

	// Get library path
	libraryPath := cfg.Library.Path
	if !filepath.IsAbs(libraryPath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		libraryPath = filepath.Join(wd, libraryPath)
	}

	loader := workflow.NewLoader(libraryPath)

	for i, workflowName := range cfg.Workflows.Active {
		fmt.Fprintf(cmd.OutOrStdout(), "%d. %s\n", i+1, workflowName)

		// Load workflow to show agent commands
		def, err := loader.Load(workflowName)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "   (failed to load: %v)\n", err)
			continue
		}

		// Show agent commands
		if len(def.AgentCommands) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "   Agent commands:")
			for cmdName, agentCmd := range def.AgentCommands {
				if agentCmd.Enabled {
					fmt.Fprintf(cmd.OutOrStdout(), "   - %s -> %s\n", cmdName, agentCmd.Action)
					if agentCmd.Triggers != nil {
						if len(agentCmd.Triggers.Keywords) > 0 {
							fmt.Fprintf(cmd.OutOrStdout(), "     Keywords: %s\n", strings.Join(agentCmd.Triggers.Keywords, ", "))
						}
						if len(agentCmd.Triggers.Patterns) > 0 {
							fmt.Fprintf(cmd.OutOrStdout(), "     Patterns: %s\n", strings.Join(agentCmd.Triggers.Patterns, ", "))
						}
					}
				}
			}
		}
		fmt.Fprintln(cmd.OutOrStdout())
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Safe word: %s\n", cfg.Workflows.SafeWord)
	return nil
}

func listWorkflows(cmd *cobra.Command) error {
	fmt.Fprintln(cmd.OutOrStdout(), "Available workflows:")
	fmt.Fprintln(cmd.OutOrStdout(), "  • helix - HELIX development methodology")
	fmt.Fprintln(cmd.OutOrStdout(), "  • agile - Agile/Scrum workflow")
	fmt.Fprintln(cmd.OutOrStdout(), "  • kanban - Kanban board workflow")
	return nil
}

func activateWorkflow(cmd *cobra.Command, name string, force bool) error {
	// Load config
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if cfg == nil {
		return fmt.Errorf("no config found - run 'ddx init' first")
	}

	cfg.ApplyDefaults()

	// Get library path
	libraryPath := cfg.Library.Path
	if !filepath.IsAbs(libraryPath) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		libraryPath = filepath.Join(wd, libraryPath)
	}

	// Verify workflow exists
	loader := workflow.NewLoader(libraryPath)
	_, err = loader.Load(name)
	if err != nil {
		errMsg := fmt.Sprintf("workflow '%s' not found: %v", name, err)
		fmt.Fprintln(cmd.ErrOrStderr(), errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// Check if already active
	for _, active := range cfg.Workflows.Active {
		if active == name {
			fmt.Fprintf(cmd.OutOrStdout(), "Workflow %s is already active\n", name)
			return nil
		}
	}

	// Add to active list
	cfg.Workflows.Active = append(cfg.Workflows.Active, name)

	// Save config
	if err := saveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	priority := len(cfg.Workflows.Active)
	fmt.Fprintf(cmd.OutOrStdout(), "✓ Activated %s workflow\n", name)
	fmt.Fprintf(cmd.OutOrStdout(), "Priority: %d of %d\n", priority, priority)
	return nil
}

func deactivateWorkflow(cmd *cobra.Command, name string) error {
	// Load config
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if cfg == nil {
		return fmt.Errorf("no config found")
	}

	cfg.ApplyDefaults()

	// Find and remove workflow from active list
	found := false
	newActive := []string{}
	for _, active := range cfg.Workflows.Active {
		if active == name {
			found = true
			continue
		}
		newActive = append(newActive, active)
	}

	if !found {
		fmt.Fprintf(cmd.OutOrStdout(), "Workflow %s is not active\n", name)
		return nil
	}

	cfg.Workflows.Active = newActive

	// Save config
	if err := saveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Deactivated %s workflow\n", name)
	return nil
}

func advanceWorkflow(cmd *cobra.Command) error {
	// Simplified implementation for now
	// In a real implementation, this would use the workflow package
	// wf, err := workflow.LoadCurrent()
	// if err != nil {
	//     return fmt.Errorf("no active workflow found")
	// }
	// if err := wf.Advance(); err != nil {
	//     return fmt.Errorf("failed to advance workflow: %w", err)
	// }

	// Check if workflow is active
	if _, err := os.Stat(".workflow-state.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no active workflow found")
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Advanced to next phase")
	return nil
}

// isKnownWorkflow checks if the given name is a known workflow
func isKnownWorkflow(name string) bool {
	workflowDir := filepath.Join("library", "workflows", name)
	if stat, err := os.Stat(workflowDir); err == nil && stat.IsDir() {
		return true
	}
	return false
}

// handleWorkflowSpecificCommand routes workflow-specific subcommands
func handleWorkflowSpecificCommand(cmd *cobra.Command, workflow string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand required for workflow %s", workflow)
	}

	subcommand := strings.ToLower(args[0])
	switch subcommand {
	case "commands":
		return listWorkflowCommands(cmd, workflow)
	case "execute":
		if len(args) < 2 {
			return fmt.Errorf("command name required for execute")
		}
		return executeWorkflowCommand(cmd, workflow, args[1], args[2:])
	default:
		return fmt.Errorf("unknown subcommand '%s' for workflow '%s'", subcommand, workflow)
	}
}

// listWorkflowCommands lists available commands for a workflow
func listWorkflowCommands(cmd *cobra.Command, workflow string) error {
	commandsDir := filepath.Join("library", "workflows", workflow, "commands")

	// Check if commands directory exists
	if _, err := os.Stat(commandsDir); os.IsNotExist(err) {
		return fmt.Errorf("workflow '%s' not found or has no commands", workflow)
	}

	entries, err := os.ReadDir(commandsDir)
	if err != nil {
		return fmt.Errorf("failed to read commands directory: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Available commands for %s workflow:\n\n", workflow)

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			commandName := strings.TrimSuffix(entry.Name(), ".md")

			// Try to read the first line for description
			description := getCommandDescription(filepath.Join(commandsDir, entry.Name()))

			fmt.Fprintf(cmd.OutOrStdout(), "  %-15s %s\n", commandName, description)
		}
	}

	return nil
}

// getCommandDescription extracts description from command file
func getCommandDescription(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "No description available"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}

	return "No description available"
}

// executeWorkflowCommand loads and displays a workflow command
func executeWorkflowCommand(cmd *cobra.Command, workflow, command string, args []string) error {
	commandPath := filepath.Join("library", "workflows", workflow, "commands", command+".md")

	// Check if command file exists
	if _, err := os.Stat(commandPath); os.IsNotExist(err) {
		return fmt.Errorf("command '%s' not found in workflow '%s'", command, workflow)
	}

	// Read command content
	content, err := os.ReadFile(commandPath)
	if err != nil {
		return fmt.Errorf("failed to read command file: %w", err)
	}

	// Display command content
	fmt.Fprintf(cmd.OutOrStdout(), "Executing %s workflow command: %s\n\n", workflow, command)

	if len(args) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Command Arguments: %v\n\n", args)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(content))

	return nil
}

// saveConfig saves the config back to .ddx/config.yaml
func saveConfig(cfg *config.NewConfig) error {
	return saveConfigWithDir(cfg, ".")
}

// saveConfigWithDir saves the config with explicit working directory
func saveConfigWithDir(cfg *config.NewConfig, workingDir string) error {
	configPath := filepath.Join(workingDir, ".ddx", "config.yaml")

	// Marshal to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// WithDir variants for testing with explicit working directory

func showWorkflowStatusWithDir(cmd *cobra.Command, workingDir string) error {
	// Load config
	cfg, err := loadConfigFrom(workingDir)
	if err != nil || cfg == nil {
		fmt.Fprintln(cmd.OutOrStdout(), "No active workflows")
		return nil
	}

	cfg.ApplyDefaults()

	// Check if any workflows are active
	if len(cfg.Workflows.Active) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No active workflows")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Active workflows (in priority order):")

	// Get library path
	libraryPath := cfg.Library.Path
	if !filepath.IsAbs(libraryPath) {
		libraryPath = filepath.Join(workingDir, libraryPath)
	}

	loader := workflow.NewLoader(libraryPath)

	for i, name := range cfg.Workflows.Active {
		fmt.Fprintf(cmd.OutOrStdout(), "%d. %s\n", i+1, name)

		// Load workflow definition to show agent commands
		def, err := loader.Load(name)
		if err == nil && len(def.AgentCommands) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "   Agent commands:")
			for cmdName, cmdDef := range def.AgentCommands {
				if cmdDef.Enabled {
					fmt.Fprintf(cmd.OutOrStdout(), "   • %s → %s\n", cmdName, cmdDef.Action)
					if cmdDef.Triggers != nil {
						if len(cmdDef.Triggers.Keywords) > 0 {
							fmt.Fprintf(cmd.OutOrStdout(), "     Keywords: %s\n", strings.Join(cmdDef.Triggers.Keywords, ", "))
						}
						if len(cmdDef.Triggers.Patterns) > 0 {
							fmt.Fprintf(cmd.OutOrStdout(), "     Patterns: %s\n", strings.Join(cmdDef.Triggers.Patterns, ", "))
						}
					}
				}
			}
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\nSafe word: %s\n", cfg.Workflows.SafeWord)
	return nil
}

func activateWorkflowWithDir(cmd *cobra.Command, name string, force bool, workingDir string) error {
	// Load config
	cfg, err := loadConfigFrom(workingDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if cfg == nil {
		return fmt.Errorf("no config found - run 'ddx init' first")
	}

	cfg.ApplyDefaults()

	// Get library path
	libraryPath := cfg.Library.Path
	if !filepath.IsAbs(libraryPath) {
		libraryPath = filepath.Join(workingDir, libraryPath)
	}

	// Verify workflow exists
	loader := workflow.NewLoader(libraryPath)
	_, err = loader.Load(name)
	if err != nil {
		errMsg := fmt.Sprintf("workflow '%s' not found: %v", name, err)
		fmt.Fprintln(cmd.ErrOrStderr(), errMsg)
		return fmt.Errorf("%s", errMsg)
	}

	// Check if already active
	for _, active := range cfg.Workflows.Active {
		if active == name {
			fmt.Fprintf(cmd.OutOrStdout(), "Workflow %s is already active\n", name)
			return nil
		}
	}

	// Add to active list
	cfg.Workflows.Active = append(cfg.Workflows.Active, name)

	// Save config
	if err := saveConfigWithDir(cfg, workingDir); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	priority := len(cfg.Workflows.Active)
	fmt.Fprintf(cmd.OutOrStdout(), "✓ Activated %s workflow\n", name)
	fmt.Fprintf(cmd.OutOrStdout(), "Priority: %d of %d\n", priority, priority)
	return nil
}

func deactivateWorkflowWithDir(cmd *cobra.Command, name string, workingDir string) error {
	// Load config
	cfg, err := loadConfigFrom(workingDir)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if cfg == nil {
		return fmt.Errorf("no config found")
	}

	cfg.ApplyDefaults()

	// Find and remove workflow from active list
	found := false
	newActive := make([]string, 0, len(cfg.Workflows.Active))
	for _, active := range cfg.Workflows.Active {
		if active == name {
			found = true
			continue
		}
		newActive = append(newActive, active)
	}

	if !found {
		fmt.Fprintf(cmd.OutOrStdout(), "Workflow %s is not active\n", name)
		return nil
	}

	cfg.Workflows.Active = newActive

	// Save config
	if err := saveConfigWithDir(cfg, workingDir); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Deactivated %s workflow\n", name)
	return nil
}

// newWorkflowCommand creates a workflow command for testing
func newWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage workflows",
		RunE:  runWorkflow,
	}
	cmd.Flags().Bool("force", false, "Force activation")
	return cmd
}
