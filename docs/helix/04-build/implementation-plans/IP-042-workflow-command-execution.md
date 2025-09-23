---
title: "Implementation Plan - Workflow Command Execution"
type: implementation-plan
feature_id: FEAT-005
user_story: US-042
workflow_phase: build
artifact_type: implementation-plan
tags:
  - helix/build
  - helix/artifact/implementation
  - helix/phase/build
  - workflow
  - command-execution
related:
  - "[[US-042-workflow-command-execution]]"
  - "[[TS-005-workflow-command-execution]]"
  - "[[FEAT-005-workflow-execution-engine]]"
status: ready
priority: P0
created: 2025-01-22
updated: 2025-01-22
---

# Implementation Plan: US-042 Workflow Command Execution

## Overview

This implementation plan defines the step-by-step approach for implementing workflow command execution functionality. The implementation must follow TDD principles, making failing tests pass incrementally.

## Current State Analysis

### Existing Implementation
- ✅ Basic workflow command exists in `cli/cmd/workflow.go`
- ✅ Supports `status`, `list`, `activate`, `advance` subcommands
- ✅ Contract tests exist in `mcp_contract_test.go`
- ❌ No support for workflow-specific subcommands
- ❌ No dynamic command discovery from library
- ❌ No command execution functionality

### Required Changes
1. Extend workflow command routing to support workflow-specific subcommands
2. Implement command discovery from `library/workflows/<name>/commands/`
3. Implement command execution with argument passing
4. Add comprehensive error handling
5. Maintain backward compatibility with existing workflow commands

## Implementation Strategy

### Phase 1: TDD Red Phase - Write Failing Tests

#### Step 1.1: Create Workflow Command Test File
**File**: `cli/cmd/workflow_test.go`

```go
package cmd

import (
    "bytes"
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestWorkflowCommandDiscovery tests workflow command discovery functionality
func TestWorkflowCommandDiscovery(t *testing.T) {
    // This test will initially fail - no implementation exists
    tests := []struct {
        name     string
        workflow string
        setup    func(t *testing.T) string
        expected []string
    }{
        {
            name:     "discover_helix_commands",
            workflow: "helix",
            setup:    setupHelixWorkflowCommands,
            expected: []string{"build-story", "continue", "status"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            workDir := tt.setup(t)
            defer os.Chdir(workDir)

            commands, err := discoverWorkflowCommands(tt.workflow)

            require.NoError(t, err)
            for _, expectedCmd := range tt.expected {
                assert.Contains(t, commands, expectedCmd)
            }
        })
    }
}

// TestWorkflowCommandExecution tests workflow command execution
func TestWorkflowCommandExecution(t *testing.T) {
    // This test will initially fail - no implementation exists
    tests := []struct {
        name     string
        workflow string
        command  string
        args     []string
        setup    func(t *testing.T) string
        expected string
    }{
        {
            name:     "execute_build_story",
            workflow: "helix",
            command:  "build-story",
            args:     []string{"US-001"},
            setup:    setupHelixWorkflowCommands,
            expected: "HELIX Command: Build Story",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            workDir := tt.setup(t)
            defer os.Chdir(workDir)

            output, err := executeWorkflowCommand(tt.workflow, tt.command, tt.args)

            require.NoError(t, err)
            assert.Contains(t, output, tt.expected)
        })
    }
}

// Helper function to setup helix workflow commands
func setupHelixWorkflowCommands(t *testing.T) string {
    workDir := t.TempDir()
    require.NoError(t, os.Chdir(workDir))

    commandsDir := filepath.Join(workDir, "library", "workflows", "helix", "commands")
    require.NoError(t, os.MkdirAll(commandsDir, 0755))

    // Create build-story command
    buildStoryContent := `# HELIX Command: Build Story

You are a HELIX workflow executor tasked with implementing work on a specific user story.`
    require.NoError(t, os.WriteFile(
        filepath.Join(commandsDir, "build-story.md"),
        []byte(buildStoryContent), 0644))

    return workDir
}
```

#### Step 1.2: Add Acceptance Tests
**File**: `cli/cmd/acceptance_test.go` (extend existing file)

Add the `TestAcceptance_US042_WorkflowCommandExecution` test from TS-005 specification.

#### Step 1.3: Extend Contract Tests
**File**: `cli/cmd/mcp_contract_test.go` (extend existing file)

