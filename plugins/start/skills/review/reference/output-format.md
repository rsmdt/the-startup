# Output Format Reference

Full report template, table guidelines, and verdict-based next steps for the review skill.

---

## Report Template

```markdown
## Code Review: [target]

**Verdict**: 🔴 REQUEST CHANGES | 🟡 APPROVE WITH COMMENTS | ✅ APPROVE

### Summary

| Category | Critical | High | Medium | Low |
|----------|----------|------|--------|-----|
| 🔐 Security | X | X | X | X |
| 🔧 Simplification | X | X | X | X |
| ⚡ Performance | X | X | X | X |
| 📝 Quality | X | X | X | X |
| 🧪 Testing | X | X | X | X |
| **Total** | X | X | X | X |

*🔴 Critical & High Findings (Must Address)*

| ID | Finding | Remediation |
|----|---------|-------------|
| C1 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |
| H1 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |

#### Code Examples for Critical Fixes

**[C1] Title**
```language
// Before
old code

// After
new code
```

*🟡 Medium Findings (Should Address)*

| ID | Finding | Remediation |
|----|---------|-------------|
| M1 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |

*⚪ Low Findings (Consider)*

| ID | Finding | Remediation |
|----|---------|-------------|
| L1 | Brief title *(file:line)* | Specific fix recommendation *(concise issue description)* |

### Strengths

- ✅ [Positive observation with specific code reference]
- ✅ [Good patterns noticed]

### Verdict Reasoning

[Why this verdict was chosen based on findings]
```

---

## Table Column Guidelines

- **ID**: Severity letter + number (C1 = Critical #1, H2 = High #2, M1 = Medium #1, L1 = Low #1)
- **Finding**: Brief title + location in italics (e.g., `Missing null check *(auth/service.ts:42)*`)
- **Remediation**: Fix recommendation + issue context in italics (e.g., `Add null guard *(query result accessed without check)*`)

## Code Example Rules

- **REQUIRED** for all Critical findings (before/after style)
- **Include** for High findings when the fix is non-obvious
- **Omit** for Medium/Low findings (table-only format)

---

## Verdict-Based Next Steps

Use `AskUserQuestion` with options based on verdict:

**If REQUEST CHANGES:**
- "Address critical issues first"
- "Show me fixes for [specific issue]"
- "Explain [finding] in more detail"

**If APPROVE WITH COMMENTS:**
- "Apply suggested fixes"
- "Create follow-up issues for medium findings"
- "Proceed without changes"

**If APPROVE:**
- "Add to PR comments (if PR review)"
- "Done"
