# User Story: [US-039] - Developer Loading Personas for Session

**Story ID**: US-030
**Feature**: FEAT-011 (AI Persona System)
**Priority**: P1
**Status**: Defined

## User Story

**As a** developer working on a project with defined personas,
**I want to** load all my project's bound personas with a single command,
**So that** my AI assistant can switch between appropriate personalities based on the context of my requests.

## Acceptance Criteria

1. **Single Command Loading**
   - Running `ddx persona load` without parameters loads all bound personas
   - Personas are loaded based on project's `.ddx.yml` configuration
   - Command provides feedback on which personas were loaded

2. **CLAUDE.md Integration**
   - All loaded personas are injected into project's CLAUDE.md
   - Each persona is clearly labeled with its role
   - Existing CLAUDE.md content is preserved

3. **Multiple Active Personas**
   - Multiple personas can be active simultaneously
   - AI can switch between personas based on task context
   - Clear structure showing available personas and their roles

4. **Status Visibility**
   - `ddx persona status` shows all currently loaded personas
   - Display shows role → persona mappings
   - Indicates if personas are from default bindings or overrides

5. **Unload Capability**
   - `ddx persona unload` removes all personas from CLAUDE.md
   - Original CLAUDE.md content remains intact
   - Confirmation that personas were unloaded

## Example Scenario

```bash
# Developer starts work on project
$ ddx persona load
Loading bound personas:
  - code-reviewer: strict-code-reviewer
  - test-engineer: test-engineer-tdd
  - architect: architect-systems
Updated CLAUDE.md with 3 personas

# Developer checks status
$ ddx persona status
Active personas:
  - code-reviewer: strict-code-reviewer
  - test-engineer: test-engineer-tdd
  - architect: architect-systems

# Developer works with AI - it automatically uses appropriate personas
# "Review this code" → uses code-reviewer persona
# "Write tests for this" → uses test-engineer persona
# "Design the API" → uses architect persona

# End of session
$ ddx persona unload
Removed all personas from CLAUDE.md
```

## Technical Notes

- Personas are read from `/personas/` directory
- Bindings come from project's `.ddx.yml` file
- CLAUDE.md modifications are non-destructive
- Support for both project and user-level persona directories

## Dependencies

- Persona definition files exist in `/personas/`
- Project has `.ddx.yml` with `persona_bindings` section
- CLAUDE.md exists or can be created

## Related Stories

- US-031: Team Lead Binding Personas to Roles
- US-032: Workflow Author Requiring Roles
- US-034: Developer Discovering Personas