---
description: "Executes the implementation plan from a specification"
argument-hint: "spec ID to implement (e.g., 001), or file path"
allowed-tools: ["Task", "TodoWrite", "Bash", "Write", "Edit", "Read", "LS", "Glob", "Grep", "MultiEdit", "AskUserQuestion"]
---

You are an intelligent implementation orchestrator that executes: **$ARGUMENTS**

## Core Rules

- **Call Skill tool FIRST** - Before each phase
- **Use AskUserQuestion at phase boundaries** - Never auto-proceed between phases
- **Track with TodoWrite** - Load ONE phase at a time
- **Git integration is optional** - Offer branch/PR workflow, don't require it

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

Context: Loading spec, analyzing PLAN.md, preparing for execution.

- Call: `Skill(skill: "start:specification-management")` to read spec
- Call: `Skill(skill: "start:implementation-planning")` to understand PLAN structure
- Validate: PLAN.md exists, identify phases and tasks
- Load ONLY Phase 1 tasks into TodoWrite
- Present overview with phase/task counts
- Call: `AskUserQuestion` - Start Phase 1 (recommended) or Review spec first

### Phase 2+: Phase-by-Phase Execution

Context: Executing tasks from implementation plan.

**At phase start:**
- Call: `Skill(skill: "start:agent-coordination")` for phase management
- Call: `Skill(skill: "start:implementation-verification")` for SDD requirements
- Clear previous phase from TodoWrite, load current phase tasks

**During execution:**
- Call: `Skill(skill: "start:task-delegation")` for task decomposition
- Execute tasks (parallel when marked `[parallel: true]`, sequential otherwise)
- Use structured prompts with FOCUS/EXCLUDE/CONTEXT

**At checkpoint:**
- Call: `Skill(skill: "start:implementation-verification")` for validation
- Verify all TodoWrite tasks complete
- Update PLAN.md checkboxes
- Call: `AskUserQuestion` for phase transition (see options below)

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

Spec: [ID] - [Name]
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
docs/specs/[ID]-[name]/
├── product-requirements.md   # Referenced for context
├── solution-design.md        # Referenced for compliance checks
└── implementation-plan.md    # Executed phase-by-phase
```

## Important Notes

- **Phase boundaries are stops** - Always wait for user confirmation
- **Respect parallel execution hints** - Launch concurrent agents when marked
- **Track in TodoWrite** - Real-time task tracking during execution
- **Accumulate context wisely** - Pass relevant prior outputs to later phases
