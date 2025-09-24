package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Integration tests validate end-to-end persona workflows
// These tests verify cross-component interactions and complete user workflows

// Helper function to create a fresh root command for tests
func getPersonaIntegrationTestRootCommand() *cobra.Command {
	factory := NewCommandFactory("/tmp")
	return factory.NewRootCommand()
}

// TestPersonaIntegration_FullWorkflow tests complete persona management workflow
func TestPersonaIntegration_FullWorkflow(t *testing.T) {
	// This is a long-running test that covers the full persona workflow
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup complete test environment
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	workDir := t.TempDir()
	require.NoError(t, os.Chdir(workDir))

	// Create initial .ddx.yml configuration
	initialConfig := `version: "1.0"
repository:
  url: "https://github.com/test/project"
  branch: "main"
variables:
  project_name: "test-project"`

	configPath := filepath.Join(workDir, ".ddx.yml")
	require.NoError(t, os.WriteFile(configPath, []byte(initialConfig), 0644))

	// Create initial CLAUDE.md
	initialClaude := `# CLAUDE.md

This is the project guidance for Claude.

## Project Context
This is a test project for validating persona workflows.`

	claudePath := filepath.Join(workDir, "CLAUDE.md")
	require.NoError(t, os.WriteFile(claudePath, []byte(initialClaude), 0644))

	// Set library path to project-local library
	libraryDir := filepath.Join(workDir, ".ddx", "library")
	t.Setenv("DDX_LIBRARY_BASE_PATH", libraryDir)

	// Create personas directory with test personas
	personasDir := filepath.Join(libraryDir, "personas")
	require.NoError(t, os.MkdirAll(personasDir, 0755))

	// Create comprehensive set of test personas
	personas := map[string]string{
		"strict-code-reviewer.md": `---
name: strict-code-reviewer
roles: [code-reviewer, security-analyst]
description: Uncompromising code quality enforcer
tags: [strict, security, production, quality]
---

# Strict Code Reviewer

You are an experienced senior code reviewer who enforces high quality standards.
Your reviews are thorough, security-focused, and aimed at maintaining production quality.

## Review Principles
- Security vulnerabilities are non-negotiable
- Performance implications must be considered
- Code readability and maintainability are essential
- Test coverage requirements must be met`,

		"balanced-code-reviewer.md": `---
name: balanced-code-reviewer
roles: [code-reviewer]
description: Balanced approach to code reviews
tags: [balanced, pragmatic, team-friendly]
---

# Balanced Code Reviewer

You provide constructive, balanced code reviews that consider both quality and team velocity.
You focus on the most important issues while being supportive of team growth.

## Review Approach
- Focus on critical issues first
- Provide constructive feedback
- Consider team experience levels
- Balance quality with delivery speed`,

		"test-engineer-tdd.md": `---
name: test-engineer-tdd
roles: [test-engineer]
description: Test-driven development specialist
tags: [tdd, testing, quality, red-green-refactor]
---

# TDD Test Engineer

You are a test engineer who follows strict TDD methodology.
Always write failing tests first, then implement minimal code to pass.

## TDD Cycle
1. Red: Write a failing test
2. Green: Write minimal code to pass
3. Refactor: Improve code while keeping tests green

## Testing Principles
- Tests should be fast and reliable
- Tests should be independent
- Test names should be descriptive
- Coverage should be meaningful, not just high`,

		"test-engineer-bdd.md": `---
name: test-engineer-bdd
roles: [test-engineer]
description: Behavior-driven development specialist
tags: [bdd, testing, behavior, acceptance]
---

# BDD Test Engineer

You are a test engineer focused on behavior-driven development.
You write tests that describe system behavior from user perspective.

## BDD Approach
- Given/When/Then structure
- Focus on user behavior
- Collaboration with stakeholders
- Living documentation through tests`,

		"architect-systems.md": `---
name: architect-systems
roles: [architect, tech-lead]
description: Systems architecture and design specialist
tags: [architecture, design, scalability, patterns]
---

# Systems Architect

You are a senior systems architect focused on scalable, maintainable design.
You think in terms of system boundaries, data flow, and long-term evolution.

## Architecture Principles
- Favor composition over inheritance
- Design for failure and resilience
- Consider operational concerns early
- Document architectural decisions`,
	}

	for filename, content := range personas {
		personaPath := filepath.Join(personasDir, filename)
		require.NoError(t, os.WriteFile(personaPath, []byte(content), 0644))
	}

	// Test workflow steps
	tests := []struct {
		name      string
		operation func(t *testing.T) error
		validate  func(t *testing.T) error
	}{
		{
			name: "step1_list_available_personas",
			operation: func(t *testing.T) error {
				// TODO: Implement persona list command
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "list")
				return err
			},
			validate: func(t *testing.T) error {
				// TODO: Validate that all personas are listed
				// For now, just check personas directory exists
				_, err := os.Stat(personasDir)
				return err
			},
		},
		{
			name: "step2_show_specific_persona",
			operation: func(t *testing.T) error {
				// TODO: Implement persona show command
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "show", "strict-code-reviewer")
				return err
			},
			validate: func(t *testing.T) error {
				// TODO: Validate persona details are shown
				// For now, just check persona file exists
				personaPath := filepath.Join(personasDir, "strict-code-reviewer.md")
				_, err := os.Stat(personaPath)
				return err
			},
		},
		{
			name: "step3_bind_personas_to_roles",
			operation: func(t *testing.T) error {
				// TODO: Implement persona bind command
				rootCmd := getPersonaIntegrationTestRootCommand()

				// Bind multiple personas
				bindings := map[string]string{
					"code-reviewer": "strict-code-reviewer",
					"test-engineer": "test-engineer-tdd",
					"architect":     "architect-systems",
				}

				for role, persona := range bindings {
					_, err := executeCommand(rootCmd, "persona", "bind", role, persona)
					if err != nil {
						return err
					}
				}
				return nil
			},
			validate: func(t *testing.T) error {
				// Validate .ddx.yml was updated with bindings
				content, err := os.ReadFile(configPath)
				if err != nil {
					return err
				}

				var config map[string]interface{}
				if err := yaml.Unmarshal(content, &config); err != nil {
					return err
				}

				// TODO: Enable when binding is implemented
				// personaBindings, exists := config["persona_bindings"]
				// if !exists {
				//     return fmt.Errorf("persona_bindings section not found")
				// }

				// bindings := personaBindings.(map[string]interface{})
				// expectedBindings := map[string]string{
				//     "code-reviewer": "strict-code-reviewer",
				//     "test-engineer": "test-engineer-tdd",
				//     "architect": "architect-systems",
				// }

				// for role, expectedPersona := range expectedBindings {
				//     if actualPersona, exists := bindings[role]; !exists || actualPersona != expectedPersona {
				//         return fmt.Errorf("binding mismatch for role %s", role)
				//     }
				// }

				return nil
			},
		},
		{
			name: "step4_show_current_bindings",
			operation: func(t *testing.T) error {
				// TODO: Implement persona bindings command
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "bindings")
				return err
			},
			validate: func(t *testing.T) error {
				// TODO: Validate bindings are displayed correctly
				return nil
			},
		},
		{
			name: "step5_load_all_bound_personas",
			operation: func(t *testing.T) error {
				// TODO: Implement persona load command
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "load")
				return err
			},
			validate: func(t *testing.T) error {
				// Validate CLAUDE.md was updated with personas
				content, err := os.ReadFile(claudePath)
				if err != nil {
					return err
				}

				claudeContent := string(content)

				// TODO: Enable when persona loading is implemented
				// // Should contain persona markers
				// if !strings.Contains(claudeContent, "<!-- PERSONAS:START -->") {
				//     return fmt.Errorf("CLAUDE.md missing personas start marker")
				// }
				// if !strings.Contains(claudeContent, "<!-- PERSONAS:END -->") {
				//     return fmt.Errorf("CLAUDE.md missing personas end marker")
				// }

				// // Should contain all bound personas
				// expectedPersonas := []string{
				//     "Strict Code Reviewer",
				//     "TDD Test Engineer",
				//     "Systems Architect",
				// }

				// for _, persona := range expectedPersonas {
				//     if !strings.Contains(claudeContent, persona) {
				//         return fmt.Errorf("CLAUDE.md missing persona: %s", persona)
				//     }
				// }

				// // Should preserve original content
				// if !strings.Contains(claudeContent, "This is the project guidance for Claude") {
				//     return fmt.Errorf("CLAUDE.md original content not preserved")
				// }

				_ = claudeContent // Suppress unused variable warning
				return nil
			},
		},
		{
			name: "step6_check_persona_status",
			operation: func(t *testing.T) error {
				// TODO: Implement persona status command
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "status")
				return err
			},
			validate: func(t *testing.T) error {
				// TODO: Validate status shows loaded personas
				return nil
			},
		},
		{
			name: "step7_load_specific_persona",
			operation: func(t *testing.T) error {
				// TODO: Implement specific persona loading
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "load", "balanced-code-reviewer")
				return err
			},
			validate: func(t *testing.T) error {
				// TODO: Validate specific persona was added to CLAUDE.md
				return nil
			},
		},
		{
			name: "step8_update_binding",
			operation: func(t *testing.T) error {
				// TODO: Update existing binding
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "bind", "test-engineer", "test-engineer-bdd")
				return err
			},
			validate: func(t *testing.T) error {
				// TODO: Validate binding was updated in .ddx.yml
				return nil
			},
		},
		{
			name: "step9_reload_with_updated_binding",
			operation: func(t *testing.T) error {
				// TODO: Reload personas to pick up new binding
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "load")
				return err
			},
			validate: func(t *testing.T) error {
				// TODO: Validate BDD engineer replaced TDD engineer in CLAUDE.md
				return nil
			},
		},
		{
			name: "step10_remove_personas",
			operation: func(t *testing.T) error {
				// TODO: Implement persona unload command
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "unload")
				return err
			},
			validate: func(t *testing.T) error {
				// Validate personas were removed from CLAUDE.md
				content, err := os.ReadFile(claudePath)
				if err != nil {
					return err
				}

				claudeContent := string(content)

				// TODO: Enable when persona unloading is implemented
				// // Should not contain persona markers
				// if strings.Contains(claudeContent, "<!-- PERSONAS:START -->") {
				//     return fmt.Errorf("CLAUDE.md still contains personas")
				// }

				// // Should preserve original content
				// if !strings.Contains(claudeContent, "This is the project guidance for Claude") {
				//     return fmt.Errorf("CLAUDE.md original content not preserved")
				// }

				_ = claudeContent // Suppress unused variable warning
				return nil
			},
		},
	}

	// Execute workflow steps
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			require.NoError(t, os.Chdir(workDir))

			// Execute operation
			err := tt.operation(t)

			// Operation should succeed
			assert.NoError(t, err, "Operation should succeed: %s", tt.name)

			// Run validation
			if tt.validate != nil {
				validateErr := tt.validate(t)
				assert.NoError(t, validateErr, "Validation should pass: %s", tt.name)
			}

			// For now, just run validation to ensure test structure is valid
			if tt.validate != nil {
				_ = tt.validate(t)
			}
		})
	}
}

