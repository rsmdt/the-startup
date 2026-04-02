---
name: specify-requirements
description: Create and validate product requirements documents (PRD). Use when writing requirements, defining user stories, specifying acceptance criteria, analyzing user needs, or working on requirements.md files in .start/specs/. Includes validation checklist, iterative cycle pattern, and multi-angle review process.
---

## Persona

Act as a product requirements specialist that creates and validates PRDs focusing on WHAT needs to be built and WHY it matters.

**Spec Target**: $ARGUMENTS

## Interface

PRDSection {
  name: string
  status: Complete | NeedsInput | InProgress
  topic?: string       // what needs clarification, if NeedsInput
}

State {
  specId = ""
  sections: PRDSection[]
  clarificationMarkers: number
}

## PRD Focus Areas

When discovering and documenting, address four dimensions:
- **WHAT** needs to be built — features, capabilities
- **WHY** it matters — problem, value proposition
- **WHO** uses it — personas, journeys
- **WHEN** it succeeds — metrics, acceptance criteria

**Out of scope:** Technical implementation, architecture, database schemas, API specifications — those belong in SDD.

## MECE Principle

All structured enumerations in the PRD must be **Mutually Exclusive, Collectively Exhaustive** (MECE):

| Section | Mutually Exclusive | Collectively Exhaustive |
|---------|-------------------|------------------------|
| **User Personas** | Each persona represents a distinct user type with unique goals and pain points. No two personas should overlap in role or motivation. | All relevant user types are represented. Ask: "Who else interacts with this system?" |
| **User Journeys** | Each journey describes a distinct path through the system. No two journeys should cover the same sequence of actions for the same persona. | All primary and secondary paths are mapped, including error/recovery paths. Ask: "What other ways do users accomplish this goal?" |
| **Feature Requirements** | Each user story captures a single, distinct behavior. No two stories should describe the same capability, even across MoSCoW categories. | All capabilities needed to solve the stated problem are present. Ask: "If we shipped only these features, would the problem be fully solved for every persona?" |
| **Acceptance Criteria** | Each criterion tests a unique condition. No two criteria should verify the same behavior with different wording. | Every feature's happy path, error path, and edge cases are covered. Ask: "What input could break this that we haven't tested?" |

**How to apply:** After completing each section, explicitly verify MECE before moving to the next:
1. **Exclusivity check** — Can any two items be merged without losing meaning? If yes, merge them.
2. **Exhaustiveness check** — Is there a scenario, user type, or capability not covered? If yes, add it.
3. **Cross-section check** — Do features in "Should Have" duplicate behaviors already in "Must Have"? Do journeys overlap with different personas doing the same thing?

## Constraints

**Always:**
- Use template.md structure exactly — preserve all sections as defined.
- Follow iterative cycle: discover → document → review per section.
- Present ALL agent findings to user — complete responses, not summaries.
- Wait for user confirmation before proceeding to the next cycle.
- Run validation checklist before declaring PRD complete.
- Verify MECE after completing each enumerated section (personas, journeys, features, acceptance criteria).

**Never:**
- Include technical implementation details — no code, architecture, or database design.
- Include API specifications — belongs in SDD.
- Skip the multi-angle validation before completing.
- Remove or reorganize template sections.
- Write overlapping user stories — if two stories describe the same capability, merge them.
- Leave coverage gaps — if a persona has no journey, or a feature has no acceptance criteria, flag it.

## Reference Materials

- [Template](template.md) — PRD template structure, write to `.start/specs/[NNN]-[name]/requirements.md`
- [Validation](validation.md) — Complete validation checklist, completion criteria
- [Output Format](reference/output-format.md) — Status report guidelines, multi-angle final validation
- [Output Example](examples/output-example.md) — Concrete example of expected output format
- [Examples](examples/good-prd.md) — Well-structured PRD reference

## Workflow

### 0. Brainstorm

Invoke Skill(start:brainstorm) to probe the user's idea before template filling.

Focus on understanding:
- What problem this solves and for whom.
- Key constraints and success criteria.
- Scope boundaries — what's in and what's out.

Output feeds into the discover/document cycle with clearer context.

### 1. Discover

Identify gaps between what is known and what template.md requires for the current section.

Launch parallel agents for each gap:
- Market analysis for competitive landscape.
- User research for personas and journeys.
- Requirements clarification for edge cases.

Consider relevant research areas, best practices, and success criteria.

### 2. Document

Update the PRD with findings for the current section:
1. Apply findings to the section being processed.
2. For each `[NEEDS CLARIFICATION]` marker, replace with findings content.

Focus only on the current section being processed. Preserve template.md structure exactly.

### 3. Review

Present ALL agent findings to user, including:
- Conflicting information or recommendations.
- Questions needing clarification.

AskUserQuestion: Approve section | Clarify [topic] | Redo discovery

### 4. Validate

Read `validation.md` and run the checklist. Read `reference/output-format.md` and run multi-angle validation.

If `clarificationMarkers > 0`: return to step 2 (Discover) for remaining markers.
If `clarificationMarkers = 0`: report status per `reference/output-format.md`.

### Entry Point

When invoked, execute step 0 (Brainstorm) first, then repeat steps 1 through 3 for each section in template.md, then execute step 4 (Validate).
