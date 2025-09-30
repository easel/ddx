# TS-015: Meta-Prompt Injection System Test Specification

**Test Spec ID**: TS-015
**Feature**: FEAT-015 - Meta-Prompt Injection System
**Status**: Draft
**Created**: 2025-01-30
**Updated**: 2025-01-30
**Author**: Core Team

## Overview

This document specifies the test requirements for implementing automatic meta-prompt synchronization to CLAUDE.md, mirroring the persona injection system.

## Test Objectives

1. **Verify Automatic Injection**: Confirm meta-prompts inject on `ddx init`
2. **Validate Automatic Sync**: Ensure meta-prompts sync on `ddx update`
3. **Confirm Sync Detection**: Verify `ddx doctor` detects drift accurately
4. **Check Config Integration**: Ensure config changes trigger re-sync
5. **Prevent Regressions**: No corruption of CLAUDE.md or other content

## Test Categories

### 1. Injection Logic Tests

#### Test 1.1: Initial Injection on Init
```yaml
test_id: TS-015-001
category: injection
priority: P0
description: Verify meta-prompt injected during ddx init

setup:
  - Empty project directory with git initialized
  - Library contains default prompt: claude/system-prompts/focused.md

test_steps:
  - Execute: ddx init
  - Read CLAUDE.md
  - Verify meta-prompt section exists
  - Verify source comment present
  - Verify prompt content matches library

expected_result:
  - CLAUDE.md created with meta-prompt section
  - Content between <!-- DDX-META-PROMPT:START/END --> markers
  - Source: claude/system-prompts/focused.md comment present
  - Prompt content identical to library file

pass_criteria:
  - markers_present == true
  - source_comment_correct == true
  - content_matches_library == true
```

#### Test 1.2: Injection into Existing CLAUDE.md
```yaml
test_id: TS-015-002
category: injection
priority: P0
description: Verify injection preserves existing CLAUDE.md content

setup:
  - Existing CLAUDE.md with project-specific content
  - No meta-prompt section currently present

test_steps:
  - Note existing content
  - Execute: ddx init --force
  - Read CLAUDE.md
  - Verify existing content preserved
  - Verify meta-prompt section added

expected_result:
  - Original content intact
  - Meta-prompt section appended
  - No content loss

pass_criteria:
  - original_content_preserved == true
  - meta_prompt_added == true
  - no_content_corruption == true
```

#### Test 1.3: Replace Existing Meta-Prompt
```yaml
test_id: TS-015-003
category: injection
priority: P0
description: Verify existing meta-prompt replaced on re-injection

setup:
  - CLAUDE.md with old meta-prompt section
  - Library prompt updated with new content

test_steps:
  - Execute: ddx update
  - Read CLAUDE.md
  - Verify old prompt replaced
  - Verify new prompt content present

expected_result:
  - Old meta-prompt removed
  - New meta-prompt injected
  - Single meta-prompt section (no duplicates)

pass_criteria:
  - old_content_removed == true
  - new_content_present == true
  - section_count == 1
```

#### Test 1.4: Injection with Custom Prompt Path
```yaml
test_id: TS-015-004
category: injection
priority: P0
description: Verify injection works with custom meta-prompt path

setup:
  - Config: system.meta_prompt = "claude/system-prompts/strict.md"
  - Library contains strict.md

test_steps:
  - Execute: ddx init
  - Read CLAUDE.md
  - Verify strict.md content injected (not focused.md)

expected_result:
  - Custom prompt injected
  - Source comment shows correct path

pass_criteria:
  - correct_prompt_injected == true
  - source_comment_matches == "claude/system-prompts/strict.md"
```

#### Test 1.5: Injection Disabled
```yaml
test_id: TS-015-005
category: injection
priority: P1
description: Verify no injection when meta_prompt is null

setup:
  - Config: system.meta_prompt = null

test_steps:
  - Execute: ddx init
  - Read CLAUDE.md (if created)
  - Verify no meta-prompt section

expected_result:
  - No meta-prompt section in CLAUDE.md
  - No error thrown

pass_criteria:
  - meta_prompt_section_absent == true
  - init_succeeds == true
```

