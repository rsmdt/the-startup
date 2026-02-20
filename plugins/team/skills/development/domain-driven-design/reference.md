# Aggregate Design Guide

Advanced heuristics for defining aggregate boundaries, sizing aggregates correctly, and choosing consistency strategies. Load this when the SKILL.md's high-level rules are insufficient for the decision at hand.

---

## Aggregate Boundaries

Aggregate boundaries exist to protect invariants — the business rules that must always be true. The primary question when drawing a boundary is: **which objects must change together to keep a business rule intact?**

### The Invariant Test

For each proposed aggregate boundary, identify the invariants it protects:

```
Example: Should OrderItem be inside Order or its own aggregate?

Invariant: "An order's total cannot exceed the customer's credit limit"

Test:
- Can we check this rule with only the Order? YES — Order knows its items and total.
- Does changing an OrderItem require validating Order-level rules? YES.
- Can OrderItem exist independently with its own lifecycle? NO — it only exists for an order.

Decision: OrderItem is INSIDE the Order aggregate.
```

```
Example: Should Order be inside Customer or its own aggregate?

Invariant: "A customer cannot have more than 10 active orders simultaneously"

Test:
- Must Order and Customer change together to protect this rule? NO.
- Can we check this rule by querying the count of Orders for a CustomerId? YES.
- Does Customer have a distinct lifecycle from Order? YES.

Decision: Order is a SEPARATE aggregate, referenced from Customer by ID.
The invariant is enforced via a domain service or application-level check before creating an Order.
```

### Aggregate Design Canvas

When designing a new aggregate, work through this canvas:

```
┌─────────────────────────────────────────────────────────────┐
│ AGGREGATE NAME:                                             │
├─────────────────────────────────────────────────────────────┤
│ INVARIANTS (rules that must always hold):                   │
│   1.                                                        │
│   2.                                                        │
├─────────────────────────────────────────────────────────────┤
│ AGGREGATE ROOT (single entry point for all changes):        │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ INSIDE THE BOUNDARY (change together to protect invariants):│
│   Entities:                                                 │
│   Value Objects:                                            │
├─────────────────────────────────────────────────────────────┤
│ OUTSIDE THE BOUNDARY (referenced by ID):                    │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ DOMAIN EVENTS (what this aggregate announces):              │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ CONSISTENCY TYPE:                                           │
│   [ ] Transactional — invariants protected within boundary  │
│   [ ] Eventual — consistency with other aggregates via events│
└─────────────────────────────────────────────────────────────┘
```

### Common Boundary Mistakes

**Grouping by noun instead of invariant**

The mistake is asking "what belongs to an order?" rather than "which objects must change together to enforce order rules?"

```
// WRONG: Grouped by association
class Order {
  customer: Customer;      // Does Customer's email change when Order changes? No.
  payments: Payment[];     // Does Payment share Order's invariants? No.
  shipments: Shipment[];   // Can Shipment change independently? Yes.
  reviews: Review[];       // No shared invariants with Order.
}

// CORRECT: Grouped by invariant
class Order {              // Invariant: item total <= approved budget
  items: OrderItem[];      // Must change with Order to enforce the rule
  discountCode: DiscountCode;  // Applied at order level, affects total
}

// Separate aggregates, referenced by OrderId
class Payment { orderId: OrderId; ... }
class Shipment { orderId: OrderId; ... }
class Review { orderId: OrderId; ... }
```

**Including historical data**

Historical records rarely share invariants with the current state. They grow unbounded and should be separate.

```
// WRONG: History inside aggregate
class Account {
  balance: Money;
  transactions: Transaction[];  // Could be millions — never inspected for invariants
}

// CORRECT: Ledger as separate aggregate
class Account {
  balance: Money;              // Current state
  // Invariant: balance = sum of credits - sum of debits
  // Enforced by debit() and credit() methods, not by inspecting transaction history
}

class Transaction {
  accountId: AccountId;       // Reference by ID
  amount: Money;
  type: TransactionType;
  occurredAt: Date;
}
```

---

## Aggregate Sizing

Start with the smallest possible aggregate — usually a single entity — and expand only when an invariant cannot be protected otherwise.

### Size Signals

**Signals that an aggregate is too large:**

| Signal | What It Means |
|--------|--------------|
| Optimistic lock conflicts on concurrent edits | Multiple users editing unrelated parts simultaneously |
| Loading thousands of rows for a simple operation | Aggregate includes unbounded collections |
| Transactional failures involving unrelated data | Scope is wider than invariants require |
| Slow reconstitution from the repository | Too many child objects |
| Cross-cutting edits by different bounded contexts | Boundary is in the wrong place |

