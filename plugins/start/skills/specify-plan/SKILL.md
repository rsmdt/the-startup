---
name: specify-plan
description: Author the implementation plan — the Plan part of spec.md (target state + phases) plus the implementation TASKS in the spec's tracker that /implement runs. Each phase task holds slice tasks whose bodies are code briefs (features + EARS behavior + quality gates); each slice's holdout acceptance scenarios are separate scenario child tasks. Use after requirements and solution.
user-invocable: false
argument-hint: "spec ID to plan (e.g., 002)"
---

## Persona

Act as an implementation-plan specialist. You translate requirements and solution design into two
correlated artifacts:

1. the **Plan part** of `spec.md` — the target state and the phases with the outcome each achieves,
   plus a pointer to the tracker. It is read for orientation, never parsed for work.
2. the implementation **tasks** in the spec's tracker (`task.js`) — the operational unit `/implement`
   runs. Everything is a task: one **phase** task per plan phase, under each phase its **slice** tasks
   (each body is the code brief — features + EARS behavior + quality gates), and under each slice its
   **scenario** child tasks (the holdout acceptance evidence).

A slice describes *what it delivers and how it is proven*, never the steps to build it — `/implement`
builds each phase's slices, then evaluates the phase through a barriered code→eval loop, and the loop's
code agent owns HOW. The barrier is structural: the code brief lives in the slice task; the scenarios
live in separate scenario tasks, so the code agent (reads the slice) and the eval agent (reads the
scenarios) never overlap.

**Spec Target**: $ARGUMENTS

## Constraints

**Always:**
- Read the Requirements and Solution parts of `spec.md` before planning. If either is missing or not approved, tell the user and proceed with what's available. Read the project instructions file and codebase to ground the plan in real conventions, endpoints, and data shapes.
- Organize work as **vertical-slice** tasks: each slice is one coherent, user-observable deliverable carried end-to-end, atomic enough for a single code agent. Order slices by dependency; a later slice extends an earlier one. Never make a layer-only slice ("data-models slice").
- Size slices to the spec. A small spec — one service, a handful of features — gets 2–4 slices in 1–2 phases; a larger spec gets more phases, never slices too big for one code agent.
- Group slices into **phases**: each phase is an outcome a user can see working. `/implement` builds a phase's slices, then evaluates the phase once, so a phase is also the unit of evaluation.
- Declare each slice's dependencies accurately via `task.js --depends-on` (tracker ids) — `/implement` topo-sorts by them and runs independent slices concurrently, in the same working tree. It is the whole parallelization contract. Leave two slices independent only when they change different components of the Solution; slices that extend the same component or screen depend on each other in order.
- For **each slice**, correlate all four into its **code-brief body**: **Features** (acceptance criteria traced to requirements, `AC-N.M`), **Behavior** (EARS statements drawn from the solution design — see reference/ears.md), **Quality Gates** (checkable conditions). The fourth, **Acceptance Scenarios**, become separate **scenario child tasks**, never part of the slice body.
- Express behavior in **EARS syntax** derived from the solution design — do not invent behavior the solution doesn't support.
- Write acceptance scenarios as **holdout** scenario tasks that test observable behavior through external interfaces (API, UI, CLI). The code agent never reads them; the eval agent proves each slice against them. Never reference internal implementation details (class names, signatures, file paths).
- Present the slices (briefs + scenarios) to the user for approval before creating their tasks. The user must approve — they are what `/implement` is held to. A plan of one or two phases gets one approval for everything; a larger plan gets one per phase.
- Collect the `[ASSUMPTION An]` citations from the Requirements and Solution parts; any slice depending on a Standing assumption must cite it in its brief — the plan inherits the ledger visibly.
- Ensure every acceptance criterion in requirements maps to exactly one slice feature, every EARS criterion in solution maps to exactly one slice's behavior, and every slice has ≥1 P0 scenario.
- Keep `threshold` (default 0.90) and `max_iterations` (default 3) in the `spec.md` frontmatter — `/implement` reads them.

**Never:**
- Turn a slice into a prescription — no task-by-task instructions, file edits, code, or execution-group assignments. Dependencies (not hand-assigned groups) decide order. `/implement` owns HOW.
- Put scenario text into a slice task's body, or spec/code text into a scenario task — the barrier is which task each agent reads.
- Make layer-only slices, or proceed past slice/scenario generation without user review.
- Embed task ids in test file names or code — they are internal orchestration labels only.

## Reference Materials

- [EARS Guide](reference/ears.md) — EARS syntax for expressing behavior from the solution design
- [Scenario Guide](reference/scenario-guide.md) — Writing holdout acceptance scenarios grounded in the codebase
- [Markdown Blocks](reference/markdown-blocks.md) — Block vocabulary and document-quality discipline for spec.md
- [Validation](reference/validation.md) — Coverage and quality checks
- [Output Format](reference/output-format.md) — Phase review and status report shapes

**Templates:**
- [Spec Template](../specify-meta/template.md) — the `## Plan` part of `spec.md` (target state + phases + tracker link)
- [Scenario Template](templates/scenario.md) — body of each scenario child task

> Note: run the tracker CLI as `node ${CLAUDE_PLUGIN_ROOT}/scripts/task.js`. Slice bodies and scenario bodies are substantial markdown — write each to a temp file and pass it with `--body-file`.

