# User Story: [US-032] - Workflow Author Requiring Roles

**Story ID**: US-032
**Feature**: FEAT-011 (AI Persona System)
**Priority**: P1
**Status**: Defined

## User Story

**As a** workflow author creating reusable development workflows,
**I want to** specify required roles for phases and artifacts,
**So that** appropriate expertise is applied regardless of which specific persona each project chooses to use.

## Acceptance Criteria

1. **Phase-Level Requirements**
   - Can specify `required_role` in workflow phase definition
   - Role is abstract (e.g., "test-engineer" not specific persona)
   - Multiple phases can require same or different roles

2. **Artifact-Level Requirements**
   - Individual artifacts can specify `required_role`
   - Artifact role overrides phase role if specified
   - Clear precedence rules for role selection

3. **Workflow Validation**
   - Workflow validates that required roles are recognized
   - Warning if role has no common personas available
   - Documentation of which roles workflow expects

4. **Execution Behavior**
   - When workflow runs, appropriate persona is selected
   - Uses project's binding for the required role
   - Falls back to prompting if no binding exists

5. **Documentation**
   - Workflow documentation shows required roles
   - Clear explanation of each role's purpose
   - Guidance on selecting appropriate personas

## Example Workflow Definition

```yaml
# workflows/helix/workflow.yml
name: helix
phases:
  - id: frame
    name: Frame the Problem
    required_role: architect  # Architect role for problem framing

  - id: design
    name: Design the Solution
    required_role: architect  # Continue with architect

  - id: test
    name: Write Tests First
    required_role: test-engineer  # Test engineer for test phase

  - id: build
    name: Build to Pass Tests
    required_role: developer  # Developer for implementation
```

```yaml
# workflows/helix/phases/03-test/artifacts/test-plan/metadata.yml
name: test-plan
type: artifact
required_role: test-engineer  # Specific role for this artifact
description: Comprehensive test plan following TDD principles
```

## Example Execution

```bash
# Author defines workflow with roles
$ ddx workflow validate helix
Workflow 'helix' requires roles:
  - architect (phases: frame, design)
  - test-engineer (phase: test)
  - developer (phase: build)
  - code-reviewer (artifact: code-review)

# User runs workflow (personas selected automatically)
$ ddx workflow run helix --phase test
Phase 'test' requires role: test-engineer
Using persona: test-engineer-tdd (from project bindings)
Generating artifacts...

# If no binding exists
$ ddx workflow run helix --phase test
Phase 'test' requires role: test-engineer
No persona bound to role 'test-engineer'

Available personas for this role:
1. test-engineer-tdd
2. test-engineer-bdd
3. test-engineer-comprehensive

Select persona (1-3): 1
Using persona: test-engineer-tdd
```

## Technical Notes

- Roles defined at workflow level, not hardcoded personas
- Enables workflow portability across projects
- Supports both phase and artifact level requirements
- Graceful fallback when bindings don't exist

## Dependencies

- Workflow definition format supports `required_role` field
- Persona system can resolve role → persona mappings
- Workflow execution engine integrates with persona system

## Related Stories

- US-030: Developer Loading Personas for Session
- US-031: Team Lead Binding Personas to Roles
- US-033: Developer Contributing Personas