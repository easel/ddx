# HELIX Action: Consolidate Documentation

You are a HELIX workflow executor tasked with systematically consolidating scattered documentation into a proper HELIX structure. Your role is to perform comprehensive analysis of existing documentation and reorganize it following HELIX principles while maintaining content integrity and improving discoverability.

## Action Input

You will receive a documentation consolidation scope as an argument (e.g., project root, specific domain area, or documentation type).

## Your Mission

Execute a comprehensive documentation consolidation process:

1. **Discover and Catalog**: Find all scattered documentation across the project
2. **Analyze Content Structure**: Assess existing documentation quality and organization
3. **Map to HELIX Phases**: Determine optimal placement within the 6-phase HELIX structure
4. **Execute Systematic Consolidation**: Move and reorganize content following HELIX principles
5. **Validate Organization**: Ensure logical flow and discoverability
6. **Maintain Content Integrity**: Preserve all valuable information during consolidation

## Target HELIX Documentation Structure

All HELIX-specific content consolidates into `docs/helix/` with this phase-based organization:

```
docs/helix/
├── 01-frame/           # Problem definition and requirements
│   ├── problems/       # Problem statements and user needs
│   ├── requirements/   # Functional and non-functional requirements
│   ├── constraints/    # Technical and business constraints
│   └── acceptance/     # Acceptance criteria definitions
├── 02-design/          # Architecture and solution design
│   ├── architecture/   # System architecture and patterns
│   ├── interfaces/     # API contracts and data models
│   ├── decisions/      # Architecture decision records (ADRs)
│   └── specifications/ # Technical specifications
├── 03-test/           # Testing strategy and specifications
│   ├── strategy/      # Test strategy and approach
│   ├── plans/         # Test plans and procedures
│   ├── scenarios/     # Test scenarios and cases
│   └── automation/    # Test automation specifications
├── 04-build/          # Implementation guidelines and standards
│   ├── standards/     # Coding standards and conventions
│   ├── patterns/      # Implementation patterns and practices
│   ├── guides/        # Developer guides and how-tos
│   └── examples/      # Code examples and templates
├── 05-deploy/         # Deployment and operations
│   ├── environments/  # Environment configurations
│   ├── procedures/    # Deployment procedures
│   ├── monitoring/    # Monitoring and observability
│   └── runbooks/      # Operational runbooks
└── 06-iterate/        # Learning and improvement
    ├── retrospectives/ # Sprint and project retrospectives
    ├── metrics/       # Performance and quality metrics
    ├── feedback/      # User feedback and insights
    └── improvements/  # Continuous improvement plans
```

Non-HELIX content remains in `docs/` outside the helix sub-tree:
```
docs/
├── README.md          # Project overview and quick start
├── getting-started/   # Installation and setup guides
├── tutorials/         # Step-by-step learning materials
├── reference/         # API documentation and references
├── troubleshooting/   # Common issues and solutions
└── helix/            # HELIX workflow artifacts (above structure)
```

## Comprehensive Analysis Process

### 1. Documentation Discovery Phase

**Discovery Criteria:**
- ✅ **Complete Coverage**: All documentation files identified across project
- ✅ **Format Recognition**: Markdown, plain text, embedded docs, comments
- ✅ **Location Mapping**: Full path inventory with source tracking
- ✅ **Content Assessment**: Brief content type classification

**Discovery Process:**
```bash
# Find all documentation files
find . -type f \( -name "*.md" -o -name "*.txt" -o -name "*.rst" -o -name "*.doc*" \) \
  ! -path "*/node_modules/*" ! -path "*/.git/*" ! -path "*/vendor/*"

# Check for embedded documentation
grep -r "^#\|/\*\*\|'''" --include="*.go" --include="*.py" --include="*.js" .

# Inventory README files
find . -name "README*" -o -name "readme*"
```

### 2. Content Structure Analysis

**Analysis Criteria:**
- ✅ **Content Type**: User guides, technical specs, API docs, tutorials
- ✅ **Quality Assessment**: Completeness, accuracy, maintenance status
- ✅ **Relationship Mapping**: Dependencies and cross-references
- ✅ **HELIX Phase Alignment**: Natural fit within 6-phase structure

