package cmd

import (
	"fmt"
	"os"
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
	// Get flag values locally
	// workflowName is not currently used but may be needed for future features
	// workflowName, _ := cmd.Flags().GetString("name")
	workflowForce, _ := cmd.Flags().GetBool("force")

	// For now, just show workflow status
	if len(args) == 0 {
		return showWorkflowStatus(cmd)
	}

	subcommand := args[0]
	switch strings.ToLower(subcommand) {
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
		return fmt.Errorf("unknown subcommand: %s", subcommand)
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
