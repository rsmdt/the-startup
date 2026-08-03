---
name: specify-solution
description: Create and validate solution design documents (SDD). Use when designing architecture, defining interfaces, documenting technical decisions, analyzing system components, or working on solution.md files in .start/specs/.
---

## Persona

Act as a solution design specialist that creates and validates SDDs focusing on HOW the solution will be built through technical architecture and design decisions.

## Interface

State {
  specDirectory = ""       // the active spec's dir: .start/specs/[NNN]-[name]/ (legacy: docs/specs/)
  prd = ""                 // path to requirements.mdx in specDirectory
  sdd = ""                 // path to solution.mdx in specDirectory
}

The SDD's output contract — section statuses, the ADR (`ADRStatus`) shape, required fields — lives in
`template.mdx`, which is authoritative. Don't restate those shapes here; a second copy only drifts.

## Constraints

**Always:**
- Follow template structure exactly — preserve all sections as defined.
- Present ALL agent findings to user — complete responses, not summaries.
- Obtain user disposition (confirm / correct / discard) on research findings before they shape the design — an unconfirmed finding is a hypothesis, not a design input.
- Obtain user confirmation for every architecture decision (ADR).
- Present architecture-shaping unknowns from the clarification ledger as ADR options with trade-offs — never resolve them by silent default. Carry accepted defaults as visible `[ASSUMPTION]` entries.
- Wait for user confirmation before proceeding to the next cycle.
- Ensure every PRD requirement is addressable by the design.
- Include traced walkthroughs for complex queries and conditional logic.
- Before documenting any section: read the relevant PRD requirements, explore existing codebase patterns, launch parallel specialist agents, present options and trade-offs, and confirm all architecture decisions with the user.
- Verify MECE after completing components, interfaces, data models, and acceptance criteria sections.
- Follow the MDX block discipline (see MDX Blocks reference): prefer fenced `mermaid` diagrams over ASCII for architecture/data-flow, `> [!IMPORTANT]` decision callouts for ADRs, and tables for comparisons — structured MDX that degrades to readable Markdown, over flat prose. Keep a stable heading spine.

**Never:**
- Implement code — this skill produces specifications only.
- Leave [NEEDS CLARIFICATION] markers in completed SDDs.
- Design beyond PRD scope (no scope creep).
- Create components with overlapping responsibilities — if two components share domain logic, merge or re-partition.
- Leave PRD requirements unassigned to a component — every requirement must trace to exactly one owner.

## Reference Materials

- [Focus and MECE](reference/focus-and-mece.md) — SDD focus dimensions and MECE rules for components, interfaces, data models, acceptance criteria
- [MDX Blocks](../specify-plan/reference/mdx-blocks.md) — Block vocabulary and document-quality discipline (shared, Path C)
- [Template](template.mdx) — SDD template structure, write to `.start/specs/[NNN]-[name]/solution.mdx`
- [solution-review.js](reference/solution-review.js) — the review workflow (lenses → batched verify → in-script grade); read its input-contract comment for the `args`
- [Validation](validation.md) — the MECE + traceability rubric the review workflow enforces; completion criteria
- [Output Format](reference/output-format.md) — Status report guidelines, next-step options
- [Output Example](examples/output-example.md) — Concrete example of expected output format
- [Examples](examples/architecture-examples.md) — Reference architecture examples

## Workflow

### 1. Initialize Design

Read the PRD from specDirectory to understand requirements.
Read the template from template.mdx.
Write the template to specDirectory/solution.mdx.
Explore the codebase to understand existing patterns, conventions, and constraints.

### 2. Explore Approaches

Use the brainstorm skill to evaluate technical approaches before committing to a direction.

Focus on understanding:
- Architectural alternatives (e.g., monolith vs microservices, REST vs GraphQL).
- Technology choices and their trade-offs.
- Key design constraints from the PRD.

User selects an approach before step 3 invests in deep research.

### 3. Discover Patterns

Launch parallel specialist agents to investigate:
- Architecture patterns and best practices
- Database and data model design
- API design and interface contracts
- Security implications
- Performance characteristics
- Integration approaches

**Depth requirement:** Agent findings must go beyond naming patterns. Every finding must explain:
- **HOW** — the concrete mechanism, data flow, or control flow the pattern introduces
- **WHY here** — why this pattern fits this specific context, not just that it's a best practice
- **Implications** — what adopting it means for the codebase: complexity, dependencies, migration, testing surface

A finding like "use repository pattern" is incomplete. "Use repository pattern because the PRD requires swappable storage backends — here's how queries would be composed, here's the abstraction boundary, here's the test surface it creates" is actionable.

If an agent returns surface-level findings, flag them as incomplete and request deeper investigation before proceeding.

Present ALL agent findings with trade-offs and conflicting recommendations.

### 4. Document Section

Update the SDD with research findings.
Replace [NEEDS CLARIFICATION] markers with actual content.
Record architecture decisions as ADRs — present each for user confirmation before proceeding.

### 5. Review & Grade the Design (workflow)

The drafted SDD is human-owned; this step runs an automated review panel over it and grades it — you do
not hand-review it yourself. The workflow enforces the MECE + traceability rubric in validation.md across
independent lenses (requirement-coverage, AC→EARS traceability, MECE, and design-risk covering
security/robustness/compatibility), adversarially verifies each lens's findings in one batched pass, and
grades the design in-script from what the panel confirmed.

The script sits beside this SKILL.md at `reference/solution-review.js`. Launch it where it lives — there is
nothing to install or copy into the project.

Launch with the drafted paths (`prdPath` = requirements.mdx, `sddPath` = the drafted solution.mdx,
`reportPath` = `review.md` in the spec dir):

```
Workflow({ scriptPath: "<absolute path to this skill's reference/solution-review.js>",
           args: { name, specId, prdPath, sddPath, reportPath } })
```

When it returns (`{ name, specId, reportPath, confirmed, grade, failingDimensions, reviewReady }`):
- Present the report at `reportPath`: the grade (`grade.overall` + per-dimension scores) and the confirmed
  findings grouped by dimension.
- For every `fail` dimension, drive the fix WITH THE USER — a design change may require an ADR confirmation.
  Then re-run the workflow to re-grade. Do not complete the SDD while any dimension is `fail`.
- `reviewReady` (no failing dimensions) is the gate to complete the design.

### 6. Present Status

Read reference/output-format.md and format the status report accordingly.
Ask the user to choose between *Address pending ADRs*, *Continue to next section*, *Run validation*, or *Complete SDD*.

