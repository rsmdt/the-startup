# Anti-Pattern Catalog

Detailed examples of common code anti-patterns: what they look like in practice, why they cause problems, and how to refactor them. Use this when a SKILL.md table entry needs deeper analysis or when explaining a finding to a developer.

---

## God Object

### What It Looks Like

A class that accumulates responsibilities over time until it knows too much and does too much. Often starts as a reasonable service class and grows through expedient additions.

```
UserManager:
  createUser()
  updateUser()
  deleteUser()
  getUserById()
  getUserByEmail()
  sendWelcomeEmail()
  sendPasswordResetEmail()
  validateUserPermissions()
  generateAuthToken()
  revokeAuthToken()
  refreshAuthToken()
  logUserActivity()
  getUserActivityReport()
  exportUserDataToCsv()
  importUsersFromCsv()
  syncUsersWithLdap()
```

### Why It Is Problematic

- Every team member who touches users touches this file — merge conflicts are constant
- Changing authentication logic risks breaking CSV export
- The class is untestable in isolation because it depends on email, database, LDAP, and auth systems simultaneously
- New developers cannot infer what "the system does with users" because the answer is "everything"
- The class has no single reason to change — it changes for every user-related requirement

### How to Refactor

Identify natural responsibility clusters. Each cluster becomes a focused class:

```
UserRepository       — persistence only (create, update, delete, findBy*)
UserEmailService     — email communication (sendWelcome, sendPasswordReset)
AuthTokenService     — token lifecycle (generate, revoke, refresh)
PermissionService    — authorization checks (validatePermissions)
UserActivityLogger   — audit trail (logActivity, getReport)
UserImportExport     — CSV operations (import, export)
LdapSyncService      — directory sync (syncUsers)
```

Each class has a single reason to change. `UserRepository` changes when the persistence layer changes. `UserEmailService` changes when email templates change. They can be tested independently.

Refactor incrementally — extract one cluster at a time, run tests between each extraction.

---

## Spaghetti Code

### What It Looks Like

Control flow that jumps unpredictably through conditionals, early returns, exceptions used for control flow, and functions that call back into their callers. The execution path is difficult to trace.

```python
def process_payment(order, user, payment_method):
    if payment_method == 'card':
        if not user.has_saved_card:
            if order.total > 0:
                try:
                    result = stripe.charge(order.total, user.card_token)
                    if result.status == 'succeeded':
                        order.status = 'paid'
                        send_receipt(user, order)
                        if user.loyalty_enabled:
                            points = calculate_points(order.total)
                            if points > 0:
                                user.loyalty_points += points
                                if user.loyalty_points >= 1000:
                                    send_loyalty_reward(user)
                    else:
                        order.status = 'failed'
                        log_failure(result)
                        if user.retry_enabled:
                            retry_payment(order, user)
                except stripe.CardError as e:
                    handle_card_error(e, order, user)
                except Exception as e:
                    raise PaymentException("unexpected") from e
        else:
            raise ValueError("no saved card")
    elif payment_method == 'invoice':
        # 40 more lines
```

### Why It Is Problematic

- Adding a new payment method requires understanding the entire existing structure
- The loyalty points logic is buried inside payment processing — impossible to test in isolation
- Cyclomatic complexity is high (10+ paths) — almost certainly under-tested
- When a bug is found in the card error path, the fix risks destabilizing the invoice path
- Reading the function requires holding many open conditionals in working memory

### How to Refactor

Separate concerns. Each logical step becomes a named function at the same abstraction level:

```python
def process_payment(order, user, payment_method):
    validate_payment_prerequisites(order, user, payment_method)
    result = charge_payment(order, user, payment_method)
    apply_payment_result(result, order, user)

def validate_payment_prerequisites(order, user, payment_method):
    if payment_method == 'card' and not user.has_saved_card:
        raise MissingPaymentMethodError("User has no saved card on file")

def charge_payment(order, user, payment_method):
    if payment_method == 'card':
        return _charge_card(order, user)
    if payment_method == 'invoice':
        return _issue_invoice(order, user)
    raise UnsupportedPaymentMethodError(payment_method)

def apply_payment_result(result, order, user):
    order.status = 'paid' if result.succeeded else 'failed'
    if result.succeeded:
        _handle_successful_payment(order, user)
    else:
        _handle_failed_payment(result, order, user)

def _handle_successful_payment(order, user):
    send_receipt(user, order)
    _award_loyalty_points_if_eligible(user, order.total)

def _award_loyalty_points_if_eligible(user, amount):
    if not user.loyalty_enabled:
        return
    points = calculate_points(amount)
    user.loyalty_points += points
    if user.loyalty_points >= LOYALTY_REWARD_THRESHOLD:
        send_loyalty_reward(user)
```