### 2. Sync Logic Tests

#### Test 2.1: Sync on Update
```yaml
test_id: TS-015-006
category: sync
priority: P0
description: Verify meta-prompt syncs on ddx update

setup:
  - Project with old meta-prompt version
  - Library prompt updated

test_steps:
  - Execute: ddx update
  - Read CLAUDE.md
  - Verify updated prompt content

expected_result:
  - Meta-prompt updated to match library
  - Sync happens even if no git changes pulled

pass_criteria:
  - content_updated == true
  - matches_library == true
```

#### Test 2.2: Sync Without Library Changes
```yaml
test_id: TS-015-007
category: sync
priority: P0
description: Verify sync happens even without library changes

setup:
  - Project with out-of-sync meta-prompt
  - No pending library updates

test_steps:
  - Execute: ddx update
  - Read CLAUDE.md
  - Verify meta-prompt updated

expected_result:
  - Meta-prompt synced to library version
  - Update command succeeds

pass_criteria:
  - sync_occurred == true
  - command_succeeds == true
```

#### Test 2.3: Sync After Config Change
```yaml
test_id: TS-015-008
category: sync
priority: P0
description: Verify sync after changing meta_prompt config

setup:
  - Project with focused.md injected
  - Config: system.meta_prompt = "claude/system-prompts/focused.md"

test_steps:
  - Execute: ddx config set system.meta_prompt "claude/system-prompts/strict.md"
  - Read CLAUDE.md
  - Verify strict.md content now present

expected_result:
  - Old prompt removed
  - New prompt injected
  - Source comment updated

pass_criteria:
  - prompt_changed == true
  - source_comment_updated == true
```

#### Test 2.4: Remove on Disable
```yaml
test_id: TS-015-009
category: sync
priority: P1
description: Verify meta-prompt removed when disabled

setup:
  - Project with meta-prompt section

test_steps:
  - Execute: ddx config set system.meta_prompt null
  - Read CLAUDE.md
  - Verify meta-prompt section removed

expected_result:
  - Meta-prompt section absent
  - Rest of CLAUDE.md intact

pass_criteria:
  - section_removed == true
  - other_content_preserved == true
```

### 3. Sync Detection Tests

#### Test 3.1: Detect In Sync
```yaml
test_id: TS-015-010
category: doctor
priority: P0
description: Verify doctor reports in-sync when current

setup:
  - Project with up-to-date meta-prompt

test_steps:
  - Execute: ddx doctor
  - Check meta-prompt sync status

expected_result:
  - Status: "ok"
  - Message: "Meta-prompt is in sync with library"

pass_criteria:
  - status == "ok"
  - in_sync_reported == true
```

#### Test 3.2: Detect Out of Sync
```yaml
test_id: TS-015-011
category: doctor
priority: P0
description: Verify doctor detects out-of-sync meta-prompt

setup:
  - Project with old meta-prompt version
  - Library prompt has been updated

test_steps:
  - Execute: ddx doctor
  - Check meta-prompt sync status

expected_result:
  - Status: "warning"
  - Message: "Meta-prompt is out of sync with library"
  - Fix: "Run 'ddx update' to sync meta-prompt"

pass_criteria:
  - status == "warning"
  - out_of_sync_detected == true
  - fix_suggested == true
```

#### Test 3.3: Detect Missing Meta-Prompt
```yaml
test_id: TS-015-012
category: doctor
priority: P1
description: Verify doctor handles missing meta-prompt section

setup:
  - CLAUDE.md exists but has no meta-prompt section

test_steps:
  - Execute: ddx doctor
  - Check meta-prompt sync status

expected_result:
  - Status: "warning"
  - Message indicates meta-prompt not found
  - Suggests running ddx update

pass_criteria:
  - status == "warning"
  - missing_detected == true
```

#### Test 3.4: Whitespace Tolerance
```yaml
test_id: TS-015-013
category: doctor
priority: P1
description: Verify sync detection ignores whitespace differences

setup:
  - CLAUDE.md with meta-prompt (extra whitespace in content)
  - Library prompt (normal whitespace)
  - Content semantically identical

test_steps:
  - Execute: ddx doctor
  - Check sync status

expected_result:
  - Status: "ok" (whitespace differences ignored)

pass_criteria:
  - status == "ok"
  - whitespace_normalized == true
```

