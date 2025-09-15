# Solution Design: [SD-011] - Persona System Implementation

**Design ID**: SD-011
**Feature**: FEAT-011 (AI Persona System)
**Status**: Approved
**Created**: 2025-01-15
**Updated**: 2025-01-15

## Executive Summary

This solution design details the technical implementation of the AI Persona System for DDX. The system enables defining reusable AI personalities as markdown files, binding them to abstract roles, and integrating them with both interactive sessions and automated workflows. The design prioritizes simplicity, portability, and community contribution.

## Architecture Overview

### System Components

```
┌─────────────────────────────────────────────────────┐
│                    CLI Layer                        │
│  ┌──────────────────────────────────────────────┐  │
│  │         persona.go Command Handler           │  │
│  │  (list, show, bind, load, status, etc.)     │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                          │
┌─────────────────────────────────────────────────────┐
│                 Core Services                       │
│  ┌────────────────┐  ┌──────────────────────────┐ │
│  │ Persona Loader │  │   Binding Manager        │ │
│  │                │  │   (.ddx.yml)             │ │
│  └────────────────┘  └──────────────────────────┘ │
│  ┌────────────────┐  ┌──────────────────────────┐ │
│  │ CLAUDE.md      │  │   Workflow Integration   │ │
│  │ Injector       │  │                          │ │
│  └────────────────┘  └──────────────────────────┘ │
└─────────────────────────────────────────────────────┘
                          │
┌─────────────────────────────────────────────────────┐
│                  Data Layer                         │
│  ┌────────────────┐  ┌──────────────────────────┐ │
│  │   /personas/   │  │     .ddx.yml             │ │
│  │   (MD files)   │  │   (bindings)             │ │
│  └────────────────┘  └──────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

## Data Models

### Persona Definition

```go
type Persona struct {
    Name        string   `yaml:"name"`
    Roles       []string `yaml:"roles"`
    Description string   `yaml:"description"`
    Tags        []string `yaml:"tags"`
    Content     string   // Markdown content (body)
}
```

### Persona Binding

```go
type PersonaBinding struct {
    Role    string
    Persona string
}

type PersonaConfig struct {
    Bindings  map[string]string            `yaml:"persona_bindings"`
    Overrides map[string]map[string]string `yaml:"overrides,omitempty"`
}
```

### File Format

```markdown
---
name: strict-code-reviewer
roles: [code-reviewer, security-analyst]
description: Uncompromising code quality enforcer
tags: [strict, security, production]
---

# Strict Code Reviewer

You are an experienced code reviewer who enforces high standards...
[Rest of persona content]
```

## Component Design

### CLI Command Structure

```bash
ddx persona
├── list      [--role ROLE] [--tag TAG]
├── show      PERSONA_NAME
├── bind      ROLE PERSONA_NAME
├── bindings
├── unbind    ROLE
├── load      [PERSONA_NAME] [--role ROLE]
├── unload
└── status
```

### Persona Loader Service

```go
type PersonaLoader interface {
    // Load and parse persona from file
    LoadPersona(name string) (*Persona, error)

    // List all available personas
    ListPersonas() ([]*Persona, error)

    // Filter personas by role
    FindByRole(role string) ([]*Persona, error)

    // Filter personas by tags
    FindByTags(tags []string) ([]*Persona, error)
}
```

### Binding Manager Service

```go
type BindingManager interface {
    // Get persona for role
    GetBinding(role string) (string, error)

    // Set binding
    SetBinding(role, persona string) error

    // Get all bindings
    GetAllBindings() (map[string]string, error)

    // Remove binding
    RemoveBinding(role string) error

    // Get workflow override
    GetOverride(workflow, role string) (string, error)
}
```

### CLAUDE.md Injector

```go
type ClaudeInjector interface {
    // Inject single persona
    InjectPersona(persona *Persona, role string) error

    // Inject multiple personas
    InjectMultiple(personas map[string]*Persona) error

    // Remove all personas
    RemovePersonas() error

    // Check if personas are loaded
    GetLoadedPersonas() ([]string, error)
}
```

## Implementation Details

### Directory Structure

```
/personas/
├── README.md                     # Documentation
├── strict-code-reviewer.md       # Persona definitions
├── balanced-code-reviewer.md
├── test-engineer-tdd.md
├── test-engineer-bdd.md
├── architect-systems.md
├── developer-senior-golang.md
└── ...

