# API Contract: Utility Commands [FEAT-001]

**Contract ID**: CLI-003
**Feature**: FEAT-001
**Type**: CLI
**Status**: Approved
**Version**: 1.0.0

*Command-line interface contracts for DDx utility commands*

## Command: diagnose

**Purpose**: Analyze project health, identify issues, and suggest improvements based on DDx best practices
**Usage**: `ddx diagnose [checks] [options]`

### Arguments
- `checks` (optional): Specific checks to run (comma-separated) or "all" (default: all)
  - Available checks: config, dependencies, structure, templates, patterns, security, performance, quality

### Options
- `--fix`: Attempt automatic fixes for issues
- `--format <format>`: Output format: "text" (default), "json", "markdown", "html"
- `--severity <level>`: Minimum severity to report: "info", "warning", "error", "critical"
- `--verbose, -v`: Include detailed diagnostics
- `--quiet, -q`: Only show errors and critical issues
- `--report <file>`: Save report to file
- `--baseline <file>`: Compare against baseline report

### Diagnostic Checks

#### Configuration Health
- Valid .ddx.yml syntax
- Required fields present
- Repository connectivity
- Git subtree health
- Variable definitions

#### Dependency Analysis
- Missing dependencies
- Version conflicts
- Outdated packages
- Security vulnerabilities
- License compatibility

#### Project Structure
- DDx conventions followed
- File organization
- Naming conventions
- Directory structure

#### Template/Pattern Compliance
- Proper variable usage
- Template syntax valid
- Pattern implementation
- Customization tracking

### Output Format (Text - Default)
```
DDx Project Diagnostics Report
==============================
Project: my-awesome-app
Path: /Users/jane/projects/my-awesome-app
DDx Version: 1.2.0
Date: 2025-01-15 14:30:00

Configuration Check
─────────────────
✓ .ddx.yml valid and complete
✓ Repository connection successful
⚠ 2 undefined variables found

Dependency Analysis
─────────────────
✓ All required dependencies installed
⚠ 3 packages outdated
✗ 1 high severity vulnerability

Project Structure
─────────────────
✓ Follows DDx conventions
✓ Template structure intact
⚠ 5 files don't follow naming convention

Template Compliance
─────────────────
✓ Template: nextjs (90% compliance)
⚠ Modified files: 3
  - package.json (expected)
  - src/app/page.tsx (review needed)
  - next.config.js (may affect updates)

Security Scan
─────────────────
✗ Hardcoded API key found in .env.example
⚠ Missing security headers configuration
✓ No exposed secrets in git history

Performance Analysis
─────────────────
✓ Build size within limits
⚠ Large node_modules (523 MB)
✓ Efficient file structure

SUMMARY
─────────────────
Critical: 0
Errors: 2
Warnings: 7
Info: 12

Recommended Actions:
1. Run 'npm audit fix' to resolve vulnerability
2. Remove API key from .env.example
3. Run 'ddx diagnose --fix' to auto-fix 4 issues

Health Score: 78/100 (Good)
```

### Output Format (JSON)
```json
{
  "project": {
    "name": "my-awesome-app",
    "path": "/Users/jane/projects/my-awesome-app",
    "ddx_version": "1.2.0",
    "timestamp": "2025-01-15T14:30:00Z"
  },
  "checks": {
    "configuration": {
      "status": "warning",
      "issues": [
        {
          "severity": "warning",
          "message": "Undefined variable: {{api_key}}",
          "file": ".ddx.yml",
          "line": 15,
          "fixable": true
        }
      ]
    },
    "dependencies": {
      "status": "error",
      "vulnerabilities": [
        {
          "package": "lodash",
          "severity": "high",
          "cve": "CVE-2024-1234"
        }
      ]
    }
  },
  "summary": {
    "critical": 0,
    "errors": 2,
    "warnings": 7,
    "info": 12,
    "health_score": 78,
    "fixable": 4
  }
}
```

### Exit Codes
- `0`: All checks passed
- `1`: Warnings found
- `2`: Errors found
- `3`: Critical issues found
- `4`: Diagnostic check failed
- `5`: Invalid check specified