### 4. Marker Handling Tests

#### Test 4.1: Correct Marker Placement
```yaml
test_id: TS-015-014
category: markers
priority: P0
description: Verify markers placed correctly

setup:
  - Empty CLAUDE.md

test_steps:
  - Inject meta-prompt
  - Verify START marker before content
  - Verify END marker after content
  - Verify source comment between START marker and content

expected_result:
  - Markers in correct positions
  - Content wrapped properly

pass_criteria:
  - start_marker_present == true
  - end_marker_present == true
  - source_comment_present == true
  - markers_ordered_correctly == true
```

#### Test 4.2: Malformed Section Recovery
```yaml
test_id: TS-015-015
category: markers
priority: P1
description: Verify recovery from malformed meta-prompt section

setup:
  - CLAUDE.md with START marker but missing END marker

test_steps:
  - Execute: ddx update
  - Verify section cleaned up
  - Verify valid section injected

expected_result:
  - Malformed section removed
  - Valid section created

pass_criteria:
  - malformed_section_removed == true
  - valid_section_created == true
  - no_corruption == true
```

#### Test 4.3: Multiple Sections Prevention
```yaml
test_id: TS-015-016
category: markers
priority: P0
description: Verify only single meta-prompt section exists

setup:
  - Attempt to inject multiple times

test_steps:
  - Inject meta-prompt
  - Inject again
  - Count meta-prompt sections

expected_result:
  - Only one section exists
  - Second injection replaces first

pass_criteria:
  - section_count == 1
  - no_duplicates == true
```

### 5. Error Handling Tests

#### Test 5.1: Prompt File Not Found
```yaml
test_id: TS-015-017
category: error_handling
priority: P0
description: Verify graceful handling when prompt file missing

setup:
  - Config: system.meta_prompt = "claude/system-prompts/nonexistent.md"

test_steps:
  - Execute: ddx init
  - Verify error message helpful
  - Verify init continues (with warning)

expected_result:
  - Error message: "failed to read meta-prompt from {path}"
  - Warning shown to user
  - Init completes successfully

pass_criteria:
  - error_message_clear == true
  - init_succeeds == true
  - warning_displayed == true
```

#### Test 5.2: CLAUDE.md Read-Only
```yaml
test_id: TS-015-018
category: error_handling
priority: P1
description: Verify error when CLAUDE.md is read-only

setup:
  - CLAUDE.md exists with read-only permissions (444)

test_steps:
  - Execute: ddx update
  - Verify error message helpful

expected_result:
  - Error message mentions permission issue
  - Suggests checking file permissions

pass_criteria:
  - error_mentions_permissions == true
  - helpful_guidance == true
```

#### Test 5.3: Prompt Too Large
```yaml
test_id: TS-015-019
category: error_handling
priority: P1
description: Verify handling of oversized prompt file

setup:
  - Create prompt file > 512KB

test_steps:
  - Attempt to inject
  - Verify size limit enforced

expected_result:
  - Error: "meta-prompt too large: {size} bytes (max {max})"
  - Injection fails gracefully

pass_criteria:
  - size_limit_enforced == true
  - error_clear == true
```

#### Test 5.4: Library Directory Missing
```yaml
test_id: TS-015-020
category: error_handling
priority: P0
description: Verify handling when .ddx/library missing

setup:
  - Project initialized but library not synced

test_steps:
  - Execute: ddx update
  - Verify error message helpful

expected_result:
  - Error indicates library not found
  - Suggests running: ddx update

pass_criteria:
  - error_clear == true
  - recovery_suggested == true
```

### 6. Content Preservation Tests

#### Test 6.1: Preserve Content Before Section
```yaml
test_id: TS-015-021
category: content_preservation
priority: P0
description: Verify content before meta-prompt section preserved

setup:
  - CLAUDE.md with content before meta-prompt section

test_steps:
  - Note content before section
  - Execute: ddx update
  - Verify content still present and unchanged

expected_result:
  - Content before section intact
  - No modifications to original content

pass_criteria:
  - content_preserved == true
  - no_modifications == true
```

