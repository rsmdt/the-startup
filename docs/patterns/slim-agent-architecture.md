---
description: Structure Claude Code sub-agents to maximize effectiveness while minimizing context window usage by separating agent prompts from skills
status: established
category: agent-design
applies_to: plugins/team/agents
last_updated: 2025-11-30
---

# Slim Agent Architecture Pattern

## Overview

This pattern describes how to structure Claude Code sub-agents to maximize effectiveness while minimizing context window usage. It establishes a clear separation between agent prompts (always loaded) and skills (loaded on-demand via progressive disclosure).

## Problem Statement

When agents contain detailed procedural knowledge that duplicates skill content:
- Context window fills with redundant information
- Updates require changes in multiple places
- Agents become bloated and hard to maintain
- Skills don't provide their intended value

## Solution

Structure agents using a **slim template** that delegates procedural knowledge to skills:

- **Agents** define WHO (role), WHAT (focus areas), and WHEN (deliverables)
- **Skills** provide HOW (procedural knowledge, patterns, checklists, examples)

## Official Documentation Basis

### Source References

- **Sub-agents**: https://code.claude.com/docs/en/sub-agents
- **Skills**: https://code.claude.com/docs/en/skills

### Key Official Guidance

**On description vs body separation:**
> - **Description field**: Signals *when* the subagent should be invoked
> - **Body/System prompt**: Details *how* the subagent operates—its role, capabilities, approach, best practices, and constraints

**On agent design:**
> - "Create subagents with single, clear responsibilities"
> - "Include specific instructions, examples, and constraints in your system prompts"

**On skills:**
> - "Skills implement progressive disclosure to manage context efficiently"
> - "Claude autonomously decides when to use them based on your request and the Skill's description"

## Template Structure

```markdown
---
name: agent-name
description: Clear purpose statement with usage examples that signal WHEN to invoke
skills: skill1, skill2, skill3
model: inherit
---

{1-2 sentence role introduction - WHO you are}

## Focus Areas

{4-6 bullet points of WHAT this agent specializes in}

## Approach

{3-5 high-level methodology steps - HOW at 30,000 feet}
{Explicit skill references: "Leverage {skill-name} for detailed patterns"}

## Deliverables

{4-6 concrete outputs this agent produces}

## Quality Standards

{Non-negotiable quality criteria}
{Include: "Don't create documentation files unless explicitly instructed"}

{1 sentence closing philosophy that guides edge-case decisions}
```

## Section-by-Section Rationale

### YAML Frontmatter

| Field | Purpose | Example |
|-------|---------|---------|
| `name` | Unique identifier | `api-development` |
| `description` | Signals WHEN to invoke (include examples) | Detailed with `<example>` blocks |
| `skills` | Auto-loaded procedural knowledge | `api-design-patterns, documentation-creation` |
| `model` | Model selection | `inherit` (from parent) |

### Role Introduction (1-2 sentences)

Establishes WHO the agent is and their mindset. Sets tone without preamble.

```markdown
You are a pragmatic API architect who designs interfaces developers love to use.
```

### Focus Areas (4-6 bullets)

Defines WHAT the agent does - capabilities and scope. These are outcomes, not procedures.

```markdown
## Focus Areas

- Design clear, consistent API contracts with well-defined schemas
- Generate comprehensive documentation directly from code
- Create interactive testing environments with live examples
```

### Approach (3-5 steps + skill references)

High-level methodology with **explicit skill references**. This is critical for progressive disclosure.

```markdown
## Approach

1. Define use cases before designing endpoints
2. Establish consistent naming conventions and error scenarios
3. Create schemas with validation rules
4. Leverage api-design-patterns skill for REST/GraphQL patterns
5. Leverage documentation-creation skill for guides and SDK examples
```

**Why skill references matter**: They tell the agent WHERE to find detailed procedural knowledge instead of duplicating it in the agent body.

### Deliverables (4-6 outputs)

Concrete artifacts the agent produces. Helps define "done".

```markdown
## Deliverables

1. Complete API specification
2. Request/response schemas with examples
3. Interactive documentation
4. Error catalog with troubleshooting steps
```

### Quality Standards (non-negotiables)

Brief, scannable constraints that guide quality decisions.

```markdown
## Quality Standards

- Include working examples for every endpoint
- Apply security best practices
- Don't create documentation files unless explicitly instructed
```

### Closing Philosophy (1 sentence)

Reinforces character and decision-making lens for edge cases.

```markdown
You approach API development with the mindset that great APIs are intuitive and delightful.
```

## Before/After Comparison

### Before: Bloated Agent (~90 lines)

```markdown
## API Development Methodology

1. **Design Phase:**
   - Define use cases and user journeys before designing endpoints
   - Map resource hierarchies and relationships
   - Create consistent naming conventions across all endpoints
   - Establish error scenarios and edge cases upfront
   - Design for API evolution and future extensibility

2. **Contract Definition:**
   - Define clear request/response schemas with validation rules
   - Apply proper HTTP semantics and status codes for REST
   - Design efficient type systems for GraphQL avoiding N+1 problems
   ...

[50+ more lines of procedural detail]
```

**Problem**: This duplicates what's in the `api-design-patterns` skill (539 lines of comprehensive patterns).

### After: Slim Agent (~47 lines)

```markdown
## Approach

1. Define use cases and user journeys before designing endpoints
2. Establish consistent naming conventions and error scenarios
3. Create request/response schemas with validation rules
4. Generate testable documentation with interactive playgrounds
5. Leverage api-design-patterns skill for REST/GraphQL patterns
6. Leverage documentation-creation skill for guides and SDK examples
```

**Solution**: High-level steps with explicit skill references. Procedural detail lives in skills.

## Metrics

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Average agent size | ~100 lines | ~50 lines | **50% reduction** |
| Total agent content (27 agents) | ~2,745 lines | ~1,396 lines | **49% reduction** |
| Duplication with skills | High | None | **Eliminated** |
| Skill utilization | Underutilized | Explicit references | **Improved** |

## Key Principles

1. **No Duplication**: Agent content should NOT repeat what skills provide
2. **Explicit Skill References**: Agents tell Claude which skills to leverage
3. **Progressive Disclosure**: Skills load on-demand; agent body always loads
4. **Single Source of Truth**: Update patterns in skills, not across 27 agents
5. **Focused Scope**: Each agent excels at one activity

## Anti-Patterns to Avoid

### ❌ Duplicating Skill Content in Agents

```markdown
## Approach

### REST API Patterns
- Use plural nouns for collections
- Use HTTP verbs correctly
- Return appropriate status codes
...
```

This duplicates `api-design-patterns` skill content.

### ❌ Missing Skill References

```markdown
## Approach

1. Design the API
2. Write documentation
3. Test endpoints
```

No guidance on WHERE to find detailed patterns.

### ✅ Correct: High-Level + Skill Reference

```markdown
## Approach

1. Design resource hierarchies and relationships
2. Establish naming conventions and error handling
3. Leverage api-design-patterns skill for REST/GraphQL implementation details
```

## When to Apply This Pattern

- Creating new Claude Code sub-agents
- Refactoring existing verbose agents
- Integrating agents with a skills library
- Reducing context window usage
- Establishing maintainable agent architecture

## Related Patterns

- Skills Library Architecture (see `plugins/team/skills/README.md`)

## References

- Official Sub-agents Documentation: https://code.claude.com/docs/en/sub-agents
- Official Skills Documentation: https://code.claude.com/docs/en/skills
- Team Plugin README: `plugins/team/README.md`
- Skills Library Index: `plugins/team/skills/README.md`
