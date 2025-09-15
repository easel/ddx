# ADR-002: Use Git Subtree for DDX Resource Distribution

**Date**: 2025-01-12
**Status**: Accepted
**Deciders**: DDX Development Team
**Related Feature(s)**: Cross-cutting - Distribution and Synchronization
**Confidence Level**: High

## Context

DDX needs to distribute templates, prompts, patterns, and configurations across multiple projects while maintaining version control, enabling local customization, and supporting bidirectional contribution flow.

### Problem Statement

How do we distribute and synchronize DDX resources across projects in a way that allows local customization, preserves version history, supports offline development, and enables community contributions without requiring additional tooling or infrastructure?

### Current State

Most tools use package managers, git submodules, or direct downloads for resource distribution. Each approach has limitations: package managers create language silos, submodules add complexity, and downloads lose version history.

### Requirements Driving This Decision
- Must work with standard Git installations without additional tools
- Support local customization without breaking updates
- Enable bidirectional contribution flow
- Preserve full Git history for resources
- Work reliably offline
- Integrate with existing Git workflows
- Minimal setup complexity for users

## Decision

We will use Git subtree to manage DDX resources as a subdirectory within each project, pulling from and pushing to a central DDX repository.

### Key Points
- Resources stored in `.ddx/` subdirectory via git subtree
- Full history preservation from DDX repository
- Local modifications possible without fork
- Updates can be pulled selectively
- Contributions can be pushed back upstream
- No additional runtime dependencies required
- Works with all standard Git tools and GUIs

## Alternatives Considered

### Option 1: Git Submodules
**Description**: Use Git submodules to link DDX resources as a separate repository reference

**Pros**:
- Clear separation between project and DDX code
- Explicit version pinning
- Smaller repository size
- Standard Git feature

**Cons**:
- Complex workflow for developers
- Requires manual submodule initialization
- Detached HEAD states confuse users
- Poor support in many Git GUIs
- Difficult to make local modifications

**Evaluation**: Rejected due to complexity and poor developer experience

### Option 2: Package Manager Distribution
**Description**: Distribute DDX resources as packages through npm/pip/cargo

**Pros**:
- Familiar distribution mechanism
- Semantic versioning support
- Dependency resolution
- Wide tooling support

**Cons**:
- Language-specific, limiting cross-language usage
- Requires package registry infrastructure
- Harder to customize resources locally
- Publishing overhead for contributions
- Version lock-in issues

**Evaluation**: Rejected because it creates language silos and infrastructure complexity

### Option 3: Git Subtree (Selected)
**Description**: Use Git subtree to embed DDX resources with full history

**Pros**:
- Native Git feature, no additional tools
- Full history preservation
- Local independence
- Clean integration in project
- Bidirectional flow support
- No submodule complexity

**Cons**:
- Repository size increases - mitigated by selective pulls
- Merge complexity for conflicts - addressed with CLI tooling
- Less familiar to some developers - solved with documentation

**Evaluation**: Selected for best balance of power and simplicity

## Consequences

### Positive Consequences
- **Zero Dependencies**: Works with standard Git installation
- **Full History**: Complete commit history preserved
- **Local Flexibility**: Projects can customize without forking
- **Offline Work**: Full functionality without network
- **Clean Integration**: Resources appear as normal project files
- **Contribution Path**: Clear mechanism for sharing improvements

### Negative Consequences
- **Repository Size**: DDX resources increase project size by ~5-10MB
- **Merge Complexity**: Subtree merges can be complex during conflicts
- **Learning Curve**: Developers need to understand subtree commands
- **History Mixing**: DDX commits appear in project history

### Neutral Consequences
- **Directory Location**: DDX resources in `.ddx/` subdirectory
- **Command Wrapping**: CLI tool abstracts subtree complexity
- **Update Control**: Projects decide when to pull updates

## Implementation Impact

### Development Impact
- **Effort**: Low - Git subtree is existing functionality
- **Time**: 1 week for CLI wrapper implementation
- **Skills Required**: Git subtree knowledge, Go for CLI

### Operational Impact
- **Performance**: Minimal - standard Git operations
- **Scalability**: Excellent - distributed by nature
- **Maintenance**: Low - leverages Git infrastructure

### Security Impact
- Standard Git security model applies
- No additional attack surface
- SSH/HTTPS for repository access
- GPG signing supported

## Risks and Mitigation

| Risk | Probability | Impact | Mitigation Strategy |
|------|------------|--------|-------------------|
| Subtree merge conflicts | Medium | Medium | CLI tooling to assist resolution |
| Repository bloat | Low | Low | Selective resource pulling |
| User confusion | Medium | Low | Clear documentation and CLI abstraction |
| History pollution | Low | Low | Squash options available |

## Dependencies

### Technical Dependencies
- Git 1.7.11+ (subtree support)
- Standard Git installation
- Network access for initial setup (not runtime)

### Decision Dependencies
- ADR-001: Defines resource structure to distribute
- ADR-003: CLI implementation wraps subtree commands

## Validation

### How We'll Know This Was Right
- Setup time < 1 minute for new projects
- Update success rate > 95%
- Contribution flow works bidirectionally
- No additional tooling required
- Offline development fully supported

### Review Triggers
This decision should be reviewed if:
- Git introduces breaking changes to subtree
- Better distribution mechanism emerges
- Repository size becomes problematic (>100MB)
- Merge conflicts exceed 20% of updates

## References

### Internal References
- [DDX Architecture Overview](/docs/architecture/architecture-overview.md)
- [DDX CLI Design](/docs/architecture/cli-architecture.md)
- Related ADRs: ADR-001, ADR-003

### External References
- [Git Subtree Documentation](https://git-scm.com/book/en/v2/Git-Tools-Advanced-Merging#_subtree_merge)
- [Git Subtree vs Submodule](https://stackoverflow.com/questions/31769820/differences-between-git-submodule-and-subtree)
- [Atlassian Git Subtree Tutorial](https://www.atlassian.com/git/tutorials/git-subtree)

## Notes

### Meeting Notes
- Team consensus on avoiding external dependencies
- Discussion on submodule pain points from past projects
- Agreement on need for bidirectional contribution flow

### Future Considerations
- Consider git sparse-checkout for large resource sets
- Explore git worktree for development workflows
- Investigate incremental update strategies
- Monitor git subtree evolution and alternatives

### Lessons Learned
*To be filled after 6 months of production use*

---

## Decision History

### 2025-01-12 - Initial Decision
- Status: Proposed
- Author: DDX Development Team
- Notes: Evaluation of distribution options

### 2025-01-12 - Review and Acceptance
- Status: Accepted
- Reviewers: DDX Core Team
- Changes: Added mitigation strategies for identified risks

### Post-Implementation Review
- *To be scheduled after Q3 2025*

---
*This ADR documents a significant architectural decision and its rationale for future reference.*