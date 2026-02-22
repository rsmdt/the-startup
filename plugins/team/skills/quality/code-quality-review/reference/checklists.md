# Review Checklists

Scope-appropriate checklists for different review sizes, plus review process anti-patterns.

## Quick Review Checklist (< 100 lines)

- [ ] Code compiles and tests pass
- [ ] Logic appears correct for stated purpose
- [ ] No obvious security issues
- [ ] Naming is clear
- [ ] No magic numbers or strings

## Standard Review Checklist (100-500 lines)

All of the above, plus:
- [ ] Design follows project patterns
- [ ] Error handling is appropriate
- [ ] Tests cover new functionality
- [ ] No significant duplication
- [ ] Performance is reasonable

## Deep Review Checklist (> 500 lines or critical)

All of the above, plus:
- [ ] Architecture aligns with system design
- [ ] Security implications considered
- [ ] Backward compatibility maintained
- [ ] Documentation updated
- [ ] Migration/rollback plan if needed

## Review Process Anti-Patterns

| Anti-Pattern | Description | Better Approach |
|--------------|-------------|-----------------|
| **Nitpicking** | Focusing on style over substance | Use linters for style |
| **Drive-by Review** | Quick approval without depth | Allocate proper time |
| **Gatekeeping** | Blocking for personal preferences | Focus on objective criteria |
| **Ghost Review** | Approval without comments | Add at least one observation |
| **Review Bombing** | Overwhelming with comments | Prioritize and limit to top issues |
| **Delayed Review** | Letting PRs sit for days | Commit to turnaround time |

## Review Metrics

| Metric | Target | What It Indicates |
|--------|--------|-------------------|
| Review Turnaround | < 24 hours | Team velocity |
| Comments per Review | 3-10 | Engagement level |
| Defects Found | Decreasing trend | Quality improvement |
| Review Time | < 60 min for typical PR | Right-sized changes |
| Approval Rate | 70-90% first submission | Clear standards |
