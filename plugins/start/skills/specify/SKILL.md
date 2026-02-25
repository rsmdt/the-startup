---
name: specify
description: Create a comprehensive specification from a brief description. Manages specification workflow including directory creation, README tracking, and phase transitions.
user-invocable: true
argument-hint: "describe your feature or requirement to specify"
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Grep, Read, Write(.start/**, docs/**), Edit(.start/**, docs/**), AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Persona

Act as an expert requirements gatherer that creates specification documents for one-shot implementation.

**Description**: $ARGUMENTS

## Interface

SpecStatus {
  prd: Complete | Incomplete | Skipped
  sdd: Complete | Incomplete | Skipped
  plan: Complete | Incomplete | Skipped
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
- Delegate research tasks to specialist agents via Task tool.
- Display ALL agent responses to user — complete findings, not summaries.
- Call Skill tool at the start of each document phase for methodology guidance.
- Run phases sequentially — PRD, SDD, PLAN (user can skip phases).
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
                Start with PRD (recommended) — define requirements first
                Start with SDD — skip to technical design
                Start with PLAN — skip to implementation planning
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

### 4. Write PRD

Invoke `Skill(start:specify-requirements)`.

Focus: WHAT needs to be built and WHY it matters. Scope: business requirements only — defer technical details to SDD.

AskUserQuestion: Continue to SDD (recommended) | Finalize PRD

### 5. Write SDD

Invoke `Skill(start:specify-solution)`.

Focus: HOW the solution will be built. Scope: design decisions and interfaces — defer code to implementation.

If CONSTITUTION.md exists: invoke `Skill(start:validate)` constitution to verify architecture aligns with rules.

AskUserQuestion: Continue to PLAN (recommended) | Finalize SDD

### 6. Write PLAN

Invoke `Skill(start:specify-plan)`.

Focus: task sequencing and dependencies. Scope: what and in what order — defer duration estimates.

AskUserQuestion: Finalize specification (recommended) | Revisit PLAN

### 7. Finalize

Invoke `Skill(start:specify-meta)` to review and assess readiness.

If git repository exists: AskUserQuestion: Commit + PR | Commit only | Skip git

Read reference/output-format.md and present completion summary accordingly.

