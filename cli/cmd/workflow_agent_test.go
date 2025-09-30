package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWorkflowActivate tests workflow activation
func TestWorkflowActivate(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(*testing.T, *TestEnvironment)
		workflow    string
		wantErr     bool
		wantContain string
	}{
		{
			name: "activate valid workflow",
			setupFunc: func(t *testing.T, env *TestEnvironment) {
				createConfigWithoutWorkflows(t, env)
				createHelixWorkflow(t, env)
			},
			workflow:    "helix",
			wantErr:     false,
			wantContain: "✓ Activated helix workflow",
		},
		{
			name: "activate nonexistent workflow",
			setupFunc: func(t *testing.T, env *TestEnvironment) {
				createConfigWithoutWorkflows(t, env)
			},
			workflow:    "nonexistent",
			wantErr:     true,
			wantContain: "not found",
		},
		{
			name: "activate already active workflow",
			setupFunc: func(t *testing.T, env *TestEnvironment) {
				createConfigWithWorkflow(t, env, "helix")
				createHelixWorkflow(t, env)
			},
			workflow:    "helix",
			wantErr:     false,
			wantContain: "already active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewTestEnvironment(t, WithGitInit(false))
			tt.setupFunc(t, env)

			// Run activate command via CLI
			output, err := env.RunCommand("workflow", "activate", tt.workflow)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.wantContain != "" && !strings.Contains(output, tt.wantContain) {
				t.Errorf("output = %q, want to contain %q", output, tt.wantContain)
			}

			// Verify config was updated (if no error)
			if !tt.wantErr && err == nil && !strings.Contains(output, "already active") {
				verifyWorkflowActive(t, env.Dir, tt.workflow)
			}
		})
	}
}

// TestWorkflowDeactivate tests workflow deactivation
func TestWorkflowDeactivate(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(*testing.T, *TestEnvironment)
		workflow    string
		wantErr     bool
		wantContain string
	}{
		{
			name: "deactivate active workflow",
			setupFunc: func(t *testing.T, env *TestEnvironment) {
				createConfigWithWorkflow(t, env, "helix")
				createHelixWorkflow(t, env)
			},
			workflow:    "helix",
			wantErr:     false,
			wantContain: "✓ Deactivated helix workflow",
		},
		{
			name: "deactivate inactive workflow",
			setupFunc: func(t *testing.T, env *TestEnvironment) {
				createConfigWithoutWorkflows(t, env)
			},
			workflow:    "helix",
			wantErr:     false,
			wantContain: "not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewTestEnvironment(t, WithGitInit(false))
			tt.setupFunc(t, env)

			// Run deactivate command via CLI
			output, err := env.RunCommand("workflow", "deactivate", tt.workflow)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.wantContain != "" && !strings.Contains(output, tt.wantContain) {
				t.Errorf("output = %q, want to contain %q", output, tt.wantContain)
			}

			// Verify config was updated (if deactivated)
			if !tt.wantErr && err == nil && !strings.Contains(output, "not active") {
				verifyWorkflowNotActive(t, env.Dir, tt.workflow)
			}
		})
	}
}

// TestWorkflowActivate_Priority tests workflow priority ordering
func TestWorkflowActivate_Priority(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithoutWorkflows(t, env)
	createHelixWorkflow(t, env)
	createKanbanWorkflow(t, env)

	// Activate helix (should be priority 1)
	output, err := env.RunCommand("workflow", "activate", "helix")
	if err != nil {
		t.Fatalf("failed to activate helix: %v", err)
	}

	if !strings.Contains(output, "Priority: 1 of 1") {
		t.Errorf("expected Priority: 1 of 1, got %q", output)
	}

	// Activate kanban (should be priority 2)
	output, err = env.RunCommand("workflow", "activate", "kanban")
	if err != nil {
		t.Fatalf("failed to activate kanban: %v", err)
	}

	if !strings.Contains(output, "Priority: 2 of 2") {
		t.Errorf("expected Priority: 2 of 2, got %q", output)
	}

	// Verify order in config
	verifyWorkflowOrder(t, env.Dir, []string{"helix", "kanban"})
}

// TestWorkflowDeactivate_PreservesOrder tests that deactivation preserves order
func TestWorkflowDeactivate_PreservesOrder(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithWorkflows(t, env, []string{"helix", "kanban", "tdd"})
	createHelixWorkflow(t, env)
	createKanbanWorkflow(t, env)

	// Deactivate middle workflow
	_, err := env.RunCommand("workflow", "deactivate", "kanban")
	if err != nil {
		t.Fatalf("failed to deactivate: %v", err)
	}

	// Verify remaining order: helix, tdd
	verifyWorkflowOrder(t, env.Dir, []string{"helix", "tdd"})
}

