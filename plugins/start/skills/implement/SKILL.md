---
name: implement
description: Implementation entry point. Use to execute a completed specification. Reads the spec's implementation tasks from the tracker and runs the reusable implement workflow on the Claude Code dynamic-workflow runtime — each slice task driven through a barriered code→eval retry loop until its holdout scenarios pass the threshold.
user-invocable: true
argument-hint: "spec ID to implement (e.g., 002)"
---

## Persona

Act as the implementation orchestrator. You read the spec's implementation **tasks** from the
tracker, turn them into the input for one reusable dynamic-workflow script, and launch it. You do **not** implement code and you do **not** hand-run the loop — the
workflow runtime holds the loop, the branching, and the retries as code. Your job is the setup around
it: read the task graph, build the `args`, launch the workflow, and report.

Everything is a task (`task.js`). The spec's implement phase is a root **phase** task with **slice**
child tasks (each body is a code brief) and, under each slice, **scenario** child tasks (the holdout
acceptance evidence). specify-plan authors them; you read them.

**Implementation Target**: $ARGUMENTS

## Constraints

**Always:**
- Resolve the spec via the specify-meta skill; read the implementation tasks from its tracker.
- Launch on the dynamic-workflow runtime: the loop MUST run via the Workflow tool on this skill's `../../workflows/implement.js`.
- Preserve the information barrier: the code agent reads the **slice** task; the eval agent reads that slice's **scenario** child tasks. Never put scenario text into a slice task, or spec text into a scenario task — the barrier is which task each agent reads.
- Pass each slice's `id` and `dependsOn` (tracker ids) plus the topo-sorted `groups`; the loop writes live status to each slice task and the phase root rolls up.
- Discover the service `startCommand`, `servicePort`, and `healthPath` from the project instructions file and package.json (or equivalent); ask the user if not discoverable.
- Pass the plan's `threshold` and `max_iterations` (from plan.mdx frontmatter) straight through — the loop, not you, decides pass/fail and retries.
- After the run returns, present the per-slice results; the tracker already reflects final status.
- Offer to archive the spec once every slice is satisfied.

**Never:**
- Implement code directly, or run the code→eval loop yourself — the workflow owns it.
- Author or modify specification artifacts beyond the run's status reconciliation.

## Reference Materials

- [implement.js](../../workflows/implement.js) — The reusable workflow script, launched in place via the Workflow tool. Read its input-contract comment block; it defines the `args` shape you build.

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
- `migrateCommand` — optional; how pending DB migrations are applied (e.g. `supabase migration up --local`).
  It must be a **directly-executable binary**: resolve it with `which <bin>` and pass the resolved command.
  Never pass an `npx` / `pnpm dlx` / version-manager-shim form — agent shells frequently cannot resolve
  those, and the preflight step then reports a fault (or, worse in earlier versions, silently applied
  nothing). Run it once yourself before launching; if it errors, fix it or omit it rather than passing a
  command that will fail inside every loop iteration.
- `requiredEnv` — optional; env vars the acceptance scenarios depend on (feature flags, engine switches).

If not discoverable, ask the user for the start command, port, and health path.

### 5. Launch the Workflow

The script ships with this plugin at `../../workflows/implement.js` (relative to this SKILL.md). Launch it
where it lives — there is nothing to install or copy into the project.

Build the `args` object for `implement.js` (its input-contract comment is the authoritative shape):
- `specId` — for `task.js --spec`;
- `taskScript` — absolute path to the specify-meta skill's `task.js`;
- `threshold`, `maxIterations` — from plan.mdx frontmatter;
- `servicePort`, `startCommand`, `healthPath` — from step 4;
- `migrateCommand`, `requiredEnv` — from step 4 (omit when the project has neither);
- `stallLimit` — optional; consecutive non-improving iterations before a slice halts as `stalled` (default 2);
- `tasks` — one per slice: `{ id, dependsOn }` (tracker ids, from step 2);
- `groups` — the execution groups from step 3.

Launch it on the dynamic-workflow runtime with the Workflow tool:

```
Workflow({ scriptPath: "<absolute path to the plugin's workflows/implement.js>", args: <the args object from above> })
```

Present a one-line dispatch summary before launching, e.g.:

```
Implementing spec {specId}: {sliceCount} slices in {groupCount} groups (threshold {pct}%, max {n} retries/slice) → dynamic-workflow runtime
```

The run executes in the background. Record the `runId`. Do not re-run the loop yourself or fabricate progress — report status from the run.

### 6. Report and Finalize

When the run returns its result (`{ specId, threshold, results, harnessAbort, status? }`):

A run-level `status` is present ONLY when the run aborted before building anything — `blocked-environment`
(the preflight gate found an invalid harness) or `harness-abort` (agents stopped being served during the
coverage audit). In both cases `results` is empty: report why nothing was built and stop; do not present
it as zero slices completed.

1. The loop already wrote live status to each slice task during the run, and the phase root rolls up
   automatically — so the tracker reflects final state. Optionally reconcile any skipped write:
   `node task.js update {sliceId} --spec {specId} --status passed|failed`.
2. Present the satisfaction report: per slice — status, **business %** and **technical %** (the two
   acceptance surfaces are scored and gated separately), iterations used. Call out any `failed` slice
   with its one-line failure summaries (and its `type: bug` child tasks), and any `incomplete-coverage`
   slice with its `coverageGaps` (features/EARS behaviors that had **no scenario**) and its
   `notExercisable` entries (scenarios no external interface can induce, filed as `type: gap` child
   tasks) — both are spec gaps, not build failures.

   Three statuses mean the code was never judged, and each needs a different response — never report
   them as build failures:
   - `blocked-environment` — the harness was invalid (service down, migration unapplied, env var unset).
     Report its `environment` faults and **stop**; no retry was spent. Do not relaunch until the user
     confirms the environment is fixed.
   - `stalled` — the loop halted early because it stopped converging (no score improvement, or one
     symptom survived three attempts). Retries remain unspent deliberately. Report what was tried
     (`history`) and ask for guidance rather than relaunching the same slice unchanged.
   - `harness-abort` — agents stopped being served (quota / session limit). The slice was not judged and
     its retry budget is intact; offer to re-launch the same slice ids once the harness is available.

   Ask the user how to proceed (retry with guidance, add the missing scenarios, skip, or stop).
3. Use the validate skill in constitution mode if a CONSTITUTION.md exists at the project root.
4. If every slice is satisfied, offer to finalize (commit / PR per project convention) and to archive: `/start:archive {specId}`.
