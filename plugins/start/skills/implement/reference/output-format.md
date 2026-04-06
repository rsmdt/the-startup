# Output Format Reference

Guidelines for factory loop output. See `examples/output-example.md` for concrete rendered examples.

---

## Report Types

Four report types used during the factory loop:

1. **Manifest Discovery** — after parsing manifest.md, before execution starts
2. **Unit Result** — after each code agent + evaluation agent cycle completes
3. **Group Summary** — after all units in an execution group are resolved
4. **Completion Summary** — after all execution groups complete (or aborted with progress)

## Manifest Discovery Guidelines

Always include:
- Spec ID and feature name
- Threshold and max iterations
- Units with statuses (completed units will be skipped)
- Execution groups with their modes (parallel/sequential)
- Service start command and port
- Next group to execute

## Unit Result Guidelines

Always include:
- Unit ID and title
- Iteration number (e.g., "iteration 2 of 5")
- Satisfaction percentage and threshold
- Passed scenario count
- Failed scenarios with one-line summaries
- Decision: completed, queued for retry, or escalated

For completed units, also include:
- Files changed by the code agent
- Test results (passing count)

## Group Summary Guidelines

Always include:
- Group number and mode (parallel/sequential)
- Units completed / total in group
- Per-unit satisfaction percentages
- Total iterations used across all units in group
- Files changed across all units
- Constitution check results (if applicable)

## Completion Summary Guidelines

Always include:
- Spec ID and feature name
- Units completed / total units
- Total iterations across all units
- Final satisfaction percentage per unit
- Files changed (total count)
- Manifest status update (completed or failed)

## Manifest Status Updates

When updating manifest.md, document the changes:
- **Starting implementation**: frontmatter `status: pending` -> `status: in_progress`
- **Unit passes**: checkbox `- [ ] {id}:` -> `- [x] {id}:`
- **All units pass**: frontmatter `status: in_progress` -> `status: completed`
- **Aborted or failed**: frontmatter `status: in_progress` -> `status: failed`
