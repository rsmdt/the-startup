# Output Format Reference

Status report and multi-angle validation for PRD work.

---

## PRD Status Report

Present after PRD work:

```markdown
📝 PRD Status: [spec-id]-[name]

Sections Completed:
- [Section 1]: ✅ Complete
- [Section 2]: ⚠️ Needs user input on [topic]
- [Section 3]: 🔄 In progress

Validation Status:
- [X] items passed
- [Y] items pending

Next Steps:
- [What needs to happen next]
```

---

## Multi-Angle Final Validation

Before completing the PRD, validate from multiple perspectives:

### Context Review
Verify:
- Problem statement clarity — is it specific and measurable?
- User persona completeness — do we understand our users?
- Value proposition strength — is it compelling?

### Gap Analysis
Identify:
- Gaps in user journeys
- Missing edge cases
- Unclear acceptance criteria
- Contradictions between sections

### User Input
Based on gaps found:
- Formulate specific questions using AskUserQuestion
- Probe alternative scenarios
- Validate priority trade-offs
- Confirm success criteria

### Coherence Validation
Confirm:
- Requirements completeness
- Feasibility assessment
- Alignment with stated goals
- Edge case coverage
