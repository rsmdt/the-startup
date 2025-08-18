# Command Prompt Style Guide

A comprehensive guide for writing concise, effective command prompts for AI orchestration systems.

## Executive Summary

This guide documents best practices for writing command prompts that maximize AI capabilities while minimizing complexity. Based on real-world optimization that achieved 95% reduction in prompt size (1974 → 95 lines) while maintaining full functionality, these principles enable clear, maintainable, and effective AI orchestration.

## Core Principles

### 1. Trust AI Capabilities
**Principle**: AI systems understand their tools and capabilities. Don't teach them what they already know.

**Do:**
- State what needs to be done
- Define success criteria
- Specify constraints

**Don't:**
- Include tool pseudo-code or usage examples
- Explain how to use basic features
- Provide verbose error handling templates

**Example:**
```markdown
# Bad (verbose, untrusting)
Use the Read tool with the following pattern:
- First check if file exists using LS tool
- Then use Read tool: Read(file_path="/path/to/file")
- Handle errors by checking if response contains "error"
- If error, display message: "Failed to read file: [error details]"

# Good (concise, trusting)
Read the specification documents from `docs/specs/$ID/`
```

### 2. Natural Language Over Pseudo-Code
**Principle**: Clear natural language is more maintainable and equally effective as pseudo-code.

**Do:**
- Use simple, declarative sentences
- Describe the logic flow naturally
- Let AI interpret the implementation details

**Don't:**
- Write IF/ELSE blocks in prompts
- Use pseudo-code patterns
- Over-specify control flow

**Example:**
```markdown
# Bad (pseudo-code style)
IF task.has_review_metadata THEN
  extract review_areas FROM metadata
  select_reviewer BASED ON review_areas
  IF reviewer == "architect" THEN
    invoke architect_reviewer
  ELSE IF reviewer == "security" THEN
    invoke security_reviewer
  END IF
END IF

# Good (natural language)
When task has review metadata, analyze the specified areas and select an appropriate reviewer.
```

### 3. Eliminate Redundancy
**Principle**: One clear instruction beats multiple examples or repeated explanations.

**Do:**
- State each requirement once
- Trust AI to generalize from single instructions
- Remove duplicate guidance

**Don't:**
- Provide multiple examples of the same pattern
- Repeat instructions in different sections
- Include "just in case" redundancies

**Example:**
```markdown
# Bad (redundant)
## Review Process
Handle reviews by checking metadata...
### For Architecture Reviews
When reviewing architecture...
### For Security Reviews  
When reviewing security...
### Review Summary
After any review type...

# Good (single clear instruction)
## Review Process
When `[review: areas]` metadata is present, select appropriate reviewer based on specified areas and handle feedback with max 3 revision cycles.
```

### 4. Simplify Metadata
**Principle**: Metadata should be minimal and intuitive. Presence implies action.

**Do:**
- Use presence to indicate boolean states
- Combine related metadata
- Keep formats simple and readable

**Don't:**
- Use verbose boolean flags
- Separate tightly coupled metadata
- Create complex parsing requirements

**Example:**
```markdown
# Bad (verbose metadata)
[agent: architect]
[needs_review: true]
[review_focus: security, performance]
[validation_required: true]

# Good (simplified)
[agent: architect]
[review: security, performance]
```

### 5. Dynamic Over Static Selection
**Principle**: Let the orchestrator choose the best agent/approach based on context.

**Do:**
- Make agent assignment optional
- Allow dynamic selection based on task requirements
- Provide selection criteria, not mappings

**Don't:**
- Hard-code agent assignments
- Create rigid task-to-agent mappings
- Prevent flexible decision-making

**Example:**
```markdown
# Bad (static)
For database tasks: ALWAYS use the-database-engineer
For UI tasks: ALWAYS use the-frontend-developer
For API tasks: ALWAYS use the-backend-developer

# Good (dynamic)
Select the best specialist based on task requirements. Consider using [agent: name] metadata as a hint, but choose based on actual needs.
```

## Structure Patterns

### 1. Frontmatter First
Start with clear metadata that tools can parse:

```yaml
---
description: "Brief description of command purpose"
argument-hint: "what user should provide"
allowed-tools: ["Tool1", "Tool2", "Tool3"]
---
```

### 2. Role Definition
Single sentence establishing the orchestrator's identity:

```markdown
You are an intelligent [type] orchestrator that [primary purpose]: **$ARGUMENTS**
```

### 3. Core Rules Section
Bullet points of non-negotiable behaviors:

```markdown
## Core Rules
- **You are an orchestrator** - Delegate tasks to specialist agents
- **Work sequentially** - Complete each phase before moving to next
- **Track everything** - Use TodoWrite for task management
- **Display commentary** - Show all agent personality blocks
```

