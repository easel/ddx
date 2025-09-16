package mcp_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/easel/ddx/internal/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstaller(t *testing.T) {
	// Setup test registry
	registryPath, cleanup := setupTestRegistry(t)
	defer cleanup()

	// Set DDX_LIBRARY_BASE_PATH for tests to use the temp directory
	libDir := filepath.Dir(filepath.Dir(registryPath))
	t.Setenv("DDX_LIBRARY_BASE_PATH", libDir)

	t.Run("install new server", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "claude.json")

		installer := mcp.NewInstaller()
		opts := mcp.InstallOptions{
			ConfigPath: configPath,
			Environment: map[string]string{
				"GITHUB_PERSONAL_ACCESS_TOKEN": "ghp_test123456789",
			},
			NoBackup: true,
		}

		err := installer.Install("github", opts)
		require.NoError(t, err)

		// Verify config was written
		data, err := os.ReadFile(configPath)
		require.NoError(t, err)

		var config map[string]interface{}
		err = json.Unmarshal(data, &config)
		require.NoError(t, err)

		servers, ok := config["mcpServers"].(map[string]interface{})
		require.True(t, ok)
		require.Contains(t, servers, "github")
	})

	t.Run("server already installed", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "claude.json")

		// Create existing config
		existingConfig := map[string]interface{}{
			"mcpServers": map[string]interface{}{
				"github": map[string]interface{}{
					"command": "npx",
				},
			},
		}
		data, _ := json.Marshal(existingConfig)
		os.WriteFile(configPath, data, 0600)

		installer := mcp.NewInstaller()
		opts := mcp.InstallOptions{
			ConfigPath: configPath,
			Environment: map[string]string{
				"GITHUB_PERSONAL_ACCESS_TOKEN": "test",
			},
			NoBackup: true,
		}

		err := installer.Install("github", opts)
		assert.Error(t, err)
		assert.ErrorIs(t, err, mcp.ErrAlreadyInstalled)
	})

	t.Run("missing required environment", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "claude.json")

		installer := mcp.NewInstaller()
		opts := mcp.InstallOptions{
			ConfigPath:  configPath,
			Environment: map[string]string{}, // Missing required token
			NoBackup:    true,
		}

		err := installer.Install("github", opts)
		assert.Error(t, err)
		assert.ErrorIs(t, err, mcp.ErrMissingRequired)
	})

	t.Run("dry run", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "claude.json")

		installer := mcp.NewInstaller()
		opts := mcp.InstallOptions{
			ConfigPath: configPath,
			Environment: map[string]string{
				"GITHUB_PERSONAL_ACCESS_TOKEN": "test",
			},
			DryRun: true,
		}

		err := installer.Install("github", opts)
		require.NoError(t, err)

		// Config should not be written
		_, err = os.Stat(configPath)
		assert.True(t, os.IsNotExist(err))
	})
}

func TestConfigManager(t *testing.T) {
	t.Run("load and save config", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "claude.json")

		// Create initial config
		initialConfig := map[string]interface{}{
			"mcpServers": map[string]interface{}{
				"existing": map[string]interface{}{
					"command": "test",
				},
			},
			"otherField": "preserved",
		}
		data, _ := json.MarshalIndent(initialConfig, "", "  ")
		os.WriteFile(configPath, data, 0600)

		// Load and modify
		cm := mcp.NewConfigManager(configPath)
		err := cm.Load()
		require.NoError(t, err)

		// Add new server
		err = cm.AddServer("new", mcp.ServerConfig{
			Command: "npx",
			Args:    []string{"test"},
		})
		require.NoError(t, err)

		// Save
		err = cm.Save()
		require.NoError(t, err)

		// Verify
		data, _ = os.ReadFile(configPath)
		var saved map[string]interface{}
		json.Unmarshal(data, &saved)

		// Check that other fields are preserved
		assert.Equal(t, "preserved", saved["otherField"])

		// Check new server was added
		servers := saved["mcpServers"].(map[string]interface{})
		assert.Contains(t, servers, "new")
		assert.Contains(t, servers, "existing")
	})

	t.Run("backup and restore", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "claude.json")

		// Create initial config
		initialData := []byte(`{"test": "data"}`)
		os.WriteFile(configPath, initialData, 0600)

		cm := mcp.NewConfigManager(configPath)

		// Create backup
		err := cm.Backup()
		require.NoError(t, err)

		// Modify original
		os.WriteFile(configPath, []byte(`{"modified": true}`), 0600)

		// Restore
		err = cm.Restore()
		require.NoError(t, err)

		// Verify restored
		data, _ := os.ReadFile(configPath)
		assert.Equal(t, initialData, data)
	})
}

func TestValidator(t *testing.T) {
	v := mcp.NewValidator()

	t.Run("validate server name", func(t *testing.T) {
		tests := []struct {
			name    string
			input   string
			wantErr bool
		}{
			{"valid", "github", false},
			{"with hyphen", "github-enterprise", false},
			{"numbers", "server123", false},
			{"empty", "", true},
			{"uppercase", "GitHub", true},
			{"spaces", "git hub", true},
			{"path traversal", "../etc/passwd", true},
			{"path separator", "servers/github", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateServerName(tt.input)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("validate environment", func(t *testing.T) {
		tests := []struct {
			name    string
			env     map[string]string
			wantErr bool
		}{
			{"valid", map[string]string{"TOKEN": "value"}, false},
			{"underscore", map[string]string{"API_TOKEN": "value"}, false},
			{"lowercase key", map[string]string{"token": "value"}, true},
			{"shell injection", map[string]string{"TOKEN": "$(rm -rf /)"}, true},
			{"backticks", map[string]string{"TOKEN": "`echo hacked`"}, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidateEnvironment(tt.env)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("validate path", func(t *testing.T) {
		tests := []struct {
			name    string
			path    string
			wantErr bool
		}{
			{"absolute", "/home/user/config.json", false},
			{"relative", "config.json", true},
			{"path traversal", "/home/../etc/passwd", true},
			{"double dots", "/home/user/../../../etc", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := v.ValidatePath(tt.path)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestMaskSensitive(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		sensitive bool
		expected  string
	}{
		{"not sensitive", "regular", false, "regular"},
		{"sensitive long", "ghp_secrettoken123", true, "ghp_***"},
		{"sensitive short", "secret", true, "***"},
		{"empty", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mcp.MaskSensitive(tt.value, tt.sensitive)
			assert.Equal(t, tt.expected, result)
		})
	}
}
