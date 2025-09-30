# Solution Design Generation Prompt

Transform business requirements into a concrete technical approach that bridges Frame outputs to Design artifacts.

## Storage Location

Store the solution design at: `docs/helix/02-design/solution-design.md`

## Purpose

The Solution Design translates business language into technical language, documenting HOW to implement the system using decisions from ADRs and technologies from Tech Spikes.

## Key Principles

### 1. Requirements-First Thinking
Start with what the business needs, not what technology you want to use. Review every requirement from the specification and understand the "why" behind each.

### 2. Multiple Approaches
Always consider alternatives. Document 2-3 different solution approaches, evaluate trade-offs objectively, and explain why you selected or rejected each.

### 3. Leverage Existing Decisions
Build on established architectural decisions and technology selections:
- Reference supporting ADRs that justify the approach
- Use technologies selected through Tech Spikes
- Don't re-argue architectural decisions - focus on implementation

### 4. Clear Traceability
Every design decision must trace back to requirements. Map each requirement to specific components and show how NFRs influence architecture.

## Artifact Relationships

Solution Designs complete the design flow:

```
ADR (Why) → Tech Spike (What) → Solution Design (How)
```

**Focus on**:
- Component architecture and interactions
- Data flow and processing patterns
- Integration strategies between components
- Deployment and scaling architecture
- Error handling and recovery patterns
- Security implementation approach

**Cross-reference**: Reference the ADR establishing the architectural approach and the Tech Spike selecting technologies.

## Process

1. **Analyze Frame Outputs**: Read feature specs, user stories, PRD, and principles
2. **Identify Technical Implications**: Determine what capabilities, patterns, and approaches are needed
3. **Model the Domain**: Extract core entities, relationships, business rules, and boundaries
4. **Decompose into Components**: Identify natural boundaries, minimize dependencies, ensure single responsibility
5. **Design Integration**: Define how components interact, data flows, and error handling

## Critical Questions

- Does every functional requirement have a technical solution?
- Are all non-functional requirements addressed?
- Are user workflows supported end-to-end?
- Why this architecture over alternatives?
- How does it support future growth?
- What are the failure modes?

## Common Pitfalls

❌ **Technology-First Thinking**: "Let's use microservices because they're modern"
✅ **Requirements-First**: "Our requirement for independent scaling justifies microservices"

❌ **Over-Engineering**: Complex architecture for simple requirements
✅ **Right-Sized**: Simplest architecture that meets current needs

❌ **Under-Specifying**: "We'll figure out the details later"
✅ **Clear Decisions**: Documented rationale and clear direction

## Quality Checklist

- All requirements addressed
- Domain model captures business logic
- Component responsibilities clear
- Technology stack defined
- Decisions are justified
- Risks are identified
- Can be built with available resources
- Follows project principles

## Remember

The Solution Design bridges "what we need" and "how we'll build it". It should be comprehensive enough to guide implementation, clear enough for stakeholder approval, and traceable back to business needs.