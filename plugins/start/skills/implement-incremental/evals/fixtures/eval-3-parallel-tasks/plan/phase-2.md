---
title: "Phase 2: API and Adapters"
status: pending
version: "1.0"
phase: 2
---

# Phase 2: API and Adapters

## Phase Context

**GATE**: Read all referenced files before starting this phase.
**Specification References**: `[ref: solution/Section 4]`

## Tasks

- [ ] **T2.1 REST controller** `[activity: api-development]` `[parallel: true]`

  1. Prime: Read solution.md API section
  2. Test: POST returns 201, validation returns 422
  3. Implement: Create `src/controllers/orders.ts`
  4. Validate: API contract tests pass
  5. Success: AC-1 verified `[ref: requirements/AC-1]`

- [ ] **T2.2 Stripe adapter** `[activity: integration]` `[parallel: true]`

  1. Prime: Read solution.md adapters section
  2. Test: Charge created with correct amount, webhook validated
  3. Implement: Create `src/adapters/stripe.ts`
  4. Validate: Integration tests pass
  5. Success: AC-2 verified `[ref: requirements/AC-2]`

- [ ] **T2.3 Phase Validation** `[activity: validate]`

  Run all phase 2 tests, lint, typecheck. Verifies parallel-task outputs cohere.
