---
title: "Team Invitations"
status: completed
version: "1.0"
---

# Product Requirements Document

## Problem

Workspace admins add teammates by sharing one password; there is no way to invite a person to their own account.

## Outcome

Admins invite teammates by email; invitees join with their own account in under a minute.

## Who

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| Workspace admin | uses it | Shares credentials to add people | Invites by email |
| Invitee | uses it | Has no account of their own | Joins from the invite link |

## Features

### Must have

#### FR-1: Invite by email
- **User story:** As a workspace admin, I want to invite a teammate by email so that they get their own account.
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given an admin, When they invite a new email address, Then an invite email with a join link is sent
  - [ ] **AC-1.2** — Given an email that is already a member, When the admin invites it, Then the invite is rejected with "already a member"

#### FR-2: Accept an invite
- **User story:** As an invitee, I want to accept an invite so that I can join the workspace.
- **Acceptance criteria:**
  - [ ] **AC-2.1** — Given a valid invite link, When the invitee signs up, Then they join the workspace as a member
  - [ ] **AC-2.2** — Given an invite older than 7 days, When the link is opened, Then it is shown as expired

### Should have

None.

### Won't have (this phase)

- Bulk invites, role selection on invite

## Quality gates

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Security | Invite links cannot be guessed | Tokens ≥ 128 bits of entropy; single use |
| NFR-2 | Performance | Invite email arrives promptly | 95% delivered < 60 s |

## Constraints & dependencies

- Uses the existing transactional email provider.

## Glossary

| Term | Definition |
|------|------------|
| Invite | A single-use, expiring offer to join one workspace |

## Open questions

None — all decisions settled.
