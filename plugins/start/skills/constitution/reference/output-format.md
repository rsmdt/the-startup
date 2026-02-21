# Output Format Reference

Templates for proposed rules presentation and constitution summaries.

---

## Proposed Rules Format

Present discovered rules grouped by category with level classification:

```markdown
📜 Proposed Constitution

## Security (N rules)
- L1: No hardcoded secrets
- L1: No eval usage
- L2: Sanitize user input

## Architecture (N rules)
- L1: Repository pattern for data access
- L2: Service layer for business logic

## Code Quality (N rules)
- L2: No console.log in production
- L3: Functions under 25 lines
- L3: Named exports preferred

## Testing (N rules)
- L1: No .only in tests
- L3: Test file recommended
```

Follow with `AskUserQuestion`: Approve rules | Modify before saving | Cancel

---

## Constitution Summary

After creation or update:

```markdown
📜 Constitution [Created/Updated]

File: CONSTITUTION.md
Total Rules: [N]

Categories:
├── Security: [N] rules
├── Architecture: [N] rules
├── Code Quality: [N] rules
├── Testing: [N] rules
└── [Custom]: [N] rules

Level Distribution:
- L1 (Must, Autofix): [N]
- L2 (Should, Manual): [N]
- L3 (May, Advisory): [N]

Integration Points:
- /start:validate constitution - Check compliance
- /start:implement - Active enforcement
- /start:review - Code review checks
- /start:specify - SDD alignment
```

---

## Update Mode Options

When constitution exists, present:
- Add new rules (to existing or new category)
- Modify existing rules
- Remove rules
- View current constitution
