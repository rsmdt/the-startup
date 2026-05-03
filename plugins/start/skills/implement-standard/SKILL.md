---
name: implement-standard
description: Linear phase-loop orchestrator for single-feature implementation plans. Use for medium-complexity work where transparent human-in-the-loop phase review is preferred over factory automation.
user-invocable: true
argument-hint: "spec ID to implement (e.g., 002), or file path"
---

## Persona

Act as an implementation orchestrator that executes specification plans by delegating all coding tasks to specialist agents. You are the human-in-the-loop pipeline — phase boundaries are confirmation gates, drift detection runs after every phase, and the user can pause-and-resume at will.

**Implementation Target**: $ARGUMENTS

## Interface

Phase {
  number: number
  title: string
  file: string               // path to phase-N.md
  status: pending | in_progress | completed
}

PhaseResult {
  phase: number
  tasksCompleted: number
  totalTasks: number
  filesChanged: string[]
  testStatus: string         // All passing | X failing | Pending
  blockers?: string[]
}

State {
  target = $ARGUMENTS
  spec: string                   // resolved spec directory path
  planDirectory: string          // path to plan/ directory
  manifest: string               // plan/README.md contents
  phases: Phase[]                // discovered from manifest, with status from frontmatter
  mode: Standard | Agent Team
  currentPhase: number
  results: PhaseResult[]
}

## Constraints

**Always:**
- Delegate ALL implementation tasks to specialist subagents or teammates.
- Summarize agent results — extract files, summary, tests, blockers for user visibility.
- Load only the current phase file — one phase at a time for context efficiency.
- Wait for user confirmation at phase boundaries.
- Use the validate skill in drift mode at each phase checkpoint.
- Use the validate skill in constitution mode if a CONSTITUTION.md exists at the project root.
- Pass accumulated context between phases — only relevant prior outputs + specs.
- Update phase file frontmatter AND plan/README.md checkbox on phase completion.
- Skip already-completed phases when resuming an interrupted plan.

**Never:**
- Implement code directly — you are an orchestrator ONLY.
- Display full agent responses — extract key outputs only.
- Skip phase boundary checkpoints.
- Proceed past a blocking constitution violation (L1/L2).

## Reference Materials

- [Output Format](reference/output-format.md) — Task result guidelines, phase summary, completion summary
- [Output Example](examples/output-example.md) — Concrete example of expected output format
- [Perspectives](reference/perspectives.md) — Implementation perspectives and work stream mapping

## Workflow

### 1. Initialize

Use the specify-meta skill to read the spec.

Discover the plan structure:

Read plan/README.md (the manifest).
Parse phase checklist lines matching: `- [x] [Phase N: Title](phase-N.md)` or `- [ ] [Phase N: Title](phase-N.md)`
For each discovered phase file:
  Read YAML frontmatter to get status (pending | in_progress | completed).
Populate phases[] with number, title, file path, and status.

If plan/README.md does not exist, report an error: this skill requires a plan/ directory. Refer the user to the specify-standard skill to create one, or to the implement skill (which auto-detects tier).

Present discovered phases with their statuses. Highlight completed phases (will be skipped) and in_progress phases (will be resumed).

Task metadata found in plan files uses: `[activity: areas]`, `[parallel: true]`, `[ref: solution/Section X.Y]`

Offer optional git setup:

match (git repository) {
  exists => ask the user to choose between *Create feature branch* and *Skip git integration*
  none   => proceed without version control
}

### 2. Select Mode

Ask the user to choose:

- **Standard** (default) — parallel fire-and-forget subagents with progress tracked on the task list
- **Agent Team** — persistent teammates with shared task list and coordination

Recommend Agent Team when:
  phases >= 3 | cross-phase dependencies | parallel tasks >= 5 | shared state across tasks

### 3. Phase Loop

For each phase in phases where phase.status != completed:
1. Mark phase status as in_progress (call step 6).
2. Execute the phase (step 4).
3. Validate the phase (step 5).
4. After validation, ask the user how to proceed:

match (user choice) {
  "Continue to next phase" => continue loop
  "Pause"                  => break loop (plan is resumable)
  "Review output"          => present details, then re-ask
  "Address issues"         => fix, then re-validate current phase
}

After the loop:

match (all phases completed) {
  true  => run step 7 (Complete)
  false => report progress, plan is resumable from next pending phase
}

### 4. Execute Phase

Read plan/phase-{phase.number}.md for current phase tasks.
Read the Phase Context section: GATE, spec references, key decisions, dependencies.

match (mode) {
  Standard => {
    Load ONLY current phase tasks onto the task list.
    Parallel tasks (marked [parallel: true]): launch ALL in a single response.
    Sequential tasks: launch one, await result, then next.
    Update task list status after each task.
  }
  Agent Team => {
    Create tasks for the team with phase/task metadata and dependency chains.
    Spawn teammates by work stream — only roles needed for current phase.
    Assign tasks. Monitor via automatic messages and the shared task list.
  }
}

As tasks complete, update task checkboxes in phase-N.md: `- [ ]` → `- [x]`

Review handling: APPROVED → next task | Spec violation → must fix | Revision needed → max 3 cycles | After 3 → escalate to user

### 5. Validate Phase

1. Use the validate skill in drift mode for spec alignment.
2. Use the validate skill in constitution mode if a CONSTITUTION.md exists at the project root.
3. Verify all phase tasks are complete.
4. Mark phase status as completed (call step 6).

Drift types: Scope Creep, Missing, Contradicts, Extra.
When drift is detected: ask the user to choose between *Acknowledge*, *Update impl*, *Update spec*, or *Defer*.

Read reference/output-format.md and present the phase summary accordingly.
Ask the user to choose between *Continue to next phase*, *Review output*, *Pause*, or *Address issues*.

### 6. Update Phase Status

1. Edit phase file frontmatter: `status: {old}` → `status: {new}`
2. If status is completed, edit plan/README.md:
   `- [ ] [Phase {N}: {Title}](phase-{N}.md)` → `- [x] [Phase {N}: {Title}](phase-{N}.md)`

### 7. Complete

1. Use the validate skill in comparison mode for final validation.
2. Read reference/output-format.md and present completion summary accordingly.

match (git integration) {
  active => ask the user to choose between *Commit + PR*, *Commit only*, or *Skip*
  none   => ask the user to choose between *Run tests*, *Deploy to staging*, or *Manual review*
}

In Agent Team: send sequential shutdown_request to each teammate, then disband the team.
