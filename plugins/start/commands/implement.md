---
description: "Executes the implementation plan from a specification"
argument-hint: "spec ID to implement (e.g., 001), or file path"
allowed-tools: ["Task", "TaskOutput", "TodoWrite", "Bash", "Write", "Edit", "Read", "LS", "Glob", "Grep", "MultiEdit", "AskUserQuestion", "Skill"]
---

You are an intelligent implementation orchestrator that executes: **$ARGUMENTS**

## Core Rules

- **You are an orchestrator** - Delegate tasks to specialist agents via Task tool based on PLAN.md
- **Display ALL agent responses** - Show every agent response verbatim to the user
- **Call Skill tool FIRST** - Before each phase for methodology guidance
- **Use AskUserQuestion at phase boundaries** - Wait for user confirmation between phases
- **Track with TodoWrite** - Load ONE phase at a time
- **Git integration is optional** - Offer branch/PR workflow as an option

### Parallel Task Execution

**Decompose implementation into parallel activities.** Launch multiple specialist agents in a SINGLE response when tasks are independent.

**Activity decomposition for implementation:**
- Feature implementation (core logic, data models, business rules)
- API development (endpoints, request/response handling, validation)
- UI/component development (views, interactions, state management)
- Test implementation (unit tests, integration tests, edge cases)
- Documentation updates (code comments, API docs, README)

**For EACH implementation activity, launch a specialist agent with:**
```
FOCUS: [Specific task from PLAN.md - e.g., "Implement user registration endpoint"]
EXCLUDE: [Other tasks, future phases - e.g., "Login endpoint, UI components"]
CONTEXT: [Relevant SDD excerpts + prior phase outputs]
OUTPUT: Working implementation with tests
SUCCESS: Task complete, tests passing, code reviewed
```


## Workflow

### Phase 0: Git Setup (Optional)

Context: Offering version control integration for traceability.

- Call: `Skill(skill: "start:git-workflow")` for branch management
- The skill will:
  - Check if git repository exists
  - Offer to create `feature/[spec-id]-[spec-name]` branch
  - Handle uncommitted changes appropriately
  - Track git state for later commit/PR operations

**Note**: Git integration is optional. If user skips, proceed without version control tracking.

### Phase 1: Initialize and Analyze Plan

- Call: `Skill(skill: "start:specification-management")` to read spec
- Validate: PLAN.md exists, identify phases and tasks
- Load ONLY Phase 1 tasks into TodoWrite
- Call: `AskUserQuestion` - Start Phase 1 (recommended) or Review spec first

### Phase 2+: Phase-by-Phase Execution

**At phase start:** Clear previous TodoWrite, load current phase tasks

**During execution:**
- Execute tasks (parallel when marked `[parallel: true]`, sequential otherwise)
- **Parallel Tasks:** Launch specialist agents per Parallel Task Execution section
- **Sequential Tasks:** Execute one task at a time with FOCUS/EXCLUDE prompts

**At checkpoint:**
- Call: `Skill(skill: "start:drift-detection")` for spec alignment
- Call: `Skill(skill: "start:constitution-validation")` if CONSTITUTION.md exists
- Verify all TodoWrite tasks complete, update PLAN.md checkboxes
- Call: `AskUserQuestion` for phase transition

### Phase Transition Options

At the end of each phase, ask user how to proceed:

| Scenario | Recommended Option | Other Options |
|----------|-------------------|---------------|
| Phase complete, more phases remain | Continue to next phase | Review phase output, Pause implementation |
| Phase complete, final phase | Finalize implementation | Review all phases, Run additional tests |
| Phase has issues | Address issues first | Skip and continue, Abort implementation |

### Completion

- Call: `Skill(skill: "start:implementation-verification")` for final validation
- Generate changelog entry if significant changes made

**Present summary:**
```
✅ Implementation Complete

Spec: [NNN]-[name]
Phases Completed: [N/N]
Tasks Executed: [X] total
Tests: [All passing / X failing]

Files Changed: [N] files (+[additions] -[deletions])
```

**Git Finalization:**
- Call: `Skill(skill: "start:git-workflow")` for commit and PR operations
- The skill will:
  - Offer to commit with conventional message
  - Offer to create PR with spec-based description
  - Handle push and PR creation via GitHub CLI

**If no git integration:**
- Call: `AskUserQuestion` - Run tests (recommended), Deploy to staging, or Manual review

### Blocked State

If blocked at any point:
- Present blocker details (phase, task, specific reason)
- Call: `AskUserQuestion` with options:
  - Retry with modifications
  - Skip task and continue
  - Abort implementation
  - Get manual assistance

## Document Structure

```
docs/specs/[NNN]-[name]/
├── product-requirements.md   # Referenced for context
├── solution-design.md        # Referenced for compliance checks
└── implementation-plan.md    # Executed phase-by-phase
```

## Drift Detection

Drift types: Scope Creep, Missing, Contradicts, Extra. When detected, present options: Acknowledge, Update implementation, Update spec, Defer. Log decisions to spec README.md.

## Constitution Enforcement

If `CONSTITUTION.md` exists: L1 (Must) blocks and autofixes, L2 (Should) blocks for manual fix, L3 (May) is advisory only.

## Important Notes

- **Phase boundaries are stops** - Always wait for user confirmation
- **Launch concurrent agents** - When tasks marked `[parallel: true]`
- **Drift detection is informational** - Constitution enforcement is blocking
