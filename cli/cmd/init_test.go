package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Helper function to create a fresh root command for tests
func getInitTestRootCommand(workingDir string) *cobra.Command {
	if workingDir == "" {
		workingDir = "/tmp"
	}
	factory := NewCommandFactory(workingDir)
	return factory.NewRootCommand()
}

// Helper to execute command with captured output
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

// Helper to setup test environment
func setupTestDir(t *testing.T) (string, func()) {
	tempDir := t.TempDir()

	// Change to temp directory
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = tempDir
	require.NoError(t, gitCmd.Run())

	// Configure git user for tests
	gitCmd = exec.Command("git", "config", "user.email", "test@example.com")
	gitCmd.Dir = tempDir
	require.NoError(t, gitCmd.Run())

	gitCmd = exec.Command("git", "config", "user.name", "Test User")
	gitCmd.Dir = tempDir
	require.NoError(t, gitCmd.Run())

	cleanup := func() {
	}

	return tempDir, cleanup
}

// TestInitCommand tests the init command
func TestInitCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		setup       func(t *testing.T, dir string)
		validate    func(t *testing.T, dir string, output string, err error)
		expectError bool
	}{
		{
			name: "basic initialization",
			args: []string{"init", "--no-git"},
			setup: func(t *testing.T, dir string) {
				// Git repository already initialized by setupTestDir
				t.Setenv("DDX_TEST_MODE", "1")
			},
			validate: func(t *testing.T, dir string, output string, err error) {
				// Check .ddx/config.yaml was created
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath)

				// Verify config content
				data, err := os.ReadFile(configPath)
				require.NoError(t, err)

				var config map[string]interface{}
				err = yaml.Unmarshal(data, &config)
				require.NoError(t, err)

				assert.Contains(t, config, "version")
				assert.Contains(t, config, "library")
				if library, ok := config["library"].(map[string]interface{}); ok {
					assert.Contains(t, library, "repository")
				}
			},
			expectError: false,
		},
		{
			name: "init with force flag when config exists",
			args: []string{"init", "--force", "--no-git"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
				// Create existing config in new format
				existingConfig := `version: "0.9"
library:
  path: "./library"
  repository:
    url: "https://old.repo"
    branch: "main"
    subtree: "library"
persona_bindings: {}
`
				ddxDir := filepath.Join(dir, ".ddx")
				require.NoError(t, os.MkdirAll(ddxDir, 0755))
				configPath := filepath.Join(ddxDir, "config.yaml")
				require.NoError(t, os.WriteFile(configPath, []byte(existingConfig), 0644))
			},
			validate: func(t *testing.T, dir string, output string, err error) {
				// Config should be overwritten
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				data, err := os.ReadFile(configPath)
				require.NoError(t, err)

				var config map[string]interface{}
				err = yaml.Unmarshal(data, &config)
				require.NoError(t, err)

				// With --force flag, creates new config with default version
				assert.Equal(t, "1.0", config["version"])
			},
			expectError: false,
		},
		{
			name: "init without force when config exists",
			args: []string{"init", "--no-git"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
				// Create existing config in new format
				ddxDir := filepath.Join(dir, ".ddx")
				require.NoError(t, os.MkdirAll(ddxDir, 0755))
				configPath := filepath.Join(ddxDir, "config.yaml")
				require.NoError(t, os.WriteFile(configPath, []byte("version: \"1.0\""), 0644))
			},
			validate: func(t *testing.T, dir string, output string, err error) {
				// Should fail
				assert.Error(t, err)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh command for test isolation
			// (flags are now local to the command)

			dir, cleanup := setupTestDir(t)
			defer cleanup()

			if tt.setup != nil {
				tt.setup(t, dir)
			}

			// Create new root command for each test with the test directory
			factory := NewCommandFactory(dir)
			rootCmd := factory.NewRootCommand()

			output, err := executeCommand(rootCmd, tt.args...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validate != nil {
				tt.validate(t, dir, output, err)
			}
		})
	}
}

// TestInitCommand_Help tests the help output
func TestInitCommand_Help(t *testing.T) {
	factory := NewCommandFactory("/tmp")
	rootCmd := factory.NewRootCommand()

	output, err := executeCommand(rootCmd, "init", "--help")

	assert.NoError(t, err)
	assert.Contains(t, output, "Initialize DDx")
	assert.Contains(t, output, "--force")
	assert.Contains(t, output, "--no-git")
}

