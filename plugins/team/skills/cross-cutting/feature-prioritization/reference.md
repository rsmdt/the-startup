# Prioritization Workshop Guide

Advanced frameworks and facilitation techniques for running effective prioritization sessions. Load this when the SKILL.md frameworks are insufficient, when running a multi-stakeholder session, or when selecting between frameworks requires more context.

---

## When to Use Each Framework

The SKILL.md covers RICE, Value vs Effort, Kano, MoSCoW, Cost of Delay, and Weighted Scoring. This guide adds context on when each fits best and how to run group sessions around them.

| Situation | Best Framework | Why |
|-----------|---------------|-----|
| Large backlog (10+ features), quantitative team | RICE | Objective scores cut through opinion |
| Stakeholders disagree on value | Weighted Scoring | Makes implicit criteria explicit |
| New product area, unclear user expectations | Kano | Reveals what users expect vs. what delights |
| Release scope negotiation | MoSCoW | Produces a commitment, not just a ranking |
| Regulatory or time-sensitive deadline | Cost of Delay / CD3 | Quantifies urgency in dollar terms |
| Quick triage, team is time-constrained | Value vs Effort | Fast visual sort, no calculation needed |
| Strategic alignment review | WSJF (see below) | Connects features to business risk and economics |

---

## WSJF: Weighted Shortest Job First

WSJF is the SAFe adaptation of Cost of Delay for agile teams. It balances economic value, urgency, and job duration.

### Formula

```
WSJF = Cost of Delay / Job Duration

Cost of Delay = User-Business Value + Time Criticality + Risk Reduction / Opportunity Enablement
```

### Scoring (1–21 Fibonacci scale)

Using Fibonacci numbers forces relative sizing and avoids false precision.

| Factor | 1 | 2 | 3 | 5 | 8 | 13 | 21 |
|--------|---|---|---|---|---|----|----|
| User-Business Value | Negligible | Minor | Small | Moderate | Significant | Large | Transformative |
| Time Criticality | None | Low | Slight | Moderate | Urgent | Very urgent | Extreme |
| Risk Reduction / Opportunity | None | Minor | Small | Moderate | Significant | Large | Critical |
| Job Duration | Very small | Small | Medium-small | Medium | Medium-large | Large | Very large |

### Example

```
Feature: API rate limiting
  User-Business Value: 8 (prevents outages for paying customers)
  Time Criticality: 13 (competitor incidents increasing)
  Risk Reduction: 13 (prevents potential SLA breach)

  Cost of Delay = 8 + 13 + 13 = 34

  Job Duration: 3 (2-week sprint)

  WSJF = 34 / 3 = 11.3

Feature: Redesigned onboarding flow
  User-Business Value: 13
  Time Criticality: 3 (no deadline pressure)
  Risk Reduction: 2 (minor churn risk)

  Cost of Delay = 13 + 3 + 2 = 18

  Job Duration: 8 (6-week effort)

  WSJF = 18 / 8 = 2.25

Decision: API rate limiting first (11.3 vs 2.25)
```

### When to Use WSJF

- SAFe or PI planning environments
- When business risk and urgency are key factors
- Features with compliance or SLA implications
- When engineering effort is the constraint, not just value

---

## Kano in Practice: Running a Survey

The Kano model requires user data to classify features correctly. This section covers how to gather and interpret that data.

### Survey Design

For each feature candidate, ask both questions:

```
Functional question: "If [feature] were available, how would you feel?"
Dysfunctional question: "If [feature] were NOT available, how would you feel?"

Answer options (same for both):
  1. I like it
  2. I expect it
  3. I'm neutral
  4. I can tolerate it
  5. I dislike it
```

Target 20–50 respondents per customer segment. Results below 20 respondents are directional only.

### Interpreting Results

Tally each functional/dysfunctional response pair against the interpretation matrix in SKILL.md. Calculate the percentage of respondents in each category (M, O, A, I, R, Q).

```
Rule of thumb:
  If Must-Have > 50%    → Non-negotiable, ship before anything else
  If Attractive > 40%   → Strong differentiator, high ROI
  If Indifferent > 50%  → Deprioritize, spend effort elsewhere
  If Reverse > 20%      → Do not build — it actively harms satisfaction
```

### Segment Separately

New users and power users often rate the same feature differently. A feature that is Indifferent to new users may be a Must-Have for power users. Segment your analysis when retention strategy depends on a specific user cohort.

---

## MoSCoW in Practice: Running the Session

MoSCoW is simple to explain but breaks down without clear rules. These facilitation steps prevent scope creep.

### Pre-Session Preparation

