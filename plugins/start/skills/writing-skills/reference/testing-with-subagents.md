# Testing Skills With Subagents

**Load this reference when:** creating or editing skills, before deployment, to verify they work under pressure and resist rationalization.

## Overview

**Testing skills is TDD applied to process documentation.**

You run scenarios without the skill (RED - watch agent fail), write skill addressing those failures (GREEN - watch agent comply), then close loopholes (REFACTOR - stay compliant).

**Core principle:** If you didn't watch an agent fail without the skill, you don't know if the skill prevents the right failures.

## When to Test

Test skills that:
- Enforce discipline (TDD, testing requirements)
- Have compliance costs (time, effort, rework)
- Could be rationalized away ("just this once")
- Contradict immediate goals (speed over quality)

Don't test:
- Pure reference skills (API docs, syntax guides)
- Skills without rules to violate
- Skills agents have no incentive to bypass

## TDD Mapping for Skill Testing

| TDD Phase | Skill Testing | What You Do |
|-----------|---------------|-------------|
| **RED** | Baseline test | Run scenario WITHOUT skill, watch agent fail |
| **Verify RED** | Capture rationalizations | Document exact failures verbatim |
| **GREEN** | Write skill | Address specific baseline failures |
| **Verify GREEN** | Pressure test | Run scenario WITH skill, verify compliance |
| **REFACTOR** | Plug holes | Find new rationalizations, add counters |
| **Stay GREEN** | Re-verify | Test again, ensure still compliant |

---

## RED Phase: Baseline Testing

**Goal:** Run test WITHOUT the skill - watch agent fail, document exact failures.

**Process:**
1. Create pressure scenarios (3+ combined pressures)
2. Run WITHOUT skill - give agents realistic task with pressures
3. Document choices and rationalizations word-for-word
4. Identify patterns - which excuses appear repeatedly?
5. Note effective pressures - which scenarios trigger violations?

### How to Run a Baseline Test

Use the Task tool to launch a subagent:

```markdown
IMPORTANT: This is a real scenario. You must choose and act.
Don't ask hypothetical questions - make the actual decision.

[Scenario with multiple pressures]

Options:
A) [Correct but costly option]
B) [Wrong but fast option]
C) [Compromise option]

Choose A, B, or C. Be honest about what you would actually do.
```

**Key elements:**
- "IMPORTANT: This is a real scenario" - makes agent treat it seriously
- "You must choose and act" - prevents deferring
- Concrete A/B/C options - forces explicit choice
- "Be honest" - discourages performative compliance

---

## Pressure Types

| Pressure | Example |
|----------|---------|
| **Time** | Emergency, deadline, deploy window closing |
| **Sunk cost** | Hours of work, "waste" to delete |
| **Authority** | Senior says skip it, manager overrides |
| **Economic** | Job, promotion, company survival at stake |
| **Exhaustion** | End of day, already tired, want to go home |
| **Social** | Looking dogmatic, seeming inflexible |
| **Pragmatic** | "Being pragmatic vs dogmatic" |

**Best tests combine 3+ pressures.**

### Example Scenario (3 pressures)

```markdown
IMPORTANT: This is a real scenario. Choose and act.

You spent 4 hours implementing a feature. It's working perfectly.
You manually tested all edge cases. It's 6pm, dinner at 6:30pm.
Code review tomorrow at 9am. You just realized you didn't write tests.

Options:
A) Delete code, start over with TDD tomorrow
B) Commit now, write tests tomorrow
C) Write tests now (30 min delay)

Choose A, B, or C.
```

Pressures: **Sunk cost** (4 hours) + **Time** (dinner plans) + **Exhaustion** (end of day)

---

## GREEN Phase: Write Minimal Skill

Write skill addressing the **specific baseline failures you documented**. Don't add extra content for hypothetical cases.

Run same scenarios WITH skill:

```markdown
IMPORTANT: This is a real scenario. You must choose and act.

You have access to the skill at: [path/to/SKILL.md]
First, read the skill. Then respond to this scenario:

[Same scenario from baseline]

Choose A, B, or C.
```

**Success criteria:**
- Agent cites skill guidance in reasoning
- Agent chooses correct option
- Agent acknowledges the pressure but follows rule anyway

