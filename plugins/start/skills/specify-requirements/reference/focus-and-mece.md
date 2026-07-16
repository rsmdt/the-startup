# PRD Focus Areas and MECE Principle

Reference loaded by the specify-requirements skill when authoring or validating a PRD.

---

## PRD Focus Areas

When discovering and documenting, address four dimensions:

- **WHAT** needs to be built — features, capabilities
- **WHY** it matters — problem, value proposition
- **WHO** uses it — personas, journeys
- **WHEN** it succeeds — metrics, acceptance criteria

**Out of scope:** Technical implementation, architecture, database schemas, API specifications — those belong in the SDD.

---

## MECE Principle

All structured enumerations in the PRD must be **Mutually Exclusive, Collectively Exhaustive** (MECE):

| Section | Mutually Exclusive | Collectively Exhaustive |
|---------|-------------------|------------------------|
| **Stakeholders** | Each row owns a distinct concern. No two stakeholders listed for the same interest. | Everyone affected by shipping — including support, ops, compliance, and those affected negatively. Ask: "Whose week changes when this ships?" |
| **User Personas** | Each persona represents a distinct user type with unique goals and pain points. No two personas should overlap in role or motivation. | All relevant user types are represented. Ask: "Who else interacts with this system?" |
| **User Journeys** | Each journey describes a distinct path through the system. No two journeys should cover the same sequence of actions for the same persona. | All primary and secondary paths are mapped, including error/recovery paths. Ask: "What other ways do users accomplish this goal?" |
| **Feature Requirements** | Each user story captures a single, distinct behavior. No two stories should describe the same capability, even across MoSCoW categories. | All capabilities needed to solve the stated problem are present. Ask: "If we shipped only these features, would the problem be fully solved for every persona?" |
| **Acceptance Criteria** | Each criterion tests a unique condition. No two criteria should verify the same behavior with different wording. | Every feature's happy path, error path, and edge cases are covered. Ask: "What input could break this that we haven't tested?" |
| **Non-Functional Requirements** | One requirement per quality concern. No two NFRs constrain the same behavior with different thresholds. | Every quality dimension relevant to the feature is specified or explicitly rejected with a reason. Ask: "Which quality failure would make users leave even though the feature 'works'?" |
| **Evaluation Criteria** | Each gate tests a distinct behavior — no two gates pass or fail together by construction. | Every behavior whose failure would block release has a gate. Ask: "What could quietly regress that no AC would catch?" |
| **Data Requirements** | Each data item appears once with a single stated purpose. | Everything collected, derived, or retained is listed. Ask: "What do we hold that the user doesn't know about?" |

**Stable IDs:** FR-N (features), AC-N.M (criteria), NFR-N, EV-N, DR-N. IDs are permanent once assigned — downstream traceability (plan.mdx) parses them.

### How to apply

After completing each section, explicitly verify MECE before moving to the next:

1. **Exclusivity check** — Can any two items be merged without losing meaning? If yes, merge them.
2. **Exhaustiveness check** — Is there a scenario, user type, or capability not covered? If yes, add it.
3. **Cross-section check** — Do features in "Should Have" duplicate behaviors already in "Must Have"? Do journeys overlap with different personas doing the same thing?
