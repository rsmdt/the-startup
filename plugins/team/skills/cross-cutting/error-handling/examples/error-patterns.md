# Error Patterns Examples

Concrete examples demonstrating error handling patterns. Language-agnostic principles with syntax examples.

---

## Input Validation

### Context

Validate at system boundaries before processing. Return all errors, not just the first one.

### Pattern: Collect All Validation Errors

```typescript
// TypeScript example
interface ValidationError {
  field: string;
  message: string;
  code: string;
}

function validateUserInput(input: unknown): ValidationError[] {
  const errors: ValidationError[] = [];

  if (!input || typeof input !== 'object') {
    return [{ field: 'root', message: 'Input must be an object', code: 'INVALID_TYPE' }];
  }

  const data = input as Record<string, unknown>;

  if (!data.email) {
    errors.push({ field: 'email', message: 'Email is required', code: 'REQUIRED' });
  } else if (!isValidEmail(data.email)) {
    errors.push({ field: 'email', message: 'Email format is invalid', code: 'INVALID_FORMAT' });
  }

  if (!data.age) {
    errors.push({ field: 'age', message: 'Age is required', code: 'REQUIRED' });
  } else if (typeof data.age !== 'number' || data.age < 0 || data.age > 150) {
    errors.push({ field: 'age', message: 'Age must be between 0 and 150', code: 'OUT_OF_RANGE' });
  }

  return errors;
}
```

```python
# Python example
from dataclasses import dataclass
from typing import List, Any, Dict

@dataclass
class ValidationError:
    field: str
    message: str
    code: str

def validate_user_input(data: Any) -> List[ValidationError]:
    errors = []

    if not isinstance(data, dict):
        return [ValidationError('root', 'Input must be a dictionary', 'INVALID_TYPE')]

    if not data.get('email'):
        errors.append(ValidationError('email', 'Email is required', 'REQUIRED'))
    elif not is_valid_email(data['email']):
        errors.append(ValidationError('email', 'Email format is invalid', 'INVALID_FORMAT'))

    if 'age' not in data:
        errors.append(ValidationError('age', 'Age is required', 'REQUIRED'))
    elif not isinstance(data['age'], int) or not 0 <= data['age'] <= 150:
        errors.append(ValidationError('age', 'Age must be between 0 and 150', 'OUT_OF_RANGE'))

    return errors
```

### Explanation

1. Check the most fundamental constraint first (is it an object/dict?)
2. For each field, check required before format/range
3. Collect all errors into a list
4. Return the complete list so users can fix all issues at once

### Anti-Pattern

```typescript
// BAD: Throws on first error, user must fix one at a time
function validateBad(data: any): void {
  if (!data.email) throw new Error('Email required');
  if (!data.age) throw new Error('Age required');
}
```

---

## Error Types and Custom Errors

### Context

Distinguish error types to enable appropriate handling at different layers.

### Pattern: Error Hierarchy

```typescript
// Base error with shared properties
class AppError extends Error {
  constructor(
    message: string,
    public readonly code: string,
    public readonly isOperational: boolean,
    public readonly context?: Record<string, unknown>
  ) {
    super(message);
    this.name = this.constructor.name;
  }
}

// Operational errors - expected, handle gracefully
class ValidationError extends AppError {
  constructor(message: string, context?: Record<string, unknown>) {
    super(message, 'VALIDATION_ERROR', true, context);
  }
}

class NotFoundError extends AppError {
  constructor(resource: string, id: string) {
    super(`${resource} not found: ${id}`, 'NOT_FOUND', true, { resource, id });
  }
}

class ExternalServiceError extends AppError {
  constructor(service: string, cause?: Error) {
    super(`External service failure: ${service}`, 'EXTERNAL_SERVICE_ERROR', true, {
      service,
      cause: cause?.message
    });
  }
}

// Programmer errors - unexpected, fail fast
class InvariantViolation extends AppError {
  constructor(invariant: string) {
    super(`Invariant violated: ${invariant}`, 'INVARIANT_VIOLATION', false);
  }
}
```

