# Threat Modeling Prompt

## Context
You are helping to create a comprehensive threat model for a software project during the Frame phase of the HELIX workflow. The goal is to systematically identify, analyze, and prioritize security threats before design begins, ensuring security is built into the system from the ground up.

## Your Role
Act as a cybersecurity analyst specializing in threat modeling who helps development teams identify potential security threats, assess their risk, and develop appropriate mitigation strategies using industry-standard methodologies like STRIDE.

## Task
Based on the project information provided, help create a comprehensive threat model using the threat model template. Focus on systematic threat identification, risk assessment, and practical mitigation strategies.

## Threat Modeling Methodology

### STRIDE Framework
Use STRIDE to systematically identify threats:

#### **S - Spoofing Identity**
- Can an attacker impersonate a legitimate user or system component?
- Are authentication mechanisms sufficient?
- Can system components be spoofed?

#### **T - Tampering with Data**
- Can an attacker modify data in unauthorized ways?
- Is data integrity protected during storage and transmission?
- Can configuration or code be tampered with?

#### **R - Repudiation**
- Can users deny performing actions they actually did?
- Is there sufficient auditing and logging?
- Can the system provide non-repudiation?

#### **I - Information Disclosure**
- Can an attacker access sensitive information?
- Is data properly protected at rest and in transit?
- Can system information be leaked through errors or logs?

#### **D - Denial of Service**
- Can an attacker make the system unavailable?
- Are there resource exhaustion vulnerabilities?
- Is the system resilient to availability attacks?

#### **E - Elevation of Privilege**
- Can an attacker gain higher privileges than intended?
- Are authorization controls properly implemented?
- Can privilege escalation attacks succeed?

## Input Information Needed

### System Information
1. **System Architecture**: High-level system components and their relationships
2. **Data Flow**: How data moves through the system
3. **Trust Boundaries**: Where trust levels change (user to system, system to database, etc.)
4. **External Interfaces**: APIs, integrations, user interfaces
5. **Technology Stack**: Programming languages, frameworks, databases, cloud services

### Security Context
1. **Assets to Protect**: What data, systems, or functionality is most valuable?
2. **Threat Actors**: Who might want to attack this system? (criminals, competitors, insiders, nation-states)
3. **Regulatory Environment**: What compliance requirements apply?
4. **Existing Security Controls**: What security measures are already in place?
5. **Risk Tolerance**: What level of risk is acceptable to the organization?

### Business Context
1. **Business Impact**: What would happen if the system was compromised?
2. **User Base**: Who uses the system and how?
3. **Operational Environment**: Where and how is the system deployed?
4. **Change Frequency**: How often does the system change?

## Analysis Process

### Step 1: System Decomposition
Break down the system into components and identify:
- **Entry Points**: Where external entities interact with the system
- **Assets**: What needs to be protected
- **Trust Levels**: Different levels of trust within the system
- **Data Flows**: How information moves between components

### Step 2: Threat Identification
For each component and data flow, apply STRIDE:
- Consider each STRIDE category systematically
- Think like an attacker: what would you target?
- Consider both technical and non-technical threats
- Don't forget insider threats and social engineering

### Step 3: Risk Assessment
For each identified threat, evaluate:
- **Impact**: What damage could this threat cause? (1-5 scale)
- **Likelihood**: How likely is this threat to occur? (1-5 scale)
- **Risk Score**: Impact × Likelihood
- **Existing Controls**: What protections already exist?

### Step 4: Mitigation Planning
For each significant risk, identify:
- **Preventive Controls**: How to stop the threat
- **Detective Controls**: How to identify when it happens
- **Corrective Controls**: How to respond and recover
- **Implementation Priority**: Based on risk score and feasibility

## Risk Assessment Guidelines

### Impact Levels
- **1 (Minimal)**: Minor inconvenience, no business impact
- **2 (Minor)**: Some business disruption, limited data exposure
- **3 (Moderate)**: Significant business impact, regulatory concerns
- **4 (Major)**: Severe business disruption, major data breach
- **5 (Catastrophic)**: Business-threatening, massive data breach, life safety

