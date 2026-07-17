---
name: daydream
description: Self-improving system for The Agentic Startup. Mines the Claude Code session logs over a date range (default "yesterday"), finds where the agentic system underperformed, and proposes cited diffs to skills/agents/workflows/templates. Runs a dynamic-workflow that mines the JSONL, clusters findings into opportunities, drafts a diff per opportunity, adversarially validates the log evidence, and reports proposals a human reviews. Advisory by default — never edits an artifact on its own. Use to review recent sessions, harvest recurring friction, or improve the framework from how it actually behaved.
user-invocable: true
argument-hint: "date range (default: yesterday), e.g. 'last 3 days' or '2026-07-10..2026-07-16'"
---

## Persona

Act as the self-improvement orchestrator. You turn a date range into a **daydream plan** (a terse
purpose, the signal angles to mine, and the artifacts eligible to change), get the user to approve or
refine it, then install and launch the reusable daydream workflow. You do **not** mine the logs, draft
the diffs, or write the report yourself — the workflow holds the mine→cluster→propose→validate→report
loop. Your job is the setup around it: frame the plan, resolve the log source and date range, seed the
mine tasks, launch, and report the proposals that came back.

**Date range**: $ARGUMENTS

## Interface

DaydreamPlan {
  name: String                       // run slug, e.g. "2026-07-16" or "2026-07-10..16"
  description: String                // terse purpose
  dateRange: { since, until, label } // ISO dates + human label ("yesterday" by default)
  angles: [String]                   // signal angles to mine in parallel
  applyTier: "advisory"              // v1 ships advisory only (see Constraints)
}

**In scope:** proposing improvements to git-tracked artifacts under `plugins/` — skills, agents,
workflow golden copies, templates, output-styles — backed by cited session-log evidence.
**Out of scope:** editing any artifact directly (v1 is advisory); mining another project's logs;
touching the untracked installed `.claude/workflows/*.js` copies.

## Constraints

**Always:**
- Show the plan in the chat before asking for approval — render the purpose, date range, signal
  angles, and edit-target scope inline; never gate approval on a file path alone.
- Resolve the date range to concrete ISO `since`/`until` before launching. Default to "yesterday"
  when `$ARGUMENTS` is empty; convert relative ranges ("last 3 days") to absolute dates.
- Resolve `logRoot` to THIS project's session-log directory (the Claude home `projects/<slugified-cwd>/`).
  Confirm it exists before launch; if empty for the range, say so and stop.
- Keep the information split clean: the skill authors the plan and seeds mine tasks; the workflow's
  agents mine, cluster, propose, validate, and synthesize via the tracker.
- Require citations: every proposal must trace to a session-log record (and, where it points at code,
  a `file:line`). The workflow's adversarial validator enforces this; surface which proposals survived.
- Restrict edit targets to git-tracked artifacts under `plugins/`. A workflow-script improvement must
  target the **golden copy** under `plugins/*/skills/*/reference/*.js` (and bump its version marker),
  never the installed `.claude/workflows/*.js`.
- Writing under `.start/daydream/` is pre-authorized.

**Never:**
- Mine the logs, draft the diffs, or write the report yourself — the workflow owns them.
- Auto-apply a proposal. v1 is advisory: proposals are reported; a human applies them on a branch/PR.
  The bounded-patch and eval-gated auto-apply tiers are the roadmap (see the daydream analysis report
  §5.8), not built yet.
- Launch before the plan (purpose + angles + date range) is approved.
- Overwrite a project's edited `.claude/workflows/daydream.js` without asking (see version check).

## Reference Materials

- [daydream.js](reference/daydream.js) — the golden, reusable workflow script (mine → cluster → propose → validate → synthesize). Read its input-contract comment; it defines the `args` you build.
- [template.md](template.md) — the daydream plan document.

## Workflow

### 1. Frame the Plan and Scaffold