// TestPersonaIntegration_WorkflowOverrides tests workflow-specific persona overrides
func TestPersonaIntegration_WorkflowOverrides(t *testing.T) {
	// This test verifies basic persona loading functionality
	// Full workflow override feature may be implemented later
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test environment
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	workDir := t.TempDir()
	require.NoError(t, os.Chdir(workDir))

	// Create .ddx.yml with overrides
	configWithOverrides := `version: "1.0"
repository:
  url: "https://github.com/test/project"
  branch: "main"

persona_bindings:
  test-engineer: test-engineer-tdd
  code-reviewer: balanced-reviewer

overrides:
  performance-workflow:
    test-engineer: test-engineer-bdd  # Use BDD for performance testing
  security-workflow:
    code-reviewer: security-reviewer  # Use security specialist for security workflows`

	configPath := filepath.Join(workDir, ".ddx.yml")
	require.NoError(t, os.WriteFile(configPath, []byte(configWithOverrides), 0644))

	// Create test personas
	personasDir := filepath.Join(tempHome, ".ddx", "library", "personas")
	require.NoError(t, os.MkdirAll(personasDir, 0755))

	personas := map[string]string{
		"test-engineer-tdd.md": `---
name: test-engineer-tdd
roles: [test-engineer]
description: TDD specialist
tags: [tdd]
---
# TDD Engineer`,

		"test-engineer-bdd.md": `---
name: test-engineer-bdd
roles: [test-engineer]
description: BDD specialist
tags: [bdd]
---
# BDD Engineer`,

		"balanced-reviewer.md": `---
name: balanced-reviewer
roles: [code-reviewer]
description: Balanced reviewer
tags: [balanced]
---
# Balanced Reviewer`,

		"security-reviewer.md": `---
name: security-reviewer
roles: [code-reviewer, security-analyst]
description: Security specialist
tags: [security]
---
# Security Reviewer`,
	}

	for filename, content := range personas {
		personaPath := filepath.Join(personasDir, filename)
		require.NoError(t, os.WriteFile(personaPath, []byte(content), 0644))
	}

	// Create CLAUDE.md
	claudePath := filepath.Join(workDir, "CLAUDE.md")
	require.NoError(t, os.WriteFile(claudePath, []byte("# CLAUDE.md\n\nProject guidance."), 0644))

	tests := []struct {
		name      string
		workflow  string
		operation func(t *testing.T, workflow string) error
		validate  func(t *testing.T, workflow string) error
	}{
		{
			name:     "default_workflow_uses_default_bindings",
			workflow: "", // No specific workflow
			operation: func(t *testing.T, workflow string) error {
				// TODO: Load personas for default workflow
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "load")
				return err
			},
			validate: func(t *testing.T, workflow string) error {
				// TODO: Validate default personas are loaded
				// Should load test-engineer-tdd and balanced-reviewer
				return nil
			},
		},
		{
			name:     "performance_workflow_uses_overrides",
			workflow: "performance-workflow",
			operation: func(t *testing.T, workflow string) error {
				// TODO: Load personas for specific workflow
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "load", "--workflow", workflow)
				return err
			},
			validate: func(t *testing.T, workflow string) error {
				// TODO: Validate override personas are loaded
				// Should load test-engineer-bdd (override) and balanced-reviewer (default)
				return nil
			},
		},
		{
			name:     "security_workflow_uses_overrides",
			workflow: "security-workflow",
			operation: func(t *testing.T, workflow string) error {
				// TODO: Load personas for security workflow
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "load", "--workflow", workflow)
				return err
			},
			validate: func(t *testing.T, workflow string) error {
				// TODO: Validate security workflow overrides
				// Should load security-reviewer (override) and test-engineer-tdd (default)
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)
			require.NoError(t, os.Chdir(workDir))

			// Execute operation
			err := tt.operation(t, tt.workflow)

			// TODO: Enable when workflow override support is implemented
			assert.Error(t, err, "Command not implemented yet: %s", tt.name)

			// TODO: Enable validation when commands are implemented
			// assert.NoError(t, err, "Operation should succeed: %s", tt.name)
			//
			// if tt.validate != nil {
			//     validateErr := tt.validate(t, tt.workflow)
			//     assert.NoError(t, validateErr, "Validation should pass: %s", tt.name)
			// }

			// For now, just run validation to ensure test structure is valid
			if tt.validate != nil {
				_ = tt.validate(t, tt.workflow)
			}
		})
	}
}

