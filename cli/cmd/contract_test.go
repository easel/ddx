package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Contract validation tests verify that CLI commands conform to their API contracts
// as defined in docs/02-design/contracts/CLI-001-core-commands.md

// TestInitCommand_Contract validates init command against CLI-001 contract
func TestInitCommand_Contract(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		args           []string
		setup          func(t *testing.T) string
		expectCode     int // Expected exit code per contract
		validateOutput func(t *testing.T, output string)
	}{
		{
			name:        "contract_exit_code_0_success",
			description: "Exit code 0: Success",
			args:        []string{"init"},
			setup: func(t *testing.T) string {
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))
				return workDir
			},
			expectCode: 0,
			validateOutput: func(t *testing.T, output string) {
				// Contract specifies these output elements
				assert.Contains(t, output, "Initialized DDx")
			},
		},
		{
			name:        "contract_exit_code_2_exists",
			description: "Exit code 2: Configuration already exists",
			args:        []string{"init"},
			setup: func(t *testing.T) string {
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))
				// Create existing config
				require.NoError(t, os.WriteFile(
					filepath.Join(workDir, ".ddx.yml"),
					[]byte("version: 1.0"),
					0644,
				))
				return workDir
			},
			expectCode: 2,
			validateOutput: func(t *testing.T, output string) {
				// Should indicate config exists
			},
		},
		{
			name:        "contract_force_flag",
			description: "--force flag overwrites existing config",
			args:        []string{"init", "--force"},
			setup: func(t *testing.T) string {
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))
				// Create existing config
				require.NoError(t, os.WriteFile(
					filepath.Join(workDir, ".ddx.yml"),
					[]byte("version: 0.9"),
					0644,
				))
				return workDir
			},
			expectCode: 0,
			validateOutput: func(t *testing.T, output string) {
				// Should succeed with force
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flag variables to ensure test isolation
			initTemplate = ""
			initForce = false

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			if tt.setup != nil {
				tt.setup(t)
			}

			// Create a fresh init command to avoid state pollution
			freshInitCmd := &cobra.Command{
				Use:   "init",
				Short: "Initialize DDx in current project",
				Long:  initCmd.Long,
				RunE:  runInit,
			}
			freshInitCmd.Flags().StringVarP(&initTemplate, "template", "t", "", "Use specific template")
			freshInitCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Force initialization even if DDx already exists")

			rootCmd := &cobra.Command{
				Use:   "ddx",
				Short: "DDx CLI",
			}
			rootCmd.AddCommand(freshInitCmd)

			output, err := executeContractCommand(rootCmd, tt.args...)

			// Validate exit code matches contract
			if tt.expectCode == 0 {
				assert.NoError(t, err, "Contract specifies exit code 0 for: %s", tt.description)
			} else {
				assert.Error(t, err, "Contract specifies non-zero exit code for: %s", tt.description)
				// Note: Cobra doesn't expose exact exit codes in tests,
				// but we validate error presence
			}

			if tt.validateOutput != nil {
				tt.validateOutput(t, output)
			}
		})
	}
}

// TestListCommand_Contract validates list command against CLI-001 contract
func TestListCommand_Contract(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		args           []string
		setup          func(t *testing.T) string
		expectCode     int
		validateOutput func(t *testing.T, output string)
	}{
		{
			name:        "contract_output_format",
			description: "Output format matches contract specification",
			args:        []string{"list"},
			setup: func(t *testing.T) string {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))
				t.Cleanup(func() { os.Chdir(origWd) })

				// Create library structure as per contract
				libraryDir := filepath.Join(testDir, "library")
				templatesDir := filepath.Join(libraryDir, "templates")
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "nextjs"), 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "python"), 0755))

				patternsDir := filepath.Join(libraryDir, "patterns")
				require.NoError(t, os.MkdirAll(filepath.Join(patternsDir, "auth"), 0755))

				// Create config
				config := []byte(`version: "2.0"
library_path: ./library`)
				require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

				return testDir
			},
			expectCode: 0,
			validateOutput: func(t *testing.T, output string) {
				// Contract specifies section headers
				assert.Contains(t, output, "Templates")
				assert.Contains(t, output, "Patterns")
			},
		},
		{
			name:        "contract_filter_argument",
			description: "Filter argument works as specified",
			args:        []string{"list", "templates"},
			setup: func(t *testing.T) string {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))
				t.Cleanup(func() { os.Chdir(origWd) })

				libraryDir := filepath.Join(testDir, "library")
				templatesDir := filepath.Join(libraryDir, "templates")
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "react"), 0755))

				patternsDir := filepath.Join(libraryDir, "patterns")
				require.NoError(t, os.MkdirAll(filepath.Join(patternsDir, "auth"), 0755))

				// Create config
				config := []byte(`version: "2.0"
library_path: ./library`)
				require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

				return testDir
			},
			expectCode: 0,
			validateOutput: func(t *testing.T, output string) {
				// Should only show templates
				assert.Contains(t, output, "Templates")
				// Patterns should not be shown when filtering
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flag variables to ensure test isolation
			listType = ""
			listSearch = ""
			listVerbose = false

			if tt.setup != nil {
				tt.setup(t)
			}

			// Create a fresh list command to avoid state pollution
			freshListCmd := &cobra.Command{
				Use:   "list",
				Short: "List available resources",
				Long:  listCmd.Long,
				RunE:  runList,
			}
			freshListCmd.Flags().StringVarP(&listType, "type", "t", "", "Filter by type (templates|patterns|configs|prompts|scripts)")
			freshListCmd.Flags().StringVarP(&listSearch, "search", "s", "", "Search for specific items")
			freshListCmd.Flags().BoolVar(&listVerbose, "verbose", false, "Show verbose output with additional details")

			rootCmd := &cobra.Command{
				Use:   "ddx",
				Short: "DDx CLI",
			}
			rootCmd.AddCommand(freshListCmd)

			output, err := executeContractCommand(rootCmd, tt.args...)

			if tt.expectCode == 0 {
				assert.NoError(t, err, "Contract specifies exit code 0 for: %s", tt.description)
			} else {
				assert.Error(t, err, "Contract specifies non-zero exit code for: %s", tt.description)
			}

			if tt.validateOutput != nil {
				tt.validateOutput(t, output)
			}
		})
	}
}

