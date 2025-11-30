# Test Pyramid Implementation Guide

## Context

This guide provides practical examples for implementing tests at each level of the test pyramid. Use these patterns when designing test suites to achieve optimal coverage with fast feedback loops.

## Test Pyramid Overview

```
        /\
       /  \        E2E: 5-10%
      /    \       Validate critical user journeys
     /------\      across the entire system
    /        \
   /          \    Integration: 20-30%
  /            \   Verify component interactions
 /--------------\  and external dependencies
/                \
/   Unit: 60-70%  \ Test isolated business logic
/==================\ Fast, deterministic, focused
```

### Why This Distribution?

| Level       | Speed     | Stability | Confidence | Cost    |
|-------------|-----------|-----------|------------|---------|
| Unit        | Very Fast | Very High | Low-Medium | Low     |
| Integration | Medium    | Medium    | Medium     | Medium  |
| E2E         | Slow      | Low       | High       | High    |

## Unit Tests

### Purpose

- Test single functions, methods, or classes in isolation
- Verify business logic correctness
- Provide fast feedback during development
- Enable refactoring with confidence

### Characteristics

- Execute in under 100ms
- No external dependencies (database, network, filesystem)
- Deterministic (same input always produces same output)
- Can run in parallel without interference

### Jest Example: Order Total Calculation

```typescript
// src/domain/Order.ts
export class Order {
  private items: OrderItem[] = [];

  addItem(item: OrderItem): void {
    this.items.push(item);
  }

  get subtotal(): number {
    return this.items.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0
    );
  }

  get discount(): number {
    if (this.subtotal >= 100) return this.subtotal * 0.1;
    if (this.subtotal >= 50) return this.subtotal * 0.05;
    return 0;
  }

  get total(): number {
    return this.subtotal - this.discount;
  }
}

// src/domain/Order.test.ts
describe('Order', () => {
  describe('subtotal', () => {
    it('returns zero for empty order', () => {
      const order = new Order();

      expect(order.subtotal).toBe(0);
    });

    it('sums item prices multiplied by quantities', () => {
      const order = new Order();
      order.addItem({ sku: 'A', price: 10, quantity: 2 });
      order.addItem({ sku: 'B', price: 25, quantity: 1 });

      expect(order.subtotal).toBe(45);
    });
  });

  describe('discount', () => {
    it('returns zero when subtotal is under 50', () => {
      const order = new Order();
      order.addItem({ sku: 'A', price: 10, quantity: 4 });

      expect(order.discount).toBe(0);
    });

    it('applies 5% discount when subtotal is 50 or more', () => {
      const order = new Order();
      order.addItem({ sku: 'A', price: 50, quantity: 1 });

      expect(order.discount).toBe(2.5);
    });

    it('applies 10% discount when subtotal is 100 or more', () => {
      const order = new Order();
      order.addItem({ sku: 'A', price: 100, quantity: 1 });

      expect(order.discount).toBe(10);
    });
  });

  describe('total', () => {
    it('returns subtotal minus discount', () => {
      const order = new Order();
      order.addItem({ sku: 'A', price: 50, quantity: 2 });

      expect(order.total).toBe(90); // 100 - 10% discount
    });
  });
});
```

### Pytest Example: Password Validation

