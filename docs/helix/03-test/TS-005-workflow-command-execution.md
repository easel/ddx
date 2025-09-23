---
title: "Test Specification - Workflow Command Execution"
type: test-specification
feature_id: FEAT-005
user_story: US-042
workflow_phase: test
artifact_type: test-specification
tags:
  - helix/test
  - helix/artifact/test
  - helix/phase/test
  - workflow
  - command-execution
related:
  - "[[FEAT-005-workflow-execution-engine]]"
  - "[[US-042-workflow-command-execution]]"
  - "[[SD-042-workflow-command-execution]]"
status: draft
priority: P0
created: 2025-01-22
updated: 2025-01-22
---

# Test Specification: US-042 Workflow Command Execution

## Test Strategy Overview

This test specification defines comprehensive test scenarios for workflow command execution. Following Test-Driven Development principles, these tests must be written and failing BEFORE any implementation begins.

### Test Categories

1. **Acceptance Tests** - User story validation (Red phase requirement)
2. **Integration Tests** - Workflow command discovery and execution
3. **Contract Tests** - CLI command interface testing
4. **Error Handling Tests** - Invalid workflow and command scenarios

## Acceptance Test Specifications

### 1. US-042: Workflow Command Execution

#### Test Suite: `TestAcceptance_US042_WorkflowCommandExecution`

```go
// Path: cli/cmd/acceptance_test.go

func TestAcceptance_US042_WorkflowCommandExecution(t *testing.T) {
    tests := []struct {
        name     string
        scenario string
        given    func(t *testing.T) string
        when     func(t *testing.T, workDir string) (string, error)
        then     func(t *testing.T, workDir string, output string, err error)
    }{
        {
            name:     "list_helix_commands",
            scenario: "AC-001: Command Discovery",
            given: func(t *testing.T) string {
                // Given: I have the HELIX workflow available
                workDir := t.TempDir()
                require.NoError(t, os.Chdir(workDir))

                // Create library structure with helix commands
                commandsDir := filepath.Join(workDir, "library", "workflows", "helix", "commands")
                require.NoError(t, os.MkdirAll(commandsDir, 0755))

                // Create build-story command
                buildStoryContent := `# HELIX Command: Build Story

You are a HELIX workflow executor...`
                require.NoError(t, os.WriteFile(
                    filepath.Join(commandsDir, "build-story.md"),
                    []byte(buildStoryContent), 0644))

                // Create continue command
                continueContent := `# HELIX Command: Continue

Continue work on current story...`
                require.NoError(t, os.WriteFile(
                    filepath.Join(commandsDir, "continue.md"),
                    []byte(continueContent), 0644))

                return workDir
            },
            when: func(t *testing.T, workDir string) (string, error) {
                // When: I run `ddx workflow helix commands`
                rootCmd := getTestRootCommand()
                buf := new(bytes.Buffer)
                rootCmd.SetOut(buf)
                rootCmd.SetErr(buf)
                rootCmd.SetArgs([]string{"workflow", "helix", "commands"})

                err := rootCmd.Execute()
                return buf.String(), err
            },
            then: func(t *testing.T, workDir string, output string, err error) {
                // Then: I see a list of available commands with descriptions
                assert.NoError(t, err)
                assert.Contains(t, output, "Available commands for helix workflow:")
                assert.Contains(t, output, "build-story")
                assert.Contains(t, output, "continue")
            },
        },
        {
            name:     "execute_build_story_command",
            scenario: "AC-002: Command Execution",
            given: func(t *testing.T) string {
                // Given: I have a workflow with commands available
                workDir := t.TempDir()
                require.NoError(t, os.Chdir(workDir))

                commandsDir := filepath.Join(workDir, "library", "workflows", "helix", "commands")
                require.NoError(t, os.MkdirAll(commandsDir, 0755))

                buildStoryContent := `# HELIX Command: Build Story

You are a HELIX workflow executor tasked with implementing work on a specific user story.

## Command Input

