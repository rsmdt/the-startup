---
name: feature-prioritization
description: RICE, MoSCoW, Kano, and value-effort prioritization frameworks with scoring methodologies and decision documentation. Use when prioritizing features, evaluating competing initiatives, creating roadmaps, or making build vs defer decisions.
---

# Prioritization Framework Application

Systematic frameworks for making objective prioritization decisions that balance value, effort, and strategic alignment.

## When to Activate

- Prioritizing feature backlogs
- Evaluating competing initiatives
- Making build vs defer decisions
- Creating product roadmaps
- Allocating limited resources
- Justifying prioritization decisions to stakeholders

## RICE Framework

Quantitative scoring for comparing initiatives objectively.

### Formula

```
RICE Score = (Reach × Impact × Confidence) / Effort
```

### Components

| Factor | Description | Scale |
|--------|-------------|-------|
| **Reach** | How many users affected per quarter | Actual number (100, 1000, 10000) |
| **Impact** | Effect on each user | 0.25 (Minimal) to 3 (Massive) |
| **Confidence** | How sure are we | 50% (Low) to 100% (High) |
| **Effort** | Person-months required | Actual estimate (0.5, 1, 3, 6) |

### Impact Scale

| Score | Label | Description |
|-------|-------|-------------|
| 3 | Massive | Life-changing for users, core workflow transformation |
| 2 | High | Major improvement, significant time savings |
| 1 | Medium | Noticeable improvement, minor friction reduction |
| 0.5 | Low | Slight improvement, nice-to-have |
| 0.25 | Minimal | Barely noticeable difference |

### Confidence Scale

| Score | Label | Basis |
|-------|-------|-------|
| 100% | High | User research + validated data + successful tests |
| 80% | Medium | Some data + team experience + analogous examples |
| 50% | Low | Intuition only, no supporting data |

### Example Calculation

```
Feature: One-click reorder

Reach: 5,000 (customers who reorder monthly)
Impact: 2 (High - saves significant time)
Confidence: 80% (Based on support ticket analysis)
Effort: 1 person-month

RICE = (5000 × 2 × 0.8) / 1 = 8000

Feature: Dark mode

Reach: 20,000 (all active users)
Impact: 0.5 (Low - preference, not productivity)
Confidence: 50% (No data, user requests only)
Effort: 2 person-months

RICE = (20000 × 0.5 × 0.5) / 2 = 2500

Decision: One-click reorder scores higher, prioritize first
```

### RICE Template

| Feature | Reach | Impact | Confidence | Effort | Score | Rank |
|---------|-------|--------|------------|--------|-------|------|
| Feature A | 5000 | 2 | 80% | 1 | 8000 | 1 |
| Feature B | 20000 | 0.5 | 50% | 2 | 2500 | 2 |

## Value vs Effort Matrix

Visual framework for quick categorization.

### The Matrix

```
              High Value
                   │
    ┌──────────────┼──────────────┐
    │              │              │
    │  QUICK WINS  │  STRATEGIC   │
    │  Do First    │  Plan & Do   │
    │              │              │
    ├──────────────┼──────────────┤ High
Low │              │              │ Effort
Effort             │              │
    │  FILL-INS    │  TIME SINKS  │
    │  If Spare    │  Avoid       │
    │  Capacity    │              │
    │              │              │
    └──────────────┼──────────────┘
                   │
              Low Value
```

### Quadrant Actions

| Quadrant | Characteristics | Action |
|----------|-----------------|--------|
| **Quick Wins** | High value, low effort | Do immediately |
| **Strategic** | High value, high effort | Plan carefully, staff appropriately |
| **Fill-Ins** | Low value, low effort | Do when nothing else is ready |
| **Time Sinks** | Low value, high effort | Don't do (or simplify drastically) |

### Estimation Guidance

**Value Assessment:**
- Revenue impact
- Cost reduction
- User satisfaction improvement
- Strategic alignment
- Risk reduction

