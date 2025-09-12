---
tags: [prompt, adr, architecture, development-workflow, decisions, ai-assisted]
references: template.md
created: 2025-01-12
modified: 2025-01-12
---

# Architecture Decision Record Creation Assistant

This prompt helps you create comprehensive Architecture Decision Records (ADRs) using the DDX ADR template. Work through each section systematically to document architectural decisions effectively.

## Using This Prompt

1. Identify the architectural decision that needs documentation
2. Answer the guiding questions for each section
3. Use the provided frameworks for consistent analysis
4. Ensure all stakeholders and alternatives are considered
5. Review the completed ADR with technical team

## Template Reference

This prompt uses the ADR template at [[template|template.md]]. You can also reference the comprehensive ADR guide at [[README|README.md]] and explore examples in the [[examples|examples/]] directory.

## Decision Identification

**Before creating an ADR, ensure this decision:**
- Has significant architectural impact
- Affects multiple system components
- Involves trade-offs between alternatives
- Will influence future development
- Needs to be understood by future team members

**Questions to Answer First:**
- What specific architectural challenge are we addressing?
- Why is this decision important enough to document?
- Who are the key stakeholders affected by this decision?
- What's the timeline for making and implementing this decision?

## Section-by-Section Guidance

### Title and Metadata

**ADR Number:**
- Use sequential numbering (001, 002, etc.)
- Include leading zeros for sorting
- Consider topic-based prefixes if helpful

**Title Format:**
- Keep it concise (under 10 words)
- Use active voice when possible
- Focus on the decision, not the problem
- Examples: "Use PostgreSQL for Primary Database", "Adopt Microservices Architecture"

**Status Guidelines:**
- **Proposed**: Still gathering input, not yet decided
- **Accepted**: Approved and ready for implementation
- **Rejected**: Considered but decided against
- **Superseded**: Replaced by a later ADR (reference it)

### Context Section

**Problem Statement - Address:**
- What specific architectural challenge exists?
- Why does this need to be solved now?
- What happens if we don't address this?
- What constraints or requirements drive this need?

**Forces Analysis Framework:**
Use this framework to identify competing forces:

**Technical Forces:**
- Performance requirements
- Scalability needs
- Security requirements
- Maintainability concerns
- Integration complexity
- Technology constraints

**Business Forces:**
- Time to market pressures
- Budget limitations
- Team expertise
- Risk tolerance
- Compliance requirements
- Strategic direction

**Organizational Forces:**
- Team size and structure
- Existing systems
- Operational capabilities
- Change management capacity
- Support requirements

**Questions to Explore:**
- What are the competing priorities or requirements?
- Which forces are most important to stakeholders?
- What trade-offs are we willing to make?
- What are the consequences of inaction?

### Decision Section

**Decision Documentation:**
- State the decision clearly and concisely
- Avoid implementation details (focus on "what" not "how")
- Use definitive language ("We will..." not "We should...")
- Make it specific enough to be actionable

**Rationale Deep-dive Questions:**
- Why is this the best option given our constraints?
- What key factors drove this decision?
- How does this align with our architectural principles?
- What evidence supports this choice?
- What assumptions are we making?

**Decision Quality Checklist:**
- [ ] Decision is clearly stated
- [ ] Rationale is well-reasoned
- [ ] Key stakeholders agree
- [ ] Implementation is feasible
- [ ] Risks are acceptable
- [ ] Aligns with business goals

### Alternatives Analysis

**For Each Alternative, Document:**

**Complete Description:**
- What exactly is this alternative?
- How would it be implemented?
- What would it look like in practice?

**Thorough Pros/Cons Analysis:**
- List all significant advantages
- Identify all major drawbacks
- Consider both short-term and long-term impacts
- Include technical, business, and operational factors

**Rejection Rationale:**
- What was the primary reason for rejection?
- Were there secondary factors?
- Under what circumstances might this be reconsidered?

**Alternative Evaluation Framework:**

**Technical Evaluation:**
- Performance implications
- Complexity and maintainability
- Integration requirements
- Scalability potential
- Security considerations

**Business Evaluation:**
- Cost implications
- Time to implement
- Risk assessment
- Strategic alignment
- Competitive advantage

**Operational Evaluation:**
- Support requirements
- Monitoring needs
- Deployment complexity
- Team training needed

### Consequences Analysis

**Positive Consequences - Consider:**
- Immediate benefits
- Long-term advantages
- Enabling capabilities
- Risk reductions
- Efficiency gains

**Negative Consequences - Consider:**
- Immediate costs or complexity
- Long-term maintenance burden
- New risks introduced
- Limitations created
- Dependencies added

**Neutral Consequences - Consider:**
- Changes that are neither clearly good nor bad
- Trade-offs that balance benefits and costs
- Shifts in complexity or responsibility
- Different approaches to familiar problems

