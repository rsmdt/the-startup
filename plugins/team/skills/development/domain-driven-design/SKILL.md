---
name: domain-driven-design
description: Domain-Driven Design tactical and strategic patterns including entities, value objects, aggregates, bounded contexts, and consistency strategies. Use when modeling business domains, designing aggregate boundaries, implementing business rules, or planning data consistency.
---

# Domain-Driven Design Patterns

Patterns for modeling complex business domains with clear boundaries, enforced invariants, and appropriate consistency strategies.

## When to Activate

- Modeling business domains and entities
- Designing aggregate boundaries
- Implementing complex business rules
- Planning data consistency strategies
- Establishing bounded contexts
- Designing domain events and integration

## Strategic Patterns

### Bounded Context

A bounded context defines the boundary within which a domain model applies. The same term can mean different things in different contexts.

```
Example: "Customer" in different contexts

┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│    Sales        │  │    Support      │  │    Billing      │
│    Context      │  │    Context      │  │    Context      │
├─────────────────┤  ├─────────────────┤  ├─────────────────┤
│ Customer:       │  │ Customer:       │  │ Customer:       │
│ - Leads         │  │ - Tickets       │  │ - Invoices      │
│ - Opportunities │  │ - SLA           │  │ - Payment       │
│ - Proposals     │  │ - Satisfaction  │  │ - Credit Limit  │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

#### Context Identification

Ask these questions to find context boundaries:
- Where does the ubiquitous language change?
- Which teams own which concepts?
- Where do integration points naturally occur?
- What could be deployed independently?

### Context Mapping

Define how bounded contexts integrate:

| Pattern | Description | Use When |
|---------|-------------|----------|
| **Shared Kernel** | Shared code between contexts | Close collaboration, same team |
| **Customer-Supplier** | Upstream/downstream relationship | Clear dependency direction |
| **Conformist** | Downstream adopts upstream model | No negotiation power |
| **Anti-Corruption Layer** | Translation layer between models | Protecting domain from external models |
| **Open Host Service** | Published API for integration | Multiple consumers |
| **Published Language** | Shared interchange format | Industry standards exist |

### Ubiquitous Language

The shared vocabulary between developers and domain experts:

```
Building Ubiquitous Language:

1. EXTRACT terms from domain expert conversations
2. DOCUMENT in a glossary with precise definitions
3. ENFORCE in code - class names, method names, variables
4. EVOLVE as understanding deepens

Example Glossary Entry:
┌─────────────────────────────────────────────────────────────┐
│ Term: Order                                                  │
│ Definition: A confirmed request from a customer to purchase │
│             one or more products at agreed prices.          │
│ NOT: A shopping cart (which is an Intent, not an Order)     │
│ Context: Sales                                              │
└─────────────────────────────────────────────────────────────┘
```

## Tactical Patterns

### Entities

Objects with identity that persists over time. Equality is based on identity, not attributes.

```
Characteristics:
- Has a unique identifier
- Mutable state
- Lifecycle (created, modified, archived)
- Equality by ID

Example:
┌─────────────────────────────────────────┐
│ Entity: Order                           │
├─────────────────────────────────────────┤
│ Identity: orderId (UUID)                │
│ State: status, items, total             │
│ Behavior: addItem(), submit(), cancel() │
└─────────────────────────────────────────┘

class Order {
  private readonly id: OrderId;      // Identity - immutable
  private status: OrderStatus;        // State - mutable
  private items: OrderItem[];         // State - mutable

  constructor(id: OrderId) {
    this.id = id;
    this.status = OrderStatus.Draft;
    this.items = [];
  }

  equals(other: Order): boolean {
    return this.id.equals(other.id);  // Equality by identity
  }
}
```

### Value Objects

Objects without identity. Equality is based on attributes. Always immutable.

```
Characteristics:
- No unique identifier
- Immutable (all properties readonly)
- Equality by attributes
- Self-validating

Example:
┌─────────────────────────────────────────┐
│ Value Object: Money                     │
├─────────────────────────────────────────┤
│ Attributes: amount, currency            │
│ Behavior: add(), subtract(), format()   │
│ Invariant: amount >= 0                  │
└─────────────────────────────────────────┘

class Money {
  constructor(
    public readonly amount: number,
    public readonly currency: Currency
  ) {
    if (amount < 0) throw new Error('Amount cannot be negative');
  }

