# Output Format Reference

Guidelines for specification output. See `examples/output-example.md` for concrete rendered examples.

---

## Documentation Structure

Specifications live in `.start/specs/[NNN]-[name]/`:
- `spec.md` — one document: the Requirements (what and why), Solution (how), and Plan (phases and their slices) parts, each closed by its gate in the frontmatter, and the Ledger (open questions, assumptions, decision log, accepted trade-offs)
- `review.md` — the solution review report
- `.tracker/` — Phase tasks, their slice tasks (code briefs), and each slice's scenario child tasks (holdout evidence)

## Decision Logging

When the user skips a phase or makes a non-default choice, log it in the Decision log of the `spec.md` Ledger with date, decision, and rationale.

## Assumption Ledger Logging

The Clarify Intent phase writes the ledger to the Assumptions table of `spec.md` (ID | Assumption | Impact if wrong | Observable | Status) and updates statuses as they change. The document cites Standing entries as `[ASSUMPTION An]` in the affected sections instead of restating them. The completion summary lists all Standing assumptions; any architecture-shaping Standing assumption caps readiness at MEDIUM.