**Consequence Assessment Questions:**
- How will this affect system performance?
- What new capabilities will this enable?
- What new risks or complexities are introduced?
- How will this impact the development team?
- What operational changes are required?
- How might this affect future architectural decisions?

### Implementation Planning

**Required Changes Analysis:**
- Which system components need modification?
- What new components need to be created?
- Are there process changes required?
- What infrastructure changes are needed?

**Migration Strategy:**
- Can this be implemented incrementally?
- What's the rollback plan if issues arise?
- How will we minimize disruption?
- What dependencies must be addressed first?

**Success Metrics:**
- How will we know this decision was successful?
- What should we measure and when?
- What are the target values or thresholds?
- Who will be responsible for monitoring?

### Risk Assessment

**Risk Identification Framework:**

**Technical Risks:**
- Integration challenges
- Performance issues
- Security vulnerabilities
- Scalability limitations
- Complexity management

**Business Risks:**
- Schedule delays
- Budget overruns
- Strategic misalignment
- Competitive disadvantage
- Compliance failures

**Operational Risks:**
- Support challenges
- Monitoring gaps
- Deployment issues
- Team capability gaps
- Process disruptions

**For Each Risk:**
- What's the probability of occurrence?
- What would be the impact if it occurs?
- How can we prevent or mitigate it?
- What's our contingency plan?

### Monitoring and Review

**Key Success Indicators:**
- Performance metrics to track
- Quality metrics to monitor
- Business metrics to measure
- User experience indicators

**Review Planning:**
- When should we formally review this decision?
- What conditions would trigger an early review?
- Who should be involved in the review?
- What criteria will we use to evaluate success?

## Template

```markdown
{{include: template.md}}
```

## Common ADR Patterns

### Technology Selection ADRs
Focus on: evaluation criteria, comparison matrix, proof-of-concept results

### Architecture Pattern ADRs  
Focus on: system boundaries, interaction patterns, data flow implications

### Process Change ADRs
Focus on: current state problems, process improvement goals, adoption strategy

### Infrastructure ADRs
Focus on: operational requirements, cost implications, scaling characteristics

## Quality Assurance

### Before Publishing, Verify:

**Completeness:**
- [ ] All required sections are filled out
- [ ] Decision is clearly stated
- [ ] Rationale is well-documented
- [ ] Alternatives were seriously considered
- [ ] Consequences are realistic

**Clarity:**
- [ ] Technical decisions are explained in business terms
- [ ] Assumptions are explicitly stated
- [ ] Trade-offs are clearly articulated
- [ ] Future readers will understand the context

**Accuracy:**
- [ ] Technical details are correct
- [ ] Timeline and constraints are accurate
- [ ] Stakeholder input is properly represented
- [ ] References and links are valid

### Common ADR Pitfalls to Avoid

1. **Post-hoc Documentation**: Writing ADRs after decisions are implemented loses valuable context
2. **Predetermined Outcomes**: Going through the motions when the decision is already made
3. **Insufficient Analysis**: Not thoroughly considering alternatives or consequences  
4. **Vague Decisions**: Statements that could be interpreted multiple ways
5. **Missing Context**: Not capturing the forces and constraints that influenced the decision
6. **Orphaned ADRs**: Creating ADRs that are never referenced or updated

## Review Process

### Stakeholder Review
- Technical leads should review for accuracy
- Architects should review for strategic alignment
- Operations should review for feasibility
- Security should review for compliance

### Review Questions
- Is the problem statement clear and compelling?
- Are all viable alternatives documented?
- Is the decision rationale convincing?
- Are the consequences realistic?
- Is the implementation plan feasible?

## Integration with Development Process

### Linking to Other Artifacts
- Reference relevant PRDs that drove the need
- Link to feature specifications that implement the decision
- Connect to test plans that validate the implementation
- Reference related ADRs for context

### Implementation Tracking
- Update ADR status as implementation progresses
- Document actual consequences as they emerge
- Record lessons learned for future decisions
- Update related documentation to reference the ADR

## Examples and References

- [[examples/database-migration|Database Migration ADR]]: Technology selection example
- [[examples/microservices-adoption|Microservices Adoption ADR]]: Architecture pattern example  
- [[../prd/README|PRD Documentation]]: How ADRs relate to product requirements
- [[../feature-spec/README|Feature Specifications]]: How ADRs inform detailed specs

## Need Help?

- For ADR best practices: [[README|ADR Overview]]
- For architecture principles: [[docs/architecture/principles|Architecture Principles]]
- For decision-making frameworks: [[prompts/common/technical_decisions|Technical Decision Guide]]
- For examples: Check the `examples/` directory

Remember: The goal is to capture the reasoning behind architectural decisions so future team members can understand not just what was decided, but why it was decided and what alternatives were considered.