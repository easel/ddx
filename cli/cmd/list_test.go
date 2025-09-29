package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListCommand tests the list command
func TestListCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		setup       func(t *testing.T) string
		validate    func(t *testing.T, output string, err error)
		expectError bool
	}{
		{
			name: "list all resources",
			args: []string{"list"},
			setup: func(t *testing.T) string {
				// Setup test directory with library
				testDir := t.TempDir()

				// Create library structure in .ddx/library
				libraryDir := filepath.Join(testDir, ".ddx", "library")

				// Create workflow directories
				workflowsDir := filepath.Join(libraryDir, "workflows")
				require.NoError(t, os.MkdirAll(filepath.Join(workflowsDir, "helix"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(workflowsDir, "helix", "workflow.yml"), []byte("name: helix"), 0644))

				// Create mcp-servers directories
				mcpDir := filepath.Join(libraryDir, "mcp-servers")
				require.NoError(t, os.MkdirAll(filepath.Join(mcpDir, "github"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(mcpDir, "github", "server.json"), []byte("{}"), 0644))

				// Create prompts directories
				promptsDir := filepath.Join(libraryDir, "prompts")
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "claude"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(promptsDir, "claude", "prompt.md"), []byte("# Prompt"), 0644))

				// Create .ddx/config.yaml config pointing to library
				config := []byte(`version: "1.0"
library:
  path: .ddx/library
  repository:
    url: https://github.com/easel/ddx-library
    branch: main`)
				require.NoError(t, os.WriteFile(filepath.Join(testDir, ".ddx", "config.yaml"), config, 0644))

				return testDir
			},
			validate: func(t *testing.T, output string, err error) {
				assert.Contains(t, output, "Workflows")
				// Note: Actual output format depends on implementation
			},
			expectError: false,
		},
		{
			name: "list specific resource type",
			args: []string{"list", "workflows"},
			setup: func(t *testing.T) string {
				testDir := t.TempDir()

				libraryDir := filepath.Join(testDir, ".ddx", "library")
				workflowsDir := filepath.Join(libraryDir, "workflows")
				require.NoError(t, os.MkdirAll(filepath.Join(workflowsDir, "helix"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(workflowsDir, "helix", "workflow.yml"), []byte("name: helix"), 0644))
				require.NoError(t, os.MkdirAll(filepath.Join(workflowsDir, "kanban"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(workflowsDir, "kanban", "workflow.yml"), []byte("name: kanban"), 0644))

				// Create .ddx/config.yaml config
				config := []byte(`version: "1.0"
library:
  path: .ddx/library
  repository:
    url: https://github.com/easel/ddx-library
    branch: main`)
				require.NoError(t, os.WriteFile(filepath.Join(testDir, ".ddx", "config.yaml"), config, 0644))

				return testDir
			},
			validate: func(t *testing.T, output string, err error) {
				assert.Contains(t, output, "Workflows")
			},
			expectError: false,
		},
		{
			name: "list with no library",
			args: []string{"list"},
			setup: func(t *testing.T) string {
				testDir := t.TempDir()
				// Don't create library or config - should fail gracefully
				return testDir
			},
			validate: func(t *testing.T, output string, err error) {
				// Should handle gracefully
				assert.NotNil(t, output)
			},
			expectError: false,
		},
		{
			name: "list with json flag",
			args: []string{"list", "--json"},
			setup: func(t *testing.T) string {
				testDir := t.TempDir()

				libraryDir := filepath.Join(testDir, "library")
				workflowsDir := filepath.Join(libraryDir, "workflows", "test")
				require.NoError(t, os.MkdirAll(workflowsDir, 0755))

				// Add a README to the workflow
				readme := filepath.Join(workflowsDir, "README.md")
				require.NoError(t, os.WriteFile(readme, []byte("# Test Workflow"), 0644))

				// Add actual workflow file
				workflow := filepath.Join(workflowsDir, "workflow.yml")
				require.NoError(t, os.WriteFile(workflow, []byte("name: test"), 0644))

				// Create config
				config := []byte(`version: "2.0"
library_path: ./library`)
				require.NoError(t, os.WriteFile(filepath.Join(testDir, ".ddx.yml"), config, 0644))

				return testDir
			},
			validate: func(t *testing.T, output string, err error) {
				// JSON output should include proper structure
				assert.Contains(t, output, "resources")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh command for test isolation
			// (flags are now local to the command)

			var testDir string
			if tt.setup != nil {
				testDir = tt.setup(t)
			} else {
				testDir = t.TempDir()
			}

			rootCmd := &cobra.Command{
				Use:   "ddx",
				Short: "DDx CLI",
			}

			// Create fresh list command to avoid state pollution
			factory := NewCommandFactory(testDir)
			freshListCmd := &cobra.Command{
				Use:   "list",
				Short: "List available DDx resources",
				Long: `List all available resources in the DDx toolkit.

You can filter by type or search for specific items.`,
				RunE: factory.runList,
			}
			freshListCmd.Flags().StringP("filter", "f", "", "Filter resources by name")
			freshListCmd.Flags().Bool("json", false, "Output results as JSON")
			freshListCmd.Flags().Bool("tree", false, "Display resources in tree format")

			rootCmd.AddCommand(freshListCmd)

			output, err := executeCommand(rootCmd, tt.args...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validate != nil {
				tt.validate(t, output, err)
			}
		})
	}
}

// TestListCommand_Help tests the help output
func TestListCommand_Help(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "ddx",
		Short: "DDx CLI",
	}

	// Create fresh list command
	factory := NewCommandFactory("/tmp")
	freshListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available DDx resources",
		RunE:  factory.runList,
	}
	freshListCmd.Flags().StringP("filter", "f", "", "Filter resources by name")
	freshListCmd.Flags().Bool("json", false, "Output results as JSON")
	freshListCmd.Flags().Bool("tree", false, "Display resources in tree format")

	rootCmd.AddCommand(freshListCmd)

	output, err := executeCommand(rootCmd, "list", "--help")

	assert.NoError(t, err)
	assert.Contains(t, output, "List available")
	assert.Contains(t, output, "filter")
}
