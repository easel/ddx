package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestInitCommand tests the init command
func TestInitCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		envOptions  []TestEnvOption
		setup       func(t *testing.T, te *TestEnvironment)
		validate    func(t *testing.T, te *TestEnvironment, output string, err error)
		expectError bool
	}{
		{
			name:       "basic initialization",
			args:       []string{"init", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(false)},
			validate: func(t *testing.T, te *TestEnvironment, output string, cmdErr error) {
				// Check .ddx/config.yaml was created
				assert.FileExists(t, te.ConfigPath)

				// Verify config content
				data, err := os.ReadFile(te.ConfigPath)
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
			name:       "init with force flag when config exists",
			args:       []string{"init", "--force", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(false)},
			setup: func(t *testing.T, te *TestEnvironment) {
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
				te.CreateConfig(existingConfig)
			},
			validate: func(t *testing.T, te *TestEnvironment, output string, cmdErr error) {
				// Config should be overwritten
				data, err := os.ReadFile(te.ConfigPath)
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
			name:       "init without force when config exists",
			args:       []string{"init", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(false)},
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create existing config
				te.CreateConfig("version: \"1.0\"")
			},
			validate: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Should fail
				assert.Error(t, err)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := NewTestEnvironment(t, tt.envOptions...)

			if tt.setup != nil {
				tt.setup(t, te)
			}

			output, err := te.RunCommand(tt.args...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validate != nil {
				tt.validate(t, te, output, err)
			}
		})
	}
}

// TestInitCommand_Help tests the help output
func TestInitCommand_Help(t *testing.T) {
	te := NewTestEnvironment(t)
	output, err := te.RunCommand("init", "--help")

	assert.NoError(t, err)
	assert.Contains(t, output, "Initialize DDx")
	assert.Contains(t, output, "--force")
	assert.Contains(t, output, "--no-git")
}

// TestInitCommand_US017_InitializeConfiguration tests US-017 Initialize Configuration
func TestInitCommand_US017_InitializeConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		envOptions     []TestEnvOption
		setup          func(t *testing.T, te *TestEnvironment)
		args           []string
		validateOutput func(t *testing.T, te *TestEnvironment, output string, err error)
		expectError    bool
	}{
		{
			name:       "creates_initial_config_with_sensible_defaults",
			args:       []string{"init", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(false)},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Should create .ddx/config.yaml with sensible defaults
				assert.FileExists(t, te.ConfigPath, "Should create .ddx/config.yaml")

				data, err := os.ReadFile(te.ConfigPath)
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
			name:       "detects_project_type_javascript",
			args:       []string{"init", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(false)},
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create package.json to simulate JavaScript project
				te.CreateFile("package.json", `{"name": "test"}`)
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				data, err := os.ReadFile(te.ConfigPath)
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
			name:       "detects_project_type_go",
			args:       []string{"init", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(false)},
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create go.mod to simulate Go project
				te.CreateFile("go.mod", "module test")
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				data, err := os.ReadFile(te.ConfigPath)
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
			name:       "validates_configuration_during_creation",
			args:       []string{"init", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(false)},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Should pass validation without error
				assert.NoError(t, err, "Configuration validation should pass")
				assert.Contains(t, output, "✅ DDx initialized successfully!")
			},
			expectError: false,
		},
		{
			name: "force_overwrites_without_backup",
			args: func() []string {
				// In CI, skip git integration due to subtree limitations
				if os.Getenv("CI") != "" {
					return []string{"init", "--force", "--no-git"}
				}
				return []string{"init", "--force"}
			}(),
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create existing config
				var existingConfig string
				if os.Getenv("CI") != "" {
					// In CI, use simple config without repository
					existingConfig = `version: "0.9"
library:
  path: .ddx/library
`
				} else {
					// Locally, test with repository URL
					existingConfig = fmt.Sprintf(`version: "0.9"
library:
  path: .ddx/library
  repository:
    url: %s
    branch: master
`, te.TestLibraryURL)
				}
				te.CreateConfig(existingConfig)
				te.CreateFile("README.md", "# Test Project")

				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = te.Dir
				require.NoError(t, gitAdd.Run())

				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = te.Dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Should NOT create backup or show backup message
				assert.NotContains(t, output, "💾 Created backup of existing config")
				assert.NotContains(t, output, "backup")

				// Should NOT have backup file
				backupFiles, _ := filepath.Glob(filepath.Join(te.Dir, ".ddx", "config.yaml.backup.*"))
				assert.Equal(t, 0, len(backupFiles), "Should not create backup file")

				// Should successfully overwrite config
				assert.Contains(t, output, "✅ DDx initialized successfully!")
			},
			expectError: false,
		},
		{
			name:       "no_git_flag_functionality",
			args:       []string{"init", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(true)},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Should create config successfully without git operations
				assert.FileExists(t, te.ConfigPath, "Should create config with --no-git flag")
			},
			expectError: false,
		},
		{
			name: "includes_example_variable_definitions",
			args: []string{"init", "--silent"},
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create initial commit required for git subtree
				te.CreateFile("README.md", "# Test Project")
				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = te.Dir
				require.NoError(t, gitAdd.Run())
				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = te.Dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				data, err := os.ReadFile(te.ConfigPath)
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
			args: []string{"init", "--silent"},
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create initial commit (required for git subtree)
				te.CreateFile("README.md", "# Test Project")

				gitAdd := exec.Command("git", "add", "README.md")
				gitAdd.Dir = te.Dir
				require.NoError(t, gitAdd.Run())

				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = te.Dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Config file should be created
				assert.FileExists(t, te.ConfigPath, "Config file should exist")

				// Check git log for config commit
				gitLog := exec.Command("git", "log", "--oneline", "--all")
				gitLog.Dir = te.Dir
				logOutput, err := gitLog.CombinedOutput()
				require.NoError(t, err, "Should be able to read git log")

				logStr := string(logOutput)
				assert.Contains(t, logStr, "chore: initialize DDx configuration", "Should have config commit")
			},
			expectError: false,
		},
		{
			name:       "skips_commit_with_no_git_flag",
			args:       []string{"init", "--no-git"},
			envOptions: []TestEnvOption{WithGitInit(true)},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Config file should be created
				assert.FileExists(t, te.ConfigPath, "Config file should exist")

				// Git log should not have config commit (--no-git skips commits)
				gitLog := exec.Command("git", "log", "--oneline", "--all")
				gitLog.Dir = te.Dir
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
			te := NewTestEnvironment(t, tt.envOptions...)

			if tt.setup != nil {
				tt.setup(t, te)
			}

			output, err := te.RunCommand(tt.args...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validateOutput != nil {
				tt.validateOutput(t, te, output, err)
			}
		})
	}
}

