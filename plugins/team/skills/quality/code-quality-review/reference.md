# Code Quality Review Reference

Expanded criteria, specific signals, and severity guidance for each review dimension. Load this when a review requires deeper analysis than the dimension tables in SKILL.md provide.

---

## Dimension Deep Dives

### Readability

Readable code is understood correctly on first read. The goal is not brevity — it is clarity.

#### Naming

Good names eliminate the need to trace code to understand it.

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Cryptic abbreviations | `usr`, `d`, `tmp`, `val`, `res` as standalone names | MEDIUM |
| Misleading names | `isEnabled` that returns a count; `getUser` that saves | HIGH |
| Generic names | `data`, `info`, `obj`, `result`, `item` in non-trivial code | LOW |
| Boolean non-predicates | `user.active` instead of `user.isActive` | LOW |
| Inconsistent vocabulary | `fetch` in one function, `get`, `retrieve`, `load` for same concept elsewhere | MEDIUM |
| Noun-verb confusion | Functions named as nouns (`userSave()`); variables named as verbs | MEDIUM |

Specific names to scrutinize:

- Loop variables: `i`, `j` are acceptable for simple index loops; not acceptable as meaningful identifiers
- Boolean parameters: `createUser(true)` — what does `true` mean? Should be `createUser({ sendWelcomeEmail: true })`
- Return value names: a function returning a filtered list should not store the result in `list2`

#### Comments

The right comment explains why code does something, not what it does. The code itself should explain the what.

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Redundant comments | `// increment counter` above `count++` | LOW |
| Stale comments | Comment describes behavior that code no longer implements | HIGH |
| Missing intent comments | Complex algorithm, regex, or workaround with no explanation | MEDIUM |
| Commented-out code | Dead code left in place — creates confusion about intent | MEDIUM |
| TODO without ticket | `// TODO: fix this later` with no owner or reference | LOW |

Comments that add value:

```
// Stripe requires amounts in cents, not dollars
// See: https://stripe.com/docs/currencies#zero-decimal
const amount = Math.round(price * 100);

// We skip soft-deleted records here because the report
// only counts billable events — deleted users are not billed.
const events = await Event.where({ deletedAt: null });
```

#### Complexity

Cyclomatic complexity above 10 is a warning. Above 15 is a blocker.

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Nesting depth > 3 | Three or more levels of if/for/while nesting | MEDIUM |
| Boolean explosion | Conditions with 4+ AND/OR clauses | MEDIUM |
| Negated negatives | `if (!isNotActive)` — double negation obscures intent | MEDIUM |
| Long methods | Functions exceeding 20 lines doing multiple things | MEDIUM |
| Flag parameters | `processOrder(order, true, false, true)` — positional booleans | HIGH |

Measuring complexity in practice:

- Count the number of independent code paths through a function
- Each `if`, `else if`, `for`, `while`, `case`, `&&`, `||` adds one path
- A function with complexity of 8 needs 8 distinct test cases to cover all paths

---

### Maintainability

Maintainable code is easy to change safely. The test is: can a developer who did not write this code modify it confidently six months later?

#### Duplication

DRY violations increase the cost of every future change and create drift between copies.

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Copy-paste blocks | Identical or near-identical logic in 2+ places | MEDIUM |
| Repeated conditionals | Same `if (user.role === 'admin')` check in 5 places | HIGH |
| Structural duplication | Different data, same shape — suggests a missing abstraction | MEDIUM |
| Literal duplication | Same string constant typed in multiple files | LOW |

The threshold: duplicated logic that appears twice is a candidate for extraction. Appearing three or more times is a requirement to extract.

Exception: premature deduplication that requires complex parameterization to handle slight variations is worse than duplication. Evaluate whether extraction actually simplifies.

#### Coupling

Tightly coupled code breaks in unexpected places when changed.

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Feature Envy | A method references another class's fields more than its own | MEDIUM |
| Inappropriate Intimacy | Class A accesses private internals of Class B | HIGH |
| Deep import chains | `import { x } from '../../../../core/utils/helpers/string'` | LOW |
| Circular dependency | Module A imports from B, B imports from A | HIGH |
| Law of Demeter violations | `order.customer.address.city` — chained navigation | MEDIUM |

