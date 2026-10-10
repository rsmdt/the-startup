---
title: "Team Invitations"
status: draft
version: "1.0"
---

# Solution Design Document

## Approach

A new `invitations` module in the existing Rails monolith stores invites in Postgres and sends join links through the existing mailer. Acceptance reuses the current sign-up flow with an invite token parameter.

- **Cross-cutting:** tokens stored hashed; invite actions logged to the audit log
- **Deployment:** No change

## Read first

| Source | Why it matters |
|--------|----------------|
| `app/controllers/registrations_controller.rb` | Sign-up flow being extended |
| `app/mailers/application_mailer.rb` | Mailer conventions |

## Boundaries

- **Must preserve:** existing sign-up without an invite
- **Can modify:** registrations controller
- **Must not touch:** billing seat counting

## Components

| Component | Responsibility | Owns |
|-----------|----------------|------|
| InvitationsController | Create invites; reject existing members | FR-1 |
| InviteAcceptance | Validate token and expiry; attach the new user to the workspace | FR-2 |
| InvitationMailer | Send the join-link email | FR-1 delivery |

## Interfaces & data

| Name | Direction | Type / format | Auth | Contract |
|------|-----------|---------------|------|----------|
| `POST /workspaces/:id/invitations` | inbound | HTTPS / JSON | admin session | `{email}` → `201` or `422 already_member` |
| `GET /invitations/:token` | inbound | HTTPS / HTML | none | join page or expired page |

```yaml
Table: invitations (NEW)
  id: bigint primary key
  workspace_id: foreign key → workspaces.id
  email: citext not null
  token_digest: text not null unique
  expires_at: timestamptz not null
  accepted_at: timestamptz null
```

## Runtime

[NEEDS CLARIFICATION: primary invite → accept flow and error handling, pending ADR-1 and ADR-2]

## Decisions

> [!IMPORTANT]
> **Decision (ADR-1): Invite token format → Random token stored as SHA-256 digest**
>
> | Option | Trade-off / why (not) |
> |--------|-----------------------|
> | **Random token, hashed at rest** | Revocable; DB leak exposes no usable links |
> | Signed JWT | Stateless, but cannot be revoked before expiry |
>
> **Rationale:** Single-use and revocation need server state · **Accepts:** one DB lookup per link open
>
> - [ ] User confirmed

> [!IMPORTANT]
> **Decision (ADR-2): Expired invite handling → Show expired page with "ask for a new invite"**
>
> | Option | Trade-off / why (not) |
> |--------|-----------------------|
> | **Expired page** | Clear to the invitee; admin re-invites |
> | Auto-resend on open | Lets anyone holding an old link trigger email |
>
> **Rationale:** Keeps admins in control of who is invited · **Accepts:** an extra admin action
>
> - [ ] User confirmed

## Behavior

[NEEDS CLARIFICATION: EARS statements per AC, after ADRs are confirmed]

## Risks & gotchas

[NEEDS CLARIFICATION]

## Open questions

None — all decisions settled.
