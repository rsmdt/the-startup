---
name: test-strategy-design
description: Test pyramid strategy, test design patterns, coverage analysis, and quality gate configuration. Use when designing test strategies, improving coverage, setting up automation, or defining quality gates. Covers unit, integration, and E2E testing across frameworks.
---

# Test Strategy Design

Systematic patterns for designing comprehensive test strategies that balance coverage, speed, and maintainability.

## When to Activate

- Designing test strategy for new projects
- Improving existing test coverage
- Setting up test automation frameworks
- Configuring CI/CD quality gates
- Analyzing test gaps and priorities
- Establishing testing standards for teams

## The Test Pyramid

The test pyramid guides the distribution of testing effort:

```
                 ┌─────────┐
                 │   E2E   │  ~10% - Critical user journeys
                 │ Tests   │  Slow, expensive, high confidence
                 ├─────────┤
                 │ Integr- │  ~20% - Component interactions
                 │  ation  │  Medium speed, real dependencies
                 ├─────────┤
                 │  Unit   │  ~70% - Business logic
                 │  Tests  │  Fast, isolated, cheap
                 └─────────┘
```

### Layer Characteristics

| Layer | Speed | Isolation | Confidence | Maintenance |
|-------|-------|-----------|------------|-------------|
| Unit | < 10ms | Complete | Low per test | Low |
| Integration | < 1s | Partial | Medium | Medium |
| E2E | < 30s | None | High | High |

### When to Test at Each Layer

**Unit Tests** - Test here when:
- Pure business logic
- Algorithmic code
- Data transformations
- Validation rules
- Edge cases and boundaries

**Integration Tests** - Test here when:
- Database interactions
- API contracts
- Service communication
- File system operations
- Cache behavior

**E2E Tests** - Test here when:
- Critical user journeys
- Cross-system workflows
- Authentication flows
- Payment processes
- Smoke tests for deployments

## Test Design Patterns

### Unit Test Patterns

#### Arrange-Act-Assert (AAA)

```
// Arrange: Set up test data and dependencies
const user = createUser({ email: 'test@example.com' });
const validator = new EmailValidator();

// Act: Execute the behavior under test
const result = validator.validate(user.email);

// Assert: Verify the expected outcome
expect(result.isValid).toBe(true);
```

#### One Assertion Per Test

Focus each test on a single behavior:

```
// Bad: Multiple unrelated assertions
test('user validation', () => {
  expect(user.isValid()).toBe(true);
  expect(user.email).toContain('@');
  expect(user.createdAt).toBeDefined();
});

// Good: Separate tests for each behavior
test('validates user with proper email', () => {
  expect(user.isValid()).toBe(true);
});

test('email contains @ symbol', () => {
  expect(user.email).toContain('@');
});
```

#### Descriptive Test Names

Test names should describe the expected behavior:

```
// Bad
test('email test', () => { ... });

// Good
test('rejects email without domain', () => { ... });
test('accepts email with subdomain', () => { ... });
test('trims whitespace from email input', () => { ... });
```

#### Test Data Builders

Create readable test data:

```
// Builder pattern for test data
const user = UserBuilder.create()
  .withEmail('test@example.com')
  .withRole('admin')
  .withActiveSubscription()
  .build();
```

### Integration Test Patterns

#### Test Through Public Interfaces

```
// Bad: Testing internal implementation
test('repository uses correct SQL', () => {
  const sql = userRepo.buildQuery();
  expect(sql).toContain('SELECT * FROM users');
});

// Good: Testing behavior through public API
test('finds users by email domain', async () => {
  await userRepo.save(createUser({ email: 'a@corp.com' }));
  await userRepo.save(createUser({ email: 'b@corp.com' }));

  const users = await userRepo.findByDomain('corp.com');

  expect(users).toHaveLength(2);
});
```

#### Use Real Dependencies Where Practical

```
// Prefer real database for repository tests
beforeAll(async () => {
  testDb = await createTestDatabase();
});

afterEach(async () => {
  await testDb.truncateAll();
});

// Mock only external services
jest.mock('./paymentGateway', () => ({
  charge: jest.fn().mockResolvedValue({ success: true })
}));
```

#### Contract Testing

Verify API contracts between services:

```
// Consumer contract
describe('User API Contract', () => {
  test('GET /users/:id returns user shape', async () => {
    const response = await api.get('/users/123');

    expect(response).toMatchSchema({
      id: expect.any(String),
      email: expect.any(String),
      createdAt: expect.any(String),
    });
  });
});
```

### E2E Test Patterns

#### Critical Path Focus

Only test the most important user journeys:

```
// Priority 1: Revenue-critical paths
test('complete purchase flow', async () => {
  await page.addToCart(product);
  await page.checkout();
  await page.enterPayment(validCard);
  await page.confirmOrder();

  await expect(page.orderConfirmation).toBeVisible();
});

// Priority 2: Core functionality
test('user registration and login', async () => {
  await page.register(newUser);
  await page.login(newUser);

  await expect(page.dashboard).toBeVisible();
});
```

#### Page Object Pattern

Encapsulate page interactions:

```
class CheckoutPage {
  constructor(page) {
    this.page = page;
  }

  async enterShippingAddress(address) {
    await this.page.fill('#street', address.street);
    await this.page.fill('#city', address.city);
    await this.page.fill('#zip', address.zip);
  }

  async submitOrder() {
    await this.page.click('[data-testid="submit-order"]');
    await this.page.waitForURL('**/confirmation');
  }
}
```

#### Resilient Selectors

Use stable selectors that survive UI changes:

```
// Fragile: Based on CSS/position
await page.click('.btn-primary');
await page.click('div > div:nth-child(2) > button');

// Stable: Based on test IDs or roles
await page.click('[data-testid="submit-button"]');
await page.getByRole('button', { name: 'Submit' });
```

## Coverage Strategy

### Coverage Types

| Type | Measures | Target |
|------|----------|--------|
| Line | Lines executed | 80% |
| Branch | Decision paths taken | 75% |
| Function | Functions called | 90% |
| Statement | Statements executed | 80% |

### Coverage Prioritization

Focus coverage on high-risk areas:

```
High Priority (100% coverage):
├── Payment processing
├── Authentication/Authorization
├── Data validation
├── Business rule enforcement
└── Security-sensitive code

Medium Priority (80% coverage):
├── Core business logic
├── API endpoints
├── Data transformations
└── Error handling

Lower Priority (60% coverage):
├── UI components
├── Configuration
├── Logging/Monitoring
└── Admin features
```

### Meaningful Coverage

Coverage quantity vs quality:

```
// 100% coverage, low value
test('getter returns value', () => {
  user.name = 'John';
  expect(user.name).toBe('John');
});

// 80% coverage, high value
test('rejects order when inventory insufficient', () => {
  const inventory = createInventory({ quantity: 5 });
  const order = createOrder({ quantity: 10 });

  expect(() => processOrder(order, inventory))
    .toThrow('Insufficient inventory');
});
```

## Quality Gates

### Gate Configuration

| Gate | Threshold | Enforcement |
|------|-----------|-------------|
| Unit Tests | 100% pass | Block merge |
| Integration Tests | 100% pass | Block merge |
| E2E Tests | 100% pass | Block deploy |
| Line Coverage | ≥ 80% | Block merge |
| Branch Coverage | ≥ 75% | Block merge |
| New Code Coverage | ≥ 90% | Warn |
| Test Duration | < 10 min | Warn |

### Flaky Test Management

```
Flaky Test Protocol:
1. DETECT: Track test stability over time
2. QUARANTINE: Move flaky tests to separate suite
3. FIX: Prioritize fixing within 1 week
4. REINTEGRATE: Return to main suite after 10 consecutive passes
```

### Test Performance Budgets

| Test Type | Target Duration | Action if Exceeded |
|-----------|-----------------|-------------------|
| Single Unit Test | < 10ms | Investigate |
| Unit Suite | < 30s | Parallelize |
| Single Integration | < 1s | Review setup/teardown |
| Integration Suite | < 2 min | Parallelize |
| Single E2E | < 30s | Optimize waits |
| E2E Suite | < 10 min | Reduce scope |

## Test Organization

### File Structure

```
src/
├── user/
│   ├── User.ts
│   ├── User.test.ts           # Unit tests co-located
│   ├── UserRepository.ts
│   └── UserRepository.test.ts
├── order/
│   ├── Order.ts
│   └── Order.test.ts
tests/
├── integration/
│   ├── api/
│   │   └── users.test.ts      # API integration tests
│   └── database/
│       └── repositories.test.ts
└── e2e/
    ├── checkout.spec.ts       # E2E tests by journey
    └── registration.spec.ts
```

### Test Naming Conventions

```
# Unit tests: [unit].test.ts
User.test.ts
OrderCalculator.test.ts

# Integration tests: [feature].integration.test.ts
users.integration.test.ts
payment-gateway.integration.test.ts

# E2E tests: [journey].e2e.ts or [journey].spec.ts
checkout.e2e.ts
user-onboarding.spec.ts
```

## Test Data Management

### Test Data Principles

1. **Isolation**: Each test manages its own data
2. **Cleanup**: Tests clean up after themselves
3. **Determinism**: Same test, same data, same result
4. **Minimal**: Only create data needed for the test

### Factories and Fixtures

```
// Factory: Dynamic test data
const createUser = (overrides = {}) => ({
  id: faker.string.uuid(),
  email: faker.internet.email(),
  name: faker.person.fullName(),
  ...overrides
});

// Fixture: Static reference data
const ADMIN_USER = {
  id: 'admin-001',
  email: 'admin@test.com',
  role: 'admin'
};
```

### Database Test Patterns

```
// Transaction rollback pattern
beforeEach(async () => {
  transaction = await db.beginTransaction();
});

afterEach(async () => {
  await transaction.rollback();
});

// Truncate pattern (slower but more isolated)
afterEach(async () => {
  await db.truncateAll();
});
```

## CI/CD Integration

### Pipeline Configuration

```yaml
test:
  stages:
    - lint-and-typecheck    # Fast feedback first
    - unit-tests            # Parallel execution
    - integration-tests     # With test database
    - e2e-tests             # Against staging environment

  parallelization:
    unit: 4 workers
    integration: 2 workers
    e2e: 1 worker (sequential for stability)

  artifacts:
    - coverage reports
    - test results (JUnit XML)
    - E2E screenshots/videos on failure
```

### Failure Handling

```
On Test Failure:
1. Capture: Screenshots, logs, database state
2. Report: Notify team channel
3. Block: Prevent merge/deploy
4. Retry: Once for E2E (flaky detection)
```

## Best Practices

### Do

- Test behavior, not implementation
- Keep tests independent and isolated
- Use meaningful, descriptive test names
- Maintain test code like production code
- Run tests frequently during development
- Fix flaky tests immediately

### Avoid

- Testing private methods directly
- Mocking everything (use real dependencies where practical)
- Ignoring slow tests (fix the root cause)
- Writing tests after bugs (write them before)
- Sharing state between tests
- Over-relying on E2E tests

## References

- [Test Pattern Examples](examples/test-patterns.md) - Code examples for each pattern
- [Framework-Specific Guides](reference.md) - Jest, Pytest, Go testing patterns
