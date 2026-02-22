# Output Format Reference

Guidelines for validation output. See `examples/output-example.md` for a concrete rendered example.

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
