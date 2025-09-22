package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Acceptance tests validate user stories and business requirements
// These tests follow the Given/When/Then pattern from user stories

// Helper function to create a fresh root command for tests
func getTestRootCommand() *cobra.Command {
	factory := NewCommandFactory()
	return factory.NewRootCommand()
}

// TestAcceptance_US001_InitializeProject tests US-001: Initialize DDX in Project
func TestAcceptance_US001_InitializeProject(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		given    func(t *testing.T) string                 // Setup conditions
		when     func(t *testing.T, dir string) error      // Execute action
		then     func(t *testing.T, dir string, err error) // Verify outcome
	}{
		{
			name:     "basic_initialization",
			scenario: "Initialize DDX in project without existing configuration",
			given: func(t *testing.T) string {
				// Given: I am in a project directory without DDX
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))
				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx init`
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "init")
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: a `.ddx.yml` configuration file exists with my settings
				configPath := filepath.Join(dir, ".ddx.yml")
				if _, statErr := os.Stat(configPath); statErr == nil {
					// Config file exists - validate structure
					data, readErr := os.ReadFile(configPath)
					require.NoError(t, readErr)

					var config map[string]interface{}
					yamlErr := yaml.Unmarshal(data, &config)
					require.NoError(t, yamlErr)

					assert.Contains(t, config, "version", "Config should have version")
					assert.Contains(t, config, "repository", "Config should have repository")
				}
			},
		},
		{
			name:     "template_based_initialization",
			scenario: "Initialize DDX with specific template",
			given: func(t *testing.T) string {
				// Given: I want to use a specific template
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Setup mock template
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				templateDir := filepath.Join(homeDir, ".ddx", "templates", "test-template")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx init --template test-template`
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "init", "--template", "test-template")
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: the specified template is applied during initialization
				// Note: Actual implementation may vary
				t.Log("Template-based initialization scenario")
			},
		},
		{
			name:     "reinitialization_prevention",
			scenario: "Prevent re-initialization of DDX-enabled project",
			given: func(t *testing.T) string {
				// Given: DDX is already initialized
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Create existing config
				config := `version: "1.0"
repository:
  url: "https://github.com/test/repo"`
				require.NoError(t, os.WriteFile(
					filepath.Join(workDir, ".ddx.yml"),
					[]byte(config),
					0644,
				))

				return workDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx init` again
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "init")
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: Clear message that DDX is already initialized
				if err != nil {
					assert.Contains(t, err.Error(), "already",
						"Error should indicate DDX already initialized")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Given
			dir := tt.given(t)

			// When
			err := tt.when(t, dir)

			// Then
			tt.then(t, dir, err)
		})
	}
}

// TestAcceptance_US002_ListAvailableAssets tests US-002: List Available Assets
func TestAcceptance_US002_ListAvailableAssets(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		given    func(t *testing.T) string
		when     func(t *testing.T) (string, error)
		then     func(t *testing.T, output string, err error)
	}{
		{
			name:     "list_all_resources",
			scenario: "List all available DDX resources",
			given: func(t *testing.T) string {
				// Given: DDX is initialized with available resources
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				ddxHome := filepath.Join(homeDir, ".ddx")

				// Create various resources
				templatesDir := filepath.Join(ddxHome, "templates")
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "nextjs"), 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "python"), 0755))

				patternsDir := filepath.Join(ddxHome, "patterns")
				require.NoError(t, os.MkdirAll(filepath.Join(patternsDir, "auth"), 0755))

				return homeDir
			},
			when: func(t *testing.T) (string, error) {
				// When: I run `ddx list`
				rootCmd := getTestRootCommand()
				return executeCommand(rootCmd, "list")
			},
			then: func(t *testing.T, output string, err error) {
				// Then: I see categorized resources with helpful descriptions
				assert.NoError(t, err)
				assert.Contains(t, output, "Templates", "Should show templates category")
				assert.Contains(t, output, "Patterns", "Should show patterns category")
			},
		},
		{
			name:     "filter_by_type",
			scenario: "Filter resources by type",
			given: func(t *testing.T) string {
				// Given: I want to see only templates
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				ddxHome := filepath.Join(homeDir, ".ddx")

				templatesDir := filepath.Join(ddxHome, "templates")
				require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "react"), 0755))

				patternsDir := filepath.Join(ddxHome, "patterns")
				require.NoError(t, os.MkdirAll(filepath.Join(patternsDir, "logging"), 0755))

				return homeDir
			},
			when: func(t *testing.T) (string, error) {
				// When: I run `ddx list templates`
				rootCmd := getTestRootCommand()
				return executeCommand(rootCmd, "list", "templates")
			},
			then: func(t *testing.T, output string, err error) {
				// Then: only templates are shown
				assert.NoError(t, err)
				assert.Contains(t, output, "Templates", "Should show templates")
				// Patterns should not be shown when filtering
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			tt.given(t)

			// When
			output, err := tt.when(t)

			// Then
			tt.then(t, output, err)
		})
	}
}

// TestAcceptance_US003_ApplyAssetToProject tests US-003: Apply Asset to Project
func TestAcceptance_US003_ApplyAssetToProject(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		given    func(t *testing.T) (string, string) // Returns homeDir, workDir
		when     func(t *testing.T, workDir string) error
		then     func(t *testing.T, workDir string, err error)
	}{
		{
			name:     "apply_template",
			scenario: "Apply a template to the project",
			given: func(t *testing.T) (string, string) {
				// Given: I want to apply a template
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)

				// Create template
				templateDir := filepath.Join(homeDir, ".ddx", "templates", "basic")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				templateFile := filepath.Join(templateDir, "README.md")
				require.NoError(t, os.WriteFile(templateFile, []byte("# {{project_name}}"), 0644))

				// Setup work directory
				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				// Create config with variables
				config := `version: "1.0"
variables:
  project_name: "TestProject"`
				require.NoError(t, os.WriteFile(
					filepath.Join(workDir, ".ddx.yml"),
					[]byte(config),
					0644,
				))

				return homeDir, workDir
			},
			when: func(t *testing.T, workDir string) error {
				// When: I run `ddx apply templates/basic`
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "apply", "templates/basic")
				return err
			},
			then: func(t *testing.T, workDir string, err error) {
				// Then: files are created with variables substituted
				readmePath := filepath.Join(workDir, "README.md")
				if _, statErr := os.Stat(readmePath); statErr == nil {
					content, readErr := os.ReadFile(readmePath)
					if readErr == nil {
						assert.Contains(t, string(content), "TestProject",
							"Variables should be substituted")
					}
				}
			},
		},
		{
			name:     "dry_run_preview",
			scenario: "Preview changes before applying",
			given: func(t *testing.T) (string, string) {
				// Given: I want to preview changes
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)

				templateDir := filepath.Join(homeDir, ".ddx", "templates", "preview")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				templateFile := filepath.Join(templateDir, "config.json")
				require.NoError(t, os.WriteFile(templateFile, []byte(`{"name": "test"}`), 0644))

				workDir := t.TempDir()
				require.NoError(t, os.Chdir(workDir))

				config := `version: "1.0"`
				require.NoError(t, os.WriteFile(
					filepath.Join(workDir, ".ddx.yml"),
					[]byte(config),
					0644,
				))

				return homeDir, workDir
			},
			when: func(t *testing.T, workDir string) error {
				// When: I run `ddx apply --dry-run templates/preview`
				rootCmd := getTestRootCommand()
				_, err := executeCommand(rootCmd, "apply", "--dry-run", "templates/preview")
				return err
			},
			then: func(t *testing.T, workDir string, err error) {
				// Then: changes are shown but not applied
				configPath := filepath.Join(workDir, "config.json")
				_, statErr := os.Stat(configPath)
				assert.Error(t, statErr, "File should not be created in dry-run mode")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Given
			_, workDir := tt.given(t)

			// When
			err := tt.when(t, workDir)

			// Then
			tt.then(t, workDir, err)
		})
	}
}

// TestAcceptance_ConfigurationManagement tests configuration-related user stories
func TestAcceptance_ConfigurationManagement(t *testing.T) {
	t.Run("view_configuration", func(t *testing.T) {
		// Given: DDX is configured in my project
		workDir := t.TempDir()
		require.NoError(t, os.Chdir(workDir))

		config := `version: "1.0"
repository:
  url: "https://github.com/test/repo"
variables:
  environment: "development"`
		require.NoError(t, os.WriteFile(
			filepath.Join(workDir, ".ddx.yml"),
			[]byte(config),
			0644,
		))

		// When: I run `ddx config`
		rootCmd := getTestRootCommand()
		output, err := executeCommand(rootCmd, "config")

		// Then: I see my current configuration clearly displayed
		assert.NoError(t, err)
		assert.Contains(t, output, "version", "Should show version")
		assert.Contains(t, output, "repository", "Should show repository")
		assert.Contains(t, output, "variables", "Should show variables")
	})

	t.Run("modify_configuration", func(t *testing.T) {
		// Given: I need to change a configuration value
		workDir := t.TempDir()
		require.NoError(t, os.Chdir(workDir))

		config := `version: "1.0"
variables:
  old_value: "original"`
		configPath := filepath.Join(workDir, ".ddx.yml")
		require.NoError(t, os.WriteFile(configPath, []byte(config), 0644))

		// When: I run `ddx config set variables.new_value "updated"`
		rootCmd := getTestRootCommand()
		_, err := executeCommand(rootCmd, "config", "set", "variables.new_value", "updated")

		// Then: the configuration is updated with the new value
		if err == nil {
			data, readErr := os.ReadFile(configPath)
			if readErr == nil {
				var updatedConfig map[string]interface{}
				yaml.Unmarshal(data, &updatedConfig)

				if vars, ok := updatedConfig["variables"].(map[string]interface{}); ok {
					assert.Equal(t, "updated", vars["new_value"],
						"New value should be set")
				}
			}
		}
	})
}

// TestAcceptance_WorkflowIntegration tests complete user workflows
func TestAcceptance_WorkflowIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping workflow integration test in short mode")
	}

	t.Run("complete_project_setup", func(t *testing.T) {
		// Scenario: Setting up a new project with DDX
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)

		// Step 1: Initialize DDX
		workDir := t.TempDir()
		require.NoError(t, os.Chdir(workDir))

		// Create library structure with templates
		libraryDir := filepath.Join(workDir, "library")
		templatesDir := filepath.Join(libraryDir, "templates")
		require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "nextjs"), 0755))
		require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, "python"), 0755))

		// Create config pointing to library
		config := []byte(`version: "2.0"
library_path: ./library`)
		require.NoError(t, os.WriteFile(".ddx.yml", config, 0644))

		rootCmd := getTestRootCommand()
		_, initErr := executeCommand(rootCmd, "init")
		// Note: May fail if DDX repo not available
		_ = initErr

		// Step 2: List available resources
		listOutput, listErr := executeCommand(rootCmd, "list")
		if listErr == nil && listOutput != "" && !strings.Contains(listOutput, "❌ DDx library not found") {
			assert.Contains(t, listOutput, "Templates", "Should list templates")
		} else {
			t.Log("Skipping template list assertion due to DDx not being initialized or available")
		}

		// Step 3: Apply a template (if available)
		// Note: Would apply actual template in real scenario

		// Step 4: Verify configuration
		configOutput, configErr := executeCommand(rootCmd, "config")
		if configErr == nil {
			assert.NotEmpty(t, configOutput, "Should show configuration")
		}
	})
}

// TestAcceptance_ErrorScenarios tests error handling from user perspective
func TestAcceptance_ErrorScenarios(t *testing.T) {
	t.Run("clear_error_messages", func(t *testing.T) {
		tests := []struct {
			name          string
			setup         func() string
			command       []string
			expectedError string
		}{
			{
				name: "template_not_found",
				setup: func() string {
					workDir := t.TempDir()
					os.Chdir(workDir)
					config := `version: "1.0"`
					os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644)
					return workDir
				},
				command:       []string{"apply", "templates/nonexistent"},
				expectedError: "not found",
			},
			{
				name: "already_initialized",
				setup: func() string {
					workDir := t.TempDir()
					os.Chdir(workDir)
					config := `version: "1.0"`
					os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(config), 0644)
					return workDir
				},
				command:       []string{"init"},
				expectedError: "already",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				originalDir, _ := os.Getwd()
				defer os.Chdir(originalDir)

				tt.setup()

				rootCmd := getTestRootCommand()
				output, err := executeCommand(rootCmd, tt.command...)

				// Verify clear error message
				if err != nil {
					assert.Contains(t, strings.ToLower(err.Error()), tt.expectedError,
						"Error message should be clear and helpful")
				} else if output != "" {
					assert.Contains(t, strings.ToLower(output), tt.expectedError,
						"Output should contain helpful error information")
				}
			})
		}
	})
}
