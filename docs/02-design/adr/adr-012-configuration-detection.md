# ADR-012: Claude Configuration Detection Strategy

## Status
Proposed

## Context

The MCP server management feature needs to automatically detect and configure Claude Code and Claude Desktop installations across different platforms. Each has different configuration file locations and formats:

### Claude Code
- **macOS/Linux**: `~/.claude/settings.local.json`
- **Windows**: `%USERPROFILE%\.claude\settings.local.json`
- **Format**: JSON with specific schema

### Claude Desktop
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Format**: JSON with different schema than Claude Code

### Requirements
1. Auto-detect installed Claude variants
2. Handle multiple installations gracefully
3. Support manual path override
4. Validate configuration formats
5. Preserve existing settings
6. Work across all platforms

### Constraints
- Cannot modify Claude application files
- Must respect user permissions
- Should not require admin/root access
- Must handle missing installations gracefully

## Decision

We will implement a **multi-strategy detection system** with the following approach:

### Detection Priority Order

1. **Environment Variables** (highest priority)
   ```bash
   CLAUDE_CODE_CONFIG=/custom/path/settings.local.json
   CLAUDE_DESKTOP_CONFIG=/custom/path/claude_desktop_config.json
   ```

2. **Standard Paths** (automatic detection)
   ```go
   type ClaudeLocation struct {
       Type       ClaudeType // Code or Desktop
       ConfigPath string
       Platform   string
       Priority   int
   }
   
   var StandardLocations = []ClaudeLocation{
       // Claude Code
       {Code, "~/.claude/settings.local.json", "darwin", 1},
       {Code, "~/.claude/settings.local.json", "linux", 1},
       {Code, "%USERPROFILE%\.claude\settings.local.json", "windows", 1},
       
       // Claude Desktop
       {Desktop, "~/Library/Application Support/Claude/claude_desktop_config.json", "darwin", 2},
       {Desktop, "~/.config/Claude/claude_desktop_config.json", "linux", 2},
       {Desktop, "%APPDATA%\Claude\claude_desktop_config.json", "windows", 2},
   }
   ```

3. **User Configuration** (`.ddx.yml`)
   ```yaml
   mcp:
     claude_code_path: /custom/path/to/claude/code/config.json
     claude_desktop_path: /custom/path/to/claude/desktop/config.json
     preferred: code  # or "desktop" or "both"
   ```

4. **Interactive Selection** (if multiple found)
   ```
   Multiple Claude installations detected:
   1. Claude Code at ~/.claude/settings.local.json
   2. Claude Desktop at ~/Library/Application Support/Claude/...
   
   Which would you like to configure? (1/2/both): 
   ```

### Detection Algorithm

```go
func DetectClaude() (*ClaudeConfig, error) {
    // 1. Check environment variables
    if path := os.Getenv("CLAUDE_CODE_CONFIG"); path != "" {
        return detectFromPath(path, ClaudeCode)
    }
    
    // 2. Check DDx configuration
    if config.Has("mcp.claude_code_path") {
        return detectFromPath(config.GetString("mcp.claude_code_path"), ClaudeCode)
    }
    
    // 3. Check standard locations
    installations := []ClaudeConfig{}
    for _, loc := range StandardLocations {
        if runtime.GOOS != loc.Platform {
            continue
        }
        path := expandPath(loc.ConfigPath)
        if exists(path) {
            installations = append(installations, detectFromPath(path, loc.Type))
        }
    }
    
    // 4. Handle results
    switch len(installations) {
    case 0:
        return nil, ErrNoClaudeFound
    case 1:
        return &installations[0], nil
    default:
        return selectInstallation(installations)
    }
}
```

### Version Detection

```go
type ClaudeVersion struct {
    Type    ClaudeType
    Version string
    Schema  string
}

func detectVersion(configPath string) (*ClaudeVersion, error) {
    content, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }
    
    // Check for version markers
    if strings.Contains(string(content), "mcpServers") {
        // Both types use mcpServers now
        if strings.Contains(configPath, "desktop") {
            return &ClaudeVersion{Desktop, "1.0+", "desktop-v1"}, nil
        }
        return &ClaudeVersion{Code, "0.2+", "code-v1"}, nil
    }
    
    return nil, ErrUnknownFormat
}
```

## Alternatives Considered

### Alternative 1: Registry-Based Detection (Windows)

**Pros:**
- Reliable on Windows
- Can detect installation paths
- Version information available

**Cons:**
- Windows-only solution
- Requires registry access
- May need elevated permissions
- Not portable to other platforms

**Verdict:** Rejected - Not cross-platform

### Alternative 2: Process Detection

**Pros:**
- Can detect running instances
- Direct validation possible
- Real-time information

