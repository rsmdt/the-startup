# SDD Validation Checklist

Use this checklist to validate SDD completeness before proceeding to the implementation plan. Step 5's
review panel covers the coverage, traceability, MECE, and design-risk items; check every other item
yourself before completing the SDD.

## Structure Validation

- [ ] **All sections are complete** - No empty or placeholder sections
- [ ] **No [NEEDS CLARIFICATION] markers remain** - All markers replaced with content
- [ ] **Template structure preserved** - The heading spine is intact; the multi-component section (`reference/conditional-sections.md`) added only when the feature spans components

## Approach

- [ ] **Architecture pattern stated with rationale** - A named pattern and why it fits the PRD's constraints
- [ ] **Technical constraints captured** - Anything the PRD does not carry (language, framework, deployment limits)
- [ ] **Cross-cutting concerns addressed** - Security, error handling, performance, observability each stated or explicitly not applicable
- [ ] **Deployment stated** - "No change" or what changes, including rollback

## Context

- [ ] **Read-first sources are real paths** - Discovered from the codebase and docs, not assumed
- [ ] **Each source says why it matters** - One line, specific to this feature
- [ ] **Boundaries are explicit** - Must preserve, can modify, must not touch

## Components

- [ ] **Every component has a directory mapping** - Where its code lives, marked NEW or MODIFY
- [ ] **Every component in the diagram** - No orphan components
- [ ] **Component names consistent** - Same names across diagram, text, directory map, and Runtime
- [ ] **MECE: Components** - Each component has a single distinct responsibility; every PRD requirement maps to exactly one component; no two components own the same domain logic

## Interfaces & Data

- [ ] **Every interface this feature adds or changes is listed** - Direction, type, auth, and a contract or a reference to existing documentation
- [ ] **Request and response shapes defined** - Inline or by reference; error responses included
- [ ] **Schema changes specified when present** - Tables, columns, indexes, relationships
- [ ] **Entities modeled** - Every datum the interfaces and Behavior reference lives in an entity
- [ ] **Examples use actual column names** - Verified against migration files, never pseudocode (e.g. `product_bundle_id` not `bundle_id`, `settled_at` not `period_start`)
- [ ] **MECE: Interfaces & Data** - No two interfaces serve the same consumer-to-provider path; all communication paths documented; each entity owns a distinct slice of the domain

## Runtime

- [ ] **Primary flow documented** - Trigger, validation, processing, outcome, with a sequence diagram
- [ ] **Error handling covers all types** - Invalid input, upstream failure, business-rule violation
- [ ] **Non-obvious logic has a traced example** - Real schema names and concrete rows showing how the logic evaluates

## Decisions

- [ ] **Every architecture-shaping choice is a Decision block** - Options weighed, rationale, accepted trade-off
- [ ] **All Decisions confirmed by user** - `[x] User confirmed` on each

## Behavior

- [ ] **EARS criteria defined** - Every statement uses an EARS pattern and names a concrete response
- [ ] **MECE: Behavior** - Each EARS criterion specifies a unique behavior and names the AC it refines; every edge and error state an AC implies is stated by the AC or an EARS criterion; none invents behavior the PRD does not ask for
- [ ] **Happy, error, and edge paths covered** - For every component the feature touches

## PRD Alignment

- [ ] **All PRD requirements addressable** - Every requirement maps to the design
- [ ] **Quality gates realizable** - The design can meet every PRD `NFR-N` gate; the SDD references them rather than restating them
- [ ] **No scope creep** - Nothing designed beyond PRD scope

## Risks & Consistency

- [ ] **Risks and gotchas listed** - Known limitations, technical debt in the touched area, timing and environment quirks
- [ ] **No contradictions between sections** - Names, shapes, and flows agree throughout
- [ ] **A developer could implement from this** - Clear enough for someone new to the codebase

## Completion Criteria

✅ **SDD is complete when:**
- All checklist items pass
- All Decisions have user confirmation
- No open questions remain
- A developer could implement without further clarification
- Ready for implementation planning (PLAN)
