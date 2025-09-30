# US-045: Meta-Prompt Automatic Synchronization

**Feature**: FEAT-015 - Meta-Prompt Injection System
**Status**: Draft
**Priority**: P2
**Created**: 2025-01-30
**Updated**: 2025-01-30

## User Story

As a **DDx user**,
I want **meta-prompts to automatically sync to CLAUDE.md**,
So that **I always have the latest behavioral guidance without manual updates**.

## Current Problem

Currently, meta-prompts (behavioral instructions from `.ddx/library/prompts/claude/system-prompts/`) must be manually copied into CLAUDE.md between the `<!-- DDX-META-PROMPT:START -->` and `<!-- DDX-META-PROMPT:END -->` markers. This leads to:

1. **Out-of-sync prompts** - Library updates don't reflect in CLAUDE.md
2. **Manual maintenance burden** - Users must remember to update prompts
3. **Inconsistent behavior** - Different projects use different prompt versions
4. **Discovery problems** - Users don't know when prompts have been updated

The persona system already solves this problem for personas (automatic injection with `<!-- PERSONAS:START/END -->` markers), but meta-prompts lack equivalent automation.

## Desired Behavior

### Automatic Injection on Init
```bash
# User initializes DDx
$ ddx init

# DDx automatically:
1. Creates CLAUDE.md if missing
2. Reads system.meta_prompt from config (default: "claude/system-prompts/focused.md")
3. Loads prompt from .ddx/library/prompts/claude/system-prompts/focused.md
4. Injects between <!-- DDX-META-PROMPT:START/END --> markers
5. Adds source comment: <!-- Source: claude/system-prompts/focused.md -->
```

### Automatic Sync on Update
```bash
# User updates library
$ ddx update

# DDx automatically:
1. Pulls latest library content
2. Re-reads system.meta_prompt from config
3. Re-injects updated prompt content to CLAUDE.md
4. Even if no git changes pulled, sync happens
```

### Sync Check in Doctor
```bash
# User checks project health
$ ddx doctor

# DDx checks:
✓ DDx configuration valid
✓ Library path exists
✓ Meta-prompt in sync with library
✗ Meta-prompt out of sync
  └─ CLAUDE.md has older version of focused.md
  └─ Fix: Run 'ddx update' to sync
```

### Config Change Triggers Re-sync
```bash
# User changes meta-prompt setting
$ ddx config set system.meta_prompt "claude/system-prompts/strict.md"

# DDx automatically:
1. Removes old meta-prompt from CLAUDE.md
2. Loads new prompt from library
3. Injects new prompt content
```

## Acceptance Criteria

### AC1: Automatic Injection on Init
- [ ] `ddx init` loads meta-prompt from config
- [ ] Default prompt is "claude/system-prompts/focused.md"
- [ ] Prompt content injected between markers
- [ ] Source comment added
- [ ] CLAUDE.md created if doesn't exist

### AC2: Automatic Sync on Update
- [ ] `ddx update` re-syncs meta-prompt
- [ ] Sync happens even if no library changes
- [ ] Updated prompt content replaces old content
- [ ] Source comment updated
- [ ] Config setting respected

### AC3: Doctor Health Check
- [ ] `ddx doctor` reports sync status
- [ ] Detects when CLAUDE.md prompt differs from library
- [ ] Suggests `ddx update` to fix
- [ ] Reports "✓ Meta-prompt in sync" when current

### AC4: Config Change Handling
- [ ] `ddx config set system.meta_prompt <path>` triggers re-sync
- [ ] Setting to null/empty removes meta-prompt section
- [ ] Invalid paths reported as error
- [ ] Validates prompt file exists in library

### AC5: Same Pattern as Personas
- [ ] Uses `<!-- DDX-META-PROMPT:START/END -->` markers
- [ ] Preserves CLAUDE.md content outside markers
- [ ] Similar API to persona injection
- [ ] Code structure mirrors persona system

