# Test Suite Structure

**Project**: [Project Name]
**Generated**: [Date]
**Coverage Target**: [Percentage]%
**Test Framework**: [Framework]

## Test Organization

```
tests/
├── contract/           # API endpoint tests
├── integration/        # Component interaction tests
├── unit/              # Business logic tests
├── e2e/               # End-to-end user journeys
├── fixtures/          # Test data
├── factories/         # Data generators
├── mocks/             # Service mocks
├── helpers/           # Test utilities
└── setup/             # Test configuration
```

## Contract Tests

### API Endpoints Coverage

| Endpoint | Method | Test File | Status |
|----------|--------|-----------|--------|
| /api/[resource] | GET | [resource].get.test.ts | ❌ Red |
| /api/[resource] | POST | [resource].post.test.ts | ❌ Red |
| /api/[resource]/:id | GET | [resource].get-by-id.test.ts | ❌ Red |
| /api/[resource]/:id | PUT | [resource].update.test.ts | ❌ Red |
| /api/[resource]/:id | DELETE | [resource].delete.test.ts | ❌ Red |

### Contract Test Template

```typescript
// tests/contract/[resource].[method].test.ts
import { request } from '../helpers/test-client';
import { fixtures } from '../fixtures/[resource]';

describe('[METHOD] /api/[resource]', () => {
  describe('Success Cases', () => {
    it('should [expected behavior] with valid data', async () => {
      // Arrange
      const validData = fixtures.valid[ResourceName]();

      // Act - This will fail (no implementation)
      const response = await request
        .[method]('/api/[resource]')
        .send(validData);

      // Assert
      expect(response.status).toBe([expectedStatus]);
      expect(response.body).toMatchObject({
        // Expected response structure
      });
    });
  });

  describe('Validation Cases', () => {
    it('should return 400 when [field] is missing', async () => {
      const invalidData = fixtures.missing[Field]();

      const response = await request
        .[method]('/api/[resource]')
        .send(invalidData);

      expect(response.status).toBe(400);
      expect(response.body.error).toContain('[field]');
    });

    it('should return 400 when [field] is invalid', async () => {
      const invalidData = fixtures.invalid[Field]();

      const response = await request
        .[method]('/api/[resource]')
        .send(invalidData);

      expect(response.status).toBe(400);
      expect(response.body.error).toContain('[validation message]');
    });
  });

  describe('Authorization Cases', () => {
    it('should return 401 when not authenticated', async () => {
      const response = await request
        .[method]('/api/[resource]')
        .send(fixtures.validData());

      expect(response.status).toBe(401);
    });

    it('should return 403 when user lacks permission', async () => {
      const response = await request
        .[method]('/api/[resource]')
        .auth(fixtures.userWithoutPermission())
        .send(fixtures.validData());

      expect(response.status).toBe(403);
    });
  });

  describe('Error Cases', () => {
    it('should return 404 when resource not found', async () => {
      const response = await request
        .get('/api/[resource]/non-existent-id');

      expect(response.status).toBe(404);
    });

    it('should return 409 when conflict occurs', async () => {
      const conflictingData = fixtures.conflicting[Resource]();

      const response = await request
        .[method]('/api/[resource]')
        .send(conflictingData);

      expect(response.status).toBe(409);
    });
  });
});
```

## Integration Tests

### Component Interactions

| Component A | Component B | Test Scenario | Test File | Status |
|------------|-------------|---------------|-----------|--------|
| [Service] | [Repository] | Data persistence | [service]-persistence.test.ts | ❌ Red |
| [Service] | [External API] | API integration | [service]-api.test.ts | ❌ Red |
| [Controller] | [Service] | Request handling | [controller]-flow.test.ts | ❌ Red |

### Integration Test Template

