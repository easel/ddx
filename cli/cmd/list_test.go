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
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))
				t.Cleanup(func() { os.Chdir(origWd) })

				// Create library structure
				libraryDir := filepath.Join(testDir, "library")

				// Create template directories
				templatesDir := filepath.Join(libraryDir, "templates")
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "nextjs"), 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "python"), 0755))

				// Create pattern directories
				patternsDir := filepath.Join(libraryDir, "patterns")
				require.NoError(t, os.MkdirAll(filepath.Join(patternsDir, "auth"), 0755))

				// Create prompts directories
				promptsDir := filepath.Join(libraryDir, "prompts")
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "claude"), 0755))

				// Create .ddx.yml config pointing to library
				config := []byte(`version: "2.0"
library_path: ./library
repository:
  url: https://github.com/easel/ddx
  branch: main`)
				require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

				return testDir
			},
			validate: func(t *testing.T, output string, err error) {
				assert.Contains(t, output, "Templates")
				// Note: Actual output format depends on implementation
			},
			expectError: false,
		},
		{
			name: "list specific resource type",
			args: []string{"list", "templates"},
			setup: func(t *testing.T) string {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))
				t.Cleanup(func() { os.Chdir(origWd) })

				libraryDir := filepath.Join(testDir, "library")
				templatesDir := filepath.Join(libraryDir, "templates")
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "react"), 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "vue"), 0755))

				// Create .ddx.yml config
				config := []byte(`version: "2.0"
library_path: ./library`)
				require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

				return testDir
			},
			validate: func(t *testing.T, output string, err error) {
				assert.Contains(t, output, "Templates")
			},
			expectError: false,
		},
		{
			name: "list with no library",
			args: []string{"list"},
			setup: func(t *testing.T) string {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))
				t.Cleanup(func() { os.Chdir(origWd) })
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
			name: "list with verbose flag",
			args: []string{"list", "--verbose"},
			setup: func(t *testing.T) string {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))
				t.Cleanup(func() { os.Chdir(origWd) })

				libraryDir := filepath.Join(testDir, "library")
				templatesDir := filepath.Join(libraryDir, "templates", "test")
				require.NoError(t, os.MkdirAll(templatesDir, 0755))

				// Add a README to the template
				readme := filepath.Join(templatesDir, "README.md")
				require.NoError(t, os.WriteFile(readme, []byte("# Test Template"), 0644))

				// Create config
				config := []byte(`version: "2.0"
library_path: ./library`)
				require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

				return testDir
			},
			validate: func(t *testing.T, output string, err error) {
				// Verbose output should include more details
				assert.NotEmpty(t, output)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh command for test isolation
			// (flags are now local to the command)

			if tt.setup != nil {
				tt.setup(t)
			}

			rootCmd := &cobra.Command{
				Use:   "ddx",
				Short: "DDx CLI",
			}

			// Create fresh list command to avoid state pollution
			freshListCmd := &cobra.Command{
				Use:   "list",
				Short: "List available templates, patterns, and configurations",
				Long: `List all available resources in the DDx toolkit.

You can filter by type or search for specific items.`,
				RunE: runList,
			}
			freshListCmd.Flags().StringP("type", "t", "", "Filter by type (templates|patterns|configs|prompts|scripts)")
			freshListCmd.Flags().StringP("search", "s", "", "Search for specific items")
			freshListCmd.Flags().Bool("verbose", false, "Show verbose output with additional details")

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
	freshListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates, patterns, and configurations",
		RunE:  runList,
	}
	freshListCmd.Flags().StringP("type", "t", "", "Filter by type (templates|patterns|configs|prompts|scripts)")
	freshListCmd.Flags().StringP("search", "s", "", "Search for specific items")
	freshListCmd.Flags().Bool("verbose", false, "Show verbose output with additional details")

	rootCmd.AddCommand(freshListCmd)

	output, err := executeCommand(rootCmd, "list", "--help")

	assert.NoError(t, err)
	assert.Contains(t, output, "List available templates")
	assert.Contains(t, output, "search")
}
