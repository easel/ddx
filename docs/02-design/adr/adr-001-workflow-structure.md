---
tags: [architecture, adr, decision, workflow, design]
aliases: ["ADR-001", "Workflow Architecture Decision", "Template Prompt Pattern"]
created: 2025-01-12
modified: 2025-01-12
status: accepted
---

# ADR-001: Workflow Structure and Template-Prompt-Pattern Architecture

**Status**: Accepted  
**Date**: 2025-01-12  
**Deciders**: DDX Architecture Team  
**Technical Story**: [Design DDX workflow system architecture](https://github.com/ddx-dev/ddx/issues/workflow-architecture)

## Summary

We have decided to structure DDX workflows using a three-layer architecture: Templates as Structure, Prompts as Intelligence, and Patterns as Methodology. This decision establishes the fundamental design principles for how DDX workflows are created, organized, and used.

## Context

During the development of DDX (Document-Driven Development eXperience), we needed to design a system that would:

1. **Support AI-Assisted Development**: Enable effective collaboration between humans and AI assistants
2. **Maintain Flexibility**: Allow workflows to be used with or without AI assistance
3. **Ensure Reusability**: Make workflow components reusable across different contexts
4. **Enable Progressive Enhancement**: Support users from beginner to expert level
5. **Preserve Knowledge**: Capture and share development best practices
6. **Follow DDX Medical Metaphor**: Align with the diagnostic/treatment protocol theme

### Key Requirements

- Workflows must work both with AI assistants and for manual completion
- Components should be independently useful and composable
- Structure should be intuitive for developers from different backgrounds
- System should scale from individual use to enterprise adoption
- Architecture should support community contributions and sharing

### Constraints

- Must integrate with existing development tools and processes
- Should leverage standard formats (Markdown, YAML) for broad compatibility
- Must support version control and collaborative editing
- Architecture should be extensible for future enhancements

## Decision Drivers

### 1. AI Integration Philosophy
The rise of AI coding assistants changed how we think about documentation and process guidance. We needed a structure that maximizes AI effectiveness while preserving human agency.

### 2. Separation of Concerns
Through experimentation, we found that mixing structure, intelligence, and methodology in single documents created several problems:
- Templates became bloated with instructions
- AI prompts were cluttered with structural concerns  
- Documentation was hard to maintain and update
- Reusability was limited

### 3. Progressive Enhancement Principle
We wanted to support users at different skill levels:
- **Beginners**: Need step-by-step guidance
- **Intermediates**: Want customizable templates
- **Experts**: Need building blocks for custom solutions

### 4. Community Scalability
For DDX to succeed as a community platform, we needed:
- Clear contribution guidelines
- Consistent quality standards
- Easy discoverability and sharing
- Composable, reusable components

## Decision

We adopt a **three-layer workflow architecture**:

### Layer 1: Templates as Structure (The "What")

Templates are pure structural documents - skeletons with placeholders that define WHAT needs to be created without any intelligence about HOW to create it.

**Characteristics:**
- Pure Markdown with placeholders
- No AI-specific instructions
- Usable standalone (human-fillable)
- Focused on information architecture
- Reusable across different workflows

**Example:**
```markdown
# Product Requirements Document: [Product Name]

## Problem Statement
[Description of the problem this product solves]

## Target Users
[Primary user personas and their needs]

## Solution Overview
[High-level description of the proposed solution]

## Success Metrics
[How success will be measured]
```

### Layer 2: Prompts as Intelligence (The "How")

Prompts embed or reference templates and add the intelligence layer. They know HOW to gather information and fill templates appropriately, providing guidance for both AI assistants and humans.

**Characteristics:**
- Include or reference the template
- Contain AI-specific instructions
- Provide question frameworks
- Include best practices and quality criteria
- Guide information gathering process

**Example:**
```markdown
# PRD Creation Assistant

## Template Structure
This prompt helps you complete a Product Requirements Document using the template below.

## Information Gathering

### Problem Statement
To complete this section effectively:
- Interview 3-5 potential users about their current pain points
- Research existing solutions and their limitations  
- Quantify the problem impact (time, cost, frustration)

### Target Users
Consider these questions:
- Who experiences this problem most acutely?
- What are their current workarounds?
- What's their decision-making process?

## Quality Criteria
A strong PRD should:
- Be specific enough that engineering can estimate effort
- Include measurable success criteria
- Address potential edge cases and failure modes

## Template
{{include: prd-template.md}}
```

### Layer 3: Patterns as Methodology (The "Why")

Patterns document the overall approach, explaining WHY this workflow exists, WHEN to use it, and HOW it fits into the larger development process.

**Characteristics:**
- Explain the methodology and philosophy
- Provide context for when to use the workflow
- Document best practices and lessons learned
- Connect to related workflows and practices
- Include case studies and examples

**Example:**
```markdown
# Product Requirements Pattern

## Purpose
Product Requirements Documents (PRDs) serve as the single source of truth for what we're building and why. They bridge the gap between business needs and technical implementation.

## When to Use
Create a PRD when:
- Building a new feature or product
- Making significant changes to existing functionality
- Multiple teams need to coordinate on implementation
- Stakeholder alignment is critical for success

## Methodology
The PRD process follows these principles:
1. **User-Centric**: Start with user problems, not solutions
2. **Data-Driven**: Base decisions on research and metrics
3. **Iterative**: Refine based on feedback and learning
4. **Collaborative**: Involve all relevant stakeholders

## Case Studies
- **E-commerce Checkout**: How Acme Corp reduced cart abandonment by 30%
- **Mobile Onboarding**: Streamlining new user experience
```

### Supporting Components

#### Metadata and Automation (workflow.yml)
YAML files define workflow structure, phases, dependencies, and automation hooks.

#### Examples and Case Studies
Real-world implementations that demonstrate successful usage patterns.

#### Phase Documentation
Detailed guidance for each workflow phase, including entry/exit criteria and activities.

## Rationale

### Why This Three-Layer Architecture?

#### 1. **Separation of Concerns**
Each layer has a single, clear responsibility:
- Templates handle structure only
- Prompts handle intelligence only  
- Patterns handle methodology only

This separation makes each component easier to create, maintain, and reuse.

#### 2. **Progressive Enhancement**
The architecture supports multiple usage modes:
- **Manual completion**: Use templates directly
- **AI-assisted**: Use prompts with AI assistants
- **Full automation**: Use patterns for process guidance

Users can engage at their comfort level and gradually adopt more sophisticated approaches.

#### 3. **Composability**
Components can be mixed and matched:
- One template can be used with multiple prompts
- One prompt can reference multiple templates
- Patterns can combine multiple template/prompt pairs

#### 4. **AI Optimization**
This structure maximizes AI assistant effectiveness:
- Templates provide clear structure for output
- Prompts provide context and guidance for generation
- Separation prevents AI confusion between structure and content

#### 5. **Human Usability**
The architecture remains human-friendly:
- Templates are readable and fillable without AI
- Prompts provide clear guidance for manual completion
- Patterns explain the reasoning and context

### Why Not Alternative Approaches?

#### Monolithic Documents
**Rejected**: Combining structure, instructions, and methodology in single documents
- **Problem**: Hard to maintain, reuse, and customize
- **Problem**: Overwhelming for beginners, limiting for experts
- **Problem**: Poor AI performance due to mixed concerns

#### Pure AI Prompts
**Rejected**: Making everything AI-specific prompts
- **Problem**: Not usable without AI assistance
- **Problem**: Hard to version control and review
- **Problem**: Limited customization and reusability

#### Traditional Documentation
**Rejected**: Static documentation with examples
- **Problem**: Doesn't scale for complex processes
- **Problem**: Doesn't integrate with AI workflows
- **Problem**: Hard to keep current and consistent

#### Procedural Scripts
**Rejected**: Fully automated generation tools
- **Problem**: Too rigid for creative/strategic work
- **Problem**: No room for human insight and adaptation
- **Problem**: Hard to debug and modify

## Consequences

### Positive Consequences

#### For Users
- **Flexibility**: Can use components independently or together
- **Learning Path**: Clear progression from manual to AI-assisted workflows
- **Customization**: Easy to adapt components for specific needs
- **Consistency**: Standardized structure across all workflows

#### For Contributors
- **Clear Guidelines**: Obvious how to create new workflow components
- **Reusability**: Components can be shared across workflows
- **Quality Standards**: Separation of concerns improves quality
- **Community**: Easier to review and improve contributions

#### For AI Assistants
- **Better Performance**: Clear structure improves AI output quality
- **Context Clarity**: Prompts provide necessary context and constraints
- **Consistency**: Standardized format improves reliability
- **Extensibility**: Easy to add new AI capabilities

#### For Organizations
- **Scalability**: Architecture scales from individual to enterprise use
- **Integration**: Works with existing tools and processes
- **Knowledge Capture**: Systematic way to document best practices
- **Training**: Clear structure for onboarding new team members

### Negative Consequences

#### Complexity
- **More Files**: Three-layer architecture requires more files per workflow
- **Learning Curve**: Users must understand the three-layer concept
- **Maintenance**: More components require more maintenance effort

*Mitigation*: Clear documentation, good tooling, and community support reduce complexity burden.

#### Potential Inconsistency
- **Component Drift**: Templates and prompts might become inconsistent over time
- **Quality Variation**: Different contributors might interpret standards differently

*Mitigation*: Automated validation tools, clear guidelines, and review processes maintain consistency.

#### Tool Dependency
- **DDX CLI Requirement**: Full benefits require DDX tooling
- **Template Engine**: Variable substitution requires processing

*Mitigation*: Components work standalone; DDX CLI enhances but doesn't replace manual usage.

### Trade-offs Accepted

#### Simplicity vs. Flexibility
We chose flexibility over simplicity, accepting that the three-layer architecture is more complex than alternatives but provides significantly more value.

#### Immediate vs. Long-term Value
We optimized for long-term community value over immediate ease of adoption, believing that proper architecture would compound benefits over time.

#### Human vs. AI Optimization
We designed for both human and AI usage, accepting some complexity to avoid forcing users to choose between approaches.

## Implementation Details

### Directory Structure
```
workflows/{workflow-name}/
├── README.md           # Pattern documentation
├── workflow.yml        # Metadata and automation
├── GUIDE.md           # Comprehensive guide
├── {artifact}/        # Each major artifact
│   ├── README.md     # About this artifact
│   ├── template.md   # The skeleton
│   ├── prompt.md     # AI assistance
│   └── examples/     # Real examples
└── phases/           # Phase documentation
```

### File Naming Conventions
- **Templates**: `template.md` (consistent naming)
- **Prompts**: `prompt.md` (pairs with template)
- **Patterns**: `README.md` (overview document)
- **Examples**: `descriptive-name.md` (specific context)

### Cross-Reference Standards
- Templates are referenced by prompts
- Prompts include or embed templates
- Patterns link to both templates and prompts
- Examples demonstrate real usage

### Validation Requirements
- Templates must have valid placeholder syntax
- Prompts must reference their templates
- Patterns must explain usage context
- All components must be cross-linked

## Validation

### Success Metrics

#### Adoption Metrics
- Number of workflows created using this architecture
- Community contributions following the pattern
- User feedback on ease of use and effectiveness

#### Quality Metrics
- Consistency across workflow components
- AI assistant performance with prompt/template pairs
- User success rate with workflow completion

#### Maintenance Metrics
- Time to create new workflows
- Effort required to maintain existing workflows
- Community contribution rate and quality

### Testing Approach

#### Component Testing
- Templates generate valid output when filled
- Prompts work effectively with AI assistants
- Patterns provide clear usage guidance

#### Integration Testing
- Complete workflows function end-to-end
- Components work together seamlessly
- Customization doesn't break functionality

#### User Testing
- New users can successfully follow workflows
- Experienced users can customize effectively
- AI assistants produce quality output

### Monitoring Plan

#### Continuous Feedback
- Regular user surveys and interviews
- Community forum monitoring
- Usage analytics where available

#### Periodic Reviews
- Quarterly architecture review meetings
- Annual major version planning
- Community retrospectives

## Related Decisions

### Future ADRs
This decision establishes foundation for future architectural decisions:
- **ADR-002**: Template variable substitution system
- **ADR-003**: Workflow validation and quality standards  
- **ADR-004**: Community contribution and review processes
- **ADR-005**: AI assistant integration patterns

### Dependencies
This decision builds on:
- DDX medical metaphor concept
- Markdown-based documentation standards
- Git-based version control requirements
- Community-driven development model

## References

### Research and Analysis
- Analysis of existing workflow tools (GitHub templates, Cookiecutter, Yeoman)
- Study of AI assistant prompt engineering best practices
- Review of software development process documentation
- Community feedback from early DDX adopters

### Standards and Guidelines
- [Markdown specification](https://spec.commonmark.org/)
- [YAML specification](https://yaml.org/spec/)
- [Architecture Decision Records](https://adr.github.io/)
- [Semantic versioning](https://semver.org/)

### Tools and Technologies
- Obsidian for knowledge management
- Git for version control
- GitHub for collaboration
- Various AI assistants for testing

## Medical Metaphor Alignment

This architectural decision aligns with DDX's medical metaphor:

- **Templates** = Medical Forms (standard structure for capturing information)
- **Prompts** = Clinical Guidelines (expert guidance for completion)
- **Patterns** = Treatment Protocols (proven methodologies for specific conditions)
- **Workflows** = Care Pathways (complete treatment plans)
- **Examples** = Case Studies (real-world applications and outcomes)

Just as medical practice separates forms, guidelines, and protocols while ensuring they work together seamlessly, DDX separates structure, intelligence, and methodology while maintaining tight integration.

## Revision History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-01-12 | Initial version based on DDX workflow architecture discussions |

---

**Next Review Date**: 2025-04-12 (3 months)  
**Review Trigger**: Major community feedback or technical limitations discovered