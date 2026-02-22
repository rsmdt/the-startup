# Output Format Reference

Guidelines for writing-skills output. See `examples/output-example.md` for concrete rendered examples of all three report types (created, audited, converted).

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