#### Test 6.2: Preserve Content After Section
```yaml
test_id: TS-015-022
category: content_preservation
priority: P0
description: Verify content after meta-prompt section preserved

setup:
  - CLAUDE.md with content after meta-prompt section

test_steps:
  - Note content after section
  - Execute: ddx update
  - Verify content still present and unchanged

expected_result:
  - Content after section intact
  - Proper spacing maintained

pass_criteria:
  - content_preserved == true
  - spacing_correct == true
```

#### Test 6.3: Preserve Personas Section
```yaml
test_id: TS-015-023
category: content_preservation
priority: P0
description: Verify persona section unaffected by meta-prompt operations

setup:
  - CLAUDE.md with both personas and meta-prompt sections

test_steps:
  - Note personas section content
  - Execute: ddx update
  - Verify personas section unchanged

expected_result:
  - Personas section completely intact
  - No cross-contamination

pass_criteria:
  - personas_preserved == true
  - no_interference == true
```

### 7. Integration Tests

#### Test 7.1: Full Lifecycle - Init to Update
```yaml
test_id: TS-015-024
category: integration
priority: P0
description: Test complete lifecycle from init through update

test_steps:
  - Execute: ddx init
  - Verify meta-prompt injected
  - Manually edit library prompt
  - Execute: ddx update
  - Verify meta-prompt updated

expected_result:
  - Init injects prompt
  - Update syncs changes
  - Full workflow functions

pass_criteria:
  - init_injection_success == true
  - update_sync_success == true
  - no_errors == true
```

#### Test 7.2: Config Change Workflow
```yaml
test_id: TS-015-025
category: integration
priority: P0
description: Test config change triggers re-sync

test_steps:
  - Execute: ddx init (default focused.md)
  - Verify focused.md injected
  - Execute: ddx config set system.meta_prompt "claude/system-prompts/strict.md"
  - Verify strict.md injected
  - Execute: ddx config set system.meta_prompt null
  - Verify meta-prompt removed

expected_result:
  - Each config change triggers appropriate action
  - Content updates correctly

pass_criteria:
  - initial_injection_success == true
  - config_change_triggers_sync == true
  - removal_works == true
```

#### Test 7.3: Doctor → Update → Doctor
```yaml
test_id: TS-015-026
category: integration
priority: P0
description: Test doctor detects issue, update fixes, doctor confirms

test_steps:
  - Create out-of-sync condition
  - Execute: ddx doctor (should warn)
  - Execute: ddx update
  - Execute: ddx doctor (should be ok)

expected_result:
  - First doctor detects issue
  - Update resolves issue
  - Second doctor confirms fix

pass_criteria:
  - initial_warning == true
  - update_fixes == true
  - final_ok == true
```

### 8. Performance Tests

#### Test 8.1: Injection Speed
```yaml
test_id: TS-015-027
category: performance
priority: P1
description: Measure meta-prompt injection performance

test_steps:
  - Time ddx init execution
  - Isolate meta-prompt injection time
  - Verify acceptable performance

expected_result:
  - Injection time < 15ms (typical)
  - No noticeable slowdown

pass_criteria:
  - injection_time_ms < 50
  - acceptable_performance == true
```

#### Test 8.2: Sync Detection Speed
```yaml
test_id: TS-015-028
category: performance
priority: P1
description: Measure sync detection performance

test_steps:
  - Time ddx doctor execution
  - Isolate sync check time
  - Verify acceptable performance

expected_result:
  - Sync check time < 10ms (typical)
  - No noticeable slowdown

pass_criteria:
  - check_time_ms < 25
  - acceptable_performance == true
```

#### Test 8.3: Memory Usage
```yaml
test_id: TS-015-029
category: performance
priority: P1
description: Verify memory usage is acceptable

test_steps:
  - Monitor memory during injection
  - Monitor memory during sync check

expected_result:
  - Peak memory < 1MB additional
  - No memory leaks

pass_criteria:
  - memory_usage_mb < 1
  - no_leaks == true
```

