---
title: "{{FEATURE_TITLE}}"
status: draft
threshold: 0.90
version: 0.1.0
---

# Plan: {{FEATURE_TITLE}}

> [!NOTE]
> This plan is the **completion contract**: it defines the target state, the quality
> gates, and the completion criteria that prove the work is done. It does **not**
> prescribe implementation steps, sequencing, or decomposition — `/implement` chooses
> how to execute (direct, iterative, or decompose).

## Target State

{{Outcome-first prose: what the system does and looks like when this spec is fully
delivered. Self-contained — a reader with no chat history understands it. Reference
real components, interfaces, and data shapes from solution.md. Describe the destination,
not the route. No step list, no task order.}}

## Quality Gates

Every implementation must clear these before the spec is considered done:

- [ ] All existing tests pass; new behavior is covered by tests
- [ ] Lint and type checks pass
- [ ] {{coverage target / performance budget / security gate / a11y — as the solution implies}}

## Completion Criteria

Observable behaviors that prove the target state, grouped by capability. Each links to a
holdout scenario — the evaluation set `/implement` validates against. Priority tags:
`P0` critical path, `P1` important, `P2` edge case.

### {{Capability A}}

- [P0] [{{scenario title}}](scenarios/{{capability-a}}/{{scenario-name}}.md)
- [P1] [{{scenario title}}](scenarios/{{capability-a}}/{{scenario-name}}.md)

### {{Capability B}}

- [P0] [{{scenario title}}](scenarios/{{capability-b}}/{{scenario-name}}.md)

## Open Questions

{{If any decision that changes what "done" means is still open, list it here as a single
question with a recommended default. Otherwise state: "None — all completion criteria are
settled." Do not enumerate open questions anywhere else in the document.}}
