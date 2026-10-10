---
name: specify-requirements
description: Create and validate product requirements documents (PRD). Use when writing requirements, defining user stories, specifying acceptance criteria, analyzing user needs, or working on requirements.md files in .start/specs/.
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

## Constraints

**Always:**
- Use template.md structure exactly — preserve all sections as defined.
- Follow iterative cycle: discover → document → review per section.
- Present ALL agent findings to user — complete responses, not summaries.
- Obtain user disposition (confirm / correct / discard) on research findings before folding them into the PRD — an unconfirmed finding is a hypothesis, not a requirement.
- Render Standing assumptions from the clarification ledger as visible `[ASSUMPTION]` entries in the affected sections.
- Wait for user confirmation before proceeding to the next cycle.
- Run validation checklist before declaring PRD complete.
- Verify MECE after completing each enumerated section (who, features, acceptance criteria, quality gates, data).
- Assign stable IDs as sections are written — FR-N, AC-N.M, NFR-N, DR-N — and never renumber; the plan's slice tasks trace them.
- Give every quality gate a measurable, implementation-blind pass condition (threshold, or reference dataset + threshold, or behavior + condition).
- Follow the Markdown block discipline (see Markdown Blocks reference): structured, scannable Markdown — `> [!NOTE]` callouts, GitHub task lists, tables — over flat prose. Keep a stable heading spine.

**Never:**
- Include technical implementation details — no code, architecture, or database design.
- Include API specifications — belongs in SDD.
- Write quality gates as solution design — state observable behavior and thresholds (WHAT); mechanisms, harnesses, and architectures are the SDD's (HOW).
- Skip the multi-angle validation before completing.
- Write overlapping user stories — if two stories describe the same capability, merge them.
- Leave coverage gaps — if a role's goal is served by no feature, or a feature has no acceptance criteria, flag it.

## Reference Materials

- [Focus and MECE](reference/focus-and-mece.md) — PRD focus areas and MECE rules for who, features, acceptance criteria, quality gates, data
- [Markdown Blocks](../specify-plan/reference/markdown-blocks.md) — Block vocabulary and document-quality discipline (shared)
- [Template](template.md) — PRD template structure, write to `.start/specs/[NNN]-[name]/requirements.md`
- [Validation](validation.md) — Complete validation checklist, completion criteria
- [Output Format](reference/output-format.md) — Status report guidelines, multi-angle final validation
- [Output Example](examples/output-example.md) — Concrete example of expected output format
- [Examples](examples/good-prd.md) — Well-structured PRD reference

## Workflow

### 0. Brainstorm

If the specify orchestrator already ran its Clarify Intent phase, consume its ledger (answered unknowns, Standing assumptions) instead of re-probing — ask only about new gaps the template exposes. Never re-ask what clarification already answered.

Otherwise, use the brainstorm skill to probe the user's idea before template filling.

Focus on understanding:
- What problem this solves and for whom.
- Key constraints and success criteria.
- Scope boundaries — what's in and what's out.

Output feeds into the discover/document cycle with clearer context.

### 1. Discover

Identify gaps between what is known and what template.md requires for the current section.

Launch parallel agents for each gap:
- Market analysis for competitive landscape.
- User research for roles, goals, and pain points.
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

Ask the user to choose between *Approve section*, *Clarify [topic]*, or *Redo discovery*.

### 4. Validate

Read `validation.md` and run the checklist. Read `reference/output-format.md` and run multi-angle validation.

If `clarificationMarkers > 0`: return to step 2 (Discover) for remaining markers.
If `clarificationMarkers = 0`: report status per `reference/output-format.md`.

### Entry Point

Read `reference/focus-and-mece.md` for the four PRD dimensions (WHAT/WHY/WHO/WHEN) and MECE rules. Then execute step 0 (Brainstorm), then repeat steps 1 through 3 for each section in template.md, then execute step 4 (Validate).