// TestInitCommand_US017_InitializeConfiguration tests US-017 Initialize Configuration
func TestInitCommand_US017_InitializeConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T, dir string)
		args           []string
		validateOutput func(t *testing.T, dir, output string, err error)
		expectError    bool
	}{
		{
			name: "creates_initial_config_with_sensible_defaults",
			args: []string{"init", "--no-git"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should create .ddx/config.yaml with sensible defaults
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath, "Should create .ddx/config.yaml")

				data, err := os.ReadFile(configPath)
				require.NoError(t, err)

				var config map[string]interface{}
				err = yaml.Unmarshal(data, &config)
				require.NoError(t, err)

				assert.Contains(t, config, "version")
				assert.Contains(t, config, "library")
				if library, ok := config["library"].(map[string]interface{}); ok {
					assert.Contains(t, library, "repository")
				}
				// New config format uses persona_bindings instead of variables
				// Library path is nested under library object, not at root
			},
			expectError: false,
		},
		{
			name: "detects_project_type_javascript",
			args: []string{"init", "--no-git"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
				// Create package.json to simulate JavaScript project
				require.NoError(t, os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name": "test"}`), 0644))
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				data, err := os.ReadFile(configPath)
				require.NoError(t, err)

				var config map[string]interface{}
				err = yaml.Unmarshal(data, &config)
				require.NoError(t, err)

				// Project type detection has been removed - init just creates basic config
				assert.Contains(t, config, "version")
				assert.Contains(t, config, "library")
			},
			expectError: false,
		},
		{
			name: "detects_project_type_go",
			args: []string{"init", "--no-git"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
				// Create go.mod to simulate Go project
				require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0644))
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				data, err := os.ReadFile(configPath)
				require.NoError(t, err)

				var config map[string]interface{}
				err = yaml.Unmarshal(data, &config)
				require.NoError(t, err)

				// Project type detection has been removed - init just creates basic config
				assert.Contains(t, config, "version")
				assert.Contains(t, config, "library")
			},
			expectError: false,
		},
		{
			name: "validates_configuration_during_creation",
			args: []string{"init", "--no-git"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should pass validation without error
				assert.NoError(t, err, "Configuration validation should pass")
				assert.Contains(t, output, "✅ DDx initialized successfully!")
			},
			expectError: false,
		},
		{
			name: "force_overwrites_without_backup",
			args: []string{"init", "--force"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
				// Create existing config
				existingConfig := `version: "0.9"
repository:
  url: "https://old.repo"
`
				ddxDir := filepath.Join(dir, ".ddx")
				require.NoError(t, os.MkdirAll(ddxDir, 0755))
				configPath := filepath.Join(ddxDir, "config.yaml")
				require.NoError(t, os.WriteFile(configPath, []byte(existingConfig), 0644))
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should NOT create backup or show backup message
				assert.NotContains(t, output, "💾 Created backup of existing config")
				assert.NotContains(t, output, "backup")

				// Should NOT have backup file
				backupFiles, _ := filepath.Glob(filepath.Join(dir, ".ddx", "config.yaml.backup.*"))
				assert.Equal(t, 0, len(backupFiles), "Should not create backup file")

				// Should successfully overwrite config
				assert.Contains(t, output, "✅ DDx initialized successfully!")
			},
			expectError: false,
		},
		{
			name: "no_git_flag_functionality",
			args: []string{"init", "--no-git"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should create config successfully without git operations
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath, "Should create config with --no-git flag")
			},
			expectError: false,
		},
		{
			name: "includes_example_variable_definitions",
			args: []string{"init"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				data, err := os.ReadFile(configPath)
				require.NoError(t, err)

				var config map[string]interface{}
				err = yaml.Unmarshal(data, &config)
				require.NoError(t, err)

				// Variable definitions have been removed - init creates minimal config
				assert.Contains(t, config, "version")
				assert.Contains(t, config, "library")
			},
			expectError: false,
		},
		{
			name: "commits_config_file_to_git",
			args: []string{"init"},
			setup: func(t *testing.T, dir string) {
				// Unset DDX_TEST_MODE - we want real git behavior
				os.Unsetenv("DDX_TEST_MODE")

				// Create initial commit (required for git subtree)
				readmePath := filepath.Join(dir, "README.md")
				require.NoError(t, os.WriteFile(readmePath, []byte("# Test Project"), 0644))

				gitAdd := exec.Command("git", "add", "README.md")
				gitAdd.Dir = dir
				require.NoError(t, gitAdd.Run())

				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Config file should be created
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath, "Config file should exist")

				// Check git log for config commit
				gitLog := exec.Command("git", "log", "--oneline", "--all")
				gitLog.Dir = dir
				logOutput, err := gitLog.CombinedOutput()
				require.NoError(t, err, "Should be able to read git log")

				logStr := string(logOutput)
				assert.Contains(t, logStr, "chore: initialize DDx configuration", "Should have config commit")
			},
			expectError: false,
		},
		{
			name: "skips_commit_in_test_mode",
			args: []string{"init"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Config file should be created
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath, "Config file should exist")

				// Git log should not have config commit (test mode skips commits)
				gitLog := exec.Command("git", "log", "--oneline", "--all")
				gitLog.Dir = dir
				logOutput, _ := gitLog.CombinedOutput()
				logStr := string(logOutput)

				// In test mode, no commits should be made at all
				assert.Empty(t, logStr, "Should have no commits in test mode")
			},
			expectError: false,
		},
		{
			name: "skips_commit_with_no_git_flag",
			args: []string{"init", "--no-git"},
			setup: func(t *testing.T, dir string) {
				// Don't set test mode, but use --no-git flag
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Config file should be created
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath, "Config file should exist")

				// Git log should not have config commit (--no-git skips commits)
				gitLog := exec.Command("git", "log", "--oneline", "--all")
				gitLog.Dir = dir
				logOutput, _ := gitLog.CombinedOutput()
				logStr := string(logOutput)

				// With --no-git, no commits should be made
				assert.Empty(t, logStr, "Should have no commits with --no-git flag")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := setupTestDir(t)
			defer cleanup()

			if tt.setup != nil {
				tt.setup(t, dir)
			}

			// Create new root command for each test with the test directory
			factory := NewCommandFactory(dir)
			rootCmd := factory.NewRootCommand()

			output, err := executeCommand(rootCmd, tt.args...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validateOutput != nil {
				tt.validateOutput(t, dir, output, err)
			}
		})
	}
}

