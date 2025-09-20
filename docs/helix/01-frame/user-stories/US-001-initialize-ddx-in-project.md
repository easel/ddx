# User Story: US-001 - Initialize DDX in Project

**Story ID**: US-001
**Feature**: FEAT-001 - Core CLI Framework
**Priority**: P0
**Status**: Draft
**Created**: 2025-01-14
**Updated**: 2025-01-14

## Story

**As a** developer
**I want** to initialize DDX in my existing project
**So that** I can start using shared assets immediately

## Acceptance Criteria

- [ ] **Given** I am in a project directory without DDX, **when** I run `ddx init`, **then** a `.ddx` directory structure is created with proper organization
- [ ] **Given** I run `ddx init`, **when** the command executes, **then** interactive prompts guide me through configuration options
- [ ] **Given** I want to use a specific template, **when** I run `ddx init --template <name>`, **then** the specified template is applied during initialization
- [ ] **Given** initialization completes, **when** I check the project, **then** a `.ddx.yml` configuration file exists with my settings
- [ ] **Given** DDX is being initialized, **when** the process runs, **then** git subtree connection is established to the master repository
- [ ] **Given** my project has specific requirements, **when** I run `ddx init`, **then** project compatibility is validated before proceeding
- [ ] **Given** initialization succeeds or fails, **when** the process completes, **then** clear success or failure feedback is provided with next steps

## Definition of Done

- [ ] Init command implemented with all flags and options
- [ ] Interactive prompt flow completed and user-friendly
- [ ] Configuration file generation working correctly
- [ ] Git subtree initialization successful
- [ ] Project validation logic implemented
- [ ] Unit tests written and passing (>80% coverage)
- [ ] Integration tests for initialization scenarios
- [ ] Documentation updated with initialization examples
- [ ] Error handling for all edge cases

## Technical Notes

### Implementation Considerations
- Must handle existing git repositories gracefully
- Should detect and warn about conflicting configurations
- Need to validate template compatibility before applying
- Consider providing rollback if initialization fails partway
- Interactive prompts should have sensible defaults

### Error Scenarios
- Project already initialized with DDX
- Incompatible project structure
- Template not found or invalid
- Git not installed or configured
- Insufficient permissions to create directories
- Network issues when fetching templates

## Validation Scenarios

### Scenario 1: Basic Initialization
1. Navigate to a project without DDX
2. Run `ddx init`
3. Answer prompts with default values
4. **Expected**: DDX initialized successfully with default configuration

### Scenario 2: Template-Based Initialization
1. Navigate to a new project directory
2. Run `ddx init --template nextjs`
3. **Expected**: DDX initialized with NextJS template applied

### Scenario 3: Re-initialization Attempt
1. Navigate to DDX-enabled project
2. Run `ddx init` again
3. **Expected**: Clear message that DDX is already initialized with option to reconfigure

### Scenario 4: Non-Interactive Mode
1. Run `ddx init --yes` to accept all defaults
2. **Expected**: DDX initialized without any prompts using default values

## User Persona

### Primary: New DDX User
- **Role**: Developer new to DDX ecosystem
- **Goals**: Quickly get started with DDX in their project
- **Pain Points**: Complex setup processes, unclear configuration options
- **Technical Level**: Comfortable with CLI tools but new to DDX

## Dependencies

- Git must be installed and configured
- Write permissions in project directory
- Network connectivity for template fetching (optional)

## Related Stories

- US-002: List Available Assets (to discover what's available after init)
- US-003: Apply Asset to Project (natural next step after initialization)
- US-007: Configure DDX Settings (for post-init configuration)

---
*This user story is part of FEAT-001: Core CLI Framework*