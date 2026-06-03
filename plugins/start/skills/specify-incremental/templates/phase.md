---
title: "Phase {{PHASE_NUMBER}}: {{PHASE_TITLE}}"
status: pending
version: "1.0"
phase: {{PHASE_NUMBER}}
---

# Phase {{PHASE_NUMBER}}: {{PHASE_TITLE}}

## Phase Context

**GATE**: Read all referenced files before starting this phase.

**Specification References**:
- [NEEDS CLARIFICATION: List solution.md sections relevant to this phase]
- `[ref: solution/Section X; lines: Y-Z]`

**Key Decisions**:
- [NEEDS CLARIFICATION: Extract decisions from solution.md ADRs that affect this phase]

**Dependencies**:
- [NEEDS CLARIFICATION: List phases or tasks that must complete before this phase]

---

## Tasks

[NEEDS CLARIFICATION: Name the user-observable capability this phase delivers end-to-end - "Lets a user do X". A phase is a vertical slice spanning the layers below, not a single layer.]

- [ ] **T{{PHASE_NUMBER}}.1 [Domain/data for this capability]** `[activity: domain-modeling]`

  1. Prime: Read [entity/component] specification `[ref: solution/Section X; lines: Y-Z]`
  2. Test: [Behavior 1]; [Behavior 2]; [Edge case handling]
  3. Implement: Create/extend `src/[path]/[File].ts` (+ migration if this slice owns shared foundation)
  4. Validate: Unit tests pass; lint clean; types check

- [ ] **T{{PHASE_NUMBER}}.2 [Endpoint/interface for this capability]** `[activity: backend-api]`

  1. Prime: Review the API/interface contract `[ref: solution/Section X]`
  2. Test: [Request → response]; [Validation/error handling]
  3. Implement: Add the route/handler that exposes this behavior
  4. Validate: Contract tests pass; lint clean; types check

- [ ] **T{{PHASE_NUMBER}}.3 [UI + end-to-end for this capability]** `[activity: frontend-ui]`

  1. Prime: Read the UI spec and user journey `[ref: solution/Section X; requirements/Section Y]`
  2. Test: [Renders]; [Submits/updates]; end-to-end from user action to observable outcome
  3. Implement: Build the UI; wire it through the endpoint; add the E2E spec
  4. Validate: Component + E2E tests pass
  5. Success: [User can <capability> end-to-end] `[ref: requirements/AC-X.Y]`

- [ ] **T{{PHASE_NUMBER}}.4 Phase Validation** `[activity: validate]`

  - Run all Phase {{PHASE_NUMBER}} tests. Verify this phase satisfies its acceptance criterion on its own. Lint and typecheck pass.