Each function now has a single responsibility and a complexity of 1-3. The loyalty logic can be tested without processing a payment.

---

## Copy-Paste Programming

### What It Looks Like

Identical or near-identical logic duplicated across the codebase. Often spotted by similar variable names, same comment blocks, or adjacent functions that differ only in the data they operate on.

```javascript
function validateUserRegistration(data) {
  const errors = [];
  if (!data.email || !data.email.includes('@')) {
    errors.push('Email is invalid');
  }
  if (!data.password || data.password.length < 8) {
    errors.push('Password must be at least 8 characters');
  }
  if (!data.name || data.name.trim().length === 0) {
    errors.push('Name is required');
  }
  return errors;
}

function validateUserUpdate(data) {
  const errors = [];
  if (data.email && !data.email.includes('@')) {
    errors.push('Email is invalid');
  }
  if (data.password && data.password.length < 8) {
    errors.push('Password must be at least 8 characters');
  }
  if (data.name !== undefined && data.name.trim().length === 0) {
    errors.push('Name is required');
  }
  return errors;
}
```

### Why It Is Problematic

- A bug fix in one copy must be found and applied to all copies — and copies are often not found
- The two functions have slightly different semantics (required vs optional fields) but the duplication obscures the intentional difference from the accidental duplication
- When the email validation rule changes (e.g., to use a proper regex), it will be changed in one place and missed in the other
- The next developer will copy-paste again rather than extract, deepening the problem

### How to Refactor

Extract the shared validation rules. Parameterize what varies:

```javascript
const emailRule = (required = true) => ({
  validate: (value) => !required && !value ? null : (value?.includes('@') ? null : 'Email is invalid'),
});

const passwordRule = (required = true) => ({
  validate: (value) => !required && !value ? null : (value?.length >= 8 ? null : 'Password must be at least 8 characters'),
});

const nameRule = (required = true) => ({
  validate: (value) => !required && value === undefined ? null : (value?.trim().length > 0 ? null : 'Name is required'),
});

function validate(data, rules) {
  return Object.entries(rules)
    .map(([field, rule]) => rule.validate(data[field]))
    .filter(Boolean);
}

const registrationRules = {
  email: emailRule(true),
  password: passwordRule(true),
  name: nameRule(true),
};

const updateRules = {
  email: emailRule(false),
  password: passwordRule(false),
  name: nameRule(false),
};

function validateUserRegistration(data) {
  return validate(data, registrationRules);
}

function validateUserUpdate(data) {
  return validate(data, updateRules);
}
```

The intentional difference (required vs optional) is now explicit. The validation logic exists once.

---

## Premature Optimization

### What It Looks Like

Complex code that sacrifices readability for performance before there is evidence of a performance problem. Often includes caching for operations that run once, bit manipulation for clarity, or micro-optimizations that modern compilers perform automatically.

```python
# Processing a list of ~50 user records from a form submission
def process_users(users):
    # "Optimized" with manual hash map to avoid O(n) lookup
    user_index = {}
    for i in range(len(users)):
        user_index[users[i]['id']] = i

    # Bit manipulation instead of modulo
    batch_size = 16
    batch_mask = batch_size - 1

    results = [None] * len(users)  # Pre-allocated for "performance"
    batch_buffer = []

    for i in range(len(users)):
        batch_buffer.append(users[i])
        if (i & batch_mask) == batch_mask or i == len(users) - 1:
            # Process batch
            for user in batch_buffer:
                idx = user_index[user['id']]
                results[idx] = transform_user(user)
            batch_buffer = []

    return results
```

