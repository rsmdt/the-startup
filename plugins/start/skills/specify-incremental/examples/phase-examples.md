# Implementation Phase Examples

Reference examples for structuring implementation phases as **vertical slices** with grouped task tracking.

## Core Principle

**Each phase is a vertical slice — one user-observable capability delivered end-to-end.** A phase touches whatever layers the behavior needs (domain → API → UI → test) and is validatable against a requirements.md acceptance criterion on its own. Phases are named after capabilities, never after layers.

**Track logical units that produce verifiable outcomes.** Each checkbox represents a deliverable, not a step in the process. The TDD cycle (Prime → Test → Implement → Validate) is embedded guidance within each task.

> The payment feature below is sliced into capabilities (*create*, *track status*, *refund*). Each phase delivers something a user can do, end-to-end. Contrast with horizontal slicing (a "schema phase", an "API phase", a "UI phase"), where nothing works — and no acceptance criterion can be verified — until the final phase.

---

## Example: First Capability (owns shared foundation)

```markdown
### Phase 1: Create a Payment

Delivers the end-to-end create-payment flow: a customer submits payment details and receives a confirmation. This is the first slice, so it builds the Payment entity and schema that later phases extend.

- [ ] **T1.1 Payment Entity & Schema** `[activity: domain-modeling]`

  **Prime**: Read payment interface contracts and validation rules `[ref: solution/Section 4.2; lines: 145-200]`

  **Test**: Entity rejects negative amounts; supports multiple currencies; persists and reloads via repository

  **Implement**: Create `src/domain/Payment.ts` with Amount value object; create `payments` table migration and `src/repositories/PaymentRepository.ts`

  **Validate**: Unit + repository integration tests pass; lint; typecheck

- [ ] **T1.2 Create Payment Endpoint** `[activity: backend-api]`

  **Prime**: Read API specification and authentication requirements `[ref: solution/Section 4.4]`

  **Test**: POST /payments creates a payment and returns 201; invalid amount returns 422 with details; charges Stripe in test mode

  **Implement**: Create `src/controllers/PaymentController.ts` (create route) wired to `StripePaymentAdapter`

  **Validate**: API contract tests pass; authentication enforced

- [ ] **T1.3 Payment Form** `[activity: frontend-ui]`

  **Prime**: Read UI specifications and form validation rules `[ref: solution/Section 5.1]`

  **Test**: Form validates card number format; submits on valid input; shows confirmation on 201; shows error states on 422

  **Implement**: Create `src/components/PaymentForm.tsx` with Stripe Elements, calling POST /payments

  **Validate**: Component tests pass; accessibility audit passes

- [ ] **T1.4 Create-Payment E2E** `[activity: e2e-test]`

  **Prime**: Review the create-payment user journey `[ref: requirements/Section 3.1]`

  **Test**: From form submission to confirmation page; declined card shows a recoverable error

  **Implement**: Wire form → endpoint → confirmation; add E2E spec

  **Validate**: E2E passes
  - Success: A customer can submit valid payment details and receive a confirmation `[ref: requirements/AC-3.1]`
```

**Task count**: 4 tracked items — one vertical slice spanning domain, API, UI, and E2E. Acceptance criterion AC-3.1 is verifiable at the end of *this* phase, not deferred.

---

## Example: Next Capability (extends Phase 1)

```markdown
### Phase 2: Track Payment Status

Delivers end-to-end payment-status visibility: a customer can see whether a payment is pending, settled, or failed. Depends on Phase 1; extends the Payment entity with status transitions.

- [ ] **T2.1 Status Transitions** `[activity: domain-modeling]` `[depends: Phase 1]`

  **Prime**: Read the payment status state machine `[ref: solution/Section 4.5]`

  **Test**: Valid transitions (pending→settled, pending→failed) succeed; invalid transitions are rejected; Stripe webhook updates status

  **Implement**: Extend `src/domain/Payment.ts` with status field + transition guards; add migration for the status column

  **Validate**: Unit + integration tests pass

- [ ] **T2.2 Status Endpoint** `[activity: backend-api]`

  **Prime**: Read the status API contract `[ref: solution/Section 4.4]`

  **Test**: GET /payments/:id returns current status or 404; reflects webhook-driven updates

  **Implement**: Add the get/status route to `PaymentController.ts`

  **Validate**: API contract tests pass

- [ ] **T2.3 Status Display & E2E** `[activity: frontend-ui]`

  **Prime**: Read the status UI spec `[ref: solution/Section 5.2]`

  **Test**: Status badge renders pending/settled/failed; polls and updates; E2E from payment to settled status

  **Implement**: Add `PaymentStatus.tsx`; wire to GET /payments/:id; add E2E spec

  **Validate**: Component + E2E tests pass
  - Success: A customer can view the live status of a payment `[ref: requirements/AC-3.3]`
```