**Signals that an aggregate is too small:**

| Signal | What It Means |
|--------|--------------|
| Business rules scattered across domain services | Invariant spans objects not in the same boundary |
| Application code enforcing consistency across multiple saves | Transaction script replacing domain logic |
| "Eventual" consistency used where the business demands immediate | Boundary was split too aggressively |

### Sizing Heuristics

**1. Default to single-entity aggregates**

Most aggregates start as a single entity. Add children only when a specific invariant requires it.

```
// Start here
class Order {
  id: OrderId;
  status: OrderStatus;
  customerId: CustomerId;     // Reference — not inside
}

// Expand when: "total must be recalculated across all items atomically"
class Order {
  id: OrderId;
  status: OrderStatus;
  customerId: CustomerId;
  items: OrderItem[];         // Now inside — total invariant requires it
}
```

**2. Limit unbounded collections**

If a collection can grow without bound, it must not be inside the aggregate.

```
// BAD: Blog with all comments inside
class BlogPost {
  title: string;
  body: string;
  comments: Comment[];   // Could be 10,000 — all loaded for every post operation
}

// GOOD: Comment is its own aggregate
class BlogPost {
  title: string;
  body: string;
  commentCount: number;  // Denormalized count for display — updated via event
}

class Comment {
  postId: PostId;        // Reference — not embedded in post
  body: string;
  approvedAt: Date | null;
}
```

**3. Prefer eventual consistency at aggregate boundaries**

Business rarely demands immediate consistency across aggregates. Clarify with domain experts whether "immediately" means within the same transaction or "shortly thereafter."

```
Question to ask: "If placing an order and reserving inventory happened within
2 seconds of each other — would that be acceptable to the business?"

If YES → eventual consistency, separate aggregates, domain events
If NO  → investigate why. Usually it's a UI concern, not a true invariant.
```

---

## Consistency Rules

### Transactional Consistency

Use when: **invariants must be true the instant a command completes.**

Scope: within a single aggregate, in a single transaction.

```
Application layer — one aggregate per command:

async function addItemToOrder(command: AddItemCommand): Promise<void> {
  const order = await this.orderRepo.findById(command.orderId);
  order.addItem(command.productId, command.quantity, command.price);
  await this.orderRepo.save(order);   // Atomic: load → mutate → save
}
```

### Eventual Consistency

Use when: **consistency between aggregates can be achieved asynchronously.**

Scope: across aggregates, within or across services.

```
Pattern: publish → subscribe → update own aggregate

// Step 1: Order aggregate raises an event
class Order {
  submit(): void {
    this.status = OrderStatus.Placed;
    this.events.push(new OrderPlaced(this.id, this.items, this.total));
  }
}

// Step 2: Repository publishes after successful save
class PostgresOrderRepository {
  async save(order: Order): Promise<void> {
    await this.db.save(this.decompose(order));
    await this.eventBus.publishAll(order.pullEvents());
  }
}

// Step 3: Inventory aggregate handles event in its own transaction
class InventoryProjection {
  async on(event: OrderPlaced): Promise<void> {
    for (const item of event.items) {
      const stock = await this.stockRepo.findByProduct(item.productId);
      stock.reserve(item.quantity);
      await this.stockRepo.save(stock);   // Own aggregate, own transaction
    }
  }
}
```

### Choosing Between Transactional and Eventual

```
Decision tree:

1. Are the objects part of the same aggregate?
   YES → Transactional (ACID within aggregate boundary)
   NO  → Continue ↓

2. Does the business genuinely require atomicity across these objects?
   YES → Reconsider the aggregate boundary — they may belong together
   NO  → Continue ↓

3. Is latency measured in seconds acceptable?
   YES → Eventual consistency via domain events
   NO  → Investigate. Usually a UI expectation, not a true business requirement.

4. Does failure in the secondary update require rollback of the primary?
   YES → Saga with compensation (not distributed transaction)
   NO  → Eventual consistency, idempotent handlers, retry on failure
```

---

## Eventual Consistency Between Aggregates

### Idempotent Event Handlers

Handlers must be safe to call more than once. Message brokers deliver at-least-once; handlers must not double-apply effects.

```typescript
class InventoryHandler {
  async on(event: OrderPlaced): Promise<void> {
    // Guard: skip if already processed this event
    const processed = await this.processedEvents.contains(event.eventId);
    if (processed) return;

    for (const item of event.items) {
      const stock = await this.stockRepo.findByProduct(item.productId);
      stock.reserve(item.quantity);
      await this.stockRepo.save(stock);
    }

    await this.processedEvents.record(event.eventId);
  }
}
```

### Handling Failures

