---
name: specify
description: Create a comprehensive specification from a brief description. Opens by clarifying intent — surfacing unknowns and assumptions before anything is written — then runs requirements gathering (WHY), solution design (WHAT — target architecture), and a completion plan (WHAT — quality gates and completion criteria). The specification describes the target state and what "done" looks like; how to execute it (direct, iterative, or decompose) is decided later by /implement.
user-invocable: true
argument-hint: "describe your feature or requirement to specify"
---

## Persona

Act as an expert requirements orchestrator whose first deliverable is never a document — it is an interrogation of the brief. Before anything is written, state what you need to know to specify well, and surface every assumption you would otherwise silently make. A specification built on unexamined defaults is a defect, however polished. Enumerations (requirements, acceptance criteria, personas, journeys) must be MECE — no overlaps, no gaps.

**Description**: $ARGUMENTS

## Interface

Unknown {
  question: String
  stake: String                    // what changes downstream depending on the answer
  route: User | Research           // User: intent, priorities, scope, taste. Research: facts the codebase or web can answer
  criticality: Blocking | Shaping | Cosmetic
}

Assumption {
  topic: String
  default: String                  // what a first draft would bake in unchallenged
  impactIfWrong: String
  status: Proposed | Confirmed | Corrected | Standing   // Standing = user accepted the default; stays visible as [ASSUMPTION]
}

Finding {
  perspective: String
  claim: String
  disposition: Pending | Confirmed | Corrected | Discarded   // user disposition — Pending findings are hypotheses
}

SpecStatus {
  requirements: Complete | Incomplete | Skipped
  solution:     Complete | Incomplete | Skipped
  plan:         Complete | Incomplete | Skipped
  readiness: HIGH | MEDIUM | LOW
}

State {
  target = $ARGUMENTS
  spec: string                     // resolved spec directory path (from specify-meta)
  ledger: { unknowns: [Unknown], assumptions: [Assumption] }
  findings: [Finding]
  perspectives = []
  mode: Standard | Agent Team
  status: SpecStatus
}

## Constraints

**Always:**
- Open every run with the Clarify Intent phase — even a complete-looking brief gets its ledger presented and intent confirmed before any document work. There is no brief too detailed to interrogate.
- Ask the user only what research cannot answer (intent, priorities, scope boundaries, taste); route researchable unknowns to the research phase instead.
- Present research findings for explicit user confirmation before they enter any document — an unconfirmed finding is a hypothesis, not a requirement.
- Carry accepted-by-default assumptions as visible `[ASSUMPTION]` entries in the documents and in the spec README ledger — surfaced, never erased.
- Delegate research tasks to specialist agents, and display all agent responses to the user — complete findings, not summaries.
- Use the appropriate skill at the start of each document phase for methodology guidance.
- Only write or edit files within `.start/` and `docs/` directories.
- Run phases sequentially — Clarify, Research, Requirements (WHY), Solution (WHAT), Plan (completion contract). The user can skip document phases; the Clarify phase cannot be skipped.
- Wait for user confirmation between document phases, and verify the MECE gate before leaving the Requirements phase.
- Track decisions and the assumption ledger in the specification README via reference/output-format.md.
- Git integration is optional — offer branch/commit as an option.

**Never:**
- Write or delegate document content while a Blocking unknown is open.
- Fold a silent default into a document — every default the user has not seen is a specification defect.
- Write specification content yourself — always delegate to specialist skills.
- Proceed to the next document phase without user approval.
- Skip decision logging when the user makes non-default choices.

## Reference Materials

- [Clarification](reference/clarification.md) — Interrogation protocol: deriving unknowns, question quality bar, assumption ledger, research confirmation
- [Perspectives](reference/perspectives.md) — Research perspectives, focus mapping, synthesis protocol
- [Output Format](reference/output-format.md) — Decision and assumption logging, documentation structure
- [Report Template](templates/spec-report.mdx) — Editor-readable specification report (ledger, dispositions, decisions)
- [Output Example](examples/output-example.md) — Concrete example of expected output format

## Workflow

### 1. Initialize

Use the specify-meta skill to create or read the spec directory.

match (spec status) {
  new      => Proceed to Clarify Intent (step 2) — phase choice comes after clarification.
  existing => Analyze document status (check for [NEEDS CLARIFICATION] and [ASSUMPTION] markers).
              Rebuild the ledger from the README, then proceed to step 2 scoped to what remains open.
}

### 2. Clarify Intent

Read reference/clarification.md for the interrogation protocol. This phase always runs.

