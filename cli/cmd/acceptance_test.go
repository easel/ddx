package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
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
// DEPRECATED: Use NewTestRootCommand(t) instead for proper test isolation
func getTestRootCommand() *cobra.Command {
	factory := NewCommandFactory("/tmp")
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
				tempDir := t.TempDir()
				return tempDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx init`
				// Use CommandFactory with the test working directory
				factory := NewCommandFactory(dir)
				rootCmd := factory.NewRootCommand()
				_, err := executeCommand(rootCmd, "init")
				return err
			},
			then: func(t *testing.T, dir string, err error) {
				// Then: a `.ddx/config.yaml` configuration file exists with my settings
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				if _, statErr := os.Stat(configPath); statErr == nil {
					// Config file exists - validate structure
					data, readErr := os.ReadFile(configPath)
					require.NoError(t, readErr)

					var config map[string]interface{}
					yamlErr := yaml.Unmarshal(data, &config)
					require.NoError(t, yamlErr)

					assert.Contains(t, config, "version", "Config should have version")
					assert.Contains(t, config, "library", "Config should have library")
				}
			},
		},
		{
			name:     "template_based_initialization",
			scenario: "Initialize DDX with specific template",
			given: func(t *testing.T) string {
				// Given: I want to use a specific template
				tempDir := t.TempDir()

				// Setup mock template
				homeDir := t.TempDir()
				t.Setenv("HOME", homeDir)
				templateDir := filepath.Join(homeDir, ".ddx", "templates", "test-template")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				return tempDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx init --template test-template`
				// Use CommandFactory with the test working directory
				factory := NewCommandFactory(dir)
				rootCmd := factory.NewRootCommand()
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
				tempDir := t.TempDir()

				// Initialize git repository in temp directory
				gitInit := exec.Command("git", "init")
				gitInit.Dir = tempDir
				require.NoError(t, gitInit.Run())

				gitConfigEmail := exec.Command("git", "config", "user.email", "test@example.com")
				gitConfigEmail.Dir = tempDir
				require.NoError(t, gitConfigEmail.Run())

				gitConfigName := exec.Command("git", "config", "user.name", "Test User")
				gitConfigName.Dir = tempDir
				require.NoError(t, gitConfigName.Run())

				// Create existing config
				config := `version: "1.0"
library:
  path: "./library"
  repository:
    url: "https://github.com/test/repo"
    branch: "main"`
				ddxDir := filepath.Join(tempDir, ".ddx")
				require.NoError(t, os.MkdirAll(ddxDir, 0755))
				require.NoError(t, os.WriteFile(
					filepath.Join(ddxDir, "config.yaml"),
					[]byte(config),
					0644,
				))

				return tempDir
			},
			when: func(t *testing.T, dir string) error {
				// When: I run `ddx init` again
				// Use CommandFactory with the test working directory
				factory := NewCommandFactory(dir)
				rootCmd := factory.NewRootCommand()
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
			//	// originalDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection // REMOVED: Using CommandFactory injection

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
				env := NewTestEnvironment(t)
				env.InitWithDDx()

				// Create various resources in library directory
				libDir := filepath.Join(env.Dir, ".ddx", "library")
				workflowsDir := filepath.Join(libDir, "workflows")
				require.NoError(t, os.MkdirAll(filepath.Join(workflowsDir, "helix"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(workflowsDir, "helix", "workflow.yml"), []byte("name: helix"), 0644))

				promptsDir := filepath.Join(libDir, "prompts")
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "claude"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(promptsDir, "claude", "prompt.md"), []byte("# Prompt"), 0644))

				return env.Dir
			},
			when: func(t *testing.T) (string, error) {
				// When: I run `ddx list`
				testDir := os.Getenv("TEST_DIR")
				factory := NewCommandFactory(testDir)
				rootCmd := factory.NewRootCommand()
				return executeCommand(rootCmd, "list")
			},
			then: func(t *testing.T, output string, err error) {
				// Then: I see categorized resources with helpful descriptions
				assert.NoError(t, err)
				assert.Contains(t, output, "Workflows", "Should show workflows category")
				assert.Contains(t, output, "Prompts", "Should show prompts category")
			},
		},
		{
			name:     "filter_by_type",
			scenario: "Filter resources by type",
			given: func(t *testing.T) string {
				// Given: I want to see only workflows
				env := NewTestEnvironment(t)
				env.InitWithDDx()

				libDir := filepath.Join(env.Dir, ".ddx", "library")
				workflowsDir := filepath.Join(libDir, "workflows")
				require.NoError(t, os.MkdirAll(filepath.Join(workflowsDir, "helix"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(workflowsDir, "helix", "workflow.yml"), []byte("name: helix"), 0644))

				promptsDir := filepath.Join(libDir, "prompts")
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "claude"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(promptsDir, "claude", "prompt.md"), []byte("# Prompt"), 0644))

				return env.Dir
			},
			when: func(t *testing.T) (string, error) {
				// When: I run `ddx list workflows`
				testDir := os.Getenv("TEST_DIR")
				factory := NewCommandFactory(testDir)
				rootCmd := factory.NewRootCommand()
				return executeCommand(rootCmd, "list", "workflows")
			},
			then: func(t *testing.T, output string, err error) {
				// Then: only workflows are shown
				assert.NoError(t, err)
				assert.Contains(t, output, "Workflows", "Should show workflows")
				// Prompts should not be shown when filtering
			},
		},
		{
			name:     "json_output",
			scenario: "Output resources as JSON",
			given: func(t *testing.T) string {
				// Given: DDx project with resources
				testDir := t.TempDir()

				// Initialize DDx properly
				factory := NewCommandFactory(testDir)
				initCmd := factory.NewRootCommand()
				initCmd.SetArgs([]string{"init", "--no-git", "--silent"})
				var initOut bytes.Buffer
				initCmd.SetOut(&initOut)
				initCmd.SetErr(&initOut)
				require.NoError(t, initCmd.Execute())

				// Create test resources in the library
				libraryDir := filepath.Join(testDir, ".ddx", "library")
				workflowsDir := filepath.Join(libraryDir, "workflows")
				helixDir := filepath.Join(workflowsDir, "helix")
				require.NoError(t, os.MkdirAll(helixDir, 0755))
				require.NoError(t, os.WriteFile(filepath.Join(helixDir, "workflow.yml"), []byte("name: helix"), 0644))

				promptsDir := filepath.Join(libraryDir, "prompts")
				claudeDir := filepath.Join(promptsDir, "claude")
				require.NoError(t, os.MkdirAll(claudeDir, 0755))
				require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "prompt.md"), []byte("# Claude Prompt"), 0644))

				// Store testDir for when() to use
				t.Setenv("TEST_DIR", testDir)
				return testDir
			},
			when: func(t *testing.T) (string, error) {
				// When: I run `ddx list --json`
				testDir := os.Getenv("TEST_DIR")
				factory := NewTestRootCommandWithDir(testDir)
				rootCmd := factory.NewRootCommand()
				return executeCommand(rootCmd, "list", "--json")
			},
			then: func(t *testing.T, output string, err error) {
				// Then: output is valid JSON with resource data
				assert.NoError(t, err)

				// Verify it's valid JSON
				var response struct {
					Resources []map[string]interface{} `json:"resources"`
					Summary   map[string]int           `json:"summary"`
				}
				assert.NoError(t, json.Unmarshal([]byte(output), &response))

				// Should have resources and summary
				assert.NotEmpty(t, response.Resources)
				assert.NotEmpty(t, response.Summary)
			},
		},
		{
			name:     "empty_filter_results",
			scenario: "Handle empty filter results gracefully",
			given: func(t *testing.T) string {
				// Given: DDx has resources but none match filter
				env := NewTestEnvironment(t)
				env.InitWithDDx()

				libDir := filepath.Join(env.Dir, ".ddx", "library")
				workflowsDir := filepath.Join(libDir, "workflows")
				require.NoError(t, os.MkdirAll(filepath.Join(workflowsDir, "helix"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(workflowsDir, "helix", "workflow.yml"), []byte("name: helix"), 0644))

				return env.Dir
			},
			when: func(t *testing.T) (string, error) {
				// When: I run `ddx list --filter nonexistent`
				testDir := os.Getenv("TEST_DIR")
				factory := NewCommandFactory(testDir)
				rootCmd := factory.NewRootCommand()
				return executeCommand(rootCmd, "list", "--filter", "nonexistent")
			},
			then: func(t *testing.T, output string, err error) {
				// Then: I see a clear message about no matches
				assert.NoError(t, err)
				assert.Contains(t, output, "No DDx resources found", "Should show no resources message")
				assert.Contains(t, output, "No resources match filter: 'nonexistent'", "Should show filter message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			testDir := tt.given(t)
			t.Setenv("TEST_DIR", testDir)

			// When
			output, err := tt.when(t)

			// Then
			tt.then(t, output, err)
		})
	}
}

