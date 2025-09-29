# TS-014: Agent-Based Workflow Enforcement Test Specification

**Test Spec ID**: TS-014
**Feature**: FEAT-014 - Agent-Based Workflow Enforcement
**Status**: Draft
**Created**: 2025-01-20
**Updated**: 2025-01-20
**Author**: Core Team

## Overview

This document specifies the test requirements for refactoring HELIX workflow enforcement from passive CLAUDE.md instructions to an active agent system.

## Test Objectives

1. **Verify Token Reduction**: Confirm 45% reduction in base context
2. **Validate Functionality**: Ensure all workflow commands work identically
3. **Confirm Agent Activation**: Verify agent activates on workflow requests
4. **Check Path Consistency**: Ensure all library paths use `.ddx/library/*` format
5. **Regression Prevention**: No functionality changes

## Test Categories

### 1. CLAUDE.md Structure Tests

#### Test 1.1: Line Count Reduction
```yaml
test_id: TS-014-001
category: structure
priority: P0
description: Verify CLAUDE.md reduced to target size

setup:
  - Measure current CLAUDE.md line count

test_steps:
  - Count lines in refactored CLAUDE.md
  - Verify line count ≤ 230 lines (target 220, allow 10 line buffer)

expected_result:
  - Line count reduced from ~400 to ~220 lines
  - Token reduction of ~45%

pass_criteria:
  - new_line_count <= 230
  - reduction_percentage >= 40%
```

#### Test 1.2: Required Sections Present
```yaml
test_id: TS-014-002
category: structure
priority: P0
description: Verify all required sections exist in minimal CLAUDE.md

test_steps:
  - Check for "Project Overview" section
  - Check for "Development Commands" section
  - Check for "Workflow System" section (new minimal section)
  - Check for "Architectural Principles" section
  - Check for "Testing Requirements" section

expected_result:
  - All required sections present
  - Each section appropriately sized

pass_criteria:
  - all_required_sections_present == true
  - no_enforcement_sections_remain == true
```

#### Test 1.3: Enforcement Sections Removed
```yaml
test_id: TS-014-003
category: structure
priority: P0
description: Verify workflow enforcement sections removed

test_steps:
  - Search for "HELIX Workflow Enforcement" section
  - Search for "Activation Instructions" section
  - Search for "Phase Detection" section
  - Search for "Auto-Prompts" section
  - Search for "Workflow Auto-Continuation" section

expected_result:
  - No enforcement sections found
  - No auto-prompt sections found

pass_criteria:
  - enforcement_section_count == 0
```

#### Test 1.4: Path Consistency
```yaml
test_id: TS-014-004
category: structure
priority: P0
description: Verify all library paths use .ddx/library/* format

test_steps:
  - Grep for "library/" pattern
  - Check each match has ".ddx/library/" format
  - Verify no standalone "library/" references (except in examples)

expected_result:
  - All library references use ".ddx/library/*" format
  - No ambiguous "library/*" references

pass_criteria:
  - all_paths_use_ddx_prefix == true
  - no_ambiguous_references == true
```

### 2. Agent Activation Tests

#### Test 2.1: Workflow Command Activation
```yaml
test_id: TS-014-005
category: agent
priority: P0
description: Verify agent activates on workflow commands

test_commands:
  - "ddx workflow helix execute build-story US-001"
  - "ddx workflow helix execute continue"
  - "ddx workflow helix execute status"
  - "ddx workflow helix execute next"

test_steps:
  - Execute each command
  - Monitor for Task tool invocation
  - Verify agent initialization

expected_result:
  - Task tool launches helix-workflow agent
  - Agent loads coordinator.md
  - Agent executes appropriate action

pass_criteria:
  - agent_activated == true for all commands
  - task_tool_invoked == true
```

