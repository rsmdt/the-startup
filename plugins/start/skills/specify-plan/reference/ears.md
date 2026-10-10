# EARS — Easy Approach to Requirements Syntax

EARS gives each behavior statement a clear trigger and an unambiguous required response. Use it in each slice task's **Behavior — EARS** section (part of its code brief) to express the technical requirements drawn from the solution design. Each statement is testable: a validation scenario should be able to exercise its trigger and check its response.

## The five patterns

| Pattern | Template | Use for |
|---------|----------|---------|
| **Ubiquitous** | The system shall `<response>`. | An always-true requirement with no trigger. |
| **Event-driven** | **When** `<trigger>`, the system shall `<response>`. | A response to a discrete event or input. |
| **State-driven** | **While** `<state>`, the system shall `<response>`. | Behavior that holds throughout a state. |
| **Unwanted behavior** | **If** `<condition>`, **then** the system shall `<response>`. | Error handling, validation, failure modes. |
| **Optional feature** | **Where** `<feature is present>`, the system shall `<response>`. | Behavior that exists only when a feature/config is enabled. |

Patterns compose: *"While the cache is warm, when a lookup misses, the system shall fetch from the origin and populate the cache."*

## Writing good EARS

- **One requirement per statement.** Split compound behaviors.
- **Name a concrete response.** "the system shall return HTTP 409 with an `error.code` of `conflict`" beats "the system shall handle the conflict".
- **Ground the trigger.** Reference real events, endpoints, states, and fields from the solution design — not hypothetical ones.
- **Derive, don't invent.** Every EARS statement must trace to technical detail in solution.md. The plan expresses the solution's behavior; it does not add new behavior.
- **Keep it observable where possible.** EARS behavior that surfaces through an external interface can be paired directly with a validation scenario.

## Examples

```
- The system shall persist every accepted order to the orders table.            (ubiquitous)
- When a client POSTs /api/v1/orders with a valid body, the system shall
  return 201 with the created order's id.                                       (event-driven)
- While an order is in `pending_payment`, the system shall reject edits to its
  line items with HTTP 409.                                                     (state-driven)
- If the payment provider returns a timeout, then the system shall mark the
  order `payment_failed` and emit an `order.payment_failed` event.             (unwanted behavior)
- Where the loyalty module is enabled, the system shall apply tier discounts
  before computing tax.                                                         (optional feature)
```

## EARS vs. acceptance scenarios

EARS states the **required behavior** (drawn from the solution). An **acceptance scenario** is the holdout check that the behavior is actually present, exercised through an external interface. They pair up: most P0/P1 scenarios verify one EARS statement. EARS lives in the slice task's code brief; scenarios live as separate holdout **scenario child tasks** under each slice.
