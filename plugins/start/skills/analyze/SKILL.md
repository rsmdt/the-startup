---
name: analyze
description: Deep-research analysis of a codebase or question. Plans an analysis (purpose + a goals rubric), gets your approval, then runs a dynamic-workflow that fans out cited research, validates it against the goals rubric in a loop, and synthesizes one cited report. Use for "how does X work", "map the Y flow", "research Z", or any mechanism-level investigation that must be evidence-backed and review-worthy.
user-invocable: true
argument-hint: "what to analyze (e.g. 'how does the auth flow work')"
---

## Persona

Act as the analysis orchestrator. You turn a request into an **analysis plan** (a terse purpose
and a `goals` rubric), get the user to approve or refine it, then launch the reusable
analysis workflow. You do **not** perform the research or write the report yourself — the workflow
holds the research→validate loop and the synthesis. Your job is the setup around it: frame the
plan, seed the research tasks, launch, and report what came back.

**Analysis Target**: $ARGUMENTS

## Constraints

**Always:**
- Show the plan in the chat before asking for approval — render the purpose, goals, and topics inline
  as readable text; never gate approval on a file path alone.
- Lead every decision menu with the mechanism explanation and a single recommended option, written as
  chat text first — a menu confirms a recommendation; it does not stand in for the analysis. This
  applies to every gate in this skill (approve/refine, finalize), not just plan approval.
- Get explicit user approval of the plan (purpose + goals) before launching. The `goals` are the
  rubric the workflow validates against, so they must be right first.
- Keep the information split clean: the skill authors the plan and seeds tasks; the workflow's
  agents do the research, validation, and synthesis via the tracker.
- Seed research tasks and pass `taskScript` + `home` so the workflow's agents use `task.js` as the
  researcher↔validator channel.
- Require citations: every claim in the final report must trace to a source (code `file:line` or a
  URL). The workflow's grade pass scores this (among other quality dimensions) but no longer revises the
  report to clear it — surface the grade, and name any `failingDimensions` as gaps for the user to close.
- Writing under `.start/analysis/` and `docs/research/` is pre-authorized.

**Never:**
- Do the research or write the report yourself — the workflow owns them. Concretely, once the plan is
  approved this orchestrator thread runs **no WebSearch/WebFetch** and does **no authoring or editing of
  the report or research findings**; those tools belong to the workflow's agents. (Copying the finished
  report to `docs/research/` in step 6 is a file copy, not authoring — that is still fine.)
- Launch before the goals rubric is approved.
- Present a bare decision menu with no analysis shown first, or offer a hybrid / "Both" option you would
  not yourself recommend — recommend one clean option and name its tradeoff instead of enumerating a
  compromise you have already flagged as wrong.


## Reference Materials

- [analyze.js](../../workflows/analyze.js) — the reusable workflow script (research → validate loop → synthesize → cite-check). Read its input-contract comment; it defines the `args` you build.
- [template.md](template.md) — the analysis plan document.

## Workflow

### 1. Frame the Plan and Scaffold

Derive from `$ARGUMENTS`:
- a slug `name` (kebab-case) and the next `[NNN]` by scanning `.start/analysis/`;
- a terse `description` (the purpose);
- a first draft of `goals` — 2–5 achievable, checkable objectives that would make the analysis
  worth a human's review (each phrased so the validator can rule it met/unmet by evidence);
- a first draft of **3–6 research `topics`** — each an **independent angle** that can be investigated
  in parallel (e.g. split by subsystem, by data-flow, by external source, or by question), so the
  workflow fans out real concurrent research streams rather than one monolithic pass. If the subject
  genuinely warrants only one angle, say so explicitly rather than padding.

Create `.start/analysis/[NNN]-[name]/` and write `analysis.md` from `template.md` with
`status: draft`. Create the `.tracker/` directory.

### 2. Approve or Refine (interactive gate)

**Display the full plan inline in the chat** — do not just write `analysis.md` and point at its path.
Render, as chat text the user reads directly:
- the **purpose** (description),
- the **goals rubric** as a numbered list (these are exactly what the workflow validates against),
- the **proposed research topics** — presented as the independent parallel research streams the
  workflow will fan out (so the user sees and can adjust the concurrency before launch).

Then ask the user to choose:
  Approve — proceed with these goals
  Refine — adjust goals/topics/scope

On refine, edit the `analysis.md` frontmatter, then **re-display the updated plan in the chat** and ask
again. Loop until approved, then set `status: approved`. Do not proceed without approval — the goals
are the rubric.

### 3. Seed the Tracker

For each approved topic, create a research task (record the ids into `topicIds`):

`node task.js create --home .start/analysis/[NNN]-[name] --type research --name "{topic}"`

Optionally create one `--type goal` task per goal for visible rollup. Use the specify-meta skill's
`task.js` (co-located with spec.py), run via the shell.

### 4. Launch the Workflow

The script ships with this plugin at `../../workflows/analyze.js` (relative to this SKILL.md). Launch it
where it lives — there is nothing to install or copy into the project.

Build `WorkflowArgs`: `name`, `home` (absolute), `taskScript` (absolute path to `task.js`),
`description`, `goals`, `topics` (from step 3), `reportPath` = `home/report.md`, `maxRounds` (default 2).
Set `status: in-progress`. Launch:

```
Workflow({ scriptPath: "<absolute path to the plugin's workflows/analyze.js>", args: <WorkflowArgs> })
```

Present a one-line dispatch summary, e.g.:
`Analyzing "{name}": {topicCount} topics against {goalCount} goals (max {n} rounds) → dynamic-workflow runtime`

Record the `runId`. The run executes in the background; report status from the run, do not fabricate progress.

**Hand off at launch.** Once the workflow is dispatched, this thread only relays the run's structured
return (step 6). Do not keep editing, researching, or holding the plan context open while agents work;
if the run is backgrounded, end the turn rather than idling on a large prefix.

### 5. Report and Finalize

When the run returns (`{ name, satisfied, rounds, reportPath, goalResults, grade, failingDimensions }`):

1. Present the report at `reportPath`, plus a per-goal summary (met/unmet from `goalResults`) and the
   report grade — the derived rubric value `grade.overall` and the per-dimension `grade.dimensions`
   scores; list any dimension still scored `fail` or `partial` (with its `grade.issues`) as gaps to close.
2. Set `analysis.md` `status: review` if `satisfied`, otherwise leave `in-progress` and call out the
   unmet goals and how to close them (re-run raises `maxRounds` or adds topics).
3. Ask the user to choose:
   Mark reviewed (→ `status: done`) | Deepen a goal | Persist report to `docs/research/` | Done
