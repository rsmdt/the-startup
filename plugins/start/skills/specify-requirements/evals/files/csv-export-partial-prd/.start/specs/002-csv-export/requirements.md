---
title: "CSV Export"
status: draft
version: "1.0"
---

# Product Requirements Document

## Problem

Account managers copy report tables by hand into spreadsheets every week for customer reviews, about two hours each, and the copies often contain mistakes.

## Outcome

Account managers download any report as a CSV that opens correctly in Excel and Google Sheets, in one click.

## Who

| Role | Relationship | Goal or concern today | What success looks like for them |
|------|--------------|-----------------------|----------------------------------|
| Account manager | uses it | Copies tables by hand | Downloads a CSV in one click |
| Data protection officer | affected | Exports could leak customer PII | Exports respect existing report permissions |

## Features

### Must have

#### FR-1: Export a report as CSV
- **User story:** As an account manager, I want to export the report I am viewing as CSV so that I can share it in a spreadsheet.
- **Acceptance criteria:**
  - [ ] **AC-1.1** — Given a report with applied filters, When the user clicks Export CSV, Then a CSV with exactly the filtered rows and visible columns downloads
  - [ ] **AC-1.2** — Given a report the user cannot view, When an export is requested, Then it is refused
  - [ ] **AC-1.3** — [NEEDS CLARIFICATION: What happens when a report exceeds the maximum export size — truncate, refuse, or email the file later?]

### Should have

None.

### Won't have (this phase)

- Excel (.xlsx) format
- Scheduled exports

## Quality gates

| ID | Concern | Requirement (user-observable) | Gate (measurable) |
|----|---------|-------------------------------|-------------------|
| NFR-1 | Performance | Export starts promptly | Download begins < 5 s for reports ≤ 10,000 rows |
| NFR-2 | Compatibility | File opens without garbled characters | UTF-8 with BOM; non-ASCII names render in Excel and Sheets |

## Data

| ID | Data item | Purpose | Retention | Classification |
|----|-----------|---------|-----------|----------------|
| DR-1 | Export audit entry (user, report, time) | Detect bulk data extraction | [NEEDS CLARIFICATION: How long are export audit entries kept?] | internal |

- **User rights:** Audit entries are visible to the user's admin on request.

## Constraints & dependencies

- [NEEDS CLARIFICATION: Is there a contractual or legal limit on exporting customer contact fields?]

## Glossary

| Term | Definition |
|------|------------|
| Report | A saved, filterable table view in the analytics area |

## Open questions

None — all decisions settled.
