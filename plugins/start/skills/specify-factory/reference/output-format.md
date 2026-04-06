# Output Format Reference

Guidelines for specify-factory output. See `examples/output-example.md` for a concrete rendered example.

---

## Status Report

After each phase, report:

| Field | Description |
|-------|-------------|
| Units | Count, with dependency summary |
| Scenarios | Count per unit, total, coverage % |
| Manifest | Execution groups, threshold, max_iterations |
| Gaps | Requirements without units, units without scenarios |

## Next-Step Options

Use `AskUserQuestion` based on state:

**After unit decomposition:**
- "Approve units"
- "Adjust decomposition"
- "Add/remove units"

**After scenario generation:**
- "Approve scenarios"
- "Edit scenarios"
- "Add missing scenarios"
- "Regenerate"

**After manifest assembly:**
- "Approve manifest"
- "Adjust execution order"
- "Change settings"

**After validation:**
- "Finalize"
- "Revisit units"
- "Revisit scenarios"
- "Adjust manifest"
