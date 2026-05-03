# Implementation Phase Examples

Reference examples for structuring implementation phases with grouped task tracking.

## Core Principle

**Track logical units that produce verifiable outcomes.** Each checkbox represents a deliverable, not a step in the process. The TDD cycle (Prime → Test → Implement → Validate) is embedded guidance within each task.

---

## Example: Foundation Phase

```markdown
### Phase 1: Payment Domain Foundation

Establishes core domain entities, repository patterns, and database schema.

- [ ] **T1.1 Payment Entity** `[activity: domain-modeling]`

  **Prime**: Read payment interface contracts and validation rules `[ref: solution/Section 4.2; lines: 145-200]`

  **Test**: Entity validation rejects negative amounts; supports multiple currencies; calculates fees correctly; handles partial refunds

  **Implement**: Create `src/domain/Payment.ts` with Amount value object and validation logic

  **Validate**: Run unit tests, lint, typecheck

- [ ] **T1.2 Payment Repository** `[activity: data-architecture]`

  **Prime**: Review existing repository patterns and database schema `[ref: existing repository patterns]` `[ref: solution/Section 4.3]`

  **Test**: CRUD operations work correctly; queries filter by status and date range; handles concurrent updates

  **Implement**: Create `src/repositories/PaymentRepository.ts` with PostgreSQL adapter; create migration for payments table

  **Validate**: Run integration tests against test database

- [ ] **T1.3 Phase Validation** `[activity: validate]`

  Run all unit and integration tests. Verify entity matches solution.md data model. Lint and typecheck pass.
```

**Task count**: 3 tracked items (vs ~12 with granular tracking)

---

## Example: Parallel Tasks

```markdown
### Phase 2: API and Integration Layer

API endpoints and external service integration. These can be developed in parallel.

- [ ] **T2.1 Payment API Endpoints** `[parallel: true]` `[component: backend]`

  **Prime**: Read API specification and authentication requirements `[ref: solution/Section 4.4]`

  **Test**: POST /payments creates payment and returns 201; GET /payments/:id returns payment or 404; validation errors return 422 with details

  **Implement**: Create `src/controllers/PaymentController.ts` with create, get, list routes

  **Validate**: API contract tests pass; authentication enforced

- [ ] **T2.2 Stripe Integration** `[parallel: true]` `[component: backend]`

  **Prime**: Read Stripe integration pattern and webhook handling `[ref: existing integration patterns]`

  **Test**: Charges created with correct amount; webhook validates signature; handles declined cards gracefully

  **Implement**: Create `src/adapters/StripePaymentAdapter.ts` with charge and refund methods

  **Validate**: Integration tests pass with Stripe test mode

- [ ] **T2.3 Phase Validation** `[activity: validate]`

  Run all API and integration tests. Verify endpoints match OpenAPI spec. Lint and typecheck pass.
```

**Task count**: 3 tracked items (vs ~10 with granular tracking)

---

## Example: Multi-Component Feature

```markdown
### Phase 3: Frontend Integration

UI components and state management for the payment flow.

- [ ] **T3.1 Payment Form Component** `[component: frontend]`

  **Prime**: Read UI specifications and form validation rules `[ref: solution/Section 5.1]`

  **Test**: Form renders with all fields; validates card number format; submits on valid input; shows error states

  **Implement**: Create `src/components/PaymentForm.tsx` with Stripe Elements integration

  **Validate**: Component tests pass; accessibility audit passes

- [ ] **T3.2 Payment State Management** `[component: frontend]`

  **Prime**: Read state management pattern and async handling `[ref: existing state management patterns]`

  **Test**: Loading states during API calls; success state updates UI; error state shows message; retry logic works

  **Implement**: Create `src/store/paymentSlice.ts` with async thunks for API calls

  **Validate**: State transitions match flow diagram; reducer tests pass

- [ ] **T3.3 Payment Flow Integration** `[depends: T3.1, T3.2]`

  **Prime**: Review complete user journey from cart to confirmation

  **Test**: E2E flow from form submission to confirmation page; handles payment failures; shows receipt

  **Implement**: Wire PaymentForm to payment state; connect to backend API; add confirmation page

  **Validate**: E2E test passes; manual QA of happy path and error cases
```

**Task count**: 3 tracked items (vs ~12 with granular tracking)

---

## Example: Final Validation Phase

```markdown
### Phase 4: Integration & Validation

Full system validation with all components working together.

- [ ] **T4.1 Integration Testing** `[activity: integration-test]`

  Verify cross-component integration: API ↔ Database, Frontend ↔ API, Payment ↔ Stripe

- [ ] **T4.2 E2E User Flows** `[activity: e2e-test]`

  Verify complete user journeys: happy path payment `[ref: requirements/Section 3.1]`; payment failure handling `[ref: requirements/Section 3.2]`; payment history display `[ref: requirements/Section 3.3]`

- [ ] **T4.3 Quality Gates** `[activity: validate]`

  Performance: API response < 200ms p95 `[ref: solution/Section 10]`; Security: input validation verified; Coverage: > 80% line coverage

- [ ] **T4.4 Specification Compliance** `[activity: business-acceptance]`

  All requirements.md acceptance criteria verified `[ref: requirements/Section 4]`; implementation follows solution.md design; documentation updated for API changes
```

**Task count**: 4 tracked items (vs ~15 with granular tracking)

---

## Activity Type Reference

Common activity types for specialist selection:

| Activity | Description |
|----------|-------------|
| `domain-modeling` | Business entity and rule design |
| `data-architecture` | Database schema, migrations, queries |
| `api-development` | REST/GraphQL endpoint implementation |
| `component-development` | UI component implementation |
| `backend-implementation` | General backend code |
| `frontend-implementation` | General frontend code |
| `integration-test` | Cross-component integration testing |
| `e2e-test` | End-to-end user flow testing |
| `validate` | Quality gates and compliance checks |
| `business-acceptance` | requirements.md criteria verification |

---

## What Makes Good Implementation Plans

1. **Logical Unit Tracking** - Each checkbox produces a verifiable outcome
2. **Embedded TDD Guidance** - Prime/Test/Implement/Validate as text, not checkboxes
3. **Specification Links** - Every task traces to requirements.md or solution.md
4. **Parallel Opportunities** - Independent work clearly marked
5. **No Time Estimates** - Focus on sequence, not duration
6. **Activity Hints** - Guide specialist selection
7. **Phase Validation** - Quality gates at phase boundaries
