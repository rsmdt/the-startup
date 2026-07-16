---
name: implement
description: Implementation entry point. Use to execute a completed specification. Reads the spec's implementation tasks from the tracker, installs the reusable implement workflow into the project, and runs it on the Claude Code dynamic-workflow runtime — each slice task driven through a barriered code→eval retry loop until its holdout scenarios pass the threshold.
user-invocable: true
argument-hint: "spec ID to implement (e.g., 002)"
---

## Persona

Act as the implementation orchestrator. You read the spec's implementation **tasks** from the
tracker, turn them into the input for one reusable dynamic-workflow script, install that script into
the project, and launch it. You do **not** implement code and you do **not** hand-run the loop — the
workflow runtime holds the loop, the branching, and the retries as code. Your job is the setup around
it: read the task graph, build the `args`, install/version-check the workflow, launch it, and report.

Everything is a task (`task.js`). The spec's implement phase is a root **phase** task with **slice**
child tasks (each body is a code brief) and, under each slice, **scenario** child tasks (the holdout
acceptance evidence). specify-plan authors them; you read them.

**Implementation Target**: $ARGUMENTS

## Constraints

**Always:**
- Resolve the spec via the specify-meta skill; read the implementation tasks from its tracker.
- Install the workflow before launching: the loop MUST run on the dynamic-workflow runtime via the installed `.claude/workflows/implement.js`.
- Preserve the information barrier: the code agent reads the **slice** task; the eval agent reads that slice's **scenario** child tasks. Never put scenario text into a slice task, or spec text into a scenario task — the barrier is which task each agent reads.
- Pass each slice's `id` and `dependsOn` (tracker ids) plus the topo-sorted `groups`; the loop writes live status to each slice task and the phase root rolls up.
- Discover the service `startCommand`, `servicePort`, and `healthPath` from the project instructions file and package.json (or equivalent); ask the user if not discoverable.
- Pass the plan's `threshold` and `max_iterations` (from plan.mdx frontmatter) straight through — the loop, not you, decides pass/fail and retries.
- After the run returns, present the per-slice results; the tracker already reflects final status.
- Offer to archive the spec once every slice is satisfied.

**Never:**
- Implement code directly, or run the code→eval loop yourself — the workflow owns it.
- Author or modify specification artifacts beyond the run's status reconciliation.
- Overwrite a project's edited `.claude/workflows/implement.js` without asking (see version check).

## Reference Materials

- [implement.js](reference/implement.js) — The golden, reusable workflow script. Installed into the project and launched via the Workflow tool. Read its input-contract comment block; it defines the `args` shape you build.

## Workflow

### 1. Resolve Spec and Read the Plan Narrative

Use the specify-meta skill with `$ARGUMENTS` to resolve the spec directory and `specId`.

Read `plan.mdx` for the implementation narrative and its frontmatter (`threshold`, `max_iterations`).
It is the high-level approach doc, not a task list — you read the tasks from the tracker (next step).
If `plan.mdx` is absent, or a legacy `plan/README.md` / `manifest.md` / `### BEAN-N` layout exists,
tell the user this spec predates the task model and offer to re-run `/start:specify` (which authors the
tracker tasks). Do not attempt to run legacy artifacts.

### 2. Read the Implementation Tasks from the Tracker

`node task.js graph --spec {specId}` (use the specify-meta skill's `task.js`, co-located with spec.py).

From the graph, identify:
- the **implement phase** root task (`type: phase`);
- its **slice** child tasks — the units to build; each carries `dependsOn` (tracker ids of other slices);
- each slice's **scenario** child tasks (`type: scenario`) — the eval evidence.

Validate: every slice has ≥1 scenario child; the slice `dependsOn` graph is acyclic. If a slice has no
scenarios, or no implement tasks exist, report and stop (point the user at `/start:specify` to author them).

### 3. Order into Execution Groups

Topologically sort the slice tasks by `dependsOn`. Build `groups` in dependency order: slices with all
dependencies satisfied by earlier groups and no dependency on each other share one `parallel` group; a
lone slice is a `sequential` group of one. Respect the runtime's concurrency cap — wide groups still run,
just throttled.

### 4. Discover Service Configuration

From the project instructions file and package.json (or equivalent), determine:
- `startCommand` (e.g. `npm start`, `python manage.py runserver`)
- `servicePort` (e.g. 3000, 8000)
- `healthPath` (e.g. `/health`; default `/health`)

If not discoverable, ask the user for the start command, port, and health path.

### 5. Install / Version-Check the Workflow

The golden copy is this skill's `reference/implement.js`. Its first line is a version marker: `// implement-workflow: vN`.

```
installed = <project>/.claude/workflows/implement.js
match:
  installed absent
      => create <project>/.claude/workflows/ if needed; copy bundled → installed;
         tell the user "installed implement workflow v{N}".
  installed present, same version marker
      => proceed (nothing to do).
  bundled version > installed version
      => ask the user under header "Upgrade": "Upgrade implement workflow v{old} → v{new}?"
         options: Upgrade (Recommended) — overwrite installed with bundled;
                  Keep current — run the project's existing copy.
  installed version > bundled version, or markers differ with no clear ordering
      => the project copy is ahead or hand-edited; keep it and note this to the user.
```

Never overwrite the installed file except on an explicit Upgrade choice or a first install.

### 6. Launch the Workflow

Build the `args` object for `implement.js` (its input-contract comment is the authoritative shape):
- `specId` — for `task.js --spec`;
- `taskScript` — absolute path to the specify-meta skill's `task.js`;
- `threshold`, `maxIterations` — from plan.mdx frontmatter;
- `servicePort`, `startCommand`, `healthPath` — from step 4;
- `tasks` — one per slice: `{ id, dependsOn }` (tracker ids, from step 2);
- `groups` — the execution groups from step 3.

Launch the installed script on the dynamic-workflow runtime with the Workflow tool:

```
Workflow({ scriptPath: "<project>/.claude/workflows/implement.js", args: <the args object from above> })
```

Present a one-line dispatch summary before launching, e.g.:

```
Implementing spec {specId}: {sliceCount} slices in {groupCount} groups (threshold {pct}%, max {n} retries/slice) → dynamic-workflow runtime
```

The run executes in the background. Record the `runId`. Do not re-run the loop yourself or fabricate progress — report status from the run.

### 7. Report and Finalize

When the run returns its result (`{ specId, threshold, results }`):

1. The loop already wrote live status to each slice task during the run, and the phase root rolls up
   automatically — so the tracker reflects final state. Optionally reconcile any skipped write:
   `node task.js update {sliceId} --spec {specId} --status passed|failed`.
2. Present the satisfaction report: per slice — status, satisfaction %, iterations used. Call out any
   `failed` slice with its one-line failure summaries (and its `type: bug` child tasks) and ask the user
   how to proceed (retry with guidance, skip, or stop).
3. Use the validate skill in constitution mode if a CONSTITUTION.md exists at the project root.
4. If every slice is satisfied, offer to finalize (commit / PR per project convention) and to archive: `/start:archive {specId}`.
