# Implementation Plan: Configuration Validation

**IP ID**: IP-003
**Feature**: FEAT-003 (Configuration Management)
**Related Documents**: ADR-005, SD-003
**Created**: 2025-01-24
**Status**: Ready

## Overview

Implement two-phase validation system for DDx configuration files, migrating from `.ddx.yml` to `.ddx/config.yaml` with robust error handling and user-friendly messages.

## Prerequisites

- Go 1.21+ installed
- Access to DDx CLI repository
- Understanding of ADR-005 architectural decisions
- Understanding of SD-003 solution design

## Implementation Tasks

### Phase 1: Schema Definition (2 hours)

#### Task 1.1: Create JSON Schema File
- **File**: `internal/config/schema/config.schema.json`
- **Purpose**: Define validation rules for .ddx/config.yaml
- **Schema Structure**:
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "DDx Configuration Schema",
  "type": "object",
  "required": ["version"],
  "properties": {
    "version": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+$",
      "description": "Configuration schema version"
    },
    "library_base_path": {
      "type": "string",
      "default": "./library",
      "description": "Path to DDx library relative to config.yaml"
    },
    "repository": {
      "type": "object",
      "properties": {
        "url": {"type": "string", "format": "uri"},
        "branch": {"type": "string"},
        "subtree_prefix": {"type": "string"}
      },
      "additionalProperties": false
    },
    "variables": {
      "type": "object",
      "additionalProperties": {"type": "string"}
    }
  },
  "additionalProperties": false
}
```

#### Task 1.2: Embed Schema in Binary
- Add `go:embed` directive to include schema in compiled binary
- Create schema loading utility function
- Add schema validation for the schema file itself

#### Task 1.3: Schema Documentation
- Document each schema property with examples
- Add common validation scenarios
- Include migration examples from old format

### Phase 2: Validator Implementation (4 hours)

#### Task 2.1: Add JSON Schema Dependency
```bash
go get github.com/santhosh-tekuri/jsonschema/v5@latest
```

#### Task 2.2: Create Validator Interface
- **File**: `internal/config/validator.go`
- **Core Interface**:
```go
type Validator interface {
    Validate(content []byte) error
    ValidateFile(path string) error
}

type ConfigValidator struct {
    schema *jsonschema.Schema
}

func NewValidator() (*ConfigValidator, error)
func (v *ConfigValidator) Validate(content []byte) error
```

#### Task 2.3: Implement Two-Phase Validation
- **Phase 1**: YAML syntax validation using gopkg.in/yaml.v3
- **Phase 2**: Schema validation using jsonschema
- Error aggregation and formatting
- Performance optimization with schema caching

#### Task 2.4: Create Error Types
- **File**: `internal/config/errors.go`
- **Error Types**:
```go
type ValidationError struct {
    Phase       string   // "syntax" or "schema"
    Message     string   // User-friendly message
    Line        int      // Line number (syntax errors)
    Column      int      // Column number (syntax errors)
    FieldPath   string   // Field path (schema errors)
    Details     string   // Technical details
    Suggestions []string // Helpful suggestions
}

func (e *ValidationError) Error() string
func formatSchemaErrors(err error) []string
func generateSuggestions(err error) []string
```

### Phase 3: Config Loading Updates (3 hours)

#### Task 3.1: Update Configuration Loading
- **File**: `internal/config/config.go`
- **New Loading Logic**:
  1. Check for `.ddx/config.yaml` first
  2. Fall back to `.ddx.yml` (legacy support)
  3. Validate before unmarshaling into structs
  4. Apply defaults for missing values
  5. Show deprecation warnings for old format

#### Task 3.2: Update Configuration Struct
```go
type Config struct {
    Version         string            `yaml:"version" json:"version"`
    LibraryBasePath string            `yaml:"library_base_path" json:"library_base_path"`
    Repository      *RepositoryConfig `yaml:"repository,omitempty" json:"repository,omitempty"`
    Variables       map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"`
}

