# Test Suites Generation Prompt

Create comprehensive test suites that define system behavior BEFORE any implementation exists. These tests will fail initially (Red phase of TDD) and drive the implementation in the Build phase.

## Storage Location

Test suites should be organized in: `tests/` directory at project root

## Core Principle: Tests Define the Specification

**Tests are the executable specification.** Every test you write defines what the system MUST do. No production code should exist yet - only tests that describe desired behavior.

## Test Organization Structure

### 1. Contract Tests (`tests/contract/`)
**Purpose**: Define external API behavior
- Test all public interfaces from the API contracts
- Verify request/response formats
- Test error responses and status codes
- Cover authentication and authorization
- No implementation should exist yet

### 2. Integration Tests (`tests/integration/`)
**Purpose**: Define component interactions
- Test data flow between components
- Verify state management behavior
- Test external service integrations
- Database operations (using test database)
- Message queue interactions

### 3. Unit Tests (`tests/unit/`)
**Purpose**: Define internal logic
- Test business logic rules
- Data transformation functions
- Validation logic
- Utility functions
- Edge cases and boundary conditions

### 4. End-to-End Tests (`tests/e2e/`)
**Purpose**: Define user journeys
- Critical user workflows
- Multi-step processes
- Cross-system interactions
- Real-world scenarios

## Writing Tests in the Red Phase

### Key Requirements
1. **Tests Must Fail**: Since no implementation exists, all tests should fail
2. **Tests Must Be Executable**: Use mocks/stubs where needed to make tests runnable
3. **Tests Must Be Specific**: Clear assertions about expected behavior
4. **Tests Must Be Complete**: Cover happy path, edge cases, and error scenarios

### Test Structure Pattern

```javascript
describe('Component/Feature Name', () => {
  describe('Scenario/Method', () => {
    it('should [specific behavior] when [condition]', () => {
      // Arrange - Set up test data and conditions

      // Act - Call the function/API (will fail - no implementation)

      // Assert - Define expected outcome
    });

    it('should handle [error case] when [invalid condition]', () => {
      // Test error handling
    });
  });
});
```

## Coverage Requirements

Based on the test specifications from Design phase:

### Minimum Coverage Targets
- **Contract Tests**: 100% of all API endpoints
- **Integration Tests**: All component boundaries
- **Unit Tests**: All business logic functions
- **E2E Tests**: All critical user paths (P0 requirements)

### What to Test

#### From User Stories
- Every acceptance criterion needs a test
- Each "Given/When/Then" becomes a test case
- Definition of Done requirements

#### From Specifications
- Functional requirements → Contract/Integration tests
- Non-functional requirements → Performance/Load tests
- Edge cases → Unit tests with boundary values
- Error handling → Negative test cases

#### From Contracts
- Every endpoint → Contract test
- Every request format → Validation test
- Every response format → Assertion test
- Every error code → Error test

## Test Data Management

### Fixtures (`tests/fixtures/`)
- Valid request payloads
- Expected response data
- Test user accounts
- Sample database records

### Factories (`tests/factories/`)
- Functions to generate test data
- Builders for complex objects
- Random data generators for property testing

### Mocks (`tests/mocks/`)
- External service mocks
- Database mocks (if needed)
- API client mocks
- Time/Date mocks

## Practical Examples

### Contract Test Example
```typescript
// tests/contract/users.test.ts
describe('POST /api/users', () => {
  it('should create a user with valid data', async () => {
    const userData = {
      email: 'test@example.com',
      name: 'Test User',
      role: 'standard'
    };

    // This will fail - no implementation yet
    const response = await request(app)
      .post('/api/users')
      .send(userData);

    expect(response.status).toBe(201);
    expect(response.body).toMatchObject({
      id: expect.any(String),
      email: userData.email,
      name: userData.name
    });
  });

  it('should reject invalid email format', async () => {
    const userData = {
      email: 'invalid-email',
      name: 'Test User'
    };

    const response = await request(app)
      .post('/api/users')
      .send(userData);

    expect(response.status).toBe(400);
    expect(response.body.error).toContain('email');
  });
});
```

### Integration Test Example
```typescript
// tests/integration/user-service.test.ts
describe('UserService', () => {
  it('should save user to database and send welcome email', async () => {
    const userData = { email: 'new@example.com', name: 'New User' };

    // These will fail - no implementation
    const user = await userService.create(userData);

    expect(user.id).toBeDefined();
    expect(mockEmailService.sendWelcome).toHaveBeenCalledWith(userData.email);
    expect(await userRepository.findById(user.id)).toEqual(user);
  });
});
```

## Test Naming Conventions

### File Names
- Contract tests: `[resource].contract.test.ts`
- Integration tests: `[feature].integration.test.ts`
- Unit tests: `[module].unit.test.ts`
- E2E tests: `[journey].e2e.test.ts`

### Test Descriptions
- Use "should" for behavior: "should return user when ID exists"
- Be specific about conditions: "when invalid ID provided"
- Include expected outcome: "returns 404 error"

## CI/CD Integration

### Test Execution Order
1. Unit tests (fastest, run first)
2. Integration tests
3. Contract tests
4. E2E tests (slowest, run last)

### Pipeline Configuration
```yaml
test:
  stage: test
  script:
    - npm run test:unit
    - npm run test:integration
    - npm run test:contract
    - npm run test:e2e
  coverage: '/Coverage: \d+\.\d+%/'
```

## Quality Checklist

Before completing test suite generation:
- [ ] All API endpoints have contract tests
- [ ] All user stories have corresponding tests
- [ ] Edge cases are covered
- [ ] Error scenarios are tested
- [ ] Tests are properly organized by type
- [ ] Test data is managed through fixtures/factories
- [ ] All tests fail (red phase verified)
- [ ] Tests are added to CI pipeline
- [ ] Coverage targets are defined

## Anti-Patterns to Avoid

### ❌ Don't Write Implementation
- No production code in this phase
- Only test code and test infrastructure

### ❌ Don't Skip Error Cases
- Error handling tests are as important as happy path
- Test all validation rules

### ❌ Don't Write Vague Tests
- Bad: "should work correctly"
- Good: "should return 201 and user object when valid email provided"

### ❌ Don't Ignore Test Organization
- Keep tests organized by type
- Maintain clear separation between test levels

### ❌ Don't Forget About Test Data
- Set up proper fixtures and factories
- Avoid hardcoded test data scattered in tests

## Next Phase: Build

Once all tests are written and failing:
1. Build phase begins with clear definition of "done"
2. Implementation proceeds to make tests pass
3. No code without a failing test demanding it
4. Refactor only when tests are green

Remember: **In TDD, the test comes first. The test IS the specification.**