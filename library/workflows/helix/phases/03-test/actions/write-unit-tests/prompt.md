# Write Unit Tests Prompt

Create comprehensive unit tests for all business logic, pure functions, and isolated components. These tests should be fast, focused, and independent, forming the foundation of the testing pyramid.

## Test Output Location

Generate unit tests in: `tests/unit/`

Organize by component type:
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
- **Arrange**: Set up test data and dependencies
- **Act**: Execute the function/method being tested
- **Assert**: Verify the outcome

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

## Mocking Strategy

### When to Mock
- External services (APIs, databases)
- File system operations
- Network requests
- Time-dependent functions
- Random number generators

### How to Mock
Create test doubles that return predictable values. Inject mocks in tests to isolate the unit being tested.

## Test Organization

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
- ✅ Group related tests
- ✅ Test public interfaces

### DON'T
- ❌ Test implementation details
- ❌ Make tests dependent on each other
- ❌ Use production data
- ❌ Test external libraries
- ❌ Write complex test logic
- ❌ Ignore flaky tests
- ❌ Test private methods directly

## Quality Checklist

Before considering unit tests complete:
- All public methods/functions have tests
- Edge cases are covered
- Error scenarios are tested
- Tests are independent and isolated
- Tests run quickly (< 5 seconds for suite)
- Clear test names describe behavior
- Mocks are properly cleaned up
- Coverage targets are met

---

Remember: Unit tests are the foundation of your test suite. They should be numerous, fast, and reliable, giving developers confidence to refactor and extend the codebase.