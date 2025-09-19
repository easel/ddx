# Build Procedures

*Step-by-step guide for implementing code using Test-Driven Development*

## TDD Implementation Overview

This document defines HOW to systematically make failing tests pass during the Build phase, following the Red-Green-Refactor cycle.

## Build Phase Entry Criteria

### Prerequisites Checklist
- [ ] All tests written and failing (Red phase complete)
- [ ] Test coverage targets defined
- [ ] Development environment configured
- [ ] Implementation plan reviewed
- [ ] Coding standards documented

### Initial State Verification
```bash
# Verify all tests are failing
npm test 2>&1 | grep -E "(passing|failing)"
# Expected: 0 passing, N failing

# Check coverage baseline
npm test -- --coverage
# Expected: 0% coverage

# Verify no implementation exists
find src -name "*.js" -o -name "*.ts" | wc -l
# Expected: 0 or minimal scaffold files
```

## Implementation Procedures

### Procedure 1: Contract Test Implementation

#### Objective
Make external interface tests pass with minimal implementation.

#### Step-by-Step Process

1. **Select Next Failing Test**
   ```bash
   # Run tests to see failures
   npm test tests/contract/ --verbose

   # Pick highest priority failing test
   # Example: POST /api/users endpoint
   ```

2. **Create Minimal Implementation**
   ```javascript
   // src/api/users.js
   export function createUser(req, res) {
     // Minimal code to make test pass
     res.status(201).json({
       id: '123',
       email: req.body.email
     });
   }
   ```

3. **Wire Up Implementation**
   ```javascript
   // src/app.js
   import { createUser } from './api/users.js';

   app.post('/api/users', createUser);
   ```

4. **Run Test - Verify Green**
   ```bash
   npm test tests/contract/api/users.test.js
   # Expected: 1 passing
   ```

5. **Commit Immediately**
   ```bash
   git add .
   git commit -m "✅ Implement POST /api/users - contract test passing"
   ```

6. **Refactor (If Needed)**
   ```javascript
   // Improve code quality while keeping test green
   export async function createUser(req, res) {
     try {
       const user = await userService.create(req.body);
       res.status(201).json(user);
     } catch (error) {
       res.status(400).json({ error: error.message });
     }
   }
   ```

7. **Verify Test Still Passes**
   ```bash
   npm test tests/contract/api/users.test.js
   # Expected: Still passing
   ```

8. **Commit Refactoring**
   ```bash
   git add .
   git commit -m "♻️ Refactor user creation - add error handling"
   ```

### Procedure 2: Integration Test Implementation

#### Objective
Implement component interactions to make integration tests pass.

#### Step-by-Step Process

1. **Identify Failed Integration Test**
   ```bash
   npm test tests/integration/ --verbose
   # Pick next priority test
   ```

2. **Implement Service Layer**
   ```javascript
   // src/services/UserService.js
   export class UserService {
     constructor(database, emailService) {
       this.db = database;
       this.email = emailService;
     }

     async create(userData) {
       // Save to database
       const user = await this.db.users.create(userData);

       // Send welcome email
       await this.email.send({
         to: user.email,
         subject: 'Welcome',
         body: 'Welcome to our service!'
       });

       return user;
     }
   }
   ```

3. **Implement Dependencies**
   ```javascript
   // src/repositories/UserRepository.js
   export class UserRepository {
     constructor(database) {
       this.db = database;
     }

     async create(data) {
       const result = await this.db.query(
         'INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *',
         [data.email, data.password]
       );
       return result.rows[0];
     }
   }
   ```

4. **Run Integration Test**
   ```bash
   npm test tests/integration/UserService.test.js
   # Expected: Passing
   ```

5. **Commit Working Integration**
   ```bash
   git add .
   git commit -m "✅ Implement UserService with database and email integration"
   ```

### Procedure 3: Unit Test Implementation

#### Objective
Implement business logic to make unit tests pass.

#### Step-by-Step Process

1. **Select Failing Unit Test**
   ```bash
   npm test tests/unit/ --verbose
   ```

2. **Implement Business Logic**
   ```javascript
   // src/validators/passwordValidator.js
   export function validatePassword(password) {
     if (password.length < 8) return false;

     const hasUpper = /[A-Z]/.test(password);
     const hasLower = /[a-z]/.test(password);
     const hasNumber = /[0-9]/.test(password);
     const hasSpecial = /[!@#$%^&*]/.test(password);

     const complexity = [hasUpper, hasLower, hasNumber, hasSpecial]
       .filter(Boolean).length;

     return complexity >= 3;
   }
   ```

3. **Run Unit Test**
   ```bash
   npm test tests/unit/passwordValidator.test.js
   # Expected: Passing
   ```

4. **Refactor for Clarity**
   ```javascript
   const PASSWORD_RULES = {
     minLength: 8,
     patterns: [
       { regex: /[A-Z]/, name: 'uppercase' },
       { regex: /[a-z]/, name: 'lowercase' },
       { regex: /[0-9]/, name: 'number' },
       { regex: /[!@#$%^&*]/, name: 'special' }
     ],
     minComplexity: 3
   };

   export function validatePassword(password) {
     if (password.length < PASSWORD_RULES.minLength) {
       return false;
     }

     const complexity = PASSWORD_RULES.patterns
       .filter(rule => rule.regex.test(password))
       .length;

     return complexity >= PASSWORD_RULES.minComplexity;
   }
   ```

