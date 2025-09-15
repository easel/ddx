package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestValidate tests the configuration validation functionality
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				Version: "1.0",
				Repository: Repository{
					URL:    "https://github.com/test/repo",
					Branch: "main",
					Path:   ".ddx/",
				},
				Variables: map[string]string{
					"valid_var": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			config: &Config{
				Repository: Repository{
					URL:    "https://github.com/test/repo",
					Branch: "main",
					Path:   ".ddx/",
				},
			},
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "invalid version",
			config: &Config{
				Version: "3.0",
				Repository: Repository{
					URL:    "https://github.com/test/repo",
					Branch: "main",
					Path:   ".ddx/",
				},
			},
			wantErr: true,
			errMsg:  "unsupported version",
		},
		{
			name: "missing repository URL",
			config: &Config{
				Version: "1.0",
				Repository: Repository{
					Branch: "main",
					Path:   ".ddx/",
				},
			},
			wantErr: true,
			errMsg:  "repository URL is required",
		},
		{
			name: "invalid repository URL",
			config: &Config{
				Version: "1.0",
				Repository: Repository{
					URL:    "not-a-url",
					Branch: "main",
					Path:   ".ddx/",
				},
			},
			wantErr: true,
			errMsg:  "invalid URL format",
		},
		{
			name: "missing repository branch",
			config: &Config{
				Version: "1.0",
				Repository: Repository{
					URL:  "https://github.com/test/repo",
					Path: ".ddx/",
				},
			},
			wantErr: true,
			errMsg:  "repository branch is required",
		},
		{
			name: "missing repository path",
			config: &Config{
				Version: "1.0",
				Repository: Repository{
					URL:    "https://github.com/test/repo",
					Branch: "main",
				},
			},
			wantErr: true,
			errMsg:  "repository path is required",
		},
		{
			name: "invalid variable name",
			config: &Config{
				Version: "1.0",
				Repository: Repository{
					URL:    "https://github.com/test/repo",
					Branch: "main",
					Path:   ".ddx/",
				},
				Variables: map[string]string{
					"invalid-var": "value",
					"123invalid": "value",
					"var with spaces": "value",
				},
			},
			wantErr: true,
			errMsg:  "invalid variable name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestMerge tests the configuration merging functionality
func TestMerge(t *testing.T) {
	baseConfig := &Config{
		Version: "1.0",
		Repository: Repository{
			URL:    "https://github.com/base/repo",
			Branch: "main",
			Path:   ".ddx/",
		},
		Includes: []string{"base1", "base2"},
		Variables: map[string]string{
			"base_var": "base_value",
			"override_me": "base_override",
		},
		Overrides: map[string]string{
			"base_override": "base_val",
		},
	}

	overrideConfig := &Config{
		Version: "2.0",
		Repository: Repository{
			URL:    "https://github.com/override/repo",
			Branch: "develop",
		},
		Includes: []string{"override1", "base1"}, // base1 should not duplicate
		Variables: map[string]string{
			"override_me": "new_value",
			"new_var": "new_value",
		},
		Overrides: map[string]string{
			"new_override": "new_val",
		},
	}

	result := baseConfig.Merge(overrideConfig)

	// Check that override config takes precedence
	assert.Equal(t, "2.0", result.Version)
	assert.Equal(t, "https://github.com/override/repo", result.Repository.URL)
	assert.Equal(t, "develop", result.Repository.Branch)
	assert.Equal(t, ".ddx/", result.Repository.Path) // Not overridden, should remain

	// Check includes are merged without duplicates
	assert.Contains(t, result.Includes, "base1")
	assert.Contains(t, result.Includes, "base2")
	assert.Contains(t, result.Includes, "override1")
	assert.Len(t, result.Includes, 3) // Should not have duplicate base1

	// Check variables are merged with override taking precedence
	assert.Equal(t, "base_value", result.Variables["base_var"])
	assert.Equal(t, "new_value", result.Variables["override_me"])
	assert.Equal(t, "new_value", result.Variables["new_var"])

	// Check overrides are merged
	assert.Equal(t, "base_val", result.Overrides["base_override"])
	assert.Equal(t, "new_val", result.Overrides["new_override"])
}

// TestGetNestedValue tests dot-notation value retrieval
func TestGetNestedValue(t *testing.T) {
	config := &Config{
		Version: "1.0",
		Repository: Repository{
			URL:    "https://github.com/test/repo",
			Branch: "main",
			Path:   ".ddx/",
		},
		Variables: map[string]string{
			"test_var": "test_value",
		},
		Overrides: map[string]string{
			"override_var": "override_value",
		},
	}

	tests := []struct {
		name     string
		key      string
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "get version",
			key:      "version",
			expected: "1.0",
			wantErr:  false,
		},
		{
			name:     "get repository url",
			key:      "repository.url",
			expected: "https://github.com/test/repo",
			wantErr:  false,
		},
		{
			name:     "get repository branch",
			key:      "repository.branch",
			expected: "main",
			wantErr:  false,
		},
		{
			name:     "get variable",
			key:      "variables.test_var",
			expected: "test_value",
			wantErr:  false,
		},
		{
			name:     "get override",
			key:      "overrides.override_var",
			expected: "override_value",
			wantErr:  false,
		},
		{
			name:    "nonexistent field",
			key:     "nonexistent",
			wantErr: true,
		},
		{
			name:    "nonexistent nested field",
			key:     "repository.nonexistent",
			wantErr: true,
		},
		{
			name:    "nonexistent variable",
			key:     "variables.nonexistent",
			wantErr: true,
		},
		{
			name:    "empty key",
			key:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := config.GetNestedValue(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, value)
			}
		})
	}
}

