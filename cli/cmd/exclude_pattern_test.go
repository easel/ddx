package cmd

import (
	"path/filepath"
	"testing"

	"github.com/easel/ddx/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestExcludePatternMatching(t *testing.T) {
	// Create config with exclude pattern
	cfg := &config.Config{
		Version: "2.0",
		Resources: &config.ResourceSelection{
			Prompts: &config.ResourceFilter{
				Include: []string{"testing/*"},
				Exclude: []string{"testing/experimental/*"},
			},
		},
	}

	// Create filter engine (libPath doesn't matter for this test)
	engine := config.NewResourceFilterEngine(cfg, "/tmp")

	// Test paths that should be checked
	testPaths := []string{
		"/tmp/prompts/testing/basic.md",
		"/tmp/prompts/testing/unit.md",
		"/tmp/prompts/testing/experimental/new.md",
	}

	filtered, err := engine.FilterResources("prompts", testPaths)
	assert.NoError(t, err)

	// Convert to base names for easier checking
	var filteredNames []string
	for _, path := range filtered {
		rel := filepath.Base(path)
		// Get the part after "prompts/"
		if filepath.Dir(path) != "/tmp/prompts" {
			parent := filepath.Base(filepath.Dir(path))
			if filepath.Dir(filepath.Dir(path)) != "/tmp/prompts" {
				grandparent := filepath.Base(filepath.Dir(filepath.Dir(path)))
				rel = grandparent + "/" + parent + "/" + rel
			} else {
				rel = parent + "/" + rel
			}
		}
		filteredNames = append(filteredNames, rel)
	}

	t.Logf("Original paths: %v", testPaths)
	t.Logf("Filtered paths: %v", filtered)
	t.Logf("Filtered names: %v", filteredNames)

	// Should include testing/basic.md and testing/unit.md
	assert.Contains(t, filteredNames, "testing/basic.md", "Should include testing/basic.md")
	assert.Contains(t, filteredNames, "testing/unit.md", "Should include testing/unit.md")

	// Should exclude testing/experimental/new.md
	assert.NotContains(t, filteredNames, "experimental/new.md", "Should exclude experimental/new.md")
	assert.NotContains(t, filteredNames, "testing/experimental/new.md", "Should exclude testing/experimental/new.md")
}