### Examples
```bash
# Full diagnostic
$ ddx diagnose
Running full project diagnostics...
[output as above]

# Specific checks only
$ ddx diagnose config,security
Running diagnostics: configuration, security
✓ Configuration valid
✗ Security issue: exposed secret

# Auto-fix issues
$ ddx diagnose --fix
Found 4 fixable issues:
  ✓ Fixed: Updated .gitignore
  ✓ Fixed: Corrected file permissions
  ✓ Fixed: Normalized line endings
  ✗ Manual fix required: Remove hardcoded secret

# Generate report
$ ddx diagnose --format markdown --report health-report.md
Diagnostic report saved to health-report.md

# Compare with baseline
$ ddx diagnose --baseline last-report.json
Comparing with baseline from 2025-01-10...
Improvements:
  - Reduced warnings from 12 to 7
  - Fixed all critical issues
Regressions:
  - New security warning added
```

---

## Command: config

**Purpose**: Manage DDx configuration settings
**Usage**: `ddx config [action] [key] [value] [options]`

### Actions
- `get <key>`: Get configuration value
- `set <key> <value>`: Set configuration value
- `list`: List all configuration
- `validate`: Validate configuration
- `edit`: Open configuration in editor
- `init`: Reinitialize configuration

### Options
- `--global, -g`: Use global configuration (~/.ddx/config.yml)
- `--local, -l`: Use local configuration (.ddx.yml)
- `--format <format>`: Output format for list/get: "text", "json", "yaml"
- `--merge`: Merge with existing configuration
- `--schema`: Show configuration schema

### Configuration Scopes
1. **Local** (.ddx.yml): Project-specific settings
2. **Global** (~/.ddx/config.yml): User defaults
3. **System** (/etc/ddx/config.yml): System-wide defaults

### Output Format
```
# Get single value
$ ddx config get repository.url
https://github.com/easel/ddx

# List all configuration
$ ddx config list
DDx Configuration (.ddx.yml)
─────────────────────────────
repository:
  url: https://github.com/easel/ddx
  branch: main

includes:
  - templates
  - patterns
  - prompts

variables:
  project_name: my-app
  port: 3000

templates:
  active: nextjs
  version: 1.2.0

# Set value
$ ddx config set variables.port 3001
✓ Updated variables.port to 3001

# Validate configuration
$ ddx config validate
✓ Configuration valid
✓ Repository accessible
✓ All included resources found
```

### Configuration Schema
```yaml
# .ddx.yml schema
version: string          # DDx config version
repository:
  url: string           # Master repository URL
  branch: string        # Branch to track
  subtree_path: string  # Local subtree path

includes: array         # Resources to include
  - templates
  - patterns
  - prompts
  - configs

variables: object       # Variable substitutions
  key: value

templates:
  active: string       # Active template
  version: string      # Template version

settings:
  auto_update: boolean # Auto-update on init
  backup: boolean      # Create backups
  color: boolean       # Colored output
  verbose: boolean     # Verbose logging
```

### Exit Codes
- `0`: Success
- `1`: Configuration not found
- `2`: Invalid key
- `3`: Validation failed
- `4`: Permission denied
- `5`: Syntax error

### Examples
```bash
# Get nested value
$ ddx config get repository.branch
main

# Set nested value
$ ddx config set settings.auto_update true
✓ Enabled auto-update

# Edit in default editor
$ ddx config edit
Opening .ddx.yml in vim...

# Show global config
$ ddx config list --global
Global DDx Configuration
─────────────────────
default_repo: https://github.com/easel/ddx
author: Jane Developer
email: jane@example.com

# Validate with verbose output
$ ddx config validate -v
Validating .ddx.yml...
✓ YAML syntax valid
✓ Schema validation passed
✓ Repository https://github.com/easel/ddx accessible
✓ Branch 'main' exists
✓ All variables defined
✓ No circular references

# Initialize new config
$ ddx config init --merge
Merging with existing configuration...
✓ Configuration updated
```

---

## Command: version

**Purpose**: Display DDx version information and check for updates
**Usage**: `ddx version [options]`

### Options
- `--check`: Check for available updates
- `--format <format>`: Output format: "text" (default), "json", "short"
- `--detailed`: Show detailed build information
- `--components`: Show component versions

