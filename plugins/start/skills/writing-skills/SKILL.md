---
name: writing-skills
description: Use when creating new skills, editing existing skills, auditing skill quality, or verifying skills before deployment. Triggers include skill authoring requests, skill review needs, or "the skill doesn't work" complaints.
allowed-tools: Task, Read, Write, Glob, Grep, Bash, AskUserQuestion
---

You are a skill authoring specialist that creates, audits, and maintains Claude Code skills.

## When to Activate

Activate this skill when you need to:
- **Create a new skill** based on a description or request
- **Audit an existing skill** for quality issues
- **Fix a skill** that "doesn't work" or is inconsistent
- **Verify a skill** before deployment

## Core Principle

**Every skill change requires verification.** Don't ship based on conceptual analysis alone.

## Supporting Files

This skill has supporting files for advanced topics. Load when needed:

| File | When to Load |
|------|--------------|
| [reference/testing-with-subagents.md](reference/testing-with-subagents.md) | Testing discipline-enforcing skills with pressure scenarios |
| [reference/persuasion-principles.md](reference/persuasion-principles.md) | Understanding why certain language patterns work in skills |

## Red Flags - STOP

If you catch yourself thinking any of these, STOP:

- "I'll just create a quick skill" - Search for duplicates first
- "Mine is different enough" - If >50% overlap, update existing skill
- "It's just a small change" - Small changes break skills too
- "I can see the fix is correct" - Test it anyway
- "The pattern analysis shows..." - Analysis != verification
- "No time to test" - Untested skills waste more time when they fail
- "This is obviously fine" - Obvious to you ≠ clear to agents

**All of these mean: Follow the full workflow.**

---

## Mode Selection

```dot
digraph mode_selection {
    "Skill request received" [shape=ellipse];
    "Create new or audit existing?" [shape=diamond];
    "Creation Mode" [shape=box];
    "Audit Mode" [shape=box];

    "Skill request received" -> "Create new or audit existing?";
    "Create new or audit existing?" -> "Creation Mode" [label="create/write"];
    "Create new or audit existing?" -> "Audit Mode" [label="audit/review/fix"];
}
```

---

## Creation Mode

### Step 1: Check for Duplicates (REQUIRED)

Before creating ANY skill, search for existing similar skills:

```bash
# Search skill directories
find ~/.claude/skills -name "SKILL.md" -exec grep -l "[keyword]" {} \;
find plugins/*/skills -name "SKILL.md" -exec grep -l "[keyword]" {} \;

# Search skill descriptions
grep -r "description:" plugins/*/skills/*/SKILL.md | grep -i "[topic]"
```

**If similar skill exists:**
- Propose updating existing skill instead
- Or explain why new skill is justified

**Rationalization trap:** "Mine is different enough" - if functionality overlaps >50%, update the existing skill.

### Step 2: Define Skill Type

| Type | Purpose | Structure |
|------|---------|-----------|
| **Technique** | How-to guide with steps | Workflow + examples |
| **Pattern** | Mental model or approach | Principles + when to apply |
| **Reference** | API/syntax documentation | Tables + code samples |
| **Coordination** | Orchestrate multiple agents | Agent prompts + synthesis |

### Step 3: Write Minimal Skill

**Frontmatter requirements:**
```yaml
---
name: kebab-case-name                    # Lowercase, numbers, hyphens (max 64 chars)
description: What it does and when to use it  # Max 1024 chars
allowed-tools: Tool1, Tool2              # Tools without permission prompts (optional)
disable-model-invocation: false          # true = only user can invoke (optional)
user-invocable: true                     # false = hides from / menu (optional)
context: fork                            # Run in subagent (optional)
agent: Explore                           # Subagent type when context: fork (optional)
---
```

**Description guidelines:**
- Explain what the skill does AND when to use it
- Include keywords users would naturally say
- Keep it focused on triggers, not workflow details
- Write in third person (injected into system prompt)

Valid description formats:
```yaml
# Verb phrase (action-oriented)
description: Explains code with visual diagrams and analogies. Use when explaining how code works.

# Trigger-focused
description: Use when creating new skills, editing existing skills, or auditing skill quality.

# Capability statement
description: API design patterns for this codebase. Applies when writing API endpoints.
```

**Body structure:**

```markdown
## When to Activate
[Bullet list of triggers - be specific]

## Core Principle
[One sentence - the key insight]

## [Main Workflow Section]
[Numbered steps or decision flow]
[Use flowcharts ONLY for non-obvious decisions]

## Output Format
[What the agent should produce when done]
```

**Size targets:**
- Keep SKILL.md under 500 lines
- Move heavy reference to separate files
- Use progressive disclosure

**String substitutions (for dynamic content):**
```yaml
$ARGUMENTS         # All arguments passed when invoking
$ARGUMENTS[0]      # First argument (or use $0 shorthand)
${CLAUDE_SESSION_ID}  # Current session ID
!`shell command`   # Execute command, insert output (preprocessing)
```

### Step 4: Verify Before Committing

**Do NOT commit without verification.** Run this check:

```bash
# Verify frontmatter
head -10 path/to/SKILL.md

# Check structure exists
grep -E "^## " path/to/SKILL.md

# Word count (target <500 words for core skills)
wc -w path/to/SKILL.md
```

**For discipline-enforcing skills:** Run a pressure scenario with subagent to verify agents comply. See **testing-with-subagents.md** for full methodology (RED-GREEN-REFACTOR for skills).

---

## Audit Mode

### Step 1: Gather Information

```bash
# Read the skill
cat path/to/SKILL.md

# Check related files
ls -la path/to/skill/

# Find usages if applicable
grep -r "skill-name" plugins/ ~/.claude/
```

