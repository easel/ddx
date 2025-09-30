# SD-015: Workflow Agent Integration - Solution Design

**Feature**: FEAT-015 Workflow Agent Integration
**User Story**: US-045
**Status**: In Design
**Owner**: Development Team
**Created**: 2025-01-20

---

## Executive Summary

Design for a generic `ddx agent` command system that enables Claude to automatically detect and delegate implementation requests to active workflows using trigger-based detection and dynamic command discovery.

**Key Design Principles**:
- Zero hardcoding of workflow names
- Dynamic command discovery from workflow.yml
- First-match routing with order-based priority
- Safe word escape hatch (default: "NODDX")
- Extensible for future workflows

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        Claude Code                          │
│                                                             │
│  1. User Request → Check for triggers                       │
│  2. Call: ddx agent request "<request>"                     │
│  3. Parse output → Execute COMMAND if workflow matches      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    ddx agent command                        │
│                                                             │
│  1. Load .ddx/config.yaml → Get active workflows           │
│  2. Check for safe word prefix                             │
│  3. For each active workflow (in order):                    │
│     - Load workflow.yml                                     │
│     - Check if agent_commands.<subcommand> enabled          │
│     - Check if triggers match                               │
│     - FIRST MATCH: Return workflow instructions             │
│  4. No match: Return NO_HANDLER                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Workflow Coordinator (e.g., HELIX)             │
│                                                             │
│  Executes action specified in agent_commands                │
│  Example: frame-request, show-status, suggest-next          │
└─────────────────────────────────────────────────────────────┘
```

---

## Component Design

### 1. Configuration Schema

#### 1.1 Config Structure

**File**: `.ddx/config.yaml`

```yaml
version: "1.0"

workflows:
  active:
    - helix
    # Future: - kanban
  safe_word: "NODDX"  # Optional, defaults to "NODDX"

library:
  path: ".ddx/library"
```

#### 1.2 Go Types

**File**: `cli/internal/config/types.go`

```go
// WorkflowsConfig defines workflow activation and settings
type WorkflowsConfig struct {
    // Active workflows in priority order (first match wins)
    Active []string `yaml:"active"`

    // SafeWord prefix to bypass workflow engagement
    // Default: "NODDX"
    SafeWord string `yaml:"safe_word,omitempty"`
}

// ApplyDefaults sets default values for workflow configuration
func (w *WorkflowsConfig) ApplyDefaults() {
    if w.SafeWord == "" {
        w.SafeWord = "NODDX"
    }
    if w.Active == nil {
        w.Active = []string{}
    }
}

// Validate ensures workflow configuration is valid
func (w *WorkflowsConfig) Validate() error {
    // Validate safe word is not empty and has no spaces
    if strings.TrimSpace(w.SafeWord) == "" {
        return fmt.Errorf("safe_word cannot be empty")
    }
    if strings.Contains(w.SafeWord, " ") {
        return fmt.Errorf("safe_word cannot contain spaces")
    }

    // Validate no duplicate workflows
    seen := make(map[string]bool)
    for _, name := range w.Active {
        if seen[name] {
            return fmt.Errorf("duplicate workflow: %s", name)
        }
        seen[name] = true
    }

    return nil
}
```

#### 1.3 JSON Schema

**File**: `cli/internal/config/schema/config.schema.json`

```json
{
  "properties": {
    "workflows": {
      "type": "object",
      "properties": {
        "active": {
          "type": "array",
          "items": {
            "type": "string",
            "minLength": 1
          },
          "description": "List of active workflow names in priority order",
          "uniqueItems": true
        },
        "safe_word": {
          "type": "string",
          "minLength": 1,
          "pattern": "^[A-Z0-9]+$",
          "default": "NODDX",
          "description": "Safe word to bypass workflow engagement"
        }
      }
    }
  }
}
```

---

### 2. Workflow Package

#### 2.1 Workflow Definition Types

**File**: `cli/internal/workflow/types.go`

```go
package workflow

