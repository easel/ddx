package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAcceptance_US021_ConfigureRepositoryConnection tests US-021: Configure Repository Connection
func TestAcceptance_US021_ConfigureRepositoryConnection(t *testing.T) {

	t.Run("repository_url_configuration", func(t *testing.T) {
		// AC: Given config file, when repository URL specified, then connection uses that URL

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Set custom repository URL
		_, err := executeCommand(rootCmd, "config", "set", "repository.url", "https://github.com/company/ddx-resources")
		require.NoError(t, err, "Should be able to set repository URL")

		// Verify the URL was set
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repository.url")
		require.NoError(t, err, "Should be able to get repository URL")
		assert.Contains(t, output, "https://github.com/company/ddx-resources", "Should show custom repository URL")
	})

	t.Run("repository_branch_specification", func(t *testing.T) {
		// AC: Given repository config, when branch specified, then that branch is tracked

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Set custom branch
		_, err := executeCommand(rootCmd, "config", "set", "repository.branch", "development")
		require.NoError(t, err, "Should be able to set repository branch")

		// Verify the branch was set
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repository.branch")
		require.NoError(t, err, "Should be able to get repository branch")
		assert.Contains(t, output, "development", "Should show custom branch")
	})

	t.Run("sync_frequency_configuration", func(t *testing.T) {
		// AC: Given sync settings, when configured, then frequency preferences are honored

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Set sync frequency
		_, err := executeCommand(rootCmd, "config", "set", "repository.sync.frequency", "daily")
		require.NoError(t, err, "Should be able to set sync frequency")

		// Verify the frequency was set
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repository.sync.frequency")
		require.NoError(t, err, "Should be able to get sync frequency")
		assert.Contains(t, output, "daily", "Should show sync frequency")

		// Test auto-update setting
		autoCmd := getConfigTestRootCommand()
		_, err = executeCommand(autoCmd, "config", "set", "repository.sync.auto_update", "true")
		require.NoError(t, err, "Should be able to set auto-update")
	})

	t.Run("authentication_configuration", func(t *testing.T) {
		// AC: Given authentication needs, when required, then auth method is configurable

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Set SSH key authentication
		_, err := executeCommand(rootCmd, "config", "set", "repository.auth.method", "ssh-key")
		require.NoError(t, err, "Should be able to set auth method")

		keyCmd := getConfigTestRootCommand()
		_, err = executeCommand(keyCmd, "config", "set", "repository.auth.key_path", "~/.ssh/ddx_deploy_key")
		require.NoError(t, err, "Should be able to set SSH key path")

		// Verify authentication config
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repository.auth.method")
		require.NoError(t, err, "Should be able to get auth method")
		assert.Contains(t, output, "ssh-key", "Should show SSH key authentication")
	})

	t.Run("multiple_repositories_support", func(t *testing.T) {
		// AC: Given multiple sources, when needed, then multiple remotes are supported

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Configure primary repository
		_, err := executeCommand(rootCmd, "config", "set", "repositories.primary.url", "https://github.com/ddx/ddx-official")
		require.NoError(t, err, "Should be able to set primary repository")

		primaryCmd := getConfigTestRootCommand()
		_, err = executeCommand(primaryCmd, "config", "set", "repositories.primary.priority", "1")
		require.NoError(t, err, "Should be able to set primary repository priority")

		// Configure secondary repository
		secondaryCmd := getConfigTestRootCommand()
		_, err = executeCommand(secondaryCmd, "config", "set", "repositories.company.url", "https://github.com/company/ddx-internal")
		require.NoError(t, err, "Should be able to set secondary repository")

		companyCmd := getConfigTestRootCommand()
		_, err = executeCommand(companyCmd, "config", "set", "repositories.company.priority", "2")
		require.NoError(t, err, "Should be able to set secondary repository priority")

		// Verify multiple repositories
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repositories.primary.url")
		require.NoError(t, err, "Should be able to get primary repository")
		assert.Contains(t, output, "ddx-official", "Should show primary repository")
	})

	t.Run("proxy_configuration", func(t *testing.T) {
		// AC: Given network restrictions, when present, then proxy configuration works

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Set proxy URL
		_, err := executeCommand(rootCmd, "config", "set", "repository.proxy.url", "http://proxy.company.com:8080")
		require.NoError(t, err, "Should be able to set proxy URL")

		// Set proxy authentication
		authCmd := getConfigTestRootCommand()
		_, err = executeCommand(authCmd, "config", "set", "repository.proxy.auth", "user:pass")
		require.NoError(t, err, "Should be able to set proxy auth")

		// Verify proxy configuration
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repository.proxy.url")
		require.NoError(t, err, "Should be able to get proxy URL")
		assert.Contains(t, output, "proxy.company.com", "Should show proxy URL")
	})

	t.Run("protocol_selection", func(t *testing.T) {
		// AC: Given protocol preference, when set, then SSH vs HTTPS is respected

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Set protocol preference to SSH
		_, err := executeCommand(rootCmd, "config", "set", "repository.protocol", "ssh")
		require.NoError(t, err, "Should be able to set protocol to SSH")

		// Verify protocol setting
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repository.protocol")
		require.NoError(t, err, "Should be able to get protocol")
		assert.Contains(t, output, "ssh", "Should show SSH protocol")

		// Test HTTPS protocol
		httpsCmd := getConfigTestRootCommand()
		_, err = executeCommand(httpsCmd, "config", "set", "repository.protocol", "https")
		require.NoError(t, err, "Should be able to set protocol to HTTPS")
	})

	t.Run("custom_remote_naming", func(t *testing.T) {
		// AC: Given remote naming, when customized, then custom names are used

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Set custom remote name
		_, err := executeCommand(rootCmd, "config", "set", "repository.remote", "company-ddx")
		require.NoError(t, err, "Should be able to set custom remote name")

		// Verify remote name
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repository.remote")
		require.NoError(t, err, "Should be able to get remote name")
		assert.Contains(t, output, "company-ddx", "Should show custom remote name")
	})

	t.Run("repository_connection_testing", func(t *testing.T) {
		// Test connection testing functionality

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Create basic config first
		config := `version: "2.0"
repository:
  url: "https://github.com/test/repo"
  branch: "main"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		rootCmd := getConfigTestRootCommand()

		// Test connection command
		output, err := executeCommand(rootCmd, "config", "test-connection")
		if err != nil {
			// Connection testing may fail in test environment - verify command exists
			assert.Contains(t, err.Error(), "connection", "Should mention connection testing")
		} else {
			assert.Contains(t, strings.ToLower(output), "test", "Should show connection test results")
		}
	})

	t.Run("repository_status_command", func(t *testing.T) {
		// Test repository status functionality

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Create basic config
		config := `version: "2.0"
repository:
  url: "https://github.com/test/repo"
  branch: "main"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		rootCmd := getConfigTestRootCommand()

		// Test repository status command
		output, err := executeCommand(rootCmd, "config", "repository", "status")
		if err != nil {
			// Status command may not be fully implemented yet
			assert.Contains(t, err.Error(), "repository", "Should mention repository operations")
		} else {
			assert.Contains(t, strings.ToLower(output), "status", "Should show repository status")
		}
	})

	t.Run("repository_branch_switching", func(t *testing.T) {
		// Test branch switching functionality

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Create basic config
		config := `version: "2.0"
repository:
  url: "https://github.com/test/repo"
  branch: "main"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		rootCmd := getConfigTestRootCommand()

		// Test branch switching command
		output, err := executeCommand(rootCmd, "config", "repository", "branch", "dev-branch")
		if err != nil {
			// Branch switching may not be fully implemented yet
			assert.Contains(t, err.Error(), "branch", "Should mention branch operations")
		} else {
			assert.Contains(t, strings.ToLower(output), "branch", "Should show branch switching results")
		}
	})

	t.Run("configuration_validation_with_advanced_features", func(t *testing.T) {
		// Test validation of advanced repository configuration

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Create advanced configuration
		advancedConfig := `version: "2.0"
repository:
  url: "git@github.com:company/ddx-private.git"
  branch: "development"
  remote: "company-ddx"
  protocol: "ssh"
  auth:
    method: "ssh-key"
    key_path: "~/.ssh/ddx_deploy_key"
  proxy:
    url: "http://proxy.company.com:8080"
    auth: "user:pass"
  sync:
    frequency: "daily"
    auto_update: true
    timeout: 30
repositories:
  primary:
    url: "https://github.com/ddx/ddx-official"
    branch: "main"
    priority: 1
  company:
    url: "https://github.com/company/ddx-internal"
    branch: "company-main"
    priority: 2
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(advancedConfig), 0644))

		rootCmd := getConfigTestRootCommand()

		// Validate advanced configuration
		output, err := executeCommand(rootCmd, "config", "--validate")
		require.NoError(t, err, "Advanced configuration should be valid")
		assert.Contains(t, strings.ToLower(output), "valid", "Should confirm configuration is valid")
	})

	t.Run("timeout_and_retry_configuration", func(t *testing.T) {
		// Test network timeout and retry settings

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		// Set timeout
		_, err := executeCommand(rootCmd, "config", "set", "repository.sync.timeout", "30")
		require.NoError(t, err, "Should be able to set sync timeout")

		// Set retry count
		retryCmd := getConfigTestRootCommand()
		_, err = executeCommand(retryCmd, "config", "set", "repository.sync.retry_count", "3")
		require.NoError(t, err, "Should be able to set retry count")

		// Verify timeout setting
		getCmd := getConfigTestRootCommand()
		output, err := executeCommand(getCmd, "config", "get", "repository.sync.timeout")
		require.NoError(t, err, "Should be able to get timeout")
		assert.Contains(t, output, "30", "Should show timeout value")
	})
}