```typescript
// tests/integration/[feature].integration.test.ts
import { [Service] } from '@/services/[service]';
import { [Repository] } from '@/repositories/[repository]';
import { mock[External] } from '../mocks/[external]';

describe('[Feature] Integration', () => {
  let service: [Service];
  let repository: [Repository];

  beforeEach(() => {
    // Setup test dependencies
    repository = new [Repository](testDatabase);
    service = new [Service](repository, mock[External]);
  });

  describe('[Operation] Flow', () => {
    it('should coordinate between components correctly', async () => {
      // Arrange
      const input = fixtures.valid[Input]();
      mock[External].setup({
        response: fixtures.external[Response]()
      });

      // Act - Will fail (no implementation)
      const result = await service.[operation](input);

      // Assert
      expect(result).toMatchObject({
        // Expected result structure
      });

      // Verify component interactions
      expect(repository.[method]).toHaveBeenCalledWith(
        expect.objectContaining({
          // Expected repository call
        })
      );

      expect(mock[External].[method]).toHaveBeenCalledWith(
        // Expected external service call
      );
    });

    it('should handle component failures gracefully', async () => {
      // Setup failure scenario
      mock[External].setupFailure(new Error('Service unavailable'));

      // Act & Assert
      await expect(service.[operation](input))
        .rejects
        .toThrow('[Expected error]');

      // Verify rollback/compensation
      expect(repository.[method]).not.toHaveBeenCalled();
    });
  });

  describe('Transaction Handling', () => {
    it('should rollback on failure', async () => {
      // Test transaction rollback behavior
    });

    it('should commit on success', async () => {
      // Test transaction commit behavior
    });
  });
});
```

## Unit Tests

### Business Logic Coverage

| Module | Function | Test Coverage | Test File | Status |
|--------|----------|---------------|-----------|--------|
| [validators] | validate[Entity] | All rules | validators.test.ts | ❌ Red |
| [calculators] | calculate[Metric] | All formulas | calculators.test.ts | ❌ Red |
| [transformers] | transform[Data] | All mappings | transformers.test.ts | ❌ Red |

### Unit Test Template

```typescript
// tests/unit/[module].unit.test.ts
import { [functionName] } from '@/[module]/[file]';

describe('[functionName]', () => {
  describe('Happy Path', () => {
    it('should [expected behavior] with valid input', () => {
      // Arrange
      const input = [validInput];
      const expected = [expectedOutput];

      // Act - Will fail (no implementation)
      const result = [functionName](input);

      // Assert
      expect(result).toEqual(expected);
    });
  });

  describe('Edge Cases', () => {
    it('should handle empty input', () => {
      expect([functionName]([])).toEqual([expectedEmpty]);
    });

    it('should handle null values', () => {
      expect([functionName](null)).toEqual([expectedNull]);
    });

    it('should handle maximum values', () => {
      const maxInput = [maximumValidInput];
      expect([functionName](maxInput)).toEqual([expectedMax]);
    });
  });

  describe('Error Cases', () => {
    it('should throw error for invalid input type', () => {
      expect(() => [functionName]('invalid'))
        .toThrow('[Expected error message]');
    });

    it('should throw error for out of range values', () => {
      expect(() => [functionName](-1))
        .toThrow('[Range error message]');
    });
  });

  describe('Business Rules', () => {
    it('should apply [business rule 1]', () => {
      // Test specific business logic
    });

    it('should enforce [constraint 1]', () => {
      // Test constraint enforcement
    });
  });
});
```

## End-to-End Tests

### User Journey Coverage

| Journey | Steps | Critical Path | Test File | Status |
|---------|-------|---------------|-----------|--------|
| [User Registration] | 5 | Yes | registration.e2e.test.ts | ❌ Red |
| [Purchase Flow] | 8 | Yes | purchase.e2e.test.ts | ❌ Red |
| [Data Export] | 3 | No | export.e2e.test.ts | ❌ Red |

### E2E Test Template

