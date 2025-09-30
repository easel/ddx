# CLI-005: Agent Commands API Contract

**Feature**: FEAT-015 Workflow Agent Integration
**Related**: SD-015
**Status**: Design
**Version**: 1.0

---

## Overview

This document defines the API contracts for the `ddx agent` command family, including input/output formats, behavior contracts, and error handling.

---

## Command: `ddx agent request`

### Purpose

Process a user request and determine if it should be handled by an active workflow.

### Signature

```bash
ddx agent request <text...>
```

### Arguments

- `text...` - User request text (required, minimum 1 word)

### Output Formats

#### 1. Workflow Match

**When**: Request matches workflow triggers

```
WORKFLOW: <workflow-name>
SUBCOMMAND: request
ACTION: <action-name>
COMMAND: ddx workflow <workflow-name> execute <action-name> "<request>"
REASON: <description>
```

**Fields**:
- `WORKFLOW`: Name of matching workflow (string)
- `SUBCOMMAND`: Always "request" for this command
- `ACTION`: Workflow action to execute (string)
- `COMMAND`: Full command for Claude to execute (string, properly quoted)
- `REASON`: Human-readable explanation (string)

**Example**:
```
WORKFLOW: helix
SUBCOMMAND: request
ACTION: frame-request
COMMAND: ddx workflow helix execute frame-request "add pagination to list"
REASON: Frame user request in workflow terms and route to appropriate phase
```

#### 2. Safe Word Bypass

**When**: Request starts with safe word prefix

```
NO_HANDLER
SAFE_WORD: <safe-word>
MESSAGE: <request-without-prefix>
```

**Fields**:
- `NO_HANDLER`: Indicates no workflow should handle this
- `SAFE_WORD`: The safe word that was detected (string)
- `MESSAGE`: Request text with safe word prefix removed (string)

**Example**:
```
NO_HANDLER
SAFE_WORD: NODDX
MESSAGE: add pagination to list
```

#### 3. No Handler

**When**: No workflow matches or no workflows active

```
NO_HANDLER
```

**Single line output**: Indicates request should be handled normally by Claude.

### Exit Codes

- `0` - Success (all cases)
- `1` - Invalid arguments
- `2` - Internal error

### Error Handling

**Configuration errors**: Output `NO_HANDLER` (graceful degradation)
**Workflow loading errors**: Skip workflow, try next
**Parse errors**: Log to stderr, output `NO_HANDLER`

### Contract Guarantees

1. **Deterministic**: Same input + config → same output
2. **Fast**: <100ms for typical case (1-2 workflows)
3. **Safe**: All user input properly quoted in COMMAND
4. **Order-sensitive**: First matching workflow wins

---

## Command: `ddx workflow <name> activate`

### Purpose

Add a workflow to the active list in configuration.

### Signature

```bash
ddx workflow <name> activate [--force]
```

### Arguments

- `name` - Workflow name (required)
- `--force` - Force activation even if already active (optional)

### Output Format

**Success**:
```
✓ Activated <workflow-name> workflow
  Priority: <position> of <total>
  Agent commands:
    • <command> - <description>
    ...
```

**Already Active**:
```
Workflow '<name>' is already active
```

**Not Found**:
```
Error: workflow '<name>' not found: <error>
```

### Side Effects

- Modifies `.ddx/config.yaml`
- Appends to `workflows.active` array
- Validates workflow exists before activation

### Contract Guarantees

1. **Validation**: Workflow must exist in library
2. **Idempotent**: Activating twice has no additional effect
3. **Order preserved**: New workflow added to end of list
4. **Atomic**: Config saved or error returned, no partial state

---

## Command: `ddx workflow <name> deactivate`

### Purpose

Remove a workflow from the active list.

### Signature

```bash
ddx workflow <name> deactivate
```

### Arguments

- `name` - Workflow name (required)

### Output Format

**Success**:
```
✓ Deactivated <workflow-name> workflow
```

**Not Active**:
```
Workflow '<name>' is not active
```

### Side Effects

- Modifies `.ddx/config.yaml`
- Removes from `workflows.active` array
- Preserves order of remaining workflows

### Contract Guarantees

1. **Idempotent**: Deactivating inactive workflow has no effect
2. **Order preserved**: Remaining workflows maintain relative order
3. **Atomic**: Config saved or error returned

---

## Command: `ddx workflow status`

### Purpose

Display active workflows and their capabilities.

### Signature

```bash
ddx workflow status
```

### Output Format

**With Active Workflows**:
```
Active workflows (in priority order):

  1. <name> - <description>
     Agent commands:
       • <cmd> - <description>
         Keywords: <keyword1>, <keyword2>, ...
         Patterns: <pattern1>, <pattern2>, ...

  2. <name> - <description>
     Agent commands:
       • <cmd> - <description>

Safe word: <safe-word> (prefix to bypass workflows)
```

**No Active Workflows**:
```
No active workflows
```

### Contract Guarantees

