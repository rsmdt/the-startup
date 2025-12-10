# Implementation Phase Examples

Reference examples for structuring implementation phases.

## Example: Complete Phase Structure

```markdown
- [ ] T1 Phase 1: Payment Domain Foundation

    - [ ] T1.1 Prime Context
        - [ ] T1.1.1 Read payment interface contracts `[ref: SDD/Section 4.2; lines: 145-200]`
        - [ ] T1.1.2 Review existing repository patterns `[ref: docs/patterns/repository-pattern.md]`
        - [ ] T1.1.3 Load database schema design `[ref: SDD/Section 4.3]`

    - [ ] T1.2 Write Tests
        - [ ] T1.2.1 Unit tests for Payment entity validation `[ref: PRD/Section 4.1 - acceptance criteria]` `[activity: write-unit-tests]`
        - [ ] T1.2.2 Unit tests for PaymentRepository CRUD operations `[activity: write-unit-tests]`
        - [ ] T1.2.3 Integration tests for database persistence `[activity: write-integration-tests]`

    - [ ] T1.3 Implement
        - [ ] T1.3.1 Create Payment entity with validation logic `[activity: domain-modeling]`
        - [ ] T1.3.2 Create PaymentRepository with PostgreSQL adapter `[activity: data-architecture]`
        - [ ] T1.3.3 Create database migration for payments table `[activity: data-architecture]`

    - [ ] T1.4 Validate
        - [ ] T1.4.1 Run unit tests and verify all pass `[activity: run-tests]`
        - [ ] T1.4.2 Run linter and fix any issues `[activity: lint-code]`
        - [ ] T1.4.3 Verify entity matches SDD data model `[activity: review-code]`
```

## Example: Parallel Tasks

```markdown
- [ ] T2 Phase 2: API and Integration Layer

    - [ ] T2.1 API Development `[parallel: true]` `[component: backend]`
        - [ ] T2.1.1 Prime: Read API specification `[ref: SDD/Section 4.4]`
        - [ ] T2.1.2 Test: Controller endpoint tests `[activity: write-unit-tests]`
        - [ ] T2.1.3 Implement: PaymentController with routes `[activity: api-development]`
        - [ ] T2.1.4 Validate: API contract matches specification `[activity: review-code]`

    - [ ] T2.2 Stripe Integration `[parallel: true]` `[component: backend]`
        - [ ] T2.2.1 Prime: Read Stripe integration pattern `[ref: docs/interfaces/stripe-payment-integration.md]`
        - [ ] T2.2.2 Test: Mock Stripe client tests `[activity: write-unit-tests]`
        - [ ] T2.2.3 Implement: StripePaymentAdapter `[activity: backend-implementation]`
        - [ ] T2.2.4 Validate: Integration tests with Stripe test mode `[activity: run-tests]`
```

## Example: Multi-Component Feature

```markdown
- [ ] T3 Phase 3: Frontend Integration

    - [ ] T3.1 Payment Form Component `[component: frontend]`
        - [ ] T3.1.1 Prime: Read UI specifications `[ref: SDD/Section 5.1]`
        - [ ] T3.1.2 Test: Component render and interaction tests `[activity: write-component-tests]`
        - [ ] T3.1.3 Implement: PaymentForm React component `[activity: component-development]`
        - [ ] T3.1.4 Validate: Accessibility audit passes `[activity: accessibility-review]`

    - [ ] T3.2 State Management `[component: frontend]`
        - [ ] T3.2.1 Prime: Read state management pattern `[ref: docs/patterns/state-management.md]`
        - [ ] T3.2.2 Test: Reducer and selector tests `[activity: write-unit-tests]`
        - [ ] T3.2.3 Implement: Payment slice with async thunks `[activity: component-development]`
        - [ ] T3.2.4 Validate: State transitions match flow diagram `[activity: review-code]`

    - [ ] T3.3 Integration Point `[depends: T3.1, T3.2]`
        - [ ] T3.3.1 Wire PaymentForm to payment state
        - [ ] T3.3.2 Connect to backend API
        - [ ] T3.3.3 E2E test: Complete payment flow `[activity: write-e2e-tests]`
```

## Example: Final Validation Phase

```markdown
- [ ] T4 Integration & End-to-End Validation

    - [ ] T4.1 Cross-Component Testing
        - [ ] T4.1.1 All backend unit tests pass
        - [ ] T4.1.2 All frontend unit tests pass
        - [ ] T4.1.3 Integration tests for API ↔ Database
        - [ ] T4.1.4 Integration tests for Frontend ↔ API

    - [ ] T4.2 End-to-End Flows
        - [ ] T4.2.1 E2E: Happy path payment completion `[ref: PRD/Section 3.1]`
        - [ ] T4.2.2 E2E: Payment failure handling `[ref: PRD/Section 3.2]`
        - [ ] T4.2.3 E2E: Payment history display `[ref: PRD/Section 3.3]`

    - [ ] T4.3 Quality Gates
        - [ ] T4.3.1 Performance: API response < 200ms p95 `[ref: SDD/Section 10]`
        - [ ] T4.3.2 Security: Input validation verified `[ref: SDD/Section 8]`
        - [ ] T4.3.3 Coverage: > 80% line coverage

    - [ ] T4.4 Final Acceptance
        - [ ] T4.4.1 All PRD acceptance criteria verified `[ref: PRD/Section 4]`
        - [ ] T4.4.2 Implementation follows SDD design `[ref: SDD/Section 5]`
        - [ ] T4.4.3 Documentation updated for API changes
        - [ ] T4.4.4 Build and deployment verification
```

## Activity Type Reference

Common activity types for specialist selection:

| Activity | Description |
|----------|-------------|
| `domain-modeling` | Business entity and rule design |
| `data-architecture` | Database schema, migrations, queries |
| `api-development` | REST/GraphQL endpoint implementation |
| `component-development` | UI component implementation |
| `write-unit-tests` | Unit test creation |
| `write-integration-tests` | Integration test creation |
| `write-e2e-tests` | End-to-end test creation |
| `write-component-tests` | UI component tests |
| `run-tests` | Test execution and verification |
| `lint-code` | Code linting and style fixes |
| `format-code` | Code formatting |
| `review-code` | Code review and quality check |
| `accessibility-review` | A11y compliance check |
| `security-review` | Security assessment |
| `backend-implementation` | General backend code |
| `frontend-implementation` | General frontend code |
| `business-acceptance` | PRD criteria verification |

## What Makes Good Implementation Plans

1. **Clear Task Boundaries** - Each task is completable independently
2. **Specification Links** - Every task traces to PRD/SDD
3. **TDD Structure** - Test → Implement → Validate
4. **Parallel Opportunities** - Independent work clearly marked
5. **No Time Estimates** - Focus on sequence, not duration
6. **Activity Hints** - Guide specialist selection
7. **Final Validation** - Comprehensive quality gates
