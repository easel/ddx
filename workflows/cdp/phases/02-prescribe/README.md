# Prescribe Phase

---
tags: [cdp, workflow, phase, prescription, treatment-design, solution-contracts]
phase: 02
name: "Prescribe"
previous_phase: "[[01-diagnose]]"
next_phase: "[[03-treat]]"
artifacts: ["[[treatment-plan]]", "[[solution-specification]]", "[[interface-contracts]]", "[[complexity-assessment]]"]
---

## Overview

The Prescribe phase translates diagnosed problems into detailed treatment plans and solution specifications. This phase creates the treatment blueprint, defining solution components, interfaces, implementation contracts, and complexity assessments that will guide the treatment implementation.

## Purpose

- Transform problem specifications into treatment solutions
- Define solution architecture and component interactions
- Establish interface contracts and implementation specifications
- Create detailed treatment roadmap with complexity assessment
- Minimize implementation risks and treatment failures
- Ensure solution aligns with diagnostic criteria

## Entry Criteria

Before entering the Prescribe phase, ensure:

- [ ] Diagnose phase completed with stakeholder approval
- [ ] Problem Specification Document finalized
- [ ] Symptoms and diagnostic criteria documented
- [ ] Success metrics and resolution goals established
- [ ] Environmental constraints and assumptions identified
- [ ] Treatment team assembled and briefed on diagnosis

## Key Activities

### 1. Treatment Architecture Definition

- Design overall solution architecture and treatment approach
- Define treatment component boundaries and responsibilities
- Establish interaction patterns and communication protocols
- Document treatment flows and intervention points
- Select treatment patterns (corrective, preventive, adaptive)

### 2. Solution Technology Selection

- Evaluate and select implementation languages and frameworks
- Choose databases, caching, and storage solutions for treatment
- Select development, testing, and deployment tools
- Define third-party libraries and treatment services
- Consider maintainability, scalability, and team expertise

### 3. Interface Contract Design

- Create component and service interface specifications
- Design treatment data models and schema contracts
- Define API specifications and service contracts
- Design user interface contracts and interaction flows
- Document treatment algorithms and business logic contracts

### 4. Complexity Assessment and Risk Analysis

- Assess treatment complexity and implementation difficulty
- Evaluate performance, scalability, and availability requirements
- Plan security measures and treatment validation flows
- Design for testability and treatment verification
- Consider compliance and regulatory treatment requirements

## Artifacts Produced

### Primary Artifacts

- **[[Treatment Plan]]** - Comprehensive treatment architecture and approach
- **[[Solution Specification]]** - Detailed component and interface designs
- **[[Interface Contracts]]** - API/service contracts and data specifications
- **[[Complexity Assessment]]** - Treatment difficulty and risk evaluation

### Supporting Artifacts

- **[[Treatment Architecture Document]]** - High-level solution design decisions
- **[[Technology Stack Document]]** - Selected technologies and justifications
- **[[Component Contracts]]** - Detailed interface and behavior specifications
- **[[Treatment Mockups]]** - User interface and interaction designs
- **[[Security Treatment Plan]]** - Security measures and validation flows
- **[[Deployment Treatment Strategy]]** - Infrastructure and deployment approach

## Exit Criteria

The Prescribe phase is complete when:

- [ ] Treatment plan reviewed and approved
- [ ] Solution specifications completed with interface contracts
- [ ] API and service contracts defined and validated
- [ ] Complexity assessment completed with risk mitigation
- [ ] Technology stack decisions documented and approved
- [ ] Treatment environment requirements specified
- [ ] Solution reviews completed with treatment team
- [ ] Treatment implementation plan created and estimated
- [ ] Next phase (Treat) entry criteria satisfied

## CDP Validation Requirements

### Contract Completeness Gate

- [ ] All interfaces have defined contracts with input/output specifications
- [ ] Treatment components have clear responsibility boundaries
- [ ] Service level agreements defined for performance and availability
- [ ] Error handling and failure modes specified in contracts

### Complexity Assessment Gate

- [ ] Treatment complexity scored using standardized metrics
- [ ] Risk factors identified with mitigation strategies
- [ ] Implementation effort estimated with confidence intervals
- [ ] Dependencies and integration complexity evaluated

### Solution Validation Gate

- [ ] Treatment plan addresses all diagnosed symptoms
- [ ] Solution specifications traceable to diagnostic criteria
- [ ] Interface contracts support required treatment flows
- [ ] Technology choices justified against treatment requirements

