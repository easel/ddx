package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getPromptsTestRootCommand creates a root command for testing
func getPromptsTestRootCommand() *cobra.Command {
	factory := NewCommandFactory()
	return factory.NewRootCommand()
}

func TestPromptsCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		setup       func(t *testing.T) (cleanup func())
		validate    func(t *testing.T, output string, err error)
		expectError bool
	}{
		{
			name: "prompts list - shows available prompts",
			args: []string{"prompts", "list"},
			setup: func(t *testing.T) func() {
				// Create test library structure
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))

				// Create library with prompts
				promptsDir := filepath.Join(testDir, "library", "prompts")
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "claude"), 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "common"), 0755))

				// Create some prompt files
				require.NoError(t, os.WriteFile(
					filepath.Join(promptsDir, "claude", "code-review.md"),
					[]byte("# Code Review Prompt\nReview this code..."),
					0644,
				))
				require.NoError(t, os.WriteFile(
					filepath.Join(promptsDir, "common", "refactor.md"),
					[]byte("# Refactor Prompt\nRefactor this code..."),
					0644,
				))

				// Create .ddx.yml pointing to library
				configContent := `version: "2.0"
library_path: ./library
repository:
  url: https://github.com/easel/ddx
  branch: main`
				require.NoError(t, os.WriteFile(".ddx.yml", []byte(configContent), 0644))

				return func() {
					os.Chdir(origWd)
				}
			},
			validate: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				assert.Contains(t, output, "claude")
				assert.Contains(t, output, "common")
			},
			expectError: false,
		},
		{
			name: "prompts list verbose - shows files recursively",
			args: []string{"prompts", "list", "--verbose"},
			setup: func(t *testing.T) func() {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))

				// Create library with nested prompts
				promptsDir := filepath.Join(testDir, "library", "prompts")
				claudeDir := filepath.Join(promptsDir, "claude", "system-prompts")
				require.NoError(t, os.MkdirAll(claudeDir, 0755))

				// Create nested prompt files
				require.NoError(t, os.WriteFile(
					filepath.Join(claudeDir, "security.md"),
					[]byte("# Security Review"),
					0644,
				))
				require.NoError(t, os.WriteFile(
					filepath.Join(promptsDir, "claude", "general.md"),
					[]byte("# General Claude Prompt"),
					0644,
				))

				// Create config
				configContent := `version: "2.0"
library_path: ./library
repository:
  url: https://github.com/easel/ddx`
				require.NoError(t, os.WriteFile(".ddx.yml", []byte(configContent), 0644))

				return func() {
					os.Chdir(origWd)
				}
			},
			validate: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				// Should show files, not just directories
				assert.Contains(t, output, "security.md")
				assert.Contains(t, output, "general.md")
			},
			expectError: false,
		},
		{
			name: "prompts show - displays specific prompt",
			args: []string{"prompts", "show", "claude/code-review"},
			setup: func(t *testing.T) func() {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))

				// Create library with prompt
				promptPath := filepath.Join(testDir, "library", "prompts", "claude")
				require.NoError(t, os.MkdirAll(promptPath, 0755))

				promptContent := `# Code Review Prompt

You are a senior code reviewer. Focus on:
- Security vulnerabilities
- Performance issues
- Code maintainability`

				require.NoError(t, os.WriteFile(
					filepath.Join(promptPath, "code-review.md"),
					[]byte(promptContent),
					0644,
				))

				// Create config
				configContent := `version: "2.0"
library_path: ./library`
				require.NoError(t, os.WriteFile(".ddx.yml", []byte(configContent), 0644))

				return func() {
					os.Chdir(origWd)
				}
			},
			validate: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				assert.Contains(t, output, "Code Review Prompt")
				assert.Contains(t, output, "Security vulnerabilities")
			},
			expectError: false,
		},
		{
			name: "prompts show - error on non-existent prompt",
			args: []string{"prompts", "show", "nonexistent/prompt"},
			setup: func(t *testing.T) func() {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))

				// Create library but no prompts
				require.NoError(t, os.MkdirAll(filepath.Join(testDir, "library", "prompts"), 0755))

				configContent := `version: "2.0"
library_path: ./library`
				require.NoError(t, os.WriteFile(".ddx.yml", []byte(configContent), 0644))

				return func() {
					os.Chdir(origWd)
				}
			},
			validate: func(t *testing.T, output string, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
			},
			expectError: true,
		},
		{
			name: "prompts list with search",
			args: []string{"prompts", "list", "--search", "review"},
			setup: func(t *testing.T) func() {
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))

				// Create library with various prompts
				promptsDir := filepath.Join(testDir, "library", "prompts")
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "claude"), 0755))
				require.NoError(t, os.MkdirAll(filepath.Join(promptsDir, "common"), 0755))

				// Create prompts with different names
				require.NoError(t, os.WriteFile(
					filepath.Join(promptsDir, "claude", "code-review.md"),
					[]byte("# Code Review"),
					0644,
				))
				require.NoError(t, os.WriteFile(
					filepath.Join(promptsDir, "claude", "security-review.md"),
					[]byte("# Security Review"),
					0644,
				))
				require.NoError(t, os.WriteFile(
					filepath.Join(promptsDir, "common", "refactor.md"),
					[]byte("# Refactor"),
					0644,
				))

				configContent := `version: "2.0"
library_path: ./library`
				require.NoError(t, os.WriteFile(".ddx.yml", []byte(configContent), 0644))

				return func() {
					os.Chdir(origWd)
				}
			},
			validate: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				// Should show review-related prompts
				assert.Contains(t, output, "code-review")
				assert.Contains(t, output, "security-review")
				// Should not show non-matching prompt
				assert.NotContains(t, output, "refactor")
			},
			expectError: false,
		},
		{
			name: "prompts list - uses development library",
			args: []string{"prompts", "list"},
			setup: func(t *testing.T) func() {
				// Simulate DDx development environment
				testDir := t.TempDir()
				origWd, _ := os.Getwd()
				require.NoError(t, os.Chdir(testDir))

				// Create git repo
				require.NoError(t, os.MkdirAll(".git", 0755))

				// Create cli/main.go to identify as DDx repo
				require.NoError(t, os.MkdirAll("cli", 0755))
				require.NoError(t, os.WriteFile("cli/main.go", []byte("package main"), 0644))

				// Create library/prompts
				promptsDir := filepath.Join(testDir, "library", "prompts", "ddx")
				require.NoError(t, os.MkdirAll(promptsDir, 0755))
				require.NoError(t, os.WriteFile(
					filepath.Join(promptsDir, "workflow.md"),
					[]byte("# DDx Workflow"),
					0644,
				))

				// No .ddx.yml - should use development mode

				return func() {
					os.Chdir(origWd)
				}
			},
			validate: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				assert.Contains(t, output, "ddx")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setup(t)
			defer cleanup()

			// This test defines the expected behavior for prompts commands
			// The actual implementation will be done after tests are written

			// Execute the command
			rootCmd := getPromptsTestRootCommand()
			output, err := executeCommand(rootCmd, tt.args...)

			// Validate results
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.validate != nil {
				tt.validate(t, output, err)
			}
		})
	}
}

func TestPromptsCommand_Help(t *testing.T) {
	// This test specifies that prompts command should have help text
	// with list and show subcommands
	rootCmd := getPromptsTestRootCommand()
	output, err := executeCommand(rootCmd, "prompts", "--help")

	assert.NoError(t, err)
	assert.Contains(t, output, "Manage AI prompts")
	assert.Contains(t, output, "list")
	assert.Contains(t, output, "show")
}