#### Test 2.2: Keyword Detection
```yaml
test_id: TS-014-006
category: agent
priority: P1
description: Verify agent activates on workflow keywords

test_phrases:
  - "Work on US-001"
  - "Continue the workflow"
  - "Check workflow status"
  - "Work on next story"

test_steps:
  - Submit each phrase to Claude Code
  - Monitor for agent activation
  - Verify appropriate action taken

expected_result:
  - Agent activates for workflow-related requests
  - Appropriate action executed

pass_criteria:
  - agent_activation_rate >= 95%
```

#### Test 2.3: No Activation on Non-Workflow Tasks
```yaml
test_id: TS-014-007
category: agent
priority: P0
description: Verify agent does NOT activate for non-workflow tasks

test_phrases:
  - "Write a function to sort an array"
  - "Explain how React hooks work"
  - "Debug this error message"
  - "Create a new file"

test_steps:
  - Submit each phrase to Claude Code
  - Monitor for agent activation
  - Verify direct response without agent

expected_result:
  - No workflow agent activation
  - Normal Claude Code behavior

pass_criteria:
  - agent_not_activated == true for all non-workflow tasks
  - no_unnecessary_overhead == true
```

### 3. Workflow Functionality Tests

#### Test 3.1: build-story Command
```yaml
test_id: TS-014-008
category: functionality
priority: P0
description: Verify build-story command works identically

setup:
  - Create test user story
  - Set workflow phase to "build"

test_steps:
  - Execute: ddx workflow helix execute build-story US-TEST
  - Verify story is read
  - Verify phase enforcer loaded
  - Verify implementation executed

expected_result:
  - Command executes successfully
  - Same behavior as before refactoring

pass_criteria:
  - command_succeeds == true
  - no_regression == true
```

#### Test 3.2: continue Command
```yaml
test_id: TS-014-009
category: functionality
priority: P0
description: Verify continue command works identically

setup:
  - Establish workflow context with active story

test_steps:
  - Execute: ddx workflow helix execute continue
  - Verify current work resumed
  - Verify phase enforcement applied

expected_result:
  - Command executes successfully
  - Work continues from previous state

pass_criteria:
  - command_succeeds == true
  - context_preserved == true
```

#### Test 3.3: status Command
```yaml
test_id: TS-014-010
category: functionality
priority: P0
description: Verify status command works identically

test_steps:
  - Execute: ddx workflow helix execute status
  - Verify current phase reported
  - Verify progress information shown

expected_result:
  - Command executes successfully
  - Accurate status information

pass_criteria:
  - command_succeeds == true
  - status_accurate == true
```

#### Test 3.4: next Command
```yaml
test_id: TS-014-011
category: functionality
priority: P0
description: Verify next command works identically

setup:
  - Multiple user stories available

test_steps:
  - Execute: ddx workflow helix execute next
  - Verify next priority story selected
  - Verify work begins on story

expected_result:
  - Command executes successfully
  - Correct story selected

pass_criteria:
  - command_succeeds == true
  - correct_story_selected == true
```

### 4. Phase Enforcement Tests

#### Test 4.1: Frame Phase Enforcement
```yaml
test_id: TS-014-012
category: enforcement
priority: P0
description: Verify Frame phase prevents premature coding

setup:
  - Set workflow phase to "frame"

test_steps:
  - Attempt to write implementation code
  - Verify phase violation detected
  - Verify guidance provided

expected_result:
  - Violation detected by agent
  - User guided to Frame phase activities

pass_criteria:
  - violation_detected == true
  - guidance_provided == true
```

#### Test 4.2: Build Phase Enforcement
```yaml
test_id: TS-014-013
category: enforcement
priority: P0
description: Verify Build phase enforces test-driven development

setup:
  - Set workflow phase to "build"

test_steps:
  - Attempt to implement without failing tests
  - Verify phase enforcer guides to write tests first

expected_result:
  - Enforcer requires tests before implementation
  - TDD principles maintained

pass_criteria:
  - tdd_enforced == true
  - guidance_appropriate == true
```

