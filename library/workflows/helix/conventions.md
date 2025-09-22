# HELIX Workflow Conventions

## Overview

This document defines conventions for projects using the HELIX workflow, ensuring consistency across implementations while allowing for project-specific needs.

## Documentation Structure

### Phase-Based Organization

Projects using HELIX should organize their documentation using the `docs/helix/` convention:

```
project-root/
├── docs/
│   ├── helix/                  # HELIX phase artifacts
│   │   ├── 01-frame/          # Problem definition & requirements
│   │   ├── 02-design/         # Architecture & design decisions
│   │   ├── 03-test/           # Test strategies & plans
│   │   ├── 04-build/          # Implementation guidance
│   │   ├── 05-deploy/         # Deployment & operations
│   │   └── 06-iterate/        # Continuous improvement
│   ├── reference/             # Reference documentation
│   ├── operations/            # Operational procedures
│   └── strategy/              # Strategic planning
```

### Why This Structure?

1. **Clear Separation**: Phase artifacts are distinct from operational/reference docs
2. **Workflow Alignment**: Numbered directories match HELIX phase order
3. **Tool Support**: Consistent structure enables validation and automation
4. **Flexibility**: Non-phase documentation has dedicated locations

### Phase Directory Contents

Each phase directory contains artifacts directly (no `artifacts/` subdirectory):

```
01-frame/
├── README.md                   # Phase overview and status
├── prd/                       # Product requirements
├── user-stories/              # User stories and scenarios
├── stakeholder-map/           # Stakeholder analysis
├── compliance-requirements/   # Regulatory requirements
├── security-requirements/     # Security requirements
└── threat-model/              # Threat modeling

02-design/
├── README.md                   # Phase overview and status
├── adr/                       # Architecture Decision Records
├── architecture/              # System architecture
├── solution-design/           # Solution designs
├── contracts/                 # API contracts
├── data-design/              # Data models
└── security-architecture/    # Security architecture

03-test/
├── README.md                   # Phase overview and status
├── test-plan/                 # Test strategy
├── test-procedures/           # Test procedures
└── security-tests/            # Security test plans

04-build/
├── README.md                   # Phase overview and status
├── implementation-plan/       # Development plans
├── build-procedures/          # Build procedures
└── secure-coding/            # Secure coding guidelines

05-deploy/
├── README.md                   # Phase overview and status
├── deployment-checklist/      # Deployment procedures
├── runbook/                  # Operational runbooks
├── monitoring-setup/         # Monitoring configuration
└── security-monitoring/      # Security monitoring

06-iterate/
├── README.md                   # Phase overview and status
├── metrics-dashboard/         # Performance metrics
├── feedback-analysis/        # User feedback
├── lessons-learned/          # Retrospectives
├── improvement-backlog/      # Enhancement ideas
└── refinements/              # Story refinement logs and tracking
```

## Naming Conventions

### File Names

1. **README.md**: Each phase directory must have a README explaining its purpose and current status
2. **Artifact Names**: Use descriptive, lowercase names with hyphens (e.g., `threat-model.md`, `api-design.md`)
3. **Numbered Items**: When multiple versions exist, use semantic versioning (e.g., `prd-v1.0.md`, `prd-v1.1.md`)

### Directory Names

1. **Phase Directories**: Always use two-digit numbering (01-frame, not 1-frame)
2. **Artifact Directories**: Use lowercase with hyphens, typically plural (e.g., `user-stories`, `contracts`)
3. **No Nesting**: Avoid deep nesting; keep artifacts at most one level deep within phase directories

## Cross-References

### Linking Between Phases

Use relative paths to reference artifacts across phases:

```markdown
# In 02-design/architecture.md
See requirements in [../01-frame/prd/requirements.md](../01-frame/prd/requirements.md)

# In 03-test/test-plan.md
Based on architecture in [../02-design/architecture/system.md](../02-design/architecture/system.md)
```

### Traceability

Maintain clear traceability by:
1. Referencing source requirements in design documents
2. Linking designs to test plans
3. Connecting test results to implementation decisions
4. Tracking deployment issues back to design choices

## Non-Phase Documentation

### Reference Documentation

Place in `docs/reference/`:
- User guides
- API documentation
- Integration guides
- Glossaries

### Operational Documentation

Place in `docs/operations/`:
- Incident response procedures
- Monitoring guides
- Performance tuning
- Backup/recovery procedures

### Strategic Documentation

Place in `docs/strategy/`:
- Product vision
- Roadmaps
- Market analysis
- Competitive analysis

## Migration from Existing Documentation

When migrating existing documentation to HELIX structure:

1. **Analyze Current State**: Map existing docs to HELIX phases
2. **Extract Requirements**: Pull requirements from various sources into 01-frame
3. **Consolidate Design**: Gather architecture docs into 02-design
4. **Identify Gaps**: Note missing artifacts for each phase
5. **Create Placeholders**: Add README files marking TODOs for missing content
6. **Maintain References**: Update all cross-references after migration

## Validation

Projects should validate their documentation structure:

```bash
# Check required phase directories exist
test -d docs/helix/01-frame || echo "Missing frame phase"
test -d docs/helix/02-design || echo "Missing design phase"
# ... etc

# Verify README files in each phase
for phase in docs/helix/*/; do
  test -f "$phase/README.md" || echo "Missing README in $phase"
done

# Check for orphaned references
grep -r "\.\./" docs/helix/ | grep -v "helix"
```

## Templates

Use HELIX workflow templates to create consistent artifacts:

