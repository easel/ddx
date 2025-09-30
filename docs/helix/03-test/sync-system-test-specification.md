# Synchronization System Test Specification

> **Feature**: FEAT-002 - Upstream Synchronization System
> **Phase**: Test (TDD Red)
> **Created**: 2025-01-20
> **Status**: Failing Tests Written

## Overview

This document specifies the test requirements for the DDx synchronization system, covering user stories US-004, US-005, US-009, and US-010. These tests follow Test-Driven Development (TDD) principles and are currently in the "Red" phase - written but failing, awaiting implementation.

## Test Coverage Matrix

| User Story | Acceptance Tests | Contract Tests | Status |
|------------|-----------------|----------------|--------|
| US-004: Update Assets from Master | 7 scenarios | 6 scenarios | ❌ Failing |
| US-005: Contribute Improvements | 6 scenarios | 7 scenarios | ❌ Failing |
| US-009: Pull Updates from Upstream | 2 scenarios | 2 scenarios | ❌ Failing |
| US-010: Handle Update Conflicts | 3 scenarios | 3 scenarios | ❌ Failing |

## US-004: Update Assets from Master

### Acceptance Tests (`sync_acceptance_test.go`)

#### Test: pull_latest_changes
- **Given**: A project with DDx initialized and updates available
- **When**: User runs `ddx update`
- **Then**: Latest changes are fetched from master repository
- **Status**: ❌ Failing - init --no-git flag not implemented

#### Test: display_changelog
- **Given**: Updates are available
- **When**: Running `ddx update --check`
- **Then**: Changelog is displayed showing available updates
- **Status**: ❌ Failing - update command incomplete

#### Test: handle_merge_conflicts
- **Given**: Local changes conflict with upstream
- **When**: Updating with conflicts
- **Then**: Conflicts are detected and resolution options provided
- **Status**: ❌ Failing - conflict detection not implemented

#### Test: selective_update
- **Given**: Multiple assets available for update
- **When**: Running `ddx update templates/nextjs`
- **Then**: Only specified asset is updated
- **Status**: ❌ Failing - selective update not implemented

#### Test: preserve_local_changes
- **Given**: Local modifications exist
- **When**: Running update
- **Then**: Local changes are preserved
- **Status**: ❌ Failing - preservation logic not implemented

#### Test: force_update
- **Given**: Local changes exist that user wants to override
- **When**: Running `ddx update --force`
- **Then**: Updates are applied overriding local changes
- **Status**: ❌ Failing - force flag not implemented

#### Test: create_backup
- **Given**: Project ready for update
- **When**: Running update
- **Then**: Backup is created before changes
- **Status**: ❌ Failing - backup mechanism not implemented

### Contract Tests (`sync_contract_test.go`)

#### Exit Codes
- **0**: Success
- **3**: No configuration file
- **5**: Network error
- **Status**: ❌ Not enforced

#### Flags
- `--check`: Only check for updates without applying
- `--force`: Override local changes
- `--strategy`: Conflict resolution strategy (ours/theirs)
- **Status**: ❌ Flags not implemented

## US-005: Contribute Improvements

### Acceptance Tests

#### Test: contribute_new_template
- **Given**: User has created a new template
- **When**: Running `ddx contribute templates/my-template`
- **Then**: Contribution is prepared with proper structure
- **Status**: ❌ Failing - contribution workflow incomplete

#### Test: validate_contribution
- **Given**: User wants to contribute changes
- **When**: Contributing with `--dry-run`
- **Then**: Contribution is validated against standards
- **Status**: ❌ Failing - validation not implemented

#### Test: create_pull_request
- **Given**: Valid contribution ready
- **When**: Contributing with `--create-pr`
- **Then**: Pull request instructions are provided
- **Status**: ❌ Failing - PR creation not implemented

#### Test: prepare_contribution_branch
- **Given**: User has changes to contribute
- **When**: Preparing contribution
- **Then**: Feature branch is created with proper naming
- **Status**: ❌ Failing - branch creation not implemented

#### Test: validate_contribution_standards
- **Given**: Contribution needs validation
- **When**: Contributing
- **Then**: Standards are checked (metadata, format, etc.)
- **Status**: ❌ Failing - standards validation not implemented

