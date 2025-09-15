# User Story: View Effective Configuration

**Story ID**: US-024
**Feature**: FEAT-003 (Configuration Management)
**Priority**: P1
**Status**: Draft
**Created**: 2025-01-14
**Updated**: 2025-01-14

## Story
**As a** developer
**I want to** see the final configuration after all overrides
**So that** I understand what settings are active

## Description
This story provides visibility into the final, effective configuration after all merging, overrides, and variable substitution has occurred. When debugging configuration issues or verifying settings, developers need to see exactly what values DDX is using, where they come from, and how they were resolved.

## Acceptance Criteria
- [ ] **Given** `ddx config show`, **when** run, **then** current configuration is displayed
- [ ] **Given** each value, **when** displayed, **then** source (file/env/default) is shown
- [ ] **Given** overridden values, **when** present, **then** they are highlighted
- [ ] **Given** large configs, **when** viewed, **then** filtering by section is supported
- [ ] **Given** output needs, **when** specified, **then** JSON/YAML formats are available
- [ ] **Given** variables, **when** computed, **then** resolved values are shown
- [ ] **Given** default values, **when** used, **then** they are clearly indicated
- [ ] **Given** multiple sources, **when** values come from different places, **then** color-coding distinguishes them

## Business Value
- Enables effective debugging of configuration issues
- Provides transparency into configuration resolution process
- Reduces time spent troubleshooting configuration problems
- Helps verify that overrides are working as expected
- Supports configuration auditing and compliance

## Definition of Done
- [ ] `ddx config show` command is implemented
- [ ] Source attribution for each value is working
- [ ] Override highlighting is functional
- [ ] Section filtering is implemented
- [ ] Multiple output formats are supported
- [ ] Variable resolution display works
- [ ] Default value indication is clear
- [ ] Color-coding for different sources is implemented
- [ ] Unit tests cover all display scenarios
- [ ] Integration tests verify configuration resolution
- [ ] Documentation explains configuration display
- [ ] All acceptance criteria are met and verified

## Technical Considerations
To be defined in technical design
- Configuration resolution and tracking mechanism
- Output formatting and color support detection
- Performance with very large configurations
- Memory usage for configuration metadata

## Dependencies
- **Prerequisite**: US-017 (Initialize Configuration) must be completed
- **Related**: US-019 (Override Configuration) for override tracking
- **Related**: US-018 (Configure Variables) for variable resolution display

## Assumptions
- Terminal supports color output (with fallback)
- Configuration resolution preserves source information
- Sensitive values should be masked in display for security
- Display should show source origin and override hierarchy

## Edge Cases
- Very large configurations that exceed screen size
- Complex nested structures with many overrides
- Variables with circular or deep references
- Configuration with missing or invalid sources
- Terminal without color support
- Piped output to files or other commands

## Examples

### Basic Display
```bash
# Show complete effective configuration
ddx config show

# Show specific section
ddx config show variables

# Show in JSON format
ddx config show --format json

# Show with source details
ddx config show --verbose
```

### Sample Output with Source Attribution
```yaml
# DDX Effective Configuration
# Generated: 2025-01-14 10:30:00

project:
  name: "my-awesome-app"          # Source: .ddx.yml:3
  type: "web"                     # Source: .ddx.dev.yml:2 (override)
  version: "1.0.0"                # Source: default

variables:
  author: "John Doe"              # Source: ${GIT_AUTHOR_NAME} → env
  api_url: "http://localhost:3000" # Source: .ddx.dev.yml:8 (override)
  log_level: "debug"              # Source: --var flag (command-line)
  database:
    host: "localhost"             # Source: .ddx.yml:12
    port: 5432                    # Source: default
    name: "my-awesome-app_db"     # Source: computed from ${PROJECT_NAME}

repository:
  url: "https://github.com/company/ddx" # Source: .ddx.yml:18
  branch: "development"           # Source: .ddx.dev.yml:12 (override)
```

### Color-Coded Display (Terminal)
- **Green**: Values from base configuration
- **Yellow**: Overridden values
- **Blue**: Environment variables
- **Cyan**: Default values
- **Magenta**: Command-line flags
- **Red**: Computed/resolved variables

### Filtering Options
```bash
# Show only variables
ddx config show variables

# Show only repository settings
ddx config show repository

# Show only overridden values
ddx config show --only-overrides

# Search for specific keys
ddx config show --filter "api*"
```

### Different Output Formats

#### YAML (Default)
```yaml
variables:
  api_url: "http://localhost:3000"  # .ddx.dev.yml (override)
```

#### JSON
```json
{
  "variables": {
    "api_url": {
      "value": "http://localhost:3000",
      "source": ".ddx.dev.yml",
      "type": "override"
    }
  }
}
```

#### Table Format
```
Section    | Key     | Value                  | Source        | Type
-----------|---------|------------------------|---------------|----------
variables  | api_url | http://localhost:3000  | .ddx.dev.yml  | override
variables  | author  | John Doe               | env:GIT_USER  | env-var
```

## User Feedback
*To be collected during implementation and testing*

## Notes
- This feature is crucial for debugging configuration issues
- Should be highly readable and informative
- Consider interactive mode for exploring large configurations
- May need performance optimization for very large configs

---
*Story is part of FEAT-003 (Configuration Management)*