### 5. Path Resolution Tests

#### Test 5.1: Coordinator Loading
```yaml
test_id: TS-014-014
category: paths
priority: P0
description: Verify coordinator.md loads from correct path

test_steps:
  - Trigger agent activation
  - Monitor file access
  - Verify path: .ddx/library/workflows/helix/coordinator.md

expected_result:
  - Coordinator loaded successfully
  - Correct path used

pass_criteria:
  - coordinator_loaded == true
  - correct_path_used == true
```

#### Test 5.2: Enforcer Loading
```yaml
test_id: TS-014-015
category: paths
priority: P0
description: Verify phase enforcers load from correct paths

test_phases:
  - frame: .ddx/library/workflows/helix/phases/01-frame/enforcer.md
  - build: .ddx/library/workflows/helix/phases/04-build/enforcer.md

test_steps:
  - For each phase, trigger agent
  - Verify correct enforcer loaded

expected_result:
  - Enforcers load from .ddx/library/* paths
  - Correct phase enforcer for current phase

pass_criteria:
  - all_enforcers_load_correctly == true
  - paths_correct == true
```

### 6. Performance Tests

#### Test 6.1: Token Usage Measurement
```yaml
test_id: TS-014-016
category: performance
priority: P0
description: Measure token usage reduction

test_steps:
  - Measure tokens in original CLAUDE.md
  - Measure tokens in refactored CLAUDE.md
  - Calculate reduction percentage

expected_result:
  - Original: ~20,000 tokens
  - Refactored: ~11,000 tokens
  - Reduction: ~45%

pass_criteria:
  - token_reduction >= 40%
  - target_met == true
```

#### Test 6.2: Agent Activation Overhead
```yaml
test_id: TS-014-017
category: performance
priority: P1
description: Measure agent activation overhead

test_steps:
  - Time workflow command execution before refactoring
  - Time workflow command execution after refactoring
  - Compare execution times

expected_result:
  - Overhead < 20ms
  - No noticeable slowdown

pass_criteria:
  - overhead_ms < 20
  - acceptable_performance == true
```

### 7. Error Handling Tests

#### Test 7.1: Missing Library Directory
```yaml
test_id: TS-014-018
category: error_handling
priority: P1
description: Verify graceful handling when .ddx/library missing

setup:
  - Temporarily rename .ddx/library

test_steps:
  - Trigger workflow command
  - Verify error message helpful
  - Verify fallback behavior

expected_result:
  - Clear error message
  - Suggests running: ddx update
  - Provides manual recovery steps

pass_criteria:
  - error_message_helpful == true
  - recovery_path_clear == true
```

#### Test 7.2: Agent Activation Failure
```yaml
test_id: TS-014-019
category: error_handling
priority: P1
description: Verify graceful handling when agent fails to activate

setup:
  - Simulate Task tool failure

test_steps:
  - Trigger workflow command
  - Verify fallback behavior
  - Verify user notified

expected_result:
  - Degraded but functional behavior
  - Error logged for debugging
  - User sees helpful message

pass_criteria:
  - graceful_degradation == true
  - user_notified == true
```

### 8. Regression Tests

#### Test 8.1: All Existing Tests Pass
```yaml
test_id: TS-014-020
category: regression
priority: P0
description: Verify no existing tests break

test_steps:
  - Run full test suite
  - Compare results before/after refactoring

expected_result:
  - All tests that passed before still pass
  - No new test failures

pass_criteria:
  - no_new_failures == true
  - test_count_unchanged == true
```

#### Test 8.2: User Experience Unchanged
```yaml
test_id: TS-014-021
category: regression
priority: P0
description: Verify user-facing behavior identical

test_scenarios:
  - Work on user story
  - Continue workflow
  - Check status
  - Handle phase violations

test_steps:
  - Execute each scenario before/after refactoring
  - Compare user-visible behavior

expected_result:
  - Identical user experience
  - No visible changes

pass_criteria:
  - behavior_identical == true
  - no_user_visible_changes == true
```

