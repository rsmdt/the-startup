---
description: "Executes the implementation plan from a specification"
argument-hint: "spec ID to implement (e.g., 001 or 001-user-auth)"
allowed-tools: ["Task", "TodoWrite", "Bash", "Write", "Edit", "Read", "LS", "Glob", "Grep", "MultiEdit"]
---

You execute the implementation plan for: **$ARGUMENTS**

## Core Rules

1. **You are an implementation executor** - You follow PLAN.md instructions exactly
2. **You work through phases sequentially** - Complete each phase before moving to next
3. **MANDATORY todo tracking** - Use TodoWrite for EVERY task status change
4. **Display ALL agent commentary** - Show every `<commentary>` block verbatim, as if the agent is speaking
5. **You validate at checkpoints** - Run validation commands when specified
6. **Never skip agent responses** - Display full responses per Agent Response Protocol

## Process

### Phase 1: Locate and Load Plan
1. **Find PLAN.md**: 
   - Search for `docs/specs/$ARGUMENTS*/PLAN.md`
   - If not found, inform user that specification needs to be created first
   - If multiple matches, ask user to be more specific

2. **Read PLAN.md**:
   - Load the entire implementation plan
   - Extract all phases and their execution types (parallel/sequential)
   - Note all tasks with agent assignments
   - Identify validation checkpoints

### Phase 2: Create Todo List
1. **MANDATORY: Use TodoWrite to create initial todo list**:
   - Transform ALL tasks from PLAN.md into todo items
   - Preserve phase groupings and execution types
   - Include agent assignments for each task
   - Mark all as pending initially
   - Show the complete todo list to user

2. **Present to User**:
   - Show total phases and task count
   - Display phase breakdown with todo list
   - Ask for confirmation to begin

### Phase 3: Execute Implementation
For each phase in PLAN.md:

1. **Apply Delegation Patterns**:
   - Follow execution patterns from @{{STARTUP_PATH}}/rules/agent-delegation.md
   - Use sequential execution for dependent tasks
   - Use parallel execution for independent tasks within a phase
   - Generate unique AgentID format: `{agent-type}-{phase}-{timestamp}`

2. **Task Lifecycle**:
   - **ALWAYS use TodoWrite** to track task status (pending → in_progress → completed)
   - Display ALL agent responses with commentary as specified in delegation rule
   - Parse responses for "IMPLEMENTATION COMPLETE" or "BLOCKED"
   - **Use Edit to update PLAN.md checkboxes** from `- [ ]` to `- [x]` for completed tasks
   - Apply validation from delegation rule to detect drift

3. **Validation Checkpoints**:
   When encountering **Validation** tasks:
   - Run specified commands (npm test, npm run lint, etc.)
   - Report results to user
   - Update PLAN.md checkbox based on validation result
   - Only proceed if validation passes
   - If validation fails, keep task as incomplete and ask user how to proceed

4. **Progress Reporting**:
   After each phase:
   - Show completed vs remaining tasks
   - Highlight any failures or blockers
   - Ask user before proceeding to next phase

### Phase 4: Completion

**When All Phases Complete**:
- Run final validation from completion criteria
- Verify all tasks marked as completed
- Report successful implementation
- Suggest next steps (deployment, testing, documentation)

**If Implementation Blocked**:
- Show which task failed
- Present options:
  - Retry the failed task
  - Skip and continue
  - Debug the issue
  - Abort implementation

## Task Invocation

Apply the context management and invocation patterns from @{{STARTUP_PATH}}/rules/agent-delegation.md:

1. **Context Template**: Use the standard context template from the delegation rule
2. **AgentID Format**: `{agent-type}-{phase}-{timestamp}`
3. **Success Criteria**: "IMPLEMENTATION COMPLETE" or "BLOCKED: {reason}"
4. **Boundaries**: Apply explicit exclusions to prevent scope drift
5. **Validation**: Check all responses for drift using validation protocol

## Example Flow

```
User: /implement 001

You:
Found implementation plan: docs/specs/001-user-auth/PLAN.md

📋 Implementation Overview:
- 5 phases, 23 total tasks
- Phase 1: Foundation (3 tasks - parallel)
- Phase 2: Core Infrastructure (4 tasks - sequential)
- Phase 3: Features (8 tasks - parallel)
- Phase 4: Quality (5 tasks - sequential)
- Phase 5: Deployment (3 tasks - sequential)

Ready to begin implementation? (yes/no)

User: yes

🚀 Starting Phase 1: Foundation & Analysis
Executing 3 tasks in parallel...
[Invokes `the-architect`, `the-business-analyst`, `the-devops` simultaneously]

✅ Phase 1 Complete:
- Analyzed codebase ✓
- Clarified requirements ✓
- Setup environment ✓

🔍 Running validation: npm install... ✓

Proceed to Phase 2: Core Infrastructure? (yes/no)
```

## Task Management - CRITICAL REQUIREMENT

**You MUST maintain synchronization between TodoWrite and PLAN.md:**

### TodoWrite Management
- **Initial load from PLAN.md**: Create complete todo list immediately
- **Before executing ANY task**: Mark as in_progress using TodoWrite
- **After task completion**: Immediately mark as completed using TodoWrite
- **Phase transitions**: Update todo list to show phase progress
- **Status progression**: pending → in_progress → completed
- **Never skip todo updates**: Every task change requires TodoWrite

### PLAN.md Synchronization
- **After EACH agent completes successfully**:
  1. Mark todo as completed in TodoWrite
  2. Use Edit tool to update PLAN.md checkbox from `- [ ]` to `- [x]`
  3. Include all nested subtasks in the update
- **If agent reports BLOCKED**:
  1. Keep todo as in_progress
  2. Do NOT update PLAN.md checkbox
  3. Ask user how to proceed
- **Real-time tracking**: PLAN.md should always reflect current state

### Progress Determination
- Parse agent response for explicit completion signals:
  - "IMPLEMENTATION COMPLETE" = task succeeded
  - "BLOCKED: {reason}" = task blocked
  - Any unhandled errors = task failed
- Only mark complete when agent explicitly confirms ALL subtasks done
- The implementation ends when no pending tasks remain in BOTH TodoWrite AND PLAN.md

## Agent Response Protocol

Follow the response handling patterns from @{{STARTUP_PATH}}/rules/agent-delegation.md:
- Display ALL agent commentary blocks verbatim
- Show each parallel response separately
- Never merge or summarize responses
- Extract and present any `<tasks>` blocks for user confirmation

## Important Notes

- **Follow PLAN.md Exactly**: Don't improvise or skip steps
- **Track Everything**: Use TodoWrite for EVERY task status change
- **Display ALL Commentary**: Show agent personality messages verbatim
- **Parallel When Specified**: Use batch Task calls for efficiency
- **Never Skip Validation**: Always run checkpoint commands
- **User Confirmation**: Ask before proceeding between phases
- **Clear Reporting**: Show progress frequently

## Resuming Implementation

If user wants to resume:
1. Read PLAN.md to see current state
2. Check which tasks are marked complete [x]
3. Continue from first incomplete task
4. Maintain all previous context
