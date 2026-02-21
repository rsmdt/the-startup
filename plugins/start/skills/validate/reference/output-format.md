# Output Format Reference

Report template, assessment levels, and verdict-based next steps for the validate skill.

---

## Report Template

```markdown
## Validation: [target]

**Mode**: [Spec | File | Drift | Constitution | Comparison | Understanding]
**Assessment**: ✅ Excellent | 🟢 Good | 🟡 Needs Attention | 🔴 Critical

### Summary

| Perspective | Pass | Warn | Fail |
|-------------|------|------|------|
| ✅ Completeness | X | X | X |
| 🔗 Consistency | X | X | X |
| 📍 Alignment | X | X | X |
| 📐 Coverage | X | X | X |
| 📊 Drift | X | X | X |
| 📜 Constitution | X | X | X |
| **Total** | X | X | X |

*🔴 Failures (Must Fix)*

| ID | Finding | Recommendation |
|----|---------|----------------|
| F1 | Brief title *(file:line)* | Fix recommendation *(issue description)* |

*🟡 Warnings (Should Fix)*

| ID | Finding | Recommendation |
|----|---------|----------------|
| W1 | Brief title *(file:line)* | Fix recommendation *(issue description)* |

*✅ Passes*

| Perspective | Verified |
|-------------|----------|
| Completeness | All sections populated, no TODO markers |
| Consistency | Terminology consistent across docs |

### Verdict

[What was validated and key conclusions]
```

---

## Assessment Levels

| Level | Condition |
|-------|-----------|
| ✅ Excellent | No failures, no warnings |
| 🟢 Good | No failures, 1-3 warnings |
| 🟡 Needs Attention | No failures, 4+ warnings |
| 🔴 Critical | Any failures present |

---

## Verdict-Based Next Steps

Use `AskUserQuestion` based on validation mode and findings:

**If Constitution L1/L2 Violations:**
- "Apply autofixes (L1)" (Recommended)
- "Show me the violations"
- "Skip constitution checks"

**If Drift Detected:**
- "Acknowledge and continue" (log drift, proceed)
- "Update implementation" (implement missing, remove extra)
- "Update specification" (modify spec to match reality)
- "Defer decision" (mark for later review)

**If Spec/File Issues:**
- "Address failures first"
- "Show detailed findings"
- "Continue anyway"
