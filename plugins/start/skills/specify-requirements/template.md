---
title: "[NEEDS CLARIFICATION: Feature title]"
status: draft
version: "1.0"
---

# Product Requirements Document

> [!NOTE]
> This PRD answers **WHY** and **WHAT** — the problem, who has it, and what "done" looks like.
> It carries no technical design (that is `solution.md`) and no implementation steps.
> Unresolved decisions live in one place: the **Open Questions** table at the foot of the document.
> **Stable IDs:** features `FR-N`, criteria `AC-N.M`, quality gates `NFR-N`, data `DR-N` — permanent
> once assigned; the plan's slice tasks trace them, so renumbering breaks traceability.

## Problem

[NEEDS CLARIFICATION: What specific problem do people face today, who has it, and what is the evidence? What does it cost them — and us — to leave it unsolved?]

## Outcome

[NEEDS CLARIFICATION: One paragraph. What is true for users once this ships, and why will they choose it over what they do now? This is the "done" the whole spec is measured against.]

## Who

Everyone with a stake — those who use it and those whose week changes when it ships (support, ops, compliance, the negatively affected). One row per distinct goal or concern.

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| [NEEDS CLARIFICATION: e.g. Claims handler] | uses it | [what they are trying to do and what gets in the way] | [their acceptance condition] |
| [e.g. Support team] | affected | [what this feature changes for them] | [their acceptance condition] |

## Features

Each feature is one distinct capability, stated once. Acceptance criteria are Gherkin — `Given / When / Then`, one testable condition each, covering the happy path, the error path, and the edge case.

### Must have

#### FR-1: [NEEDS CLARIFICATION: Feature name]
- **User story:** As a [role], I want to [action] so that [benefit]
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given [state], When [action], Then [observable outcome]
  - [ ] **AC-1.2** — Given [state], When [action], Then [observable outcome]
  - [ ] **AC-1.3** — Given [edge case], When [action], Then [correct behavior]

#### FR-2: [NEEDS CLARIFICATION: Feature name]
[Repeat the structure]

### Should have

[NEEDS CLARIFICATION: What significantly improves the outcome but is not required to deliver it? Same FR / AC structure, or "None".]

### Won't have (this phase)

[NEEDS CLARIFICATION: What is explicitly out of scope — especially capabilities a reader could confuse with the ones above?]

## Quality gates

How well it must work, from the user's and the business's view, and what must be demonstrably true before release. One row per quality concern — performance, reliability, security and privacy, usability and accessibility, compliance, and any behavior that could quietly regress without an AC catching it. Every gate is measurable by a tester who does not know the implementation; the mechanism belongs in `solution.md`.

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | [NEEDS CLARIFICATION: e.g. Performance] | [e.g. Status lookups feel instant] | [e.g. p95 response < 300 ms at 1,000 concurrent users] |
| NFR-2 | [e.g. Behavior under test] | [e.g. Acts when it should, silent when it shouldn't] | [e.g. acts in ≥ N of M trigger scenarios; zero unprompted actions on the no-trigger set] |

## Data

What the feature collects and keeps, as an obligation: purpose, retention, sensitivity, and the rights users hold over it. The technical data model belongs in `solution.md`. Omit this section only when the feature holds no personal or sensitive data.

| ID | Data item | Purpose | Retention | Classification |
|----|-----------|---------|-----------|----------------|
| DR-1 | [NEEDS CLARIFICATION: e.g. Customer email] | [e.g. receipt delivery] | [e.g. life of account + 90 days] | [PII / sensitive / internal / public] |

- **User rights:** [NEEDS CLARIFICATION: How do users view, export, correct, or delete this data, and how fast must each request be honored?]

## Constraints & dependencies

- [NEEDS CLARIFICATION: Budget, timeline, platform, legal, or compliance limits on the solution space]
- [ASSUMPTION] entries from the clarification ledger appear here — surfaced, never erased

| Dependency | Owner | Needed by | Impact if late |
|------------|-------|-----------|----------------|
| [NEEDS CLARIFICATION: e.g. Vendor API access] | [team / vendor] | [phase or date] | [what slips or degrades] |

## Glossary

One term, one meaning. Define only the domain terms and acronyms a new team member would stumble on; two names for the same concept is how overlapping requirements sneak in. "None" is a valid entry.

| Term | Definition |
|------|------------|
| [NEEDS CLARIFICATION: term] | [Single authoritative definition] |

## Open questions

The single home for every unresolved decision in this PRD. One row per question, each with a recommended default so the spec can proceed if unanswered. Resolve a row by checking it and recording the choice; "None — all decisions settled." is a valid, good state.

| Question | Options | Recommended default | ☐ Resolved |
|----------|---------|---------------------|:----------:|
| [What needs stakeholder input?] | [Option A / Option B] | [Default + one-line why] | [ ] |
