# Define Phase

---
tags: [development, workflow, phase, requirements, planning]
phase: 01
name: "Define"
next_phase: "[[02-design]]"
artifacts: ["[[product-requirements-document]]", "[[user-stories]]", "[[acceptance-criteria]]"]
---

## Overview

The Define phase is the foundation of the DDX development workflow where project requirements are established, user needs are analyzed, and success criteria are defined. This phase transforms ideas into actionable specifications that guide all subsequent development activities.

## Purpose

- Establish clear project objectives and scope
- Define user requirements and success metrics
- Create alignment among stakeholders
- Provide foundation for technical design decisions
- Minimize scope creep and rework in later phases

## Entry Criteria

Before entering the Define phase, ensure:

- [ ] Project sponsor or stakeholder identified
- [ ] Basic problem statement or opportunity articulated
- [ ] Initial resources and timeline constraints understood
- [ ] Key stakeholders available for requirements gathering
- [ ] DDX toolkit initialized in project repository (`ddx init`)

## Key Activities

### 1. Requirements Gathering

- Conduct stakeholder interviews and workshops
- Document functional and non-functional requirements
- Identify user personas and use cases
- Define system boundaries and constraints
- Capture business rules and compliance requirements

### 2. User Story Development

- Write user stories in standard format: "As a [user], I want [goal] so that [benefit]"
- Prioritize stories using MoSCoW method (Must, Should, Could, Won't)
- Define acceptance criteria for each story
- Estimate story complexity and effort

### 3. Success Metrics Definition

- Establish measurable objectives (KPIs)
- Define quality attributes and performance targets
- Set user experience and usability goals
- Document technical constraints and limitations

### 4. Risk Assessment

- Identify potential project risks and dependencies
- Assess technical feasibility and complexity
- Document assumptions and external dependencies
- Create initial risk mitigation strategies

## Artifacts Produced

### Primary Artifacts

- **[[Product Requirements Document (PRD)]]** - Comprehensive requirements specification
- **[[User Stories]]** - Prioritized list of user requirements
- **[[Acceptance Criteria]]** - Testable conditions for story completion

### Supporting Artifacts

- **[[Stakeholder Analysis]]** - Key people and their interests
- **[[Risk Register]]** - Identified risks and mitigation plans
- **[[Assumptions Log]]** - Documented assumptions and constraints
- **[[Success Metrics]]** - Measurable project objectives

## Exit Criteria

The Define phase is complete when:

- [ ] Product Requirements Document approved by stakeholders
- [ ] User stories prioritized and estimated
- [ ] Acceptance criteria defined for priority stories
- [ ] Success metrics and quality goals established
- [ ] Major risks identified and assessed
- [ ] Stakeholder sign-off obtained on requirements
- [ ] Next phase (Design) entry criteria satisfied

## Common Challenges and Solutions

### Challenge: Unclear or Changing Requirements

**Solutions:**
- Use iterative requirements gathering with frequent stakeholder reviews
- Implement change control process early
- Document decisions and rationale clearly
- Use prototypes or wireframes to clarify requirements

### Challenge: Stakeholder Misalignment

**Solutions:**
- Facilitate requirements workshops with all key stakeholders
- Create shared vision documents
- Use acceptance criteria to resolve ambiguity
- Establish regular communication cadence

### Challenge: Scope Creep

**Solutions:**
- Define clear project boundaries and constraints
- Use MoSCoW prioritization to manage feature requests
- Implement formal change request process
- Maintain traceability between requirements and business objectives

## Tips and Best Practices

### Requirements Elicitation

- Ask "why" questions to understand underlying needs
- Use multiple elicitation techniques (interviews, workshops, observation)
- Document both what users say and what they actually need
- Validate requirements with multiple stakeholders

### Documentation Standards

- Use templates and standards for consistency
- Keep requirements atomic and testable
- Link requirements to business objectives
- Maintain traceability throughout development

### Stakeholder Management

- Identify decision makers and influencers early
- Create RACI matrix for requirements decisions
- Schedule regular review cycles
- Use collaborative tools for transparency

### Quality Assurance

- Review requirements for completeness and consistency
- Validate requirements against business objectives
- Ensure all requirements are measurable or testable
- Create requirements baseline before proceeding

## DDX Integration

### Using DDX Templates

Apply relevant DDX templates for requirements documentation:

```bash
ddx apply templates/requirements/prd-template
ddx apply templates/requirements/user-stories-template
ddx apply patterns/requirements/acceptance-criteria
```

### Configuration Management

Ensure `.ddx.yml` includes requirements artifacts:

```yaml
artifacts:
  requirements:
    - requirements/prd.md
    - requirements/user-stories.md
    - requirements/acceptance-criteria.md
```

### Quality Gates

Use DDX diagnostics to validate phase completion:

```bash
ddx diagnose --phase define
```

## Next Phase

Upon successful completion of the Define phase, proceed to **[[02-design|Design Phase]]** where technical architecture and system design will be developed based on the established requirements.

---

*This document is part of the DDX Development Workflow. For questions or improvements, use `ddx contribute` to share feedback with the community.*