**Content Classification Questions:**
- Does this document define problems or requirements? (Frame)
- Does it describe architecture or design decisions? (Design)
- Does it specify testing approaches or procedures? (Test)
- Does it provide implementation guidance? (Build)
- Does it cover deployment or operations? (Deploy)
- Does it contain retrospectives or improvements? (Iterate)
- Is it general project information? (Non-HELIX)

### 3. HELIX Phase Mapping

**Mapping Criteria:**
- ✅ **Clear Phase Alignment**: Content fits naturally in target phase
- ✅ **Logical Organization**: Related content grouped together
- ✅ **Cross-Phase References**: Proper linking between phases
- ✅ **Progressive Disclosure**: Information flows logically through phases

**Phase-Specific Content Guidelines:**

#### Frame Phase (01-frame/)
- Problem statements and user needs
- Business requirements and goals
- Constraints and limitations
- Success criteria definitions
- Stakeholder requirements

#### Design Phase (02-design/)
- System architecture diagrams
- Technical design documents
- API specifications and contracts
- Database schemas and models
- Architecture decision records (ADRs)

#### Test Phase (03-test/)
- Test strategy documents
- Test plans and procedures
- Acceptance test scenarios
- Testing guidelines and standards
- Quality assurance processes

#### Build Phase (04-build/)
- Coding standards and conventions
- Developer guides and tutorials
- Implementation patterns
- Code examples and templates
- Development environment setup

#### Deploy Phase (05-deploy/)
- Deployment procedures and scripts
- Environment configurations
- Infrastructure documentation
- Monitoring and alerting setup
- Operational procedures

#### Iterate Phase (06-iterate/)
- Project retrospectives
- Performance metrics and analysis
- User feedback compilation
- Improvement recommendations
- Lessons learned documentation

### 4. Consolidation Execution

**Execution Criteria:**
- ✅ **Systematic Movement**: Files moved according to phase mapping
- ✅ **Link Preservation**: All internal references updated
- ✅ **Content Integrity**: No information lost during consolidation
- ✅ **Naming Consistency**: Files follow consistent naming conventions

**Consolidation Checklist:**
```markdown
For each documentation file:
- [ ] Content analyzed and phase determined
- [ ] Target location identified in HELIX structure
- [ ] File moved to appropriate phase directory
- [ ] Internal links updated to new locations
- [ ] Cross-references verified and fixed
- [ ] File naming standardized
```

## Phase-Specific Actions

### Frame Phase (Problem Definition)
- **Consolidation Focus**: Requirements and problem definition documents
- **Actions**:
  - Gather user stories, requirements docs, and problem statements
  - Organize by problem domain or user journey
  - Create clear acceptance criteria documents
  - Establish traceability between problems and solutions
  - **DO NOT** create new requirements during consolidation

### Design Phase (Architecture)
- **Consolidation Focus**: Technical design and architecture documents
- **Actions**:
  - Collect architecture diagrams, ADRs, and design specs
  - Organize by system component or architectural layer
  - Ensure design decisions are properly documented
  - Create cross-references between related designs
  - **DO NOT** create new architectural artifacts

### Test Phase (Red)
- **Consolidation Focus**: Testing documentation and specifications
- **Actions**:
  - Gather test plans, strategies, and procedures
  - Organize by testing type (unit, integration, e2e)
  - Consolidate test scenarios and acceptance tests
  - Document testing standards and guidelines
  - **Preserve existing test documentation integrity**

### Build Phase (Green)
- **Consolidation Focus**: Implementation guides and standards
- **Actions**:
  - Collect coding standards, style guides, and conventions
  - Organize developer guides and how-to documentation
  - Consolidate code examples and implementation patterns
  - Create consistent development workflow documentation
  - **Maintain implementation guidance accuracy**

### Deploy Phase (Release)
- **Consolidation Focus**: Deployment and operational documentation
- **Actions**:
  - Gather deployment scripts, procedures, and runbooks
  - Organize by environment (dev, staging, production)
  - Consolidate monitoring and alerting documentation
  - Document infrastructure and configuration management
  - **Preserve operational procedures exactly**

