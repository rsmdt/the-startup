# Output Format Reference

Guidelines for writing-skills output. See `examples/output-example.md` for concrete rendered examples of all three report types (created, audited, converted).

---

## Audit Checklist

| Check | Question | Pass/Fail |
|-------|----------|-----------|
| **Frontmatter** | Valid YAML? name + description present? | |
| **Description** | Explains what + when? Includes keywords? No workflow? | |
| **PICS Structure** | Persona, Interface, Constraints, Workflow sections present? | |
| **Constraints** | Always/Never lists? No mirrored duplicates? | |
| **Workflow** | Numbered ### headings? Entry Point only if non-linear? | |
| **Types** | Valid TypeScript types (string, number, boolean, Type[])? | |
| **Tools** | Tool usage referenced in workflow steps? | |
| **Output** | Clear output format defined or referenced? | |
| **Size** | Under 500 lines? Heavy content in reference/? | |

---

## Issue Categories

| Symptom | Category | Fix Approach |
|---------|----------|--------------|
| "Agent doesn't follow it" | Execution gap | Add imperative ### step headings |
| "Agent skips sections" | Discovery problem | Restructure for progressive disclosure |
| "Agent does wrong thing" | Ambiguity | Add explicit constraints in Never list |
| "Agent can't find skill" | CSO problem | Improve description keywords |
| "Skill is too long" | Bloat | Extract to reference/ directory |
| "Agent treats blocks as code" | Code fence problem | Remove ``` wrappers around skill syntax |
