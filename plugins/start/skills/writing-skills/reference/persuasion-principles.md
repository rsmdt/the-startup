# Persuasion Principles for Skill Design

**Load this reference when:** designing discipline-enforcing skills, understanding why certain language patterns work better, or making skills more effective under pressure.

## Overview

LLMs respond to the same persuasion principles as humans. Understanding this psychology helps you design more effective skills - not to manipulate, but to ensure critical practices are followed even under pressure.

**Research foundation:** Meincke et al. (2025) tested 7 persuasion principles with N=28,000 AI conversations. Persuasion techniques more than doubled compliance rates (33% → 72%, p < .001).

---

## Principles That Work for Skills

### 1. Authority

**What it is:** Deference to expertise and non-negotiable rules.

**How to use:**
- Imperative language: "YOU MUST", "Never", "Always"
- Non-negotiable framing: "No exceptions"
- Eliminates decision fatigue and rationalization

**Example:**
```markdown
✅ Write code before test? Delete it. Start over. No exceptions.
❌ Consider writing tests first when feasible.
```

### 2. Commitment

**What it is:** Consistency with prior statements and explicit choices.

**How to use:**
- Force explicit choices: "Choose A, B, or C"
- Require announcements: "First, state what you will do"
- Use tracking: Checklists that must be completed

**Example:**
```markdown
✅ Before proceeding, confirm: "I have checked for duplicates: [Yes/No]"
❌ You might want to check for duplicates.
```

### 3. Scarcity / Urgency

**What it is:** Time-bound requirements that prevent procrastination.

**How to use:**
- "Before proceeding" requirements
- "Immediately after X" sequences
- Sequential dependencies

**Example:**
```markdown
✅ After completing a task, IMMEDIATELY run verification before proceeding.
❌ You can verify when convenient.
```

### 4. Social Proof

**What it is:** Establishing what's normal and expected.

**How to use:**
- Universal patterns: "Every time", "Always"
- Failure modes: "X without Y = failure"
- "All skills follow this pattern"

**Example:**
```markdown
✅ Skills without testing = skills that fail in production. Every time.
❌ Some people find testing helpful.
```

---

## Principles to Avoid

### Liking
Don't use "be friendly" to get compliance. Creates sycophancy and conflicts with honest feedback.

### Reciprocity
Rarely needed and can feel manipulative.

---

## Principle Combinations by Skill Type

| Skill Type | Use | Avoid |
|------------|-----|-------|
| **Discipline-enforcing** | Authority + Commitment + Social Proof | Liking |
| **Technique/guidance** | Moderate Authority | Heavy authority |
| **Reference** | Clarity only | All persuasion |

---

## Why This Works

**Bright-line rules reduce rationalization:**
- "YOU MUST" removes decision fatigue
- Absolute language eliminates "is this an exception?" questions
- Explicit anti-rationalization counters close specific loopholes

**LLMs are parahuman:**
- Trained on human text containing these patterns
- Authority language precedes compliance in training data
- Commitment sequences (statement → action) frequently modeled

---

## Ethical Use

**Legitimate:**
- Ensuring critical practices are followed
- Creating effective documentation
- Preventing predictable failures

**Test:** Would this technique serve the user's genuine interests if they fully understood it?

---

## Quick Reference

When designing a discipline skill:

1. **Use Authority** - "YOU MUST", "No exceptions", imperative language
2. **Use Commitment** - Force explicit choices, require announcements
3. **Use Social Proof** - "Every time", "All skills do X"
4. **Use Scarcity** - "Before proceeding", "Immediately"
5. **Avoid Liking** - Don't use friendliness for compliance

## Research Citations

- Cialdini, R. B. (2021). *Influence: The Psychology of Persuasion.* Harper Business.
- Meincke, L., et al. (2025). Call Me A Jerk: Persuading AI to Comply. University of Pennsylvania.
