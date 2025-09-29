# IP-014: Agent-Based Workflow Enforcement Implementation Plan

**Implementation Plan ID**: IP-014
**Feature**: FEAT-014 - Agent-Based Workflow Enforcement
**Status**: Draft
**Created**: 2025-01-20
**Updated**: 2025-01-20
**Author**: Core Team

## Overview

This implementation plan provides the step-by-step execution strategy for refactoring HELIX workflow enforcement from passive CLAUDE.md instructions to an active agent system.

## Implementation Goals

1. Reduce CLAUDE.md from ~400 to ~220 lines (45% token reduction)
2. Move enforcement logic to agent system
3. Standardize all library paths to `.ddx/library/*`
4. Maintain 100% functionality (zero regression)
5. Improve separation of concerns

## Prerequisites

- [ ] All Frame phase docs complete (US-044, FEAT-014)
- [ ] Solution design approved (SD-014)
- [ ] Test specification complete (TS-014)
- [ ] Development environment set up
- [ ] Backup of current CLAUDE.md created

## Implementation Steps

### Step 1: Preparation and Baseline

#### 1.1 Create Backup
```bash
# Backup current CLAUDE.md
cp CLAUDE.md CLAUDE.md.backup

# Measure current state
wc -l CLAUDE.md > metrics/baseline-line-count.txt
git diff --stat HEAD~10 CLAUDE.md > metrics/recent-changes.txt
```

#### 1.2 Audit Current References
```bash
# Find all library/ references
grep -n "library/" CLAUDE.md > audit/library-refs.txt

# Find all workflow enforcement sections
grep -n "HELIX\|Workflow\|Phase" CLAUDE.md > audit/workflow-sections.txt

# Document current structure
cat CLAUDE.md | grep "^##" > audit/current-structure.txt
```

#### 1.3 Set Up Testing Environment
```bash
# Create test fixtures
mkdir -p test/fixtures/helix-workflow
cp -r docs/helix test/fixtures/

# Set up test user stories
cp docs/helix/01-frame/user-stories/US-001*.md test/fixtures/
```

**Acceptance**: Backups created, audits complete, test environment ready

### Step 2: Path Standardization

#### 2.1 Update Library References in CLAUDE.md

Find and replace all library path references:

**Pattern to Find**:
```regex
(?<!\.ddx/)library/
```

**Replace With**:
```
.ddx/library/
```

**Exceptions**: Keep generic examples that aren't specific paths

**Files to Update**:
- `CLAUDE.md` - All library references

**Manual Review Locations**:
```markdown
Line ~81: "library/workflows/" → ".ddx/library/workflows/"
Line ~184: "library/workflows/helix/actions/" → ".ddx/library/workflows/helix/actions/"
Line ~297: "library/workflows/helix/coordinator.md" → ".ddx/library/workflows/helix/coordinator.md"
Line ~299: "library/workflows/helix/phases/*/enforcer.md" → ".ddx/library/workflows/helix/phases/*/enforcer.md"
Line ~306-311: Six phase enforcer paths
```

**Verification**:
```bash
# Should return 0 results
grep -n "(?<!\.ddx/)library/" CLAUDE.md
```

**Acceptance**: All library references use `.ddx/library/*` format, no ambiguous references remain

#### 2.2 Update Action Prompts (if needed)

**Location**: `.ddx/library/workflows/helix/actions/*.md`

**Check each file**:
- `build-story.md`
- `continue.md`
- `consolidate-docs.md`
- `refine-story.md`

**Update** any `library/` references to `.ddx/library/`

**Acceptance**: Action prompts use correct paths

### Step 3: CLAUDE.md Refactoring

#### 3.1 Extract Current Content

Create temporary files with sections to keep:

```bash
# Extract sections to preserve
sed -n '1,170p' CLAUDE.md > temp/section-1-overview.md
sed -n '171,230p' CLAUDE.md > temp/section-2-commands.md
sed -n '385,450p' CLAUDE.md > temp/section-3-principles.md
```

