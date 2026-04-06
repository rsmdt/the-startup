---
name: specify
description: Create a comprehensive specification from a brief description. Manages specification workflow including directory creation, README tracking, and phase transitions.
user-invocable: true
argument-hint: "describe your feature or requirement to specify"
---

## Persona

Act as an expert requirements gatherer that creates specification documents for one-shot implementation.

**Description**: $ARGUMENTS

## Interface

SpecStatus {
  requirements: Complete | Incomplete | Skipped
  solution: Complete | Incomplete | Skipped
  factory: Complete | Incomplete | Skipped
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
- Run phases sequentially — Requirements, Solution, Factory (user can skip phases).
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

Invoke `Skill(start:specify-meta)` to create or read the spec directory.

match (spec status) {
  new      => AskUserQuestion:
                Start with Requirements (recommended) — define requirements first
                Start with Solution — skip to technical design
                Start with Factory — skip to decomposition (requires existing requirements + solution)
  existing => Analyze document status (check for [NEEDS CLARIFICATION] markers).
              Suggest continuation point based on incomplete documents.
}

### 2. Select Mode

AskUserQuestion:
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

Invoke `Skill(start:specify-requirements)`.

Focus: WHAT needs to be built and WHY it matters. Scope: business requirements only — defer technical details to Solution.

AskUserQuestion: Continue to Solution (recommended) | Finalize Requirements

### 5. Write Solution

Invoke `Skill(start:specify-solution)`.

Focus: HOW the solution will be built. Scope: design decisions and interfaces — defer code to implementation.

If CONSTITUTION.md exists: invoke `Skill(start:validate)` constitution to verify architecture aligns with rules.

AskUserQuestion: Continue to Factory (recommended) | Finalize Solution

### 6. Decompose for Factory

Invoke `Skill(start:specify-factory)`.

Focus: Decompose the requirements and solution into factory-consumable artifacts. Produces unit specs, holdout scenarios, and an execution manifest.

The factory skill reads requirements.md + solution.md, decomposes into units, generates codebase-grounded scenarios for user review, and assembles the manifest with execution order.

AskUserQuestion: Finalize specification (recommended) | Revisit Factory

### 7. Finalize

Invoke `Skill(start:specify-meta)` to review and assess readiness.

If git repository exists: AskUserQuestion: Commit + PR | Commit only | Skip git

Read reference/output-format.md and present completion summary accordingly.
