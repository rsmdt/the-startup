# Common Pattern Recognition Examples

## Context

This document provides concrete examples of how to recognize and apply patterns in real codebases. Use these examples as reference when analyzing unfamiliar projects.

## Example 1: React Component Naming Patterns

### Pattern Discovery

When encountering a React codebase, examine existing components:

```
src/components/
  Button.tsx         # PascalCase filename
  UserProfile.tsx    # Multi-word PascalCase
  user-profile.css   # kebab-case for styles
```

### Pattern Recognition

From these files, recognize:
- Component files use PascalCase
- Style files use kebab-case matching component name
- Co-located styles with components

### Pattern Application

When creating a new component:
```
# CORRECT - follows pattern
src/components/
  OrderSummary.tsx
  order-summary.css

# INCORRECT - violates pattern
src/components/
  orderSummary.tsx      # Wrong: camelCase
  OrderSummary.styles.ts # Wrong: different style file convention
```

## Example 2: Python Function Naming

### Pattern Discovery

Survey existing functions:

```python
# Found in existing codebase
def get_user_by_id(user_id: int) -> User:
    ...

def create_order(items: list[Item]) -> Order:
    ...

def update_inventory_count(product_id: int, delta: int) -> None:
    ...
```

### Pattern Recognition

From these functions, recognize:
- snake_case for function names
- get/create/update verb prefixes
- Type hints consistently used
- Descriptive names over abbreviations

### Pattern Application

When adding a new function:
```python
# CORRECT - follows pattern
def delete_user_session(session_id: str) -> bool:
    ...

# INCORRECT - violates pattern
def deleteSession(id):  # Wrong: camelCase, abbreviated, no types
    ...
```

## Example 3: Test File Organization

### Pattern Discovery

Examine test directory structure:

```
src/
  services/
    UserService.ts
    OrderService.ts
tests/
  services/
    UserService.test.ts
    OrderService.test.ts
  fixtures/
    users.json
    orders.json
```

### Pattern Recognition

From this structure, recognize:
- Tests mirror source directory structure
- Test files use `.test.ts` suffix
- Shared fixtures in dedicated directory
- Test filename matches source filename

### Pattern Application

When adding tests for a new service:
```
# CORRECT - follows pattern
src/services/PaymentService.ts
tests/services/PaymentService.test.ts
tests/fixtures/payments.json

# INCORRECT - violates pattern
src/services/PaymentService.ts
src/services/__tests__/payment.spec.ts  # Wrong: different structure and naming
```

## Example 4: API Endpoint Naming

### Pattern Discovery

Review existing API routes:

```typescript
// Found in existing codebase
router.get('/users/:userId', getUser);
router.post('/users', createUser);
router.put('/users/:userId', updateUser);
router.delete('/users/:userId', deleteUser);

router.get('/orders/:orderId/items', getOrderItems);
router.post('/orders/:orderId/items', addOrderItem);
```

### Pattern Recognition

From these routes, recognize:
- Plural nouns for resources (`users` not `user`)
- camelCase for path parameters (`userId` not `user_id`)
- Nested resources use parent path (`orders/:orderId/items`)
- HTTP methods map to CRUD operations

### Pattern Application

When adding a new endpoint:
```typescript
// CORRECT - follows pattern
router.get('/products/:productId/reviews', getProductReviews);
router.post('/products/:productId/reviews', createProductReview);

// INCORRECT - violates pattern
router.get('/product/:product_id/review', getReview);  # Wrong: singular, snake_case
```

## Example 5: Error Handling Patterns

### Pattern Discovery

Examine how errors are handled:

```typescript
// Found in existing codebase
class NotFoundError extends AppError {
  constructor(resource: string, id: string) {
    super(`${resource} with id ${id} not found`, 404);
  }
}

class ValidationError extends AppError {
  constructor(field: string, message: string) {
    super(`Validation failed for ${field}: ${message}`, 400);
  }
}

// Usage
if (!user) {
  throw new NotFoundError('User', userId);
}
```

### Pattern Recognition

