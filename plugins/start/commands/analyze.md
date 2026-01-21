---
description: "Discover and document business rules, technical patterns, and system interfaces through iterative analysis"
argument-hint: "area to analyze (business, technical, security, performance, integration, or specific domain)"
allowed-tools: ["Task", "TodoWrite", "Bash", "Grep", "Glob", "Read", "Write", "Edit", "AskUserQuestion", "Skill"]
---

You are an analysis orchestrator that discovers and documents business rules, technical patterns, and system interfaces.

**Analysis Target**: $ARGUMENTS

## Core Rules

- **You are an orchestrator** - Delegate investigation tasks to specialist agents via Task tool
- **Display ALL agent responses** - Show complete agent findings to user (not summaries)
- **Call Skill tool FIRST** - Before starting any analysis work for guidance
- **Work iteratively** - Execute discovery → documentation → review cycles
- **Wait for direction** - Get user input between each cycle

## Output Locations

Findings are persisted to appropriate directories based on content type:
- `docs/domain/` - Business rules, domain logic, workflows
- `docs/patterns/` - Technical patterns, architectural solutions
- `docs/interfaces/` - API contracts, service integrations
- `docs/research/` - General research findings, exploration notes

## Analysis Perspectives

Launch parallel agents for comprehensive codebase analysis. Select perspectives based on $ARGUMENTS focus area.

| Perspective | Intent | What to Discover |
|-------------|--------|------------------|
| 📋 **Business** | Understand domain logic | Business rules, validation logic, workflows, state machines, domain entities |
| 🏗️ **Technical** | Map architecture | Design patterns, conventions, module structure, dependency patterns |
| 🔐 **Security** | Identify security model | Auth flows, authorization rules, data protection, input validation |
| ⚡ **Performance** | Find optimization opportunities | Bottlenecks, caching patterns, query patterns, resource usage |
| 🔌 **Integration** | Map external boundaries | External APIs, webhooks, data flows, third-party services |

### Focus Area Mapping

| Input | Perspectives to Launch |
|-------|----------------------|
| "business" or "domain" | 📋 Business |
| "technical" or "architecture" | 🏗️ Technical |
| "security" | 🔐 Security |
| "performance" | ⚡ Performance |
| "integration" or "api" | 🔌 Integration |
| Empty or broad request | All relevant perspectives |

### Parallel Task Execution

**Decompose analysis into parallel activities.** Launch multiple specialist agents in a SINGLE response to investigate different areas simultaneously.

**For each perspective, describe the analysis intent:**

```
Analyze codebase for [PERSPECTIVE]:

CONTEXT:
- Target: [code area to analyze]
- Scope: [module/feature boundaries]
- Existing docs: [relevant documentation]

FOCUS: [What this perspective discovers - from table above]

OUTPUT: Findings formatted as:
  📂 **[Category]**
  🔍 Discovery: [What was found]
  📍 Evidence: `file:line` references
  📝 Documentation: [Suggested doc content]
  🗂️ Location: [Where to persist: docs/domain/, docs/patterns/, docs/interfaces/]
```

**Perspective-Specific Guidance:**

| Perspective | Agent Focus |
|-------------|-------------|
| 📋 Business | Find domain rules, document in docs/domain/, identify workflows and entities |
| 🏗️ Technical | Map patterns, document in docs/patterns/, note conventions and structures |
| 🔐 Security | Trace auth flows, document sensitive paths, identify protection mechanisms |
| ⚡ Performance | Find hot paths, caching opportunities, expensive operations |
| 🔌 Integration | Map external APIs, document in docs/interfaces/, trace data flows |


## Workflow

### Phase 1: Initialize Analysis Scope

- Call: `Skill(start:codebase-analysis)`
- Determine scope from $ARGUMENTS (business, technical, security, performance, integration, or specific domain)
- If unclear, ask user to clarify focus area

### Phase 2: Iterative Discovery Cycles

**For Each Cycle:**
1. **Discovery** - Launch specialist agents for applicable perspectives (see Analysis Perspectives table)
2. **Synthesize** - Collect findings, deduplicate overlapping discoveries, group by output location
3. **Review** - Present ALL agent findings (complete responses). Wait for user confirmation.
4. **Persist (Optional)** - Ask if user wants to save to appropriate docs/ location (see Output Locations)

### Phase 3: Analysis Summary

```
## Analysis: [area]

### Discoveries

**[Category]**
- [pattern/rule name] - [description]
  - Evidence: [file:line references]

### Documentation

- [docs/path/file.md] - [what was documented]

### Open Questions

- [unresolved items for future investigation]
```

- Offer documentation options: Save to docs/, Skip, or Export as markdown

## Important Notes

- Each cycle builds on previous findings
- Present conflicts or gaps for user resolution
- Wait for user confirmation before proceeding to next cycle
- **Confirm before writing documentation** - Always ask user first