## Technical Notes

### Files to Create
1. **`cli/internal/metaprompt/injector.go`** - Injection logic
   - `MetaPromptInjector` interface
   - `InjectMetaPrompt(promptPath string) error`
   - `RemoveMetaPrompt() error`
   - `IsInSync() (bool, error)`
   - `GetCurrentMetaPrompt() (string, error)`

2. **`cli/internal/metaprompt/injector_test.go`** - Unit tests

### Files to Modify
1. **`cli/cmd/init.go`** - Add meta-prompt injection after config creation
2. **`cli/cmd/update.go`** - Add meta-prompt sync after library update
3. **`cli/cmd/doctor.go`** - Add sync status check
4. **`cli/cmd/config.go`** - Add re-sync on `system.meta_prompt` change

### Markers Used
```markdown
<!-- DDX-META-PROMPT:START -->
<!-- Source: claude/system-prompts/focused.md -->
# System Instructions

[Prompt content here]
<!-- DDX-META-PROMPT:END -->
```

### Config Schema
Already exists in `cli/internal/config/types.go`:
```go
type SystemConfig struct {
    MetaPrompt *string `yaml:"meta_prompt,omitempty"`
}

func (c *NewConfig) GetMetaPrompt() string {
    if c.System == nil || c.System.MetaPrompt == nil {
        return "claude/system-prompts/focused.md" // Default
    }
    return *c.System.MetaPrompt
}
```

### Implementation Pattern
Follow persona system architecture:
- `cli/internal/persona/claude.go` - Reference implementation
- `cli/internal/persona/types.go` - Interface patterns
- Same marker/section removal/injection logic
- Same file I/O patterns

## Testing Strategy

### Unit Tests
1. Test `InjectMetaPrompt()` with various prompts
2. Test marker detection and removal
3. Test source comment generation
4. Test sync detection logic
5. Test null/empty config handling

### Integration Tests
1. Test `ddx init` creates and injects prompt
2. Test `ddx update` re-syncs prompt
3. Test `ddx doctor` detects sync status
4. Test config change triggers re-sync
5. Test switching between different prompts

### Acceptance Tests
1. Full workflow: init → update → verify sync
2. Doctor reports out-of-sync correctly
3. Config changes reflected in CLAUDE.md
4. Manual edits outside markers preserved

## Dependencies

- Existing persona injection system (reference pattern)
- Config system with `system.meta_prompt` support (already exists)
- Library structure with prompts (already exists)

## Risks

### Risk 1: CLAUDE.md Corruption
**Description**: Faulty marker detection could corrupt CLAUDE.md
**Mitigation**: Copy persona system's robust marker handling, extensive tests
**Severity**: Medium

### Risk 2: Sync Detection False Positives
**Description**: Doctor might incorrectly report out-of-sync
**Mitigation**: Hash-based or content comparison, whitespace normalization
**Severity**: Low

### Risk 3: Performance Impact
**Description**: Reading/writing CLAUDE.md on every update could be slow
**Mitigation**: File operations are fast, minimal impact expected
**Severity**: Low

## Success Metrics

- **Automation**: 100% of prompts synced automatically (zero manual updates)
- **Accuracy**: Doctor sync detection 100% accurate
- **Reliability**: Zero CLAUDE.md corruption incidents
- **Adoption**: Feature works in all existing projects after `ddx update`

## Related Documents

- Feature Spec: `docs/helix/01-frame/features/FEAT-015-meta-prompt-injection.md`
- Reference: `cli/internal/persona/claude.go` (persona injection system)
- Config: `cli/internal/config/types.go` (SystemConfig.MetaPrompt)

## Notes

This feature brings meta-prompt management to parity with the persona system, eliminating manual synchronization and ensuring all projects use current behavioral guidance. It's a quality-of-life improvement that reduces maintenance burden and improves consistency across projects.