From this code, recognize:
- Custom error classes extend base `AppError`
- Error names describe the condition (NotFoundError, ValidationError)
- Errors include HTTP status codes
- Constructors accept contextual parameters
- Error messages are descriptive and include context

### Pattern Application

When adding new error handling:
```typescript
// CORRECT - follows pattern
class ConflictError extends AppError {
  constructor(resource: string, conflictReason: string) {
    super(`${resource} conflict: ${conflictReason}`, 409);
  }
}

// INCORRECT - violates pattern
throw new Error('conflict');  # Wrong: generic Error, no context, no status
```

## Example 6: Configuration Patterns

### Pattern Discovery

Review how configuration is structured:

```typescript
// config/database.ts
export const databaseConfig = {
  host: process.env.DB_HOST || 'localhost',
  port: parseInt(process.env.DB_PORT || '5432'),
  name: process.env.DB_NAME || 'app_dev',
};

// config/auth.ts
export const authConfig = {
  jwtSecret: process.env.JWT_SECRET || 'dev-secret',
  tokenExpiry: parseInt(process.env.TOKEN_EXPIRY || '3600'),
  refreshExpiry: parseInt(process.env.REFRESH_EXPIRY || '86400'),
};
```

### Pattern Recognition

From these files, recognize:
- Config split by domain (database, auth)
- Environment variables with fallback defaults
- Type coercion for non-string values
- Exported as named objects
- Env var naming: SCREAMING_SNAKE_CASE

### Pattern Application

When adding new configuration:
```typescript
// config/email.ts - CORRECT
export const emailConfig = {
  smtpHost: process.env.SMTP_HOST || 'localhost',
  smtpPort: parseInt(process.env.SMTP_PORT || '587'),
  senderAddress: process.env.SENDER_ADDRESS || 'noreply@example.com',
};
```

## Example 7: State Management Patterns

### Pattern Discovery

Examine how state is managed:

```typescript
// store/users/types.ts
interface UserState {
  items: User[];
  loading: boolean;
  error: string | null;
}

// store/users/actions.ts
const fetchUsers = createAsyncThunk('users/fetch', async () => {...});
const updateUser = createAsyncThunk('users/update', async (user: User) => {...});

// store/users/slice.ts
const usersSlice = createSlice({
  name: 'users',
  initialState,
  reducers: {...},
  extraReducers: {...},
});
```

### Pattern Recognition

From this code, recognize:
- Feature-based organization (store/users/)
- Separate files for types, actions, slice
- Consistent state shape (items, loading, error)
- Async thunk naming: `domain/action`
- Slice naming matches domain

### Pattern Application

When adding new state:
```
store/
  orders/
    types.ts      # OrderState interface
    actions.ts    # fetchOrders, updateOrder thunks
    slice.ts      # ordersSlice
```

## Variations

### When Patterns Conflict

If you find conflicting patterns (e.g., some files use one style, others use another):

1. Check for date/author patterns - newer code may represent intended direction
2. Look for documentation or style guides
3. Follow the pattern in the specific area you are modifying
4. When in doubt, ask the team or document your choice

### When No Clear Pattern Exists

If a pattern is not established:

1. Check similar projects in the same ecosystem
2. Follow language/framework conventions
3. Document the pattern you establish
4. Be prepared to refactor if team prefers different approach

## Anti-Patterns

### Mixing Conventions

```typescript
// WRONG: Mixed naming in same file
const getUserData = () => {...};        // camelCase verb
const fetch_user_profile = () => {...}; // snake_case verb
const retrieveUserInfo = () => {...};   // different verb pattern
```

### Inconsistent Organization

```
// WRONG: Inconsistent test placement
src/
  UserService.ts
  UserService.test.ts      # Co-located
  OrderService.ts
tests/
  OrderService.test.ts     # Separate directory
```

### Pattern Drift

```typescript
// WRONG: Old pattern in new code
// Existing: Modern async/await
async function getUser(id) {
  const user = await db.findUser(id);
  return user;
}

// New code using outdated callback pattern
function getOrder(id, callback) {
  db.findOrder(id, (err, order) => {
    callback(err, order);
  });
}
```
