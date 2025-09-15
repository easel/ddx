# User Story: Configure Resource Selection

**Story ID**: US-020
**Feature**: FEAT-003 (Configuration Management)
**Priority**: P1
**Status**: Draft
**Created**: 2025-01-14
**Updated**: 2025-01-14

## Story
**As a** developer
**I want to** specify which resources to include
**So that** I only get relevant assets for my project

## Description
This story allows developers to configure which DDX resources (templates, patterns, prompts, workflows) should be available for their project. Not every project needs all available resources, and selective inclusion helps keep projects focused and reduces clutter. This feature supports patterns, wildcards, and explicit include/exclude rules.

## Acceptance Criteria
- [ ] **Given** configuration, **when** resources specified, **then** selection is honored during operations
- [ ] **Given** patterns, **when** used, **then** wildcards work for resource selection
- [ ] **Given** resources, **when** configured, **then** include/exclude rules are applied correctly
- [ ] **Given** resource types, **when** organized, **then** grouping by category is supported
- [ ] **Given** resource dependencies, **when** present, **then** they are automatically included
- [ ] **Given** resource config, **when** complete, **then** preview of selection is available
- [ ] **Given** resource paths, **when** specified, **then** availability is validated
- [ ] **Given** resources, **when** listed, **then** tree view shows structure

## Business Value
- Reduces project complexity by including only needed resources
- Speeds up operations by working with smaller resource sets
- Prevents accidental use of inappropriate resources
- Enables project-specific resource curation
- Supports different resource needs across project types

## Definition of Done
- [ ] Resource selection configuration syntax is implemented
- [ ] Wildcard pattern matching works correctly
- [ ] Include/exclude logic is functional
- [ ] Resource grouping by category is supported
- [ ] Dependency resolution is implemented
- [ ] Preview functionality shows selected resources
- [ ] Resource availability validation works
- [ ] Tree view display is implemented
- [ ] Unit tests cover selection logic
- [ ] Integration tests verify resource filtering
- [ ] Documentation explains selection patterns
- [ ] All acceptance criteria are met and verified

## Technical Considerations
[NEEDS CLARIFICATION: These will be defined in the Design phase]
- Pattern matching algorithm efficiency
- Dependency graph resolution
- Resource discovery mechanism
- Caching of resource metadata

## Dependencies
- **Prerequisite**: US-017 (Initialize Configuration) must be completed
- **Related**: Affects all resource application commands

## Assumptions
- Resources are organized in a predictable structure
- Resource names follow consistent conventions
- [NEEDS CLARIFICATION: Are resource dependencies explicitly defined?]
- [NEEDS CLARIFICATION: How are new resources discovered?]

## Edge Cases
- Conflicting include/exclude rules
- Circular dependencies between resources
- Non-existent resources in configuration
- Wildcard patterns matching nothing
- Excluding required dependencies
- Very large resource selections
- Resources added after configuration

## Examples

### Resource Selection Configuration
```yaml
resources:
  prompts:
    include:
      - "code-review"
      - "testing/*"
      - "documentation/api"
    exclude:
      - "testing/experimental/*"

  templates:
    include:
      - "nextjs"
      - "react-*"
    exclude:
      - "react-legacy"

  patterns:
    include:
      - "error-handling"
      - "logging"
      - "auth/*"

  workflows:
    include:
      - "helix"
      - "agile/*"
```

### Wildcard Patterns
- `*` - Matches any characters except `/`
- `**` - Matches any characters including `/`
- `?` - Matches single character
- `[abc]` - Matches any character in brackets
- `{option1,option2}` - Matches any of the options

### Command Usage
```bash
# List selected resources
ddx list --selected

# Preview what would be included
ddx config resources --preview

# Apply only configured resources
ddx apply pattern error-handling  # Works if included
ddx apply pattern unused-pattern   # Fails if not included
```

## User Feedback
*To be collected during implementation and testing*

## Notes
- Resource selection significantly impacts user experience
- Should balance flexibility with simplicity
- Consider providing preset resource bundles for common use cases
- May need versioning for resource compatibility

---
*Story is part of FEAT-003 (Configuration Management)*