package workflow

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoader_Load tests loading workflow definitions
func TestLoader_Load(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		setupFunc    func(string) // Setup test files
		workflowName string
		wantErr      bool
		errContains  string
		validate     func(*testing.T, *Definition)
	}{
		{
			name: "load valid workflow",
			setupFunc: func(dir string) {
				createValidWorkflow(t, dir, "helix")
			},
			workflowName: "helix",
			wantErr:      false,
			validate: func(t *testing.T, def *Definition) {
				if def.Name != "helix" {
					t.Errorf("Name = %v, want helix", def.Name)
				}
				if def.Version == "" {
					t.Error("Version should not be empty")
				}
			},
		},
		{
			name:         "workflow not found",
			setupFunc:    func(dir string) {}, // No setup
			workflowName: "nonexistent",
			wantErr:      true,
			errContains:  "not found",
		},
		{
			name: "invalid yaml",
			setupFunc: func(dir string) {
				createInvalidWorkflow(t, dir, "broken")
			},
			workflowName: "broken",
			wantErr:      true,
			errContains:  "parse",
		},
		{
			name: "workflow with agent commands",
			setupFunc: func(dir string) {
				createWorkflowWithAgentCommands(t, dir, "helix")
			},
			workflowName: "helix",
			wantErr:      false,
			validate: func(t *testing.T, def *Definition) {
				if len(def.AgentCommands) == 0 {
					t.Error("expected agent commands")
				}

				requestCmd, exists := def.AgentCommands["request"]
				if !exists {
					t.Error("expected 'request' agent command")
				}

				if !requestCmd.Enabled {
					t.Error("request command should be enabled")
				}

				if requestCmd.Action == "" {
					t.Error("request command should have action")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			testLibDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(testLibDir, 0755); err != nil {
				t.Fatalf("failed to create test dir: %v", err)
			}

			if tt.setupFunc != nil {
				tt.setupFunc(testLibDir)
			}

			// Create loader
			loader := NewLoader(testLibDir)

			// Test Load
			def, err := loader.Load(tt.workflowName)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.wantErr && err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, want to contain %v", err.Error(), tt.errContains)
				}
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, def)
			}
		})
	}
}

