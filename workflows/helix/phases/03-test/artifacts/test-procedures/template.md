# Test Procedures

*Step-by-step guide for executing tests in the TDD Red phase*

## Test Execution Overview

This document defines HOW tests will be executed during the Test phase, ensuring all tests are properly written and failing before implementation begins.

## Pre-Test Setup

### Environment Preparation
- [ ] Development environment configured
- [ ] Test framework installed
- [ ] CI/CD pipeline ready
- [ ] Test databases initialized
- [ ] Mock services configured
- [ ] Test data prepared

### Tool Configuration
```bash
# Install test dependencies
npm install --save-dev jest @types/jest
npm install --save-dev supertest # For API testing
npm install --save-dev puppeteer # For E2E testing

# Configure test runner
npx jest --init
```

## Test Writing Procedures

### Procedure 1: Contract Test Creation

#### Purpose
Define external interface behavior before implementation.

#### Steps
1. **Identify Contract**
   - Review API specification from Design phase
   - Identify endpoint, method, inputs, outputs

2. **Create Test File**
   ```bash
   touch tests/contract/api/users.test.js
   ```

3. **Write Test Structure**
   ```javascript
   describe('POST /api/users', () => {
     beforeEach(() => {
       // Setup
     });

     afterEach(() => {
       // Cleanup
     });

     it('should create user with valid data', async () => {
       // Arrange
       const userData = { /* test data */ };

       // Act
       const response = await request(app)
         .post('/api/users')
         .send(userData);

       // Assert
       expect(response.status).toBe(201);
       expect(response.body).toHaveProperty('id');
     });

     it('should reject invalid email', async () => {
       // Test implementation
     });
   });
   ```

4. **Verify Test Fails**
   ```bash
   npm test tests/contract/api/users.test.js
   # Expected: ReferenceError: app is not defined
   ```

5. **Document Failure**
   - Record expected failure reason
   - Confirm test is properly written

### Procedure 2: Integration Test Creation

#### Purpose
Test component interactions without implementation.

#### Steps
1. **Identify Integration Points**
   - Review component diagram
   - List component dependencies

2. **Create Integration Test**
   ```javascript
   describe('User Service Integration', () => {
     let userService;
     let emailService;
     let database;

     beforeEach(() => {
       // Setup mocked dependencies
       database = createMockDatabase();
       emailService = createMockEmailService();
       userService = new UserService(database, emailService);
     });

     it('should create user and send email', async () => {
       // Test that components work together
       const user = await userService.createUser(userData);

       expect(database.save).toHaveBeenCalled();
       expect(emailService.send).toHaveBeenCalled();
     });
   });
   ```

3. **Run and Verify Failure**
   ```bash
   npm test tests/integration/
   # Expected: UserService is not defined
   ```

### Procedure 3: Unit Test Creation

#### Purpose
Test complex business logic in isolation.

#### Steps
1. **Identify Complex Logic**
   - Review business rules
   - Find algorithms needing testing

2. **Create Focused Unit Test**
   ```javascript
   describe('Password Validator', () => {
     it('should require minimum length', () => {
       expect(validatePassword('short')).toBe(false);
       expect(validatePassword('longenoughpassword')).toBe(true);
     });

     it('should require complexity', () => {
       expect(validatePassword('simple')).toBe(false);
       expect(validatePassword('Complex123!')).toBe(true);
     });
   });
   ```

3. **Verify Failure**
   ```bash
   npm test tests/unit/
   # Expected: validatePassword is not defined
   ```

## Test Execution Procedures

### Daily Test Execution

#### Morning Routine
1. **Pull Latest Code**
   ```bash
   git pull origin main
   ```

2. **Run All Tests**
   ```bash
   npm test
   # All should fail in Test phase
   ```

3. **Generate Coverage Report**
   ```bash
   npm test -- --coverage
   # Should show 0% coverage
   ```

4. **Review Failures**
   - Ensure failures are expected
   - No tests should accidentally pass

### Continuous Test Monitoring

#### Watch Mode Setup
```bash
# Run tests in watch mode during development
npm test -- --watch

# Run specific test suite
npm test -- --watch tests/contract/
```

#### CI/CD Pipeline
```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
      - run: npm ci
      - run: npm test
      - run: npm test -- --coverage
```

## Test Organization Procedures

### File Naming Convention
```
tests/
├── contract/
│   └── [feature].contract.test.js
├── integration/
│   └── [component].integration.test.js
├── unit/
│   └── [module].unit.test.js
├── e2e/
│   └── [journey].e2e.test.js
└── fixtures/
    └── [data].fixture.js
```

### Test Categorization

#### Priority Levels
- **P0**: Must pass before deployment
- **P1**: Should pass for good quality
- **P2**: Nice to have for edge cases

#### Tagging Tests
```javascript
describe('User API [P0]', () => {
  it('[CRITICAL] should authenticate user', () => {});
  it('[NORMAL] should update profile', () => {});
  it('[EDGE] should handle Unicode names', () => {});
});
```

## Test Data Procedures

### Fixture Creation
```javascript
// tests/fixtures/users.fixture.js
export const validUser = {
  email: 'test@example.com',
  password: 'SecurePass123!',
  name: 'Test User'
};

export const invalidUsers = [
  { email: 'invalid', password: 'pass' },
  { email: '', password: 'SecurePass123!' },
  { email: 'test@example.com', password: '' }
];
```

### Test Database Setup
```bash
# Create test database
createdb myapp_test

# Run migrations
npm run migrate:test

# Seed test data
npm run seed:test
```

## Validation Procedures

### Test Quality Checklist
Before marking a test complete:
- [ ] Test has clear name describing behavior
- [ ] Test follows Arrange-Act-Assert pattern
- [ ] Test fails for the right reason
- [ ] Test has appropriate assertions
- [ ] Test is independent of other tests
- [ ] Test uses appropriate test data

### Test Review Process
1. **Self Review**
   - Run test in isolation
   - Verify failure message is clear
   - Check test isn't accidentally passing

2. **Peer Review**
   - Another developer reviews test
   - Confirms test matches requirements
   - Validates test quality

3. **Sign-off**
   - Test added to tracking document
   - Marked as ready for implementation

## Troubleshooting Procedures

### Common Issues

#### Test Passing Unexpectedly
**Problem**: Test passes when it should fail
**Solution**:
- Check if implementation exists
- Verify assertions are correct
- Ensure mocks aren't too permissive

#### Unclear Failure Message
**Problem**: Test fails but reason unclear
**Solution**:
- Add descriptive assertion messages
- Break complex tests into smaller ones
- Use better test names

#### Flaky Tests
**Problem**: Test sometimes passes, sometimes fails
**Solution**:
- Remove time dependencies
- Fix race conditions
- Use proper async/await

## Handoff to Build Phase

### Completion Criteria
- [ ] All required tests written
- [ ] All tests properly failing
- [ ] Test plan coverage targets defined
- [ ] CI/CD pipeline configured
- [ ] Test documentation complete

### Deliverables
1. **Failing Test Suite**: Complete set of red tests
2. **Coverage Report**: Baseline (0%) coverage
3. **Test Inventory**: List of all tests with priorities
4. **Test Data**: Fixtures and seeds ready
5. **Environment**: Test infrastructure operational

---
*These procedures ensure tests are properly written and failing before any implementation begins.*