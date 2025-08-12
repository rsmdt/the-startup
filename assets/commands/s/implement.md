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

1. **Sequential Execution**:
   If phase marked as sequential:
   - Execute tasks one at a time in listed order
   - **ALWAYS use TodoWrite to mark as in_progress BEFORE starting**
   - Invoke assigned agent with task description and context
   - **Display agent response with full commentary (see Agent Response Protocol)**
   - **IMMEDIATELY use TodoWrite to mark as completed when done**
   - Stop phase if any task fails

2. **Parallel Execution**:
   If phase marked as parallel:
   - Identify all tasks in the phase
   - **Use TodoWrite to mark ALL as in_progress**
   - Invoke multiple agents simultaneously using batch Task calls
   - Each agent gets: task description, subtasks, spec path, document references
   - Monitor all executions
   - **Display EACH agent's response separately with commentary**
   - Wait for all to complete before proceeding
   - **Use TodoWrite to mark all as completed**

3. **Validation Checkpoints**:
   When encountering **Validation** tasks:
   - Run specified commands (npm test, npm run lint, etc.)
   - Report results to user
   - Only proceed if validation passes
   - If validation fails, ask user how to proceed

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

When invoking specialists from PLAN.md:
1. **Extract Context**:
   - Task description and any indented subtasks
   - Referenced documents (e.g., [→ PRD#auth-requirements])
   - Spec path for accessing BRD/PRD/SDD

2. **Provide Clear Instructions**:
   - Pass complete task with all subtasks
   - Include "Read docs/specs/XXX/[document].md for context"
   - Specify expected deliverables

3. **Handle Results**:
   - Update todo list based on success/failure
   - Collect any error messages
   - Proceed according to execution type

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

**You MUST use the TodoWrite tool throughout the entire workflow:**
- **Initial load from PLAN.md**: Create complete todo list immediately
- **Before executing ANY task**: Mark as in_progress using TodoWrite
- **After task completion**: Immediately mark as completed using TodoWrite
- **Phase transitions**: Update todo list to show phase progress
- **Status progression**: pending → in_progress → completed
- **Never skip todo updates**: Every task change requires TodoWrite
- **Continue until empty**: The implementation ends when no pending tasks remain

## Agent Response Protocol - MANDATORY

**You MUST display EVERY agent response completely, including their personality/commentary:**

1. **Display Commentary First**:
   - Show any `<commentary>...</commentary>` block EXACTLY as written
   - Display it immediately after "Response from [agent-name]:"
   - Include ALL formatting, emojis, personality messages
   - Add `---` separator after commentary

2. **For Parallel Execution**:
   - Display EACH agent's response separately
   - Show all commentaries in full
   - Never merge or summarize responses

3. **Example of CORRECT Display**:
   ```
   Response from the-developer:
   
   💻 **The Developer**: *cracks knuckles*
   
   Time to implement this authentication system. I've built dozens of these...
   
   ---
   
   [rest of response]
   ```

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
