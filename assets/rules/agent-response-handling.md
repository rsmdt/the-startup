# Agent Response Handling Protocol

## Commentary Display

When receiving agent responses with `<commentary>` blocks:

1. **Display the ENTIRE commentary block verbatim**
2. Show `---` separator after commentary
3. Display the rest of the response

### Example:
```
Response from the-developer:

<commentary>
💻 **The Developer**: *cracks knuckles*

Time to implement this authentication system...
</commentary>

---

[rest of response continues here]
```

## Task Extraction

When agents return `<tasks>` blocks:

1. **Extract all tasks** from the block
2. **Present tasks to user** for confirmation:
   ```
   📋 Recommended next tasks:
   1. [Task 1 description]
   2. [Task 2 description]
   
   Add these to the todo list? [Y/n]
   ```
3. **Add approved tasks** to TodoWrite
4. **Mark as in_progress** before execution
5. **Mark as completed** immediately after completion

## Parallel Execution

When multiple agents run in parallel:

- **Display EACH response separately** with clear headers
- **Never merge or summarize** responses
- **Extract tasks from ALL** responses
- **Maintain clear agent attribution**

### Example:
```
=== Response from the-architect ===
[Full response with commentary]

=== Response from the-business-analyst ===
[Full response with commentary]
```

## Error Handling

If agent returns error or is blocked:

```
⚠️ Agent blocked: [agent-name]
Reason: [error or block reason]

Options:
a) Retry with revised context
b) Skip this task
c) Reassign to different agent
d) Cancel operation

Your choice [a/b/c/d]: _
```

## Important Rules

- **NEVER skip or summarize** agent responses
- **ALWAYS show commentary blocks** exactly as written
- **PRESERVE agent personality** and style
- **MAINTAIN task continuity** through proper TodoWrite usage