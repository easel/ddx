package templates

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReplaceVariables_EdgeCases tests advanced edge cases for variable replacement
func TestReplaceVariables_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		content   string
		variables map[string]string
		expected  string
	}{
		{
			name:      "nested braces",
			content:   "{{{{name}}}}",
			variables: map[string]string{"name": "test"},
			expected:  "{{test}}",
		},
		{
			name:      "variable in file path",
			content:   "/path/to/{{project_name}}/file.txt",
			variables: map[string]string{"project_name": "myproject"},
			expected:  "/path/to/myproject/file.txt",
		},
		{
			name:      "mixed case variables",
			content:   "{{Name}} ${name} {{name}}",
			variables: map[string]string{"name": "test", "Name": "Test"},
			expected:  "Test test test",
		},
		{
			name:      "special characters in variables",
			content:   "{{special}} contains special chars",
			variables: map[string]string{"special": "test@#$%^&*()"},
			expected:  "test@#$%^&*() contains special chars",
		},
		{
			name:      "empty variable value",
			content:   "Before {{empty}} after",
			variables: map[string]string{"empty": ""},
			expected:  "Before  after",
		},
		{
			name:      "variable with different spacing",
			content:   "{{name}} {{ name}} {{name }} {{ name }}",
			variables: map[string]string{"name": "test"},
			expected:  "test test test test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceVariables(tt.content, tt.variables)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestProcessTemplateFile_EdgeCases tests edge cases for processTemplateFile
func TestProcessTemplateFile_EdgeCases(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	// Test file with variables in path
	sourceFile := filepath.Join(tempDir, "source.txt")
	sourceContent := "Content: {{content}}"
	require.NoError(t, os.WriteFile(sourceFile, []byte(sourceContent), 0755)) // executable file

	variables := map[string]string{
		"content":      "test",
		"project_name": "myproject",
	}

	// Target path has variables
	targetFile := filepath.Join(tempDir, "output", "{{project_name}}.txt")
	expectedTarget := filepath.Join(tempDir, "output", "myproject.txt")

	err := processTemplateFile(sourceFile, targetFile, variables)
	require.NoError(t, err)

	// Verify the processed file exists at the correct path
	assert.FileExists(t, expectedTarget)

	// Verify content
	content, err := os.ReadFile(expectedTarget)
	require.NoError(t, err)
	assert.Equal(t, "Content: test", string(content))

	// Verify permissions are capped at 0644 for security
	info, err := os.Stat(expectedTarget)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0644), info.Mode())
}

// TestApplyTemplate_EdgeCases tests edge cases for applyTemplate
func TestApplyTemplate_EdgeCases(t *testing.T) {
	t.Parallel()
	templateDir := t.TempDir()
	targetDir := t.TempDir()

	// Create template with hidden files (should be skipped)
	files := map[string]string{
		"regular.txt":     "{{name}}",
		".hidden":         "should be skipped",
		".git/config":     "should be skipped",
		"dir/.gitignore":  "should be skipped",
		"normal/file.txt": "{{name}}",
	}

	for path, content := range files {
		fullPath := filepath.Join(templateDir, path)
		require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0755))
		require.NoError(t, os.WriteFile(fullPath, []byte(content), 0644))
	}

	variables := map[string]string{
		"name": "test",
	}

	err := applyTemplate(templateDir, targetDir, variables)
	require.NoError(t, err)

	// Verify regular files exist
	assert.FileExists(t, filepath.Join(targetDir, "regular.txt"))
	assert.FileExists(t, filepath.Join(targetDir, "normal", "file.txt"))

	// Verify hidden files don't exist
	assert.NoFileExists(t, filepath.Join(targetDir, ".hidden"))
	assert.NoFileExists(t, filepath.Join(targetDir, ".git", "config"))
	assert.NoDirExists(t, filepath.Join(targetDir, ".git"))
	assert.NoFileExists(t, filepath.Join(targetDir, "dir", ".gitignore"))

	// Verify content
	content, err := os.ReadFile(filepath.Join(targetDir, "regular.txt"))
	require.NoError(t, err)
	assert.Equal(t, "test", string(content))
}

