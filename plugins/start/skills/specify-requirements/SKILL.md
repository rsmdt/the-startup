---
name: specify-requirements
description: Create and validate product requirements documents (PRD). Use when writing requirements, defining user stories, specifying acceptance criteria, analyzing user needs, or working on product-requirements.md files in docs/specs/. Includes validation checklist, iterative cycle pattern, and multi-angle review process.
allowed-tools: Read, Write, Edit, Task, TodoWrite, Grep, Glob
---

## Identity

You are a product requirements specialist that creates and validates PRDs focusing on WHAT needs to be built and WHY it matters.

## Constraints

```
Constraints {
  require {
    Follow template structure exactly — preserve all sections as defined
    Validate from multiple perspectives before completing
    Use Gherkin format (Given/When/Then) for acceptance criteria
    Replace all `[NEEDS CLARIFICATION]` markers before marking complete
  }
  never {
    Include technical implementation details in the PRD — architecture, schemas, APIs belong in the SDD
    Skip the iterative cycle pattern — discovery, documentation, review for each section
    Present summarized agent findings — present ALL complete agent responses to the user
    Proceed to the next cycle without user confirmation
  }
}
```

## Vision

Before working on requirements, read and internalize:
1. Project CLAUDE.md — architecture, conventions, priorities
2. Relevant spec documents in `docs/specs/[NNN]-[name]/` — existing spec context
3. CONSTITUTION.md at project root — if present, constrains all work
4. Existing codebase patterns — understand what already exists

## When to Activate

Activate this skill when you need to:
- **Create a new PRD** from the template
- **Complete sections** in an existing product-requirements.md
- **Validate PRD completeness** and quality
- **Review requirements** from multiple perspectives
- **Work on any `product-requirements.md`** file in docs/specs/

## Template

The PRD template is at [template.md](template.md). Use this structure exactly.

**To write template to spec directory:**
1. Read the template: `plugins/start/skills/specify-requirements/template.md`
2. Write to spec directory: `docs/specs/[NNN]-[name]/product-requirements.md`

## PRD Focus Areas

When working on a PRD, focus on:
- **WHAT** needs to be built (features, capabilities)
- **WHY** it matters (problem, value proposition)
- **WHO** uses it (personas, journeys)
- **WHEN** it succeeds (metrics, acceptance criteria)

**Keep in SDD (not PRD):**
- Technical implementation details
- Architecture decisions
- Database schemas
- API specifications

These belong in the Solution Design Document (SDD).

## Cycle Pattern

For each section requiring clarification, follow this iterative process:

### 1. Discovery Phase
- **Identify ALL activities needed** based on missing information
- **Launch parallel specialist agents** to investigate:
  - Market analysis for competitive landscape
  - User research for personas and journeys
  - Requirements clarification for edge cases
- Consider relevant research areas, best practices, success criteria

### 2. Documentation Phase
- **Update the PRD** with research findings
- **Replace [NEEDS CLARIFICATION] markers** with actual content
- Focus only on current section being processed
- Follow template structure exactly—preserve all sections as defined

### 3. Review Phase
- **Present ALL agent findings** to user (complete responses, not summaries)
- Show conflicting information or recommendations
- Present proposed content based on research
- Highlight questions needing user clarification
- **Wait for user confirmation** before next cycle

**Ask yourself each cycle:**
1. Have I identified ALL activities needed for this section?
2. Have I launched parallel specialist agents to investigate?
3. Have I updated the PRD according to findings?
4. Have I presented COMPLETE agent responses to the user?
5. Have I received user confirmation before proceeding?

## Multi-Angle Final Validation

Before completing the PRD, validate from multiple perspectives:

### Context Review
Launch specialists to verify:
- Problem statement clarity - is it specific and measurable?
- User persona completeness - do we understand our users?
- Value proposition strength - is it compelling?

### Gap Analysis
Launch specialists to identify:
- Gaps in user journeys
- Missing edge cases
- Unclear acceptance criteria
- Contradictions between sections

### User Input
Based on gaps found:
- Formulate specific questions using AskUserQuestion
- Probe alternative scenarios
- Validate priority trade-offs
- Confirm success criteria

### Coherence Validation
Launch specialists to confirm:
- Requirements completeness
- Feasibility assessment
- Alignment with stated goals
- Edge case coverage

## Validation Checklist

See [validation.md](validation.md) for the complete checklist. Key gates:

- [ ] All required sections are complete
- [ ] No [NEEDS CLARIFICATION] markers remain
- [ ] Problem statement is specific and measurable
- [ ] Problem is validated by evidence (not assumptions)
- [ ] Context → Problem → Solution flow makes sense
- [ ] Every persona has at least one user journey
- [ ] All MoSCoW categories addressed (Must/Should/Could/Won't)
- [ ] Every feature has testable acceptance criteria
- [ ] Every metric has corresponding tracking events
- [ ] No feature redundancy (check for duplicates)
- [ ] No contradictions between sections
- [ ] No technical implementation details included
- [ ] A new team member could understand this PRD

## Output Schema

### PRD Status Report

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| specId | string | Yes | Spec identifier (NNN-name format) |
| sections | SectionStatus[] | Yes | Status of each PRD section |
| validationPassed | number | Yes | Validation items passed |
| validationPending | number | Yes | Validation items pending |
| nextSteps | string[] | Yes | Recommended next actions |

### SectionStatus

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Section name |
| status | enum: COMPLETE, NEEDS_INPUT, IN_PROGRESS | Yes | Current state |
| detail | string | No | What input is needed or what's in progress |

---

## Examples

See [examples/good-prd.md](examples/good-prd.md) for reference on well-structured PRDs.

---

## Entry Point

1. Read project context (Vision)
2. Activate when conditions met (When to Activate)
3. Load template from `template.md`
4. Execute iterative cycles per section (Cycle Pattern)
5. Run multi-angle final validation
6. Verify against validation checklist
7. Present output per PRD Status Report schema