### Why It Is Problematic

- The code processes at most 50 items — the "optimization" provides zero measurable benefit at this scale
- A reviewer must understand bit masking, index remapping, and manual batching to verify correctness
- The complexity introduces a real bug opportunity (off-by-one in batch boundary detection)
- Performance assumptions baked into the code are not documented — the next developer does not know why this is complex
- If the actual bottleneck is the `transform_user` function or a database call, this optimization is in the wrong place

### How to Refactor

Write the simple version. Measure if performance becomes an issue. Optimize the measured bottleneck:

```python
def process_users(users):
    return [transform_user(user) for user in users]
```

If profiling reveals this is genuinely slow at real-world scale, the optimization is warranted. Document it:

```python
def process_users(users):
    # Process in batches of 500 to stay within API rate limits.
    # Batch size determined by load testing at 10k users/minute.
    # See: docs/adr/0012-user-batch-processing.md
    BATCH_SIZE = 500
    results = []
    for batch in chunked(users, BATCH_SIZE):
        results.extend(_process_batch(batch))
    return results
```

---

## Magic Numbers

### What It Looks Like

Numeric (or string) literals embedded in logic without explanation of what they represent or why they have that value.

```javascript
function calculateShipping(orderTotal, destinationCountry) {
  if (destinationCountry === 'US') {
    if (orderTotal >= 75) {
      return 0;
    }
    return 8.99;
  }

  if (orderTotal >= 150) {
    return 0;
  }

  if (['CA', 'MX'].includes(destinationCountry)) {
    return 15.99;
  }

  return 29.99;
}

function calculateLoyaltyPoints(orderTotal) {
  return Math.floor(orderTotal * 1.5);
}

setTimeout(retryPayment, 30000);
```

### Why It Is Problematic

- `75` and `150` appear to be free shipping thresholds — but a reader cannot confirm this without business documentation
- When marketing changes the domestic free shipping threshold from $75 to $100, the developer must find every `75` in the codebase (risky with grep)
- `1.5` in loyalty points has no unit — is it 1.5 points per dollar? Per cent? The calculation is opaque
- `30000` is milliseconds — readers must know this and convert mentally
- Magic numbers in tests are particularly dangerous: they assert specific values without explaining why those values are expected

### How to Refactor

Name every value that has domain meaning:

```javascript
const FREE_SHIPPING_THRESHOLD_DOMESTIC = 75.00;  // USD
const FREE_SHIPPING_THRESHOLD_INTERNATIONAL = 150.00;  // USD
const SHIPPING_COST_DOMESTIC = 8.99;
const SHIPPING_COST_NORTH_AMERICA = 15.99;
const SHIPPING_COST_INTERNATIONAL = 29.99;
const NORTH_AMERICA_COUNTRIES = ['CA', 'MX'];

const LOYALTY_POINTS_PER_DOLLAR = 1.5;

const PAYMENT_RETRY_DELAY_MS = 30_000;  // 30 seconds

function calculateShipping(orderTotal, destinationCountry) {
  if (destinationCountry === 'US') {
    return orderTotal >= FREE_SHIPPING_THRESHOLD_DOMESTIC ? 0 : SHIPPING_COST_DOMESTIC;
  }
  if (orderTotal >= FREE_SHIPPING_THRESHOLD_INTERNATIONAL) {
    return 0;
  }
  return NORTH_AMERICA_COUNTRIES.includes(destinationCountry)
    ? SHIPPING_COST_NORTH_AMERICA
    : SHIPPING_COST_INTERNATIONAL;
}

function calculateLoyaltyPoints(orderTotal) {
  return Math.floor(orderTotal * LOYALTY_POINTS_PER_DOLLAR);
}

setTimeout(retryPayment, PAYMENT_RETRY_DELAY_MS);
```

Constants that change together (e.g., all shipping thresholds) should live in a single configuration file — not scattered across the codebase.

---

## Deep Nesting

### What It Looks Like

Control flow that indents further and further right, making it difficult to track which conditions govern a given block of code.

