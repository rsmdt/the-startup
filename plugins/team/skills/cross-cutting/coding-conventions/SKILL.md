---
name: coding-conventions
description: Apply consistent security, performance, and accessibility standards across all recommendations. Use when reviewing code, designing features, or validating implementations. Cross-cutting skill for all agents.
---

# Best Practices

A cross-cutting skill that enforces consistent security, performance, and quality standards across all agent recommendations. This skill provides actionable checklists aligned with industry standards.

## When to Use

- Reviewing code for security vulnerabilities
- Validating performance characteristics of implementations
- Ensuring accessibility compliance in UI components
- Designing error handling strategies
- Auditing existing systems for quality gaps

## Core Domains

### Security

Covers common vulnerability prevention aligned with OWASP Top 10. Apply these checks to any code that handles user input, authentication, data storage, or external communications.

See: `checklists/security-checklist.md`

### Performance

Covers optimization patterns for frontend, backend, and database operations. Apply these checks when performance is a concern or during code review.

See: `checklists/performance-checklist.md`

### Accessibility

Covers WCAG 2.1 Level AA compliance. Apply these checks to all user-facing components to ensure inclusive design.

See: `checklists/accessibility-checklist.md`

## Error Handling Patterns

All agents should recommend these error handling approaches:

### Pattern 1: Fail Fast at Boundaries

Validate inputs at system boundaries and fail immediately with clear error messages. Do not allow invalid data to propagate through the system.

```
// At API boundary
function handleRequest(input) {
  const validation = validateInput(input);
  if (!validation.valid) {
    throw new ValidationError(validation.errors);
  }
  // Process validated input
}
```

### Pattern 2: Specific Error Types

Create domain-specific error types that carry context about what failed and why. Generic errors lose valuable debugging information.

```
class PaymentDeclinedError extends Error {
  constructor(reason, transactionId) {
    super(`Payment declined: ${reason}`);
    this.reason = reason;
    this.transactionId = transactionId;
  }
}
```

### Pattern 3: User-Safe Messages

Never expose internal error details to users. Log full context internally, present sanitized messages externally.

```
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

### Pattern 4: Graceful Degradation

When non-critical operations fail, degrade gracefully rather than failing entirely. Define what is critical vs. optional.

```
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

### Pattern 5: Retry with Backoff

For transient failures (network, rate limits), implement exponential backoff with maximum attempts.

```
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

## Best Practices

- Apply security checks before performance optimization
- Make accessibility a default, not an afterthought
- Use checklists during code review, not just at the end
- Document exceptions to standards with rationale
- Automate checks where possible (linting, testing)

## References

- `checklists/security-checklist.md` - OWASP-aligned security checks
- `checklists/performance-checklist.md` - Performance optimization checklist
- `checklists/accessibility-checklist.md` - WCAG 2.1 AA compliance checklist
