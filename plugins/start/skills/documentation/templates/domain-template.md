# [Domain Concept/Entity Name]

> **Category:** Domain/Business Rules
> **Last Updated:** [Date]
> **Status:** [Active/Under Review/Deprecated]

## Overview

**What:** [What this concept represents in the business]
**Why:** [Why this exists, business justification]
**Scope:** [Where in the application this applies]

## Business Context

### Background

[Business context and history of this domain concept]

### Stakeholders

- **[Role 1]:** [How they interact with this]
- **[Role 2]:** [How they interact with this]
- **[Role 3]:** [How they interact with this]

### Business Goals

1. [Goal 1]
2. [Goal 2]
3. [Goal 3]

## Core Concepts

### [Concept 1]

**Definition:** [Clear definition]
**Examples:** [Real-world examples]
**Constraints:** [Business constraints]

### [Concept 2]

**Definition:** [Clear definition]
**Examples:** [Real-world examples]
**Constraints:** [Business constraints]

## Business Rules

### Rule 1: [Rule Name]

**Statement:** [Clear rule statement]

**Rationale:** [Why this rule exists]

**Applies to:** [Who/what this affects]

**Exceptions:** [When this rule doesn't apply]

**Example:**
```
Given: [Initial state]
When: [Action occurs]
Then: [Expected outcome]
```

### Rule 2: [Rule Name]

[Same structure as above]

## States and Transitions

### State Machine (if applicable)

```
[Initial State]
    ↓ [Event/Action]
[Next State]
    ↓ [Event/Action]
[Final State]
```

### State Definitions

**[State 1]**
- **Meaning:** [What this state represents]
- **Entry conditions:** [How entity enters this state]
- **Exit conditions:** [How entity leaves this state]
- **Allowed actions:** [What can happen in this state]

**[State 2]**
[Same structure]

### Transition Rules

**[State A] → [State B]**
- **Trigger:** [What causes transition]
- **Conditions:** [Required conditions]
- **Side effects:** [What else happens]
- **Validation:** [What must be true]

## Permissions and Access Control

### Who Can Do What

**[Role 1]:**
- ✅ Can: [Action 1, Action 2]
- ❌ Cannot: [Action 3, Action 4]
- ⚠️ Conditional: [Action 5 - under conditions]

**[Role 2]:**
[Same structure]

### Permission Rules

**Rule:** [Permission rule statement]
**Logic:**
```
IF [condition]
  AND [condition]
  THEN [permission granted/denied]
```

## Validation Rules

### Field Validations

**[Field 1]:**
- **Type:** [Data type]
- **Required:** [Yes/No]
- **Format:** [Pattern or format]
- **Range:** [Min/max values]
- **Business rule:** [Any business constraint]

**[Field 2]:**
[Same structure]

### Cross-Field Validations

**Validation 1:** [Description]
```
IF [field1] is [value]
  THEN [field2] must be [constraint]
```

**Validation 2:** [Description]
```
[Validation logic]
```

## Workflows

### Workflow 1: [Workflow Name]

**Trigger:** [What initiates this workflow]

**Steps:**
1. **[Step 1]**
   - Actor: [Who performs this]
   - Action: [What happens]
   - Validation: [What's checked]
   - Outcome: [Result]

2. **[Step 2]**
   [Same structure]

3. **[Step 3]**
   [Same structure]

**Success Criteria:** [What defines success]
**Failure Scenarios:** [What can go wrong]

## Calculations and Algorithms

### Calculation 1: [Name]

**Purpose:** [What this calculates]

**Formula:**
```
[Mathematical or logical formula]
```

**Example:**
```
Given:
  - input1 = [value]
  - input2 = [value]

Calculation:
  result = [formula applied]

Output: [result]
```

**Edge Cases:**
- [Edge case 1 and handling]
- [Edge case 2 and handling]

## Constraints and Limits

### Business Constraints

1. **[Constraint 1]:** [Description and rationale]
2. **[Constraint 2]:** [Description and rationale]
3. **[Constraint 3]:** [Description and rationale]

### System Limits

- **[Limit 1]:** [Value and reason]
- **[Limit 2]:** [Value and reason]
- **[Limit 3]:** [Value and reason]

## Edge Cases

### Edge Case 1: [Scenario]

**Situation:** [Describe the edge case]
**Business Rule:** [How to handle it]
**Example:** [Concrete example]

### Edge Case 2: [Scenario]

[Same structure]

## Compliance and Regulations

### Regulatory Requirements

**[Regulation 1]:** [How it affects this domain concept]
**[Regulation 2]:** [How it affects this domain concept]

### Audit Requirements

- **What to log:** [Events/changes to track]
- **Retention:** [How long to keep records]
- **Who can access:** [Audit log access rules]

## Reporting and Analytics

### Key Metrics

1. **[Metric 1]:** [What it measures and why it matters]
2. **[Metric 2]:** [What it measures and why it matters]
3. **[Metric 3]:** [What it measures and why it matters]

### Reporting Requirements

- **[Report 1]:** [Purpose, frequency, audience]
- **[Report 2]:** [Purpose, frequency, audience]

## Examples and Scenarios

### Scenario 1: [Happy Path]

**Description:** [Common successful scenario]

**Flow:**
```
1. [Step with data]
2. [Step with data]
3. [Step with outcome]
```

**Business Rules Applied:** [Which rules from above]

### Scenario 2: [Error Case]

**Description:** [Common error scenario]

**Flow:**
```
1. [Step with data]
2. [Error condition]
3. [Error handling per business rules]
```

**Business Rules Applied:** [Which rules from above]

### Scenario 3: [Edge Case]

**Description:** [Unusual but valid scenario]

**Flow:**
```
1. [Step with data]
2. [Edge condition]
3. [Special handling]
```

**Business Rules Applied:** [Which rules from above]

## Integration Points

### System Touchpoints

**[System 1]:**
- **Interaction:** [How they interact]
- **Data shared:** [What data flows]
- **Trigger:** [What causes interaction]

**[System 2]:**
[Same structure]

## Glossary

**[Term 1]:** [Definition in this context]
**[Term 2]:** [Definition in this context]
**[Term 3]:** [Definition in this context]

## Related Documentation

- **Patterns:** [Pattern Doc](../patterns/doc.md) - [Technical implementation]
- **Interfaces:** [Interface Doc](../interfaces/doc.md) - [External systems]
- **Specifications:** [Spec](../specs/NNN-name/PRD.md) - [Feature requirements]

## References

- [Business document or policy]
- [Industry standard or regulation]
- [Internal decision document]

## Version History

| Date | Change | Reason | Author |
|------|--------|--------|--------|
| [Date] | Initial documentation | [Why] | [Name/Tool] |
| [Date] | Updated [aspect] | [Why] | [Name/Tool] |
