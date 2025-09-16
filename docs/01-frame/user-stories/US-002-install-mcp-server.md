# US-002: Install MCP Server

## Story Overview

**Story ID**: US-002  
**Feature**: FEAT-001 (MCP Server Management)  
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
- Download server definition
- Prompt for required environment variables
- Generate Claude configuration
- Add to appropriate config file
- Confirm successful installation

### AC2: Environment Variable Handling
**Given** the GitHub server requires a personal access token  
**When** I'm prompted for the token  
**Then** the input should be masked (shown as ****)  
**And** the value should be validated before saving  
**And** stored securely in the configuration  

### AC3: Auto-Detection of Claude Type
**Given** I have Claude Code installed  
**When** I install an MCP server  
**Then** the system should automatically detect Claude Code  
**And** write to `~/.claude/settings.local.json`  
**And** use the correct configuration format  

### AC4: Configuration Backup
**Given** an existing Claude configuration file  
**When** I install a new MCP server  
**Then** the system should:
- Create a backup of the current config
- Merge new server configuration
- Preserve existing settings
- Provide rollback option if needed

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
2. Detect Claude installation type and path
3. Collect required environment variables
4. Validate all inputs
5. Backup existing configuration
6. Generate server configuration JSON
7. Merge with existing config
8. Write updated configuration
9. Verify installation
10. Provide next steps to user

### Configuration Structure
```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "<masked>"
      }
    }
  }
}
```

## Test Scenarios

### Happy Path
1. Install new server with valid credentials
2. Install server with existing config file
3. Install multiple servers sequentially
4. Non-interactive install with all parameters

### Edge Cases
1. No Claude installation found - show instructions
2. Multiple Claude installations - prompt for choice
3. Corrupted config file - offer repair options
4. Server already installed - prompt for update
5. Network timeout - retry with exponential backoff

### Error Cases
1. Invalid server name - suggest similar names
2. Missing required variables - list requirements
3. Invalid credentials - specific error messages
4. Permission denied - fix instructions
5. Disk full - check space before install

## UX Considerations

### Interactive Flow
```bash
$ ddx mcp install github
🔧 Installing GitHub MCP Server...

ℹ️ This server requires the following configuration:
  - GitHub Personal Access Token (with repo permissions)
  
🔐 Enter your GitHub Personal Access Token: ****
✅ Token format validated

📍 Detected Claude Code at: ~/.claude/
💾 Creating backup: ~/.claude/settings.local.json.backup

📦 Installing server components...
✅ Configuration written successfully

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

- US-001: List Available MCP Servers
- US-003: Configure MCP Server
- US-004: Remove MCP Server
- US-005: Update MCP Server

## Change Log

| Date | Author | Changes |
|------|--------|----------|
| 2025-01-15 | System | Initial story creation |