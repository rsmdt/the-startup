# Output Format Reference

Guidelines for specification output. See `examples/output-example.md` for concrete rendered examples.

---

## Documentation Structure

Specifications live in `.start/specs/[NNN]-[name]/`:
- `README.md` — Decisions, assumption ledger, and progress tracking
- `requirements.mdx` — What and why
- `solution.mdx` — How
- `plan.mdx` — Completion contract (target state, quality gates, completion criteria)
- `scenarios/` — Holdout validation scenarios per phase
- `spec-report.mdx` — Editor-readable report (ledger, findings, documents, decisions)

## Specification Report

At Finalize, write `spec-report.mdx` from templates/spec-report.mdx — MDX that degrades to readable Markdown (tables, callouts, task lists; only `<details>` HTML, no custom components), so it reviews cleanly in any editor and on GitHub. The chat completion summary stays concise and links the report; the report carries the full ledger, finding dispositions, and decision log. Keep its heading spine stable — other skills may grep it.

## Decision Logging

When the user skips a phase or makes a non-default choice, log it in the spec README.md decisions table with date, decision, and rationale.

## Assumption Ledger Logging

The Clarify Intent phase writes the ledger to the README's Assumption Ledger table (Topic | Default | Impact if wrong | Status) and updates statuses as they change. Standing entries additionally appear as `[ASSUMPTION]` markers in the affected document sections. The completion summary lists all Standing assumptions; any architecture-shaping Standing assumption caps readiness at MEDIUM.
