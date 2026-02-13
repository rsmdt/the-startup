# Code Review Reference

Detailed checklists, classification matrices, and agent prompts for the `/start:review` skill.

---

## Review Checklists

### Security Review Checklist

**Authentication & Authorization:**
- [ ] Proper auth checks before sensitive operations
- [ ] No privilege escalation vulnerabilities
- [ ] Session management is secure

**Injection Prevention:**
- [ ] SQL queries use parameterized statements
- [ ] XSS prevention (output encoding)
- [ ] Command injection prevention (input validation)

**Data Protection:**
- [ ] No hardcoded secrets or credentials
- [ ] Sensitive data properly encrypted
- [ ] PII handled according to policy

**Input Validation:**
- [ ] All user inputs validated
- [ ] Proper sanitization before use
- [ ] Safe deserialization practices

### Performance Review Checklist

**Database Operations:**
- [ ] No N+1 query patterns
- [ ] Efficient use of indexes
- [ ] Proper pagination for large datasets
- [ ] Connection pooling in place

**Computation:**
- [ ] Efficient algorithms (no O(n²) when O(n) possible)
- [ ] Proper caching for expensive operations
- [ ] No unnecessary recomputations

**Resource Management:**
- [ ] No memory leaks
- [ ] Proper cleanup of resources
- [ ] Async operations where appropriate
- [ ] No blocking operations in event loops

### Quality Review Checklist

**Code Structure:**
- [ ] Single responsibility principle
- [ ] Functions are focused (< 20 lines ideal)
- [ ] No deep nesting (< 4 levels)
- [ ] DRY - no duplicated logic

**Naming & Clarity:**
- [ ] Intention-revealing names
- [ ] Consistent terminology
- [ ] Self-documenting code
- [ ] Comments explain "why", not "what"

**Error Handling:**
- [ ] Errors handled at appropriate level
- [ ] Specific error messages
- [ ] No swallowed exceptions
- [ ] Proper error propagation

**Project Standards:**
- [ ] Follows coding conventions
- [ ] Consistent with existing patterns
- [ ] Proper file organization
- [ ] Type safety (if applicable)

### Test Coverage Checklist

**Coverage:**
- [ ] Happy path tested
- [ ] Error cases tested
- [ ] Edge cases tested
- [ ] Boundary conditions tested

**Test Quality:**
- [ ] Tests are independent
- [ ] Tests are deterministic (not flaky)
- [ ] Proper assertions (not just "no error")
- [ ] Mocking at appropriate boundaries

**Test Organization:**
- [ ] Tests match code structure
- [ ] Clear test names
- [ ] Proper setup/teardown
- [ ] Integration tests where needed

---

## Severity & Confidence Classification

### Severity Levels

| Level | Definition | Action |
|-------|------------|--------|
| 🔴 **CRITICAL** | Security vulnerability, data loss risk, or system crash | **Must fix before merge** |
| 🟠 **HIGH** | Significant bug, performance issue, or breaking change | **Should fix before merge** |
| 🟡 **MEDIUM** | Code quality issue, maintainability concern, or missing test | **Consider fixing** |
| ⚪ **LOW** | Style preference, minor improvement, or suggestion | **Nice to have** |

### Confidence Levels

| Level | Definition | Usage |
|-------|------------|-------|
| **HIGH** | Clear violation of established pattern or security rule | Present as definite issue |
| **MEDIUM** | Likely issue but context-dependent | Present as probable concern |
| **LOW** | Potential improvement, may not be applicable | Present as suggestion |

### Classification Matrix

| Finding Type | Severity | Confidence | Priority |
|--------------|----------|------------|----------|
| SQL Injection | CRITICAL | HIGH | Immediate |
| XSS Vulnerability | CRITICAL | HIGH | Immediate |
| Hardcoded Secret | CRITICAL | HIGH | Immediate |
| N+1 Query | HIGH | HIGH | Before merge |
| Missing Auth Check | CRITICAL | MEDIUM | Before merge |
| No Input Validation | MEDIUM | HIGH | Should fix |
| Long Function | LOW | HIGH | Nice to have |
| Missing Test | MEDIUM | MEDIUM | Should fix |