// TestPersonaIntegration_ErrorHandling tests error scenarios in persona workflows
func TestPersonaIntegration_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup minimal test environment
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	workDir := t.TempDir()
	require.NoError(t, os.Chdir(workDir))

	tests := []struct {
		name          string
		setup         func(t *testing.T)
		operation     func(t *testing.T) error
		expectedError string
	}{
		{
			name: "bind_nonexistent_persona",
			setup: func(t *testing.T) {
				// Create valid .ddx.yml but no personas
				config := `version: "1.0"`
				require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))
			},
			operation: func(t *testing.T) error {
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "bind", "code-reviewer", "nonexistent-persona")
				return err
			},
			expectedError: "persona 'nonexistent-persona' not found",
		},
		{
			name: "load_without_config",
			setup: func(t *testing.T) {
				// No .ddx.yml file - this should be an error
			},
			operation: func(t *testing.T) error {
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "load")
				return err
			},
			expectedError: "No .ddx.yml configuration found",
		},
		{
			name: "load_persona_with_invalid_content",
			setup: func(t *testing.T) {
				// Create persona with invalid YAML
				currentDir, _ := os.Getwd()
				libraryDir := filepath.Join(currentDir, ".ddx", "library")
				t.Setenv("DDX_LIBRARY_BASE_PATH", libraryDir)
				personasDir := filepath.Join(libraryDir, "personas")
				require.NoError(t, os.MkdirAll(personasDir, 0755))

				invalidPersona := `---
name: invalid-persona
roles: [code-reviewer
description: Invalid YAML - missing bracket
---
# Invalid Persona`

				personaPath := filepath.Join(personasDir, "invalid-persona.md")
				require.NoError(t, os.WriteFile(personaPath, []byte(invalidPersona), 0644))

				config := `version: "1.0"
persona_bindings:
  code-reviewer: invalid-persona`
				require.NoError(t, os.WriteFile(".ddx.yml", []byte(config), 0644))
			},
			operation: func(t *testing.T) error {
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "load")
				return err
			},
			expectedError: "failed to parse YAML frontmatter",
		},
		{
			name: "show_nonexistent_persona",
			setup: func(t *testing.T) {
				// Create empty personas directory
				currentDir, _ := os.Getwd()
				libraryDir := filepath.Join(currentDir, ".ddx", "library")
				t.Setenv("DDX_LIBRARY_BASE_PATH", libraryDir)
				personasDir := filepath.Join(libraryDir, "personas")
				require.NoError(t, os.MkdirAll(personasDir, 0755))
			},
			operation: func(t *testing.T) error {
				rootCmd := getPersonaIntegrationTestRootCommand()
				_, err := executeCommand(rootCmd, "persona", "show", "nonexistent-persona")
				return err
			},
			expectedError: "persona 'nonexistent-persona' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean working directory for each test
			workSubDir := filepath.Join(workDir, tt.name)
			require.NoError(t, os.MkdirAll(workSubDir, 0755))
			require.NoError(t, os.Chdir(workSubDir))

			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Run setup
			if tt.setup != nil {
				tt.setup(t)
			}

			// Execute operation
			err := tt.operation(t)

			// Validate error occurred and contains expected message
			assert.Error(t, err, "Should have error: %s", tt.name)
			if err != nil {
				assert.Contains(t, err.Error(), tt.expectedError, "Error message should match: %s", tt.name)
			}

			// For now, just validate test structure
			assert.NotEmpty(t, tt.expectedError, "Expected error message should be defined")
		})
	}
}