## Common Challenges and Solutions

### Challenge: Over-Engineering Treatment Solutions

**Solutions:**
- Focus on current symptoms, not speculative future problems
- Apply YAGNI (You Aren't Gonna Need It) principle to treatment design
- Prioritize treatment simplicity and maintainability
- Use iterative design with stakeholder feedback loops

### Challenge: Technology Selection Paralysis

**Solutions:**
- Define evaluation criteria based on treatment requirements
- Create proof-of-concept treatments for critical technology decisions
- Consider team expertise and treatment learning curve
- Document trade-offs and decision rationale

### Challenge: Treatment Scalability vs. Simplicity Trade-offs

**Solutions:**
- Design treatment for current scale plus reasonable growth factor
- Identify scaling bottlenecks and document future treatment options
- Use modular treatment design to enable incremental scaling
- Implement monitoring to validate treatment scaling assumptions

### Challenge: Integration Complexity in Treatment

**Solutions:**
- Define clear treatment interface contracts and boundaries
- Use standard protocols and data formats for treatment integration
- Design for loose coupling between treatment components
- Plan treatment integration testing strategy early

## Tips and Best Practices

### Treatment Architecture Design

- Start with simple treatment architecture and evolve as needed
- Use proven treatment patterns and avoid reinventing solutions
- Document treatment decisions and trade-offs
- Consider Conway's Law when designing team treatment boundaries

### Contract Design

- Define clear interface contracts with version management
- Follow REST principles or GraphQL best practices for APIs
- Use consistent naming conventions and error handling
- Design treatment contracts for backward compatibility

### Treatment Data Design

- Design treatment data models to reduce redundancy
- Consider read/write patterns and treatment query optimization
- Plan for data migration and treatment schema evolution
- Document treatment data relationships and constraints

### Security by Design

- Apply principle of least privilege in treatment access
- Design defense in depth with multiple treatment security layers
- Consider OWASP Top 10 vulnerabilities in treatment design
- Plan for secure treatment communication and data storage

## DDX Integration

### Using DDX Patterns

Apply relevant DDX patterns for treatment architecture and design:

```bash
ddx apply patterns/treatment/solution-architecture
ddx apply patterns/contracts/interface-specification
ddx apply patterns/treatment/complexity-assessment
ddx apply templates/treatment/solution-design
```

### Treatment Design Reviews

Use DDX diagnostics to validate treatment design quality:

```bash
ddx diagnose --phase prescribe
ddx diagnose --artifact treatment-architecture
ddx diagnose --artifact interface-contracts
ddx diagnose --artifact complexity-assessment
```

### Documentation Standards

Follow DDX documentation patterns for treatment specification:

```bash
ddx apply patterns/documentation/treatment-plan
ddx apply patterns/documentation/interface-contracts
```

## Quality Gates

### Treatment Architecture Review Checklist

- [ ] Treatment architecture aligns with diagnostic criteria
- [ ] Non-functional requirements addressed (performance, security, scalability)
- [ ] Technology choices justified and documented for treatment
- [ ] Treatment component interfaces clearly defined
- [ ] Treatment data flows and integration points documented
- [ ] Deployment and operational treatment considerations addressed

### Treatment Design Validation

- [ ] Treatment design addresses all symptoms and diagnostic criteria
- [ ] Interface contracts complete and consistent
- [ ] Treatment data schema optimized and normalized
- [ ] Security measures adequate for treatment requirements
- [ ] Treatment design reviewed by senior developers
- [ ] Implementation estimates realistic and detailed

## Risk Mitigation

### Treatment Technical Risks

- Create proof-of-concept for complex or unfamiliar treatment technologies
- Plan spike investigations for uncertain treatment design areas
- Design fallback treatment options for critical dependencies
- Validate treatment performance assumptions with prototypes

### Treatment Integration Risks

- Define clear treatment interface contracts between components
- Plan treatment integration testing strategy
- Identify external treatment dependencies and failure modes
- Design for graceful treatment degradation and error handling

## Next Phase

Upon successful completion of the Prescribe phase, proceed to **[[03-treat|Treatment Phase]]** where the detailed treatment prescriptions will be implemented into working treatment solutions.

---

*This document is part of the DDX Clinical Development Process (CDP) Workflow. For questions or improvements, use `ddx contribute` to share feedback with the community.*