package templates

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReplaceVariables tests variable substitution
func TestReplaceVariables_Basic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		content   string
		variables map[string]string
		expected  string
	}{
		{
			name:    "simple replacement",
			content: "Hello {{name}}!",
			variables: map[string]string{
				"name": "World",
			},
			expected: "Hello World!",
		},
		{
			name:    "multiple variables",
			content: "{{greeting}} {{name}}, version {{version}}",
			variables: map[string]string{
				"greeting": "Hello",
				"name":     "User",
				"version":  "1.0.0",
			},
			expected: "Hello User, version 1.0.0",
		},
		{
			name:    "variable with spaces",
			content: "Name: {{ name }}",
			variables: map[string]string{
				"name": "Test",
			},
			expected: "Name: Test",
		},
		{
			name:    "uppercase dollar syntax",
			content: "Port: ${PORT}",
			variables: map[string]string{
				"port": "8080",
			},
			expected: "Port: 8080",
		},
		{
			name:    "mixed syntax",
			content: "{{name}} uses ${PORT}",
			variables: map[string]string{
				"name": "App",
				"port": "3000",
			},
			expected: "App uses 3000",
		},
		{
			name:      "no variables",
			content:   "Plain text without variables",
			variables: map[string]string{},
			expected:  "Plain text without variables",
		},
		{
			name:    "missing variable unchanged",
			content: "Hello {{unknown}}!",
			variables: map[string]string{
				"name": "World",
			},
			expected: "Hello {{unknown}}!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceVariables(tt.content, tt.variables)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestProcessTemplateFile tests processing a single template file
func TestProcessTemplateFile_Basic(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tempDir, "source.txt")
	sourceContent := "Hello {{name}}, welcome to {{project}}!"
	require.NoError(t, os.WriteFile(sourceFile, []byte(sourceContent), 0644))

	// Process file
	targetFile := filepath.Join(tempDir, "output", "target.txt")
	variables := map[string]string{
		"name":    "Developer",
		"project": "DDx",
	}

	err := processTemplateFile(sourceFile, targetFile, variables)
	require.NoError(t, err)

	// Verify output
	assert.FileExists(t, targetFile)

	content, err := os.ReadFile(targetFile)
	require.NoError(t, err)
	assert.Equal(t, "Hello Developer, welcome to DDx!", string(content))
}

// TestApplyTemplate tests applying a template directory
func TestApplyTemplate_Basic(t *testing.T) {
	t.Parallel()
	// Create template directory
	templateDir := t.TempDir()
	targetDir := t.TempDir()

	// Create template structure
	files := map[string]string{
		"README.md":           "# {{project_name}}",
		"src/main.go":         "package main\n// {{project_name}}",
		"config/settings.yml": "name: {{project_name}}\nport: {{port}}",
	}

	for path, content := range files {
		fullPath := filepath.Join(templateDir, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}

	// Apply template
	variables := map[string]string{
		"project_name": "test-app",
		"port":         "8080",
	}

	err := applyTemplate(templateDir, targetDir, variables)
	require.NoError(t, err)

	// Verify structure
	assert.DirExists(t, filepath.Join(targetDir, "src"))
	assert.DirExists(t, filepath.Join(targetDir, "config"))

	// Verify files
	tests := []struct {
		path     string
		expected string
	}{
		{"README.md", "# test-app"},
		{"src/main.go", "package main\n// test-app"},
		{"config/settings.yml", "name: test-app\nport: 8080"},
	}

	for _, tt := range tests {
		content, err := os.ReadFile(filepath.Join(targetDir, tt.path))
		require.NoError(t, err)
		assert.Equal(t, tt.expected, string(content))
	}
}

// TestList tests listing available templates
func TestList_Basic(t *testing.T) {
	// Create temp home directory
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	// Create DDx templates directory
	templatesDir := filepath.Join(tempHome, ".ddx", "templates")
	require.NoError(t, os.MkdirAll(templatesDir, 0755))

	// Test empty directory
	templates, err := List()
	require.NoError(t, err)
	assert.Empty(t, templates)

	// Create template directories
	templateNames := []string{"template1", "template2", "template3"}
	for _, name := range templateNames {
		require.NoError(t, os.MkdirAll(filepath.Join(templatesDir, name), 0755))
	}

	// Also create a file (should not be listed)
	require.NoError(t, os.WriteFile(filepath.Join(templatesDir, "not-a-template.txt"), []byte("file"), 0644))

	// List templates
	templates, err = List()
	require.NoError(t, err)
	assert.ElementsMatch(t, templateNames, templates)
}

// TestApply tests the main Apply function
func TestApply_Basic(t *testing.T) {
	// Create temp home directory
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	// Create template
	templateName := "test-template"
	templateDir := filepath.Join(tempHome, ".ddx", "templates", templateName)
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Add template files
	templateFile := filepath.Join(templateDir, "app.txt")
	require.NoError(t, os.WriteFile(templateFile, []byte("App: {{app_name}}"), 0644))

	// Apply template
	targetDir := t.TempDir()
	variables := map[string]string{
		"app_name": "MyApp",
	}

	err := Apply(templateName, targetDir, variables)
	require.NoError(t, err)

	// Verify result
	content, err := os.ReadFile(filepath.Join(targetDir, "app.txt"))
	require.NoError(t, err)
	assert.Equal(t, "App: MyApp", string(content))

	// Test non-existent template
	err = Apply("non-existent", targetDir, variables)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// TestApply_EdgeCases tests edge cases for Apply function
func TestApply_EdgeCases(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	// Test empty template name
	err := Apply("", t.TempDir(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template name cannot be empty")

	// Test empty target directory
	err = Apply("test", "", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "target directory cannot be empty")

	// Test template that is a file, not directory
	templateName := "file-template"
	templatePath := filepath.Join(tempHome, ".ddx", "templates", templateName)
	require.NoError(t, os.MkdirAll(filepath.Dir(templatePath), 0755))
	require.NoError(t, os.WriteFile(templatePath, []byte("content"), 0644))

	err = Apply(templateName, t.TempDir(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not a directory")

	// Test nil variables (should work)
	templateDir := filepath.Join(tempHome, ".ddx", "templates", "nil-vars")
	require.NoError(t, os.MkdirAll(templateDir, 0755))
	templateFile := filepath.Join(templateDir, "test.txt")
	require.NoError(t, os.WriteFile(templateFile, []byte("no variables"), 0644))

	err = Apply("nil-vars", t.TempDir(), nil)
	assert.NoError(t, err)
}
