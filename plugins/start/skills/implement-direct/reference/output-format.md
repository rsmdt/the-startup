# Output Format Reference

Guidelines for direct-mode implementation output. See `examples/output-example.md` for a concrete rendered example.

---

## Report Types

Three report types used during direct implementation:

1. **Discovery** — after reading context documents, before decomposition
2. **Unit Result** — after each delegated subagent completes (success or blocked)
3. **Completion Summary** — after delegation and validation are done

## Discovery Guidelines

Always include:
- Spec ID and feature name (if applicable, otherwise the brief title)
- Context documents found (requirements.md, solution.md, AGENTS.md, others)
- Context documents missing (note as advisory, not blocking)
- Recommended decomposition (number of delivery units)

If decomposition would yield more than 3 units, recommend the user re-run `/start:specify` and select Standard or Factory tier.

## Unit Result Guidelines

Always include:
- Unit title
- Area (backend, frontend, data, infra, docs)
- Files changed (list paths)
- Test status (passing count or failure details)
- Blockers (if any)

## Completion Summary Guidelines

Always include:
- Spec ID or brief title
- Delivery units completed / total
- Total files changed (count with additions/deletions)
- Test status (all passing or failure details)
- Drift detection results
- Constitution check results (if applicable)
- Mode used (Standard or Team)
