# Solution Design: Contribution Command (US-005)

*Technical implementation design for `ddx contribute` command*

**User Story**: US-005 - Contribute Improvements
**Feature**: FEAT-001 - Core CLI Framework
**Dependencies**: FEAT-002 (git.SubtreePush primitive)
**Created**: 2025-01-15

## Overview

This document specifies the technical implementation of the `ddx contribute` command, which enables developers to share improvements back to the community. The design focuses on the command implementation, integrating with existing git subtree primitives from FEAT-002.

## Requirements from US-005

### Acceptance Criteria to Implement
1. Contribution workflow initiated with clear steps
2. Changes formatted appropriately for submission
3. System checks contribution meets standards
4. Contribution sent to upstream repository
5. Clear submission status and next steps
6. Contribution guidelines readily available
7. Credentials handled securely
8. Selective asset contribution support

## Command Architecture

### Command Flow

```
ddx contribute [path] [flags]
│
├─→ Parse flags and validate inputs
├─→ Load configuration (.ddx.yml)
├─→ Validate authentication (git credential helpers)
├─→ Detect changes to contribute
├─→ Run validation checks
│   ├─ Secret detection
│   ├─ Format validation
│   └─ Standards compliance
│
├─→ [DRY RUN MODE]
│   ├─ Display what would be contributed
│   ├─ Show validation results
│   └─ Exit (no git operations)
│
├─→ [EXECUTE MODE]
│   ├─ Collect metadata (message, description)
│   ├─ Execute git subtree push
│   │   └─ Call git.SubtreePush(".ddx/library", repoURL, branch)
│   ├─ [if --create-pr] Generate PR instructions
│   └─ Display success message with next steps
│
└─→ Return result
```

### Component Integration

```
┌─────────────────┐
│ contribute.go   │ ← Cobra command handler
│ (existing)      │
└────────┬────────┘
         │
         ├─→ validateContribution() ← Validation logic
         ├─→ collectMetadata()      ← User prompts
         └─→ executeContributionInDir() ← NEEDS IMPLEMENTATION
                    │
                    ├─→ git.SubtreePush()  ← FEAT-002 primitive
                    └─→ createPRInstructions() ← GitHub guidance
```

## Validation Specification

### Standards Checks

```go
type ValidationResult struct {
    Passed       bool
    Errors       []ValidationError
    Warnings     []ValidationWarning
}

func validateContribution(path string) (*ValidationResult, error) {
    // 1. Secret detection
    if hasSecrets := detectSecrets(path); hasSecrets {
        return fail("Secrets detected in contribution")
    }

    // 2. Format validation
    if !validateFormat(path) {
        return fail("Invalid format or structure")
    }

    // 3. Documentation check
    if !hasDocumentation(path) {
        return warn("Missing or incomplete documentation")
    }

    // 4. Size limits
    if exceedsSizeLimit(path) {
        return fail("Contribution exceeds size limit")
    }

    return passed()
}
```

### Validation Rules

| Check | Type | Action on Failure |
|-------|------|------------------|
| Secrets detected | Error | Block contribution |
| Invalid file format | Error | Block contribution |
| Missing documentation | Warning | Allow with warning |
| Exceeds size limit | Error | Block contribution |
| No tests (for code) | Warning | Allow with warning |

## Git Operations

### Subtree Push Implementation

```go
func executeContributionInDir(workingDir string, cfg *config.Config, opts *ContributeOptions) (*ContributeResult, error) {
    // 1. Prepare contribution
    if err := validateAuthentication(); err != nil {
        return nil, fmt.Errorf("authentication required: %w", err)
    }

    // 2. Execute git subtree push
    prefix := ".ddx/library"
    remote := cfg.Library.Repository.URL
    branch := cfg.Library.Repository.ContributionBranch // from config

    err := git.SubtreePush(
        workingDir,
        prefix,
        remote,
        branch,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to push contribution: %w", err)
    }

    // 3. Build result
    result := &ContributeResult{
        Success: true,
        Message: "Contribution submitted successfully!",
        Branch:  branch,
    }

    // 4. Add PR instructions if requested
    if opts.CreatePR {
        result.PRInfo = generatePRInstructions(cfg, opts)
    }

    return result, nil
}
```

### Push Strategy

**Approach**: Direct push to contribution branch (no feature branch creation)

**Rationale**:
- Simplest implementation for MVP
- Matches existing config: `subtree.contribution_branch`
- Users can create branches manually if needed
- Future enhancement: auto-create feature branches

