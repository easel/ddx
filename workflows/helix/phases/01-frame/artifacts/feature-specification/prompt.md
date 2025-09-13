# Specification Generation Prompt

Create a complete feature specification with explicit uncertainty markers.

## Critical Instructions

### 1. Use [NEEDS CLARIFICATION] Markers
**DO NOT GUESS OR ASSUME**. When information is not explicitly provided:
- Mark it with `[NEEDS CLARIFICATION: specific question]`
- Be specific about what needs clarification
- It's better to have many clarification markers than hidden assumptions

### 2. Focus on WHAT, not HOW
- ✅ Describe user needs and business requirements
- ✅ Define success criteria and constraints
- ❌ Don't specify implementation details
- ❌ Don't choose technologies or architectures

### 3. Make Everything Testable
Every requirement must be:
- **Specific**: No ambiguous terms like "fast" or "user-friendly"
- **Measurable**: Include concrete metrics
- **Observable**: Can be verified through testing

## Required Sections

1. **Problem Statement**: What problem does this solve?
2. **Functional Requirements**: What must the system do?
3. **Non-Functional Requirements**: How well must it perform?
4. **User Stories**: Who uses it and why?
5. **Edge Cases**: What could go wrong?
6. **Success Metrics**: How do we measure success?

## Quality Checklist

Before completing the specification:
- [ ] Every ambiguity is marked with [NEEDS CLARIFICATION]
- [ ] All requirements are testable
- [ ] Success metrics are measurable
- [ ] User stories have clear acceptance criteria
- [ ] Edge cases are identified
- [ ] No implementation details included

## Example Clarification Markers

Good:
- `[NEEDS CLARIFICATION: Maximum response time in milliseconds?]`
- `[NEEDS CLARIFICATION: Which user roles need access?]`
- `[NEEDS CLARIFICATION: Data retention period?]`

Bad:
- `[NEEDS CLARIFICATION: Performance requirements?]` (too vague)
- "The system should be fast" (missing clarification marker)
- "Use PostgreSQL for data storage" (implementation detail)

Remember: The specification drives everything that follows. Ambiguity here cascades into problems throughout development.