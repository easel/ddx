# CLI-001: MCP Commands Contract

## Contract Overview

**Contract ID**: CLI-001  
**Component**: MCP Management CLI  
**Version**: 1.0.0  
**Status**: In Design  

## Command Structure

The MCP management functionality is exposed through the `ddx mcp` command with multiple subcommands.

```
ddx mcp <subcommand> [options] [arguments]
```

## Commands

### 1. List Command

#### Signature
```bash
ddx mcp list [options]
```

#### Purpose
Display available MCP servers from the registry with filtering and search capabilities.

#### Options
| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--category` | `-c` | string | all | Filter by category (development, database, filesystem, productivity) |
| `--search` | `-s` | string | - | Search term for name/description |
| `--installed` | `-i` | bool | false | Show only installed servers |
| `--available` | `-a` | bool | false | Show only not-installed servers |
| `--verbose` | `-v` | bool | false | Show detailed information |
| `--format` | `-f` | string | table | Output format (table, json, yaml) |
| `--no-cache` | - | bool | false | Force registry refresh |

#### Output Format (Table)
```
📋 Available MCP Servers (15 total, 3 installed)

Development:
  ✅ github         - GitHub integration for repository access
  ⬜ gitlab         - GitLab integration for repository access
  
Database:
  ✅ postgres       - PostgreSQL database integration
  ⬜ mysql          - MySQL database integration
```

#### Output Format (JSON)
```json
{
  "total": 15,
  "installed": 3,
  "servers": [
    {
      "name": "github",
      "description": "GitHub integration for repository access",
      "category": "development",
      "installed": true,
      "version": "1.0.0"
    }
  ]
}
```

#### Exit Codes
- `0`: Success
- `1`: Registry loading error
- `2`: Invalid options

#### Examples
```bash
# List all servers
ddx mcp list

# List development servers
ddx mcp list --category development

# Search for git-related servers
ddx mcp list --search git

# Show only installed servers with details
ddx mcp list --installed --verbose

# Output as JSON
ddx mcp list --format json
```

### 2. Install Command

#### Signature
```bash
ddx mcp install <server-name> [options]
```

#### Purpose
Install and configure an MCP server for Claude Code/Desktop.

#### Arguments
- `server-name` (required): Name of the server to install

#### Options
| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--env` | `-e` | string[] | - | Environment variables (KEY=VALUE) |
| `--yes` | `-y` | bool | false | Skip confirmation prompts |
| `--claude-type` | `-t` | string | auto | Target Claude type (code, desktop, both) |
| `--config-path` | `-p` | string | auto | Custom config file path |
| `--no-backup` | - | bool | false | Skip configuration backup |
| `--dry-run` | - | bool | false | Show what would be done |

#### Interactive Flow
```
🔧 Installing GitHub MCP Server...

ℹ️ This server requires:
  - GitHub Personal Access Token (with repo permissions)
  
🔐 Enter GitHub Personal Access Token: ****
✅ Token validated

📍 Detected: Claude Code at ~/.claude/
💾 Backup created: ~/.claude/settings.local.json.backup

📦 Configuring server...
✅ GitHub MCP server installed successfully!

🚀 Next steps:
  1. Restart Claude Code
  2. Look for GitHub in MCP section
  3. Test with: "Show my recent commits"
```

#### Non-Interactive Mode
```bash
ddx mcp install github --env GITHUB_TOKEN=ghp_xxx --yes
```

#### Exit Codes
- `0`: Success
- `1`: Server not found
- `2`: Invalid environment variables
- `3`: Configuration write error
- `4`: Claude not detected
- `5`: Validation error

#### Examples
```bash
# Interactive installation
ddx mcp install github

# Non-interactive with environment
ddx mcp install postgres --env DATABASE_URL=postgresql://... --yes

# Install for specific Claude type
ddx mcp install github --claude-type desktop

# Dry run to preview changes
ddx mcp install github --dry-run
```

### 3. Configure Command

#### Signature
```bash
ddx mcp configure <server-name> [options]
```

#### Purpose
Update configuration for an installed MCP server.

#### Arguments
- `server-name` (required): Name of installed server

#### Options
| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--env` | `-e` | string[] | - | New environment variables |
| `--add-env` | - | string[] | - | Add environment variables |
| `--remove-env` | - | string[] | - | Remove environment variables |
| `--reset` | - | bool | false | Reset to defaults |

#### Exit Codes
- `0`: Success
- `1`: Server not installed
- `2`: Invalid configuration
- `3`: Write error

#### Examples
```bash
# Update token
ddx mcp configure github --env GITHUB_TOKEN=new_token

# Add new environment variable
ddx mcp configure postgres --add-env POOL_SIZE=10

# Reset to defaults
ddx mcp configure github --reset
```

### 4. Remove Command

#### Signature
```bash
ddx mcp remove <server-name> [options]
```

#### Purpose
Remove an installed MCP server configuration.

#### Arguments
- `server-name` (required): Name of server to remove

#### Options
| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--yes` | `-y` | bool | false | Skip confirmation |
| `--no-backup` | - | bool | false | Skip backup creation |
| `--purge` | - | bool | false | Remove all related data |

#### Exit Codes
- `0`: Success
- `1`: Server not installed
- `2`: Removal error

#### Examples
```bash
# Remove with confirmation
ddx mcp remove github

# Remove without confirmation
ddx mcp remove github --yes

# Remove and purge all data
ddx mcp remove github --purge
```

### 5. Status Command

#### Signature
```bash
ddx mcp status [server-name] [options]
```

#### Purpose
Show status of installed MCP servers.

#### Arguments
- `server-name` (optional): Specific server to check

