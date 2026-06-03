# Deepening

How to deepen a cluster of shallow modules safely, given its dependencies. Assumes the vocabulary in language.md — **module**, **interface**, **seam**, **adapter**.

## Dependency categories

When assessing a candidate for deepening, classify its dependencies. The category determines how the deepened module is tested across its seam, and which tag the candidate card carries.

### 1. In-process

Pure computation, in-memory state, no I/O. Always deepenable — merge the modules and test through the new interface directly. No adapter needed.

### 2. Local-substitutable

Dependencies that have local test stand-ins (PGLite for Postgres, an in-memory filesystem). Deepenable if the stand-in exists. The deepened module is tested with the stand-in running in the test suite. The seam is internal; no port at the module's external interface.

### 3. Remote but owned (ports & adapters)

Your own services across a network boundary (microservices, internal APIs). Define a **port** (interface) at the seam. The deep module owns the logic; the transport is injected as an **adapter**. Tests use an in-memory adapter. Production uses an HTTP/gRPC/queue adapter.

Recommendation shape: *"Define a port at the seam, implement an HTTP adapter for production and an in-memory adapter for testing, so the logic sits in one deep module even though it's deployed across a network."*

### 4. True external (mock)

Third-party services (Stripe, Twilio, etc.) you don't control. The deepened module takes the external dependency as an injected port; tests provide a mock adapter.

## Seam discipline

- **One adapter means a hypothetical seam. Two adapters means a real one.** Don't introduce a port unless at least two adapters are justified (typically production + test). A single-adapter seam is just indirection.
- **Internal seams vs external seams.** A deep module can have internal seams (private to its implementation, used by its own tests) as well as the external seam at its interface. Don't expose internal seams through the interface just because tests use them.

## Testing strategy: replace, don't layer

- Old unit tests on shallow modules become waste once tests at the deepened module's interface exist — delete them.
- Write new tests at the deepened module's interface. The **interface is the test surface**.
- Tests assert on observable outcomes through the interface, not internal state.
- Tests should survive internal refactors — they describe behaviour, not implementation. If a test has to change when the implementation changes, it's testing past the interface.

## Domain glossary (CONTEXT.md)

This repo does not require a `CONTEXT.md`. When the grilling loop names a deepened module after a concept, capture it so future reviews share the vocabulary. Create the file lazily at the repo root if it doesn't exist:

```markdown
# Context — domain glossary

## <Term>
One paragraph in domain language: what it is, its invariants, what it is not.
```

Add a term when a deepened module is named after a concept not already in the glossary; sharpen a term in place when the conversation clarifies a fuzzy one.

## Architecture decision records (docs/adr/)

When the user rejects a candidate for a load-bearing reason — one a future explorer would need in order to avoid re-suggesting the same deepening — offer to record it. Write to `docs/adr/NNNN-<slug>.md`, creating the directory lazily:

```markdown
# ADR-NNNN: <decision title>

- Status: Accepted
- Date: YYYY-MM-DD

## Context
What friction prompted considering the deepening.

## Decision
What was decided (e.g., "keep these as separate modules") and why.

## Consequences
What this rules in and out for future architecture reviews.
```

Skip the ADR for ephemeral reasons ("not worth it right now") and self-evident ones — record only decisions worth not re-litigating.
