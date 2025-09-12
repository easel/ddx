# CDP Validators

This directory contains Starlark-based validators that enforce Continuous Documentation-driven Process (CDP) practices in DDx workflows.

## Overview

The CDP validation system ensures that development follows documentation-first, test-driven practices by validating:

1. **Specifications Before Code** - Requirements and design documents must exist before implementation
2. **Test-First Development** - Tests must be written and initially fail before implementation
3. **Complexity Constraints** - Limits concurrent work and enforces complexity boundaries
4. **Documentation Quality** - Ensures documentation completeness and accuracy

## How Validators Work

Each validator is a Starlark script that implements specific validation rules:

- **Input**: Context object containing file changes, git history, and project metadata
- **Output**: List of violations with severity levels and remediation suggestions
- **Integration**: Called by DDx CLI during pre-commit hooks and CI/CD pipelines

### Context Object Structure

```python
ctx = {
    "files": {
        "added": ["path1", "path2"],
        "modified": ["path3", "path4"],
        "deleted": ["path5"]
    },
    "git": {
        "branch": "feature/new-api",
        "commits": [...],
        "diff": "git diff output"
    },
    "project": {
        "root": "/path/to/project",
        "config": {...}
    }
}
```

### Violation Object Structure

```python
violation = {
    "rule": "spec_before_code",
    "severity": "error|warning|info",
    "message": "Human readable description",
    "file": "path/to/violating/file",
    "suggestion": "Recommended fix",
    "line": 42  # optional line number
}
```

## Available Validators

### spec_validator.star
Ensures specifications exist before code implementation.

**Rules:**
- New code files require corresponding specification documents
- Specification documents must be updated when modifying existing code
- Enforces documentation patterns based on file types

### test_validator.star
Enforces test-first development practices.

**Rules:**
- New implementation code requires corresponding tests
- Tests must initially fail (red phase of TDD)
- Test coverage requirements for modified code

### complexity_validator.star
Manages development complexity and concurrent work.

**Rules:**
- Limits number of concurrent features in development
- Calculates and enforces complexity scores
- Prevents overly complex changes in single commits

## Usage

### Direct Validation
```bash
# Run all validators
ddx validate

# Run specific validator
ddx validate --validator=spec_validator

# Validate specific files
ddx validate --files=src/api.go,tests/api_test.go
```

### Git Hook Integration
```bash
# Install pre-commit hook
ddx install-hooks

# Manual pre-commit validation
ddx pre-commit
```

### CI/CD Integration
```yaml
steps:
  - name: CDP Validation
    run: ddx validate --strict
```

## Configuration

Validators can be configured via `.ddx.yml`:

```yaml
cdp:
  validators:
    spec_validator:
      enabled: true
      spec_patterns:
        - "docs/specs/*.md"
        - "docs/design/*.md"
    test_validator:
      enabled: true
      coverage_threshold: 80
    complexity_validator:
      enabled: true
      max_concurrent_features: 3
      complexity_threshold: 10
```

## Extending Validators

Create new validators by:

1. Creating a new `.star` file in this directory
2. Implementing validation functions that return violation lists
3. Using shared utilities from `lib/` directory
4. Adding configuration options to project `.ddx.yml`

## Shared Libraries

### lib/common.star
Common utilities for file operations, string manipulation, and data structures.

### lib/git.star
Git-specific operations like parsing diffs, analyzing commit history, and file change detection.