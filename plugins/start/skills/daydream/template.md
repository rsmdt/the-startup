---
name: {{slug}}
description: {{terse one-line purpose of this daydream run}}
status: draft   # draft | approved | in-progress | review | done
dateRange:
  since: {{YYYY-MM-DD}}
  until: {{YYYY-MM-DD}}
  label: {{yesterday | last 3 days | ...}}
applyTier: advisory   # v1 ships advisory only (bounded / eval tiers are roadmap)
angles:
  - {{signal angle 1 — e.g. tool-errors}}
  - {{signal angle 2 — e.g. user-corrections}}
---

# Daydream: {{Title}}

{{One paragraph framing the run — which range of sessions it reviews and what it is trying to improve.}}

## Signal angles (parallel mine streams)

Each angle is mined independently and in parallel over the date range:

- {{angle 1}}
- {{angle 2}}
- {{angle 3}}

## Edit-target scope

Proposals may target git-tracked artifacts under `plugins/` only:

- Skills — `plugins/*/skills/*/SKILL.md` (+ `reference/`, `templates/`)
- Agents — `plugins/team/agents/the-*/*.md`
- Workflow golden copies — `plugins/*/skills/*/reference/*.js` (bump the version marker)
- Templates — `plugins/*/skills/*/template*.md`
- Output styles — `plugins/*/output-styles/*.md`

Never the untracked installed `.claude/workflows/*.js`. **Advisory**: proposals are reported; a human
applies them on a branch/PR.

## Tracker

Granular mine and proposal tasks (and their status) live in the tracker:

- Location: `.start/daydream/{{slug}}/.tracker/`
- Query the graph: `node <path-to>/task.js graph --home .start/daydream/{{slug}}`
- Overview (roots): `node <path-to>/task.js list --home .start/daydream/{{slug}}`

The proposals report is written to `report.md` in this directory once the run completes.