type RepositoryConfig struct {
    URL           string `yaml:"url,omitempty" json:"url,omitempty"`
    Branch        string `yaml:"branch,omitempty" json:"branch,omitempty"`
    SubtreePrefix string `yaml:"subtree_prefix,omitempty" json:"subtree_prefix,omitempty"`
}
```

#### Task 3.3: Create Defaults Application
- **File**: `internal/config/defaults.go`
- **Default Values**:
  - `library_base_path`: "./library"
  - `repository.url`: "https://github.com/easel/ddx"
  - `repository.branch`: "main"
  - `repository.subtree_prefix`: "library"

### Phase 4: Migration Support (2 hours)

#### Task 4.1: Create Migration Command
- **File**: `cmd/config_migrate.go`
- **Command**: `ddx config migrate`
- **Functionality**:
  - Detect `.ddx.yml` files
  - Convert to new schema format
  - Create `.ddx/` directory structure
  - Move config to `.ddx/config.yaml`
  - Backup original file

#### Task 4.2: Add Auto-Migration Logic
- **Integration Point**: `internal/config/config.go`
- **Auto-Migration**:
  - Detect legacy format during config loading
  - Prompt user for migration (or auto-migrate with flag)
  - Show migration success/failure messages
  - Log migration actions

#### Task 4.3: Legacy Support Warnings
- Show deprecation warnings when loading `.ddx.yml`
- Provide clear migration instructions
- Set timeline for legacy format removal

### Phase 5: Test Updates (4 hours)

#### Task 5.1: Update Test Fixtures
- Move all test configs from `.ddx.yml` to `.ddx/config.yaml`
- Update file paths in test files
- Fix directory structure in test setups
- Update expected outputs for new format

#### Task 5.2: Add Validation Tests
- **File**: `internal/config/validator_test.go`
- **Test Scenarios**:
  - Valid configurations
  - YAML syntax errors (indentation, quotes, etc.)
  - Schema violations (missing required fields, wrong types)
  - Default value application
  - Error message formatting
  - Performance benchmarks

#### Task 5.3: Add Migration Tests
- **File**: `cmd/config_migrate_test.go`
- **Test Scenarios**:
  - Successful migration from `.ddx.yml`
  - Edge cases (malformed configs, missing directories)
  - Backup and rollback functionality
  - Auto-migration prompt behavior

#### Task 5.4: Fix Failing Tests
- Update test expectations for new validation
- Fix path references throughout test suite
- Update helper functions for new config structure
- Ensure all tests pass with new validation

### Phase 6: Documentation Updates (1 hour)

#### Task 6.1: Update User Documentation
- Update configuration examples in README
- Add migration guide to documentation
- Update CLI help text for new format

#### Task 6.2: Update CHANGELOG
- Document breaking changes
- Provide migration instructions
- Set deprecation timeline for legacy format

#### Task 6.3: Update Development Documentation
- Update CLAUDE.md with new config patterns
- Add validation examples to development docs

## Implementation Order

### Day 1 (6 hours):
1. ✅ Documentation updates (ADR-005, SD-003)
2. ⏳ Phase 1: Schema Definition (2 hours)
3. ⏳ Phase 2: Validator Implementation (4 hours)

### Day 2 (6 hours):
4. ⏳ Phase 3: Config Loading Updates (3 hours)
5. ⏳ Phase 4: Migration Support (2 hours)
6. ⏳ Phase 6: Documentation Updates (1 hour)

### Day 3 (4 hours):
7. ⏳ Phase 5: Test Updates (4 hours)

## Success Criteria

- [ ] All tests pass with new config validation
- [ ] Validation completes in <20ms for typical configs
- [ ] Migration works seamlessly for existing projects
- [ ] Clear, actionable error messages for invalid configs
- [ ] Backwards compatibility maintained during transition
- [ ] No breaking changes for users who migrate properly
- [ ] Performance benchmarks meet targets

## Risk Mitigation

### Risk: Breaking Existing Projects
- **Mitigation**: Auto-migration with legacy support
- **Fallback**: Manual migration command
- **Communication**: Clear migration documentation

### Risk: Complex Error Messages
- **Mitigation**: Custom error formatting with suggestions
- **Testing**: User experience testing with invalid configs
- **Documentation**: Error message examples in docs

### Risk: Performance Degradation
- **Mitigation**: Benchmarking with realistic configs
- **Target**: <20ms validation time
- **Optimization**: Schema caching and lazy loading

### Risk: Schema Too Restrictive
- **Mitigation**: Start with minimal required fields
- **Evolution**: Use `additionalProperties: false` carefully
- **Flexibility**: Allow custom variables section

## Validation Approach

### Testing Strategy
1. **Unit Tests**: Individual validator components
2. **Integration Tests**: End-to-end config loading
3. **Performance Tests**: Validation speed benchmarks
4. **Migration Tests**: Legacy format conversion
5. **Error Message Tests**: User experience validation

### Acceptance Criteria
- All existing functionality preserved
- New validation catches configuration errors
- Migration path works for all existing projects
- Performance meets targets (<20ms)
- Error messages provide actionable guidance

## Dependencies

### Internal Dependencies
- `internal/config/config.go` - Core configuration loading
- `cmd/init.go` - Project initialization
- `cmd/config.go` - Configuration management commands

### External Dependencies
- `github.com/santhosh-tekuri/jsonschema/v5` - Schema validation
- `gopkg.in/yaml.v3` - YAML parsing (already present)
- `github.com/spf13/viper` - Configuration management (already present)

## Monitoring and Metrics

### Performance Metrics
- Configuration validation time
- Schema compilation time
- Memory usage during validation
- File I/O performance

### Success Metrics
- Migration success rate
- User error reporting (validation catches issues)
- Support ticket reduction for config issues
- Adoption rate of new format