---
tags: [workflow, cdp, clinical-development-protocol, pattern, methodology, software-development]
aliases: ["Clinical Development Protocol", "CDP Workflow", "Medical Development Process", "CDP Pattern"]
created: 2025-01-12
modified: 2025-01-12
---

# Clinical Development Protocol (CDP) Pattern

## Overview

The Clinical Development Protocol (CDP) is DDX's systematic workflow for rigorous software development that follows medical best practices. It implements a patient-centered approach to building software, treating each project as a patient requiring careful diagnosis, treatment, and ongoing care.

This workflow follows the medical principle: **Diagnose → Prescribe → Treat → Monitor → Follow-up**

## When to Use This Workflow

Use the Clinical Development Protocol when:
- Building mission-critical software systems
- Working in regulated industries (healthcare, finance, aerospace)
- Developing software where quality and safety are paramount
- Establishing rigorous development practices for high-stakes projects
- Building products that require comprehensive validation and documentation

## Core Philosophy

### Patient-Centered Development
This workflow treats every software project as a patient requiring careful attention, proper diagnosis, and appropriate treatment. Each phase represents a critical step in the care continuum.

### Evidence-Based Practice
Every decision is backed by documentation and validation, just as medical treatments are based on clinical evidence. Documentation isn't bureaucracy—it's patient safety.

### Continuous Monitoring
Like ongoing patient care, software requires continuous monitoring and adjustment based on real-world performance and changing conditions.

### Quality as Life-Critical
Quality isn't optional—it's essential for patient (user) safety and wellbeing. Every phase includes validation gates to ensure standards are met.

## CDP Principles

### 1. Primum Non Nocere (First, Do No Harm)
Every change must be carefully evaluated for potential negative impact before implementation.

### 2. Evidence-Based Decision Making
All architectural and design decisions must be documented with clear rationale and supporting evidence.

### 3. Comprehensive Validation
Each phase requires validation that criteria are met before proceeding to the next treatment stage.

### 4. Continuous Care
Software, like patients, requires ongoing monitoring and care throughout its lifecycle.

## Workflow Phases

### 1. Diagnose Phase
**Purpose**: Thoroughly analyze the problem and establish clear requirements
**Key Artifact**: [[prd/README|Product Requirements Document (Diagnosis)]]
**Duration**: 1-2 weeks typically

During this phase, you:
- Conduct thorough problem analysis
- Identify all stakeholders and their needs
- Establish measurable success criteria
- Document the comprehensive diagnosis

### 2. Prescribe Phase
**Purpose**: Design the treatment plan and technical architecture
**Key Artifact**: [[architecture/README|Architecture Treatment Plan]]
**Duration**: 1-2 weeks typically

During this phase, you:
- Design comprehensive system architecture
- Create detailed treatment specifications
- Document all architectural decisions
- Plan implementation approach

### 3. Treat Phase
**Purpose**: Implement the prescribed solution with precision
**Key Artifacts**: [[feature-spec/README|Treatment Implementation Records]]
**Duration**: Variable based on complexity

During this phase, you:
- Execute treatment plan systematically
- Document each implementation step
- Conduct peer reviews (medical consultations)
- Maintain comprehensive treatment records

### 4. Monitor Phase
**Purpose**: Validate treatment effectiveness through comprehensive testing
**Key Artifact**: [[test-plan/README|Monitoring and Assessment Plans]]
**Duration**: 25-35% of treatment time

During this phase, you:
- Execute comprehensive monitoring protocols
- Validate all treatment objectives
- Document any adverse effects (bugs)
- Ensure patient safety criteria are met

### 5. Follow-up Phase
**Purpose**: Deploy to production with ongoing care planning
**Key Artifact**: [[release/README|Treatment Summary and Care Plan]]
**Duration**: 1-3 days typically

During this phase, you:
- Execute careful deployment protocol
- Establish ongoing monitoring systems
- Document treatment outcomes
- Plan follow-up care schedule

### 6. Continuing Care Phase
**Purpose**: Ongoing monitoring and iterative improvements
**Duration**: Continuous

During this phase, you:
- Monitor patient (system) health
- Analyze treatment effectiveness
- Identify improvement opportunities
- Plan next treatment cycle

## Validation Gates

Each phase transition requires validation through specific gates:

### Pre-Prescribe Gate
- Diagnosis completeness verified
- All stakeholders have reviewed and approved
- Success criteria are measurable and realistic
- Technical feasibility confirmed

