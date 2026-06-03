# Decomposition Guide

How to break a solution design into factory-sized units.

---

## Slice Vertically, Not Horizontally

Every unit is a **vertical slice** (tracer-bullet): a thin, end-to-end path through whatever layers a behavior needs (data → logic → interface → tests) that delivers one user-observable capability. Each unit's holdout scenarios must be able to pass when only that unit (plus its declared dependencies) is built.

❌ **Do NOT slice by layer.** A "data models" unit, a "service layer" unit, and a "middleware" unit are horizontal layers: none is independently runnable against an observable behavior, and the feature only works once they integrate. Splitting one endpoint into core-logic + sanitization + rate-limit units is horizontal slicing in disguise — the sanitization and rate-limit units can't pass scenarios without the endpoint.

✅ **Slice by capability.** `validate-a-query` and `explain-a-query` are vertical: each owns its endpoint end-to-end and can be evaluated alone.

**Unit independence test:** *Could this unit's scenarios pass if only this unit (and its dependencies) were built?* If no, it's a layer, not a slice — fold it into the capability it serves.

---

## Unit Sizing Principle

A unit is the right size when a single code agent can implement it in one session:
- **Too small**: A single function or a config change. Not worth the factory overhead.
- **Right size**: An endpoint + its data model + its tests. A self-contained feature slice.
- **Too large**: An entire service or multiple unrelated features. Should be split.

## Identifying Unit Boundaries

Start from the **requirements' user-observable capabilities**, then map each to the SDD building blocks it touches. A unit is one capability, sliced vertically:

1. **Capability slice**: One API endpoint (or user action) with its controller, service, the data model it needs, and its tests = one unit.
2. **Second capability that reuses the first**: A new user-observable behavior that builds on an existing slice = its own unit, depending on the first.
3. **Integration point**: A connection to an external service, exercised through an observable behavior, = one unit.

**Cross-cutting concerns (auth, sanitization, rate limiting, logging) are NOT their own units.** Fold each into the capability slice it protects — its scenarios test that the protected behavior is safe. Carve out a separate unit only if the concern is itself an independently user-observable capability (e.g. an auth *endpoint*), in which case it is a vertical slice in its own right.

## Dependency Rules

- Dependencies must be explicit: `dependencies: [vq1]` means "vq1 must complete before this unit starts."
- **Shared foundation belongs to the first slice, not its own unit.** When several units need the same entity, schema, or scaffolding, the **first** unit that needs it builds that foundation as part of its own end-to-end slice; later units depend on that unit and **extend** the foundation. Never extract shared foundation into a standalone data-only unit — a dependency edge must point to a real, independently-validatable slice, never to a horizontal layer.
- No circular dependencies. If A depends on B and B depends on A, merge them into one unit.
- Independent units (units that share no dependency, or depend only on an already-built foundation slice) can run in parallel.

## Unit IDs

Short alphanumeric, position-independent. Name them after the **capability**, not a layer:
- `vq1` — validate a query
- `eq1` — explain a query
- `oq1` — optimize a query

Adding a capability later is just a new ID with a dependency declaration. No renumbering. (Avoid IDs like `dm1`/"data models" or `rl1`/"rate limiting" — those name layers, not capabilities, and signal horizontal slicing.)

## Coverage Matrix

After decomposition, verify:
- Every requirement in requirements.md maps to at least one unit
- No requirement maps to zero units (gap)
- If a requirement maps to multiple units, clarify which unit owns the primary implementation