**Effort Assessment:**
- Development time
- Design complexity
- Testing requirements
- Deployment complexity
- Ongoing maintenance

## Kano Model

Categorize features by their impact on satisfaction.

### Categories

```
Satisfaction
     ▲
     │         ╱ Delighters
     │       ╱   (Unexpected features)
     │     ╱
─────┼────●──────────────────────────► Feature
     │    │╲                           Implementation
     │    │  ╲ Performance
     │    │    (More is better)
     │    │
     │    └── Must-Haves
     │         (Expected, dissatisfaction if missing)
     ▼
```

### Category Definitions

| Category | Present | Absent | Example |
|----------|---------|--------|---------|
| **Must-Have** | Neutral | Very dissatisfied | Login functionality |
| **Performance** | More = better | Less = worse | Page load speed |
| **Delighter** | Very satisfied | Neutral | Personalized recommendations |
| **Indifferent** | No effect | No effect | Backend tech choice |
| **Reverse** | Dissatisfied | Satisfied | Forced tutorials |

### Kano Survey Questions

For each feature, ask two questions:

```
Functional: "If [feature] were present, how would you feel?"
Dysfunctional: "If [feature] were absent, how would you feel?"

Answer Options:
1. I like it
2. I expect it
3. I'm neutral
4. I can tolerate it
5. I dislike it
```

### Interpretation Matrix

|  | Like | Expect | Neutral | Tolerate | Dislike |
|--|------|--------|---------|----------|---------|
| **Like** | Q | A | A | A | O |
| **Expect** | R | I | I | I | M |
| **Neutral** | R | I | I | I | M |
| **Tolerate** | R | I | I | I | M |
| **Dislike** | R | R | R | R | Q |

Key: M=Must-Have, O=One-dimensional, A=Attractive, I=Indifferent, R=Reverse, Q=Questionable

## MoSCoW Method

Simple categorization for scope definition.

### Categories

| Category | Definition | Negotiability |
|----------|------------|---------------|
| **Must** | Critical for success, release blocked without | Non-negotiable |
| **Should** | Important but not critical | Can defer to next release |
| **Could** | Nice to have, minor impact | First to cut if needed |
| **Won't** | Explicitly excluded from scope | Not this release |

### Application Rules

```
Budget Allocation (Recommended):
- Must: 60% of capacity
- Should: 20% of capacity
- Could: 20% of capacity (buffer)
- Won't: 0% (explicitly excluded)

Why the buffer matters:
- Must items often take longer than estimated
- Should items may become Must if requirements change
- Could items fill capacity at sprint end
```

### Example

```
Feature: User Registration

MUST:
✓ Email/password signup
✓ Email verification
✓ Password requirements enforcement

SHOULD:
○ Social login (Google)
○ Remember me functionality
○ Password strength indicator

COULD:
◐ Social login (Facebook, Apple)
◐ Profile picture upload
◐ Username suggestions

WON'T (this release):
✗ Two-factor authentication
✗ SSO integration
✗ Biometric login
```

## Cost of Delay

Prioritize by economic impact of waiting.

### CD3 Formula

```
CD3 = Cost of Delay / Duration

Cost of Delay: Weekly value lost by not having the feature
Duration: Weeks to implement
```

### Delay Cost Types

| Type | Description | Calculation |
|------|-------------|-------------|
| **Revenue** | Sales not captured | Lost deals × average value |
| **Cost** | Ongoing expenses | Weekly operational cost |
| **Risk** | Penalty or loss potential | Probability × impact |
| **Opportunity** | Market window | Revenue × time sensitivity |

### Urgency Profiles

```
                Value
                  │
Standard:         │────────────────
                  │
                  └──────────────────► Time

Urgent:           │╲
                  │  ╲
                  │    ╲──────────
                  │
                  └──────────────────► Time

Deadline:         │
                  │────────┐
                  │        │
                  │        └─ (drops to zero)
                  └──────────────────► Time
```

### Example

