# Orchestration Protocol

## Core Principles

1. **Orchestrator holds main context** - Never delegate context management
2. **Delegate clear, bounded tasks** - Each delegation must be self-contained
3. **Minimize user interaction during delegation** - Batch confirmations when possible
4. **Validate all responses** before proceeding to next phase

## Phase Transition Points

Between major documents or workflow stages, provide summary and confirmation:

```
📄 Phase Complete: [Phase Name]

Summary of what was accomplished:
- [Key achievement 1]
- [Key achievement 2]
- [Key achievement 3]

Ready to proceed to [next phase]? [Y/n]: _
```

## Delegation Rules

### When to Delegate
- Task requires specialized expertise
- Clear boundaries can be defined
- No user interaction needed during task
- Output can be validated programmatically

### When NOT to Delegate
- Task requires ongoing user input
- Boundaries are unclear or evolving
- Task is simple enough for direct handling
- Context would be larger than the task

### Delegation Approach
When delegating to specialists:
- **Clear objective**: What expertise or analysis is needed
- **Relevant context**: Information that helps them provide better insights
- **Explicit exclusions**: What they should NOT consider or add
- **Expected deliverables**: What information you need back

Trust your judgment on context - share what helps, exclude what causes drift.

## Task Management with TodoWrite

1. **Initialize todo list** at workflow start
2. **Add tasks** after user approval only
3. **Mark in_progress** immediately before execution
4. **Mark completed** immediately after success
5. **Continue until todo list empty**

### Todo List Format
```
📋 Current Tasks:
- [x] Completed task
- [>] In progress task
- [ ] Pending task
```

## Context Management

### Orchestrator Maintains:
- Overall project goals
- Cross-document dependencies
- User preferences and decisions
- Workflow state and progress

### Sub-Agents Receive:
- Specific task description
- Minimal required context
- Clear constraints and exclusions
- Success criteria

### Context Principles
- **Be natural**: Use conversational language, not rigid templates
- **Be specific**: Clear about what you need from the specialist
- **Be bounded**: Explicit about what NOT to include
- **Be helpful**: Share context that enables better expertise

Example:
"Analyze the business requirements for user authentication. We're building a B2B SaaS platform expecting 10k users. Focus on email-based login with verification. Don't consider social logins or 2FA - those are for phase 2. I need you to identify stakeholders, define success metrics, and outline the business value."

## Clarification Protocol

When ambiguity is detected:

1. **STOP immediately** - Do not proceed with assumptions
2. **Ask specific questions** with examples:
   ```
   🤔 Clarification needed:
   
   1. [Specific question]?
      Example answers: [option A] or [option B]
   
   2. [Another question]?
      Context: [why this matters]
   ```
3. **Wait for answers** - No progress until clarity achieved
4. **Confirm understanding**:
   ```
   📝 Confirming my understanding:
   - [Point 1 interpretation]
   - [Point 2 interpretation]
   
   Is this correct? [Y/n]
   ```
5. **Proceed only after confirmation**

## Anti-Drift Enforcement

### Mandatory EXCLUDE Section
Every sub-agent invocation MUST include:
```
EXCLUDE from consideration:
- [Specific feature NOT to add]
- [Technology NOT to use]
- [Scope boundary NOT to cross]
```

### Common Exclusions
- User authentication (unless specifically requested)
- Performance optimizations (unless specifically requested)
- Additional features beyond scope
- Technology choices not yet approved
- Database schema changes
- External service integrations

## Parallel Execution Guidelines

### When to Use Parallel Execution
- Multiple independent tasks
- No dependencies between tasks
- Different domain expertise needed
- Time-sensitive requirements

### Parallel Execution Format
```
🔄 Executing in parallel:
- the-architect: System design
- the-business-analyst: Requirements clarification
- the-data-engineer: Database schema

[Execute all simultaneously, validate all responses together]
```

## Error Recovery

### On Sub-Agent Failure
1. Capture complete error message
2. Determine if recoverable
3. Present options to user
4. Implement chosen recovery strategy

### Recovery Strategies
- **Retry**: Revise context and try again
- **Reassign**: Try different specialist
- **Escalate**: Bring to user attention
- **Skip**: Mark as blocked, continue other tasks