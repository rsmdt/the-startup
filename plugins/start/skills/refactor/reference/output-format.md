# Output Format Reference

Templates for baseline, analysis, error recovery, and completion reports.

---

## Baseline Template

```markdown
📊 Refactoring Baseline

Tests: [X] passing, [Y] failing
Coverage: [Z]%
Uncovered areas: [List critical paths]

Baseline Status: [READY / TESTS FAILING / COVERAGE GAP]
```

## Analysis Summary Template

```markdown
## Refactoring Analysis: [target]

### Summary

| Perspective | High | Medium | Low |
|-------------|------|--------|-----|
| 🔧 Code Smells | X | X | X |
| 🔗 Dependencies | X | X | X |
| 🧪 Test Coverage | X | X | X |
| 🏗️ Patterns | X | X | X |
| ⚠️ Risk | X | X | X |
| **Total** | X | X | X |

*🔴 High Impact Issues*

| ID | Finding | Remediation | Risk |
|----|---------|-------------|------|
| H1 | Brief title *(file:line)* | Specific technique *(problem description)* | Risk level |

*🟡 Medium Impact Issues*

| ID | Finding | Remediation | Risk |
|----|---------|-------------|------|
| M1 | Brief title *(file:line)* | Specific technique *(problem description)* | Risk level |

*⚪ Low Impact Issues*

| ID | Finding | Remediation | Risk |
|----|---------|-------------|------|
| L1 | Brief title *(file:line)* | Specific technique *(problem description)* | Risk level |
```

## Error Recovery Template

```markdown
⚠️ Refactoring Failed

Refactoring: [Name]
Reason: Tests failing

Reverted: ✓ Working state restored

Options:
1. Try alternative approach
2. Add missing tests first
3. Skip this refactoring
4. Get guidance

Awaiting your decision...
```

## Completion Summary Template

```markdown
## Refactoring Complete: [target]

**Status**: [Complete / Partial - reason]

### Before / After

| File | Before | After | Technique |
|------|--------|-------|-----------|
| billing.ts | 75-line method | 4 functions, 20 lines each | Extract Method |

### Verification

- Tests: [X] passing (baseline: [Y])
- Behavior: Preserved ✓
- Coverage: [Z]% (baseline: [W]%)

### Quality Improvements

- [Improvement 1]
- [Improvement 2]

### Skipped

- [file:line] - [reason] (e.g., no test coverage, user declined)
```

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
