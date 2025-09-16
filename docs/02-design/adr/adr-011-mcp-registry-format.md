# ADR-011: MCP Server Registry Format

## Status
Proposed

## Context

The MCP server management feature requires a registry format to store and distribute server definitions. The registry must be:
- Human-readable and editable
- Machine-parseable with validation
- Versionable in Git
- Extensible for future needs
- Secure against malicious content

### Requirements
1. Store server metadata (name, description, author, version)
2. Define execution commands and arguments
3. Specify environment variables with validation rules
4. Support categorization and search
5. Enable offline usage with caching
6. Allow community contributions

### Constraints
- Must work with existing DDx YAML infrastructure
- Should minimize binary size impact
- Need to support cross-platform paths
- Must validate against security threats

## Decision

We will use **YAML format** for MCP server definitions with the following structure:

```yaml
# mcp-servers/servers/github.yml
name: github
description: GitHub integration for repository access
category: development
author: modelcontextprotocol
version: 1.0.0
tags: [git, repository, collaboration, vcs]

command:
  executable: npx
  args: ["-y", "@modelcontextprotocol/server-github"]
  
environment:
  - name: GITHUB_PERSONAL_ACCESS_TOKEN
    description: GitHub personal access token with repo permissions
    required: true
    sensitive: true
    validation: "^ghp_[a-zA-Z0-9]{36}$"
    
documentation:
  setup: "Generate token at https://github.com/settings/tokens"
  permissions: [repo, read:user]
  examples:
    - "Show my recent commits"
    - "Create a new issue in repo X"
    
compatibility:
  platforms: [darwin, linux, windows]
  claude_versions: [desktop, code]
  min_ddx_version: "0.2.0"
  
security:
  sandbox: recommended
  network_access: required
  file_access: none
```

With a registry index file:

```yaml
# mcp-servers/registry.yml
version: 1.0.0
updated: 2025-01-15T00:00:00Z

servers:
  - name: github
    file: servers/github.yml
    checksum: sha256:abc123...
    
  - name: postgres
    file: servers/postgres.yml
    checksum: sha256:def456...

categories:
  development:
    description: Development and version control tools
    servers: [github, gitlab, bitbucket]
    
  database:
    description: Database and data storage integrations
    servers: [postgres, mysql, mongodb, redis]
```

## Alternatives Considered

### Alternative 1: JSON Format

**Pros:**
- Native Go unmarshaling
- Faster parsing
- Stricter schema enforcement

**Cons:**
- Not human-friendly for editing
- No comments support
- Verbose for complex structures
- Poor Git diff experience

**Verdict:** Rejected - Developer experience is crucial for community contributions

### Alternative 2: TOML Format

**Pros:**
- Human-readable
- Supports comments
- Good for configuration

**Cons:**
- Less familiar to developers
- Limited nesting support
- Smaller ecosystem of tools
- Array handling is awkward

**Verdict:** Rejected - YAML is more familiar and flexible

### Alternative 3: Database (SQLite)

**Pros:**
- Fast queries
- ACID compliance
- Complex relationships

**Cons:**
- Not human-editable
- Binary format in Git
- Requires migration tools
- Overkill for our needs

**Verdict:** Rejected - Adds unnecessary complexity

### Alternative 4: Protocol Buffers

**Pros:**
- Type-safe
- Efficient serialization
- Cross-language support

**Cons:**
- Requires compilation step
- Not human-readable
- Complex toolchain
- Over-engineered for our use case

**Verdict:** Rejected - Too complex for simple registry

## Consequences

### Positive

1. **Developer-Friendly**: Easy to read, write, and contribute
2. **Git-Compatible**: Text-based, diffable, mergeable
3. **Flexible**: Supports complex structures and relationships
4. **Ecosystem**: Extensive tooling and library support
5. **Consistent**: Aligns with existing DDx configuration patterns
6. **Extensible**: Easy to add new fields without breaking compatibility
7. **Comments**: Can document directly in the file

### Negative

1. **Security Risk**: YAML parsing vulnerabilities (mitigated by validation)
2. **Performance**: Slower than binary formats (mitigated by caching)
3. **Size**: More verbose than binary formats (acceptable tradeoff)
4. **Type Safety**: Weaker than schema-based formats (mitigated by validation)
5. **Indentation**: Sensitive to whitespace errors (mitigated by linting)

### Neutral

1. **Learning Curve**: Most developers know YAML basics
2. **Tooling Requirements**: Need YAML parser library
3. **Validation Needs**: Must implement schema validation

## Implementation Notes

### Security Measures

1. **Schema Validation**: Strict validation against defined schema
2. **Size Limits**: Maximum file size of 100KB per server definition
3. **Sanitization**: Strip dangerous content before parsing
4. **Sandboxing**: Parse in restricted environment
5. **Checksum Verification**: Validate file integrity

### Caching Strategy

```go
type RegistryCache struct {
    Servers   map[string]*MCPServer
    Updated   time.Time
    TTL       time.Duration // 15 minutes default
    Checksum  string
}
```

### Validation Rules

```go
type RegistryValidator struct {
    MaxFileSize     int64  // 100KB
    MaxServers      int    // 1000
    AllowedExecutables []string // ["npx", "node", "python"]
    RequiredFields  []string // ["name", "description", "command"]
}
```

## Migration Path

1. **Phase 1**: Implement YAML parser with validation
2. **Phase 2**: Create initial registry with core servers
3. **Phase 3**: Add caching layer
4. **Phase 4**: Enable community contributions
5. **Phase 5**: Add registry signing/verification

## Success Metrics

1. **Parse Time**: <50ms for full registry
2. **Cache Hit Rate**: >90% for list operations
3. **Contribution Rate**: >5 new servers per month
4. **Security Incidents**: Zero from registry content
5. **User Satisfaction**: >4.5/5 for ease of use

## References

- [YAML 1.2 Specification](https://yaml.org/spec/1.2/spec.html)
- [YAML Security Best Practices](https://github.com/yaml/yaml-spec/wiki/Security)
- [Go YAML v3 Documentation](https://pkg.go.dev/gopkg.in/yaml.v3)
- [MCP Server Examples](https://github.com/modelcontextprotocol/servers)

## Decision Outcome

**Chosen Option:** YAML format with strict validation and caching

This decision provides the best balance of developer experience, security, and functionality for the MCP server registry. The human-readable format encourages community contributions while the validation layer ensures security.

## Review Schedule

This decision will be reviewed:
- After 6 months of production use
- If security incidents occur
- If performance becomes an issue
- If MCP protocol significantly changes

---

**Date:** 2025-01-15  
**Deciders:** DDx Architecture Team  
**Status:** Proposed → Accepted → Superseded by ADR-XXX (if applicable)