```python
# src/domain/password_validator.py
import re
from dataclasses import dataclass

@dataclass
class ValidationResult:
    valid: bool
    errors: list[str]

class PasswordValidator:
    MIN_LENGTH = 8
    MAX_LENGTH = 128

    def validate(self, password: str) -> ValidationResult:
        errors = []

        if len(password) < self.MIN_LENGTH:
            errors.append(f"Password must be at least {self.MIN_LENGTH} characters")

        if len(password) > self.MAX_LENGTH:
            errors.append(f"Password must not exceed {self.MAX_LENGTH} characters")

        if not re.search(r'[A-Z]', password):
            errors.append("Password must contain at least one uppercase letter")

        if not re.search(r'[a-z]', password):
            errors.append("Password must contain at least one lowercase letter")

        if not re.search(r'\d', password):
            errors.append("Password must contain at least one digit")

        if not re.search(r'[!@#$%^&*]', password):
            errors.append("Password must contain at least one special character")

        return ValidationResult(valid=len(errors) == 0, errors=errors)

# tests/domain/test_password_validator.py
import pytest
from src.domain.password_validator import PasswordValidator, ValidationResult

class TestPasswordValidator:
    @pytest.fixture
    def validator(self):
        return PasswordValidator()

    def test_accepts_valid_password(self, validator):
        result = validator.validate("SecurePass1!")

        assert result.valid is True
        assert result.errors == []

    @pytest.mark.parametrize("password,expected_error", [
        ("Short1!", "Password must be at least 8 characters"),
        ("a" * 129 + "A1!", "Password must not exceed 128 characters"),
        ("lowercase1!", "Password must contain at least one uppercase letter"),
        ("UPPERCASE1!", "Password must contain at least one lowercase letter"),
        ("NoDigits!", "Password must contain at least one digit"),
        ("NoSpecial1", "Password must contain at least one special character"),
    ])
    def test_rejects_invalid_passwords(self, validator, password, expected_error):
        result = validator.validate(password)

        assert result.valid is False
        assert expected_error in result.errors

    def test_collects_multiple_validation_errors(self, validator):
        result = validator.validate("bad")

        assert result.valid is False
        assert len(result.errors) == 5  # Too short + missing uppercase + lowercase + digit + special
```

## Integration Tests

### Purpose

- Verify interactions between components
- Test database operations and queries
- Validate API contracts and responses
- Ensure external service integrations work correctly

### Characteristics

- Execute in 1-5 seconds
- May use real databases (often in-memory or containers)
- Test multiple units working together
- Verify data persistence and retrieval

### Jest Example: User Repository with Database

```typescript
// tests/integration/repositories/UserRepository.test.ts
import { UserRepository } from '@/infrastructure/UserRepository';
import { PrismaClient } from '@prisma/client';
import { execSync } from 'child_process';

describe('UserRepository', () => {
  let prisma: PrismaClient;
  let repository: UserRepository;

  beforeAll(async () => {
    // Use test database
    process.env.DATABASE_URL = 'postgresql://test:test@localhost:5432/test_db';
    prisma = new PrismaClient();
    await prisma.$connect();
  });

  afterAll(async () => {
    await prisma.$disconnect();
  });

  beforeEach(async () => {
    // Clean database before each test
    await prisma.user.deleteMany();
    repository = new UserRepository(prisma);
  });

  describe('create', () => {
    it('persists user to database and returns with generated id', async () => {
      const userData = {
        email: 'test@example.com',
        name: 'Test User',
        hashedPassword: 'hashed123',
      };

      const created = await repository.create(userData);

      expect(created.id).toBeDefined();
      expect(created.email).toBe(userData.email);

      // Verify persistence
      const found = await prisma.user.findUnique({
        where: { id: created.id },
      });
      expect(found).not.toBeNull();
      expect(found?.email).toBe(userData.email);
    });

    it('throws ConflictError when email already exists', async () => {
      const userData = {
        email: 'duplicate@example.com',
        name: 'First User',
        hashedPassword: 'hashed123',
      };
      await repository.create(userData);

      await expect(repository.create({
        ...userData,
        name: 'Second User',
      })).rejects.toThrow('User with this email already exists');
    });
  });

  describe('findByEmail', () => {
    it('returns user when found', async () => {
      const created = await repository.create({
        email: 'findme@example.com',
        name: 'Find Me',
        hashedPassword: 'hashed123',
      });

      const found = await repository.findByEmail('findme@example.com');

      expect(found).toEqual(created);
    });

    it('returns null when not found', async () => {
      const found = await repository.findByEmail('nonexistent@example.com');

      expect(found).toBeNull();
    });
  });
});
```

