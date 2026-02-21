# Output Format Reference

Guidelines for implementation output. See `examples/output-example.md` for concrete rendered examples.

---

## Report Types

Three report types used during implementation:

1. **Task Result** — after each agent completes a task (success or blocked)
2. **Phase Summary** — at each phase checkpoint before user confirmation
3. **Completion Summary** — after all phases complete

## Phase Summary Guidelines

Always include:
- Task completion ratio
- Files changed (list paths)
- Test status (passing count or failure details)
- Blockers (if any)
- Drift detection results