## Test Data Requirements

### Test User Stories
Create test stories for validation:
- US-TEST-001: Simple build task
- US-TEST-002: Complex multi-step task
- US-TEST-003: Story requiring phase transition

### Test Workflow States
Phase is inferred from artifacts, not state files:
- Frame phase: docs/helix/01-frame/ exists, docs/helix/02-design/ missing
- Design phase: docs/helix/02-design/ exists, no test files
- Build phase: Test files exist and failing
- Deploy phase: Tests passing, not deployed

### Test Library Content
- Complete .ddx/library structure
- Missing .ddx/library (error case)
- Corrupted coordinator.md (error case)

## Test Environment Setup

### Prerequisites
1. DDx CLI installed
2. Test project with .ddx.yml
3. .ddx/library synced from ddx-library repo
4. Test user stories created
5. Workflow state files initialized

### Environment Variables
```bash
DDX_TEST_MODE=1                  # Enable test mode
DDX_VERBOSE=true                 # Verbose logging
DDX_LIBRARY_PATH=.ddx/library    # Override path for testing
```

**Note**: Phase detection is artifact-based, not state-file-based. Tests should manipulate actual artifacts (docs, tests) to change phases, not state files.

## Test Execution Plan

### Phase 1: Structure Tests (TS-014-001 to TS-014-004)
- Automated via script
- Runs on every commit
- Fast feedback (<1 second)

### Phase 2: Agent Tests (TS-014-005 to TS-014-007)
- Manual testing initially
- Automated where possible
- Moderate execution time (~30 seconds)

### Phase 3: Functionality Tests (TS-014-008 to TS-014-011)
- Automated integration tests
- Runs before merge
- Execution time (~2 minutes)

### Phase 4: All Other Tests
- Run in CI/CD pipeline
- Before release deployment
- Full execution time (~5 minutes)

## Pass/Fail Criteria

### Must Pass (P0)
All P0 tests must pass before release:
- Structure tests (4 tests)
- Agent activation tests (2 tests)
- Functionality tests (4 tests)
- Phase enforcement tests (2 tests)
- Path resolution tests (2 tests)
- Performance tests (1 test)
- Regression tests (2 tests)

**Total P0 tests: 17**

### Should Pass (P1)
P1 tests should pass but can be addressed post-release:
- Keyword detection (1 test)
- Performance overhead (1 test)
- Error handling tests (2 tests)

**Total P1 tests: 4**

### Overall Success Criteria
- All P0 tests pass: 100% required
- All P1 tests pass: 75% required
- Token reduction achieved: ≥40%
- No user-visible regressions

## Test Reporting

### Test Results Format
```yaml
test_run:
  date: 2025-01-20
  duration: 7m 23s
  total_tests: 21
  passed: 21
  failed: 0
  skipped: 0

metrics:
  line_count_before: 412
  line_count_after: 218
  token_reduction: 47.1%
  agent_activation_success: 100%
  command_success_rate: 100%

status: PASS
```

### Failure Investigation
For each failed test:
1. Capture full error output
2. Identify root cause
3. Create fix or document workaround
4. Retest to verify fix

## Test Maintenance

### Update Triggers
Update tests when:
- CLAUDE.md structure changes
- New workflow commands added
- Agent activation logic changes
- Path resolution changes

### Review Schedule
- After each implementation change
- Before each release
- Quarterly comprehensive review

## Related Documents

- **Feature Spec**: `docs/helix/01-frame/features/FEAT-014-agent-workflow-enforcement.md`
- **Solution Design**: `docs/helix/02-design/solution-designs/SD-014-agent-workflow-enforcement.md`
- **User Story**: `docs/helix/01-frame/user-stories/US-044-workflow-enforcer-agent.md`

---
*Status: Ready for implementation*