### Pytest Example: API Integration Tests

```python
# tests/integration/api/test_users_api.py
import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from src.main import app
from src.database import Base, get_db

# Test database setup
TEST_DATABASE_URL = "sqlite:///./test.db"
engine = create_engine(TEST_DATABASE_URL, connect_args={"check_same_thread": False})
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

@pytest.fixture(scope="function")
def db_session():
    Base.metadata.create_all(bind=engine)
    session = TestingSessionLocal()
    try:
        yield session
    finally:
        session.close()
        Base.metadata.drop_all(bind=engine)

@pytest.fixture(scope="function")
def client(db_session):
    def override_get_db():
        yield db_session

    app.dependency_overrides[get_db] = override_get_db
    with TestClient(app) as test_client:
        yield test_client
    app.dependency_overrides.clear()

class TestUsersAPI:
    def test_create_user_returns_201_with_user_data(self, client):
        response = client.post("/api/users", json={
            "email": "new@user.com",
            "password": "SecurePass1!",
            "name": "New User"
        })

        assert response.status_code == 201
        data = response.json()
        assert data["email"] == "new@user.com"
        assert data["name"] == "New User"
        assert "id" in data
        assert "password" not in data  # Should not expose password

    def test_create_user_returns_400_for_invalid_email(self, client):
        response = client.post("/api/users", json={
            "email": "invalid-email",
            "password": "SecurePass1!",
            "name": "New User"
        })

        assert response.status_code == 400
        assert "email" in response.json()["detail"].lower()

    def test_create_user_returns_409_for_duplicate_email(self, client):
        user_data = {
            "email": "duplicate@user.com",
            "password": "SecurePass1!",
            "name": "First User"
        }
        client.post("/api/users", json=user_data)

        response = client.post("/api/users", json=user_data)

        assert response.status_code == 409

    def test_get_user_returns_200_when_found(self, client):
        create_response = client.post("/api/users", json={
            "email": "existing@user.com",
            "password": "SecurePass1!",
            "name": "Existing User"
        })
        user_id = create_response.json()["id"]

        response = client.get(f"/api/users/{user_id}")

        assert response.status_code == 200
        assert response.json()["email"] == "existing@user.com"

    def test_get_user_returns_404_when_not_found(self, client):
        response = client.get("/api/users/nonexistent-id")

        assert response.status_code == 404

    def test_list_users_returns_paginated_results(self, client):
        # Create 15 users
        for i in range(15):
            client.post("/api/users", json={
                "email": f"user{i}@example.com",
                "password": "SecurePass1!",
                "name": f"User {i}"
            })

        response = client.get("/api/users?page=1&limit=10")

        assert response.status_code == 200
        data = response.json()
        assert len(data["items"]) == 10
        assert data["total"] == 15
        assert data["page"] == 1
        assert data["pages"] == 2
```

## E2E Tests

### Purpose

- Validate critical user journeys through the entire system
- Ensure the application works from the user's perspective
- Catch integration issues that other tests miss
- Verify deployment and infrastructure configuration

### Characteristics

- Execute in 10-60 seconds
- Run against deployed application (staging or local)
- Use real browser or API client
- Test complete user workflows

### Playwright Example: User Registration Flow