// TestLoader_MatchesTriggers tests trigger matching logic
func TestLoader_MatchesTriggers(t *testing.T) {
	tmpDir := t.TempDir()
	createWorkflowWithAgentCommands(t, tmpDir, "helix")

	loader := NewLoader(tmpDir)
	def, err := loader.Load("helix")
	if err != nil {
		t.Fatalf("failed to load test workflow: %v", err)
	}

	tests := []struct {
		name       string
		subcommand string
		text       string
		wantMatch  bool
	}{
		// Keyword matches
		{
			name:       "keyword at start",
			subcommand: "request",
			text:       "add pagination to list",
			wantMatch:  true,
		},
		{
			name:       "keyword in middle",
			subcommand: "request",
			text:       "please add pagination",
			wantMatch:  true,
		},
		{
			name:       "keyword at end",
			subcommand: "request",
			text:       "pagination should be added",
			wantMatch:  false, // "added" not "add"
		},
		{
			name:       "keyword exact match",
			subcommand: "request",
			text:       "add",
			wantMatch:  true,
		},
		{
			name:       "keyword case insensitive",
			subcommand: "request",
			text:       "Add pagination",
			wantMatch:  true,
		},
		{
			name:       "keyword partial word no match",
			subcommand: "request",
			text:       "adding pagination",
			wantMatch:  false, // "adding" not "add"
		},
		// Pattern matches
		{
			name:       "pattern US- prefix",
			subcommand: "request",
			text:       "work on US-042",
			wantMatch:  true,
		},
		{
			name:       "pattern work on",
			subcommand: "request",
			text:       "work on the feature",
			wantMatch:  true,
		},
		{
			name:       "pattern case insensitive",
			subcommand: "request",
			text:       "Work On US-042",
			wantMatch:  true,
		},
		// No matches
		{
			name:       "no trigger match",
			subcommand: "request",
			text:       "should we consider pagination?",
			wantMatch:  false,
		},
		{
			name:       "subcommand not found",
			subcommand: "nonexistent",
			text:       "add pagination",
			wantMatch:  false,
		},
		// Multiple keywords
		{
			name:       "multiple keywords fix",
			subcommand: "request",
			text:       "fix the bug in update.go",
			wantMatch:  true,
		},
		{
			name:       "multiple keywords implement",
			subcommand: "request",
			text:       "implement pagination feature",
			wantMatch:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := loader.MatchesTriggers(def, tt.subcommand, tt.text)

			if match != tt.wantMatch {
				t.Errorf("MatchesTriggers() = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}

// TestDefinition_Validate tests workflow definition validation
func TestDefinition_Validate(t *testing.T) {
	tests := []struct {
		name    string
		def     Definition
		wantErr bool
	}{
		{
			name: "valid definition",
			def: Definition{
				Name:    "helix",
				Version: "1.0",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			def: Definition{
				Version: "1.0",
			},
			wantErr: true,
		},
		{
			name: "missing version",
			def: Definition{
				Name: "helix",
			},
			wantErr: true,
		},
		{
			name: "enabled command without action",
			def: Definition{
				Name:    "helix",
				Version: "1.0",
				AgentCommands: map[string]AgentCommand{
					"request": {
						Enabled: true,
						Action:  "", // Missing action
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.def.Validate()

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestDefinition_SupportsAgentCommand tests agent command support checking
func TestDefinition_SupportsAgentCommand(t *testing.T) {
	def := Definition{
		AgentCommands: map[string]AgentCommand{
			"request":  {Enabled: true, Action: "frame-request"},
			"status":   {Enabled: true, Action: "show-status"},
			"disabled": {Enabled: false, Action: "something"},
		},
	}

	tests := []struct {
		name       string
		subcommand string
		want       bool
	}{
		{"enabled command", "request", true},
		{"another enabled", "status", true},
		{"disabled command", "disabled", false},
		{"nonexistent command", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := def.SupportsAgentCommand(tt.subcommand)
			if got != tt.want {
				t.Errorf("SupportsAgentCommand(%q) = %v, want %v", tt.subcommand, got, tt.want)
			}
		})
	}
}

// Helper functions to create test workflows
func createValidWorkflow(t *testing.T, baseDir, name string) {
	workflowDir := filepath.Join(baseDir, "workflows", name)
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatalf("failed to create workflow dir: %v", err)
	}

	content := `name: ` + name + `
version: 1.0.0
description: Test workflow
phases:
  - id: test
    order: 1
    name: Test Phase
`

	workflowFile := filepath.Join(workflowDir, "workflow.yml")
	if err := os.WriteFile(workflowFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write workflow file: %v", err)
	}
}

func createInvalidWorkflow(t *testing.T, baseDir, name string) {
	workflowDir := filepath.Join(baseDir, "workflows", name)
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatalf("failed to create workflow dir: %v", err)
	}

	content := `invalid: yaml: content:
  - broken
  indentation`

	workflowFile := filepath.Join(workflowDir, "workflow.yml")
	if err := os.WriteFile(workflowFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write workflow file: %v", err)
	}
}

func createWorkflowWithAgentCommands(t *testing.T, baseDir, name string) {
	workflowDir := filepath.Join(baseDir, "workflows", name)
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatalf("failed to create workflow dir: %v", err)
	}

	content := `name: ` + name + `
version: 1.0.0
description: Test workflow with agent commands
coordinator: coordinator.md

agent_commands:
  request:
    enabled: true
    triggers:
      keywords:
        - add
        - implement
        - fix
        - create
        - build
      patterns:
        - "US-"
        - "work on"
    action: frame-request
    description: Frame user request

  status:
    enabled: true
    action: show-status
    description: Show workflow status

phases:
  - id: test
    order: 1
    name: Test Phase
`

	workflowFile := filepath.Join(workflowDir, "workflow.yml")
	if err := os.WriteFile(workflowFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write workflow file: %v", err)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