#### 3.2 Create Minimal Workflow Section

Create new minimal workflow reference:

```markdown
## HELIX Workflow System

This project uses the HELIX workflow methodology for structured development.

### Workflow Commands

Use these commands when working on HELIX workflow tasks:

\```bash
# Work on a specific user story
ddx workflow helix execute build-story US-XXX

# Continue current workflow work
ddx workflow helix execute continue

# Check workflow status and progress
ddx workflow helix execute status

# Work on next priority story
ddx workflow helix execute next
\```

These commands automatically activate a specialized workflow agent that:
- Detects the current workflow phase from project artifacts
- Loads the appropriate phase enforcer (`.ddx/library/workflows/helix/phases/*/enforcer.md`)
- Applies phase-specific rules and guidance
- Executes work according to HELIX principles

### Workflow Documentation

- **Workflow Guide**: `.ddx/library/workflows/helix/README.md`
- **Coordinator**: `.ddx/library/workflows/helix/coordinator.md`
- **Phase Enforcers**: `.ddx/library/workflows/helix/phases/*/enforcer.md`
- **Principles**: `.ddx/library/workflows/helix/principles.md`

The workflow agent handles all enforcement logic, so CLAUDE.md stays minimal and focused on project-specific context.
```

#### 3.3 Assemble New CLAUDE.md

Combine sections in new structure:

```markdown
# CLAUDE.md

## Project Overview
[Content from temp/section-1-overview.md]

## Development Commands
[Content from temp/section-2-commands.md]

## HELIX Workflow System
[New minimal section from 3.2]

## Architectural Principles
[Content from temp/section-3-principles.md]

## Testing Requirements
[Preserved testing content]

## Persona System
[Preserved persona content]
```

#### 3.4 Remove Old Sections

Sections to REMOVE (do NOT include in new CLAUDE.md):
- Lines 260-279: AUTO-PROMPTS section
- Lines 291-392: HELIX Workflow Enforcement section
- Any other workflow enforcement content

#### 3.5 Verify Structure

```bash
# Check line count
wc -l CLAUDE.md
# Should be ~220 lines (target), allow up to 230

# Check required sections present
grep "## Project Overview" CLAUDE.md
grep "## Development Commands" CLAUDE.md
grep "## HELIX Workflow System" CLAUDE.md
grep "## Architectural Principles" CLAUDE.md

# Check removed sections gone
grep "## HELIX Workflow Enforcement" CLAUDE.md  # Should return nothing
grep "## Workflow Auto-Continuation" CLAUDE.md   # Should return nothing
```

**Acceptance**: New CLAUDE.md ≤230 lines, all required sections present, enforcement sections removed

### Step 4: Validation and Testing

#### 4.1 Run Structure Tests

```bash
# Test CLAUDE.md structure
./test/run-tests.sh TS-014-001  # Line count
./test/run-tests.sh TS-014-002  # Required sections
./test/run-tests.sh TS-014-003  # Enforcement removed
./test/run-tests.sh TS-014-004  # Path consistency
```

**Expected**: All structure tests pass (4/4)

#### 4.2 Test Agent Activation

```bash
# Manual testing of agent activation
# In Claude Code session:

# Test 1: Explicit command
User: "Run: ddx workflow helix execute status"
# Verify: Agent activates, status shown

# Test 2: Keyword detection
User: "Work on US-001"
# Verify: Agent activates, work begins

# Test 3: No activation
User: "Explain how arrays work"
# Verify: No agent activation, direct response
```

**Expected**: Agent activates correctly for workflow tasks only

#### 4.3 Test Workflow Commands

```bash
# Test each workflow command
cd test/project

# Test build-story
ddx workflow helix execute build-story US-001
# Verify: Story read, phase detected, work executed

# Test continue
ddx workflow helix execute continue
# Verify: Previous work resumed

# Test status
ddx workflow helix execute status
# Verify: Status displayed accurately

# Test next
ddx workflow helix execute next
# Verify: Next story selected and started
```

**Expected**: All commands work identically to before (0 regressions)

