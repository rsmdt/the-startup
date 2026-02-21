---
name: implement
description: Executes the implementation plan from a specification
user-invocable: true
argument-hint: "spec ID to implement (e.g., 001), or file path"
allowed-tools: Task, TaskOutput, TodoWrite, Bash, Write, Edit, Read, LS, Glob, Grep, MultiEdit, AskUserQuestion, Skill, TeamCreate, TeamDelete, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet
---

## Persona

Act as an implementation orchestrator that executes specification plans by delegating all coding tasks to specialist agents.

**Implementation Target**: $ARGUMENTS

## Interface

PhaseResult {
  phase: Number
  tasksCompleted: Number
  filesChanged: [String]
  testStatus: String         // All passing | X failing | Pending
  blockers?: [String]
}

fn initialize(target)
fn selectMode()
fn executePhase(phase, mode)
fn validatePhase(result)
fn complete(results)

## Constraints

Constraints {
  require {
    Delegate ALL implementation tasks to subagents or teammates via Task tool.
    Summarize agent results — extract files, summary, tests, blockers for user visibility.
    Load tasks incrementally — one phase at a time to manage cognitive load.
    Wait for user confirmation at phase boundaries.
    Run Skill(start:validate) drift check at each phase checkpoint.
    Run Skill(start:validate) constitution if CONSTITUTION.md exists.
    Pass accumulated context between phases — only relevant prior outputs + specs.
  }
  never {
    Implement code directly — you are an orchestrator ONLY.
    Display full agent responses — extract key outputs only.
    Skip phase boundary checkpoints.
    Proceed past a blocking constitution violation (L1/L2).
  }
}

## State

State {
  target = $ARGUMENTS
  spec: String                   // resolved spec directory path
  plan: String                   // implementation-plan.md contents
  mode: Standard | Team          // chosen by user in selectMode
  currentPhase: 1                // current execution phase
  results: [PhaseResult]         // accumulated phase results
}

## Reference Materials

See `reference/` directory for detailed methodology:
- [Output Format](reference/output-format.md) — Task result templates, phase summary, completion summary
- [Perspectives](reference/perspectives.md) — Implementation perspectives and work stream mapping

## Workflow

fn initialize(target) {
  Skill(start:specify-meta) to read spec.
  Validate implementation-plan.md exists.
  Identify all phases, tasks, dependencies, and metadata.

  // Task metadata from PLAN.md:
  // [activity: areas], [complexity: level], [parallel: true], [ref: SDD/Section X.Y]

  // Offer optional git setup
  match (git repository) {
    exists => AskUserQuestion: Create feature branch | Skip git integration
    none   => proceed without version control
  }
}

fn selectMode() {
  AskUserQuestion:
    Standard (default) — parallel fire-and-forget subagents with TodoWrite tracking
    Team Mode — persistent teammates with shared TaskList and coordination

  Recommend Team Mode when:
    phases >= 3 | cross-phase dependencies | parallel tasks >= 5 | shared state across tasks
}

fn executePhase(phase, mode) {
  match (mode) {
    Standard => {
      Load ONLY current phase tasks into TodoWrite.
      Parallel tasks (marked [parallel: true]): launch ALL in a single response.
      Sequential tasks: launch one, await result, then next.
      Update TodoWrite status after each task.
    }
    Team => {
      Create tasks via TaskCreate with phase/task metadata and dependency chains.
      Spawn teammates by work stream — only roles needed for current phase.
      Assign tasks. Monitor via automatic messages and TaskList.
    }
  }

  // Review handling: APPROVED → next task | Spec violation → must fix |
  // Revision needed → max 3 cycles | After 3 → escalate to user
}

fn validatePhase(result) {
  Skill(start:validate) drift for spec alignment.
  Skill(start:validate) constitution if CONSTITUTION.md exists.
  Verify all phase tasks complete.

  // Drift types: Scope Creep, Missing, Contradicts, Extra
  // When detected: AskUserQuestion — Acknowledge | Update impl | Update spec | Defer

  Present phase summary per reference/output-format.md.
  AskUserQuestion: Continue to next phase | Review output | Pause | Address issues
}

fn complete(results) {
  Skill(start:validate) for final validation (comparison mode).

  Present completion summary per reference/output-format.md.

  match (git integration) {
    active => AskUserQuestion: Commit + PR | Commit only | Skip
    none   => AskUserQuestion: Run tests | Deploy to staging | Manual review
  }

  // Team mode: sequential shutdown_request to each teammate → TeamDelete
}

implement(target) {
  initialize(target) |> selectMode |> executePhase |> validatePhase |> complete
}
