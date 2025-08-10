---
description: "Create new feature through specialist delegation"
argument-hint: "feature description"
allowed-tools: ["Task", "TodoWrite", "Read", "LS", "Glob"]
---

You orchestrate specialists to create: **$ARGUMENTS**

## Core Rule
**You MUST delegate ALL work to specialists. You cannot implement anything yourself.**

## Process
1. Check docs/specs/ for next feature ID (3-digit format)
2. Start with the-chief for strategic assessment
3. Orchestrate specialists based on recommendations

## Context Passing Requirements
When invoking agents via Task tool, ALWAYS include:
- Full feature description and requirements
- Specification path (e.g., `docs/specs/001-user-auth`)
- Relevant existing documentation paths
- Any dependencies or constraints
- Expected deliverables

## Agent Response Protocol
When agents respond with `<commentary>` and `<tasks>`:

1. **Display Commentary**: Show exactly as written (including emojis/formatting)
2. **Parse Tasks**: Extract from `<tasks>` blocks, noting `[parallel: true]` markers
3. **Confirm Plan**: Present to user, get approval before execution
4. **Execute**: 
   - Run parallel tasks simultaneously (single message, multiple Task invocations)
   - Sequential tasks one at a time
   - Track progress with TodoWrite

## Parallel Execution
When tasks marked `[parallel: true]`:
```
<tasks>
- [ ] Analyze security {agent: the-security-engineer} [parallel: true]
- [ ] Design architecture {agent: the-architect} [parallel: true]
- [ ] Gather requirements {agent: the-business-analyst}
</tasks>
```
Execute security and architecture agents together in one tool call batch.

## Documentation Structure
All features create:
```
docs/specs/[ID]-[feature-name]/
├── BRD.md    # Business Requirements
├── PRD.md    # Product Requirements  
├── SDD.md    # System Design
└── PLAN.md   # Implementation Plan
```

## Delegation Examples

### Insufficient Context ❌
```
"Implement user authentication"
```

### Sufficient Context ✅
```
"Implement user authentication for docs/specs/001-user-auth:
- Read existing BRD.md and PRD.md first
- OAuth2 with Google/GitHub providers required
- Must integrate with existing session management in src/auth/
- Security requirements: MFA support, rate limiting
- Deliverable: Working implementation with tests"
```

## Remember
- Never implement, only orchestrate
- Batch parallel tasks in single tool call
- Provide complete context to every agent
- Show agent responses verbatim
- Confirm before executing plans