### 9. Edge Case Tests

#### Test 9.1: Empty Prompt File
```yaml
test_id: TS-015-030
category: edge_cases
priority: P1
description: Handle empty prompt file gracefully

setup:
  - Create empty prompt file in library

test_steps:
  - Attempt injection
  - Verify handling

expected_result:
  - Empty section created (or error, depending on design decision)
  - No crash

pass_criteria:
  - no_crash == true
  - behavior_defined == true
```

#### Test 9.2: Very Large CLAUDE.md
```yaml
test_id: TS-015-031
category: edge_cases
priority: P1
description: Handle large CLAUDE.md files

setup:
  - Create CLAUDE.md with 500KB of content

test_steps:
  - Execute injection
  - Verify performance acceptable
  - Verify content preserved

expected_result:
  - Injection succeeds
  - Performance acceptable (<100ms)

pass_criteria:
  - injection_succeeds == true
  - performance_acceptable == true
```

#### Test 9.3: Special Characters in Prompt
```yaml
test_id: TS-015-032
category: edge_cases
priority: P1
description: Handle prompts with special characters

setup:
  - Prompt contains: <, >, &, quotes, markdown markers

test_steps:
  - Inject prompt
  - Read back
  - Verify special characters preserved

expected_result:
  - All characters preserved correctly
  - No escaping issues

pass_criteria:
  - content_identical == true
  - no_corruption == true
```

### 10. Regression Tests

#### Test 10.1: Persona System Unaffected
```yaml
test_id: TS-015-033
category: regression
priority: P0
description: Verify persona injection still works

test_steps:
  - Execute: ddx persona bind code-reviewer strict-code-reviewer
  - Execute: ddx update (trigger meta-prompt sync)
  - Verify personas section intact
  - Verify personas still load

expected_result:
  - Persona system functions normally
  - No interference from meta-prompt system

pass_criteria:
  - personas_work == true
  - no_interference == true
```

#### Test 10.2: Existing Tests Pass
```yaml
test_id: TS-015-034
category: regression
priority: P0
description: Verify no existing tests break

test_steps:
  - Run full test suite
  - Compare results before/after feature

expected_result:
  - All previously passing tests still pass
  - No new failures

pass_criteria:
  - no_new_failures == true
  - test_count_unchanged == true
```

## Test Data Requirements

### Test Prompts
Create test prompt files:
- `test-prompts/minimal.md`: Minimal prompt (1 line)
- `test-prompts/standard.md`: Standard prompt (~20 lines)
- `test-prompts/large.md`: Large prompt (400 lines)
- `test-prompts/special-chars.md`: Prompt with special characters

### Test CLAUDE.md Files
- `test-claude/empty.md`: Empty file
- `test-claude/basic.md`: Basic content, no meta-prompt
- `test-claude/with-meta.md`: Content with meta-prompt section
- `test-claude/with-personas.md`: Content with personas section
- `test-claude/malformed.md`: Malformed meta-prompt section

### Test Config Files
```yaml
# Default config
system:
  meta_prompt: "claude/system-prompts/focused.md"

# Custom prompt
system:
  meta_prompt: "claude/system-prompts/strict.md"

# Disabled
system:
  meta_prompt: null
```

## Test Environment Setup

### Prerequisites
1. DDx CLI installed
2. Test project with .ddx/config.yaml
3. .ddx/library synced with test prompts
4. Git initialized (for init tests)

### Environment Variables
```bash
DDX_TEST_MODE=1                  # Enable test mode
DDX_VERBOSE=true                 # Verbose logging
DDX_LIBRARY_PATH=.ddx/library    # Override path for testing
```

### Test Helpers
```go
// setupTestEnvironment creates isolated test environment
func setupTestEnvironment(t *testing.T) *TestEnv

// setupTestLibrary creates library with test prompts
func setupTestLibrary(t *testing.T, testEnv *TestEnv)

// createTestCLAUDE creates CLAUDE.md with specific content
func createTestCLAUDE(t *testing.T, testEnv *TestEnv, content string)

// verifyMetaPromptSection checks meta-prompt section validity
func verifyMetaPromptSection(t *testing.T, claudePath string) bool

// comparePromptContent compares CLAUDE.md prompt with library
func comparePromptContent(t *testing.T, claudePath, libraryPromptPath string) bool
```