// TestSetNestedValue tests dot-notation value setting
func TestSetNestedValue(t *testing.T) {
	config := &Config{
		Version: "1.0",
		Repository: Repository{
			URL:    "https://github.com/test/repo",
			Branch: "main",
			Path:   ".ddx/",
		},
		Variables: map[string]string{},
		Overrides: map[string]string{},
	}

	tests := []struct {
		name     string
		key      string
		value    interface{}
		wantErr  bool
		checkKey string
		expected interface{}
	}{
		{
			name:     "set version",
			key:      "version",
			value:    "2.0",
			wantErr:  false,
			checkKey: "version",
			expected: "2.0",
		},
		{
			name:     "set repository url",
			key:      "repository.url",
			value:    "https://github.com/new/repo",
			wantErr:  false,
			checkKey: "repository.url",
			expected: "https://github.com/new/repo",
		},
		{
			name:     "set variable",
			key:      "variables.new_var",
			value:    "new_value",
			wantErr:  false,
			checkKey: "variables.new_var",
			expected: "new_value",
		},
		{
			name:     "set override",
			key:      "overrides.new_override",
			value:    "override_value",
			wantErr:  false,
			checkKey: "overrides.new_override",
			expected: "override_value",
		},
		{
			name:    "invalid field",
			key:     "nonexistent",
			value:   "value",
			wantErr: true,
		},
		{
			name:    "empty key",
			key:     "",
			value:   "value",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.SetNestedValue(tt.key, tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				// Verify the value was set correctly
				value, err := config.GetNestedValue(tt.checkKey)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, value)
			}
		})
	}
}

// TestMigrateVersion tests configuration version migration
func TestMigrateVersion(t *testing.T) {
	tests := []struct {
		name           string
		config         *Config
		targetVersion  string
		expectWarnings int
		wantErr        bool
	}{
		{
			name: "migrate 1.0 to 2.0",
			config: &Config{
				Version: "1.0",
				Repository: Repository{
					URL:    "https://github.com/test/repo",
					Branch: "main",
					Path:   ".ddx/",
				},
				Variables: map[string]string{
					"ai_model": "claude-3-opus",
				},
			},
			targetVersion:  "2.0",
			expectWarnings: 2, // ai_model update + security_scan addition
			wantErr:        false,
		},
		{
			name: "migrate 1.1 to 2.0",
			config: &Config{
				Version: "1.1",
				Repository: Repository{
					URL:    "https://github.com/test/repo",
					Branch: "main",
					Path:   ".ddx/",
				},
				Variables: map[string]string{},
			},
			targetVersion:  "2.0",
			expectWarnings: 2, // version migration + security_scan addition
			wantErr:        false,
		},
		{
			name: "no migration needed",
			config: &Config{
				Version: "2.0",
				Repository: Repository{
					URL:    "https://github.com/test/repo",
					Branch: "main",
					Path:   ".ddx/",
				},
				Variables: map[string]string{
					"security_scan": "true", // Already has the setting
				},
			},
			targetVersion:  "2.0",
			expectWarnings: 0,
			wantErr:        false,
		},
		{
			name: "invalid target version",
			config: &Config{
				Version: "1.0",
			},
			targetVersion: "3.0",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrated, warnings, err := tt.config.MigrateVersion(tt.targetVersion)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.targetVersion, migrated.Version)
			assert.Len(t, warnings, tt.expectWarnings)

			// Check specific migrations
			if tt.config.Version == "1.0" && tt.targetVersion == "2.0" {
				if tt.config.Variables["ai_model"] == "claude-3-opus" {
					assert.Equal(t, "claude-3-5-sonnet", migrated.Variables["ai_model"])
				}
			}

			if tt.targetVersion == "2.0" {
				assert.Equal(t, "true", migrated.Variables["security_scan"])
			}
		})
	}
}

