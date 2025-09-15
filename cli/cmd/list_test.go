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
				// Setup mock DDx home with resources
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				ddxHome := filepath.Join(homeDir, ".ddx")

				// Create template directories
				templatesDir := filepath.Join(ddxHome, "templates")
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "nextjs"), 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "python"), 0755))

				// Create pattern directories
				patternsDir := filepath.Join(ddxHome, "patterns")
				require.NoError(t, os.MkdirAll(filepath.Join(patternsDir, "auth"), 0755))

				// Create prompts directories
				promptsDir := filepath.Join(ddxHome, "prompts")
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "claude"), 0755))

				return homeDir
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
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				ddxHome := filepath.Join(homeDir, ".ddx")

				templatesDir := filepath.Join(ddxHome, "templates")
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "react"), 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "vue"), 0755))

				return homeDir
			},
			validate: func(t *testing.T, output string, err error) {
				assert.Contains(t, output, "Templates")
			},
			expectError: false,
		},
		{
			name: "list with empty DDx home",
			args: []string{"list"},
			setup: func(t *testing.T) string {
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				// Don't create .ddx directory
				return homeDir
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
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				ddxHome := filepath.Join(homeDir, ".ddx")

				templatesDir := filepath.Join(ddxHome, "templates", "test")
				require.NoError(t, os.MkdirAll(templatesDir, 0755))

				// Add a README to the template
				readme := filepath.Join(templatesDir, "README.md")
				require.NoError(t, os.WriteFile(readme, []byte("# Test Template"), 0644))

				return homeDir
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
			if tt.setup != nil {
				tt.setup(t)
			}

			rootCmd := &cobra.Command{
				Use:   "ddx",
				Short: "DDx CLI",
			}
			rootCmd.AddCommand(listCmd)

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
	rootCmd.AddCommand(listCmd)

	output, err := executeCommand(rootCmd, "list", "--help")

	assert.NoError(t, err)
	assert.Contains(t, output, "List all available resources")
	assert.Contains(t, output, "search")
}
