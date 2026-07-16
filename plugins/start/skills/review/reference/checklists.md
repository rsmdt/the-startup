# Review Checklists

Detailed checklists for each core review perspective (scoped mode).

A checklist hit is a lead, not a finding. Before reporting, apply the consequence test from classification.md: name the concrete bug class, future change, or comprehension cost at file:line. Hits without a nameable victim go to the Accepted list.

---

## Security Review Checklist

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

---

## Performance Review Checklist

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

---

## Quality Review Checklist

**Code Structure:**
- [ ] Single responsibility principle (report only with a named change that gets harder)
- [ ] Functions are focused — mixed abstraction levels or policy tangled with I/O
- [ ] No deep nesting (< 4 levels)
- [ ] DRY - no duplicated logic that must co-evolve (identical-looking but independently-changing code is fine)

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

---

## Test Coverage Checklist

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
