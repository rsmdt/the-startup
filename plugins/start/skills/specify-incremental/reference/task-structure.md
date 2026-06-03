# Task Structure Reference

Phase slicing axis, task granularity principle, TDD phase pattern, and metadata annotations.

---

## Phase Slicing Axis: Vertical, Not Horizontal

Slice phases by **capability**, not by **layer**. Each phase is a tracer-bullet: a thin, end-to-end path through whatever layers the behavior needs, delivering one user-observable outcome that can be validated against a requirements.md acceptance criterion on its own.

✅ **Vertical (correct)** — phase = capability:
- *Phase 1: Create a payment* → Payment entity + `POST /payments` + payment form + e2e creation test
- *Phase 2: Refund a payment* → extends entity + `POST /payments/:id/refund` + refund button + e2e refund test

❌ **Horizontal (wrong)** — phase = layer:
- *Phase 1: Database schema* (no API, no UI, satisfies no acceptance criterion alone)
- *Phase 2: API endpoints* (cannot be validated without the UI)
- *Phase 3: Frontend* (feature only works once all three land)

**Why vertical:** every phase is independently reviewable and PRD-validatable. Horizontal slicing defers all acceptance verification to the final phase — a database-only phase *cannot* satisfy "user can submit a payment" because that behavior needs all layers.

**Shared foundation:** when several phases need the same entity, schema, or scaffolding, the first phase that needs it builds it as part of its own end-to-end slice; later phases depend on that phase and extend the foundation. Never carve out a foundation-only phase.

**Activities are specialist hints, not phase boundaries.** A single vertical phase normally spans several activities (`domain-modeling` + `backend-api` + `frontend-ui` + `e2e-test`). Do not group all tasks of one activity into their own phase — that recreates horizontal slicing.

---

## Task Granularity Principle

**Track logical units that produce verifiable outcomes.** The TDD cycle is the execution method, not separate tracked items.

### Good Tracking Units (produces outcome)
- "Payment Entity" — Produces: working entity with tests
- "Stripe Adapter" — Produces: working integration with tests
- "Payment Form Component" — Produces: working UI with tests

### Bad Tracking Units (too granular)
- "Read payment interface contracts" — Preparation, not deliverable
- "Test Payment.validate() rejects negative amounts" — Part of larger outcome
- "Run linting" — Validation step, not deliverable

---

## Task Structure Pattern

```markdown
- [ ] **T1.1 Payment Entity** `[activity: domain-modeling]`

  **Prime**: Read payment interface contracts `[ref: solution/Section 4.2; lines: 145-200]`

  **Test**: Entity validation rejects negative amounts; supports currency conversion; handles refunds

  **Implement**: Create `src/domain/Payment.ts` with validation logic

  **Validate**: Run unit tests, lint, typecheck
```

The checkbox tracks "Payment Entity" as a unit. Prime/Test/Implement/Validate are embedded guidance.

---

## TDD Phase Structure

Every task follows red-green-refactor within this pattern:

### 1. Prime Context
- Read relevant specification sections
- Understand interfaces and contracts
- Load patterns and examples

### 2. Write Tests (Red)
- Test behavior before implementation
- Reference requirements.md acceptance criteria
- Cover happy path and edge cases

### 3. Implement (Green)
- Build to pass tests
- Follow solution.md architecture
- Use discovered patterns

### 4. Validate (Refactor)
- Run automated tests
- Check code quality (lint, format)
- Verify specification compliance

---

## Task Metadata

Use these annotations in the plan:

```markdown
- [ ] T1.2.1 [Task description] `[ref: solution/Section 5; lines: 100-150]` `[activity: backend-api]`
```

| Metadata | Description |
|----------|-------------|
| `[parallel: true]` | Tasks that can run concurrently |
| `[component: name]` | For multi-component features |
| `[ref: doc/section; lines: X-Y]` | Links to specifications |
| `[activity: type]` | Hint for specialist selection |

---

## Specification Compliance

Every phase should include a validation task:

```markdown
- [ ] **T1.3 Phase Validation** `[activity: validate]`

  Run all phase tests, linting, type checking. Verify against solution.md patterns and requirements.md acceptance criteria.
```

For complex phases, validation is embedded in each task's **Validate** step.

### Deviation Protocol

When implementation requires changes from the specification:
1. Document the deviation with clear rationale
2. Obtain approval before proceeding
3. Update solution.md when the deviation improves the design
4. Record all deviations in the relevant phase file for traceability