```typescript
// tests/e2e/[journey].e2e.test.ts
import { test, expect } from '@playwright/test';
import { fixtures } from '../fixtures/e2e';

describe('[User Journey Name]', () => {
  test('should complete [journey] successfully', async ({ page }) => {
    // Step 1: [Initial action]
    await page.goto('/[starting-point]');
    await expect(page).toHaveTitle('[Expected title]');

    // Step 2: [User action]
    await page.fill('[selector]', fixtures.user.email);
    await page.fill('[selector]', fixtures.user.password);
    await page.click('[submit-button]');

    // Step 3: [Verification]
    await expect(page).toHaveURL('/[expected-route]');
    await expect(page.locator('[selector]')).toContainText('[expected]');

    // Step 4: [Core action]
    await page.click('[action-button]');
    await page.waitForResponse('**/api/[endpoint]');

    // Step 5: [Final verification]
    await expect(page.locator('[success-indicator]')).toBeVisible();
    await expect(page.locator('[result]')).toContainText('[expected-result]');
  });

  test('should handle errors during [journey]', async ({ page }) => {
    // Test error scenarios in user journey
  });
});
```

## Test Data Management

### Fixtures Structure

```typescript
// tests/fixtures/[entity].ts
export const [entity]Fixtures = {
  valid: {
    minimal: () => ({
      // Minimum required fields
    }),
    complete: () => ({
      // All fields populated
    }),
    withRelations: () => ({
      // Including related entities
    })
  },
  invalid: {
    missingRequired: () => ({
      // Missing required field
    }),
    invalidFormat: () => ({
      // Invalid field format
    }),
    outOfRange: () => ({
      // Values outside valid range
    })
  },
  edge: {
    empty: () => ({}),
    maxLength: () => ({
      // Maximum length values
    }),
    special: () => ({
      // Special characters
    })
  }
};
```

### Factory Pattern

```typescript
// tests/factories/[entity].factory.ts
export class [Entity]Factory {
  static build(overrides = {}) {
    return {
      id: faker.datatype.uuid(),
      name: faker.name.fullName(),
      email: faker.internet.email(),
      createdAt: faker.date.recent(),
      ...overrides
    };
  }

  static buildList(count: number, overrides = {}) {
    return Array.from({ length: count }, () =>
      this.build(overrides)
    );
  }

  static buildWithRelations(overrides = {}) {
    return {
      ...this.build(overrides),
      [relation]: [RelatedEntity]Factory.buildList(3)
    };
  }
}
```

## Coverage Report Template

### Current Coverage Status

| Test Type | Files | Lines | Branches | Functions | Status |
|-----------|-------|-------|----------|-----------|--------|
| Contract | 0/[total] | 0% | 0% | 0% | ❌ Red |
| Integration | 0/[total] | 0% | 0% | 0% | ❌ Red |
| Unit | 0/[total] | 0% | 0% | 0% | ❌ Red |
| E2E | 0/[total] | N/A | N/A | N/A | ❌ Red |

### Coverage Targets

| Metric | Target | Current | Gap |
|--------|--------|---------|-----|
| Overall Line Coverage | 80% | 0% | 80% |
| Contract Test Coverage | 100% | 0% | 100% |
| Critical Path Coverage | 100% | 0% | 100% |
| Error Handling Coverage | 90% | 0% | 90% |

## Test Execution Plan

### Local Development
```bash
# Run all tests (will fail - no implementation)
npm test

# Run specific test type
npm run test:contract
npm run test:integration
npm run test:unit
npm run test:e2e

# Run with coverage
npm run test:coverage

# Watch mode for TDD
npm run test:watch
```

### CI Pipeline Configuration
```yaml
test-suite:
  stage: test
  parallel:
    matrix:
      - TEST_TYPE: unit
      - TEST_TYPE: integration
      - TEST_TYPE: contract
      - TEST_TYPE: e2e
  script:
    - npm run test:${TEST_TYPE}
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage/cobertura-coverage.xml
```

## Definition of Done

- [ ] All API endpoints have contract tests (failing)
- [ ] All components have integration tests (failing)
- [ ] All business logic has unit tests (failing)
- [ ] All user journeys have E2E tests (failing)
- [ ] Test data fixtures are prepared
- [ ] Test factories are implemented
- [ ] Mocks for external services are ready
- [ ] Tests are organized by type
- [ ] CI pipeline is configured
- [ ] Coverage targets are defined
- [ ] All tests are in RED state (failing)

## Next Phase: Build

With all tests defined and failing, the Build phase can begin with:
1. Clear definition of "done" (tests passing)
2. Incremental implementation to satisfy tests
3. Red-Green-Refactor cycle
4. No production code without a failing test