---
name: implement
description: Implementation entry point. Use to execute a completed specification. Reads the spec's slice and scenario tasks from the tracker, then runs a dynamic workflow that builds each phase's slices and evaluates the phase through a barriered code→eval retry loop until every slice's holdout scenarios pass the threshold.
user-invocable: true
argument-hint: "spec ID to implement (e.g., 002)"
---

## Persona

Act as the implementation orchestrator. You resolve the spec, gather everything the build loop needs,
then author and run **one dynamic workflow** (the Workflow tool with an inline script you write) that
builds every slice. You do **not** write product code and you do **not** hand-run the loop — the
workflow holds the loop, the branching, and the retries as code.

Everything is a task (`task.js`). Each plan phase is a root **phase** task with **slice** child tasks
(each body is a code brief) and, under each slice, **scenario** child tasks (the holdout acceptance
evidence). specify-plan authors them; you read their structure.

**Implementation Target**: $ARGUMENTS

## Constraints

**Always:**
- Resolve the spec via the specify-meta skill; read the task graph with `task.js graph`.
- Read task **structure only** — ids, `dependsOn`, scenario counts. Never open a slice or scenario body yourself: you write every agent's prompt, so whatever you read can leak across the barrier.
- Preserve the information barrier inside the workflow: the code agent reads only its **slice** task; the eval agent reads only that slice's **scenario** child tasks.
- Take `threshold` and `max_iterations` from the frontmatter of the spec's `plan_doc` (`spec.md`; `plan.md` in a legacy spec); the workflow, not you, decides pass/fail and retries.
- Ask the user for service configuration you cannot discover — before launching, never mid-run.
- Present the per-slice results when the workflow returns; the tracker already reflects final status.

**Never:**
- Implement code directly, or run the code→eval loop outside the workflow.
- Author or modify specification artifacts beyond the run's status reconciliation.
- Launch with a slice that has no scenario child task, or with an unresolved service configuration.

## Reference Materials

- [Loop Contract](reference/loop.md) — the workflow's stages, agent return schemas, decision rules, and result shape
- [Code Agent](reference/code-agent.md) — prompt for the build stage
- [Eval Agent](reference/eval-agent.md) — prompt for the evaluation stage
- [Coverage Agent](reference/coverage-agent.md) — prompt for the spec-coverage audit
- [Preflight Agent](reference/preflight-agent.md) — prompt for the harness assertion
- [Service Agent](reference/service-agent.md) — prompt for restarting the service before each evaluation

## Workflow

### 1. Resolve Spec

Use the specify-meta skill with `$ARGUMENTS` to resolve the spec directory and `specId`. If no spec
exists, say so and suggest `/start:specify`; stop.

Read the `plan_doc` frontmatter for `threshold` (default 0.90) and `max_iterations` (default 3). If the
spec has no `plan_doc`, tell the user the spec has no plan and offer to run `/start:specify`; stop.

### 2. Read the Task Graph

`node ${CLAUDE_PLUGIN_ROOT}/scripts/task.js graph --spec {specId}`. Pass that absolute `task.js` path to the
workflow — its agents call it.

Identify the **phase** roots (id, `dependsOn`), each phase's **slice** children (id, `dependsOn`,
status), and each slice's count of **scenario** children. Leave out slices whose status is `passed`
(already built) or `building` (another run is building them), and phases left with no slice — the plan
may still be growing while earlier phases are built. Stop and point the user at `/start:specify` if
there is no phase root, no slice to build, a slice with zero scenarios, or a dependency cycle.

### 3. Order into Phases and Groups

Order the phases by their `dependsOn`. Within each phase, topologically sort its slices by `dependsOn`:
slices whose dependencies are all satisfied by earlier groups share one parallel group; a dependency on
a slice of an earlier phase, or on a non-slice id, does not block.

### 4. Discover Service Configuration

From the project instructions file and package.json (or equivalent) determine `startCommand`,
`servicePort`, `healthPath` (default `/health`), `requiredEnv` (env vars the scenarios depend on), and
`migrateCommand` when the project has database migrations. `migrateCommand` must be a directly
executable binary resolved with `which` — never an `npx` / `pnpm dlx` / version-manager shim. Run it once
yourself; fix or omit it if it errors. Ask the user for anything you could not read from a file.

Start the service once, configured, before launching. Every eval agent evaluates against it read-only.

### 5. Run the Dynamic Workflow

These instructions are the opt-in: call the Workflow tool directly with an inline script. Do not load the workflow-authoring reference or assemble the script through shell commands first — runs that detour through either stall before launching. Read reference/loop.md and the five agent references. Author one workflow script that implements the
loop contract exactly, with each agent's prompt taken from its reference and its placeholders filled
from steps 1–4. Launch it with the Workflow tool, passing the resolved values as `args`.

Choose `model` and `effort` per stage from what the stage does. Mechanical stages — running a command,
checking a URL, mapping text into a fixed schema — can run on a smaller, faster model at lower effort.
Stages that write code or judge another agent's work need the strongest model. A judging stage never
runs on a weaker model than the stage it judges. When unsure, omit `model` so the stage inherits the
session model.

Present a one-line dispatch summary before launching:

```
Implementing spec {specId}: {sliceCount} slices in {phaseCount} phases (threshold {pct}%, max {n} iterations/slice) → dynamic workflow
```

The run executes in the background. Do not re-run the loop yourself or fabricate progress.

### 6. Report and Finalize

When the workflow returns (shape in reference/loop.md):

1. A run-level `status` means nothing was built — report why and stop; never present it as zero slices
   completed.
2. Present per slice: status, **business %** and **technical %**, iterations used. For `failed`, show the
   symptoms and its `type: bug` child tasks. For `incomplete-coverage`, show `coverageGaps` and
   `notExercisable` — both are spec gaps, not build failures.
   An `incomplete-coverage` slice cleared both thresholds: it is built. Record its gaps as known gaps in
   the Ledger part of `spec.md` and do not relaunch it. Recommend adding scenarios and relaunching only for
   `business` gaps — a missing proof of a requirement; `technical` gaps are recorded, not re-run.
3. Three statuses mean the code was never judged:
   - `blocked-environment` — report the faults and stop; do not relaunch until the user confirms the
     environment is fixed.
   - `stalled` — the loop stopped converging with retries left. Report the attempt history and ask for
     guidance rather than relaunching unchanged.
   - `harness-abort` — agents stopped being served. Offer to relaunch the same slices once available.
4. Ask how to proceed: retry with guidance, add the missing scenarios, skip, or stop.
5. Use the validate skill in constitution mode if a CONSTITUTION.md exists at the project root.
6. If every slice is `completed`, offer to finalize (commit / PR per project convention) and to archive:
   `/start:archive {specId}`.
