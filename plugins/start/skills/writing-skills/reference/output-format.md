# Output Format Reference

Status report templates for skill creation, audit, and conversion.

---

## Skill Created

```
Skill Created

Path: [full path to SKILL.md]
Name: [skill name]
Type: [Technique/Pattern/Reference/Coordination]
Size: [line count] lines

Verification:
- [ ] Duplicate check passed
- [ ] Frontmatter valid
- [ ] PICS structure complete
- [ ] SudoLang conventions applied
- [ ] Entry-point pipe chain present
- [ ] Output format defined
- [ ] [For discipline skills] Pressure test passed

Ready for deployment: [Yes/No]
```

---

## Skill Audited

```
Skill Audit Complete

Skill: [name]
Path: [path]

Issues Found: [N]
| Issue | Category | Severity | Recommendation |
|-------|----------|----------|----------------|
| [issue] | [category] | [High/Medium/Low] | [fix] |

Root Cause: [Why agents don't follow this skill]

Recommended Actions:
1. [Action 1]
2. [Action 2]

Verified: [Yes - tested fix / No - conceptual only]
```

---

## Skill Converted

```
Skill Converted

Skill: [name]
Path: [path]

Transformation Applied:
- [ ] Code fences removed
- [ ] PICS structure applied
- [ ] Constraints split into require/never
- [ ] fn declarations + entry-point pipe chain
- [ ] Enums inlined, type aliases removed
- [ ] State initialized with concrete defaults
- [ ] Heavy content externalized to reference/

Before: [line count] lines
After: [line count] lines
Reference files: [count] files ([total lines] lines)

Content verification: [All logic preserved / Issues found]
```

---

## Audit Checklist

| Check | Question | Pass/Fail |
|-------|----------|-----------|
| **Frontmatter** | Valid YAML? name + description present? | |
| **Description** | Explains what + when? Includes keywords? No workflow? | |
| **PICS Structure** | Persona, Interface, Constraints, State sections present? | |
| **SudoLang** | No code fences? require/never? fn convention? | |
| **Entry Point** | Pipe chain present? Name matches skill? No `fn` keyword? | |
| **Tools** | Tool usage referenced in workflow functions? | |
| **Output** | Clear output format defined or referenced? | |
| **Size** | Under 500 lines? Heavy content in reference/? | |

---

## Issue Categories

| Symptom | Category | Fix Approach |
|---------|----------|--------------|
| "Agent doesn't follow it" | Execution gap | Add imperative fn definitions |
| "Agent skips sections" | Discovery problem | Restructure for progressive disclosure |
| "Agent does wrong thing" | Ambiguity | Add explicit constraints in never {} |
| "Agent can't find skill" | CSO problem | Improve description keywords |
| "Skill is too long" | Bloat | Extract to reference/ directory |
| "Agent treats SudoLang as code" | Code fence problem | Remove ``` wrappers |