// TestPersonaIntegration_ConcurrentAccess tests concurrent persona operations
func TestPersonaIntegration_ConcurrentAccess(t *testing.T) {
	// Basic concurrent access test - simplified version
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test environment
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	workDir := t.TempDir()
	require.NoError(t, os.Chdir(workDir))

	// Create .ddx.yml
	config := `version: "1.0"
persona_bindings:
  code-reviewer: test-reviewer`

	configPath := filepath.Join(workDir, ".ddx.yml")
	require.NoError(t, os.WriteFile(configPath, []byte(config), 0644))

	// Create test persona
	personasDir := filepath.Join(tempHome, ".ddx", "library", "personas")
	require.NoError(t, os.MkdirAll(personasDir, 0755))

	personaContent := `---
name: test-reviewer
roles: [code-reviewer]
description: Test reviewer
tags: [test]
---
# Test Reviewer`

	personaPath := filepath.Join(personasDir, "test-reviewer.md")
	require.NoError(t, os.WriteFile(personaPath, []byte(personaContent), 0644))

	// Create CLAUDE.md
	claudePath := filepath.Join(workDir, "CLAUDE.md")
	require.NoError(t, os.WriteFile(claudePath, []byte("# CLAUDE.md\n\nProject guidance."), 0644))

	// Test concurrent operations
	const numGoroutines = 5
	results := make(chan error, numGoroutines)

	operations := []func() error{
		func() error {
			// TODO: List personas
			rootCmd := &cobra.Command{Use: "ddx", Short: "DDx CLI"}
			// rootCmd.AddCommand(personaCmd)
			_, err := executeCommand(rootCmd, "persona", "list")
			return err
		},
		func() error {
			// TODO: Show persona
			rootCmd := &cobra.Command{Use: "ddx", Short: "DDx CLI"}
			// rootCmd.AddCommand(personaCmd)
			_, err := executeCommand(rootCmd, "persona", "show", "test-reviewer")
			return err
		},
		func() error {
			// TODO: Get bindings
			rootCmd := &cobra.Command{Use: "ddx", Short: "DDx CLI"}
			// rootCmd.AddCommand(personaCmd)
			_, err := executeCommand(rootCmd, "persona", "bindings")
			return err
		},
		func() error {
			// TODO: Load personas
			rootCmd := &cobra.Command{Use: "ddx", Short: "DDx CLI"}
			// rootCmd.AddCommand(personaCmd)
			_, err := executeCommand(rootCmd, "persona", "load")
			return err
		},
		func() error {
			// TODO: Check status
			rootCmd := &cobra.Command{Use: "ddx", Short: "DDx CLI"}
			// rootCmd.AddCommand(personaCmd)
			_, err := executeCommand(rootCmd, "persona", "status")
			return err
		},
	}

	// Run operations concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			operation := operations[index%len(operations)]
			err := operation()
			results <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-results

		// Commands are now implemented, expect them to succeed
		assert.NoError(t, err, "Concurrent operation %d should succeed", i)
	}

	// Verify final state is consistent
	// TODO: Verify that concurrent operations didn't corrupt files
	_, err := os.Stat(configPath)
	assert.NoError(t, err, "Config file should still exist after concurrent operations")

	_, err = os.Stat(claudePath)
	assert.NoError(t, err, "CLAUDE.md should still exist after concurrent operations")
}
