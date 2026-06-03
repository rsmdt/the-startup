# Example: Decomposition Manifest

Example manifest for the SQL Tools API. Each unit is a vertical capability slice; the first slice owns the shared foundation and the next two extend it in parallel.

---

## Manifest File (manifest.md)

```yaml
---
title: "SQL Tools API"
status: pending
threshold: 0.90
max_iterations: 5
---
```

```markdown
# Decomposition Manifest

## Units
- [ ] vq1: Validate a query — no dependencies
- [ ] eq1: Explain a query — after: vq1
- [ ] oq1: Optimize a query — after: vq1

## Execution Order
Group 1 (parallel): vq1
Group 2 (parallel): eq1, oq1
```

## Reading the Manifest

- **Checkboxes**: `[ ]` = pending, `[x]` = completed (updated by implement skill)
- **Dependencies**: `after: vq1` means this unit cannot start until vq1 is complete. The edge points to a real, independently-validatable slice (the validate-a-query capability), never to a "data layer" unit.
- **Execution groups**: Units in the same group with `(parallel)` can have code agents spawned concurrently. `eq1` and `oq1` both depend only on `vq1`, so once the foundation slice lands they build in parallel.
- **Threshold**: 90% of scenarios must pass for a unit to be marked complete
- **Max iterations**: Code agent retries up to 5 times per unit before escalating

> **Vertical-slicing check.** Each unit's holdout scenarios target its own endpoint (`vq1` → /validate, `eq1` → /explain, `oq1` → /optimize) and pass once that unit plus its dependency is built. If a unit's scenarios needed a *sibling* unit to exist, it would be a horizontal layer and the decomposition should be re-sliced.
