package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUS019_EnvironmentConfigFiles tests US-019: Override Configuration
func TestUS019_EnvironmentConfigFiles(t *testing.T) {
	// Save original directory and restore at end
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	t.Run("environment_specific_config_files_supported", func(t *testing.T) {
		// AC: Given environment configs, when present, then .ddx.dev.yml, .ddx.staging.yml are supported

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Create base config
		baseConfig := `version: "2.0"
repository:
  url: "https://github.com/base/repo"
  branch: "main"
variables:
  api_url: "http://localhost:3000"
  log_level: "info"
  cache_enabled: "false"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		// Create development override
		devConfig := `variables:
  log_level: "debug"
  cache_enabled: "false"
`
		require.NoError(t, os.WriteFile(".ddx.dev.yml", []byte(devConfig), 0644))

		// Create staging override
		stagingConfig := `variables:
  api_url: "https://api.staging.com"
  log_level: "warn"
`
		require.NoError(t, os.WriteFile(".ddx.staging.yml", []byte(stagingConfig), 0644))

		// Test dev environment
		t.Setenv("DDX_ENV", "dev")
		config, err := Load()
		require.NoError(t, err)
		assert.Equal(t, "debug", config.Variables["log_level"], "Dev config should override log_level")
		assert.Equal(t, "http://localhost:3000", config.Variables["api_url"], "Dev config should preserve base api_url")

		// Test staging environment
		t.Setenv("DDX_ENV", "staging")
		config, err = Load()
		require.NoError(t, err)
		assert.Equal(t, "warn", config.Variables["log_level"], "Staging config should override log_level")
		assert.Equal(t, "https://api.staging.com", config.Variables["api_url"], "Staging config should override api_url")
	})

	t.Run("ddx_env_variable_selects_override_file", func(t *testing.T) {
		// AC: Given DDX_ENV variable, when set, then appropriate override file is selected

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		baseConfig := `version: "2.0"
variables:
  env_name: "base"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		// Create different environment overrides
		prodConfig := `variables:
  env_name: "production"
  api_url: "https://api.production.com"
`
		require.NoError(t, os.WriteFile(".ddx.prod.yml", []byte(prodConfig), 0644))

		testConfig := `variables:
  env_name: "test"
  debug: "true"
`
		require.NoError(t, os.WriteFile(".ddx.test.yml", []byte(testConfig), 0644))

		// Test production environment
		t.Setenv("DDX_ENV", "prod")
		config, err := Load()
		require.NoError(t, err)
		assert.Equal(t, "production", config.Variables["env_name"])
		assert.Equal(t, "https://api.production.com", config.Variables["api_url"])

		// Test test environment
		t.Setenv("DDX_ENV", "test")
		config, err = Load()
		require.NoError(t, err)
		assert.Equal(t, "test", config.Variables["env_name"])
		assert.Equal(t, "true", config.Variables["debug"])

		// Test no environment (should use base)
		os.Unsetenv("DDX_ENV")
		config, err = Load()
		require.NoError(t, err)
		assert.Equal(t, "base", config.Variables["env_name"])
		assert.NotContains(t, config.Variables, "api_url")
		assert.NotContains(t, config.Variables, "debug")
	})

	t.Run("overrides_merge_correctly_with_base", func(t *testing.T) {
		// AC: Given overrides, when loaded, then they merge correctly with base configuration

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		baseConfig := `version: "2.0"
repository:
  url: "https://github.com/base/repo"
  branch: "main"
  path: ".ddx/"
includes:
  - "base1"
  - "base2"
variables:
  api_url: "http://localhost:3000"
  log_level: "info"
  cache_enabled: "false"
  base_only: "value"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		envConfig := `repository:
  branch: "develop"
includes:
  - "env1"
variables:
  log_level: "debug"
  env_only: "env_value"
`
		require.NoError(t, os.WriteFile(".ddx.dev.yml", []byte(envConfig), 0644))

		t.Setenv("DDX_ENV", "dev")
		config, err := Load()
		require.NoError(t, err)

		// Repository should merge
		assert.Equal(t, "https://github.com/base/repo", config.Repository.URL, "Base URL preserved")
		assert.Equal(t, "develop", config.Repository.Branch, "Override branch applied")
		assert.Equal(t, ".ddx/", config.Repository.Path, "Base path preserved")

		// Includes should merge
		assert.Contains(t, config.Includes, "base1")
		assert.Contains(t, config.Includes, "base2")
		assert.Contains(t, config.Includes, "env1")

		// Variables should merge with override precedence
		assert.Equal(t, "http://localhost:3000", config.Variables["api_url"], "Base value preserved")
		assert.Equal(t, "debug", config.Variables["log_level"], "Override value applied")
		assert.Equal(t, "false", config.Variables["cache_enabled"], "Base value preserved")
		assert.Equal(t, "value", config.Variables["base_only"], "Base-only value preserved")
		assert.Equal(t, "env_value", config.Variables["env_only"], "Override-only value added")
	})

	t.Run("override_takes_precedence_over_base", func(t *testing.T) {
		// AC: Given any config value, when overridden, then override takes precedence

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		baseConfig := `version: "1.0"
repository:
  url: "https://github.com/base/repo"
  branch: "main"
variables:
  all_override: "base_value"
  partial_override: "base_value"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		overrideConfig := `version: "2.0"
repository:
  url: "https://github.com/override/repo"
  branch: "feature"
variables:
  all_override: "override_value"
  partial_override: "override_value"
`
		require.NoError(t, os.WriteFile(".ddx.prod.yml", []byte(overrideConfig), 0644))

		t.Setenv("DDX_ENV", "prod")
		config, err := Load()
		require.NoError(t, err)

		// All overridden values should use override
		assert.Equal(t, "2.0", config.Version)
		assert.Equal(t, "https://github.com/override/repo", config.Repository.URL)
		assert.Equal(t, "feature", config.Repository.Branch)
		assert.Equal(t, "override_value", config.Variables["all_override"])
		assert.Equal(t, "override_value", config.Variables["partial_override"])
	})

	t.Run("partial_configs_only_change_specified_values", func(t *testing.T) {
		// AC: Given partial configs, when used as overrides, then only specified values are changed

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		baseConfig := `version: "2.0"
repository:
  url: "https://github.com/base/repo"
  branch: "main"
  path: ".ddx/"
variables:
  var1: "base1"
  var2: "base2"
  var3: "base3"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		// Partial override - only changes var2
		partialConfig := `variables:
  var2: "overridden2"
`
		require.NoError(t, os.WriteFile(".ddx.partial.yml", []byte(partialConfig), 0644))

		t.Setenv("DDX_ENV", "partial")
		config, err := Load()
		require.NoError(t, err)

		// Only var2 should be overridden
		assert.Equal(t, "2.0", config.Version, "Version should remain from base")
		assert.Equal(t, "https://github.com/base/repo", config.Repository.URL, "URL should remain from base")
		assert.Equal(t, "main", config.Repository.Branch, "Branch should remain from base")
		assert.Equal(t, ".ddx/", config.Repository.Path, "Path should remain from base")
		assert.Equal(t, "base1", config.Variables["var1"], "var1 should remain from base")
		assert.Equal(t, "overridden2", config.Variables["var2"], "var2 should be overridden")
		assert.Equal(t, "base3", config.Variables["var3"], "var3 should remain from base")
	})

	t.Run("missing_environment_override_file", func(t *testing.T) {
		// Edge case: Missing environment override file when DDX_ENV is set

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		baseConfig := `version: "2.0"
variables:
  env: "base"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		// Set environment but don't create the file
		t.Setenv("DDX_ENV", "nonexistent")
		config, err := Load()

		// Should load successfully with just base config
		require.NoError(t, err, "Should handle missing environment file gracefully")
		assert.Equal(t, "base", config.Variables["env"])
	})

	t.Run("empty_override_file", func(t *testing.T) {
		// Edge case: Empty override files

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		baseConfig := `version: "2.0"
variables:
  test: "base"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		// Create empty override file
		require.NoError(t, os.WriteFile(".ddx.empty.yml", []byte(""), 0644))

		t.Setenv("DDX_ENV", "empty")
		config, err := Load()
		require.NoError(t, err, "Should handle empty override file")
		assert.Equal(t, "base", config.Variables["test"], "Base config should be preserved")
	})

	t.Run("invalid_override_file_syntax", func(t *testing.T) {
		// Edge case: Invalid override file syntax

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		baseConfig := `version: "2.0"
variables:
  test: "base"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		// Create invalid YAML override file
		invalidConfig := `version: "2.0"
variables:
  test: [invalid: yaml: syntax
`
		require.NoError(t, os.WriteFile(".ddx.invalid.yml", []byte(invalidConfig), 0644))

		t.Setenv("DDX_ENV", "invalid")
		_, err := Load()

		// Should fail with clear error
		require.Error(t, err, "Should fail on invalid YAML syntax")
		assert.Contains(t, err.Error(), "parse", "Error should mention parsing issue")
	})
}

