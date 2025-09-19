# PRD Generation Prompt

Create a comprehensive Product Requirements Document that defines the business vision and requirements.

## Storage Location

Store the PRD at: `docs/helix/01-frame/prd.md`

This central location ensures the PRD is easily discoverable and remains the single source of truth for product requirements.

## PRD Purpose

The PRD is a **business document** that:
- Communicates the vision to stakeholders
- Aligns team on goals and success metrics
- Defines the problem and opportunity
- Establishes scope and priorities

## Key Principles

### 1. Focus on the WHY and WHAT, not HOW
- ✅ "Users need to track their expenses"
- ❌ "Users click a button that calls the API endpoint"

### 2. Be Specific About Success
- Define measurable metrics with targets
- Specify timeline for achieving goals
- Include method for measurement

### 3. Know Your Users
- Create detailed personas based on research
- Understand their goals and pain points
- Design for specific users, not "everyone"

### 4. Prioritize Ruthlessly
- P0 = Cannot ship without
- P1 = Should have for good experience
- P2 = Nice to have if time permits

## Section-by-Section Guidance

### Executive Summary
Write this LAST. Summarize:
- The problem and why it matters
- The solution and its impact
- Key success metrics
- Timeline and scope

### Problem Statement
- Start with user pain, not solution
- Quantify the problem if possible
- Explain why solving this matters now

### Goals and Success Metrics
- Link metrics directly to user value
- Make metrics specific and measurable
- Include baseline and target

### Users and Personas
- Base on actual user research/data
- Include specific scenarios
- Focus on primary persona for MVP

### Requirements Overview
- Start with user needs, not features
- Group by priority (P0, P1, P2)
- Keep requirements high-level

### Risks and Mitigation
- Be honest about uncertainties
- Include technical and business risks
- Provide specific mitigation strategies

## Common Pitfalls to Avoid

### ❌ Solution Masquerading as Problem
**Bad**: "Users need a dashboard"
**Good**: "Users can't track project progress"

### ❌ Vague Success Metrics
**Bad**: "Improve user satisfaction"
**Good**: "Increase NPS from 30 to 50 within 6 months"

### ❌ Feature Laundry List
**Bad**: Long list of features without context
**Good**: User needs with priority and rationale

### ❌ Ignoring Constraints
**Bad**: Assuming unlimited resources
**Good**: Acknowledging technical, business, and user constraints

## Quality Checklist

Before completing the PRD:
- [ ] Would a new team member understand the vision?
- [ ] Are success metrics specific and measurable?
- [ ] Is the primary persona clearly defined?
- [ ] Are requirements prioritized (P0, P1, P2)?
- [ ] Are risks identified with mitigation plans?
- [ ] Is scope clearly defined (including non-goals)?

## Remember

The PRD is about **alignment and clarity**. It should answer:
1. Why are we building this?
2. Who are we building it for?
3. What does success look like?
4. What are we building (and not building)?
5. What could go wrong?

A good PRD prevents misalignment, scope creep, and building the wrong thing.