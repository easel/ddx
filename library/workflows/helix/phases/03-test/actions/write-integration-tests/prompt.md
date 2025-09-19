# Write Integration Tests Prompt

Create integration tests that verify component interactions, service boundaries, and data flow between different parts of the system. Following Specification-Driven Development principles, use real services and dependencies wherever possible.

## Test Output Location

Generate integration tests in: `tests/integration/`

Organize tests by integration type:
- `tests/integration/api/` - API endpoint tests
- `tests/integration/database/` - Database operation tests
- `tests/integration/services/` - Service integration tests
- `tests/integration/external/` - External service tests

## Purpose

Integration tests verify that:
- Components work together correctly
- Services integrate properly with databases
- APIs meet their contracts
- Data flows correctly through the system
- External dependencies are properly integrated

## Key Principles (from SDD Article IX)

### Integration-First Testing
- **Use real services over mocks** whenever possible
- **Test in production-like environments**
- **Verify actual integration points**, not simulations
- **Maintain test databases** that mirror production structure
- **Use containerization** for consistent test environments

## Test Scope

### Service Integration
Test interactions between:
- Application services and databases
- Microservices communication
- Message queue integrations
- Cache layer interactions
- External API integrations

### Data Flow Testing
Verify:
- Data persistence and retrieval
- Transaction boundaries
- Data transformation pipelines
- Event propagation
- State synchronization

### API Contract Testing
Ensure:
- Request/response formats match specifications
- Error responses follow standards
- Authentication/authorization works correctly
- Rate limiting functions properly
- Versioning compatibility

## Test Structure

### Setup and Teardown
```javascript
describe('UserService Integration', () => {
  let database;
  let service;

  beforeAll(async () => {
    // Setup real database connection
    database = await createTestDatabase();
    await database.migrate();
    service = new UserService(database);
  });

  afterAll(async () => {
    // Clean up resources
    await database.close();
  });

  beforeEach(async () => {
    // Reset data between tests
    await database.truncate(['users', 'sessions']);
  });

  it('should create and retrieve user with sessions', async () => {
    // Test with real database operations
    const user = await service.createUser({ name: 'Test User' });
    const session = await service.createSession(user.id);

    const retrieved = await service.getUserWithSessions(user.id);
    expect(retrieved.sessions).toHaveLength(1);
    expect(retrieved.sessions[0].id).toBe(session.id);
  });
});
```

## Database Integration Tests

### Test Database Management
```javascript
// Use real database, not mocks
const testDb = {
  host: process.env.TEST_DB_HOST || 'localhost',
  port: process.env.TEST_DB_PORT || 5432,
  database: `test_${process.env.JEST_WORKER_ID}`, // Parallel test support
  user: process.env.TEST_DB_USER,
  password: process.env.TEST_DB_PASSWORD
};

// Ensure clean state
async function setupTestDatabase() {
  await db.migrate.latest();
  await db.seed.run(); // Minimal seed data
}

async function cleanupTestDatabase() {
  await db.raw('TRUNCATE TABLE users, orders, products RESTART IDENTITY CASCADE');
}
```

### Transaction Testing
```javascript
it('should rollback on error', async () => {
  const trx = await db.transaction();

  try {
    await orderService.createOrder({ /* valid data */ }, trx);
    await paymentService.processPayment({ /* invalid data */ }, trx);
    await trx.commit();
  } catch (error) {
    await trx.rollback();
  }

  // Verify nothing was persisted
  const orders = await db('orders').select();
  expect(orders).toHaveLength(0);
});
```

## API Integration Tests

### REST API Testing
```javascript
describe('POST /api/users', () => {
  it('should create user and return 201', async () => {
    const response = await request(app)
      .post('/api/users')
      .send({
        name: 'John Doe',
        email: 'john@example.com'
      })
      .expect(201);

    expect(response.body).toMatchObject({
      id: expect.any(Number),
      name: 'John Doe',
      email: 'john@example.com'
    });

    // Verify in database
    const user = await db('users').where({ email: 'john@example.com' }).first();
    expect(user).toBeDefined();
  });
});
```

