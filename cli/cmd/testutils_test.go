package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/easel/ddx/internal/config"
	"github.com/stretchr/testify/require"
)

// TestEnvironment provides isolated testing environment for .ddx/config.yaml
type TestEnvironment struct {
	Dir        string
	ConfigPath string
	t          *testing.T
}

// NewTestEnvironment creates a clean test environment with temp directory
func NewTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()

	tempDir := t.TempDir()
	ddxDir := filepath.Join(tempDir, ".ddx")
	configPath := filepath.Join(ddxDir, "config.yaml")

	// Create .ddx directory
	require.NoError(t, os.MkdirAll(ddxDir, 0755))

	return &TestEnvironment{
		Dir:        tempDir,
		ConfigPath: configPath,
		t:          t,
	}
}

// CreateConfig creates a config file with the given content
func (te *TestEnvironment) CreateConfig(content string) {
	te.t.Helper()
	require.NoError(te.t, os.WriteFile(te.ConfigPath, []byte(content), 0644))
}

// CreateDefaultConfig creates a minimal valid config file
func (te *TestEnvironment) CreateDefaultConfig() {
	te.t.Helper()
	content := `version: "1.0"
library_base_path: "./library"
repository:
  url: "https://github.com/easel/ddx"
  branch: "main"
  subtree_prefix: "library"
variables: {}
`
	te.CreateConfig(content)
}

// LoadConfig loads the config using ConfigLoader
func (te *TestEnvironment) LoadConfig() (*config.Config, error) {
	loader, err := config.NewConfigLoaderWithWorkingDir(te.Dir)
	if err != nil {
		return nil, err
	}
	return loader.LoadConfig()
}

// CreateFile creates any file in the test environment
func (te *TestEnvironment) CreateFile(relativePath, content string) {
	te.t.Helper()
	fullPath := filepath.Join(te.Dir, relativePath)
	dir := filepath.Dir(fullPath)
	require.NoError(te.t, os.MkdirAll(dir, 0755))
	require.NoError(te.t, os.WriteFile(fullPath, []byte(content), 0644))
}