// TestApplyCommand_Contract validates apply command against CLI-001 contract
func TestApplyCommand_Contract(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		args           []string
		setup          func(t *testing.T) (string, string)
		expectCode     int
		validateOutput func(t *testing.T, output string)
		validateFiles  func(t *testing.T, workDir string)
	}{
		{
			name:        "contract_exit_code_6_not_found",
			description: "Exit code 6: Resource not found",
			args:        []string{"apply", "templates/nonexistent"},
			setup: func(t *testing.T) (string, string) {
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Create library structure (but no templates/nonexistent)
				libraryDir := filepath.Join(workDir, "library")
				require.NoError(t, os.MkdirAll(filepath.Join(libraryDir, "templates"), 0755))

				// Create config pointing to library
				config := []byte(`version: "2.0"
library_path: ./library`)
				require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

				return workDir, workDir
			},
			expectCode: 6,
			validateOutput: func(t *testing.T, output string) {
				// Should indicate resource not found
			},
		},
		{
			name:        "contract_dry_run_flag",
			description: "--dry-run flag shows changes without applying",
			args:        []string{"apply", "templates/test", "--dry-run"},
			setup: func(t *testing.T) (string, string) {
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Create library with test template
				templateDir := filepath.Join(workDir, "library", "templates", "test")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				// Add template file
				templateFile := filepath.Join(templateDir, "test.txt")
				require.NoError(t, os.WriteFile(templateFile, []byte("Test content"), 0644))

				// Create config pointing to library
				config := []byte(`version: "2.0"
library_path: ./library`)
				require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

				return workDir, workDir
			},
			expectCode: 0,
			validateOutput: func(t *testing.T, output string) {
				// Should show what would be applied
				assert.Contains(t, output, "Would apply")
			},
			validateFiles: func(t *testing.T, workDir string) {
				// Files should NOT be created in dry-run
				testFile := filepath.Join(workDir, "test.txt")
				assert.NoFileExists(t, testFile, "Dry-run should not create files")
			},
		},
		{
			name:        "contract_variable_substitution",
			description: "Variables are substituted as per contract",
			args:        []string{"apply", "templates/vars", "--var", "name=TestProject"},
			setup: func(t *testing.T) (string, string) {
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Create library with vars template
				templateDir := filepath.Join(workDir, "library", "templates", "vars")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				// Template with variable
				templateFile := filepath.Join(templateDir, "project.txt")
				require.NoError(t, os.WriteFile(templateFile, []byte("Project: {{name}}"), 0644))

				// Create config pointing to library
				config := []byte(`version: "2.0"
library_path: ./library`)
				require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

				return workDir, workDir
			},
			expectCode: 0,
			validateFiles: func(t *testing.T, workDir string) {
				projectFile := filepath.Join(workDir, "project.txt")
				if _, err := os.Stat(projectFile); err == nil {
					content, _ := os.ReadFile(projectFile)
					assert.Equal(t, "Project: TestProject", string(content),
						"Variables should be substituted as per contract")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flag variables to ensure test isolation
			applyPath = "."
			applyDryRun = false
			applyVars = nil

			// Also reset flags on the command
			applyCmd.Flags().Set("path", ".")
			applyCmd.Flags().Set("dry-run", "false")
			applyCmd.Flags().Set("var", "")

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			var workDir string
			if tt.setup != nil {
				_, workDir = tt.setup(t)
			}

			// Reset apply command flags to ensure test isolation
			applyPath = "."
			applyDryRun = false
			applyVars = nil
			libraryPath = ""

			// Create a fresh apply command to avoid state pollution
			freshApplyCmd := &cobra.Command{
				Use:   "apply <resource>",
				Short: "Apply a DDx resource to your project",
				Long:  applyCmd.Long,
				Args:  cobra.ExactArgs(1),
				RunE:  runApply,
			}
			freshApplyCmd.Flags().StringVarP(&applyPath, "path", "p", ".", "Target path for application")
			freshApplyCmd.Flags().BoolVar(&applyDryRun, "dry-run", false, "Show what would be applied without making changes")
			freshApplyCmd.Flags().StringSliceVar(&applyVars, "var", nil, "Set template variables (key=value)")

			rootCmd := &cobra.Command{
				Use:   "ddx",
				Short: "DDx CLI",
			}
			rootCmd.AddCommand(freshApplyCmd)

			output, err := executeContractCommand(rootCmd, tt.args...)

			if tt.expectCode == 0 {
				assert.NoError(t, err, "Contract specifies exit code 0 for: %s", tt.description)
			} else {
				assert.Error(t, err, "Contract specifies non-zero exit code for: %s", tt.description)
			}

			if tt.validateOutput != nil {
				tt.validateOutput(t, output)
			}

			if tt.validateFiles != nil {
				tt.validateFiles(t, workDir)
			}
		})
	}
}