### 4. Process Flow
Numbered steps with clear progression:

```markdown
## Process

### Step 1: Context Loading
- Find specification at path
- Read documents
- Display summary

### Step 2: Task Creation
- Parse plan for tasks
- Create todo list
- Confirm with user

### Step 3: Orchestrated Execution
- Delegate tasks to agents
- Handle responses
- Update tracking
```

### 5. Reference External Rules
Link to detailed patterns rather than embedding:

```markdown
## Delegation Guidelines
You MUST FOLLOW patterns from @{{STARTUP_PATH}}/rules/agent-delegation.md for all task delegations.
```

## Language Guidelines

### 1. Use Active Voice
**Do:** "Delegate tasks to specialist agents"  
**Don't:** "Tasks should be delegated to specialist agents"

### 2. Be Specific Without Being Verbose
**Do:** "Mark task as `in_progress` before execution"  
**Don't:** "Before you begin executing any task, you must first update the TodoWrite tool to mark the task status as in_progress"

### 3. Use Consistent Terminology
**Do:** Always use "orchestrator" for main agent, "specialist" for sub-agents  
**Don't:** Mix terms like coordinator, manager, worker, sub-agent

### 4. Prefer Imperatives
**Do:** "Validate responses before proceeding"  
**Don't:** "You should validate responses before you proceed"

### 5. Minimize Modal Verbs
**Do:** "Update TodoWrite immediately"  
**Don't:** "You must always remember to update TodoWrite immediately"

## Metadata Design

### 1. Task Metadata Format
Keep it minimal and semantic:

```markdown
[agent: optional-name]     # Hint for agent selection
[review: areas]            # Presence indicates review needed
[parallel: group-name]     # Indicates parallelizable tasks
```

### 2. Status Indicators
Use simple, visual markers:

```markdown
- [ ] Todo
- [~] In progress  
- [x] Complete
- [!] Blocked
```

### 3. Phase Markers
Clear visual separation:

```markdown
### Step 1: Context Loading
📁 Found spec: 001-user-auth

### Step 2: Task Creation
📊 Creating task list

### Step 3: Execution
🚀 Starting implementation
```

## Orchestration Patterns

### 1. Delegation Pattern
```markdown
For each task:
- Extract agent from metadata or select best fit
- Mark as in_progress
- Delegate with minimal context
- Validate response
- Update status
```

### 2. Parallel Execution
```markdown
Batch parallel tasks in single response when:
- Tasks are independent
- Different expertise needed
- No dependencies exist
```

### 3. Review Cycles
```markdown
When review metadata present:
- Analyze implementation
- Select reviewer based on areas
- Handle feedback with max 3 cycles
- Escalate if not resolved
```

### 4. Validation Checkpoints
```markdown
At phase boundaries:
- Run validation commands
- Only proceed if passing
- Show clear status to user
```

## Anti-Patterns to Avoid

### 1. Tool Pseudo-Code
**Bad:**
```markdown
Use Grep tool like this:
Grep(pattern="TODO", path="./src", type="typescript")
If results.length > 0:
  For each result in results:
    Process(result)
```

**Why it's bad:** AI already knows how to use its tools. This adds complexity without value.

### 2. Excessive Examples
**Bad:**
```markdown
Examples of task delegation:
- Example 1: When delegating to architect...
- Example 2: When delegating to developer...
- Example 3: When delegating to tester...
[10 more examples]
```

**Why it's bad:** AI can generalize from principles. Examples should only clarify ambiguous cases.

### 3. Defensive Over-Specification
**Bad:**
```markdown
IMPORTANT: You MUST ALWAYS check if the file exists before reading.
CRITICAL: NEVER forget to validate responses.
WARNING: Make sure to ALWAYS update the todo list.
REMEMBER: You are an orchestrator, not an implementer.
```

**Why it's bad:** Defensive language and repetition suggest distrust and add noise.

### 4. Nested Conditionals
**Bad:**
```markdown
If task requires review:
  If review type is architecture:
    If architecture is complex:
      Use senior architect
    Else:
      Use regular architect
  Else if review type is security:
    [more nesting]
```

**Why it's bad:** Complex nesting is hard to maintain. Let AI make intelligent decisions.

### 5. Verbose Templates
**Bad:**
```markdown
Display this exact format:
╔════════════════════════════════════╗
║  TASK EXECUTION REPORT             ║
╠════════════════════════════════════╣
║ Status: [STATUS]                   ║
║ Agent: [AGENT]                     ║
║ Time: [TIME]                       ║
╚════════════════════════════════════╝
```

