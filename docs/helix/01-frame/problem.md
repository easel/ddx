# Problem Statement: HELIX Workflow Auto-Continuation

## Problem Description

Currently, when using HELIX workflow with Claude Code, the workflow stalls between tasks because Claude naturally completes discrete tasks and waits for explicit instructions on what to do next. This breaks the flow of continuous development and forces users to manually prompt Claude for each subsequent action.

## Core Issues

1. **Manual Context Restoration**: Each time Claude finishes a task, the user must manually remind Claude of the workflow state and next steps
2. **Workflow Discontinuity**: Natural stopping points break the HELIX methodology's intended continuous flow
3. **Context Loss**: Claude doesn't retain workflow state between task completions
4. **Phase Awareness Gap**: Claude lacks awareness of current HELIX phase and appropriate next actions

## Impact

- **Reduced Productivity**: Manual intervention required between each workflow step
- **Methodology Violations**: Temptation to skip phases or take shortcuts when flow is broken
- **Cognitive Load**: Users must track workflow state manually instead of focusing on the work
- **Inconsistent Results**: Different users may interpret "next steps" differently

## Success Criteria

When solved, the system should:
- Automatically suggest the next logical action after task completion
- Maintain workflow state and phase awareness across sessions
- Provide seamless continuation without manual intervention
- Enforce HELIX methodology naturally through automation

## Scope

**In Scope**:
- DDx workflow command integration
- Dynamic CLAUDE.md context generation
- Automatic phase progression logic
- Task completion detection

**Out of Scope**:
- Claude Code modification
- Git hook integration (initial version)
- Real-time file watching (initial version)