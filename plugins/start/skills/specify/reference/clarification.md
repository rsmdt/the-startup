# Clarification Protocol

Loaded by the specify skill's Clarify Intent phase. The governing stance: before answering, say what you need to know to answer well, and point out the assumptions you would otherwise make. Writing from half-clarified intent produces confident, wrong documents.

---

## Deriving Unknowns

Sweep the brief against these lenses; every gap becomes an unknown with a stake:

| Lens | Ask yourself |
|------|--------------|
| WHAT | Which capabilities are actually requested vs. implied? Where does scope end? |
| WHY | What problem motivates this, and how would the requester know it's solved? |
| WHO | Which user types touch this? Who is explicitly *not* served? |
| WHEN | What are the success metrics and acceptance conditions? |
| Constraints | Budget, deadline, compliance, tech stack mandates, backwards compatibility? |
| Integrations | What existing systems, data, or contracts does this touch? |
| Non-functionals | Load, latency, availability, security expectations? |

An unknown without a stake is not worth asking — if you cannot say what changes downstream depending on the answer, drop it or assume it.

## Criticality

| Level | Meaning | Handling |
|-------|---------|----------|
| **Blocking** | The answer changes the shape of the requirements or architecture | Must be resolved before any document work |
| **Shaping** | The answer changes priorities, ordering, or details within a settled shape | Ask if budget allows; otherwise assume with a visible default |
| **Cosmetic** | Naming, formatting, presentation | Never ask — assume, and surface in the ledger |

## Routing

Ask the user only what research cannot answer.

| Route to **User** | Route to **Research** |
|---|---|
| Intent — the goal behind the ask | What the codebase currently does |
| Priorities and trade-off preferences | Existing patterns, conventions, contracts |
| Scope boundaries — in vs. out | Vendor/API capabilities and limits |
| Taste — UX, tone, style choices | Domain facts, standards, prior art |
| Success criteria only they can define | Feasibility and effort signals |

Research-routed unknowns seed the research perspectives. When findings return, they are **not** accepted silently — see Research Confirmation below.

## Question Quality Bar

- Every question names its stake: "X or Y? — this decides whether requirements include Z."
- One decision per question — no stacked multi-part questions.
- Batched, at most five per round, ordered by leverage (Blocking first).
- Each question carries a proposed default, so the user can answer with one word or accept all defaults in a single action.
- No question the codebase, git history, or a web search could answer.

## The Assumption Ledger

Present before any questions, maintain throughout the run, and record in the spec README (see output-format.md):

| Topic | Default (what I'd otherwise assume) | Impact if wrong | Status |
|-------|--------------------------------------|-----------------|--------|
| Auth model | Reuse existing session auth | New auth flow = rework of every AC | Proposed |
| Notification channels | In-app only, no email | Email adds vendor + template scope | Standing |

Status lifecycle: **Proposed** (surfaced, unanswered) → **Confirmed** (user validated the default) / **Corrected** (user gave a different answer) / **Standing** (user accepted the default without validating it — carried into documents as a visible `[ASSUMPTION]` entry, never erased).

## Research Confirmation

Findings return as claims, not facts. Present each material claim with: the claim, its evidence (file, doc, or source), and a proposed disposition. The user chooses:

- **Confirm** — the finding closes its unknown and may enter documents.
- **Correct** — the user's version supersedes; record the correction.
- **Discard** — the finding never enters a document.
- **Defer** — becomes a Standing `[ASSUMPTION]` with the finding as its default.

A finding no user has dispositioned is Pending and must not appear in requirements, solution, or plan.

## MECE Tie-In

Clarification feeds the exhaustiveness half of MECE: every Standing assumption and every closed unknown must trace to a requirement, an acceptance criterion, or an explicit out-of-scope statement. The mechanics of MECE verification for personas, journeys, features, and acceptance criteria are owned by the specify-requirements skill (reference/focus-and-mece.md) — this protocol supplies its inputs, not its rules.
