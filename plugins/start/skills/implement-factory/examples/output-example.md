# Example Factory Loop Output

## Manifest Discovery

Spec: 002-query-engine-enhancements
Feature: Query Engine Enhancements
Threshold: 90% | Max Iterations: 5

Units:
- dm1: Data models — no dependencies — pending
- ve1: Validation endpoint — after: dm1 — pending
- ve2: Input sanitization — after: ve1 — pending
- rl1: Rate limiting — after: ve1 — pending
- md1: Metrics dashboard — after: dm1 — pending

Execution Order:
- Group 1 (parallel): dm1
- Group 2 (parallel): ve1, md1
- Group 3 (sequential): ve2
- Group 4 (sequential): rl1

Service: `npm start` on port 3000
Starting from Group 1.

---

## Unit Result — Passed

dm1: Data models (iteration 1 of 5)

Satisfaction: 3/3 scenarios (100%)
Threshold: 90%

All scenarios passed. Unit complete.

Files: src/models/query.ts, src/models/result.ts, src/migrations/001-query-tables.ts
Tests: 14 passing

Manifest updated: `- [x] dm1: Data models`

---

## Unit Result — Failed, Queuing Retry

ve1: Validation endpoint (iteration 1 of 5)

Satisfaction: 3/5 scenarios (60%)
Threshold: 90%

Passed:
- Valid SQL accepted: 3/3
- Column extraction: 3/3
- Complexity score: 2/3

Failed:
- SQL injection detection: endpoint returned 500 instead of 400
- Empty input handling: no validation response

Queuing retry (iteration 2). Failure summaries extracted for code agent.

---

## Unit Result — Passed on Retry

ve1: Validation endpoint (iteration 2 of 5)

Satisfaction: 5/5 scenarios (100%)
Threshold: 90%

All scenarios passed. Unit complete.

Files: src/controllers/validate.ts, src/services/sql-parser.ts, src/middleware/input-guard.ts
Tests: 22 passing

Manifest updated: `- [x] ve1: Validation endpoint`

---

## Unit Result — Max Iterations Reached

rl1: Rate limiting (iteration 5 of 5)

Satisfaction: 2/3 scenarios (67%)
Threshold: 90%

Failed:
- Burst traffic throttling: requests not rejected after limit exceeded

Max iterations reached. Escalating to user.
Options: Retry with guidance | Skip unit | Abort

---

## Group Summary

Group 2 Complete (parallel): ve1, md1

Units: 2/2 completed
- ve1: 100% satisfaction (2 iterations)
- md1: 100% satisfaction (1 iteration)

Total iterations: 3
Files changed: 8 files
Constitution: No violations

---

## Completion Summary

Implementation Complete

Spec: 002-query-engine-enhancements
Units Completed: 4/5
- dm1: 100% (1 iteration)
- ve1: 100% (2 iterations)
- ve2: 100% (1 iteration)
- md1: 100% (1 iteration)
- rl1: 67% (5 iterations, skipped)

Total Iterations: 10
Files Changed: 24 files (+1,204 -89)

Manifest status: `status: failed` (rl1 did not meet threshold)

---

## Paused Summary (when user chooses Abort)

Implementation Paused

Spec: 002-query-engine-enhancements
Units Completed: 3/5

Resume with: `/start:implement 002`
The factory loop will pick up from the next pending unit automatically.