```typescript
function processOrderItem(order, item, user) {
  if (order !== null) {
    if (order.status === 'active') {
      if (item !== null) {
        if (item.inStock) {
          if (user !== null) {
            if (user.isVerified) {
              if (item.price <= user.creditLimit) {
                // The actual logic — buried at level 7
                addItemToOrder(order, item);
                deductFromCreditLimit(user, item.price);
                updateInventory(item);
                return { success: true };
              } else {
                return { success: false, reason: 'insufficient_credit' };
              }
            } else {
              return { success: false, reason: 'unverified_user' };
            }
          } else {
            return { success: false, reason: 'no_user' };
          }
        } else {
          return { success: false, reason: 'out_of_stock' };
        }
      } else {
        return { success: false, reason: 'no_item' };
      }
    } else {
      return { success: false, reason: 'inactive_order' };
    }
  } else {
    return { success: false, reason: 'no_order' };
  }
}
```

### Why It Is Problematic

- The happy path (the actual logic) is invisible — buried at the deepest level of nesting
- Each closing brace must be mentally matched to its opening — readers count braces
- Adding a new guard condition requires restructuring the entire block
- The function is 30+ lines for logic that occupies 3 lines
- Error handling is scattered — the `insufficient_credit` case and the `no_order` case look structurally identical but carry completely different meaning

### How to Refactor

Invert conditions with guard clauses. Return early on invalid cases. The happy path stays at the left margin:

```typescript
function processOrderItem(order, item, user) {
  if (order === null) return { success: false, reason: 'no_order' };
  if (order.status !== 'active') return { success: false, reason: 'inactive_order' };
  if (item === null) return { success: false, reason: 'no_item' };
  if (!item.inStock) return { success: false, reason: 'out_of_stock' };
  if (user === null) return { success: false, reason: 'no_user' };
  if (!user.isVerified) return { success: false, reason: 'unverified_user' };
  if (item.price > user.creditLimit) return { success: false, reason: 'insufficient_credit' };

  addItemToOrder(order, item);
  deductFromCreditLimit(user, item.price);
  updateInventory(item);

  return { success: true };
}
```

Same logic, same behavior, zero nesting. The happy path is immediately visible. New guard conditions are added as a new line at the top, not a new level of indentation.

---

## Feature Envy

### What It Looks Like

A method that is more interested in another class's data than its own. It reaches into another object repeatedly to do work that belongs to that other object.

```python
class Order:
    def __init__(self, items, customer):
        self.items = items
        self.customer = customer

    def calculate_discount(self):
        # This method "envies" the Customer class
        if self.customer.membership_tier == 'gold':
            if self.customer.years_as_customer >= 5:
                return 0.20
            return 0.15
        if self.customer.membership_tier == 'silver':
            if self.customer.years_as_customer >= 3:
                return 0.10
            return 0.05
        if self.customer.total_lifetime_spend > 10000:
            return 0.03
        return 0

class InvoiceGenerator:
    def format_shipping_address(self, order):
        # This method knows too much about Customer's address structure
        customer = order.customer
        line1 = customer.address.street_number + ' ' + customer.address.street_name
        if customer.address.apartment:
            line1 += ', Apt ' + customer.address.apartment
        city_line = customer.address.city + ', ' + customer.address.state
        if customer.address.country != 'US':
            city_line += ', ' + customer.address.country
        return line1 + '\n' + city_line + '\n' + customer.address.postal_code
```

### Why It Is Problematic

- Discount rules live in `Order` but are entirely about `Customer` properties — when discount rules change, `Order` changes
- If `Customer.address` gains a new field, `InvoiceGenerator` must be updated — a coupling that should not exist
- The knowledge about how to format an address is now in two places: the generator and wherever else code accesses these fields
- Testing `calculate_discount` requires constructing full `Customer` objects instead of testing the discount logic alone

### How to Refactor

Move the method to the class whose data it uses:

