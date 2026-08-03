# Output Format Reference

Guidelines for the review report format. See `examples/output-example.md` for a concrete rendered example.

---

## Breaking Changes Section

Render this section immediately after the verdict and before every severity table, but only when at least one finding has `breaking` set.
Omit the whole section when there are none - an empty "no breaking changes" block trains readers to skip past it.

```markdown
### ⚠️ Breaking Changes

**[ID] Title** *(location)*
- **Breaks**: what depends on this contract
- **Migration**: what consumers must do to adapt
```

Breaking findings still appear in their severity table so the counts stay honest.
Prefix the table row with ⚠️ and keep the row terse - the detail already lives in the section above.

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
- "Accept the breaking change and draft the migration note" (only when breaking findings exist)

**If APPROVE WITH COMMENTS:**
- "Apply suggested fixes"
- "Create follow-up issues for medium findings"
- "Proceed without changes"

**If APPROVE:**
- "Add to PR comments (if PR review)"
- "Done"