### GraphQL Testing
```javascript
it('should resolve nested user data', async () => {
  const query = `
    query GetUser($id: ID!) {
      user(id: $id) {
        id
        name
        posts {
          id
          title
          comments {
            id
            content
          }
        }
      }
    }
  `;

  const response = await graphqlRequest(query, { id: 1 });
  expect(response.data.user.posts).toBeDefined();
  expect(response.data.user.posts[0].comments).toBeDefined();
});
```

## Message Queue Integration

### Event Publishing Tests
```javascript
it('should publish order event to queue', async () => {
  const queue = new MessageQueue(rabbitMqConnection);
  const orderService = new OrderService(database, queue);

  const order = await orderService.createOrder({ /* data */ });

  // Verify message was published
  const message = await queue.consume('orders.created', { timeout: 5000 });
  expect(message.orderId).toBe(order.id);
  expect(message.timestamp).toBeDefined();
});
```

## External Service Integration

### Using Real Services When Possible
```javascript
describe('Payment Gateway Integration', () => {
  // Use sandbox/test environment of real service
  const paymentGateway = new StripeGateway({
    apiKey: process.env.STRIPE_TEST_KEY,
    webhookSecret: process.env.STRIPE_TEST_WEBHOOK_SECRET
  });

  it('should process payment with real gateway', async () => {
    const payment = await paymentGateway.charge({
      amount: 1000,
      currency: 'usd',
      source: 'tok_visa' // Test token
    });

    expect(payment.status).toBe('succeeded');
    expect(payment.id).toMatch(/^ch_test_/);
  });
});
```

### When Mocking is Necessary
```javascript
// Only mock when real service is not feasible
const mockEmailService = {
  send: jest.fn().mockResolvedValue({ messageId: 'test-123' })
};

// But verify the integration point
it('should call email service with correct parameters', async () => {
  await userService.sendWelcomeEmail(user);

  expect(mockEmailService.send).toHaveBeenCalledWith({
    to: user.email,
    template: 'welcome',
    data: { name: user.name }
  });
});
```

## Test Data Management

### Data Builders for Complex Scenarios
```javascript
class TestDataBuilder {
  async createUserWithOrders(orderCount = 3) {
    const user = await db('users').insert({ /* user data */ }).returning('*');

    const orders = [];
    for (let i = 0; i < orderCount; i++) {
      const order = await db('orders').insert({
        user_id: user.id,
        /* order data */
      }).returning('*');
      orders.push(order);
    }

    return { user, orders };
  }
}
```

## Environment Configuration

### Docker Compose for Test Environment
```yaml
# docker-compose.test.yml
version: '3.8'
services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: testdb
      POSTGRES_USER: testuser
      POSTGRES_PASSWORD: testpass
    ports:
      - "5433:5432"

  redis:
    image: redis:7
    ports:
      - "6380:6379"

  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5673:5672"
      - "15673:15672"
```

## Quality Checklist

Before integration tests are complete:
- [ ] All service boundaries have integration tests
- [ ] Database operations are tested with real database
- [ ] API contracts are validated
- [ ] Message queue integrations are verified
- [ ] External service integrations are tested (sandbox/mocked)
- [ ] Transaction boundaries are tested
- [ ] Error propagation is verified
- [ ] Performance is acceptable (< 30s for suite)
- [ ] Tests can run in parallel
- [ ] Test data is properly isolated

## Common Pitfalls to Avoid

### ❌ Over-Mocking
- Don't mock the database - use a real test database
- Don't mock internal services - test real interactions
- Don't mock message queues - use real queues in Docker

### ✅ Better Approach
- Use containerized services for testing
- Maintain separate test databases per developer/CI
- Use transaction rollback for test isolation
- Leverage database migrations for schema consistency

Remember: Integration tests build confidence that your system works as a whole. They should test real integrations, not mocked approximations.