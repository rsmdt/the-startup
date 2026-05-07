# Output Format Reference

Guidelines for implementation output. See `examples/output-example.md` for concrete rendered examples.

---

## Report Types

Four report types used during implementation:

1. **Plan Discovery** — after reading plan/README.md, before execution starts
2. **Task Result** — after each agent completes a task (success or blocked)
3. **Phase Summary** — at each phase checkpoint before user confirmation
4. **Completion Summary** — after all phases complete (or paused with progress)

## Plan Discovery Guidelines

Always include:
- Spec ID and feature name
- Total phases discovered and their statuses
- Which phases will be skipped (already completed)
- Which phase will be resumed (if in_progress)
- Next phase to execute

## Phase Summary Guidelines

Always include:
- Phase number and title
- Task completion ratio (completed/total)
- Files changed (list paths)
- Test status (passing count or failure details)
- Blockers (if any)
- Drift detection results
- Updated status: phase file frontmatter and README checkbox

## Completion Summary Guidelines

Always include:
- Spec ID and feature name
- Total phases completed / total phases
- Total tasks executed across all phases
- Test status (all passing or failure details)
- Files changed (total count with additions/deletions)
- Mode used (Standard or Team)

## Phase Status Updates

When a phase transitions status, document the file changes:
- **Starting phase**: `phase-N.md` frontmatter `status: pending` → `status: in_progress`
- **Completing phase**: `phase-N.md` frontmatter `status: in_progress` → `status: completed`, `plan/README.md` checkbox `- [ ]` → `- [x]`
- **Pausing mid-plan**: Report which phases are completed and which remain pending