**Why it's bad:** Overly specific formatting adds complexity. Simple, clear output is better.

## Examples: Before and After

### Example 1: Task Delegation

**Before (436 lines):**
```markdown
## Detailed Delegation Process

When you need to delegate a task to a specialist agent, follow these steps:

1. First, identify the task that needs to be delegated
2. Determine which specialist agent is best suited
3. Prepare the context package with all necessary information
4. Use the Task tool with the following format:
   Task(
     task="the task description",
     context={
       "requirements": "...",
       "constraints": "...",
       "dependencies": "..."
     },
     agentId="the-agent-name-uniqueid"
   )
5. Wait for the response
6. Parse the response for success indicators
7. Check if response contains "IMPLEMENTATION COMPLETE"
8. If blocked, check for "BLOCKED:" message
9. Handle errors appropriately
10. Update your todo list

[... continues with examples and edge cases ...]
```

**After (3 lines):**
```markdown
## Delegation
Apply patterns from @{{STARTUP_PATH}}/rules/agent-delegation.md when invoking specialists.
```

### Example 2: Review Process

**Before (287 lines):**
```markdown
## Review and Validation Process

### Identifying Review Requirements
Check each task for review metadata...
[50 lines of detection logic]

### Selecting Reviewers
Based on the review areas specified...
[75 lines of selection criteria]

### Review Execution
When executing a review...
[100 lines of review process]

### Handling Feedback
Process review feedback...
[60 lines of feedback handling]
```

**After (4 lines):**
```markdown
When task has `[review: areas]` metadata:
- Select appropriate reviewer based on areas
- Invoke with full context
- Handle feedback with max 3 revision cycles
```

### Example 3: Error Handling

**Before (198 lines):**
```markdown
## Comprehensive Error Handling

### File Not Found Errors
When a file is not found:
- Display: "ERROR: File not found at [path]"
- Suggest: "Please check the file path"
- Options: retry, skip, abort

### Permission Errors
When permission is denied:
- Display: "ERROR: Permission denied for [operation]"
[... continues with 15 more error types ...]
```

**After (1 line):**
```markdown
Handle errors appropriately and inform the user of issues and recovery options.
```

## Implementation Checklist

When writing or reviewing a command prompt:

- [ ] **Frontmatter**: Clear metadata with description and tools?
- [ ] **Role**: Single-sentence orchestrator definition?
- [ ] **Rules**: Bullet points of core behaviors?
- [ ] **Process**: Clear, numbered steps?
- [ ] **References**: External rules linked, not embedded?
- [ ] **Language**: Active voice, imperatives, consistent terms?
- [ ] **Metadata**: Minimal, semantic, presence-based?
- [ ] **Trust**: No tool pseudo-code or over-specification?
- [ ] **Clarity**: Each instruction stated once?
- [ ] **Flexibility**: Dynamic selection enabled?

## Maintenance Guidelines

### 1. Regular Audits
- Review prompts quarterly for redundancy
- Remove instructions that AI handles well naturally
- Consolidate similar patterns

### 2. Performance Metrics
Track and optimize for:
- Token usage (aim for <500 tokens per prompt)
- Task completion rate
- User intervention frequency
- Error recovery success

### 3. Evolution Strategy
- Start with minimal viable instructions
- Add specificity only when errors occur repeatedly
- Prefer external rule files over embedded instructions
- Document why complexity was added

## Conclusion

Effective command prompts trust AI capabilities, use natural language, eliminate redundancy, and maintain flexibility. By following these principles, we achieve prompts that are:

- **Concise**: 95% smaller while maintaining functionality
- **Maintainable**: Easy to understand and modify
- **Effective**: Clear instructions that AI can execute reliably
- **Flexible**: Adaptable to various contexts and requirements

Remember: The goal is not to control every aspect of AI behavior, but to provide clear objectives and let AI leverage its capabilities to achieve them efficiently.

## Quick Reference Card

```markdown
# Command Prompt Template

---
description: "What this command does"
argument-hint: "what user provides"
allowed-tools: ["Tool1", "Tool2"]
---

You are an intelligent [type] orchestrator that [purpose]: **$ARGUMENTS**

## Core Rules
- **Rule name** - Brief explanation
- **Another rule** - Brief explanation

## Process

### Step 1: [Phase Name]
- Action 1
- Action 2
- Validation

### Step 2: [Phase Name]  
- Delegate using @{{STARTUP_PATH}}/rules/[rule].md
- Handle responses
- Update tracking

## Important Notes
Key reminders without redundancy.
```

*Last Updated: 2025-01-18*
*Based on: /s:implement optimization (1974 → 95 lines)*