---

## Agent Prompt Templates

### Security Reviewer

```
FOCUS: Security review of the provided code changes
  - Identify authentication/authorization issues
  - Check for injection vulnerabilities (SQL, XSS, command, LDAP)
  - Look for hardcoded secrets or credentials
  - Verify input validation and sanitization
  - Check for insecure data handling (encryption, PII)
  - Review session management
  - Check for CSRF vulnerabilities in forms

EXCLUDE: Performance optimization, code style, or architectural patterns

CONTEXT:
  - Files changed: [list]
  - Changes: [the diff or code]
  - Full file context: [surrounding code]

OUTPUT: Security findings in this format:
  FINDING:
  - severity: CRITICAL | HIGH | MEDIUM | LOW
  - confidence: HIGH | MEDIUM | LOW
  - title: Brief title (max 40 chars)
  - location: file:line
  - issue: One sentence describing what's wrong
  - fix: Actionable recommendation
  - code_example: (Optional, for CRITICAL/HIGH)

SUCCESS: All security concerns identified with remediation steps
TERMINATION: Analysis complete OR code context insufficient
```

### Performance Reviewer

```
FOCUS: Performance review of the provided code changes
  - Identify N+1 query patterns
  - Check for unnecessary re-renders or recomputations
  - Look for blocking operations in async code
  - Identify memory leaks or resource cleanup issues
  - Check algorithm complexity (avoid O(n²) when O(n) possible)
  - Review caching opportunities
  - Check for proper pagination

EXCLUDE: Security vulnerabilities, code style, or naming conventions

CONTEXT:
  - Files changed: [list]
  - Changes: [the diff or code]
  - Full file context: [surrounding code]

OUTPUT: Performance findings in this format:
  FINDING:
  - severity: CRITICAL | HIGH | MEDIUM | LOW
  - confidence: HIGH | MEDIUM | LOW
  - title: Brief title (max 40 chars)
  - location: file:line
  - issue: One sentence describing what's wrong
  - fix: Optimization strategy

SUCCESS: All performance concerns identified with optimization strategies
TERMINATION: Analysis complete OR code context insufficient
```

### Quality Reviewer

```
FOCUS: Code quality review of the provided code changes
  - Check adherence to project coding standards
  - Identify code smells (long methods, duplication, complexity)
  - Verify proper error handling
  - Check naming conventions and code clarity
  - Identify missing or inadequate documentation
  - Verify consistent patterns with existing codebase
  - Check for proper abstractions

EXCLUDE: Security vulnerabilities or performance optimization

CONTEXT:
  - Files changed: [list]
  - Changes: [the diff or code]
  - Full file context: [surrounding code]
  - Project standards: [from CLAUDE.md, .editorconfig]

OUTPUT: Quality findings in this format:
  FINDING:
  - severity: CRITICAL | HIGH | MEDIUM | LOW
  - confidence: HIGH | MEDIUM | LOW
  - title: Brief title (max 40 chars)
  - location: file:line
  - issue: One sentence describing what's wrong
  - fix: Improvement suggestion

SUCCESS: All quality concerns identified with clear improvements
TERMINATION: Analysis complete OR code context insufficient
```

### Test Coverage Reviewer

```
FOCUS: Test coverage review of the provided code changes
  - Identify new code paths that need tests
  - Check if existing tests cover the changes
  - Look for test quality issues (flaky, incomplete assertions)
  - Verify edge cases are covered
  - Check for proper mocking at boundaries
  - Identify integration test needs
  - Verify test naming and organization

EXCLUDE: Implementation details not related to testing

CONTEXT:
  - Files changed: [list]
  - Changes: [the diff or code]
  - Full file context: [surrounding code]
  - Related test files: [existing tests]

OUTPUT: Test coverage findings in this format:
  FINDING:
  - severity: CRITICAL | HIGH | MEDIUM | LOW
  - confidence: HIGH | MEDIUM | LOW
  - title: Brief title (max 40 chars)
  - location: file:line
  - issue: One sentence describing what's wrong
  - fix: Suggested test case with code example

SUCCESS: All testing gaps identified with specific test recommendations
TERMINATION: Analysis complete OR code context insufficient
```