// TestWorkflowStatus tests workflow status display
func TestWorkflowStatus(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(*testing.T, *TestEnvironment)
		wantContain []string
	}{
		{
			name: "no active workflows",
			setupFunc: func(t *testing.T, env *TestEnvironment) {
				createConfigWithoutWorkflows(t, env)
			},
			wantContain: []string{"No active workflows"},
		},
		{
			name: "single workflow",
			setupFunc: func(t *testing.T, env *TestEnvironment) {
				createConfigWithWorkflow(t, env, "helix")
				createHelixWorkflow(t, env)
			},
			wantContain: []string{
				"Active workflows (in priority order):",
				"1. helix",
				"Agent commands:",
				"request",
				"frame-request",
				"Safe word: NODDX",
			},
		},
		{
			name: "multiple workflows",
			setupFunc: func(t *testing.T, env *TestEnvironment) {
				createConfigWithWorkflows(t, env, []string{"helix", "kanban"})
				createHelixWorkflow(t, env)
				createKanbanWorkflow(t, env)
			},
			wantContain: []string{
				"1. helix",
				"2. kanban",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := NewTestEnvironment(t, WithGitInit(false))
			tt.setupFunc(t, env)

			// Run status command via CLI
			output, err := env.RunCommand("workflow", "status")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(output, want) {
					t.Errorf("output = %q, want to contain %q", output, want)
				}
			}
		})
	}
}

// TestWorkflowStatus_ShowsTriggers tests that status shows trigger keywords
func TestWorkflowStatus_ShowsTriggers(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithWorkflow(t, env, "helix")
	createHelixWorkflow(t, env)

	// Run status command via CLI
	output, err := env.RunCommand("workflow", "status")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should show triggers for request command
	wantTriggers := []string{
		"Keywords:",
		"add",
		"implement",
		"fix",
		"Patterns:",
		"US-",
		"work on",
	}

	for _, want := range wantTriggers {
		if !strings.Contains(output, want) {
			t.Errorf("output = %q, want to contain %q", output, want)
		}
	}
}

// Helper functions

func createConfigWithoutWorkflows(t *testing.T, env *TestEnvironment) {
	env.CreateConfig(`version: "1.0"
library:
  path: ".ddx/library"
`)
}

func verifyWorkflowActive(t *testing.T, tmpDir, workflow string) {
	configPath := filepath.Join(tmpDir, ".ddx", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "active:") {
		t.Error("config should contain active: section")
	}

	if !strings.Contains(content, "- "+workflow) {
		t.Errorf("config should contain workflow %s in active list", workflow)
	}
}

func verifyWorkflowNotActive(t *testing.T, tmpDir, workflow string) {
	configPath := filepath.Join(tmpDir, ".ddx", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	content := string(data)

	// Check that workflow is not in active list
	lines := strings.Split(content, "\n")
	inActiveSection := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "active:" {
			inActiveSection = true
			continue
		}
		if inActiveSection {
			if strings.HasPrefix(trimmed, "-") {
				if strings.Contains(trimmed, workflow) {
					t.Errorf("workflow %s should not be in active list", workflow)
				}
			} else if !strings.HasPrefix(trimmed, "#") && trimmed != "" {
				// End of active section
				break
			}
		}
	}
}

func verifyWorkflowOrder(t *testing.T, tmpDir string, expectedOrder []string) {
	configPath := filepath.Join(tmpDir, ".ddx", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	var actualOrder []string
	inActiveSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "active:" {
			inActiveSection = true
			continue
		}
		if inActiveSection {
			if strings.HasPrefix(trimmed, "- ") {
				workflow := strings.TrimPrefix(trimmed, "- ")
				workflow = strings.TrimSpace(workflow)
				actualOrder = append(actualOrder, workflow)
			} else if !strings.HasPrefix(trimmed, "#") && trimmed != "" {
				// End of active section
				break
			}
		}
	}

	if len(actualOrder) != len(expectedOrder) {
		t.Errorf("workflow count = %d, want %d", len(actualOrder), len(expectedOrder))
	}

	for i, expected := range expectedOrder {
		if i >= len(actualOrder) {
			t.Errorf("missing workflow at position %d: %s", i, expected)
			continue
		}
		if actualOrder[i] != expected {
			t.Errorf("workflow at position %d = %s, want %s", i, actualOrder[i], expected)
		}
	}
}
