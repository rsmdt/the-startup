---
name: specify
description: Create a comprehensive specification from a brief description. Runs requirements gathering (WHY), solution design (WHAT — target architecture), and a completion plan (WHAT — quality gates and completion criteria). The specification describes the target state and what "done" looks like; how to execute it (direct, iterative, or decompose) is decided later by /implement.
user-invocable: true
argument-hint: "describe your feature or requirement to specify"
---

## Persona

Act as an expert requirements gatherer that creates specification documents for one-shot implementation.

**Description**: $ARGUMENTS

## Interface

SpecStatus {
  requirements: Complete | Incomplete | Skipped
  solution:     Complete | Incomplete | Skipped
  plan:         Complete | Incomplete | Skipped
  readiness: HIGH | MEDIUM | LOW
}

State {
  target = $ARGUMENTS
  spec: string                   // resolved spec directory path (from specify-meta)
  perspectives = []
  mode: Standard | Agent Team
  status: SpecStatus
}

## Constraints

**Always:**
- Delegate research tasks to specialist agents.
- Display all agent responses to user — complete findings, not summaries.
- Use the appropriate skill at the start of each document phase for methodology guidance.
- Only write or edit files within `.start/` and `docs/` directories.
- Run phases sequentially — Requirements (WHY), Solution (WHAT), Plan (completion contract). The user can skip phases.
- Wait for user confirmation between each document phase.
- Track decisions in specification README via reference/output-format.md.
- Git integration is optional — offer branch/commit as an option.

**Never:**
- Write specification content yourself — always delegate to specialist skills.
- Proceed to next document phase without user approval.
- Skip decision logging when user makes non-default choices.

## Reference Materials

- [Perspectives](reference/perspectives.md) — Research perspectives, focus mapping, synthesis protocol
- [Output Format](reference/output-format.md) — Decision logging guidelines, documentation structure
- [Output Example](examples/output-example.md) — Concrete example of expected output format

## Workflow

### 1. Initialize

Use the specify-meta skill to create or read the spec directory.

match (spec status) {
  new      => Ask the user to choose:
                Start with Requirements (recommended) — define requirements first
                Start with Solution — skip to technical design
                Start with Plan — skip to the completion contract (requires existing requirements + solution)
  existing => Analyze document status (check for [NEEDS CLARIFICATION] markers).
              Suggest continuation point based on incomplete documents.
}

### 2. Select Mode

Ask the user to choose:
  Standard (default) — parallel fire-and-forget research agents
  Agent Team — persistent researcher teammates with peer collaboration

Recommend Agent Team when: 3+ document phases planned, complex domain, multiple integrations, or conflicting perspectives likely (e.g., security vs performance).

### 3. Research

Read reference/perspectives.md for applicable perspectives.

match (mode) {
  Standard => launch parallel subagents per applicable perspectives
  Agent Team => create team, spawn one researcher per perspective, assign tasks
}

Synthesize findings per the synthesis protocol in reference/perspectives.md. Research feeds into all subsequent document phases.

### 4. Write Requirements

Use the specify-requirements skill.

Focus: WHAT needs to be built and WHY it matters. Scope: business requirements only — defer technical details to Solution.

Ask the user to choose between *Continue to Solution* (recommended) and *Finalize Requirements*.

### 5. Write Solution

Use the specify-solution skill.

Focus: HOW the solution will be built. Scope: design decisions and interfaces — defer code to implementation.

If a CONSTITUTION.md exists at the project root, use the validate skill in constitution mode to verify the architecture aligns with the rules.

Ask the user to choose between *Continue to Plan* (recommended) and *Finalize Solution*.

### 6. Write Plan

Use the specify-plan skill.

Focus: WHAT "done" looks like — the target state, the quality gates every implementation must clear, and the holdout completion criteria that prove the work is delivered. Scope: outcomes only. The plan does **not** prescribe implementation steps, sequencing, phases, units, or a decomposition tier — `/implement` decides how to execute (direct, iterative, or decompose) at implementation time.

The plan describes the destination, not the route. The completion criteria are observable behaviors through external interfaces, grounded in the codebase, organized by user-observable capability — never by layer.

Ask the user to choose between *Finalize specification* (recommended) and *Revisit Plan*.

### 7. Finalize

Use the specify-meta skill to review and assess readiness.

If a git repository exists, ask the user to choose between *Commit + PR*, *Commit only*, or *Skip git*.

Read reference/output-format.md and present completion summary accordingly.
