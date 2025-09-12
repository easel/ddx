# Architecture Decision Records (ADRs)

Architecture Decision Records (ADRs) are documents that capture important architectural decisions made during development. They provide context, rationale, and consequences of architectural choices to help future maintainers understand why certain decisions were made.

## Why Use ADRs?

- **Context Preservation**: Capture the reasoning behind architectural decisions
- **Knowledge Transfer**: Help new team members understand the system's design
- **Decision Tracking**: Maintain a historical record of architectural evolution
- **Trade-off Documentation**: Record the alternatives considered and why they were rejected
- **Future Reference**: Support future architectural decisions with historical context

## When to Create an ADR

Create an ADR for decisions that:
- Have significant impact on the system's architecture
- Affect multiple components or subsystems
- Involve trade-offs between competing alternatives
- Establish architectural patterns or conventions
- Address non-functional requirements (performance, security, scalability)
- Impact development processes or tooling choices

## ADR Workflow

1. **Identify Decision**: Recognize when an architectural decision needs to be made
2. **Research Options**: Investigate alternative approaches and their implications
3. **Create ADR**: Use the template to document the decision
4. **Review**: Get input from relevant stakeholders and team members
5. **Decide**: Make the final decision and update the ADR status
6. **Implement**: Execute the decision and reference the ADR in implementation
7. **Learn**: Update the ADR with actual consequences as they emerge

## ADR Structure

Each ADR follows a consistent structure:

- **Title**: Brief description of the decision
- **Status**: Current state (Proposed, Accepted, Rejected, Superseded)
- **Context**: Background information and forces influencing the decision
- **Decision**: The architectural decision and its rationale
- **Consequences**: Expected and actual outcomes of the decision

## Files in This Directory

- **[template.md](template.md)**: Standard ADR template structure
- **[prompt.md](prompt.md)**: Guided prompts for creating comprehensive ADRs
- **examples/**: Sample ADRs demonstrating best practices

## Numbering Convention

ADRs are typically numbered sequentially (e.g., `001-use-database-migrations.md`, `002-adopt-microservices.md`). This helps maintain chronological order and makes references easier.

## ADR Lifecycle

### Status Definitions

- **Proposed**: Decision is under consideration
- **Accepted**: Decision has been approved and should be implemented
- **Rejected**: Decision was considered but not adopted
- **Superseded**: Decision was replaced by a later ADR (reference the new ADR)

### Updating ADRs

ADRs should be immutable once accepted. If circumstances change:
- Create a new ADR that supersedes the old one
- Update the old ADR's status to "Superseded by ADR-XXX"
- Link to the superseding ADR

## Best Practices

### Writing Effective ADRs

- **Be Clear**: Use plain language and avoid unnecessary jargon
- **Be Concise**: Focus on the essential information
- **Include Trade-offs**: Document alternatives and why they were rejected
- **Use Data**: Support decisions with metrics, benchmarks, or evidence
- **Update Consequences**: Record actual outcomes as they become known

### Common Pitfalls

- **Too Late**: Writing ADRs after implementation loses valuable context
- **Too Vague**: Decisions without clear rationale don't help future readers
- **Missing Alternatives**: Not documenting rejected options loses valuable information
- **Status Unclear**: Ambiguous status makes it unclear if decisions are active

## Integration with Development Workflow

ADRs integrate with other development artifacts:
- **PRDs**: Architecture decisions implement requirements from Product Requirements Documents
- **Feature Specs**: Detailed specifications reference relevant ADRs
- **Test Plans**: Architecture decisions influence testing strategies
- **Implementation**: Code should reference relevant ADRs in comments

## Tools and Templates

- Use the [template.md](template.md) for consistent structure
- Use the [prompt.md](prompt.md) for guided ADR creation
- Reference examples for inspiration and best practices
- Store ADRs in version control alongside code

## Getting Started

1. Review the [template structure](template.md)
2. Use the [guided prompts](prompt.md) to create your first ADR
3. Check the examples directory for reference implementations
4. Integrate ADR creation into your development workflow

Remember: The goal of ADRs is to capture the reasoning behind architectural decisions for future reference. Focus on why decisions were made, not just what was decided.