## Test Execution Plan

### Phase 1: Unit Tests (TS-015-001 to TS-015-020)
- Automated via Go test framework
- Runs on every commit
- Fast feedback (<5 seconds)
- Focus: Injection logic, marker handling, error cases

### Phase 2: Integration Tests (TS-015-021 to TS-015-026)
- Automated integration tests
- Runs before merge
- Execution time (~30 seconds)
- Focus: Command integration, workflows

### Phase 3: Performance Tests (TS-015-027 to TS-015-029)
- Automated benchmarks
- Runs in CI/CD pipeline
- Execution time (~10 seconds)
- Focus: Speed, memory usage

### Phase 4: Edge Cases & Regression (TS-015-030 to TS-015-034)
- Mix of automated and manual
- Runs before release
- Execution time (~1 minute)
- Focus: Edge cases, no regressions

## Pass/Fail Criteria

### Must Pass (P0)
All P0 tests must pass before release:
- Injection tests (4 tests)
- Sync tests (3 tests)
- Doctor tests (2 tests)
- Marker tests (2 tests)
- Error handling tests (2 tests)
- Content preservation tests (3 tests)
- Integration tests (3 tests)
- Regression tests (2 tests)

**Total P0 tests: 21**

### Should Pass (P1)
P1 tests should pass but can be addressed post-release:
- Injection disabled test (1 test)
- Sync removal test (1 test)
- Doctor missing test (1 test)
- Doctor whitespace test (1 test)
- Malformed recovery test (1 test)
- Permission error test (1 test)
- Size limit test (1 test)
- Performance tests (3 tests)
- Edge case tests (3 tests)

**Total P1 tests: 13**

### Overall Success Criteria
- All P0 tests pass: 100% required
- All P1 tests pass: 80% required
- No CLAUDE.md corruption: 0 incidents
- No user-visible regressions
- Performance acceptable: <50ms injection

## Test Reporting

### Test Results Format
```yaml
test_run:
  date: 2025-01-30
  duration: 1m 45s
  total_tests: 34
  passed: 34
  failed: 0
  skipped: 0

metrics:
  injection_avg_ms: 12.3
  sync_check_avg_ms: 8.7
  doctor_check_avg_ms: 9.2
  claude_corruption_count: 0

test_coverage:
  unit: 95.2%
  integration: 87.6%

status: PASS
```

### Failure Investigation
For each failed test:
1. Capture full error output
2. Capture CLAUDE.md state before/after
3. Identify root cause
4. Create fix or document workaround
5. Retest to verify fix

## Test Maintenance

### Update Triggers
Update tests when:
- Meta-prompt injection logic changes
- Config schema changes
- New prompt types added
- Marker handling changes
- Doctor check logic updates

### Review Schedule
- After each implementation change
- Before each release
- Quarterly comprehensive review

## Test Automation

### CI/CD Integration
```yaml
# .github/workflows/test-meta-prompt-injection.yml
name: Meta-Prompt Injection Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up Go
        uses: actions/setup-go@v2
      - name: Run unit tests
        run: go test ./internal/metaprompt/...
      - name: Run integration tests
        run: go test ./cmd/... -run "TestMetaPrompt"
      - name: Run performance tests
        run: go test -bench=. ./internal/metaprompt/...
```

### Pre-commit Hooks
```bash
# Run meta-prompt tests before commit
go test ./internal/metaprompt/... -v
```

## Related Documents

- **Feature Spec**: `docs/helix/01-frame/features/FEAT-015-meta-prompt-injection.md`
- **Solution Design**: `docs/helix/02-design/solution-designs/SD-015-meta-prompt-injection.md`
- **User Story**: `docs/helix/01-frame/user-stories/US-045-meta-prompt-auto-sync.md`
- **Reference Implementation**: `cli/internal/persona/claude.go`

---
*Status: Ready for implementation*