You will receive a user story ID as an argument (e.g., US-001, US-042, etc.).`
                require.NoError(t, os.WriteFile(
                    filepath.Join(commandsDir, "build-story.md"),
                    []byte(buildStoryContent), 0644))

                return workDir
            },
            when: func(t *testing.T, workDir string) (string, error) {
                // When: I run `ddx workflow helix execute build-story US-001`
                rootCmd := getTestRootCommand()
                buf := new(bytes.Buffer)
                rootCmd.SetOut(buf)
                rootCmd.SetErr(buf)
                rootCmd.SetArgs([]string{"workflow", "helix", "execute", "build-story", "US-001"})

                err := rootCmd.Execute()
                return buf.String(), err
            },
            then: func(t *testing.T, workDir string, output string, err error) {
                // Then: The build-story command prompt is loaded and displayed
                assert.NoError(t, err)
                assert.Contains(t, output, "HELIX Command: Build Story")
                assert.Contains(t, output, "You are a HELIX workflow executor")
                assert.Contains(t, output, "Command Arguments: [US-001]")
            },
        },
        {
            name:     "invalid_workflow_error",
            scenario: "AC-003: Error Handling - Invalid Workflow",
            given: func(t *testing.T) string {
                // Given: I specify a non-existent workflow
                workDir := t.TempDir()
                require.NoError(t, os.Chdir(workDir))
                return workDir
            },
            when: func(t *testing.T, workDir string) (string, error) {
                // When: I run `ddx workflow invalid commands`
                rootCmd := getTestRootCommand()
                buf := new(bytes.Buffer)
                rootCmd.SetOut(buf)
                rootCmd.SetErr(buf)
                rootCmd.SetArgs([]string{"workflow", "invalid", "commands"})

                err := rootCmd.Execute()
                return buf.String(), err
            },
            then: func(t *testing.T, workDir string, output string, err error) {
                // Then: I receive an error message about the workflow not being found
                assert.Error(t, err)
                assert.Contains(t, err.Error(), "workflow 'invalid' not found")
            },
        },
        {
            name:     "invalid_command_error",
            scenario: "AC-004: Error Handling - Invalid Command",
            given: func(t *testing.T) string {
                // Given: I specify a non-existent command
                workDir := t.TempDir()
                require.NoError(t, os.Chdir(workDir))

                commandsDir := filepath.Join(workDir, "library", "workflows", "helix", "commands")
                require.NoError(t, os.MkdirAll(commandsDir, 0755))

                return workDir
            },
            when: func(t *testing.T, workDir string) (string, error) {
                // When: I run `ddx workflow helix execute invalid-command`
                rootCmd := getTestRootCommand()
                buf := new(bytes.Buffer)
                rootCmd.SetOut(buf)
                rootCmd.SetErr(buf)
                rootCmd.SetArgs([]string{"workflow", "helix", "execute", "invalid-command"})

                err := rootCmd.Execute()
                return buf.String(), err
            },
            then: func(t *testing.T, workDir string, output string, err error) {
                // Then: I receive an error about the command not being found
                assert.Error(t, err)
                assert.Contains(t, err.Error(), "command 'invalid-command' not found")
            },
        },
        {
            name:     "command_with_arguments",
            scenario: "AC-005: Command Arguments",
            given: func(t *testing.T) string {
                // Given: A command requires arguments
                workDir := t.TempDir()
                require.NoError(t, os.Chdir(workDir))

                commandsDir := filepath.Join(workDir, "library", "workflows", "helix", "commands")
                require.NoError(t, os.MkdirAll(commandsDir, 0755))

                buildStoryContent := `# HELIX Command: Build Story

Command accepts arguments for user story processing.`
                require.NoError(t, os.WriteFile(
                    filepath.Join(commandsDir, "build-story.md"),
                    []byte(buildStoryContent), 0644))

                return workDir
            },
            when: func(t *testing.T, workDir string) (string, error) {
                // When: I execute it with arguments
                rootCmd := getTestRootCommand()
                buf := new(bytes.Buffer)
                rootCmd.SetOut(buf)
                rootCmd.SetErr(buf)
                rootCmd.SetArgs([]string{"workflow", "helix", "execute", "build-story", "US-001", "--verbose"})

                err := rootCmd.Execute()
                return buf.String(), err
            },
            then: func(t *testing.T, workDir string, output string, err error) {
                // Then: The arguments are passed to the command context
                assert.NoError(t, err)
                assert.Contains(t, output, "Command Arguments: [US-001 --verbose]")
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            workDir := tt.given(t)
            output, err := tt.when(t, workDir)
            tt.then(t, workDir, output, err)
        })
    }
}
```

### 2. Integration Test Specifications

#### Test Suite: `TestWorkflowCommandIntegration`

