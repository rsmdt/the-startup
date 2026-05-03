# SDD Focus and MECE Principle

Reference loaded by the specify-solution skill when authoring or validating an SDD.

---

## SDD Focus

When designing, address four dimensions:

- **HOW** it will be built — architecture, patterns, approach
- **WHERE** code lives — directory structure, components, layers
- **WHAT** interfaces exist — APIs, data models, integrations
- **WHY** decisions were made — ADRs with rationale and trade-offs

---

## MECE Principle

All structural decompositions in the SDD must be **Mutually Exclusive, Collectively Exhaustive** (MECE):

| Section | Mutually Exclusive | Collectively Exhaustive |
|---------|-------------------|------------------------|
| **Components** | Each component has a single, distinct responsibility. No two components should own the same domain logic or serve the same purpose. | All system capabilities from the PRD are assigned to exactly one component. Ask: "Which component handles X?" — if the answer is ambiguous, the decomposition has overlap. |
| **Interfaces** | Each interface serves a distinct purpose. No two interfaces should expose the same operation or data shape. | All communication paths between components, external systems, and data stores are documented. Ask: "How does component A talk to component B?" — if undocumented, there's a gap. |
| **Data Models** | Each entity owns a distinct slice of the domain. No two entities should store the same business data. | All data required by the components and interfaces is modeled. Ask: "Where is X stored?" — if unanswerable, there's a gap. |
| **Acceptance Criteria (EARS)** | Each criterion specifies a unique system behavior. No two criteria should verify the same thing with different triggers. | Every PRD acceptance scenario has a corresponding system-level criterion. Ask: "How does the system satisfy PRD/AC-X.Y?" — if unanswerable, there's a gap. |

### How to apply

During validation, explicitly run MECE checks:

1. **Responsibility matrix** — Map each PRD requirement to exactly one component. Flag requirements mapped to multiple components (overlap) or zero components (gap).
2. **Interface deduplication** — Verify no two interfaces serve the same consumer-to-provider path.
3. **Criteria traceability** — Verify 1:1 mapping between PRD acceptance criteria and EARS criteria.
