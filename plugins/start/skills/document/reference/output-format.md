# Output Format Reference

Summary template and coverage metrics for the document skill.

---

## Documentation Summary Template

Present after all documentation is generated:

```markdown
## Documentation Complete

**Target**: [what was documented]

### Changes Made

| File | Action | Coverage |
|------|--------|----------|
| `path/file.ts` | Added JSDoc | 15 functions |
| `docs/api.md` | Created | 8 endpoints |
| `README.md` | Updated | 3 sections |

### Coverage Metrics

| Area | Before | After |
|------|--------|-------|
| Code | X% | Y% |
| API | X% | Y% |
| README | Partial | Complete |

### Next Steps

- [Remaining gaps to address]
- [Stale docs to review]
```

---

## Next Steps Options

Use `AskUserQuestion` after summary:
- "Address remaining gaps"
- "Review stale documentation"
- "Done"
