---
tags: [adr, architecture, configuration, yaml, ddx, schema]
template: false
version: 1.0.0
---

# ADR-005: Configuration Management Architecture

**Date**: 2025-01-14
**Updated**: 2025-01-24
**Status**: Accepted
**Deciders**: DDX Development Team
**Technical Story**: Define the configuration management architecture and file structure for DDX project configuration

## Context

### Problem Statement
DDX needs a robust configuration management system that:
- Defines project-specific DDX settings and resource selections
- Supports environment-specific configurations
- Enables configuration inheritance and composition
- Provides clear validation and error reporting
- Maintains backward compatibility across versions
- Supports both simple and advanced use cases

### Forces at Play
- **Simplicity**: New users need minimal configuration to get started
- **Flexibility**: Advanced users need powerful customization options
- **Validation**: Configuration errors should be caught early with clear messages
- **Portability**: Configurations should work across different environments
- **Extensibility**: Must support future features without breaking changes
- **Documentation**: Configuration options must be self-documenting
- **Version Control**: Configurations must be diff-friendly and mergeable

### Constraints
- Must use YAML format (established in project conventions)
- Cannot break existing .ddx.yml files in production
- Must validate quickly (< 100ms)
- Should support JSON Schema for tooling integration
- Must handle missing/partial configurations gracefully
- Need to support environment variable expansion

## Decision

### Chosen Approach
Implement a layered configuration system with:
1. **Schema-based validation** using JSON Schema
2. **Multi-level configuration** with inheritance
3. **Environment-specific overrides**
4. **Smart defaults** with progressive disclosure
5. **Semantic versioning** for configuration compatibility

### Configuration File Location
**Decision**: Move configuration from `.ddx.yml` to `.ddx/config.yaml`

**Rationale**:
- Organizes all DDX files in a dedicated directory
- Separates configuration from project root
- Provides clear structure for library, cache, and state files
- Uses standard `.yaml` extension

### Directory Structure
```
project/
├── .ddx/
│   ├── config.yaml         # DDX configuration (moved from .ddx.yml)
│   ├── library/            # Git subtree synced with upstream
│   ├── cache/              # Local cache files
│   ├── state/              # Workflow state
│   └── local/              # User customizations
└── src/                    # Project code
```

### Configuration Hierarchy
```
1. Built-in defaults (lowest priority)
2. Global user config (~/.ddx/config.yaml)
3. Project config (.ddx/config.yaml)
4. Environment overrides (.ddx/config.{env}.yaml)
5. Environment variables (DDX_*)
6. Command-line flags (highest priority)
```

### Library Management Architecture
**Decision**: Use git subtree for library synchronization to `.ddx/library/` subfolder only

**Rationale**:
- Separates synced library content from local configuration and state
- Allows local customizations without sync conflicts
- Maintains git history for library changes
- Provides clear boundary between upstream and local content

### Library Path Resolution
**Decision**: Make `library_base_path` a top-level configuration property with explicit default

**Rationale**:
- Library location is architectural, not just a variable
- Default value (`./library`) should be explicit and documented
- Provides clear override mechanism for advanced users
- Simplifies configuration for typical use cases

### Simplified Configuration Schema
**Decision**: Simplify configuration schema, removing complex resource selection

```yaml
# .ddx/config.yaml schema version 1.0
version: "1.0"                      # Required: Configuration schema version

# Library configuration
library:
  path: "./library"                 # Path to DDX library relative to config.yaml (optional, default)
  repository:
    url: string                     # DDX repository URL (default: github.com/easel/ddx)
    branch: string                  # Branch to sync (default: main)
    subtree: string                 # What to sync from upstream (default: library)
```

**Rationale for Simplification**:
- Most users don't need complex resource selection
- Git subtree handles inclusion/exclusion at the repository level
- Simpler schema is easier to understand and validate
- Advanced users can still customize via repository settings
- Reduces cognitive load for new users
- Library configuration is grouped logically under single namespace

### Migration Path
**Decision**: Provide backwards compatibility and automatic migration

**Strategy**:
1. Auto-detect `.ddx.yml` and migrate to `.ddx/config.yaml`
2. Support both formats during transition period
3. Show deprecation warnings for old format
4. Remove legacy support in v2.0

## Consequences

### Positive
- **Better Organization**: All DDX files in dedicated `.ddx/` directory
- **Clearer Separation**: Library, config, cache, and state are distinct
- **Simpler Configuration**: Fewer required fields, sensible defaults
- **Git Subtree Isolation**: Only library syncs, preserving local customizations
- **Future-Proof**: Room for additional DDX metadata and state

### Negative
- **Breaking Change**: Requires migration for existing projects
- **Learning Curve**: Users must understand new directory structure
- **Migration Complexity**: Must handle edge cases during transition

### Neutral
- **File Format**: Remains YAML, no syntax changes
- **Validation**: Still uses JSON Schema (implementation detail)
- **Git Integration**: Still uses git subtree (different target path)

## Alternatives Considered

### Option 1: TOML Configuration
**Description**: Use TOML format for configuration files

**Pros**:
- Simpler syntax than YAML
- Less ambiguous parsing
- Better date/time support
- No indentation issues

**Cons**:
- Less familiar to developers
- Poor support for complex nested structures
- Limited tooling compared to YAML
- Not as human-readable for deep nesting

**Why rejected**: YAML's widespread adoption and better support for complex structures outweigh TOML's simplicity benefits

### Option 2: JSON Configuration
**Description**: Use JSON for configuration files

**Pros**:
- Universal parsing support
- No ambiguity in syntax
- Direct schema validation
- Excellent tooling

**Cons**:
- Not human-friendly (no comments, verbose)
- Poor git diff readability
- Prone to syntax errors (trailing commas)
- No multiline strings

