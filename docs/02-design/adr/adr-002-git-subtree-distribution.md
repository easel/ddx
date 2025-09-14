---
tags: [adr, architecture, development-workflow, decisions, ddx, git-subtree]
template: false
version: 1.0.0
---

# ADR-001: Use Git Subtree for DDX Resource Distribution

**Date**: 2025-01-12  
**Status**: Accepted  
**Deciders**: DDX Development Team  
**Technical Story**: Design the core mechanism for distributing and synchronizing DDX resources across projects

## Context

### Problem Statement
DDX needs to distribute templates, prompts, patterns, and configurations across multiple projects while allowing for:
- Local customization without breaking updates
- Bidirectional contribution flow
- Version control integration
- Minimal tooling dependencies
- Reliable and predictable behavior

### Forces at Play
- **Simplicity**: Developers want minimal complexity and setup
- **Flexibility**: Projects need to customize resources locally
- **Contribution**: Users should easily share improvements
- **Versioning**: Need full history and rollback capabilities
- **Independence**: Projects shouldn't depend on external services
- **Compatibility**: Must work with existing Git workflows

### Constraints
- Must work with standard Git installations
- Cannot require additional runtime dependencies
- Must preserve full Git history
- Should integrate with existing CI/CD pipelines
- Need to support offline development

## Decision

### Chosen Approach
Use Git subtree to manage DDX resources as a subdirectory within each project, pulling from and pushing to a central DDX repository.

### Rationale
- **Native Git Feature**: Git subtree is built into Git, requiring no additional tools
- **Full History**: Preserves complete commit history from the DDX repository
- **Local Independence**: Projects can work offline and make local modifications
- **Clean Integration**: Resources appear as normal files in the project repository
- **Bidirectional Flow**: Supports both pulling updates and contributing changes
- **No Submodule Complexity**: Avoids the complexities and gotchas of Git submodules

## Alternatives Considered

### Option 1: Git Submodules
**Description**: Use Git submodules to link DDX resources as a separate repository reference

**Pros**:
- Clear separation between project and DDX code
- Explicit version pinning
- Smaller repository size

**Cons**:
- Complex workflow for developers
- Requires manual submodule initialization and updates
- Detached HEAD states confuse users
- Poor support in many Git GUIs
- Difficult to make local modifications

**Why rejected**: Too complex for the target user base and creates friction in daily development workflow

### Option 2: Package Manager Distribution (npm/pip/cargo)
**Description**: Distribute DDX resources as packages through language-specific package managers

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

**Why rejected**: Creates language silos and adds infrastructure complexity that conflicts with the simplicity goal

### Option 3: Direct File Copying/Syncing
**Description**: Use a custom sync mechanism or rsync to copy files from a master repository

**Pros**:
- Very simple conceptually
- Fast initial setup
- No Git complexity

**Cons**:
- Loses version history
- No built-in conflict resolution
- Difficult to track what's been modified
- No standard contribution path
- Custom tooling required

**Why rejected**: Loses critical version control benefits and requires custom tooling

### Option 4: Monorepo Approach
**Description**: Maintain all projects in a single repository with shared resources

**Pros**:
- Single source of truth
- Atomic updates across projects
- Simplified dependency management
- Consistent tooling

**Cons**:
- Not feasible for independent projects
- Massive repository size
- Requires organizational buy-in
- Complex permissions management
- CI/CD becomes complicated

**Why rejected**: Incompatible with the distributed nature of independent development projects

## Consequences

### Positive Consequences
- **Zero Additional Dependencies**: Works with standard Git installation
- **Full Version Control**: Complete history preservation enables debugging and rollback
- **Local Flexibility**: Projects can customize without breaking update flow
- **Familiar Workflow**: Uses standard Git commands developers already know
- **Contribution Path**: Clear mechanism for sharing improvements via subtree push
- **Offline Support**: Full functionality without network access

### Negative Consequences
- **Repository Size**: DDX resources add to project repository size
- **Merge Complexity**: Subtree merges can be complex when conflicts arise
- **Learning Curve**: Developers need to understand subtree commands
- **History Pollution**: DDX commits appear in project history
- **Squash Limitations**: Squashing can lose granular history

### Neutral Consequences
- **Directory Structure**: DDX resources live in `.ddx/` subdirectory
- **Command Abstraction**: CLI tool wraps complex subtree commands
- **Update Frequency**: Projects control when to pull updates
- **Contribution Process**: Requires explicit push to contribute

## Implementation

### Required Changes
- CLI tool implements subtree wrapper commands
- Documentation for subtree workflow
- `.ddx.yml` configuration for resource selection
- Contribution guidelines for subtree push
- Conflict resolution documentation

### Migration Strategy
For new projects:
1. `ddx init` runs `git subtree add`
2. Resources immediately available in `.ddx/`
3. Configuration via `.ddx.yml`

For existing projects:
1. Remove any existing DDX installation
2. Run `ddx init` to add via subtree
3. Migrate local customizations

### Success Metrics
- **Setup Time**: < 1 minute for initial installation
- **Update Success Rate**: > 95% of updates without conflicts
- **Contribution Rate**: > 10 contributions per month
- **User Satisfaction**: > 80% prefer subtree over alternatives

## Compliance

### Security Requirements
- No credentials or secrets in DDX repository
- Git's standard security model applies
- SSH/HTTPS for repository access

### Performance Requirements
- Initial clone: < 30 seconds on standard connection
- Update pull: < 10 seconds for typical update
- No runtime performance impact

### Regulatory Requirements
- Open source licensing compatibility
- No export control restrictions
- GDPR compliance for any user data

## Monitoring and Review

### Key Indicators to Watch
- Subtree merge conflict frequency
- Repository size growth rate
- Update adoption rate
- Contribution success rate
- User feedback on complexity

### Review Date
Q2 2025 - After 6 months of production use

### Review Triggers
- Major Git version with subtree changes
- > 20% of updates causing conflicts
- Repository size exceeds 100MB
- User satisfaction drops below 70%

## Related Decisions

### Dependencies
- ADR-002: DDX CLI Architecture - CLI wraps subtree complexity
- ADR-003: Configuration File Format - YAML configuration for resource selection

### Influenced By
- Git's subtree implementation and limitations
- Target user base (developers familiar with Git)
- Simplicity as core design principle

### Influences
- Future decisions about resource organization
- Contribution workflow design
- Update notification mechanisms
- Conflict resolution strategies

## References

### Documentation
- [Git Subtree Documentation](https://git-scm.com/book/en/v2/Git-Tools-Advanced-Merging#_subtree_merge)
- [DDX Architecture Overview](/docs/architecture/README.md)
- [DDX PRD](/docs/product/prd-ddx-v1.md)

### External Resources
- [Atlassian Git Subtree Tutorial](https://www.atlassian.com/git/tutorials/git-subtree)
- [Git Subtree vs Submodule Comparison](https://stackoverflow.com/questions/31769820/differences-between-git-submodule-and-subtree)

### Discussion History
- Initial architecture discussion in project kickoff
- Community feedback on resource distribution approaches
- Performance testing results for various approaches

## Notes

The git subtree approach aligns well with DDX's medical metaphor - like distributing medical knowledge and best practices across different hospitals while allowing each to maintain their own procedures and customizations. The subtree acts as a "knowledge transfer" mechanism that preserves the full context and history of decisions.

Key insight: By abstracting subtree complexity behind the DDX CLI, we get the power of subtree without exposing users to its complexity. This is similar to how modern package managers abstract complex dependency resolution.

---

**Last Updated**: 2025-01-12  
**Next Review**: 2025-07-01