# User Story: [US-040] - Team Lead Binding Personas to Roles

**Story ID**: US-031
**Feature**: FEAT-011 (AI Persona System)
**Priority**: P1
**Status**: Defined

## User Story

**As a** team lead managing a development project,
**I want to** bind specific personas to roles in my project configuration,
**So that** all team members use consistent AI personalities and our AI interactions maintain quality standards.

## Acceptance Criteria

1. **Binding Commands**
   - `ddx persona bind <role> <persona>` creates a role → persona binding
   - Binding is saved to project's `.ddx.yml` file
   - Command validates that persona exists and supports the role

2. **View Bindings**
   - `ddx persona bindings` shows all current bindings
   - Display includes role name and bound persona
   - Indicates if there are workflow-specific overrides

3. **Unbind Capability**
   - `ddx persona unbind <role>` removes a binding
   - Configuration file is updated appropriately
   - Confirmation message shows unbinding succeeded

4. **Validation**
   - Cannot bind non-existent persona
   - Warning if persona doesn't declare support for role
   - Suggest alternatives if requested persona not found

5. **Override Support**
   - Can define workflow-specific overrides
   - Overrides are clearly marked in configuration
   - Bindings show which are defaults vs overrides

## Example Scenario

```bash
# Team lead explores available personas
$ ddx persona list --role code-reviewer
Available personas for role 'code-reviewer':
  - strict-code-reviewer: Uncompromising quality standards
  - balanced-code-reviewer: Pragmatic approach with teaching focus
  - security-code-reviewer: Security-first review approach

# Bind strict reviewer for the team
$ ddx persona bind code-reviewer strict-code-reviewer
Bound 'strict-code-reviewer' to role 'code-reviewer' in .ddx.yml

# Bind other roles
$ ddx persona bind test-engineer test-engineer-tdd
$ ddx persona bind architect architect-systems

# View current bindings
$ ddx persona bindings
Project Persona Bindings:
  Role            → Persona
  ─────────────────────────────────────
  code-reviewer   → strict-code-reviewer
  test-engineer   → test-engineer-tdd
  architect       → architect-systems
  developer       → (unbound)

Workflow Overrides:
  helix/test-engineer → test-engineer-bdd

# Later, change binding
$ ddx persona bind code-reviewer balanced-code-reviewer
Updated binding: 'code-reviewer' now uses 'balanced-code-reviewer'
```

## Configuration Result

After bindings, `.ddx.yml` contains:
```yaml
persona_bindings:
  code-reviewer: strict-code-reviewer
  test-engineer: test-engineer-tdd
  architect: architect-systems

  overrides:
    helix:
      test-engineer: test-engineer-bdd
```

## Technical Notes

- Bindings are stored in `.ddx.yml` under `persona_bindings`
- Validation against available personas in `/personas/`
- Support for nested overrides per workflow
- Changes are immediately reflected in git diff

## Dependencies

- `.ddx.yml` configuration file exists
- Personas defined in `/personas/` directory
- Write permissions to project configuration

## Related Stories

- US-030: Developer Loading Personas for Session
- US-032: Workflow Author Requiring Roles
- US-035: Developer Overriding Workflow Personas