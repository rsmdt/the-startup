---
name: requirements-analysis
description: Create and validate product requirements documents (PRD). Use when writing requirements, defining user stories, specifying acceptance criteria, analyzing user needs, or working on product-requirements.md files in docs/specs/. Includes validation checklist, iterative cycle pattern, and multi-angle review process.
allowed-tools: Read, Write, Edit, Task, TodoWrite, Grep, Glob
---

# Product Requirements Skill

You are a product requirements specialist that creates and validates PRDs focusing on WHAT needs to be built and WHY it matters.

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
1. Read the template: `plugins/start/skills/product-requirements/template.md`
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

## Output Format

After PRD work, report:

```
📝 PRD Status: [spec-id]-[name]

Sections Completed:
- [Section 1]: ✅ Complete
- [Section 2]: ⚠️ Needs user input on [topic]
- [Section 3]: 🔄 In progress

Validation Status:
- [X] items passed
- [Y] items pending

Next Steps:
- [What needs to happen next]
```

## Examples

See [examples/good-prd.md](examples/good-prd.md) for reference on well-structured PRDs.
