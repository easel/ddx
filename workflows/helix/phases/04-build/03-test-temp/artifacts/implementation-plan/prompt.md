# Implementation Plan Generation Prompt

Create a detailed implementation plan that respects project principles and follows test-first development.

## Critical Constraints

### 1. Principles Gate
You MUST check compliance with all project principles:
- Start with ≤3 components (add more only with justification)
- Design as libraries with CLI interfaces
- Plan for test-first implementation
- Use frameworks directly without abstraction layers

### 2. Test-First Planning
The plan must support writing tests before code:
- Define contracts/interfaces first
- Specify test scenarios before implementation details
- Plan for failing tests before passing tests

### 3. Simplicity First
- Start with the minimal viable architecture
- Add complexity only when requirements demand it
- Document any complexity with justification

## Planning Process

### Step 1: Review Requirements
- Re-read the specification and user stories
- Identify the minimal feature set
- List non-negotiable requirements

### Step 2: Design Minimal Architecture
- What are the absolutely necessary components?
- Can this be done with 1-2 components instead of 3?
- What's the simplest data flow that works?

### Step 3: Define Interfaces First
- What are the public APIs?
- What CLI commands are needed?
- What are the input/output formats?

### Step 4: Plan Test Strategy
- What contracts need testing?
- What integration scenarios are critical?
- What unit tests are actually necessary?

### Step 5: Choose Technology
- Pick boring, proven technology
- Prefer standard library over external dependencies
- Use the framework's built-in features

## Anti-Patterns to Avoid

❌ **Over-Engineering**
- Don't create abstraction layers "just in case"
- Don't add features not in requirements
- Don't optimize prematurely

❌ **Test-After Thinking**
- Don't plan implementation before tests
- Don't skip contract definitions
- Don't assume interfaces

❌ **Complexity Creep**
- Don't exceed 3 components without justification
- Don't add unnecessary dependencies
- Don't create deep inheritance hierarchies

## Quality Checklist

Before finalizing the plan:
- [ ] Does it pass all principle gates?
- [ ] Are interfaces defined before implementation?
- [ ] Is the architecture minimal but sufficient?
- [ ] Are technology choices justified?
- [ ] Is the test strategy clear?
- [ ] Are risks identified and mitigated?

## Example: Good vs Bad Architecture

### ❌ Bad: Over-Engineered
```
Repository -> Service -> Controller -> DTO -> Mapper -> 
Validator -> Logger -> Cache -> Queue -> ...
```

### ✅ Good: Simple and Direct
```
CLI -> Library -> Database
```

Start simple. Add complexity only when proven necessary.