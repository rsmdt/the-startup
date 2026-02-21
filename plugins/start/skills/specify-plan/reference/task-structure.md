# Task Structure Reference

Task granularity principle, TDD phase pattern, and metadata annotations.

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

  **Prime**: Read payment interface contracts `[ref: SDD/Section 4.2; lines: 145-200]`

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
- Reference PRD acceptance criteria
- Cover happy path and edge cases

### 3. Implement (Green)
- Build to pass tests
- Follow SDD architecture
- Use discovered patterns

### 4. Validate (Refactor)
- Run automated tests
- Check code quality (lint, format)
- Verify specification compliance

---

## Task Metadata

Use these annotations in the plan:

```markdown
- [ ] T1.2.1 [Task description] `[ref: SDD/Section 5; lines: 100-150]` `[activity: backend-api]`
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

  Run all phase tests, linting, type checking. Verify against SDD patterns and PRD acceptance criteria.
```

For complex phases, validation is embedded in each task's **Validate** step.

### Deviation Protocol

When implementation requires changes from the specification:
1. Document the deviation with clear rationale
2. Obtain approval before proceeding
3. Update SDD when the deviation improves the design
4. Record all deviations in the plan for traceability