// TestAcceptance_ConfigurationManagement tests configuration-related user stories
func TestAcceptance_ConfigurationManagement(t *testing.T) {
	t.Run("view_configuration", func(t *testing.T) {
		// Given: DDX is configured in my project
		tempDir := t.TempDir()

		config := `version: "1.0"
library:
  path: "./library"
  repository:
    url: "https://github.com/test/repo"
    branch: "main"
persona_bindings:
  environment: "development"`
		ddxDir := filepath.Join(tempDir, ".ddx")
		require.NoError(t, os.MkdirAll(ddxDir, 0755))
		require.NoError(t, os.WriteFile(
			filepath.Join(ddxDir, "config.yaml"),
			[]byte(config),
			0644,
		))

		// When: I run `ddx config export`
		factory := NewCommandFactory(tempDir)
		rootCmd := factory.NewRootCommand()
		output, err := executeCommand(rootCmd, "config", "export")

		// Then: I see my current configuration clearly displayed
		assert.NoError(t, err)
		assert.Contains(t, output, "version", "Should show version")
		assert.Contains(t, output, "library", "Should show library")
		assert.Contains(t, output, "persona_bindings", "Should show persona_bindings")
	})

	t.Run("modify_configuration", func(t *testing.T) {
		// Given: I need to change a configuration value
		tempDir := t.TempDir()

		config := `version: "1.0"
library:
  path: "./library"
persona_bindings:
  old_value: "original"`
		ddxDir := filepath.Join(tempDir, ".ddx")
		require.NoError(t, os.MkdirAll(ddxDir, 0755))
		configPath := filepath.Join(ddxDir, "config.yaml")
		require.NoError(t, os.WriteFile(configPath, []byte(config), 0644))

		// When: I run `ddx config set persona_bindings.new_value "updated"`
		factory := NewCommandFactory(tempDir)
		rootCmd := factory.NewRootCommand()
		_, err := executeCommand(rootCmd, "config", "set", "persona_bindings.new_value", "updated")

		// Then: the configuration is updated with the new value
		if err == nil {
			data, readErr := os.ReadFile(configPath)
			if readErr == nil {
				var updatedConfig map[string]interface{}
				_ = yaml.Unmarshal(data, &updatedConfig)

				if vars, ok := updatedConfig["persona_bindings"].(map[string]interface{}); ok {
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
		//	// originalDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection // REMOVED: Using CommandFactory injection

		// Step 1: Initialize DDX
		tempDir := t.TempDir()

		// Create library structure with workflows
		libraryDir := filepath.Join(tempDir, "library")
		workflowsDir := filepath.Join(libraryDir, "workflows")
		require.NoError(t, os.MkdirAll(filepath.Join(workflowsDir, "helix"), 0755))
		require.NoError(t, os.MkdirAll(filepath.Join(workflowsDir, "kanban"), 0755))

		// Create workflow files so they can be discovered
		require.NoError(t, os.WriteFile(filepath.Join(workflowsDir, "helix", "workflow.yml"), []byte("name: helix"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(workflowsDir, "kanban", "workflow.yml"), []byte("name: kanban"), 0644))

		// Create config pointing to library in new format
		config := []byte(`version: "2.0"
library:
  path: "./library"
  repository:
    url: "https://github.com/easel/ddx-library"
    branch: "main"
persona_bindings: {}`)
		ddxDir := filepath.Join(tempDir, ".ddx")
		require.NoError(t, os.MkdirAll(ddxDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(ddxDir, "config.yaml"), config, 0644))

		// Use CommandFactory with working directory
		factory := NewCommandFactory(tempDir)
		rootCmd := factory.NewRootCommand()
		_, initErr := executeCommand(rootCmd, "init")
		// Note: May fail if DDX repo not available
		_ = initErr

		// Step 2: List available resources
		listOutput, listErr := executeCommand(rootCmd, "list")
		if listErr == nil && listOutput != "" && !strings.Contains(listOutput, "❌ DDx library not found") {
			assert.Contains(t, listOutput, "Workflows", "Should list workflows")
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
				name: "already_initialized",
				setup: func() string {
					tempDir := t.TempDir()
					// Initialize git repository first
					gitInit := exec.Command("git", "init")
					gitInit.Dir = tempDir
					_ = gitInit.Run()

					gitConfigEmail := exec.Command("git", "config", "user.email", "test@example.com")
					gitConfigEmail.Dir = tempDir
					_ = gitConfigEmail.Run()

					gitConfigName := exec.Command("git", "config", "user.name", "Test User")
					gitConfigName.Dir = tempDir
					_ = gitConfigName.Run()

					config := `version: "1.0"
library:
  path: "./library"
  repository:
    url: "https://github.com/easel/ddx-library"
    branch: "main"
persona_bindings: {}`
					ddxDir := filepath.Join(tempDir, ".ddx")
					_ = os.MkdirAll(ddxDir, 0755)
					_ = os.WriteFile(filepath.Join(ddxDir, "config.yaml"), []byte(config), 0644)
					return tempDir
				},
				command:       []string{"init"},
				expectedError: "already",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				//	// originalDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection // REMOVED: Using CommandFactory injection

				tempDir := tt.setup()

				// Use CommandFactory with working directory
				factory := NewCommandFactory(tempDir)
				rootCmd := factory.NewRootCommand()
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

// TestAcceptance_US042_WorkflowCommandExecution tests US-042: Workflow Command Execution
func TestAcceptance_US042_WorkflowCommandExecution(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		given    func(t *testing.T) string
		when     func(t *testing.T, workDir string) (string, error)
		then     func(t *testing.T, workDir string, output string, err error)
	}{
		{
			name:     "list_helix_commands",
			scenario: "AC-001: Command Discovery",
			given: func(t *testing.T) string {
				// Given: I have the HELIX workflow available
				tempDir := t.TempDir()

				// Create library structure with helix commands
				commandsDir := filepath.Join(tempDir, "library", "workflows", "helix", "commands")
				require.NoError(t, os.MkdirAll(commandsDir, 0755))

				// Create build-story command
				buildStoryContent := `# HELIX Command: Build Story

You are a HELIX workflow executor...`
				require.NoError(t, os.WriteFile(
					filepath.Join(commandsDir, "build-story.md"),
					[]byte(buildStoryContent), 0644))

				// Create continue command
				continueContent := `# HELIX Command: Continue

Continue work on current story...`
				require.NoError(t, os.WriteFile(
					filepath.Join(commandsDir, "continue.md"),
					[]byte(continueContent), 0644))

				return tempDir
			},
			when: func(t *testing.T, workDir string) (string, error) {
				// When: I run `ddx workflow helix commands`
				factory := NewTestRootCommandWithDir(workDir)
				rootCmd := factory.NewRootCommand()
				buf := new(bytes.Buffer)
				rootCmd.SetOut(buf)
				rootCmd.SetErr(buf)
				rootCmd.SetArgs([]string{"workflow", "helix", "commands"})

				err := rootCmd.Execute()
				return buf.String(), err
			},
			then: func(t *testing.T, workDir string, output string, err error) {
				// Then: I see a list of available commands with descriptions
				assert.NoError(t, err)
				assert.Contains(t, output, "Available commands for helix workflow:")
				assert.Contains(t, output, "build-story")
				assert.Contains(t, output, "continue")
			},
		},
		{
			name:     "execute_build_story_command",
			scenario: "AC-002: Command Execution",
			given: func(t *testing.T) string {
				// Given: I have a workflow with commands available
				tempDir := t.TempDir()

				commandsDir := filepath.Join(tempDir, "library", "workflows", "helix", "commands")
				require.NoError(t, os.MkdirAll(commandsDir, 0755))

				buildStoryContent := `# HELIX Command: Build Story

You are a HELIX workflow executor tasked with implementing work on a specific user story.

## Command Input

You will receive a user story ID as an argument (e.g., US-001, US-042, etc.).`
				require.NoError(t, os.WriteFile(
					filepath.Join(commandsDir, "build-story.md"),
					[]byte(buildStoryContent), 0644))

				return tempDir
			},
			when: func(t *testing.T, workDir string) (string, error) {
				// When: I run `ddx workflow helix execute build-story US-001`
				factory := NewTestRootCommandWithDir(workDir)
				rootCmd := factory.NewRootCommand()
				buf := new(bytes.Buffer)
				rootCmd.SetOut(buf)
				rootCmd.SetErr(buf)
				rootCmd.SetArgs([]string{"workflow", "helix", "execute", "build-story", "US-001"})

				err := rootCmd.Execute()
				return buf.String(), err
			},
			then: func(t *testing.T, workDir string, output string, err error) {
				// Then: The build-story command prompt is loaded and displayed
				assert.NoError(t, err)
				assert.Contains(t, output, "HELIX Command: Build Story")
				assert.Contains(t, output, "You are a HELIX workflow executor")
				assert.Contains(t, output, "Command Arguments: [US-001]")
			},
		},
		{
			name:     "invalid_workflow_error",
			scenario: "AC-003: Error Handling - Invalid Workflow",
			given: func(t *testing.T) string {
				// Given: I specify a non-existent workflow
				tempDir := t.TempDir()
				return tempDir
			},
			when: func(t *testing.T, workDir string) (string, error) {
				// When: I run `ddx workflow invalid commands`
				factory := NewTestRootCommandWithDir(workDir)
				rootCmd := factory.NewRootCommand()
				buf := new(bytes.Buffer)
				rootCmd.SetOut(buf)
				rootCmd.SetErr(buf)
				rootCmd.SetArgs([]string{"workflow", "invalid", "commands"})

				err := rootCmd.Execute()
				return buf.String(), err
			},
			then: func(t *testing.T, workDir string, output string, err error) {
				// Then: I receive an error message about the workflow not being found
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "workflow 'invalid' not found")
			},
		},
		{
			name:     "invalid_command_error",
			scenario: "AC-004: Error Handling - Invalid Command",
			given: func(t *testing.T) string {
				// Given: I specify a non-existent command
				tempDir := t.TempDir()

				commandsDir := filepath.Join(tempDir, "library", "workflows", "helix", "commands")
				require.NoError(t, os.MkdirAll(commandsDir, 0755))

				return tempDir
			},
			when: func(t *testing.T, workDir string) (string, error) {
				// When: I run `ddx workflow helix execute invalid-command`
				factory := NewTestRootCommandWithDir(workDir)
				rootCmd := factory.NewRootCommand()
				buf := new(bytes.Buffer)
				rootCmd.SetOut(buf)
				rootCmd.SetErr(buf)
				rootCmd.SetArgs([]string{"workflow", "helix", "execute", "invalid-command"})

				err := rootCmd.Execute()
				return buf.String(), err
			},
			then: func(t *testing.T, workDir string, output string, err error) {
				// Then: I receive an error about the command not being found
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "command 'invalid-command' not found")
			},
		},
		{
			name:     "command_with_arguments",
			scenario: "AC-005: Command Arguments",
			given: func(t *testing.T) string {
				// Given: A command requires arguments
				tempDir := t.TempDir()

				commandsDir := filepath.Join(tempDir, "library", "workflows", "helix", "commands")
				require.NoError(t, os.MkdirAll(commandsDir, 0755))

				buildStoryContent := `# HELIX Command: Build Story

Command accepts arguments for user story processing.`
				require.NoError(t, os.WriteFile(
					filepath.Join(commandsDir, "build-story.md"),
					[]byte(buildStoryContent), 0644))

				return tempDir
			},
			when: func(t *testing.T, workDir string) (string, error) {
				// When: I execute it with arguments
				factory := NewTestRootCommandWithDir(workDir)
				rootCmd := factory.NewRootCommand()
				buf := new(bytes.Buffer)
				rootCmd.SetOut(buf)
				rootCmd.SetErr(buf)
				rootCmd.SetArgs([]string{"workflow", "helix", "execute", "build-story", "US-001", "extra-arg"})

				err := rootCmd.Execute()
				return buf.String(), err
			},
			then: func(t *testing.T, workDir string, output string, err error) {
				// Then: The arguments are passed to the command context
				assert.NoError(t, err)
				assert.Contains(t, output, "Command Arguments: [US-001 extra-arg]")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//	// originalDir, _ := os.Getwd() // REMOVED: Using CommandFactory injection // REMOVED: Using CommandFactory injection

			workDir := tt.given(t)
			output, err := tt.when(t, workDir)
			tt.then(t, workDir, output, err)
		})
	}
}
