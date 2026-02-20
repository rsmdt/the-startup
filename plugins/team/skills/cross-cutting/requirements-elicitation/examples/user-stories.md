# User Story Examples

## Context

This document provides annotated examples of well-crafted user stories with acceptance criteria, alongside examples of common mistakes and how to fix them. Use these as reference when writing or reviewing stories.

The format covered in SKILL.md is:

```
As a [role],
I want [capability],
So that [benefit].
```

The acceptance criteria format is Given-When-Then. This document shows what that looks like in practice across different domains.

---

## Good vs. Bad User Stories

Understanding what makes a story weak helps you write stronger ones from the start.

### Example 1: The Vague Role

**Bad:**

```
As a user,
I want to see notifications,
So that I know what's happening.
```

Problems:
- "User" is too broad — which type of user? The story means nothing without knowing who it serves.
- "See notifications" is undefined — what kind, from where, in what format?
- "Know what's happening" has no testable outcome.

**Good:**

```
As an order fulfillment coordinator,
I want to receive a real-time alert when an order's shipping status changes,
So that I can proactively contact customers before they reach out to us.
```

Why this works:
- The role is specific enough to know who has this need and why.
- The capability describes an observable action (receive an alert, when something specific happens).
- The benefit explains the business value and the decision the coordinator needs to make.

---

### Example 2: The Solution-First Story

**Bad:**

```
As a manager,
I want a dashboard with pie charts,
So that I can see the data.
```

Problems:
- "Pie charts" is a solution, not a capability. The user's real need is to understand data — the visualization type is an implementation choice.
- "See the data" says nothing about what data, why it matters, or what the manager will do with it.

**Good:**

```
As a customer success manager,
I want to see each account's health score and its trend over the past 30 days,
So that I can identify at-risk accounts before they churn.
```

Why this works:
- The capability describes what the manager needs to know, not how it's displayed.
- The benefit ties directly to a business outcome (preventing churn).
- This leaves implementation open — a table, a chart, or a sorted list all become valid choices.

---

### Example 3: The Missing Benefit

**Bad:**

```
As a developer,
I want API documentation,
So that I have it.
```

Problems:
- The benefit "so that I have it" restates the want and adds no value.
- Removing the benefit makes this story unvalidatable — you can't measure success without knowing why.

**Good:**

```
As a developer integrating with the payment API,
I want interactive endpoint documentation with request/response examples,
So that I can implement and test my integration without needing support from the platform team.
```

Why this works:
- The role is specific (developer integrating, not developers in general).
- The benefit is a measurable outcome — reduced support requests, independent implementation.

---

### Example 4: The Epic Masquerading as a Story

**Bad:**

```
As a shopper,
I want to be able to browse, filter, search, add items to my cart, check out, and track my order,
So that I can buy things online.
```

Problems:
- This is an entire e-commerce system in one story.
- It cannot be completed in a single sprint.
- There is no way to write meaningful acceptance criteria for something this large.

**Good — Split into stories:**

```
As a shopper,
I want to filter products by category and price range,
So that I can find relevant items without browsing everything.
```

```
As a shopper,
I want to add items to a cart and adjust quantities before checkout,
So that I can review my order before committing to purchase.
```

```
As a shopper,
I want to receive an email confirmation with my order number after purchase,
So that I have a record and can track my order if there are issues.
```

Each story is independently deliverable and testable.

---

## Well-Crafted Stories with Acceptance Criteria

### Story 1: Authentication — Password Reset

```
As a registered user who has forgotten their password,
I want to reset my password via a link sent to my email,
So that I can regain access to my account without contacting support.
```

**Acceptance Criteria:**

```
Scenario: Request a password reset
Given I am on the login page
When I click "Forgot password" and enter my registered email address
Then I receive an email containing a password reset link within 2 minutes
And the link expires after 24 hours

Scenario: Complete a password reset
Given I have received a valid, unexpired reset link
When I click the link and enter a new password that meets the requirements
Then my password is updated
And I am redirected to the login page with a confirmation message
And the reset link is invalidated so it cannot be reused

Scenario: Reset link has expired
Given I have a reset link that is more than 24 hours old
When I click the link
Then I see a message that the link has expired
And I am offered the option to request a new reset link

Scenario: Email address is not registered
Given I am on the "Forgot password" page
When I enter an email address that has no account
Then I see the same success message as if the email were registered
And no email is sent
```

Note on the last scenario: showing the same message regardless prevents user enumeration — an attacker cannot determine which emails are registered by trying different addresses.

---

### Story 2: Search — Finding Records

```
As a support agent,
I want to search for customer accounts by name, email, or account ID,
So that I can quickly locate the right account when a customer contacts us.
```

**Acceptance Criteria:**

```
Scenario: Search returns results
Given I am on the customer search page
When I type at least 3 characters into the search field
Then results appear within 500ms
And each result shows the customer name, email, and account status

Scenario: Search by partial name
Given there are customers named "Sarah Johnson" and "Sarah Williams"
When I search for "Sarah"
Then both customers appear in the results

Scenario: Search by email
Given a customer with email "jane.doe@example.com"
When I search for "jane.doe"
Then that customer appears in the results

Scenario: No results found
Given no customers match the search term
When I complete a search
Then I see a "No results found" message
And I see a suggestion to check the spelling or try a different identifier

Scenario: Search is too short
Given I have typed fewer than 3 characters
When I look at the search field
Then no results are shown
And I see a hint indicating I need to type at least 3 characters
```