When a handler fails, the business needs a clear answer: is the failure retryable or fatal?

```
Retryable failures (retry with backoff):
- Network timeout
- Downstream service unavailable
- Optimistic lock conflict

Fatal failures (dead letter queue + alert):
- Invalid event schema — event is malformed
- Business rule violation — inventory does not exist for the product
- Permanent downstream rejection — payment provider blacklisted the account

Pattern:
  try {
    await handler.on(event)
  } catch (RetryableError) {
    await queue.nack(event)          // Return to queue for retry
  } catch (FatalError) {
    await deadLetterQueue.send(event) // Human review required
    await alerting.notify(error)
  }
```

### Eventual Consistency and User Experience

When the UI must reflect a consistent state before eventual handlers complete, use optimistic updates:

```
Pattern: Update read model immediately, correct if event fails

1. User places order → UI shows "Order Placed" immediately
2. OrderPlaced event → InventoryHandler runs asynchronously
3. If inventory is insufficient → OrderFailed event raised
4. UI receives OrderFailed notification → shows error, reverts display

This preserves responsiveness without requiring synchronous cross-aggregate updates.
```

---

## Saga Pattern for Multi-Step Processes

Sagas coordinate a sequence of aggregate updates with explicit compensation when steps fail.

### Choreography-Based Saga

Each aggregate reacts to events and publishes its own. No central coordinator.

```
OrderPlaced
  → InventoryHandler: reserves stock → InventoryReserved
  → InventoryReserved → PaymentHandler: charges card → PaymentCharged
  → PaymentCharged → FulfillmentHandler: ships order → OrderShipped

Compensation chain on PaymentFailed:
  PaymentFailed → InventoryHandler: releases reservation → InventoryReleased
  InventoryReleased → OrderHandler: cancels order → OrderCancelled

Tradeoff: Simple to implement, difficult to trace the full process.
```

### Orchestration-Based Saga

A saga object drives the steps and handles compensation explicitly.

```typescript
class OrderFulfillmentSaga {
  private step: FulfillmentStep = FulfillmentStep.ReserveInventory;

  async handle(event: DomainEvent): Promise<void> {
    switch (this.step) {
      case FulfillmentStep.ReserveInventory:
        if (event instanceof OrderPlaced) {
          await this.inventoryService.reserve(event.items);
          this.step = FulfillmentStep.ChargePayment;
        }
        break;

      case FulfillmentStep.ChargePayment:
        if (event instanceof InventoryReserved) {
          await this.paymentService.charge(event.orderId, event.amount);
          this.step = FulfillmentStep.Ship;
        }
        if (event instanceof InventoryReservationFailed) {
          await this.orderService.cancel(event.orderId, CancelReason.OutOfStock);
          this.step = FulfillmentStep.Compensating;
        }
        break;

      case FulfillmentStep.Ship:
        if (event instanceof PaymentCharged) {
          await this.fulfillmentService.ship(event.orderId);
          this.step = FulfillmentStep.Complete;
        }
        if (event instanceof PaymentFailed) {
          await this.inventoryService.release(event.orderId);
          await this.orderService.cancel(event.orderId, CancelReason.PaymentFailed);
          this.step = FulfillmentStep.Compensating;
        }
        break;
    }
  }
}
```

### When to Use Each

| Approach | Use When |
|----------|----------|
| **Choreography** | Simple 2-3 step processes, teams can coordinate on event contracts |
| **Orchestration** | Complex multi-step processes, compensation logic is non-trivial |
| **Neither** | Steps are fast, can be synchronous, business requires atomicity |

---

## Rules of Thumb

These are practical defaults, not laws. Use them as starting points and adjust based on invariants.

```
1. Default to one entity per aggregate.
   Expand only when you can name the specific invariant that requires it.

2. If you need to update two aggregates simultaneously, question the boundary.
   Either they should be one aggregate, or eventual consistency is acceptable.

3. An aggregate that is never loaded alone is a smell.
   If you always load Order with Customer together, Customer may belong inside Order.
   Or the operation should live in a different bounded context.

4. Prefer many small aggregates over one large one.
   A large aggregate is a coordination bottleneck under concurrent load.

5. A repository method that returns a list of thousands is wrong.
   Either the collection shouldn't be in the aggregate, or you need a read model.

6. Use value objects for anything described by its attributes, not its identity.
   If two instances with the same data are interchangeable, it's a value object.

7. Events describe what happened, not what should happen next.
   "OrderPlaced" is correct. "SendConfirmationEmail" is a command, not an event.

8. Aggregate boundaries are not permanent.
   Model based on current understanding. Refactor when invariants become clearer.
```
