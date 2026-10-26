# PRD Validation Checklist

Use this checklist to validate PRD completeness before proceeding to SDD.

## Structure Validation

- [ ] **All sections are complete** - No empty or placeholder sections; Data may be omitted
- [ ] **No [NEEDS CLARIFICATION] markers remain** - All markers replaced with content
- [ ] **Template structure preserved** - No sections added, removed, or reorganized

## Content Quality

### Problem & Outcome
- [ ] **Problem is specific** - Who has it and what it costs
- [ ] **Outcome states "done"** - A reader can tell what is true for users once this ships and why they would choose it
- [ ] **Problem → Outcome → Features flow makes sense** - Logical narrative a new team member could follow

### Who
- [ ] **Every row has a goal or concern and a success condition** - Not a demographic sketch
- [ ] **Every role's goal is served by at least one feature** - No role without a feature; no feature without a role

### Features
- [ ] **Every feature has a user story** - As a [role], I want [action] so that [benefit]
- [ ] **Every feature has testable acceptance criteria** - Given / When / Then with an observable outcome, covering happy, error, and edge paths
- [ ] **Must / Should / Won't all addressed** - Won't-have names the confusable capabilities
- [ ] **No feature redundancy** - No capability stated twice, even across priorities
- [ ] **No contradictions between sections** - Consistent throughout

### Quality Gates
- [ ] **Every gate is measurable without knowing the implementation** - Threshold, reference dataset + threshold, or behavior + pass condition
- [ ] **Every quality concern the brief or user raises is gated** - No gate for a concern neither raised
- [ ] **No gate is solution design** - WHAT must hold, never HOW it is achieved

### Data & Dependencies
- [ ] **Every data item has a purpose** - Nothing held without a stated why
- [ ] **Dependencies have owner, need-by, and impact-if-late** - Not buried in assumptions
- [ ] **Constraints identified** - Budget, timeline, platform, compliance; Standing ledger entries cited as [ASSUMPTION An]

### Identifiers & Terms
- [ ] **Stable IDs assigned** - FR-N, AC-N.M, NFR-N, DR-N; unique and never renumbered
- [ ] **Glossary defines every ambiguous term** - One term, one meaning; acronyms expanded on first use

## Boundary Validation

- [ ] **No technical implementation details included** - No code, architecture, or database design
- [ ] **No API specifications** - Belongs in SDD
- [ ] **Focus on WHAT and WHY, not HOW** - Business requirements only

## Completion Criteria

✅ **PRD is complete when:**
- All checklist items pass
- User has reviewed and approved content
- No open questions remain
- Ready for technical specification (SDD)
