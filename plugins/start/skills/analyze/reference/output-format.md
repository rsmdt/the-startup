# Output Format Reference

Cycle summary template, documentation format, and analysis summary for the analyze skill.

---

## Cycle Summary Template

Present after each discovery cycle:

```markdown
🔍 Discovery Cycle [N] Complete

Area: [Analysis area]
Agents Launched: [N]

Key Findings:
1. [Finding with evidence]
2. [Finding with evidence]
3. [Finding with evidence]

Patterns Identified:
- [Pattern name]: [Brief description]

Documentation Created/Updated:
- docs/[category]/[file.md]

Questions for Clarification:
1. [Question about ambiguous finding]

Should I continue to [next area] or investigate [finding] further?
```

---

## Analysis Summary Template

Present at the end of all cycles:

```markdown
## Analysis: [area]

### Discoveries

**[Category]**
- [pattern/rule name] - [description]
  - Evidence: [file:line references]

### Documentation

- [docs/path/file.md] - [what was documented]

### Open Questions

- [unresolved items for future investigation]
```

---

## Next Steps Options

Use `AskUserQuestion` based on analysis state:

**After each cycle:**
- "Continue to next area"
- "Investigate [finding] further"
- "Persist findings to docs/"
- "Complete analysis"

**After final summary:**
- "Save documentation to docs/"
- "Skip documentation"
- "Export as markdown"