---

### Story 3: Permissions — Role-Based Access

```
As an organization administrator,
I want to assign team members to roles with different permission levels,
So that each person only has access to the features and data their job requires.
```

**Acceptance Criteria:**

```
Scenario: Assign a role to a team member
Given I am an administrator on the organization account
When I navigate to Team Settings and select a team member
Then I can assign them one of the available roles: Viewer, Editor, or Admin
And the change takes effect immediately

Scenario: Viewer cannot edit content
Given a team member has the Viewer role
When they view a record
Then all fields are read-only
And the "Save" and "Delete" buttons are not visible

Scenario: Editor can edit but not manage users
Given a team member has the Editor role
When they navigate to Team Settings
Then they can view the team list
But they cannot add, remove, or change roles for other members

Scenario: Removing a team member
Given I am an administrator
When I remove a team member from the organization
Then they immediately lose access to all organization data
And they receive an email notifying them of the change

Scenario: Last administrator cannot be demoted
Given an organization has exactly one administrator
When that administrator attempts to change their own role to Editor
Then the action is blocked
And they see an error: "Organizations must have at least one administrator"
```

---

### Story 4: Data Export — CSV Download

```
As a finance team member,
I want to export transaction records as a CSV file filtered by date range,
So that I can import them into our accounting software for monthly reconciliation.
```

**Acceptance Criteria:**

```
Scenario: Export with a valid date range
Given I am on the Transactions page
When I select a start date and end date and click "Export CSV"
Then a CSV file downloads containing all transactions within that range
And the file includes columns: transaction ID, date, description, amount, status

Scenario: Export with no results
Given no transactions exist in the selected date range
When I click "Export CSV"
Then a CSV file downloads containing only the header row
And no error is shown

Scenario: Export is limited to my permissions
Given I am a Viewer with access to only my team's transactions
When I export a CSV
Then the file contains only transactions I am authorized to see
And no transactions from other teams appear

Scenario: Large export
Given there are more than 10,000 transactions in the selected range
When I click "Export CSV"
Then the export is queued
And I receive an email with a download link when it is ready
And the link expires after 24 hours

Scenario: Invalid date range
Given I set the start date after the end date
When I attempt to export
Then the export button is disabled
And I see an inline error: "Start date must be before end date"
```

---

### Story 5: Notifications — Real-Time Alerts

```
As a warehouse operations manager,
I want to receive an alert when inventory for a SKU drops below its reorder threshold,
So that I can initiate a purchase order before stock runs out.
```

**Acceptance Criteria:**

```
Scenario: Alert triggered when threshold is crossed
Given SKU-1234 has a reorder threshold of 50 units
And current stock is 51 units
When a sale reduces stock to 50 units or fewer
Then an alert appears in my notification center within 60 seconds
And the alert includes: SKU, product name, current quantity, and threshold

Scenario: Alert is not repeated for the same condition
Given an alert has already been sent for SKU-1234 being below threshold
When subsequent sales occur without restocking
Then no additional alerts are sent for the same threshold breach

Scenario: Alert clears when stock is replenished
Given SKU-1234 is in an alerted state
When stock is updated to a quantity above the threshold
Then the alert is dismissed from the notification center
And a new alert will be sent if stock drops again in the future

Scenario: Manager can configure thresholds
Given I am viewing a SKU's settings
When I update the reorder threshold and save
Then future alerts use the new threshold value
And the change is logged with my name and a timestamp

Scenario: Alert delivery via email
Given I have email notifications enabled for inventory alerts
When an inventory threshold is crossed
Then I receive an email within 5 minutes
And the email contains a direct link to the affected SKU's page
```

---

## Acceptance Criteria Patterns

Certain patterns appear repeatedly in well-written acceptance criteria. Recognizing them helps you write criteria faster and more completely.

### The Happy Path + Variations Pattern

Always start with the success scenario, then add variations:

```
Scenario: [Feature] succeeds under normal conditions
Scenario: [Feature] when input is at the boundary (minimum/maximum)
Scenario: [Feature] when the user has limited permissions
Scenario: [Feature] fails gracefully with a clear error
```

### The Security Variant Pattern

For any feature that handles data or access:

```
Scenario: [User A] cannot access [User B]'s data
Scenario: Unauthenticated user is redirected to login
Scenario: Error messages do not expose sensitive information
```

### The Scale Variant Pattern

For features involving lists, searches, or exports:

```
Scenario: Returns correct results when list is empty
Scenario: Returns correct results when list has one item
Scenario: Handles large datasets without timeout or data loss
Scenario: Pagination works correctly at boundaries (first page, last page)
```

---

## Anti-Patterns in Acceptance Criteria

| Anti-Pattern | Problem | Fix |
|---|---|---|
| **"Works correctly"** | Not testable — correct according to whom? | Define the specific expected output |
| **"Should be fast"** | Relative and unmeasurable | "Loads within 2 seconds on a 4G connection" |
| **"The system handles errors"** | What errors? What does handling mean? | Specify the error condition and the exact user-visible response |
| **Only happy path** | Misses the majority of real-world usage | Add scenarios for empty state, boundaries, failures, and unauthorized access |
| **Implementation detail in criteria** | Ties acceptance to a technical choice | Describe the observable outcome, not how it's achieved |
| **Passive voice throughout** | Obscures who is doing what | "The user sees X" or "The system sends Y" — name the actor |