**Configuration** (`.ddx.yml`):
```yaml
library:
  repository:
    url: https://github.com/easel/ddx
    branch: main
    contribution_branch: contributions  # Where to push
```

## GitHub Integration

### PR Creation

**Phase 1** (Current): Generate instructions
```go
func generatePRInstructions(cfg *config.Config, opts *ContributeOptions) *PRInfo {
    repoURL := strings.TrimSuffix(cfg.Library.Repository.URL, ".git")
    baseBranch := cfg.Library.Repository.Branch
    headBranch := cfg.Library.Repository.ContributionBranch

    compareURL := fmt.Sprintf("%s/compare/%s...%s",
        repoURL, baseBranch, headBranch)

    return &PRInfo{
        URL:    compareURL,
        Title:  opts.Message,
        Branch: headBranch,
        Description: "Visit the URL above to create a pull request",
    }
}
```

**Phase 2** (Future): GitHub API integration
- Requires GitHub token
- Creates PR automatically
- Sets title, description, labels
- Out of scope for initial implementation

## Error Handling

### Error Scenarios and Messages

| Scenario | Error Message | User Action |
|----------|--------------|-------------|
| Not authenticated | "Git authentication required. Run: git credential approve" | Configure git credentials |
| Validation failed | "Contribution validation failed: [specific errors]" | Fix issues and retry |
| Network error | "Network error pushing contribution. Check connection." | Retry when online |
| Push conflict | "Push rejected: branch has diverged. Update and retry." | Pull latest changes |
| No changes detected | "No changes detected in [path]. Nothing to contribute." | Make changes first |

### Implementation

```go
// Wrap git errors with user-friendly messages
func wrapGitError(err error) error {
    if strings.Contains(err.Error(), "authentication") {
        return fmt.Errorf("authentication required: %w\n\nConfigure git credentials with your GitHub token", err)
    }
    if strings.Contains(err.Error(), "rejected") {
        return fmt.Errorf("push rejected: %w\n\nYour contribution conflicts with recent changes. Pull latest and retry", err)
    }
    return err
}
```

## Flag Behavior

### Supported Flags

```bash
ddx contribute [path] [flags]
```

| Flag | Type | Default | Behavior |
|------|------|---------|----------|
| `--message, -m` | string | (prompt) | Contribution message |
| `--dry-run` | bool | false | Preview without executing |
| `--create-pr` | bool | false | Generate PR instructions |
| `--branch` | string | (from config) | Override contribution branch |

### Dry Run Implementation

```go
if opts.DryRun {
    fmt.Println("DRY RUN: Would contribute the following:")
    fmt.Printf("  Path: %s\n", opts.ResourcePath)
    fmt.Printf("  Message: %s\n", opts.Message)
    fmt.Printf("  Validation: %s\n", validationResult)
    fmt.Printf("  Push to: %s/%s\n", cfg.Repository.URL, branch)
    return &ContributeResult{Success: true, DryRun: true}, nil
}
```

## Implementation Checklist

### Wire Existing Code
- [x] Command structure exists (`cmd/contribute.go`)
- [x] Validation framework exists
- [x] Result structures defined
- [x] User prompts implemented
- [ ] **Wire `executeContributionInDir()` to `git.SubtreePush()`** ← KEY WORK
- [ ] Add error wrapping for user-friendly messages
- [ ] Implement PR instruction generation

### Test Coverage
- [ ] Unit tests for validation logic
- [ ] Integration tests with mock git remote
- [ ] Dry-run behavior tests
- [ ] Error scenario tests

## Success Metrics

From US-005:
- Time from change to submission < 5 minutes
- Validation catches 95% of issues before submission
- Successful submission rate > 90%
- Clear error messages for failures

## Future Enhancements

**Not in scope for initial implementation:**
1. Automatic feature branch creation
2. GitHub API PR creation (requires token management)
3. Multi-platform contribution (GitLab, Bitbucket)
4. Automated quality checks beyond secrets/format
5. Contribution preview with diff

## References

- **US-005**: User story with acceptance criteria
- **FEAT-002**: Upstream synchronization system (provides git primitives)
- **ADR-008**: Community contribution governance (review process)
- **SD-002**: Overall synchronization system architecture
- **Implementation**: `cli/cmd/contribute.go`, `cli/internal/git/git.go`

---
*This design focuses on the minimal viable implementation of `ddx contribute` to satisfy US-005 acceptance criteria while keeping complexity low.*