# Decomposition Guide

How to break a solution design into factory-sized units.

---

## Unit Sizing Principle

A unit is the right size when a single code agent can implement it in one session:
- **Too small**: A single function or a config change. Not worth the factory overhead.
- **Right size**: An endpoint + its data model + its tests. A self-contained feature slice.
- **Too large**: An entire service or multiple unrelated features. Should be split.

## Identifying Unit Boundaries

Start from the SDD's building blocks (components, interfaces, data models). Natural boundaries:

1. **Vertical slices**: One API endpoint with its controller, service, model, and tests = one unit.
2. **Data model + consumers**: A new entity and all code that reads/writes it = one unit.
3. **Cross-cutting concern**: A middleware, interceptor, or shared utility = one unit.
4. **Integration point**: Connection to an external service = one unit.

## Dependency Rules

- Dependencies must be explicit: `dependencies: [dm1]` means "dm1 must complete before this unit starts."
- Minimize dependencies. If unit A and unit B share a data model, make the data model its own unit (dm1) and both depend on it.
- No circular dependencies. If A depends on B and B depends on A, merge them into one unit.
- Independent units (no shared dependencies) can run in parallel.

## Unit IDs

Short alphanumeric, position-independent:
- `dm1` — data models
- `ve1` — validation endpoint
- `rl1` — rate limiting
- `si1` — new unit inserted later

Adding a unit between `ve1` and `rl1` is just a new ID with a dependency declaration. No renumbering.

## Coverage Matrix

After decomposition, verify:
- Every requirement in requirements.md maps to at least one unit
- No requirement maps to zero units (gap)
- If a requirement maps to multiple units, clarify which unit owns the primary implementation