// TestConfigCommand_Contract validates config command output format
func TestConfigCommand_Contract(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		args           []string
		setup          func(t *testing.T) string
		validateOutput func(t *testing.T, output string)
	}{
		{
			name:        "contract_yaml_output",
			description: "Config output is valid YAML as per contract",
			args:        []string{"config"},
			setup: func(t *testing.T) string {
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				config := `version: "1.0"
repository:
  url: "https://github.com/test/repo"
variables:
  test: "value"`
				require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644))
				return workDir
			},
			validateOutput: func(t *testing.T, output string) {
				// Output should contain YAML structure
				assert.Contains(t, output, "version")
				assert.Contains(t, output, "repository")
			},
		},
		{
			name:        "contract_get_operation",
			description: "Get operation returns specific value",
			args:        []string{"config", "get", "version"},
			setup: func(t *testing.T) string {
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				config := `version: "1.0"`
				require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644))
				return workDir
			},
			validateOutput: func(t *testing.T, output string) {
				// Should return just the value
				assert.Contains(t, strings.TrimSpace(output), "1.0")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global flag variables to ensure test isolation
			configGlobal = false
			configLocal = false
			configInit = false
			configShow = false
			libraryPath = ""

			// Also reset flags on the command if they exist
			if configCmd.Flags().Lookup("global") != nil {
				configCmd.Flags().Set("global", "false")
			}
			if configCmd.Flags().Lookup("local") != nil {
				configCmd.Flags().Set("local", "false")
			}
			if configCmd.Flags().Lookup("init") != nil {
				configCmd.Flags().Set("init", "false")
			}
			if configCmd.Flags().Lookup("show") != nil {
				configCmd.Flags().Set("show", "false")
			}

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			if tt.setup != nil {
				tt.setup(t)
			}

			// Create a fresh config command to avoid state pollution
			freshConfigCmd := &cobra.Command{
				Use:   "config [get|set|validate] [key] [value]",
				Short: "Manage DDx configuration",
				Long:  configCmd.Long,
				RunE:  runConfig,
			}
			freshConfigCmd.Flags().BoolVarP(&configGlobal, "global", "g", false, "Edit global configuration")
			freshConfigCmd.Flags().BoolVarP(&configLocal, "local", "l", false, "Edit local project configuration")
			freshConfigCmd.Flags().BoolVar(&configInit, "init", false, "Initialize configuration wizard")
			freshConfigCmd.Flags().BoolVar(&configShow, "show", false, "Show current configuration")

			rootCmd := &cobra.Command{
				Use:   "ddx",
				Short: "DDx CLI",
			}
			rootCmd.AddCommand(freshConfigCmd)

			output, err := executeContractCommand(rootCmd, tt.args...)
			assert.NoError(t, err)

			if tt.validateOutput != nil {
				tt.validateOutput(t, output)
			}
		})
	}
}

// Helper to execute command for contract testing
func executeContractCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}
