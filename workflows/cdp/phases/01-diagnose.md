# Diagnose Phase

---
tags: [cdp, workflow, phase, diagnosis, problem-specification, symptoms]
phase: 01
name: "Diagnose"
next_phase: "[[02-prescribe]]"
artifacts: ["[[symptom-analysis]]", "[[problem-specification]]", "[[diagnostic-criteria]]"]
---

## Overview

The Diagnose phase is the foundation of the CDP (Clinical Development Process) workflow where problems are systematically identified, symptoms are analyzed, and diagnostic criteria are established. This phase transforms observed issues into measurable specifications that guide all subsequent treatment activities.

## Purpose

- Identify and analyze observable symptoms and problems
- Establish measurable diagnostic criteria
- Define problem boundaries and impact assessment
- Create alignment among stakeholders on problem definition
- Provide foundation for solution prescription decisions
- Minimize misdiagnosis and inappropriate treatments in later phases

## Entry Criteria

Before entering the Diagnose phase, ensure:

- [ ] Problem stakeholder or issue reporter identified
- [ ] Basic symptom observation or issue manifestation documented
- [ ] Initial constraints and environmental context understood
- [ ] Key stakeholders available for symptom analysis
- [ ] DDX toolkit initialized in project repository (`ddx init`)

## Key Activities

### 1. Symptom Analysis

- Conduct stakeholder interviews to identify observable symptoms
- Document functional and behavioral manifestations
- Identify affected user personas and use cases
- Define system boundaries and environmental constraints
- Capture compliance and regulatory symptom requirements

### 2. Problem Specification Development

- Write problem statements in standard format: "When [context], [stakeholder] experiences [symptom] resulting in [impact]"
- Prioritize symptoms using severity assessment (Critical, High, Medium, Low)
- Define measurable diagnostic criteria for each symptom
- Estimate symptom complexity and treatment effort
- Establish success criteria for symptom resolution

### 3. Measurable Symptoms Definition

- Establish quantifiable symptom indicators (KPIs)
- Define quality degradation and performance impact metrics
- Set user experience and usability impact measurements
- Document technical constraints and system limitations
- Create baseline measurements for symptom severity

### 4. Risk and Impact Assessment

- Identify potential risks of leaving symptoms untreated
- Assess technical complexity and diagnostic uncertainty
- Document assumptions and external dependencies
- Create initial risk mitigation strategies
- Evaluate business and operational impact

## Artifacts Produced

### Primary Artifacts

- **[[Symptom Analysis]]** - Comprehensive documentation of observable problems
- **[[Problem Specification]]** - Detailed problem statements with measurable criteria
- **[[Diagnostic Criteria]]** - Testable conditions for problem validation

### Supporting Artifacts

- **[[Stakeholder Impact Analysis]]** - Affected parties and severity assessment
- **[[Risk Register]]** - Identified risks and mitigation plans
- **[[Environmental Assessment]]** - Context and constraints documentation
- **[[Success Metrics]]** - Measurable problem resolution objectives

## Exit Criteria

The Diagnose phase is complete when:

- [ ] Problem Specification Document approved by stakeholders
- [ ] Symptoms prioritized and severity assessed
- [ ] Diagnostic criteria defined for priority symptoms
- [ ] Success metrics and resolution goals established
- [ ] Major risks identified and assessed
- [ ] Stakeholder sign-off obtained on problem definition
- [ ] Next phase (Prescribe) entry criteria satisfied

## CDP Validation Requirements

### Symptom Measurability Gate

- [ ] All symptoms have quantifiable indicators
- [ ] Baseline measurements captured for severity assessment
- [ ] Success criteria defined in measurable terms
- [ ] Diagnostic tests can validate symptom presence/absence

### Specification Completeness Gate

- [ ] Problem statements follow standard CDP format
- [ ] Impact assessment includes business and technical dimensions
- [ ] Environmental constraints fully documented
- [ ] Stakeholder alignment verified through formal sign-off

### Quality Assurance Gate

- [ ] Symptoms traced to observable manifestations
- [ ] Diagnostic criteria validated for reliability
- [ ] Problem boundaries clearly defined
- [ ] Resolution success criteria established

## Common Challenges and Solutions

### Challenge: Vague or Subjective Symptoms

**Solutions:**
- Use quantitative measurement techniques where possible
- Implement objective diagnostic tests and validation
- Document symptoms with concrete examples and scenarios
- Use multiple stakeholder perspectives to validate observations

### Challenge: Stakeholder Misalignment on Problem Definition

**Solutions:**
- Facilitate diagnostic workshops with all key stakeholders
- Create shared symptom documentation and evidence
- Use measurable criteria to resolve ambiguity
- Establish regular validation and review cycles

### Challenge: Scope Creep in Problem Definition

**Solutions:**
- Define clear symptom boundaries and exclusion criteria
- Use severity prioritization to manage symptom scope
- Implement formal change request process for symptom modifications
- Maintain traceability between symptoms and business objectives

## Tips and Best Practices

### Symptom Elicitation

- Ask "what" and "when" questions to understand symptom manifestations
- Use multiple diagnostic techniques (observation, measurement, analysis)
- Document both reported symptoms and underlying root causes
- Validate symptoms with multiple affected stakeholders

### Documentation Standards

- Use templates and standards for diagnostic consistency
- Keep symptom descriptions atomic and measurable
- Link symptoms to business and operational impact
- Maintain traceability throughout treatment process

### Stakeholder Management

- Identify problem reporters and affected parties early
- Create RACI matrix for diagnostic decisions
- Schedule regular review cycles for symptom validation
- Use collaborative tools for transparent documentation

### Quality Assurance

- Review symptoms for completeness and measurability
- Validate problem specifications against business objectives
- Ensure all symptoms are testable or measurable
- Create diagnostic baseline before proceeding

## DDX Integration

### Using DDX Templates

Apply relevant DDX templates for symptom documentation:

```bash
ddx apply templates/diagnosis/symptom-analysis
ddx apply templates/diagnosis/problem-specification
ddx apply patterns/diagnosis/diagnostic-criteria
```

### Configuration Management

Ensure `.ddx.yml` includes diagnostic artifacts:

```yaml
artifacts:
  diagnosis:
    - diagnosis/symptom-analysis.md
    - diagnosis/problem-specification.md
    - diagnosis/diagnostic-criteria.md
```

### Quality Gates

Use DDX diagnostics to validate phase completion:

```bash
ddx diagnose --phase diagnose
ddx diagnose --artifact symptom-measurability
ddx diagnose --artifact specification-completeness
```

## Next Phase

Upon successful completion of the Diagnose phase, proceed to **[[02-prescribe|Prescribe Phase]]** where treatment solutions and implementation approaches will be designed based on the established problem specifications and diagnostic criteria.

---

*This document is part of the DDX Clinical Development Process (CDP) Workflow. For questions or improvements, use `ddx contribute` to share feedback with the community.*