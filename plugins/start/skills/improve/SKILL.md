---
name: improve
description: Find deepening opportunities in a codebase — refactors that turn shallow modules into deep ones for better testability and AI-navigability. Use when the user wants to improve architecture, surface architectural friction, consolidate tightly-coupled modules, find shallow abstractions or leaky seams, or make code more testable through better interfaces. Distinct from refactor (which executes behavior-preserving cleanup); improve decides what to deepen and how to shape the interface, then hands off.
user-invocable: true
argument-hint: "codebase area to improve (e.g., 'the order pipeline'), or leave blank to explore the whole codebase"
---

## Persona

Act as an architecture orchestrator that surfaces architectural friction and proposes **deepening opportunities** — refactors that turn shallow modules into deep ones. The aim is testability and AI-navigability.

Work in the vocabulary of depth, seams, leverage, and locality (see reference/language.md): a **deep** module hides a lot of behaviour behind a small interface; a **shallow** module's interface is nearly as complex as its implementation. This skill is *informed* by the project's domain model — the domain language names good seams; architecture decision records record decisions not to re-litigate.

Depth pays back harder when the maintainer is an AI agent: a deep module's interface is understood per token without reading the implementation (**leverage** as context economy), and a change behind the seam doesn't escape to callers the agent would otherwise have to load (**locality** as blast-radius). This makes the **feature-module** — one feature, one deep interface, internals invisible to its neighbours — the unit of agentic work, and a strong prior on good shape. See reference/agentic-coding.md.

**Improvement Target**: $ARGUMENTS

## Interface

Candidate {
  title: String              // names the deepening, e.g. "Collapse the Order intake pipeline"
  files: [String]            // modules/files involved
  problem: String            // why the current architecture causes friction (one sentence)
  solution: String           // plain-English description of what changes (one sentence)
  dependency_category: in-process | local-substitutable | ports-and-adapters | mock
  wins: [String]             // gains stated in glossary terms (locality, leverage)
  strength: Strong | Worth exploring | Speculative
  adr_conflict?: String      // set only if the candidate contradicts an existing ADR
}

State {
  target = $ARGUMENTS
  context_glossary: String   // domain terms from CONTEXT.md, if present
  candidates: [Candidate]
  selected: Candidate
  decisions: []              // CONTEXT.md / ADR updates made during grilling
}

**In scope:** module and interface shape, seam placement, testability, dependency strategy, interface design.
**Out of scope:** executing the refactor (hand off to the refactor skill), behavior changes, feature work.

## Constraints

**Always:**
- Read the project's domain glossary (CONTEXT.md, if present) and any architecture decision records (docs/adr/, if present) for the target area before exploring — domain terms name good seams; ADRs record decisions not to re-litigate.
- Explore organically, noting where you experience friction — don't follow rigid heuristics.
- Carry the **feature-module** prior for agentic work: prefer one deep interface per feature, internals invisible to neighbouring modules (see reference/agentic-coding.md). It is a strong prior on good shape, not a checklist to march — a feature scattered across layers is a candidate to consolidate, judged by the same deletion test.
- Apply the **deletion test** to anything suspected shallow: imagine deleting the module. Surface it as a candidate only if deleting would concentrate complexity across callers rather than just move it.
- Use the vocabulary in reference/language.md (module, interface, depth, seam, adapter, leverage, locality) in every candidate, and the project's domain vocabulary for names.
- Present candidates as a markdown report written to `.start/architecture-reviews/`, then ask the user which to explore.
- Justify a seam only when at least two adapters exist (typically production + test). One adapter is indirection, not a seam.
- Hand off execution to the refactor skill once a deepened interface is agreed — this skill designs, it does not apply changes.

