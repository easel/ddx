# User Stories Generation Prompt

Create comprehensive user stories with testable acceptance criteria.

## Storage Location

Store user stories at: `docs/helix/01-frame/user-stories/US-{number}-{title-with-hyphens}.md`

## Naming Convention

Follow this format:
- **File Format**: `US-{number}-{title-with-hyphens}.md`
- **Number**: Zero-padded 3-digit sequence (001, 002, 003...)
- **Title**: Descriptive, lowercase with hyphens

Examples: `US-001-user-login-authentication.md`, `US-002-project-template-selection.md`

## Core Principles

### 1. User-Centric Focus
- **WHO** will use this (be specific about user types)
- **WHAT** they need to do (specific actions)
- **WHY** it matters (business value)

### 2. Testable Acceptance Criteria
Each criterion must follow the Given-When-Then format:
- **Given** [initial context/state]
- **When** [action taken]
- **Then** [expected outcome]

### 3. Independent and Valuable
Each story should:
- Deliver value on its own
- Be completable in one iteration
- Not depend on unwritten stories

## Required Elements

1. **Clear User Type**: Specific role, not just "user"
2. **Specific Functionality**: Concrete action, not vague desire
3. **Business Value**: Why this matters
4. **Acceptance Criteria**: 3-5 specific, testable criteria
5. **Definition of Done**: Clear completion checklist

## Anti-Patterns

❌ **Too Large**: "As a user, I want a complete system"
✅ **Right-Sized**: "As a user, I want to reset my password"

❌ **Vague Criteria**: "System should be fast"
✅ **Specific Criteria**: "Password reset email arrives within 30 seconds"

❌ **Technical Tasks**: "As a developer, I want to refactor the database"
✅ **User Value**: "As a user, I want faster search results"

## Story Splitting

If a story is too large, split by:
1. **Workflow Steps**: Each step becomes a story
2. **Business Rules**: Simple case first, complex cases later
3. **Data Types**: One entity type per story
4. **Operations**: CRUD operations as separate stories

## Priority Guidelines

- **P0 - Must Have**: Core functionality without which the product has no value
- **P1 - Should Have**: Important features that users expect
- **P2 - Nice to Have**: Enhancements that delight users

## Validation Questions

For each story, ask:
1. Can we demo this to a user?
2. Can we write automated tests for this?
3. Does this deliver value even if nothing else is built?
4. Is this small enough to complete in one iteration?
5. Are the acceptance criteria unambiguous?

---

Remember: Good user stories enable teams to deliver value incrementally while maintaining clear traceability to user needs.