// Definition represents a workflow definition from workflow.yml
type Definition struct {
    Name          string                  `yaml:"name"`
    Version       string                  `yaml:"version"`
    Description   string                  `yaml:"description"`
    Author        string                  `yaml:"author,omitempty"`
    Created       string                  `yaml:"created,omitempty"`
    Tags          []string                `yaml:"tags,omitempty"`

    // Coordinator markdown file (e.g., "coordinator.md")
    Coordinator   string                  `yaml:"coordinator,omitempty"`

    // Agent command definitions for Claude integration
    AgentCommands map[string]AgentCommand `yaml:"agent_commands,omitempty"`

    // Workflow phases (existing)
    Phases        []Phase                 `yaml:"phases"`

    // Variables (existing)
    Variables     []Variable              `yaml:"variables,omitempty"`
}

// AgentCommand defines a command that Claude can invoke
type AgentCommand struct {
    // Enabled controls whether this command is active
    Enabled     bool     `yaml:"enabled"`

    // Triggers define when this command should be invoked
    Triggers    *Triggers `yaml:"triggers,omitempty"`

    // Action is the workflow action to execute (e.g., "frame-request")
    Action      string   `yaml:"action"`

    // Description explains what this command does
    Description string   `yaml:"description"`
}

// Triggers define patterns that activate a command
type Triggers struct {
    // Keywords are single words that trigger the command
    Keywords []string `yaml:"keywords,omitempty"`

    // Patterns are regex-like patterns (simpler: just contains checks)
    Patterns []string `yaml:"patterns,omitempty"`
}

// Phase represents a workflow phase (existing type)
type Phase struct {
    ID               string   `yaml:"id"`
    Order            int      `yaml:"order"`
    Name             string   `yaml:"name"`
    Description      string   `yaml:"description"`
    RequiredRole     string   `yaml:"required_role,omitempty"`
    ExitCriteria     []string `yaml:"exit_criteria,omitempty"`
    EstimatedDuration string  `yaml:"estimated_duration,omitempty"`
}

// Variable represents a workflow variable (existing type)
type Variable struct {
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    Prompt      string `yaml:"prompt,omitempty"`
    Required    bool   `yaml:"required,omitempty"`
}

// Validate ensures the workflow definition is valid
func (d *Definition) Validate() error {
    if d.Name == "" {
        return fmt.Errorf("workflow name is required")
    }
    if d.Version == "" {
        return fmt.Errorf("workflow version is required")
    }

    // Validate agent commands
    for cmdName, cmd := range d.AgentCommands {
        if cmd.Enabled && cmd.Action == "" {
            return fmt.Errorf("agent command %s: action is required when enabled", cmdName)
        }
    }

    return nil
}

// SupportsAgentCommand checks if workflow supports a given agent subcommand
func (d *Definition) SupportsAgentCommand(subcommand string) bool {
    cmd, exists := d.AgentCommands[subcommand]
    return exists && cmd.Enabled
}

// GetAgentCommand returns the agent command definition if it exists and is enabled
func (d *Definition) GetAgentCommand(subcommand string) (*AgentCommand, bool) {
    cmd, exists := d.AgentCommands[subcommand]
    if !exists || !cmd.Enabled {
        return nil, false
    }
    return &cmd, true
}
```

#### 2.2 Workflow Loader

**File**: `cli/internal/workflow/loader.go`

```go
package workflow

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "gopkg.in/yaml.v3"
)

// Loader loads workflow definitions from the library
type Loader struct {
    libraryPath string
}

// NewLoader creates a new workflow loader
func NewLoader(libraryPath string) *Loader {
    return &Loader{
        libraryPath: libraryPath,
    }
}

