# Test Specifications Generation Prompt

Create comprehensive test specifications that define all tests BEFORE writing any implementation code.

## Test-First Philosophy

### Why Test-First?
1. **Forces clear thinking** about behavior before implementation
2. **Prevents untestable code** by designing for testability
3. **Documents expected behavior** through executable specifications
4. **Enables confident refactoring** with safety net in place

### The Red-Green-Refactor Cycle
1. **Red**: Write tests that fail (no implementation yet)
2. **Green**: Write minimal code to make tests pass
3. **Refactor**: Improve code while tests stay green

## Test Pyramid Strategy

### Level 1: Contract Tests (Foundation - Most Tests)
Test the **external behavior** that users/systems depend on:
- CLI commands and their outputs
- API endpoints and responses
- Library public interfaces
- File formats and protocols

### Level 2: Integration Tests (Middle - Some Tests)
Test **component interactions**:
- Database operations
- Service communications
- File system operations
- Third-party integrations

### Level 3: Unit Tests (Top - Few Tests)
Test **complex internal logic** only:
- Algorithms with multiple paths
- Business rule calculations
- Data transformations
- Error handling logic

## Writing Good Test Specifications

### For Contract Tests
Define:
- **Exact input format** (including edge cases)
- **Expected output format** (including errors)
- **Side effects** (files created, data stored)
- **Performance requirements** (response time, throughput)

Example:
```
Test: Parse JSON Command
Input: {"key": "value"}
Output: Validated and formatted JSON
Error Cases: 
  - Malformed JSON → Error code 1
  - Invalid schema → Error code 2
Performance: < 100ms for files under 1MB
```

### For Integration Tests
Define:
- **Component boundaries** being tested
- **Real dependencies** to use (not mocks)
- **Failure scenarios** to handle
- **Concurrency behavior** if applicable

### For Unit Tests
Define:
- **Function signature** and purpose
- **Input/output pairs** covering all paths
- **Edge cases** and boundaries
- **Error conditions** and exceptions

## Test Data Strategy

### Use Realistic Data
- Real-world examples, not "foo/bar"
- Production-like volumes
- Actual edge cases from users

### Categories to Cover
1. **Happy Path**: Normal, expected inputs
2. **Boundary Values**: Min, max, empty, full
3. **Invalid Inputs**: Wrong types, formats, ranges
4. **Malicious Inputs**: Injection attempts, overflows
5. **Special Characters**: Unicode, quotes, nulls

## Performance Test Specifications

Define specific, measurable requirements:
- **Response Time**: 95th percentile < X ms
- **Throughput**: X operations per second
- **Resource Usage**: Memory < X MB, CPU < X%
- **Scalability**: Linear up to X concurrent users

## Security Test Specifications

Cover OWASP Top 10:
- Input validation (injection prevention)
- Authentication and session management
- Access control and authorization
- Security misconfiguration detection
- Sensitive data exposure prevention

## Test Organization

### Naming Convention
```
test_[what]_[condition]_[expected_result]
```
Examples:
- `test_parse_valid_json_returns_formatted_output`
- `test_api_unauthorized_request_returns_401`
- `test_calculator_divide_by_zero_throws_error`

### Test Grouping
- Group by feature or component
- Separate fast tests from slow tests
- Isolate tests requiring external resources

## Anti-Patterns to Avoid

### ❌ Testing Implementation Details
**Bad**: Test that specific method was called
**Good**: Test observable behavior/output

### ❌ Excessive Mocking
**Bad**: Mock every dependency
**Good**: Use real dependencies where possible

### ❌ Brittle Tests
**Bad**: Tests that break with valid refactoring
**Good**: Tests that only break when behavior changes

### ❌ Non-Deterministic Tests
**Bad**: Tests that sometimes pass/fail
**Good**: Tests that always produce same result

## Quality Checklist

Before finalizing test specifications:
- [ ] Do tests cover all user stories?
- [ ] Are contract tests comprehensive?
- [ ] Are edge cases identified?
- [ ] Are performance requirements testable?
- [ ] Are tests independent of each other?
- [ ] Can tests run in any order?
- [ ] Are test names self-documenting?

## Remember

**Tests are specifications**. They define what the system should do. Write them to:
1. Fail initially (proving they detect issues)
2. Pass when implementation is correct
3. Serve as living documentation
4. Enable fearless refactoring

Good tests make development faster, not slower, by catching issues early and enabling confident changes.