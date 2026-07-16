# PRD Validation Checklist

Use this checklist to validate PRD completeness before proceeding to SDD.

## Structure Validation

- [ ] **All required sections are complete** - No empty or placeholder sections
- [ ] **No [NEEDS CLARIFICATION] markers remain** - All markers replaced with content
- [ ] **Template structure preserved** - No sections added, removed, or reorganized

## Content Quality

### Problem Definition
- [ ] **Problem statement is specific and measurable** - Clear metrics for problem impact
- [ ] **Problem is validated by evidence** - Data, user research, or market analysis (not assumptions)
- [ ] **Context → Problem → Solution flow makes sense** - Logical narrative

### User Understanding
- [ ] **Every persona has demographics** - Age, role, technical expertise
- [ ] **Every persona has goals** - What success looks like for them
- [ ] **Every persona has pain points** - Current frustrations
- [ ] **Every persona has at least one user journey** - End-to-end flow

### Requirements Quality
- [ ] **All MoSCoW categories addressed** - Must/Should/Could/Won't all defined
- [ ] **Every feature has a user story** - As a [user], I want [action] so that [benefit]
- [ ] **Every feature has testable acceptance criteria** - Specific, verifiable conditions
- [ ] **No feature redundancy** - Check for duplicates or overlapping features
- [ ] **No contradictions between sections** - Consistent throughout

### Success Criteria
- [ ] **KPIs defined for adoption** - User acquisition/activation targets
- [ ] **KPIs defined for engagement** - Usage frequency/depth targets
- [ ] **KPIs defined for quality** - Error rate, satisfaction targets
- [ ] **KPIs defined for business impact** - Revenue, retention targets
- [ ] **Every metric has corresponding tracking events** - How to measure

### Non-Functional Requirements & Evaluation Gates
- [ ] **Every NFR has a measurable fit criterion** - Verifiable without knowing the implementation
- [ ] **Rejected NFR categories carry a one-line reason** - Considered-and-rejected, not silently dropped
- [ ] **Every evaluation criterion is business-observable** - Dataset/threshold or behavior/pass-condition; no harness or mechanism design
- [ ] **EV gates cover every release-blocking behavior** - Nothing that could quietly regress is ungated

### Stakeholders, Data & Dependencies
- [ ] **Non-user stakeholders identified** - Support, ops, compliance, legal where affected
- [ ] **Every data item has purpose, retention, and classification** - Nothing held without a stated why
- [ ] **User rights addressed** - View/export/correct/delete obligations stated
- [ ] **Dependencies have owner, need-by, and impact-if-late** - Not buried in assumptions

### Constraints & Risks
- [ ] **Constraints identified** - Budget, timeline, technical, compliance
- [ ] **Assumptions documented** - Explicit about what we're assuming; Standing ledger entries visible as [ASSUMPTION]
- [ ] **Risks identified with mitigations** - What could go wrong and how to handle

### Identifiers & Traceability
- [ ] **Stable IDs assigned** - FR-N, AC-N.M, NFR-N, EV-N, DR-N; unique and never renumbered
- [ ] **Glossary defines every domain term** - One term, one meaning

## Boundary Validation

- [ ] **No technical implementation details included** - No code, architecture, or database design
- [ ] **No API specifications** - Belongs in SDD
- [ ] **Focus on WHAT and WHY, not HOW** - Business requirements only

## Clarity Validation

- [ ] **A new team member could understand this PRD** - Self-contained and clear
- [ ] **Jargon is defined** - Domain terms explained
- [ ] **Acronyms are expanded** - First use includes full form

## Cross-Section Consistency

Run these checks across the entire document:
- [ ] **User personas match user journeys** - All personas have journeys
- [ ] **Features align with user goals** - Each feature maps to a persona goal
- [ ] **Metrics map to features** - Success measured for key features
- [ ] **Risks relate to requirements** - Identified risks are relevant

## Completion Criteria

✅ **PRD is complete when:**
- All checklist items pass
- User has reviewed and approved content
- No open questions remain
- Ready for technical specification (SDD)