  add(other: Money): Money {
    if (!this.currency.equals(other.currency)) {
      throw new Error('Cannot add different currencies');
    }
    return new Money(this.amount + other.amount, this.currency);
  }

  equals(other: Money): boolean {
    return this.amount === other.amount &&
           this.currency.equals(other.currency);
  }
}
```

#### When to Use Value Objects

| Use Value Object | Use Entity |
|------------------|------------|
| No need to track over time | Need to track lifecycle |
| Interchangeable instances | Unique identity matters |
| Defined by attributes | Defined by continuity |
| Examples: Money, Address, DateRange | Examples: User, Order, Account |

### Aggregates

A cluster of entities and value objects with a defined boundary. One entity is the aggregate root.

```
Aggregate Design Rules:

1. PROTECT invariants at aggregate boundary
2. REFERENCE other aggregates by identity only
3. UPDATE one aggregate per transaction
4. DESIGN small aggregates (prefer single entity)

Example:
┌─────────────────────────────────────────────────────────────┐
│ Aggregate: Order                                            │
│ Root: Order (entity)                                        │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐                                        │
│  │ Order (Root)    │◄── Aggregate Root                      │
│  │ - orderId       │                                        │
│  │ - customerId ───┼──► Reference by ID only                │
│  │ - status        │                                        │
│  └────────┬────────┘                                        │
│           │                                                 │
│  ┌────────▼────────┐                                        │
│  │ OrderItem       │◄── Inside aggregate                    │
│  │ - productId ────┼──► Reference by ID only                │
│  │ - quantity      │                                        │
│  │ - price (Money) │◄── Value Object                        │
│  └─────────────────┘                                        │
└─────────────────────────────────────────────────────────────┘
```

#### Aggregate Sizing

```
Start Small:
- Begin with single-entity aggregates
- Expand only when invariants require it

Signs of Too-Large Aggregate:
- Frequent optimistic lock conflicts
- Loading too much data for simple operations
- Multiple users editing simultaneously
- Transactional failures across unrelated data

Signs of Too-Small Aggregate:
- Invariants not protected
- Business rules scattered across services
- Eventual consistency where immediate is required
```

### Domain Events

Represent something that happened in the domain. Immutable facts about the past.

```
Event Structure:
┌─────────────────────────────────────────┐
│ Event: OrderPlaced                      │
├─────────────────────────────────────────┤
│ eventId: UUID                           │
│ occurredAt: DateTime                    │
│ aggregateId: orderId                    │
│ payload:                                │
│   - customerId                          │
│   - items                               │
│   - totalAmount                         │
└─────────────────────────────────────────┘

Naming Convention:
- Past tense (OrderPlaced, not PlaceOrder)
- Domain language (not technical)
- Include all relevant data (event is immutable)

class OrderPlaced implements DomainEvent {
  readonly eventId = uuid();
  readonly occurredAt = new Date();

  constructor(
    readonly orderId: OrderId,
    readonly customerId: CustomerId,
    readonly items: OrderItemData[],
    readonly totalAmount: Money
  ) {}
}
```

#### Event Patterns

| Pattern | Description | Use Case |
|---------|-------------|----------|
| **Event Notification** | Minimal data, query for details | Loose coupling |
| **Event-Carried State** | Full data in event | Performance, offline |
| **Event Sourcing** | Events as source of truth | Audit, temporal queries |

### Repositories

Abstract persistence, providing collection-like access to aggregates.

```
Repository Principles:
- One repository per aggregate
- Returns aggregate roots only
- Hides persistence mechanism
- Supports aggregate reconstitution

interface OrderRepository {
  findById(id: OrderId): Promise<Order | null>;
  findByCustomer(customerId: CustomerId): Promise<Order[]>;
  save(order: Order): Promise<void>;
  delete(order: Order): Promise<void>;
}

// Implementation hides persistence details
class PostgresOrderRepository implements OrderRepository {
  async findById(id: OrderId): Promise<Order | null> {
    const row = await this.db.query('SELECT * FROM orders WHERE id = $1', [id]);
    return row ? this.reconstitute(row) : null;
  }

  private reconstitute(row: OrderRow): Order {
    // Rebuild aggregate from persistence
  }
}
```

## Consistency Strategies

### Transactional Consistency (ACID)

Use for invariants within an aggregate:

```
Rule: One aggregate per transaction