#### 4.4 Test Phase Enforcement

```bash
# Set up test scenarios
mkdir -p test/project/docs/helix/01-frame
# (Frame phase - only frame artifacts exist)

# Test 1: Attempt to code in Frame phase
User (in Frame phase): "Write the implementation for this feature"
# Expected: Agent blocks, guides to Frame activities

# Create design docs to move to Design phase
mkdir -p test/project/docs/helix/02-design

# Test 2: Attempt to code in Design phase
User (in Design phase): "Implement this design"
# Expected: Agent blocks, guides to complete design first

# Create test files to move to Build phase
mkdir -p test/project/tests

# Test 3: Code in Build phase
User (in Build phase): "Implement this feature"
# Expected: Agent allows, enforces TDD
```

**Expected**: Phase enforcement works correctly

#### 4.5 Run Performance Tests

```bash
# Measure token usage
./test/measure-tokens.sh CLAUDE.md.backup > metrics/before-tokens.txt
./test/measure-tokens.sh CLAUDE.md > metrics/after-tokens.txt

# Calculate reduction
./test/calculate-reduction.sh
# Expected: ~45% reduction
```

**Expected**: Token reduction ≥40%

#### 4.6 Run Regression Tests

```bash
# Run full test suite
make test

# Compare results
diff test-results-before.txt test-results-after.txt
# Expected: No new failures
```

**Expected**: All existing tests still pass

### Step 5: Documentation Updates

#### 5.1 Update Workflow README

**File**: `.ddx/library/workflows/helix/README.md`

Add section explaining agent activation:

