# Solution Design Generation Prompt

Transform business requirements into a concrete technical approach that bridges Frame outputs to Design artifacts.

## Storage Location

Store the solution design at: `docs/02-design/solution-design.md`

## Purpose

The Solution Design is the critical transformation layer that:
- Translates business language into technical language
- Maps requirements to architectural decisions
- Documents the rationale for technical choices
- Creates traceability from needs to implementation

## Key Principles

### 1. Requirements-First Thinking
Start with what the business needs, not what technology you want to use:
- Review every requirement from the specification
- Understand the "why" behind each requirement
- Consider business constraints before technical preferences

### 2. Multiple Approaches
Always consider alternatives:
- Document at least 2-3 different solution approaches
- Evaluate trade-offs objectively
- Explain why you selected or rejected each approach

### 3. Clear Traceability
Every design decision must trace back to requirements:
- Map each requirement to specific components
- Show how NFRs influence architecture
- Identify any requirements not fully addressed

## Process to Follow

### Step 1: Analyze Frame Outputs
Read and understand:
- Feature Specification (all requirements)
- User Stories (user needs and workflows)
- PRD (business context and constraints)
- Principles (technical constraints)

### Step 2: Identify Technical Implications
For each requirement, determine:
- What technical capability is needed?
- What patterns or approaches could work?
- What are the performance/scale implications?
- What are the security considerations?

### Step 3: Model the Domain
Extract from user stories and requirements:
- Core business entities
- Relationships between entities
- Business rules and invariants
- Transaction boundaries

### Step 4: Decompose into Components
Group related functionality:
- Identify natural boundaries
- Minimize dependencies between components
- Ensure single responsibility
- Map back to requirements

### Step 5: Evaluate Technology Options
For each technology choice:
- How does it meet the requirements?
- What are the trade-offs?
- What's the team's expertise?
- What's the long-term maintenance impact?

## Critical Questions to Answer

### Requirements Coverage
- [ ] Does every functional requirement have a technical solution?
- [ ] Are all non-functional requirements addressed?
- [ ] Are user workflows supported end-to-end?
- [ ] Are edge cases handled?

### Architecture Decisions
- [ ] Why this architecture over alternatives?
- [ ] How does it support future growth?
- [ ] What are the failure modes?
- [ ] How complex is it to understand and maintain?

### Technology Selection
- [ ] Why these specific technologies?
- [ ] What alternatives were considered?
- [ ] What are the risks of these choices?
- [ ] Do we have the skills to implement?

## Common Pitfalls to Avoid

### ❌ Technology-First Thinking
**Bad**: "Let's use microservices because they're modern"
**Good**: "Our requirement for independent scaling justifies microservices"

### ❌ Over-Engineering
**Bad**: Complex architecture for simple requirements
**Good**: Simplest architecture that meets current needs

### ❌ Under-Specifying
**Bad**: "We'll figure out the details later"
**Good**: Clear decisions with documented rationale

### ❌ Ignoring Constraints
**Bad**: Design that violates project principles
**Good**: Design that respects all constraints or documents exceptions

## Quality Checklist

Before finalizing the solution design:

### Completeness
- [ ] All requirements addressed
- [ ] Domain model captures business logic
- [ ] Component responsibilities clear
- [ ] Technology stack defined

### Clarity
- [ ] Non-technical stakeholders can understand approach
- [ ] Technical team knows what to build
- [ ] Decisions are justified
- [ ] Risks are identified

### Feasibility
- [ ] Can be built with available resources
- [ ] Timeline is realistic
- [ ] Skills exist or can be acquired
- [ ] Risks have mitigation strategies

### Alignment
- [ ] Follows project principles
- [ ] Supports business goals
- [ ] Enables required user workflows
- [ ] Allows for future evolution

## Output Expectations

The solution design should provide:
1. **Clear mapping** from requirements to technical approach
2. **Justified decisions** with alternatives considered
3. **Domain model** that captures business concepts
4. **Component architecture** with clear boundaries
5. **Technology rationale** based on requirements
6. **Risk assessment** with mitigation strategies

## Remember

The Solution Design is the bridge between "what we need" and "how we'll build it". It should be:
- **Comprehensive** enough to guide implementation
- **Clear** enough for stakeholder approval
- **Flexible** enough to accommodate learning
- **Traceable** back to business needs

This artifact prevents the common problem of jumping straight from requirements to code without thinking through the technical approach and its implications.