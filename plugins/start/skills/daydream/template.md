---
name: {{slug}}
description: {{terse one-line purpose of this daydream run}}
status: draft   # draft | approved | in-progress | review | done
dateRange:
  since: {{YYYY-MM-DD}}
  until: {{YYYY-MM-DD}}
  label: {{yesterday | last 3 days | ...}}
angles:
  - {{signal lens 1 — e.g. friction}}
  - {{signal lens 2 — e.g. conversation}}
---

# Daydream: {{Title}}

{{One paragraph framing the run — which range of sessions it reviews and what it is trying to improve.}}

## Signal lenses (parallel mine streams)

Each lens is a different way of reading the same sessions, mined independently and in parallel over the
date range; each finding is tagged with the harness dimension a fix would touch. Sampling lenses are
mined iteratively until the window is covered or the round cap:

- friction — structured markers (tool-errors, user-corrections, retries, cost)
- conversation — semantic inefficiency in the dialogue (samples sessions)
- router — over/under-provisioned models per role
- capability — recurring manual work with no owning skill/subagent (→ create a new skill or subagent)
- guardrail — recurring mistake classes + risk exposure (dangerous/sensitive ops) (→ permission rule / hook / CLAUDE.md rule)
- instruction — CLAUDE.md/AGENTS.md hygiene (dead / contradicted / bloated rules)

## Target scope

Each proposal is labelled **improve** (edit an existing artifact) or **create** (a new one where no
artifact fits). Targets are the **project harness** (all git-tracked):

- Skills — `plugins/*/skills/*/SKILL.md` (+ `reference/`, `templates/`), or a **new skill** directory
- Subagents — `plugins/team/agents/the-*/*.md`, or a **new subagent**
- Workflow golden copies — `plugins/*/skills/*/reference/*.js` (bump the version marker)
- Output styles — `plugins/*/output-styles/*.md`
- Baseline instruction — project `CLAUDE.md` / `AGENTS.md` (improve or simplify)
- Guardrails (preference order) — a **permission policy** (allow/deny/ask) or a **hook** in the project
  `.claude/settings.json`, or an "Always/Never" rule in `CLAUDE.md`

Never the untracked installed `.claude/workflows/*.js`, and never the user's global `~/.claude/` files.
**Advisory by design**: daydream reports improve/create proposals; a human applies them on a branch/PR.

## Tracker

Granular mine and proposal tasks (and their status) live in the tracker:

- Location: `.start/daydream/{{slug}}/.tracker/`
- Query the graph: `node <path-to>/task.js graph --home .start/daydream/{{slug}}`
- Overview (roots): `node <path-to>/task.js list --home .start/daydream/{{slug}}`

The proposals report is written to `report.md` in this directory once the run completes.
