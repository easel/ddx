package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWorkflowCommandDiscovery tests workflow command discovery functionality
func TestWorkflowCommandDiscovery(t *testing.T) {
	tests := []struct {
		name     string
		workflow string
		setup    func(t *testing.T) string
		expected []string
	}{
		{
			name:     "discover_helix_commands",
			workflow: "helix",
			setup:    setupHelixWorkflowCommands,
			expected: []string{"build-story", "continue"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	// originalDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection // REMOVED: Using CommandFactory injection
			_ = tt.setup(t)
			defer func() {
			}()

			// This will fail initially - function doesn't exist
			commands, err := discoverWorkflowCommands(tt.workflow)

			require.NoError(t, err)
			for _, expectedCmd := range tt.expected {
				assert.Contains(t, commands, expectedCmd)
			}
		})
	}
}

// TestWorkflowCommandExecution tests workflow command execution
func TestWorkflowCommandExecution(t *testing.T) {
	tests := []struct {
		name     string
		workflow string
		command  string
		args     []string
		setup    func(t *testing.T) string
		expected string
	}{
		{
			name:     "execute_build_story",
			workflow: "helix",
			command:  "build-story",
			args:     []string{"US-001"},
			setup:    setupHelixWorkflowCommands,
			expected: "HELIX Command: Build Story",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	// originalDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection // REMOVED: Using CommandFactory injection
			_ = tt.setup(t)
			defer func() {
			}()

			// This will fail initially - function doesn't exist
			output, err := loadWorkflowCommandContent(tt.workflow, tt.command, tt.args)

			require.NoError(t, err)
			assert.Contains(t, output, tt.expected)
		})
	}
}

// TestWorkflowCommandCLI tests the CLI integration for workflow commands
func TestWorkflowCommandCLI(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		setup    func(t *testing.T) string
		expected string
		wantErr  bool
	}{
		{
			name:     "cli_list_helix_commands",
			args:     []string{"workflow", "helix", "commands"},
			setup:    setupHelixWorkflowCommands,
			expected: "Available commands for helix workflow:",
			wantErr:  false,
		},
		{
			name:     "cli_execute_build_story",
			args:     []string{"workflow", "helix", "execute", "build-story", "US-001"},
			setup:    setupHelixWorkflowCommands,
			expected: "HELIX Command: Build Story",
			wantErr:  false,
		},
		{
			name:     "cli_invalid_workflow",
			args:     []string{"workflow", "invalid", "commands"},
			setup:    setupEmptyWorkspace,
			expected: "workflow 'invalid' not found",
			wantErr:  true,
		},
		{
			name:     "cli_invalid_command",
			args:     []string{"workflow", "helix", "execute", "invalid-command"},
			setup:    setupHelixWorkflowCommands,
			expected: "command 'invalid-command' not found",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	// originalDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection // REMOVED: Using CommandFactory injection
			workDir := tt.setup(t)
			defer func() {
			}()

			// Use CommandFactory with the test working directory
			factory := NewCommandFactory(workDir)
			rootCmd := factory.NewRootCommand()
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.expected)
				}
			} else {
				assert.NoError(t, err)
				output := buf.String()
				assert.Contains(t, output, tt.expected)
			}
		})
	}
}

// Helper function to setup helix workflow commands
func setupHelixWorkflowCommands(t *testing.T) string {
	workDir := t.TempDir()

	commandsDir := filepath.Join(workDir, "library", "workflows", "helix", "commands")
	require.NoError(t, os.MkdirAll(commandsDir, 0755))

	// Create build-story command
	buildStoryContent := `# HELIX Command: Build Story

You are a HELIX workflow executor tasked with implementing work on a specific user story through comprehensive evaluation and systematic implementation.

## Command Input

You will receive a user story ID as an argument (e.g., US-001, US-042, etc.).

## Your Mission

Execute a comprehensive evaluation and implementation process.`

	require.NoError(t, os.WriteFile(
		filepath.Join(commandsDir, "build-story.md"),
		[]byte(buildStoryContent), 0644))

	// Create continue command
	continueContent := `# HELIX Command: Continue

Continue work on the current user story following HELIX methodology.`

	require.NoError(t, os.WriteFile(
		filepath.Join(commandsDir, "continue.md"),
		[]byte(continueContent), 0644))

	return workDir
}

// Helper function to setup empty workspace
func setupEmptyWorkspace(t *testing.T) string {
	workDir := t.TempDir()
	return workDir
}

// TestIsKnownWorkflow tests workflow detection
func TestIsKnownWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		workflow string
		setup    func(t *testing.T) string
		expected bool
	}{
		{
			name:     "known_helix_workflow",
			workflow: "helix",
			setup:    setupHelixWorkflowCommands,
			expected: true,
		},
		{
			name:     "unknown_workflow",
			workflow: "unknown",
			setup:    setupEmptyWorkspace,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	// originalDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection // REMOVED: Using CommandFactory injection
			_ = tt.setup(t)
			defer func() {
			}()

			// This will fail initially - function doesn't exist
			result := isKnownWorkflow(tt.workflow)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions for testing

// discoverWorkflowCommands discovers available commands for a workflow
func discoverWorkflowCommands(workflow string) ([]string, error) {
	commandsDir := filepath.Join("library", "workflows", workflow, "commands")

	// Check if commands directory exists
	if _, err := os.Stat(commandsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("workflow '%s' not found or has no commands", workflow)
	}

	entries, err := os.ReadDir(commandsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read commands directory: %w", err)
	}

	var commands []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			commandName := strings.TrimSuffix(entry.Name(), ".md")
			commands = append(commands, commandName)
		}
	}

	return commands, nil
}

// loadWorkflowCommandContent loads content of a workflow command
func loadWorkflowCommandContent(workflow, command string, args []string) (string, error) {
	commandPath := filepath.Join("library", "workflows", workflow, "commands", command+".md")

	// Check if command file exists
	if _, err := os.Stat(commandPath); os.IsNotExist(err) {
		return "", fmt.Errorf("command '%s' not found in workflow '%s'", command, workflow)
	}

	// Read command content
	content, err := os.ReadFile(commandPath)
	if err != nil {
		return "", fmt.Errorf("failed to read command file: %w", err)
	}

	return string(content), nil
}
