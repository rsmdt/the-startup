---
description: Execute the implementation for the provided specification
allowed-tools:
  - LS
  - Glob
  - Grep
  - Read
  - Write
  - MultiEdit
  - Edit
  - Bash
  - WebSearch
  - WebFetch
  - Task
  - TodoWrite
argument-hint: </path/to/spec directory or spec ID>
---

Execute the passed specification for: **$ARGUMENTS**. You orchestrate specification implementation through specialized sub-agents, ensuring systematic execution, proper validation, and clear progress tracking.

## Process

### Phase 1: Initialization
1. **Locate Specification**: Find the PRD using the provided argument
   - If ID (e.g., "001"): Search `docs/specs/001-*/`
   - If path: Use directly
   - If name: Search `docs/specs/` for matching title

2. **Prepare Context**: 
   - Read CLAUDE.md files (global and project)
   - Use `the-project-manager` sub-agent to analyze all files in the specification directory:
     1. PRD.md - Product Specification Document
     2. SDD.md - Solution Design Documentation
     3. PLAN.md - Implementation Plan

### Phase 2: Implementation Loop
Execute implementation phases as identified by the task manager:

1. **Task Planning**: 
   Use `the-project-manager` to:
   - Identify next tasks to execute
   - Determine which tasks can run in parallel
   - Mark identified tasks as `[~]` in PRD before execution

2. **Parallel Implementation**:
   For each group of independent tasks, identify the most-suitable sub-agent and spawn multiple instances simultaneously, each with:
   - Specific and detailed task assignment for clear execution
   - Relevant PRD, SDD, and PLAN context (patterns, specs, examples, todo's)
   - Clear scope boundaries
   
   While implementers run, provide progress updates in chat:
   ```
   🔄 Starting parallel execution:
   - Task 1: Component structure
   - Task 2: Data fetching logic
   - Task 3: Error handling
   ```

3. **Task Completion**:
   As implementers return results:
   - Task manager updates PLAN: `[~]` → `[x]` for success
   - Task manager updates PLAN: `[~]` → `[BLOCKED: reason]` for blockers
   - You only report status in chat

4. **Validation**:
   After all implementers in phase complete, use the most-suitable sub-agents to:
   - Run all validation commands
   - Interpret results
   - Suggest fixes for issues

5. **Progress Update**:
   Task manager provides comprehensive status and determines next phase

**You MUST CONTINUE this loop until ALL TASKS complete or blockers are encountered.**

### Phase 3: Completion Handling

#### If All Tasks Complete:
- use `the-project-manager` sub-agent to update PLAN with completion timestamp
- Provide final success report

#### If Tasks Remain Incomplete:
Enter interactive completion mode:
1. `the-project-manager` sub-agent generates comprehensive status report
2. Present blockers and options to user
3. Based on guidance, continue implementation
4. Repeat until resolution

## Orchestration Guidelines

### Parallel Execution
When task manager identifies independent tasks:
- Launch multiple implementers in single Task invocation
- Each implementer gets focused context
- Consolidate results before validation
- Update all task states together

### Error Handling
When implementers report blockers:
- Collect all blocker information
- Run validator to check current state
- Present consolidated issues to user
- Await guidance before proceeding

### Progress Tracking
Throughout execution:
- Task manager maintains source of truth in PRD
- Regular status updates in chat
- Clear visibility of completed/blocked/remaining tasks

## Example Execution Flow

```
1. Task Manager: "Phase 1 has 3 independent tasks ready"
   - Updates PRD: marks all 3 tasks as [~]
   
2. Main: Spawns 3 implementers for parallel execution
   - Reports in chat: "🔄 Starting 3 parallel tasks..."
   
3. Implementers complete:
   - Implementer 1: "TASK_COMPLETE: Created component structure"
   - Implementer 2: "TASK_COMPLETE: Added data fetching logic"
   - Implementer 3: "TASK_BLOCKED: Missing API endpoint specification"
   
4. Task Manager: Updates PRD based on results
   - Task 1: [~] → [x]
   - Task 2: [~] → [x]  
   - Task 3: [~] → [BLOCKED: Missing API endpoint]
   
5. Validator: Runs validation on completed work
   - "VALIDATION_FAILED: 2 linting errors in component file"
   
6. Main: Reports consolidated status
   - "✅ 2/3 tasks complete, 1 blocked, validation errors to fix"
   
7. Present options to user for resolution
```

## Key Principles

- **Separation of Concerns**: Each agent has a focused role
- **Parallel When Possible**: Maximize efficiency within phases
- **Clear State Tracking**: PRD is single source of truth
- **Fail Fast**: Report blockers immediately
- **Interactive Resolution**: Work with user to resolve issues
- **Quality Validation**: Never skip validation steps

## Resuming Partial Implementation

If returning to incomplete PRD:
1. Task manager reads current state
2. Validator checks code state
3. Continue from last incomplete phase
4. Maintain previous context and decisions

Begin by locating the PRD and initializing the implementation process.