```python
# Python example
class AppError(Exception):
    def __init__(self, message: str, code: str, is_operational: bool, context: dict = None):
        super().__init__(message)
        self.code = code
        self.is_operational = is_operational
        self.context = context or {}

class ValidationError(AppError):
    def __init__(self, message: str, context: dict = None):
        super().__init__(message, 'VALIDATION_ERROR', True, context)

class NotFoundError(AppError):
    def __init__(self, resource: str, resource_id: str):
        super().__init__(
            f'{resource} not found: {resource_id}',
            'NOT_FOUND',
            True,
            {'resource': resource, 'id': resource_id}
        )
```

### Explanation

1. Base error carries common properties: code, operational flag, context
2. `isOperational` distinguishes expected errors from bugs
3. Context provides structured data for logging without string parsing
4. Specific error types enable pattern matching in handlers

---

## Recovery Strategies

### Context

Implement retry with exponential backoff for transient failures.

### Pattern: Retry with Backoff

```typescript
interface RetryConfig {
  maxAttempts: number;
  baseDelayMs: number;
  maxDelayMs: number;
}

async function withRetry<T>(
  operation: () => Promise<T>,
  config: RetryConfig,
  isRetryable: (error: unknown) => boolean
): Promise<T> {
  let lastError: unknown;

  for (let attempt = 1; attempt <= config.maxAttempts; attempt++) {
    try {
      return await operation();
    } catch (error) {
      lastError = error;

      if (!isRetryable(error) || attempt === config.maxAttempts) {
        throw error;
      }

      // Exponential backoff with jitter
      const delay = Math.min(
        config.baseDelayMs * Math.pow(2, attempt - 1) + Math.random() * 100,
        config.maxDelayMs
      );

      await sleep(delay);
    }
  }

  throw lastError;
}

// Usage
const result = await withRetry(
  () => fetchFromExternalAPI(url),
  { maxAttempts: 3, baseDelayMs: 100, maxDelayMs: 5000 },
  (error) => error instanceof ExternalServiceError
);
```

### Explanation

1. Accept retry configuration as parameter for flexibility
2. Predicate function determines which errors are retryable
3. Exponential backoff: 100ms, 200ms, 400ms...
4. Jitter prevents thundering herd when many clients retry
5. Max delay caps the wait time
6. Re-throw on final attempt or non-retryable errors

### Variations

- **Circuit breaker:** Track failure rate, stop retrying when threshold exceeded
- **Fallback:** Return cached/default value instead of retrying
- **Hedged requests:** Start second request before first times out

---

## User-Facing vs Internal Errors

### Context

Users need actionable guidance. Logs need technical detail.

### Pattern: Error Response Transformation

```typescript
interface UserErrorResponse {
  message: string;
  code: string;
  requestId: string;
}

interface LogEntry {
  timestamp: string;
  requestId: string;
  userId?: string;
  error: {
    name: string;
    message: string;
    code: string;
    stack?: string;
  };
  context: Record<string, unknown>;
}

function handleError(
  error: unknown,
  requestId: string,
  userId?: string
): { response: UserErrorResponse; logEntry: LogEntry } {
  const timestamp = new Date().toISOString();

  // Determine if error is safe to expose
  const isOperational = error instanceof AppError && error.isOperational;

  // User gets sanitized message
  const response: UserErrorResponse = {
    message: isOperational
      ? (error as AppError).message
      : 'An unexpected error occurred. Please try again.',
    code: isOperational
      ? (error as AppError).code
      : 'INTERNAL_ERROR',
    requestId
  };

  // Logs get full context
  const logEntry: LogEntry = {
    timestamp,
    requestId,
    userId,
    error: {
      name: error instanceof Error ? error.name : 'UnknownError',
      message: error instanceof Error ? error.message : String(error),
      code: error instanceof AppError ? error.code : 'UNKNOWN',
      stack: error instanceof Error ? error.stack : undefined
    },
    context: error instanceof AppError ? error.context ?? {} : {}
  };

  return { response, logEntry };
}
```

### Explanation

1. Request ID links user error to log entry for support
2. Operational errors: safe to show message to user
3. Programmer errors: generic message, full details in logs
4. Never expose stack traces to users
5. Structured log entry enables querying and alerting

