# Shared Context Management Instructions

## Overview

This document defines the context management system used by all agents in the `/develop` command orchestration. All agents should follow these instructions for consistent behavior.

## Context Storage Structure

```
.the-startup/
└── [sessionId]/
    ├── main.jsonl          # Main orchestrator context log
    ├── [agentId].jsonl     # Individual agent instance contexts
    └── ...
```

## JSONL Format

All context files use a simple append-only JSONL format with two roles:

```json
{"role": "user", "content": "The task or question from the orchestrator"}
{"role": "assistant", "content": "The agent's response"}
```

### Example Entry
```json
{"role": "user", "content": "Analyze requirements for user authentication system"}
{"role": "assistant", "content": "<commentary>...</commentary>\n\n## Requirements Analysis\n..."}
```

## Agent Context Management Protocol

### 1. Receiving Context

When invoked, you will receive:
- `sessionId`: The current session identifier
- `agentId`: Your unique instance identifier for this context
- Your task/prompt from the orchestrator

### 2. Context File Location

Your context file is located at:
```
.the-startup/[sessionId]/[agentId].jsonl
```

### 3. Context Handling Process

1. **Check if context file exists**:
   - If exists: Read all previous interactions to understand context
   - If not exists: This is your first invocation for this task

2. **Process the task** with awareness of any previous context

3. **Append to context file**:
   - Add the user's prompt as: `{"role": "user", "content": "[the prompt you received]"}`
   - Add your response as: `{"role": "assistant", "content": "[your complete response]"}`

### 4. Important Notes

- **Never modify existing entries** - only append new ones
- **Include your complete response** in the assistant content, including any `<commentary>` blocks
- **Each invocation** with the same agentId means you should continue from previous work
- **Different agentId** means fresh context, even if you're the same agent type

## Example Implementation

```python
# Pseudo-code for agent context handling
context_file = f".the-startup/{sessionId}/{agentId}.jsonl"

# Read existing context if available
previous_context = []
if file_exists(context_file):
    previous_context = read_jsonl(context_file)
    
# Process with context awareness
response = process_task(prompt, previous_context)

# Append new interaction
append_jsonl(context_file, [
    {"role": "user", "content": prompt},
    {"role": "assistant", "content": response}
])
```

## Agent Instance Scenarios

### Scenario 1: Fresh Task (New agentId)
You receive agentId "a1b2c3" for the first time:
- No context file exists
- Start fresh analysis
- Create new context file

### Scenario 2: Continuation (Same agentId)
You receive agentId "a1b2c3" again:
- Context file exists with previous work
- Read and understand previous interactions
- Continue/refine based on new instructions

### Scenario 3: Parallel Work (Different agentId)
You receive agentId "x7y8z9" while "a1b2c3" exists:
- This is independent work
- Don't read "a1b2c3" context
- Create new "x7y8z9" context file

## Documentation Path Handling

When you receive a documentation path in your prompt (e.g., "Documentation path: docs/specs/001-user-auth/"), this indicates where to create your output documents, not where to store context.

- Context always goes to: `.the-startup/[sessionId]/[agentId].jsonl`
- Documentation outputs go to: `[provided documentation path]/[your-document].md`

## Key Principles

1. **Isolation**: Each agentId represents an isolated context stream
2. **Continuity**: Same agentId means building on previous work
3. **Persistence**: All interactions are preserved for debugging and audit
4. **Simplicity**: Just "user" and "assistant" roles, append-only
5. **Independence**: Agents manage their own context files