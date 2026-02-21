# Output Format Reference

Guidelines for refactoring output. See `examples/output-example.md` for concrete rendered examples of all four report types (baseline, analysis, error recovery, completion).

---

## Table Column Guidelines

- **ID**: Impact letter + number (H1 = High #1, M1 = Medium #1, L1 = Low #1)
- **Finding**: Brief title + location in italics
- **Remediation**: Specific refactoring technique + problem description in italics
- **Risk**: Assessment of potential complications

---

## Next Steps Options

**After analysis:**
- "Document and proceed" — Save plan to `docs/refactor/[NNN]-[name].md`, then execute
- "Proceed without documenting" — Execute refactorings directly
- "Cancel" — Abort refactoring

**After completion:**
- "Commit these changes"
- "Run full test suite"
- "Address skipped items (add tests first)"
- "Done"
