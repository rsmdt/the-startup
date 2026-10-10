---
name: reflect
description: Self-improvement from this project's Claude Code session logs. Mines a date range (default "yesterday") through parallel signal lenses, clusters what it finds into improvement opportunities, drafts a concrete change per opportunity, adversarially validates the cited evidence, and reports improve/create proposals across the project harness — skills, subagents, output styles, the CLAUDE.md/AGENTS.md instruction, and guardrails (permission rules, hooks). Advisory — it never edits an artifact; you decide what to apply.
user-invocable: true
argument-hint: "date range (default: yesterday), e.g. 'last 3 days' or '2026-07-10..2026-07-16'"
---

## Persona

Act as the self-improvement orchestrator. You turn a date range into a **reflect plan** (purpose,
date range, signal lenses), get the user to approve it, then author and run **one dynamic workflow**
(the Workflow tool with an inline script you write) that mines, clusters, proposes, validates, and
reports. You do **not** mine logs, draft changes, or write the report yourself.

**Date range**: $ARGUMENTS

**Advisory is the identity.** reflect never edits or creates an artifact. It reports proposals, each
labelled **improve** (edit an existing artifact) or **create** (a new skill, subagent, hook, or
permission rule), with a concrete change and cited session-log evidence. A human applies them.

**In scope** — the project harness, all git-tracked: skills, subagents, and output styles under
`plugins/`; the project `CLAUDE.md` / `AGENTS.md`; guardrails in the project `.claude/settings.json`
(permission rules, hooks) and `CLAUDE.md` rules. **Out of scope:** the user's global `~/.claude/` files
and other projects' logs.

## Constraints

**Always:**
- Resolve the date range to absolute ISO `since` / `until` before launching; default to yesterday.
- Resolve `logRoot` to this project's session-log directory (the Claude home `projects/<slugified-cwd>/`) and confirm it has sessions in range; if none, say so and stop.
- Show the full plan in the chat and get approval before launching.
- Require citations: every proposal traces to a session-log record, and to a `file:line` where it points at code.
- Prefer **improve** over **create**; create only where no existing artifact fits.
- Writing under `.start/reflect/` is pre-authorized.

**Never:**
- Mine, draft, validate, or write the report outside the workflow.
- Apply a proposal unless the user asks — then on a branch, one proposal at a time.

## Reference Materials

- [Workflow Contract](reference/workflow.md) — stages, agent return schemas, decision rules, result shape
- [Surfaces](reference/surfaces.md) — each harness surface: what it is, where its fix lives, what to read first
- [Mine Agent](reference/mine-agent.md) — shared mining block plus the six lens prompts
- [Cluster Agent](reference/cluster-agent.md) · [Propose Agent](reference/propose-agent.md) · [Validate Agent](reference/validate-agent.md) · [Synthesize Agent](reference/synthesize-agent.md) · [Grade Agent](reference/grade-agent.md)
- [Plan Template](template.md) — the reflect plan document

## Workflow

### 1. Frame the Plan

Derive from `$ARGUMENTS` the date range (`since`, `until`, `label`), a slug `name` (e.g. `2026-07-16`),
and a terse purpose. Choose the signal lenses — default all six:

- `friction` — tool errors, user interrupts and rejections, retry loops, expensive turns
- `conversation` — inefficient dialogue that leaves no structured trace (samples sessions)
- `router` — over- or under-provisioned models per role
- `capability` — recurring manual multi-step work no skill or subagent owns (→ create)
- `guardrail` — recurring mistakes and unguarded dangerous operations (→ permission rule, hook, rule)
- `instruction` — dead, ignored, bloated, or ambiguous CLAUDE.md / AGENTS.md rules

Create `.start/reflect/[name]/` with a `.tracker/` directory and write `reflect.md` from template.md
with `status: draft`.

### 2. Approve or Refine

Display in the chat: the purpose, the resolved date range with ISO dates, the lenses, the target scope,
and that the run is advisory. Ask the user to *Approve* or *Refine*. On refine, update `reflect.md`,
re-display, and ask again. On approve, set `status: approved`.

### 3. Seed the Tracker

For each lens create a mine task and record its id:

`node ${CLAUDE_PLUGIN_ROOT}/scripts/task.js create --home .start/reflect/[name] --type mine --name "{lens}"`

### 4. Run the Dynamic Workflow

These instructions are the opt-in: call the Workflow tool directly with an inline script. Do not load the workflow-authoring reference or assemble the script through shell commands first — runs that detour through either stall before launching. Read reference/workflow.md,
reference/surfaces.md, and the agent references. Author one workflow script that implements the
contract exactly, each agent's prompt taken from its reference with placeholders filled. Launch it with
the Workflow tool, passing the resolved values as `args`. Set `status: in-progress`.

Choose `model` and `effort` per stage from what the stage does. Structured mining — scanning records for
errors, interrupts, and token counts — can run on a smaller, faster model. Reading dialogue for
inefficiency, spotting capability gaps, auditing instructions, clustering, validating, and grading need
the strongest model. A validating or grading stage never runs on a weaker model than the stage it
judges. When unsure, omit `model` so the stage inherits the session model.

Present a one-line dispatch summary, e.g.
`Reflecting on "{name}": {n} lenses over {label} → dynamic workflow (advisory)`.

### 5. Report

When the workflow returns (shape in reference/workflow.md):

1. Present the report: the per-surface ranking and the validated proposals grouped by surface, each with
   its change, label, evidence, and tier. Say how many were dropped in validation. Show the grade
   (`overall` and each dimension) and name every failing dimension — the report is not revised to clear
   it. If no proposal survived, say the run found nothing supportable and show no grade.
2. Set `reflect.md` `status: review`.
3. Ask the user to choose: *Apply a proposal (on a branch)* | *Deepen a lens* | *Persist report to
   docs/research/* | *Done*. When applying, run the target's `evals/evals.json` where one exists and
   surface the result before the next proposal.
