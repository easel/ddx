# DDx Data Schema Design

> **Last Updated**: 2025-01-20
> **Status**: Active
> **Phase**: Design

## Overview

DDx uses a file-based data model with YAML configurations and JSON state management. This document defines the schemas for all data structures used in the DDx system.

## Configuration Schema

### Project Configuration (.ddx.yml)

The primary configuration file for DDx projects.

```yaml
# DDx Project Configuration Schema
version: "1.0"                   # DDx config version
name: string                      # Project name
description: string               # Project description

# Repository configuration
repository:
  url: string                     # Master repository URL
  branch: string                  # Branch name (default: main)
  subtree_path: string           # Local subtree path (default: .ddx)

# Resource selection
resources:
  templates: [string]            # Template names to include
  patterns: [string]             # Pattern names to include
  prompts: [string]              # Prompt names to include
  configs: [string]              # Config names to include
  scripts: [string]              # Script names to include

# Template variables
variables:
  project_name: string           # Used in templates
  author: string                 # Project author
  license: string                # License type
  custom: map[string]string      # Custom variables

# Persona bindings
personas:
  bindings:
    role_name: persona_name      # Map roles to personas
  overrides:
    workflow_name:
      role_name: persona_name    # Workflow-specific overrides

# MCP server configuration
mcp:
  servers:
    - name: string               # Server identifier
      version: string            # Version constraint
      config: map[string]any     # Server-specific config

# Package manager preference
package_manager: string         # npm|pnpm|yarn|bun|auto

# Workflow configuration
workflow:
  type: string                   # Workflow type (e.g., "helix")
  phase: string                  # Current phase
  features: [string]             # Active features
```

### Global Configuration (~/.ddx/config.yml)

User-level DDx configuration.

```yaml
# Global DDx Configuration Schema
version: "1.0"

# Default values for new projects
defaults:
  repository:
    url: string
    branch: string
  author: string
  license: string
  package_manager: string

# User preferences
preferences:
  editor: string                 # Preferred editor
  verbose: boolean               # Verbose output
  color: boolean                 # Colored output
  telemetry: boolean             # Usage telemetry (always false currently)

# Authentication tokens
auth:
  github_token: string           # GitHub personal access token
  gitlab_token: string           # GitLab personal access token

# Cache settings
cache:
  directory: string              # Cache directory path
  ttl_hours: integer             # Cache time-to-live
  max_size_mb: integer           # Maximum cache size
```

## State Management Schemas

### HELIX Workflow State (.helix-state.yml)

Tracks HELIX workflow progression.

```yaml
# HELIX Workflow State Schema
workflow: "helix"                # Workflow type identifier
current_phase: string            # frame|design|test|build|deploy|iterate
phases_completed: [string]       # List of completed phases

# Active features being developed
active_features:
  feature_id: string             # Feature status

# Timestamps
started_at: datetime             # Workflow start time
last_updated: datetime           # Last modification time

# Progress tracking
tasks_completed: [string]        # Completed task descriptions
next_actions: [string]           # Upcoming tasks
phase_progress:                  # Progress percentages
  phase_name: integer

# Feature-specific state
features:
  feature_id:
    status: string               # draft|in_progress|complete
    started_at: datetime
    completed_at: datetime
    artifacts: [string]          # Generated artifact paths
```

### Synchronization State (.ddx/sync-state.json)

Tracks synchronization with upstream repository.

```json
{
  "schema_version": "1.0",
  "last_sync": {
    "timestamp": "2025-01-20T10:00:00Z",
    "commit_sha": "abc123def456",
    "branch": "main",
    "status": "success|failed|conflict"
  },
  "local_changes": [
    {
      "file": "path/to/file",
      "action": "added|modified|deleted",
      "timestamp": "2025-01-20T10:00:00Z"
    }
  ],
  "conflicts": [
    {
      "file": "path/to/file",
      "type": "merge|content|deleted",
      "description": "Conflict description",
      "resolution": "pending|local|remote|manual"
    }
  ],
  "subtree": {
    "prefix": ".ddx",
    "remote": "ddx-upstream",
    "branch": "main",
    "squash": true
  }
}
```

### MCP Registry Cache (.ddx/cache/mcp-registry.json)

Caches MCP server registry for offline access.

