# DDx Comprehensive Test Specification

> **Phase**: Test (TDD Red)
> **Created**: 2025-01-20
> **Status**: Failing Tests Complete

## Executive Summary

This document provides a comprehensive specification of all tests written for the DDx project during the Test phase of the HELIX workflow. All tests are currently in the "Red" state (failing), which is the expected outcome for Test-Driven Development. These tests define the expected behavior for unimplemented features.

## Test Coverage Overview

| Feature Area | User Stories | Acceptance Tests | Contract Tests | Status |
|--------------|-------------|------------------|----------------|--------|
| Synchronization System | US-004, 005, 009-011 | 18 scenarios | 16 scenarios | ❌ Failing |
| MCP Server Management | US-036, 037 | 10 scenarios | 10 scenarios | ❌ Failing |
| Configuration Management | US-017-024 | 7 scenarios | 5 scenarios | ❌ Failing |
| Installation & Setup | US-028-035 | 8 scenarios | 4 scenarios | ❌ Failing |
| **Total** | **20 stories** | **43 scenarios** | **35 scenarios** | **78 tests** |

## Test Organization

### Test Files Created

1. **`sync_acceptance_test.go`** - Synchronization system acceptance tests
2. **`sync_contract_test.go`** - Synchronization command contracts
3. **`mcp_acceptance_test.go`** - MCP and configuration acceptance tests
4. **`mcp_contract_test.go`** - MCP, config, and installation contracts

### Test Naming Convention

- **Acceptance Tests**: `TestAcceptance_US###_Description`
- **Contract Tests**: `Test<Command>Command_Contract`
- **Test Scenarios**: `snake_case` describing the scenario

## Synchronization System Tests

### US-004: Update Assets from Master

**Purpose**: Ensure users can pull latest improvements from master repository

**Test Scenarios**:
1. `pull_latest_changes` - Basic update functionality
2. `display_changelog` - Show what's new before applying
3. `handle_merge_conflicts` - Graceful conflict resolution
4. `selective_update` - Update specific assets only
5. `preserve_local_changes` - Don't overwrite customizations
6. `force_update` - Override when explicitly requested
7. `create_backup` - Safety net before changes

**Key Failure**: Missing `--no-git` flag in init command

### US-005: Contribute Improvements

**Purpose**: Enable users to share improvements back to community

**Test Scenarios**:
1. `contribute_new_template` - Share new templates
2. `validate_contribution` - Ensure quality standards
3. `create_pull_request` - GitHub workflow integration

**Key Failure**: Contribution workflow not implemented

### US-009-011: Advanced Sync Features

**Purpose**: Handle complex synchronization scenarios

**Test Scenarios**:
- Sync with upstream
- Handle diverged branches
- Detect conflicts
- Interactive resolution
- Automatic resolution strategies
- Prepare contribution branches
- Validate contribution standards
- Push to fork

**Key Failures**: Git subtree not configured, conflict detection missing

## MCP Server Management Tests

### US-036: List Available MCP Servers

**Purpose**: Discover and browse available MCP integrations

**Test Scenarios**:
1. `display_all_available_servers` - Complete listing
2. `filter_by_category` - Category-based browsing
3. `search_functionality` - Find specific servers
4. `show_installation_status` - Visual indicators (✅/⬜)
5. `detailed_verbose_view` - Additional information

**Key Failure**: Registry not found at `mcp-servers/registry.yml`

### US-037: Install MCP Server

**Purpose**: Install and configure MCP servers locally

**Test Scenarios**:
1. `install_server_locally` - Basic installation
2. `detect_package_manager` - npm/pnpm/yarn/bun detection
3. `configure_server_environment` - Set up environment variables
4. `handle_already_installed` - Graceful re-installation
5. `validate_installation` - Connection testing

**Expected Artifacts**:
- `package.json` creation
- `.claude/settings.local.json` configuration
- Node modules installation

## Configuration Management Tests

### US-017-024: Configuration System

**Purpose**: Flexible, hierarchical configuration management

**Test Scenarios**:
1. `initialize_configuration` - Create initial config
2. `configure_variables` - Template variable management
3. `override_configuration` - Environment-specific settings
4. `configure_resource_selection` - Choose what to include
5. `validate_configuration` - Ensure correctness
6. `export_import_configuration` - Share configurations
7. `view_effective_configuration` - Merged view of all layers

