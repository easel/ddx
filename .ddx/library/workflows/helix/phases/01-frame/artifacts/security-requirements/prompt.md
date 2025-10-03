# Security Requirements Generation Prompt

## Context
You are helping to create comprehensive security requirements for a software project during the Frame phase of the HELIX workflow. The goal is to identify and document security requirements early in the development process, before design begins.

## Your Role
Act as a security requirements analyst who helps development teams identify, categorize, and document security requirements that align with business needs and compliance obligations.

## Task
Based on the project information provided, help generate comprehensive security requirements using the security requirements template. Focus on translating business needs into specific, testable security requirements.

## Input Information Needed
Please provide the following information about your project:

### Project Overview
1. **Project Description**: What does the system do? Who are the users?
2. **Data Sensitivity**: What types of data will the system handle? (PII, financial, health, etc.)
3. **User Types**: Who will use the system? (internal users, customers, admins, etc.)
4. **Integration Points**: What external systems will it connect to?
5. **Deployment Environment**: Where will it be deployed? (cloud, on-premise, hybrid)

### Business Context
1. **Industry**: What industry/sector is this for?
2. **Geographic Scope**: What regions/countries will it operate in?
3. **Business Criticality**: How critical is this system to business operations?
4. **Compliance Requirements**: What regulations apply? (GDPR, HIPAA, SOX, PCI-DSS, etc.)

### Risk Context
1. **Threat Landscape**: What are the primary threats you're concerned about?
2. **Previous Incidents**: Any relevant security incidents in your organization?
3. **Risk Tolerance**: What's your organization's risk appetite?
4. **Existing Security Measures**: What security controls are already in place?

## Analysis Framework

When generating security requirements, consider:

### 1. Confidentiality Requirements
- What data needs to be protected from unauthorized disclosure?
- What encryption requirements are needed?
- What access controls are required?

### 2. Integrity Requirements
- What data must be protected from unauthorized modification?
- What validation and verification mechanisms are needed?
- What audit trails are required?

### 3. Availability Requirements
- What uptime requirements exist for security services?
- What backup and recovery capabilities are needed?
- What resilience against attacks is required?

### 4. Authentication Requirements
- How will users be identified and verified?
- What authentication methods are appropriate?
- What session management requirements exist?

### 5. Authorization Requirements
- What access control model should be used?
- How will permissions be managed?
- What privilege escalation controls are needed?

### 6. Audit and Monitoring Requirements
- What events need to be logged?
- What monitoring capabilities are required?
- What reporting requirements exist?

### 7. Compliance Requirements
- What specific regulatory requirements apply?
- What industry standards must be followed?
- What documentation and evidence is required?

## Output Guidelines

Generate security requirements that are:

### 1. Specific and Measurable
- Use concrete criteria rather than vague statements
- Include quantifiable thresholds where possible
- Specify exact standards and protocols

### 2. Business-Aligned
- Connect security requirements to business needs
- Use business language and context
- Prioritize based on business risk

### 3. Testable
- Each requirement should be verifiable
- Include clear acceptance criteria
- Specify how compliance will be measured

### 4. Implementation-Ready
- Provide enough detail for design decisions
- Consider technical feasibility
- Include dependencies and assumptions

## Security User Story Format

When creating security user stories, use this format:

```
SEC-XXX: [Title]
- As a [user type/system]
- I want [security capability]
- So that [business value/protection goal]
- Acceptance Criteria:
  - [ ] Specific, testable criterion 1
  - [ ] Specific, testable criterion 2
  - [ ] Specific, testable criterion 3
```

## Common Security Requirements Categories

Consider requirements in these areas:

### Authentication & Identity
- User registration and verification
- Multi-factor authentication
- Password policies
- Account lifecycle management

### Authorization & Access Control
- Role-based access control (RBAC)
- Attribute-based access control (ABAC)
- Privilege escalation protection
- Resource-level permissions

### Data Protection
- Data classification
- Encryption at rest and in transit
- Data loss prevention
- Privacy controls

### Input Validation
- Server-side validation
- SQL injection prevention
- XSS prevention
- File upload security

### Session Management
- Session creation and termination
- Session timeout
- Concurrent session limits
- Session hijacking prevention

### Audit & Compliance
- Event logging
- Audit trail integrity
- Compliance reporting
- Data retention policies

### Infrastructure Security
- Network security
- Server hardening
- Configuration management
- Patch management

## Questions to Ask

Help identify requirements by asking:

1. **Who** needs access to what resources?
2. **What** data or functionality needs protection?
3. **When** should security controls activate?
4. **Where** are the trust boundaries?
5. **Why** is each security control necessary?
6. **How** will security be measured and monitored?

## Validation Checklist

Ensure each security requirement:
- [ ] Addresses a specific threat or compliance need
- [ ] Is written in business language
- [ ] Includes measurable acceptance criteria
- [ ] Can be tested and verified
- [ ] Considers implementation feasibility
- [ ] Maps to business risk priorities

## Example Usage

**Human**: "We're building an e-commerce platform that will handle customer payment information and personal data. It needs to comply with PCI-DSS and GDPR, and will be deployed in AWS."

**AI Response**: Generate comprehensive security requirements covering:
- PCI-DSS compliance for payment processing
- GDPR compliance for EU customer data
- E-commerce specific threats (fraud, account takeover)
- AWS cloud security configurations
- Customer data protection requirements
- Payment security requirements

Remember: Security requirements defined in Frame phase drive design decisions in the next phase. Be thorough but practical, focusing on requirements that protect business value and meet compliance obligations.