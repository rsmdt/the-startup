# Output Format Reference

Status report template and next-step options for the specify-solution skill.

---

## Status Report Template

Present after each design cycle and after final validation:

```markdown
🏗️ SDD Status: [spec-id]-[name]

Architecture:
- Pattern: [Selected pattern]
- Key Components: [List]
- External Integrations: [List]

Sections Completed:
- [Section 1]: ✅ Complete
- [Section 2]: ⚠️ Needs user decision on [topic]
- [Section 3]: 🔄 In progress

ADRs:
- [ADR-1]: ✅ Confirmed
- [ADR-2]: ⏳ Pending confirmation

Validation Status:
- [X] items passed
- [Y] items pending

Next Steps:
- [What needs to happen next]
```

---

## Next-Step Options

Use `AskUserQuestion` based on SDD state:

**During design cycles:**
- "Address pending ADRs"
- "Continue to next section"
- "Run validation"
- "Complete SDD"

**After validation:**
- "Fix validation failures"
- "Re-run validation"
- "Proceed to implementation plan"
