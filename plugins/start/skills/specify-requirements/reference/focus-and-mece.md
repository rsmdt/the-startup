# PRD Focus Areas and MECE Principle

Reference loaded by the specify-requirements skill when authoring or validating a PRD.

---

## PRD Focus Areas

When discovering and documenting, address four dimensions:

- **WHAT** needs to be built — features, capabilities
- **WHY** it matters — problem, outcome
- **WHO** uses it or is affected — roles, goals, concerns
- **WHEN** it succeeds — acceptance criteria, quality gates

**Out of scope:** Technical implementation, architecture, database schemas, API specifications — those belong in the SDD.

---

## MECE Principle

All structured enumerations in the PRD must be **Mutually Exclusive, Collectively Exhaustive** (MECE):

| Section | Mutually Exclusive | Collectively Exhaustive |
|---------|-------------------|------------------------|
| **Who** | Each row is a distinct role with its own goal or concern. No two rows for the same interest, and no role listed twice as user and affected. | Everyone who uses it and everyone whose week changes when it ships — support, ops, compliance, the negatively affected. Ask: "Who else interacts with this, and whose week changes?" |
| **Features** | Each user story captures a single, distinct behavior. No two stories describe the same capability, even across Must / Should. Won't-have names the confusable capabilities. | All capabilities needed to solve the stated problem are present. Ask: "If we shipped only these features, would the problem be fully solved for every role?" |
| **Acceptance Criteria** | Each criterion tests a unique condition. No two criteria verify the same behavior with different wording. | Every feature's happy path, error path, and edge cases are covered. Ask: "What input could break this that we haven't tested?" |
| **Quality Gates** | One gate per quality concern. No two gates constrain the same behavior with different thresholds, and no two pass or fail together by construction. | Every quality dimension relevant to the feature is gated or explicitly rejected with a reason, and every behavior whose regression would block release has a gate. Ask: "Which quality failure would make users leave even though the feature 'works'? What could quietly regress that no AC would catch?" |
| **Data** | Each data item appears once with a single stated purpose. | Everything collected, derived, or retained is listed. Ask: "What do we hold that the user doesn't know about?" |

**Stable IDs:** FR-N (features), AC-N.M (criteria), NFR-N (quality gates), DR-N (data). IDs are permanent once assigned — the plan's slice tasks trace them.

### How to apply

After completing each section, explicitly verify MECE before moving to the next:

1. **Exclusivity check** — Can any two items be merged without losing meaning? If yes, merge them.
2. **Exhaustiveness check** — Is there a scenario, user type, or capability not covered? If yes, add it.
3. **Cross-section check** — Do features in "Should have" duplicate behaviors already in "Must have"? Does a quality gate restate an acceptance criterion?
