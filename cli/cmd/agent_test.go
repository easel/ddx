package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestAgentRequest_NoConfig tests behavior when no config exists
func TestAgentRequest_NoConfig(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))

	factory := NewCommandFactory(env.Dir)
	cmd := factory.NewRootCommand()
	output := &bytes.Buffer{}
	cmd.SetOut(output)
	cmd.SetErr(output)
	cmd.SetArgs([]string{"agent", "request", "add", "pagination"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	got := output.String()
	if !strings.Contains(got, "NO_HANDLER") {
		t.Errorf("output = %q, want NO_HANDLER", got)
	}
}

// TestAgentRequest_NoActiveWorkflows tests behavior with config but no active workflows
func TestAgentRequest_NoActiveWorkflows(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	env.CreateConfig(`version: "1.0"
library:
  path: ".ddx/library"
`)

	output, err := env.RunCommand("agent", "request", "add", "pagination")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "NO_HANDLER") {
		t.Errorf("output = %q, want NO_HANDLER", output)
	}
}

// TestAgentRequest_SafeWord tests safe word bypass
func TestAgentRequest_SafeWord(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithWorkflow(t, env, "helix")
	createHelixWorkflow(t, env)

	tests := []struct {
		name        string
		args        []string
		wantContain []string
	}{
		{
			name: "safe word with space",
			args: []string{"agent", "request", "NODDX", "add", "pagination"},
			wantContain: []string{
				"NO_HANDLER",
				"SAFE_WORD: NODDX",
				"MESSAGE: add pagination",
			},
		},
		{
			name: "safe word with colon",
			args: []string{"agent", "request", "NODDX:", "add", "pagination"},
			wantContain: []string{
				"NO_HANDLER",
				"SAFE_WORD: NODDX",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := env.RunCommand(tt.args...)
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

// TestAgentRequest_TriggerMatch tests trigger-based workflow activation
func TestAgentRequest_TriggerMatch(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithWorkflow(t, env, "helix")
	createHelixWorkflow(t, env)

	tests := []struct {
		name        string
		args        []string
		wantContain []string
		wantAbsent  []string
	}{
		{
			name: "add keyword triggers",
			args: []string{"agent", "request", "add", "pagination", "to", "list"},
			wantContain: []string{
				"WORKFLOW: helix",
				"SUBCOMMAND: request",
				"ACTION: frame-request",
				"COMMAND:",
			},
		},
		{
			name: "fix keyword triggers",
			args: []string{"agent", "request", "fix", "bug", "in", "update.go"},
			wantContain: []string{
				"WORKFLOW: helix",
				"ACTION: frame-request",
			},
		},
		{
			name: "US- pattern triggers",
			args: []string{"agent", "request", "work", "on", "US-042"},
			wantContain: []string{
				"WORKFLOW: helix",
				"ACTION: frame-request",
			},
		},
		{
			name: "no trigger match",
			args: []string{"agent", "request", "should", "we", "add", "pagination?"},
			wantContain: []string{
				"NO_HANDLER",
			},
			wantAbsent: []string{
				"WORKFLOW:",
			},
		},
		{
			name: "discussion context",
			args: []string{"agent", "request", "what", "do", "you", "think", "about", "pagination?"},
			wantContain: []string{
				"NO_HANDLER",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := env.RunCommand(tt.args...)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(output, want) {
					t.Errorf("output = %q, want to contain %q", output, want)
				}
			}

			for _, unwanted := range tt.wantAbsent {
				if strings.Contains(output, unwanted) {
					t.Errorf("output = %q, should not contain %q", output, unwanted)
				}
			}
		})
	}
}

// TestAgentRequest_MultipleWorkflows tests first-match priority
func TestAgentRequest_MultipleWorkflows(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))

	// Create config with two workflows in specific order
	createConfigWithWorkflows(t, env, []string{"helix", "kanban"})
	createHelixWorkflow(t, env)
	createKanbanWorkflow(t, env)

	tests := []struct {
		name         string
		args         []string
		wantWorkflow string
	}{
		{
			name:         "helix trigger matches first",
			args:         []string{"agent", "request", "add", "pagination"},
			wantWorkflow: "helix", // helix is first in active list
		},
		{
			name:         "kanban-only trigger",
			args:         []string{"agent", "request", "create", "card", "for", "task"},
			wantWorkflow: "kanban",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := env.RunCommand(tt.args...)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			expectedWorkflow := "WORKFLOW: " + tt.wantWorkflow
			if !strings.Contains(output, expectedWorkflow) {
				t.Errorf("output = %q, want to contain %q", output, expectedWorkflow)
			}
		})
	}
}

// TestAgentRequest_CommandQuoting tests proper quoting in COMMAND output
func TestAgentRequest_CommandQuoting(t *testing.T) {
	env := NewTestEnvironment(t, WithGitInit(false))
	createConfigWithWorkflow(t, env, "helix")
	createHelixWorkflow(t, env)

	// Test with special characters that need quoting
	output, err := env.RunCommand("agent", "request", "add", "pagination", "with", "quotes\"and'stuff")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Command should be quoted
	if !strings.Contains(output, "COMMAND:") {
		t.Error("output should contain COMMAND:")
	}

	// Should contain quoted request (implementation will determine exact quoting)
	if !strings.Contains(output, "ddx workflow helix execute") {
		t.Error("COMMAND should contain workflow execution")
	}
}

// Helper functions

func createConfigWithWorkflow(t *testing.T, env *TestEnvironment, workflow string) {
	createConfigWithWorkflows(t, env, []string{workflow})
}

func createConfigWithWorkflows(t *testing.T, env *TestEnvironment, workflows []string) {
	configContent := `version: "1.0"

workflows:
  active:
`
	for _, wf := range workflows {
		configContent += "    - " + wf + "\n"
	}

	configContent += `  safe_word: "NODDX"

library:
  path: ".ddx/library"
`

	env.CreateConfig(configContent)
}

func createHelixWorkflow(t *testing.T, env *TestEnvironment) {
	env.CreateFile(filepath.Join(".ddx", "library", "workflows", "helix", "workflow.yml"), `name: helix
version: 1.0.0
description: HELIX workflow
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

phases:
  - id: frame
    order: 1
    name: Frame
`)
}

func createKanbanWorkflow(t *testing.T, env *TestEnvironment) {
	env.CreateFile(filepath.Join(".ddx", "library", "workflows", "kanban", "workflow.yml"), `name: kanban
version: 1.0.0
description: Kanban workflow

agent_commands:
  request:
    enabled: true
    triggers:
      keywords:
        - card
      patterns:
        - "create card"
    action: create-card
    description: Create kanban card

phases:
  - id: todo
    order: 1
    name: Todo
`)
}