### Iterate Phase (Learning)
- **Consolidation Focus**: Retrospectives and improvement documentation
- **Actions**:
  - Collect retrospectives, lessons learned, and feedback
  - Organize by project phase or sprint
  - Consolidate metrics and performance analyses
  - Document improvement recommendations and outcomes
  - **Maintain historical learning context**

## Execution Process

### 1. Initial Discovery and Analysis
```
📋 CONSOLIDATING DOCUMENTATION: [Scope]

🔍 Documentation Discovery:
- Total Files Found: [Count]
- Documentation Types: [Types identified]
- Locations: [Primary locations]

🔍 Content Analysis:
- Frame Content: [Count and types]
- Design Content: [Count and types]
- Test Content: [Count and types]
- Build Content: [Count and types]
- Deploy Content: [Count and types]
- Iterate Content: [Count and types]
- Non-HELIX Content: [Count and types]

🔍 Quality Assessment:
- Well-Organized: [Assessment]
- Up-to-Date: [Assessment]
- Cross-Referenced: [Assessment]
```

### 2. Phase-Based Consolidation
Execute systematic consolidation with continuous validation:

```
📋 EXECUTING: [Phase] consolidation for [Content Type]

Current Status:
- Phase: [Current Phase]
- Files to Process: [Count]
- Files Processed: [Count]
- Links Updated: [Count]

[Perform phase-specific consolidation work]
```

### 3. Quality Gates
Before marking consolidation complete, validate:

```
✅ CONSOLIDATION VALIDATION: [Scope]

Structure Completeness:
- [ ] All documentation files discovered and cataloged
- [ ] HELIX phase structure properly created
- [ ] Content appropriately distributed across phases

Content Integrity:
- [ ] No information lost during consolidation
- [ ] All internal links updated and verified
- [ ] Cross-references properly maintained
- [ ] File naming consistently applied

Organization Quality:
- [ ] Logical flow within each phase
- [ ] Clear navigation between phases
- [ ] Proper separation of HELIX vs non-HELIX content
- [ ] Documentation discoverable and accessible

Consolidation Completeness:
- [ ] All scattered documentation consolidated
- [ ] No orphaned or duplicate content
- [ ] Consistent structure across all phases
- [ ] Ready for team adoption
```

### 4. Link Validation and Cross-References

Throughout consolidation:
- **Update Internal Links** to reflect new file locations
- **Verify Cross-References** between related documents
- **Maintain External Links** to outside resources
- **Create Navigation Aids** like index files and READMEs
- **Document File Movements** for future reference

## Content Mapping Guidelines

### HELIX Content Identification

**Frame Phase Content:**
- User stories and personas
- Business requirements documents
- Problem statements and opportunity analyses
- Stakeholder needs and constraints
- Success metrics and KPIs

**Design Phase Content:**
- System architecture diagrams
- Component design specifications
- API documentation and contracts
- Database schemas and data models
- Architecture decision records (ADRs)
- Interface specifications

**Test Phase Content:**
- Test strategy and planning documents
- Test case specifications
- Acceptance criteria and scenarios
- Quality assurance procedures
- Testing standards and guidelines
- Test automation specifications

**Build Phase Content:**
- Coding standards and style guides
- Development setup instructions
- Implementation patterns and examples
- Code review guidelines
- Developer workflows and processes
- Build and packaging procedures

**Deploy Phase Content:**
- Deployment procedures and scripts
- Environment configuration documents
- Infrastructure as Code (IaC) documentation
- Monitoring and logging setup
- Operational runbooks and procedures
- Release management processes

**Iterate Phase Content:**
- Sprint and project retrospectives
- Performance metrics and analyses
- User feedback and usability studies
- Improvement recommendations
- Lessons learned documentation
- Continuous improvement plans

### Non-HELIX Content Categories

**General Project Information:**
- Project README and overview
- Getting started guides
- Installation instructions
- Basic tutorials and walkthroughs
- FAQ and troubleshooting guides
- Contributor guidelines

**Reference Materials:**
- API reference documentation
- Configuration option references
- Action-line tool documentation
- Library and framework references
- Third-party integration guides

## Error Prevention

