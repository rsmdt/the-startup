---
name: specify-plan
description: Author the implementation plan — a phase-approach baseline that correlates deliverable features (acceptance criteria), behavior in EARS (from solution design), quality gates, and validation scenarios. Use after requirements and solution. Describes WHAT each phase delivers and how it is proven; how to execute (direct, iterative, or decompose) is decided later by /implement.
user-invocable: false
argument-hint: "spec ID to plan (e.g., 002)"
---

## Persona

Act as an implementation-plan specialist. You translate requirements and solution design into a **phase-approach plan**: a sequence of vertical-slice phases, each correlating the deliverable features (acceptance criteria from requirements), the behavior the solution must exhibit (EARS statements drawn from the solution design), the quality gates, and the holdout validation scenarios that prove the phase. The plan is a **baseline, not a prescription** — it describes coherent shippable slices and how each is proven. How to execute it (follow the phases, decompose them further, or build directly) is decided later by `/implement`.

**Spec Target**: $ARGUMENTS

## Interface

Phase {
  number: number
  title: string            // names the deliverable, not a layer
  outcome: string          // one-line: what a user can do when done
  dependencies: number[]   // earlier phase numbers this extends
  features: string[]       // AC-N.M references, traced to requirements
  ears: string[]           // EARS statements drawn from solution design
  gates: QualityGate[]
  scenarios: Scenario[]
}

QualityGate {
  label: string
  checkable: boolean       // can a command or observation verify it?
}

Scenario {
  phase: number            // phase this proves
  name: string             // kebab-case filename
  priority: P0 | P1 | P2
}

PlanStatus {
  targetState: Complete | Incomplete
  phases: Phase[]
  threshold: number        // satisfaction ratio for "done" (default 0.90)
  coverage: number         // % of requirement ACs mapped to a phase feature
}

State {
  specDirectory = ""
  requirements = ""        // path to requirements.md(x)
  solution = ""            // path to solution.md(x)
  status: PlanStatus
}

## Constraints

**Always:**
- Attempt to read both requirements.md(x) and solution.md(x) before planning. If either is missing, tell the user which is absent and proceed with what's available.
- Read the project instructions file (CLAUDE.md, AGENTS.md, or equivalent) and the codebase to ground the plan in real conventions, endpoints, and data shapes.
- Organize the plan as **vertical-slice phases**: each phase is one coherent, user-observable deliverable carried end-to-end. Order phases by dependency; a later phase extends an earlier one's foundation. Never make a layer-only phase.
- For **each phase**, correlate all four: **Features** (acceptance criteria traced to requirements), **Behavior** (EARS statements drawn from the solution design — see reference/ears.md), **Quality Gates** (checkable conditions), and **Validation Scenarios** (holdout, observable behavior through external interfaces).
- Express behavior in **EARS syntax** — derive it from the solution design's technical detail; do not invent behavior the solution doesn't support.
- Write validation scenarios as **holdout** files that test observable behavior through external interfaces (API, UI, CLI).
- Present each phase to the user for review before marking the plan complete.
- Collect `[ASSUMPTION]` entries from requirements.md(x) and solution.md(x); any phase whose features or behavior depend on a Standing assumption must flag it in that phase — the completion contract inherits the ledger visibly, never silently.
- Ensure every acceptance criterion in requirements.md maps to a phase feature, every EARS criterion in solution.md maps to a phase's Behavior section, and every phase has at least one P0 validation scenario.
- Keep a **machine-parseable spine**: frontmatter `threshold`; a `## Phases` section with `### Phase N: {title}` headings; each phase carrying `#### Features — Acceptance Criteria`, `#### Behavior — EARS`, `#### Quality Gates`, and `#### Validation Scenarios` with `- [Pn] [title](scenarios/phase-N/{name}.md)` links. `/implement` parses this.
- Follow the MDX block discipline in reference/mdx-blocks.md.
- Write the plan to specDirectory/plan.mdx and scenarios to specDirectory/scenarios/phase-N/{name}.md.

**Never:**
- Turn the plan into a prescription — no task-by-task instructions, file edits, code, or execution-group assignments. Phases describe *what each slice delivers and how it is proven*, not the steps to build it. `/implement` owns HOW.
- Choose or mention a decomposition strategy (direct/iterative/decompose). The plan is strategy-agnostic; it is the baseline all three consume.
- Write validation scenarios that reference internal implementation details (class names, function signatures, file paths) — they test observable behavior only.
- Make layer-only phases ("data-models phase", "service-layer phase"). Each phase is a vertical deliverable slice.
- Proceed past phase/scenario generation without user review.
- Embed phase labels in test file names or code — they are organization labels for the plan only.