// Load reads and parses a workflow.yml file
func (l *Loader) Load(workflowName string) (*Definition, error) {
    // Construct path: {libraryPath}/workflows/{workflowName}/workflow.yml
    workflowPath := filepath.Join(l.libraryPath, "workflows", workflowName, "workflow.yml")

    // Check if file exists
    if _, err := os.Stat(workflowPath); err != nil {
        if os.IsNotExist(err) {
            return nil, fmt.Errorf("workflow '%s' not found at %s", workflowName, workflowPath)
        }
        return nil, fmt.Errorf("failed to access workflow file: %w", err)
    }

    // Read file
    data, err := os.ReadFile(workflowPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read workflow file: %w", err)
    }

    // Parse YAML
    var def Definition
    if err := yaml.Unmarshal(data, &def); err != nil {
        return nil, fmt.Errorf("failed to parse workflow.yml: %w", err)
    }

    // Validate definition
    if err := def.Validate(); err != nil {
        return nil, fmt.Errorf("invalid workflow definition: %w", err)
    }

    return &def, nil
}

// MatchesTriggers checks if text matches the triggers for a given agent command
func (l *Loader) MatchesTriggers(def *Definition, subcommand string, text string) bool {
    // Get agent command
    cmd, exists := def.GetAgentCommand(subcommand)
    if !exists {
        return false
    }

    // No triggers = never matches
    if cmd.Triggers == nil {
        return false
    }

    // Normalize text for matching (lowercase, trim)
    normalizedText := strings.ToLower(strings.TrimSpace(text))

    // Check keyword matches
    for _, keyword := range cmd.Triggers.Keywords {
        normalizedKeyword := strings.ToLower(keyword)

        // Match as whole word (with word boundaries)
        if matchesKeyword(normalizedText, normalizedKeyword) {
            return true
        }
    }

    // Check pattern matches (simple substring match)
    for _, pattern := range cmd.Triggers.Patterns {
        normalizedPattern := strings.ToLower(pattern)
        if strings.Contains(normalizedText, normalizedPattern) {
            return true
        }
    }

    return false
}

// matchesKeyword checks if keyword appears as a whole word in text
func matchesKeyword(text, keyword string) bool {
    // Simple word boundary check
    // Look for keyword with spaces/punctuation around it

    patterns := []string{
        keyword + " ",  // "add something"
        " " + keyword,  // "should add"
        " " + keyword + " ", // "please add this"
    }

    // Also check if text starts with keyword or is exactly keyword
    if strings.HasPrefix(text, keyword+" ") || text == keyword {
        return true
    }

    // Check patterns
    for _, pattern := range patterns {
        if strings.Contains(text, pattern) {
            return true
        }
    }

    return false
}
```

---

### 3. Agent Command

#### 3.1 Command Structure

**File**: `cli/cmd/agent.go`

```go
package cmd

import (
    "fmt"
    "strings"

    "github.com/easel/ddx/internal/config"
    "github.com/easel/ddx/internal/workflow"
    "github.com/spf13/cobra"
)

// newAgentCommand creates the agent command for AI assistant coordination
func newAgentCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "agent <subcommand> [args...]",
        Short: "Agent coordination commands for AI assistants",
        Long: `Agent coordination commands help AI assistants like Claude understand
project context and determine appropriate actions based on active workflows.

Subcommands are dynamically discovered from active workflows.`,
        Example: `  # Check if request should use workflow
  ddx agent request "add pagination to list command"

  # Show workflow status
  ddx agent status

  # Get next action suggestion
  ddx agent next`,
    }

    // Add subcommands
    cmd.AddCommand(newAgentRequestCommand())

    return cmd
}

// newAgentRequestCommand handles user implementation requests
func newAgentRequestCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "request <text...>",
        Short: "Process user request and determine appropriate workflow action",
        Long: `Analyzes user request to determine if it should be handled by an active workflow.

Checks for:
- Safe word prefix (bypasses workflow)
- Trigger keywords/patterns
- Active workflows with request handling

Returns structured output for Claude to execute.`,
        Args:  cobra.MinimumNArgs(1),
        RunE:  runAgentRequest,
    }
}

