---
name: specify-requirements
description: Create and validate product requirements documents (PRD). Use when writing requirements, defining user stories, specifying acceptance criteria, analyzing user needs, or working on the Requirements part of a spec.md in .start/specs/.
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
- Write the `## Requirements` part of `spec.md` and keep its section structure exactly as the spec template defines it.
- Draft every section, validate, then close the part with one approval gate — not one review per section.
- Obtain user disposition (confirm / correct / discard) on research findings before folding them into the PRD — an unconfirmed finding is a hypothesis, not a requirement. Ask each observable finding (one an acceptance test could tell apart) as its own question.
- Cite Standing assumptions from the clarification ledger as `[ASSUMPTION An]` in the affected sections; the assumption itself lives only in the Ledger's Assumptions table.
- Run validation checklist before declaring PRD complete.
- Verify MECE after completing each enumerated section (who, features, acceptance criteria, quality gates, data).
- Assign stable IDs as sections are written — FR-N, AC-N.M, NFR-N, DR-N — and never renumber; the plan's slice tasks trace them.
- Give every quality gate a measurable, implementation-blind pass condition (threshold, or reference dataset + threshold, or behavior + condition).
- Follow the Markdown block discipline (see Markdown Blocks reference): structured, scannable Markdown — `> [!NOTE]` callouts, GitHub task lists, tables — over flat prose. Keep a stable heading spine.

**Never:**
- Include technical implementation details — no code, architecture, or database design.
- Include API specifications — they belong in the Solution part.
- Write quality gates as solution design — state observable behavior and thresholds (WHAT); mechanisms, harnesses, and architectures are the SDD's (HOW).
- Skip the multi-angle validation before completing.
- Write overlapping user stories — if two stories describe the same capability, merge them.
- Leave coverage gaps — if a role's goal is served by no feature, or a feature has no acceptance criteria, flag it.

## Reference Materials

- [Focus and MECE](reference/focus-and-mece.md) — PRD focus areas and MECE rules for who, features, acceptance criteria, quality gates, data
- [Markdown Blocks](../specify-plan/reference/markdown-blocks.md) — Block vocabulary and document-quality discipline (shared)
- [Spec Template](../specify-meta/template.md) — the `## Requirements` part of `.start/specs/[NNN]-[name]/spec.md`
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

Identify the gaps between what is known and what each section of the Requirements part needs.

Close them yourself in the main context: competitive landscape, roles, goals, and pain points, edge cases. Launch parallel research agents only when the domain is unknown territory — an unfamiliar market, regulation, or standard needing independent investigations — and present their findings in full.

### 2. Document

Write every section of the Requirements part with the findings, replacing each `[NEEDS CLARIFICATION]` marker. Set `gates.requirements` to `draft`.

### 3. Validate

Read `validation.md` and run the checklist. Read `reference/output-format.md` and run multi-angle validation. Fix what fails; if a gap needs the user, carry it to step 4 as a question.

### 4. Approve

Present the part for one approval:
- one or two lines per section — what it now says;
- each observable item still undecided (a finding, an assumption, an open question), each as its own question with its proposed default;
- the non-observable findings and defaults together, confirmed with the approval.

Ask the user to choose between *Approve and continue*, *Walk me through section by section*, or *Change [topic]*. A walkthrough presents one section at a time for approval. On approval set `gates.requirements` to `approved` and report status per `reference/output-format.md`.

### Entry Point

Read `reference/focus-and-mece.md` for the four PRD dimensions (WHAT/WHY/WHO/WHEN) and MECE rules. Then execute step 0 (Brainstorm) and steps 1 through 4 in order.