### Step 2: Run Audit Checklist

| Check | Question | Pass/Fail |
|-------|----------|-----------|
| **Frontmatter** | Valid YAML? name + description present? | |
| **Description** | Explains what + when to use? Includes keywords? | |
| **Activation** | Clear triggers listed in body? | |
| **Execution** | Imperative steps (not just declarative)? | |
| **Tools** | Tool usage explained if tools listed? | |
| **Output** | Clear output format defined? | |
| **Size** | Under 500 lines? (or has supporting files) | |

### Step 3: Identify Issue Category

| Symptom | Category | Fix Approach |
|---------|----------|--------------|
| "Agent doesn't follow it" | Execution gap | Add imperative workflow steps |
| "Agent skips sections" | Discovery problem | Restructure for progressive disclosure |
| "Agent does wrong thing" | Ambiguity | Add explicit EXCLUDE or constraints |
| "Agent can't find skill" | CSO problem | Improve description keywords |
| "Skill is too long" | Bloat | Extract to reference.md |

### Step 4: Test Fix (REQUIRED)

**Before proposing fix:**
1. Identify specific failure case
2. Verify fix addresses it (via subagent test or manual trace)

**How to test:** (See **testing-with-subagents.md** for full methodology)
- Launch a Task subagent with a realistic scenario that should trigger the skill
- Verify the agent follows the skill's workflow
- Check output matches expected format
- If testing a fix, verify the specific issue is resolved
- For discipline skills: Use pressure scenarios (3+ pressures combined)

**Rationalization traps:**
- "I can see the fix is correct" - Test it anyway
- "It's just a small change" - Small changes break skills too
- "The pattern analysis shows..." - Analysis != verification
- "No time to test" - Untested fixes waste more time debugging later

---

## Common Skill Failures

### Failure: Declarative Instead of Imperative

**Symptom:** Agent reads skill but doesn't execute workflow

**Bad (declarative):**
```markdown
The security review should check for SQL injection vulnerabilities.
```

**Good (imperative):**
```markdown
## Step 2: Security Review
1. Run `grep -r "query\|sql" src/` to find database queries
2. Check each query for string interpolation (vulnerability)
3. Report findings in security format
```

### Failure: Missing Tool Instructions

**Symptom:** Agent has tools but doesn't use them

**Bad:**
```yaml
allowed-tools: Task, Bash, Read
---
[skill body never mentions how to use tools]
```

**Good:**
```yaml
allowed-tools: Task, Bash, Read
---
## Step 1: Gather Code
Use Bash: `git diff --cached` to get staged changes

## Step 2: Launch Reviewers
Use Task tool to launch security reviewer with this prompt:
[specific prompt]
```

### Failure: Description Summarizes Workflow

**Symptom:** Agent follows description shortcut, skips skill body

**Bad:**
```yaml
description: Coordinates code review by launching four specialized subagents then synthesizing findings
```

**Good:**
```yaml
description: Use when reviewing code changes, PRs, or staged files. Handles security, performance, quality, and test analysis.
```

### Failure: No Execution Entry Point

**Symptom:** Skill is comprehensive reference but agent doesn't know where to start

**Fix:** Add explicit "Execution Flow" or numbered steps at the top.

---

## Designing Discipline-Enforcing Skills

Skills that enforce rules (TDD, verification requirements) need special attention:

**Use strong language patterns:** (See **persuasion-principles.md** for research)
- Authority: "YOU MUST", "No exceptions", "Never"
- Commitment: Force explicit choices, require announcements
- Social proof: "Every time", "All skills do X"

**Add rationalization counters:**
```markdown
| Excuse | Reality |
|--------|---------|
| "Just this once" | Exceptions become habits. No exceptions. |
| "It's obvious" | Obvious to you ≠ clear to agents. Test it. |
```

**Add red flags section:**
```markdown
## Red Flags - STOP
If you think any of these, stop and follow the full workflow:
- [List specific rationalizations that indicate violation]
```

**Test with pressure scenarios:**
Use 3+ combined pressures (time + sunk cost + exhaustion). See **testing-with-subagents.md**.

---

## Output Formats

### Skill Created

```
Skill Created

Path: [full path to SKILL.md]
Name: [skill name]
Type: [Technique/Pattern/Reference/Coordination]
Size: [word count] words

Verification:
- [ ] Duplicate check passed
- [ ] Frontmatter valid
- [ ] Workflow is imperative
- [ ] Output format defined
- [ ] [For discipline skills] Pressure test passed

Ready for deployment: [Yes/No]
```

### Skill Audited

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

## Anti-Patterns (Never Do)

| Anti-Pattern | Why It Fails | Instead |
|--------------|--------------|---------|
| Create without duplicate check | Fragments knowledge | Search first |
| 400+ line skill | Agents skim, miss sections | Extract to reference.md |
| Audit without reading file | Miss actual issues | Always read the skill |
| Fix without testing | May not address root cause | Verify fix works |
| Declarative-only workflow | Agents don't execute | Use imperative steps |
| Description with workflow | Agents skip body | Triggers only |

---

## Quick Reference

### Skill Creation Checklist
- [ ] Searched for existing similar skills
- [ ] Defined skill type
- [ ] Wrote valid frontmatter (name + description)
- [ ] Description explains what + when to use
- [ ] Body has imperative workflow
- [ ] Output format defined
- [ ] Size under 500 lines
- [ ] Verified before commit

### Skill Audit Checklist
- [ ] Read the actual skill file
- [ ] Ran audit checklist
- [ ] Identified issue category
- [ ] Proposed specific fix
- [ ] Tested fix (not just analyzed)