**Configuration Hierarchy**:
1. Global (`~/.ddx/config.yml`)
2. Project (`.ddx.yml`)
3. Local override (`.ddx.local.yml`)
4. Environment (`.ddx.<env>.yml`)

## Installation & Setup Tests

### US-028-035: Installation Features

**Purpose**: Easy installation and maintenance of DDx

**Test Scenarios**:
1. `one_command_installation` - Simple setup
2. `automatic_path_configuration` - Shell integration
3. `installation_verification` - Doctor command
4. `package_manager_installation` - Homebrew/apt/choco
5. `upgrade_existing_installation` - Self-update
6. `uninstall_ddx` - Clean removal
7. `offline_installation` - Bundle support
8. `installation_diagnostics` - Troubleshooting

**Key Commands**:
- `ddx doctor` - Verify installation
- `ddx self-update` - Upgrade DDx
- `ddx setup path` - Configure PATH
- `ddx uninstall` - Remove DDx

## Contract Tests

### Exit Code Contracts

All commands follow consistent exit code patterns:

| Code | Meaning | Example |
|------|---------|---------|
| 0 | Success | Command completed successfully |
| 1 | General error | Unexpected failure |
| 2 | Already exists | Config/resource already present |
| 3 | No configuration | Missing .ddx.yml |
| 4 | Not found (template) | Template doesn't exist |
| 5 | Network error | Can't reach repository |
| 6 | Not found (general) | Resource doesn't exist |

### Output Format Contracts

Commands provide consistent output:
- Status indicators (✓, ✗, ⚠)
- Progress messages
- Structured data (tables, JSON, YAML)
- Color coding for clarity

### Flag Contracts

Common flags across commands:
- `--dry-run` - Preview without changes
- `--force` - Override safety checks
- `--verbose` - Detailed output
- `--check` - Verify without applying
- `--validate` - Test after operation

## Test Execution Guide

### Run All Tests
```bash
cd cli
go test ./cmd/ -v
```

### Run Specific Feature Tests
```bash
# Synchronization tests
go test ./cmd/ -run "TestAcceptance_US00[4-5]|TestAcceptance_US0[09-11]" -v

# MCP tests
go test ./cmd/ -run "TestAcceptance_US03[6-7]" -v

# Configuration tests
go test ./cmd/ -run "ConfigurationManagement" -v

# Contract tests
go test ./cmd/ -run "Contract" -v
```

### Current Test Results

All 78 test scenarios are currently **failing**, which is the expected state for TDD Red phase:

- ✅ Tests compile successfully
- ✅ Tests run and fail with meaningful errors
- ✅ Failure messages indicate what needs implementation
- ✅ Test coverage defines expected behavior

## Implementation Requirements

To transition from Red to Green phase, implement:

### 1. Synchronization System
- Git subtree configuration
- Conflict detection and resolution
- Backup mechanisms
- Changelog generation

### 2. MCP Server Management
- Registry loading from library
- Package manager detection
- Claude config generation
- Installation validation

### 3. Configuration Management
- Hierarchical config merging
- Variable substitution
- Import/export functionality
- Validation framework

### 4. Installation Features
- Self-update mechanism
- PATH configuration
- Doctor diagnostics
- Uninstall cleanup

## Success Criteria for Test Phase

✅ **Completed**:
- All user stories have corresponding tests
- Tests follow consistent naming conventions
- Contract tests define command behavior
- Acceptance tests cover user scenarios
- Tests fail with clear error messages
- Test specifications documented

## Next Phase: Build

With all tests in place and failing, we're ready to enter the Build phase where we'll:

1. Implement features to make tests pass
2. Follow TDD cycle: Red → Green → Refactor
3. Achieve >80% test coverage
4. Ensure all acceptance criteria met
5. Validate against user stories

## Test Maintenance

As we implement features:
1. Run tests frequently
2. Fix one failing test at a time
3. Refactor after tests pass
4. Add edge case tests as discovered
5. Keep tests as living documentation

## Appendix: Test Statistics

- **Total Test Files**: 4 (+2 existing)
- **Total Test Functions**: 78
- **Lines of Test Code**: ~1,500
- **User Stories Covered**: 20
- **Features Tested**: 4 major areas
- **Expected Build Time**: 2-3 sprints