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
- [NEEDS CLARIFICATION: List SDD sections relevant to this phase]
- `[ref: SDD/Section X; lines: Y-Z]`

**Key Decisions**:
- [NEEDS CLARIFICATION: Extract decisions from SDD ADRs that affect this phase]

**Dependencies**:
- [NEEDS CLARIFICATION: List phases or tasks that must complete before this phase]

---

## Tasks

[NEEDS CLARIFICATION: Describe what this phase delivers - "Establishes X capability" or "Enables Y functionality"]

- [ ] **T{{PHASE_NUMBER}}.1 [Primary deliverable name]** `[activity: domain-modeling]`

  1. Prime: Read [entity/component] specification `[ref: SDD/Section X; lines: Y-Z]`
  2. Test: [Behavior 1]; [Behavior 2]; [Edge case handling]
  3. Implement: Create `src/[path]/[File].ts` with [key capability]
  4. Validate: Unit tests pass; lint clean; types check
  5. Success: [Criterion 1] `[ref: PRD/AC-X.Y]`; [Criterion 2] `[ref: SDD/Section]`

- [ ] **T{{PHASE_NUMBER}}.2 [Secondary deliverable name]** `[activity: data-architecture]`

  1. Prime: Review [pattern/interface] requirements `[ref: SDD/Section X]`
  2. Test: [CRUD operations]; [Query patterns]; [Error handling]
  3. Implement: Create `src/[path]/[File].ts` with [key capability]
  4. Validate: Integration tests pass; lint clean; types check
  5. Success: [Criterion 1] `[ref: PRD/AC-X.Y]`; [Criterion 2] `[ref: SDD/Section]`

- [ ] **T{{PHASE_NUMBER}}.3 Phase Validation** `[activity: validate]`

  - Run all Phase {{PHASE_NUMBER}} tests. Verify against SDD patterns and PRD acceptance criteria. Lint and typecheck pass.