// Good: Single aggregate updated
async function addItemToOrder(orderId: OrderId, item: OrderItem) {
  const order = await orderRepo.findById(orderId);
  order.addItem(item);  // Business rules enforced
  await orderRepo.save(order);
}

// Bad: Multiple aggregates in one transaction
async function createOrderWithInventory() {
  await db.transaction(async (tx) => {
    await orderRepo.save(order, tx);
    await inventoryRepo.decrement(productId, quantity, tx);  // Don't do this
  });
}
```

### Eventual Consistency

Use for consistency across aggregates:

```
Pattern: Domain Events + Handlers

// Order aggregate publishes event
class Order {
  submit(): void {
    this.status = OrderStatus.Placed;
    this.addEvent(new OrderPlaced(this.id, this.customerId, this.items));
  }
}

// Separate handler updates inventory (eventually)
class InventoryHandler {
  async handle(event: OrderPlaced): Promise<void> {
    for (const item of event.items) {
      await this.inventoryService.reserve(item.productId, item.quantity);
    }
  }
}
```

### Saga Pattern

Coordinate multiple aggregates with compensation:

```
Saga: Order Fulfillment

┌─────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────┐
│ Create  │────►│ Reserve     │────►│ Charge      │────►│ Ship    │
│ Order   │     │ Inventory   │     │ Payment     │     │ Order   │
└────┬────┘     └──────┬──────┘     └──────┬──────┘     └─────────┘
     │                 │                   │
     │ Compensate:     │ Compensate:       │ Compensate:
     │ Cancel Order    │ Release Inventory │ Refund Payment
     ▼                 ▼                   ▼

On failure at any step, execute compensation in reverse order.
```

### Choosing Consistency

| Scenario | Strategy |
|----------|----------|
| Within single aggregate | Transactional (ACID) |
| Across aggregates, same service | Eventual (domain events) |
| Across services | Saga with compensation |
| Read model updates | Eventual (projection) |

## Anti-Patterns

### Anemic Domain Model

```
// Anti-pattern: Logic outside domain objects
class Order {
  id: string;
  items: Item[];
  status: string;
}

class OrderService {
  calculateTotal(order: Order): number { ... }
  validate(order: Order): boolean { ... }
  submit(order: Order): void { ... }
}

// Better: Logic inside domain objects
class Order {
  private items: OrderItem[];
  private status: OrderStatus;

  get total(): Money {
    return this.items.reduce((sum, item) => sum.add(item.subtotal), Money.zero());
  }

  submit(): void {
    this.validate();
    this.status = OrderStatus.Submitted;
  }
}
```

### Large Aggregates

```
// Anti-pattern: Everything in one aggregate
class Customer {
  orders: Order[];           // Could be thousands
  addresses: Address[];
  paymentMethods: PaymentMethod[];
  preferences: Preferences;
  activityLog: Activity[];   // Could be millions
}

// Better: Separate aggregates referenced by ID
class Customer {
  id: CustomerId;
  defaultAddressId: AddressId;
  defaultPaymentMethodId: PaymentMethodId;
}

class Order {
  customerId: CustomerId;    // Reference by ID
}
```

### Primitive Obsession

```
// Anti-pattern: Primitive types for domain concepts
function createOrder(
  customerId: string,
  productId: string,
  quantity: number,
  price: number,
  currency: string
) { ... }

// Better: Value objects
function createOrder(
  customerId: CustomerId,
  productId: ProductId,
  quantity: Quantity,
  price: Money
) { ... }
```

## Implementation Checklist

### Aggregate Design

- [ ] Single entity can be aggregate root
- [ ] Invariants are protected at boundary
- [ ] Other aggregates referenced by ID only
- [ ] Fits in memory comfortably
- [ ] One transaction per aggregate

### Entity Implementation

- [ ] Has unique identifier
- [ ] Equality based on ID
- [ ] Encapsulates business rules
- [ ] State changes through methods

### Value Object Implementation

- [ ] All properties immutable
- [ ] Equality based on attributes
- [ ] Self-validating
- [ ] Operations return new instances

### Repository Implementation

- [ ] One per aggregate
- [ ] Returns aggregate roots only
- [ ] Hides persistence details
- [ ] Supports queries needed by domain

## References

- [Pattern Implementation Examples](examples/ddd-patterns.md) - Code examples in multiple languages
- [Aggregate Design Guide](reference.md) - Detailed aggregate sizing heuristics
