# User Story: US-006 - Get Command Help

**Story ID**: US-006
**Feature**: FEAT-001 - Core CLI Framework
**Priority**: P0
**Status**: Draft
**Created**: 2025-01-14
**Updated**: 2025-01-14

## Story

**As a** developer
**I want** to get help for any command
**So that** I can learn how to use DDX effectively

## Acceptance Criteria

- [ ] **Given** I need general help, **when** I run `ddx help`, **then** I see an overview of all available commands with brief descriptions
- [ ] **Given** I need command-specific help, **when** I run `ddx help <command>`, **then** detailed help for that command is displayed
- [ ] **Given** I prefer flag syntax, **when** I run `ddx <command> --help`, **then** the same help information is shown
- [ ] **Given** I'm viewing help, **when** I read the output, **then** I see practical examples for common use cases
- [ ] **Given** a command has flags, **when** I view its help, **then** all available flags are listed with their defaults
- [ ] **Given** I need more information, **when** I view help, **then** links to online documentation are provided
- [ ] **Given** arguments are required, **when** I view help, **then** required vs optional arguments are clearly indicated
- [ ] **Given** commands have aliases, **when** I view help, **then** available aliases are shown

## Definition of Done

- [ ] Help command implemented for all commands
- [ ] Help flag (--help) working on all commands
- [ ] Examples included in help output
- [ ] Flag documentation complete with defaults
- [ ] Required/optional arguments clearly marked
- [ ] Online documentation links included
- [ ] Alias information displayed
- [ ] Unit tests for help generation
- [ ] Help text reviewed for clarity and accuracy
- [ ] Consistent formatting across all help output

## Technical Notes

### Implementation Considerations
- Help text should be maintained close to command implementation
- Consider auto-generating help from code comments
- Help should work offline (no network required)
- Terminal width should be respected for formatting
- Color coding for better readability (if terminal supports)

### Error Scenarios
- Invalid command name provided to help
- Help text missing for a command
- Terminal too narrow for formatted output
- Help for deprecated commands

## Validation Scenarios

### Scenario 1: General Help
1. Run `ddx help`
2. **Expected**: List of all commands with short descriptions

### Scenario 2: Command-Specific Help
1. Run `ddx help init`
2. **Expected**: Detailed help for init command with examples

### Scenario 3: Help Flag
1. Run `ddx apply --help`
2. **Expected**: Same output as `ddx help apply`

### Scenario 4: Invalid Command
1. Run `ddx help invalid-command`
2. **Expected**: Error message with suggestion for similar commands

### Scenario 5: Nested Command Help
1. Run `ddx help list templates`
2. **Expected**: Help specific to listing templates

## User Persona

### Primary: New User
- **Role**: Developer new to DDX
- **Goals**: Learn how to use DDX commands effectively
- **Pain Points**: Unclear syntax, missing examples
- **Technical Level**: Varies widely

### Secondary: Experienced User
- **Role**: Regular DDX user
- **Goals**: Quick reference for advanced options
- **Pain Points**: Remembering exact flag names and options
- **Technical Level**: Intermediate to advanced

## Dependencies

- Cobra framework's help generation system
- Command documentation must be maintained
- Terminal capabilities for formatting

## Related Stories

- US-001: Initialize DDX in Project (help for getting started)
- US-008: Check DDX Version (often checked with help)
- All other command stories (each needs help text)

---
*This user story is part of FEAT-001: Core CLI Framework*