```
Feature A: New payment method
- Cost of Delay: $10,000/week (lost sales to competitor)
- Duration: 4 weeks
- CD3 = 10000 / 4 = 2500

Feature B: Admin dashboard redesign
- Cost of Delay: $2,000/week (support inefficiency)
- Duration: 2 weeks
- CD3 = 2000 / 2 = 1000

Feature C: Compliance update (deadline in 6 weeks)
- Cost of Delay: $50,000/week after deadline (fines)
- Duration: 4 weeks
- CD3 = 50000 / 4 = 12500 (if started now, 0 if after deadline)

Priority: C (deadline), then A (highest CD3), then B
```

## Weighted Scoring

Custom scoring for organization-specific criteria.

### Building a Weighted Model

```
Step 1: Define Criteria
- Strategic alignment
- Revenue potential
- User demand
- Technical feasibility
- Competitive advantage

Step 2: Assign Weights (total = 100%)
| Criterion | Weight |
|-----------|--------|
| Strategic | 30% |
| Revenue | 25% |
| User demand | 20% |
| Feasibility | 15% |
| Competitive | 10% |

Step 3: Score Each Feature (1-5 scale)
| Feature | Strategic | Revenue | Demand | Feasible | Competitive | Total |
|---------|-----------|---------|--------|----------|-------------|-------|
| A | 5 | 4 | 3 | 4 | 2 | 3.95 |
| B | 3 | 5 | 5 | 3 | 3 | 3.90 |
| C | 4 | 3 | 4 | 5 | 4 | 3.85 |
```

### Calculation

```
Score = Σ (criterion_score × criterion_weight)

Feature A:
= (5 × 0.30) + (4 × 0.25) + (3 × 0.20) + (4 × 0.15) + (2 × 0.10)
= 1.5 + 1.0 + 0.6 + 0.6 + 0.2
= 3.9
```

## Decision Documentation

### Priority Decision Record

```markdown
# Priority Decision: [Feature/Initiative]

## Date: [YYYY-MM-DD]
## Decision: [Prioritize / Defer / Reject]

## Context
[What prompted this decision?]

## Evaluation

### Framework Used: [RICE / Kano / MoSCoW / Weighted]

### Scores
[Show calculations or categorization]

### Trade-offs Considered
- Option A: [description] - [pros/cons]
- Option B: [description] - [pros/cons]

## Decision Rationale
[Why this priority over alternatives?]

## Stakeholders
- Agreed: [names]
- Disagreed: [names, reasons documented]

## Review Date
[When to revisit if deferred]
```

## Framework Selection Guide

| Situation | Recommended Framework |
|-----------|----------------------|
| Comparing many similar features | RICE (quantitative) |
| Quick triage of backlog | Value vs Effort |
| Understanding user expectations | Kano Model |
| Defining release scope | MoSCoW |
| Time-sensitive decisions | Cost of Delay |
| Organization-specific criteria | Weighted Scoring |

## Anti-Patterns

| Anti-Pattern | Problem | Solution |
|--------------|---------|----------|
| **HiPPO** | Highest-paid person's opinion wins | Use data-driven frameworks |
| **Recency Bias** | Last request gets priority | Systematic evaluation of all options |
| **Squeaky Wheel** | Loudest stakeholder wins | Weight by strategic value |
| **Analysis Paralysis** | Over-analyzing decisions | Time-box evaluation |
| **Sunken Cost** | Continuing failed initiatives | Evaluate future value only |
| **Feature Factory** | Shipping without measuring | Tie features to outcomes |

## Best Practices

1. **Use multiple frameworks** - Validate with different approaches
2. **Document decisions** - Enable future learning
3. **Revisit regularly** - Priorities change as context evolves
4. **Include stakeholders** - Ensure buy-in
5. **Measure outcomes** - Validate prioritization quality

## References

- [RICE Scoring Template](examples/rice-template.md) - Spreadsheet template
- [Prioritization Workshop Guide](reference.md) - Facilitation guide
