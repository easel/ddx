---
tags: [workflows, creation, guide, tutorial, patterns]
aliases: ["Creating Workflows", "Workflow Creation Guide", "How to Create Workflows"]
created: 2025-01-12
modified: 2025-01-12
---

# Creating DDX Workflows

This comprehensive guide explains how to create effective DDX workflows, capturing the philosophy and practical approach developed through the DDX project itself.

## Philosophy: The Three Layers

DDX workflows operate on three conceptual layers, each serving a distinct purpose:

### 1. Templates as Structure (The "What")
Templates are pure structural documents - skeletons with placeholders. They define WHAT needs to be created without any intelligence about HOW to create it. Think of them as forms waiting to be filled out.

**Example Template:**
```markdown
# [Document Title]

## Section One
[Content for section one]

## Section Two
[Content for section two]
```

### 2. Prompts as Intelligence (The "How")
Prompts embed or reference templates and add the intelligence layer. They know HOW to gather information and fill templates appropriately. Prompts guide AI assistants through the process of creating content.

**Example Prompt:**
```markdown
# Document Creation Assistant

## Template
{{include: template.md}}

## Guidance
For Section One, consider:
- What is the main purpose?
- Who is the audience?
- What key points must be covered?
```

### 3. Patterns as Methodology (The "Why")
Patterns document the overall approach - explaining WHY this workflow exists and WHEN to use it. They provide context and best practices.

## Anatomy of a Workflow

Every workflow follows this standardized structure:

```
workflows/{workflow-name}/
├── README.md                    # Pattern documentation (the "why")
├── workflow.yml                 # Metadata and automation
├── GUIDE.md                    # Comprehensive usage guide
├── {artifact-name}/            # Each major artifact/output
│   ├── README.md              # About this specific artifact
│   ├── template.md            # The structural skeleton
│   ├── prompt.md              # AI guidance for filling template
│   ├── meta.yml               # Artifact metadata
│   └── examples/              # Real-world examples
│       └── example-1.md
└── phases/                     # Workflow phases
    ├── 01-phase-name.md
    ├── 02-phase-name.md
    └── ...
```

## Step-by-Step Creation Process

### Step 1: Define the Workflow Purpose

Start by answering fundamental questions:
- **What problem does this workflow solve?**
- **Who will use this workflow?**
- **When should someone use this workflow?**
- **What are the expected outcomes?**

Document these answers in your workflow's `README.md`.

### Step 2: Identify Workflow Phases

Break down the workflow into logical phases. Each phase should:
- Have clear entry criteria
- Produce specific outputs
- Have defined exit criteria
- Connect logically to the next phase

Example for a development workflow:
1. **Define** - Establish requirements and vision
2. **Design** - Create technical architecture
3. **Implement** - Build the solution
4. **Test** - Validate functionality
5. **Release** - Deploy to production
6. **Iterate** - Gather feedback and improve

### Step 3: Identify Artifacts

For each phase, identify the key artifacts (documents, code, configurations) that need to be created. Each artifact becomes a subdirectory in your workflow.

Common artifacts:
- Requirements documents
- Design specifications
- Architecture decisions
- Test plans
- Release notes

### Step 4: Create Templates

For each artifact, create a template that defines its structure:

```markdown
# {artifact-name}/template.md
---
tags: [template, {workflow-name}, {artifact-name}]
template: true
---

# [Title]

## Section 1: [Name]
[Placeholder for content]

## Section 2: [Name]
[Placeholder for content]

## Metadata
- Created: {{date}}
- Author: {{author}}
- Status: {{status}}
```

### Step 5: Create Prompts

For each template, create a corresponding prompt that adds intelligence:

```markdown
# {artifact-name}/prompt.md
---
tags: [prompt, {workflow-name}, {artifact-name}]
references: template.md
---

# {Artifact Name} Creation Assistant

## Template Structure
This prompt helps you complete the {Artifact Name} using the template at [[template.md]].

## Information Gathering

### Section 1: [Name]
To complete this section, provide:
- [Specific question 1]
- [Specific question 2]
- [Guidance for what makes a good answer]

### Section 2: [Name]
Consider the following:
- [Key consideration 1]
- [Key consideration 2]

## Best Practices
- [Specific guidance for this artifact]
- [Common pitfalls to avoid]
- [Quality criteria]

## Template
{{include: template.md}}
```

### Step 6: Create Workflow Metadata

Define the workflow structure in `workflow.yml`:

```yaml
# workflow.yml
name: {workflow-name}
version: 1.0.0
description: Brief description of the workflow
author: Your name or team
tags:
  - category1
  - category2

phases:
  - id: phase1
    name: Phase One Name
    description: What happens in this phase
    artifacts:
      - artifact1
      - artifact2
    entry_criteria:
      - Criteria 1
      - Criteria 2
    exit_criteria:
      - Criteria 1
      - Criteria 2
    next: phase2

  - id: phase2
    name: Phase Two Name
    # ... continued

artifacts:
  - id: artifact1
    name: Artifact One Name
    type: document
    template: artifact1/template.md
    prompt: artifact1/prompt.md
    required: true
    
automation:
  init_command: ddx workflow init {workflow-name}
  validate_command: ddx workflow validate
```

### Step 7: Add Examples

For each artifact, provide at least one real-world example:

```markdown
# {artifact-name}/examples/example-1.md
---
tags: [example, {workflow-name}, {artifact-name}]
context: Brief description of this example's context
---

# Real Example Title

[Actual filled-out version of the template]
```

### Step 8: Document Phases

Create detailed documentation for each phase:

