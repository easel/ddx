# FEAT-015: Meta-Prompt Injection System

**Feature ID**: FEAT-015
**Feature Name**: Meta-Prompt Automatic Synchronization
**Status**: Specified
**Priority**: P2
**Created**: 2025-01-30
**Updated**: 2025-01-30
**Owner**: Core Team

## Executive Summary

Implement automatic meta-prompt synchronization to CLAUDE.md, mirroring the existing persona injection system. This eliminates manual prompt updates, ensures consistency across projects, and keeps behavioral guidance current with library updates.

## Problem Statement

### Current Issues

1. **Manual Synchronization**: Meta-prompts must be manually copied from library to CLAUDE.md
2. **Inconsistent State**: Projects use different versions of prompts as library evolves
3. **Discovery Problems**: Users don't know when prompts have been updated
4. **Maintenance Burden**: Each project requires manual attention to stay current
5. **Asymmetric Systems**: Personas auto-inject, but meta-prompts don't

### Impact

- **Consistency**: Projects drift from current best practices for AI behavior
- **Maintenance**: Manual work required to keep CLAUDE.md current
- **Quality**: Outdated prompts mean suboptimal AI interactions
- **User Experience**: Users must remember to manually update prompts

## Goals & Non-Goals

### Goals

1. **Automatic Injection**: Meta-prompts inject during `ddx init`
2. **Automatic Sync**: Meta-prompts re-sync during `ddx update`
3. **Health Monitoring**: `ddx doctor` detects out-of-sync prompts
4. **Config Integration**: Changing `system.meta_prompt` triggers re-sync
5. **Pattern Consistency**: Use same approach as persona injection

### Non-Goals

1. Not creating new prompt formats (use existing prompts)
2. Not changing CLAUDE.md structure beyond injection
3. Not modifying library content or prompt files
4. Not adding prompt validation beyond existence checks

## User Stories

- **US-045**: Meta-Prompt Automatic Synchronization

## Solution Overview

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ DDx Commands                                                 │
├─────────────────────────────────────────────────────────────┤
│ ddx init          → MetaPromptInjector.Inject()             │
│ ddx update        → MetaPromptInjector.Sync()               │
│ ddx doctor        → MetaPromptInjector.IsInSync()           │
│ ddx config set    → MetaPromptInjector.Update()             │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ MetaPromptInjector (cli/internal/metaprompt/injector.go)   │
├─────────────────────────────────────────────────────────────┤
│ • InjectMetaPrompt(promptPath string) error                 │
│ • RemoveMetaPrompt() error                                  │
│ • IsInSync() (bool, error)                                  │
│ • GetCurrentMetaPrompt() (string, error)                    │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ CLAUDE.md                                                    │
├─────────────────────────────────────────────────────────────┤
│ [Project content]                                            │
│                                                              │
│ <!-- DDX-META-PROMPT:START -->                              │
│ <!-- Source: claude/system-prompts/focused.md -->           │
│ # System Instructions                                        │
│ [Prompt content from library]                               │
│ <!-- DDX-META-PROMPT:END -->                                │
│                                                              │
│ [More project content]                                       │
└─────────────────────────────────────────────────────────────┘
                            ↑
┌─────────────────────────────────────────────────────────────┐
│ Library (synced via git subtree)                            │
├─────────────────────────────────────────────────────────────┤
│ .ddx/library/prompts/claude/system-prompts/                 │
│ ├── focused.md       (default)                              │
│ ├── strict.md                                               │
│ ├── creative.md                                             │
│ └── ...                                                      │
└─────────────────────────────────────────────────────────────┘
```

### Integration Points

#### 1. Init Command (`cli/cmd/init.go`)
```go
// After creating config and setting up library
if err := injectInitialMetaPrompt(cfg, workingDir); err != nil {
    // Warn but don't fail - meta-prompt is optional
    fmt.Fprintf(os.Stderr, "Warning: Failed to inject meta-prompt: %v\n", err)
}
```

#### 2. Update Command (`cli/cmd/update.go`)
```go
// After library sync (even if no changes)
if err := syncMetaPrompt(cfg, workingDir); err != nil {
    // Warn but don't fail
    fmt.Fprintf(os.Stderr, "Warning: Failed to sync meta-prompt: %v\n", err)
}
```

#### 3. Doctor Command (`cli/cmd/doctor.go`)
```go
// Add health check
checks = append(checks, HealthCheck{
    Name: "Meta-prompt sync",
    Check: func() error {
        return checkMetaPromptSync(cfg, workingDir)
    },
})
```

#### 4. Config Command (`cli/cmd/config.go`)
```go
// After setting system.meta_prompt
if key == "system.meta_prompt" {
    if err := resyncMetaPrompt(cfg, workingDir); err != nil {
        return fmt.Errorf("failed to re-sync meta-prompt: %w", err)
    }
}
```

## Technical Architecture

### Core Components

#### MetaPromptInjector Interface
```go
type MetaPromptInjector interface {
    // InjectMetaPrompt injects a meta-prompt into CLAUDE.md
    InjectMetaPrompt(promptPath string) error

    // RemoveMetaPrompt removes the meta-prompt section from CLAUDE.md
    RemoveMetaPrompt() error

    // IsInSync checks if CLAUDE.md prompt matches library version
    IsInSync() (bool, error)

    // GetCurrentMetaPrompt returns the currently injected prompt info
    GetCurrentMetaPrompt() (string, error)
}
```

#### Implementation Structure
```go
type MetaPromptInjectorImpl struct {
    claudeFilePath string
    libraryPath    string
    workingDir     string
}

