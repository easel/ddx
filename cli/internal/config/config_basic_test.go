package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestDefaultConfig validates the default configuration values
func TestDefaultConfig_Basic(t *testing.T) {
	t.Parallel()
	config := DefaultConfig

	assert.Equal(t, "2.0", config.Version)
	assert.Equal(t, "https://github.com/easel/ddx", config.Repository.URL)
	assert.Equal(t, "main", config.Repository.Branch)
	assert.Equal(t, ".ddx/", config.Repository.Path)
	assert.Empty(t, config.Includes)
	assert.Empty(t, config.Variables)
}

// TestLoadConfig_DefaultOnly tests loading when no config files exist
func TestLoadConfig_DefaultOnly_Basic(t *testing.T) {
	// Create temp directory without config files
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(tempDir))

	config, err := Load()

	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, DefaultConfig.Version, config.Version)
	assert.Equal(t, DefaultConfig.Repository.URL, config.Repository.URL)
	// Check that project_name was set from directory
	assert.NotEmpty(t, config.Variables["project_name"])
}

// TestLoadConfig_LocalConfig tests loading with local .ddx.yml
func TestLoadConfig_LocalConfig_Basic(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Create local config
	localConfig := &Config{
		Version: "2.0",
		Repository: Repository{
			URL:    "https://github.com/custom/repo",
			Branch: "develop",
			Path:   "custom/",
		},
		Includes: []string{"custom/templates"},
		Variables: map[string]string{
			"project_name": "test-project",
		},
	}

	configData, err := yaml.Marshal(localConfig)
	require.NoError(t, err)

	configPath := filepath.Join(tempDir, ".ddx.yml")
	require.NoError(t, os.WriteFile(configPath, configData, 0644))

	require.NoError(t, os.Chdir(tempDir))

	// Load config
	config, err := Load()

	require.NoError(t, err)
	assert.Equal(t, "2.0", config.Version)
	assert.Equal(t, "https://github.com/custom/repo", config.Repository.URL)
	assert.Equal(t, "develop", config.Repository.Branch)
	assert.Contains(t, config.Variables, "project_name")
}

// TestLoadLocal tests LoadLocal function
func TestLoadLocal_Basic(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Create local config
	localConfig := &Config{
		Version: "1.5",
		Repository: Repository{
			URL:    "https://github.com/local/repo",
			Branch: "feature",
		},
		Variables: map[string]string{
			"test_var": "test_value",
		},
	}

	configData, err := yaml.Marshal(localConfig)
	require.NoError(t, err)

	configPath := filepath.Join(tempDir, ".ddx.yml")
	require.NoError(t, os.WriteFile(configPath, configData, 0644))

	require.NoError(t, os.Chdir(tempDir))

	// Load local config
	config, err := LoadLocal()

	require.NoError(t, err)
	assert.Equal(t, "1.5", config.Version)
	assert.Equal(t, "https://github.com/local/repo", config.Repository.URL)
	assert.Equal(t, "test_value", config.Variables["test_var"])
}

// TestSaveLocal tests SaveLocal function
func TestSaveLocal_Basic(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(tempDir))

	config := &Config{
		Version: "1.0",
		Repository: Repository{
			URL:    "https://github.com/test/repo",
			Branch: "main",
			Path:   ".ddx/",
		},
		Variables: map[string]string{
			"key1": "value1",
		},
	}

	// Save config locally
	err := SaveLocal(config)
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, ".ddx.yml")

	// Load and verify
	loadedConfig, err := LoadLocal()
	require.NoError(t, err)

	assert.Equal(t, config.Version, loadedConfig.Version)
	assert.Equal(t, config.Repository.URL, loadedConfig.Repository.URL)
	assert.Equal(t, "value1", loadedConfig.Variables["key1"])
}

// TestReplaceVariables tests the ReplaceVariables method
func TestReplaceVariables_Basic(t *testing.T) {
	t.Parallel()
	config := &Config{
		Variables: map[string]string{
			"name":    "TestProject",
			"version": "1.0.0",
			"author":  "Test Author",
		},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple replacement",
			input:    "Project: {{name}}",
			expected: "Project: TestProject",
		},
		{
			name:     "multiple replacements",
			input:    "{{name}} v{{version}} by {{author}}",
			expected: "TestProject v1.0.0 by Test Author",
		},
		{
			name:     "with spaces",
			input:    "Name: {{ name }}",
			expected: "Name: TestProject",
		},
		{
			name:     "non-existent variable",
			input:    "Unknown: {{unknown}}",
			expected: "Unknown: {{unknown}}", // Should remain unchanged
		},
		{
			name:     "no variables",
			input:    "Plain text",
			expected: "Plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.ReplaceVariables(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestLoadConfig_InvalidYAML tests handling of invalid YAML
func TestLoadConfig_InvalidYAML_Basic(t *testing.T) {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Create invalid YAML file
	invalidYAML := `
version: 1.0
repository:
  url: https://github.com/test
  branch: [this is invalid
`
	configPath := filepath.Join(tempDir, ".ddx.yml")
	require.NoError(t, os.WriteFile(configPath, []byte(invalidYAML), 0644))

	require.NoError(t, os.Chdir(tempDir))

	// Should return error
	config, err := LoadLocal()

	assert.Error(t, err)
	assert.Nil(t, config)
}