// runAgentRequest implements the request subcommand logic
func runAgentRequest(cmd *cobra.Command, args []string) error {
    userRequest := strings.Join(args, " ")

    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        // No config or error loading = no workflows active
        fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
        return nil
    }

    // Check if any workflows are active
    if len(cfg.Workflows.Active) == 0 {
        fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
        return nil
    }

    // Check for safe word (default: "NODDX")
    safeWord := cfg.Workflows.SafeWord
    if strings.HasPrefix(userRequest, safeWord+" ") || strings.HasPrefix(userRequest, safeWord+":") {
        // Remove safe word prefix
        message := strings.TrimPrefix(userRequest, safeWord+" ")
        message = strings.TrimPrefix(message, safeWord+":")
        message = strings.TrimSpace(message)

        // Return safe word bypass
        fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
        fmt.Fprintf(cmd.OutOrStdout(), "SAFE_WORD: %s\n", safeWord)
        fmt.Fprintf(cmd.OutOrStdout(), "MESSAGE: %s\n", message)
        return nil
    }

    // Create workflow loader
    loader := workflow.NewLoader(cfg.Library.Path)

    // Check each active workflow IN ORDER (first match wins)
    for _, workflowName := range cfg.Workflows.Active {
        wf, err := loader.Load(workflowName)
        if err != nil {
            // Skip workflows that can't be loaded
            continue
        }

        // Check if workflow has "request" agent command enabled
        agentCmd, exists := wf.GetAgentCommand("request")
        if !exists {
            continue
        }

        // Check if request matches triggers
        if !loader.MatchesTriggers(wf, "request", userRequest) {
            continue
        }

        // FIRST MATCH WINS - return workflow instructions
        fmt.Fprintf(cmd.OutOrStdout(), "WORKFLOW: %s\n", wf.Name)
        fmt.Fprintf(cmd.OutOrStdout(), "SUBCOMMAND: request\n")
        fmt.Fprintf(cmd.OutOrStdout(), "ACTION: %s\n", agentCmd.Action)
        fmt.Fprintf(cmd.OutOrStdout(), "COMMAND: ddx workflow %s execute %s %q\n",
            wf.Name, agentCmd.Action, userRequest)
        fmt.Fprintf(cmd.OutOrStdout(), "REASON: %s\n", agentCmd.Description)
        return nil
    }

    // No workflow matched
    fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
    return nil
}
```

#### 3.2 Output Format Specification

**Standard Output Formats**:

1. **Workflow Match**:
```
WORKFLOW: helix
SUBCOMMAND: request
ACTION: frame-request
COMMAND: ddx workflow helix execute frame-request "add pagination"
REASON: Frame user request in workflow terms and route to appropriate phase
```

2. **Safe Word Bypass**:
```
NO_HANDLER
SAFE_WORD: NODDX
MESSAGE: add pagination
```

3. **No Handler**:
```
NO_HANDLER
```

---

### 4. Workflow Commands (Extensions)

#### 4.1 Activate Command

**File**: `cli/cmd/workflow.go` (extend existing)

```go
// activateWorkflow adds a workflow to the active list
func activateWorkflow(cmd *cobra.Command, name string, force bool) error {
    // Load config
    cfg, err := config.Load()
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    // Verify workflow exists before activating
    loader := workflow.NewLoader(cfg.Library.Path)
    wf, err := loader.Load(name)
    if err != nil {
        return fmt.Errorf("workflow '%s' not found: %w", name, err)
    }

    // Check if already active
    for _, active := range cfg.Workflows.Active {
        if active == name {
            fmt.Fprintf(cmd.OutOrStdout(), "Workflow '%s' is already active\n", name)
            return nil
        }
    }

    // Add to active workflows (maintains order)
    cfg.Workflows.Active = append(cfg.Workflows.Active, name)

    // Save config
    if err := saveConfig(cfg); err != nil {
        return fmt.Errorf("failed to save config: %w", err)
    }

    fmt.Fprintf(cmd.OutOrStdout(), "✓ Activated %s workflow\n", wf.Name)
    fmt.Fprintf(cmd.OutOrStdout(), "  Priority: %d of %d\n",
        len(cfg.Workflows.Active), len(cfg.Workflows.Active))

    // Show agent commands if available
    if len(wf.AgentCommands) > 0 {
        fmt.Fprintln(cmd.OutOrStdout(), "  Agent commands:")
        for cmdName, cmdDef := range wf.AgentCommands {
            if cmdDef.Enabled {
                fmt.Fprintf(cmd.OutOrStdout(), "    • %s - %s\n", cmdName, cmdDef.Description)
            }
        }
    }

    return nil
}