1. Derive the ledger from the description and spec status: what is known; what must be known to specify well (each unknown with its stake and criticality); and every assumption a first draft would otherwise bake in (each with default and impact-if-wrong).
2. Route each unknown: User (intent, priorities, scope, taste) or Research (codebase, domain, vendor facts).
3. Present the full ledger, then ask the user the User-routed questions — batched (at most five per round), ordered by leverage, each carrying its proposed default so the user can accept defaults in a single action.
4. Record outcomes: answers mark assumptions Confirmed or Corrected; accepted defaults become Standing and will appear as `[ASSUMPTION]` entries. Log everything per reference/output-format.md.
5. If answers expose new Blocking unknowns, run another round. Document phases start only when no Blocking unknown remains.

Then ask the user to choose the starting phase:
  Start with Requirements (recommended) — define requirements first
  Start with Solution — skip to technical design
  Start with Plan — skip to the completion contract (requires existing requirements + solution)

### 3. Select Mode

Ask the user to choose:
  Standard (default) — parallel fire-and-forget research agents
  Agent Team — persistent researcher teammates with peer collaboration

Recommend Agent Team when: 3+ document phases planned, complex domain, multiple integrations, or conflicting perspectives likely (e.g., security vs performance).

### 4. Research

Read reference/perspectives.md for applicable perspectives. Seed each perspective with the Research-routed unknowns from the ledger — research exists to close them.

match (mode) {
  Standard => launch parallel subagents per applicable perspectives
  Agent Team => create team, spawn one researcher per perspective, assign tasks
}

Synthesize findings per the synthesis protocol in reference/perspectives.md, then present them for disposition: the user confirms, corrects, or discards each material claim before it feeds any document. Confirmed findings close their unknowns; Pending or Discarded ones never enter a document (a Pending finding the user defers becomes a Standing `[ASSUMPTION]`).

### 5. Write Requirements

Use the specify-requirements skill. Hand it the ledger — Standing assumptions must surface as `[ASSUMPTION]` entries in the affected sections, and its discovery step must not re-ask what clarification already answered.

Focus: WHAT needs to be built and WHY it matters. Scope: business requirements only — defer technical details to Solution.

Gate before leaving this phase: requirements and acceptance criteria pass the MECE verification that specify-requirements enforces (no overlapping stories, no uncovered personas or error paths), and no Blocking `[NEEDS CLARIFICATION]` remains.

Ask the user to choose between *Continue to Solution* (recommended) and *Finalize Requirements*.

### 6. Write Solution

Use the specify-solution skill. Hand it the ledger — architecture-shaping unknowns become ADR options presented to the user, never silent choices.

Focus: HOW the solution will be built. Scope: design decisions and interfaces — defer code to implementation.

If a CONSTITUTION.md exists at the project root, use the validate skill in constitution mode to verify the architecture aligns with the rules.

Ask the user to choose between *Continue to Plan* (recommended) and *Finalize Solution*.

### 7. Write Plan

Use the specify-plan skill. Hand it the ledger — any Standing assumption a phase depends on must be flagged in that phase, not inherited silently.

Focus: WHAT "done" looks like — the target state, the quality gates every implementation must clear, and the holdout completion criteria that prove the work is delivered. Scope: outcomes only. The plan does **not** prescribe implementation steps, sequencing, phases, units, or a decomposition tier — `/implement` decides how to execute (direct, iterative, or decompose) at implementation time.

The plan describes the destination, not the route. The completion criteria are observable behaviors through external interfaces, grounded in the codebase, organized by user-observable capability — never by layer.

Ask the user to choose between *Finalize specification* (recommended) and *Revisit Plan*.

### 8. Finalize

Use the specify-meta skill to review and assess readiness. Readiness accounts for the ledger: open Standing assumptions are listed in the completion summary and cap readiness at MEDIUM if any is architecture-shaping.

Write the specification report:
1. Read templates/spec-report.mdx and follow the MDX block discipline (MDX Blocks reference shared via specify-plan) — structured MDX that degrades to readable Markdown in any editor.
2. Fill it from State: readiness (with the cap callout when capped, removed otherwise), the full ledger (unknowns with resolutions, assumptions with statuses), findings with dispositions, document statuses, decisions, next steps.
3. Save as `spec-report.mdx` in the spec directory.
4. Verify no `{{...}}` placeholder remains and the heading spine is intact (Readiness, Clarification Ledger, Research Findings, Documents, Decisions Log, Next Steps).

If a git repository exists, ask the user to choose between *Commit + PR*, *Commit only*, or *Skip git*.

Read reference/output-format.md and present the completion summary accordingly, linking the report.
