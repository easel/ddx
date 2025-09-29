# US-044: Developer Using Workflow Enforcer Agent

**Feature**: FEAT-014 - Agent-Based Workflow Enforcement
**Status**: Draft
**Priority**: P1
**Created**: 2025-01-20
**Updated**: 2025-01-20

## User Story

As a **developer using Claude Code with DDx**,
I want **workflow enforcement handled by a specialized agent system**,
So that **token usage is minimized and enforcement happens reactively when needed**.

## Current Problem

Currently, CLAUDE.md contains ~150 lines of HELIX workflow enforcement instructions including:
- Phase detection logic
- Enforcement principles
- Example violations
- Workflow action usage patterns
- Auto-continuation instructions

These instructions are:
1. **Parsed on every message** - wasting tokens unnecessarily
2. **Passive, not reactive** - always present even when not needed
3. **Mixed with other concerns** - making CLAUDE.md bloated
4. **Harder to maintain** - logic scattered across documentation

## Desired Behavior

### Agent Activation Pattern
```bash
# User asks to work on a story
User: "Work on US-001"

# Claude Code activates workflow agent automatically
Claude: [Launches workflow agent via Task tool]
Agent: [Reads story, detects phase, applies enforcement, executes work]

# No workflow instructions parsed unless workflow work is happening
```

### Minimal CLAUDE.md
CLAUDE.md should contain only:
- Project overview (~50 lines)
- Architecture and key components (~50 lines)
- Development commands (~50 lines)
- Brief pointer to workflow system (~20 lines)
- Architectural principles (~30 lines)
- Testing requirements (~20 lines)
Total: ~220 lines (down from ~400 lines, 45% reduction)

### Workflow Enforcement in Agent
- `./library/workflows/helix/coordinator.md` - Already exists as agent prompt (master repo)
- In user projects: `.ddx/library/workflows/helix/coordinator.md` (via git subtree)
- `library/workflows/helix/actions/*.md` - Action-specific prompts
- Activated via `ddx workflow helix execute` commands
- Phase enforcers loaded dynamically from `library/workflows/helix/phases/*/enforcer.md`

**Note**: This is the master DDx repository, so paths are `./library/*`. User projects have `.ddx/library/*` via git subtree.

## Acceptance Criteria

### AC1: CLAUDE.md Refactored
- [ ] CLAUDE.md reduced to ~220 lines
- [ ] Workflow enforcement sections removed
- [ ] Brief workflow reference section added
- [ ] All essential project info retained
- [ ] Token usage reduced by ~40%

### AC2: Agent System Active
- [ ] Workflow agent activates when `ddx workflow` commands run
- [ ] Phase detection works from agent context
- [ ] Enforcement rules applied from agent
- [ ] No workflow logic in CLAUDE.md

### AC3: Commands Still Work
- [ ] `ddx workflow helix execute build-story US-XXX` works
- [ ] `ddx workflow helix execute continue` works
- [ ] `ddx workflow helix execute status` works
- [ ] `ddx workflow helix execute next` works
- [ ] Phase violations detected and blocked

### AC4: Documentation Complete
- [ ] Agent activation documented in workflow README
- [ ] CLAUDE.md changes documented
- [ ] Token optimization benefits measured

## Technical Notes

### Files to Modify
1. **CLAUDE.md** - Major refactoring
   - Remove: Lines 291-392 (HELIX Workflow Enforcement section)
   - Remove: Lines 260-279 (Auto-Prompts section)
   - Add: Brief workflow reference (~20 lines)

2. **./library/workflows/helix/README.md** - Update with agent usage

3. **./library/workflows/helix/coordinator.md** - Validate as standalone agent

### Token Savings Calculation
- Current CLAUDE.md: ~400 lines × ~50 tokens/line = ~20,000 tokens per message
- Refactored CLAUDE.md: ~220 lines × ~50 tokens/line = ~11,000 tokens per message
- **Savings: ~9,000 tokens per message (45% reduction)**

### Agent Activation Logic
Claude Code should automatically:
1. Detect workflow-related requests
2. Launch Task tool with workflow agent
3. Agent reads coordinator.md as its system prompt
4. Agent detects phase and loads appropriate enforcer
5. Agent executes work with enforcement active

## Testing Strategy

### Manual Testing
1. Execute workflow commands before/after refactoring
2. Verify phase violations still detected
3. Measure token usage reduction
4. Confirm all workflow actions function

### Automated Testing
1. Test CLAUDE.md parsing (validate structure)
2. Test agent activation patterns
3. Test workflow command execution
4. Regression test existing workflow functionality

## Dependencies

- None (internal refactoring only)

## Risks

### Risk 1: Agent Activation Reliability
**Description**: Agent might not activate reliably for workflow tasks
**Mitigation**: Add explicit agent activation patterns, document triggers clearly
**Severity**: Medium

### Risk 2: Functionality Regression
**Description**: Workflow commands might break during refactoring
**Mitigation**: Test all commands before/after, maintain backward compatibility
**Severity**: Low

### Risk 3: Token Savings Not Realized
**Description**: Agent system might use more tokens than saved
**Mitigation**: Measure actual token usage, optimize agent prompts
**Severity**: Low

## Success Metrics

- **Token Reduction**: Achieve 40%+ reduction in base token usage
- **Functionality**: 100% of workflow commands continue working
- **Agent Activation**: Agent activates on 100% of workflow requests
- **Enforcement**: Phase violations detected at same rate as before

## Related Documents

- Feature Spec: `docs/helix/01-frame/features/FEAT-014-agent-workflow-enforcement.md`
- Solution Design: `docs/helix/02-design/solution-designs/SD-014-agent-workflow-enforcement.md`
- Current coordinator: `./library/workflows/helix/coordinator.md`
- Current CLAUDE.md: `CLAUDE.md`

## Notes

This is an architectural improvement that doesn't change user-facing functionality, but significantly improves system efficiency and maintainability. It aligns with the principle that Claude Code should use specialized agents for complex, context-heavy tasks.