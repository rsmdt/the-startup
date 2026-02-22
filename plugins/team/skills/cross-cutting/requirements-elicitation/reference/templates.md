# Requirements Templates — Detailed Reference

Document templates and formats for requirements artifacts. Loaded on demand when creating deliverables.

---

## User Story Format

```
Format:
As a [role],
I want [capability],
So that [benefit].

Components:
- Role: Who benefits (be specific)
- Capability: What they can do (action, not solution)
- Benefit: Why it matters (business value)

Example:
As a customer service representative,
I want to see a customer's order history when they call,
So that I can resolve their issues without asking them to repeat information.
```

---

## Acceptance Criteria (Given-When-Then)

```
Format:
Given [context/precondition]
When [action/event]
Then [expected outcome]

Example:
Feature: Order Cancellation

Scenario: Cancel order before shipping
Given an order in "confirmed" status
And the order has not been shipped
When the customer requests cancellation
Then the order status changes to "cancelled"
And the customer receives a cancellation confirmation email
And the payment is refunded within 3-5 business days

Scenario: Cannot cancel shipped order
Given an order in "shipped" status
When the customer requests cancellation
Then the cancellation is rejected
And the customer is directed to the returns process
```

---

## Edge Cases and Exceptions

Document what happens when things go wrong:

```
Feature: Password Reset

Happy Path:
- User requests reset → Email sent → User clicks link → Sets new password

Edge Cases:
| Scenario | Expected Behavior |
|----------|-------------------|
| Email not found | Show same success message (security) |
| Link expired (>24h) | Show "link expired" with new reset option |
| Link already used | Show "link already used" message |
| Weak password | Show requirements, block submission |
| Same as old password | Show error, require different password |
| User locked out | Still send reset email (unlock flow) |
```

---

## Non-Functional Requirements Template

```
NFR Template:
┌─────────────────────────────────────────────────────────────┐
│ Category: Performance                                        │
├─────────────────────────────────────────────────────────────┤
│ Requirement: Response Time                                   │
│ Measure: 95th percentile page load time                     │
│ Target: < 2 seconds                                          │
│ Context: Desktop browser, 4G connection                     │
│ Priority: Must Have                                          │
└─────────────────────────────────────────────────────────────┘

Common Categories:
- Performance: Speed, throughput, latency
- Scalability: Users, data volume, geographic distribution
- Availability: Uptime, recovery time, disaster recovery
- Security: Authentication, authorization, encryption
- Usability: Accessibility, learnability, efficiency
- Maintainability: Modularity, testability, documentation
```

---

## Feature Request Template

```markdown
# Feature: [Name]

## Problem Statement
[What problem does this solve?]

## User Stories
- As a [role], I want [what] so that [why]

## Acceptance Criteria
- Given [context] when [action] then [outcome]

## Out of Scope
- [What this feature does NOT include]

## Dependencies
- [Other features or systems required]

## Open Questions
- [Unresolved issues needing discussion]
```

---

## Requirements Document Template

```markdown
# [Project Name] Requirements Specification

## 1. Introduction
### 1.1 Purpose
### 1.2 Scope
### 1.3 Definitions

## 2. Overall Description
### 2.1 Product Perspective
### 2.2 User Classes
### 2.3 Constraints

## 3. Functional Requirements
### 3.1 [Feature Area 1]
### 3.2 [Feature Area 2]

## 4. Non-Functional Requirements
### 4.1 Performance
### 4.2 Security
### 4.3 Usability

## 5. Appendices
### A. Stakeholder Register
### B. Traceability Matrix
```