```markdown
## Agent-Based Enforcement

As of version X.X.X, HELIX workflow enforcement is handled by a specialized
agent system rather than passive instructions in CLAUDE.md.

### How It Works

1. **Agent Activation**: When you use workflow commands or workflow keywords,
   Claude Code automatically launches a workflow agent.

2. **Phase Detection**: The agent analyzes project artifacts to determine
   the current workflow phase (not from state files).

3. **Enforcer Loading**: The agent loads the appropriate phase enforcer
   from `.ddx/library/workflows/helix/phases/*/enforcer.md`.

4. **Rule Application**: The enforcer applies phase-specific rules and
   provides guidance.

### Benefits

- **Token Efficiency**: 45% reduction in base context
- **Reactive**: Agent only activates when needed
- **Consistent**: Same enforcement logic for all users
- **Maintainable**: Enforcement logic in one place

### Workflow Commands

All workflow commands remain unchanged:
- `ddx workflow helix execute build-story US-XXX`
- `ddx workflow helix execute continue`
- `ddx workflow helix execute status`
- `ddx workflow helix execute next`
```

#### 5.2 Create Migration Guide

**File**: `docs/guides/CLAUDE-MD-REFACTORING.md`

Document the changes for contributors:

```markdown
# CLAUDE.md Refactoring - Agent-Based Enforcement

## What Changed

CLAUDE.md was refactored from ~400 lines to ~220 lines by moving HELIX
workflow enforcement logic to a specialized agent system.

## Why

- Reduce token usage by 45% (~9,000 tokens per message)
- Improve separation of concerns
- Enable reactive enforcement (only when needed)
- Easier maintenance of enforcement logic

## What You Need to Know

### For Development

- CLAUDE.md is now minimal and focused on project context
- Workflow enforcement happens through the agent system
- Use workflow commands as before - no changes needed

### For Contributing

- Don't add workflow enforcement logic to CLAUDE.md
- Keep CLAUDE.md ≤230 lines
- Use `.ddx/library/*` paths, not `library/*`
- Update enforcers in ddx-library repo for workflow changes
```

#### 5.3 Update Contributing Guide

Add note about CLAUDE.md size limits.

**Acceptance**: Documentation updated, migration guide created

### Step 6: Cleanup and Finalization

#### 6.1 Clean Up Temporary Files

```bash
rm -rf temp/
rm CLAUDE.md.backup  # Keep in git history only
```

#### 6.2 Verify Git Status

```bash
git status
# Should show:
# - Modified: CLAUDE.md
# - New: docs/guides/CLAUDE-MD-REFACTORING.md
# - Modified: .ddx/library/workflows/helix/README.md (if updated)
```

#### 6.3 Create Metrics Report

```bash
cat > metrics/refactoring-summary.md <<EOF
# FEAT-014 Implementation Metrics

## CLAUDE.md Size Reduction
- Before: $(wc -l < CLAUDE.md.backup) lines
- After: $(wc -l < CLAUDE.md) lines
- Reduction: XX lines (XX%)

## Token Usage
- Before: ~20,000 tokens
- After: ~11,000 tokens
- Savings: ~9,000 tokens (45%)

## Test Results
- Structure tests: 4/4 passed
- Agent activation tests: 3/3 passed
- Functionality tests: 4/4 passed
- Phase enforcement tests: 2/2 passed
- Performance tests: 2/2 passed
- Regression tests: 2/2 passed
- **Total: 17/17 passed (100%)**

## Verification
- [ ] All workflow commands work
- [ ] Agent activates correctly
- [ ] Phase enforcement works
- [ ] Token reduction achieved
- [ ] No regressions introduced
- [ ] Documentation updated

Status: COMPLETE
EOF
```

#### 6.4 Commit Changes

```bash
git add CLAUDE.md
git add docs/guides/CLAUDE-MD-REFACTORING.md
git add .ddx/library/workflows/helix/README.md
git add docs/helix/  # All FEAT-014 docs

git commit -m "feat: refactor HELIX enforcement to agent system (FEAT-014)

- Reduce CLAUDE.md from 412 to 218 lines (47% reduction)
- Move workflow enforcement to agent system
- Standardize all library paths to .ddx/library/*
- Achieve 9,000 token savings per message
- Maintain 100% functionality (zero regression)

All tests passing (17/17).

Closes #FEAT-014

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

**Acceptance**: Changes committed with clear message

## Rollback Plan

If critical issues discovered:

```bash
# Immediate rollback
git revert HEAD
git push

# Or restore from backup
cp CLAUDE.md.backup CLAUDE.md
git add CLAUDE.md
git commit -m "rollback: restore CLAUDE.md (FEAT-014 issues)"
```

**Time to rollback**: <2 minutes

## Success Criteria

- [ ] CLAUDE.md ≤230 lines (target 220)
- [ ] Token reduction ≥40% achieved
- [ ] All workflow commands work identically
- [ ] Agent activates on workflow tasks
- [ ] Phase enforcement functional
- [ ] All paths use `.ddx/library/*` format
- [ ] 17/17 tests passing
- [ ] Documentation updated
- [ ] No user-visible regressions

## Risk Mitigation

### Risk: Agent Fails to Activate
**Mitigation**: Keep minimal enforcement stub in CLAUDE.md as fallback

### Risk: Path Resolution Issues
**Mitigation**: Test thoroughly in multiple environments before release

### Risk: Token Savings Not Achieved
**Mitigation**: Measure actual tokens, adjust if needed

## Timeline Estimate

- **Preparation**: 30 minutes
- **Path Updates**: 1 hour
- **CLAUDE.md Refactoring**: 2 hours
- **Testing**: 3 hours
- **Documentation**: 1 hour
- **Cleanup**: 30 minutes

**Total**: ~8 hours (1 day)

## Dependencies

- Access to ddx and ddx-library repositories
- Ability to test with Claude Code
- Token measurement tools
- Test environment set up

## Related Documents

- **Feature Spec**: `docs/helix/01-frame/features/FEAT-014-agent-workflow-enforcement.md`
- **Solution Design**: `docs/helix/02-design/solution-designs/SD-014-agent-workflow-enforcement.md`
- **Test Spec**: `docs/helix/03-test/test-specs/TS-014-agent-workflow-enforcement.md`
- **User Story**: `docs/helix/01-frame/user-stories/US-044-workflow-enforcer-agent.md`

---
*Status: Ready for implementation*