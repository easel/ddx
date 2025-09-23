package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUS020_ConfigureResourceSelection tests US-020: Configure Resource Selection
func TestUS020_ConfigureResourceSelection(t *testing.T) {
	// Save original directory and restore at end
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	t.Run("resource_selection_honored_during_operations", func(t *testing.T) {
		// AC: Given configuration, when resources specified, then selection is honored during operations

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Create config with resource selection
		config := `version: "2.0"
resources:
  prompts:
    include:
      - "code-review.md"
      - "testing/*"
    exclude:
      - "testing/experimental/*"
  templates:
    include:
      - "nextjs.yml"
      - "react-*"
    exclude:
      - "react-legacy.yml"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		// Setup mock library structure
		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()

		// Test that list operation honors resource selection
		output, err := executeCommand(rootCmd, "list", "--json")
		require.NoError(t, err, "List should honor resource selection")

		// Should include selected prompts
		assert.Contains(t, output, "code-review.md", "Should include code-review prompt")
		assert.Contains(t, output, "testing/basic", "Should include testing/basic prompt")

		// Should exclude experimental testing prompts
		assert.NotContains(t, output, "testing/experimental", "Should exclude experimental testing")

		// Should include selected templates
		assert.Contains(t, output, "nextjs.yml", "Should include nextjs template")
		assert.Contains(t, output, "react-component.yml", "Should include react-component template")

		// Should exclude legacy react template
		assert.NotContains(t, output, "react-legacy.yml", "Should exclude react-legacy template")
	})

	t.Run("wildcard_patterns_work_for_selection", func(t *testing.T) {
		// AC: Given patterns, when used, then wildcards work for resource selection

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Test various wildcard patterns
		config := `version: "2.0"
resources:
  patterns:
    include:
      - "auth/*"              # Matches any auth pattern
      - "logging.md"          # Exact match
      - "error-*"             # Prefix wildcard
      - "*-validation.md"     # Suffix wildcard
      - "api/*/handlers.md"   # Nested wildcard
    exclude:
      - "auth/legacy/*"       # Exclude legacy auth patterns
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "list", "patterns", "--json")
		require.NoError(t, err, "List patterns should work with wildcards")

		// Should match auth/* but exclude auth/legacy/*
		assert.Contains(t, output, "auth/oauth", "Should include auth/oauth")
		assert.Contains(t, output, "auth/jwt", "Should include auth/jwt")
		assert.NotContains(t, output, "auth/legacy/basic", "Should exclude auth/legacy/basic")

		// Should match exact name
		assert.Contains(t, output, "logging", "Should include exact match 'logging'")

		// Should match prefix wildcards
		assert.Contains(t, output, "error-handling", "Should include error-handling")
		assert.Contains(t, output, "error-recovery", "Should include error-recovery")

		// Should match suffix wildcards
		assert.Contains(t, output, "input-validation", "Should include input-validation")
		assert.Contains(t, output, "data-validation", "Should include data-validation")

		// Should match nested wildcards
		assert.Contains(t, output, "api/v1/handlers", "Should include api/v1/handlers")
		assert.Contains(t, output, "api/v2/handlers", "Should include api/v2/handlers")
	})

	t.Run("include_exclude_rules_applied_correctly", func(t *testing.T) {
		// AC: Given resources, when configured, then include/exclude rules are applied correctly

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Test complex include/exclude logic
		config := `version: "2.0"
resources:
  prompts:
    include:
      - "code-*"
      - "testing/*"
      - "documentation/*"
    exclude:
      - "code-legacy.md"
      - "testing/experimental/*"
      - "documentation/internal/*"
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "list", "prompts", "--json")
		require.NoError(t, err, "List should apply include/exclude rules")

		// Include rules should match
		assert.Contains(t, output, "code-review", "Should include code-review")
		assert.Contains(t, output, "code-analysis", "Should include code-analysis")
		assert.Contains(t, output, "testing/unit", "Should include testing/unit")
		assert.Contains(t, output, "documentation/api", "Should include documentation/api")

		// Exclude rules should override includes
		assert.NotContains(t, output, "code-legacy", "Should exclude code-legacy despite code-* include")
		assert.NotContains(t, output, "testing/experimental", "Should exclude experimental testing")
		assert.NotContains(t, output, "documentation/internal", "Should exclude internal docs")

		// Items not in include patterns should be excluded
		assert.NotContains(t, output, "deployment", "Should exclude items not in include patterns")
	})

	t.Run("resource_grouping_by_category_supported", func(t *testing.T) {
		// AC: Given resource types, when organized, then grouping by category is supported

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		config := `version: "2.0"
resources:
  prompts:
    include: ["*"]
  templates:
    include: ["react-*", "nextjs.yml"]
  patterns:
    include: ["auth/*", "logging.md"]
  configs:
    include: ["eslint.json", "prettier.json"]
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "list", "--json")
		require.NoError(t, err, "List should support resource grouping")

		// Should group by resource type/category
		assert.Contains(t, output, `"type":"prompts"`, "Should group prompts")
		assert.Contains(t, output, `"type":"templates"`, "Should group templates")
		assert.Contains(t, output, `"type":"patterns"`, "Should group patterns")
		assert.Contains(t, output, `"type":"configs"`, "Should group configs")
	})

	t.Run("resource_dependencies_automatically_included", func(t *testing.T) {
		// AC: Given resource dependencies, when present, then they are automatically included

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		config := `version: "2.0"
resources:
  templates:
    include:
      - "nextjs-app"  # This has dependencies
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryWithDependencies(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "list", "templates", "--json")
		require.NoError(t, err, "List should include dependencies")

		// Should include the explicitly requested template
		assert.Contains(t, output, "nextjs-app", "Should include explicitly selected template")

		// Should automatically include dependencies
		assert.Contains(t, output, "eslint-config", "Should auto-include eslint-config dependency")
		assert.Contains(t, output, "prettier-config", "Should auto-include prettier-config dependency")
		assert.Contains(t, output, "react-patterns", "Should auto-include react-patterns dependency")
	})

	t.Run("preview_of_selection_available", func(t *testing.T) {
		// AC: Given resource config, when complete, then preview of selection is available

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		config := `version: "2.0"
resources:
  prompts:
    include: ["code-*", "testing/*"]
    exclude: ["testing/experimental/*"]
  templates:
    include: ["react-*"]
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "config", "resources", "--preview")
		require.NoError(t, err, "Preview should be available")

		// Should show what would be included
		assert.Contains(t, output, "Preview of selected resources", "Should show preview header")
		assert.Contains(t, output, "code-review", "Should preview included prompts")
		assert.Contains(t, output, "react-component", "Should preview included templates")
		assert.Contains(t, output, "Excluded:", "Should show excluded items")
		assert.Contains(t, output, "testing/experimental", "Should show excluded patterns")

		// Should show summary counts
		assert.Contains(t, output, "prompts:", "Should show prompt count")
		assert.Contains(t, output, "templates:", "Should show template count")
	})

	t.Run("resource_availability_validated", func(t *testing.T) {
		// AC: Given resource paths, when specified, then availability is validated

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		// Config with non-existent resources
		config := `version: "2.0"
resources:
  prompts:
    include:
      - "code-review.md"    # exists
      - "nonexistent-prompt" # doesn't exist
  templates:
    include:
      - "react-*"           # pattern matches existing
      - "missing-*"         # pattern matches nothing
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "config", "validate")

		// Should validate and report missing resources
		if err != nil {
			// Validation should catch missing resources
			assert.Contains(t, err.Error(), "nonexistent-prompt", "Should report missing prompt")
			assert.Contains(t, strings.ToLower(output), "warning", "Should show validation warnings")
		} else {
			// If no error, should show warnings in output
			assert.Contains(t, strings.ToLower(output), "warning", "Should show validation warnings")
			assert.Contains(t, output, "nonexistent-prompt", "Should warn about missing prompt")
			assert.Contains(t, output, "missing-*", "Should warn about patterns with no matches")
		}
	})

	t.Run("tree_view_shows_resource_structure", func(t *testing.T) {
		// AC: Given resources, when listed, then tree view shows structure

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		config := `version: "2.0"
resources:
  prompts:
    include: ["*"]
  patterns:
    include: ["auth/*", "error-*"]
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "list", "--tree")
		require.NoError(t, err, "Tree view should work")

		// Should show hierarchical structure
		assert.Contains(t, output, "├──", "Should show tree structure symbols")
		assert.Contains(t, output, "└──", "Should show tree structure symbols")
		assert.Contains(t, output, "📁", "Should show folder icons")
		assert.Contains(t, output, "📄", "Should show file icons")

		// Should show nested structure
		assert.Contains(t, output, "auth/", "Should show auth directory")
		assert.Contains(t, output, "  oauth", "Should show nested oauth item")
		assert.Contains(t, output, "  jwt", "Should show nested jwt item")
	})

	t.Run("conflicting_include_exclude_rules", func(t *testing.T) {
		// Edge case: Conflicting include/exclude rules (exclude should win)

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		config := `version: "2.0"
resources:
  prompts:
    include:
      - "code-*"          # Includes code-legacy
    exclude:
      - "code-legacy.md"  # But also excludes it - exclude should win
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "list", "prompts", "--json")
		require.NoError(t, err, "Should handle conflicting rules")

		// Exclude should take precedence over include
		assert.Contains(t, output, "code-review", "Should include other code-* items")
		assert.NotContains(t, output, "code-legacy", "Exclude should override include")
	})

	t.Run("wildcard_patterns_matching_nothing", func(t *testing.T) {
		// Edge case: Wildcard patterns that match nothing

		tempDir := t.TempDir()
		require.NoError(t, os.Chdir(tempDir))

		config := `version: "2.0"
resources:
  prompts:
    include:
      - "nonexistent-*"   # Matches nothing
      - "another-missing-*" # Also matches nothing
`
		require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))

		libDir := filepath.Join(tempDir, ".ddx", "library")
		setupMockLibraryForResourceSelection(t, libDir)

		rootCmd := getConfigTestRootCommand()
		output, err := executeCommand(rootCmd, "list", "prompts", "--json")
		require.NoError(t, err, "Should handle patterns matching nothing")

		// Should return empty results gracefully
		assert.Contains(t, output, `"resources":[]`, "Should return empty resources array")
	})
}

// setupMockLibraryForResourceSelection creates a mock library structure for testing
func setupMockLibraryForResourceSelection(t *testing.T, libDir string) {
	t.Helper()

	// Create directory structure
	dirs := []string{
		"prompts",
		"prompts/testing",
		"prompts/testing/experimental",
		"prompts/documentation",
		"prompts/documentation/internal",
		"templates",
		"patterns",
		"patterns/auth",
		"patterns/auth/legacy",
		"patterns/api/v1",
		"patterns/api/v2",
		"configs",
	}

	for _, dir := range dirs {
		require.NoError(t, os.MkdirAll(filepath.Join(libDir, dir), 0755))
	}

	// Create mock files
	files := map[string]string{
		"prompts/code-review.md":                    "Code review prompt",
		"prompts/code-analysis.md":                  "Code analysis prompt",
		"prompts/code-legacy.md":                    "Legacy code prompt",
		"prompts/testing/unit.md":                   "Unit testing prompt",
		"prompts/testing/basic.md":                  "Basic testing prompt",
		"prompts/testing/experimental/new.md":      "Experimental testing prompt",
		"prompts/documentation/api.md":              "API documentation prompt",
		"prompts/documentation/internal/spec.md":   "Internal documentation",
		"prompts/deployment.md":                     "Deployment prompt",
		"templates/nextjs.yml":                      "NextJS template",
		"templates/react-component.yml":             "React component template",
		"templates/react-legacy.yml":                "Legacy React template",
		"patterns/auth/oauth.md":                    "OAuth authentication pattern",
		"patterns/auth/jwt.md":                      "JWT authentication pattern",
		"patterns/auth/legacy/basic.md":             "Basic auth pattern",
		"patterns/logging.md":                       "Logging pattern",
		"patterns/error-handling.md":                "Error handling pattern",
		"patterns/error-recovery.md":                "Error recovery pattern",
		"patterns/input-validation.md":              "Input validation pattern",
		"patterns/data-validation.md":               "Data validation pattern",
		"patterns/api/v1/handlers.md":               "API v1 handlers",
		"patterns/api/v2/handlers.md":               "API v2 handlers",
		"configs/eslint.json":                       "ESLint configuration",
		"configs/prettier.json":                     "Prettier configuration",
	}

	for filePath, content := range files {
		fullPath := filepath.Join(libDir, filePath)
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}
}

// setupMockLibraryWithDependencies creates a mock library with dependency metadata
func setupMockLibraryWithDependencies(t *testing.T, libDir string) {
	t.Helper()

	// Setup basic library first
	setupMockLibraryForResourceSelection(t, libDir)

	// Create template with dependencies
	templateDir := filepath.Join(libDir, "templates", "nextjs-app")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Template metadata with dependencies
	metaContent := `name: "NextJS Application"
description: "Full-featured NextJS application template"
dependencies:
  configs:
    - "eslint-config"
    - "prettier-config"
  patterns:
    - "react-patterns"
`
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "meta.yml"), []byte(metaContent), 0644))

	// Create dependency files
	dependencyFiles := map[string]string{
		"configs/eslint-config.json":    "ESLint configuration for NextJS",
		"configs/prettier-config.json":  "Prettier configuration for NextJS",
		"patterns/react-patterns.md":    "Common React patterns",
	}

	for filePath, content := range dependencyFiles {
		fullPath := filepath.Join(libDir, filePath)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}
}