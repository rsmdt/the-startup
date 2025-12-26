---
name: test-design
description: Apply test pyramid principles, coverage targets, and framework-specific patterns. Use when designing test suites, reviewing test coverage, or implementing tests. Covers Jest, Pytest, and common testing frameworks with naming conventions and organization patterns.
---

# Testing Strategies

A development skill that provides comprehensive testing methodology including test pyramid structure, coverage targets, framework-specific patterns, and test organization conventions.

## When to Use

- Designing test suites for new features or projects
- Reviewing test coverage and identifying gaps
- Implementing unit, integration, or E2E tests
- Establishing test naming and organization conventions
- Selecting appropriate testing frameworks
- Planning test automation strategies

## Test Pyramid

The test pyramid guides test distribution for optimal feedback speed and confidence:

```
        /\
       /  \        E2E Tests (5-10%)
      /----\       - Critical user journeys
     /      \      - Slow, expensive, flaky-prone
    /--------\
   /          \    Integration Tests (20-30%)
  /  Service   \   - API contracts, database
 /--------------\  - Component interactions
/                \
/   Unit Tests    \  Unit Tests (60-70%)
/==================\ - Fast, isolated, deterministic
                     - Business logic focus
```

### Distribution Guidelines

| Test Type   | Target % | Execution Time | Scope                    |
|-------------|----------|----------------|--------------------------|
| Unit        | 60-70%   | < 100ms each   | Single function/class    |
| Integration | 20-30%   | < 5s each      | Service boundaries       |
| E2E         | 5-10%    | < 30s each     | Critical user paths      |

## Core Patterns

### Pattern 1: Test Behavior, Not Implementation

Tests should verify observable behavior, not internal implementation details.

```typescript
// CORRECT: Tests behavior
describe('ShoppingCart', () => {
  it('calculates total with quantity discounts', () => {
    const cart = new ShoppingCart();
    cart.add({ sku: 'WIDGET', price: 10, quantity: 5 });

    expect(cart.total).toBe(45); // 10% discount for 5+ items
  });
});

// INCORRECT: Tests implementation
describe('ShoppingCart', () => {
  it('calls _applyDiscount method', () => {
    const cart = new ShoppingCart();
    const spy = jest.spyOn(cart, '_applyDiscount');

    cart.add({ sku: 'WIDGET', price: 10, quantity: 5 });

    expect(spy).toHaveBeenCalledWith(0.1);
  });
});
```

### Pattern 2: Arrange-Act-Assert Structure

Every test follows the AAA pattern for clarity and consistency.

```python
def test_user_registration_sends_welcome_email():
    # Arrange
    email_service = MockEmailService()
    user_service = UserService(email_service=email_service)
    registration_data = {"email": "new@user.com", "name": "New User"}

    # Act
    user_service.register(registration_data)

    # Assert
    assert email_service.sent_emails == [
        {"to": "new@user.com", "template": "welcome"}
    ]
```

### Pattern 3: One Assertion Per Behavior

Each test verifies one specific behavior. Multiple assertions are acceptable when verifying a single logical outcome.

```typescript
// CORRECT: One behavior, multiple related assertions
it('creates order with correct initial state', () => {
  const order = createOrder(items);

  expect(order.status).toBe('pending');
  expect(order.items).toHaveLength(items.length);
  expect(order.createdAt).toBeInstanceOf(Date);
});

// INCORRECT: Multiple unrelated behaviors
it('creates and processes order', () => {
  const order = createOrder(items);
  expect(order.status).toBe('pending');

  processPayment(order);
  expect(order.status).toBe('paid');  // Different behavior
});
```

### Pattern 4: Descriptive Test Names

Test names describe the scenario and expected outcome.

```typescript
// Format: [unit]_[scenario]_[expected outcome]
// or: [action]_[condition]_[result]

// CORRECT
it('calculateTotal returns zero for empty cart')
it('validateEmail rejects addresses without @ symbol')
it('UserService throws NotFoundError when user does not exist')

// INCORRECT
it('test calculateTotal')
it('validateEmail works')
it('error handling')
```

### Pattern 5: Test Isolation

Tests must be independent and not share mutable state.

```python
# CORRECT: Fresh fixture per test
class TestOrderService:
    def setup_method(self):
        self.db = InMemoryDatabase()
        self.service = OrderService(self.db)

    def test_create_order(self):
        order = self.service.create(items=[...])
        assert order.id is not None

    def test_list_orders_empty(self):
        orders = self.service.list()
        assert orders == []

# INCORRECT: Shared state between tests
db = Database()  # Shared!

def test_create_order():
    order = OrderService(db).create(items=[...])

def test_list_orders_empty():
    orders = OrderService(db).list()  # May see order from previous test!
```

## Coverage Targets

### Recommended Coverage by Code Type

| Code Type          | Statement | Branch | Target |
|--------------------|-----------|--------|--------|
| Business Logic     | 90%       | 85%    | High   |
| API Controllers    | 80%       | 75%    | Medium |
| Utility Functions  | 95%       | 90%    | High   |
| UI Components      | 70%       | 65%    | Medium |
| Generated Code     | N/A       | N/A    | Skip   |

### Coverage Quality Over Quantity

Coverage percentage alone is insufficient. Prioritize:

1. **Critical paths**: Authentication, payments, data integrity
2. **Edge cases**: Boundaries, empty states, error conditions
3. **Regression prevention**: Previously broken functionality
4. **Complex logic**: High cyclomatic complexity areas

## Framework-Specific Patterns

### Jest (JavaScript/TypeScript)