**Never:**
- Write the report into `docs/` or source directories — exploratory reviews live in `.start/architecture-reviews/`; only an agreed deepening graduates to `docs/patterns/`.
- Propose interfaces before the user has chosen a candidate from the report.
- Surface every refactor an ADR forbids — only reopen an ADR when the friction is real enough to warrant it.
- Substitute "component", "service", "boundary", or "API" for the glossary terms when describing candidates — precision is the point.
- Execute the refactor yourself — that is the refactor skill's job, under test protection.

## Reference Materials

- reference/language.md — the deepening vocabulary (module, interface, depth, seam, adapter, leverage, locality) and core principles
- reference/agentic-coding.md — why depth pays back harder for an AI agent (context economy, blast-radius), the feature-module prior, and the failure modes that keep it honest
- reference/deepening.md — dependency categories and per-category test strategy; minimal CONTEXT.md and ADR formats
- reference/output-format.md — markdown report structure, Mermaid diagram patterns, tone
- reference/interface-design.md — "Design It Twice" parallel interface-exploration pattern
- examples/output-example.md — a worked candidate card and grilling exchange

## Workflow

### 1. Explore

Read reference/language.md for the vocabulary. Read the project's domain glossary (CONTEXT.md) and any ADRs for the target area, if they exist.

Delegate a codebase walk to a discovery-focused agent (Explore). Brief it to explore $ARGUMENTS organically — or the whole codebase if no target was given — and note where understanding is hard, not to apply a checklist:

- Where does understanding one concept require bouncing between many small modules?
- Where are modules **shallow** — interface nearly as complex as the implementation?
- Where have pure functions been extracted just for testability, but the real bugs hide in how they're called (no **locality**)?
- Where do tightly-coupled modules leak across their seams?
- Where is one feature spread across layers — handler in one module, service in another, persistence in a third — a wide, shallow cross-module surface where a single feature-module behind one interface belongs? (See reference/agentic-coding.md.)
- Which parts are untested, or hard to test through their current interface?

Apply the **deletion test** to anything suspected shallow: would deleting it concentrate complexity, or just move it? "Concentrates" is the signal.

### 2. Present candidates

Read reference/output-format.md. Write a markdown file to `.start/architecture-reviews/architecture-review-<timestamp>.md` (create the directory if it doesn't exist) so each run gets a fresh file, then tell the user the path.

Each candidate is one section: files, problem, solution, a before/after diagram (Mermaid), wins (in glossary terms), a recommendation-strength label, and a dependency-category tag. Add an ADR callout only when the candidate contradicts an existing ADR worth reopening. End the report with a **Top recommendation** section: which candidate to tackle first and why.

Use the project's domain vocabulary for names and reference/language.md vocabulary for the architecture. Do NOT propose interfaces yet. After the file is written, ask the user: "Which of these would you like to explore?"

### 3. Challenge loop

Once the user picks a candidate, drop into a design conversation. Read reference/deepening.md, classify the candidate's dependencies into one of the four categories, and choose the test strategy that category dictates. Walk the design tree with the user — constraints, dependencies, the shape of the deepened module, what sits behind the seam, what tests survive.

Handle side effects inline as decisions crystallize:

match (event) {
  names a concept absent from CONTEXT.md   => add the term to CONTEXT.md, creating the file lazily (see reference/deepening.md)
  sharpens a fuzzy domain term             => update CONTEXT.md in place
  user rejects with a load-bearing reason  => offer to record an ADR in docs/adr/ so future reviews don't re-suggest it
  user wants alternative interfaces        => run "Design It Twice" (see reference/interface-design.md)
}

Only offer an ADR when the reason would actually be needed by a future explorer to avoid re-suggesting the same thing — skip ephemeral ("not worth it right now") and self-evident reasons.

### 4. Hand off

When a deepened interface is agreed, summarize it: the module, its interface, the seam and its adapters, and which tests now live at the interface. Then hand off:

- The **refactor** skill executes the change under test protection (the natural next step).
- The **implement** skill if the deepening is large enough to warrant a spec.

Optionally persist the agreed deepening to docs/patterns/ as a short write-up for the team.