/cli/
├── cmd/
│   ├── persona.go               # Main command
│   ├── persona_list.go          # Subcommands
│   ├── persona_bind.go
│   └── ...
└── internal/
    ├── persona/
    │   ├── loader.go            # Persona loading
    │   ├── binding.go           # Binding management
    │   └── claude.go            # CLAUDE.md integration
    └── ...
```

### Persona Loading Process

1. **Discovery Phase**
   - Scan `/personas/` directory
   - Parse YAML frontmatter
   - Extract markdown content
   - Build persona index

2. **Selection Phase**
   - Check project bindings in `.ddx.yml`
   - Apply workflow overrides if present
   - Fall back to user selection if needed

3. **Injection Phase**
   - Read current CLAUDE.md
   - Inject persona content with markers
   - Preserve existing content
   - Write updated file

### CLAUDE.md Format

```markdown
# CLAUDE.md

[Original project instructions...]

<!-- PERSONAS:START -->
## Active Personas

### Code Reviewer: strict-code-reviewer
You are an experienced code reviewer...

### Test Engineer: test-engineer-tdd
You are a test engineer focused on TDD...

When responding, adopt the appropriate persona based on the task.
<!-- PERSONAS:END -->
```

### Workflow Integration

```go
func ExecuteArtifact(artifact *Artifact, config *Config) error {
    // Get required role
    role := artifact.RequiredRole
    if role == "" {
        role = artifact.Phase.RequiredRole
    }

    // Get persona for role
    personaName := config.GetBinding(role)
    if personaName == "" {
        personaName = promptUserForPersona(role)
    }

    // Load persona
    persona, err := LoadPersona(personaName)
    if err != nil {
        return err
    }

    // Combine persona + artifact prompt
    fullPrompt := combinePrompts(persona.Content, artifact.Prompt)

    // Execute with AI
    return executeWithAI(fullPrompt)
}
```

## Configuration Schema

### .ddx.yml Extension

```yaml
# Existing configuration...

persona_bindings:
  code-reviewer: strict-code-reviewer
  test-engineer: test-engineer-tdd
  architect: architect-systems
  developer: developer-senior-golang

  # Optional workflow-specific overrides
  overrides:
    helix:
      test-engineer: test-engineer-bdd
    agile:
      developer: developer-fullstack-react
```

### Workflow Definition Extension

```yaml
phases:
  - id: test
    name: Write Tests First
    required_role: test-engineer  # New field

artifacts:
  - name: test-plan
    required_role: test-engineer  # New field
```

## Error Handling

### Common Error Scenarios

1. **Persona Not Found**
   ```
   Error: Persona 'advanced-reviewer' not found
   Available personas for role 'code-reviewer':
     - strict-code-reviewer
     - balanced-code-reviewer
   ```

2. **No Binding for Role**
   ```
   Warning: No persona bound to role 'architect'
   Use 'ddx persona bind architect <persona>' to set binding
   ```

3. **Invalid Persona Format**
   ```
   Warning: Skipping invalid persona 'broken.md'
   Error: Missing required field 'name' in frontmatter
   ```

## Performance Considerations

- **Lazy Loading**: Only parse personas when needed
- **Caching**: Cache parsed personas for session
- **Incremental Updates**: Only update changed sections in CLAUDE.md
- **Pagination**: Support pagination for large persona lists

## Security Considerations

- **No Code Execution**: Personas are text only, no executable code
- **Sandboxing**: Personas cannot access system resources
- **Validation**: Validate YAML frontmatter to prevent injection
- **Git Integration**: All changes tracked in version control

## Testing Strategy

### Unit Tests
- Persona parsing and validation
- Binding management logic
- CLAUDE.md injection/removal
- Role resolution

### Integration Tests
- CLI command execution
- File system operations
- Configuration updates
- Workflow integration

### E2E Tests
- Complete persona workflow
- Multi-persona loading
- Workflow execution with personas

## Migration Path

1. **Phase 1**: Implement core persona system
2. **Phase 2**: Add workflow integration
3. **Phase 3**: Migrate existing prompts to personas
4. **Phase 4**: Community contribution features

## Future Enhancements

- Persona inheritance/composition
- Private team repositories
- Persona analytics
- AI-assisted persona creation
- Dynamic persona switching

---

*This solution design provides a complete technical blueprint for implementing the persona system while maintaining simplicity and extensibility.*