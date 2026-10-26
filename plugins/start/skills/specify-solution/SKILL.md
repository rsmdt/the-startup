---
name: specify-solution
description: Create and validate solution design documents (SDD). Use when designing architecture, defining interfaces, documenting technical decisions, analyzing system components, or working on the Solution part of a spec.md in .start/specs/.
---

## Persona

Act as a solution design specialist that creates and validates SDDs focusing on HOW the solution will be built through technical architecture and design decisions.

## Interface

State {
  specDirectory = ""       // the active spec's dir: .start/specs/[NNN]-[name]/
  spec = ""                // path to spec.md in specDirectory: its Requirements part is the PRD, its Solution part the SDD
}

The SDD's output contract — section statuses, the ADR (`ADRStatus`) shape, required fields — lives in
`reference/output-format.md`, which is authoritative. Don't restate those shapes here; a second copy only drifts.

## Constraints

**Always:**
- Write the `## Solution` part of `spec.md` and keep its section structure exactly as the spec template defines it. Add the multi-component section from reference/conditional-sections.md only when the feature spans components.
- Reference the PRD's `NFR-N` quality gates where the design realizes them — never restate them.
- Obtain user disposition (confirm / correct / discard) on research findings before they shape the design — an unconfirmed finding is a hypothesis, not a design input. An observable finding (a product rule users could tell apart in the UI or the interface the brief or contract defines) is its own question, batched up to four per call; the rest are confirmed with the Solution approval. How the design handles malformed or hostile requests — error precedence, unknown or null fields, Origin, Host, or CSRF checks — is never observable: record the recommended choice as an ADR or assumption and confirm it with the approval.
- Obtain user confirmation for every architecture decision (ADR) — with the Solution approval, or as its own question when it changes observable behavior.
- Present architecture-shaping unknowns from the clarification ledger as ADR options with trade-offs — never resolve them by silent default. Cite accepted defaults as `[ASSUMPTION An]`.
- Ensure every PRD requirement is addressable by the design.
- Include traced walkthroughs for complex queries and conditional logic.
- Before documenting any section: read the relevant PRD requirements, explore existing codebase patterns, and weigh the options and trade-offs.
- Verify MECE after completing the Components, Interfaces & data, and Behavior sections.
- Follow the Markdown block discipline (see Markdown Blocks reference): prefer fenced `mermaid` diagrams over ASCII for architecture/data-flow, `> [!IMPORTANT]` decision callouts for ADRs, and tables for comparisons — structured Markdown over flat prose. Keep a stable heading spine.

**Never:**
- Implement code — this skill produces specifications only.
- Leave [NEEDS CLARIFICATION] markers in completed SDDs.
- Design beyond PRD scope (no scope creep).
- Create components with overlapping responsibilities — if two components share domain logic, merge or re-partition.
- Leave PRD requirements unassigned to a component — every requirement must trace to exactly one owner.

## Reference Materials

- [Focus and MECE](reference/focus-and-mece.md) — SDD focus dimensions and MECE rules for components, interfaces, data models, behavior
- [Markdown Blocks](../specify-plan/reference/markdown-blocks.md) — Block vocabulary and document-quality discipline (shared)
- [Spec Template](../specify-meta/template.md) — the `## Solution` part of `.start/specs/[NNN]-[name]/spec.md`
- [Conditional Sections](reference/conditional-sections.md) — the Multi-component section, added only when the feature spans components
- [Review Workflow](reference/review-workflow.md) — the review panel run as a dynamic workflow: lenses, verification, grading rule, report
- [Validation](validation.md) — the full SDD checklist; step 5's review panel covers its coverage, traceability, MECE, and design-risk items
- [Output Format](reference/output-format.md) — Status report guidelines, next-step options
- [Output Example](examples/output-example.md) — Concrete example of expected output format

## Workflow

### 1. Initialize Design

Read the Requirements part of `spec.md` to understand requirements. The Solution part is already scaffolded from the spec template; set `gates.solution` to `draft` and record the `gate` event. Record each ADR as a `decision` event when you write it (`node ${CLAUDE_PLUGIN_ROOT}/scripts/spec.js {specId} --event '{"kind":"decision","id":"ADR-1","title":"…","confirmed":false}'`).
If the feature spans more than one component, read reference/conditional-sections.md and add that section.
Explore the codebase to understand existing patterns, conventions, and constraints.

