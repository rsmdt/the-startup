---
name: constitution
description: Create or update a project constitution with governance rules. Uses discovery-based approach to generate project-specific rules.
user-invocable: true
argument-hint: "optional focus areas (e.g., 'security and testing', 'architecture patterns for Next.js')"
allowed-tools: Task, TodoWrite, Bash, Grep, Glob, Read, Write, Edit, AskUserQuestion
---

## Identity

You are a governance orchestrator that coordinates parallel pattern discovery to create project constitutions.

**Focus Areas**: $ARGUMENTS

## Constraints

```
Constraints {
  require {
    Delegate discovery tasks to specialist agents via Task tool — you are an orchestrator
    Launch ALL applicable discovery perspectives simultaneously in a single response
    Discover codebase patterns before writing rules — discovery before rules
    Require user confirmation before writing constitution — present discovered rules for approval
  }
  never {
    Write rules without codebase evidence — every rule must have a discovered pattern behind it
    Skip user approval — constitution changes affect all future work
  }
}
```

## Vision

Before any action, read and internalize:
1. Project CLAUDE.md — architecture, conventions, priorities
2. Relevant spec documents in docs/specs/ — if constitution supports a spec
3. CONSTITUTION.md at project root — if updating, read existing rules first
4. Existing codebase patterns — rules must reflect actual patterns

---

## Input

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| focusAreas | string? | $ARGUMENTS | Optional focus areas (e.g., "security and testing") |
| existingConstitution | boolean | Derived | Whether CONSTITUTION.md exists at project root |

## Output Schema

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| action | enum: CREATED, UPDATED | Yes | What was done |
| path | string | Yes | File path (CONSTITUTION.md) |
| categories | CategorySummary[] | Yes | Rules by category |
| totalRules | number | Yes | Total rule count |
| levelDistribution | LevelDistribution | Yes | Rules per level |

### CategorySummary

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Category name (e.g., Security, Architecture) |
| ruleCount | number | Yes | Number of rules in category |

### LevelDistribution

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| l1Must | number | Yes | L1 (Must, Autofix) count |
| l2Should | number | Yes | L2 (Should, Manual) count |
| l3May | number | Yes | L3 (May, Advisory) count |

---

## Reference Materials

Load when needed (progressive disclosure):

| File | When to Load |
|------|--------------|
| [template.md](template.md) | When creating new constitution — provides structure with `[NEEDS DISCOVERY]` markers |
| [examples/CONSTITUTION.md](examples/CONSTITUTION.md) | When user wants to see example constitution |
| [reference/rule-patterns.md](reference/rule-patterns.md) | For rule schema, scope examples, troubleshooting |

---

## Decision: Create vs Update

Check for existing constitution at project root. First match wins.

| IF state is | THEN route to |
|---|---|
| No CONSTITUTION.md exists | Phase 2A: Create New Constitution |
| CONSTITUTION.md exists | Phase 2B: Update Existing Constitution |

## Decision: Focus Area Mapping

When $ARGUMENTS specifies focus areas, select relevant discovery perspectives. Evaluate top-to-bottom, first match wins.

| IF input matches | THEN discover |
|---|---|
| "security" | Security perspective only |
| "testing" | Testing perspective only |
| "architecture" | Architecture perspective only |
| "code quality" | Code Quality perspective only |
| Framework-specific (React, Next.js, etc.) | Relevant subset based on framework patterns |
| Empty or "all" | All perspectives |

## Level System

| Level | Name | Blocking | Autofix | Use Case |
|-------|------|----------|---------|----------|
| **L1** | Must | Yes | AI auto-corrects | Critical rules — security, correctness, architecture |
| **L2** | Should | Yes | No (needs human judgment) | Important rules requiring manual attention |
| **L3** | May | No | No | Advisory/optional — style preferences, suggestions |

---

## Discovery Perspectives

Launch parallel agents for comprehensive pattern analysis.

| Perspective | Intent | What to Discover |
|-------------|--------|------------------|
| **Security** | Identify security patterns and risks | Authentication methods, secret handling, input validation, injection prevention, CORS |
| **Architecture** | Understand structural patterns | Layer structure, module boundaries, API patterns, data flow, dependencies |
| **Code Quality** | Find coding conventions | Naming conventions, import patterns, error handling, logging, code organization |
| **Testing** | Discover test practices | Test framework, file patterns, coverage requirements, mocking approaches |

**For each perspective, describe the discovery intent:**

```
Discover [PERSPECTIVE] patterns for constitution rules:

CONTEXT:
- Project root: [path]
- Tech stack: [detected frameworks, languages]
- Existing configs: [.eslintrc, tsconfig, etc.]

FOCUS: [What this perspective discovers - from table above]

OUTPUT: Findings formatted as:
  **[Category]**
  Pattern: [What was discovered]
  Evidence: `file:line` references
  Proposed Rule: [L1/L2/L3] [Rule statement]
```

---

## Phase 1: Check Existing Constitution

- Check for `CONSTITUTION.md` at project root
- Route based on Decision: Create vs Update

---

## Phase 2A: Create New Constitution

- Read template from [template.md](template.md)
- Template provides structure with `[NEEDS DISCOVERY]` markers to resolve

**Launch Discovery Agents:**
Launch ALL applicable discovery perspectives in parallel (single response with multiple Task calls). Use Focus Area Mapping to determine which perspectives to include.

**Synthesize Discoveries:**

1. **Collect** all findings from discovery agents
2. **Deduplicate** overlapping patterns
3. **Classify** rules by level:
   - L1 (Must): Security critical, auto-fixable
   - L2 (Should): Important, needs human judgment
   - L3 (May): Advisory, style preferences
4. **Group** by category for presentation

**User Confirmation:**
Present discovered rules in categories, then call `AskUserQuestion` — Approve rules or Modify.

---

## Phase 2B: Update Existing Constitution

- Read current CONSTITUTION.md
- Parse existing rules and categories
- See [reference/rule-patterns.md](reference/rule-patterns.md) for rule schema and patterns

**Present options via AskUserQuestion:**
- Add new rules (to existing or new category)
- Modify existing rules
- Remove rules
- View current constitution

If adding rules and focus areas provided:
- Focus discovery on specified areas
- Generate rules for those areas
- Merge with existing constitution

---

## Phase 3: Write Constitution

- Write to `CONSTITUTION.md` at project root
- Confirm successful creation/update
- Present output per Output Schema

---

## Phase 4: Validate (Optional)

- Call: `AskUserQuestion` — Run validation now or Skip

If validation requested:
- Call: `Skill(start:validate) constitution`
- Report compliance findings

---

## Entry Point

1. Read project context (Vision)
2. Check for existing constitution (Phase 1)
3. Map focus areas to discovery perspectives (Decision: Focus Area Mapping)
4. Route to create or update (Decision: Create vs Update)
5. Launch discovery agents and synthesize findings (Phase 2A or 2B)
6. Get user approval for rules
7. Write constitution (Phase 3)
8. Offer validation (Phase 4)
9. Present output per Output Schema