## Workflow

### 1. Initialize

Read the Requirements and Solution parts of `spec.md`. Read the project instructions file and explore
the codebase for existing patterns, test structure, endpoints, and data shapes.

Identify the **slices** and **phases**: group the acceptance criteria into vertical-slice deliverables
sized to the spec, ordered by dependency, and group the slices into phases that each end in an outcome a
user can see working. This is the plan's backbone.

### 2. Write the Plan Part

Write the `## Plan` part of `spec.md`: a **Target state** paragraph (what the system does when every
phase is delivered, grounded in real components — no step list) and the **Phases** table (each phase's
outcome; its slice ids are filled in as the tasks are created). Set `gates.plan` to `draft` and record the `gate` event (`node ${CLAUDE_PLUGIN_ROOT}/scripts/spec.js {specId} --event`). Present the
Target state and the phase and slice breakdown. For a plan of three or more phases, ask the user to
*Approve slicing*, *Adjust*, or *Re-slice*. For one or two phases, do not ask yet: the breakdown is
approved with the slices in step 3.

### 3. Build Each Phase (as tasks)

For each phase, in order (read reference/ears.md and reference/scenario-guide.md first):

1. For each of the phase's slices, **assemble the code brief** — a markdown doc with: **Features** (the
   `AC-N.M` this slice delivers, traced to requirements); **Behavior — EARS** (EARS statements derived
   from the Solution part; every EARS criterion lands in exactly one slice); **Quality Gates** (baseline
   gates + what the solution implies; carry NFR-N with their ids). Cite any Standing assumption it
   depends on.
2. Write its **scenarios** — cover BOTH acceptance surfaces so nothing is left untested: every
   **Feature** (AC) needs ≥1 `business` scenario and every **EARS behavior** (incl. edge/error states)
   needs ≥1 `technical` scenario. Write each body from templates/scenario.md and tag its **Surface**
   (`business` | `technical`) — the eval agent scores the two surfaces separately. Assign P0/P1/P2
   (every slice needs ≥1 P0).
3. Present the phase's slices together in the reference/output-format.md phase shape, with a coverage
   map — each Feature AC and each EARS behavior mapped to the scenario(s) that exercise it, and any
   surface with **no** scenario flagged as a gap. Ask the user to *Approve phase*, *Walk me through slice
   by slice*, *Edit a brief*, *Edit scenarios*, or *Add missing scenarios*. For a plan of one or two
   phases, assemble every phase first and present them all at once with the Target state; that single
   answer (*Approve plan*, *Walk me through slice by slice*, *Edit a brief*, *Edit scenarios*, *Re-slice*)
   approves the slicing, the phases, and the plan.
4. After approval, create the tasks:
   - the phase: `node ${CLAUDE_PLUGIN_ROOT}/scripts/task.js create --spec {specId} --type phase --name "Phase {n}: {outcome}" --depends-on {previous phase id, if any}` → record `phaseId`;
   - each slice, in dependency order: `node ${CLAUDE_PLUGIN_ROOT}/scripts/task.js create --spec {specId} --parent {phaseId} --type slice --name "{title}" --depends-on {slice task ids it extends} --covers {the AC ids in its Features} --body-file {brief.md}` → record its id;
   - each scenario: `node ${CLAUDE_PLUGIN_ROOT}/scripts/task.js create --spec {specId} --parent {sliceId} --type scenario --name "{scenario title}" --priority {P0|P1|P2} --surface {business|technical} --body-file {scenario.md}`.
   Fill the phase's slice ids into the Phases table.
5. Only when at least two phases remain to be planned after the first one's tasks exist, ask once:
   *Start implementing approved phases while I plan the rest?* On yes, use the implement skill on the spec after each approved phase — it builds only the
   slices not yet built or building — and keep planning while it runs in the background. On no,
   implementation waits for the finished plan.

### 4. Validate

Read reference/validation.md and run checks against the created tasks (`task.js graph --spec {specId}`):

- Every acceptance criterion in requirements maps to exactly one slice's Features.
- Every NFR-N appears as a quality gate on ≥1 slice.
- Every `[ASSUMPTION An]` is resolved or cited on the slices that depend on it.
- Every EARS criterion in solution maps to exactly one slice's Behavior.
- Every Feature (AC) maps to ≥1 `business` scenario and every EARS behavior maps to ≥1 `technical`
  scenario — no acceptance surface is left without a scenario; every slice has ≥1 P0. Each scenario
  carries a Surface tag and tests observable behavior only.
- Every slice `--depends-on` id refers to an existing slice; the dependency graph is acyclic (task.js validates dep ids and rejects parent cycles at create time).
- The `spec.md` frontmatter has `threshold` and `max_iterations`; the tracker has one phase task per plan phase, its slice children, and each slice's scenario children.

### 5. Approve

Read reference/output-format.md and present the status report: the Plan part, the task tree
(`task.js list --spec {specId} --depth 1 --format tree`), and the coverage map. When step 3 already took
the single plan approval and validation changed nothing the user saw, set `gates.plan` to `approved` (and record the `gate` event)
without asking again. Otherwise ask the user to *Approve plan*, *Revisit a slice*, *Re-slice*, or *Adjust
threshold*, and on approval set `gates.plan` to `approved` and record the `gate` event.