#### Options
| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--check` | `-c` | bool | false | Verify server connectivity |
| `--verbose` | `-v` | bool | false | Show detailed information |
| `--format` | `-f` | string | table | Output format |

#### Output Format
```
📊 MCP Server Status

✅ github (v1.0.0)
   Status: Configured
   Claude: Code (~/.claude/settings.local.json)
   Environment: GITHUB_TOKEN=***

✅ postgres (v2.1.0)
   Status: Configured
   Claude: Code (~/.claude/settings.local.json)
   Environment: DATABASE_URL=***

⚠️ redis (v1.2.0)
   Status: Configuration Error
   Issue: Missing required environment variable REDIS_URL
```

#### Exit Codes
- `0`: All servers healthy
- `1`: Some servers have issues
- `2`: Check failed

#### Examples
```bash
# Check all servers
ddx mcp status

# Check specific server
ddx mcp status github

# Verify connectivity
ddx mcp status --check

# Detailed status as JSON
ddx mcp status --verbose --format json
```

### 6. Update Command

#### Signature
```bash
ddx mcp update [options]
```

#### Purpose
Update MCP server registry and definitions.

#### Options
| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--force` | `-f` | bool | false | Force update even if current |
| `--server` | `-s` | string | - | Update specific server |
| `--check` | `-c` | bool | false | Check for updates only |

#### Output
```
🔄 Updating MCP Registry...

📥 Fetching latest registry from https://github.com/easel/ddx
✅ Registry updated (15 new servers, 3 updates)

New servers:
  - anthropic-docs: Anthropic documentation access
  - slack: Slack workspace integration
  
Updated servers:
  - github: 1.0.0 → 1.1.0
  - postgres: 2.0.0 → 2.1.0
```

#### Exit Codes
- `0`: Success
- `1`: Update failed
- `2`: Network error

#### Examples
```bash
# Update registry
ddx mcp update

# Check for updates
ddx mcp update --check

# Force update
ddx mcp update --force
```

## Global Options

These options apply to all MCP commands:

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--config` | - | string | .ddx.yml | DDx configuration file |
| `--debug` | - | bool | false | Enable debug output |
| `--quiet` | `-q` | bool | false | Suppress non-error output |
| `--help` | `-h` | bool | false | Show help message |

## Error Handling

### Error Response Format
```json
{
  "error": {
    "code": "MCP_SERVER_NOT_FOUND",
    "message": "MCP server 'invalid' not found in registry",
    "details": {
      "searched": "invalid",
      "suggestions": ["github", "gitlab"]
    },
    "resolution": "Run 'ddx mcp list' to see available servers"
  }
}
```

### Common Error Codes

| Code | Description | Resolution |
|------|-------------|------------|
| `MCP_SERVER_NOT_FOUND` | Server doesn't exist | Check server name |
| `MCP_ALREADY_INSTALLED` | Server already configured | Use configure command |
| `MCP_CLAUDE_NOT_FOUND` | Claude not detected | Install Claude or specify path |
| `MCP_INVALID_ENV` | Invalid environment variable | Check format and requirements |
| `MCP_CONFIG_CORRUPT` | Configuration file corrupted | Restore from backup |
| `MCP_PERMISSION_DENIED` | Cannot write config | Check file permissions |
| `MCP_NETWORK_ERROR` | Registry unreachable | Check network connection |
| `MCP_VALIDATION_ERROR` | Input validation failed | Review requirements |

## Environment Variables

The MCP commands respect these environment variables:

| Variable | Description | Example |
|----------|-------------|----------|
| `CLAUDE_CODE_CONFIG` | Override Claude Code config path | `/custom/path/config.json` |
| `CLAUDE_DESKTOP_CONFIG` | Override Claude Desktop config path | `/custom/path/desktop.json` |
| `DDX_MCP_REGISTRY` | Custom registry URL | `https://custom.registry.com` |
| `DDX_MCP_CACHE_DIR` | Cache directory | `/tmp/ddx-mcp-cache` |
| `DDX_MCP_NO_COLOR` | Disable colored output | `1` or `true` |

## Completion Support

The MCP commands support shell completion:

```bash
# Bash
ddx mcp completion bash >> ~/.bashrc

# Zsh
ddx mcp completion zsh >> ~/.zshrc

# Fish
ddx mcp completion fish > ~/.config/fish/completions/ddx-mcp.fish
```

## Version Compatibility

| DDx Version | MCP Protocol | Claude Code | Claude Desktop |
|-------------|--------------|-------------|----------------|
| 0.2.0+ | 1.0 | 0.2+ | 1.0+ |
| 0.3.0+ | 1.1 | 0.3+ | 1.1+ |

## Security Considerations

1. **Credential Masking**: All sensitive values shown as `***`
2. **Secure Input**: Password inputs use terminal raw mode
3. **File Permissions**: Config files set to 0600 (owner only)
4. **Validation**: All inputs validated before use
5. **Audit Logging**: Operations logged without sensitive data

## Testing

### Contract Tests
```go
// Test command structure
func TestMCPCommandStructure(t *testing.T) {
    cmd := NewMCPCommand()
    assert.Equal(t, "mcp", cmd.Use)
    assert.True(t, cmd.HasSubCommands())
}

// Test list output format
func TestListCommandOutput(t *testing.T) {
    output := runCommand("mcp", "list", "--format", "json")
    var result ListResult
    json.Unmarshal(output, &result)
    assert.NotNil(t, result.Servers)
}
```

## Change Log

| Version | Date | Changes |
|---------|------|----------|
| 1.0.0 | 2025-01-15 | Initial MCP command specification |

---

*This contract defines the complete CLI interface for MCP server management in DDx.*