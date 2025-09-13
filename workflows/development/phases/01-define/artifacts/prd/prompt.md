---
tags: [prompt, prd, development-workflow, requirements, ai-assisted]
references: template.md
created: 2025-01-12
modified: 2025-01-12
---

# Product Requirements Document Creation Assistant

This prompt helps you create a comprehensive Product Requirements Document using the DDX PRD template. Work through each section systematically to ensure complete requirements coverage.

## Using This Prompt

1. Review the template structure below
2. Answer the guiding questions for each section
3. Use the provided examples as reference
4. Ensure all required sections are completed

## Template Reference

This prompt uses the PRD template at [[template|template.md]]. You can also reference the comprehensive PRD guide at [[prompts/common/product_requirements|Product Requirements Guide]].

## Section-by-Section Guidance

### Executive Summary

**Questions to Answer:**
- What is the product's primary purpose in one sentence?
- Who are the main users and what problem does it solve for them?
- What makes this solution unique or valuable?
- What is the scope of this initial version?

**Tips:**
- Keep it concise (2-3 paragraphs maximum)
- Focus on the "why" more than the "what"
- Write this section last, after completing other sections

### Problem Statement

**The Problem - Consider:**
- What specific pain point are you addressing?
- How widespread is this problem?
- What is the cost/impact of not solving it?
- Can you provide concrete examples?

**Current State - Consider:**
- How do users currently handle this problem?
- What tools or workarounds do they use?
- What are the limitations of current solutions?
- What frustrations do users experience?

**Desired State - Consider:**
- What would an ideal solution look like?
- How would users' lives be different?
- What new capabilities would they have?
- What would success look like?

### Goals and Objectives

**Primary Goals - Should be:**
- Specific and measurable
- Achievable within project constraints
- Aligned with business objectives
- Limited to 3-5 key goals

**Success Metrics - For each metric define:**
- What exactly will be measured?
- What is the target value?
- How will it be measured?
- When will it be measured?
- Who is responsible for tracking?

### Users and Personas

**For Each Persona, Define:**
- **Job Title/Role**: What do they do?
- **Experience Level**: Novice, intermediate, or expert?
- **Context of Use**: When/where will they use the product?
- **Technical Proficiency**: Comfort with technology?
- **Frequency of Use**: Daily, weekly, occasional?

**Key Questions:**
- What are their primary goals?
- What frustrates them most?
- What would make their job easier?
- How do they measure success?

### User Stories

**Format:** As a [persona], I want to [action] so that [benefit]

**Prioritization:**
- **Epic Level**: Major functionality areas
- **Story Level**: Specific user actions
- **Acceptance Criteria**: Testable conditions

**Quality Checklist:**
- [ ] Written from user perspective
- [ ] Includes clear benefit/value
- [ ] Testable and specific
- [ ] Independently valuable

### Functional Requirements

**For Each Requirement:**
- **Name**: Clear, descriptive title
- **Description**: What it does (not how)
- **User Benefit**: Why users need it
- **Priority**: P0 (must), P1 (should), P2 (nice)
- **Acceptance Criteria**: Specific, testable conditions

**Prioritization Guidelines:**
- **P0**: Product fails without this
- **P1**: Significantly degraded experience without this
- **P2**: Enhancement that adds value

### Non-Functional Requirements

**Performance:**
- Response time requirements
- Throughput requirements
- Scalability needs
- Data volume expectations

**Security:**
- Authentication requirements
- Authorization model
- Data protection needs
- Compliance requirements

**Usability:**
- Accessibility standards
- Browser/device support
- Localization needs
- User training requirements

**Reliability:**
- Uptime requirements
- Error handling approach
- Data integrity needs
- Backup/recovery requirements

### Constraints and Assumptions

**Technical Constraints:**
- Technology limitations
- Integration requirements
- Platform restrictions
- Performance boundaries

**Business Constraints:**
- Budget limitations
- Timeline requirements
- Resource availability
- Legal/regulatory requirements

**Document All Assumptions:**
- User behavior assumptions
- Technical assumptions
- Market assumptions
- Resource assumptions

### Dependencies

**For Each Dependency:**
- What is the dependency?
- Why is it needed?
- When is it needed?
- What's the impact if unavailable?
- What's the mitigation plan?

### Risks and Mitigation

**Risk Assessment Framework:**
- **Identify**: What could go wrong?
- **Assess Impact**: High/Medium/Low
- **Assess Probability**: High/Medium/Low
- **Mitigation**: How to prevent or address
- **Contingency**: Backup plan if it occurs

### Timeline and Milestones

**Key Considerations:**
- Dependencies between milestones
- Critical path items
- Buffer for unknowns
- Review/approval cycles
- Testing and iteration time

### Out of Scope

**Be Explicit About:**
- Features for future versions
- Related but separate projects
- Assumed external responsibilities
- Deliberate limitations

## Template

```markdown
{{include: template.md}}
```

## Examples and References

- [[examples/ddx-v1|DDX PRD Example]]: See how DDX used this template
- [[docs/product/prd-ddx-v1|Full DDX PRD]]: Complete PRD for reference

## Validation Checklist

Before finalizing your PRD:

- [ ] All sections completed with sufficient detail
- [ ] Success metrics are specific and measurable
- [ ] Requirements are testable
- [ ] Priorities are clear (P0/P1/P2)
- [ ] Dependencies are identified
- [ ] Risks have mitigation strategies
- [ ] Timeline is realistic
- [ ] Stakeholders are identified
- [ ] Out of scope is explicit

## Common Pitfalls to Avoid

1. **Vague Requirements**: Use specific, measurable criteria
2. **Solution Bias**: Focus on the problem, not predetermined solutions
3. **Missing NFRs**: Don't forget performance, security, usability
4. **Unrealistic Timeline**: Account for testing, iteration, and delays
5. **Assumed Knowledge**: Define all terms and concepts
6. **Kitchen Sink**: Resist including everything in v1

## Next Steps

After completing your PRD:

1. Review with stakeholders
2. Get formal approval
3. Proceed to [[../architecture/README|Architecture Design]]
4. Use PRD to guide implementation
5. Validate against PRD during testing

## Need Help?

- For PRD best practices: [[docs/usage/workflows/overview|Workflow Overview]]
- For general requirements gathering: [[prompts/common/product_requirements|Product Requirements Guide]]
- For examples: Check the `examples/` directory