// Use same marker pattern as personas
const (
    MetaPromptStartMarker = "<!-- DDX-META-PROMPT:START -->"
    MetaPromptEndMarker   = "<!-- DDX-META-PROMPT:END -->"
)
```

### Data Flow

1. **Read Config**: Get `system.meta_prompt` path (default: "claude/system-prompts/focused.md")
2. **Load Prompt**: Read from `.ddx/library/prompts/{path}`
3. **Find Markers**: Locate `<!-- DDX-META-PROMPT:START/END -->` in CLAUDE.md
4. **Replace Content**: Replace between markers with new prompt + source comment
5. **Write File**: Save updated CLAUDE.md

### Sync Detection Algorithm

```go
func (m *MetaPromptInjectorImpl) IsInSync() (bool, error) {
    // 1. Read CLAUDE.md content between markers
    currentContent, err := m.extractCurrentContent()

    // 2. Extract source path from comment
    sourcePath, err := m.extractSourcePath(currentContent)

    // 3. Read library prompt file
    libraryContent, err := m.readLibraryPrompt(sourcePath)

    // 4. Normalize and compare (ignore whitespace differences)
    return normalize(currentContent) == normalize(libraryContent), nil
}
```

## Implementation Plan

### Phase 1: Frame (Current)
- [x] Create US-045
- [x] Create FEAT-015 (this document)
- [ ] Update feature registry
- [ ] Review and approve requirements

### Phase 2: Design
- [ ] Create SD-015 solution design
- [ ] Define MetaPromptInjector interface details
- [ ] Specify error handling strategy
- [ ] Design sync detection algorithm
- [ ] Document integration points

### Phase 3: Test
- [ ] Write unit tests for injector
- [ ] Write integration tests for commands
- [ ] Write acceptance tests for workflows
- [ ] Create test fixtures (sample prompts)
- [ ] Test marker handling edge cases

### Phase 4: Build
- [ ] Implement `cli/internal/metaprompt/injector.go`
- [ ] Integrate with `init` command
- [ ] Integrate with `update` command
- [ ] Integrate with `doctor` command
- [ ] Integrate with `config` command
- [ ] Add unit tests
- [ ] Add integration tests

### Phase 5: Deploy
- [ ] Update documentation
- [ ] Add migration notes
- [ ] Test on existing projects
- [ ] Monitor for issues

### Phase 6: Iterate
- [ ] Collect user feedback
- [ ] Optimize sync detection
- [ ] Add additional prompts if needed

## Success Metrics

### Primary Metrics
- **Automation**: 100% of projects auto-sync prompts
- **Accuracy**: Doctor reports sync status with 100% accuracy
- **Reliability**: Zero CLAUDE.md corruption incidents
- **Coverage**: Works in all project types (new and existing)

### Secondary Metrics
- **Adoption**: Feature used in 100% of `ddx init` invocations
- **Maintenance**: Zero manual prompt updates needed
- **Consistency**: All projects use current library prompts

## Risks & Mitigation

### Risk 1: CLAUDE.md Corruption
**Likelihood**: Low
**Impact**: High
**Mitigation**:
- Copy robust marker handling from persona system
- Extensive unit tests for edge cases
- Preserve content outside markers
- Test on various CLAUDE.md formats

### Risk 2: Sync Detection False Positives/Negatives
**Likelihood**: Medium
**Impact**: Medium
**Mitigation**:
- Normalize whitespace before comparison
- Hash-based comparison as fallback
- Clear error messages
- Manual override option

### Risk 3: Performance Impact
**Likelihood**: Low
**Impact**: Low
**Mitigation**:
- File operations are fast (~milliseconds)
- Only runs during init/update/config
- Minimal code execution

### Risk 4: Config Migration Issues
**Likelihood**: Low
**Impact**: Medium
**Mitigation**:
- Config schema already supports `system.meta_prompt`
- Default value handles missing config
- Backward compatible

## Dependencies

### Internal
- Persona injection system (reference implementation)
- Config system with `system.meta_prompt` (already exists)
- Library structure with prompts (already exists)
- CLAUDE.md markers (already in place)

### External
- None

## Open Questions

1. **Q**: Should we validate prompt content (e.g., max size, format)?
   **A**: TBD during design phase. Start with existence check only.

2. **Q**: How to handle manual edits to injected prompts?
   **A**: Updates will overwrite. Doctor should warn about drift.

3. **Q**: Should we support multiple prompts in CLAUDE.md?
   **A**: No, single meta-prompt only. Keep it simple.

4. **Q**: What if library doesn't have the configured prompt?
   **A**: Error reported, CLAUDE.md unchanged, suggest available prompts.

## Related Features

- **FEAT-011**: Persona System (reference implementation)
- **FEAT-012**: Library Management System (prompt storage)
- **FEAT-003**: Configuration Management (config integration)

## Documentation Requirements

- Update `README.md` with auto-sync information
- Document `system.meta_prompt` config option
- Add troubleshooting guide for sync issues
- Update architecture documentation

## Acceptance Criteria Summary

From US-045:
- [ ] `ddx init` injects meta-prompt automatically
- [ ] `ddx update` syncs meta-prompt (even without changes)
- [ ] `ddx doctor` reports sync status
- [ ] `ddx config set system.meta_prompt` triggers re-sync
- [ ] Uses same marker pattern as personas
- [ ] Source comment added to injection
- [ ] Preserves CLAUDE.md content outside markers

## Notes

This feature completes the automation of CLAUDE.md management, bringing meta-prompts to parity with personas. It's a natural extension of existing patterns and eliminates a manual maintenance task that users shouldn't have to think about.

The implementation should closely follow the persona injection system architecture, reusing patterns and approaches that are already proven to work reliably.

---
*Status: Awaiting design phase*