```bash
# Apply templates from helix workflow
ddx apply workflows/helix/phases/01-frame/artifacts/prd

# Copy template structure
cp -r $DDX_HOME/workflows/helix/phases/01-frame/artifacts/prd/template.md \
      docs/helix/01-frame/prd/
```

## Best Practices

1. **Start Early**: Create the structure at project inception
2. **Keep Current**: Update documentation as the project evolves
3. **Review Regularly**: Include doc reviews in phase transitions
4. **Automate Checks**: Add structure validation to CI/CD
5. **Version Control**: Track all documentation changes in git
6. **Link Liberally**: Cross-reference related artifacts
7. **Stay Flat**: Avoid deep directory nesting
8. **Be Consistent**: Follow naming conventions strictly

## FAQ

### Q: Can I add custom directories to phases?
A: Yes, phases can have project-specific subdirectories. Document them in the phase README.

### Q: Should code live in helix/?
A: No, code belongs in the project's source directories. Documentation only in helix.

### Q: How do I handle multiple features in parallel?
A: Create feature-specific subdirectories within each artifact type (e.g., `prd/feature-a/`, `prd/feature-b/`).

### Q: What about diagrams and images?
A: Store them alongside the documents that reference them, or in a phase-level `images/` directory.

### Q: Can I skip phases?
A: While not recommended, if skipping phases, document why in the project root README.

## Story Refinement Conventions

### Refinement Documentation Structure

Story refinements are tracked in the iterate phase to maintain learning and traceability:

```
docs/helix/06-iterate/refinements/
├── README.md                           # Refinement process overview
├── US-001-refinement-001.md           # First refinement of US-001
├── US-001-refinement-002.md           # Second refinement of US-001
├── US-042-refinement-001.md           # First refinement of US-042
└── refinement-index.md                # Cross-reference index
```

### Refinement Naming Convention

**File Naming Pattern**: `{{STORY_ID}}-refinement-{{NUMBER}}.md`
- `{{STORY_ID}}`: Original user story identifier (e.g., US-001, US-042)
- `{{NUMBER}}`: Zero-padded refinement sequence (001, 002, 003...)

Examples:
- `US-001-refinement-001.md` - First refinement of US-001
- `US-042-refinement-003.md` - Third refinement of US-042

### Refinement Linking Strategy

**Story Updates**: Original user stories reference their refinements:
```markdown
## Refinement History
- [Refinement 001](../06-iterate/refinements/US-001-refinement-001.md) - Bug fixes for error handling
- [Refinement 002](../06-iterate/refinements/US-001-refinement-002.md) - Scope expansion for mobile support
```

**Cross-Phase References**: Refinement logs link to all affected documents:
```markdown
### Updated Documents
- [User Story](../01-frame/user-stories/US-001.md) - Updated acceptance criteria
- [Technical Design](../02-design/architecture/auth-service.md) - Added error handling flows
- [Test Plan](../03-test/test-procedures/US-001-tests.md) - Added regression tests
```

### Refinement Categories

**Standard Categories** for consistent tracking:
- `bugs` - Issues discovered during implementation or testing
- `requirements` - New or evolved business requirements
- `enhancement` - Improvements identified during development
- `mixed` - Combination of multiple refinement types

### Version Control Integration

**Branch Strategy** for refinements:
- Create refinement branches: `refinement/US-001-001`
- Commit refinement log first, then affected documents
- Ensure atomic commits for traceability

**Commit Message Format**:
```
refine(US-001): fix error handling specification gaps

- Add refinement log US-001-refinement-001
- Update acceptance criteria for edge cases
- Add regression test requirements
- Update error handling design patterns

Addresses bugs discovered during implementation phase.
```

### Quality Gates for Refinements

**Pre-Refinement Checklist**:
- [ ] Issues clearly documented and categorized
- [ ] Impact assessment completed
- [ ] Stakeholder approval obtained (if scope changes)
- [ ] Current implementation status captured

**Post-Refinement Validation**:
- [ ] All affected phase documents updated
- [ ] Cross-references verified and functional
- [ ] Traceability maintained from issue to resolution
- [ ] No conflicts introduced between requirements
- [ ] Team communication completed

### Refinement Index Maintenance

**Index Structure** for discoverability:
```markdown
# Story Refinement Index

## Active Stories with Refinements
- US-001: [3 refinements](US-001-refinement-001.md) - Authentication Service
- US-042: [1 refinement](US-042-refinement-001.md) - Workflow Commands

## Refinement Categories
### Bugs (High Impact)
- [US-001-refinement-001](US-001-refinement-001.md) - Critical error handling gaps
- [US-018-refinement-002](US-018-refinement-002.md) - Input validation issues

### Requirements Evolution
- [US-025-refinement-001](US-025-refinement-001.md) - Mobile support addition
- [US-042-refinement-001](US-042-refinement-001.md) - Enhanced command discovery
```

### Template Usage

Use the standard refinement template for consistency:
```bash
# Apply refinement template
ddx apply workflows/helix/templates/refinement-log.md docs/helix/06-iterate/refinements/US-001-refinement-001.md

# Fill template variables
STORY_ID=US-001 REFINEMENT_NUMBER=001 REFINEMENT_TYPE=bugs ddx apply ...
```

## Evolution

These conventions will evolve based on usage. To propose changes:

1. Document the issue with current conventions
2. Propose specific changes with rationale
3. Show examples of the new approach
4. Update this document after consensus

---

*These conventions ensure consistency while maintaining flexibility for project-specific needs. They enable tooling support and make HELIX projects more maintainable and understandable.*