# Error Handling Patterns

Standard error handling approaches that all agents should recommend.

---

## Error Classification

Distinguish between operational and programmer errors:

| Type | Examples | Response |
|------|----------|----------|
| **Operational** | Network failures, invalid input, timeouts, rate limits | Handle gracefully, log appropriately, provide user feedback, implement recovery |
| **Programmer** | Type errors, null access, failed assertions | Fail fast, log full context, alert developers - do NOT attempt recovery |

---

## Pattern 1: Fail Fast at Boundaries

Validate inputs at system boundaries and fail immediately with clear error messages. Do not allow invalid data to propagate through the system.

```javascript
// At API boundary
function handleRequest(input) {
  const validation = validateInput(input);
  if (!validation.valid) {
    throw new ValidationError(validation.errors);
  }
  // Process validated input
}
```

## Pattern 2: Specific Error Types

Create domain-specific error types that carry context about what failed and why. Generic errors lose valuable debugging information.

```javascript
class PaymentDeclinedError extends Error {
  constructor(reason, transactionId) {
    super(`Payment declined: ${reason}`);
    this.reason = reason;
    this.transactionId = transactionId;
  }
}
```

## Pattern 3: User-Safe Messages

Never expose internal error details to users. Log full context internally, present sanitized messages externally.

```javascript
try {
  await processPayment(order);
} catch (error) {
  logger.error('Payment failed', {
    error,
    orderId: order.id,
    userId: user.id
  });
  throw new UserFacingError('Payment could not be processed. Please try again.');
}
```

## Pattern 4: Graceful Degradation

When non-critical operations fail, degrade gracefully rather than failing entirely. Define what is critical vs. optional.

```javascript
async function loadDashboard() {
  const [userData, analytics, recommendations] = await Promise.allSettled([
    fetchUserData(),      // Critical - fail if missing
    fetchAnalytics(),     // Optional - show placeholder
    fetchRecommendations() // Optional - hide section
  ]);

  if (userData.status === 'rejected') {
    throw new Error('Cannot load dashboard');
  }

  return {
    user: userData.value,
    analytics: analytics.value ?? null,
    recommendations: recommendations.value ?? []
  };
}
```

## Pattern 5: Retry with Backoff

For transient failures (network, rate limits), implement exponential backoff with maximum attempts.

```javascript
async function fetchWithRetry(url, maxAttempts = 3) {
  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      return await fetch(url);
    } catch (error) {
      if (attempt === maxAttempts) throw error;
      await sleep(Math.pow(2, attempt) * 100); // 200ms, 400ms, 800ms
    }
  }
}
```

---

## Logging Levels

| Level | Use For |
|-------|---------|
| ERROR | Operational errors requiring attention |
| WARN | Recoverable issues, degraded performance |
| INFO | Significant state changes, request lifecycle |
| DEBUG | Detailed flow for troubleshooting |

**Log:** Correlation IDs, user context (sanitized), operation attempted, error type, duration.
**Never log:** Passwords, tokens, secrets, credit card numbers, PII.
