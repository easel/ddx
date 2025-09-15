---
tags: [adr, architecture, configuration, yaml, ddx, schema]
template: false
version: 1.0.0
---

# ADR-005: Configuration Management and .ddx.yml Schema

**Date**: 2025-01-14
**Status**: Proposed
**Deciders**: DDX Development Team
**Technical Story**: Define the configuration management strategy and schema for DDX project configuration

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

### Configuration Hierarchy
```
1. Built-in defaults (lowest priority)
2. Global user config (~/.ddx/config.yml)
3. Project config (.ddx.yml)
4. Environment overrides (.ddx.{env}.yml)
5. Environment variables (DDX_*)
6. Command-line flags (highest priority)
```

### Schema Definition
```yaml
# .ddx.yml schema version 1.0
version: "1.0"  # Configuration schema version

# Project metadata
project:
  name: string  # Optional project name
  description: string  # Optional description
  type: string  # Project type (web, cli, library, etc.)
  languages: [string]  # Primary languages used

# Repository configuration
repository:
  url: string  # DDX repository URL
  branch: string  # Branch to track (default: main)
  path: string  # Path within repository (default: /)

# Resource selection
resources:
  templates:
    include: [string]  # Templates to include
    exclude: [string]  # Templates to exclude
  patterns:
    include: [string]  # Patterns to include
    exclude: [string]  # Patterns to exclude
  prompts:
    include: [string]  # Prompts to include
    exclude: [string]  # Prompts to exclude
  configs:
    include: [string]  # Configs to include
    exclude: [string]  # Configs to exclude

# Variable definitions for substitution
variables:
  # User-defined variables
  key: value
  nested:
    key: value
  # Special variables
  $env: {}  # Environment variable passthrough
  $git: {}  # Git metadata injection

# Workflow configuration
workflows:
  enabled: [string]  # Enabled workflows
  disabled: [string]  # Explicitly disabled workflows
  custom: {}  # Custom workflow definitions

# Validator configuration
validators:
  enabled: boolean  # Global validator toggle
  rules: [string]  # Active validation rules
  custom: {}  # Custom validator definitions
  strict: boolean  # Fail on validation warnings

# Extension configuration
extensions:
  starlark:
    enabled: boolean
    modules: [string]  # Starlark modules to load
    timeout: number  # Execution timeout in ms
    memory_limit: string  # Memory limit (e.g., "100MB")

# Update configuration
updates:
  check: boolean  # Check for updates
  auto_pull: boolean  # Automatically pull updates
  frequency: string  # Check frequency (daily, weekly, monthly)
  channel: string  # Update channel (stable, beta, edge)

# Contribution configuration
contribution:
  enabled: boolean  # Allow contributions
  branch_prefix: string  # Prefix for contribution branches
  sign_commits: boolean  # GPG sign commits

# Local overrides
local:
  paths:
    templates: string  # Local templates directory
    patterns: string  # Local patterns directory
    prompts: string  # Local prompts directory
  ignore: [string]  # Patterns to ignore

# Hooks configuration
hooks:
  pre_apply: [string]  # Commands before apply
  post_apply: [string]  # Commands after apply
  pre_update: [string]  # Commands before update
  post_update: [string]  # Commands after update
```

### Rationale
- **YAML Format**: Human-readable, widely supported, git-friendly
- **Schema Validation**: Catches errors early, enables IDE support
- **Layered Configuration**: Supports multiple use cases from simple to complex
- **Environment Overrides**: Enables CI/CD and multi-environment workflows
- **Semantic Versioning**: Allows evolution while maintaining compatibility

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