# Example Writing-Skills Output

## Skill Created

Skill Created

Path: plugins/start/skills/deploy/SKILL.md
Name: deploy
Type: Coordination
Size: 142 lines

Verification:
- [x] Duplicate check passed
- [x] Frontmatter valid
- [x] PICS structure complete
- [x] SudoLang conventions applied
- [x] Entry-point pipe chain present
- [x] Output format defined
- [ ] [For discipline skills] Pressure test passed

Ready for deployment: Yes

---

## Skill Audited

Skill Audit Complete

Skill: review
Path: plugins/start/skills/review/SKILL.md

Issues Found: 3

| Issue | Category | Severity | Recommendation |
|-------|----------|----------|----------------|
| Output format not followed | Execution gap | High | Move template from code-fenced reference to concrete example |
| Perspectives loaded inconsistently | Discovery problem | Medium | Add explicit Read instruction in synthesize() |
| Missing entry for "staged" target | Ambiguity | Low | Add staged to gatherContext match cases |

Root Cause: Output template wrapped in code fences in reference/ — AI treats as documentation rather than format to produce

Recommended Actions:
1. Move report template to examples/output-example.md as rendered markdown
2. Add explicit "Read examples/output-example.md" step in synthesize()

Verified: Yes - tested fix with sample review invocation

---

## Skill Converted

Skill Converted

Skill: debug
Path: plugins/start/skills/debug/SKILL.md

Transformation Applied:
- [x] Code fences removed
- [x] PICS structure applied
- [x] Constraints split into require/never
- [x] fn declarations + entry-point pipe chain
- [x] Enums inlined, type aliases removed
- [x] State initialized with concrete defaults
- [x] Heavy content externalized to reference/

Before: 415 lines
After: 121 lines
Reference files: 2 files (109 lines)

Content verification: All logic preserved