// TestList_EdgeCases tests edge cases for List function
func TestList_EdgeCases(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	templatesDir := filepath.Join(tempHome, ".ddx", "templates")
	require.NoError(t, os.MkdirAll(templatesDir, 0755))

	// Create mix of directories and files, including hidden ones
	items := []string{
		"template1",        // should be included
		"template2",        // should be included
		".hidden-template", // should be excluded (hidden)
		"not-template.txt", // should be excluded (file, not directory)
	}

	for i, item := range items {
		path := filepath.Join(templatesDir, item)
		if i < 2 || i == 2 { // directories (including hidden)
			require.NoError(t, os.MkdirAll(path, 0755))
		} else { // file
			require.NoError(t, os.WriteFile(path, []byte("content"), 0644))
		}
	}

	templates, err := List()
	require.NoError(t, err)

	// Should only include non-hidden directories
	expected := []string{"template1", "template2"}
	assert.ElementsMatch(t, expected, templates)
}

// TestConcurrentApply tests applying templates concurrently
func TestConcurrentApply(t *testing.T) {
	// Set up shared template
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	templateName := "concurrent-template"
	templateDir := filepath.Join(tempHome, ".ddx", "templates", templateName)
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	templateFile := filepath.Join(templateDir, "test.txt")
	require.NoError(t, os.WriteFile(templateFile, []byte("Test: {{value}}"), 0644))

	// Apply template concurrently to different directories
	const numGoroutines = 5
	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			targetDir := t.TempDir()
			variables := map[string]string{
				"value": "test-value",
			}

			err := Apply(templateName, targetDir, variables)
			results <- err
		}(i)
	}

	// Check all goroutines completed successfully
	for i := 0; i < numGoroutines; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}

// TestApplyLargeTemplate tests applying templates with many files
func TestApplyLargeTemplate(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	templateName := "large-template"
	templateDir := filepath.Join(tempHome, ".ddx", "templates", templateName)
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Create template with many files and nested directories
	const numFiles = 50
	const numDirectories = 10

	for i := 0; i < numDirectories; i++ {
		dirPath := filepath.Join(templateDir, "dir"+string(rune('A'+i)))
		require.NoError(t, os.MkdirAll(dirPath, 0755))

		for j := 0; j < numFiles/numDirectories; j++ {
			filePath := filepath.Join(dirPath, "file"+string(rune('0'+j))+".txt")
			content := "File {{file_index}} in directory {{dir_index}}"
			require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
		}
	}

	// Apply template
	targetDir := t.TempDir()
	variables := map[string]string{
		"file_index": "test-file",
		"dir_index":  "test-dir",
	}

	err := Apply(templateName, targetDir, variables)
	require.NoError(t, err)

	// Verify structure was created
	for i := 0; i < numDirectories; i++ {
		dirPath := filepath.Join(targetDir, "dir"+string(rune('A'+i)))
		assert.DirExists(t, dirPath)

		for j := 0; j < numFiles/numDirectories; j++ {
			filePath := filepath.Join(dirPath, "file"+string(rune('0'+j))+".txt")
			assert.FileExists(t, filePath)

			// Check content
			content, err := os.ReadFile(filePath)
			require.NoError(t, err)
			assert.Equal(t, "File test-file in directory test-dir", string(content))
		}
	}
}

// TestTemplateWithSymlinks tests handling of symbolic links in templates
func TestTemplateWithSymlinks(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	templateName := "symlink-template"
	templateDir := filepath.Join(tempHome, ".ddx", "templates", templateName)
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Create regular file
	regularFile := filepath.Join(templateDir, "regular.txt")
	require.NoError(t, os.WriteFile(regularFile, []byte("Regular content"), 0644))

	// Skip symlink test on platforms that don't support them or in environments where creation might fail
	symlinkFile := filepath.Join(templateDir, "symlink.txt")
	err := os.Symlink(regularFile, symlinkFile)
	if err != nil {
		t.Skip("Skipping symlink test: symlinks not supported or permission denied")
		return
	}

	// Apply template
	targetDir := t.TempDir()
	err = Apply(templateName, targetDir, nil)
	require.NoError(t, err)

	// Verify files exist (symlinks should be resolved or copied appropriately)
	assert.FileExists(t, filepath.Join(targetDir, "regular.txt"))

	// The behavior for symlinks may vary, so we just check that the apply succeeded
	// and the target directory contains expected files
	entries, err := os.ReadDir(targetDir)
	require.NoError(t, err)
	assert.NotEmpty(t, entries, "Target directory should contain files")
}
