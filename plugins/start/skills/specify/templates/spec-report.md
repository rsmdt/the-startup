---
title: "{{SPEC_ID}}-{{SPEC_NAME}} — Specification Report"
status: "{{STATUS}}"
readiness: "{{READINESS}}"
---

# Specification Report

> [!NOTE]
> Written by the specify skill at Finalize. This is the editor-readable record of *how* the
> specification was clarified: the ledger, research dispositions, document status, and decisions.
> The documents themselves (`requirements.md`, `solution.md`, `plan.md`) remain canonical —
> nothing here duplicates their content.

## Readiness

| Field | Value |
|-------|-------|
| Readiness | {{READINESS}} |
| Phase | {{PHASE}} |
| Last updated | {{DATE}} |

> [!IMPORTANT]
> **Readiness capped:** {{CAP_REASON — remove this callout when readiness is uncapped}}

## Clarification Ledger

### Unknowns

What had to be known to specify well. User-routed rows were asked; Research-routed rows were
closed by confirmed findings.

| Question | Stake — what it changes | Route | Criticality | Resolution |
|----------|-------------------------|-------|-------------|------------|
| {{QUESTION}} | {{STAKE}} | User \| Research | Blocking \| Shaping \| Cosmetic | {{RESOLUTION or `open`}} |

### Assumptions

Standing entries appear as `[ASSUMPTION]` markers in the affected documents — surfaced, never erased.

| Topic | Default (otherwise assumed) | Impact if wrong | Status |
|-------|-----------------------------|-----------------|--------|
| {{TOPIC}} | {{DEFAULT}} | {{IMPACT}} | Proposed \| Confirmed \| Corrected \| Standing |

## Research Findings

Every claim required user disposition before entering a document. A `Pending` finding is a
hypothesis and appears in no document.

| Perspective | Claim | Evidence | Disposition |
|-------------|-------|----------|-------------|
| {{PERSPECTIVE}} | {{CLAIM — append "→ correction" when Corrected}} | {{EVIDENCE}} | Confirmed \| Corrected \| Discarded \| Pending |

## Documents

| Document | Status | Notes |
|----------|--------|-------|
| requirements.md | {{STATUS}} | {{NOTES}} |
| solution.md | {{STATUS}} | {{NOTES}} |
| plan.md | {{STATUS}} | {{NOTES}} |

## Decisions Log

| Date | Decision | Rationale |
|------|----------|-----------|
| {{DATE}} | {{DECISION}} | {{RATIONALE}} |

## Next Steps

- [ ] {{NEXT_STEP — e.g. /start:validate NNN}}
- [ ] {{NEXT_STEP — e.g. /start:implement NNN}}