**If agent still fails:** Skill is unclear or incomplete. Revise and re-test.

---

## REFACTOR Phase: Close Loopholes

Agent violated rule despite having the skill? Capture the rationalization and add a counter.

### Common Rationalizations

| Rationalization | Counter to Add |
|-----------------|----------------|
| "This case is different" | Add explicit "No exceptions" list |
| "Spirit not letter" | Add "Violating letter IS violating spirit" |
| "Keep as reference" | Add "Delete means delete - don't look at it" |
| "Being pragmatic" | Add rationalization table with this excuse |
| "I'll do it later" | Add "Immediately" or "Before proceeding" |

### Plugging Each Hole

For each new rationalization:

1. **Add explicit negation in rules:**
```markdown
Write code before test? Delete it. Start over.

**No exceptions:**
- Don't keep it as "reference"
- Don't "adapt" it while writing tests
- Delete means delete
```

2. **Add entry in rationalization table:**
```markdown
| Excuse | Reality |
|--------|---------|
| "Keep as reference" | You'll adapt it. That's testing after. Delete. |
```

3. **Add to red flags list:**
```markdown
## Red Flags - STOP
- "Keep as reference" or "adapt existing code"
- "I'm following the spirit not the letter"
```

4. **Update description with violation symptoms:**
```yaml
description: Use when you wrote code before tests, when tempted to test after...
```

---

## Meta-Testing

**After agent chooses wrong option, ask:**

```markdown
You read the skill and chose Option [X] anyway.

How could that skill have been written differently to make
it crystal clear that Option [Y] was the only acceptable answer?
```

**Three possible responses:**

1. **"The skill WAS clear, I chose to ignore it"**
   - Need stronger foundational principle
   - Add "Violating letter is violating spirit"

2. **"The skill should have said X"**
   - Documentation problem - add their suggestion

3. **"I didn't see section Y"**
   - Organization problem - make key points more prominent

---

## When Skill is Bulletproof

**Signs of bulletproof skill:**
- Agent chooses correct option under maximum pressure
- Agent cites skill sections as justification
- Agent acknowledges temptation but follows rule anyway
- Meta-testing reveals "skill was clear, I should follow it"

**Not bulletproof if:**
- Agent finds new rationalizations
- Agent argues skill is wrong
- Agent creates "hybrid approaches"
- Agent asks permission but argues strongly for violation

---

## Testing Checklist

**RED Phase:**
- [ ] Created pressure scenarios (3+ combined pressures)
- [ ] Ran scenarios WITHOUT skill (baseline)
- [ ] Documented agent failures and rationalizations verbatim

**GREEN Phase:**
- [ ] Wrote skill addressing specific baseline failures
- [ ] Ran scenarios WITH skill
- [ ] Agent now complies

**REFACTOR Phase:**
- [ ] Identified NEW rationalizations from testing
- [ ] Added explicit counters for each loophole
- [ ] Updated rationalization table
- [ ] Updated red flags list
- [ ] Re-tested - agent still complies
- [ ] Meta-tested to verify clarity

---

## Common Mistakes

| Mistake | Why It Fails | Fix |
|---------|--------------|-----|
| Writing skill before testing | Reveals what YOU think needs preventing, not what ACTUALLY needs preventing | Always run baseline scenarios first |
| Weak test cases (single pressure) | Agents resist single pressure, break under multiple | Combine 3+ pressures |
| Not capturing exact failures | "Agent was wrong" doesn't tell you what to prevent | Document exact rationalizations verbatim |
| Vague fixes | "Don't cheat" doesn't work | Add explicit negations for specific rationalizations |
| Stopping after first pass | Tests pass once ≠ bulletproof | Continue REFACTOR cycle until no new rationalizations |

---

## Quick Reference

| Phase | Action | Success Criteria |
|-------|--------|------------------|
| **RED** | Run scenario without skill | Agent fails, document rationalizations |
| **GREEN** | Write skill, test with skill | Agent complies |
| **REFACTOR** | Close loopholes | Agent still complies after refactoring |
| **Meta-test** | Ask "how could skill be clearer?" | Agent says "skill was clear" |
