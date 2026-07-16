---
name: {{slug}}
description: {{terse one-line purpose of this analysis}}
status: draft   # draft | approved | in-progress | review | done
goals:
  - {{achievable, checkable objective 1}}
  - {{achievable, checkable objective 2}}
---

# Analysis: {{Title}}

{{One paragraph framing the analysis — what question it answers and why it matters.}}

## Goals (the review rubric)

The analysis is review-worthy only when every goal below is met by cited evidence:

1. {{objective 1}}
2. {{objective 2}}

## Topics to research

Each topic is an independent angle, researched in parallel as its own stream:

- {{topic 1 — an independent angle}}
- {{topic 2 — an independent angle}}
- {{topic 3 — an independent angle}}

## Tracker

Granular research tasks (and their status) live in the tracker:

- Location: `.start/analysis/{{NNN}}-{{slug}}/.tracker/`
- Query the graph: `node <path-to>/task.js graph --home .start/analysis/{{NNN}}-{{slug}}`
- Overview (roots): `node <path-to>/task.js list --home .start/analysis/{{NNN}}-{{slug}}`

The synthesized report is written to `report.md` in this directory once all goals are met.