// TestInitCommand_US014_SynchronizationSetup tests US-014 synchronization initialization
func TestInitCommand_US014_SynchronizationSetup(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T, dir string)
		args           []string
		validateOutput func(t *testing.T, dir, output string, err error)
		expectError    bool
	}{
		{
			name: "basic_sync_initialization",
			args: []string{"init"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should show DDx initialization progress
				assert.Contains(t, output, "Initializing DDx")
				assert.Contains(t, output, "DDx initialized successfully")
				// Verify config is created
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath, "Should create config file")
			},
			expectError: false,
		},
		{
			name: "sync_initialization_with_custom_repository",
			args: []string{"init", "--force"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
				// Create existing config with custom repo
				existingConfig := `version: "1.0"
repository:
  url: "https://github.com/custom/repo"
  branch: "develop"
`
				ddxDir := filepath.Join(dir, ".ddx")
				require.NoError(t, os.MkdirAll(ddxDir, 0755))
				configPath := filepath.Join(ddxDir, "config.yaml")
				require.NoError(t, os.WriteFile(configPath, []byte(existingConfig), 0644))
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should handle custom repository successfully
				assert.Contains(t, output, "DDx initialized successfully")
				assert.NotContains(t, output, "backup", "Should not create or mention backup")
			},
			expectError: false,
		},
		{
			name: "sync_initialization_fresh_project",
			args: []string{"init"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
				// Fresh project - no existing files
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should work with new projects
				assert.Contains(t, output, "Initializing DDx")
				assert.Contains(t, output, "DDx initialized successfully")

				// Check .ddx/config.yaml was created with sync config
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath)

				data, err := os.ReadFile(configPath)
				require.NoError(t, err)

				var config map[string]interface{}
				err = yaml.Unmarshal(data, &config)
				require.NoError(t, err)

				assert.Contains(t, config, "library")
				if library, ok := config["library"].(map[string]interface{}); ok {
					assert.Contains(t, library, "repository")
					if repo, ok := library["repository"].(map[string]interface{}); ok {
						assert.Contains(t, repo, "url")
						assert.Contains(t, repo, "branch")
					}
				}
			},
			expectError: false,
		},
		{
			name: "sync_initialization_existing_project",
			args: []string{"init", "--force"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
				// Create existing project files
				require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Existing Project"), 0644))
				require.NoError(t, os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name": "test"}`), 0644))
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should handle existing project files appropriately
				assert.Contains(t, output, "Initializing DDx")
				assert.Contains(t, output, "DDx initialized successfully")

				// Existing files should remain untouched
				assert.FileExists(t, filepath.Join(dir, "README.md"))
				assert.FileExists(t, filepath.Join(dir, "package.json"))
			},
			expectError: false,
		},
		{
			name: "sync_initialization_validation_success",
			args: []string{"init"},
			setup: func(t *testing.T, dir string) {
				t.Setenv("DDX_TEST_MODE", "1")
			},
			validateOutput: func(t *testing.T, dir, output string, err error) {
				// Should validate all sync settings
				assert.Contains(t, output, "Initializing DDx")
				assert.Contains(t, output, "DDx initialized successfully")
				// Verify config file exists with proper structure
				configPath := filepath.Join(dir, ".ddx", "config.yaml")
				assert.FileExists(t, configPath, "Should create config file")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := setupTestDir(t)
			defer cleanup()

			if tt.setup != nil {
				tt.setup(t, dir)
			}

			// Create new root command for each test with the test directory
			factory := NewCommandFactory(dir)
			rootCmd := factory.NewRootCommand()

			output, err := executeCommand(rootCmd, tt.args...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validateOutput != nil {
				tt.validateOutput(t, dir, output, err)
			}
		})
	}
}
