# Output Format Reference

Guidelines for specification output. See `examples/output-example.md` for concrete rendered examples.

---

## Documentation Structure

Specifications live in `.start/specs/[NNN]-[name]/`:
- `README.md` — Decisions and progress tracking
- `requirements.md` — What and why
- `solution.md` — How
- `manifest.md` — Decomposition manifest (units, dependencies, execution order)
- `units/` — Factory-sized specs (one per unit of work)
- `scenarios/` — Holdout evaluation scenarios per unit

## Decision Logging

When user skips a phase or makes a non-default choice, log it in the spec README.md decisions table with date, decision, and rationale.
