# HELIX Action: Write Story Tests

You are a HELIX Test phase executor tasked with creating comprehensive test suites for user stories following Test-Driven Development principles. Your role is to write failing tests that define system behavior before implementation begins.

## Action Purpose

Create failing tests that serve as executable specifications for user stories, ensuring implementation matches requirements and enabling the TDD Red-Green-Refactor cycle.

## When to Use This Action

- After user stories are defined with clear acceptance criteria
- Before implementation of story functionality begins
- When TDD approach requires failing tests first
- When test specifications need to be created

## Prerequisites

- [ ] User story with clear acceptance criteria
- [ ] Design specifications available
- [ ] Test environment configured
- [ ] Testing frameworks selected and set up
- [ ] Test data strategy defined

## Action Workflow

### 1. Test Analysis and Planning

**User Story Test Analysis**:
```
🧪 TEST PLANNING FOR USER STORY: [Story ID]

Story: [Story title and description]

1. ACCEPTANCE CRITERIA ANALYSIS
   For each acceptance criterion:
   - What specific behavior does this define?
   - What are the inputs and expected outputs?
   - What edge cases need to be tested?
   - What error conditions should be handled?

2. TEST LEVEL STRATEGY
   - Unit Tests: [Component/function level tests needed]
   - Integration Tests: [Component interaction tests]
   - Contract Tests: [API/interface behavior tests]
   - End-to-End Tests: [Full user journey tests]

3. TEST DATA REQUIREMENTS
   - What test data is needed?
   - What test scenarios require specific data setups?
   - What boundary conditions need data?
   - What error scenarios need data?

4. TEST ENVIRONMENT NEEDS
   - What dependencies need to be mocked/stubbed?
   - What external services need test doubles?
   - What infrastructure is required?
   - What configuration is needed?
```

### 2. Test Design and Creation

**Test Structure Template**:
```javascript
// Test Suite for User Story: [Story ID]
describe('User Story: [Story Title]', () => {

  describe('Acceptance Criterion 1: [AC Description]', () => {

    describe('Happy Path Scenarios', () => {
      it('should [expected behavior] when [condition]', () => {
        // Arrange
        // Set up test data and conditions

        // Act
        // Execute the behavior being tested

        // Assert
        // Verify the expected outcome

        // This test should FAIL initially (Red phase)
      });
    });

    describe('Edge Cases', () => {
      it('should handle [edge case] appropriately', () => {
        // Test edge case behavior
      });
    });

    describe('Error Scenarios', () => {
      it('should [error handling] when [error condition]', () => {
        // Test error handling behavior
      });
    });
  });

  describe('Acceptance Criterion 2: [AC Description]', () => {
    // Additional test scenarios for second AC
  });
});
```

### 3. Test Implementation

**Unit Test Examples**:
```javascript
// Unit Tests - Component Behavior
describe('UserRegistrationService', () => {
  let userService;
  let mockEmailService;
  let mockUserRepository;

  beforeEach(() => {
    mockEmailService = new MockEmailService();
    mockUserRepository = new MockUserRepository();
    userService = new UserRegistrationService(mockEmailService, mockUserRepository);
  });

  describe('registerUser', () => {
    it('should create user account when valid data provided', async () => {
      // Arrange
      const userData = {
        email: 'test@example.com',
        password: 'SecurePass123!',
        name: 'John Doe'
      };

      // Act
      const result = await userService.registerUser(userData);

      // Assert
      expect(result.success).toBe(true);
      expect(result.userId).toBeDefined();
      expect(mockUserRepository.save).toHaveBeenCalledWith(
        expect.objectContaining({
          email: userData.email,
          name: userData.name
        })
      );
      expect(mockEmailService.sendWelcomeEmail).toHaveBeenCalledWith(userData.email);
    });

    it('should reject registration when email already exists', async () => {
      // Arrange
      const userData = { email: 'existing@example.com', password: 'Pass123!', name: 'Jane Doe' };
      mockUserRepository.findByEmail.mockResolvedValue({ id: 1, email: userData.email });

      // Act & Assert
      await expect(userService.registerUser(userData)).rejects.toThrow('Email already registered');
      expect(mockUserRepository.save).not.toHaveBeenCalled();
      expect(mockEmailService.sendWelcomeEmail).not.toHaveBeenCalled();
    });
  });
});
```

**Integration Test Examples**:
```javascript
// Integration Tests - Component Interaction
describe('User Registration Integration', () => {
  let app;
  let testDb;

  beforeAll(async () => {
    testDb = await setupTestDatabase();
    app = createTestApp(testDb);
  });

  afterAll(async () => {
    await cleanupTestDatabase(testDb);
  });

  describe('POST /api/users/register', () => {
    it('should register user and return success response', async () => {
      // Arrange
      const userData = {
        email: 'integration@example.com',
        password: 'SecurePass123!',
        name: 'Integration User'
      };

      // Act
      const response = await request(app)
        .post('/api/users/register')
        .send(userData)
        .expect(201);

      // Assert
      expect(response.body).toMatchObject({
        success: true,
        userId: expect.any(String),
        message: 'User registered successfully'
      });

      // Verify user was saved to database
      const savedUser = await testDb.users.findByEmail(userData.email);
      expect(savedUser).toBeDefined();
      expect(savedUser.email).toBe(userData.email);
    });

    it('should return 400 when invalid email provided', async () => {
      // Arrange
      const invalidUserData = {
        email: 'invalid-email',
        password: 'SecurePass123!',
        name: 'Invalid User'
      };

      // Act & Assert
      const response = await request(app)
        .post('/api/users/register')
        .send(invalidUserData)
        .expect(400);

      expect(response.body).toMatchObject({
        success: false,
        error: 'Invalid email format'
      });
    });
  });
});
```