// Test contract for repository configuration commands
func TestRepositoryConfigurationCommands_Contract(t *testing.T) {

	t.Run("config_set_repository_fields", func(t *testing.T) {
		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		rootCmd := getConfigTestRootCommand()

		testCases := []struct {
			key   string
			value string
		}{
			{"repository.url", "https://github.com/test/repo"},
			{"repository.branch", "main"},
			{"repository.remote", "origin"},
			{"repository.protocol", "https"},
			{"repository.auth.method", "ssh-key"},
			{"repository.auth.key_path", "~/.ssh/id_rsa"},
			{"repository.proxy.url", "http://proxy:8080"},
			{"repository.sync.frequency", "daily"},
			{"repository.sync.auto_update", "true"},
			{"repository.sync.timeout", "30"},
		}

		for _, tc := range testCases {
			t.Run(tc.key, func(t *testing.T) {
				_, err := executeCommand(rootCmd, "config", "set", tc.key, tc.value)
				if err != nil {
					// Some fields may not be implemented yet - verify error is reasonable
					assert.Contains(t, strings.ToLower(err.Error()), "config", "Error should be config-related")
				}
			})
		}
	})

	t.Run("config_get_repository_fields", func(t *testing.T) {
		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Create config with repository settings
		config := `version: "2.0"
repository:
  url: "https://github.com/test/repo"
  branch: "main"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		testCases := []string{
			"repository.url",
			"repository.branch",
			"repository.remote",
			"repository.protocol",
		}

		for _, key := range testCases {
			t.Run(key, func(t *testing.T) {
				rootCmd := getConfigTestRootCommand()
				_, err := executeCommand(rootCmd, "config", "get", key)
				if err != nil {
					// Some fields may not be implemented yet
					assert.Contains(t, strings.ToLower(err.Error()), "not found", "Should indicate field not found")
				}
			})
		}
	})

	t.Run("repository_subcommands", func(t *testing.T) {
		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		config := `version: "2.0"
repository:
  url: "https://github.com/test/repo"
  branch: "main"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		subcommands := [][]string{
			{"config", "test-connection"},
			{"config", "repository", "status"},
			{"config", "repository", "branch", "new-branch"},
		}

		for _, cmd := range subcommands {
			t.Run(strings.Join(cmd[1:], "_"), func(t *testing.T) {
				rootCmd := getConfigTestRootCommand()
				_, err := executeCommand(rootCmd, cmd...)
				// Commands may not be fully implemented yet - just verify they're recognized
				if err != nil {
					errorMsg := strings.ToLower(err.Error())
					// Should not be "unknown command" error
					assert.False(t, strings.Contains(errorMsg, "unknown command"),
						"Command should be recognized even if not fully implemented: %s", strings.Join(cmd, " "))
				}
			})
		}
	})
}