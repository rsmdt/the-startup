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
- [x] PICS + Workflow structure complete
- [x] Markdown conventions applied (Always/Never, ### headings, TypeScript types)
- [x] Entry Point present (non-linear workflow)
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
| Perspectives loaded inconsistently | Discovery problem | Medium | Add explicit Read instruction in step 3 |
| Missing entry for "staged" target | Ambiguity | Low | Add staged to step 1 match cases |

Root Cause: Output template wrapped in code fences in reference/ — AI treats as documentation rather than format to produce

Recommended Actions:
1. Move report template to examples/output-example.md as rendered markdown
2. Add explicit "Read examples/output-example.md" step in synthesis workflow

Verified: Yes - tested fix with sample review invocation

---

## Skill Converted

Skill Converted

Skill: debug
Path: plugins/start/skills/debug/SKILL.md

Transformation Applied:
- [x] PICS + Workflow structure applied
- [x] Constraints converted to Always/Never lists
- [x] Workflow converted to numbered ### headings
- [x] Entry Point removed (linear workflow)
- [x] Types converted to valid TypeScript (string, number, Type[])
- [x] Always/Never deduplicated — no mirrored rules
- [x] State merged into Interface
- [x] Heavy content externalized to reference/

Before: 415 lines
After: 121 lines
Reference files: 2 files (109 lines)

Content verification: All logic preserved
