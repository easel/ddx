# Architecture Decision Record (ADR) Generation Prompt

Document significant architectural decisions with clear rationale and trade-offs.

## Storage Location

Store ADRs at: `docs/02-design/adr/ADR-{number}-{title-with-hyphens}.md`

## Naming Convention

Follow this consistent format:
- **File Format**: `ADR-{number}-{title-with-hyphens}.md`
- **Number**: Zero-padded 3-digit sequence (001, 002, 003...)
- **Title**: Descriptive, lowercase with hyphens
- **Examples**:
  - `ADR-001-use-postgresql-for-primary-database.md`
  - `ADR-002-adopt-microservices-architecture.md`
  - `ADR-003-implement-jwt-authentication.md`

## What Deserves an ADR?

Create an ADR when the decision:
- Has significant impact on the system architecture
- Involves choosing between multiple viable alternatives
- Has long-term implications for maintenance or evolution
- Affects multiple components or teams
- Involves significant trade-offs
- Deviates from standard practices
- Has high cost of change if wrong

**Focus on fundamental architectural choices**, not specific technology selections.

### ADR Scope Guidelines

**✅ Good ADR Topics:**
- Protocol choices: "Use GraphQL for internal APIs" (not "Use Caliban for GraphQL")
- Architectural patterns: "Adopt microservices architecture" (not "Use Spring Boot")
- Data strategies: "Use event-driven architecture" (not "Use Kafka implementation")
- System boundaries: "Separate internal and external APIs" (not "Use different ports")

**❌ Not for ADRs (use other artifacts):**
- Specific library/framework selection → **Tech Spike**
- Implementation patterns and details → **Solution Design**
- Configuration and deployment procedures → **Implementation Guide**
- Tool selection for development workflow → **Process documentation**

## What Doesn't Need an ADR?

Skip ADRs for:
- Obvious choices with no alternatives
- Minor implementation details
- Decisions that can easily be changed
- Standard/conventional approaches
- **Specific technology selections** (create a **Tech Spike** instead)
- **Implementation architecture** (create a **Solution Design** instead)

## Artifact Relationships

ADRs work with other design artifacts in this flow:

```
ADR (Why) → Tech Spike (What) → Solution Design (How)
```

**Example:**
1. **ADR-012**: "Use distributed caching" → Why we need caching
2. **SPIKE-003**: "Redis vs Hazelcast evaluation" → Which technology
3. **SD-008**: "Redis cluster implementation" → How to implement

**Cross-reference related artifacts** in your ADR to maintain traceability.

## Key Principles

### 1. Document the "Why", Not Just the "What"
The most valuable part of an ADR is understanding why a decision was made:
- What problem were we solving?
- What constraints influenced us?
- What alternatives did we consider?
- Why did we reject the alternatives?

### 2. Be Honest About Trade-offs
Every architectural decision involves trade-offs:
- Document both benefits and drawbacks
- Explain how you'll mitigate negative consequences
- Acknowledge technical debt being incurred

### 3. Keep It Relevant
Focus on information that will help future readers:
- Will someone in 2 years understand why?
- Does it help prevent repeating mistakes?
- Does it explain non-obvious choices?

## Writing Process

### Step 1: Establish Context
Before documenting the decision:
- What triggered this decision?
- What requirements are driving it?
- What constraints exist?
- What's the current state?

### Step 2: Explore Alternatives
For each viable option:
- Describe it clearly
- List genuine pros and cons
- Evaluate against requirements
- Consider long-term implications

### Step 3: Make and Justify the Decision
- State the decision clearly
- Explain why this option best meets needs
- Acknowledge the trade-offs accepted
- Describe how negatives will be addressed

### Step 4: Document Consequences
Think through all impacts:
- Development effort required
- Operational implications
- Performance impacts
- Security considerations
- Maintenance burden

### Step 5: Define Success Criteria
How will we validate this decision:
- Metrics to track
- Success indicators
- Review triggers

## Quality Checklist

### Context Section
- [ ] Problem is clearly stated
- [ ] Requirements driving decision are listed
- [ ] Constraints are documented
- [ ] Current state is described (if applicable)

### Decision Section
- [ ] Decision is stated in active voice
- [ ] Decision is specific and actionable
- [ ] Key points are highlighted
- [ ] Rationale is clear

### Alternatives Section
- [ ] At least 2-3 alternatives considered
- [ ] Each alternative has genuine pros/cons
- [ ] Evaluation criteria are consistent
- [ ] Rejection reasons are clear

### Consequences Section
- [ ] Both positive and negative consequences listed
- [ ] Impact on different aspects considered
- [ ] Mitigation strategies for negatives included
- [ ] Long-term implications addressed

### Overall Quality
- [ ] Would a new team member understand this?
- [ ] Is the confidence level appropriate?
- [ ] Are risks identified and addressed?
- [ ] Can this guide similar future decisions?

## Common ADR Topics

### Technology Selection
- Programming language choice
- Framework selection
- Database technology
- Message queue selection
- Cloud provider choice

### Architectural Patterns
- Monolith vs microservices
- Synchronous vs asynchronous
- Event-driven vs request-response
- Layered vs hexagonal architecture

### Data Decisions
- Database schema approach
- Caching strategy
- Data partitioning
- Consistency vs availability trade-offs

### Security Decisions
- Authentication method
- Authorization model
- Encryption approach
- Secret management

### Operational Decisions
- Deployment strategy
- Scaling approach
- Monitoring and observability
- Disaster recovery

## Examples of Good ADR Titles

- "Use PostgreSQL for primary data storage"
- "Adopt event sourcing for audit requirements"
- "Choose monolithic architecture for MVP"
- "Implement JWT-based authentication"
- "Use Kubernetes for container orchestration"

## Red Flags to Avoid

### ❌ Vague Decisions
**Bad**: "We will use modern architecture"
**Good**: "We will use microservices with REST APIs"

### ❌ Missing Alternatives
**Bad**: Only documenting the chosen option
**Good**: Showing what else was considered and why rejected

### ❌ No Trade-offs
**Bad**: Only listing benefits
**Good**: Honest about drawbacks and mitigation

### ❌ No Success Criteria
**Bad**: No way to validate the decision
**Good**: Clear metrics and review triggers

## Remember

ADRs are for your future self and team members who will wonder:
- "Why did we build it this way?"
- "What else did we consider?"
- "Can we change this now?"
- "What were we thinking?"

Good ADRs prevent repeating past mistakes and provide confidence in technical decisions. They're an investment in your team's future productivity and understanding.