```typescript
// File structure
src/
  services/
    UserService.ts
    UserService.test.ts  // Co-located

// Test setup
describe('UserService', () => {
  let userService: UserService;
  let mockRepo: jest.Mocked<UserRepository>;

  beforeEach(() => {
    mockRepo = {
      findById: jest.fn(),
      save: jest.fn(),
    };
    userService = new UserService(mockRepo);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('getUser', () => {
    it('returns user when found', async () => {
      const expected = { id: '1', name: 'Alice' };
      mockRepo.findById.mockResolvedValue(expected);

      const result = await userService.getUser('1');

      expect(result).toEqual(expected);
      expect(mockRepo.findById).toHaveBeenCalledWith('1');
    });

    it('throws NotFoundError when user does not exist', async () => {
      mockRepo.findById.mockResolvedValue(null);

      await expect(userService.getUser('999'))
        .rejects.toThrow(NotFoundError);
    });
  });
});
```

### Pytest (Python)

```python
# File structure
src/
  services/
    user_service.py
tests/
  services/
    test_user_service.py  # Mirror structure

# conftest.py - Shared fixtures
import pytest
from unittest.mock import Mock

@pytest.fixture
def mock_user_repo():
    return Mock(spec=UserRepository)

@pytest.fixture
def user_service(mock_user_repo):
    return UserService(repository=mock_user_repo)

# test_user_service.py
class TestUserService:
    def test_get_user_returns_user_when_found(
        self, user_service, mock_user_repo
    ):
        expected = User(id="1", name="Alice")
        mock_user_repo.find_by_id.return_value = expected

        result = user_service.get_user("1")

        assert result == expected
        mock_user_repo.find_by_id.assert_called_once_with("1")

    def test_get_user_raises_not_found_when_missing(
        self, user_service, mock_user_repo
    ):
        mock_user_repo.find_by_id.return_value = None

        with pytest.raises(NotFoundError):
            user_service.get_user("999")

    @pytest.mark.parametrize("invalid_email", [
        "",
        "no-at-sign",
        "@no-local-part.com",
        "no-domain@",
    ])
    def test_validate_email_rejects_invalid_formats(
        self, user_service, invalid_email
    ):
        assert not user_service.validate_email(invalid_email)
```

### React Testing Library

```typescript
// Component test
import { render, screen, userEvent } from '@testing-library/react';
import { LoginForm } from './LoginForm';

describe('LoginForm', () => {
  it('calls onSubmit with credentials when form is submitted', async () => {
    const onSubmit = jest.fn();
    const user = userEvent.setup();

    render(<LoginForm onSubmit={onSubmit} />);

    await user.type(screen.getByLabelText('Email'), 'test@example.com');
    await user.type(screen.getByLabelText('Password'), 'secret123');
    await user.click(screen.getByRole('button', { name: 'Sign In' }));

    expect(onSubmit).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'secret123',
    });
  });

  it('displays validation error for invalid email', async () => {
    const user = userEvent.setup();

    render(<LoginForm onSubmit={jest.fn()} />);

    await user.type(screen.getByLabelText('Email'), 'invalid');
    await user.click(screen.getByRole('button', { name: 'Sign In' }));

    expect(screen.getByText('Please enter a valid email')).toBeInTheDocument();
  });
});
```

## Test Organization

### File Naming Conventions

| Framework | Test File Pattern      | Example                    |
|-----------|------------------------|----------------------------|
| Jest      | `*.test.ts`            | `UserService.test.ts`      |
| Pytest    | `test_*.py`            | `test_user_service.py`     |
| Go        | `*_test.go`            | `user_service_test.go`     |
| JUnit     | `*Test.java`           | `UserServiceTest.java`     |

### Directory Structure Patterns

**Co-located tests** (preferred for unit tests):
```
src/
  services/
    UserService.ts
    UserService.test.ts
  utils/
    validators.ts
    validators.test.ts
```

**Separate test directory** (for integration/E2E):
```
src/
  services/
    UserService.ts
tests/
  unit/
    services/
      UserService.test.ts
  integration/
    api/
      users.test.ts
  e2e/
    user-registration.spec.ts
```

## Best Practices

- Run tests before committing; never commit failing tests
- Keep unit tests under 100ms execution time
- Mock external dependencies at service boundaries only
- Use factories or fixtures for test data, not raw literals
- Delete flaky tests or fix them immediately
- Review tests during code review with same rigor as production code
- Name tests as specifications that document behavior
- Prefer real implementations over mocks when practical
- Test edge cases: nulls, empty collections, boundaries
- Avoid conditional logic in tests

## Anti-Patterns to Avoid

### Testing Implementation Details

```typescript
// WRONG: Brittle test that breaks on refactoring
it('stores user in _users array', () => {
  const service = new UserService();
  service.addUser(user);
  expect(service._users).toContain(user);
});
```

### Shared Mutable State

```python
# WRONG: Tests interfere with each other
users = []

def test_add_user():
    users.append(User())

def test_user_count():
    assert len(users) == 0  # Fails if test_add_user runs first
```

### Over-Mocking

```typescript
// WRONG: Mocking the system under test
it('processes payment', () => {
  const processor = new PaymentProcessor();
  jest.spyOn(processor, 'validate').mockReturnValue(true);
  jest.spyOn(processor, 'charge').mockResolvedValue({ success: true });

  // What are we even testing?
  const result = processor.process(payment);
});
```

### Test Duplication

```python
# WRONG: Copy-paste tests with minor variations
def test_validate_email_empty():
    assert not validate_email("")

def test_validate_email_no_at():
    assert not validate_email("invalid")

def test_validate_email_no_domain():
    assert not validate_email("user@")

# CORRECT: Parameterized test
@pytest.mark.parametrize("invalid_email", ["", "invalid", "user@"])
def test_validate_email_rejects_invalid_formats(invalid_email):
    assert not validate_email(invalid_email)
```

## References

- `examples/test-pyramid.md` - Detailed test pyramid implementation guide with framework examples
