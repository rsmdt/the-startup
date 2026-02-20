# RICE Scoring Template

A ready-to-use template for scoring and ranking features using the RICE framework. Copy the blank template, fill in your estimates, and let the scores determine priority.

## Quick Reminder: The Formula

```
RICE Score = (Reach × Impact × Confidence) / Effort
```

Higher score = higher priority.

---

## Blank Template

Copy this table for your prioritization session.

| Feature | Reach (users/qtr) | Impact (0.25–3) | Confidence (50–100%) | Effort (person-months) | RICE Score | Rank |
|---------|-------------------|-----------------|----------------------|------------------------|------------|------|
|         |                   |                 |                      |                        |            |      |
|         |                   |                 |                      |                        |            |      |
|         |                   |                 |                      |                        |            |      |
|         |                   |                 |                      |                        |            |      |
|         |                   |                 |                      |                        |            |      |

### Score Calculation

For each row:

```
RICE Score = (Reach × Impact × (Confidence / 100)) / Effort
```

Note: Convert confidence percentage to decimal before calculating (80% → 0.80).

---

## Scale Reference

### Impact Scale

| Value | Label | Description |
|-------|-------|-------------|
| 3 | Massive | Core workflow transformation, life-changing for users |
| 2 | High | Major improvement, significant time or cost savings |
| 1 | Medium | Noticeable improvement, reduces meaningful friction |
| 0.5 | Low | Slight improvement, nice-to-have quality of life |
| 0.25 | Minimal | Barely noticeable difference |

### Confidence Scale

| Value | Label | When to Use |
|-------|-------|-------------|
| 100% | High | User research + validated data + prior successful tests |
| 80% | Medium | Some data + team experience + analogous examples |
| 50% | Low | Intuition or anecdote only, no supporting data |

**Rule of thumb**: If you're debating between two confidence levels, use the lower one. Overconfidence inflates scores.

---

## Filled-In Example: SaaS Analytics Dashboard

This example scores five competing features for a B2B analytics product.

### Context

- Team capacity: 3 person-months per quarter
- User base: 25,000 monthly active users
- 8,000 users engage with the reporting section

### Scored Features

| Feature | Reach | Impact | Confidence | Effort | RICE Score | Rank |
|---------|-------|--------|------------|--------|------------|------|
| CSV export | 6,000 | 2 | 80% | 0.5 | **19,200** | 1 |
| Scheduled email reports | 3,500 | 2 | 80% | 1 | **5,600** | 2 |
| Custom date range picker | 8,000 | 1 | 100% | 0.5 | **16,000** | — |
| Dashboard sharing (public link) | 2,000 | 2 | 50% | 1.5 | **1,333** | 4 |
| Dark mode | 25,000 | 0.25 | 50% | 2 | **1,563** | 3 |

### Score Calculations

```
CSV export:
  (6,000 × 2 × 0.80) / 0.5 = 9,600 / 0.5 = 19,200

Scheduled email reports:
  (3,500 × 2 × 0.80) / 1 = 5,600

Custom date range picker:
  (8,000 × 1 × 1.00) / 0.5 = 16,000

Dashboard sharing:
  (2,000 × 2 × 0.50) / 1.5 = 2,000 / 1.5 = 1,333

Dark mode:
  (25,000 × 0.25 × 0.50) / 2 = 3,125 / 2 = 1,563
```

### Adjusted Priority

Scores alone tell most of the story, but two features need a note:

**Custom date range picker (16,000)** scored second-highest but was pre-committed to a partner. It does not compete for the open roadmap slots.

Final open-roadmap ranking with 3 person-months of capacity:

1. **CSV export** (score: 19,200) — 0.5 months. High confidence data from 62 support tickets.
2. **Scheduled email reports** (score: 5,600) — 1 month. Validated by customer interviews with 3 enterprise accounts.
3. **Dark mode** (score: 1,563) — defer to next quarter. High reach, but very low impact and confidence.
4. **Dashboard sharing** (score: 1,333) — defer. Low confidence, significant security design work needed first.

Total committed: 1.5 months, leaving 1.5 months buffer for Must items and scope growth.

---

## Common Mistakes

| Mistake | Problem | Fix |
|---------|---------|-----|
| Using 100% confidence by default | Inflates every score equally, ranking becomes meaningless | Only use 100% when you have validated data |
| Estimating Reach as total user base | Overstates impact — most features affect a subset of users | Count only users who encounter the relevant workflow |
| Ignoring Effort entirely | Low-effort features win by default regardless of value | Always estimate Effort; a 0.1 score skews results badly |
| Scoring in isolation | Individual scorers have different mental scales | Score as a group, or calibrate with one known-reference feature |
| Never revisiting scores | Context changes — last quarter's data is stale | Re-score when key inputs (user count, team size) shift significantly |

---

## Tips for Estimation Sessions

**Anchor on a reference feature.** Before scoring new items, pick one feature the team already shipped and agree on its scores. Use it as a calibration baseline for Reach and Impact.

**Score independently, then converge.** Have each team member fill in scores before comparing. This surfaces disagreements that group discussion would suppress.

**Time-box the session.** Spend no more than 5 minutes per feature. If a feature requires more debate, mark it as low confidence and move on.

**Document your assumptions.** Record the data source behind each Reach estimate and the reasoning behind each Confidence rating. Scores without sources are guesses with extra steps.