#### Test: push_to_fork
- **Given**: Contribution is ready
- **When**: Using `--push` flag
- **Then**: Changes are pushed to user's fork
- **Status**: ❌ Failing - fork workflow not implemented

### Contract Tests

#### Exit Codes
- **0**: Success
- **6**: Asset not found
- **Status**: ❌ Not enforced

#### Required Elements
- Message (prompt or flag)
- Validation before submission
- Dry-run capability
- Branch creation support
- Push to upstream support
- **Status**: ❌ Not implemented

## US-009: Pull Updates from Upstream

### Acceptance Tests

#### Test: sync_with_upstream
- **Given**: Upstream has new commits
- **When**: Running `ddx update --sync`
- **Then**: Local repository is synchronized
- **Status**: ❌ Failing - sync mechanism not implemented

#### Test: handle_diverged_branches
- **Given**: Local and upstream have diverged
- **When**: Attempting to sync
- **Then**: Divergence is detected and handled
- **Status**: ❌ Failing - divergence detection not implemented

## US-010: Handle Update Conflicts

### Acceptance Tests

#### Test: detect_conflicts
- **Given**: Conflicting changes exist
- **When**: Updating
- **Then**: Conflicts are detected and reported
- **Status**: ❌ Failing - conflict detection not implemented

#### Test: interactive_resolution
- **Given**: Conflicts need resolution
- **When**: Using `--interactive` flag
- **Then**: Interactive options are provided
- **Status**: ❌ Failing - interactive mode not implemented

#### Test: automatic_resolution_strategy
- **Given**: User wants automatic resolution
- **When**: Using `--strategy=ours` or `--strategy=theirs`
- **Then**: Conflicts are resolved automatically
- **Status**: ❌ Failing - strategy system not implemented

## Git Subtree Integration Tests

### Test: subtree_pull
- **Given**: Git repository with subtree configured
- **When**: Pulling via subtree
- **Then**: Changes are merged correctly
- **Status**: ❌ Failing - subtree not configured

### Test: subtree_push
- **Given**: Changes to push via subtree
- **When**: Pushing via subtree
- **Then**: Changes are split and pushed correctly
- **Status**: ❌ Failing - subtree push not implemented

## Implementation Requirements

To make these tests pass, the following must be implemented:

### 1. Update Command Enhancement
- Add `--check`, `--force`, `--sync`, `--strategy` flags
- Implement changelog generation
- Add backup mechanism
- Implement selective asset updates

### 2. Contribute Command Enhancement
- Add `--dry-run`, `--create-pr`, `--push` flags
- Implement validation framework
- Add branch creation logic
- Implement PR workflow

### 3. Git Subtree Integration
- Configure git subtree for .ddx directory
- Implement pull/push operations
- Add conflict detection
- Implement merge strategies

### 4. Conflict Resolution System
- Detect merge conflicts
- Provide resolution strategies
- Implement interactive mode
- Support automatic resolution

### 5. Validation Framework
- Validate contribution standards
- Check metadata requirements
- Ensure format compliance
- Validate git state

## Test Execution

Run all synchronization tests:
```bash
go test ./cmd/ -run "TestAcceptance_US00[4-5]|TestAcceptance_US0[09-11]" -v
go test ./cmd/ -run "TestUpdateCommand_Contract|TestContributeCommand_Contract" -v
```

Run specific user story tests:
```bash
go test ./cmd/ -run "TestAcceptance_US004" -v  # Update tests
go test ./cmd/ -run "TestAcceptance_US005" -v  # Contribute tests
```

## Next Steps (Build Phase)

Once tests are failing consistently (Red phase complete):
1. Implement missing command flags
2. Add git subtree configuration
3. Implement conflict detection
4. Add validation framework
5. Run tests iteratively until all pass (Green phase)
6. Refactor for code quality (Refactor phase)

## Success Criteria

All tests must pass before advancing to the Deploy phase:
- [ ] All acceptance tests passing
- [ ] All contract tests passing
- [ ] Integration tests with git subtree passing
- [ ] Error scenarios properly handled
- [ ] Documentation updated