// TestUS019_PrecedenceOrder tests the precedence order implementation
func TestUS019_PrecedenceOrder(t *testing.T) {
	// Save original directory and restore at end
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	t.Run("precedence_order_correct", func(t *testing.T) {
		// Test precedence: Command-line flags > Environment override > Base config > Defaults

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// 1. Base configuration
		baseConfig := `version: "2.0"
variables:
  test_var: "base_value"
  base_only: "base"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		// 2. Environment override
		envConfig := `variables:
  test_var: "env_value"
  env_only: "env"
`
		require.NoError(t, os.WriteFile(".ddx.test.yml", []byte(envConfig), 0644))

		t.Setenv("DDX_ENV", "test")
		config, err := Load()
		require.NoError(t, err)

		// Environment should override base
		assert.Equal(t, "env_value", config.Variables["test_var"])
		assert.Equal(t, "base", config.Variables["base_only"])
		assert.Equal(t, "env", config.Variables["env_only"])

		// TODO: Add command-line flag tests when that functionality is implemented
		// 3. Command-line flags would override environment overrides
	})
}

// TestUS019_AcceptanceCriteria tests all acceptance criteria for US-019
func TestUS019_AcceptanceCriteria(t *testing.T) {
	// Save original directory and restore at end
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	t.Run("all_acceptance_criteria", func(t *testing.T) {
		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Setup base config
		baseConfig := `version: "2.0"
repository:
  url: "https://github.com/base/repo"
  branch: "main"
variables:
  api_url: "http://localhost:3000"
  log_level: "info"
  cache_enabled: "false"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(baseConfig), 0644))

		// Setup multiple environment configs
		devConfig := `variables:
  log_level: "debug"
`
		require.NoError(t, os.WriteFile(".ddx.dev.yml", []byte(devConfig), 0644))

		prodConfig := `variables:
  api_url: "https://api.production.com"
  log_level: "error"
  cache_enabled: "true"
`
		require.NoError(t, os.WriteFile(".ddx.prod.yml", []byte(prodConfig), 0644))

		// AC: Environment configs are supported
		t.Setenv("DDX_ENV", "dev")
		devLoadedConfig, err := Load()
		require.NoError(t, err)
		assert.Equal(t, "debug", devLoadedConfig.Variables["log_level"])

		// AC: DDX_ENV variable selects appropriate file
		t.Setenv("DDX_ENV", "prod")
		prodLoadedConfig, err := Load()
		require.NoError(t, err)
		assert.Equal(t, "error", prodLoadedConfig.Variables["log_level"])
		assert.Equal(t, "https://api.production.com", prodLoadedConfig.Variables["api_url"])

		// AC: Overrides merge correctly with base
		assert.Equal(t, "https://github.com/base/repo", prodLoadedConfig.Repository.URL, "Base repository preserved")

		// AC: Override takes precedence over base
		assert.Equal(t, "true", prodLoadedConfig.Variables["cache_enabled"], "Override takes precedence")

		// AC: Partial configs only change specified values
		t.Setenv("DDX_ENV", "dev")
		devPartialConfig, err := Load()
		require.NoError(t, err)
		assert.Equal(t, "debug", devPartialConfig.Variables["log_level"], "Only log_level overridden")
		assert.Equal(t, "http://localhost:3000", devPartialConfig.Variables["api_url"], "api_url preserved from base")
		assert.Equal(t, "false", devPartialConfig.Variables["cache_enabled"], "cache_enabled preserved from base")
	})
}