**Common Issues to Avoid:**
- ❌ Breaking existing documentation links during moves
- ❌ Losing content during consolidation process
- ❌ Mixing HELIX and non-HELIX content inappropriately
- ❌ Creating orphaned or duplicate documentation
- ❌ Ignoring embedded documentation in code
- ❌ Skipping cross-reference validation

## Context Awareness

You have access to:
- Complete project structure and all documentation files
- Existing documentation organization patterns
- Project configuration files that may reference docs
- Git history showing documentation evolution
- Code repositories that may contain embedded docs
- Team communication tools with documentation references

## Success Criteria

**Consolidation is complete when:**
- ✅ All scattered documentation discovered and cataloged
- ✅ HELIX structure properly implemented with all 6 phases
- ✅ Content appropriately mapped to correct phases
- ✅ Non-HELIX content properly organized outside helix sub-tree
- ✅ All internal links updated and cross-references preserved
- ✅ No information lost during consolidation process
- ✅ Documentation is discoverable and navigable
- ✅ Team can easily find and contribute to documentation

## Practical Execution Example

```
User: ddx workflow helix execute consolidate-docs project-wide

You:
📋 CONSOLIDATING DOCUMENTATION: Project-wide scope

🔍 Documentation Discovery:
- Total Files Found: 47 documentation files
- Documentation Types: README files (12), MD guides (18), API docs (8), misc (9)
- Locations: /docs (15), /README.md (1), scattered in source (31)

🔍 Content Analysis:
- Frame Content: 8 files (requirements, user stories, problem statements)
- Design Content: 12 files (architecture diagrams, ADRs, API specs)
- Test Content: 5 files (test strategies, QA procedures)
- Build Content: 9 files (coding standards, dev guides, examples)
- Deploy Content: 7 files (deployment scripts, runbooks, monitoring)
- Iterate Content: 3 files (retrospectives, metrics)
- Non-HELIX Content: 3 files (project README, getting started, tutorials)

🔍 Quality Assessment:
- Well-Organized: ⚠️  Scattered across multiple locations
- Up-to-Date: ✅ Most content appears current
- Cross-Referenced: ❌ Many broken internal links found

📁 CREATING: HELIX documentation structure in docs/helix/

📋 EXECUTING: Frame phase consolidation
- Moving user-stories.md to docs/helix/01-frame/requirements/
- Moving business-requirements.md to docs/helix/01-frame/requirements/
- Moving acceptance-criteria.md to docs/helix/01-frame/acceptance/
- Updating 15 internal links to new locations

📋 EXECUTING: Design phase consolidation
- Moving architecture-overview.md to docs/helix/02-design/architecture/
- Moving api-specification.md to docs/helix/02-design/interfaces/
- Moving ADR-001-database-choice.md to docs/helix/02-design/decisions/
- Updating 23 internal links and cross-references

[Continues systematically through all phases...]

🔗 VALIDATING: Link integrity and cross-references
- Updated 67 internal links to new locations
- Verified 34 cross-references between phases
- Created phase navigation READMEs

✅ CONSOLIDATION VALIDATION COMPLETE
- All 47 files successfully consolidated
- HELIX structure properly implemented
- No content lost, all links functional
- Documentation now discoverable and well-organized

📚 SUMMARY: Documentation consolidation complete
- HELIX structure: docs/helix/ with 6 phase directories
- Non-HELIX content: Organized in docs/ root level
- Team can now easily navigate and contribute to documentation
- Clear separation between workflow artifacts and general docs
```

## Integration with HELIX Principles

This systematic consolidation approach ensures:
- **Phase-Based Organization**: Content aligned with HELIX workflow phases
- **Progressive Disclosure**: Information flows logically from Frame to Iterate
- **Separation of Concerns**: HELIX workflow docs separate from general project docs
- **Discoverability**: Clear structure makes documentation easy to find
- **Maintainability**: Organized structure supports ongoing documentation efforts
- **Team Adoption**: Well-organized docs encourage team contribution

Remember: Effective documentation consolidation requires systematic analysis, careful content mapping, and meticulous attention to maintaining information integrity. The goal is to create a documentation structure that supports both HELIX workflow execution and general project understanding while making information easily discoverable and maintainable for the entire team.