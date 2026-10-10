# Output Format Reference

Guidelines for PRD output. See `examples/output-example.md` for a concrete rendered example.

---

## Output Schema — PRD Status Report

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| specId | string | Yes | Spec identifier (NNN-name format) |
| title | string | Yes | Feature title |
| status | enum: `DRAFT`, `IN_REVIEW`, `COMPLETE` | Yes | Document readiness |
| sections | SectionStatus[] | Yes | Status of each PRD section |
| clarificationsRemaining | number | Yes | Count of `[NEEDS CLARIFICATION]` markers |
| acceptanceCriteria | number | Yes | Total testable acceptance criteria defined |
| openQuestions | string[] | No | Unresolved items requiring stakeholder input |

### SectionStatus

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Section name |
| status | enum: `COMPLETE`, `NEEDS_CLARIFICATION`, `IN_PROGRESS` | Yes | Current state |
| detail | string | No | What clarification is needed or what's in progress |

---

## Multi-Angle Final Validation

Before completing the PRD, validate from multiple perspectives:

### Context Review
Verify:
- Problem statement clarity — is it specific and evidenced?
- Who completeness — do we understand every role that uses it or is affected?
- Outcome strength — is "done" clear and compelling?

### Gap Analysis
Identify:
- Roles whose goal no feature serves
- Missing edge cases
- Unclear acceptance criteria
- Contradictions between sections

### User Input
Based on gaps found:
- Ask the user targeted, single-decision questions with proposed defaults
- Probe alternative scenarios
- Validate priority trade-offs
- Confirm success criteria

### Coherence Validation
Confirm:
- Requirements completeness
- Feasibility assessment
- Alignment with stated goals
- Edge case coverage