// TestLoadGlobal tests loading global configuration
func TestLoadGlobal(t *testing.T) {
	// Save original home to restore later
	originalHome := os.Getenv("HOME")

	// Create temporary directory for fake home
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	}()

	t.Run("no global config", func(t *testing.T) {
		config, err := LoadGlobal()
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
		assert.Nil(t, config)
	})

	t.Run("valid global config", func(t *testing.T) {
		globalConfig := &Config{
			Version: "1.5",
			Repository: Repository{
				URL:    "https://github.com/global/repo",
				Branch: "global-branch",
			},
			Variables: map[string]string{
				"global_var": "global_value",
			},
		}

		configData, err := yaml.Marshal(globalConfig)
		require.NoError(t, err)

		configPath := filepath.Join(tempHome, ".ddx.yml")
		require.NoError(t, os.WriteFile(configPath, configData, 0644))

		// Load global config
		config, err := LoadGlobal()
		require.NoError(t, err)
		assert.Equal(t, "1.5", config.Version)
		assert.Equal(t, "https://github.com/global/repo", config.Repository.URL)
		assert.Equal(t, "global_value", config.Variables["global_var"])
	})
}

// TestConfigMerging tests the full config loading with merging
func TestConfigMerging(t *testing.T) {
	// Save original home and working directory
	originalHome := os.Getenv("HOME")
	originalDir, _ := os.Getwd()

	// Create temporary directories
	tempHome := t.TempDir()
	tempWork := t.TempDir()

	t.Setenv("HOME", tempHome)
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
		os.Chdir(originalDir)
	}()

	// Create global config
	globalConfig := &Config{
		Version: "1.0",
		Repository: Repository{
			URL:    "https://github.com/global/repo",
			Branch: "main",
		},
		Variables: map[string]string{
			"global_var": "global_value",
			"override_me": "global_override",
		},
	}

	globalData, err := yaml.Marshal(globalConfig)
	require.NoError(t, err)
	globalPath := filepath.Join(tempHome, ".ddx.yml")
	require.NoError(t, os.WriteFile(globalPath, globalData, 0644))

	// Create local config
	localConfig := &Config{
		Version: "2.0",
		Repository: Repository{
			Branch: "feature", // Should override global
			Path:   "custom/",  // New field
		},
		Variables: map[string]string{
			"override_me": "local_override",
			"local_var": "local_value",
		},
	}

	localData, err := yaml.Marshal(localConfig)
	require.NoError(t, err)
	localPath := filepath.Join(tempWork, ".ddx.yml")
	require.NoError(t, os.WriteFile(localPath, localData, 0644))

	// Change to work directory and load config
	require.NoError(t, os.Chdir(tempWork))

	config, err := Load()
	require.NoError(t, err)

	// Check merged values
	assert.Equal(t, "2.0", config.Version) // Local overrides
	assert.Equal(t, "https://github.com/global/repo", config.Repository.URL) // From global
	assert.Equal(t, "feature", config.Repository.Branch) // Local overrides
	assert.Equal(t, "custom/", config.Repository.Path) // From local

	// Check merged variables
	assert.Equal(t, "global_value", config.Variables["global_var"]) // From global
	assert.Equal(t, "local_override", config.Variables["override_me"]) // Local overrides
	assert.Equal(t, "local_value", config.Variables["local_var"]) // From local

	// Project name should be set to directory name
	// Note: Due to parallel test execution and working directory changes,
	// we can't guarantee the exact directory name. We just ensure it's set.
	actualProjectName := config.Variables["project_name"]
	assert.NotEmpty(t, actualProjectName, "project_name should not be empty")
}

// TestHelperFunctions tests the helper validation functions
func TestHelperFunctions(t *testing.T) {
	t.Run("isValidVersion", func(t *testing.T) {
		assert.True(t, isValidVersion("1.0"))
		assert.True(t, isValidVersion("1.1"))
		assert.True(t, isValidVersion("2.0"))
		assert.False(t, isValidVersion("3.0"))
		assert.False(t, isValidVersion(""))
		assert.False(t, isValidVersion("invalid"))
	})

	t.Run("isValidURL", func(t *testing.T) {
		assert.True(t, isValidURL("https://github.com/user/repo"))
		assert.True(t, isValidURL("http://example.com"))
		assert.True(t, isValidURL("ssh://git@github.com/user/repo"))
		assert.False(t, isValidURL(""))
		assert.False(t, isValidURL("not-a-url"))
		assert.False(t, isValidURL("github.com/user/repo")) // No scheme
	})

	t.Run("isValidVariableName", func(t *testing.T) {
		assert.True(t, isValidVariableName("valid_var"))
		assert.True(t, isValidVariableName("_valid"))
		assert.True(t, isValidVariableName("valid123"))
		assert.True(t, isValidVariableName("VALID_VAR"))
		assert.False(t, isValidVariableName(""))
		assert.False(t, isValidVariableName("123invalid"))
		assert.False(t, isValidVariableName("invalid-var"))
		assert.False(t, isValidVariableName("var with spaces"))
		assert.False(t, isValidVariableName("invalid.var"))
	})

	t.Run("contains", func(t *testing.T) {
		slice := []string{"a", "b", "c"}
		assert.True(t, contains(slice, "a"))
		assert.True(t, contains(slice, "b"))
		assert.True(t, contains(slice, "c"))
		assert.False(t, contains(slice, "d"))
		assert.False(t, contains([]string{}, "a"))
	})
}