**End-to-End Test Examples**:
```javascript
// E2E Tests - Full User Journey
describe('User Registration Flow', () => {
  let browser;
  let page;

  beforeAll(async () => {
    browser = await puppeteer.launch();
    page = await browser.newPage();
  });

  afterAll(async () => {
    await browser.close();
  });

  it('should allow user to complete registration flow', async () => {
    // Arrange
    await page.goto('http://localhost:3000/register');

    // Act
    await page.type('#email', 'e2e@example.com');
    await page.type('#password', 'SecurePass123!');
    await page.type('#name', 'E2E User');
    await page.click('#register-button');

    // Wait for registration to complete
    await page.waitForSelector('#success-message', { timeout: 5000 });

    // Assert
    const successMessage = await page.$eval('#success-message', el => el.textContent);
    expect(successMessage).toContain('Registration successful');

    // Verify user is redirected to dashboard
    await page.waitForNavigation();
    expect(page.url()).toContain('/dashboard');
  });
});
```

### 4. Test Validation

**Test Quality Checklist**:
```markdown
For each test:
- [ ] Test has clear, descriptive name
- [ ] Test follows Arrange-Act-Assert pattern
- [ ] Test verifies specific acceptance criterion
- [ ] Test is independent and can run in isolation
- [ ] Test data is clearly defined
- [ ] Expected behavior is explicitly verified
- [ ] Error conditions are tested
- [ ] Test initially fails (Red phase requirement)
```

## Outputs

### Primary Artifacts
- **Unit Test Suite** → `tests/unit/[story-id]/`
- **Integration Test Suite** → `tests/integration/[story-id]/`
- **Contract Test Suite** → `tests/contract/[story-id]/`
- **End-to-End Test Suite** → `tests/e2e/[story-id]/`

### Supporting Artifacts
- **Test Plan Document** → `docs/helix/03-test/test-plans/[story-id]-test-plan.md`
- **Test Data Specifications** → `tests/data/[story-id]/`
- **Test Environment Setup** → `tests/setup/[story-id]/`

## Quality Gates

**Test Completion Criteria**:
- [ ] All acceptance criteria have corresponding tests
- [ ] Happy path scenarios are fully tested
- [ ] Edge cases and error conditions are covered
- [ ] Tests are independent and repeatable
- [ ] Test data is well-defined and reusable
- [ ] All tests initially fail (Red phase)
- [ ] Test coverage meets project standards
- [ ] Tests are readable and maintainable
- [ ] Performance test criteria defined (if applicable)
- [ ] Security test scenarios included (if applicable)

## Integration with Test Phase

This action supports the Test phase by:
- **Enabling TDD Cycle**: Creates failing tests for Red phase
- **Defining Behavior**: Tests serve as executable specifications
- **Supporting Quality Gates**: Tests validate implementation completeness
- **Facilitating Refactoring**: Tests enable safe code improvements

## Test Strategy Considerations

### Test Pyramid
- **Unit Tests (70%)**: Fast, isolated, specific
- **Integration Tests (20%)**: Component interactions
- **End-to-End Tests (10%)**: Full user journeys

### Test Types by Purpose
- **Functional Tests**: Verify business requirements
- **Performance Tests**: Validate speed and scalability
- **Security Tests**: Check for vulnerabilities
- **Usability Tests**: Confirm user experience
- **Accessibility Tests**: Ensure inclusive design

## Common Pitfalls to Avoid

❌ **Implementation-Driven Tests**: Writing tests after code is written
❌ **Brittle Tests**: Tests that break with minor code changes
❌ **Unclear Test Names**: Vague descriptions of what's being tested
❌ **Missing Edge Cases**: Only testing happy path scenarios
❌ **Dependent Tests**: Tests that require specific execution order
❌ **Over-Mocking**: Mocking everything instead of testing real interactions

## Success Criteria

This action succeeds when:
- ✅ Comprehensive test suite covers all acceptance criteria
- ✅ Tests initially fail (Red phase requirement met)
- ✅ Test structure enables clear TDD cycle
- ✅ Edge cases and error scenarios are tested
- ✅ Tests are maintainable and readable
- ✅ Test data strategy supports all scenarios
- ✅ Foundation is established for implementation phase
- ✅ Quality gates are defined and measurable

Remember: In TDD, tests are not just validation - they are the specification. Write them to fail, then make them pass.