```python
class Customer:
    def discount_rate(self):
        if self.membership_tier == 'gold':
            return 0.20 if self.years_as_customer >= 5 else 0.15
        if self.membership_tier == 'silver':
            return 0.10 if self.years_as_customer >= 3 else 0.05
        if self.total_lifetime_spend > 10000:
            return 0.03
        return 0

    def formatted_shipping_address(self):
        return self.address.formatted()

class Address:
    def formatted(self):
        line1 = f"{self.street_number} {self.street_name}"
        if self.apartment:
            line1 += f", Apt {self.apartment}"
        city_line = f"{self.city}, {self.state}"
        if self.country != 'US':
            city_line += f", {self.country}"
        return f"{line1}\n{city_line}\n{self.postal_code}"

class Order:
    def calculate_discount(self):
        return self.customer.discount_rate()

class InvoiceGenerator:
    def format_shipping_address(self, order):
        return order.customer.formatted_shipping_address()
```

Each class owns its own data and the logic that operates on it. Address formatting is in `Address`. Discount calculation is in `Customer`.

---

## Shotgun Surgery

### What It Looks Like

A single conceptual change requires edits to many files. Common when a concern is spread across the codebase rather than encapsulated.

```
Scenario: "Add a new user role: 'moderator'"

Files that must change:
  src/auth/permissions.ts        — add role to enum
  src/auth/guards/role.guard.ts  — add role to allowed list
  src/middleware/authorize.ts    — handle new role in switch
  src/db/migrations/...          — add role to DB enum
  src/db/seeds/roles.ts          — add seed record
  src/api/users/users.controller.ts — add role to DTO validation
  src/api/users/users.service.ts — add role-specific logic
  src/frontend/components/RoleBadge.tsx — add display case
  src/frontend/hooks/usePermissions.ts — add permission checks
  src/tests/fixtures/users.ts    — update test fixtures
  src/tests/auth/permissions.test.ts — add test cases
  docs/api/users.md              — update API documentation
```

### Why It Is Problematic

- Adding a role takes hours instead of minutes — not because of complexity, but because of coordination overhead
- Each file change is a merge conflict opportunity
- It is easy to miss one file — the bug only appears when that specific code path is triggered
- The codebase gives no indication that these files are related — a developer must know to look
- Reviewing the change is difficult: the reviewer must trace the conceptual change across 12 files

### How to Refactor

Encapsulate the concern. A role definition should exist in one place:

```typescript
// src/domain/roles/role.ts
export const Role = {
  ADMIN: 'admin',
  MODERATOR: 'moderator',
  USER: 'user',
} as const;

export type RoleType = typeof Role[keyof typeof Role];

export const RolePermissions: Record<RoleType, Permission[]> = {
  [Role.ADMIN]: [Permission.READ, Permission.WRITE, Permission.DELETE, Permission.MODERATE],
  [Role.MODERATOR]: [Permission.READ, Permission.WRITE, Permission.MODERATE],
  [Role.USER]: [Permission.READ, Permission.WRITE],
};

export const RoleDisplayName: Record<RoleType, string> = {
  [Role.ADMIN]: 'Administrator',
  [Role.MODERATOR]: 'Moderator',
  [Role.USER]: 'User',
};
```

Now adding a role is a single file change. Everything that depends on role definitions imports from this one location. The migration and tests still change, but the application logic does not scatter.

The test for whether Shotgun Surgery is resolved: can a developer make the conceptual change (add a role) by editing one logical location?

---

## Pattern Summary

| Anti-Pattern | Primary Dimension | Key Signal | First Refactoring Step |
|--------------|-------------------|------------|------------------------|
| God Object | Design | Class with 10+ methods spanning unrelated concerns | List responsibilities, find natural clusters |
| Spaghetti Code | Readability | Cannot trace execution path without diagramming | Identify the happy path, separate guard clauses |
| Copy-Paste Programming | Maintainability | Two functions that differ in one parameter | Find what varies, extract shared logic |
| Premature Optimization | Readability | Complexity added before measuring performance | Write the obvious version, delete the complex one |
| Magic Numbers | Readability | Literal values in logic with no named constant | Name every value that has domain meaning |
| Deep Nesting | Readability | Happy path is at indent level 4+ | Invert conditions, return early on failure |
| Feature Envy | Design | Method uses another class's data more than its own | Move method to the class whose data it uses |
| Shotgun Surgery | Maintainability | One concept requires N file changes | Find the concept, encapsulate it in one place |
