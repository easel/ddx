package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigResourcesDebug tests if config loading works correctly
func TestConfigResourcesDebug(t *testing.T) {
	// Use TestEnvironment for isolated testing
	env := NewTestEnvironment(t)

	// Create simple config (without resources since we don't support them in new format)
	configContent := `version: "1.0"
library_base_path: "./library"
repository:
  url: "https://github.com/test/repo"
  branch: "main"
  subtree_prefix: "library"
variables:
  project_name: "test-project"
`
	env.CreateConfig(configContent)

	// Load config and verify it works
	cfg, err := env.LoadConfig()
	require.NoError(t, err, "Config should load successfully")
	require.NotNil(t, cfg, "Config should not be nil")

	t.Logf("Config version: %s", cfg.Version)
	t.Logf("Config library base path: %s", cfg.LibraryBasePath)
	t.Logf("Config variables: %+v", cfg.Variables)

	// Verify basic config is loaded correctly
	assert.Equal(t, "1.0", cfg.Version, "Version should be loaded")
	assert.Equal(t, "./library", cfg.LibraryBasePath, "Library base path should be loaded")
	assert.NotNil(t, cfg.Variables, "Variables should be loaded")
	assert.Equal(t, "test-project", cfg.Variables["project_name"], "Project name should be loaded")
}

// TestBasicConfigRepository tests repository configuration
func TestBasicConfigRepository(t *testing.T) {
	// Use TestEnvironment for isolated testing
	env := NewTestEnvironment(t)

	// Create config with repository settings
	configContent := `version: "1.0"
library_base_path: "./library"
repository:
  url: "https://github.com/ddx-tools/ddx"
  branch: "main"
  subtree_prefix: "library"
variables:
  project_name: "ddx-test"
`
	env.CreateConfig(configContent)

	// Load config and verify repository settings
	cfg, err := env.LoadConfig()
	require.NoError(t, err, "Config should load successfully")
	require.NotNil(t, cfg, "Config should not be nil")
	require.NotNil(t, cfg.Repository, "Repository should be loaded")

	// Verify repository configuration
	assert.Equal(t, "https://github.com/ddx-tools/ddx", cfg.Repository.URL, "Repository URL should be loaded")
	assert.Equal(t, "main", cfg.Repository.Branch, "Repository branch should be loaded")
	assert.Equal(t, "library", cfg.Repository.SubtreePrefix, "Repository subtree prefix should be loaded")
}