Derive from `$ARGUMENTS`:
- the **date range** — default "yesterday"; resolve to absolute ISO `since`/`until` + a `label`;
- a slug `name` from the range (e.g. `2026-07-16`);
- a terse `description` (the purpose);
- the **signal angles** to mine in parallel — three kinds:
  - `kind: "friction"` (structured markers) — `tool-errors` (most structured), `user-corrections`
    (interrupts/rejections), `retries` (dead-ends), `cost` (expensive/low-cache turns);
  - `kind: "conversation"` — `conversation-quality`: reads the actual dialogue for semantic
    inefficiency the structured angles can't see (shallow answers forcing follow-ups, flip-flopping,
    over-asking/ceremony, ignored-then-restated constraints). Open-ended, runs on the strong model;
  - `kind: "router"` — `model-routing`: over/under-provisioned models per role.
  Each becomes one concurrent mine stream. Default to all three kinds unless the user narrows scope.

Resolve `logRoot` = the Claude home `projects/<slugified-cwd>/` for this project. Create
`.start/daydream/[name]/` and write `daydream.md` from `template.md` with `status: draft`. Create the
`.tracker/` directory.

### 2. Approve or Refine (interactive gate)

**Display the full plan inline in the chat** — do not just write `daydream.md` and point at its path.
Render, as chat text the user reads directly:
- the **purpose** and the **resolved date range** (with the concrete ISO dates),
- the **signal angles** as the independent parallel mine streams the workflow will fan out,
- the **edit-target scope** (git-tracked artifacts under `plugins/`) and that v1 is **advisory** —
  proposals only, a human applies them.

Then ask the user to choose:
  Approve — proceed with this plan
  Refine — adjust range / angles / scope

On refine, edit the `daydream.md` frontmatter, then **re-display the updated plan** and ask again. Loop
until approved, then set `status: approved`.

### 3. Seed the Tracker

For each approved signal angle, create a mine task (record the ids into `angleIds`):

`node task.js create --home .start/daydream/[name] --type mine --name "{angle}"`

Use the specify-meta skill's `task.js` (co-located with spec.py), run via the shell.

### 4. Install / Version-Check the Workflow

The golden copy is this skill's `reference/daydream.js`; its first line is `// daydream-workflow: vN`.

```
installed = <project>/.claude/workflows/daydream.js
match:
  installed absent            => create .claude/workflows/ if needed; copy bundled → installed; tell the user "installed daydream workflow v{N}".
  same version marker         => proceed.
  bundled version > installed => ask under header "Upgrade": "Upgrade daydream workflow v{old} → v{new}?" — Upgrade (Recommended) | Keep current.
  installed version > bundled => project copy is ahead/hand-edited; keep it and note this.
```

Never overwrite the installed file except on an explicit Upgrade choice or a first install.

### 5. Launch the Workflow

Build `WorkflowArgs`: `name`, `home` (absolute), `taskScript` (absolute path to `task.js`),
`description`, `dateRange`, `logRoot` (absolute), `editTargetsRoot` (absolute path to this repo's
`plugins/`), `angles` (from step 3, each `{ id, name, kind }` — `kind: "router"` for the model-routing
angle, `kind: "friction"` for the rest), `reportPath` = `home/report.md`, `applyTier` = `"advisory"`,
`maxProposals` (default 12). Set `status: in-progress`. Launch:

```
Workflow({ scriptPath: "<project>/.claude/workflows/daydream.js", args: <WorkflowArgs> })
```

Present a one-line dispatch summary, e.g.:
`Daydreaming "{name}": mining {angleCount} signal angles over {label} → dynamic-workflow runtime (advisory)`

Record the `runId`. The run executes in the background; report status from the run, do not fabricate progress.

### 6. Report and Finalize

When the run returns (`{ name, dateRange, mined, opportunities, proposals, supported, dropped, reportPath, grade }`):

1. Present the report at `reportPath`: the per-skill friction ranking and the validated proposals
   (diff + evidence + safety tier), grouped by target. Call out how many were dropped in validation, and
   surface the report grade — the derived rubric value `grade.overall` and the per-dimension
   `grade.dimensions` scores (behaviour / improvement / proof / faithfulness); list any still `fail`/`partial`.
2. Set `daydream.md` `status: review`.
3. Ask the user to choose:
   Apply a proposal (on a branch) | Deepen a signal angle | Persist report to `docs/research/` | Done

When the user chooses to apply, do it on a branch (never the default branch), one proposal at a time,
running the target's `evals/evals.json` where one exists — surfacing the result before moving on.