### Output Format (Text - Default)
```
DDx - Document-Driven Development CLI
Version: 1.2.0
Built: 2025-01-10 10:30:00 UTC
Go Version: 1.21.5
Platform: darwin/arm64
```

### Output Format (Detailed)
```
DDx - Document-Driven Development CLI
─────────────────────────────────────
Version: 1.2.0
Build Date: 2025-01-10 10:30:00 UTC
Git Commit: abc123def456789
Go Version: go1.21.5
Platform: darwin/arm64
Compiler: gc
Build User: ci-runner
Build Host: github-actions

Components:
─────────────────────────────────────
CLI Framework: Cobra v1.8.0
Config Manager: Viper v1.17.0
Git Library: go-git v5.11.0
Template Engine: v2.0.0

Configuration:
─────────────────────────────────────
Config File: .ddx.yml
Repository: https://github.com/easel/ddx
Branch: main
Last Update: 2025-01-15 08:00:00
```

### Output Format (JSON)
```json
{
  "version": "1.2.0",
  "build_date": "2025-01-10T10:30:00Z",
  "git_commit": "abc123def456789",
  "go_version": "go1.21.5",
  "platform": "darwin/arm64",
  "components": {
    "cobra": "1.8.0",
    "viper": "1.17.0",
    "go-git": "5.11.0"
  }
}
```

### Update Check
```
$ ddx version --check
Current Version: 1.2.0
Latest Version: 1.3.0 (released 2025-01-14)

Update available! New features:
- Improved template engine
- New patterns for microservices
- Performance improvements

Update with: curl -sSL https://ddx.dev/install.sh | sh
```

### Exit Codes
- `0`: Success
- `1`: Update check failed
- `2`: Network error during update check

### Examples
```bash
# Simple version
$ ddx version
DDx version 1.2.0

# Short format
$ ddx version --format short
1.2.0

# Check for updates
$ ddx version --check
Current: 1.2.0
Latest: 1.3.0
Update available!

# Component versions
$ ddx version --components
DDx 1.2.0
├── Cobra 1.8.0
├── Viper 1.17.0
├── go-git 5.11.0
└── Template Engine 2.0.0

# JSON output for scripts
$ ddx version --format json
{"version":"1.2.0","build_date":"2025-01-10T10:30:00Z"}
```

## Global Utility Features

### Help System
All commands support comprehensive help:
```bash
# General help
$ ddx help
DDx - Document-Driven Development CLI

Usage:
  ddx [command] [flags]

Available Commands:
  init        Initialize DDx in a project
  list        List available resources
  apply       Apply resources to project
  update      Update from master repository
  contribute  Share improvements back
  diagnose    Analyze project health
  config      Manage configuration
  version     Show version information

# Command-specific help
$ ddx help diagnose
Analyze project health and identify issues...
[detailed help]

# Quick help
$ ddx diagnose -h
[concise help]
```

### Debug Mode
All commands support debug output:
```bash
$ DDX_DEBUG=1 ddx diagnose
[DEBUG] Loading configuration from .ddx.yml
[DEBUG] Parsing YAML...
[DEBUG] Running check: configuration
[DEBUG] Connecting to repository...
```

### Quiet Mode
Suppress non-essential output:
```bash
$ ddx diagnose -q
✗ 2 errors found
```

### Non-Interactive Mode
For CI/CD environments:
```bash
$ ddx apply templates/nextjs --non-interactive --vars-file ci-vars.yml
```

## Performance Requirements

All utility commands must:
- Start execution within 100ms
- Complete typical operations within 5 seconds
- Handle projects with 10,000+ files
- Work efficiently on slow network connections
- Minimize memory usage (< 100MB typical)

## Error Message Standards

Error messages must include:
1. Clear description of what went wrong
2. Specific file/line when applicable
3. Suggested fix or next steps
4. Error code for automation

Example:
```
Error [E1234]: Configuration syntax error
File: .ddx.yml, Line 23
Problem: Unexpected indentation
Fix: Ensure consistent spacing (2 spaces recommended)
Run 'ddx config validate -v' for detailed analysis
```