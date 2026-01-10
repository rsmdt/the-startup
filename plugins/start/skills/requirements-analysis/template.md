---
title: "[NEEDS CLARIFICATION: Feature title]"
status: draft
version: "1.0"
---

# Product Requirements Document

## Validation Checklist

### CRITICAL GATES (Must Pass)

- [ ] All required sections are complete
- [ ] No [NEEDS CLARIFICATION] markers remain
- [ ] Problem statement is specific and measurable
- [ ] Every feature has testable acceptance criteria (EARS format)
- [ ] No contradictions between sections

### QUALITY CHECKS (Should Pass)

- [ ] Problem is validated by evidence (not assumptions)
- [ ] Context → Problem → Solution flow makes sense
- [ ] Every persona has at least one user journey
- [ ] All MoSCoW categories addressed (Must/Should/Could/Won't)
- [ ] Every metric has corresponding tracking events
- [ ] No feature redundancy (check for duplicates)
- [ ] No technical implementation details included
- [ ] A new team member could understand this PRD

---

## Product Overview

### Vision
[NEEDS CLARIFICATION: What is the one-sentence vision for this feature? What future state are we creating for users?]

### Problem Statement
[NEEDS CLARIFICATION: What specific problem are users facing today? Why is this painful for them? What are the consequences of not solving this?]

### Value Proposition
[NEEDS CLARIFICATION: Why will users choose this solution over alternatives? What unique value does it provide? How does it make their life better?]

## User Personas

### Primary Persona: [NEEDS CLARIFICATION: persona name]
- **Demographics:** [Age range, role/occupation, technical expertise level]
- **Goals:** [What are they trying to accomplish? What does success look like for them?]
- **Pain Points:** [What frustrates them about current solutions? What obstacles do they face?]

### Secondary Personas
[NEEDS CLARIFICATION: Are there other user types? If yes, define them. If no, remove this section]

## User Journey Maps

### Primary User Journey: [NEEDS CLARIFICATION: Journey name]
1. **Awareness:** [How do users discover they have this need?]
2. **Consideration:** [What alternatives do they evaluate? What criteria matter to them?]
3. **Adoption:** [What convinces them to try this solution?]
4. **Usage:** [What are the key steps in using the feature?]
5. **Retention:** [What keeps them coming back?]

### Secondary User Journeys
[NEEDS CLARIFICATION: Are there other user journeys? If yes, define them. If no, remove this section]

## Feature Requirements

### Must Have Features
[NEEDS CLARIFICATION: What are the absolute minimum features needed for this to be valuable to users?]

#### Feature 1: [NEEDS CLARIFICATION: Feature name]
- **User Story:** As a [user type], I want to [action] so that [benefit]
- **Acceptance Criteria (EARS Format):**

  Use the appropriate pattern for each criterion:
  - UBIQUITOUS: `THE SYSTEM SHALL [action]` - always-on behavior
  - EVENT-DRIVEN: `WHEN [trigger], THE SYSTEM SHALL [action]` - user/system events
  - STATE-DRIVEN: `WHILE [state], THE SYSTEM SHALL [action]` - mode-dependent
  - OPTIONAL: `WHERE [feature enabled], THE SYSTEM SHALL [action]` - configurable
  - COMPLEX: `IF [condition], THEN THE SYSTEM SHALL [action]` - business rules

  **Good Example:**
  - [ ] WHEN the user submits valid credentials, THE SYSTEM SHALL redirect to dashboard within 2 seconds

  **Bad Example:**
  - [ ] User can log in _(vague, not testable)_

  Criteria:
  - [ ] [EARS-formatted criterion]
  - [ ] [EARS-formatted criterion]
  - [ ] [EARS-formatted edge case]

#### Feature 2: [NEEDS CLARIFICATION: Feature name]
[Repeat structure as needed]

### Should Have Features
[NEEDS CLARIFICATION: What would significantly improve the experience but isn't critical for launch?]

### Could Have Features
[NEEDS CLARIFICATION: What nice-to-haves could we add if time and resources permit?]

### Won't Have (This Phase)
[NEEDS CLARIFICATION: What is explicitly out of scope for this phase? What are we NOT building?]

## Detailed Feature Specifications

### Feature: [NEEDS CLARIFICATION: Pick the most complex feature from above]
**Description:** [Detailed explanation of how this feature works]

**User Flow:**
1. User [First action]
2. System [Response]
3. User [Next action]

**Business Rules:**
- Rule 1: [Specification - When X happens, then Y should occur]
- Rule X: ...

**Edge Cases:**
- Scenario 1: [What could go wrong?] → Expected: [How should system handle it?]
- Scenario X: ...

## Success Metrics

### Key Performance Indicators
[NEEDS CLARIFICATION: How will we measure if this feature is successful?]

- **Adoption:** [Target number/percentage of users who try the feature]
- **Engagement:** [Target frequency of use or actions per user]
- **Quality:** [Target error rate, success rate, or satisfaction score]
- **Business Impact:** [Revenue, retention, or other business metric]

### Tracking Requirements
[NEEDS CLARIFICATION: What user actions and data points must we track to validate our success metrics and make informed product decisions?]

| Event | Properties | Purpose |
|-------|------------|---------|
| [User action] | [What data to capture] | [Why we track this] |

---

## Constraints and Assumptions

### Constraints [NEEDS CLARIFICATION: What are limiting factors]
- [Budget, timeline, or resource limitations]
- [Technical or platform constraints]
- [Legal or compliance requirements]

### Assumptions [NEEDS CLARIFICATION: What are we assuming that is not explicitly defined]
- [What are we assuming about users?]
- [What are we assuming about the market?]
- [What dependencies are we assuming will be available?]

## Risks and Mitigations
[NEEDS CLARIFICATION: What are potential risks? How do we intend to tackle them?]

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| [Risk description] | [High/Medium/Low] | [High/Medium/Low] | [How to prevent or handle] |

## Open Questions
[NEEDS CLARIFICATION: What requires more external input? What decisions do we need to continue?]

- [ ] [Question that needs stakeholder input]
- [ ] [Decision that needs to be made]
- [ ] [Information that needs to be gathered]

---

## Supporting Research

### Competitive Analysis
[NEEDS CLARIFICATION: How do competitors solve this problem? What can we learn from them?]

### User Research
[NEEDS CLARIFICATION: What user research has been done? Key findings?]

### Market Data
[NEEDS CLARIFICATION: Any relevant market size, trends, or data points?]