// deactivateWorkflow removes a workflow from the active list
func deactivateWorkflow(cmd *cobra.Command, name string) error {
    cfg, err := config.Load()
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    // Find and remove from active list
    newActive := []string{}
    found := false
    for _, w := range cfg.Workflows.Active {
        if w != name {
            newActive = append(newActive, w)
        } else {
            found = true
        }
    }

    if !found {
        fmt.Fprintf(cmd.OutOrStdout(), "Workflow '%s' is not active\n", name)
        return nil
    }

    cfg.Workflows.Active = newActive

    // Save config
    if err := saveConfig(cfg); err != nil {
        return fmt.Errorf("failed to save config: %w", err)
    }

    fmt.Fprintf(cmd.OutOrStdout(), "✓ Deactivated %s workflow\n", name)
    return nil
}

// showWorkflowStatus displays active workflows and their capabilities
func showWorkflowStatus(cmd *cobra.Command) error {
    cfg, err := config.Load()
    if err != nil || len(cfg.Workflows.Active) == 0 {
        fmt.Fprintln(cmd.OutOrStdout(), "No active workflows")
        return nil
    }

    fmt.Fprintln(cmd.OutOrStdout(), "Active workflows (in priority order):")
    fmt.Fprintln(cmd.OutOrStdout(), "")

    loader := workflow.NewLoader(cfg.Library.Path)
    for i, name := range cfg.Workflows.Active {
        wf, err := loader.Load(name)
        if err != nil {
            fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s (error loading: %v)\n", i+1, name, err)
            continue
        }

        fmt.Fprintf(cmd.OutOrStdout(), "  %d. %s - %s\n", i+1, wf.Name, wf.Description)

        // Show agent commands
        if len(wf.AgentCommands) > 0 {
            fmt.Fprintln(cmd.OutOrStdout(), "     Agent commands:")
            for cmdName, cmdDef := range wf.AgentCommands {
                if cmdDef.Enabled {
                    fmt.Fprintf(cmd.OutOrStdout(), "       • %s - %s\n", cmdName, cmdDef.Description)

                    // Show triggers for request command
                    if cmdName == "request" && cmdDef.Triggers != nil {
                        if len(cmdDef.Triggers.Keywords) > 0 {
                            fmt.Fprintf(cmd.OutOrStdout(), "         Keywords: %s\n",
                                strings.Join(cmdDef.Triggers.Keywords, ", "))
                        }
                        if len(cmdDef.Triggers.Patterns) > 0 {
                            fmt.Fprintf(cmd.OutOrStdout(), "         Patterns: %s\n",
                                strings.Join(cmdDef.Triggers.Patterns, ", "))
                        }
                    }
                }
            }
        }
        fmt.Fprintln(cmd.OutOrStdout(), "")
    }

    fmt.Fprintf(cmd.OutOrStdout(), "Safe word: %s (prefix to bypass workflows)\n", cfg.Workflows.SafeWord)

    return nil
}

// saveConfig writes configuration to disk
func saveConfig(cfg *config.Config) error {
    // Implementation depends on existing config save logic
    // Should save to .ddx/config.yaml
    loader, err := config.NewConfigLoader()
    if err != nil {
        return err
    }

    return loader.SaveConfig(cfg, ".ddx/config.yaml")
}
```

---

### 5. Workflow Definition Extension

#### 5.1 HELIX workflow.yml Extension

**File**: `library/workflows/helix/workflow.yml` (extend existing)

```yaml
name: helix
version: 1.0.0
description: AI-assisted iterative development workflow with human-AI collaboration
author: DDX Team
created: 2025-01-13
tags:
  - ai-assisted-development
  - iterative
  - human-ai-collaboration

