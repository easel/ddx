# User Stories Generation Prompt

Create comprehensive user stories with testable acceptance criteria.

## Core Principles

### 1. User-Centric Focus
- Start with WHO will use this (be specific about user types)
- Identify WHAT they need to do (specific actions)
- Explain WHY it matters (business value)

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

## Required Elements for Each Story

1. **Clear User Type**: Not just "user" but specific role
2. **Specific Functionality**: Concrete action, not vague desire
3. **Business Value**: Why this matters to the organization
4. **Acceptance Criteria**: 3-5 specific, testable criteria
5. **Definition of Done**: Clear completion checklist

## Anti-Patterns to Avoid

❌ **Too Large**: "As a user, I want a complete system"
✅ **Right-Sized**: "As a user, I want to reset my password"

❌ **Vague Criteria**: "System should be fast"
✅ **Specific Criteria**: "Password reset email arrives within 30 seconds"

❌ **Technical Tasks**: "As a developer, I want to refactor the database"
✅ **User Value**: "As a user, I want faster search results"

## Story Splitting Techniques

If a story is too large, split by:
1. **Workflow Steps**: Each step becomes a story
2. **Business Rules**: Simple case first, complex cases later
3. **Data Types**: One entity type per story
4. **Operations**: CRUD operations as separate stories

## Priority Guidelines

**P0 - Must Have**: Core functionality, without this the product has no value
**P1 - Should Have**: Important features that users expect
**P2 - Nice to Have**: Enhancements that delight users

## Validation Questions

For each story, ask:
1. Can we demo this to a user?
2. Can we write automated tests for this?
3. Does this deliver value even if nothing else is built?
4. Is this small enough to complete in one iteration?
5. Are the acceptance criteria unambiguous?

## Example Story

### Good Example:
**As a** registered customer  
**I want** to save items to a wishlist  
**So that** I can purchase them later without searching again

**Acceptance Criteria:**
- [ ] Given I am logged in, when I click "Add to Wishlist" on a product, then the product appears in my wishlist
- [ ] Given I have items in my wishlist, when I view my wishlist, then I see all saved items with current prices
- [ ] Given an item is in my wishlist, when I click "Remove", then it is removed immediately

This story is specific, testable, and delivers clear value.