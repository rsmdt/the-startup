# Gold-Standard Skill Conventions

The definitive reference for skill structure. Apply when creating, converting, or auditing skills.

---

## Skill Anatomy

skills/[skill-name]/
├── SKILL.md           # Core logic (always loaded, <500 lines, <25 KB)
├── reference/         # Advanced protocols (loaded on demand)
├── templates/         # Document templates
├── examples/          # Real-world scenarios
└── validation.md      # Quality checklists

**Progressive disclosure**: Only `SKILL.md` is loaded into context. Reference files, templates, and examples are loaded when the skill explicitly requests them — keeping context lean for simple invocations.

---

## Frontmatter

Required fields:

name: kebab-case-name                        // Lowercase, numbers, hyphens (max 64 chars)
description: What it does and when to use it  // Max 1024 chars

Optional fields:

allowed-tools: Task, Bash, Read              // Tools without permission prompts
user-invocable: true                         // false = hides from / menu
argument-hint: "description of arguments"    // Shown in / menu
disable-model-invocation: false              // true = only user can invoke
context: fork                                // Run in subagent
agent: Explore                               // Subagent type when context: fork

### Description Guidelines

- Explain WHAT the skill does AND WHEN to use it
- Include keywords users would naturally say
- Keep it focused on triggers, not implementation details
- Write in third person (injected into system prompt)
- NEVER describe the workflow in the description — agents will follow it as a shortcut and skip the body

### String Substitutions

$ARGUMENTS              // All arguments passed when invoking
$ARGUMENTS[0]           // First argument
${CLAUDE_SESSION_ID}    // Current session ID
!`shell command`        // Execute command, insert output (preprocessing)

### Security Note
Never combine `!`shell command`` preprocessing with `$ARGUMENTS` — this executes user input as a shell command at skill load time. Use `Bash()` in the Workflow section instead, where the AI mediates the execution.

---

## Skill Body: PICS + Workflow

Section order — each section is a `## ` heading:

## Persona               // Role and expertise frame
## Interface             // Data shapes, then fn signatures
## Constraints           // require {} and never {} blocks
## State                 // Concrete defaults with origin comments
## Reference Materials   // Optional — links to progressive disclosure files
## Workflow              // fn definitions + entry-point pipe chain

### Persona

Sets the AI's role and expertise frame. Keep enforcement rules out — those go in Constraints.

### Interface

Data shapes first, then function signatures. Inline enum values directly — no `type` aliases.

Finding {
  severity: CRITICAL | HIGH | MEDIUM | LOW
  title: String
  fix: String
}

fn gatherContext(target)       // forward declaration — body in Workflow
fn synthesize(findings)

### Constraints

Split into `require {}` (must do) and `never {}` (must not do). Move enforcement-worthy rules from Persona into `never {}`.

Constraints {
  require {
    Every finding must have a specific, implementable fix.
    Provide full file context to reviewers, not just diffs.
  }
  never {
    Review code yourself — always delegate to specialist agents.
    Present findings without actionable fix recommendations.
  }
}

### State

Concrete defaults with comments explaining which function populates them. No `infer()`.

State {
  target = $ARGUMENTS
  perspectives = []              // populated by gatherContext
  mode: Standard | Team          // chosen by user in selectMode
  findings: [Finding]            // collected from agents
}

### Reference Materials

Links to progressive disclosure files. Only include when the skill has a `reference/` directory.

### Workflow

Define each step as `fn` (definition = not executed), then chain in an entry-point function without `fn` (execution = runs immediately). Entry-point function name matches skill name.

fn gatherContext(target) {
  match (target) {
    /^\d+$/       => gh pr diff $target
    "staged"      => git diff --cached
    default       => git diff main...$target
  }
}

fn synthesize(findings) {
  findings |> deduplicate |> sort |> buildSummaryTable
}

review(target) {
  gatherContext(target) |> selectMode |> launchReviews |> synthesize |> nextSteps
}

**The fn convention**:

| Pattern | Meaning | Where |
|---------|---------|-------|
| `fn foo()` | Forward declaration (no body) | Interface section |
| `fn foo() { ... }` | Definition (has body, not executed) | Workflow section |
| `foo() { ... }` | Execution (runs immediately) | Entry point |

Functions can contain nested `Constraints {}` blocks for per-step rules.

---

## Skill Types

| Type | Purpose | Structure |
|------|---------|-----------|
| **Technique** | How-to guide with steps | Workflow + examples |
| **Pattern** | Mental model or approach | Principles + when to apply |
| **Reference** | API/syntax documentation | Tables + code samples |
| **Coordination** | Orchestrate multiple agents | Perspectives + synthesis |

---

## Discipline-Enforcing Skills

Skills that enforce rules (TDD, verification) need special attention:
- Use strong language: "YOU MUST", "No exceptions", "Never"
- Add rationalization counters (excuse → reality table)
- Add Red Flags section listing rationalizations that indicate violation
- Test with 3+ combined pressure scenarios (see testing-with-subagents.md)
- See persuasion-principles.md for research on language patterns

---

## Transformation Checklist

When converting an existing skill to these conventions:

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
- [ ] Confirm no content/logic was lost in transformation

---

## Canonical Example

See `../examples/canonical-skill.md` for a fully annotated skill demonstrating all conventions.
