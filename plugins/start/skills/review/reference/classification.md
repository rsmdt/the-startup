# Severity, Confidence & the Consequence Test

Definitions, ranking rules, and example findings. Applies to both scoped reviews and audits.

---

## The Consequence Test (gate for every finding)

A finding is valid only if it names the concrete future change, bug class, or comprehension cost it causes.

| Not a finding | A finding |
|---------------|-----------|
| "This violates SRP" | "SRP violation at billing/invoice.ts:120 — adding a payment provider forces edits in four unrelated modules" |
| "Function is too long" | "90-line function at sync/run.ts:40 mixes retry policy with I/O — retry changes require re-testing the network layer" |
| "Magic number" | "Hardcoded 24h expiry at auth/config.ts:12 — ops cannot tune session length without a deploy" |

If no victim can be named, the observation goes to **CONTEXTUALLY ACCEPTABLE** with one line of reasoning — not to findings, not silently dropped. Principles (SOLID, DRY, etc.) may explain a finding, never constitute one.

---

## Severity Levels

| Level | Definition | Action |
|-------|------------|--------|
| 🔴 **CRITICAL** | Security vulnerability, data loss risk, or system crash | **Must fix before merge** |
| 🟠 **HIGH** | Significant bug, performance issue, or breaking change | **Should fix before merge** |
| 🟡 **MEDIUM** | Code quality issue, maintainability concern, or missing test | **Consider fixing** |
| ⚪ **LOW** | Style preference, minor improvement, or suggestion | **Nice to have** |

Breaking changes are orthogonal to severity: set the finding's `breaking` field and let severity describe the defect itself.
A clean, well-written API rename is not a bug, but it still forces REQUEST CHANGES and still leads the report.

## Confidence Levels

| Level | Definition | Usage |
|-------|------------|-------|
| **HIGH** | Clear violation of established pattern or security rule | Present as definite issue |
| **MEDIUM** | Likely issue but context-dependent | Present as probable concern |
| **LOW** | Potential improvement, may not be applicable | Present as suggestion |

---

## Ranking

**Scoped mode**: sort by severity descending, then confidence descending.

**Audit mode**: rank by likelihood × cost.

| | Cost when it materializes: LOW | MEDIUM | HIGH |
|---|---|---|---|
| **Likelihood HIGH** (touched monthly, or on any feature in a hot area) | MEDIUM | HIGH | CRITICAL |
| **Likelihood MEDIUM** (touched on a plausible roadmap item) | LOW | MEDIUM | HIGH |
| **Likelihood LOW** (only on rare events: vendor swap, rewrite of a stable area) | Accepted | LOW | MEDIUM |

Likelihood comes from the co-change data and orientation map (is the smell in a volatile or stable area?). Cost is the amplification measurement (files touched, bug class severity).

---

## Classification Matrix (scoped mode anchors)

| Finding Type | Severity | Confidence |
|--------------|----------|------------|
| SQL Injection | CRITICAL | HIGH |
| XSS Vulnerability | CRITICAL | HIGH |
| Hardcoded Secret | CRITICAL | HIGH |
| Missing Auth Check | CRITICAL | MEDIUM |
| N+1 Query | HIGH | HIGH |
| No Input Validation | MEDIUM | HIGH |
| Missing Test | MEDIUM | MEDIUM |
| Long Function (with named victim) | LOW–MEDIUM | HIGH |
| Long Function (no victim) | → Contextually Acceptable | — |

---

## Example Findings

### Critical Security Finding

```
[🔐 Security] **SQL Injection Vulnerability** (CRITICAL, HIGH confidence)
Location: src/api/users.ts:45
Consequence: any caller can read or destroy the users table via crafted id
Fix: parameterized query
```

```diff
- const result = db.query(`SELECT * FROM users WHERE id = ${req.params.id}`)
+ const result = db.query('SELECT * FROM users WHERE id = $1', [req.params.id])
```

### Medium Quality Finding (passes consequence test)

```
[📝 Quality] **Policy mixed with I/O** (MEDIUM, HIGH confidence)
Location: src/services/export.ts:23-89
Consequence: pricing-rule changes cannot be unit-tested without mocking S3;
             the team has 3 open tickets touching pricing rules this quarter
Fix: extract computePrice() pure function; keep upload in place
```

### Contextually Acceptable (correctly NOT a finding)

```
ACCEPTED: 7-parameter constructor in test fixture builder (tests/factory.ts:12)
Reasoning: builder is only called by tests via named helpers; no production
caller exists to suffer the parameter-ordering bug class.
```
