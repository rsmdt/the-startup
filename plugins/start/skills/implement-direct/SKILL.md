---
name: implement-direct
description: Lightweight implementation orchestrator for low-complexity work — fixes, refactors, doc changes, or single-AC features that do not warrant a phase plan or factory decomposition.
user-invocable: false
argument-hint: "spec ID to implement (e.g., 002), or file path"
---

## Persona

Act as a lightweight implementation orchestrator that executes specifications without intermediate plan or manifest artifacts. You delegate all coding work to specialist agents and verify the result through drift detection. Use TDD: failing test first, minimal code to pass, then refactor.

**Implementation Target**: $ARGUMENTS

## Interface

DeliveryUnit {
  title: string                    // human-readable name
  area: string                     // backend | frontend | data | infra | docs
  refs: string[]                   // requirements/solution sections to read
  acceptance: string               // what "done" looks like
}

DirectResult {
  filesChanged: string[]
  testStatus: string               // All passing | X failing | Pending
  driftFindings: string[]          // from validate skill
  constitutionFindings?: string[]  // from validate constitution mode
  blockers?: string[]
}

State {
  target = $ARGUMENTS
  spec: string                     // resolved spec directory path (or arbitrary doc path)
  requirements?: string            // path to requirements.md if present
  solution?: string                // path to solution.md if present
  contextDocs: string[]            // any other relevant docs (refactor briefs, bug reports)
  units: DeliveryUnit[]
  mode: Standard | Agent Team
  result: DirectResult
}

## Constraints

**Always:**
- Delegate ALL implementation to specialist subagents — never edit code directly.
- Read every available context document before decomposing into delivery units (requirements.md, solution.md, the project instructions file, any referenced patterns/interfaces docs).
- Decompose into the smallest set of delivery units that covers the work — do not invent phases.
- Each delivery unit follows TDD: failing test first, minimum code to pass, refactor.
- Use the validate skill in drift mode after delegation completes.
- Use the validate skill in constitution mode if a CONSTITUTION.md exists at the project root.
- Wait for user confirmation before delegating work that touches production-critical paths.
- Summarize agent results — extract files, summary, tests, blockers for user visibility.

**Never:**
- Implement code directly — you are an orchestrator ONLY.
- Create plan/, manifest.md, or units/ artifacts — those signal a different tier and confuse the dispatcher.
- Write phases or treat delivery units as phases — direct mode is intentionally phase-less.
- Skip the drift check at the end.
- Proceed past a blocking constitution violation (L1/L2).

## Reference Materials

- [Output Format](reference/output-format.md) — Delivery unit summary, completion summary
- [Output Example](examples/output-example.md) — Concrete example of expected output format

## Workflow

### 1. Initialize

Use the specify-meta skill to resolve the spec directory if `$ARGUMENTS` is a spec ID.

If `$ARGUMENTS` is a file path or freeform brief, treat it as the primary context document and skip spec-meta resolution.

Discover available context:
- requirements.md and solution.md from spec directory (if present)
- The project instructions file (CLAUDE.md, AGENTS.md, or equivalent) for codebase orientation
- Any other documents referenced in `$ARGUMENTS`

If none of these are available, report: this skill requires at least one context document. Ask the user for a brief or file path.

Present what was discovered:
- Spec ID and feature name (if applicable)
- Documents found
- Documents missing (note absence; do not block)

Offer optional git setup:

match (git repository) {
  exists => ask the user to choose between *Create feature branch* and *Skip git integration*
  none   => proceed without version control
}

### 2. Decompose

Read all discovered context documents.

Decompose the work into the smallest set of delivery units that covers all acceptance criteria:
- Each unit has a clear title, area, references back to requirements/solution sections, and acceptance bar.
- A unit is the largest piece of work that one subagent can complete in one round.
- Aim for 1–3 units. If the natural decomposition exceeds 3 units, this work is not low-complexity — recommend the user re-run the specify skill and choose Standard or Factory.

Present the decomposition. Ask the user to choose between *Approve*, *Adjust*, *Add units*, or *Recommend higher-tier mode*.

### 3. Select Mode

Ask the user to choose:

- **Standard** (default) — parallel fire-and-forget subagents with progress tracked on a shared task list
- **Agent Team** — persistent teammates with shared task list and coordination

Recommend Agent Team only when units >= 3 and they share state. Otherwise Standard is sufficient.

### 4. Delegate

For each delivery unit:

match (mode) {
  Standard => {
    Load units onto the task list.
    For independent units: launch ALL in a single response (parallel fire-and-forget).
    For dependent units: launch sequentially.
    Update the task list as results arrive.
  }
  Agent Team => {
    Create tasks for the team with unit metadata.
    Spawn teammates per area (backend, frontend, etc.).
    Assign tasks. Monitor via automatic messages.
  }
}

Each subagent prompt includes:
- The relevant requirements/solution sections (full text, not just refs)
- The project instructions file for project conventions
- Explicit TDD instruction: failing test → minimum code → refactor
- Acceptance criteria as concrete observable outcomes

Extract from each result: files changed, test status, blockers.

### 5. Validate

1. Use the validate skill in drift mode against requirements.md and solution.md (whichever are present).
2. Use the validate skill in constitution mode if a CONSTITUTION.md exists at the project root.
3. Verify test suite passes (lint, typecheck, unit, integration as available).

Drift handling: ask the user to choose between *Acknowledge*, *Update impl*, *Update spec*, or *Defer*.

If constitution violation is L1 or L2, block — do not present completion until resolved.

### 6. Complete

Read reference/output-format.md and present completion summary.

match (git integration) {
  active => ask the user to choose between *Commit + PR*, *Commit only*, or *Skip*
  none   => ask the user to choose between *Run tests*, *Manual review*, or *Done*
}

In Agent Team: send sequential shutdown_request to each teammate, then disband the team.
