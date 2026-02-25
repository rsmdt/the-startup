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
### 3. Security Review

Search for database queries using Grep.
Check each query for string interpolation (vulnerability).
Report findings in Finding format.
```

---

## Failure: Missing Tool Instructions

**Symptom:** Agent has tools listed but doesn't use them

**Bad:** `allowed-tools: Task, Bash, Read` with no mention of tools in body.

**Good:** Workflow steps reference tools explicitly:
```markdown
### 1. Gather Context

match (target) {
  /^\d+$/  => gh pr diff $target    // uses Bash
  "staged" => git diff --cached     // uses Bash
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

## Failure: No Clear Starting Point

**Symptom:** Skill is comprehensive reference but agent doesn't know where to start

**Fix for linear workflows:** Use numbered `### N. Step Name` headings — execution order is implicit.

**Fix for non-linear workflows:** Add `### Entry Point` with a match block:
```markdown
### Entry Point

match (mode) {
  Create  => steps 2, 3, 7
  Audit   => steps 4, 7
  Convert => steps 5, 7
}
```

---

## Failure: Code Fences Around Skill Syntax

**Symptom:** Agent treats Interface, Constraints, or Workflow blocks as inert code rather than following them

**Fix:** Remove all ` ``` ` wrappers around Interface definitions, Constraints, State, and Workflow content. These blocks are instructions, not code samples.

---

## Anti-Patterns

| Anti-Pattern | Why It Fails | Instead |
|--------------|--------------|---------|
| Create without duplicate check | Fragments knowledge | Search first |
| 400+ line SKILL.md | Agents skim, miss sections | Extract to reference/ |
| Audit without reading file | Miss actual issues | Always read the skill |
| Fix without testing | May not address root cause | Verify fix works |
| Declarative-only workflow | Agents don't execute | Use imperative ### step headings |
| Description with workflow | Agents skip body | Triggers only in description |
| `type` aliases for enums | Extra indirection | Inline values into interface fields |
| Mirrored Always/Never rules | Wastes tokens | Each rule once, most natural framing |
