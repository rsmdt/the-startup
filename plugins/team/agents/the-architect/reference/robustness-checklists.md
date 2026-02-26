# Robustness Review Checklists

Detailed checklists for complexity and concurrency review. Load when the review-robustness agent needs specific patterns to evaluate.

---

## Complexity: Abstraction Challenge Table

When you see a new abstraction, challenge it. First match wins.

| If You See | Ask | Expected Justification |
|------------|-----|----------------------|
| New interface | "How many implementations exist TODAY?" | 2+ concrete implementations |
| Factory pattern | "Is there more than one product RIGHT NOW?" | Multiple products in use |
| Abstract class | "What behavior is actually shared?" | Concrete shared methods |
| Generic type parameter | "What concrete types are used TODAY?" | 2+ distinct type usages |
| Configuration option | "Has anyone ever changed this from default?" | Evidence of variation |
| Event/callback system | "Could a direct function call work?" | Multiple listeners needed |
| Microservice extraction | "Does this NEED to scale independently?" | Different scaling profile proven |

## Complexity: Code-Level Simplification

- [ ] Functions under 20 lines? If not, WHY?
- [ ] Nesting under 3 levels? Demand guard clauses and early returns
- [ ] No flag variables? Replace with early returns
- [ ] Positive conditionals only? No `if (!notReady)` double negatives
- [ ] Complex expressions named? `const isEligible = x && y && z`
- [ ] No dead code? Unused variables, unreachable branches removed
- [ ] No commented-out code? That's what version control is for

## Complexity: Architecture-Level Simplification

- [ ] Every abstraction justified by CURRENT need (not future speculation)?
- [ ] No pass-through layers? (method just calls another method)
- [ ] No over-engineering? (factory for single implementation)
- [ ] No premature generics? (`Repository<T>` with only one T)
- [ ] Dependencies proportional to functionality?
- [ ] Layer count justified? Can any layer be collapsed?

## Complexity: Anti-Pattern Detection

- [ ] No Lasagna Code? (too many thin layers)
- [ ] No Interface Bloat? (interfaces with unused methods)
- [ ] No Inheritance Addiction? (> 2 levels of inheritance)
- [ ] No Callback Hell? (use async/await)
- [ ] No Ternary Chains? (`a ? b : c ? d : e` → use if/else or switch)
- [ ] No Regex Golf? (unreadable regex → multiple simple checks)
- [ ] No Metaprogramming? (when direct code works fine)

---

## Concurrency: Race Conditions

- [ ] Shared state protected by synchronization?
- [ ] Check-then-act operations atomic? (no TOCTOU vulnerabilities)
- [ ] Compound operations properly locked?
- [ ] Read AND write operations protected? (not just writes)
- [ ] Loop variables captured correctly in closures?
- [ ] Lazy initialization thread-safe?

## Concurrency: Async/Await Patterns

- [ ] All promises awaited or intentionally fire-and-forget?
- [ ] Promise.all used for independent operations?
- [ ] No await in loops when Promise.all would work?
- [ ] Proper error handling for async operations?
- [ ] Async cleanup in finally blocks or using patterns?
- [ ] No mixing callbacks and promises inconsistently?

## Concurrency: Deadlock Prevention

- [ ] Consistent lock ordering maintained?
- [ ] No nested locks that could deadlock?
- [ ] Timeouts on blocking operations?
- [ ] No circular wait conditions?
- [ ] Resources acquired in consistent order?

## Concurrency: Resource Management

- [ ] Async resources properly cleaned up?
- [ ] Event listeners removed when no longer needed?
- [ ] Subscriptions unsubscribed on teardown?
- [ ] Connection pools configured with limits?
- [ ] Timeouts set on external calls?
- [ ] Graceful shutdown handles in-flight operations?

## Concurrency: Database

- [ ] Appropriate transaction isolation level?
- [ ] Optimistic locking where applicable?
- [ ] No long-running transactions holding locks?
- [ ] Batch operations instead of row-by-row?
- [ ] Connection returned to pool promptly?

## Concurrency: Event Handling

- [ ] Event handlers idempotent?
- [ ] No duplicate event processing?
- [ ] Event ordering handled correctly?
- [ ] Backpressure handled for high-volume events?
- [ ] Dead letter queues for failed events?

---

## Concurrency: Common Patterns to Flag

| Pattern | Issue | Fix |
|---------|-------|-----|
| `if (cache[key]) return cache[key]` | TOCTOU race | Use atomic get-or-set |
| `await` inside `forEach` | Sequential, not parallel | Use `Promise.all` with `map` |
| `async () => { fetch(...) }` | Unhandled promise | Add error handling or await |
| Shared mutable object | Race condition | Immutable or synchronized |
| `setTimeout` without cleanup | Memory leak | Store and clear timeout ID |
