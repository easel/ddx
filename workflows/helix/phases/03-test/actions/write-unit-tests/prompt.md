# Write Unit Tests Prompt

Create comprehensive unit tests for all business logic, pure functions, and isolated components. These tests should be fast, focused, and independent, forming the foundation of the testing pyramid.

## Test Output Location

Generate unit tests in: `tests/unit/`

Organize tests by component type:
- `tests/unit/models/` - Data model tests
- `tests/unit/services/` - Business logic tests
- `tests/unit/utils/` - Utility function tests
- `tests/unit/components/` - Component tests

## Purpose

Unit tests verify that individual components work correctly in isolation. They should:
- Test single units of functionality
- Run quickly (< 100ms per test)
- Have no external dependencies
- Provide immediate feedback during development
- Enable safe refactoring

## Test Requirements

### Coverage Targets
- Minimum 80% code coverage for business logic
- 100% coverage for critical algorithms
- 100% coverage for data transformation functions
- Edge cases and error conditions covered
- All public methods/functions tested

### Test Structure

Follow the AAA (Arrange-Act-Assert) pattern:

```javascript
describe('ComponentName', () => {
  it('should perform expected behavior when given specific input', () => {
    // Arrange - Set up test data and dependencies
    const input = { /* test data */ };
    const expected = { /* expected result */ };

    // Act - Execute the function/method being tested
    const actual = functionUnderTest(input);

    // Assert - Verify the outcome
    expect(actual).toEqual(expected);
  });
});
```

## What to Test

### Pure Functions
- Input/output transformations
- Calculation accuracy
- Algorithm correctness
- Data validation logic
- Formatting functions

### Class Methods
- State changes
- Method interactions
- Constructor initialization
- Getter/setter behavior
- Private method effects (through public interface)

### Error Handling
- Invalid input handling
- Null/undefined checks
- Exception throwing
- Error message accuracy
- Recovery behavior

### Edge Cases
- Boundary values
- Empty collections
- Maximum/minimum values
- Special characters
- Concurrent operations

## Mocking Strategy

### When to Mock
- External services (APIs, databases)
- File system operations
- Network requests
- Time-dependent functions
- Random number generators

### How to Mock
```javascript
// Example: Mocking a database service
const mockDb = {
  query: jest.fn().mockResolvedValue([{ id: 1, name: 'Test' }]),
  insert: jest.fn().mockResolvedValue({ id: 2 }),
  update: jest.fn().mockResolvedValue(true),
  delete: jest.fn().mockResolvedValue(true)
};

// Inject mock in test
const service = new UserService(mockDb);
```

## Test Organization

### File Structure
```
tests/unit/
├── models/
│   ├── user.test.js
│   └── product.test.js
├── services/
│   ├── auth.test.js
│   └── payment.test.js
├── utils/
│   ├── validators.test.js
│   └── formatters.test.js
└── components/
    ├── button.test.js
    └── form.test.js
```

### Naming Conventions
- Test files: `{component}.test.{ext}`
- Test suites: Describe the component/module
- Test cases: Start with "should" and describe behavior
- Use descriptive names that explain the scenario

## Best Practices

### DO
- ✅ Keep tests simple and focused
- ✅ Test one thing per test
- ✅ Use descriptive test names
- ✅ Make tests deterministic
- ✅ Clean up after tests
- ✅ Use test data builders
- ✅ Group related tests
- ✅ Test public interfaces

### DON'T
- ❌ Test implementation details
- ❌ Make tests dependent on each other
- ❌ Use production data
- ❌ Test external libraries
- ❌ Write complex test logic
- ❌ Ignore flaky tests
- ❌ Leave console.logs in tests
- ❌ Test private methods directly

## Common Patterns

### Data Builders
```javascript
class UserBuilder {
  constructor() {
    this.user = {
      id: 1,
      name: 'Default User',
      email: 'user@example.com'
    };
  }

  withName(name) {
    this.user.name = name;
    return this;
  }

  withEmail(email) {
    this.user.email = email;
    return this;
  }

  build() {
    return this.user;
  }
}

// Usage in tests
const user = new UserBuilder()
  .withName('Test User')
  .withEmail('test@example.com')
  .build();
```

### Parameterized Tests
```javascript
describe.each([
  [1, 1, 2],
  [2, 3, 5],
  [3, 5, 8],
])('add(%i, %i)', (a, b, expected) => {
  test(`returns ${expected}`, () => {
    expect(add(a, b)).toBe(expected);
  });
});
```

## Quality Checklist

Before considering unit tests complete:
- [ ] All public methods/functions have tests
- [ ] Edge cases are covered
- [ ] Error scenarios are tested
- [ ] Tests are independent and isolated
- [ ] Tests run quickly (< 5 seconds for suite)
- [ ] No hardcoded values or magic numbers
- [ ] Clear test names describe behavior
- [ ] Mocks are properly cleaned up
- [ ] Coverage targets are met
- [ ] Tests are maintainable

## Framework-Specific Guidance

Based on your technology stack, consider:
- **JavaScript/TypeScript**: Jest, Mocha, Vitest
- **Python**: pytest, unittest
- **Go**: Built-in testing package
- **Java**: JUnit, TestNG
- **C#**: NUnit, xUnit
- **Ruby**: RSpec, Minitest

Remember: Unit tests are the foundation of your test suite. They should be numerous, fast, and reliable, giving developers confidence to refactor and extend the codebase.