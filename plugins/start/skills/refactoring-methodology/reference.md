# Refactoring Reference

## Code Smells Catalog

### Method-Level Smells

| Smell | Symptoms | Refactorings |
|-------|----------|--------------|
| Long Method | >20 lines, multiple responsibilities | Extract Method, Decompose Conditional |
| Long Parameter List | >3-4 parameters | Introduce Parameter Object, Preserve Whole Object |
| Duplicate Code | Same/similar logic in multiple places | Extract Method, Pull Up Method, Form Template Method |
| Complex Conditionals | Nested if/else, switch statements | Decompose Conditional, Replace with Guard Clauses |
| Feature Envy | Method uses more data from other class | Move Method, Move Field |
| Data Clumps | Same group of variables passed together | Extract Class, Introduce Parameter Object |
| Speculative Generality | Unused abstractions "for the future" | Collapse Hierarchy, Inline Class, Remove Parameter |
| Dead Code | Unreachable or unused code | Remove Dead Code |

### Class-Level Smells

| Smell | Symptoms | Refactorings |
|-------|----------|--------------|
| Large Class | >200 lines, multiple responsibilities | Extract Class, Extract Subclass |
| God Class | Knows too much, does too much | Extract Class, Move Method |
| Data Class | Only getters/setters, no behavior | Move Method into class, Encapsulate Field |
| Primitive Obsession | Overuse of primitives for domain concepts | Replace Primitive with Object, Extract Class |
| Refused Bequest | Subclass doesn't use inherited members | Replace Inheritance with Delegation |
| Lazy Class | Class that doesn't do enough | Inline Class, Collapse Hierarchy |
| Middle Man | Class that only delegates | Remove Middle Man, Inline Method |
| Parallel Inheritance | Every time you subclass A, you must subclass B | Move Method, Move Field |

### Architecture-Level Smells

| Smell | Symptoms | Refactorings |
|-------|----------|--------------|
| Circular Dependencies | A→B→C→A | Dependency Inversion, Extract Interface |
| Inappropriate Intimacy | Classes too coupled | Move Method, Hide Delegate, Change Bidirectional to Unidirectional |
| Shotgun Surgery | One change requires many file edits | Move Method, Inline Class |
| Divergent Change | One class changed for different reasons | Extract Class |
| Message Chains | a.getB().getC().getD() | Hide Delegate, Extract Method |
| Comments (excessive) | Comments explaining bad code | Extract Method, Rename Method, Introduce Assertion |

---

## Safe Refactoring Patterns

### Extract Method

**When**: Long method with embedded logic blocks

**Before**:
```javascript
function processOrder(order) {
  // Validate order
  if (!order.items || order.items.length === 0) {
    throw new Error('Order must have items');
  }
  if (!order.customer) {
    throw new Error('Order must have customer');
  }

  // Calculate total
  let total = 0;
  for (const item of order.items) {
    total += item.price * item.quantity;
  }

  // Apply discount
  if (order.discountCode) {
    total = total * 0.9;
  }

  return { ...order, total };
}
```

**After**:
```javascript
function processOrder(order) {
  validateOrder(order);
  const total = calculateTotal(order.items, order.discountCode);
  return { ...order, total };
}

function validateOrder(order) {
  if (!order.items || order.items.length === 0) {
    throw new Error('Order must have items');
  }
  if (!order.customer) {
    throw new Error('Order must have customer');
  }
}

function calculateTotal(items, discountCode) {
  let total = items.reduce((sum, item) => sum + item.price * item.quantity, 0);
  if (discountCode) {
    total = total * 0.9;
  }
  return total;
}
```

**Safety**: Run tests after each extraction

---

### Rename (Variable/Method/Class)

**When**: Names don't reveal intent

**Before**:
```javascript
const d = new Date();
function calc(x, y) {
  return x * y * 0.15;
}
```

**After**:
```javascript
const orderDate = new Date();
function calculateTax(subtotal, taxRate) {
  return subtotal * taxRate;
}
```

**Safety**: Use IDE refactoring tools for automatic updates across codebase

---

### Move Method/Field

**When**: Method uses more data from another class than its own

