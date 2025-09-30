package cmd

import (
	"strings"
	"testing"
)

// TestIntegration_FullWorkflow tests complete workflow: activate -> request -> execute
func TestIntegration_FullWorkflow(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithoutWorkflows(t, env)
	createHelixWorkflow(t, env)

	// Step 1: Activate workflow
	activateOutput, err := env.RunCommand("workflow", "activate", "helix")
	if err != nil {
		t.Fatalf("failed to activate workflow: %v", err)
	}

	if !strings.Contains(activateOutput, "✓ Activated helix workflow") {
		t.Errorf("activation failed: %s", activateOutput)
	}

	// Step 2: Check status
	statusStr, err := env.RunCommand("workflow", "status")
	if err != nil {
		t.Fatalf("failed to show status: %v", err)
	}

	if !strings.Contains(statusStr, "1. helix") {
		t.Errorf("status should show helix: %s", statusStr)
	}

	// Step 3: Make agent request
	requestStr, err := env.RunCommand("agent", "request", "add", "pagination")
	if err != nil {
		t.Fatalf("agent request failed: %v", err)
	}
	expectedOutputs := []string{
		"WORKFLOW: helix",
		"SUBCOMMAND: request",
		"ACTION: frame-request",
		"COMMAND: ddx workflow helix execute frame-request",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(requestStr, expected) {
			t.Errorf("agent output missing %q, got: %s", expected, requestStr)
		}
	}

	// Step 4: Deactivate workflow
	_, err = env.RunCommand("workflow", "deactivate", "helix")
	if err != nil {
		t.Fatalf("failed to deactivate: %v", err)
	}

	// Step 5: Verify agent returns NO_HANDLER after deactivation
	requestOutput2, err := env.RunCommand("agent", "request", "add", "pagination")
	if err != nil {
		t.Fatalf("agent request failed: %v", err)
	}

	if !strings.Contains(requestOutput2, "NO_HANDLER") {
		t.Errorf("should return NO_HANDLER after deactivation, got: %s", requestOutput2)
	}
}

// TestIntegration_MultipleWorkflowPriority tests priority handling
func TestIntegration_MultipleWorkflowPriority(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithoutWorkflows(t, env)
	createHelixWorkflow(t, env)
	createKanbanWorkflow(t, env)

	// Activate in specific order: helix then kanban
	_, err := env.RunCommand("workflow", "activate", "helix")
	if err != nil {
		t.Fatalf("failed to activate helix: %v", err)
	}

	_, err = env.RunCommand("workflow", "activate", "kanban")
	if err != nil {
		t.Fatalf("failed to activate kanban: %v", err)
	}

	// Request with trigger that matches both
	// "add" matches helix, but helix should win (first in list)
	got, err := env.RunCommand("agent", "request", "add", "feature")
	if err != nil {
		t.Fatalf("agent request failed: %v", err)
	}

	if !strings.Contains(got, "WORKFLOW: helix") {
		t.Errorf("helix should handle request (first match), got: %s", got)
	}

	// Now test kanban-specific trigger
	// "card" only matches kanban
	got2, err := env.RunCommand("agent", "request", "create", "card")
	if err != nil {
		t.Fatalf("agent request failed: %v", err)
	}

	if !strings.Contains(got2, "WORKFLOW: kanban") {
		t.Errorf("kanban should handle card request, got: %s", got2)
	}
}

// TestIntegration_SafeWordBypass tests safe word in real scenario
func TestIntegration_SafeWordBypass(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithWorkflow(t, env, "helix")
	createHelixWorkflow(t, env)

	// Normal request triggers workflow
	normalStr, err := env.RunCommand("agent", "request", "add", "pagination")
	if err != nil {
		t.Fatalf("normal request failed: %v", err)
	}

	if !strings.Contains(normalStr, "WORKFLOW: helix") {
		t.Errorf("normal request should trigger workflow: %s", normalStr)
	}

	// Safe word request bypasses workflow
	safeStr, err := env.RunCommand("agent", "request", "NODDX", "add", "pagination")
	if err != nil {
		t.Fatalf("safe word request failed: %v", err)
	}
	expectedSafeOutput := []string{
		"NO_HANDLER",
		"SAFE_WORD: NODDX",
		"MESSAGE: add pagination",
	}

	for _, expected := range expectedSafeOutput {
		if !strings.Contains(safeStr, expected) {
			t.Errorf("safe word output missing %q, got: %s", expected, safeStr)
		}
	}
}

// TestIntegration_CustomSafeWord tests custom safe word configuration
func TestIntegration_CustomSafeWord(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))

	// Create config with custom safe word
	createConfigWithCustomSafeWord(t, env, "helix", "SKIP")
	createHelixWorkflow(t, env)

	// Test custom safe word
	got, err := env.RunCommand("agent", "request", "SKIP", "add", "pagination")
	if err != nil {
		t.Fatalf("safe word request failed: %v", err)
	}

	if !strings.Contains(got, "SAFE_WORD: SKIP") {
		t.Errorf("should use custom safe word SKIP, got: %s", got)
	}

	// Test that NODDX doesn't work with custom safe word
	got2, err := env.RunCommand("agent", "request", "NODDX", "add", "pagination")
	if err != nil {
		t.Fatalf("NODDX request failed: %v", err)
	}

	// NODDX should be treated as regular text, so workflow should trigger
	if !strings.Contains(got2, "WORKFLOW: helix") {
		t.Errorf("NODDX should not be safe word with custom config, got: %s", got2)
	}
}

// TestIntegration_InvalidWorkflowGracefulDegradation tests error handling
func TestIntegration_InvalidWorkflowGracefulDegradation(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))

	// Create config with nonexistent workflow
	createConfigWithWorkflow(t, env, "nonexistent")

	// Agent request should gracefully return NO_HANDLER
	got, err := env.RunCommand("agent", "request", "add", "pagination")
	if err != nil {
		t.Fatalf("agent should not error on invalid workflow: %v", err)
	}

	if !strings.Contains(got, "NO_HANDLER") {
		t.Errorf("should gracefully degrade to NO_HANDLER, got: %s", got)
	}
}

// Helper function for custom safe word config
func createConfigWithCustomSafeWord(t *testing.T, env *TestEnvironment, workflow, safeWord string) {
	configContent := `version: "1.0"

workflows:
  active:
    - ` + workflow + `
  safe_word: "` + safeWord + `"

library:
  path: ".ddx/library"
`

	env.CreateConfig(configContent)
}
