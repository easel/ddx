# Workflow System Technical Overview

**Version**: 1.0  
**Status**: Specification  
**Date**: 2025-09-12  
**Related**: [Usage Documentation](../../../usage/workflows/overview.md), [PRD](../../prd-ddx-v1.md)

## Executive Summary

A workflow in DDX is a sequence of phases that guide you through complex development tasks. Each phase has specific entry requirements, work to be done, and exit criteria.

## Core Concept: Workflows Are Phases

**A workflow is a collection of phases.** Each phase represents a distinct stage of work with three key components:

### 1. Input Gates
**Prerequisites to enter the phase**
- Required conditions that must be satisfied
- Dependencies from previous phases
- Resources or approvals needed
- Validation commands to verify readiness

### 2. Work Items
The actual work performed in a phase consists of two types:

#### Artifacts (Template-Based Outputs)
- One or more files created from **templates + a prompt**
- Each file has a fixed structure defined by its template
- A single prompt guides the creation of all template files
- Examples: 
  - PRD (might generate: overview.md, requirements.md, success-metrics.md)
  - Architecture (might generate: system-design.md, api-spec.yml, database-schema.sql)
  - Test Suite (might generate: test-plan.md, unit-tests/*, integration-tests/*)
- Templates ensure consistency and completeness

#### Actions (Arbitrary Operations)  
- Operations that perform arbitrary tasks without fixed templates
- Defined by **prompts only** (no templates)
- Examples: "Refactor authentication across codebase", "Update all API documentation", "Generate CRUD endpoints for all models"
- More flexible but less structured than artifacts

### 3. Exit Criteria
**Requirements to complete the phase**
- Required artifacts must be generated
- Actions must be completed
- Quality gates must be passed
- Validation checks must succeed

## System Architecture

### Phase-Centric Design

The entire workflow system revolves around phases. Everything else - templates, prompts, artifacts, automation - exists to support phase execution.

```
Workflow
    ↓
phases/
├── 01-define/
│   ├── input-gates.yml
│   ├── artifacts/
│   └── actions/
├── 02-design/
│   ├── input-gates.yml
│   ├── artifacts/
│   └── actions/
└── 03-implement/
    ├── input-gates.yml
    ├── artifacts/
    └── actions/
```

### Phase Ordering

Phases are ordered using a simple numbering scheme:
- **Sequential phases**: 01-define → 02-design → 03-implement
- **Parallel phases**: 03a-implement, 03b-security-review (same number, different suffix)
- **Sub-phases**: 03.1-backend, 03.2-frontend (decimal notation for sub-phases)

The filesystem naturally sorts phases in execution order, making workflows self-documenting.

### Data Structures

#### Workflow Metadata (`workflow.yml`)
```yaml
name: string                    # Unique workflow identifier
version: semver                 # Semantic version
description: string             # Human-readable description
author: string                  # Creator identification
tags: [string]                 # Categorization tags
category: string               # Primary workflow category

phases:                        # The core of any workflow
  - id: string                 # Phase identifier
    order: number              # Execution order (1, 2, 3...)
    name: string               # Display name
    description: string        # What this phase accomplishes
    
    # Input gates defined in phases/{order}-{id}/input-gates.yml
    # Artifacts discovered from phases/{order}-{id}/artifacts/
    # Actions discovered from phases/{order}-{id}/actions/
    
    exit_criteria:            # Requirements for phase completion
      - requirement: string   # What must be done
        validation: string    # How to verify completion
    
    timeout: duration         # Maximum execution time

# With nested structure, artifacts and actions are defined by their location
# in the filesystem rather than in workflow.yml

# Input gates file structure (phases/{order}-{id}/input-gates.yml):
input_gates:
  - criteria: string           # What must be true/complete
    validation: string         # Optional command to verify
    source: string            # Where this comes from (prior phase, external)
    required: boolean         # Must be satisfied to proceed

# Artifact metadata (phases/{order}-{id}/artifacts/{name}/meta.yml):
artifact:
  name: string                # Display name
  type: enum                  # document|code|config|data
  templates:                  # List of template files
    - path: templates/overview.md
      output: docs/prd/overview.md
    - path: templates/requirements.md  
      output: docs/prd/requirements.md
  required: boolean           # Must be completed
  validation: object          # Validation rules

# Action metadata (phases/{order}-{id}/actions/{name}/meta.yml):
action:
  name: string                # Display name
  description: string         # What this action accomplishes
  affects:                    # What this action modifies
    - files: [glob]          # File patterns affected
    - artifacts: [string]    # Artifact IDs affected
    
variables:                     # Workflow-specific variables
  - name: string               # Variable identifier
    prompt: string             # User prompt for value
    default: any               # Default value
    type: enum                 # string|number|boolean|array|object
    required: boolean          # Must be provided
    validation: regex          # Validation pattern

automation:                    # CLI integration points
  init_command: string         # Initialization command
  validate_command: string     # Validation command
  complete_command: string     # Completion command

integrations:                  # External tool integrations
  - type: string               # Integration type (github, jira, etc.)
    trigger: string            # When to execute
    action: string             # What action to perform
    config: object             # Integration-specific config
```

#### Artifact Metadata (`{artifact}/meta.yml`)
```yaml
id: string                     # Unique artifact identifier
name: string                   # Human-readable name
description: string            # Purpose and usage
type: enum                     # document|code|config|data
version: semver                # Artifact version
status: enum                   # draft|review|approved|deprecated

dependencies:                  # Artifact dependencies
  - artifact: string           # Dependent artifact ID
    relationship: enum         # requires|extends|references
    version: semver           # Required version

template:                      # Template configuration
  file: path                   # Template file path
  variables: [string]          # Required variables
  conditionals: [string]       # Conditional sections
  
prompt:                        # Prompt configuration
  file: path                   # Prompt file path
  model_hints: object          # Model-specific guidance
  context_window: number       # Recommended context size

validation:                    # Validation rules
  required_sections: [string]  # Must include these sections
  min_length: number           # Minimum content length
  format: enum                 # markdown|yaml|json|text
  custom_rules: [object]       # Custom validation logic

examples:                      # Reference examples
  - file: path                 # Example file path
    context: string            # Example context/scenario
    quality: enum              # basic|intermediate|advanced
```

### Storage Architecture

#### File System Layout
```
workflows/
├── {workflow-name}/
│   ├── README.md              # Workflow documentation
│   ├── workflow.yml           # Phase definitions and configuration
│   └── phases/               # Phases in execution order
│       ├── 01-define/        # Phase 1: Definition
│       │   ├── README.md     # Phase documentation
│       │   ├── input-gates.yml  # Entry criteria
│       │   ├── artifacts/    # Artifacts created in this phase
│       │   │   └── prd/
│       │   │       ├── templates/
│       │   │       │   ├── overview.md
│       │   │       │   ├── requirements.md
│       │   │       │   └── metrics.md
│       │   │       ├── prompt.md
│       │   │       └── examples/
│       │   └── actions/      # Actions performed in this phase
│       │       └── gather-requirements/
│       │           └── prompt.md
│       ├── 02-design/        # Phase 2: Design
│       │   ├── README.md
│       │   ├── input-gates.yml
│       │   ├── artifacts/
│       │   │   └── architecture/
│       │   │       ├── templates/
│       │   │       │   ├── system-design.md
│       │   │       │   └── api-spec.yml
│       │   │       ├── prompt.md
│       │   │       └── examples/
│       │   └── actions/
│       └── 03-implement/     # Phase 3: Implementation
│           ├── README.md
│           ├── input-gates.yml
│           ├── artifacts/
│           └── actions/
```

#### Key Principles
1. **Phase ordering through name-mangling**: Prefix directories with order numbers (01, 02, 03...)
2. **Self-contained phases**: Each phase directory contains all its artifacts and actions
3. **Natural hierarchy**: The filesystem structure mirrors workflow execution
4. **Parallel phases**: Use same prefix with letter suffix (03a-implement, 03b-security-review)

#### State Management
- **Workflow State**: Stored in `.ddx/workflows/{name}/state.yml`
- **Phase Progress**: Individual phase completion status
- **Artifact Status**: Generation and validation states
- **Variable Context**: Resolved variable values
- **Session Recovery**: Support for interrupted workflows

### Integration Points

#### CLI Command Interface
```bash
# Primary workflow commands
ddx workflow list               # List available workflows
ddx workflow create <name>      # Create new workflow
ddx workflow apply <name>       # Apply existing workflow
ddx workflow validate <name>    # Validate workflow structure
ddx workflow status <name>      # Check workflow progress

# Phase management
ddx workflow phase list <name>     # List workflow phases
ddx workflow phase start <name> <phase> # Start specific phase
ddx workflow phase complete <name> <phase> # Mark phase complete
ddx workflow phase reset <name> <phase>   # Reset phase state

# Artifact operations
ddx workflow artifact list <name>         # List workflow artifacts
ddx workflow artifact generate <name> <artifact> # Generate artifact
ddx workflow artifact validate <name> <artifact> # Validate artifact
```

#### API Interfaces

##### WorkflowManager Interface
```go
type WorkflowManager interface {
    List(criteria ListCriteria) ([]WorkflowInfo, error)
    Load(name string) (*Workflow, error)
    Create(spec WorkflowSpec) error
    Apply(name string, options ApplyOptions) (*WorkflowExecution, error)
    Validate(name string) (*ValidationResult, error)
    Delete(name string) error
}
```

##### PhaseManager Interface
```go
type PhaseManager interface {
    GetPhases(workflowName string) ([]Phase, error)
    StartPhase(workflowName, phaseID string) error
    CompletePhase(workflowName, phaseID string, artifacts []Artifact) error
    GetPhaseStatus(workflowName, phaseID string) (*PhaseStatus, error)
    ValidateTransition(from, to string) error
}
```

##### ArtifactManager Interface
```go
type ArtifactManager interface {
    GenerateArtifact(workflowName, artifactID string, context Context) (*Artifact, error)
    ValidateArtifact(artifact *Artifact) (*ValidationResult, error)
    SaveArtifact(artifact *Artifact, path string) error
    LoadArtifact(path string) (*Artifact, error)
    GetDependencies(artifactID string) ([]Dependency, error)
}
```

#### Template Engine Integration
- **Variable Resolution**: Context-aware variable substitution
- **Conditional Logic**: Template sections based on metadata
- **Include Mechanism**: Template composition and reuse
- **Validation**: Schema-based template validation

#### AI Integration
- **Model Abstraction**: Support for multiple AI providers
- **Context Management**: Efficient context window utilization
- **Prompt Templates**: Reusable prompt components
- **Response Processing**: Structured output parsing

### Security Considerations

#### Access Control
- Workflow execution permissions
- Template modification restrictions  
- Artifact access controls
- Integration authorization

#### Data Protection
- Variable value encryption for sensitive data
- Artifact access logging
- Secure external integrations
- PII handling in templates

#### Validation Security
- Template injection prevention
- Command injection protection
- Safe variable substitution
- Malicious workflow detection

### Performance Requirements

#### Scalability Targets
- **Workflow Catalog**: Support 1000+ workflows
- **Concurrent Executions**: 50+ simultaneous workflow runs
- **Artifact Generation**: Sub-second template processing
- **Phase Transitions**: <100ms validation and state updates

#### Resource Management
- **Memory**: Efficient template caching
- **Storage**: Incremental artifact updates
- **Network**: Optimized external integrations
- **CPU**: Parallel phase execution

### Error Handling

#### Error Categories
1. **Configuration Errors**: Invalid workflow.yml, missing templates
2. **Execution Errors**: Phase failures, artifact generation issues
3. **Validation Errors**: Schema violations, requirement failures
4. **Integration Errors**: External system failures
5. **Recovery Errors**: State corruption, incomplete workflows

#### Recovery Strategies
- **Checkpoint Recovery**: Resume from last completed phase
- **Partial Rollback**: Undo specific artifacts or phases
- **State Repair**: Automatic state consistency checking
- **Manual Recovery**: Administrative override capabilities

### Monitoring and Observability

#### Metrics
- Workflow execution success rates
- Phase completion times
- Artifact generation metrics
- Template usage statistics
- Error rate tracking

#### Logging
- Structured workflow execution logs
- Phase transition events
- Artifact generation tracking
- Error and exception logging
- Performance timing data

#### Health Checks
- Workflow catalog integrity
- Template availability
- External integration status
- State persistence health

## Examples: Artifacts vs Actions

### Artifact Example: Product Requirements Document
```
phases/01-define/artifacts/prd/
├── templates/
│   ├── overview.md          # Template for project overview
│   ├── requirements.md      # Template for requirements list
│   └── success-metrics.md   # Template for success criteria
├── prompt.md                # Single prompt that fills all templates
├── meta.yml                 # Artifact metadata
└── examples/               # Example outputs
```

**meta.yml:**
```yaml
artifact:
  name: Product Requirements Document
  type: document
  templates:
    - path: templates/overview.md
      output: docs/prd/overview.md
    - path: templates/requirements.md
      output: docs/prd/requirements.md
    - path: templates/success-metrics.md
      output: docs/prd/success-metrics.md
  required: true
```
Result: Three structured documents, each following its template.

### Action Example: Refactor Authentication
```
phases/03-implement/actions/refactor-auth/
├── prompt.md               # Instructions for the refactoring
└── meta.yml               # Action metadata
```

**meta.yml:**
```yaml
action:
  name: Refactor Authentication System
  description: Update all authentication code to use OAuth2
  affects:
    - files: ["src/**/*.ts", "src/**/*.js"]
    - artifacts: []  # Doesn't affect specific artifacts
```
Result: Arbitrary changes across many existing files.

### Key Differences
- **Artifacts**: Templates define structure → Predictable outputs with fixed formats
- **Actions**: No templates → Flexible operations that can do anything
- **Location**: Both live under their respective phase directories, making ownership clear

## Implementation Priorities

### Phase 1: Core Infrastructure (MVP)
- Phase execution engine with input gates and exit criteria
- Artifact generation from templates + prompts
- Action execution for multi-file operations
- State management and phase transitions
- CLI commands for workflow operations

### Phase 2: Advanced Features
- Parallel phase execution
- Cross-phase dependencies
- Validation framework for gates and criteria
- Error recovery and rollback
- Progress tracking and reporting

### Phase 3: Enterprise Features
- External system integrations
- Advanced automation hooks
- Performance optimization
- Analytics and metrics
- Multi-user collaboration

## Testing Strategy

### Unit Testing
- Template processing logic
- Variable substitution
- Phase transition validation
- Artifact generation
- Error handling paths

### Integration Testing
- End-to-end workflow execution
- CLI command interface
- External system integrations
- State persistence
- Recovery scenarios

### Performance Testing
- Large workflow catalog handling
- Concurrent execution scalability
- Memory usage optimization
- Template processing speed

### Security Testing
- Input validation
- Template injection attacks
- Access control enforcement
- Data protection compliance

## Related Documentation

- [Workflow Commands Specification](./commands.md)
- [Create Command Specification](./create-command.md)
- [Apply Command Specification](./apply-command.md)
- [Usage Documentation](../../../usage/workflows/overview.md)
- [Architecture Documentation](../../../architecture/system-design.md)