```typescript
// tests/e2e/user-registration.spec.ts
import { test, expect } from '@playwright/test';

test.describe('User Registration', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/register');
  });

  test('allows new user to register and login', async ({ page }) => {
    const uniqueEmail = `test-${Date.now()}@example.com`;

    // Fill registration form
    await page.getByLabel('Email').fill(uniqueEmail);
    await page.getByLabel('Password').fill('SecurePass1!');
    await page.getByLabel('Confirm Password').fill('SecurePass1!');
    await page.getByLabel('Full Name').fill('Test User');
    await page.getByRole('button', { name: 'Create Account' }).click();

    // Verify redirect to login with success message
    await expect(page).toHaveURL('/login');
    await expect(page.getByText('Account created successfully')).toBeVisible();

    // Login with new credentials
    await page.getByLabel('Email').fill(uniqueEmail);
    await page.getByLabel('Password').fill('SecurePass1!');
    await page.getByRole('button', { name: 'Sign In' }).click();

    // Verify logged in state
    await expect(page).toHaveURL('/dashboard');
    await expect(page.getByText('Welcome, Test User')).toBeVisible();
  });

  test('shows validation errors for invalid input', async ({ page }) => {
    // Submit with empty form
    await page.getByRole('button', { name: 'Create Account' }).click();

    // Verify validation messages
    await expect(page.getByText('Email is required')).toBeVisible();
    await expect(page.getByText('Password is required')).toBeVisible();
  });

  test('prevents registration with existing email', async ({ page }) => {
    // Attempt to register with known existing email
    await page.getByLabel('Email').fill('existing@user.com');
    await page.getByLabel('Password').fill('SecurePass1!');
    await page.getByLabel('Confirm Password').fill('SecurePass1!');
    await page.getByLabel('Full Name').fill('Another User');
    await page.getByRole('button', { name: 'Create Account' }).click();

    // Verify error message
    await expect(page.getByText('An account with this email already exists')).toBeVisible();
  });
});
```

### Cypress Example: E-commerce Checkout

```typescript
// cypress/e2e/checkout.cy.ts
describe('Checkout Flow', () => {
  beforeEach(() => {
    cy.loginAs('test@customer.com');
    cy.clearCart();
  });

  it('completes purchase with valid payment', () => {
    // Add items to cart
    cy.visit('/products/widget-pro');
    cy.get('[data-testid="add-to-cart"]').click();
    cy.get('[data-testid="cart-count"]').should('contain', '1');

    // Proceed to checkout
    cy.visit('/cart');
    cy.get('[data-testid="checkout-button"]').click();

    // Fill shipping address
    cy.get('#shipping-address').type('123 Test Street');
    cy.get('#shipping-city').type('Test City');
    cy.get('#shipping-zip').type('12345');
    cy.get('[data-testid="continue-to-payment"]').click();

    // Fill payment (test card)
    cy.getStripeElement('cardNumber').type('4242424242424242');
    cy.getStripeElement('cardExpiry').type('1225');
    cy.getStripeElement('cardCvc').type('123');
    cy.get('[data-testid="place-order"]').click();

    // Verify success
    cy.url().should('include', '/order-confirmation');
    cy.get('[data-testid="order-number"]').should('exist');
    cy.get('[data-testid="order-total"]').should('contain', '$49.99');

    // Verify email sent (via test mailbox API)
    cy.task('getLastEmail', 'test@customer.com').then((email) => {
      expect(email.subject).to.include('Order Confirmation');
    });
  });

  it('handles payment failure gracefully', () => {
    cy.visit('/products/widget-pro');
    cy.get('[data-testid="add-to-cart"]').click();
    cy.visit('/cart');
    cy.get('[data-testid="checkout-button"]').click();

    // Fill shipping
    cy.get('#shipping-address').type('123 Test Street');
    cy.get('#shipping-city').type('Test City');
    cy.get('#shipping-zip').type('12345');
    cy.get('[data-testid="continue-to-payment"]').click();

    // Use declining test card
    cy.getStripeElement('cardNumber').type('4000000000000002');
    cy.getStripeElement('cardExpiry').type('1225');
    cy.getStripeElement('cardCvc').type('123');
    cy.get('[data-testid="place-order"]').click();

    // Verify error handling
    cy.get('[data-testid="payment-error"]')
      .should('contain', 'Your card was declined');
    cy.url().should('include', '/checkout/payment');
    cy.get('[data-testid="cart-count"]').should('contain', '1');
  });
});
```

