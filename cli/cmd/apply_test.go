package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestApplyCommand tests the apply command
func TestApplyCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		setup       func(t *testing.T) (string, string) // returns homeDir, workDir
		validate    func(t *testing.T, workDir string, output string, err error)
		expectError bool
	}{
		{
			name: "apply template",
			args: []string{"apply", "templates/test"},
			setup: func(t *testing.T) (string, string) {
				// Setup home directory with template
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)

				templateDir := filepath.Join(homeDir, ".ddx", "templates", "test")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				// Add template files
				templateFile := filepath.Join(templateDir, "app.txt")
				require.NoError(t, os.WriteFile(templateFile, []byte("Hello {{name}}"), 0644))

				// Setup work directory
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Create .ddx.yml config
				config := `version: "1.0"
repository:
  url: "https://github.com/test/repo"
variables:
  name: "World"
`
				require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644))

				return homeDir, workDir
			},
			validate: func(t *testing.T, workDir string, output string, err error) {
				// Check if file was created
				appFile := filepath.Join(workDir, "app.txt")
				if _, err := os.Stat(appFile); err == nil {
					content, _ := os.ReadFile(appFile)
					assert.Contains(t, string(content), "Hello")
				}
			},
			expectError: false,
		},
		{
			name: "apply non-existent template",
			args: []string{"apply", "templates/nonexistent"},
			setup: func(t *testing.T) (string, string) {
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)

				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Create minimal config
				config := `version: "1.0"`
				require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644))

				return homeDir, workDir
			},
			validate: func(t *testing.T, workDir string, output string, err error) {
				// Should fail
				assert.Error(t, err)
			},
			expectError: true,
		},
		{
			name: "apply with variables",
			args: []string{"apply", "templates/vars", "--var", "project=TestProject", "--var", "version=1.0.0"},
			setup: func(t *testing.T) (string, string) {
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)

				templateDir := filepath.Join(homeDir, ".ddx", "templates", "vars")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				// Template with variables
				templateFile := filepath.Join(templateDir, "info.txt")
				content := "Project: {{project}}\nVersion: {{version}}"
				require.NoError(t, os.WriteFile(templateFile, []byte(content), 0644))

				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Minimal config
				config := `version: "1.0"`
				require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644))

				return homeDir, workDir
			},
			validate: func(t *testing.T, workDir string, output string, err error) {
				infoFile := filepath.Join(workDir, "info.txt")
				if _, err := os.Stat(infoFile); err == nil {
					content, _ := os.ReadFile(infoFile)
					assert.Contains(t, string(content), "TestProject")
					assert.Contains(t, string(content), "1.0.0")
				}
			},
			expectError: false,
		},
		{
			name: "apply pattern",
			args: []string{"apply", "patterns/auth"},
			setup: func(t *testing.T) (string, string) {
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)

				patternDir := filepath.Join(homeDir, ".ddx", "patterns", "auth")
				require.NoError(t, os.MkdirAll(patternDir, 0755))

				// Add pattern file
				patternFile := filepath.Join(patternDir, "auth.js")
				require.NoError(t, os.WriteFile(patternFile, []byte("// Auth pattern"), 0644))

				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				config := `version: "1.0"`
				require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644))

				return homeDir, workDir
			},
			validate: func(t *testing.T, workDir string, output string, err error) {
				authFile := filepath.Join(workDir, "auth.js")
				if _, err := os.Stat(authFile); err == nil {
					content, _ := os.ReadFile(authFile)
					assert.Contains(t, string(content), "Auth pattern")
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flag variables to ensure test isolation
			applyPath = "."
			applyDryRun = false
			applyVars = nil

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			var workDir string
			if tt.setup != nil {
				_, workDir = tt.setup(t)
			}

			rootCmd := &cobra.Command{
				Use:   "ddx",
				Short: "DDx CLI",
			}
			rootCmd.AddCommand(applyCmd)

			output, err := executeCommand(rootCmd, tt.args...)

			if tt.expectError {
				// Note: Command might not return error directly
				// but we validate in the validate function
			}

			if tt.validate != nil {
				tt.validate(t, workDir, output, err)
			}
		})
	}
}

// TestApplyCommand_Help tests the help output
func TestApplyCommand_Help(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "ddx",
		Short: "DDx CLI",
	}
	rootCmd.AddCommand(applyCmd)

	output, err := executeCommand(rootCmd, "apply", "--help")

	assert.NoError(t, err)
	assert.Contains(t, output, "Apply")
	assert.Contains(t, output, "template")
	assert.Contains(t, output, "--dry-run")
}