# NEW: Agent command integration
coordinator: coordinator.md

agent_commands:
  request:
    enabled: true
    triggers:
      keywords:
        - add
        - implement
        - create
        - build
        - fix
        - refactor
        - update
        - remove
        - delete
      patterns:
        - "US-"
        - "work on"
        - "make this change"
    action: frame-request
    description: "Frame user request in workflow terms and route to appropriate phase"

  status:
    enabled: true
    action: show-status
    description: "Show current HELIX phase and progress"

  next:
    enabled: true
    action: suggest-next
    description: "Suggest next action based on workflow state"

# ... existing phases, variables, etc.
```

---

## Algorithms

### Algorithm 1: Safe Word Detection

```
function detectSafeWord(request, safeWord):
  // Normalize
  trimmedRequest = trim(request)

  // Check prefixes
  if startsWith(trimmedRequest, safeWord + " "):
    message = removePrefix(trimmedRequest, safeWord + " ")
    return (true, message)

  if startsWith(trimmedRequest, safeWord + ":"):
    message = removePrefix(trimmedRequest, safeWord + ":")
    return (true, message)

  return (false, request)
```

### Algorithm 2: Trigger Matching

```
function matchesTriggers(definition, subcommand, text):
  // Get agent command
  agentCmd = definition.getAgentCommand(subcommand)
  if agentCmd == null or not agentCmd.enabled:
    return false

  if agentCmd.triggers == null:
    return false

  // Normalize text
  normalized = lowercase(trim(text))

  // Check keywords (whole word match)
  for each keyword in agentCmd.triggers.keywords:
    normalizedKeyword = lowercase(keyword)
    if matchesKeyword(normalized, normalizedKeyword):
      return true

  // Check patterns (substring match)
  for each pattern in agentCmd.triggers.patterns:
    normalizedPattern = lowercase(pattern)
    if contains(normalized, normalizedPattern):
      return true

  return false

function matchesKeyword(text, keyword):
  // Match as whole word with boundaries
  if text == keyword:
    return true

  if startsWith(text, keyword + " "):
    return true

  if endsWith(text, " " + keyword):
    return true

  if contains(text, " " + keyword + " "):
    return true

  return false
```

### Algorithm 3: Workflow Routing (First-Match)

```
function routeToWorkflow(request, config, loader):
  // Check safe word first
  (isSafe, message) = detectSafeWord(request, config.workflows.safeWord)
  if isSafe:
    return NoHandler(safeWord=config.workflows.safeWord, message=message)

  // Check active workflows in order
  for each workflowName in config.workflows.active:
    definition = loader.load(workflowName)
    if definition == null:
      continue  // Skip invalid workflows

    // Check if workflow supports "request" subcommand
    if matchesTriggers(definition, "request", request):
      // FIRST MATCH WINS
      agentCmd = definition.getAgentCommand("request")
      return WorkflowMatch(
        workflow=definition.name,
        subcommand="request",
        action=agentCmd.action,
        command="ddx workflow " + definition.name + " execute " + agentCmd.action + " " + quote(request)
      )

  // No workflow matched
  return NoHandler()
```

---

## Data Flow Diagrams

### Request Flow

```
User Request
     │
     ▼
Claude detects trigger keywords
     │
     ▼
ddx agent request "<request>"
     │
     ├─> Load .ddx/config.yaml
     │   └─> Get workflows.active, safe_word
     │
     ├─> Check safe word prefix
     │   ├─> YES: Return NO_HANDLER + SAFE_WORD
     │   └─> NO: Continue
     │
     └─> For each workflow in active (order):
         ├─> Load workflow.yml
         ├─> Check agent_commands.request.enabled
         ├─> Match triggers
         ├─> MATCH: Return WORKFLOW instructions
         └─> NO MATCH: Try next workflow

     If no match: Return NO_HANDLER
