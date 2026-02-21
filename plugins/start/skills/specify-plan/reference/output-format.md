# Output Format Reference

Status report template and next-step options for the specify-plan skill.

---

## Status Report Template

Present after each phase definition and after final validation:

```markdown
📋 PLAN Status: [spec-id]-[name]

Phases Defined:
- Phase 1 [Name]: ✅ Complete (X tasks)
- Phase 2 [Name]: 🔄 In progress
- Phase 3 [Name]: ⏳ Pending

Task Summary:
- Total tasks: [N]
- Parallel groups: [N]
- Dependencies: [List key dependencies]

Specification Coverage:
- PRD requirements mapped: [X/Y]
- SDD components covered: [X/Y]

Next Steps:
- [What needs to happen next]
```

---

## Next-Step Options

Use `AskUserQuestion` based on plan state:

**During phase definition:**
- "Define next phase"
- "Run validation"
- "Address gaps in current phase"
- "Complete PLAN"

**After validation:**
- "Fix validation failures"
- "Re-run validation"
- "Proceed to implementation"