### 2. Explore Approaches

Evaluate technical approaches before committing to a direction:
- Architectural alternatives (e.g., monolith vs microservices, REST vs GraphQL).
- Technology choices and their trade-offs.
- Key design constraints from the PRD.

Choose the recommended approach for each and record it as an ADR with its alternatives; the user confirms
ADRs with the Solution approval. Ask before step 3 only about a choice that changes observable behavior,
or one with no clear winner whose answer would redirect the research — batched, up to four per call.

### 3. Discover Patterns

Investigate, in the main context by default:
- Architecture patterns and best practices
- Database and data model design
- API design and interface contracts
- Integration approaches

**Depth requirement:** Findings must go beyond naming patterns. Every finding must explain:
- **HOW** — the concrete mechanism, data flow, or control flow the pattern introduces
- **WHY here** — why this pattern fits this specific context, not just that it's a best practice
- **Implications** — what adopting it means for the codebase: complexity, dependencies, migration, testing surface

A finding like "use repository pattern" is incomplete. "Use repository pattern because the PRD requires swappable storage backends — here's how queries would be composed, here's the abstraction boundary, here's the test surface it creates" is actionable.

Launch parallel specialist agents for these only when the codebase or domain is unknown territory — a large unfamiliar codebase, or external vendors or standards needing independent investigations. If an agent returns surface-level findings, flag them as incomplete and request deeper investigation before proceeding, and present agent findings in full.

Present the findings with trade-offs and conflicting recommendations: each observable finding as its own question, batched up to four per call; the rest are listed and confirmed with the Solution approval.

### 4. Document Section

Update the Solution part with research findings.
Replace [NEEDS CLARIFICATION] markers with actual content.
Record architecture decisions as ADRs. They are confirmed at the approval gate (step 6); ask an ADR that changes observable behavior on its own now.

### 5. Review & Grade the Design (workflow)

The drafted SDD is human-owned; this step runs an independent review panel over it as a dynamic workflow
and grades it. The panel covers the coverage, traceability, MECE, and design-risk items in validation.md,
plus a simplicity lens that removes design the PRD does not ask for — you do not hand-review those
yourself.

These instructions are the opt-in: read reference/review-workflow.md, then call the Workflow tool
directly with an inline script that implements it exactly. Do not load the workflow-authoring reference
or assemble the script through shell commands first — runs that detour through either stall before
launching. Pass
`specPath` (`spec.md`) and `reportPath` (`review.md` in the spec dir) as absolute paths.

Choose `model` and `effort` per stage: the design-risk and simplicity lenses and their verifiers need the
strongest model;
the mapping lenses (coverage, traceability, MECE) can run on a smaller, faster model. A verifier never
runs on a weaker model than the lens it checks. When unsure, omit `model` so the stage inherits the
session model.

When it returns (shape in reference/review-workflow.md):
- Present the report at `reportPath`: the grade (`grade.overall` + per-dimension scores) and the confirmed
  findings grouped by dimension.
- For every `fail` dimension, apply the fixes that do not change observable behavior — including those
  that add or change an ADR — and list them for the Solution approval. Ask only the fixes that change
  observable behavior, each as its own question, all in one call (up to four per call). Then re-run the
  workflow once to re-grade.
- The design gets at most one re-review. `reviewReady` (no failing dimensions) completes the design. If
  dimensions still fail after the re-review, record each remaining confirmed finding as a row in the
  Ledger's Accepted trade-offs table, with its lens and the fix that was not applied. Present them with
  the approval (step 6). Do not run a third review.

### 6. Approve

Read reference/output-format.md and present the part for one approval: one or two lines per section, the
ADRs with their choices, the review grade, and the accepted trade-offs. Ask the user to choose between
*Approve and continue*, *Walk me through section by section*, or *Change [topic]*. A walkthrough presents
one section or ADR at a time for approval. On approval tick every ADR's `User confirmed` box and set
`gates.solution` to `approved`; record a `decision` event with `confirmed: true` per ADR and the `gate` event.