### Simplification Reviewer

```
FOCUS: Complexity review - aggressively challenge unnecessary complexity
  - Identify YAGNI violations (You Aren't Gonna Need It)
  - Find over-engineered solutions
  - Spot premature abstractions
  - Look for dead code paths
  - Challenge "clever" code that should be obvious
  - Find unnecessary indirection
  - Identify code that could be deleted

EXCLUDE: Security vulnerabilities or performance optimization

CONTEXT:
  - Files changed: [list]
  - Changes: [the diff or code]
  - Full file context: [surrounding code]

OUTPUT: Simplification findings in this format:
  FINDING:
  - severity: CRITICAL | HIGH | MEDIUM | LOW
  - confidence: HIGH | MEDIUM | LOW
  - title: Brief title (max 40 chars)
  - location: file:line
  - issue: Why this is more complex than needed
  - fix: Simpler alternative

SUCCESS: All complexity issues identified with simpler alternatives
TERMINATION: Analysis complete OR code context insufficient
```

---

## Synthesis Protocol

### Deduplication

If multiple agents flag the same issue:
1. Keep the finding with highest severity
2. Merge context from all agents
3. Note which perspectives flagged it

Example:
```
[🔐+⚡ Security/Performance] **Unvalidated User Input** (CRITICAL)
📍 Location: `src/api/search.ts:34`
🔍 Flagged by: Security Reviewer, Performance Reviewer
❌ Issue:
  - Security: Potential injection vulnerability
  - Performance: Unvalidated input could cause DoS
✅ Fix: Add input validation and length limits
```

### Grouping

Group findings for readability:
1. **By Severity** (Critical → Low) - default
2. **By File** (for file-focused reviews)
3. **By Category** (for category-focused reports)

---

## Example Findings

### Critical Security Finding

```
[🔐 Security] **SQL Injection Vulnerability** (CRITICAL)
📍 Location: `src/api/users.ts:45`
🔍 Confidence: HIGH
❌ Issue: User input directly interpolated into SQL query
✅ Fix: Use parameterized queries

```diff
- const result = db.query(`SELECT * FROM users WHERE id = ${req.params.id}`)
+ const result = db.query('SELECT * FROM users WHERE id = $1', [req.params.id])
```
```

### High Performance Finding

```
[⚡ Performance] **N+1 Query Pattern** (HIGH)
📍 Location: `src/services/orders.ts:78-85`
🔍 Confidence: HIGH
❌ Issue: Each order fetches its items in a separate query
✅ Fix: Use eager loading or batch fetch

```diff
- const orders = await Order.findAll()
- for (const order of orders) {
-   order.items = await OrderItem.findByOrderId(order.id)
- }
+ const orders = await Order.findAll({ include: [OrderItem] })
```
```

### Medium Quality Finding

```
[📝 Quality] **Function Exceeds Recommended Length** (MEDIUM)
📍 Location: `src/utils/validator.ts:23-89`
🔍 Confidence: HIGH
❌ Issue: Function is 66 lines, exceeding 20-line recommendation
✅ Fix: Extract validation logic into separate focused functions

Suggested breakdown:
- validateEmail() - lines 25-40
- validatePhone() - lines 42-55
- validateAddress() - lines 57-85
```

### Low Suggestion

```
[🧪 Testing] **Edge Case Not Tested** (LOW)
📍 Location: `src/utils/date.ts:12` (formatDate function)
🔍 Confidence: MEDIUM
❌ Issue: No test for invalid date input
✅ Fix: Add test case for null/undefined/invalid dates

```javascript
it('should handle invalid date input', () => {
  expect(formatDate(null)).toBe('')
  expect(formatDate('invalid')).toBe('')
})
```
```