**Why rejected**: JSON's lack of comments and poor human readability make it unsuitable for configuration files

### Option 3: HCL (HashiCorp Configuration Language)
**Description**: Use HCL for configuration

**Pros**:
- Designed for configuration
- Good interpolation support
- Clean syntax
- Terraform familiarity

**Cons**:
- Less familiar outside HashiCorp ecosystem
- Requires custom parser
- Limited tooling support
- Smaller community

**Why rejected**: Limited adoption and tooling support compared to YAML

### Option 4: INI Files
**Description**: Use INI format for simple configuration

**Pros**:
- Very simple syntax
- Wide support
- Easy to parse
- Familiar format

**Cons**:
- No support for nested structures
- No arrays/lists
- Limited data types
- Too simple for complex configuration

**Why rejected**: Insufficient expressiveness for DDX's configuration needs

### Option 5: Code-based Configuration (Go/Python)
**Description**: Use actual code files for configuration

**Pros**:
- Full programming language power
- Type safety (with Go)
- Complex logic possible
- IDE support

**Cons**:
- Security risks (code execution)
- Not declarative
- Harder to validate
- Platform-specific
- Not diff-friendly

**Why rejected**: Security concerns and loss of declarative simplicity

## Consequences

### Positive Consequences
- **Progressive Complexity**: Simple projects need minimal config, complex projects have full control
- **Clear Validation**: Schema validation provides immediate feedback on errors
- **Tool Support**: JSON Schema enables IDE autocomplete and validation
- **Environment Flexibility**: Easy to manage dev/staging/prod configurations
- **Backward Compatibility**: Versioned schema allows graceful evolution
- **Self-Documenting**: Schema serves as configuration documentation

### Negative Consequences
- **YAML Complexity**: YAML parsing edge cases can surprise users
- **Schema Maintenance**: Must maintain schema definitions and validators
- **Migration Burden**: Schema changes require migration strategies
- **Learning Curve**: Advanced features add cognitive load
- **Validation Overhead**: Schema validation adds startup time

### Neutral Consequences
- **File Proliferation**: Multiple config files for different environments
- **Override Complexity**: Understanding precedence requires documentation
- **Tooling Dependency**: Relies on YAML/JSON Schema tooling

## Implementation

### Required Changes
1. Implement JSON Schema validator in Go
2. Create schema definitions for each version
3. Build configuration loading hierarchy
4. Implement environment variable expansion
5. Add configuration migration utilities
6. Create validation error reporting
7. Build configuration debugging tools
8. Document all configuration options

### Migration Strategy
For existing .ddx.yml files:
1. Detect schema version (assume 0.x if missing)
2. Run migration transformers sequentially
3. Validate against target schema
4. Backup original if changes made
5. Report migration actions to user

### Success Metrics
- **Validation Speed**: < 50ms for typical config
- **Error Clarity**: 95% of errors understood without documentation
- **Migration Success**: 100% of v0.x configs migrate automatically
- **Schema Coverage**: 100% of options documented in schema
- **User Satisfaction**: > 80% find configuration intuitive

## Compliance

### Security Requirements
- No code execution in configuration
- Validate all external references
- Sanitize environment variable expansion
- Prevent path traversal in file references
- Secure handling of sensitive variables

### Performance Requirements
- Schema validation < 50ms
- Configuration loading < 100ms total
- Efficient caching of validated configs
- Lazy loading of optional sections

### Regulatory Requirements
- GDPR compliance for any PII in config
- No secrets in committed configurations
- Audit trail for configuration changes

## Monitoring and Review

### Key Indicators to Watch
- Configuration validation error rates
- Migration success rates
- Time spent debugging config issues
- Feature adoption rates
- Schema version distribution

### Review Date
Q2 2025 - After initial production usage

### Review Triggers
- Major YAML specification changes
- > 10% validation error rate
- User feedback on complexity
- Performance degradation
- Security vulnerability discovered

## Related Decisions

### Dependencies
- ADR-001: Three-layer architecture defines resource types
- ADR-002: Git subtree defines repository configuration
- ADR-003: Go implementation determines parser choice
- ADR-004: Starlark integration requires extension config

### Influenced By
- YAML adoption in DevOps tools
- JSON Schema standardization
- Kubernetes configuration patterns
- Docker Compose configuration design

### Influences
- Variable substitution design (future ADR)
- Environment management patterns
- CI/CD integration approach
- Plugin configuration format

## References

### Documentation
- [YAML Specification](https://yaml.org/spec/1.2/spec.html)
- [JSON Schema](https://json-schema.org/)
- [Viper Configuration](https://github.com/spf13/viper)
- [12-Factor App Config](https://12factor.net/config)

### External Resources
- [YAML Best Practices](https://www.yaml.org/spec/1.2/best-practices.html)
- [Configuration Complexity](https://research.swtch.com/deps)
- [Schema Evolution Patterns](https://martinfowler.com/articles/schemaEvolution.html)

### Discussion History
- Initial configuration format discussion
- Schema versioning strategy review
- Environment override design session
- User feedback on configuration complexity

## Notes

The configuration system follows DDX's medical metaphor - like a patient's medical record that captures all relevant information in a structured, validated format. The schema acts as the form template ensuring completeness and correctness.

Key insight: By providing smart defaults and progressive disclosure, we can make simple things simple while keeping complex things possible. The schema versioning ensures we can evolve without breaking existing users.

Implementation tip: Start with a minimal required configuration and expand based on actual user needs. Resist the temptation to add configuration options "just in case."

The layered override system mirrors medical treatment protocols - general guidelines (defaults) are specialized for specific conditions (project config) and further adjusted for individual patients (environment overrides).

---

**Last Updated**: 2025-01-14
**Next Review**: 2025-04-14