**Cons:**
- Privacy concerns
- Requires process listing permissions
- Claude might not be running
- Complex implementation

**Verdict:** Rejected - Too invasive and unreliable

### Alternative 3: File System Scanning

**Pros:**
- Comprehensive detection
- Finds non-standard installations
- No user input needed

**Cons:**
- Very slow performance
- Privacy concerns
- May trigger antivirus
- Resource intensive

**Verdict:** Rejected - Performance and privacy issues

### Alternative 4: Manual Configuration Only

**Pros:**
- Simple implementation
- No detection logic needed
- User has full control

**Cons:**
- Poor user experience
- Error-prone
- Requires documentation
- Higher support burden

**Verdict:** Rejected - Poor UX for common case

## Consequences

### Positive

1. **Automatic Discovery**: Works out-of-the-box for most users
2. **Flexible Override**: Power users can customize paths
3. **Cross-Platform**: Single approach works everywhere
4. **Fast Detection**: Check specific paths, not scanning
5. **Privacy-Friendly**: Only checks expected locations
6. **Graceful Fallback**: Multiple strategies ensure success
7. **Version Aware**: Can adapt to format changes

### Negative

1. **Maintenance Burden**: Must track Claude path changes
2. **Platform Complexity**: Different paths per OS
3. **Version Compatibility**: Must handle format evolution
4. **Edge Cases**: Non-standard installations need manual config

### Neutral

1. **User Prompts**: May need interaction for multiple installations
2. **Configuration Files**: Adds MCP section to .ddx.yml
3. **Error Messages**: Must be clear about detection failures

## Implementation Notes

### Platform-Specific Path Expansion

```go
func expandPath(path string) string {
    if strings.HasPrefix(path, "~") {
        home, _ := os.UserHomeDir()
        path = filepath.Join(home, path[2:])
    }
    return os.ExpandEnv(path)
}
```

### Configuration Validation

```go
type ConfigValidator interface {
    ValidateFormat(content []byte) error
    ValidateSchema(config map[string]interface{}) error
    CanMerge(existing, new map[string]interface{}) bool
}

type ClaudeCodeValidator struct{}
type ClaudeDesktopValidator struct{}
```

### Error Handling

```go
var (
    ErrNoClaudeFound = errors.New("no Claude installation detected")
    ErrInvalidConfig = errors.New("invalid Claude configuration format")
    ErrPermissionDenied = errors.New("cannot access Claude configuration")
    ErrMultipleInstalls = errors.New("multiple Claude installations found")
)
```

### Caching Strategy

```go
type DetectionCache struct {
    Installations []ClaudeConfig
    DetectedAt    time.Time
    TTL           time.Duration // 5 minutes
}
```

## Migration Path

1. **Phase 1**: Implement basic path detection
2. **Phase 2**: Add version detection
3. **Phase 3**: Add user preference handling
4. **Phase 4**: Add caching layer
5. **Phase 5**: Add installation validation

## Success Metrics

1. **Detection Rate**: >95% automatic success
2. **Detection Speed**: <100ms for detection
3. **Error Rate**: <1% false positives
4. **Platform Coverage**: 100% of supported OS
5. **User Satisfaction**: <5% manual configuration needed

## Testing Strategy

### Test Scenarios

1. **No Installation**: Appropriate error message
2. **Single Installation**: Automatic detection
3. **Multiple Installations**: User prompt
4. **Custom Path**: Environment variable override
5. **Missing Permissions**: Graceful failure
6. **Corrupted Config**: Validation error

### Platform Matrix

| Platform | Claude Code | Claude Desktop | Both |
|----------|------------|----------------|------|
| macOS | ✓ | ✓ | ✓ |
| Linux | ✓ | ✓ | ✓ |
| Windows | ✓ | ✓ | ✓ |

## References

- [Claude Code Documentation](https://docs.anthropic.com/claude-code)
- [Claude Desktop Documentation](https://docs.anthropic.com/claude-desktop)
- [Go filepath Package](https://pkg.go.dev/path/filepath)
- [XDG Base Directory Spec](https://specifications.freedesktop.org/basedir-spec/)

## Decision Outcome

**Chosen Option:** Multi-strategy detection with priority order

This approach provides the best balance of automatic detection, user control, and cross-platform compatibility. The fallback chain ensures success in most scenarios while respecting user preferences.

## Review Schedule

This decision will be reviewed:
- When Claude changes configuration paths
- If detection success rate falls below 90%
- After major Claude version updates
- Based on user feedback

---

**Date:** 2025-01-15  
**Deciders:** DDx Architecture Team  
**Status:** Proposed → Accepted → Superseded by ADR-XXX (if applicable)