```json
{
  "schema_version": "1.0",
  "updated_at": "2025-01-20T10:00:00Z",
  "source": "https://ddx-tools.github.io/mcp-registry/",
  "servers": [
    {
      "name": "filesystem",
      "package": "@modelcontextprotocol/server-filesystem",
      "version": "1.0.0",
      "description": "File system access for Claude",
      "author": "Anthropic",
      "tags": ["core", "filesystem"],
      "installation": {
        "npm": "@modelcontextprotocol/server-filesystem",
        "binary": null,
        "docker": null
      },
      "configuration": {
        "command": "npx",
        "args": ["@modelcontextprotocol/server-filesystem", "$PWD"],
        "env": {}
      },
      "requirements": {
        "node": ">=18.0.0",
        "platforms": ["darwin", "linux", "win32"]
      }
    }
  ]
}
```

## Resource Schemas

### Template Metadata (templates/*/metadata.yml)

Describes template properties and requirements.

```yaml
# Template Metadata Schema
name: string                     # Template identifier
version: "1.0.0"                 # Semantic version
description: string              # Template description
author: string                   # Template author
tags: [string]                   # Categorization tags

# Template requirements
requirements:
  ddx_version: string            # Minimum DDx version
  dependencies: [string]         # Required tools/packages

# Template variables
variables:
  - name: string                 # Variable name
    description: string          # Variable purpose
    default: string              # Default value
    required: boolean            # Is required?
    pattern: string              # Validation regex

# File operations
files:
  include: [string]              # Glob patterns to include
  exclude: [string]              # Glob patterns to exclude

# Post-application hooks
hooks:
  post_apply: [string]           # Commands to run after applying
```

### Persona Definition (personas/*.md)

AI persona specifications with YAML frontmatter.

```yaml
---
# Persona Frontmatter Schema
name: string                     # Unique persona identifier
roles: [string]                  # Roles this persona can fulfill
description: string              # Brief description
tags: [string]                   # Discovery tags
version: "1.0.0"                 # Persona version
author: string                   # Persona author
created: date                    # Creation date
updated: date                    # Last update date

# Persona capabilities
expertise:
  domains: [string]              # Knowledge domains
  tools: [string]                # Tool proficiencies
  languages: [string]            # Programming languages

# Behavioral traits
traits:
  communication_style: string    # formal|casual|technical
  response_length: string        # concise|detailed|adaptive
  proactivity: string            # low|medium|high
---

# Persona content in markdown...
```

## Claude Integration Schemas

### Claude Settings (.claude/settings.local.json)

Local Claude Code/Desktop configuration.

```json
{
  "mcpServers": {
    "server_name": {
      "command": "string",
      "args": ["array", "of", "strings"],
      "env": {
        "KEY": "value"
      },
      "disabled": false
    }
  },
  "defaultModelId": "string",
  "gitCommitMessageFormat": "string"
}
```

### CLAUDE.md Persona Section

Structured format for persona injection.

```markdown
<!-- PERSONAS:START -->
## Active Personas

### Role Name: persona-identifier
[Full persona content from persona file]

When responding, adopt the appropriate persona based on the task.
<!-- PERSONAS:END -->
```

## Data Validation Rules

### Configuration Validation
1. **Version compatibility**: Config version must be supported
2. **Required fields**: name, version are mandatory
3. **URL validation**: Repository URLs must be valid git URLs
4. **Path validation**: File paths must be valid for the OS
5. **Variable naming**: Must match pattern `^[a-zA-Z_][a-zA-Z0-9_]*$`

### State File Validation
1. **Phase progression**: Phases must follow HELIX order
2. **Timestamp format**: ISO 8601 format required
3. **Status values**: Must be from allowed enum values
4. **File integrity**: State files must be valid YAML/JSON

### Security Validation
1. **No secrets in config**: Tokens only in global config
2. **Path traversal prevention**: No `../` in paths
3. **Command injection prevention**: Sanitize all user inputs
4. **File permissions**: Respect system file permissions

## Schema Evolution

### Versioning Strategy
- **Major version**: Breaking changes to schema structure
- **Minor version**: New optional fields added
- **Patch version**: Documentation or validation updates

### Migration Support
- DDx CLI includes migration logic for schema upgrades
- Backward compatibility maintained for 2 major versions
- Automatic backups before schema migrations

### Deprecation Process
1. Mark field as deprecated in documentation
2. Add warning when deprecated field is used
3. Maintain support for 2 minor versions
4. Remove in next major version

## Related Documentation

- [[architecture]] - System architecture overview
- [[contracts/CLI-001-core-commands]] - CLI command contracts
- [[security-architecture]] - Security considerations
- [[implementation/config-management]] - Configuration implementation