## Test Selection Guidelines

### When to Write Unit Tests

- Pure functions with business logic
- Data transformations and calculations
- Validation and parsing logic
- State machines and reducers
- Utility and helper functions

### When to Write Integration Tests

- Database queries and transactions
- API endpoints and middleware
- Service-to-service communication
- Message queue producers/consumers
- Cache operations

### When to Write E2E Tests

- User registration and authentication
- Payment and checkout flows
- Critical business workflows
- Multi-step forms and wizards
- Features involving multiple services

## Variations

### Contract Testing (Alternative to E2E)

For microservices, consider contract tests instead of E2E:

```typescript
// consumer/tests/UserServiceContract.test.ts
import { Pact } from '@pact-foundation/pact';

const provider = new Pact({
  consumer: 'OrderService',
  provider: 'UserService',
});

describe('UserService Contract', () => {
  it('returns user details for valid user id', async () => {
    await provider.addInteraction({
      state: 'user 123 exists',
      uponReceiving: 'a request for user 123',
      withRequest: {
        method: 'GET',
        path: '/users/123',
      },
      willRespondWith: {
        status: 200,
        body: {
          id: '123',
          email: Matchers.email(),
          name: Matchers.string(),
        },
      },
    });

    const client = new UserServiceClient(provider.mockService.baseUrl);
    const user = await client.getUser('123');

    expect(user.id).toBe('123');
  });
});
```

### Visual Regression Testing

For UI-heavy applications:

```typescript
// tests/visual/Button.visual.test.ts
import { test, expect } from '@playwright/test';

test.describe('Button Visual Regression', () => {
  test('primary button states', async ({ page }) => {
    await page.goto('/storybook/button--primary');

    // Default state
    await expect(page.locator('.button')).toHaveScreenshot('button-default.png');

    // Hover state
    await page.locator('.button').hover();
    await expect(page.locator('.button')).toHaveScreenshot('button-hover.png');

    // Focus state
    await page.locator('.button').focus();
    await expect(page.locator('.button')).toHaveScreenshot('button-focus.png');
  });
});
```

## Anti-Patterns

### Testing at Wrong Level

```typescript
// WRONG: Using E2E test for simple calculation
test('calculates discount correctly', async ({ page }) => {
  await page.goto('/cart');
  await page.fill('#quantity', '10');
  await expect(page.locator('#discount')).toHaveText('10%');
});

// CORRECT: Unit test for calculation logic
test('applies 10% discount for orders of 10+ items', () => {
  expect(calculateDiscount(10)).toBe(0.1);
});
```

### Missing Critical Paths

```
// WRONG: Skipping E2E for checkout
Unit tests: 95% coverage
Integration tests: API endpoints tested
E2E tests: None for payment flow

// Payment bugs slip through because the complete flow
// (UI -> API -> Payment Provider -> Webhook -> DB) is never tested
```

### Excessive E2E Tests

```
// WRONG: E2E test for every validation rule
e2e/
  registration-empty-email.spec.ts
  registration-invalid-email.spec.ts
  registration-short-password.spec.ts
  registration-no-uppercase.spec.ts
  registration-no-number.spec.ts
  ... 20 more files

// CORRECT: Unit tests for validation, single E2E for happy path
unit/
  PasswordValidator.test.ts  // All validation rules
e2e/
  registration.spec.ts       // One happy path + one error case
```

### Ignoring Test Pyramid Distribution

```
// WRONG: Inverted pyramid
E2E Tests: 500 (slow, flaky)
Integration Tests: 50
Unit Tests: 100

// Build time: 45 minutes
// Flaky test rate: 15%
// Developer confidence: Low

// CORRECT: Following pyramid
Unit Tests: 2000 (fast, stable)
Integration Tests: 300
E2E Tests: 50

// Build time: 8 minutes
// Flaky test rate: 1%
// Developer confidence: High
```
