# Deepening for Agentic Coding

Why depth pays back harder when the maintainer is an AI agent — and the organizing prior that follows. Vocabulary from language.md: **module**, **interface**, **seam**, **depth**, **leverage**, **locality**, **adapter**.

## Why depth matters more for an agent

A human navigates a codebase as a graph — jumping, holding partial knowledge, filling gaps by inference. An agent works against a finite, degrading context budget and cannot cheaply verify what it hasn't loaded. Three properties this skill already names stop being niceties and become load-bearing:

- **Leverage is paid in tokens.** A deep module lets an agent load the **interface** and skip the **implementation** — much behaviour understood per token spent. A shallow module's interface is nearly as wide as its body, so the agent must read the internals to use it at all. Depth is context economy.
- **Locality is blast-radius the agent can't see.** A human verifies a change by reviewing callers; an agent must load every caller into context to know its edit is safe. When callers cross only the **interface**, an internal change stays local and the agent need not read the rest of the tree to be sure. Every leak across the seam is one more caller the agent has to find and hold.
- **Navigability is directory jumps.** A feature whose logic is spread across layers — handler in one module, service in another, persistence in a third — costs one context load per hop. A feature that lives behind one **interface** costs one.

## The organizing prior: feature-modules

This skill stays exploratory — it finds friction, it does not march a checklist (see SKILL.md). But for agentic work it carries a strong prior about good shape:

**Prefer feature-modules: each feature is one deep module with a single interface; its internals are invisible to the modules around it.**

This is the depth principle applied to organization. Put a feature behind one interface and only that interface need be visible — internals stay private, the **seam** is a real place, and an agent assigned the feature loads one entry point and works outward without spilling into its neighbours. Spread the same feature across layers and every layer must be visible to the others: the cross-module surface is wide and shallow, and an agent must hold all of it to change any of it. (This is the package-by-feature / package-by-layer split, in depth terms — named only to locate it.)

So the skill gains a candidate type: **a feature scattered across layers is a shallow, leaky cross-module surface.** Apply the deletion test — delete any one layer and its complexity reappears in the others, so the pieces were earning their keep as logic but not as separate **modules**. Consolidate them into one deep module behind one interface.

## The agentic seam test

The interface is the **sole surface** a caller crosses. Callers depend only on it; reaching past it into internals is **leakage**. The test, in agentic terms:

> An agent should be able to change a module's internals without any other module noticing. If it can't, the seam is leaking or in the wrong place.

This is the deletion test's cousin: a real seam contains change; a leaky one lets it escape to callers the agent never loaded.

## Keep the prior honest

A prior is not a mandate. The failure modes this skill already guards still apply, plus two the organizing opinion invites:

- **Classitis (premature consolidation).** A hundred feature-modules each with a one-line interface and nothing behind it is shallowness wearing a feature's name — it inverts depth. Consolidate only where the deletion test concentrates complexity.
- **God-module.** A feature-module that absorbs unrelated concerns is deep to callers but incoherent inside. Use **internal seams** (see deepening.md) to keep the implementation comprehensible without widening the interface.
- **Over-deep opacity.** An interface so terse a caller cannot reason about failure trades simplicity for un-debuggability. Depth includes the interface's invariants, ordering, and error modes (see language.md) — an agent leans entirely on that contract, so it must be precise, not merely small.
- **One adapter is not a seam** (see language.md). A feature-module with a single implementation does not need a port — that is indirection, not depth.