```markdown
# phases/01-{phase-name}.md
---
tags: [phase, {workflow-name}, {phase-name}]
order: 1
---

# Phase: {Phase Name}

## Overview
What this phase accomplishes and why it's important.

## Entry Criteria
- [ ] Prerequisite 1
- [ ] Prerequisite 2

## Activities
1. **Activity 1**: Description
2. **Activity 2**: Description

## Artifacts Produced
- [[../{artifact1}/README|Artifact 1]]: Purpose
- [[../{artifact2}/README|Artifact 2]]: Purpose

## Exit Criteria
- [ ] Completion criteria 1
- [ ] Completion criteria 2

## Common Challenges
- Challenge 1: How to address
- Challenge 2: How to address

## Next Phase
[[02-{next-phase}|Next Phase Name]]
```

### Step 9: Create the Comprehensive Guide

Write a `GUIDE.md` that ties everything together:

```markdown
# {Workflow Name}: Comprehensive Guide
---
tags: [guide, {workflow-name}]
---

## Introduction
Detailed explanation of this workflow's purpose and value.

## Prerequisites
What users need before starting this workflow.

## Workflow Overview
Visual or textual representation of the entire flow.

## Phase-by-Phase Walkthrough
### Phase 1: {Name}
Detailed guidance for executing this phase...

## Tips and Best Practices
Lessons learned from real-world usage.

## Troubleshooting
Common issues and their solutions.

## Case Studies
Real examples of this workflow in action.
```

## Design Principles

### 1. Separation of Concerns
- **Templates** handle structure only
- **Prompts** handle intelligence only
- **Documentation** handles explanation only

### 2. Progressive Enhancement
- Templates work standalone (can be filled manually)
- Prompts enhance templates with AI assistance
- Automation enhances the entire workflow

### 3. Reusability
- Templates can be reused across workflows
- Prompts can reference common patterns
- Examples serve as learning resources

### 4. Discoverability
- Clear naming conventions
- Comprehensive tagging
- Cross-linking between related content

## Validation Checklist

Before publishing your workflow, ensure:

- [ ] **Structure**: All required directories and files exist
- [ ] **Documentation**: README, GUIDE, and phase docs are complete
- [ ] **Templates**: All templates have clear placeholders
- [ ] **Prompts**: All prompts reference their templates
- [ ] **Examples**: At least one example per artifact
- [ ] **Metadata**: workflow.yml is valid and complete
- [ ] **Cross-links**: All internal links work
- [ ] **Tags**: Consistent tagging throughout
- [ ] **Testing**: Workflow has been tested end-to-end

## Using the Bootstrap Prompt

DDX provides a bootstrap prompt to help create new workflows:

```bash
ddx apply prompts/ddx/create-workflow
```

This interactive prompt will:
1. Gather information about your workflow
2. Create the directory structure
3. Generate starter templates
4. Set up metadata files

## Common Patterns

### Sequential Workflow
Phases must be completed in order:
```yaml
phases:
  - id: phase1
    next: phase2
  - id: phase2
    next: phase3
    requires: [phase1]
```

### Parallel Workflow
Some phases can be done simultaneously:
```yaml
phases:
  - id: planning
    next: [design, research]
  - id: design
    requires: [planning]
  - id: research
    requires: [planning]
```

### Iterative Workflow
Phases can loop back:
```yaml
phases:
  - id: implement
    next: test
  - id: test
    next: [release, implement]  # Can loop back
```

## Workflow Naming Conventions

### Directory Names
- Use lowercase with hyphens: `my-workflow`
- Be descriptive: `incident-response` not `ir`
- Avoid version numbers in names

### File Names
- Templates: `template.md`
- Prompts: `prompt.md`
- Examples: `descriptive-name.md`
- Phases: `01-phase-name.md` (numbered for order)

### Tag Conventions
- Workflow tags: `workflow/{name}`
- Artifact tags: `{workflow-name}/{artifact}`
- Phase tags: `phase/{workflow-name}/{phase}`

## Integration with DDX Commands

Your workflow can integrate with DDX CLI commands:

### Initialization
```bash
ddx workflow init {workflow-name}
```

### Application
```bash
ddx workflow apply {workflow-name}:{artifact}
```

### Validation
```bash
ddx workflow validate {workflow-name}
```

## Contributing Your Workflow

Once created, share your workflow with the community:

1. **Test Thoroughly**: Ensure it works end-to-end
2. **Document Completely**: All sections filled out
3. **Add Examples**: Real-world usage examples
4. **Submit PR**: Via `ddx contribute`

See [[contributing-workflows|Contributing Workflows Guide]] for details.

## Advanced Topics

### Conditional Logic
Use metadata to define conditional paths:
```yaml
artifacts:
  - id: security-review
    condition: "project.type == 'public'"
```

### Variable Substitution
Support project-specific values:
```yaml
variables:
  - name: project_name
    prompt: "What is your project name?"
    default: "my-project"
```

### External Integrations
Connect to external tools:
```yaml
integrations:
  - type: github
    trigger: "phase.release.complete"
    action: "create_release"
```

## Troubleshooting

### Common Issues

**Templates too rigid**: Add more placeholders and sections
**Prompts too vague**: Add specific questions and examples
**Phases unclear**: Better define entry/exit criteria
**Examples missing context**: Add background information

### Validation Errors

Run validation to catch issues:
```bash
ddx workflow validate {workflow-name}
```

Common errors:
- Missing required files
- Broken cross-references
- Invalid YAML syntax
- Inconsistent metadata

## Next Steps

- Study the [[workflows/development/README|Development Workflow]] as a complete example
- Use the [[prompts/ddx/create-workflow|Bootstrap Prompt]] to start
- Read [[using-workflows|Using Workflows]] to understand the user perspective
- Join the community to share ideas and get feedback

Remember: The best workflows come from real-world experience. Start simple, use it yourself, iterate based on feedback, and share your improvements with the community.