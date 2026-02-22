# Common Skill Failures

Failure patterns to watch for when creating, auditing, or converting skills.

---

## Failure: Declarative Instead of Imperative

**Symptom:** Agent reads skill but doesn't execute workflow

**Bad (declarative):**
```markdown
The security review should check for SQL injection vulnerabilities.
```

**Good (imperative):**
```markdown
fn securityReview() {
  Search for database queries using Grep.
  Check each query for string interpolation (vulnerability).
  Report findings in Finding format.
}
```

---

## Failure: Missing Tool Instructions

**Symptom:** Agent has tools listed but doesn't use them

**Bad:** `allowed-tools: Task, Bash, Read` with no mention of tools in body.

**Good:** Workflow functions reference tools explicitly:
```
fn gatherContext(target) {
  match (target) {
    /^\d+$/  => gh pr diff $target    // uses Bash
    "staged" => git diff --cached     // uses Bash
  }
}
```

---

## Failure: Description Summarizes Workflow

**Symptom:** Agent follows description shortcut, skips skill body

**Bad:**
```yaml
description: Coordinates code review by launching four specialized subagents then synthesizing findings
```

**Good:**
```yaml
description: Use when reviewing code changes, PRs, or staged files. Handles security, performance, quality, and test analysis.
```

---

## Failure: No Execution Entry Point

**Symptom:** Skill is comprehensive reference but agent doesn't know where to start

**Fix:** Add entry-point function (no `fn` keyword) with pipe chain:
```
review(target) {
  gatherContext(target) |> selectMode |> launchReviews |> synthesize |> nextSteps
}
```

---

## Failure: Code Fences Around SudoLang

**Symptom:** Agent treats SudoLang as inert code block, doesn't follow it

**Fix:** Remove all ` ``` ` wrappers around Interface, Constraints, State, and Workflow blocks.

---

## Anti-Patterns

| Anti-Pattern | Why It Fails | Instead |
|--------------|--------------|---------|
| Create without duplicate check | Fragments knowledge | Search first |
| 400+ line SKILL.md | Agents skim, miss sections | Extract to reference/ |
| Audit without reading file | Miss actual issues | Always read the skill |
| Fix without testing | May not address root cause | Verify fix works |
| Declarative-only workflow | Agents don't execute | Use imperative fn definitions |
| Description with workflow | Agents skip body | Triggers only in description |
| `type` aliases for enums | Extra indirection | Inline values into interface fields |
| `infer()` in State | Opaque initialization | Concrete defaults with comments |
| `### Phase N` headings | Markdown structure, not SudoLang | fn definitions + pipe chain |