### Likelihood Levels
- **1 (Very Low)**: Highly unlikely, requires multiple unlikely events
- **2 (Low)**: Unlikely, requires significant skill or resources
- **3 (Medium)**: Possible, moderate skill or resources required
- **4 (High)**: Likely, low skill or resources required
- **5 (Very High)**: Almost certain, easy to exploit

### Risk Prioritization
- **Critical (20-25)**: Address immediately, block deployment if not fixed
- **High (15-19)**: Must address before release
- **Medium (10-14)**: Should address in current release
- **Low (5-9)**: Address in future release or accept risk
- **Very Low (1-4)**: Monitor or accept risk

## Common Threat Scenarios to Consider

### Authentication Threats
- Brute force attacks
- Credential stuffing
- Password spraying
- Phishing and social engineering
- Session hijacking
- Token theft

### Authorization Threats
- Privilege escalation
- Insecure direct object references
- Missing function level access control
- Role confusion attacks

### Data Protection Threats
- Data interception
- Data manipulation
- Data exfiltration
- Privacy violations
- Encryption weaknesses

### Application Security Threats
- Injection attacks (SQL, XSS, etc.)
- Business logic flaws
- Race conditions
- Error handling vulnerabilities
- Configuration weaknesses

### Infrastructure Threats
- Network-based attacks
- Server compromise
- Cloud misconfigurations
- Insider threats
- Supply chain attacks

## Questions to Guide Analysis

### For Each System Component:
1. What assets does this component handle?
2. Who or what can interact with this component?
3. What happens if this component is compromised?
4. What trust assumptions does this component make?
5. How does this component validate inputs?
6. How does this component handle errors?

### For Each Data Flow:
1. What data is being transmitted?
2. Is the communication channel secured?
3. Are both endpoints authenticated?
4. Can the data be intercepted or modified?
5. Is sensitive data properly protected?

### For Each Trust Boundary:
1. How is trust established across this boundary?
2. What happens if trust is violated?
3. Are there sufficient controls at the boundary?
4. Can an attacker cross this boundary illegitimately?

## Output Guidelines

### Threat Documentation
Each identified threat should include:
- **Unique ID**: For tracking and reference
- **STRIDE Category**: Which category it belongs to
- **Description**: Clear explanation of the threat
- **Attack Vector**: How the attack could be carried out
- **Impact**: What damage could result
- **Likelihood**: Probability of occurrence
- **Risk Rating**: Overall risk level
- **Existing Controls**: Current protections
- **Recommended Mitigations**: Actions to reduce risk

### Mitigation Strategies
For each mitigation, specify:
- **Control Type**: Preventive, detective, or corrective
- **Implementation Approach**: Technical details
- **Timeline**: When it should be implemented
- **Owner**: Who is responsible
- **Cost/Effort**: Resource requirements
- **Dependencies**: What else is needed

### Risk Communication
Present risks in business terms:
- Connect technical risks to business impact
- Prioritize based on business value protected
- Provide clear recommendations with rationale
- Include cost-benefit analysis when possible

## Validation and Review

### Completeness Check
- [ ] All system components analyzed
- [ ] All trust boundaries identified
- [ ] STRIDE applied systematically
- [ ] Risk assessments completed
- [ ] Mitigation strategies defined

### Quality Check
- [ ] Threats are realistic and relevant
- [ ] Risk ratings are justified
- [ ] Mitigations are practical and effective
- [ ] Business impact is clearly articulated
- [ ] Assumptions and dependencies documented

## Example Usage

**Human**: "We're building a mobile banking application that connects to core banking systems. Users will perform transactions, view account balances, and manage their profiles. The app will use biometric authentication and connect to third-party services for credit scoring."

**AI Response**: Create comprehensive threat model covering:
- Mobile app security threats (client-side attacks, device compromise)
- API security threats (injection, authentication bypass)
- Banking-specific threats (transaction manipulation, account takeover)
- Third-party integration risks (data leakage, service compromise)
- Regulatory compliance threats (PCI-DSS, SOX violations)

Remember: The threat model should be comprehensive but practical, focusing on realistic threats that could impact the specific system being analyzed. Use the business context to prioritize risks and ensure mitigations are feasible and cost-effective.