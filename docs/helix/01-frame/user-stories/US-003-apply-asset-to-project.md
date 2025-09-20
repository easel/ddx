# User Story: US-003 - Apply Asset to Project

**Story ID**: US-003
**Feature**: FEAT-001 - Core CLI Framework
**Priority**: P0
**Status**: Draft
**Created**: 2025-01-14
**Updated**: 2025-01-14

## Story

**As a** developer
**I want** to apply an asset to my current project
**So that** I can use proven solutions

## Acceptance Criteria

- [ ] **Given** I have identified an asset to use, **when** I run `ddx apply <asset-path>`, **then** the specified asset is applied to my project
- [ ] **Given** an asset requires specific dependencies, **when** I attempt to apply it, **then** compatibility is validated before any changes are made
- [ ] **Given** an asset has variables, **when** it's applied, **then** variable substitution occurs using values from my configuration
- [ ] **Given** I want to preview changes, **when** I run `ddx apply <asset> --dry-run`, **then** I see what would change without modifying files
- [ ] **Given** an operation could overwrite files, **when** I apply an asset, **then** I am prompted to confirm destructive operations
- [ ] **Given** files already exist, **when** applying an asset, **then** file conflicts are detected and handled gracefully with clear options
- [ ] **Given** something goes wrong, **when** an apply operation fails, **then** rollback instructions are provided to restore previous state
- [ ] **Given** an asset is successfully applied, **when** the operation completes, **then** the local asset cache is updated for future reference

## Definition of Done

- [ ] Apply command implemented with all options
- [ ] Asset validation logic complete
- [ ] Variable substitution system working
- [ ] Dry-run mode shows accurate preview
- [ ] Conflict resolution implemented
- [ ] Confirmation prompts for destructive operations
- [ ] Rollback mechanism documented
- [ ] Unit tests written and passing (>80% coverage)
- [ ] Integration tests for apply scenarios
- [ ] Documentation updated with apply examples

## Technical Notes

### Implementation Considerations
- Must preserve file permissions and attributes
- Should create backups before overwriting files
- Need atomic operations (all-or-nothing)
- Variable substitution should be recursive
- Consider supporting conditional application based on project type

### Error Scenarios
- Asset not found
- Incompatible asset for project type
- File write permissions denied
- Disk space insufficient
- Template syntax errors
- Variable values missing
- Circular variable references

## Validation Scenarios

### Scenario 1: Simple Asset Application
1. Run `ddx apply templates/config/eslint`
2. **Expected**: ESLint configuration files added to project

### Scenario 2: Dry Run Preview
1. Run `ddx apply patterns/auth/jwt --dry-run`
2. **Expected**: See list of files that would be created/modified without actual changes

### Scenario 3: Conflict Resolution
1. Apply an asset that would overwrite existing file
2. **Expected**: Prompted with options (overwrite, skip, backup, merge)

### Scenario 4: Variable Substitution
1. Apply template with variables like `{{PROJECT_NAME}}`
2. **Expected**: Variables replaced with values from .ddx.yml

### Scenario 5: Failed Application
1. Apply asset when disk is full
2. **Expected**: Operation fails cleanly with rollback instructions

## User Persona

### Primary: Efficient Developer
- **Role**: Developer wanting to quickly add functionality
- **Goals**: Apply proven solutions without reinventing the wheel
- **Pain Points**: Manual copying and adapting of code, missing dependencies
- **Technical Level**: Intermediate to advanced

## Dependencies

- DDX must be initialized
- Asset must exist in DDX repository
- Configuration file for variable values
- Write permissions in project directory

## Related Stories

- US-001: Initialize DDX in Project (prerequisite)
- US-002: List Available Assets (to discover what to apply)
- US-004: Update Assets from Master (to get latest versions)
- US-007: Configure DDX Settings (for variable values)

---
*This user story is part of FEAT-001: Core CLI Framework*