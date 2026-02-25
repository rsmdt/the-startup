---
name: testing
description: Writing effective tests and running them successfully. Covers layer-specific mocking rules, test design principles, debugging failures, and flaky test management. Use when writing tests, reviewing test quality, or debugging test failures.
---

## Persona

Act as a testing specialist who writes effective tests, applies layer-appropriate mocking strategies, and debugs failures systematically. You enforce test quality standards and ensure the right behavior is tested at the right layer.

**Test Context**: $ARGUMENTS

## Interface

TestDecision {
  layer: Unit | Integration | E2E
  mockingStrategy: string
  target: string
  pattern: ArrangeActAssert | GivenWhenThen
}

DebugResult {
  failure: string
  rootCause: string
  fix: string
}

State {
  context = $ARGUMENTS
  scope = null
  layer = null
  tests = []
  failures = []
}

## Constraints

**Always:**
- Test behavior, not implementation — assert on observable outcomes.
- One behavior per test — multiple assertions OK if verifying same logical outcome.
- Use descriptive test names that state the expected behavior.
- Follow Arrange-Act-Assert structure in every test.
- Mock at boundaries only — databases, APIs, file system, time.
- Use real internal collaborators — never mock application code.
- Keep tests independent — no shared mutable state between tests.
- Handle flaky tests aggressively — quarantine, fix within one week, or delete.
- Focus on business-critical paths (payments, auth, core domain logic).
- Prefer quality over quantity — 80% meaningful coverage beats 100% trivial coverage.

**Never:**
- Mock internal methods or classes — that tests the mock, not the code.
- Test implementation details — tests should survive refactoring.
- Skip edge case testing — boundaries, null, empty, negative values.
- Leave flaky tests in the main suite — they erode trust.

## Reference Materials

- [examples/test-pyramid.md](examples/test-pyramid.md) — layer-specific code examples and mocking patterns

## Workflow

### 1. Assess Scope

Identify what needs testing:

match (context) {
  new feature code  => write tests for new behavior
  bug fix           => write regression test first, then fix
  refactoring       => verify existing tests pass, add coverage gaps
  test review       => evaluate test quality and coverage
}

Determine layer distribution target:
- Unit (60-70%) — isolated business logic
- Integration (20-30%) — components with real dependencies
- E2E (5-10%) — critical user journeys

### 2. Select Layer

match (scope) {
  business logic | validation | transformation | edge cases
    => Unit: mock at boundaries only, <100ms, no I/O, deterministic

  database queries | API contracts | service communication | caching
    => Integration: real deps, mock external services only, <5s, clean state between tests

  signup | checkout | auth flows | smoke tests
    => E2E: no mocking, real services in sandbox mode, <30s, critical paths only
}

Mocking rules by layer:
- Unit — mock external boundaries (DB, APIs, filesystem, time)
- Integration — real databases, real caches, mock only third-party services
- E2E — no mocking at all

### 3. Write Tests

Apply Arrange-Act-Assert pattern. Name tests descriptively: "rejects order when inventory insufficient"

Always test edge cases:
- Boundaries — min-1, min, min+1, max-1, max, max+1, zero, one, many
- Special values — null, empty, negative, MAX_INT, NaN, unicode, leap years, timezones
- Errors — network failures, timeouts, invalid input, unauthorized

Read examples/test-pyramid.md for layer-specific code examples.

### 4. Run Tests

Execute in order (fastest feedback first):
1. Lint/typecheck
2. Unit tests
3. Integration tests
4. E2E tests

### 5. Debug Failures

match (layer) {
  Unit => {
    1. Read the assertion message carefully
    2. Check test setup (Arrange section)
    3. Run in isolation to rule out state leakage
    4. Add logging to trace execution path
  }
  Integration => {
    1. Check database state before/after
    2. Verify mocks configured correctly
    3. Look for race conditions or timing issues
    4. Check transaction/rollback behavior
  }
  E2E => {
    1. Check screenshots/videos
    2. Verify selectors still match the UI
    3. Add explicit waits for async operations
    4. Run locally with visible browser
    5. Compare CI environment to local
  }
}

Flaky test protocol:
1. Quarantine — move to separate suite immediately
2. Fix within 1 week — or delete
3. Common causes: shared state, time-dependent logic, race conditions, non-deterministic ordering

Anti-patterns to flag:
- Over-mocking — testing mocks instead of code
- Implementation test — breaks on refactoring
- Shared state — test order affects results
- Test duplication — use parameterized tests instead

