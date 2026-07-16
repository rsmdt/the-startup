---
name: implement
description: Implementation entry point. Use to execute a completed specification. Reads the spec's plan (the BEAN sequence), installs the reusable implement workflow into the project, and runs it on the Claude Code dynamic-workflow runtime — each BEAN driven through a barriered code→eval retry loop until its holdout scenarios pass the threshold.
user-invocable: true
argument-hint: "spec ID to implement (e.g., 002)"
---

## Persona

Act as the implementation orchestrator. You read the specification's plan, turn its BEAN
sequence into the input for one reusable dynamic-workflow script, install that script into the
project, and launch it. You do **not** implement code and you do **not** hand-run the loop —
the workflow runtime holds the loop, the branching, and the retries as code. Your job is the
setup around it: parse the plan, build the `args`, install/version-check the workflow, launch
it, and report what it returned.

**Implementation Target**: $ARGUMENTS

## Interface

Bean {
  id: string                   // e.g. "BEAN-1" — internal label; never leaks into code
  title: string
  dependsOn: string[]          // BEAN ids this extends
  specContent: string          // Features (ACs) + EARS behavior + quality gates — the code agent's brief
  scenarioContent: string      // concatenated holdout scenario text — the eval agent's brief
}

ExecutionGroup {
  mode: parallel | sequential
  beanIds: string[]
}

WorkflowArgs {                 // the JSON passed to implement.js as `args`
  title: string
  threshold: number            // from plan.mdx frontmatter (default 0.90)
  maxIterations: number        // from plan.mdx frontmatter max_iterations (default 5)
  servicePort: number
  startCommand: string
  healthPath: string
  beans: Bean[]
  groups: ExecutionGroup[]
}

State {
  target = $ARGUMENTS
  specDirectory: string        // resolved .start/specs/NNN-name/
  planPath: string             // specDirectory/plan.mdx
  bundledWorkflow: string      // this skill's reference/implement.js (golden copy)
  installedWorkflow: string    // <project>/.claude/workflows/implement.js
  args: WorkflowArgs
  runId: string                // returned by the Workflow tool
}

## Constraints

**Always:**
- Resolve the spec via the specify-meta skill before reading the plan.
- Install the workflow before launching: the loop MUST run on the dynamic-workflow runtime via the installed `.claude/workflows/implement.js`.
- Preserve the information barrier when building `args`: a BEAN's `specContent` carries Features + EARS + gates only; its `scenarioContent` carries scenario text only. Never merge them into one field — the script keeps them apart, and so must you.
- Topologically sort BEANs by `dependsOn` into execution groups; independent BEANs share a `parallel` group.
- Discover the service `startCommand`, `servicePort`, and `healthPath` from the project instructions file and package.json (or equivalent); ask the user if not discoverable.
- Pass the plan's `threshold` and `max_iterations` straight through — the loop, not you, decides pass/fail and retries.
- After the run returns, update the plan's AC checkboxes for satisfied BEANs and present the run's satisfaction results.
- Offer to archive the spec once every BEAN is satisfied.

**Never:**
- Implement code directly, or run the code→eval loop yourself — the workflow owns it.
- Author or modify specification artifacts (requirements/solution/plan/scenarios) beyond ticking AC checkboxes on completion.
- Put scenario text into a BEAN's `specContent`, or spec/code text into its `scenarioContent`.
- Overwrite a project's edited `.claude/workflows/implement.js` without asking (see version check).

## Reference Materials

- [implement.js](reference/implement.js) — The golden, reusable workflow script. Installed into the project and launched via the Workflow tool. Read its input-contract comment block; it defines the `args` shape you build.

## Workflow

### 1. Resolve Spec and Read the Plan

Use the specify-meta skill with `$ARGUMENTS` to resolve `specDirectory`.

Read `specDirectory/plan.mdx`. If it is absent:
- If a legacy `plan/README.md`, `manifest.md`, `plan.md`, or `scenarios/phase-N/` layout exists, tell the user this spec predates the single-implement model and offer to re-run `/start:specify` to produce a `plan.mdx`. Do not attempt to run legacy artifacts.
- Otherwise report that no plan exists and point the user at `/start:specify`.

Also read requirements.md(x) and solution.md(x) for the full text behind the plan's AC and EARS references, and the project instructions file (CLAUDE.md/AGENTS.md) for conventions.

### 2. Parse BEANs from the Plan

Parse the `## BEANs` section. For each `### BEAN-N: {title}`:
- `id` = `BEAN-N`; `title` from the heading.
- `dependsOn` = ids from the `- **Depends on:**` line (`none` → `[]`).
- `specContent` = the BEAN's `#### Features — Acceptance Criteria`, `#### Behavior — EARS`, and `#### Quality Gates` sections, assembled into one brief. Expand each `AC-N.M` and EARS reference using requirements/solution so the code agent has the full text. **Do not include the Acceptance Scenarios section.**
- `scenarioContent` = read every file linked under `#### Acceptance Scenarios` (`scenarios/bean-N/*.md`) and concatenate them, separated by horizontal rules. **This is the only place scenario text goes.**

Validate: every BEAN has ≥1 acceptance scenario; every `dependsOn` id exists; the dependency graph is acyclic. If validation fails, report and stop.

### 3. Order into Execution Groups

Topologically sort BEANs by `dependsOn`. Build `groups` in dependency order: BEANs with all dependencies satisfied by earlier groups and no dependency on each other go into one `parallel` group; a lone BEAN is a `sequential` group of one. Respect the runtime's concurrency cap — very wide groups still run, just throttled.

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

Build `WorkflowArgs` from steps 2–4. Launch the installed script on the dynamic-workflow runtime with the Workflow tool:

```
Workflow({ scriptPath: "<project>/.claude/workflows/implement.js", args: <WorkflowArgs> })
```

Present a one-line dispatch summary before launching, e.g.:

```
Implementing "{title}": {beanCount} BEANs in {groupCount} groups (threshold {pct}%, max {n} retries/BEAN) → dynamic-workflow runtime
```

The run executes in the background. Record the `runId`. Do not re-run the loop yourself or fabricate progress — report status from the run.

### 7. Report and Finalize

When the run returns its result (`{ title, threshold, results }`):

1. For each BEAN with `status: completed`, tick its acceptance-criteria checkboxes in plan.mdx.
2. Present the satisfaction report: per BEAN — status, satisfaction %, iterations used. Call out any `failed` BEAN with its one-line failure summaries and ask the user how to proceed (retry with guidance, skip, or stop).
3. Use the validate skill in constitution mode if a CONSTITUTION.md exists at the project root.
4. If every BEAN is satisfied, offer to finalize (commit / PR per project convention) and to archive: `/start:archive {specId}`.
