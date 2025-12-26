---
name: error-recovery
description: Consistent error patterns, validation approaches, and recovery strategies. Use when implementing input validation, designing error responses, handling failures gracefully, or establishing logging practices. Covers operational vs programmer errors, user-facing vs internal errors, and recovery mechanisms.
---

# Error Handling

## When to Use

- Implementing input validation at system boundaries
- Designing error responses for APIs or user interfaces
- Building recovery mechanisms for transient failures
- Establishing logging and monitoring patterns
- Distinguishing between errors that need user action vs system intervention

## Philosophy

Errors are not exceptional - they are expected. Good error handling treats errors as first-class citizens of the system design, not afterthoughts. The goal is to fail safely, provide actionable feedback, and enable recovery.

## Error Classification

### Operational Errors

Runtime problems that occur during normal operation. These are expected and must be handled.

**Characteristics:**
- External system failures (network, database, filesystem)
- Invalid user input
- Resource exhaustion (memory, disk, connections)
- Timeout conditions
- Rate limiting

**Response:** Handle gracefully, log appropriately, provide user feedback, implement recovery.

### Programmer Errors

Bugs in the code that should not happen if the code is correct.

**Characteristics:**
- Type errors caught at runtime
- Null/undefined access on required values
- Failed assertions on invariants
- Invalid internal state

**Response:** Fail fast, log full context, alert developers. Do not attempt recovery - fix the bug.

## Core Patterns

### Input Validation

Validate early, validate completely, provide specific feedback.

```
PATTERN: Fail-Fast Validation

1. Validate at system boundaries (API entry, user input, file reads)
2. Check all constraints before processing
3. Return ALL validation errors, not just the first one
4. Include field name, actual value (if safe), and expected format
5. Never trust data from external sources
```

**Validation checklist:**
- Required fields present
- Types correct
- Values within allowed ranges
- Formats match expectations (email, URL, date)
- Business rules satisfied

### Error Messages

Different audiences need different information.

**User-facing errors:**
- Clear action the user can take
- No technical jargon or stack traces
- Consistent tone and format
- Localization-ready

**Internal/logged errors:**
- Full technical context
- Request/correlation IDs
- Timestamp and service identifier
- Stack trace for programmer errors
- Sanitized sensitive data

### Recovery Strategies

**Retry with backoff:**
- Transient failures (network timeouts, rate limits)
- Exponential backoff with jitter
- Maximum retry count with circuit breaker

**Fallback:**
- Degraded functionality over complete failure
- Cached data when live data unavailable
- Default values when configuration missing

**Compensation:**
- Undo partial operations on failure
- Maintain consistency in distributed operations
- Saga pattern for multi-step processes

### Logging Levels

| Level | Use For |
|-------|---------|
| ERROR | Operational errors requiring attention |
| WARN | Recoverable issues, degraded performance |
| INFO | Significant state changes, request lifecycle |
| DEBUG | Detailed flow for troubleshooting |

**What to log:**
- Correlation/request ID
- User context (sanitized)
- Operation being attempted
- Error type and message
- Duration and timing

**What NOT to log:**
- Passwords, tokens, secrets
- Full credit card numbers
- Personal identifiable information (PII)
- Raw request/response bodies containing sensitive data

## Best Practices

- Fail fast on programmer errors - do not mask bugs
- Handle operational errors gracefully with recovery options
- Provide correlation IDs for tracing requests across services
- Use structured logging (JSON) for machine parseability
- Centralize error handling logic - avoid scattered try/catch blocks
- Test error paths as rigorously as success paths
- Monitor error rates and set alerts for anomalies

## Anti-Patterns

- Catching all exceptions silently (`catch {}`)
- Logging sensitive data in error messages
- Returning generic "Something went wrong" without context
- Retrying non-idempotent operations without safeguards
- Mixing validation errors with system errors in responses
- Treating all errors the same regardless of recoverability

## References

- [examples/error-patterns.md](examples/error-patterns.md) - Concrete examples across languages