**Task count**: 3 tracked items — a second vertical slice. It *extends* the entity from Phase 1 rather than waiting on a separate "schema phase".

---

## Example: Parallel-Eligible Capability (also extends Phase 1)

```markdown
### Phase 3: Refund a Payment

Delivers the end-to-end refund flow. Depends on Phase 1; independent of Phase 2, so the two can be sequenced in either order or worked in parallel by separate developers.

- [ ] **T3.1 Refund Logic** `[activity: domain-modeling]` `[depends: Phase 1]`

  **Prime**: Read refund rules and partial-refund handling `[ref: solution/Section 4.6]`

  **Test**: Full and partial refunds adjust balance correctly; over-refund is rejected; Stripe refund issued in test mode

  **Implement**: Extend `Payment.ts` with refund behavior; add `refund` method to `StripePaymentAdapter`

  **Validate**: Unit + integration tests pass

- [ ] **T3.2 Refund Endpoint & UI & E2E** `[activity: backend-api]`

  **Prime**: Read refund API + UI specs `[ref: solution/Section 4.4; Section 5.3]`

  **Test**: POST /payments/:id/refund returns 200 and updated balance; refund button disabled when ineligible; E2E from refund action to updated receipt

  **Implement**: Add refund route to `PaymentController.ts`; add refund control to the payment detail view; add E2E spec

  **Validate**: API + component + E2E tests pass
  - Success: A customer can refund a payment in full or in part `[ref: requirements/AC-3.2]`
```

**Task count**: 2 tracked items — a third vertical slice, parallel-eligible with Phase 2 because both depend only on Phase 1.

---

## Example: Final Cross-Cutting Phase

The one phase that is *not* a single capability. It hardens behavior already delivered and verified by earlier phases — it is not the first place any acceptance criterion is checked.

```markdown
### Phase 4: Hardening & Cross-Cutting Validation

Full-system validation across all delivered capabilities.

- [ ] **T4.1 Full Regression & Integration** `[activity: integration-test]`

  Run the complete suite across create/status/refund; verify Stripe webhook resilience (retries, signature validation, idempotency)

- [ ] **T4.2 Quality Gates** `[activity: validate]`

  Performance: API response < 200ms p95 `[ref: solution/Section 10]`; Security: input validation + auth verified; Coverage: > 80% line coverage

- [ ] **T4.3 Acceptance Re-Verification** `[activity: business-acceptance]`

  Confirm all requirements.md acceptance criteria still hold together `[ref: requirements/Section 4]`; update API docs
```

**Task count**: 3 tracked items.

---

## Activity Type Reference

Activity types are **specialist-selection hints, not phase boundaries**. A single vertical phase normally spans several activities — do not collect all tasks of one activity into a phase of its own.

| Activity | Description |
|----------|-------------|
| `domain-modeling` | Business entity and rule design |
| `data-architecture` | Database schema, migrations, queries |
| `backend-api` | REST/GraphQL endpoint implementation |
| `frontend-ui` | UI component implementation |
| `integration-test` | Cross-component integration testing |
| `e2e-test` | End-to-end user flow testing |
| `validate` | Quality gates and compliance checks |
| `business-acceptance` | requirements.md criteria verification |

---

## What Makes Good Implementation Plans

1. **Vertical Slices** - Each phase is one user-observable capability delivered end-to-end, validatable against an acceptance criterion on its own — never a single layer
2. **Shared Foundation in the First Slice** - The first phase to need a shared entity/schema builds it; later phases extend it. No foundation-only phase
3. **Logical Unit Tracking** - Each checkbox produces a verifiable outcome
4. **Embedded TDD Guidance** - Prime/Test/Implement/Validate as text, not checkboxes
5. **Specification Links** - Every task traces to requirements.md or solution.md
6. **Parallel Opportunities** - Capability slices that share only an early foundation are parallel-eligible; mark them
7. **No Time Estimates** - Focus on sequence, not duration
8. **Phase Validation** - Each phase verifies its own acceptance criteria; the final phase re-verifies them together
