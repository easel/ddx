package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/easel/ddx/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigResourcesDebug tests if config loading works correctly
func TestConfigResourcesDebug(t *testing.T) {
	// Save original directory and restore at end
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	tempDir := t.TempDir()
	require.NoError(t, os.Chdir(tempDir))

	// Create simple config with resources
	configContent := `version: "2.0"
resources:
  templates:
    include:
      - "nextjs"
    exclude:
      - "react-legacy"
`
	require.NoError(t, os.WriteFile(".ddx.yml", []byte(configContent), 0644))

	// Load config and verify it has resources
	cfg, err := config.Load()
	require.NoError(t, err, "Config should load successfully")
	require.NotNil(t, cfg, "Config should not be nil")

	t.Logf("Config version: %s", cfg.Version)
	t.Logf("Config resources: %+v", cfg.Resources)

	if cfg.Resources != nil {
		t.Logf("Templates config: %+v", cfg.Resources.Templates)
		if cfg.Resources.Templates != nil {
			t.Logf("Templates include: %v", cfg.Resources.Templates.Include)
			t.Logf("Templates exclude: %v", cfg.Resources.Templates.Exclude)
		}
	}

	// Verify resources are loaded correctly
	assert.NotNil(t, cfg.Resources, "Resources section should be loaded")
	assert.NotNil(t, cfg.Resources.Templates, "Templates section should be loaded")
	assert.Contains(t, cfg.Resources.Templates.Include, "nextjs", "Should include nextjs")
	assert.Contains(t, cfg.Resources.Templates.Exclude, "react-legacy", "Should exclude react-legacy")
}

// TestResourceFilterEngine tests the filter engine directly
func TestResourceFilterEngine(t *testing.T) {
	// Create config manually
	cfg := &config.Config{
		Version: "2.0",
		Resources: &config.ResourceSelection{
			Templates: &config.ResourceFilter{
				Include: []string{"nextjs.yml", "react-*"},
				Exclude: []string{"react-legacy.yml"},
			},
		},
	}

	// Create temp library structure
	libDir := t.TempDir()
	templateDir := filepath.Join(libDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Create test files
	files := map[string]string{
		"nextjs.yml":       "NextJS template",
		"react-legacy.yml": "Legacy React template",
		"react-modern.yml": "Modern React template",
		"vue.yml":          "Vue template",
	}

	for name, content := range files {
		path := filepath.Join(templateDir, name)
		require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	}

	// Create filter engine
	engine := config.NewResourceFilterEngine(cfg, libDir)

	// Test filtering
	allPaths := []string{
		filepath.Join(templateDir, "nextjs.yml"),
		filepath.Join(templateDir, "react-legacy.yml"),
		filepath.Join(templateDir, "react-modern.yml"),
		filepath.Join(templateDir, "vue.yml"),
	}

	filtered, err := engine.FilterResources("templates", allPaths)
	require.NoError(t, err, "Filtering should work")

	t.Logf("All paths: %v", allPaths)
	t.Logf("Filtered paths: %v", filtered)

	// Convert to base names for easier checking
	var baseNames []string
	for _, path := range filtered {
		baseNames = append(baseNames, filepath.Base(path))
	}

	t.Logf("Filtered base names: %v", baseNames)

	// Verify filtering worked
	assert.Contains(t, baseNames, "nextjs.yml", "Should include nextjs.yml")
	assert.Contains(t, baseNames, "react-modern.yml", "Should include react-modern.yml (matches react-*)")
	assert.NotContains(t, baseNames, "react-legacy.yml", "Should exclude react-legacy.yml")
	assert.NotContains(t, baseNames, "vue.yml", "Should not include vue.yml (not in include patterns)")
}