### Pre-Treatment Gate
- Architecture treatment plan approved by senior practitioners
- All technology choices justified with evidence
- Risk assessment completed and mitigation planned
- Implementation approach validated

### Pre-Monitor Gate
- All prescribed treatments implemented
- Peer review (code review) completed
- Unit-level validation successful
- Implementation documentation complete

### Pre-Follow-up Gate
- Comprehensive monitoring completed
- All critical and high-severity issues resolved
- Performance and safety criteria validated
- Treatment effectiveness demonstrated

### Continuing Care Gate
- Production deployment successful
- Monitoring systems operational
- Initial health indicators positive
- Care transition plan documented

## Key Principles

### 1. Validation Gates
Each phase has rigorous entry and exit criteria with validation requirements. No phase may be bypassed without proper justification and approval.

### 2. Living Medical Records
All artifacts serve as comprehensive medical records that evolve throughout the project lifecycle. They provide complete treatment history.

### 3. Chain of Evidence
Requirements trace to designs, designs trace to implementations, implementations trace to tests. This creates an unbroken chain of evidence.

### 4. Continuous Assessment
Each phase validates the work of previous phases through systematic assessment and peer review.

## Success Metrics

A successful CDP implementation results in:
- **Predictable Outcomes**: Consistent, reliable delivery of high-quality solutions
- **Zero Critical Failures**: Mission-critical systems operate without life-threatening issues
- **Complete Traceability**: Full audit trail from requirements to deployment
- **Knowledge Preservation**: Comprehensive records enable future care decisions
- **Continuous Improvement**: Each treatment cycle improves patient outcomes

## Anti-Patterns to Avoid

### Rushing to Treatment
Don't skip thorough diagnosis to "start coding faster"—incomplete diagnosis leads to ineffective treatment.

### Post-Treatment Documentation
Don't treat documentation as post-treatment paperwork. Medical records are created during treatment, not after.

### Ignoring Adverse Events
Don't ignore bugs or issues found during monitoring. Address all adverse events immediately.

### Abandoning Follow-up Care
Don't consider treatment complete at deployment. Ongoing care is essential for patient health.

## Customization Points

This protocol can be adapted for:
- **Criticality Level**: Adjust validation rigor based on system criticality
- **Team Experience**: More oversight for junior practitioners
- **Regulatory Environment**: Additional compliance requirements
- **Patient Population**: Different user needs require different care approaches

## Getting Started

1. **Initialize the protocol**: `ddx workflow init cdp`
2. **Begin with Diagnosis**: Create comprehensive problem analysis
3. **Follow validation gates**: Don't skip assessment checkpoints
4. **Use the templates**: Leverage clinical templates and protocols
5. **Document everything**: Maintain complete medical records

## Artifacts Overview

| Artifact | Phase | Purpose |
|----------|-------|---------|
| [[prd/README\|Diagnosis]] | Diagnose | Comprehensive problem analysis |
| [[architecture/README\|Treatment Plan]] | Prescribe | Detailed technical treatment approach |
| [[feature-spec/README\|Treatment Records]] | Treat | Implementation documentation |
| [[test-plan/README\|Monitoring Plan]] | Monitor | Comprehensive validation strategy |
| [[release/README\|Care Summary]] | Follow-up | Treatment outcomes and care plan |

## Integration with DDX

This protocol integrates with DDX's medical metaphor:
- **Templates**: Clinical templates for consistent documentation
- **Prompts**: AI assistance following medical best practices
- **Patterns**: Evidence-based development patterns
- **Tools**: CLI commands for protocol automation

## Related Protocols

- **Emergency Response Protocol**: For critical production issues
- **Preventive Care Protocol**: For maintenance and optimization
- **Specialist Consultation Protocol**: For complex technical decisions

## Further Reading

- [[GUIDE|Clinical Development Guide]]: Comprehensive implementation guide
- [[phases/01-diagnose|Phase Protocols]]: Detailed phase documentation
- [[docs/usage/workflows/using-workflows|Using Protocols]]: General protocol usage
- [[docs/product/prd-ddx-v1|DDX Diagnosis]]: Example diagnosis using this protocol

## Medical Metaphor Integration

In DDX's clinical approach:
- **This Protocol** = Comprehensive Patient Care Plan
- **Phases** = Treatment Stages
- **Artifacts** = Medical Records and Documentation
- **Validation Gates** = Clinical Checkpoints
- **Continuing Care** = Long-term Patient Management

Just as medical protocols ensure consistent patient care and safety, the Clinical Development Protocol ensures consistent software quality and reliability.