```go
func TestWorkflowCommands_Contract(t *testing.T) {
    // ... existing tests ...

    t.Run("contract_workflow_helix_commands", func(t *testing.T) {
        // Given: HELIX workflow with commands available
        origDir, _ := os.Getwd()
        defer os.Chdir(origDir)

        tempDir := t.TempDir()
        os.Chdir(tempDir)
        setupHelixWorkflowCommands(t)

        // When: Listing HELIX commands
        cmd := getFreshRootCmd()
        buf := new(bytes.Buffer)
        cmd.SetOut(buf)
        cmd.SetErr(buf)
        cmd.SetArgs([]string{"workflow", "helix", "commands"})

        err := cmd.Execute()

        // Then: Should list available commands
        if err == nil {
            output := buf.String()
            assert.Contains(t, output, "Available commands")
            assert.Contains(t, output, "build-story")
        }
    })

    t.Run("contract_workflow_helix_execute", func(t *testing.T) {
        // Given: HELIX workflow with commands available
        origDir, _ := os.Getwd()
        defer os.Chdir(origDir)

        tempDir := t.TempDir()
        os.Chdir(tempDir)
        setupHelixWorkflowCommands(t)

        // When: Executing HELIX command
        cmd := getFreshRootCmd()
        buf := new(bytes.Buffer)
        cmd.SetOut(buf)
        cmd.SetErr(buf)
        cmd.SetArgs([]string{"workflow", "helix", "execute", "build-story", "US-001"})

        err := cmd.Execute()

        // Then: Should execute command
        if err == nil {
            output := buf.String()
            assert.Contains(t, output, "HELIX Command")
        }
    })
}
```

### Phase 2: TDD Green Phase - Implement Functionality

#### Step 2.1: Modify Workflow Command Router
**File**: `cli/cmd/workflow.go`

```go
// Enhanced runWorkflow function to support workflow-specific commands
func runWorkflow(cmd *cobra.Command, args []string) error {
    workflowForce, _ := cmd.Flags().GetBool("force")

    if len(args) == 0 {
        return showWorkflowStatus(cmd)
    }

    firstArg := strings.ToLower(args[0])

    // Check if first argument is a known workflow name
    if isKnownWorkflow(firstArg) && len(args) > 1 {
        return handleWorkflowSpecificCommand(cmd, firstArg, args[1:])
    }

    // Handle generic workflow commands (existing functionality)
    switch firstArg {
    case "status":
        return showWorkflowStatus(cmd)
    case "list":
        return listWorkflows(cmd)
    case "activate":
        if len(args) < 2 {
            return fmt.Errorf("workflow name required")
        }
        return activateWorkflow(cmd, args[1], workflowForce)
    case "advance":
        return advanceWorkflow(cmd)
    default:
        return fmt.Errorf("unknown subcommand: %s", firstArg)
    }
}
```

#### Step 2.2: Implement Workflow Detection
**File**: `cli/cmd/workflow.go`

```go
// isKnownWorkflow checks if the given name is a known workflow
func isKnownWorkflow(name string) bool {
    workflowDir := filepath.Join("library", "workflows", name)
    if stat, err := os.Stat(workflowDir); err == nil && stat.IsDir() {
        return true
    }
    return false
}
```

#### Step 2.3: Implement Workflow-Specific Command Handler
**File**: `cli/cmd/workflow.go`

```go
// handleWorkflowSpecificCommand routes workflow-specific subcommands
func handleWorkflowSpecificCommand(cmd *cobra.Command, workflow string, args []string) error {
    if len(args) == 0 {
        return fmt.Errorf("subcommand required for workflow %s", workflow)
    }

    subcommand := strings.ToLower(args[0])
    switch subcommand {
    case "commands":
        return listWorkflowCommands(cmd, workflow)
    case "execute":
        if len(args) < 2 {
            return fmt.Errorf("command name required for execute")
        }
        return executeWorkflowCommand(cmd, workflow, args[1], args[2:])
    default:
        return fmt.Errorf("unknown subcommand '%s' for workflow '%s'", subcommand, workflow)
    }
}
```

#### Step 2.4: Implement Command Discovery
**File**: `cli/cmd/workflow.go`

```go
// listWorkflowCommands lists available commands for a workflow
func listWorkflowCommands(cmd *cobra.Command, workflow string) error {
    commandsDir := filepath.Join("library", "workflows", workflow, "commands")

    // Check if commands directory exists
    if _, err := os.Stat(commandsDir); os.IsNotExist(err) {
        return fmt.Errorf("workflow '%s' not found or has no commands", workflow)
    }

    entries, err := os.ReadDir(commandsDir)
    if err != nil {
        return fmt.Errorf("failed to read commands directory: %w", err)
    }

    fmt.Fprintf(cmd.OutOrStdout(), "Available commands for %s workflow:\n\n", workflow)

    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
            commandName := strings.TrimSuffix(entry.Name(), ".md")

            // Try to read the first line for description
            description := getCommandDescription(filepath.Join(commandsDir, entry.Name()))

            fmt.Fprintf(cmd.OutOrStdout(), "  %-15s %s\n", commandName, description)
        }
    }

    return nil
}

// getCommandDescription extracts description from command file
func getCommandDescription(filePath string) string {
    content, err := os.ReadFile(filePath)
    if err != nil {
        return "No description available"
    }

    lines := strings.Split(string(content), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "# ") {
            return strings.TrimPrefix(line, "# ")
        }
    }

    return "No description available"
}
```

