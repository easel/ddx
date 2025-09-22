# US-037: Install MCP Server

## Story Overview

**Story ID**: US-037
**Feature**: FEAT-014 (MCP Server Management)  
**Priority**: P0  
**Points**: 5  
**Sprint**: 1  

## User Story

**As a** developer using Claude Code  
**I want to** install an MCP server with a single command  
**So that** I can quickly add new capabilities to my AI assistant without manual configuration  

## Acceptance Criteria

### AC1: Basic Installation
**Given** an MCP server exists in the registry
**When** I run `ddx mcp install github`
**Then** the system should:
- Load server definition from registry
- Prompt for required environment variables
- Build Claude CLI command with proper arguments
- Execute `claude mcp add` or `claude mcp add-json`
- Confirm successful installation

### AC2: Environment Variable Handling
**Given** the GitHub server requires a personal access token
**When** I'm prompted for the token
**Then** the input should be masked (shown as ****)
**And** the value should be validated before saving
**And** passed securely to Claude CLI with `-e` flag  

### AC3: Claude CLI Detection
**Given** I have Claude Code installed with CLI
**When** I install an MCP server
**Then** the system should automatically detect `claude` command availability
**And** use appropriate Claude CLI commands for installation
**And** handle CLI errors gracefully  

### AC4: Server Conflict Detection
**Given** an MCP server is already installed
**When** I attempt to install the same server
**Then** the system should:
- Check existing installation via `claude mcp list`
- Warn about duplicate installation
- Offer to update or skip installation
- Provide removal option via `claude mcp remove`

### AC5: Validation and Error Handling
**Given** I provide invalid environment values  
**When** the system validates the input  
**Then** I should see specific error messages  
**And** be prompted to correct the values  
**And** have the option to cancel installation  

### AC6: Non-Interactive Installation
**Given** I want to automate installation  
**When** I run `ddx mcp install github --env GITHUB_TOKEN=ghp_xxx`  
**Then** the server should install without prompts  
**And** validate all provided values  
**And** fail with clear errors if values are invalid  

## Definition of Done

- [ ] Install command implemented with all options
- [ ] Environment variable prompting with masking
- [ ] Configuration file generation and merging
- [ ] Backup and rollback functionality
- [ ] Unit tests for all installation paths
- [ ] Integration tests with mock Claude configs
- [ ] Documentation with examples
- [ ] Security validation for credentials

## Technical Notes

### Implementation Flow
1. Load server definition from registry
2. Verify `claude` CLI command availability
3. Check for existing server installation
4. Collect required environment variables
5. Validate all inputs
6. Build Claude CLI command with arguments
7. Execute `claude mcp add` or `claude mcp add-json`
8. Verify installation via `claude mcp list`
9. Provide next steps to user

### Claude CLI Commands
```bash
# Simple server without environment variables
claude mcp add github 'npx' -- '-y' '@modelcontextprotocol/server-github'

# Server with environment variables
claude mcp add-json github '{"type":"stdio","command":"npx","args":["-y","@modelcontextprotocol/server-github"],"env":{"GITHUB_PERSONAL_ACCESS_TOKEN":"ghp_xxx"}}'
```

## Test Scenarios

### Happy Path
1. Install new server with valid credentials
2. Install server with existing config file
3. Install multiple servers sequentially
4. Non-interactive install with all parameters

### Edge Cases
1. No Claude CLI found - show installation instructions
2. Claude CLI outdated - warn about MCP support
3. Server already installed - detected via `claude mcp list`
4. Server name conflicts - suggest removal first
5. Network timeout during installation - Claude CLI handles retries

### Error Cases
1. Invalid server name - suggest similar names
2. Missing required variables - list requirements
3. Invalid credentials - specific error messages
4. Claude CLI execution failed - show CLI output
5. Command not found - Claude CLI installation guide

## UX Considerations

### Interactive Flow
```bash
$ ddx mcp install github
🔧 Installing GitHub MCP Server...

ℹ️ This server requires the following configuration:
  - GitHub Personal Access Token (with repo permissions)
  
🔐 Enter your GitHub Personal Access Token: ****
✅ Token format validated

📍 Detected Claude CLI available
🔧 Executing: claude mcp add-json github {...}

📦 Installing server components...
✅ Server added to Claude Code successfully

🎉 GitHub MCP server installed!

🚀 Next steps:
  1. Restart Claude Code to load the new server
  2. Look for the GitHub icon in the MCP section
  3. Test with: "Show me my recent GitHub commits"
```

### Non-Interactive Flow
```bash
$ ddx mcp install github --env GITHUB_TOKEN=ghp_xxx --yes
✅ GitHub MCP server installed successfully
```

## Security Considerations

1. **Credential Masking**: Never display tokens in plaintext
2. **Validation**: Check token format before saving
3. **Secure Storage**: Use appropriate file permissions (0600)
4. **Audit Logging**: Log installations without sensitive data
5. **Backup Recovery**: Encrypted backups of configurations

## Related Stories

- US-036: List Available MCP Servers
- US-038: Configure MCP Server (future)
- US-039: Remove MCP Server (future)
- US-040: Update MCP Server (future)

## Change Log

| Date | Author | Changes |
|------|--------|----------|
| 2025-01-15 | System | Initial story creation |