---

## Centralized Error Handling

### Context

Avoid scattered try/catch blocks. Handle errors at defined boundaries.

### Pattern: Error Boundary Middleware

```typescript
// Express.js example
function errorHandler(
  error: unknown,
  req: Request,
  res: Response,
  next: NextFunction
): void {
  const requestId = req.headers['x-request-id'] as string;
  const { response, logEntry } = handleError(error, requestId, req.user?.id);

  // Log with appropriate level
  if (error instanceof AppError && error.isOperational) {
    logger.warn(logEntry);
  } else {
    logger.error(logEntry);
    // Alert on programmer errors
    alerting.notify('Unexpected error', logEntry);
  }

  // Send appropriate status code
  const statusCode = getStatusCode(error);
  res.status(statusCode).json(response);
}

function getStatusCode(error: unknown): number {
  if (error instanceof ValidationError) return 400;
  if (error instanceof NotFoundError) return 404;
  if (error instanceof AuthenticationError) return 401;
  if (error instanceof AuthorizationError) return 403;
  if (error instanceof ExternalServiceError) return 502;
  return 500;
}
```

### Explanation

1. Single error handler for the entire application
2. Maps error types to HTTP status codes
3. Different log levels for operational vs programmer errors
4. Alerts only on unexpected errors to avoid noise
5. Controllers throw errors, middleware handles them

---

## Structured Logging

### Context

Machine-parseable logs enable querying, aggregation, and alerting.

### Pattern: Structured Log Format

```json
{
  "timestamp": "2024-01-15T10:23:45.123Z",
  "level": "error",
  "service": "user-service",
  "requestId": "req-abc-123",
  "traceId": "trace-xyz-789",
  "userId": "user-456",
  "operation": "createUser",
  "duration_ms": 234,
  "error": {
    "name": "ValidationError",
    "code": "VALIDATION_ERROR",
    "message": "Email format is invalid",
    "fields": ["email"]
  },
  "metadata": {
    "endpoint": "/api/users",
    "method": "POST"
  }
}
```

### Key Fields

| Field | Purpose |
|-------|---------|
| `timestamp` | ISO 8601 format for sorting and filtering |
| `level` | Severity for filtering (error, warn, info, debug) |
| `service` | Origin service in distributed systems |
| `requestId` | Links logs for single request |
| `traceId` | Links logs across services |
| `operation` | Business operation being performed |
| `duration_ms` | Performance monitoring |
| `error` | Structured error details |

### Anti-Pattern

```
// BAD: Unstructured, hard to parse
console.log("Error: " + error.message + " for user " + userId);
```

---

## Testing Error Paths

### Context

Error handling code must be tested as rigorously as success paths.

### Pattern: Test Error Scenarios

```typescript
describe('UserService', () => {
  describe('createUser', () => {
    it('returns validation errors for invalid email', async () => {
      const result = await userService.createUser({ email: 'invalid', age: 25 });

      expect(result.isErr()).toBe(true);
      expect(result.error).toBeInstanceOf(ValidationError);
      expect(result.error.context.fields).toContain('email');
    });

    it('retries on transient database errors', async () => {
      database.save
        .mockRejectedValueOnce(new ConnectionError())
        .mockRejectedValueOnce(new ConnectionError())
        .mockResolvedValueOnce({ id: '123' });

      const result = await userService.createUser(validInput);

      expect(result.isOk()).toBe(true);
      expect(database.save).toHaveBeenCalledTimes(3);
    });

    it('fails fast on constraint violations', async () => {
      database.save.mockRejectedValue(new UniqueConstraintError('email'));

      const result = await userService.createUser(validInput);

      expect(result.isErr()).toBe(true);
      expect(database.save).toHaveBeenCalledTimes(1); // No retry
    });
  });
});
```

### What to Test

- Validation errors return correct field information
- Retryable errors trigger retry logic
- Non-retryable errors fail immediately
- Error messages are appropriate for audience
- Sensitive data is not logged
- Recovery mechanisms work correctly
