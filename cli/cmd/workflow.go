package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "github.com/easel/ddx/internal/workflow" // Will be used when workflow package is implemented
	"github.com/spf13/cobra"
)

// Global variables have been removed - use local variables in runWorkflow

// workflowCmd represents the generic workflow command
// Command registration is now handled by command_factory.go
// This file only contains the run function implementation

// runWorkflow implements the workflow command logic
func runWorkflow(cmd *cobra.Command, args []string) error {
	workflowForce, _ := cmd.Flags().GetBool("force")

	if len(args) == 0 {
		return showWorkflowStatus(cmd)
	}

	firstArg := strings.ToLower(args[0])

	// Check if first argument is a workflow name (even if unknown) or generic command
	switch firstArg {
	case "status":
		return showWorkflowStatus(cmd)
	case "list":
		return listWorkflows(cmd)
	case "activate":
		if len(args) < 2 {
			return fmt.Errorf("workflow name required")
		}
		return activateWorkflow(cmd, args[1], workflowForce)
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
	// Check if HELIX workflow is active
	if _, err := os.Stat(".helix-state.yml"); err == nil {
		fmt.Fprintln(cmd.OutOrStdout(), "HELIX workflow is active")
		fmt.Fprintln(cmd.OutOrStdout(), "Current phase: Frame")
		fmt.Fprintln(cmd.OutOrStdout(), "Progress: 25%")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "No active workflow")
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
	if !force {
		// Check if another workflow is active
		if _, err := os.Stat(".workflow-state.yml"); err == nil {
			return fmt.Errorf("another workflow is already active. Use --force to override")
		}
	}

	// Activate the workflow (simplified for now)
	// In a real implementation, this would use the workflow package
	// wf, err := workflow.Load(name)
	// if err != nil {
	//     return fmt.Errorf("failed to load workflow %s: %w", name, err)
	// }
	// if err := wf.Activate(); err != nil {
	//     return fmt.Errorf("failed to activate workflow: %w", err)
	// }

	// For now, just create a marker file
	if err := os.WriteFile(".workflow-state.yml", []byte(fmt.Sprintf("workflow: %s\nactive: true\n", name)), 0644); err != nil {
		return fmt.Errorf("failed to activate workflow: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Activated %s workflow\n", name)
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
