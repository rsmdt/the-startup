# Severity & Confidence Classification

Definitions, classification matrix, and example findings at each severity level.

---

## Severity Levels

| Level | Definition | Action |
|-------|------------|--------|
| 🔴 **CRITICAL** | Security vulnerability, data loss risk, or system crash | **Must fix before merge** |
| 🟠 **HIGH** | Significant bug, performance issue, or breaking change | **Should fix before merge** |
| 🟡 **MEDIUM** | Code quality issue, maintainability concern, or missing test | **Consider fixing** |
| ⚪ **LOW** | Style preference, minor improvement, or suggestion | **Nice to have** |

## Confidence Levels

| Level | Definition | Usage |
|-------|------------|-------|
| **HIGH** | Clear violation of established pattern or security rule | Present as definite issue |
| **MEDIUM** | Likely issue but context-dependent | Present as probable concern |
| **LOW** | Potential improvement, may not be applicable | Present as suggestion |

## Classification Matrix

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

## Example Findings

### Critical Security Finding

```
[🔐 Security] **SQL Injection Vulnerability** (CRITICAL)
Location: `src/api/users.ts:45`
Confidence: HIGH
Issue: User input directly interpolated into SQL query
Fix: Use parameterized queries
```

```diff
- const result = db.query(`SELECT * FROM users WHERE id = ${req.params.id}`)
+ const result = db.query('SELECT * FROM users WHERE id = $1', [req.params.id])
```

### High Performance Finding

```
[⚡ Performance] **N+1 Query Pattern** (HIGH)
Location: `src/services/orders.ts:78-85`
Confidence: HIGH
Issue: Each order fetches its items in a separate query
Fix: Use eager loading or batch fetch
```

```diff
- const orders = await Order.findAll()
- for (const order of orders) {
-   order.items = await OrderItem.findByOrderId(order.id)
- }
+ const orders = await Order.findAll({ include: [OrderItem] })
```

### Medium Quality Finding

```
[📝 Quality] **Function Exceeds Recommended Length** (MEDIUM)
Location: `src/utils/validator.ts:23-89`
Confidence: HIGH
Issue: Function is 66 lines, exceeding 20-line recommendation
Fix: Extract validation logic into separate focused functions

Suggested breakdown:
- validateEmail() - lines 25-40
- validatePhone() - lines 42-55
- validateAddress() - lines 57-85
```

### Low Suggestion

```
[🧪 Testing] **Edge Case Not Tested** (LOW)
Location: `src/utils/date.ts:12` (formatDate function)
Confidence: MEDIUM
Issue: No test for invalid date input
Fix: Add test case for null/undefined/invalid dates
```

```javascript
it('should handle invalid date input', () => {
  expect(formatDate(null)).toBe('')
  expect(formatDate('invalid')).toBe('')
})
```
