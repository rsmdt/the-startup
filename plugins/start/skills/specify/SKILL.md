---
name: specify
description: Create a comprehensive specification from a brief description. Manages specification workflow including directory creation, README tracking, and phase transitions.
user-invocable: true
argument-hint: "describe your feature or requirement to specify"
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Grep, Read, Write(docs/**), Edit(docs/**), AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
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

fn initialize(target)
fn selectMode()
fn research(mode)
fn writePRD()
fn writeSDD()
fn writePLAN()
fn finalize(status)

## Constraints

Constraints {
  require {
    Delegate research tasks to specialist agents via Task tool.
    Display ALL agent responses to user — complete findings, not summaries.
    Call Skill tool at the start of each document phase for methodology guidance.
    Phases are sequential — PRD, SDD, PLAN (user can skip phases).
    Wait for user confirmation between each document phase.
    Track decisions in specification README via reference/output-format.md.
    Git integration is optional — offer branch/commit as an option.
  }
  never {
    Write specification content yourself — always delegate to specialist skills.
    Proceed to next document phase without user approval.
    Skip decision logging when user makes non-default choices.
  }
}

## State

State {
  target = $ARGUMENTS
  spec: String                   // resolved spec directory path (from specify-meta)
  perspectives = []              // determined by research scope
  mode: Standard | Team          // chosen by user in selectMode
  status: SpecStatus             // tracked across phases
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Perspectives](reference/perspectives.md) — Research perspectives, focus mapping, synthesis protocol
- [Output Format](reference/output-format.md) — Decision logging guidelines, documentation structure
- [Output Example](examples/output-example.md) — Concrete example of expected output format

## Workflow

fn initialize(target) {
  Skill(start:specify-meta) to create or read spec directory.

  match (spec status) {
    new => AskUserQuestion:
      Start with PRD (recommended) — define requirements first
      Start with SDD — skip to technical design
      Start with PLAN — skip to implementation planning
    existing => {
      Analyze document status (check for [NEEDS CLARIFICATION] markers).
      Suggest continuation point based on incomplete documents.
    }
  }
}

fn selectMode() {
  AskUserQuestion:
    Standard (default) — parallel fire-and-forget research agents
    Team Mode — persistent researcher teammates with peer collaboration

  Recommend Team Mode when:
    3+ document phases planned | complex domain | multiple integrations |
    conflicting perspectives likely (security vs performance)
}

fn research(mode) {
  // Launch research per applicable perspectives from reference/perspectives.md
  match (mode) {
    Standard => launch parallel subagents per applicable perspectives
    Team     => create team, spawn one researcher per perspective, assign tasks
  }

  Synthesize per reference/perspectives.md synthesis protocol.
  // Research feeds into all subsequent document phases.
}

fn writePRD() {
  Skill(start:specify-requirements)
  // Focus: WHAT needs to be built and WHY it matters
  // Scope: Business requirements only (defer technical details to SDD)

  AskUserQuestion: Continue to SDD (recommended) | Finalize PRD
}

fn writeSDD() {
  Skill(start:specify-solution)
  // Focus: HOW the solution will be built
  // Scope: Design decisions and interfaces (defer code to implementation)

  // Constitution alignment (if CONSTITUTION.md exists):
  match (CONSTITUTION.md) {
    exists => Skill(start:validate) constitution — verify architecture aligns with rules
  }

  AskUserQuestion: Continue to PLAN (recommended) | Finalize SDD
}

fn writePLAN() {
  Skill(start:specify-plan)
  // Focus: Task sequencing and dependencies
  // Scope: What and in what order (defer duration estimates)

  AskUserQuestion: Finalize specification (recommended) | Revisit PLAN
}

fn finalize(status) {
  Skill(start:specify-meta) to review and assess readiness.

  match (git repository) {
    exists => AskUserQuestion: Commit + PR | Commit only | Skip git
  }

  Present completion summary per reference/output-format.md.
}

specify(target) {
  initialize(target) |> selectMode |> research |> writePRD |> writeSDD |> writePLAN |> finalize
}
