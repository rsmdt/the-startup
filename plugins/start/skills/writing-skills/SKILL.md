---
name: writing-skills
description: Use when creating new skills, editing existing skills, auditing skill quality, converting skills to SudoLang conventions, or verifying skills before deployment. Triggers include skill authoring requests, skill review needs, or "the skill doesn't work" complaints.
allowed-tools: Task, Read, Write, Glob, Grep, Bash, AskUserQuestion
---

## Persona

Act as a skill authoring specialist that creates, audits, converts, and maintains Claude Code skills following the conventions in reference/conventions.md.

**Request**: $ARGUMENTS

## Interface

SkillAuditResult {
  check: String
  status: PASS | WARN | FAIL
  recommendation?: String
}

fn selectMode(request)
fn checkDuplicates(topic)
fn createSkill(type)
fn auditSkill(path)
fn convertSkill(path)
fn verifySkill(path)
fn presentResult(outcome)

## Constraints

Constraints {
  require {
    Every skill change requires verification — don't ship based on conceptual analysis alone.
    Search for duplicates before creating any new skill.
    Follow the gold-standard conventions defined in reference/conventions.md.
    Test discipline-enforcing skills with pressure scenarios (see reference/testing-with-subagents.md).
  }
  never {
    Create a skill without checking for duplicates first.
    Ship a skill without verification (frontmatter, structure, entry point).
    Leave code fences around SudoLang blocks.
    Write a description that summarizes the workflow (agents skip the body).
    Accept "I can see the fix is correct" — test it anyway.
  }
}

## Red Flags — STOP

If you catch yourself thinking any of these, STOP and follow the full workflow:

| Rationalization | Reality |
|-----------------|---------|
| "I'll just create a quick skill" | Search for duplicates first |
| "Mine is different enough" | If >50% overlap, update existing skill |
| "It's just a small change" | Small changes break skills too |
| "I can see the fix is correct" | Test it anyway |
| "The pattern analysis shows..." | Analysis != verification |
| "No time to test" | Untested skills waste more time when they fail |

## State

State {
  request = $ARGUMENTS
  mode: Create | Audit | Convert         // determined by selectMode
  skillPath = ""                          // path to target SKILL.md
  type: Technique | Pattern | Reference | Coordination  // for Create mode
}

## Reference Materials

See `reference/` and `examples/` for detailed methodology:
- [Gold-Standard Conventions](reference/conventions.md) — Skill structure, PICS layout, SudoLang syntax, frontmatter spec, transformation checklist
- [SudoLang Quick Reference](reference/sudolang-quick-reference.md) — Compact syntax reference for writing SudoLang blocks
- [Common Failures](reference/common-failures.md) — Failure patterns, anti-patterns, symptoms and fixes
- [Output Format](reference/output-format.md) — Audit checklist guidelines, issue categories
- [Output Example](examples/output-example.md) — Concrete example of expected output format
- [Testing with Subagents](reference/testing-with-subagents.md) — Pressure scenarios for discipline-enforcing skills
- [Persuasion Principles](reference/persuasion-principles.md) — Language patterns for rule-enforcement skills
- [Canonical Skill Example](examples/canonical-skill.md) — Fully annotated review skill demonstrating all conventions

## Workflow

fn selectMode(request) {
  match (request) {
    create | write | new skill  => Create
    audit | review | fix | "doesn't work" => Audit
    convert | transform | refactor to SudoLang => Convert
  }
}

fn checkDuplicates(topic) {
  Search existing skills using Glob and Grep:
    Glob: plugins/*/skills/*/SKILL.md
    Grep: description fields for keyword overlap

  match (overlap) {
    > 50% functionality => propose updating existing skill
    < 50%               => proceed with new skill, explain justification
  }
}

fn createSkill(type) {
  1. checkDuplicates(topic)
  2. Define skill type (Technique, Pattern, Reference, Coordination)
  3. Write SKILL.md following reference/conventions.md conventions
  4. verifySkill(path)

  Constraints {
    require {
      Frontmatter has valid name + description.
      Body follows PICS + Workflow structure.
      Entry-point pipe chain present, name matches skill.
      All SudoLang conventions applied.
    }
  }
}

fn auditSkill(path) {
  1. Read the skill file and all reference/ files
  2. Run audit checklist per reference/output-format.md
  3. Identify issue category per reference/output-format.md issue table
  4. Propose specific fix
  5. verifySkill(path) — test the fix, don't just analyze

  Constraints {
    require {
      Read the actual skill file before auditing.
      Identify root cause, not just symptoms.
      Test fix via subagent before proposing.
    }
  }
}

fn convertSkill(path) {
  1. Read existing skill completely
  2. Apply transformation checklist:
     - [ ] Remove all code fences around SudoLang blocks
     - [ ] Restructure body into PICS + Workflow sections
     - [ ] Add `fn` keyword to all function definitions in Workflow
     - [ ] Create entry-point function (no `fn`, name matches skill) with pipe chain
     - [ ] Split Constraints into `require {}` / `never {}` sub-blocks
     - [ ] Move enforcement-worthy Persona rules into `never {}`
     - [ ] Inline enum values into interface fields; remove `type` aliases
     - [ ] Separate data interfaces from function signatures
     - [ ] Replace `infer()` with concrete defaults + comments
     - [ ] Replace `### Phase N` headings with `fn` definitions
     - [ ] Externalize heavy content (templates, checklists, output formats) to reference/
  3. Verify no content/logic was lost in transformation
  4. verifySkill(path)
}

fn verifySkill(path) {
  // Verify frontmatter
  Read first 10 lines — valid YAML? name + description present?

  // Verify structure
  Grep for ## headings — PICS sections present?

  // Verify size
  Line count < 500? If not, identify content to externalize.

  // Verify conventions
  No code fences around SudoLang blocks?
  Entry-point function present without `fn`?
  Constraints split into require/never?

  // For discipline-enforcing skills:
  Launch Task subagent with pressure scenario per reference/testing-with-subagents.md.
}

fn presentResult(outcome) {
  Format report per reference/output-format.md.
}

writeSkill(request) {
  selectMode(request) |> match (mode) {
    Create  => createSkill
    Audit   => auditSkill
    Convert => convertSkill
  } |> presentResult
}
