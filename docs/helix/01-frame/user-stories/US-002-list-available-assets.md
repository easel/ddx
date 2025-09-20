# User Story: US-002 - List Available Assets

**Story ID**: US-002
**Feature**: FEAT-001 - Core CLI Framework
**Priority**: P0
**Status**: Draft
**Created**: 2025-01-14
**Updated**: 2025-01-14

## Story

**As a** developer
**I want** to browse available assets by category
**So that** I can discover useful resources

## Acceptance Criteria

- [ ] **Given** I have DDX initialized, **when** I run `ddx list`, **then** I see all available asset categories with counts
- [ ] **Given** I want to see prompts, **when** I run `ddx list prompts`, **then** all available prompts are displayed with descriptions
- [ ] **Given** I want to see templates, **when** I run `ddx list templates`, **then** all available templates are shown with descriptions
- [ ] **Given** I want to see patterns, **when** I run `ddx list patterns`, **then** all available patterns are listed with descriptions
- [ ] **Given** assets are listed, **when** I view the output, **then** each item shows name, description, and relevant tags
- [ ] **Given** I want to filter results, **when** I run `ddx list --filter <keyword>`, **then** only matching assets are displayed
- [ ] **Given** I need machine-readable output, **when** I run `ddx list --json`, **then** results are formatted as valid JSON
- [ ] **Given** usage tracking is enabled, **when** I list assets, **then** usage statistics are displayed if available

## Definition of Done

- [ ] List command implemented with category filtering
- [ ] Output formatting for human and machine readability
- [ ] Filter functionality working across all fields
- [ ] JSON output mode implemented and validated
- [ ] Asset metadata properly displayed
- [ ] Unit tests written and passing (>80% coverage)
- [ ] Integration tests for various listing scenarios
- [ ] Documentation updated with list command examples
- [ ] Performance acceptable for large asset collections

## Technical Notes

### Implementation Considerations
- Must handle large asset collections efficiently
- Should cache asset listings for performance
- Need to support multiple output formats (table, list, JSON)
- Consider pagination for very long lists
- Tags and categories should be searchable

### Error Scenarios
- No assets available in category
- Invalid category specified
- Corrupted asset metadata
- Filter returns no results
- Asset directory not accessible

## Validation Scenarios

### Scenario 1: List All Categories
1. Run `ddx list` without arguments
2. **Expected**: See all categories with asset counts

### Scenario 2: List Specific Category
1. Run `ddx list templates`
2. **Expected**: All templates listed with descriptions and metadata

### Scenario 3: Filtered Listing
1. Run `ddx list prompts --filter "claude"`
2. **Expected**: Only prompts containing "claude" in name/description/tags

### Scenario 4: JSON Output
1. Run `ddx list patterns --json`
2. **Expected**: Valid JSON array of pattern objects

### Scenario 5: Empty Category
1. Run `ddx list` for a category with no assets
2. **Expected**: Clear message that no assets are available in this category

## User Persona

### Primary: Exploring Developer
- **Role**: Developer looking for reusable assets
- **Goals**: Discover relevant templates, patterns, and prompts
- **Pain Points**: Hard to find what's available, poor discoverability
- **Technical Level**: Varies from beginner to expert

## Dependencies

- DDX must be initialized in project
- Asset metadata files must be properly formatted
- File system access to asset directories

## Related Stories

- US-001: Initialize DDX in Project (prerequisite)
- US-003: Apply Asset to Project (natural next action after discovery)
- US-004: Update Assets from Master (to get latest assets)

---
*This user story is part of FEAT-001: Core CLI Framework*