#### Cohesion

A class or module should have one reason to change. When it has multiple, changes in one area risk breaking another.

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| God Object | Class with 10+ public methods spanning unrelated concerns | HIGH |
| Utility dumping ground | `utils.ts` or `helpers.py` growing without domain organization | MEDIUM |
| Mixed abstraction levels | High-level orchestration mixed with low-level string parsing | MEDIUM |
| Unrelated exports | Module exporting types for two different bounded contexts | MEDIUM |

---

### Testability

If code is hard to test, it is usually hard to understand, change, or reason about. Poor testability is a design signal, not just a testing concern.

#### Dependency Injection

Code that instantiates its own dependencies is hard to test in isolation.

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Hard-coded instantiation | `const db = new Database()` inside a service constructor | HIGH |
| Static method calls | `UserService.getCurrentUser()` called deep inside business logic | HIGH |
| Global state access | Direct access to `process.env`, singleton registries inside functions | MEDIUM |
| Date/time coupling | `new Date()` or `Date.now()` called inside functions that need deterministic tests | MEDIUM |

The test: can you run a unit test for this function without starting a database, making a network call, or reading the filesystem? If no, the dependencies are not properly injected.

#### Observability

Code needs to expose enough surface area to write meaningful assertions.

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Hidden side effects | Function sends email, writes file, or mutates global state without returning a testable signal | HIGH |
| Void returns on complex logic | Functions doing significant work but returning nothing | MEDIUM |
| Non-determinism | Functions that depend on random values, current time, or external state without injection | MEDIUM |
| Private everything | Classes with no public interface except the final output | MEDIUM |

#### Test Coverage Gaps

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| No tests for new code | PR adds logic with zero corresponding test file changes | HIGH |
| Tests for implementation only | Tests that break on every refactor without logic changes | MEDIUM |
| Happy path only | Tests that never simulate errors, empty inputs, or boundary values | HIGH |
| Missing boundary tests | Off-by-one errors, empty collections, zero values, maximum values | MEDIUM |
| Integration tests masking unit gaps | Every test hits the database; no unit-level isolation | MEDIUM |

---

### Error Handling

Error handling is part of correctness. Code that fails silently or loses context on failure is as broken as code with wrong logic.

#### Failure Coverage

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Swallowed exceptions | `catch (e) { /* ignore */ }` or `catch (e) { return null; }` | HIGH |
| Generic catch-all | Catching `Exception` or `Error` instead of specific types | MEDIUM |
| Missing error propagation | Callers cannot distinguish success from failure | HIGH |
| Optimistic code | Database calls, API calls, file reads with no error handling | HIGH |

#### Error Quality

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Context-free messages | `throw new Error('invalid input')` — which field? what value? | MEDIUM |
| Exposing internals | Stack traces or SQL errors returned to API callers | CRITICAL |
| Type information lost | Re-throwing as a different error type without preserving original | MEDIUM |
| Logging without acting | `console.error(e)` followed by normal code execution | HIGH |

What good error handling looks like:

```
// Specific error type
// Context in the message (what failed, what was expected)
// Original error preserved for debugging
throw new ValidationError(
  `User age must be between 18 and 120, got ${age}`,
  { field: 'age', received: age, min: 18, max: 120 }
);
```

#### Null and Undefined Safety

| Signal | What to Look For | Severity |
|--------|------------------|----------|
| Unchecked optional chaining | Accessing `.property` on a value that could be null/undefined | HIGH |
| Missing null guards at boundaries | Data from external APIs, user input, or database queries used without null check | HIGH |
| Implicit truthiness checks | `if (user)` when the check should be `if (user !== null)` | LOW |
| Nullable return not handled | Function documented to return null used without null check at call site | HIGH |

---

### Naming Specifics

Naming deserves its own expanded section because it is the most frequent source of LOW and MEDIUM findings.

#### Functions and Methods