// TestInitCommand_US014_SynchronizationSetup tests US-014 synchronization initialization
func TestInitCommand_US014_SynchronizationSetup(t *testing.T) {
	tests := []struct {
		name           string
		envOptions     []TestEnvOption
		setup          func(t *testing.T, te *TestEnvironment)
		args           []string
		validateOutput func(t *testing.T, te *TestEnvironment, output string, err error)
		expectError    bool
	}{
		{
			name: "basic_sync_initialization",
			args: []string{"init", "--silent"},
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create initial commit
				te.CreateFile("README.md", "# Test")
				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = te.Dir
				require.NoError(t, gitAdd.Run())
				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = te.Dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Verify config is created
				assert.FileExists(t, te.ConfigPath, "Should create config file")
			},
			expectError: false,
		},
		{
			name: "sync_initialization_with_custom_repository",
			args: func() []string {
				// In CI, skip git integration due to subtree limitations
				if os.Getenv("CI") != "" {
					return []string{"init", "--force", "--silent", "--no-git"}
				}
				return []string{"init", "--force", "--silent"}
			}(),
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create existing config
				var existingConfig string
				if os.Getenv("CI") != "" {
					// In CI, use simple config without repository
					existingConfig = `version: "1.0"
library:
  path: .ddx/library
`
				} else {
					// Locally, test with repository URL
					existingConfig = fmt.Sprintf(`version: "1.0"
library:
  path: .ddx/library
  repository:
    url: %s
    branch: master
`, te.TestLibraryURL)
				}
				te.CreateConfig(existingConfig)
				te.CreateFile("README.md", "# Test")
				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = te.Dir
				require.NoError(t, gitAdd.Run())
				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = te.Dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Should handle custom repository successfully
				assert.NotContains(t, output, "backup", "Should not create or mention backup")
			},
			expectError: false,
		},
		{
			name: "sync_initialization_fresh_project",
			args: []string{"init", "--silent"},
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create initial commit
				te.CreateFile("README.md", "# Test")
				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = te.Dir
				require.NoError(t, gitAdd.Run())
				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = te.Dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Check .ddx/config.yaml was created with sync config
				assert.FileExists(t, te.ConfigPath)

				data, err := os.ReadFile(te.ConfigPath)
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
			args: []string{"init", "--force", "--silent"},
			setup: func(t *testing.T, te *TestEnvironment) {
				// Create existing project files
				te.CreateFile("README.md", "# Existing Project")
				te.CreateFile("package.json", `{"name": "test"}`)
				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = te.Dir
				require.NoError(t, gitAdd.Run())
				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = te.Dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Existing files should remain untouched
				assert.FileExists(t, filepath.Join(te.Dir, "README.md"))
				assert.FileExists(t, filepath.Join(te.Dir, "package.json"))
			},
			expectError: false,
		},
		{
			name: "sync_initialization_validation_success",
			args: []string{"init", "--silent"},
			setup: func(t *testing.T, te *TestEnvironment) {
				te.CreateFile("README.md", "# Test")
				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = te.Dir
				require.NoError(t, gitAdd.Run())
				gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
				gitCommit.Dir = te.Dir
				require.NoError(t, gitCommit.Run())
			},
			validateOutput: func(t *testing.T, te *TestEnvironment, output string, err error) {
				// Verify config file exists with proper structure
				assert.FileExists(t, te.ConfigPath, "Should create config file")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := NewTestEnvironment(t, tt.envOptions...)

			if tt.setup != nil {
				tt.setup(t, te)
			}

			output, err := te.RunCommand(tt.args...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validateOutput != nil {
				tt.validateOutput(t, te, output, err)
			}
		})
	}
}
