# US-036: List Available MCP Servers

## Story Overview

**Story ID**: US-036
**Feature**: FEAT-014 (MCP Server Management)  
**Priority**: P0  
**Points**: 3  
**Sprint**: 1  

## User Story

**As a** developer using Claude Code  
**I want to** list available MCP servers with filtering and search  
**So that** I can discover useful integrations for my development workflow  

## Acceptance Criteria

### AC1: Display All Available Servers
**Given** the MCP server registry is available  
**When** I run `ddx mcp list`  
**Then** I should see all available MCP servers with:
- Server name
- Brief description
- Category
- Installation status (✅ installed, ⬜ not installed)

### AC2: Filter by Category
**Given** multiple categories of MCP servers exist  
**When** I run `ddx mcp list --category development`  
**Then** I should only see servers in the "development" category  
**And** the output should indicate the active filter  

### AC3: Search Functionality
**Given** I want to find servers related to "git"  
**When** I run `ddx mcp list --search git`  
**Then** I should see servers with "git" in name or description  
**And** the search term should be highlighted in output  

### AC4: Show Installation Status
**Given** I have some MCP servers already installed  
**When** I run `ddx mcp list`  
**Then** installed servers should show ✅ indicator  
**And** not-installed servers should show ⬜ indicator  

### AC5: Detailed View
**Given** I want more information about a server  
**When** I run `ddx mcp list --verbose`  
**Then** I should see additional details:
- Required environment variables
- Author/maintainer
- Version
- Documentation link

### AC6: Empty Results Handling
**Given** no servers match my search criteria  
**When** I run `ddx mcp list --search nonexistent`  
**Then** I should see a helpful message: "No MCP servers found matching 'nonexistent'"  
**And** suggestions for broadening the search  

## Definition of Done

- [ ] Command implemented with all flags
- [ ] Unit tests cover all acceptance criteria
- [ ] Integration tests verify registry loading
- [ ] Help documentation includes examples
- [ ] Error messages are helpful and actionable
- [ ] Performance: results display in <100ms
- [ ] Works offline with cached registry data

## Technical Notes

### Implementation Details
- Load registry from `mcp-servers/registry.yml`
- Cache registry data for 15 minutes
- Support multiple output formats (table, json, yaml)
- Use color coding for better readability

### Dependencies
- Registry YAML parser
- Terminal output formatter
- Cache management system
- Claude configuration detector

## Test Scenarios

### Happy Path
1. List all servers - verify complete list
2. Filter by valid category - verify filtered results
3. Search for known term - verify matches found
4. List with verbose flag - verify detailed output

### Edge Cases
1. Empty registry - show appropriate message
2. Invalid category - show available categories
3. Network offline - use cached data
4. Corrupted cache - regenerate from registry
5. No Claude installation - still show list

### Error Cases
1. Invalid registry format - error with fix instructions
2. Permission denied on cache - fallback to temp directory
3. Multiple filters conflict - clear precedence rules

## UX Considerations

### Output Format
```
📋 Available MCP Servers (15 total, 3 installed)

Development (5 servers):
  ✅ github         - GitHub integration for repository access
  ✅ gitlab         - GitLab integration for repository access  
  ⬜ bitbucket      - Bitbucket integration for repository access
  ⬜ code-search    - Search across multiple code repositories
  ⬜ pull-request   - Manage pull requests across platforms

Database (4 servers):
  ✅ postgres       - PostgreSQL database integration
  ⬜ mysql          - MySQL database integration
  ⬜ mongodb        - MongoDB database integration
  ⬜ redis          - Redis cache and database integration

[Showing 2 of 5 categories. Use --all to see everything]
```

### Interactive Features
- Pagination for long lists
- Column sorting (name, category, status)
- Quick install prompt: "Install github? (y/n)"

## Related Stories

- US-037: Install MCP Server
- US-038: Search MCP Server Details (future)
- US-039: Update Registry Cache (future)

## Change Log

| Date | Author | Changes |
|------|--------|----------|
| 2025-01-15 | System | Initial story creation |