5. **Verify and Commit**
   ```bash
   npm test tests/unit/passwordValidator.test.js
   git commit -m "♻️ Refactor password validator - improve readability"
   ```

## Build Patterns and Techniques

### Pattern: Incremental Implementation

#### Start Simple
```javascript
// Step 1: Hardcoded response
function getUser(id) {
  return { id: '123', name: 'John' };
}

// Step 2: Add parameter usage
function getUser(id) {
  return { id: id, name: 'John' };
}

// Step 3: Add database
async function getUser(id) {
  return await db.users.findById(id);
}

// Step 4: Add error handling
async function getUser(id) {
  const user = await db.users.findById(id);
  if (!user) throw new NotFoundError('User not found');
  return user;
}
```

### Pattern: Dependency Injection

#### Make Components Testable
```javascript
// Bad: Hard dependencies
class UserService {
  constructor() {
    this.db = new Database(); // Hard to test
    this.email = new EmailService(); // Hard to mock
  }
}

// Good: Injected dependencies
class UserService {
  constructor(database, emailService) {
    this.db = database;
    this.email = emailService;
  }
}

// Easy to test
const service = new UserService(mockDb, mockEmail);
```

## Daily Build Workflow

### Morning Routine

1. **Review Test Status**
   ```bash
   # See what's left to implement
   npm test -- --listTests | grep failing
   ```

2. **Pick Next Test Set**
   - Choose related tests to implement together
   - Prioritize P0 tests first

3. **Set Up Watch Mode**
   ```bash
   # Watch specific test file
   npm test -- --watch tests/contract/users.test.js
   ```

### Implementation Cycle

```bash
# 1. See test fail
npm test path/to/test.js
# RED ❌

# 2. Write minimal code
vim src/implementation.js

# 3. See test pass
npm test path/to/test.js
# GREEN ✅

# 4. Commit immediately
git add . && git commit -m "✅ Make [test] pass"

# 5. Refactor if needed
vim src/implementation.js

# 6. Ensure still green
npm test path/to/test.js
# STILL GREEN ✅

# 7. Commit refactoring
git add . && git commit -m "♻️ Refactor [what]"
```

### End of Day Review

```bash
# Check progress
npm test 2>&1 | grep -E "(passing|failing)"
# Example: 45 passing, 5 failing

# Check coverage increase
npm test -- --coverage --coverageReporters=text-summary

# Review commits
git log --oneline --since="9am"

# Push to remote
git push origin feature-branch
```

## Common Implementation Scenarios

### Scenario: Database Operations

#### Making Database Tests Pass
```javascript
// Step 1: In-memory implementation
const users = [];

export function createUser(data) {
  const user = { ...data, id: users.length + 1 };
  users.push(user);
  return user;
}

// Step 2: Real database
import { pool } from './database.js';

export async function createUser(data) {
  const result = await pool.query(
    'INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *',
    [data.email, data.password]
  );
  return result.rows[0];
}
```

### Scenario: External API Integration

#### Making API Tests Pass
```javascript
// Step 1: Mock response
export function callExternalAPI() {
  return { status: 'success', data: {} };
}

// Step 2: Real integration
import axios from 'axios';

export async function callExternalAPI(params) {
  const response = await axios.post(API_URL, params);
  return response.data;
}
```

## Troubleshooting Guide

### Problem: Test Won't Pass

#### Diagnosis Steps
1. **Verify Test Is Correct**
   - Review test expectations
   - Check test data

2. **Debug Implementation**
   ```javascript
   console.log('Input:', input);
   console.log('Output:', output);
   console.log('Expected:', expected);
   ```

3. **Simplify**
   - Remove complexity
   - Make test pass with hardcoded values first

### Problem: Tests Pass But Break Others

#### Solution
1. **Run All Tests**
   ```bash
   npm test
   ```

2. **Find Breaking Change**
   ```bash
   git diff HEAD~1
   ```

3. **Fix Without Breaking**
   - Adjust implementation
   - May need to refactor shared code

### Problem: Can't Make Test Pass

#### Escalation Path
1. **Pair Programming**
   - Get another developer
   - Work through together

2. **Review Test**
   - Test might be wrong
   - Requirements might have changed

3. **Spike Solution**
   - Time-boxed exploration
   - Throw away and reimplement

## Performance Monitoring

### Build Metrics

Track daily progress:
```markdown
| Date | Tests Passing | Coverage | Commits | Notes |
|------|--------------|----------|---------|-------|
| Day 1 | 10/100 | 15% | 25 | Core infrastructure |
| Day 2 | 35/100 | 40% | 30 | API endpoints |
| Day 3 | 70/100 | 65% | 28 | Business logic |
| Day 4 | 95/100 | 82% | 20 | Edge cases |
| Day 5 | 100/100 | 90% | 15 | Complete! |
```

### Quality Gates

Before marking implementation complete:
- [ ] All P0 tests passing
- [ ] Coverage target met (e.g., 80%)
- [ ] No linting errors
- [ ] Code reviewed
- [ ] Documentation updated

## Handoff to Deploy Phase

### Completion Checklist
- [ ] All tests passing
- [ ] Coverage targets met
- [ ] Code refactored and clean
- [ ] Performance benchmarks met
- [ ] Security scan passed
- [ ] Documentation complete

### Build Artifacts
1. **Source Code**: Complete implementation
2. **Test Results**: All green
3. **Coverage Report**: Meeting targets
4. **Build Output**: Compiled/bundled application
5. **Release Notes**: What was implemented

---
*These procedures ensure disciplined TDD implementation throughout the Build phase.*