```

---

## Error Handling

### Configuration Errors

1. **No config file**: Return `NO_HANDLER` (graceful degradation)
2. **Invalid config**: Return error message to stderr, `NO_HANDLER` to stdout
3. **Invalid safe word**: Use default "NODDX"

### Workflow Loading Errors

1. **Workflow not found**: Skip and try next workflow
2. **Invalid workflow.yml**: Skip and try next workflow
3. **Parse error**: Log warning, skip workflow

### Agent Command Errors

1. **Missing action**: Skip agent command
2. **No triggers defined**: Never matches
3. **Malformed triggers**: Log warning, skip

---

## Security Considerations

1. **Path Traversal**: Validate workflow names (no `..`, `/`, `\`)
2. **Command Injection**: Quote all user input in COMMAND output
3. **YAML Parsing**: Use safe YAML parser (gopkg.in/yaml.v3)
4. **File Access**: Restrict to library path only

---

## Performance Considerations

1. **Lazy Loading**: Only load workflows when checking
2. **Caching**: Consider caching workflow definitions (future)
3. **Short-Circuit**: First-match wins, skip remaining workflows
4. **Fast Path**: Safe word check before loading workflows

**Expected Performance**:
- Safe word check: <1ms
- Load single workflow.yml: <10ms
- Trigger matching: <1ms
- Total: <50ms for typical case (1-2 active workflows)

---

## Testing Strategy

### Unit Tests

1. **Config Validation**:
   - Valid config with workflows
   - Empty workflows.active
   - Custom safe word
   - Invalid safe word (with spaces)
   - Duplicate workflows

2. **Workflow Loader**:
   - Load valid workflow.yml
   - Load missing workflow
   - Parse errors
   - Validate agent_commands

3. **Trigger Matching**:
   - Keyword matching (whole word)
   - Pattern matching (substring)
   - Case insensitivity
   - No triggers defined
   - Empty text

4. **Safe Word Detection**:
   - Safe word prefix with space
   - Safe word prefix with colon
   - Safe word in middle of text (should not match)
   - Default safe word

5. **Routing Logic**:
   - First workflow matches
   - Second workflow matches
   - No workflow matches
   - Safe word bypass

### Integration Tests

1. **End-to-End Flow**:
   - Activate workflow → agent request → workflow match
   - Multiple workflows with priority
   - Safe word usage
   - No workflows active

2. **CLI Commands**:
   - `ddx workflow helix activate`
   - `ddx workflow helix deactivate`
   - `ddx workflow status`
   - `ddx agent request`

---

## Open Issues

1. **Trigger Ambiguity**: How to handle "we should add..." vs "add..."?
   - Resolution: Context heuristics or rely on safe word

2. **Multiple Subcommands**: Should agent dynamically discover all subcommands?
   - Resolution: Start with `request` only, add others as needed

3. **Caching**: Should workflow definitions be cached?
   - Resolution: Defer until performance measurement shows need

---

## Dependencies

- `gopkg.in/yaml.v3` - YAML parsing
- Existing config package
- Existing workflow structure

---

## Migration Path

1. Add config schema (backward compatible)
2. Add workflow package (new)
3. Add agent command (new)
4. Extend workflow.yml (backward compatible)
5. Update CLAUDE.md

**No breaking changes to existing functionality.**

---

## Future Enhancements

1. **Workflow Composition**: Multiple workflows handling same request
2. **Context-Aware Triggers**: Consider previous conversation
3. **Learning**: Track trigger accuracy, adjust over time
4. **Workflow Dependencies**: One workflow can require another
5. **Priority Weights**: More sophisticated than just order

---

**Status**: Design Complete
**Next Phase**: Test (write failing tests)
**Reviewers**: TBD
**Approval Date**: TBD