# PRD Generation Prompt

Create a comprehensive Product Requirements Document that defines the business vision and requirements.

## Storage Location

Store the PRD at: `docs/helix/01-frame/prd.md`

## PRD Purpose

The PRD is a **business document** that:
- Communicates the vision to stakeholders
- Aligns team on goals and success metrics
- Defines the problem and opportunity
- Establishes scope and priorities

## Key Principles

### 1. Focus on WHY and WHAT, not HOW
✅ "Users need to track their expenses"
❌ "Users click a button that calls the API endpoint"

### 2. Be Specific About Success
- Define measurable metrics with targets
- Specify timeline for achieving goals
- Include method for measurement

### 3. Know Your Users
- Create detailed personas based on research
- Understand their goals and pain points
- Design for specific users, not "everyone"

### 4. Prioritize Ruthlessly
- **P0** = Cannot ship without
- **P1** = Should have for good experience
- **P2** = Nice to have if time permits

## Section Guidance

**Executive Summary**: Write LAST. Summarize the problem, solution, key metrics, timeline, and scope.

**Problem Statement**: Start with user pain, quantify if possible, explain why now.

**Goals and Success Metrics**: Link metrics to user value, make them specific and measurable, include baseline and target.

**Users and Personas**: Base on actual research/data, include specific scenarios, focus on primary persona for MVP.

**Requirements Overview**: Start with user needs not features, group by priority, keep high-level.

**Risks and Mitigation**: Be honest about uncertainties, include technical and business risks, provide specific mitigation strategies.

## Common Pitfalls

❌ **Solution Masquerading as Problem**: "Users need a dashboard"
✅ **True Problem**: "Users can't track project progress"

❌ **Vague Success Metrics**: "Improve user satisfaction"
✅ **Specific Metrics**: "Increase NPS from 30 to 50 within 6 months"

❌ **Feature Laundry List**: Long list of features without context
✅ **Prioritized Needs**: User needs with priority and rationale

## Quality Checklist

- Would a new team member understand the vision?
- Are success metrics specific and measurable?
- Is the primary persona clearly defined?
- Are requirements prioritized (P0, P1, P2)?
- Are risks identified with mitigation plans?
- Is scope clearly defined (including non-goals)?

---

Remember: The PRD is about alignment and clarity. It should answer: Why are we building this? Who for? What does success look like? What are we building (and not building)? What could go wrong?