| Pattern | Problem | Better |
|---------|---------|--------|
| `getData()` | What data? Where from? | `fetchUserProfileFromCache()` |
| `process()` | Process what? How? | `normalizeIncomingWebhookPayload()` |
| `handle()` | Vague event handler | `handlePaymentFailedEvent()` |
| `check()` | Returns bool? Throws? Logs? | `validateEmailFormat()` or `assertEmailIsValid()` |
| `update()` | Updates one field? All fields? | `updateUserEmailAddress()` |
| `calculate()` | Calculate and return? Side effects? | `computeOrderSubtotal()` |

#### Variables

| Pattern | Problem | Better |
|---------|---------|--------|
| `flag` | Flag for what? | `isEmailVerified` |
| `list` | List of what? | `pendingOrderIds` |
| `count` | Count of what? | `failedLoginAttemptCount` |
| `temp` | Temporary what? | `intermediateCalculationResult` (or extract a function) |
| `data` | Data from where? | `rawApiResponse` or `parsedUserRecord` |
| `config` | Config for what? | `databaseConnectionConfig` |

#### Boolean Naming

Booleans should always be readable as a yes/no question:

- `isActive` not `active`
- `hasPermission` not `permission`
- `canEdit` not `editable`
- `shouldRetry` not `retry`
- `wasDeleted` not `deleted`

---

## Severity Reference

The SKILL.md severity matrix covers broad categories. This section provides finer guidance for borderline cases.

### CRITICAL

Reserve CRITICAL for findings that represent immediate risk if the code ships:

- Any code path that could expose secrets, credentials, or PII
- SQL/command/script injection without sanitization
- Missing authentication on endpoints that modify data
- Logic that can corrupt or permanently destroy data
- Breaking changes to a public API or shared contract without versioning

Do not use CRITICAL for style issues, even severe ones.

### HIGH

Use HIGH for issues that will cause problems in production — not might, but will, given normal usage:

- Logic error affecting the stated purpose of the code
- Missing error handling for a failure mode that occurs in the real world (network timeout, disk full, invalid user input)
- N+1 queries in endpoints that handle realistic data volumes
- Race condition in code that runs concurrently
- Architectural violation that will require significant rework later

### MEDIUM

Use MEDIUM for issues that reduce quality, increase risk over time, or add friction without causing immediate breakage:

- Code duplication that will drift
- Missing tests for new, non-trivial logic
- Naming that requires reading surrounding code to understand
- Cyclomatic complexity between 10-15
- Missing documentation for public APIs
- Hard-coded values that belong in configuration

### LOW

Use LOW for findings where the current code works correctly but could be improved:

- Style inconsistencies not caught by linters
- Mild optimization opportunities not in hot paths
- Naming that could be slightly clearer
- Comments that could be more precise
- Minor structural improvements with no behavioral impact

### Nitpick (not a severity level)

Anything caught by a configured linter should not appear as a review finding. Flag it once if the linter is misconfigured, then stop.

---

## Severity Escalation Rules

Some findings warrant escalating from their initial severity based on context:

| Situation | Escalation |
|-----------|-----------|
| MEDIUM issue in security-critical path | Escalate to HIGH |
| LOW naming issue on public API boundary | Escalate to MEDIUM (it will outlive this PR) |
| HIGH issue but only reached by authenticated admins | May stay HIGH, document the mitigating control |
| MEDIUM duplication in code with frequent change history | Escalate to HIGH (drift will occur) |
| Any finding in code with zero test coverage | Escalate one level (harder to catch regressions) |

---

## Dimension Interaction

Dimensions are not fully independent. Findings often span multiple:

| Finding | Primary Dimension | Secondary Dimensions |
|---------|-------------------|----------------------|
| God Object | Design | Testability, Maintainability |
| Hardcoded API key | Security | Correctness (will break in different environments) |
| N+1 query | Performance | Correctness (may timeout in production) |
| Swallowed exception | Correctness | Readability (hides failure signals) |
| Missing null check | Correctness | Readability (reader cannot tell if null is valid) |

When a finding spans dimensions, report it under the highest-severity dimension and note the secondary impact in the recommendation.