1. **Read-only**: No side effects
2. **Complete**: Shows all active workflows
3. **Ordered**: Displays in priority order (1 = highest)
4. **Detailed**: Shows agent commands and triggers for "request" command

---

## Workflow Package API

### Loader.Load()

```go
func (l *Loader) Load(workflowName string) (*Definition, error)
```

**Purpose**: Load and parse a workflow.yml file

**Arguments**:
- `workflowName` - Name of workflow (must be valid filename)

**Returns**:
- `*Definition` - Parsed workflow definition
- `error` - Error if workflow not found, parse failed, or validation failed

**Contract**:
- Validates workflow name (no path traversal)
- Checks file exists before reading
- Parses YAML strictly
- Validates definition before returning
- Returns specific error types for different failures

**Errors**:
- `ErrWorkflowNotFound` - Workflow doesn't exist
- `ErrInvalidWorkflow` - Parse or validation failed

### Loader.MatchesTriggers()

```go
func (l *Loader) MatchesTriggers(def *Definition, subcommand string, text string) bool
```

**Purpose**: Check if text matches triggers for an agent command

**Arguments**:
- `def` - Workflow definition
- `subcommand` - Agent command name (e.g., "request")
- `text` - User input text

**Returns**:
- `bool` - true if text matches any trigger

**Contract**:
- Case-insensitive matching
- Keywords match as whole words
- Patterns match as substrings
- Returns false if no triggers defined
- Returns false if subcommand not found or disabled

**Matching Rules**:
- Keywords: Whole word boundaries required
- Patterns: Simple substring match
- OR logic: Any keyword OR any pattern matches

### Definition.GetAgentCommand()

```go
func (d *Definition) GetAgentCommand(subcommand string) (*AgentCommand, bool)
```

**Purpose**: Get agent command definition if enabled

**Arguments**:
- `subcommand` - Command name

**Returns**:
- `*AgentCommand` - Command definition
- `bool` - true if command exists and is enabled

**Contract**:
- Returns nil, false if command doesn't exist
- Returns nil, false if command exists but disabled
- Returns command, true only if exists AND enabled

---

## Configuration API

### WorkflowsConfig.Validate()

```go
func (w *WorkflowsConfig) Validate() error
```

**Purpose**: Validate workflow configuration

**Returns**:
- `error` - Validation error or nil

**Validation Rules**:
1. Safe word not empty
2. Safe word has no spaces
3. No duplicate workflows in active list
4. Workflow names valid (no special chars)

### WorkflowsConfig.ApplyDefaults()

```go
func (w *WorkflowsConfig) ApplyDefaults()
```

**Purpose**: Set default values for missing fields

**Side Effects**:
- Sets `SafeWord` to "NODDX" if empty
- Initializes `Active` to empty slice if nil

**Contract**:
- Idempotent: Can be called multiple times safely
- Non-destructive: Doesn't overwrite existing values

---

## Error Handling Contracts

### Graceful Degradation

All agent commands must gracefully degrade:
1. Missing config → NO_HANDLER
2. Invalid config → NO_HANDLER + warning to stderr
3. Workflow load error → Skip workflow, try next
4. Parse error → Skip, log warning

### User-Facing Errors

User-facing errors must be clear and actionable:
- Include what went wrong
- Include how to fix it (when possible)
- Use consistent format

**Example**:
```
Error: workflow 'helix' not found
  Path checked: .ddx/library/workflows/helix/workflow.yml
  Suggestion: Run 'ddx workflow list' to see available workflows
```

---

## Performance Contracts

### Response Time Guarantees

- `ddx agent request`: <100ms (typical, 1-2 workflows)
- `ddx workflow activate`: <500ms
- `ddx workflow deactivate`: <500ms
- `ddx workflow status`: <200ms

### Resource Usage

- Memory: <10MB per workflow definition
- File I/O: One read per workflow (no caching in v1)
- CPU: Minimal (simple string matching)

---

## Backward Compatibility

### Version 1.0 Promises

1. **Output format stability**: WORKFLOW/NO_HANDLER format won't change
2. **Config schema stability**: `workflows.active` array format won't change
3. **Safe word behavior**: Prefix checking won't change
4. **Trigger matching**: Keyword/pattern logic won't change

### Breaking Changes (Future)

If breaking changes needed:
1. Increment major version
2. Provide migration guide
3. Support old format for 2 versions

---

## Testing Requirements

### Contract Tests

Each command must have contract tests:
1. Valid input → expected output format
2. Invalid input → appropriate error
3. Edge cases → documented behavior
4. Performance → within guarantees

### Integration Tests

1. Full workflow: activate → request → execute
2. Priority handling: multiple workflows
3. Safe word usage
4. Error scenarios

---

## Documentation Requirements

### User Documentation

1. Command reference with examples
2. Output format specification
3. Error messages and solutions
4. Workflow author guide

### Developer Documentation

1. API reference for workflow package
2. Contract test examples
3. Error handling patterns
4. Performance guidelines

---

**Status**: Design Complete
**Next**: Write contract tests (Test Phase)
**Reviewed By**: TBD
**Approved Date**: TBD