---
name: specify-solution
description: Create and validate solution design documents (SDD). Use when designing architecture, defining interfaces, documenting technical decisions, analyzing system components, or working on solution.md files in .start/specs/.
---

## Persona

Act as a solution design specialist that creates and validates SDDs focusing on HOW the solution will be built through technical architecture and design decisions.

## Interface

SddSection {
  status: Complete | NeedsDecision | InProgress | Pending
  adrs?: ArchitectureDecision[]
}

ArchitectureDecision {
  id: string               // ADR-1, ADR-2, ...
  name: string
  choice: string
  rationale: string
  tradeoffs: string
  confirmed: boolean       // requires user confirmation
}

State {
  specDirectory = ""       // .start/specs/[NNN]-[name]/ (or legacy docs/specs/)
  prd = ""                 // path to requirements.md (or product-requirements.md)
  sdd = ""                 // path to solution.md (or solution-design.md)
  sections: SddSection[]
  adrs: ArchitectureDecision[]
}

## Constraints

**Always:**
- Focus exclusively on research, design, and documentation — never implementation.
- Follow template structure exactly — preserve all sections as defined.
- Present ALL agent findings to user — complete responses, not summaries.
- Obtain user confirmation for every architecture decision (ADR).
- Wait for user confirmation before proceeding to the next cycle.
- Ensure every PRD requirement is addressable by the design.
- Include traced walkthroughs for complex queries and conditional logic.
- Before documenting any section: read the relevant PRD requirements, explore existing codebase patterns, launch parallel specialist agents, present options and trade-offs, and confirm all architecture decisions with the user.
- Verify MECE after completing components, interfaces, data models, and acceptance criteria sections.

**Never:**
- Implement code — this skill produces specifications only.
- Skip user confirmation on architecture decisions.
- Remove or reorganize template sections.
- Leave [NEEDS CLARIFICATION] markers in completed SDDs.
- Design beyond PRD scope (no scope creep).
- Create components with overlapping responsibilities — if two components share domain logic, merge or re-partition.
- Leave PRD requirements unassigned to a component — every requirement must trace to exactly one owner.

## Reference Materials

- [Focus and MECE](reference/focus-and-mece.md) — SDD focus dimensions and MECE rules for components, interfaces, data models, acceptance criteria
- [Template](template.md) — SDD template structure, write to `.start/specs/[NNN]-[name]/solution.md`
- [Validation](validation.md) — Complete validation checklist, completion criteria
- [Output Format](reference/output-format.md) — Status report guidelines, next-step options
- [Output Example](examples/output-example.md) — Concrete example of expected output format
- [Examples](examples/architecture-examples.md) — Reference architecture examples

## Workflow

### 1. Initialize Design

Read the PRD from specDirectory to understand requirements.
Read the template from template.md.
Write the template to specDirectory/solution.md.
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

### 5. Validate Design

Read validation.md and run the full checklist, focusing on:

**MECE Validation (run first):**

Mutually Exclusive — no overlap:
- Component responsibilities — does any PRD requirement map to more than one component? If yes, re-partition.
- Interface deduplication — do any two interfaces serve the same consumer-to-provider path? If yes, merge.
- Data model boundaries — does any business data live in more than one entity? If yes, pick a single owner.
- Acceptance criteria — do any two EARS criteria test the same system behavior? If yes, consolidate.

Collectively Exhaustive — no gaps:
- PRD coverage — does every PRD requirement map to at least one component? If not, assign it.
- Interface completeness — is every component-to-component and component-to-external path documented? If not, add it.
- Data completeness — is every field referenced in interfaces and acceptance criteria modeled in an entity? If not, add it.
- Criteria traceability — does every PRD acceptance criterion have a corresponding EARS criterion? If not, write it.

**Structural Validation:**

Boundary validation:
- Layer separation — presentation, business, data properly separated?
- Dependency direction — no circular dependencies?
- Integration points — all system boundaries documented?

Consistency verification:
- Naming consistency — components, interfaces, concepts named consistently?
- Pattern adherence — architectural patterns applied consistently?
- Cross-cutting concerns — security, error handling, logging, performance?

### 6. Present Status

Read reference/output-format.md and format the status report accordingly.
Ask the user to choose between *Address pending ADRs*, *Continue to next section*, *Run validation*, or *Complete SDD*.

