# Example: Decomposition Manifest

Example manifest for a SQL validation feature.

---

## Manifest File (manifest.md)

```yaml
---
title: "SQL Validation Feature"
status: pending
threshold: 0.90
max_iterations: 5
---
```

```markdown
# Decomposition Manifest

## Units
- [ ] dm1: Data models — no dependencies
- [ ] ve1: Validation endpoint — after: dm1
- [ ] ve2: Input sanitization — after: ve1
- [ ] rl1: Rate limiting — after: ve1

## Execution Order
Group 1 (parallel): dm1
Group 2 (parallel): ve1
Group 3 (parallel): ve2, rl1
```

## Reading the Manifest

- **Checkboxes**: `[ ]` = pending, `[x]` = completed (updated by implement skill)
- **Dependencies**: `after: dm1` means this unit cannot start until dm1 is complete
- **Execution groups**: Units in the same group with `(parallel)` can have code agents spawned concurrently
- **Threshold**: 90% of scenarios must pass for a unit to be marked complete
- **Max iterations**: Code agent retries up to 5 times per unit before escalating