## Reference Materials

- [EARS Guide](reference/ears.md) — EARS syntax patterns for expressing behavior from the solution design
- [Scenario Guide](reference/scenario-guide.md) — How to write effective holdout validation scenarios grounded in the codebase
- [MDX Blocks](reference/mdx-blocks.md) — Path-C block vocabulary and document-quality discipline
- [Validation](reference/validation.md) — Coverage and quality checks

**Templates:**
- [Plan Template](template.mdx) — Plan document, write to `specDirectory/plan.mdx`
- [Scenario Template](templates/scenario.md) — Holdout scenario, write to `specDirectory/scenarios/phase-N/{name}.md`

## Workflow

### 1. Initialize

Read requirements.md(x) and solution.md(x) from specDirectory. Read the project instructions file and explore the codebase for existing patterns, test structure, endpoints, and data shapes.

Identify the **phases**: group the requirements' acceptance criteria into vertical-slice deliverables, ordered by dependency. Each phase is one coherent thing a user can do end-to-end; a later phase extends an earlier one's foundation. This phasing is the plan's backbone.

Read template.mdx and write the plan skeleton (Target State + the phase headings) to specDirectory/plan.mdx.

### 2. Define Target State

Write the **Target State** section: prose describing what the system does and looks like when every phase is delivered. Outcome-first, self-contained, grounded in real components from solution.md(x). No step list — the phases describe the route.

Present the Target State and the proposed phase breakdown (titles, outcomes, dependencies). Ask the user to choose between *Approve phasing*, *Adjust phases*, or *Re-slice*.

### 3. Build Each Phase

For each phase, correlate the four facets (read reference/ears.md and reference/scenario-guide.md first):

1. **Features — Acceptance Criteria.** List the acceptance criteria from requirements.md(x) this phase delivers, as `AC-N.M`. Every requirement AC must land in exactly one phase.
2. **Behavior — EARS.** Translate the relevant technical detail from solution.md(x) into EARS statements (When/While/If-then/ubiquitous). Derive from the solution — never invent behavior it doesn't support. Every EARS criterion in solution.md(x) must land in exactly one phase.
3. **Quality Gates.** Checkable conditions for the phase. Always include the project's baseline gates (tests pass, lint/type checks clean); add what the solution implies (coverage, performance budget, security, a11y, migration safety). Carry the requirements' NFR fit criteria (NFR-N) and Evaluation & Validation gates (EV-N) into the phases that deliver the behavior they measure — cite the ID.
4. **Validation Scenarios.** Generate holdout scenarios that test the phase's observable behavior through external interfaces. Assign priorities: P0 (critical path), P1 (important), P2 (edge case). Every phase needs at least one P0. Write each using templates/scenario.md to specDirectory/scenarios/phase-N/{name}.md and link it under the phase's `#### Validation Scenarios` heading.

Present each phase for review with a coverage map (which requirement ACs and which solution behaviors it covers) and any gaps. Ask the user to choose between *Approve phase*, *Edit features/EARS/gates*, *Edit scenarios*, or *Add missing scenarios*.

CRITICAL: The user must approve the phases and their validation scenarios before the plan is finalized — they are the baseline `/implement` is held to.

### 4. Set Satisfaction Threshold

Set the `threshold` in plan.mdx frontmatter (default 0.90) — the ratio of validation scenarios that must pass for the spec to be considered done. Present the default and let the user adjust.

### 5. Validate

Read reference/validation.md and run checks:

Coverage checks:
- Every acceptance criterion in requirements.md(x) maps to exactly one phase feature.
- Every NFR fit criterion (NFR-N) and evaluation gate (EV-N) in requirements.md(x) appears as a quality gate on at least one phase.
- Every `[ASSUMPTION]` carried from requirements.md(x) or solution.md(x) is either resolved or flagged on the phases that depend on it.
- Every EARS criterion in solution.md(x) maps to exactly one phase's Behavior section, and every phase has at least one EARS behavior statement traceable to solution.md(x).
- The `## Traceability` index lists every requirement AC and every solution EARS criterion exactly once, and each row's phase matches where that criterion actually appears in the phases below.
- Every phase has at least one P0 validation scenario.
- Every validation scenario tests observable behavior (no internal implementation references).

Format checks:
- plan.mdx frontmatter has `threshold`; the `## Phases` section follows the parseable spine (each phase carries Features / Behavior — EARS / Quality Gates / Validation Scenarios).
- Scenarios follow the template and live under scenarios/phase-N/.

### 6. Present Status

Read reference/output-format.md and present the status report. Ask the user to choose between *Finalize*, *Revisit a phase*, *Re-slice phases*, or *Adjust threshold*.