**Before**:
```javascript
class Order {
  calculateShipping(customer) {
    // Uses customer data extensively
    if (customer.isPremium) return 0;
    if (customer.country === 'US') return 5;
    return 15;
  }
}
```

**After**:
```javascript
class Customer {
  calculateShippingFor(order) {
    if (this.isPremium) return 0;
    if (this.country === 'US') return 5;
    return 15;
  }
}

class Order {
  calculateShipping(customer) {
    return customer.calculateShippingFor(this);
  }
}
```

**Safety**: Update all callers, run tests

---

### Replace Conditional with Guard Clauses

**When**: Nested conditionals obscure the main path

**Before**:
```javascript
function getPayAmount(employee) {
  let result;
  if (employee.isSeparated) {
    result = { amount: 0, reason: 'separated' };
  } else {
    if (employee.isRetired) {
      result = { amount: 0, reason: 'retired' };
    } else {
      result = { amount: employee.salary, reason: 'active' };
    }
  }
  return result;
}
```

**After**:
```javascript
function getPayAmount(employee) {
  if (employee.isSeparated) {
    return { amount: 0, reason: 'separated' };
  }
  if (employee.isRetired) {
    return { amount: 0, reason: 'retired' };
  }
  return { amount: employee.salary, reason: 'active' };
}
```

**Safety**: Ensure all paths tested before refactoring

---

### Introduce Parameter Object

**When**: Same group of parameters passed together frequently

**Before**:
```javascript
function createBooking(startDate, endDate, roomType, guestName, guestEmail) {
  // ...
}

function validateBooking(startDate, endDate, roomType) {
  // ...
}

function calculatePrice(startDate, endDate, roomType) {
  // ...
}
```

**After**:
```javascript
class BookingRequest {
  constructor(startDate, endDate, roomType, guestName, guestEmail) {
    this.startDate = startDate;
    this.endDate = endDate;
    this.roomType = roomType;
    this.guestName = guestName;
    this.guestEmail = guestEmail;
  }
}

function createBooking(request) {
  // ...
}

function validateBooking(request) {
  // ...
}

function calculatePrice(request) {
  // ...
}
```

**Safety**: Update all call sites, run tests

---

### Extract Class

**When**: One class has multiple responsibilities

**Before**:
```javascript
class User {
  name;
  email;

  // User data methods
  getName() { return this.name; }
  setName(name) { this.name = name; }

  // Email methods (separate concern)
  sendWelcomeEmail() { /* ... */ }
  sendPasswordReset() { /* ... */ }
  validateEmail() { /* ... */ }
}
```

**After**:
```javascript
class User {
  name;
  email;
  emailService;

  getName() { return this.name; }
  setName(name) { this.name = name; }
}

class EmailService {
  sendWelcomeEmail(user) { /* ... */ }
  sendPasswordReset(user) { /* ... */ }
  validateEmail(email) { /* ... */ }
}
```

**Safety**: Update all dependencies, run tests

---

## Behavior Preservation Checklist

Before ANY refactoring:

- [ ] Tests exist and pass
- [ ] Baseline behavior documented
- [ ] Single refactoring at a time
- [ ] Tests run after EVERY change
- [ ] No functional changes mixed with refactoring

---

## Refactoring Decision Matrix

| Situation | Action |
|-----------|--------|
| Tests passing, clear smell | Proceed with refactoring |
| Tests passing, unclear benefit | Skip or discuss with team |
| Tests failing | Fix tests first, then refactor |
| No tests for area | Add tests first OR skip refactoring |
| Behavior change required | Not refactoring - this is a feature change |

---

## Anti-Patterns to Avoid

### Don't Mix Refactoring with Feature Changes

**Wrong**:
```
Commit: "Refactor user service and add email validation"
```

**Right**:
```
Commit 1: "Refactor: Extract email methods from UserService"
Commit 2: "Feature: Add email validation to registration"
```

### Don't Refactor Without Tests

If code isn't covered by tests:
1. Add characterization tests first
2. Then refactor
3. Or skip and document technical debt

### Don't Refactor Everything at Once

Prioritize by:
1. Code you're actively working on
2. Code with highest change frequency
3. Code with most bugs