#### Step 2.5: Implement Command Execution
**File**: `cli/cmd/workflow.go`

```go
// executeWorkflowCommand loads and displays a workflow command
func executeWorkflowCommand(cmd *cobra.Command, workflow, command string, args []string) error {
    commandPath := filepath.Join("library", "workflows", workflow, "commands", command+".md")

    // Check if command file exists
    if _, err := os.Stat(commandPath); os.IsNotExist(err) {
        return fmt.Errorf("command '%s' not found in workflow '%s'", command, workflow)
    }

    // Read command content
    content, err := os.ReadFile(commandPath)
    if err != nil {
        return fmt.Errorf("failed to read command file: %w", err)
    }

    // Display command content
    fmt.Fprintf(cmd.OutOrStdout(), "Executing %s workflow command: %s\n\n", workflow, command)

    if len(args) > 0 {
        fmt.Fprintf(cmd.OutOrStdout(), "Command Arguments: %v\n\n", args)
    }

    fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(content))

    return nil
}
```

### Phase 3: Test Integration and Validation

#### Step 3.1: Run Tests and Fix Issues
```bash
cd cli
go test ./cmd -run TestWorkflow -v
go test ./cmd -run TestAcceptance_US042 -v
```

#### Step 3.2: Validate Contract Tests
```bash
go test ./cmd -run TestWorkflowCommands_Contract -v
```

#### Step 3.3: Integration Testing
```bash
# Test the actual CLI commands
./build/ddx workflow helix commands
./build/ddx workflow helix execute build-story US-001
```

## Error Handling Implementation

### Validation Functions

```go
// validateWorkflowExists checks if workflow directory exists
func validateWorkflowExists(workflow string) error {
    workflowDir := filepath.Join("library", "workflows", workflow)
    if _, err := os.Stat(workflowDir); os.IsNotExist(err) {
        return fmt.Errorf("workflow '%s' not found", workflow)
    }
    return nil
}

// validateCommandExists checks if command file exists
func validateCommandExists(workflow, command string) error {
    commandPath := filepath.Join("library", "workflows", workflow, "commands", command+".md")
    if _, err := os.Stat(commandPath); os.IsNotExist(err) {
        return fmt.Errorf("command '%s' not found in workflow '%s'", command, workflow)
    }
    return nil
}
```

## Backward Compatibility

### Ensure Existing Functionality Works
- All existing workflow commands must continue working
- `ddx workflow status`, `ddx workflow list`, etc. remain unchanged
- No breaking changes to existing command structure

### Migration Strategy
- New functionality is additive only
- Existing tests must continue passing
- CLI help text updated to include new subcommands

## Performance Considerations

### File System Operations
- Cache workflow discovery results
- Minimize file system calls
- Use efficient directory traversal

### Error Handling
- Fail fast on invalid workflows
- Provide clear error messages
- Handle missing files gracefully

## Testing Strategy

### Test Execution Order
1. Unit tests for individual functions
2. Integration tests for command routing
3. Acceptance tests for user scenarios
4. Contract tests for CLI interface

### Test Coverage Goals
- 100% of acceptance criteria covered
- All error paths tested
- Command discovery edge cases
- File system error scenarios

## Definition of Done

- [ ] All tests pass (Red → Green phase complete)
- [ ] `ddx workflow helix commands` lists available commands
- [ ] `ddx workflow helix execute build-story US-001` loads and displays command
- [ ] Error handling for invalid workflows and commands
- [ ] Backward compatibility maintained
- [ ] Documentation updated
- [ ] Integration tests validate end-to-end functionality

## Implementation Timeline

1. **Day 1**: Write failing tests (TDD Red phase)
2. **Day 2**: Implement command routing and discovery
3. **Day 3**: Implement command execution and error handling
4. **Day 4**: Integration testing and refinement
5. **Day 5**: Documentation and final validation

This implementation plan follows TDD principles and ensures robust, well-tested functionality that meets all acceptance criteria for US-042.