1. Define the constraint upfront: "We have 6 weeks of engineering capacity."
2. Identify who has authority to commit. Observers attend; decision-makers vote.
3. Prepare a feature list. Do not allow new items during the session.

### Session Steps

```
Step 1: Classify independently (10 minutes)
  Each participant assigns M/S/C/W to every feature.
  No discussion yet.

Step 2: Reveal and identify splits (10 minutes)
  Show a tally of votes per feature.
  Flag features with significant disagreement (>30% split).

Step 3: Discuss splits only (20–30 minutes)
  Do not revisit consensus items — they are done.
  For each split, ask: "What would have to be true for this to be a Must?"

Step 4: Apply the budget rule (10 minutes)
  Total Must items. If they exceed 60% of capacity, someone must move an item.
  The product owner has final say when consensus fails.

Step 5: Confirm Won't list explicitly (5 minutes)
  Read the Won't list aloud. Agreement here prevents later scope creep.
```

### Common Failure Modes

| Problem | Symptom | Fix |
|---------|---------|-----|
| Must inflation | >80% of items classified as Must | Enforce the 60% capacity rule — force trade-offs |
| Scope creep | New items added after session | Lock the list before the session starts |
| Silent dissent | Everyone agrees in the room, feature gets cut later | Explicitly read Won't list, get verbal confirmation |
| No Won't items | Team avoids hard conversations | Ask: "What are we explicitly NOT doing this cycle?" |

---

## Value vs Effort: Running a Dot Vote Session

The Value vs Effort matrix works best as a participatory exercise, not a solo analysis.

### Materials

- Whiteboard or digital canvas (Miro, Figma)
- Sticky notes or virtual cards for each feature (one per card)
- Dot stickers or voting tokens (3–5 per participant)

### Session Format

```
Step 1: Draw the matrix (2 minutes)
  Label axes: Value (Low → High) and Effort (Low → High).
  Mark quadrant names: Quick Wins, Strategic, Fill-Ins, Time Sinks.

Step 2: Silent placement (10 minutes)
  Each participant places features on the board without discussion.
  Everyone works simultaneously. No lobbying.

Step 3: Discuss outliers (15 minutes)
  Find features with the widest placement spread.
  Ask: "Why did you place this here?" Identify assumption gaps.
  Move to consensus position.

Step 4: Dot vote on Quick Wins (5 minutes)
  Give each participant 3 votes. They place dots on Quick Win quadrant items.
  Top-voted items become first priorities.
```

### Effort Calibration Tip

Effort disagreements usually reflect different assumptions about scope, not different engineering estimates. Before placing items, agree on a reference point: "Feature X took 2 weeks — that is Medium effort." Calibrate everything relative to that anchor.

---

## Facilitation Tips for Any Session

### Pre-Session

- Send the feature list 24 hours before so participants can think independently.
- Define the decision boundary: what is in scope for this session, what is locked.
- Time-box ruthlessly. State the duration at the start and stick to it.

### During the Session

- Separate generating from evaluating. Do not debate while placing items.
- Name the HiPPO dynamic explicitly if it emerges: "Let's hear from others before we discuss your view."
- Use silence deliberately. After asking a question, wait 10 seconds before speaking.
- Park tangents in a visible "parking lot" — acknowledge them without letting them derail.

### Post-Session

- Publish the ranked output within 24 hours. Delay breeds doubt.
- Document the rationale for top decisions, not just the ranking.
- Set a review date for deferred items so they are not forgotten.

---

## Combining Frameworks

No single framework covers every dimension. Combining two frameworks often produces better decisions than using one exclusively.

### Recommended Combinations

| Goal | Primary Framework | Validation Framework |
|------|-------------------|----------------------|
| Quarterly roadmap | RICE (objective ranking) | Value vs Effort (gut-check on outliers) |
| Release scope | MoSCoW (commitment) | Cost of Delay (urgency validation) |
| New market entry | Kano (user expectations) | Weighted Scoring (strategic fit) |
| Resource allocation | WSJF (economic priority) | MoSCoW (feasibility check) |

### Validation Rule

If two frameworks produce the same top 3, confidence is high. If they diverge, investigate why before committing. The divergence usually points to a hidden assumption worth surfacing.

---

## Decision Quality Checklist

Before finalizing a prioritization output, verify:

- [ ] Multiple frameworks used or at least one framework applied rigorously
- [ ] Confidence levels are explicit — not everything rated 80%+ by default
- [ ] Effort estimates reviewed by at least one engineer
- [ ] Won't / defer items documented with a review date
- [ ] Decision rationale captured in a Priority Decision Record (see SKILL.md)
- [ ] Stakeholders who disagreed are documented, not silenced
- [ ] Outcome metrics defined — how will you know if this decision was right?
