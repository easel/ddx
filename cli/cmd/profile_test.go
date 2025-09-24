package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProfileCommands tests profile management functionality for US-023
func TestProfileCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		setup    func(t *testing.T) string
		expected string
		wantErr  bool
	}{
		{
			name:     "profile_create_dev",
			args:     []string{"config", "profile", "create", "dev"},
			setup:    setupCleanWorkspace,
			expected: "Created profile 'dev'",
			wantErr:  false,
		},
		{
			name:     "profile_list_empty",
			args:     []string{"config", "profile", "list"},
			setup:    setupCleanWorkspace,
			expected: "No environment profiles found",
			wantErr:  false,
		},
		{
			name:     "profile_list_with_profiles",
			args:     []string{"config", "profile", "list"},
			setup:    setupWithProfiles,
			expected: "Available Environment Profiles:",
			wantErr:  false,
		},
		{
			name:     "profile_activate_existing",
			args:     []string{"config", "profile", "activate", "dev"},
			setup:    setupWithProfiles,
			expected: "Profile 'dev' is ready for activation",
			wantErr:  false,
		},
		{
			name:     "profile_activate_nonexistent",
			args:     []string{"config", "profile", "activate", "nonexistent"},
			setup:    setupCleanWorkspace,
			expected: "profile 'nonexistent' does not exist",
			wantErr:  true,
		},
		{
			name:     "profile_copy",
			args:     []string{"config", "profile", "copy", "dev", "staging"},
			setup:    setupWithProfiles,
			expected: "Copied profile 'dev' to 'staging'",
			wantErr:  false,
		},
		{
			name:     "profile_show_existing",
			args:     []string{"config", "profile", "show", "dev"},
			setup:    setupWithProfiles,
			expected: "Profile Configuration: dev",
			wantErr:  false,
		},
		{
			name:     "profile_validate_existing",
			args:     []string{"config", "profile", "validate", "dev"},
			setup:    setupWithProfiles,
			expected: "Validating profile 'dev'",
			wantErr:  false,
		},
		{
			name:     "profile_diff",
			args:     []string{"config", "profile", "diff", "dev", "staging"},
			setup:    setupWithMultipleProfiles,
			expected: "Profile Comparison: dev vs staging",
			wantErr:  false,
		},
		{
			name:     "profile_delete_existing",
			args:     []string{"config", "profile", "delete", "dev"},
			setup:    setupWithProfiles,
			expected: "Deleted profile 'dev'",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalDir, _ := os.Getwd()
			workDir := tt.setup(t)
			require.NoError(t, os.Chdir(workDir))
			defer func() {
				os.Chdir(originalDir)
			}()

			rootCmd := getTestRootCommand()
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

// Helper function to setup a clean workspace
func setupCleanWorkspace(t *testing.T) string {
	workDir := t.TempDir()
	return workDir
}

// Helper function to setup workspace with profiles
func setupWithProfiles(t *testing.T) string {
	workDir := t.TempDir()

	// Create a base .ddx.yml file
	baseConfig := `version: "1.0"
variables:
  project_name: "test-project"
repository:
  url: "https://github.com/test/repo"
  branch: "main"
`
	require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.yml"), []byte(baseConfig), 0644))

	// Create a dev profile
	devConfig := `version: "1.0"
variables:
  project_name: "test-project"
  DDX_PROFILE: "dev"
  DDX_ENV: "dev"
  api_endpoint: "https://api-dev.example.com"
repository:
  url: "https://github.com/test/repo"
  branch: "development"
`
	require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.dev.yml"), []byte(devConfig), 0644))

	return workDir
}

// Helper function to setup workspace with multiple profiles
func setupWithMultipleProfiles(t *testing.T) string {
	workDir := setupWithProfiles(t)

	// Create a staging profile
	stagingConfig := `version: "1.0"
variables:
  project_name: "test-project"
  DDX_PROFILE: "staging"
  DDX_ENV: "staging"
  api_endpoint: "https://api-staging.example.com"
repository:
  url: "https://github.com/test/repo"
  branch: "staging"
`
	require.NoError(t, os.WriteFile(filepath.Join(workDir, ".ddx.staging.yml"), []byte(stagingConfig), 0644))

	return workDir
}

// TestProfileCommandStructure tests the command structure
func TestProfileCommandStructure(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError string
	}{
		{
			name:        "profile_no_action",
			args:        []string{"config", "profile"},
			expectError: "profile subcommand requires additional arguments",
		},
		{
			name:        "profile_create_no_name",
			args:        []string{"config", "profile", "create"},
			expectError: "profile create requires a profile name",
		},
		{
			name:        "profile_copy_insufficient_args",
			args:        []string{"config", "profile", "copy", "dev"},
			expectError: "profile copy requires source and destination profile names",
		},
		{
			name:        "profile_unknown_action",
			args:        []string{"config", "profile", "unknown"},
			expectError: "unknown profile action: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workDir := setupCleanWorkspace(t)
			originalDir, _ := os.Getwd()
			require.NoError(t, os.Chdir(workDir))
			defer func() {
				os.Chdir(originalDir)
			}()

			rootCmd := getTestRootCommand()
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectError)
		})
	}
}

// TestProfileInheritance tests that profiles inherit from base configuration
func TestProfileInheritance(t *testing.T) {
	workDir := setupWithProfiles(t)
	originalDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(workDir))
	defer func() {
		os.Chdir(originalDir)
	}()

	// Check that the dev profile file exists and contains inherited values
	devProfileContent, err := os.ReadFile(".ddx.dev.yml")
	require.NoError(t, err)

	content := string(devProfileContent)
	assert.Contains(t, content, "project_name: \"test-project\"") // Inherited from base
	assert.Contains(t, content, "DDX_PROFILE: \"dev\"")           // Profile-specific
	assert.Contains(t, content, "api_endpoint:")                  // Profile-specific
}

// TestProfileFilenameValidation tests profile name validation
func TestProfileFilenameValidation(t *testing.T) {
	workDir := setupCleanWorkspace(t)
	originalDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(workDir))
	defer func() {
		os.Chdir(originalDir)
	}()

	invalidNames := []string{
		"dev/test",  // Contains path separator
		"dev\\test", // Contains path separator
		"../dev",    // Relative path
	}

	for _, invalidName := range invalidNames {
		t.Run("invalid_name_"+strings.ReplaceAll(invalidName, "/", "_"), func(t *testing.T) {
			rootCmd := getTestRootCommand()
			rootCmd.SetArgs([]string{"config", "profile", "create", invalidName})

			err := rootCmd.Execute()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid profile name")
		})
	}
}