```go
func TestWorkflowCommandIntegration(t *testing.T) {
    t.Run("workflow_command_discovery", func(t *testing.T) {
        // Test dynamic command discovery from library structure
        workDir := setupWorkflowLibrary(t)

        result := executeWorkflowCommand(workDir, []string{"workflow", "helix", "commands"})

        assert.Equal(t, 0, result.ExitCode)
        assert.Contains(t, result.Output, "build-story")
        assert.Contains(t, result.Output, "continue")
        assert.Contains(t, result.Output, "status")
    })

    t.Run("workflow_command_execution", func(t *testing.T) {
        // Test command file loading and display
        workDir := setupWorkflowLibrary(t)

        result := executeWorkflowCommand(workDir, []string{"workflow", "helix", "execute", "build-story", "US-001"})

        assert.Equal(t, 0, result.ExitCode)
        assert.Contains(t, result.Output, "HELIX Command: Build Story")
    })
}
```

## Contract Test Specifications

### CLI Command Contracts

#### Test Suite: `TestWorkflowCommands_Contract`

```go
func TestWorkflowCommands_Contract(t *testing.T) {
    t.Run("contract_workflow_commands_list", func(t *testing.T) {
        // Given: HELIX workflow with commands available
        setupWorkflowWithCommands(t)

        // When: Listing workflow commands
        cmd := getFreshRootCmd()
        buf := new(bytes.Buffer)
        cmd.SetOut(buf)
        cmd.SetErr(buf)
        cmd.SetArgs([]string{"workflow", "helix", "commands"})

        err := cmd.Execute()

        // Then: Commands are listed with proper format
        if err == nil {
            output := buf.String()
            assert.Contains(t, output, "Available commands")
            assert.Contains(t, output, "workflow")
        }
    })

    t.Run("contract_workflow_command_execution", func(t *testing.T) {
        // Given: Valid workflow and command
        setupWorkflowWithCommands(t)

        // When: Executing workflow command
        cmd := getFreshRootCmd()
        buf := new(bytes.Buffer)
        cmd.SetOut(buf)
        cmd.SetErr(buf)
        cmd.SetArgs([]string{"workflow", "helix", "execute", "build-story", "US-001"})

        err := cmd.Execute()

        // Then: Command executes successfully
        if err == nil {
            output := buf.String()
            assert.Contains(t, output, "HELIX Command")
            assert.Contains(t, output, "US-001")
        }
    })
}
```

## Test Environment Setup

### Workflow Library Structure

```go
func setupWorkflowLibrary(t *testing.T) string {
    workDir := t.TempDir()
    require.NoError(t, os.Chdir(workDir))

    // Create helix workflow commands
    commandsDir := filepath.Join(workDir, "library", "workflows", "helix", "commands")
    require.NoError(t, os.MkdirAll(commandsDir, 0755))

    // Create build-story command
    buildStoryContent := `# HELIX Command: Build Story

You are a HELIX workflow executor tasked with implementing work on a specific user story.

## Command Input
You will receive a user story ID as an argument.

## Your Mission
Execute comprehensive evaluation and implementation process.`

    require.NoError(t, os.WriteFile(
        filepath.Join(commandsDir, "build-story.md"),
        []byte(buildStoryContent), 0644))

    // Create continue command
    continueContent := `# HELIX Command: Continue

Continue work on the current user story.`

    require.NoError(t, os.WriteFile(
        filepath.Join(commandsDir, "continue.md"),
        []byte(continueContent), 0644))

    return workDir
}
```

## Success Criteria

### Test Execution Requirements

1. **All Tests Initially Fail (Red Phase)**
   - Tests compile successfully
   - Tests execute and fail with clear error messages
   - Failure messages indicate what needs implementation

2. **Incremental Implementation (Green Phase)**
   - Implement features to make tests pass one by one
   - No feature implementation without corresponding test
   - Maintain test isolation and reliability

3. **Test Coverage Targets**
   - 100% of user story acceptance criteria covered
   - All error scenarios validated
   - Command discovery and execution verified

### Quality Gates

- Command discovery works for any workflow with commands directory
- Command execution properly loads and displays command content
- Error handling provides clear, actionable messages
- Command arguments are properly passed to command context
- Multiple workflows are properly isolated

---

This test specification serves as the definitive guide for implementing and validating the workflow command execution system. All tests must be written and failing before any implementation begins, following strict TDD principles.