# Agent Delegation Rule

## Overview
Guidelines for delegating tasks to sub-agents, with focus on parallel execution patterns, validation, and response handling. Use these patterns when invoking specialist agents to enhance efficiency and maintain quality.

## Core Principles

1. **Clear boundaries** - Specify what's in and out of scope
2. **Minimal context** - Pass only what's needed for the task
3. **Validate responses** - Check for drift before proceeding
4. **Track progress** - Use TodoWrite throughout delegation

## Parallel Execution

Execute tasks in parallel when they are:
- Independent of each other
- Require different expertise
- Can be validated separately
- Would benefit from simultaneous processing

**Implementation:**
1. Mark all tasks as `in_progress` in TodoWrite
2. Generate unique AgentID for each
3. Invoke all agents simultaneously using multiple Task tool calls
4. Validate each response independently
5. Mark tasks as `completed` based on results

**Synchronization:**
- Handle conflicts between responses
- Consolidate findings before proceeding

## Context Management

### What to Include
- Specific task requirements
- Relevant constraints
- Success criteria
- Dependencies (if any)

### What to Exclude
Be explicit about boundaries to prevent scope drift:
- Features not requested
- Optimizations not needed
- Future phase items
- Other agents' responsibilities

### AgentID Format
Use consistent format for tracking: `{agent}-{short-id}`

Example: `the-architect-3xy87q`

## Response Handling

### Commentary Display
Always display agent personality and commentary:
```
<commentary>
[Agent's personality/thoughts]
</commentary>

---

[Response content]
```

### Parallel Responses
Display each response separately - never merge or summarize:
```
=== Response from agent-1 ===
[Full response]

=== Response from agent-2 ===
[Full response]
```

### Task Extraction
When agents return `<tasks>`:
1. Extract all tasks
2. Present for confirmation
3. Add approved tasks to TodoWrite

## Validation & Drift Detection

### Check Every Response
```
🔍 Validating Response from [agent-name]
├─ Scope adherence: [status]
├─ Complexity: [status]
└─ Result: [PASS/DRIFT]
```

### Drift Categories

**Auto-accept (PASS):**
- Error handling improvements
- Input validation
- Security best practices
- Documentation additions

**Ask user (MINOR DRIFT):**
- Helpful additions not requested
- Better patterns suggested
- Extra test coverage

**Requires approval (MAJOR DRIFT):**
- New features
- Database changes
- External integrations
- Performance optimizations

### Handling Drift
For major drift:
```
⚠️ Drift detected: [description]

Options:
a) Accept additions
b) Reject and revise
c) Partial accept

Your choice: _
```

If rejecting, re-invoke with stricter boundaries.

## Error Recovery

### When Agent is Blocked
```
⚠️ Agent blocked: [reason]

Options:
a) Retry with revised context
b) Skip task
c) Reassign to different agent

Your choice: _
```

### Recovery Strategies
- **Retry**: Add clarifications and stricter boundaries
- **Reassign**: Try different specialist
- **Skip**: Mark as blocked, continue other tasks

## TodoWrite Integration

### Lifecycle
- Before delegation: Add task
- Starting: Mark `in_progress`
- After validation: Mark `completed` or keep `in_progress` if blocked
- Update immediately - don't batch updates

### Parallel Tasks
Update all parallel tasks based on their individual results:
```
📋 Parallel execution complete:
- [x] Agent 1: Success
- [x] Agent 2: Success
- [>] Agent 3: Blocked
```

## Phase Transitions

Between major phases:
```
📄 Phase Complete: [Name]

Summary:
- [Key outcome 1]
- [Key outcome 2]

Continue? [Y/n]
```

## Best Practices

**DO:**
- Use parallel execution for independent tasks
- Validate every response
- Display all commentary
- Update TodoWrite immediately
- Be explicit about exclusions

**DON'T:**
- Skip validation
- Merge parallel responses
- Pass unnecessary context
- Allow unchecked drift
- Forget task tracking

## Quick Reference

### Parallel Pattern
```
# Mark all in_progress
# Launch simultaneously
# Validate each
# Update todos based on results
```

### Context Boundaries
```
FOCUS: [what to do]
EXCLUDE: [what not to do]
```

### Validation Flow
```
Response → Validate → Pass/Drift → Action
```

Remember: This rule provides patterns to enhance your natural delegation abilities. Apply them when they add value, using your judgment for the specific situation.
